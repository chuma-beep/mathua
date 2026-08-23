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

	probes := map[string]bool{ // concept -> expects display math
		"discrete.sequences.recurrence": true,
		"nt.adv.diophantine":            true,
	}
	for _, l := range lib.All() {
		for _, c := range l.Concepts {
			wantDisplay, probed := probes[c]
			if !probed {
				continue
			}
			if strings.Contains(l.Body, "\\begin{equation") {
				t.Errorf("%s: raw equation environment survived ingestion", c)
			}
			if wantDisplay && !strings.Contains(l.Body, "$$") {
				t.Errorf("%s: expected $$ display math after canonicalization", c)
			}
			delete(probes, c)
		}
	}
	for c := range probes {
		t.Errorf("expected lesson for %s was not loaded", c)
	}
}
