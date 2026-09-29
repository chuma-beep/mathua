package engine

import (
	"testing"

	"github.com/chuma-beep/mathua/internal/storage"
)

func TestMasteryPctBands(t *testing.T) {
	mk := func(status string, streak int) *storage.ConceptProgress {
		return &storage.ConceptProgress{Status: status, Streak: streak}
	}
	cases := []struct {
		name     string
		p        *storage.ConceptProgress
		req      int
		wantBand string
	}{
		{"nil is unseen", nil, 10, "unseen"},
		{"unseen", mk("UNSEEN", 0), 10, "unseen"},
		{"learning partial", mk("LEARNING", 2), 10, "barely understood"},
		{"practicing", mk("PRACTICING", 6), 10, "developing"},
		{"practicing strong", mk("PRACTICING", 9), 10, "strong"},
		{"mastered", mk("MASTERED", 10), 10, "well retained"},
	}
	for _, c := range cases {
		pct := MasteryPct(c.p, c.req, 0, 30)
		if got := MasteryBand(pct); got != c.wantBand {
			t.Errorf("%s: MasteryPct=%v band=%q want %q", c.name, pct, got, c.wantBand)
		}
	}
}

func TestMasteryPctDecaying(t *testing.T) {
	p := &storage.ConceptProgress{Status: "MASTERED", Streak: 10}
	fresh := MasteryPct(p, 10, 1, 30)
	stale := MasteryPct(p, 10, 60, 30)
	if stale >= fresh {
		t.Errorf("decayed mastery %v should be below fresh %v", stale, fresh)
	}
}
