package diagnostic

import (
	"sort"
	"strings"

	"github.com/chuma-beep/mathua/internal/concepts"
)

// DiagnosticReport is the adaptive-exam summary (improve.md:126).
type DiagnosticReport struct {
	PlacementCourseID string `json:"placement_course_id"`
	FrontierIdx       int    `json:"frontier_idx"`
	FrontierLabel     string `json:"frontier_label"`
	// FrontierConditional marks a provisional frontier: the frontier concept
	// was barely passed (MA "conditionally completed"), tasks assume it
	// known but must fall back along its prerequisites on struggle.
	FrontierConditional bool `json:"frontier_conditional"`
	// GapsByDomain lists concepts with enough evidence to call a gap: at
	// least two probes and belief below the known cutoff. A concept that was
	// never probed is not a gap, and neither is one with a single probe.
	GapsByDomain map[string][]string `json:"gaps_by_domain"`
	// UnassessedByDomain lists covered concepts the diagnostic never probed.
	// They are unknown, not demonstrated gaps.
	UnassessedByDomain map[string][]string `json:"unassessed_by_domain"`
	// InsufficientEvidenceByDomain lists concepts that were probed once but
	// did not reach the confidence needed to call a gap.
	InsufficientEvidenceByDomain map[string][]string `json:"insufficient_evidence_by_domain"`
	// ConditionallyCompleted lists probed concepts with belief just above
	// the known cutoff (MA: "barely place out"). Assumed known initially;
	// struggle on dependents falls back to these first.
	ConditionallyCompleted []string           `json:"conditionally_completed"`
	MasteryLevels          map[string]float64 `json:"mastery_levels"`
	Confidence             map[string]float64 `json:"confidence"`
	AvgConfidence          float64            `json:"avg_confidence"`
	// CompletionEstimates maps an example daily XP pace to the estimated time to
	// finish. It is empty here on purpose: the CAT knows nothing about a
	// destination, so there is no honest denominator for "complete" at this
	// layer. It used to carry a hardcoded {50: "≈5 days", 100: "≈10 days",
	// 300: "≈30 days"} that no code computed and that no learner's pace
	// matched. The handler fills it from the planner when the learner has a
	// saved destination, and leaves it null otherwise.
	CompletionEstimates map[int]string `json:"completion_estimates"`
	TotalQuestions      int            `json:"total_questions"`
	// Latent ability estimate under 2PL Item Response Theory.
	Theta float64 `json:"theta,omitempty"`
	// Standard error of measurement SE(theta).
	SE float64 `json:"se,omitempty"`
}

// Report builds the done-branch summary for a session.
func (e *Engine) Report(s *Session) *DiagnosticReport {
	s.Lock()
	defer s.Unlock()

	rep := &DiagnosticReport{
		FrontierIdx:                  -1,
		GapsByDomain:                 map[string][]string{},
		UnassessedByDomain:           map[string][]string{},
		InsufficientEvidenceByDomain: map[string][]string{},
		MasteryLevels:                map[string]float64{},
		Confidence:                   map[string]float64{},
		TotalQuestions:               s.totalAsked,
		Theta:                        s.theta,
		SE:                           s.se,
	}

	var confSum float64
	var confCount int
	conditional := map[string]bool{}
	for _, c := range s.order {
		b := s.beliefs[c.ID]
		conf := s.confidence[c.ID]
		if conf > 0 && b >= 0.05 && b <= 0.95 {
			rep.MasteryLevels[c.ID] = b
		}
		rep.Confidence[c.ID] = conf
		confSum += conf
		confCount++
		switch {
		case conf == 0:
			rep.UnassessedByDomain[c.Domain] = append(rep.UnassessedByDomain[c.Domain], c.Label)
		case b < beliefThreshold && conf < confidenceThreshold:
			rep.InsufficientEvidenceByDomain[c.Domain] = append(rep.InsufficientEvidenceByDomain[c.Domain], c.Label)
		case b < beliefThreshold:
			rep.GapsByDomain[c.Domain] = append(rep.GapsByDomain[c.Domain], c.Label)
		}
		// MA "conditionally completed": probed, barely above the known
		// cutoff, assumed known, but first in line for fall-back.
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
	// Deterministic order for every report bucket.
	for _, m := range []map[string][]string{
		rep.GapsByDomain,
		rep.UnassessedByDomain,
		rep.InsufficientEvidenceByDomain,
	} {
		for d := range m {
			sort.Strings(m[d])
		}
	}
	return rep
}

// domainToCourse maps the frontier concept's domain to a placement course id
// (ids match data/courses.json catalog). Advanced subdomains place into the
// *2 courses: calculus.series/sequences (Calc II: series), discrete
// graph_theory/proof/recurrence (Discrete II: Eulerian, strong induction),
// abstract rings/fields/structures (Abstract II: polynomial rings, fields).
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
		switch c.Subdomain {
		case "calculus.series", "calculus.sequences":
			return "calc2"
		}
		return "calc1"
	case strings.HasPrefix(d, "statistics"):
		return "probstat"
	case strings.HasPrefix(d, "linear_algebra"):
		return "linalg"
	case strings.HasPrefix(d, "geometry"):
		return "geo"
	case strings.HasPrefix(d, "discrete"):
		switch c.Subdomain {
		case "discrete_math.graph_theory", "discrete_math.proof", "discrete_math.recurrence":
			return "discrete2"
		}
		return "discrete"
	case strings.HasPrefix(d, "number_theory"):
		return "nt"
	case strings.HasPrefix(d, "complex"):
		return "a2"
	case strings.HasPrefix(d, "differential_equations"):
		return "diffeq"
	case strings.HasPrefix(d, "abstract"):
		switch c.Subdomain {
		case "abstract_algebra.rings", "abstract_algebra.fields", "abstract_algebra.structures":
			return "abstract2"
		}
		return "abstract"
	case strings.HasPrefix(d, "topology"):
		return "topo"
	case strings.HasPrefix(d, "machine_learning"):
		return "ml"
	case strings.HasPrefix(d, "real_analysis"):
		return "ra"
	case strings.HasPrefix(d, "graph_theory"):
		return "graph"
	case strings.HasPrefix(d, "numerical_analysis"):
		return "numa"
	case strings.HasPrefix(d, "probability_theory"):
		return "prob"
	}
	return ""
}
