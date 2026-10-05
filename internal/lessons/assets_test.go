package lessons

import (
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// Every figure a lesson references must be addressable as data.
//
// This is the property that could not be expressed before: an image lived only as literal
// markdown inside an opaque body, so a figure that failed to resolve and a figure that
// rendered correctly were indistinguishable. If this test fails, a body has grown a
// reference shape the extractor does not understand — and the consequence would be a
// diagram that silently stops appearing on `/learn`, which is exactly the bug that started
// this.
func TestAssetExtractionCoversEveryImageReference(t *testing.T) {
	root := "../../data/lessons"
	if _, err := os.Stat(root); err != nil {
		t.Skipf("lesson corpus not present: %v", err)
	}
	var bodies, refs, extracted int
	broken := map[string]bool{}

	err := filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(p, ".md") {
			return err
		}
		raw, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		body := string(raw)
		bodies++

		// Counted per paragraph, the same way the extractor finds them, so the two numbers
		// are comparable. Two earlier counting methods were both wrong: per line disagrees for
		// any reference whose alt text wraps, and whole-body disagrees in the other direction
		// — a single match can swallow the figures in between and report 148 for a corpus that
		// holds 292.
		for _, chunk := range paragraphChunks(body) {
			refs += len(imageRe.FindAllStringIndex(chunk.text, -1))
		}
		extracted += len(ExtractAssets(body, []string{"x"}))
		for _, off := range BrokenImageRefs(body) {
			broken[relative(p)+":"+itoa(lineOf(body, off))] = true
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk corpus: %v", err)
	}

	if refs == 0 {
		t.Skip("no image references in the corpus")
	}
	if len(broken) > 0 {
		sample := make([]string, 0, 5)
		for k := range broken {
			sample = append(sample, k)
			if len(sample) == 5 {
				break
			}
		}
		sort.Strings(sample)
		t.Errorf("%d BROKEN image references (bare '![', alt and path lost by the scrape): %v\n"+
			"These render as a literal '![' in the lesson. Four were repaired against figures "+
			"that ship but were never referenced, so a new one means a figure was lost again.",
			len(broken), sample)
	}
	if extracted != refs {
		t.Errorf("extracted %d assets for %d well-formed references across %d lesson files — "+
			"a figure the extractor misses silently disappears from /learn",
			extracted, refs, bodies)
	}
	t.Logf("%d lesson files, %d figure references, %d extracted", bodies, refs, extracted)
}

// Every figure a lesson references must resolve to a file that is actually shipped.
//
// This is the gate `scripts/audit_lessons.py` could not write. It validated the
// hand-maintained concept->diagram map in `internal/engine`, which says nothing about the
// figures the lesson bodies actually cite — so a body could reference a missing diagram
// indefinitely and the audit would report clean. The resolution rule is duplicated from
// `web/next-app/components/KatexContent.tsx:181-192`, and that duplication is deliberate
// rather than incidental: the Go side needs it to gate the corpus, the client needs it to
// render, and each has a test pinning its own copy.
func TestEveryAssetResolvesToAShippedDiagram(t *testing.T) {
	root := "../../data/lessons"
	pub := "../../web/next-app/public"
	if _, err := os.Stat(root); err != nil {
		t.Skipf("lesson corpus not present: %v", err)
	}
	if _, err := os.Stat(pub); err != nil {
		t.Skipf("public dir not present: %v", err)
	}

	// Mirrors the client rule: an absolute or remote URL is used as written, anything else
	// is a bare filename resolved against the algebrica diagram directory.
	resolve := func(src string) string {
		if strings.HasPrefix(src, "http") || strings.HasPrefix(src, "/") {
			return src
		}
		return "/diagrams/algebrica/" + path.Base(src)
	}

	var total, checked, remote int
	var broken []string
	err := filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(p, ".md") {
			return err
		}
		raw, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		for _, a := range ExtractAssets(string(raw), []string{"x"}) {
			total++
			url := resolve(a.Src)
			if strings.HasPrefix(url, "http") {
				remote++
				continue
			}
			checked++
			if _, err := os.Stat(filepath.Join(pub, filepath.FromSlash(url))); err != nil {
				broken = append(broken, relative(p)+" -> "+a.Src+" (resolves to "+url+")")
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk corpus: %v", err)
	}
	if total == 0 {
		t.Skip("no assets in the corpus")
	}
	sort.Strings(broken)
	if len(broken) > 0 {
		limit := len(broken)
		if limit > 10 {
			limit = 10
		}
		t.Errorf("%d of %d assets resolve to no shipped file, e.g. %v\n"+
			"A broken figure renders as a browser placeholder inside the lesson prose.",
			len(broken), total, broken[:limit])
	}
	t.Logf("%d assets: %d resolve to a shipped diagram, %d remote, %d broken",
		total, checked-len(broken), remote, len(broken))
}

// Alt text may contain brackets, and the extractor must agree with the renderer about that.
//
// A caption of "The interval [a, b] has been divided into 6 subintervals" is ordinary Markdown
// — CommonMark allows balanced brackets inside alt text, and react-markdown renders it. An
// extractor that stops at the first `]` therefore disagrees with the renderer, and reports a
// reference the learner sees as a working image as a corrupt one.
func TestAssetAltTextMayContainBrackets(t *testing.T) {
	body := "The rule:\n\n![The interval [a, b] has been divided into 6 subintervals of width h]" +
		"(/diagrams/algebrica/numerical-integration-1.svg)\n\nThen the sum."
	got := ExtractAssets(body, []string{"c"})
	if len(got) != 1 {
		t.Fatalf("extracted %d assets, want 1: an alt containing [a, b] is valid Markdown", len(got))
	}
	if !strings.Contains(got[0].Alt, "[a, b]") {
		t.Errorf("alt = %q, want it to keep the bracketed interval", got[0].Alt)
	}
	if got[0].Src != "/diagrams/algebrica/numerical-integration-1.svg" {
		t.Errorf("src = %q", got[0].Src)
	}
	if len(BrokenImageRefs(body)) != 0 {
		t.Errorf("a bracketed alt was reported as a broken reference")
	}
}

// A figure with neither alt text nor a caption is invisible to a screen reader and
// unattributable. The vendored corpus has a lot of these, which is worth knowing rather
// than assuming.
func TestAssetAccessibilityIsMeasured(t *testing.T) {
	root := "../../data/lessons"
	if _, err := os.Stat(root); err != nil {
		t.Skipf("lesson corpus not present: %v", err)
	}
	var total, withAlt, withCaption, withNeither int
	_ = filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(p, ".md") {
			return err
		}
		raw, _ := os.ReadFile(p)
		for _, a := range ExtractAssets(string(raw), []string{"x"}) {
			total++
			alt, cap := a.Alt != "", a.Caption != ""
			if alt {
				withAlt++
			}
			if cap {
				withCaption++
			}
			if !alt && !cap {
				withNeither++
			}
		}
		return nil
	})
	if total == 0 {
		t.Skip("no assets in the corpus")
	}
	t.Logf("%d assets: %d have alt text, %d have a caption, %d have neither",
		total, withAlt, withCaption, withNeither)
}

