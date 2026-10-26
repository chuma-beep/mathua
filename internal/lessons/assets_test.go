package lessons

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// Every figure a lesson references must be addressable as data.
//
// This is the property that could not be expressed before: an image lived only as literal
// markdown inside an opaque body, so a figure that failed to resolve and a figure that
// rendered correctly were indistinguishable. If this test fails, a body has grown a
// reference shape the extractor does not understand and the consequence would be a
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
		//, a single match can swallow the figures in between and report 148 for a corpus that
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
		t.Errorf("extracted %d assets for %d well-formed references across %d lesson files, "+
			"a figure the extractor misses silently disappears from /learn",
			extracted, refs, bodies)
	}
	t.Logf("%d lesson files, %d figure references, %d extracted", bodies, refs, extracted)
}

// Every figure a lesson references must resolve to a file that is actually shipped.
//
// This is the gate `scripts/audit_lessons.py` could not write. It validated the
// hand-maintained concept->diagram map in `internal/engine`, which says nothing about the
// figures the lesson bodies actually cite, so a body could reference a missing diagram
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
// , CommonMark allows balanced brackets inside alt text, and react-markdown renders it. An
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

// A diagram's accessible name must describe the diagram.
//
// `LessonDiagram` wraps an inlined SVG in `role="img"` with an `aria-label` from the alt text,
// so a screen reader is unaffected either way, but `sanitizeSvg` returns early when the SVG
// already has a `<title>`, so whatever that title says is what a browser shows as the hover
// tooltip over the figure. The vendored algebrica SVGs carried `<title>Group 42</title>` and
// friends from whatever batch generated them: 68 of the 69 that had a title, with the numbers
// reused across unrelated diagrams (seven files all titled "Group 6"). A learner hovering the
// distributive-law figure in the sets lesson was told "Group 2".
//
// Those titles are stripped so the runtime injects one from the alt text, which is the
// figure's own words. This gate exists so the next content wave cannot put them back.
func TestDiagramTitlesDescribeTheirDiagram(t *testing.T) {
	root := "../../web/next-app/public/diagrams"
	if _, err := os.Stat(root); err != nil {
		t.Skipf("diagram directory not present: %v", err)
	}
	// "Group 42" is the batch artifact. A bare filename is the same failure in milder form:
	// accurate, and useless as a description.
	placeholder := regexp.MustCompile(`(?i)^group\s*\d*$`)
	var bad []string
	var titled int

	err := filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || filepath.Ext(p) != ".svg" {
			return err
		}
		raw, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		m := titleRe.FindSubmatch(raw)
		if m == nil {
			return nil
		}
		titled++
		text := strings.TrimSpace(string(m[1]))
		stem := strings.TrimSuffix(filepath.Base(p), ".svg")
		if placeholder.MatchString(text) || text == stem {
			bad = append(bad, fmt.Sprintf("%s: <title>%s</title>", relative(p), text))
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk diagrams: %v", err)
	}
	sort.Strings(bad)
	if len(bad) > 0 {
		limit := len(bad)
		if limit > 10 {
			limit = 10
		}
		t.Errorf("%d of %d diagram titles are placeholders or bare filenames, e.g. %v\n"+
			"A placeholder title is what a browser shows on hover, and it suppresses the "+
			"accessible title the renderer would otherwise inject from the alt text.",
			len(bad), titled, bad[:limit])
	}
	t.Logf("%d diagram SVGs, %d carry a <title>, %d are placeholders", len(allSvgs(t, root)), titled, len(bad))
}

var titleRe = regexp.MustCompile(`(?s)<title[^>]*>(.*?)</title>`)

func allSvgs(t *testing.T, root string) []string {
	t.Helper()
	var out []string
	_ = filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() && filepath.Ext(p) == ".svg" {
			out = append(out, p)
		}
		return nil
	})
	return out
}

