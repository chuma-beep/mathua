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
// Assembly is fail-soft by construction. A concept with no schema, or a
// placeholder with no fact behind it, resolves to the generator's own
// Problem.Explanation. A step whose fact is absent is simply not rendered. A
// learner is never shown a half-filled or partly-invented explanation, and the
// contract test fails the build when a schema and its generator disagree.
package solutions

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
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
// is the authored prose explaining it.
//
// Fact is load-bearing. A step is only rendered when its fact is present, so a
// schema can cover both questions a multi-shape generator asks. Mark the step
// optional when its absence is that design — "classify this angle" versus "find
// the complement" — and leave it required when a missing fact would mean the
// schema and its generator had drifted apart.
type Step struct {
	Fact string `json:"fact"`
	Say  string `json:"say"`
	// Optional marks a step that belongs to only one of the question shapes a
	// multi-shape generator asks. When its fact is absent the step is skipped
	// and nothing is reported, because its absence is the design rather than
	// drift. A step without this flag must have its fact, and a missing one
	// drops the whole schema — that is what catches a renamed fact.
	Optional bool `json:"optional,omitempty"`
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

// Render substitutes [[name]] placeholders in a schema and reports any
// placeholder it could not resolve, in order of first appearance.
//
// A step whose declared fact is absent is skipped rather than rendered: one
// generator often asks two different questions (find the complement of an
// angle, or classify it) and a single schema covers both, with each question's
// steps applying only to the one it came from. A step is never rendered
// without its fact, because that would be a claim we cannot support.
//
// An unresolved *placeholder* is different: it means the schema and the
// generator disagree about a name, and the whole schema is dropped. See
// Assemble.
//
// It is exported so the contract test can assert that every placeholder in
// every shipped schema resolves against a real generated problem, which is the
// only thing keeping an authored schema and its generator in agreement.
func Render(s *Schema, facts map[string]string) (string, []string) {
	return render(s, facts, false)
}

// Missing lists everything a schema needs that facts does not supply: the
// placeholder names it interpolates, plus the fact each step declares. Unlike
// Render it treats an absent step fact as a finding, so the contract test is
// told when a schema and its generator have drifted apart even though the
// runtime degrades gracefully.
func Missing(s *Schema, facts map[string]string) []string {
	_, missing := render(s, facts, true)
	return missing
}

func render(s *Schema, facts map[string]string, strictFacts bool) (string, []string) {
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
	var b strings.Builder
	if s.Setup != "" {
		b.WriteString(lookup(s.Setup, facts, note))
	}
	for _, step := range s.Steps {
		if step.Fact != "" {
			if _, ok := facts[step.Fact]; !ok {
				if strictFacts && !step.Optional {
					note(step.Fact)
				}
				continue
			}
		}
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

// placeholderRe matches a placeholder's contents: a snake_case fact name. The
// restriction is not cosmetic — the row-echelon form a systems question is
// graded against is literally written [[a,b,c],[0,d,e]], so an unrestricted
// [[...]] would chew up the answer it is supposed to be embedding. A typo like
// [[Num]] is therefore left as literal text and caught by the leak test rather
// than silently treated as a fact.
var placeholderRe = regexp.MustCompile(`^[a-z_][a-z0-9_]*$`)

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
		if !placeholderRe.MatchString(name) {
			// Not a placeholder (matrix notation, or a malformed name): pass the
			// whole span through untouched.
			b.WriteString(text[:i+j+2])
			text = text[i+j+2:]
			continue
		}
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

// Unresolved returns the placeholder-looking spans still present in text after
// assembly — that is, a test for a schema that referenced a fact its generator
// does not publish. Matrix notation such as [[4,3,6],[0,1,1]] is not a
// placeholder and is not reported.
func Unresolved(text string) []string {
	var out []string
	for i := 0; ; {
		j := strings.Index(text[i:], "[[")
		if j < 0 {
			return out
		}
		i += j
		k := strings.Index(text[i:], "]]")
		if k < 0 {
			return out
		}
		if name := text[i+2 : i+k]; placeholderRe.MatchString(name) {
			out = append(out, name)
		}
		i += k + 2
	}
}
