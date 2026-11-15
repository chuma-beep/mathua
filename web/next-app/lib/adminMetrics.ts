// One definition per metric used by the admin analytics surface.
//
// A metric's label, meaning and window live here, not inside a chart or a KPI
// card. That is what stops "Active learners" from silently meaning practice
// activity on one screen and browsing on another, and it is what lets the
// interface mark snapshot data as Current instead of pretending the range
// selector applies to it.
export type AdminMetricKind = 'event' | 'current' | 'derived'
export type AdminMetricFormat = 'number' | 'percent'

// Reset and delete remove attempts, XP events, progress and quiz rows. Every
// metric derived from those tables describes the rows still present, not a
// permanent history, and the info popover has to say so.
const RETAINED_ROWS = 'Retained rows only; reset or delete removes history.'

export interface AdminMetricDefinition {
  id: string
  label: string
  definition: string
  kind: AdminMetricKind
  window: string
  format: AdminMetricFormat
  caveat?: string
}

const metrics: AdminMetricDefinition[] = [
  {
    id: 'online_now',
    label: 'Online now',
    definition: 'Signed-in learners seen on any request in the last five minutes.',
    kind: 'current',
    window: 'Live, five minutes',
    format: 'number',
    caveat: 'Browsing presence, not practice activity. Resets on process restart.',
  },
  {
    id: 'practicing_now',
    label: 'Practicing now (15m)',
    definition: 'Signed-in learners with practice activity in the last fifteen minutes.',
    kind: 'current',
    window: 'Live, fifteen minutes',
    format: 'number',
    caveat: 'A practice session is recorded when a question is served or answered.',
  },
  {
    id: 'guest_practice_now',
    label: 'Guest practice (15m)',
    definition: 'Guest learners with practice activity in the last fifteen minutes.',
    kind: 'current',
    window: 'Live, fifteen minutes',
    format: 'number',
  },
  {
    id: 'answers_15m',
    label: 'Answers, 15m',
    definition: 'Answers submitted in the last fifteen minutes.',
    kind: 'event',
    window: 'Live, fifteen minutes',
    format: 'number',
  },
  {
    id: 'active_sessions',
    label: 'Active practice sessions',
    definition: 'Practice sessions touched in the last fifteen minutes.',
    kind: 'current',
    window: 'Live, fifteen minutes',
    format: 'number',
  },
  {
    id: 'unique_learners',
    label: 'Unique learners',
    definition: 'Distinct learners with at least one answer in the selected window.',
    kind: 'event',
    window: 'Selected range',
    format: 'number',
    caveat: 'One distinct count over the whole window, never the sum of daily uniques.',
  },
  {
    id: 'answers',
    label: 'Answers',
    definition: 'Answers submitted in the selected window.',
    kind: 'event',
    window: 'Selected range',
    format: 'number',
  },
  {
    id: 'accuracy',
    label: 'Accuracy',
    definition: 'Correct answers divided by all answers, with the answer count beside it.',
    kind: 'derived',
    window: 'Selected range',
    format: 'percent',
    caveat: 'A rate without n is not shown as a ranking.',
  },
  {
    id: 'new_registrations',
    label: 'New registrations',
    definition: 'Registered accounts created in the selected window.',
    kind: 'event',
    window: 'Selected range',
    format: 'number',
    caveat: 'Guest and unclassified rows are excluded.',
  },
  {
    id: 'registered_present',
    label: 'Registered accounts present',
    definition: 'Registered accounts currently present, not a historical total.',
    kind: 'current',
    window: 'Current',
    format: 'number',
    caveat: 'Accounts deleted or reset are absent from this count.',
  },
  {
    id: 'guest_activity',
    label: 'Guest activity',
    definition: 'Distinct guest learners with an answer in the bucket.',
    kind: 'event',
    window: 'Selected range',
    format: 'number',
  },
  {
    id: 'unclassified_activity',
    label: 'Unclassified activity',
    definition: 'Distinct unclassified rows with an answer in the bucket.',
    kind: 'event',
    window: 'Selected range',
    format: 'number',
    caveat: 'UUID guest rows without credentials. Never folded into registered or guest.',
  },
  {
    id: 'answers_per_learner',
    label: 'Answers per learner',
    definition: 'Answers divided by the range-wide distinct learner count.',
    kind: 'derived',
    window: 'Selected range',
    format: 'number',
  },
  {
    id: 'quiz_completions',
    label: 'Quiz completions',
    definition: 'Mastery-check quiz completions in the selected window.',
    kind: 'event',
    window: 'Selected range',
    format: 'number',
  },
  {
    id: 'timed_quiz_completions',
    label: 'Timed quiz completions',
    definition: 'Timed quiz runs completed in the selected window.',
    kind: 'event',
    window: 'Selected range',
    format: 'number',
  },
  {
    id: 'xp_awarded',
    label: 'XP awarded',
    definition: 'Sum of positive XP ledger events in the selected window.',
    kind: 'event',
    window: 'Selected range',
    format: 'number',
  },
  {
    id: 'xp_deducted',
    label: 'XP deducted',
    definition: 'Absolute sum of negative XP ledger events in the selected window.',
    kind: 'event',
    window: 'Selected range',
    format: 'number',
  },
  {
    id: 'net_xp',
    label: 'Net XP',
    definition: 'Awarded XP minus deducted XP in the selected window.',
    kind: 'derived',
    window: 'Selected range',
    format: 'number',
  },
  {
    id: 'mastered_current',
    label: 'Currently mastered',
    definition: 'Concepts currently mastered, grouped by the date they were mastered.',
    kind: 'current',
    window: 'Current',
    format: 'number',
    caveat: 'Re-mastery overwrites the date. This is not an award history.',
  },
  {
    id: 'mastery_by_domain',
    label: 'Mastery by domain',
    definition: 'Concepts currently mastered in each domain.',
    kind: 'current',
    window: 'Current',
    format: 'number',
  },
  {
    id: 'active_learning_days',
    label: 'Active learning days',
    definition: 'Days with at least one answer, counted per active learner.',
    kind: 'derived',
    window: 'Selected range',
    format: 'number',
  },
  {
    id: 'diagnostic_completed',
    label: 'Diagnostic completed',
    definition: 'Learners whose placement was completed in the selected window.',
    kind: 'event',
    window: 'Selected range',
    format: 'number',
    caveat: 'Diagnostic starts are not measured: the session row is deleted when a run finishes.',
  },
  {
    id: 'registered_diagnostic_cohort',
    label: 'Registered in window, diagnostic completed to date',
    definition: 'Accounts registered in the selected window that have completed a diagnostic by the time of this request.',
    kind: 'current',
    window: 'Registered in selected window; completion current',
    format: 'number',
    caveat: 'A cohort, not an event rate: completion may happen after the registration window closes. Visitor-to-signup is unavailable because no visitor event history exists.',
  },
  {
    id: 'domain_coverage',
    label: 'Domain coverage',
    definition: 'Concepts mastered by at least one learner over all concepts in the domain.',
    kind: 'current',
    window: 'Current',
    format: 'percent',
  },
  {
    id: 'domain_activity',
    label: 'Domain activity',
    definition: 'Answers submitted by domain in the selected window.',
    kind: 'event',
    window: 'Selected range',
    format: 'number',
  },
  {
    id: 'domain_accuracy',
    label: 'Domain accuracy',
    definition: 'Correct answers by domain, with the attempt count beside it.',
    kind: 'derived',
    window: 'Selected range',
    format: 'percent',
    caveat: 'Domains below fifty attempts are labelled as insufficient data.',
  },
  {
    id: 'reports_opened',
    label: 'Reports opened',
    definition: 'Question reports created in the selected window.',
    kind: 'event',
    window: 'Selected range',
    format: 'number',
  },
  {
    id: 'currently_resolved',
    label: 'Currently resolved',
    definition: 'Reports currently resolved, grouped by the resolution date.',
    kind: 'current',
    window: 'Current',
    format: 'number',
    caveat: 'Current state, not full transition history.',
  },
  {
    id: 'report_status',
    label: 'Report status',
    definition: 'Current count of reports by status.',
    kind: 'current',
    window: 'Current',
    format: 'number',
  },
  {
    id: 'top_concepts',
    label: 'Top concepts',
    definition: 'Concepts ranked by attempts, with accuracy and report counts.',
    kind: 'derived',
    window: 'Selected range',
    format: 'number',
    caveat: 'Accuracy ranking requires at least twenty attempts.',
  },
  {
    id: 'delivery_status',
    label: 'Delivery status',
    definition: 'Notification delivery counts by channel and status.',
    kind: 'current',
    window: 'All time',
    format: 'number',
    caveat: 'Reuses the existing delivery endpoint; it has no date filter.',
  },
  {
    id: 'delivery_problems',
    label: 'Delivery problems',
    definition: 'Recent failed or suppressed notification deliveries.',
    kind: 'current',
    window: 'Current',
    format: 'number',
  },
]

// Event metrics and the derived rates built from them describe retained rows.
// Current-state metrics do not: they describe the rows present right now.
export const ADMIN_METRICS: AdminMetricDefinition[] = metrics.map(metric => {
  const retained = metric.kind === 'event' || (metric.kind === 'derived' && metric.id !== 'registered_diagnostic_cohort')
  if (!retained) return metric
  return { ...metric, caveat: [metric.caveat, RETAINED_ROWS].filter(Boolean).join(' ') }
})

const byID = new Map(ADMIN_METRICS.map(metric => [metric.id, metric]))

export function adminMetric(id: string): AdminMetricDefinition | undefined {
  return byID.get(id)
}