// The rune-at-a-time scan in BrokenImageRefs depends on two properties, and both are asserted
// here rather than left as reasoning.
//
// This is *not* a regression test for a bug: the previous byte-at-a-time loop was lossless,
// it compared bytes for equality and appended the index it was standing on, so it never
// truncated anything. ADR-038 records the same false alarm on `cmd/latexfix`, where a test
// written to prove byte iteration corrupted UTF-8 passed against the unfixed code. What the
// change buys is a guarantee that is local to the loop instead of resting on the encoding.
//
// The two properties: a continuation byte can never be `[`, and the offsets returned are byte
// offsets, because the caller turns them into line numbers by counting newlines.
func TestBrokenImageRefsReportsByteOffsetsAroundMultibyteText(t *testing.T) {
	body := "α β γ, an em dash and 2π/3, plus !π not a bracket.\n" +
		"\n" +
		"![legit](/diagrams/algebrica/a.svg)\n" +
		"\n" +
		"a bare ![ that never closes\n"

	off := BrokenImageRefs(body)
	if len(off) != 1 {
		t.Fatalf("got %d broken references, want 1: %v", len(off), off)
	}
	// The `!π` on line 1 must not be mistaken for `![`, and the well-formed reference on
	// line 3 must not be reported.
	if got := string(body[off[0] : off[0]+2]); got != "![" {
		t.Errorf("reported offset points at %q, want %q", got, "![")
	}
	if line := strings.Count(body[:off[0]], "\n") + 1; line != 5 {
		t.Errorf("offset resolved to line %d, want 5, offsets must be byte offsets", line)
	}
}

// A continuation byte is >= 0x80, and `[` is 0x5B, so a multi-byte character can never be
// mistaken for the opening bracket of an image reference. This is the property the loop's
// `body[i+1]` lookup rests on, so it is asserted rather than argued.
func TestContinuationByteIsNeverAnOpeningBracket(t *testing.T) {
	for r := rune(0x80); r <= 0x10FFFF; r++ {
		if r >= 0xD800 && r <= 0xDFFF { // surrogates are not encodable in UTF-8
			continue
		}
		for _, b := range []byte(string(r)) {
			if b >= 0x80 && b == '[' {
				t.Fatalf("rune %U encodes a byte equal to '['", r)
			}
		}
	}
}

