package main

import (
	"context"
	"errors"
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

// openRepo opens the configured store and reports which path it used.
//
// Extracted because `mathua admin …` needs a store and does not need a generator registry, a
// lesson loader or a concept graph. Loading the whole corpus to change one column meant the
// command printed "loaded 657 concepts" on its way to a two-line operation, and made a typo
// in the corpus fail a maintenance task it had no business touching.
//
// The store's own Migrate has already run by the time this returns, which is what makes the
// admin command work against a database created before the role column existed.
func openRepo() (storage.Repository, string, error) {
	dsn := os.Getenv("DATABASE_URL")
	if strings.HasPrefix(dsn, "postgres://") || strings.HasPrefix(dsn, "postgresql://") {
		fmt.Println("connecting to PostgreSQL...")
		r, err := storage.NewPostgresStore(dsn)
		return r, dsn, err
	}
	path := dsn
	if path == "" {
		path = "mathua.db"
	}
	r, err := storage.NewSQLiteStore(path)
	return r, path, err
}

// adminCommand is the operator's bootstrap. It is the only way the first administrator comes
// into existence, and it is a shell command on purpose.
//
// Every alternative was worse. An HTTP endpoint means anyone who can reach the server can ask
// to become an administrator, which is not a bootstrap, it is an open door. A role in the
// signup body means a client chooses its own privilege. An env var naming the first admin
// needs a database row to point at and silently stops working once it does. Shell access to
// the box is the one credential that already implies full control of the deployment, so it is
// the right place for the decision that someone should have it.
// adminActorCLI is the actor recorded for a shell-initiated promotion. See the call site for
// why it is not the operator's email.
const adminActorCLI = "cli"

func adminCommand(args []string) error {
	if len(args) == 0 {
		return errors.New("usage: mathua admin promote <email>")
	}
	switch args[0] {
	case "promote":
		if len(args) != 2 {
			return errors.New("usage: mathua admin promote <email>")
		}
		email := strings.TrimSpace(args[1])
		if email == "" {
			return errors.New("email is required")
		}
		repo, _, err := openRepo()
		if err != nil {
			return err
		}
		defer repo.Close()

		st, err := repo.FindByEmail(email)
		if err != nil {
			return err
		}
		if st == nil {
			// Refusing to create the account is the important half. A bootstrap that can sign
			// someone up would be reachable by anyone who can guess an address, and an
			// administrator with no password is an administrator nobody can attribute.
			return fmt.Errorf("no account with email %q — sign up through the app first, then promote", email)
		}

		// The same primitive the web admin uses, so there is exactly one path that can change a
		// role and exactly one place the audit event is written. A second implementation here
		// would be a second set of invariants to keep — and it is the one nobody would remember
		// to update.
		//
		// `cli` is the actor because that is the truth: shell access carries no account identity,
		// and the argument for trusting this command is that it needs the machine, not a
		// credential. Writing the operator's email here would be inventing an identity the code
		// does not have.
		res, err := repo.SetRoleAudited(adminActorCLI, st.ID, storage.RoleAdmin)
		switch {
		case errors.Is(err, storage.ErrRoleUnchanged):
			fmt.Printf("%s (%s) is already an admin; nothing to do\n", email, st.ID)
			return nil
		case err != nil:
			return err
		}
		fmt.Printf("promoted %s (%s) to admin — audit %s by %s\n",
			email, st.ID, res.Action, adminActorCLI)
		return nil
	default:
		return fmt.Errorf("unknown admin command %q. Available: promote <email>", args[0])
	}
}

func main() {
	// Accepted and ignored: the server always starts, and it did before this flag meant
	// anything — both arms of the branch it used to select opened the same store. The Dockerfile
	// passes it, so removing it would break a deployment for no gain. Kept because deleting it
	// silently would be the only real change, and it is not one worth making at the same time
	// as adding administration.
	flag.Bool("serve", false, "run web server (accepted for compatibility; the server always starts)")
	port := flag.Int("port", 8080, "web server port")
	noAuth := flag.Bool("no-auth", false, "disable authentication (dev mode)")
	repairGrading := flag.Bool("repair-grading", false, "re-grade persisted attempts, repair grader false negatives, then exit")
	repairDryRun := flag.Bool("repair-grading-dry-run", false, "report grading repairs without writing, then exit")
	flag.Parse()

	// Subcommands run before anything else is loaded. They are operator tools: someone with
	// shell access to the box, not someone with a token. That distinction is the whole reason
	// the first administrator is created here and not through an endpoint — there is no
	// /register-admin route, no ?role=admin, and no request a browser can make that reaches
	// this code.
	if len(flag.Args()) > 0 {
		switch flag.Arg(0) {
		case "admin":
			if err := adminCommand(flag.Args()[1:]); err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				os.Exit(1)
			}
			return
		default:
			fmt.Fprintf(os.Stderr, "unknown command %q. Available: admin\n", flag.Arg(0))
			os.Exit(2)
		}
	}

	dag, err := concepts.LoadDir(dataPath("concepts"))
	if err != nil {
		log.Fatalf("load concepts: %v", err)
	}
	fmt.Printf("loaded %d concepts across %d domains\n", dag.Count(), len(dag.Domains()))

	repo, dsn, err := openRepo()
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
		// Regular files only. `os.Stat` succeeds for a directory too, and a Next.js export
		// creates one for any route that has child routes — `/admin` because of /admin/users,
		// `/docs` because of /docs/architecture. Handing a directory to http.FileServer makes
		// it redirect to the trailing-slash form and then render a directory listing, so
		// `/admin` answered 301 and its target answered 200 with a list of files instead of the
		// administration page. The page exists at `out/admin.html`; only the `.html` fallback
		// below finds it, and this branch was shadowing it.
		if st, err := os.Stat(localPath); err == nil && !st.IsDir() {
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
