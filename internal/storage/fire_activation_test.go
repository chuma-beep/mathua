package storage

import (
	"sync"
	"testing"
	"time"
)

// TestFireActivationRevisionAdvancesOncePerCommit catches the lost-update
// shape the original SELECT-then-UPDATE had: two operators read revision N and
// both wrote N+1. Revisions are part of the durable audit trail, so every
// accepted change must get a distinct, higher value.
func TestFireActivationRevisionAdvancesOncePerCommit(t *testing.T) {
	for name, repo := range creditStores(t) {
		t.Run(name, func(t *testing.T) {
			store, ok := repo.(FireActivationStore)
			if !ok {
				t.Fatalf("%T does not implement FireActivationStore", repo)
			}
			initial, err := store.SetFireActivation(FireActivation{
				Mode: "shadow", AcknowledgedBy: "test", UpdatedAt: time.Now().UTC(),
			})
			if err != nil {
				t.Fatalf("initial activation: %v", err)
			}
			baseRevision := initial.Revision
			const workers = 6
			start := make(chan struct{})
			revisions := make(chan int64, workers)
			var wg sync.WaitGroup
			for i := 0; i < workers; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					<-start
					act, err := store.SetFireActivation(FireActivation{
						Mode: "active", Cohort: "pilot", AcknowledgedBy: "test",
						UpdatedAt: time.Now().UTC(),
					})
					if err != nil {
						t.Errorf("set activation: %v", err)
						return
					}
					revisions <- act.Revision
				}()
			}
			close(start)
			wg.Wait()
			close(revisions)

			seen := map[int64]bool{}
			for revision := range revisions {
				if seen[revision] {
					t.Fatalf("revision %d was reused", revision)
				}
				seen[revision] = true
			}
			final, err := store.GetFireActivation()
			if err != nil {
				t.Fatalf("get activation: %v", err)
			}
			if final == nil {
				t.Fatal("activation row missing")
			}
			if final.Revision != baseRevision+workers {
				t.Fatalf("final revision = %d, want %d", final.Revision, baseRevision+workers)
			}
		})
	}
}
