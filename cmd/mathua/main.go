package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/chuma-beep/mathua/internal/auth"
	"github.com/chuma-beep/mathua/internal/concepts"
	"github.com/chuma-beep/mathua/internal/diagnostic"
	"github.com/chuma-beep/mathua/internal/engine"
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/abstract"
	"github.com/chuma-beep/mathua/internal/grader"
	"github.com/chuma-beep/mathua/internal/generator/algebra"
	"github.com/chuma-beep/mathua/internal/generator/arithmetic"
	"github.com/chuma-beep/mathua/internal/generator/calculus"
	"github.com/chuma-beep/mathua/internal/generator/complex"
	"github.com/chuma-beep/mathua/internal/generator/counting"
	"github.com/chuma-beep/mathua/internal/generator/discrete"
	"github.com/chuma-beep/mathua/internal/generator/fractions"
	"github.com/chuma-beep/mathua/internal/generator/geometry"
	"github.com/chuma-beep/mathua/internal/generator/linalg"
	"github.com/chuma-beep/mathua/internal/generator/machinelearning"
	"github.com/chuma-beep/mathua/internal/generator/numtheory"
	"github.com/chuma-beep/mathua/internal/generator/odes"
	"github.com/chuma-beep/mathua/internal/generator/prealgebra"
	"github.com/chuma-beep/mathua/internal/generator/statistics"
	"github.com/chuma-beep/mathua/internal/generator/topology"
	"github.com/chuma-beep/mathua/internal/generator/trigonometry"
	"github.com/chuma-beep/mathua/internal/lessons"
	"github.com/chuma-beep/mathua/internal/mastery"
	"github.com/chuma-beep/mathua/internal/planning"
	"github.com/chuma-beep/mathua/internal/scoring"
	"github.com/chuma-beep/mathua/internal/server"
	"github.com/chuma-beep/mathua/internal/storage"
	"github.com/chuma-beep/mathua/ui/tui"
)

