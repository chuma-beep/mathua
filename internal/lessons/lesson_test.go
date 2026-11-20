package lessons

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoad(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "lessons.json"), `[
		{"concept_id": "trig.basics.sin_cos_def", "source": "advanced/sine-and-cosine.md"},
		{"concept_id": "trig.basics.unit_circle", "source": "advanced/unit-circle.md"}
	]`)
	os.MkdirAll(filepath.Join(dir, "advanced"), 0755)
	writeFile(t, filepath.Join(dir, "advanced/sine-and-cosine.md"), "# Sine and Cosine\n\nSine and cosine are fundamental trigonometric functions.")
	writeFile(t, filepath.Join(dir, "advanced/unit-circle.md"), "# Unit Circle\n\nThe unit circle is a circle of radius 1.")

	loader, err := Load(dir)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if loader.Count() != 2 {
		t.Errorf("expected 2 lessons, got %d", loader.Count())
	}

	l := loader.Lesson("trig.basics.sin_cos_def")
	if l == nil {
		t.Fatal("expected lesson for trig.sin_cos_def")
	}
	if l.Title != "Sine and Cosine" {
		t.Errorf("expected title 'Sine and Cosine', got %q", l.Title)
	}
	if l.Body == "" {
		t.Error("expected non-empty body")
	}
}

func TestLoad_MissingMappingFile(t *testing.T) {
	dir := t.TempDir()
	_, err := Load(dir)
	if err == nil {
		t.Fatal("expected error for missing lessons.json")
	}
}

func TestLoad_MissingLessonFile(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "lessons.json"), `[
		{"concept_id": "c", "source": "nonexistent.md"}
	]`)
	loader, err := Load(dir)
	if err != nil {
		t.Fatal("expected no error, missing files are skipped")
	}
	if loader.Lesson("c") != nil {
		t.Error("expected nil lesson for missing file")
	}
}

func TestLoad_SharedSource(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "lessons.json"), `[
		{"concept_id": "a", "source": "shared.md"},
		{"concept_id": "b", "source": "shared.md"}
	]`)
	writeFile(t, filepath.Join(dir, "shared.md"), "# Shared\n\nContent.")
	loader, err := Load(dir)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if loader.Count() != 2 {
		t.Errorf("expected 2 concept entries, got %d", loader.Count())
	}
	if loader.Lesson("a") != loader.Lesson("b") {
		t.Error("expected shared lesson object for concepts with same source")
	}
}

// Deterministic multi-source pick: teaching/* beats algebrica/* per concept.
func TestLoad_MultiSource_TeachingWins(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "lessons.json"), `[
		{"concept_id": "c", "source": "algebrica/functions/x.md"},
		{"concept_id": "c", "source": "teaching/c.md"}
	]`)
	writeFile(t, filepath.Join(dir, "algebrica/functions/x.md"), "# Algebrica X\n\nCorpus.")
	writeFile(t, filepath.Join(dir, "teaching/c.md"), "# Teaching C\n\nAuthored.")
	loader, err := Load(dir)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	l := loader.Lesson("c")
	if l == nil {
		t.Fatal("expected lesson")
	}
	if l.Title != "Teaching C" {
		t.Errorf("expected authored teaching lesson to win, got %q", l.Title)
	}
}

func TestLoad_MultiSource_Deterministic(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "lessons.json"), `[
		{"concept_id": "c", "source": "a.md"},
		{"concept_id": "c", "source": "b.md"}
	]`)
	writeFile(t, filepath.Join(dir, "a.md"), "# A\n\nA.")
	writeFile(t, filepath.Join(dir, "b.md"), "# B\n\nB.")
	l1, _ := Load(dir)
	l2, _ := Load(dir)
	if l1.Lesson("c").Title != l2.Lesson("c").Title {
		t.Error("expected deterministic pick across loads")
	}
	if l1.Lesson("c").Title != "B" {
		t.Errorf("expected reverse-lexicographic winner B, got %q", l1.Lesson("c").Title)
	}
}

func TestLoad_NoTitle(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "lessons.json"), `[
		{"concept_id": "c", "source": "notitle.md"}
	]`)
	writeFile(t, filepath.Join(dir, "notitle.md"), "Just some content without a heading.")
	loader, err := Load(dir)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	l := loader.Lesson("c")
	if l.Title != "Notitle" {
		t.Errorf("expected %q, got %q", "Notitle", l.Title)
	}
}

func TestLoad_NotFound(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "lessons.json"), `[]`)
	loader, _ := Load(dir)
	if loader.Lesson("nonexistent") != nil {
		t.Error("expected nil for unknown concept")
	}
}