// Figure captions left behind by a dropped image may not grow.
//
// 153 paragraphs in the teaching corpus are a figure caption with no figure beside it: the
// shape `ingest_pretext.mjs` emits for a `<figure>` when it keeps the caption and drops the
// image. `scripts/prune_orphan_captions.py` removed the 35 that nothing refers to. The 153 that
// remain are kept on purpose, and the reasons are recorded here because the number looks
// indefensible without them:
//
//	49  the lesson's prose cites the figure by number, so the caption cannot go without
//	    leaving the prose pointing at nothing. These returned with the Precalculus 2e
//	    re-ingest of 2026-10-17, which restored the two sections' captions along with their
//	    text. They are work rather than a judgement: a cited caption is one to reclaim, and
//	    `docs/figure-reclaim.md` section 3 is where they are tracked.
//	74  no picture vocabulary, so the text could as readily be emphasis. Several are real
//	    captions ("Share of all income held by the top 1%, United States") and the rule that
//	    would catch them also catches content that must stay.
//	30  bare section labels such as `*Solution*`, which head a worked answer. The audit's
//	    "duplicated in this lesson" signal rates them `strong`, because a lesson with two
//	    worked answers has two of them, so a tier-based rule deletes a heading from 15
//	    geometry lessons.
//
// The third reason is the lesson's prose citing the figure, by number ("as shown in
// Figure 1.2") or by position ("the graph above illustrates"), where removing the caption would
// leave the prose pointing at nothing. It kept no row between 2026-10-16 and 2026-10-17,
// because every caption a lesson cited by number had been reclaimed and the tier was empty; the
// Precalculus 2e re-ingest then restored 49 such captions with the two trigonometry sections,
// and they are the first bullet above. See `docs/figure-reclaim.md` section 3.
//
// The figure is not the same set of rows it was. Ninety-eight moves have changed it: reclaiming
// Figures 1.2 and then 1.7-1.9 removed 24 rows from the first bucket, and exempting
// auto-numbered labels from the list/table/quote edge test in `orphanCaptions` put back 24 that
// had been hidden for sitting above a table or a bullet list, so the total stayed at 324 while
// its composition changed; and reclaiming Figures 1.4, 1.5 and 1.6 removed 18 more, to 306,
// Figure 1.12 removed 6, to 300, Figures 4.2, 4.3 and 4.5 removed 18, to 282, Figures 6.3 and
// 6.5 removed 6, to 276, Figure 5.5 removed 4, to 272, Figure 5.6 removed 3, to 269, Figure 5.7
// removed 5, to 264, Figure 2.12 removed 1, to 263, Figures 2.13 and 2.14 removed 2, to 261, and
// Figures 2.15 and 2.16 removed 2, to 259, Figures 2.17, 2.18 and 2.19 removed 3, to 256, and
// Figures 2.21 and 2.23 removed 2, to 254, and Figures 3.25, 3.26, 3.27 and 3.33 removed 4, to
// 250, Figures 4.40, 4.42 and 4.45 removed 3, to 247, and Figures 4.50, 4.51 and 4.52 removed 3
// more, to 244, and Figures 4.53, 4.54 and 4.56 removed 3 more, to 241, and Figures 4.58 and
// 4.59 removed 2 more, to 239, and Figures 4.3, 5.32 and 5.33 removed 3 more, to 236, and Figures
// 5.35 and 5.36 removed 2 more, to 234, Figures 6.25 and 6.27 removed 2 more, to 232, Figure
// 6.26 removed 1 more, to 231, Figures 6.28 and 6.29 removed 2 more, to 229, Figure 6.30
// removed 1 more, to 228, Figure 6.31 removed 1 more, to 227, Figure 6.32 removed 1 more, to 226,
// Figure 6.33 removed 1 more, to 225, and Figures 6.35 and 6.36 removed 2 more, to 223, which
// empties the redraw tier, and ten fetched Prealgebra drawings removed 47 more, to 176, and nine
// more fetched Prealgebra drawings removed 31 more, to 145, and twenty-five fetched Calculus
// Volume 1 drawings removed 33 more, to 112, and three fetched Calculus photographs removed 3
// more, to 109, and two fetched Prealgebra photographs removed 5 more, to 104. The Precalculus 2e
// re-ingest of 2026-10-17 then restored 49 cited captions with the two trigonometry sections, to
// 153. The audit's tiers -- 49 `certain`, 33 `strong`, 71 `likely` on 2026-10-17 -- are the
// measurement, and this ceiling is only a ratchet on it, so it moves down with each batch. It
// moved up once, by a deliberate re-ingest, and the sentence you are reading is the record of
// that single movement.
//
// This is a ceiling on a known, explained remainder, not an assertion that 153 is correct.
// The audit is the tool for that judgement, see docs/lesson-media.md and it is a report,
// because the remaining cases are per instance.
func TestOrphanCaptionsDoNotGrow(t *testing.T) {
	root := "../../data/lessons"
	if _, err := os.Stat(root); err != nil {
		t.Skipf("lesson corpus not present: %v", err)
	}
	const known = 153

	var files, total int
	var worst []string
	err := filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(p, ".md") {
			return err
		}
		raw, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		n := len(orphanCaptions(string(raw)))
		if n == 0 {
			return nil
		}
		files++
		total += n
		if n > 0 {
			worst = append(worst, relative(p))
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk corpus: %v", err)
	}
	if total > known {
		t.Errorf("%d orphaned figure captions across %d lessons, ceiling is %d.\n"+
			"A new one means an ingestion run dropped another figure's image and kept its "+
			"caption. `scripts/prune_orphan_captions.py --dry-run` lists them.", total, files, known)
	}
	t.Logf("%d orphaned figure captions across %d lessons (ceiling %d)", total, files, known)
}

// orphanCaptions finds paragraphs that are nothing but emphasis and that sit next to no image.
//
// This deliberately mirrors `section_a_candidates` in scripts/audit_lesson_media.py rather than
// being a convenient second version of the rule. An earlier permissive version counted 416
// against the audit's 324, and the ceiling below is only meaningful if the two agree, so
// agreement is itself the cross-check, and a drift shows up as a failing number rather than as
// a silently weaker gate. The pruning script imports the audit for the same reason: two
// versions of this rule disagreed by 74, and the extra 74 included content that had to stay.
//
// Drift can also come from the image pattern rather than from the caption rule, and it did:
// `imageRe` here has carried a `\]` alternative since it was written, for the two vendored alts
// that contain a bracket, while the audit's pattern stopped at the first `]` until Figure
// 2.16's alt, "over [-2.5, 5]", made the difference visible on 2026-10-16, at 259 here against
// 260 there on the same corpus. The ceiling is an upper bound, so neither number failed and the
// disagreement was found by recomputing both plans by hand; the audit's pattern is what changed.
//
// `numberedCaptionRe` is the audit's `NUMBERED_CAPTION_RE`: an auto-numbered "Figure 1.5" or
// "Table 1.4" label, which PreTeXt emits for a figure and which survives ingestion.
var numberedCaptionRe = regexp.MustCompile(`(?i)^(figure|table|graph|diagram)\s+\d+(\.\d+)*\b`)

