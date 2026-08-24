// latexdump walks data/lessons and prints every lesson body after LaTeX
// canonicalization — exactly what the API serves — as JSON. Used by the
// frontend corpus gate (web/next-app/test/latexCorpus.test.ts) so KaTeX is
// tested against production content rather than raw scrape dialects.
package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/chuma-beep/mathua/internal/latex"
)

func main() {
	root := "data/lessons"
	if len(os.Args) > 1 {
		root = os.Args[1]
	}
	out := map[string]string{}
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(d.Name(), ".md") {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		raw, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		content := latex.Canonicalize(string(raw), latex.ForSource(filepath.ToSlash(rel), string(raw)))
		out[filepath.ToSlash(rel)] = content
		return nil
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "latexdump:", err)
		os.Exit(1)
	}
	enc := json.NewEncoder(os.Stdout)
	if err := enc.Encode(out); err != nil {
		fmt.Fprintln(os.Stderr, "latexdump encode:", err)
		os.Exit(1)
	}
}
