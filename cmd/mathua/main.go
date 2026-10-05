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
	"runtime/debug"
	"strings"
	"syscall"
	"time"

	"github.com/chuma-beep/mathua/internal/auth"
	"github.com/chuma-beep/mathua/internal/concepts"
	"github.com/chuma-beep/mathua/internal/engine"
	"github.com/chuma-beep/mathua/internal/generator/all"
	"github.com/chuma-beep/mathua/internal/lessons"
	"github.com/chuma-beep/mathua/internal/planning"
	"github.com/chuma-beep/mathua/internal/repair"
	"github.com/chuma-beep/mathua/internal/server"
	"github.com/chuma-beep/mathua/internal/solutions"
	"github.com/chuma-beep/mathua/internal/storage"
)

// dataPath resolves a file inside the data directory.
//
// Every data path used to be a bare relative string, which made the server depend on its
// working directory: started anywhere but the repo root, `data/courses.json` was missing,
// the planner came back nil, and — because that error was discarded — nothing said so.
// Resolution order is $MATHUA_DATA_DIR, then the executable's own directory, then the
// working directory, so an installed binary works without configuration and the repo-root
// case keeps working unchanged.
func dataPath(name string) string {
	dirs := []string{}
	if d := os.Getenv("MATHUA_DATA_DIR"); d != "" {
		dirs = append(dirs, d)
	}
	if exe, err := os.Executable(); err == nil {
		if exe, err = filepath.EvalSymlinks(exe); err == nil {
			dirs = append(dirs, filepath.Dir(exe))
		}
	}
	dirs = append(dirs, "data")
	for _, d := range dirs {
		p := filepath.Join(d, name)
		if st, err := os.Stat(p); err == nil && !st.IsDir() {
			return p
		}
	}
	// Nothing found: return the historical cwd-relative path so the caller's error message
	// names the location it looked in.
	return filepath.Join("data", name)
}

