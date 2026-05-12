package lessons

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Lesson struct {
	Title    string
	Body     string
	Concepts []string
}

type mapping struct {
	ConceptID string `json:"concept_id"`
	Source    string `json:"source"`
}

type Loader struct {
	concepts map[string]*Lesson
}

func Load(lessonsDir string) (*Loader, error) {
	mapPath := filepath.Join(lessonsDir, "lessons.json")
	data, err := os.ReadFile(mapPath)
	if err != nil {
		return nil, fmt.Errorf("read lessons.json: %w", err)
	}
	var mappings []mapping
	if err := json.Unmarshal(data, &mappings); err != nil {
		return nil, fmt.Errorf("parse lessons.json: %w", err)
	}

	concepts := make(map[string]*Lesson, len(mappings))

	for _, m := range mappings {
		mdPath := filepath.Join(lessonsDir, m.Source)
		body, err := os.ReadFile(mdPath)
		if err != nil {
			return nil, fmt.Errorf("read lesson %q: %w", m.Source, err)
		}
		title := extractTitle(string(body))
		lesson := &Lesson{
			Title:    title,
			Body:     string(body),
			Concepts: append(lessonConcepts(mappings, m.Source, m.ConceptID), m.ConceptID),
		}
		concepts[m.ConceptID] = lesson
		// Share the same lesson object for all concepts pointing to same source.
		for _, m2 := range mappings {
			if m2.Source == m.Source && m2.ConceptID != m.ConceptID {
				concepts[m2.ConceptID] = lesson
			}
		}
	}

	return &Loader{concepts: concepts}, nil
}

func (l *Loader) Lesson(conceptID string) *Lesson {
	return l.concepts[conceptID]
}

func (l *Loader) Count() int {
	return len(l.concepts)
}

func extractTitle(md string) string {
	lines := strings.SplitN(md, "\n", 3)
	if len(lines) > 0 && strings.HasPrefix(lines[0], "# ") {
		return strings.TrimPrefix(lines[0], "# ")
	}
	return ""
}

func lessonConcepts(mappings []mapping, source, skip string) []string {
	var out []string
	for _, m := range mappings {
		if m.Source == source && m.ConceptID != skip {
			out = append(out, m.ConceptID)
		}
	}
	return out
}
