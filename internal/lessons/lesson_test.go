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
	// No kp/ dir at all — fallback shard carries the lesson title.
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
