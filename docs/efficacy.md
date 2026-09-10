# Efficacy Instrumentation

Mathua records every graded attempt and derives efficacy metrics from that log.
This document defines the metrics, the endpoints that expose them, and the
honest limits of what they can and cannot claim.

## What is recorded

Every answer writes one `attempts` row (`internal/storage/store.go`
`AttemptEntry`): student, concept, correctness, elapsed seconds, timestamp.
Nothing is sampled or aggregated at write time — the metrics below are all
computed from the raw log.

## Point-in-time metrics

`GET /api/efficacy` (authenticated, per student) and `GET /api/efficacy/all`
(public, product-wide) return `EfficacyReport`
(`internal/engine/efficacy.go`):

| Field | Meaning |
|---|---|
| `concepts_touched` | distinct (student, concept) pairs seen |
| `first_pass_rate` | fraction of concepts answered correctly on the **first** attempt |
| `second_pass_rate` | fraction correct within the **first two** attempts |
| `avg_attempts_per_concept` | mean attempts per concept |
| `total_attempts` | raw attempt count |
| `students_tracked` | distinct students (aggregate only) |

These mirror the Math Academy-style first-pass / second-pass framing
(`improve.md:82`). They are **observational**: they describe learners who
chose to use Mathua, not a controlled trial.

## Longitudinal metrics

`GET /api/efficacy/trend` (public, rate-limited) returns `EfficacyTrend`
(`internal/engine/efficacy.go` `EfficacyTrend`):

```json
{
  "weeks": [
    { "week_start": "2026-09-07", "attempts": 128, "concepts_touched": 41,
      "first_pass_rate": 0.61, "second_pass_rate": 0.83, "active_students": 12 }
  ],
  "total_students": 40,
  "returning_students": 17,
  "retention_rate": 0.425,
  "first_pass_trend": 0.08
}
```

- **Weekly buckets** run Monday–Sunday in UTC (`weekStartUTC`). Each week's
  pass rates use the same definition as the point-in-time report, computed
  over that week's attempts only.
- **`retention_rate`** = students active in two or more distinct weeks ÷ all
  students who ever attempted. It is a coarse cohort signal, not a
  week-over-week retention curve.
- **`first_pass_trend`** = last week's first-pass rate minus the first week's.
  Positive means accuracy improved across the observed window. It is noisy at
  small N and should be read alongside `active_students`.

The Profile page renders the last 12 weeks as a first-pass sparkline plus the
retention figure (`web/next-app/app/profile/page.tsx`).

## How to read it

1. **Compare within a student, not across.** A learner's own weekly first-pass
   rate is a more honest signal than the product-wide number, which mixes
   learners at very different frontiers.
2. **Watch `active_students` before the rate.** A week with three attempts can
   show 0% or 100% by accident.
3. **Second-pass is the recovery metric.** A gap between first- and second-pass
   rates means material is learnable with feedback but not yet automatic.

## Limitations

- No control group and no randomisation: these are usage metrics, not an
  efficacy claim. Mathua does **not** assert the 4× / 93% / 98% figures from
  Math Academy's marketing.
- Retention is binary (returned / did not) over the whole window; it does not
  model decay or churn timing.
- Buckets are UTC weeks; a learner near a timezone boundary can appear in two
  weeks for one study session.
- The attempt log is the source of truth. If it is cleared, all metrics reset.

## Reproducing an analysis

```bash
# product-wide point-in-time
curl -s localhost:8080/api/efficacy/all | jq
# weekly trend + retention
curl -s localhost:8080/api/efficacy/trend | jq
```

Both read only the local database; no external service or API key is involved.
