package diagnostic

import (
	"sort"
	"strings"

	"github.com/chuma-beep/mathua/internal/concepts"
)

// DiagnosticReport is the adaptive-exam summary (improve.md:126).
type DiagnosticReport struct {
	PlacementCourseID string              `json:"placement_course_id"`
	FrontierIdx       int                 `json:"frontier_idx"`
	FrontierLabel     string              `json:"frontier_label"`
	// FrontierConditional marks a provisional frontier: the frontier concept
	// was barely passed (MA "conditionally completed") — tasks assume it
	// known but must fall back along its prerequisites on struggle.
	FrontierConditional bool                `json:"frontier_conditional"`
	GapsByDomain        map[string][]string `json:"gaps_by_domain"`
	// ConditionallyCompleted lists probed concepts with belief just above
	// the known cutoff (MA: "barely place out"). Assumed known initially;
	// struggle on dependents falls back to these first.
	ConditionallyCompleted []string           `json:"conditionally_completed"`
	MasteryLevels          map[string]float64 `json:"mastery_levels"`
	Confidence             map[string]float64 `json:"confidence"`
	AvgConfidence          float64            `json:"avg_confidence"`
	CompletionEstimates map[int]string    `json:"completion_estimates"`
	TotalQuestions    int                 `json:"total_questions"`
}

// Report builds the done-branch summary for a session.
func (e *Engine) Report(s *Session) *DiagnosticReport {
	s.Lock()
	defer s.Unlock()

	rep := &DiagnosticReport{
		FrontierIdx:         -1,
		GapsByDomain:        map[string][]string{},
		MasteryLevels:       map[string]float64{},
		Confidence:          map[string]float64{},
		CompletionEstimates: map[int]string{150: "≈5 days", 300: "≈10 days", 900: "≈30 days"},
		TotalQuestions:      s.totalAsked,
	}

	var confSum float64
	var confCount int
	conditional := map[string]bool{}
	for _, c := range s.order {
		b := s.beliefs[c.ID]
		if b >= 0.05 && b <= 0.95 {
			rep.MasteryLevels[c.ID] = b
		}
		rep.Confidence[c.ID] = s.confidence[c.ID]
		confSum += s.confidence[c.ID]
		confCount++
		if b < beliefThreshold {
			rep.GapsByDomain[c.Domain] = append(rep.GapsByDomain[c.Domain], c.Label)
		}
		// MA "conditionally completed": probed, barely above the known
		// cutoff — assumed known, but first in line for fall-back.
		if s.confidence[c.ID] > 0 && b >= beliefThreshold && b < beliefThreshold+0.15 {
			conditional[c.ID] = true
			rep.ConditionallyCompleted = append(rep.ConditionallyCompleted, c.Label)
		}
		if b >= beliefThreshold && s.index(c.ID) > rep.FrontierIdx {
			rep.FrontierIdx = s.index(c.ID)
			rep.FrontierLabel = c.Label
		}
	}
	if confCount > 0 {
		rep.AvgConfidence = float64(int(confSum/float64(confCount)*100)) / 100
	}
	if rep.FrontierIdx >= 0 && rep.FrontierIdx < len(s.order) {
		rep.PlacementCourseID = domainToCourse(s.order[rep.FrontierIdx])
		rep.FrontierConditional = conditional[s.order[rep.FrontierIdx].ID]
	}
	// Deterministic gap order.
	for d := range rep.GapsByDomain {
		sort.Strings(rep.GapsByDomain[d])
	}
	return rep
}

// domainToCourse maps the frontier concept's domain to a placement course id
// (ids match data/courses.json catalog).
func domainToCourse(c *concepts.Concept) string {
	d := c.Domain
	switch {
	case strings.HasPrefix(d, "arithmetic"):
		return "4"
	case strings.HasPrefix(d, "fractions"):
		return "5"
	case strings.HasPrefix(d, "prealgebra"):
		return "pa"
	case strings.HasPrefix(d, "algebra"):
		return "a1"
	case strings.HasPrefix(d, "trigonometry"):
		return "a2"
	case strings.HasPrefix(d, "precalculus"):
		return "pc"
	case strings.HasPrefix(d, "calculus"):
		return "calc1"
	case strings.HasPrefix(d, "statistics"):
		return "probstat"
	case strings.HasPrefix(d, "linear_algebra"):
		return "linalg"
	case strings.HasPrefix(d, "geometry"):
		return "geo"
	case strings.HasPrefix(d, "discrete"):
		return "discrete"
	case strings.HasPrefix(d, "number_theory"):
		return "nt"
	case strings.HasPrefix(d, "complex"):
		return "a2"
	case strings.HasPrefix(d, "differential_equations"):
		return "diffeq"
	case strings.HasPrefix(d, "abstract"):
		return "abstract"
	case strings.HasPrefix(d, "topology"):
		return "topo"
	case strings.HasPrefix(d, "machine_learning"):
		return "ml"
	}
	return ""
}