func TestLoad_RealData(t *testing.T) {
	if _, err := os.Stat("../../data/lessons/algebrica/trigonometry/sine-and-cosine.md"); err != nil {
		if os.Getenv("CI") != "" {
			t.Fatalf("Algebrica corpus missing in CI: run git submodule update --init (%v)", err)
		}
		t.Skip("Algebrica submodule not checked out: run git submodule update --init")
	}
	loader, err := Load("../../data/lessons")
	if err != nil {
		t.Fatalf("load real data: %v", err)
	}
	if loader.Count() == 0 {
		t.Fatal("expected lessons to be loaded")
	}
	// Verify a known mapping works.
	l := loader.Lesson("trig.basics.sin_cos_def")
	if l == nil {
		t.Fatal("expected trig.sin_cos_def lesson")
	}
	if l.Title == "" {
		t.Error("expected non-empty title")
	}
	if l.Body == "" {
		t.Error("expected non-empty body")
	}
	// Shared source test: trig.sin_cos_def and trig.graph_sin share sine-and-cosine.md.
	l2 := loader.Lesson("trig.graph.sin")
	if l2 != l {
		t.Error("expected same lesson object for trig.sin_cos_def and trig.graph_sin")
	}
}

func TestLoad_KPs(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "lessons.json"), `[
		{"concept_id": "arith.add.single", "source": "teaching/arith.add.single.md"},
		{"concept_id": "arith.add.carry", "source": "teaching/arith.add.carry.md"}
	]`)
	os.MkdirAll(filepath.Join(dir, "teaching"), 0755)
	writeFile(t, filepath.Join(dir, "teaching/arith.add.single.md"), "## Add Whole Numbers\n\nContent.\n")
	writeFile(t, filepath.Join(dir, "teaching/arith.add.carry.md"), "## Add Whole Numbers\n\nMore.\n")
	writeFile(t, filepath.Join(dir, "kp/arith.add.single.json"), `[
		{"label": "Use addition notation", "section": "Add Whole Numbers",
		 "subgoals": ["Identify the addends", "Read 3 + 4 as three plus four"]}
	]`)

	loader, err := Load(dir)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	kps := loader.KPs("arith.add.single")
	if len(kps) != 1 {
		t.Fatalf("expected 1 shard, got %d", len(kps))
	}
	if kps[0].Label != "Use addition notation" || len(kps[0].Subgoals) != 2 {
		t.Errorf("unexpected shard: %+v", kps[0])
	}
}

func TestLoad_KPs_FallbackToLessonTitle(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "lessons.json"), `[
		{"concept_id": "c", "source": "shared.md"}
	]`)
	writeFile(t, filepath.Join(dir, "shared.md"), "# Shared\n\nContent.")
	loader, err := Load(dir)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	// No kp/ dir at all, fallback shard carries the lesson title.
	kps := loader.KPs("c")
	if len(kps) != 1 {
		t.Fatalf("expected fallback shard, got %d", len(kps))
	}
	if kps[0].Label != "Shared" {
		t.Errorf("expected fallback label 'Shared', got %q", kps[0].Label)
	}
	if loader.KPs("missing") != nil {
		t.Error("expected nil KPs for unknown concept")
	}
}

func TestLoad_KPSectionBody(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "lessons.json"), `[
		{"concept_id": "c", "source": "lesson.md"}
	]`)
	writeFile(t, filepath.Join(dir, "lesson.md"), `# Lesson

## Use Addition Notation

The operation of addition combines numbers.

### Example

Add 3 and 4.

## Next Section

Different content.
`)
	loader, err := Load(dir)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	body, ok := loader.KPSectionBody("c", "Use Addition Notation")
	if !ok {
		t.Fatal("expected section found")
	}
	if !strings.Contains(body, "operation of addition") || strings.Contains(body, "Next Section") {
		t.Errorf("expected section body bounded by next heading, got:\n%s", body)
	}
	_, ok = loader.KPSectionBody("c", "Does Not Exist")
	if ok {
		t.Error("expected not-found for unknown heading")
	}
}

func TestExtractSection_NestedSubheadingBoundary(t *testing.T) {
	body := `# L

### A

One.

#### Sub

Two.

### B

Three.
`
	sec, ok := extractSection(body, "A")
	if !ok {
		t.Fatal("expected section found")
	}
	// A sub-heading (####) belongs to the section; the next ### terminates it.
	if !strings.Contains(sec, "One.") || !strings.Contains(sec, "Two.") || strings.Contains(sec, "Three.") {
		t.Errorf("unexpected boundary handling, got:\n%s", sec)
	}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("write file: %v", err)
	}
}

