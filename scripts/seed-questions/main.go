package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/chuma-beep/mathua/internal/concepts"
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/abstract"
	"github.com/chuma-beep/mathua/internal/generator/algebra"
	"github.com/chuma-beep/mathua/internal/generator/arithmetic"
	"github.com/chuma-beep/mathua/internal/generator/calculus"
	"github.com/chuma-beep/mathua/internal/generator/complex"
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
	"github.com/chuma-beep/mathua/internal/storage"
)

type externalQ struct {
	ConceptID   string  `json:"concept_id"`
	Question    string  `json:"question"`
	Answer      string  `json:"answer"`
	Explanation string  `json:"explanation"`
	Source      string  `json:"source"`
	Difficulty  float64 `json:"difficulty"`
}

func main() {
	dbPath := flag.String("db", "mathua.db", "path to SQLite database")
	count := flag.Int("n", 8, "questions per concept")
	importJSON := flag.String("import", "", "import questions from JSON file instead of generating")
	conceptsDir := flag.String("concepts", "data/concepts", "path to concepts directory")
	flag.Parse()

	repo, err := storage.NewSQLiteStore(*dbPath)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer repo.Close()

	if *importJSON != "" {
		data, err := os.ReadFile(*importJSON)
		if err != nil {
			log.Fatalf("read json: %v", err)
		}
		var external []externalQ
		if err := json.Unmarshal(data, &external); err != nil {
			log.Fatalf("parse json: %v", err)
		}
		questions := make([]storage.Question, len(external))
		for i, eq := range external {
			questions[i] = storage.Question{
				ConceptID:   eq.ConceptID,
				Question:    eq.Question,
				Answer:      eq.Answer,
				Explanation: eq.Explanation,
				Source:      eq.Source,
				Difficulty:  eq.Difficulty,
			}
		}
		if err := repo.ImportQuestions(questions); err != nil {
			log.Fatalf("import: %v", err)
		}
		fmt.Printf("imported %d questions from %s\n", len(questions), *importJSON)
		return
	}

	dag, err := concepts.LoadDir(*conceptsDir)
	if err != nil {
		log.Fatalf("load concepts: %v", err)
	}

	reg := generator.NewRegistry()
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

	fmt.Printf("loaded %d concepts, %d generators\n", dag.Count(), reg.Count())

	var seeded int
	for _, c := range dag.Order() {
		if !reg.Has(c.ID) {
			continue
		}
		existing, err := repo.GetQuestionCount(c.ID)
		if err != nil {
			log.Printf("check count for %s: %v", c.ID, err)
			continue
		}
		if existing >= *count {
			continue
		}
		needed := *count - existing
		problems, err := reg.BatchGenerate(c.ID, needed, 0.5)
		if err != nil {
			log.Printf("generate for %s: %v", c.ID, err)
			continue
		}
		questions := make([]storage.Question, len(problems))
		for i, p := range problems {
			questions[i] = storage.Question{
				ConceptID:   c.ID,
				Question:    p.Question,
				Answer:      p.Answer,
				Explanation: p.Explanation,
				Source:      "generator",
				Difficulty:  0.5,
			}
		}
		if err := repo.ImportQuestions(questions); err != nil {
			log.Printf("import for %s: %v", c.ID, err)
			continue
		}
		seeded += len(questions)
		fmt.Printf("  seeded %d questions for %s (%s)\n", len(questions), c.ID, c.Label)
	}

	fmt.Printf("\ndone! seeded %d total questions\n", seeded)
}

func init() {
	// Change working directory to project root (two levels up from script)
	// This expects the script to be run from the project root via `go run scripts/seed-questions/main.go`
	dir, err := os.Getwd()
	if err != nil {
		return
	}
	if filepath.Base(dir) == "seed-questions" {
		os.Chdir(filepath.Dir(filepath.Dir(dir)))
	}
}
