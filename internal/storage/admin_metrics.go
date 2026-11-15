package storage

import (
	"strings"
	"time"
)

// StudentClass is the derived account classification used by admin analytics.
// It is deliberately not a database column: the stores already carry enough
// evidence to tell a guest from a registered account, and a new column would
// need a backfill that the old rows cannot answer reliably.
type StudentClass string

const (
	StudentClassRegistered   StudentClass = "registered"
	StudentClassGuest        StudentClass = "guest"
	StudentClassUnclassified StudentClass = "unclassified"
)

// ClassifyStudent labels one student row.
//
// The guest prefix is the only reliable marker for a browser-created guest.
// A row without that prefix is registered when it carries any credential or
// identity. Everything else is unclassified: those are the UUID guest rows
// created by callers that did not use the guest prefix, and folding them into
// either bucket would make the audience mix lie.
func ClassifyStudent(id, username, passwordHash, email, googleID string) StudentClass {
	if strings.HasPrefix(id, "guest_") {
		return StudentClassGuest
	}
	if strings.TrimSpace(username) != "" ||
		strings.TrimSpace(passwordHash) != "" ||
		strings.TrimSpace(email) != "" ||
		strings.TrimSpace(googleID) != "" {
		return StudentClassRegistered
	}
	return StudentClassUnclassified
}

// AdminActivityPoint is one bucket of learner and answer activity. Start is
// Unix seconds at the bucket boundary. Registered, Guests and Unclassified are
// distinct learners in the bucket, not answer counts.
type AdminActivityPoint struct {
	Start        int64 `json:"start"`
	Registered   int   `json:"registered"`
	Guests       int   `json:"guests"`
	Unclassified int   `json:"unclassified"`
	Answers      int   `json:"answers"`
	Correct      int   `json:"correct"`
}

// AdminSeriesPoint is a single count in one bucket.
type AdminSeriesPoint struct {
	Start int64 `json:"start"`
	Value int   `json:"value"`
}

// AdminXPPoint is awarded and deducted XP in one bucket. Deducted is a
// positive magnitude so the chart can plot the two directions separately.
type AdminXPPoint struct {
	Start    int64 `json:"start"`
	Awarded  int   `json:"awarded"`
	Deducted int   `json:"deducted"`
}

// AdminReportPoint is reports opened and resolved in one bucket. Resolution is
// the current resolved_at value, not a transition history.
type AdminReportPoint struct {
	Start    int64 `json:"start"`
	Opened   int   `json:"opened"`
	Resolved int   `json:"resolved"`
}

// AdminConceptMetric is one concept's aggregate for the selected window.
type AdminConceptMetric struct {
	ConceptID string `json:"concept_id"`
	Attempts  int    `json:"attempts"`
	Correct   int    `json:"correct"`
	Reports   int    `json:"reports"`
}

// AdminAccountCounts is the current classification of every student row.
type AdminAccountCounts struct {
	Total        int `json:"total"`
	Registered   int `json:"registered"`
	Guest        int `json:"guest"`
	Unclassified int `json:"unclassified"`
}

// AdminMetricTotals are the range totals. Current-state counts live on
// AdminMetrics, not here, so a previous-window comparison cannot accidentally
// claim that mastery or report status was reconstructed historically.
type AdminMetricTotals struct {
	UniqueLearners       int `json:"unique_learners"`
	UniqueRegistered     int `json:"unique_registered"`
	UniqueGuests         int `json:"unique_guests"`
	UniqueUnclassified   int `json:"unique_unclassified"`
	Answers              int `json:"answers"`
	Correct              int `json:"correct"`
	NewRegistrations     int `json:"new_registrations"`
	QuizCompletions      int `json:"quiz_completions"`
	TimedQuizCompletions int `json:"timed_quiz_completions"`
	XPAwarded            int `json:"xp_awarded"`
	XPDeducted           int `json:"xp_deducted"`
	DiagnosticCompleted  int `json:"diagnostic_completed"`
	ReportsOpened        int `json:"reports_opened"`
	ReportsResolved      int `json:"reports_resolved"`
}

// AdminMetrics is one historical window for the admin dashboard. Current-state
// fields are labelled as such by the API; they are not part of Totals.
type AdminMetrics struct {
	Accounts          AdminAccountCounts   `json:"accounts"`
	Totals            AdminMetricTotals    `json:"totals"`
	Activity          []AdminActivityPoint `json:"activity"`
	Signups           []AdminSeriesPoint   `json:"signups"`
	Mastery           []AdminSeriesPoint   `json:"mastery"`
	XP                []AdminXPPoint       `json:"xp"`
	Reports           []AdminReportPoint   `json:"reports"`
	Quizzes           []AdminSeriesPoint   `json:"quizzes"`
	TimedQuizzes      []AdminSeriesPoint   `json:"timed_quizzes"`
	Diagnostics       []AdminSeriesPoint   `json:"diagnostics"`
	TopConcepts       []AdminConceptMetric `json:"top_concepts"`
	MostReported      []AdminConceptMetric `json:"most_reported"`
	MasteredByConcept map[string]int       `json:"mastered_by_concept"`
	MasteredCurrent   int                  `json:"mastered_current"`
	ReportStatus      map[string]int       `json:"report_status"`
}

// AdminLiveMetrics is the uncached live strip. It deliberately carries no
// history: these values are current state, and the dashboard must not imply
// that a restart-preserved trend exists.
type AdminLiveMetrics struct {
	Answers15m     int `json:"answers_15m"`
	ActiveSessions int `json:"active_sessions"`
}

// AdminMetricsStore is separate from Repository so existing test doubles do
// not need every analytics method. The server type-asserts its repository to
// this interface and returns a clear error when a test double does not
// implement it.
type AdminMetricsStore interface {
	AdminMetrics(since, until time.Time, bucket time.Duration) (*AdminMetrics, error)
	AdminLiveMetrics(since time.Time) (AdminLiveMetrics, error)
}

// adminBucketStarts walks aligned bucket boundaries from since to until.
// Since may be mid-bucket for rolling ranges; the first bucket then contains
// only data at or after the filter, which the API reports as the requested
// window rather than inventing earlier data.
// adminMergeReport combines the opened and resolved queries into one point
// per bucket. The two queries are separate so each can use its own timestamp
// column and index.
func adminMergeReport(points []AdminReportPoint, p AdminReportPoint) []AdminReportPoint {
	for i := range points {
		if points[i].Start == p.Start {
			points[i].Opened += p.Opened
			points[i].Resolved += p.Resolved
			return points
		}
	}
	return append(points, p)
}

func adminBucketStarts(since, until time.Time, bucket time.Duration) []int64 {
	if bucket <= 0 {
		bucket = 24 * time.Hour
	}
	start := since.UTC().Truncate(bucket)
	starts := make([]int64, 0, int(until.Sub(since)/bucket)+2)
	for t := start; t.Before(until.UTC()); t = t.Add(bucket) {
		starts = append(starts, t.Unix())
		if len(starts) > 10000 {
			break
		}
	}
	return starts
}