func orphanCaptions(body string) []string {
	lines := strings.Split(body, "\n")
	// All four prefixes the audit checks. Omitting "*" counted 398 against the audit's 324,
	// which is exactly the number of captions that sit next to another emphasis line.
	edge := func(l string) bool {
		return strings.HasPrefix(l, "-") || strings.HasPrefix(l, "*") ||
			strings.HasPrefix(l, "|") || strings.HasPrefix(l, ">")
	}
	nonBlank := func(from, step int) string {
		for j := from; j >= 0 && j < len(lines); j += step {
			if t := strings.TrimSpace(lines[j]); t != "" {
				return t
			}
		}
		return ""
	}
	var out []string
	for i, ln := range lines {
		t := strings.TrimSpace(ln)
		// A paragraph that is nothing but emphasis. Excludes a bullet (`* item`) and a
		// bold-only line (`**Header**`), which look the same and are neither.
		if len(t) <= 2 || !strings.HasPrefix(t, "*") || !strings.HasSuffix(t, "*") ||
			strings.HasPrefix(t, "* ") || strings.HasPrefix(t[1:], "*") {
			continue
		}
		// No well-formed image on this line, or on the nearest non-blank line either side.
		if imageRe.MatchString(t) || imageRe.MatchString(nonBlank(i-1, -1)) {
			continue
		}
		// A caption introduced by a list, table or blockquote edge is part of that structure
		// rather than a dropped figure -- unless it carries an auto-numbered Figure/Table
		// label, because that label is decisive and a caption legitimately sits above a table
		// or a list: Figure 1.5 is above the place-value table, Figure 1.6 above the bulleted
		// list of digit values, Figures 4.50 and 4.51 are two captions in a row. Applying the
		// edge test to labelled lines as well hid 24 rows and 10 distinct figures from the
		// audit, measured 2026-10-15; the exemption is in `section_a_candidates` too, and the
		// two have to agree for the ceiling below to mean anything.
		if !numberedCaptionRe.MatchString(t[1:len(t)-1]) && (edge(nonBlank(i-1, -1)) || edge(nonBlank(i+1, 1))) {
			continue
		}
		out = append(out, fmt.Sprintf("%d: %s", i+1, t))
	}
	return out
}

func TestConsecutiveImageRefsAreSeparateAssets(t *testing.T) {
	body := "![The interval [a, b] has two parts](/diagrams/graphs/a.svg)\n" +
		"![The second figure](/diagrams/graphs/b.svg)\n" +
		"![The third figure](/diagrams/graphs/c.svg)"
	got := ExtractAssets(body, nil)
	if len(got) != 3 {
		t.Fatalf("extracted %d assets from three consecutive image references: %v", len(got), got)
	}
	for i, want := range []string{"a.svg", "b.svg", "c.svg"} {
		if got[i].ID != want {
			t.Errorf("asset %d is %q, want %q", i, got[i].ID, want)
		}
	}
}

func TestGraphBasicsLessonKeepsItsFigures(t *testing.T) {
	body, err := os.ReadFile("../../data/lessons/teaching/discrete.graphs.basics.md")
	if err != nil {
		t.Fatalf("read lesson: %v", err)
	}
	refs := regexp.MustCompile(`!\[[^\]]*\]\(/diagrams/graphs/[^)]+\.svg\)`).FindAllString(string(body), -1)
	if len(refs) < 25 {
		t.Fatalf("graph basics lesson has %d figures, want at least 25: the text refers to graphs that are not drawn", len(refs))
	}
	assets := ExtractAssets(string(body), nil)
	if len(assets) != len(refs) {
		t.Fatalf("extracted %d assets for %d figure references: consecutive figures are being swallowed by the image pattern", len(assets), len(refs))
	}
	for _, want := range []string{
		"sec_gt-intro-2-2-1.svg",
		"sec_gt-intro-31-1-2-1.svg",
		"ex-gt-non-connected.svg",
		"sec_gt-intro-50-11-1.svg",
	} {
		if !strings.Contains(string(body), want) {
			t.Fatalf("graph basics lesson lost its %s figure", want)
		}
	}
}
