package lessons

import (
	"regexp"
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

// Guard: every KP shard section must resolve to a real slice of the
// concept's canonical lesson body. Unresolved sections silently fall back
// to the FULL body (duplicated study content) — e.g. raw `\\(...)`
// sections vs canonical `$...$` headings (calc.integral.weierstrass_sub).
func TestRealDataKPSectionsResolve(t *testing.T) {
	lib, err := Load("../../data/lessons")
	if err != nil {
		t.Skipf("dataset not present: %v", err)
	}
	checked := 0
	for cid, kps := range lib.kps {
		lesson := lib.concepts[cid]
		if lesson == nil {
			continue
		}
		multiSection := strings.Count(lesson.Body, "\n## ") > 0
		for _, kp := range kps {
			if kp.Section == "" {
				continue
			}
			checked++
			body, ok := lib.KPSectionBody(cid, kp.Section)
			if !ok {
				t.Errorf("%s: section %q unresolved (would serve full body)", cid, kp.Section)
				continue
			}
			if multiSection && body == lesson.Body {
				t.Errorf("%s: section %q fell back to full body", cid, kp.Section)
			}
		}
	}
	if checked == 0 {
		t.Fatal("no KP sections checked")
	}
}

// Guard: every served lesson title must be plain renderable text — the
// title slots (study list/detail, concept page) never run math or markdown.
// Catches paragraph-titles (deep-note H6 picked by a naive first-heading
// scan), raw delimiters, and pasted markdown links.
func TestRealDataTitlesSane(t *testing.T) {
	lib, err := Load("../../data/lessons")
	if err != nil {
		t.Skipf("dataset not present: %v", err)
	}
	badBackslash := regexp.MustCompile(`\\[a-zA-Z]`)
	checked := 0
	for _, l := range lib.All() {
		for _, c := range l.Concepts {
			checked++
			if l.Title == "" {
				t.Errorf("%s: empty lesson title", c)
			}
			if len(l.Title) > 120 {
				t.Errorf("%s: suspiciously long title (%d chars): %q", c, len(l.Title), l.Title[:60])
			}
			if strings.Contains(l.Title, "$") || strings.Contains(l.Title, "](") {
				t.Errorf("%s: raw markup in title: %q", c, l.Title)
			}
			if badBackslash.MatchString(l.Title) {
				t.Errorf("%s: raw TeX command in title: %q", c, l.Title)
			}
		}
	}
	if checked == 0 {
		t.Fatal("no lessons checked")
	}
}
