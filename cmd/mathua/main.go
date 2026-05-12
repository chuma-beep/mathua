package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	tea "github.com/charmbracelet/bubbletea"

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
	"github.com/chuma-beep/mathua/ui/tui"
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
		m := tui.New(eng, repo)
		if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
			log.Fatalf("tui: %v", err)
		}
	}
}
