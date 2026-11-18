package lessons

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func visualAllowlist(t *testing.T) map[string]bool {
	t.Helper()
	types, err := LoadVisualTypes("../../data/visualizations/types.json")
	if err != nil {
		t.Fatalf("load allowlist: %v", err)
	}
	return types
}

func TestLoadVisualTypes(t *testing.T) {
	types := visualAllowlist(t)
	for _, id := range []string{
		"linalg.matrix2d",
		"calc.tangent",
		"calc.riemann",
		"complex.multiply",
		"stat.normal",
		"geo.transform2d",
	} {
		if !types[id] {
			t.Errorf("allowlist is missing %q", id)
		}
	}
}

func TestLoadVisualTypesRejectsBadFiles(t *testing.T) {
	dir := t.TempDir()
	cases := map[string]string{
		"empty.json":      `[]`,
		"no-id.json":      `[{"id":"","status":"stable"}]`,
		"bad-status.json": `[{"id":"x","status":"maybe"}]`,
		"not-json.json":   `{`,
	}
	for name, body := range cases {
		path := filepath.Join(dir, name)
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
		if _, err := LoadVisualTypes(path); err == nil {
			t.Errorf("%s: expected an error", name)
		}
	}
}

func TestVisualValidate(t *testing.T) {
	types := map[string]bool{"calc.tangent": true}
	ok := &Visual{
		Type:   "calc.tangent",
		Title:  "Tangent",
		Params: json.RawMessage(`{"coeffs":[0,0,1],"x":1}`),
	}
	if err := ok.Validate(types); err != nil {
		t.Fatalf("valid visual rejected: %v", err)
	}
	if err := (&Visual{Type: "calc.tangent"}).Validate(types); err != nil {
		t.Fatalf("a visual with no params should use the template defaults: %v", err)
	}

	oversized := `{"blob":"` + strings.Repeat("x", MaxVisualParamsBytes) + `"}`
	cases := []struct {
		name string
		v    *Visual
	}{
		{"unknown type", &Visual{Type: "calc.nope"}},
		{"empty type", &Visual{Type: "  "}},
		{"params not an object", &Visual{Type: "calc.tangent", Params: json.RawMessage(`[1,2]`)}},
		{"params invalid", &Visual{Type: "calc.tangent", Params: json.RawMessage(`{`)}},
		{"params oversized", &Visual{Type: "calc.tangent", Params: json.RawMessage(oversized)}},
		{"title too long", &Visual{Type: "calc.tangent", Title: strings.Repeat("t", maxVisualTitleRunes+1)}},
		{"caption too long", &Visual{Type: "calc.tangent", Caption: strings.Repeat("c", maxVisualCaptionRunes+1)}},
	}
	for _, tc := range cases {
		if err := tc.v.Validate(types); err == nil {
			t.Errorf("%s: expected an error", tc.name)
		}
	}
}

func TestLoadKPsDropsInvalidVisual(t *testing.T) {
	dir := t.TempDir()
	body := `[
	  {"label":"Good","section":"Good","subgoals":[],"visual":{"type":"calc.tangent","params":{"x":1}}},
	  {"label":"Bad","section":"Bad","subgoals":[],"visual":{"type":"made.up","params":{}}},
	  {"label":"Oversized","section":"Big","subgoals":[],"visual":{"type":"calc.tangent","params":{"blob":"` +
		strings.Repeat("x", MaxVisualParamsBytes) + `"}}}
	]`
	kpDir := filepath.Join(dir, "kp")
	if err := os.MkdirAll(kpDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(kpDir, "concept.json"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	kps := loadKPs(dir, map[string]bool{"calc.tangent": true})
	if len(kps["concept"]) != 3 {
		t.Fatalf("expected all three steps to survive, got %d", len(kps["concept"]))
	}
	if kps["concept"][0].Visual == nil {
		t.Error("valid visual was dropped")
	}
	if kps["concept"][1].Visual != nil {
		t.Error("unknown type was kept")
	}
	if kps["concept"][2].Visual != nil {
		t.Error("oversized params were kept")
	}
}

// The curriculum is the real test: every shard that names a visual must pass the
// same validation a deployment applies, so a typo cannot reach production.
func TestCurriculumVisualsAreValid(t *testing.T) {
	types := visualAllowlist(t)
	files, err := filepath.Glob("../../data/lessons/kp/*.json")
	if err != nil {
		t.Fatal(err)
	}
	count := 0
	for _, path := range files {
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read %s: %v", path, err)
		}
		var kps []KP
		if err := json.Unmarshal(data, &kps); err != nil {
			t.Fatalf("parse %s: %v", path, err)
		}
		for i, kp := range kps {
			if kp.Visual == nil {
				continue
			}
			count++
			if err := kp.Visual.Validate(types); err != nil {
				t.Errorf("%s step %d: %v", filepath.Base(path), i+1, err)
			}
		}
	}
	if count < 6 {
		t.Fatalf("expected at least the six authored visuals, found %d", count)
	}
}

func TestPredictValidate(t *testing.T) {
	types := map[string]bool{"calc.tangent": true}
	good := &Visual{
		Type:    "calc.tangent",
		Predict: &Predict{Question: "Which is steeper?", Options: []string{"the tangent", "the secant"}, Answer: 0, Explain: "At h > 0 the secant is below."},
	}
	if err := good.Validate(types); err != nil {
		t.Fatalf("valid prediction rejected: %v", err)
	}
	// A prediction must validate even when the figure relies on template defaults,
	// because the params block is allowed to be empty.
	noParams := &Visual{Type: "calc.tangent", Predict: &Predict{Question: "Q?", Options: []string{"a", "b"}, Answer: 1}}
	if err := noParams.Validate(types); err != nil {
		t.Fatalf("prediction without params rejected: %v", err)
	}

	cases := []struct {
		name string
		p    *Predict
	}{
		{"empty question", &Predict{Question: "  ", Options: []string{"a", "b"}}},
		{"one option", &Predict{Question: "Q?", Options: []string{"a"}}},
		{"five options", &Predict{Question: "Q?", Options: []string{"a", "b", "c", "d", "e"}}},
		{"empty option", &Predict{Question: "Q?", Options: []string{"a", "  "}}},
		{"answer negative", &Predict{Question: "Q?", Options: []string{"a", "b"}, Answer: -1}},
		{"answer past end", &Predict{Question: "Q?", Options: []string{"a", "b"}, Answer: 2}},
		{"question too long", &Predict{Question: strings.Repeat("q", maxPredictQuestionRunes+1), Options: []string{"a", "b"}}},
		{"option too long", &Predict{Question: "Q?", Options: []string{strings.Repeat("o", maxPredictOptionRunes+1), "b"}}},
		{"explain too long", &Predict{Question: "Q?", Options: []string{"a", "b"}, Explain: strings.Repeat("e", maxPredictExplainRunes+1)}},
	}
	for _, tc := range cases {
		v := &Visual{Type: "calc.tangent", Predict: tc.p}
		if err := v.Validate(types); err == nil {
			t.Errorf("%s: expected an error", tc.name)
		}
	}
}
