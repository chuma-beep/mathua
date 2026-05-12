package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/chuma-beep/mathua/internal/concepts"
	"github.com/chuma-beep/mathua/internal/engine"
	"github.com/chuma-beep/mathua/internal/generator"
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
	// Domain generators are registered via sub-package Register() calls.
	// arithmetic.Register(reg)
	// counting.Register(reg)
	// ...

	eng := engine.New(repo, dag, reg)
	fmt.Printf("engine ready — %d generators registered\n", reg.Count())

	if *serve {
		fmt.Printf("starting web server on :%d\n", *port)
		_ = eng
		// startAPIServer(eng, *port)
	} else {
		fmt.Println("TUI mode — not yet implemented")
		_ = eng
	}
}
