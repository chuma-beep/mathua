package xp

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// The landing page advertised a 150 XP mastery check for as long as the gate was 50. Nothing
// broke: the page is static, the number is a string literal, and ADR-020's rescale touched Go
// and the daily-goal backfill but not the marketing copy. It is the one number on that page
// that build-graph.mjs cannot derive from the corpus, so it is the one that can drift, and it
// did.
//
// The check is on the generated artifact rather than on the page, because the page reads it:
// app/home/data.ts exports graphMeta.quizGateXP and app/home/sections.tsx renders that. Editing
// the JSON by hand fails here; changing the engine constant without rebuilding the artifact
// fails here too.
func TestLandingQuizGateMatchesEngine(t *testing.T) {
	path := filepath.Join("..", "..", "web", "next-app", "data", "graph.meta.json")
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v (run `npm run graph:build`)", path, err)
	}
	var meta struct {
		QuizGateXP *int `json:"quizGateXP"`
	}
	if err := json.Unmarshal(raw, &meta); err != nil {
		t.Fatalf("parse %s: %v", path, err)
	}
	if meta.QuizGateXP == nil {
		t.Fatalf("%s has no quizGateXP; the landing page falls back to a hand-written number", path)
	}
	if *meta.QuizGateXP != QuizGateXP {
		t.Errorf("landing page advertises a %d XP mastery check, the engine gates it at %d XP\n"+
			"fix: set QUIZ_GATE_XP in web/next-app/lib/graphPayload.ts to %d and run `npm run graph:build`",
			*meta.QuizGateXP, QuizGateXP, QuizGateXP)
	}
}
