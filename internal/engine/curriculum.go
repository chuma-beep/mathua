package engine

import (
	"fmt"
	"time"

	"github.com/chuma-beep/mathua/internal/concepts"
	"github.com/chuma-beep/mathua/internal/scheduler"
	"github.com/chuma-beep/mathua/internal/storage"
)

// DomainSummary is one domain's counts for the curriculum picker. Every number
// is derived from the DAG and the learner's snapshots through the scheduler's
// eligibility rule, so the browser never reconstructs a state from raw progress.
type DomainSummary struct {
	ID              string `json:"id"`
	ConceptCount    int    `json:"conceptCount"`
	MasteredCount   int    `json:"masteredCount"`
	UnlockedCount   int    `json:"unlockedCount"`
	InProgressCount int    `json:"inProgressCount"`
}

// TopicSummary is one concept in a domain, with the state the server decided.
type TopicSummary struct {
	ConceptID string `json:"conceptId"`
	Title     string `json:"title"`
	State     string `json:"state"`
	Depth     int    `json:"depth"`
}

// DomainCurriculum is a domain's concepts in curriculum order.
type DomainCurriculum struct {
	Domain string         `json:"domain"`
	Topics []TopicSummary `json:"topics"`
}

// CurriculumDomains lists every domain with per-state counts. Order is the DAG's
// (alphabetical); the client applies its own pedagogical ordering for display.
func (e *Engine) CurriculumDomains(studentID string) ([]DomainSummary, error) {
	snapshots, err := e.learnerSnapshots(studentID)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	byID := make(map[string]*DomainSummary, len(e.dag.Domains()))
	for _, id := range e.dag.Domains() {
		byID[id] = &DomainSummary{ID: id}
	}
	for _, c := range e.dag.Order() {
		d := byID[c.Domain]
		if d == nil {
			continue
		}
		d.ConceptCount++
		switch scheduler.EligibilityOf(snapshots, now, c).State {
		case scheduler.StateMastered, scheduler.StateDueForReview:
			d.MasteredCount++
		case scheduler.StateInProgress:
			d.InProgressCount++
		case scheduler.StateUnlocked:
			d.UnlockedCount++
		}
	}
	out := make([]DomainSummary, 0, len(byID))
	for _, id := range e.dag.Domains() {
		out = append(out, *byID[id])
	}
	return out, nil
}

// CurriculumDomain lists one domain's concepts with their state and depth.
func (e *Engine) CurriculumDomain(studentID, domainID string) (*DomainCurriculum, error) {
	snapshots, err := e.learnerSnapshots(studentID)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	depths := conceptDepths(e.dag)
	out := &DomainCurriculum{Domain: domainID, Topics: []TopicSummary{}}
	for _, c := range e.dag.Order() {
		if c.Domain != domainID {
			continue
		}
		out.Topics = append(out.Topics, TopicSummary{
			ConceptID: c.ID,
			Title:     c.Label,
			State:     scheduler.EligibilityOf(snapshots, now, c).State,
			Depth:     depths[c.ID],
		})
	}
	if len(out.Topics) == 0 {
		return nil, fmt.Errorf("domain %q not found", domainID)
	}
	return out, nil
}

// learnerSnapshots is the shared load: progress plus the scheduler's view of it.
// An empty student id (anonymous) yields all-unseen snapshots.
func (e *Engine) learnerSnapshots(studentID string) (map[string]*scheduler.ConceptSnapshot, error) {
	var progress map[string]*storage.ConceptProgress
	if studentID != "" {
		p, err := e.repo.GetAllProgress(studentID)
		if err != nil {
			return nil, fmt.Errorf("load progress: %w", err)
		}
		progress = p
	}
	return e.conceptSnapshotsFrom(studentID, progress), nil
}

// conceptDepths is the longest prerequisite chain length ending at each concept,
// computed once per call. It is a display hint (how far into a domain a topic
// sits), not a gate; the DAG is acyclic by construction, so there are no cycles
// to guard against.
func conceptDepths(dag *concepts.DAG) map[string]int {
	depths := make(map[string]int, dag.Count())
	var visit func(c *concepts.Concept) int
	visit = func(c *concepts.Concept) int {
		if d, ok := depths[c.ID]; ok {
			return d
		}
		depth := 0
		for _, pid := range c.Prerequisites {
			if p := dag.Concept(pid); p != nil {
				if d := visit(p) + 1; d > depth {
					depth = d
				}
			}
		}
		depths[c.ID] = depth
		return depth
	}
	for _, c := range dag.Order() {
		visit(c)
	}
	return depths
}
