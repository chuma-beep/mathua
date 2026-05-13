package leaderboard

import (
	"testing"

	"github.com/chuma-beep/mathua/internal/levels"
)

func TestComputeLevel(t *testing.T) {
	cases := []struct {
		mastered int
		expected string
	}{
		{0, "Novice"},
		{31, "Novice"},
		{32, "Apprentice"},
		{64, "Student"},
		{96, "Scholar"},
		{128, "Adept"},
		{160, "Expert"},
		{192, "Master"},
		{224, "Grandmaster"},
		{256, "Math Architect"},
		{284, "Math Architect"},
	}
	for _, tc := range cases {
		got := levels.Compute(tc.mastered)
		if got != tc.expected {
			t.Errorf("mastered=%d: expected %q, got %q", tc.mastered, tc.expected, got)
		}
	}
}