func main() {
	serve := flag.Bool("serve", false, "run web server")
	port := flag.Int("port", 8080, "web server port")
	noAuth := flag.Bool("no-auth", false, "disable authentication (dev mode)")
	flag.Parse()

	dag, err := concepts.LoadDir("data/concepts")
	if err != nil {
		log.Fatalf("load concepts: %v", err)
	}
	fmt.Printf("loaded %d concepts across %d domains\n", dag.Count(), len(dag.Domains()))

	dsn := os.Getenv("DATABASE_URL")
	if strings.HasPrefix(dsn, "postgres://") || strings.HasPrefix(dsn, "postgresql://") {
		fmt.Println("warning: PostgreSQL backend not yet implemented, falling back to SQLite")
		dsn = "mathua.db"
	}
	var repo storage.Repository
	if dsn != "" || *serve {
		if dsn == "" {
			dsn = "mathua.db"
		}
		repo, err = storage.NewSQLiteStore(dsn)
	} else {
		repo, err = storage.NewSQLiteStore("mathua.db")
	}
	if err != nil {
		log.Fatalf("storage: %v", err)
	}
	defer repo.Close()

	reg := generator.NewRegistry()
	counting.Register(reg)
	arithmetic.Register(reg)
	fractions.Register(reg)
	geometry.Register(reg)
	prealgebra.Register(reg)
	algebra.Register(reg)
	trigonometry.Register(reg)
	statistics.Register(reg)
	numtheory.Register(reg)
	complex.Register(reg)
	linalg.Register(reg)
	machinelearning.Register(reg)
	discrete.Register(reg)
	calculus.Register(reg)
	odes.Register(reg)
	abstract.Register(reg)
	topology.Register(reg)

	ll, err := lessons.Load("data/lessons")
	if err != nil {
		fmt.Printf("lessons not loaded: %v (continuing without lessons)\n", err)
	}

	planner, _ := planning.Load("data/courses.json", dag)
	eng := engine.New(repo, dag, reg, ll, planner)
	if ll != nil {
		fmt.Printf("%d lessons loaded\n", ll.Count())
	}
	fmt.Printf("engine ready — %d generators registered\n", reg.Count())

	if *serve {
		fmt.Printf("starting web server on :%d", *port)
		var authSvc *auth.AuthService
		if !*noAuth {
			authSvc = auth.New(repo)
			fmt.Print(" (auth enabled)")
		} else {
			fmt.Print(" (no auth)")
		}
		fmt.Println()
		srv := server.New(eng, repo, authSvc)
		mux := http.NewServeMux()
		srv.Register(mux)
		if info, err := os.Stat("web/next-app/out"); err == nil && info.IsDir() {
			mux.Handle("/", nextStaticFS("web/next-app/out"))
			fmt.Println("serving static frontend from web/next-app/out")
		}
		httpSrv := &http.Server{
			Addr:         fmt.Sprintf(":%d", *port),
			Handler:      mux,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  60 * time.Second,
		}
		go func() {
			if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("server: %v", err)
			}
		}()
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		fmt.Println("\nshutting down gracefully...")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := httpSrv.Shutdown(ctx); err != nil {
			log.Fatalf("forced shutdown: %v", err)
		}
		eng.Close()
		fmt.Println("server stopped")
	} else {
		m := tui.New()
		var studentID, sessionID string
		var diagSession *diagnostic.Session
		var diagCount int

		m.OnNext = func() tui.SessionMsg {
			if sessionID == "" {
				return tui.SessionMsg{}
			}
			q, err := eng.NextQuestion(sessionID, studentID)
			if err != nil || q == nil {
				return tui.SessionMsg{}
			}
			c := dag.Concept(q.ConceptID)
			domain, subdomain := "", ""
			if c != nil {
				domain = c.Domain
				subdomain = c.Subdomain
			}
			prog, _ := eng.GetProgress(studentID)
			masteryPct := 0.0
			streak := 0
		streakNeeded := 10
		status := "UNSEEN"
			if p, ok := prog[q.ConceptID]; ok {
				streak = p.Streak
				status = p.Status
				if c != nil && c.MasteryThreshold.Streak > 0 {
					masteryPct = float64(streak) / float64(c.MasteryThreshold.Streak)
					streakNeeded = c.MasteryThreshold.Streak
				}
				if masteryPct > 1 {
					masteryPct = 1
				}
			}
			timeLimit := 60.0
			if c != nil {
				timeLimit = c.MasteryThreshold.AvgTimeSeconds
			}
			lessonTitle := ""
			if q.Lesson != nil {
				lessonTitle = q.Lesson.Title
			}
			return tui.SessionMsg{
				ConceptID:    q.ConceptID,
				Domain:       domain,
				Subdomain:    subdomain,
				Label:        q.ConceptName,
				Question:     q.Question,
				Difficulty:   0.5,
				TimeLimit:    timeLimit,
				MasteryPct:   masteryPct,
				Streak:       streak,
				StreakNeeded: streakNeeded,
				Status:       status,
				LessonTitle:  lessonTitle,
			}
		}

		m.OnSubmit = func(conceptID, answer string, elapsed float64) tui.SubmitResultMsg {
			result, err := eng.SubmitAnswer(sessionID, studentID, answer, elapsed)
			if err != nil {
				return tui.SubmitResultMsg{Correct: false, Explanation: err.Error()}
			}
			masteryAchieved := result.NewStatus == mastery.StatusMastered
			c := dag.Concept(conceptID)
			streakNeeded := 10
			if c != nil {
				streakNeeded = c.MasteryThreshold.Streak
			}
			newMastery := 0.0
			if streakNeeded > 0 {
				newMastery = float64(result.Streak) / float64(streakNeeded)
				if newMastery > 1 {
					newMastery = 1
				}
			}
			var unlocked []string
			if masteryAchieved {
				for _, dep := range dag.DependentsOf(conceptID) {
					unlocked = append(unlocked, dep.ID)
				}
			}
			return tui.SubmitResultMsg{
				Correct:          result.Correct,
				UserAnswer:       answer,
				CorrectAnswer:    result.ExpectedAnswer,
				Explanation:      result.Explanation,
				ElapsedSecs:      elapsed,
				NewMastery:       newMastery,
				NewStreak:        result.Streak,
				NewStatus:        string(result.NewStatus),
				MasteryAchieved:  masteryAchieved,
				UnlockedConcepts: unlocked,
			}
		}

		m.OnProgress = func() tui.ProgressMsg {
			prog, _ := eng.GetProgress(studentID)
			domains := make(map[string]*tui.DomainProgress)
			subdomains := make(map[string]map[string]*tui.SubdomainProgress)

			for _, c := range dag.Order() {
				d, ok := domains[c.Domain]
				if !ok {
					d = &tui.DomainProgress{Domain: c.Domain}
					domains[c.Domain] = d
					subdomains[c.Domain] = make(map[string]*tui.SubdomainProgress)
				}
				d.Total++
				if p, ok := prog[c.ID]; ok && p.Status == "MASTERED" {
					d.Mastered++
				}
				sd, ok := subdomains[c.Domain][c.Subdomain]
				if !ok {
					sd = &tui.SubdomainProgress{Name: c.Subdomain}
					subdomains[c.Domain][c.Subdomain] = sd
				}
				sd.Total++
				if p, ok := prog[c.ID]; ok && p.Status == "MASTERED" {
					sd.Mastered++
				}
			}

			var dl []tui.DomainProgress
			for _, d := range domains {
				for _, sd := range subdomains[d.Domain] {
					d.Subdomains = append(d.Subdomains, *sd)
				}
				dl = append(dl, *d)
			}
			return tui.ProgressMsg{Domains: dl}
		}

		m.OnStudy = func() []tui.StudyLesson {
		ll := eng.GetLessonLoader()
		if ll == nil {
			return nil
		}
		byDomain := ll.LessonsByDomain()
		var out []tui.StudyLesson
		for domain, lessons := range byDomain {
			for _, l := range lessons {
				for _, cid := range l.Concepts {
					out = append(out, tui.StudyLesson{
						Title:     l.Title,
						Body:      l.Body,
						ConceptID: cid,
						Domain:    domain,
					})
				}
			}
		}
		return out
	}

	m.OnBrowse = func() tui.ConceptTreeMsg {
		nodes := eng.ConceptTree(studentID)
		var domains []tui.DomainNode
		for _, d := range nodes {
			var subs []tui.SubdomainNode
			for _, sd := range d.Subdomains {
				var concepts []tui.ConceptNode
				for _, c := range sd.Concepts {
					concepts = append(concepts, tui.ConceptNode{
						ID:         c.ID,
						Label:      c.Label,
						Unlocked:   c.Unlocked,
						MasteryPct: c.MasteryPct,
						Streak:     c.Streak,
						Status:     c.Status,
					})
				}
				subs = append(subs, tui.SubdomainNode{Name: sd.Name, Concepts: concepts})
			}
			domains = append(domains, tui.DomainNode{Name: d.Name, Subdomains: subs})
		}
		return tui.ConceptTreeMsg{Domains: domains}
	}

	m.OnSelectConcept = func(conceptID string) tui.SessionMsg {
		if sessionID == "" {
			return tui.SessionMsg{}
		}
		q, err := eng.PracticeConcept(sessionID, studentID, conceptID)
		if err != nil || q == nil {
			return tui.SessionMsg{}
		}
		c := dag.Concept(conceptID)
		domain, subdomain := "", ""
		if c != nil {
			domain = c.Domain
			subdomain = c.Subdomain
		}
		prog, _ := eng.GetProgress(studentID)
		masteryPct := 0.0
		streak := 0
		streakNeeded := 10
		status := "UNSEEN"
		if p, ok := prog[conceptID]; ok {
			streak = p.Streak
			status = p.Status
			if c != nil && c.MasteryThreshold.Streak > 0 {
				masteryPct = float64(streak) / float64(c.MasteryThreshold.Streak)
				streakNeeded = c.MasteryThreshold.Streak
			}
			if masteryPct > 1 {
				masteryPct = 1
			}
		}
		timeLimit := 60.0
		if c != nil {
			timeLimit = c.MasteryThreshold.AvgTimeSeconds
		}
		lessonTitle := ""
		if q.Lesson != nil {
			lessonTitle = q.Lesson.Title
		}
		return tui.SessionMsg{
			ConceptID:    conceptID,
			Domain:       domain,
			Subdomain:    subdomain,
			Label:        q.ConceptName,
			Question:     q.Question,
			TimeLimit:    timeLimit,
			MasteryPct:   masteryPct,
			Streak:       streak,
			StreakNeeded: streakNeeded,
			Status:       status,
			LessonTitle:  lessonTitle,
		}
	}

	m.OnConceptDetail = func(conceptID string) tui.ConceptDetailMsg {
		detail, err := eng.ConceptDetail(studentID, conceptID)
		if err != nil {
			return tui.ConceptDetailMsg{}
		}
		prereqs := make([]tui.PrereqInfo, len(detail.Prerequisites))
		for i, p := range detail.Prerequisites {
			prereqs[i] = tui.PrereqInfo{
				ID:         p.ID,
				Label:      p.Label,
				Status:     p.Status,
				MasteryPct: p.MasteryPct,
			}
		}
		msg := tui.ConceptDetailMsg{
			ConceptID:    detail.Concept.ID,
			Label:        detail.Concept.Label,
			Domain:       detail.Concept.Domain,
			Subdomain:    detail.Concept.Subdomain,
			Prerequisites: prereqs,
			Unlocked:     detail.Unlocked,
		}
		if detail.Lesson != nil {
			msg.LessonTitle = detail.Lesson.Title
			msg.LessonBody = detail.Lesson.Body
		}
		if detail.Progress != nil {
			msg.Status = detail.Progress.Status
			msg.Streak = detail.Progress.Streak
			msg.StreakNeeded = detail.Progress.RequiredStreak
			msg.MasteryPct = detail.Progress.MasteryPct
		}
		return msg
	}

	m.OnDiagStart = func() tui.DiagQuestionMsg {
		if diagSession == nil || eng.IsDiagnosticComplete(diagSession) {
			diagSession = eng.StartDiagnostic()
			diagCount = 0
		}
		prob, cid, err := eng.NextDiagnosticQuestion(diagSession)
		if err != nil || prob == nil {
			return tui.DiagQuestionMsg{}
		}
		diagCount++
		c := dag.Concept(cid)
		label := cid
		if c != nil {
			label = c.Label
		}
		return tui.DiagQuestionMsg{
			ConceptID:   cid,
			ConceptName: label,
			Question:    prob.Question,
			Count:       diagCount,
		}
	}

	m.OnDiagSubmit = func(conceptID, answer string, elapsed float64) tui.DiagFeedbackMsg {
		prob := diagSession.LastProblem
		correct := false
		correctAnswer := ""
		explanation := ""
		if prob != nil {
			c := dag.Concept(conceptID)
			gt := grader.GradingNumeric
			if c != nil {
				gt = grader.GradingType(c.GradingType)
			}
			grResult := eng.GetGrader().Grade(gt, prob.Answer, answer)
			correct = grResult.Correct
			correctAnswer = prob.Answer
			explanation = prob.Explanation
		}
		eng.SubmitDiagnosticAnswer(diagSession, conceptID, correct, elapsed < 10.0)
		done := eng.IsDiagnosticComplete(diagSession)
		return tui.DiagFeedbackMsg{
			Correct:       correct,
			UserAnswer:    answer,
			CorrectAnswer: correctAnswer,
			Explanation:   explanation,
			Done:          done,
		}
	}

	m.OnWelcome = func() tui.WelcomeStatsMsg {
			if studentID == "" {
				st, err := eng.CreateStudent("")
				if err != nil {
					return tui.WelcomeStatsMsg{TotalConcepts: dag.Count()}
				}
				studentID = st.ID
				sess, _ := repo.CreateSession(studentID)
				sessionID = sess.ID
			}
			scores, _ := scoring.NewUpdater(dag, repo).Compute(studentID)
			prog, _ := eng.GetProgress(studentID)

			domains := make(map[string]*tui.DomainProgress)
			for _, c := range dag.Order() {
				d, ok := domains[c.Domain]
				if !ok {
					d = &tui.DomainProgress{Domain: c.Domain}
					domains[c.Domain] = d
				}
				d.Total++
				if p, ok := prog[c.ID]; ok && p.Status == "MASTERED" {
					d.Mastered++
				}
			}
			var dl []tui.DomainProgress
			for _, d := range domains {
				dl = append(dl, *d)
			}

			return tui.WelcomeStatsMsg{
				TotalMastered:  scores.ConceptsMastered,
				TotalConcepts:  dag.Count(),
				DayStreak:      scores.CurrentStreak,
				Level:          scores.Level,
				LevelName:      scores.Level,
				WeeklyScore:    scores.WeeklyScore,
				DomainProgress: dl,
			}
		}

		if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
			log.Fatalf("tui: %v", err)
		}
	}
}

// nextStaticFS wraps http.FileServer to support Next.js static exports
// where routes like /session map to session.html files.
func nextStaticFS(root string) http.Handler {
	fs := http.FileServer(http.Dir(root))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			fs.ServeHTTP(w, r)
			return
		}
		localPath := filepath.Join(root, r.URL.Path)
		if _, err := os.Stat(localPath); err == nil {
			fs.ServeHTTP(w, r)
			return
		}
		if filepath.Ext(r.URL.Path) == "" {
			htmlPath := r.URL.Path + ".html"
			if _, err := os.Stat(filepath.Join(root, htmlPath)); err == nil {
				r.URL.Path = htmlPath
				fs.ServeHTTP(w, r)
				return
			}
		}
		fs.ServeHTTP(w, r)
	})
}
