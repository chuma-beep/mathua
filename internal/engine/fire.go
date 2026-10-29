package engine

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"runtime/debug"
	"sort"
	"sync/atomic"
	"time"

	"github.com/chuma-beep/mathua/internal/fire"
	"github.com/chuma-beep/mathua/internal/storage"
)

// FIReMetrics records the process-local operation counters exposed by the FIRe
// status report. They are intentionally plain atomics: FIRe is a single-process
// feature on the current deployment and these counters are for operators, not
// for billing or durable accounting.
type FIReMetrics struct {
	Allocations   atomic.Int64
	Applications  atomic.Int64
	Conflicts     atomic.Int64
	Replays       atomic.Int64
	Retries       atomic.Int64
	Failures      atomic.Int64
	DueShifts     atomic.Int64
	LatencyMicros atomic.Int64
	LatencyCount  atomic.Int64
}

// FIReMetricsSnapshot is the JSON-friendly view of FIReMetrics.
type FIReMetricsSnapshot struct {
	Allocations   int64 `json:"allocations"`
	Applications  int64 `json:"applications"`
	Conflicts     int64 `json:"conflicts"`
	Replays       int64 `json:"replays"`
	Retries       int64 `json:"retries"`
	Failures      int64 `json:"failures"`
	DueShifts     int64 `json:"due_shifts"`
	LatencyMicros int64 `json:"latency_micros"`
	LatencyCount  int64 `json:"latency_count"`
}

func (m *FIReMetrics) Snapshot() FIReMetricsSnapshot {
	return FIReMetricsSnapshot{
		Allocations:   m.Allocations.Load(),
		Applications:  m.Applications.Load(),
		Conflicts:     m.Conflicts.Load(),
		Replays:       m.Replays.Load(),
		Retries:       m.Retries.Load(),
		Failures:      m.Failures.Load(),
		DueShifts:     m.DueShifts.Load(),
		LatencyMicros: m.LatencyMicros.Load(),
		LatencyCount:  m.LatencyCount.Load(),
	}
}

// FIReMode returns the configured mode. It is active only when the operator
// configuration passed the fail-closed activation gate.
func (e *Engine) FIReMode() fire.Mode { return e.fireMode }

// FIReConfig returns a copy of the operator configuration.
func (e *Engine) FIReConfig() fire.Config { return e.fireConfig }

// FIReMetrics returns the process counters for the status report.
func (e *Engine) FIReMetrics() FIReMetricsSnapshot { return e.fireStats.Snapshot() }

// FIReActivation returns the durable activation row, when one exists.
func (e *Engine) FIReActivation() (*storage.FireActivation, error) {
	store, ok := e.repo.(storage.FireActivationStore)
	if !ok {
		return nil, nil
	}
	return store.GetFireActivation()
}

// SetFIReActivation writes the durable activation row. The caller is
// responsible for the acknowledgement, cohort and allowlist checks. It does not
// change the running engine configuration; activation takes effect on restart.
func (e *Engine) SetFIReActivation(act storage.FireActivation) (storage.FireActivation, error) {
	store, ok := e.repo.(storage.FireActivationStore)
	if !ok {
		return storage.FireActivation{}, fmt.Errorf("fire: storage does not support activation")
	}
	return store.SetFireActivation(act)
}

// FIReDryRunRow is one proposed due-date shift for a banked ordinal.
type FIReDryRunRow struct {
	StudentID        string  `json:"student_id"`
	TargetConceptID  string  `json:"target_concept_id"`
	Epoch            int64   `json:"epoch"`
	Ordinal          int     `json:"ordinal"`
	Eligible         bool    `json:"eligible"`
	Reason           string  `json:"reason,omitempty"`
	ActualDue        string  `json:"actual_due,omitempty"`
	ProposedDue      string  `json:"proposed_due,omitempty"`
	DeltaDays        int     `json:"delta_days"`
	StabilityBefore  float64 `json:"stability_before"`
	StabilityAfter   float64 `json:"stability_after"`
	DifficultyBefore float64 `json:"difficulty_before"`
	DifficultyAfter  float64 `json:"difficulty_after"`
	Conflict         bool    `json:"conflict"`
	ReviewsDeferred  int     `json:"reviews_deferred"`
}

// FIReDryRunTotals summarises a dry-run report.
type FIReDryRunTotals struct {
	Targets         int `json:"targets"`
	Eligible        int `json:"eligible"`
	Skipped         int `json:"skipped"`
	Conflicts       int `json:"conflicts"`
	ReviewsDeferred int `json:"reviews_deferred"`
}

// FIReDryRunReport is the read-only projection of what active mode would do to
// banked ordinals. It never calls ApplyCredit or ProjectCredit.
type FIReDryRunReport struct {
	Mode        fire.Mode           `json:"mode"`
	Cohort      string              `json:"cohort"`
	Revision    string              `json:"revision"`
	GeneratedAt string              `json:"generated_at"`
	Checksum    string              `json:"checksum"`
	Totals      FIReDryRunTotals    `json:"totals"`
	Rows        []FIReDryRunRow     `json:"rows"`
	Metrics     FIReMetricsSnapshot `json:"metrics"`
}

