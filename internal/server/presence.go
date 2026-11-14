package server

import (
	"sync"
	"time"
)

// onlineWindow is how long one authenticated request keeps a learner counted
// as online. It is shorter than the database-backed active window because it
// answers a different question: who is browsing right now, not who answered a
// question in the last fifteen minutes.
const onlineWindow = 5 * time.Minute

// presenceTracker is an in-memory view of signed-in learners seen recently.
// The deployment runs one process, so no shared store is needed. A restart
// clears the map, which is honest: nobody is known to be online until they
// make another request.
type presenceTracker struct {
	mu      sync.Mutex
	seen    map[string]time.Time
	now     func() time.Time
	touches int
}

func newPresenceTracker(now func() time.Time) *presenceTracker {
	if now == nil {
		now = time.Now
	}
	return &presenceTracker{seen: make(map[string]time.Time), now: now}
}

// Touch records a request from a signed-in learner. Empty ids are ignored so a
// guest flow cannot inflate the signed-in count.
func (p *presenceTracker) Touch(studentID string) {
	if p == nil || studentID == "" {
		return
	}
	now := p.now().UTC()
	p.mu.Lock()
	p.seen[studentID] = now
	p.touches++
	// Prune amortised rather than on every request: the map is read by the
	// admin overview far less often than it is written by traffic.
	if p.touches%256 == 0 {
		p.pruneLocked(now)
	}
	p.mu.Unlock()
}

// Count reports how many learners have made a request inside the window. It
// prunes first so a quiet period cannot leave stale ids counted.
func (p *presenceTracker) Count(now time.Time) int {
	if p == nil {
		return 0
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.pruneLocked(now.UTC())
	return len(p.seen)
}

func (p *presenceTracker) pruneLocked(now time.Time) {
	cutoff := now.Add(-onlineWindow)
	for id, at := range p.seen {
		if !at.After(cutoff) {
			delete(p.seen, id)
		}
	}
}