func lineOf(body string, off int) int {
	return strings.Count(body[:off], "\n") + 1
}

func relative(p string) string {
	if i := strings.Index(p, "data/lessons/"); i >= 0 {
		return p[i:]
	}
	return filepath.Base(p)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var b []byte
	for n > 0 {
		b = append([]byte{byte('0' + n%10)}, b...)
		n /= 10
	}
	return string(b)
}

// `Source` is provenance for the prose, not a path and not an attribution. Getting it wrong
// is worse than leaving it empty: `algebrica/functions/floor-…md` is a per-topic directory,
// and returning it makes the field look like a corpus while differing per lesson.
func TestCorpusOfReportsTheCorpusNotTheDirectory(t *testing.T) {
	cases := map[string]string{
		"data/lessons/algebrica/functions/floor-and-ceiling-functions.md": "algebrica",
		"data/lessons/algebrica/trigonometry/unit-circle.md":              "algebrica",
		"data/lessons/teaching/arith.add.word.md":                         "teaching",
		"arith.add.word.md": "authored",
		"":                  "authored",
	}
	for in, want := range cases {
		if got := corpusOf(in); got != want {
			t.Errorf("corpusOf(%q) = %q, want %q", in, got, want)
		}
	}
}

// The section a figure was placed in is recomputed from the written body, not taken on trust
// from whatever wrote it. This is the check that a placement is right: a figure that landed a
// section early or late has a different nearest-preceding-heading, so a wrong placement
// changes the value and fails here rather than being noticed by a reader.
func TestPlacedFiguresAreInTheSectionTheirPlacementNamed(t *testing.T) {
	root := "../../data/lessons"
	if _, err := os.Stat(root); err != nil {
		t.Skipf("lesson corpus not present: %v", err)
	}
	type want struct{ lesson, file, section string }
	expect := []want{
		{"roots-of-unity.md", "roots-of-unity-1.svg", "geometric interpretation"},
		{"cauchy-theorem.md", "cauchy-theorem-1.svg", "introduction"},
		{"exponential-function.md", "exponential-function-1.svg", "properties for $a$ greater than one"},
		{"exponential-function.md", "exponential-function-2.svg", "properties for $a$ between zero and one"},
		{"exponential-function.md", "exponential-function-3.svg", "properties for $a$ equal to one"},
		{"polynomial-function.md", "polynomial-function-1.svg", "degree 1: linear functions"},
		{"polynomial-function.md", "polynomial-function-2.svg", "degree 2: quadratic functions"},
		{"polynomial-function.md", "polynomial-function-3.svg", "degree 3: cubic functions"},
		{"trigonometric-inequalities.md", "trigonometric-inequalities-1.svg", "inequalities involving sine"},
		{"trigonometric-inequalities.md", "trigonometric-inequalities-2.svg", "inequalities involving cosine"},
		{"trigonometric-inequalities.md", "trigonometric-inequalities-3.svg", "inequalities involving tangent"},
		{"indefinite-integrals.md", "indefinite-integrals-1.svg", "primitives"},
		{"indefinite-integrals.md", "indefinite-integrals-2.svg", "primitives"},
		{"integration-by-parts.md", "integration-by-parts-1.svg", "derivation of the formula"},
		{"numerical-integration.md", "numerical-integration-1.svg", ""},
		{"numerical-integration.md", "numerical-integration-2.svg", ""},
		{"polynomials.md", "polynomials-1.svg", "degree of a polynomial and its geometric interpretation"},
		{"polynomials.md", "polynomials-2.svg", "polynomial equations"},
		{"polynomials.md", "polynomials-3.svg", "end behavior of polynomial"},
		{"trinomials.md", "trinomials-1.svg", "classification of trinomials"},
		{"logarithms.md", "logarithms-1.svg", "logarithmic function"},
		{"logarithms.md", "logarithms-2.svg", "logarithmic function"},
		{"logarithms.md", "logarithms-3.svg", "fundamental inequality for the natural logarithm"},
		{"radicals.md", "radicals-1.svg", "definition of radicals"},
		{"radicals.md", "radicals-2.svg", "geometric construction of the segment \\(\\sqrt{a}\\)"},
		{"euler-number-limit-sequence.md", "euler-number-limit-sequence-1.svg", ""},
		{"principle-of-mathematical-induction.md", "principle-of-mathematical-induction-1.svg", "mathematical induction"},
		{"power-series.md", "power-series-1.svg", "radius of convergence"},
		{"power-series.md", "power-series-2.svg", "radius of convergence"},
		{"integers.md", "integers-1.svg", "definition"},
		{"integers.md", "integers-2.svg", "definition"},
		{"sets.md", "sets-1.svg", "set operations"},
		{"sets.md", "sets-2.svg", "set operations"},
		{"sets.md", "sets-3.svg", "set operations"},
		{"sets.md", "sets-4.svg", "set operations"},
		{"sets.md", "sets-5.svg", "set operations"},
		{"sets.md", "sets-6.svg", "properties of set operations"},
		{"eigenvalues-and-eigenvectors.md", "eigenvalues-and-eigenvectors-1.svg", "definition"},
	}

	byName := map[string]string{}
	for _, e := range expect {
		found := ""
		_ = filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() || filepath.Base(p) != e.lesson {
				return err
			}
			raw, rerr := os.ReadFile(p)
			if rerr != nil {
				return rerr
			}
			for _, a := range ExtractAssets(string(raw), []string{"x"}) {
				if path.Base(a.Src) == e.file {
					found = a.Section
				}
			}
			return nil
		})
		if found == "" {
			t.Errorf("%s: %s is not referenced by any lesson body", e.lesson, e.file)
			continue
		}
		byName[e.file] = found
		// An empty expectation means the figure's section was not pinned deliberately, which is
		// the case for the three that replaced a paragraph already in the right place.
		if e.section == "" {
			continue
		}
		// Compared through normSectionKey and case-insensitively: the expectation is written
		// as a human reads the heading ("Properties for $a$ greater than one") while the
		// extracted value has been through the same normalisation the section matcher uses.
		if !strings.EqualFold(normSectionKey(e.section), normSectionKey(found)) {
			t.Errorf("%s: %s is in section %q, want %q", e.lesson, e.file, found, e.section)
		}
	}
	if len(expect) != 38 {
		t.Fatalf("expectation list has %d entries; it must cover all 38 placed figures", len(expect))
	}
}
