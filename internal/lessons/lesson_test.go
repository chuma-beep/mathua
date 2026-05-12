package lessons

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "lessons.json"), `[
		{"concept_id": "trig.sin_cos_def", "source": "advanced/sine-and-cosine.md"},
		{"concept_id": "trig.unit_circle", "source": "advanced/unit-circle.md"}
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

	l := loader.Lesson("trig.sin_cos_def")
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
	_, err := Load(dir)
	if err == nil {
		t.Fatal("expected error for missing lesson file")
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
	if l.Title != "" {
		t.Errorf("expected empty title, got %q", l.Title)
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
	if _, err := os.Stat("../../data/lessons/advanced/trigonometry/sine-and-cosine.md"); err != nil {
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
	l := loader.Lesson("trig.sin_cos_def")
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
	l2 := loader.Lesson("trig.graph_sin")
	if l2 != l {
		t.Error("expected same lesson object for trig.sin_cos_def and trig.graph_sin")
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
