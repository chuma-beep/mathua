package engine

import (
	"testing"
	"time"

	"github.com/chuma-beep/mathua/internal/storage"
)

func at(y, m, d int) time.Time {
	return time.Date(y, time.Month(m), d, 12, 0, 0, 0, time.UTC)
}

func TestWeekStartUTC_Monday(t *testing.T) {
	cases := []struct {
		in   time.Time
		want string
	}{
		{at(2026, 9, 7), "2026-09-07"},  // Monday stays
		{at(2026, 9, 10), "2026-09-07"}, // Thursday -> Monday
		{at(2026, 9, 13), "2026-09-07"}, // Sunday -> previous Monday
		{at(2026, 9, 14), "2026-09-14"}, // next Monday
	}
	for _, c := range cases {
		if got := weekStartUTC(c.in).Format("2006-01-02"); got != c.want {
			t.Errorf("weekStartUTC(%s) = %s, want %s", c.in.Format("2006-01-02"), got, c.want)
		}
	}
}

func TestComputeEfficacyTrend_WeeksAndRetention(t *testing.T) {
	attempts := []storage.AttemptEntry{
		// Week 1: student s1 touches c1 (correct first try), s2 touches c2 (wrong then right)
		{StudentID: "s1", ConceptID: "c1", Correct: true, Timestamp: at(2026, 9, 8)},
		{StudentID: "s2", ConceptID: "c2", Correct: false, Timestamp: at(2026, 9, 8)},
		{StudentID: "s2", ConceptID: "c2", Correct: true, Timestamp: at(2026, 9, 8)},
		// Week 2: s1 returns (correct), s3 appears (correct)
		{StudentID: "s1", ConceptID: "c3", Correct: true, Timestamp: at(2026, 9, 15)},
		{StudentID: "s3", ConceptID: "c4", Correct: true, Timestamp: at(2026, 9, 15)},
	}
	trend := computeEfficacyTrend(attempts)
	if len(trend.Weeks) != 2 {
		t.Fatalf("expected 2 weeks, got %d", len(trend.Weeks))
	}
	if trend.Weeks[0].WeekStart != "2026-09-07" || trend.Weeks[1].WeekStart != "2026-09-14" {
		t.Errorf("weeks not Monday-sorted: %+v", trend.Weeks)
	}
	if trend.Weeks[0].ActiveStudents != 2 {
		t.Errorf("week1 active students = %d, want 2", trend.Weeks[0].ActiveStudents)
	}
	if trend.Weeks[1].ActiveStudents != 2 {
		t.Errorf("week2 active students = %d, want 2", trend.Weeks[1].ActiveStudents)
	}
	if trend.TotalStudents != 3 {
		t.Errorf("total students = %d, want 3", trend.TotalStudents)
	}
	// s1 active both weeks; s2 and s3 only one week each.
	if trend.ReturningStudents != 1 {
		t.Errorf("returning students = %d, want 1", trend.ReturningStudents)
	}
	if got := trend.RetentionRate; got < 0.33 || got > 0.34 {
		t.Errorf("retention rate = %f, want ~0.333", got)
	}
	// Week 1 first-pass: c1 correct, c2 wrong-then-right -> 1/2.
	if trend.Weeks[0].FirstPassRate != 0.5 {
		t.Errorf("week1 first-pass = %f, want 0.5", trend.Weeks[0].FirstPassRate)
	}
	// Week 2 first-pass: both correct -> 1.0, so the trend is +0.5.
	if trend.Weeks[1].FirstPassRate != 1.0 {
		t.Errorf("week2 first-pass = %f, want 1.0", trend.Weeks[1].FirstPassRate)
	}
	if trend.FirstPassTrend != 0.5 {
		t.Errorf("first-pass trend = %f, want 0.5", trend.FirstPassTrend)
	}
}

func TestComputeEfficacyTrend_Empty(t *testing.T) {
	trend := computeEfficacyTrend(nil)
	if len(trend.Weeks) != 0 || trend.TotalStudents != 0 || trend.RetentionRate != 0 {
		t.Errorf("expected empty trend, got %+v", trend)
	}
}
