package storage

import (
	"errors"
	"time"

	"github.com/chuma-beep/mathua/internal/repetition"
)

// Shadow credit accounting. Nothing here writes a production scheduling column:
// a projection stores a counterfactual due date on the ordinal and leaves
// concept_progress.next_review_due exactly as the explicit scheduler set it.
// That is what makes the whole feature reversible by dropping tables.
var (
	ErrCreditOrdinalMissing      = errors.New("credit: ordinal not found")
	ErrCreditOrdinalStale        = errors.New("credit: ordinal is stale")
	ErrCreditOrdinalStateChanged = errors.New("credit: ordinal state changed")
	ErrCreditClosedEpoch         = errors.New("credit: epoch is closed")
	ErrCreditFundsProjection     = errors.New("credit: void would leave a banked ordinal unfunded")
	ErrCreditDueConflict         = errors.New("credit: the due date moved during projection")
)

// CreditCandidate is one target the attempt is believed to exercise, with the
// raw credit before the per-attempt cap is applied.
type CreditCandidate struct {
	ConceptID string
	RawUnits  int64
}

// CreditEvent is one target's share of one attempt's credit, after the cap.
type CreditEvent struct {
	AttemptID      string
	ConceptID      string
	Epoch          int64
	ScaledUnits    int64
	RetainedUnits  int64
	DiscardedUnits int64
	DiscardReason  string
	Banked         int
}

// CreditAllocation is the committed record of one attempt. Raw and allocated
// units are different numbers on purpose: raw is the pre-cap sum, allocated is
// what was actually emitted, so a replay can be checked without recomputing.
type CreditAllocation struct {
	StudentID      string
	AttemptID      string
	Epoch          int64
	RawUnits       int64
	AllocatedUnits int64
	TargetCount    int
	DroppedCount   int
	Events         []CreditEvent
}

// CreditOrdinal is one banked or projected repetition.
type CreditOrdinal struct {
	ConceptID         string
	Epoch             int64
	Ordinal           int
	State             string
	BaselineDue       string
	CounterfactualDue string
	Attempts          int
}

// CreditOrdinalDetail adds the owner to a banked ordinal for the read-only
// dry-run report. It is separate so existing callers that only have a student
// ID are unchanged.
type CreditOrdinalDetail struct {
	StudentID string
	CreditOrdinal
	// Target progress fields let an active caller avoid opening a transaction
	// for a target whose due date is still in the future. ApplyCredit remains
	// the authority; these fields are only a prefilter.
	HasProgress     bool
	TargetStatus    string
	TargetDue       string
	PlacementSeeded int
	IsProvisional   int
}

// CreditAccount is the derived balance for one target in one epoch. Retained is
// the sum of non-voided retained units; Banked counts every ordinal in the
// epoch, including projected and stale ones, because MaxBank is a lifetime cap.
type CreditAccount struct {
	ConceptID     string
	Epoch         int64
	RetainedUnits int64
	Banked        int
	Projected     int
	Stale         int
}

// Balance is the fractional credit left after banking: always in [0, Micro).
func (a CreditAccount) Balance() int64 {
	return a.RetainedUnits - int64(a.Banked)*repetition.Micro
}

// Pending returns the number of unprojected banked ordinals waiting to satisfy reviews.
func (a CreditAccount) Pending() int {
	p := a.Banked - a.Projected - a.Stale
	if p < 0 {
		return 0
	}
	return p
}

// CreditStore is the slice of the repository this stage needs. Declared
// separately so both backends implement exactly these operations and no
// caller can reach a production scheduling column through this door.
type CreditStore interface {
	// AllocateCredit applies the per-attempt cap, admits what capacity allows,
	// banks whole repetitions and commits everything in one transaction. The
	// bool is false when this attempt was already allocated, in which case the
	// committed rows are returned unchanged: a replay never banks, stales or
	// projects.
	AllocateCredit(studentID, attemptID string, epoch int64, candidates []CreditCandidate) (CreditAllocation, bool, error)
	// AccountOf recomputes the balance from the ledger rather than trusting a
	// stored counter.
	AccountOf(studentID, conceptID string, epoch int64) (CreditAccount, error)
	// ProjectCredit turns banked ordinals into counterfactual due dates chained
	// from the real one. It never moves the real due date.
	ProjectCredit(studentID, conceptID string, epoch int64, stepDays int, now time.Time) ([]CreditOrdinal, error)
	// ResetCreditEpoch stales banked ordinals and returns the new epoch.
	ResetCreditEpoch(studentID, conceptID string) (int64, error)
	// VoidCreditEvent removes one event's credit without deleting the row.
	VoidCreditEvent(studentID, attemptID, conceptID, reason string) error
	// GetBankedOrdinals returns the count of unprojected banked ordinals per concept for a student.
	GetBankedOrdinals(studentID string) (map[string]int, error)
}

// CreditApplyRequest asks the store to apply at most one banked ordinal to one
// target. The store verifies mode-independent eligibility and the due-date
// baseline in the same transaction as the progress update.
type CreditApplyRequest struct {
	StudentID       string
	TargetConceptID string
	SourceConceptID string
	SourceAttemptID string
	Epoch           int64
	Ordinal         int
	// Cohort is the durable activation cohort the caller is running under.
	// ApplyCredit checks it against fire_activations in the same transaction.
	Cohort string
	Now    time.Time
}

// CreditApplication is the committed or replayed result of applying one
// ordinal. An ineligible request returns Applied=false with a Reason and no
// write. The due dates are RFC3339 strings so both stores report one shape.
type CreditApplication struct {
	StudentID        string
	TargetConceptID  string
	SourceConceptID  string
	SourceAttemptID  string
	Epoch            int64
	Ordinal          int
	Applied          bool
	Replayed         bool
	Reason           string
	ActualDue        string
	ProposedDue      string
	StabilityBefore  float64
	StabilityAfter   float64
	DifficultyBefore float64
	DifficultyAfter  float64
	ReviewsDeferred  int
}

// CreditOrdinalStore exposes the banked ordinal detail active mode needs to
// derive a stable application key. It is optional so a shadow-only fake does
// not have to implement it.
type CreditOrdinalStore interface {
	BankedOrdinalDetails(studentID string) ([]CreditOrdinalDetail, error)
	AllBankedOrdinalDetails() ([]CreditOrdinalDetail, error)
}

// CreditApplyStore is the active-mode application surface. It is separate from
// CreditStore so shadow-only stores and tests keep compiling without a second
// implementation.
type CreditApplyStore interface {
	ApplyCredit(req CreditApplyRequest) (CreditApplication, error)
	CreditApplications(studentID, targetConceptID string, limit int) ([]CreditApplication, error)
}

// CreditPreviewStore computes the same decision as ApplyCredit without writing.
// It exists only for the dry-run command.
type CreditPreviewStore interface {
	PreviewCredit(req CreditApplyRequest) (CreditApplication, error)
}

// FireActivation is the durable, server-side activation record. Environment
// variables can request Active, but only this record plus the fire.Config
// allowlist can authorize application.
type FireActivation struct {
	Mode           string
	Cohort         string
	AcknowledgedBy string
	AcknowledgedAt time.Time
	UpdatedAt      time.Time
	Revision       int64
}

// FireActivationStore persists the activation record and its audit history.
type FireActivationStore interface {
	GetFireActivation() (*FireActivation, error)
	SetFireActivation(a FireActivation) (FireActivation, error)
}