// FIReDryRun computes the proposed active-mode effect without writing. It reads
// every banked ordinal and asks CreditPreviewStore what would happen now. A
// repository without preview support yields an empty report rather than an
// error, so the command is safe against an older database.
func (e *Engine) FIReDryRun(now time.Time) (FIReDryRunReport, error) {
	report := FIReDryRunReport{
		Mode:        e.fireMode,
		Cohort:      e.fireConfig.Cohort,
		Revision:    buildRevision(),
		GeneratedAt: now.UTC().Format(time.RFC3339),
		Metrics:     e.fireStats.Snapshot(),
		Rows:        []FIReDryRunRow{},
	}
	ordinals, ok := e.repo.(storage.CreditOrdinalStore)
	if !ok {
		return report, nil
	}
	preview, ok := e.repo.(storage.CreditPreviewStore)
	if !ok {
		return report, nil
	}
	details, err := ordinals.AllBankedOrdinalDetails()
	if err != nil {
		return report, fmt.Errorf("fire: read banked ordinals: %w", err)
	}
	type group struct {
		student string
		target  string
		epoch   int64
		lowest  storage.CreditOrdinalDetail
		count   int
	}
	groups := map[string]*group{}
	for _, d := range details {
		if e.fireMode == fire.Active && !e.fireConfig.ActiveFor(d.StudentID) {
			continue
		}
		key := fmt.Sprintf("%s\x00%s\x00%d", d.StudentID, d.ConceptID, d.Epoch)
		g := groups[key]
		if g == nil {
			g = &group{student: d.StudentID, target: d.ConceptID, epoch: d.Epoch, lowest: d, count: 1}
			groups[key] = g
			continue
		}
		g.count++
		if d.Ordinal < g.lowest.Ordinal {
			g.lowest = d
		}
	}
	keys := make([]string, 0, len(groups))
	for key := range groups {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		g := groups[key]
		res, err := preview.PreviewCredit(storage.CreditApplyRequest{
			StudentID:       g.student,
			TargetConceptID: g.target,
			Epoch:           g.epoch,
			Ordinal:         g.lowest.Ordinal,
			Now:             now.UTC(),
		})
		if err != nil {
			report.Rows = append(report.Rows, FIReDryRunRow{
				StudentID: g.student, TargetConceptID: g.target, Epoch: g.epoch,
				Ordinal: g.lowest.Ordinal, Reason: "preview_error",
			})
			continue
		}
		row := FIReDryRunRow{
			StudentID: g.student, TargetConceptID: g.target, Epoch: g.epoch,
			Ordinal: res.Ordinal, Eligible: res.Applied, Reason: res.Reason,
			ActualDue: res.ActualDue, ProposedDue: res.ProposedDue,
			StabilityBefore: res.StabilityBefore, StabilityAfter: res.StabilityAfter,
			DifficultyBefore: res.DifficultyBefore, DifficultyAfter: res.DifficultyAfter,
			ReviewsDeferred: res.ReviewsDeferred,
		}
		if res.Applied {
			// The preview covers the next application. Count every banked
			// ordinal for this target because active mode may apply one per
			// serving operation.
			row.ReviewsDeferred = g.count
		}
		if res.Reason == "stale_baseline" || res.Reason == "due_moved" || res.Reason == "epoch_closed" {
			row.Conflict = true
		}
		if from, okFrom := parseRFC3339(res.ActualDue); okFrom {
			if to, okTo := parseRFC3339(res.ProposedDue); okTo {
				row.DeltaDays = int(to.Sub(from).Hours() / 24)
			}
		}
		report.Rows = append(report.Rows, row)
	}
	for _, row := range report.Rows {
		report.Totals.Targets++
		if row.Eligible {
			report.Totals.Eligible++
			report.Totals.ReviewsDeferred += row.ReviewsDeferred
		} else {
			report.Totals.Skipped++
		}
		if row.Conflict {
			report.Totals.Conflicts++
		}
	}
	sum, err := dryRunChecksum(report.Rows)
	if err != nil {
		return report, err
	}
	report.Checksum = sum
	return report, nil
}

func dryRunChecksum(rows []FIReDryRunRow) (string, error) {
	payload, err := json.Marshal(rows)
	if err != nil {
		return "", fmt.Errorf("fire: checksum rows: %w", err)
	}
	sum := sha256.Sum256(payload)
	return hex.EncodeToString(sum[:]), nil
}

func parseRFC3339(raw string) (time.Time, bool) {
	if raw == "" {
		return time.Time{}, false
	}
	t, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return time.Time{}, false
	}
	return t, true
}

func buildRevision() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, s := range info.Settings {
			if s.Key == "vcs.revision" && len(s.Value) >= 7 {
				return s.Value[:7]
			}
		}
	}
	return "dev"
}
