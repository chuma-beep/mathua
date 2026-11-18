package engine

import (
	"fmt"
	"os"
	"time"

	"github.com/chuma-beep/mathua/internal/scheduler"
	"github.com/chuma-beep/mathua/internal/storage"
)

// ConceptEligibility is the learner-facing answer to "may I start this, and if
// not, why". The decision is made in internal/scheduler (EligibilityOf), which
// reuses the same prerequisite rule the recommendation engine and question
// scheduler use; this type only renders it.
type ConceptEligibility struct {
	ConceptID     string              `json:"conceptId"`
	Title         string              `json:"title"`
	Eligible      bool                `json:"eligible"`
	State         string              `json:"state"`
	Reason        string              `json:"reason,omitempty"`
	Prerequisites []EligibilityPrereq `json:"prerequisites"`
}

// EligibilityPrereq names one prerequisite and whether it is satisfied. `Met`
// is the server's verdict; the client renders it and never recomputes it.
type EligibilityPrereq struct {
	ConceptID string `json:"conceptId"`
	Title     string `json:"title"`
	State     string `json:"state"`
	Met       bool   `json:"met"`
}

// EligibilityFor reports whether the learner may enter the /learn loop for a
// concept. It is the one place the practice endpoint, the study-answer endpoint
// and the curriculum surfaces ask, so a locked topic cannot be started through
// any of them (invariant 3).
func (e *Engine) EligibilityFor(studentID, conceptID string) (*ConceptEligibility, error) {
	c := e.dag.Concept(conceptID)
	if c == nil {
		return nil, fmt.Errorf("%w: %q", ErrUnknownConcept, conceptID)
	}
	var progress map[string]*storage.ConceptProgress
	if studentID != "" {
		if p, err := e.repo.GetAllProgress(studentID); err == nil {
			progress = p
		}
	}
	snapshots := e.conceptSnapshotsFrom(studentID, progress)
	now := time.Now()
	el := scheduler.EligibilityOf(snapshots, now, c)

	missing := make(map[string]bool, len(el.Missing))
	for _, id := range el.Missing {
		missing[id] = true
	}

	out := &ConceptEligibility{
		ConceptID:     c.ID,
		Title:         c.Label,
		Eligible:      el.Eligible,
		State:         el.State,
		Reason:        el.Reason,
		Prerequisites: make([]EligibilityPrereq, 0, len(c.Prerequisites)),
	}
	for _, pid := range c.Prerequisites {
		pre := e.dag.Concept(pid)
		title := pid
		state := ""
		if pre != nil {
			title = pre.Label
			state = scheduler.EligibilityOf(snapshots, now, pre).State
		}
		out.Prerequisites = append(out.Prerequisites, EligibilityPrereq{
			ConceptID: pid,
			Title:     title,
			State:     state,
			Met:       !missing[pid],
		})
	}

	// MATHUA_DEV_UNLOCK=1 opens every concept on a local instance, so a reviewer
	// can read a lesson without first mastering its prerequisites. It lives here
	// because this is the single place the practice, study-answer and curriculum
	// surfaces ask, so the override cannot be applied to one surface and not
	// another. It is read per call and is never set in a deployment.
	if os.Getenv("MATHUA_DEV_UNLOCK") == "1" {
		out.Eligible = true
		out.State = "unlocked"
		out.Reason = ""
		for i := range out.Prerequisites {
			out.Prerequisites[i].Met = true
			out.Prerequisites[i].State = "mastered"
		}
	}
	return out, nil
}