func main() {
	serve := flag.Bool("serve", false, "run web server")
	port := flag.Int("port", 8080, "web server port")
	noAuth := flag.Bool("no-auth", false, "disable authentication (dev mode)")
	repairGrading := flag.Bool("repair-grading", false, "re-grade persisted attempts, repair grader false negatives, then exit")
	repairDryRun := flag.Bool("repair-grading-dry-run", false, "report grading repairs without writing, then exit")
	flag.Parse()

	dag, err := concepts.LoadDir(dataPath("concepts"))
	if err != nil {
		log.Fatalf("load concepts: %v", err)
	}
	fmt.Printf("loaded %d concepts across %d domains\n", dag.Count(), len(dag.Domains()))

	dsn := os.Getenv("DATABASE_URL")
	var repo storage.Repository
	if strings.HasPrefix(dsn, "postgres://") || strings.HasPrefix(dsn, "postgresql://") {
		fmt.Println("connecting to PostgreSQL...")
		repo, err = storage.NewPostgresStore(dsn)
	} else {
		if dsn != "" || *serve {
			if dsn == "" {
				dsn = "mathua.db"
			}
			repo, err = storage.NewSQLiteStore(dsn)
		} else {
			repo, err = storage.NewSQLiteStore("mathua.db")
		}
	}
	if err != nil {
		log.Fatalf("storage: %v", err)
	}
	defer repo.Close()

	// Persist an auto-generated JWT secret next to the SQLite file so logins
	// survive restarts without a JWT_SECRET env (single-machine volumes).
	// Postgres deployments must set JWT_SECRET (no shared file to use).
	dbPath := strings.TrimPrefix(dsn, "sqlite:")
	if dbPath == "" {
		dbPath = "mathua.db"
	}
	if dbPath != ":memory:" && !strings.HasPrefix(dbPath, "postgres") {
		auth.SetSecretFile(filepath.Join(filepath.Dir(dbPath), ".jwt_secret"))
	}

	reg := all.Registry()

	ll, err := lessons.Load("data/lessons")
	if err != nil {
		fmt.Printf("lessons not loaded: %v (continuing without lessons)\n", err)
	}

	// Corpus solution schemas: authored prose per concept, interpolated with
	// the facts each generator publishes. Optional — without them every problem
	// keeps the generator's own explanation.
	if sol, err := solutions.Load("data/lessons/solutions"); err != nil {
		fmt.Printf("solution schemas not loaded: %v (using generator explanations)\n", err)
	} else {
		reg.SetSolutions(sol)
	}

	// Purge stored questions for concepts that have live generators: rows
	// seeded by older generator versions would otherwise be served forever
	// (INSERT OR IGNORE makes them immortal). Live generation is cheap and
	// always current, so reproducible content never needs a DB copy.
	if purger, ok := repo.(interface {
		PurgeGeneratedQuestions(map[string]bool) (int64, error)
	}); ok {
		idSet := make(map[string]bool)
		for _, id := range reg.Concepts() {
			idSet[id] = true
		}
		if n, err := purger.PurgeGeneratedQuestions(idSet); err != nil {
			log.Printf("warning: question purge failed: %v", err)
		} else if n > 0 {
			fmt.Printf("purged %d stale generated questions\n", n)
		}
	}

	// The learning planner is optional in the sense that the rest of the server runs without
	// it — but it was loaded with `_`, so a missing or unreadable catalog left `planner` nil
	// with nothing logged, and every /api/destinations/{id}/estimate then answered 404
	// "planner unavailable". From the learner's side that was a planner page with an empty
	// dropdown and no message. Fail loudly instead.
	planner, err := planning.Load(dataPath("courses.json"), dag)
	if err != nil {
		log.Printf("WARNING: course catalog unavailable (%v) — the learning planner and every "+
			"destination estimate will be unavailable. Set MATHUA_DATA_DIR if the data "+
			"directory is not next to the binary or the working directory", err)
	}
	if planner != nil {
		if ds, err := planning.LoadDestinations(dataPath("destinations.json")); err != nil {
			log.Printf("WARNING: destinations unavailable (%v) — destination estimates will "+
				"return 404 until this is fixed", err)
		} else {
			planner.SetDestinations(ds)
		}
	}
	eng := engine.New(repo, dag, reg, ll, planner)

	// Data repair: re-grade persisted attempts and fix grader false negatives
	// (streak/weakness corruption). Exit without starting the server.
	if *repairGrading || *repairDryRun {
		if _, err := repair.Run(eng, repo, dag, *repairDryRun, os.Stdout); err != nil {
			eng.Close()
			log.Fatalf("repair-grading: %v", err)
		}
		eng.Close()
		return
	}
	if ll != nil {
		fmt.Printf("%d lessons loaded\n", ll.Count())
	}
	rev := "dev"
	if bi, ok := debug.ReadBuildInfo(); ok {
		for _, s := range bi.Settings {
			if s.Key == "vcs.revision" && len(s.Value) >= 7 {
				rev = s.Value[:7]
			}
		}
	}
	lessonCount := 0
	if ll != nil {
		lessonCount = ll.Count()
	}
	fmt.Printf("mathua build %s — %d generators, %d lessons ingested\n", rev, reg.Count(), lessonCount)

	fmt.Printf("starting web server on :%d", *port)
	var authSvc *auth.AuthService
	if !*noAuth {
		authSvc = auth.New(repo)
		fmt.Print(" (auth enabled)")
		if n, err := authSvc.BackfillUsernames(); err != nil {
			log.Printf("username backfill failed: %v", err)
		} else if n > 0 {
			fmt.Printf(" (assigned %d usernames)", n)
		}
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
		Handler:      server.GzipMiddleware(securityHeaders(mux)),
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
}

// nextStaticFS wraps http.FileServer to support Next.js static exports
// where routes like /session map to session.html files.
func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("X-XSS-Protection", "0")
		next.ServeHTTP(w, r)
	})
}

func nextStaticFS(root string) http.Handler {
	fs := http.FileServer(http.Dir(root))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cleaned := filepath.Clean(r.URL.Path)
		if strings.Contains(cleaned, "..") {
			http.NotFound(w, r)
			return
		}
		if r.URL.Path == "/" {
			fs.ServeHTTP(w, r)
			return
		}
		localPath := filepath.Join(root, cleaned)
		if _, err := os.Stat(localPath); err == nil {
			fs.ServeHTTP(w, r)
			return
		}
		if filepath.Ext(cleaned) == "" {
			htmlPath := cleaned + ".html"
			if _, err := os.Stat(filepath.Join(root, htmlPath)); err == nil {
				r.URL.Path = htmlPath
				fs.ServeHTTP(w, r)
				return
			}
		}
		fs.ServeHTTP(w, r)
	})
}
