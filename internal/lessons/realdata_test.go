package lessons

import (
	"fmt"
	"regexp"
	"sort"
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

// Guard: lesson ownership is 1:1 — every concept resolves to exactly one
// distinct lesson, and it is the server pick (Loader.Lesson). The Study
// client resolves ?concept= first-wins over the domain-grouped lessons; as
// long as this holds, first-wins IS the server pick and chips always land
// on a lesson that actually serves the concept. If content ever maps one
// concept to several sources, this fails and forces the `primary` field
// follow-up (client cannot replicate reverse-lex source order).
func TestRealDataSingleLessonPerConcept(t *testing.T) {
	lib, err := Load("../../data/lessons")
	if err != nil {
		t.Skipf("dataset not present: %v", err)
	}
	owners := map[string]map[*Lesson]bool{}
	for _, l := range lib.All() {
		for _, c := range l.Concepts {
			if owners[c] == nil {
				owners[c] = map[*Lesson]bool{}
			}
			owners[c][l] = true
		}
	}
	checked := 0
	for cid, set := range owners {
		checked++
		if len(set) != 1 {
			t.Errorf("%s: served by %d distinct lessons, want 1", cid, len(set))
			continue
		}
		for l := range set {
			if l != lib.Lesson(cid) {
				t.Errorf("%s: sole lesson is not the server pick", cid)
			}
		}
	}
	if checked == 0 {
		t.Fatal("no lessons checked")
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

// The worked example `/learn` serves is never the lesson body.
//
// This is the invariant whose absence let a real defect ship. When a shard's `section` was
// empty the server substituted the whole lesson body, so 69 of 1,971 steps across 43 concepts
// served the entire reference article — 33,784 words, and discrete.logic.propositions served
// its 2,724-word article three times over, against a concept whose mastery threshold is twelve
// seconds. Nothing caught it, for two compounding reasons:
//
//   - `TestRealDataKPSectionsResolve` above `continue`s on `kp.Section == ""`, so all 69 of
//     them sat outside the one gate whose own comment names this failure mode.
//   - Every e2e test mocks `**/api/lessons**`, so the `worked_example` under test was a
//     hand-written literal and never a served body.
//
// So this asserts the served string is *bounded*, not that every section resolves. An authored
// slice cannot reintroduce the article, and neither can a shard whose section stops matching.
func TestWorkedExampleIsBoundedAndNeverTheWholeBody(t *testing.T) {
	lib, err := Load("../../data/lessons")
	if err != nil {
		t.Skipf("dataset not present: %v", err)
	}
	if lib.Count() == 0 {
		t.Skip("empty library")
	}

	var total, sectioned, authored, minimal, firstCardMinimal int
	var overBudget []string
	conceptsWithMinimal := map[string]bool{}

	for _, f := range lib.SliceAudit() {
		total++
		lesson := lib.Lesson(f.ConceptID)
		if lesson == nil {
			continue
		}
		if f.Words > MaxSliceWords {
			// Reported, not failed. A section over budget is normally a section that genuinely
			// needs the room — `abstract.group.dihedral`'s "Two generators and three relations"
			// is 534 words and cannot be truncated without breaking the argument. The budget
			// gates what is *authored* (see TestAuthoredSlicesMeetTheirBoundaries); for content
			// that already exists it is a number worth knowing rather than a rule to enforce,
			// because the only mechanical fix would be cutting a proof in half.
			overBudget = append(overBudget, fmt.Sprintf("%s/%s (%d)", f.ConceptID, f.KPLabel, f.Words))
		}
		if f.FirstCard && f.Level == SliceLevelMinimal {
			firstCardMinimal++
		}
		switch f.Level {
		case SliceLevelAuthored:
			authored++
		case SliceLevelMinimal:
			minimal++
			conceptsWithMinimal[f.ConceptID] = true
		default:
			sectioned++
			// A sectioned step must not be the whole article. This is the direct check on the
			// substitution that caused the defect, and it holds even if the article is short.
			if len(strings.Fields(lesson.Body)) == f.Words {
				t.Errorf("%s: worked example is exactly the lesson body (%d words)",
					f.ConceptID, f.Words)
			}
		}
	}

	t.Logf("%d knowledge points: %d served a lesson section, %d an authored slice, %d a minimal placeholder",
		total, sectioned, authored, minimal)
	if len(overBudget) > 0 {
		sort.Strings(overBudget)
		t.Logf("%d steps exceed the %d-word guideline, e.g. %v (reported, not enforced)",
			len(overBudget), MaxSliceWords, overBudget[:min(4, len(overBudget))])
	}
	t.Logf("%d concepts have at least one minimal placeholder; %d show one on the first card",
		len(conceptsWithMinimal), firstCardMinimal)
	t.Logf("content backlog: %d steps await an authored slice", minimal)
}

// An authored slice is only worth having if it teaches. These are curriculum boundaries rather
// than a word-count gate: a slice with no rule teaches nothing, and one long past the budget
// has stopped being the minimum needed for the next practice item.
func TestAuthoredSlicesMeetTheirBoundaries(t *testing.T) {
	lib, err := Load("../../data/lessons")
	if err != nil {
		t.Skipf("dataset not present: %v", err)
	}
	for _, f := range lib.SliceAudit() {
		if f.Level != SliceLevelAuthored {
			continue
		}
		if f.Words < MinSliceWords {
			t.Errorf("%s %q: slice is %d words, below MinSliceWords=%d — it teaches nothing",
				f.ConceptID, f.KPLabel, f.Words, MinSliceWords)
		}
		if f.Words > MaxSliceWords {
			t.Errorf("%s %q: slice is %d words, over MaxSliceWords=%d",
				f.ConceptID, f.KPLabel, f.Words, MaxSliceWords)
		}
	}
}
