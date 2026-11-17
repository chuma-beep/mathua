package geometry

import (
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

func TestPointsLinesServesSingleChoiceOptions(t *testing.T) {
	for i := 0; i < 50; i++ {
		p := pointsLinesProblem()
		if !strings.Contains(p.Question, "Which of these statements about points, lines and planes is true?") {
			t.Fatalf("stem does not ask for a single true statement: %q", p.Question)
		}
		options := generator.ChoiceOptions(p.Question)
		if len(options) != 4 {
			t.Fatalf("expected four options, got %d in %q", len(options), p.Question)
		}
		found := false
		for _, text := range options {
			if text == p.Answer {
				found = true
			}
		}
		if !found {
			t.Fatalf("answer %q is not one of the served options %v", p.Answer, options)
		}
		if !strings.Contains(p.Explanation, "true") || !strings.Contains(p.Explanation, "false") {
			t.Fatalf("explanation does not mark every statement: %q", p.Explanation)
		}
	}
}
