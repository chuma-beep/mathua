package lessons

import (
	"strings"
	"testing"
)

// Guard: lessons loaded from the real dataset must come out of the latex
// ingestion pipeline canonicalized ($$/$ delimiters, no raw environments).
func TestRealDataCanonicalized(t *testing.T) {
	lib, err := Load("../../data/lessons")
	if err != nil {
		t.Skipf("dataset not present: %v", err)
	}
	if lib.Count() == 0 {
		t.Skip("empty library")
	}

	badArtifacts := []string{"\\begin{align}\\n", "sqrt}{", "\\amp", "<span class=\"math-", "\\begin{equation"}
	expectDisplay := map[string]string{
		"discrete.sequences.recurrence": "",
		"nt.adv.diophantine":            "",
		"calc.integral.indefinite":      "",
		"complex.ops.add_sub":           "\\begin{aligned}",
	}
	seen := map[string]bool{}
	for _, l := range lib.All() {
		for _, c := range l.Concepts {
			for _, bad := range badArtifacts {
				if strings.Contains(l.Body, bad) {
					t.Errorf("%s: stale artifact %q survived ingestion", c, bad)
				}
			}
			want, probed := expectDisplay[c]
			if !probed {
				continue
			}
			if !strings.Contains(l.Body, "$$") {
				t.Errorf("%s: expected $$ display math after canonicalization", c)
			}
			if want != "" && !strings.Contains(l.Body, want) {
				t.Errorf("%s: expected aligned body %q", c, want)
			}
			seen[c] = true
		}
	}
	for c := range expectDisplay {
		if !seen[c] {
			t.Errorf("expected lesson for %s was not loaded", c)
		}
	}
}
