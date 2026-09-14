package diagnostic

import (
	"strings"
	"testing"
)

// RetryServe voids a surprising miss and restores the exact prior state:
// same question re-served, beliefs restored, counts decremented.
func TestRetry_EligibleVoidRestores(t *testing.T) {
	e := NewEngine(testDAG(t), mustRegistry(t))
	s := e.Start()

	p1, cid1, err := e.NextQuestion(s)
	if err != nil {
		t.Fatalf("next question: %v", err)
	}
	// Established knowledge: high pre-answer belief.
	s.beliefs[cid1] = 0.85
	e.RecordAnswerTimed(s, cid1, false, 5.0, 10.0)
	if got := s.beliefs[cid1]; got >= 0.85 {
		t.Fatalf("expected belief drop after miss, got %v", got)
	}

	if _, _, err := e.NextQuestion(s); err != nil {
		t.Fatalf("next question: %v", err)
	}
	if !e.RetryAvailableFor(s, cid1, s.LastConceptID) {
		t.Fatal("expected retry offer for a surprising miss")
	}

	rp, rcid, _, err := e.RetryServe(s, cid1)
	if err != nil {
		t.Fatalf("retry: %v", err)
	}
	if rcid != cid1 || rp.Question != p1.Question {
		t.Errorf("expected restored %q (%q), got %q (%q)", cid1, p1.Question, rcid, rp.Question)
	}
	if got := s.beliefs[cid1]; got != 0.85 {
		t.Errorf("expected belief restored to 0.85, got %v", got)
	}
	if len(s.Attempts) != 0 {
		t.Errorf("expected voided attempts, got %d", len(s.Attempts))
	}
	if s.doneSet[cid1] {
		t.Errorf("expected %q unsettled after void", cid1)
	}
	if s.LastConceptID != cid1 {
		t.Errorf("expected pending %q, got %q", cid1, s.LastConceptID)
	}

	// Once per question.
	if _, _, _, err := e.RetryServe(s, cid1); err == nil {
		t.Error("expected error on second retry of the same question")
	} else if !strings.Contains(err.Error(), "already used") {
		t.Errorf("expected already-used error, got %v", err)
	}
}

// Each trigger in isolation, via crafted misses (no CAT randomness).
func TestRetry_Triggers(t *testing.T) {
	e := NewEngine(testDAG(t), mustRegistry(t))
	newMiss := func(concept string, belief float64, probes int, dontKnow bool) *Session {
		s := e.Start()
		s.probeCounts[concept] = probes
		s.Attempts = append(s.Attempts, Attempt{
			ConceptID:      concept,
			Correct:        false,
			ElapsedSeconds: 5.0,
			BeliefSnapshot: map[string]float64{concept: belief},
			DontKnow:       dontKnow,
		})
		s.totalCount[concept] = 1
		s.totalAsked = 1
		return s
	}

	// Level drop: low-belief miss on "c", staged next "a" sits earlier.
	// (testDAG is a linear a→b→c→d… chain.)
	if s := newMiss("c", 0.4, 2, false); !e.RetryAvailableFor(s, "c", "a") {
		t.Error("level drop should qualify a low-belief miss")
	}
	// Same miss staged forward: no trigger fires.
	if s := newMiss("c", 0.4, 2, false); e.RetryAvailableFor(s, "c", "d") {
		t.Error("low-belief miss on hard material staged forward must not qualify")
	}
	// Really simple: first-probe miss at low belief still qualifies.
	if s := newMiss("c", 0.4, 1, false); !e.RetryAvailableFor(s, "c", "d") {
		t.Error("first-probe (base difficulty) miss should qualify")
	}
	// Admits never qualify, even at high belief.
	if s := newMiss("c", 0.9, 1, true); e.RetryAvailableFor(s, "c", "d") {
		t.Error("admitted unknowns must never qualify")
	}
	// Established knowledge baseline.
	if s := newMiss("c", 0.85, 3, false); !e.RetryAvailableFor(s, "c", "d") {
		t.Error("established knowledge should qualify")
	}
}

// Genuine unknowns on hard material never get the offer.
func TestRetry_IneligibleUnknown(t *testing.T) {
	e := NewEngine(testDAG(t), mustRegistry(t))
	s := e.Start()

	cid, err := func() (string, error) {
		_, cid, err := e.NextQuestion(s)
		return cid, err
	}()
	if err != nil {
		t.Fatalf("next question: %v", err)
	}
	// Re-probe tier: not "really simple", belief untouched at 0.5.
	s.probeCounts[cid] = 2
	e.RecordAnswerTimed(s, cid, false, 5.0, 10.0)
	if _, _, err := e.NextQuestion(s); err != nil {
		t.Fatalf("next question: %v", err)
	}
	if e.RetryAvailableFor(s, cid, s.LastConceptID) {
		t.Error("expected no retry offer for a genuine unknown on hard material")
	}
	if _, _, _, err := e.RetryServe(s, cid); err == nil {
		t.Error("expected error retrying a genuine unknown")
	} else if !strings.Contains(err.Error(), "slip") {
		t.Errorf("expected not-a-slip error, got %v", err)
	}
}

func TestRetry_Rejections(t *testing.T) {
	e := NewEngine(testDAG(t), mustRegistry(t))

	// Correct answers cannot be retried.
	s := e.Start()
	cid, err := func() (string, error) {
		_, cid, err := e.NextQuestion(s)
		return cid, err
	}()
	if err != nil {
		t.Fatalf("next question: %v", err)
	}
	s.beliefs[cid] = 0.9
	e.RecordAnswerTimed(s, cid, true, 5.0, 10.0)
	if _, _, err := e.NextQuestion(s); err != nil {
		t.Fatalf("next question: %v", err)
	}
	if _, _, _, err := e.RetryServe(s, cid); err == nil {
		t.Error("expected error retrying a correct answer")
	} else if !strings.Contains(err.Error(), "correct") {
		t.Errorf("expected correct-answer error, got %v", err)
	}

	// Nothing pending to retry on a fresh session.
	s2 := e.Start()
	if _, _, _, err := e.RetryServe(s2, "a"); err == nil {
		t.Error("expected error with no attempts")
	}

	// Complete sessions stay complete.
	s3 := e.Start()
	s3.State = StateDone
	if _, _, _, err := e.RetryServe(s3, "a"); err == nil {
		t.Error("expected error on a complete session")
	}

	// Stale concept: the voidable attempt belongs to another question.
	s4 := e.Start()
	if _, cid4, err := e.NextQuestion(s4); err != nil {
		t.Fatalf("next question: %v", err)
	} else {
		s4.beliefs[cid4] = 0.9
		e.RecordAnswerTimed(s4, cid4, false, 5.0, 10.0)
	}
	if _, _, err := e.NextQuestion(s4); err != nil {
		t.Fatalf("next question: %v", err)
	}
	if _, _, _, err := e.RetryServe(s4, "zzz"); err == nil {
		t.Error("expected error on concept mismatch")
	} else if !strings.Contains(err.Error(), "match") {
		t.Errorf("expected mismatch error, got %v", err)
	}
}
