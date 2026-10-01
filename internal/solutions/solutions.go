// Package solutions assembles a problem's explanation from a corpus-authored
// schema and the facts its generator published.
//
// The split is the whole point of this package, and it is why a schema never
// re-implements the arithmetic:
//
//   - The numbers stay in the generator, which remains the single source of
//     truth for the answer. A schema names the values the generator already had
//     in scope; it cannot recompute them, so it cannot disagree with them.
//   - The prose lives in data/lessons/solutions/<concept>.json, where a
//     content author can read and edit it as prose.
//
// A placeholder is [[name]]. Single braces are LaTeX and the corpus is full of
// them, so they cannot be the delimiter; [[name]] cannot appear in LaTeX or
// ordinary prose, which makes the substitution unambiguous and leaves the
// authored text readable as English with a few obvious slots.
//
// Assembly is fail-soft by construction. A concept with no schema, a step
// whose fact the generator no longer publishes, or a placeholder with no fact
// behind it all resolve to the generator's own Problem.Explanation. A learner
// is never shown a half-filled or partly-invented explanation.
package solutions

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Schema is the corpus-authored explanation for one concept. It is a template
// over the instance's facts, not a worked example of a specific problem: the
// same schema renders every instance of the concept, so the learner sees the
// same reasoning shape with different numbers each time.
type Schema struct {
	Concept string `json:"concept"`
	// Setup frames the problem in words, with the instance's values inlined.
	Setup string `json:"setup"`
	// Steps are the derivation, in order. Each rests on one named fact.
	Steps []Step `json:"steps"`
	// Answer closes on the final result. Required: an explanation that does not
	// end on the answer leaves the learner to infer it.
	Answer string `json:"answer"`
}

// Step is one derivation step. Fact names the fact the step rests on, and Say
// is the authored prose explaining it. Fact is load-bearing, not decoration:
// a step whose fact is missing is a stale schema, and Assemble refuses the
// whole schema rather than render a step that has lost its footing.
type Step struct {
	Fact string `json:"fact"`
	Say  string `json:"say"`
}

// Loader holds the schemas for every concept that has one. A nil *Loader is
// valid and assembles nothing, which is what the engine gets before any
// content is authored.
type Loader struct {
	byConcept map[string]*Schema
}

// Load reads every <concept>.json in dir. A missing or unreadable directory
// yields an empty Loader rather than an error: solutions are an enhancement to
// the generator's own explanation, never a dependency of serving a question.
func Load(dir string) (*Loader, error) {
	l := &Loader{byConcept: map[string]*Schema{}}
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return l, nil
		}
		return l, fmt.Errorf("read solutions dir: %w", err)
	}
	var problems []string
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		path := filepath.Join(dir, e.Name())
		raw, err := os.ReadFile(path)
		if err != nil {
			problems = append(problems, fmt.Sprintf("%s: %v", e.Name(), err))
			continue
		}
		var s Schema
		if err := json.Unmarshal(raw, &s); err != nil {
			problems = append(problems, fmt.Sprintf("%s: %v", e.Name(), err))
			continue
		}
		if s.Answer == "" {
			problems = append(problems, fmt.Sprintf("%s: no answer step", e.Name()))
			continue
		}
		id := s.Concept
		if id == "" {
			id = strings.TrimSuffix(e.Name(), ".json")
		}
		if _, dup := l.byConcept[id]; dup {
			problems = append(problems, fmt.Sprintf("%s: duplicate schema for %q", e.Name(), id))
			continue
		}
		l.byConcept[id] = &s
	}
	if len(problems) > 0 {
		sort.Strings(problems)
		return l, fmt.Errorf("unusable solution schemas: %s", strings.Join(problems, "; "))
	}
	return l, nil
}

// Concepts returns the ids that have a schema, sorted.
func (l *Loader) Concepts() []string {
	if l == nil {
		return nil
	}
	ids := make([]string, 0, len(l.byConcept))
	for id := range l.byConcept {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}

// Schema returns the schema for a concept, if it has one.
func (l *Loader) Schema(conceptID string) (*Schema, bool) {
	if l == nil {
		return nil, false
	}
	s, ok := l.byConcept[conceptID]
	return s, ok
}

// Assemble renders conceptID's schema against facts, or returns fallback when
// it cannot do so faithfully. See Render for the rules and Missing for why.
func (l *Loader) Assemble(conceptID string, facts map[string]string, fallback string) string {
	if l == nil {
		return fallback
	}
	s, ok := l.byConcept[conceptID]
	if !ok {
		return fallback
	}
	rendered, missing := Render(s, facts)
	if len(missing) > 0 {
		// A schema that no longer matches its generator is a content bug, not
		// a reason to teach something wrong. Fall back and say so once.
		log.Printf("solutions: %q schema skipped, unresolved %v", conceptID, missing)
		return fallback
	}
	return rendered
}

// Render substitutes [[name]] placeholders in a schema and reports any name it
// could not resolve, in order of first appearance.
//
// It is exported so the contract test can assert that every placeholder in
// every shipped schema resolves against a real generated problem, which is the
// only thing keeping an authored schema and its generator in agreement.
func Render(s *Schema, facts map[string]string) (string, []string) {
	if s == nil {
		return "", nil
	}
	var missing []string
	seen := map[string]bool{}
	note := func(name string) {
		if !seen[name] {
			seen[name] = true
			missing = append(missing, name)
		}
	}
	// A step's declared fact must still exist, or the step has lost its
	// footing and rendering it would be a claim we cannot support.
	for _, step := range s.Steps {
		if step.Fact == "" {
			continue
		}
		if _, ok := facts[step.Fact]; !ok {
			note(step.Fact)
		}
	}
	var b strings.Builder
	if s.Setup != "" {
		b.WriteString(lookup(s.Setup, facts, note))
	}
	for _, step := range s.Steps {
		say := lookup(step.Say, facts, note)
		if say == "" {
			continue
		}
		if b.Len() > 0 {
			b.WriteString("\n\n")
		}
		b.WriteString(say)
	}
	answer := lookup(s.Answer, facts, note)
	if answer == "" {
		return "", append(missing, "<answer>")
	}
	if b.Len() > 0 {
		b.WriteString("\n\n")
	}
	b.WriteString(answer)
	if len(missing) > 0 {
		return "", missing
	}
	return b.String(), nil
}

// Missing lists the placeholder names a schema needs that facts does not
// supply. The contract test uses it to report exactly which concepts still
// need their generator to publish a fact.
func Missing(s *Schema, facts map[string]string) []string {
	_, missing := Render(s, facts)
	return missing
}

// lookup substitutes every [[name]] in text. An unresolved name is left
// verbatim so a partially-resolved render is visibly wrong in tests and logs
// rather than silently plausible.
func lookup(text string, facts map[string]string, note func(string)) string {
	if !strings.Contains(text, "[[") {
		return text
	}
	var b strings.Builder
	for {
		i := strings.Index(text, "[[")
		if i < 0 {
			b.WriteString(text)
			return b.String()
		}
		j := strings.Index(text[i:], "]]")
		if j < 0 {
			b.WriteString(text)
			return b.String()
		}
		name := text[i+2 : i+j]
		b.WriteString(text[:i])
		if v, ok := facts[name]; ok {
			b.WriteString(v)
		} else {
			note(name)
			b.WriteString(text[i : i+j+2])
		}
		text = text[i+j+2:]
	}
}