func TestExtractSection_MathHeadingMatchesRawShardSection(t *testing.T) {
	// Served headings are canonicalized (`$...$`); shard sections stay raw
	// Algebrica (`\\(...\\)`). They must resolve to the section, not "".
	body := "## Intro\n\nplain\n\n## The substitution $ t = \\tan(x/2)$\n\nwork\n\n## Next\n\ntail\n"
	sec, ok := extractSection(body, `The substitution \\( t = \tan(x/2) \\)`)
	if !ok {
		t.Fatal("raw shard section did not resolve against canonical heading")
	}
	if !strings.Contains(sec, "work") || strings.Contains(sec, "tail") {
		t.Errorf("wrong slice bounds, got:\n%s", sec)
	}
}

func TestNormSectionKey_PlainHeadingsUnchanged(t *testing.T) {
	if got := normSectionKey("Translating sine, cosine, and the differential"); got != "Translating sine, cosine, and the differential" {
		t.Errorf("plain heading altered: %q", got)
	}
}

func TestExtractTitle_SkipsDeepNoteHeadings(t *testing.T) {
	// H4+ headings are styled notes (PreTeXt small print), never titles:
	// the old code returned the whole note paragraph as the lesson title.
	md := "## Introduction\n\n###### Each of these types will be examined.\n\n## Recall of continuity\n"
	if got := extractTitle(md); got != "Recall of continuity" {
		t.Errorf("got %q", got)
	}
}

func TestExtractTitle_StripsMathAndLinks(t *testing.T) {
	md := "## Integrals of trigonometric powers with \\(n\\) even\n"
	if got := extractTitle(md); got != "Integrals of trigonometric powers with n even" {
		t.Errorf("got %q", got)
	}
	md = "## Proof of [Rolle's Theorem](<../rolles-theorem/>)\n"
	if got := extractTitle(md); got != "Proof of Rolle's Theorem" {
		t.Errorf("got %q", got)
	}
}

