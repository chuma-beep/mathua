package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	entries, err := os.ReadDir("data/concepts")
	if err != nil {
		fmt.Fprintf(os.Stderr, "read dir: %v\n", err)
		os.Exit(1)
	}
	var all []any
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		data, err := os.ReadFile(filepath.Join("data/concepts", e.Name()))
		if err != nil {
			fmt.Fprintf(os.Stderr, "read %s: %v\n", e.Name(), err)
			os.Exit(1)
		}
		var batch []any
		if err := json.Unmarshal(data, &batch); err != nil {
			fmt.Fprintf(os.Stderr, "parse %s: %v\n", e.Name(), err)
			os.Exit(1)
		}
		all = append(all, batch...)
	}
	out, err := json.MarshalIndent(all, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "marshal: %v\n", err)
		os.Exit(1)
	}
	if err := os.WriteFile("web/next-app/data/concepts.json", out, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "write: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("OK: %d concepts merged into web/next-app/data/concepts.json\n", len(all))
}
