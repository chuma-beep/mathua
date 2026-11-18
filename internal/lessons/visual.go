package lessons

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"unicode/utf8"
)

// Interactive visuals attached to a knowledge-point shard.
//
// A visual is data, not code: the shard names a template from a closed
// allowlist and supplies its parameters. The Go server never interprets those
// parameters, it only refuses a payload that is oversized, malformed, or names a
// template that does not exist. The browser validates the parameters against the
// template's own schema before rendering, which keeps an authored mistake from
// taking down a lesson.
//
// The allowlist lives in data/visualizations/types.json rather than in this
// package so a template can be added to the frontend registry and the curriculum
// in one change without a second Go constant drifting beside it.

// MaxVisualParamsBytes bounds one template's parameters. It is generous for the
// six shipped templates, whose payloads are a few hundred bytes, and small
// enough that a shard cannot smuggle an essay or a data URL through a visual.
const MaxVisualParamsBytes = 8 * 1024

const (
	maxVisualTitleRunes   = 160
	maxVisualCaptionRunes = 600

	// A prediction is one question with two to four answers. The limits are small
	// on purpose: this is a prompt to commit to an expectation before touching the
	// figure, not a second explanation of the concept.
	minPredictOptions       = 2
	maxPredictOptions       = 4
	maxPredictQuestionRunes = 300
	maxPredictOptionRunes   = 160
	maxPredictExplainRunes  = 600
)

// Predict is an authored question shown before the figure is used.
//
// The proposal's fifth pedagogical step is "ask learners to predict the outcome
// of an interaction", and the figures shipped without it: every one responded to
// a control without ever asking the learner what they expected. The block lives
// on the visual rather than on each template so a template never has to know it
// is being used inside a lesson, and so a figure added later gets the behaviour
// for free.
type Predict struct {
	Question string   `json:"question"`
	Options  []string `json:"options"`
	// Answer indexes Options. The browser compares the learner's choice to it and
	// never sends the choice anywhere: a prediction is not evidence and must not
	// reach mastery or scheduling.
	Answer  int    `json:"answer"`
	Explain string `json:"explain,omitempty"`
}

// Visual is the declarative specification on a KP. Params stays raw JSON: the
// server has no reason to know what a matrix or a normal distribution looks
// like, and a schema here would have to be duplicated in TypeScript anyway.
type Visual struct {
	Type    string          `json:"type"`
	Title   string          `json:"title,omitempty"`
	Caption string          `json:"caption,omitempty"`
	Params  json.RawMessage `json:"params,omitempty"`
	// Predict is optional. A visual without one is unchanged, which is what the
	// eleven shipped visuals were before this field existed.
	Predict *Predict `json:"predict,omitempty"`
}

type visualTypeEntry struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

// LoadVisualTypes reads the allowlist. A missing file is an error, because a
// silent empty allowlist would drop every visual on a deployment with no way to
// tell an authoring mistake from a broken path.
func LoadVisualTypes(path string) (map[string]bool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read visual types: %w", err)
	}
	var entries []visualTypeEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, fmt.Errorf("parse visual types: %w", err)
	}
	out := make(map[string]bool, len(entries))
	for i, e := range entries {
		id := strings.TrimSpace(e.ID)
		if id == "" {
			return nil, fmt.Errorf("visual types entry %d has no id", i)
		}
		switch e.Status {
		case "stable", "experimental":
		default:
			return nil, fmt.Errorf("visual type %q has unknown status %q", id, e.Status)
		}
		out[id] = true
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("visual types allowlist is empty")
	}
	return out, nil
}

// Validate checks one visual against the allowlist and the payload budget.
func (v *Visual) Validate(types map[string]bool) error {
	if v == nil {
		return fmt.Errorf("visual is nil")
	}
	if strings.TrimSpace(v.Type) == "" {
		return fmt.Errorf("visual type is empty")
	}
	if !types[v.Type] {
		return fmt.Errorf("unknown visual type %q", v.Type)
	}
	if n := utf8.RuneCountInString(v.Title); n > maxVisualTitleRunes {
		return fmt.Errorf("visual title is %d runes, limit %d", n, maxVisualTitleRunes)
	}
	if n := utf8.RuneCountInString(v.Caption); n > maxVisualCaptionRunes {
		return fmt.Errorf("visual caption is %d runes, limit %d", n, maxVisualCaptionRunes)
	}
	// The prediction is validated before the early return for an empty parameter
	// block, because a visual may legitimately rely on its template defaults and
	// still ask a question.
	if err := v.Predict.validate(); err != nil {
		return err
	}
	params := bytes.TrimSpace(v.Params)
	if len(params) == 0 {
		return nil
	}
	if len(params) > MaxVisualParamsBytes {
		return fmt.Errorf("visual params are %d bytes, limit %d", len(params), MaxVisualParamsBytes)
	}
	if params[0] != '{' {
		return fmt.Errorf("visual params must be a JSON object")
	}
	if !json.Valid(params) {
		return fmt.Errorf("visual params are not valid JSON")
	}
	return nil
}

// validate checks one prediction. A nil prediction is valid: the block is
// optional and most visuals may never need one.
func (p *Predict) validate() error {
	if p == nil {
		return nil
	}
	question := strings.TrimSpace(p.Question)
	if question == "" {
		return fmt.Errorf("visual prediction has no question")
	}
	if n := utf8.RuneCountInString(question); n > maxPredictQuestionRunes {
		return fmt.Errorf("visual prediction question is %d runes, limit %d", n, maxPredictQuestionRunes)
	}
	if len(p.Options) < minPredictOptions || len(p.Options) > maxPredictOptions {
		return fmt.Errorf("visual prediction has %d options, want %d to %d", len(p.Options), minPredictOptions, maxPredictOptions)
	}
	for i, option := range p.Options {
		trimmed := strings.TrimSpace(option)
		if trimmed == "" {
			return fmt.Errorf("visual prediction option %d is empty", i+1)
		}
		if n := utf8.RuneCountInString(trimmed); n > maxPredictOptionRunes {
			return fmt.Errorf("visual prediction option %d is %d runes, limit %d", i+1, n, maxPredictOptionRunes)
		}
	}
	if p.Answer < 0 || p.Answer >= len(p.Options) {
		return fmt.Errorf("visual prediction answer %d is outside 1 to %d", p.Answer, len(p.Options))
	}
	if n := utf8.RuneCountInString(p.Explain); n > maxPredictExplainRunes {
		return fmt.Errorf("visual prediction explanation is %d runes, limit %d", n, maxPredictExplainRunes)
	}
	return nil
}
