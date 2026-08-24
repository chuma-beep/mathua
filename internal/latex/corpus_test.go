package latex

import (
	"encoding/json"
	"flag"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

var updateBaseline = flag.Bool("update", false, "rewrite the corpus warning baseline")

// warningKind normalizes a warning to its stable kind so cosmetic changes
// (e.g. the counted number of stray dollars) don't churn the baseline.
var warningCountRe = regexp.MustCompile(`\(\d+\)`)

func warningKind(w string) string {
	return warningCountRe.ReplaceAllString(w, "(N)")
}

// TestCorpusGate walks the real lesson corpus and enforces two properties:
//
//  1. Canonicalize reaches a fixed point (C(C(x)) == C(C(C(x)))) — the first
//     application may legitimately keep rewriting (e.g. heading demotion runs
//     per call), but content must stabilize.
//  2. Validation warnings on the stabilized form match the committed baseline
//     exactly: no file may gain a warning kind, and a file that becomes clean
//     must be dropped from the baseline (update with -update and commit).
//
// Run `go test ./internal/latex -run TestCorpusGate -update` to regenerate
// the baseline after intentional changes.
func firstDiff(a, b string) int {
	n := len(a)
	if len(b) > n {
		n = len(b)
	}
	for i := 0; i < n; i++ {
		var ca, cb byte
		if i < len(a) {
			ca = a[i]
		}
		if i < len(b) {
			cb = b[i]
		}
		if ca != cb {
			return i
		}
	}
	return n
}

func TestCorpusGate(t *testing.T) {
	const corpusRoot = "../../data/lessons"
	if _, err := os.Stat(corpusRoot); err != nil {
		t.Skip("lesson corpus not present")
	}

	var files []string
	err := filepath.WalkDir(corpusRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".md") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk corpus: %v", err)
	}
	if len(files) == 0 {
		t.Fatal("no lesson files found")
	}
	sort.Strings(files)

	baselinePath := filepath.Join("testdata", "warning_baseline.json")
	baseline := map[string][]string{}
	if b, err := os.ReadFile(baselinePath); err == nil {
		if err := json.Unmarshal(b, &baseline); err != nil {
			t.Fatalf("parse baseline: %v", err)
		}
	} else if !*updateBaseline {
		t.Fatal("warning baseline missing; run with -update to create it")
	}

	fresh := map[string][]string{}
	idempotenceFailures := 0

	for _, path := range files {
		rel, err := filepath.Rel(corpusRoot, path)
		if err != nil {
			t.Fatalf("rel %s: %v", path, err)
		}
		rel = filepath.ToSlash(rel)
		raw, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read %s: %v", path, err)
		}
		content := string(raw)
		adapter := ForSource(rel, content)

		c1 := Canonicalize(content, adapter)
		c2 := Canonicalize(c1, adapter)
		c3 := Canonicalize(c2, adapter)
		if c2 != c3 {
			idempotenceFailures++
			at := firstDiff(c2, c3)
			lo := at - 100
			if lo < 0 {
				lo = 0
			}
			hi := at + 100
			if hi > len(c2) {
				hi = len(c2)
			}
			hj := at + 100
			if hj > len(c3) {
				hj = len(c3)
			}
			t.Errorf("%s: Canonicalize does not stabilize at %d:\nC2: %q\nC3: %q", rel, at, c2[lo:hi], c3[lo:hj])
		}

		var kinds []string
		seen := map[string]bool{}
		for _, w := range Validate(c2) {
			k := warningKind(w)
			if !seen[k] {
				seen[k] = true
				kinds = append(kinds, k)
			}
		}
		sort.Strings(kinds)
		if len(kinds) > 0 {
			fresh[rel] = kinds
		}
	}

	if *updateBaseline {
		b, err := json.MarshalIndent(fresh, "", "  ")
		if err != nil {
			t.Fatalf("marshal baseline: %v", err)
		}
		if err := os.MkdirAll("testdata", 0o755); err != nil {
			t.Fatalf("mkdir testdata: %v", err)
		}
		if err := os.WriteFile(baselinePath, append(b, '\n'), 0o644); err != nil {
			t.Fatalf("write baseline: %v", err)
		}
		t.Logf("baseline rewritten: %d files with warnings", len(fresh))
		return
	}

	// New warnings: in fresh but not baseline.
	for file, kinds := range fresh {
		base, ok := baseline[file]
		if !ok {
			t.Errorf("%s: new warnings (file was clean): %v", file, kinds)
			continue
		}
		baseSet := map[string]bool{}
		for _, k := range base {
			baseSet[k] = true
		}
		for _, k := range kinds {
			if !baseSet[k] {
				t.Errorf("%s: new warning kind %q", file, k)
			}
		}
	}
	// Cleanups: in baseline but now clean — require a deliberate baseline update.
	for file := range baseline {
		if _, ok := fresh[file]; !ok {
			t.Errorf("%s: no longer warns; remove it from the baseline (-update)", file)
		}
	}
}
