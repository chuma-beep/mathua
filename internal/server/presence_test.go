package server

import (
	"testing"
	"time"
)

func TestPresenceTrackerCountsRecentLearners(t *testing.T) {
	base := time.Date(2026, 11, 14, 10, 0, 0, 0, time.UTC)
	now := base
	p := newPresenceTracker(func() time.Time { return now })

	p.Touch("a")
	p.Touch("b")
	p.Touch("")
	if got := p.Count(now); got != 2 {
		t.Fatalf("Count = %d, want 2", got)
	}

	now = base.Add(onlineWindow - time.Second)
	if got := p.Count(now); got != 2 {
		t.Fatalf("Count inside window = %d, want 2", got)
	}

	now = base.Add(onlineWindow)
	if got := p.Count(now); got != 0 {
		t.Fatalf("Count at window edge = %d, want 0", got)
	}
}

func TestPresenceTrackerKeepsAFreshLearnerWhenAnotherExpires(t *testing.T) {
	base := time.Date(2026, 11, 14, 10, 0, 0, 0, time.UTC)
	now := base
	p := newPresenceTracker(func() time.Time { return now })

	p.Touch("old")
	now = base.Add(4 * time.Minute)
	p.Touch("new")
	now = base.Add(6 * time.Minute)

	if got := p.Count(now); got != 1 {
		t.Fatalf("Count = %d, want 1", got)
	}
}

func TestPresenceTrackerNilIsSafe(t *testing.T) {
	var p *presenceTracker
	p.Touch("a")
	if got := p.Count(time.Now()); got != 0 {
		t.Fatalf("Count on nil tracker = %d, want 0", got)
	}
}
