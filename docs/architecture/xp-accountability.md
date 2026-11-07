# Mathua XP Accountability

## Overview

Mathua uses two independent accountability mechanisms to maintain a fair XP economy:

1. **Performance accountability:** Evaluates the learner's performance across completed tasks.
2. **Effort accountability:** Detects behaviours that undermine meaningful learning, independently of task completion.

Both mechanisms can produce XP adjustments, but neither directly changes the learner's balance. All adjustments must pass through the XP ledger.

## 1. XP Accounting Model

Mathua distinguishes three values:

| Value | Definition |
|---|---|
| Gross XP | XP earned from correct answers. |
| Accountability adjustment | Actual XP added or removed by accountability rules. |
| Net XP | Gross XP plus the actual accountability adjustment. |

Formula:

`Net XP = Gross XP + Actual Accountability Adjustment`

The distinction between calculated and actual adjustments is essential.

A policy may calculate a negative adjustment, and the ledger applies it in full. Lifetime XP may therefore become negative; there is no artificial zero balance floor.

For example:

- Gross XP: +3
- Calculated adjustment: -5
- Actual adjustment: -5
- Net XP: -2
- Resulting lifetime balance: previous balance minus 5

At zero balance, the full calculated deduction is still applied, so the balance becomes negative. On a first application the actual adjustment equals the calculated adjustment; a duplicate replay applies zero because the event key already exists. The API response, learner-facing explanation and audit record must report the actual applied adjustment.

## 2. Performance Accountability

Performance accountability evaluates completed tasks against expected performance for their difficulty.

It can produce the following outcomes:

- Full XP: No performance deduction.
- Acceptable performance: Removes 25% of task-earned XP.
- Poor performance: Removes all task-earned XP.
- Extremely poor performance: Removes all task-earned XP and applies an additional 2 XP penalty.
- Genuine difficulty protection: Hard tasks can lose earned XP at most, without the additional negative penalty. Sufficiently difficult tasks with no correct answers can be exempt entirely.

Extremely poor performance requires:

- At least six graded attempts.
- Expected-correct total of at least 2.5.
- Performance ratio below the configured negative threshold.

### Completion requirements

Performance accountability only applies when the engine has a reliable completion signal.

| Task | Completion condition |
|---|---|
| Practice | Concept requirements completed. |
| Study | Concept reaches mastery during the task. |
| Regular quiz | Final graded question submitted. |
| Timed quiz | Final graded item completed. |
| Review | Final review answer, when `NextReviewQuestion` reports no questions remain. A durable `review_runs` row is finalized once for the whole run. |

Abandoned, unfinished and incomplete tasks do not receive performance deductions.

## 3. Effort Accountability

Effort accountability evaluates behaviour rather than overall task performance.

It can apply before task completion.

Triggers include:

- Repeated unusually fast incorrect answers.
- Repeated identical incorrect answers within the configured short-answer window.
- Farming already mastered material.
- Legacy rush rule: two very fast incorrect answers trigger a 5 XP penalty.

A single ordinary incorrect answer does not trigger an effort deduction.

## 4. Exemptions and Protections

Performance and effort accountability must respect applicable exemptions, including:

- Diagnostic assessments.
- Unfinished and abandoned tasks for performance accountability.
- Final skips that do not produce graded completion.
- Genuine-difficulty protection.

Performance and effort decisions remain independent.

## 5. Ledger and Idempotency

All XP adjustments must be applied through the established accounting path.

Requirements:

- Preserve atomicity between XP adjustments and associated task submissions.
- Prevent duplicate deductions on replay.
- Apply negative adjustments in full so lifetime XP may become negative.
- Persist the actual applied adjustment.
- Generate learner-facing explanations from the actual applied amount.
- Preserve the original adjustment and explanation on duplicate submissions.

## 6. Learner-Facing Transparency

Where task-level summaries are available, show:

- Gross XP earned.
- Actual accountability adjustment.
- Net XP change.

Never present a calculated deduction as an actual XP loss.

## 7. Architectural Invariants

1. Performance accountability requires reliable task completion.
2. Effort accountability does not require task completion.
3. Performance and effort are independent decisions.
4. The ledger is the authoritative source of actual XP movements.
5. Lifetime XP and daily net XP may both be negative.
6. Replaying a submission cannot apply another adjustment.
7. Learner-facing explanations must match actual ledger movements.

## Implementation References

- `internal/xp/accountability.go`
- `internal/engine/effort.go`
- `internal/engine/engine.go`
- `internal/engine/accountability_metrics.go`
- `internal/storage/effort.go`
- `internal/storage/submissions.go`
- `internal/storage/timed_quiz.go`
- `web/next-app/components/XPAdjustmentNote.tsx`
- `web/next-app/components/XPTaskSummary.tsx`

This document describes the XP accountability architecture and its behavioural invariants. Implementation changes must preserve these invariants unless a deliberate policy change is approved.
