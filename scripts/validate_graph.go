package main

import (
	"fmt"
	"os"

	"github.com/chuma-beep/mathua/internal/concepts"
)

func main() {
	dag, err := concepts.LoadDir("data/concepts")
	if err != nil {
		fmt.Fprintf(os.Stderr, "FAIL: %v\n", err)
		os.Exit(1)
	}
	allow := concepts.ReviewedDuplicateLabels()
	if err := concepts.CheckReviewedDuplicateLabelsExist(dag, allow); err != nil {
		fmt.Fprintf(os.Stderr, "FAIL: %v\n", err)
		os.Exit(1)
	}
	if err := concepts.CheckDuplicateLabels(dag, allow); err != nil {
		fmt.Fprintf(os.Stderr, "FAIL: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("OK: %d concepts across %d domains, %d reviewed label collisions\n",
		dag.Count(), len(dag.Domains()), len(allow))
}
