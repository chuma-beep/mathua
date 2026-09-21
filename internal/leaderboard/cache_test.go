package leaderboard

import (
	"testing"
	"time"

	"github.com/chuma-beep/mathua/internal/storage"
)

type countingRepo struct {
	storage.Repository
	calls int
}

func (r *countingRepo) GetWeeklyLeaderboard() ([]storage.LeaderboardRow, error) {
	r.calls++
	return r.Repository.GetWeeklyLeaderboard()
}

func TestWeekly_CachesWithinTTL(t *testing.T) {
	store, err := storage.NewSQLiteStore(":memory:")
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	if _, err := store.CreateStudent("loadtest"); err != nil {
		t.Fatalf("seed student: %v", err)
	}
	repo := &countingRepo{Repository: store}

	c := NewComputerWithTTL(repo, 5*time.Minute)
	first, err := c.Weekly()
	if err != nil {
		t.Fatalf("first Weekly: %v", err)
	}
	second, err := c.Weekly()
	if err != nil {
		t.Fatalf("second Weekly: %v", err)
	}
	if repo.calls != 1 {
		t.Fatalf("expected 1 store hit for 2 calls within TTL, got %d", repo.calls)
	}
	if len(first) != 1 || len(second) != 1 {
		t.Fatalf("expected 1 entry per call, got %d/%d", len(first), len(second))
	}
	if first[0] != second[0] {
		t.Fatalf("cached entries differ: %+v vs %+v", first[0], second[0])
	}
}

func TestWeekly_RefetchesAfterTTL(t *testing.T) {
	store, err := storage.NewSQLiteStore(":memory:")
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	if _, err := store.CreateStudent("loadtest"); err != nil {
		t.Fatalf("seed student: %v", err)
	}
	repo := &countingRepo{Repository: store}

	c := NewComputerWithTTL(repo, 20*time.Millisecond)
	if _, err := c.Weekly(); err != nil {
		t.Fatalf("first Weekly: %v", err)
	}
	time.Sleep(50 * time.Millisecond)
	if _, err := c.Weekly(); err != nil {
		t.Fatalf("second Weekly: %v", err)
	}
	if repo.calls != 2 {
		t.Fatalf("expected 2 store hits across TTL expiry, got %d", repo.calls)
	}
}
