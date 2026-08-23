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

	"github.com/chuma-beep/mathua/internal/auth"
	"github.com/chuma-beep/mathua/internal/concepts"
	"github.com/chuma-beep/mathua/internal/engine"
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/abstract"
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
	"github.com/chuma-beep/mathua/internal/generator/precalculus"
	"github.com/chuma-beep/mathua/internal/generator/statistics"
	"github.com/chuma-beep/mathua/internal/generator/topology"
	"github.com/chuma-beep/mathua/internal/generator/trigonometry"
	"github.com/chuma-beep/mathua/internal/lessons"
	"github.com/chuma-beep/mathua/internal/planning"
	"github.com/chuma-beep/mathua/internal/server"
	"github.com/chuma-beep/mathua/internal/storage"
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
	precalculus.Register(reg)
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

	planner, _ := planning.Load("data/courses.json", dag)
	eng := engine.New(repo, dag, reg, ll, planner)
	if ll != nil {
		fmt.Printf("%d lessons loaded\n", ll.Count())
	}
	fmt.Printf("engine ready — %d generators registered\n", reg.Count())

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
		Handler:      securityHeaders(mux),
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
