package scoring

import (
	"testing"
	"time"

	"github.com/chuma-beep/mathua/internal/concepts"
	"github.com/chuma-beep/mathua/internal/storage"
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
		got := computeLevel(tc.mastered)
		if got != tc.expected {
			t.Errorf("mastered=%d: expected %q, got %q", tc.mastered, tc.expected, got)
		}
	}
}

func TestComputeSpeedBonus(t *testing.T) {
	p := &storage.ConceptProgress{AvgResponseTime: 5.0}
	bonus := computeSpeedBonus(p)
	if bonus != 10 {
		t.Errorf("expected 10, got %f", bonus)
	}
	p.AvgResponseTime = 0
	if computeSpeedBonus(p) != 0 {
		t.Error("expected 0 for zero time")
	}
	p.AvgResponseTime = 0.5
	bonus = computeSpeedBonus(p)
	if bonus != 25 {
		t.Errorf("expected capped at 25, got %f", bonus)
	}
}

func TestComputeCurrentStreak_Empty(t *testing.T) {
	s := computeCurrentStreak(map[string]*storage.ConceptProgress{})
	if s != 0 {
		t.Errorf("expected 0, got %d", s)
	}
}

func TestComputeCurrentStreak_Today(t *testing.T) {
	now := time.Now().UTC()
	s := computeCurrentStreak(map[string]*storage.ConceptProgress{
		"a": {LastAttempted: &now},
	})
	if s != 1 {
		t.Errorf("expected 1, got %d", s)
	}
}

func TestComputeCurrentStreak_Yesterday(t *testing.T) {
	yesterday := time.Now().UTC().AddDate(0, 0, -1)
	s := computeCurrentStreak(map[string]*storage.ConceptProgress{
		"a": {LastAttempted: &yesterday},
	})
	// No activity today, so streak is 0.
	if s != 0 {
		t.Errorf("expected 0 (no today), got %d", s)
	}
}

func TestCompute(t *testing.T) {
	d, _ := concepts.Load("../../data/concepts.json")
	store, _ := storage.NewSQLiteStore(":memory:")
	defer store.Close()
	st, _ := store.CreateStudent("tester")
	now := time.Now().UTC()
	_ = store.UpsertProgress(&storage.ConceptProgress{
		StudentID: st.ID, ConceptID: "count.objects",
		Status: "MASTERED", AvgResponseTime: 5.0, MasteredAt: &now, LastAttempted: &now,
	})
	_ = store.UpsertProgress(&storage.ConceptProgress{
		StudentID: st.ID, ConceptID: "count.cardinality",
		Status: "LEARNING", AvgResponseTime: 8.0, LastAttempted: &now,
	})

	u := NewUpdater(d, store)
	scores, err := u.Compute(st.ID)
	if err != nil {
		t.Fatalf("compute: %v", err)
	}
	if scores.ConceptsMastered != 1 {
		t.Errorf("expected 1 mastered, got %d", scores.ConceptsMastered)
	}
	if scores.Level != "Novice" {
		t.Errorf("expected Novice, got %s", scores.Level)
	}
	if scores.WeeklyScore == 0 {
		t.Errorf("expected non-zero weekly score, got %d", scores.WeeklyScore)
	}
}

func TestWeekStart(t *testing.T) {
	// 2026-05-12 is a Tuesday. Monday is May 11.
	ts := time.Date(2026, 5, 12, 15, 30, 0, 0, time.UTC)
	m := weekStart(ts)
	if m.Year() != 2026 || m.Month() != 5 || m.Day() != 11 {
		t.Errorf("expected 2026-05-11, got %s", m.Format("2006-01-02"))
	}
	if m.Hour() != 0 || m.Minute() != 0 {
		t.Error("expected midnight")
	}
}