func TestStripTitleMarkup(t *testing.T) {
	cases := map[string]string{
		"The $p$-Integral Test":                   "The p-Integral Test",
		"See [limits](<../limits/>) next":         "See limits next",
		"Plain title":                             "Plain title",
		"The substitution \\( t = \\tan(x/2) \\)": "The substitution t = \\tan(x/2)",
	}
	// NOTE: lone `$` prices would lose their sign (`Cost $5` -> `Cost 5`).
	// No lesson heading in the corpus carries a price (verified), so the
	// simple strip stands; revisit if one ever appears.
	for in, want := range cases {
		if got := stripTitleMarkup(in); got != want {
			t.Errorf("stripTitleMarkup(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestLoader_AssetsForAndKPAssets(t *testing.T) {
	lessonAssets := []Asset{
		{ID: "fig-1", Caption: "First Figure", Section: "Overview"},
		{ID: "fig-2", Caption: "Second Figure", Section: "Details"},
		{ID: "fig-3", Caption: "Third Figure", Section: "Overview"},
	}

	l := &Loader{
		concepts: map[string]*Lesson{
			"calc.derivatives": {
				Title:  "Derivatives",
				Assets: lessonAssets,
			},
		},
		kps: map[string][]KP{},
	}

	// AssetsFor existing and nonexistent
	if got := l.AssetsFor("calc.derivatives"); len(got) != 3 {
		t.Fatalf("expected 3 assets for calc.derivatives, got %d", len(got))
	}
	if got := l.AssetsFor("nonexistent"); got != nil {
		t.Fatalf("expected nil assets for nonexistent concept, got %v", got)
	}

	// KPAssets with nonexistent concept
	if got := l.KPAssets("nonexistent", KP{}); got != nil {
		t.Fatalf("expected nil for nonexistent concept in KPAssets, got %v", got)
	}

	// KPAssets with authored asset_ids
	kpWithIDs := KP{
		AssetIDs: []string{"fig-3", "fig-1", "nonexistent-id"},
	}
	gotIDs := l.KPAssets("calc.derivatives", kpWithIDs)
	if len(gotIDs) != 2 || gotIDs[0].ID != "fig-3" || gotIDs[1].ID != "fig-1" {
		t.Fatalf("expected fig-3 and fig-1 matching in order, got %+v", gotIDs)
	}

	// KPAssets with empty section
	kpEmptySec := KP{Section: "   "}
	if got := l.KPAssets("calc.derivatives", kpEmptySec); got != nil {
		t.Fatalf("expected nil assets for empty section, got %v", got)
	}

	// KPAssets with matching section
	kpSec := KP{Section: "Overview"}
	gotSec := l.KPAssets("calc.derivatives", kpSec)
	if len(gotSec) != 2 || gotSec[0].ID != "fig-1" || gotSec[1].ID != "fig-3" {
		t.Fatalf("expected Overview assets fig-1 and fig-3, got %+v", gotSec)
	}
}

func TestLoader_LessonByTitleAndLessonsByDomain(t *testing.T) {
	lessonA := &Lesson{
		Title:    "Quadratic Equations",
		Concepts: []string{"algebra.quadratics.factoring", "algebra.quadratics.formula"},
	}
	lessonB := &Lesson{
		Title:    "Single Part Concept",
		Concepts: []string{"single"},
	}
	lessonC := &Lesson{
		Title:    "No Concepts Registered",
		Concepts: nil,
	}

	l := &Loader{
		concepts: map[string]*Lesson{
			"algebra.quadratics.factoring": lessonA,
			"algebra.quadratics.formula":   lessonA, // duplicate pointer sharing lesson
			"single":                       lessonB,
			"empty":                        lessonC,
		},
	}

	// LessonByTitle
	if got := l.LessonByTitle("Quadratic Equations"); got != lessonA {
		t.Fatalf("expected lessonA for title lookup, got %v", got)
	}
	if got := l.LessonByTitle("Nonexistent Title"); got != nil {
		t.Fatalf("expected nil for nonexistent title lookup, got %v", got)
	}

	// All()
	all := l.All()
	if len(all) != 4 {
		t.Fatalf("expected 4 entries in All(), got %d", len(all))
	}

	// LessonsByDomain
	byDomain := l.LessonsByDomain()
	if len(byDomain["algebra.quadratics"]) != 1 || byDomain["algebra.quadratics"][0] != lessonA {
		t.Fatalf("expected deduplicated lessonA under algebra.quadratics, got %+v", byDomain["algebra.quadratics"])
	}
	if len(byDomain["general"]) != 2 {
		t.Fatalf("expected 2 general lessons (single and empty), got %d", len(byDomain["general"]))
	}
}

func TestCorpusOf(t *testing.T) {
	cases := []struct {
		path string
		want string
	}{
		{"data/lessons/algebrica/functions/log.md", "algebrica"},
		{"/root/lessons/teaching/trig.md", "teaching"},
		{"lessons/root-file.md", "authored"},
		{"standalone/overview.md", "authored"},
		{"", "authored"},
	}
	for _, c := range cases {
		if got := corpusOf(c.path); got != c.want {
			t.Errorf("corpusOf(%q) = %q, want %q", c.path, got, c.want)
		}
	}
}

func TestExtractTitle_EdgeCases(t *testing.T) {
	if got := extractTitle(""); got != "" {
		t.Errorf("expected empty string for empty markdown, got %q", got)
	}
	if got := extractTitle("No headings in this body at all."); got != "" {
		t.Errorf("expected empty string when no headings exist, got %q", got)
	}
	if got := extractTitle("###### Only h6 headings\nParagraph"); got != "" {
		t.Errorf("expected empty string when only h6 exists, got %q", got)
	}
	if got := extractTitle("#NoSpaceHeading\nParagraph"); got != "" {
		t.Errorf("expected empty string when heading hash lacks space, got %q", got)
	}
}

func TestLoadKPs_DirectoriesAndCorruptFiles(t *testing.T) {
	dir := t.TempDir()
	kpDir := filepath.Join(dir, "kp")
	os.MkdirAll(kpDir, 0755)

	// Valid shard
	writeFile(t, filepath.Join(kpDir, "calc.derivatives.json"), `[{"label": "Derivative Definition"}]`)
	// Empty shard array
	writeFile(t, filepath.Join(kpDir, "calc.empty.json"), `[]`)
	// Corrupt json shard
	writeFile(t, filepath.Join(kpDir, "calc.corrupt.json"), `not json`)
	// Non-json file
	writeFile(t, filepath.Join(kpDir, "readme.txt"), `instructions`)
	// Subdirectory
	os.MkdirAll(filepath.Join(kpDir, "subdir"), 0755)

	kps := loadKPs(dir, map[string]bool{})
	if len(kps) != 1 {
		t.Fatalf("expected 1 valid shard loaded, got %d", len(kps))
	}
	if len(kps["calc.derivatives"]) != 1 || kps["calc.derivatives"][0].Label != "Derivative Definition" {
		t.Fatalf("unexpected loaded kp: %+v", kps["calc.derivatives"])
	}

	// Nonexistent kp directory returns nil
	if missing := loadKPs(t.TempDir(), map[string]bool{}); missing != nil {
		t.Fatalf("expected nil for missing kp directory, got %v", missing)
	}
}
