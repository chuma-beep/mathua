package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/chuma-beep/mathua/internal/concepts"
	"github.com/chuma-beep/mathua/internal/engine"
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/algebra"
	"github.com/chuma-beep/mathua/internal/generator/arithmetic"
	"github.com/chuma-beep/mathua/internal/generator/counting"
	"github.com/chuma-beep/mathua/internal/generator/fractions"
	"github.com/chuma-beep/mathua/internal/generator/geometry"
	"github.com/chuma-beep/mathua/internal/generator/prealgebra"
	"github.com/chuma-beep/mathua/internal/lessons"
	"github.com/chuma-beep/mathua/internal/server"
	"github.com/chuma-beep/mathua/internal/storage"
)

func main() {
	serve := flag.Bool("serve", false, "run web server")
	port := flag.Int("port", 8080, "web server port")
	flag.Parse()

	dag, err := concepts.Load("data/concepts.json")
	if err != nil {
		log.Fatalf("load concepts: %v", err)
	}
	fmt.Printf("loaded %d concepts across %d domains\n", dag.Count(), len(dag.Domains()))

	dsn := os.Getenv("DATABASE_URL")
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

	ll, err := lessons.Load("data/lessons")
	if err != nil {
		fmt.Printf("lessons not loaded: %v (continuing without lessons)\n", err)
	}

	eng := engine.New(repo, dag, reg, ll)
	if ll != nil {
		fmt.Printf("%d lessons loaded\n", ll.Count())
	}
	fmt.Printf("engine ready — %d generators registered\n", reg.Count())

	if *serve {
		fmt.Printf("starting web server on :%d\n", *port)
		srv := server.New(eng, repo)
		mux := http.NewServeMux()
		srv.Register(mux)
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), mux))
	} else {
		runCLI(eng, repo)
	}
}

func runCLI(eng *engine.Engine, repo storage.Repository) {
	fmt.Println("\n=== Mathua — Math Understanding Agent ===")
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Your name [learner]: ")
	scanner.Scan()
	name := strings.TrimSpace(scanner.Text())
	if name == "" {
		name = "learner"
	}

	st, err := eng.CreateStudent(name)
	if err != nil {
		log.Fatalf("create student: %v", err)
	}

	sess, err := repo.CreateSession(st.ID)
	if err != nil {
		log.Fatalf("create session: %v", err)
	}
	fmt.Printf("Welcome, %s!  Type 'quit', 'stats', or 'skip'.\n\n", name)

	for {
		q, err := eng.NextQuestion(sess.ID, st.ID)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			break
		}
		if q == nil {
			fmt.Println("All available concepts mastered! Nothing left to practice.")
			fmt.Println("Come back tomorrow for reviews or add more generators.")
			printStats(eng, st.ID)
			break
		}

		if q.Lesson != nil {
			fmt.Printf("[Lesson: %s]\n", q.Lesson.Title)
		}
		status := engCurrentStatus(eng, st.ID, q.ConceptID)
		fmt.Printf("[%s] %s: %s\n", status, q.ConceptName, q.Question)
		fmt.Print("> ")

		t0 := time.Now()
		if !scanner.Scan() {
			break
		}
		elapsed := time.Since(t0).Seconds()

		input := strings.TrimSpace(scanner.Text())
		switch strings.ToLower(input) {
		case "quit", "q":
			printStats(eng, st.ID)
			return
		case "stats", "s":
			printStats(eng, st.ID)
			continue
		case "skip":
			_ = repo.RecordAttempt(storage.AttemptEntry{
				SessionID: sess.ID, StudentID: st.ID,
				ConceptID: q.ConceptID, Answer: "(skipped)",
				Expected: "(skipped)", Correct: false,
				ElapsedSeconds: elapsed, Timestamp: time.Now().UTC(),
			})
			continue
		case "":
			continue
		}

		result, err := eng.SubmitAnswer(sess.ID, st.ID, input, elapsed)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			break
		}
		if result.Correct {
			fmt.Printf("Correct! (%.1fs, streak %d/%d → %s)\n",
				elapsed, result.Streak, result.RequiredStreak, result.NewStatus)
		} else {
			fmt.Printf("Wrong — %s\n  %s\n", result.Feedback, result.Explanation)
		}
	}
}

func engCurrentStatus(eng *engine.Engine, studentID, conceptID string) string {
	prog, _ := eng.GetProgress(studentID)
	if p, ok := prog[conceptID]; ok {
		return p.Status
	}
	return "UNSEEN"
}

func printStats(eng *engine.Engine, studentID string) {
	scores, err := eng.GetScores(studentID)
	if err != nil {
		fmt.Printf("stats error: %v\n", err)
		return
	}
	fmt.Printf("\n--- %s ---\n", scores.Level)
	fmt.Printf("Mastered: %d  Streak: %dd  Score: %d\n",
		scores.ConceptsMastered, scores.CurrentStreak, scores.WeeklyScore)
	prog, _ := eng.GetProgress(studentID)
	inProgress := 0
	for _, p := range prog {
		if p.Status != "MASTERED" && p.Status != "UNSEEN" {
			inProgress++
		}
	}
	fmt.Printf("In progress: %d  Unseen: %d\n", inProgress, 284-scores.ConceptsMastered-inProgress)
}
