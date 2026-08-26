package lessons

import (
	"encoding/json"
	"fmt"
	"github.com/chuma-beep/mathua/internal/latex"
	"log"
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

	// Group concept IDs by source file
	sourceConcepts := make(map[string][]string)
	for _, m := range mappings {
		sourceConcepts[m.Source] = append(sourceConcepts[m.Source], m.ConceptID)
	}

	concepts := make(map[string]*Lesson, len(mappings))

	for source, ids := range sourceConcepts {
		mdPath := filepath.Join(lessonsDir, source)
		body, err := os.ReadFile(mdPath)
		if err != nil {
			log.Printf("warning: lesson file not found, skipping %q: %v", source, err)
			continue
		}
		content := latex.Canonicalize(string(body), latex.ForSource(source, string(body)))
		for _, w := range latex.Validate(content, latex.ForSource(source, string(body))) {
			log.Printf("latex warning in %q: %s", source, w)
		}
		title := extractTitle(content)
		if title == "" {
			title = titleFromFilename(source)
		}
		lesson := &Lesson{
			Title:    title,
			Body:     content,
			Concepts: ids,
		}
		for _, id := range ids {
			concepts[id] = lesson
		}
	}

	return &Loader{concepts: concepts}, nil
}

func (l *Loader) Lesson(conceptID string) *Lesson {
	return l.concepts[conceptID]
}

// LessonByTitle returns the first lesson whose Title matches exactly.
// Titles are unique per lessons.json registration in practice; if
// duplicates exist, the domain-grouped order determines the winner.
func (l *Loader) LessonByTitle(title string) *Lesson {
	for _, lesson := range l.concepts {
		if lesson.Title == title {
			return lesson
		}
	}
	return nil
}

func (l *Loader) All() map[string]*Lesson {
	return l.concepts
}

func (l *Loader) LessonsByDomain() map[string][]*Lesson {
	byDomain := make(map[string][]*Lesson)
	seen := make(map[*Lesson]bool)
	for _, lesson := range l.concepts {
		if seen[lesson] {
			continue
		}
		seen[lesson] = true
		domain := "general"
		if len(lesson.Concepts) > 0 {
			parts := strings.SplitN(lesson.Concepts[0], ".", 3)
			if len(parts) >= 2 {
				domain = parts[0] + "." + parts[1]
			}
		}
		byDomain[domain] = append(byDomain[domain], lesson)
	}
	return byDomain
}

func (l *Loader) Count() int {
	return len(l.concepts)
}

func extractTitle(md string) string {
	var first string
	for _, line := range strings.Split(md, "\n") {
		trimmed := strings.TrimLeft(line, "#")
		if len(trimmed) < len(line) && strings.HasPrefix(trimmed, " ") {
			title := strings.TrimSpace(trimmed)
			if first == "" {
				first = title
			}
			if !isGenericTitle(title) {
				return title
			}
		}
	}
	if first != "" {
		return first
	}
	return ""
}

// Generic opening headings ("Introduction") describe position, not content;
// prefer the first substantive heading for the lesson title.
func isGenericTitle(t string) bool {
	switch strings.ToLower(strings.Trim(t, " :")) {
	case "introduction", "intro", "overview", "contents", "summary", "prerequisites":
		return true
	}
	return false
}

func titleFromFilename(path string) string {
	base := filepath.Base(path)
	base = strings.TrimSuffix(base, ".md")
	name := strings.ReplaceAll(base, "-", " ")
	if len(name) > 0 {
		name = strings.ToUpper(name[:1]) + name[1:]
	}
	return name
}
