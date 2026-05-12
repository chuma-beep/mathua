package main

import (
	"fmt"
	"os"

	"github.com/chuma-beep/mathua/internal/concepts"
)

func main() {
	dag, err := concepts.Load("data/concepts.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "FAIL: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("OK: %d concepts across %d domains\n", dag.Count(), len(dag.Domains()))
}
