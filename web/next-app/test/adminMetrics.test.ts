import { describe, expect, it } from 'vitest'
import { AdminAnalyticsSchema, AdminAnalyticsLiveSchema } from '../lib/api'
import { ADMIN_METRICS, adminMetric } from '../lib/adminMetrics'

const emptyPoint = { start: 0, value: 0 }

const analyticsFixture = {
  range: '7d',
  from: '2026-11-08T00:00:00Z',
  to: '2026-11-15T12:00:00Z',
  bucket_seconds: 86400,
  generated_at: '2026-11-15T12:00:01Z',
  partial_last_bucket: true,
  previous: {
    unique_learners: 1,
    unique_registered: 1,
    unique_guests: 0,
    unique_unclassified: 0,
    answers: 2,
    correct: 1,
    new_registrations: 0,
    registered_diagnostic_cohort: 0,
    quiz_completions: 0,
    timed_quiz_completions: 0,
    xp_awarded: 10,
    xp_deducted: 0,
    diagnostic_completions: 0,
    reports_opened: 0,
    reports_resolved: 0,
  },
  accounts: { total: 3, registered: 1, guest: 1, unclassified: 1 },
  totals: {
    unique_learners: 3,
    unique_registered: 1,
    unique_guests: 1,
    unique_unclassified: 1,
    answers: 3,
    correct: 2,
    new_registrations: 1,
    registered_diagnostic_cohort: 1,
    quiz_completions: 0,
    timed_quiz_completions: 0,
    xp_awarded: 10,
    xp_deducted: 0,
    diagnostic_completions: 1,
    reports_opened: 0,
    reports_resolved: 0,
  },
  activity: [{ start: 0, registered: 1, guests: 1, unclassified: 1, answers: 3, correct: 2 }],
  signups: [emptyPoint],
  mastery: [emptyPoint],
  xp: [{ start: 0, awarded: 10, deducted: 0 }],
  reports: [{ start: 0, opened: 0, resolved: 0 }],
  quizzes: [emptyPoint],
  timed_quizzes: [emptyPoint],
  diagnostic_completions: [emptyPoint],
  top_concepts: [{ concept_id: 'a', attempts: 3, correct: 2, reports: 0 }],
  most_reported: [],
  mastered_by_concept: { a: 1 },
  mastered_current: 1,
  report_status: { open: 0, resolved: 0 },
}

describe('admin metric registry', () => {
  it('has unique ids and every headline metric', () => {
    const ids = ADMIN_METRICS.map(metric => metric.id)
    expect(new Set(ids).size).toBe(ids.length)
    for (const id of [
      'online_now',
      'practicing_now',
      'unique_learners',
      'answers',
      'new_registrations',
      'accuracy',
      'mastered_current',
      'domain_coverage',
      'reports_opened',
      'delivery_status',
    ]) {
      expect(adminMetric(id), id).toBeDefined()
    }
    expect(adminMetric('not_a_metric')).toBeUndefined()
  })

  it('labels snapshot metrics as current, not range-driven', () => {
    for (const id of ['mastered_current', 'mastery_by_domain', 'domain_coverage', 'report_status', 'delivery_status']) {
      expect(adminMetric(id)?.kind, id).toBe('current')
      expect(['Current', 'All time'], id).toContain(adminMetric(id)?.window)
    }
  })

  it('marks the guest split as a derived classification', () => {
    expect(adminMetric('guest_activity')?.definition).toMatch(/guest/i)
    expect(adminMetric('unclassified_activity')?.caveat).toMatch(/never folded/i)
  })

  it('puts the retained-rows caveat on every event metric', () => {
    for (const metric of ADMIN_METRICS.filter(metric => metric.kind === 'event')) {
      expect(metric.caveat, metric.id).toMatch(/retained rows only/i)
    }
  })

  it('labels the diagnostic cohort as current state, not an event rate', () => {
    const metric = adminMetric('registered_diagnostic_cohort')
    expect(metric?.kind).toBe('current')
    expect(metric?.label).toMatch(/completed to date/i)
    expect(metric?.caveat).toMatch(/cohort, not an event rate/i)
  })
})

describe('admin analytics schema', () => {
  it('accepts the documented envelope', () => {
    const parsed = AdminAnalyticsSchema.parse(analyticsFixture)
    expect(parsed.range).toBe('7d')
    expect(parsed.accounts.unclassified).toBe(1)
    expect(parsed.previous?.unique_learners).toBe(1)
  })

  it('accepts a missing previous window', () => {
    const { previous: _previous, ...withoutPrevious } = analyticsFixture
    expect(AdminAnalyticsSchema.parse(withoutPrevious).previous).toBeUndefined()
  })

  it('allows a section response without totals', () => {
    const { totals: _totals, ...sectionOnly } = analyticsFixture
    expect(AdminAnalyticsSchema.parse(sectionOnly).totals).toBeUndefined()
  })

  it('accepts the live strip', () => {
    const parsed = AdminAnalyticsLiveSchema.parse({
      online_now: 2,
      online_window_seconds: 300,
      practicing_now: 3,
      guest_practice_now: 1,
      answers_15m: 12,
      active_sessions: 4,
      generated_at: '2026-11-15T12:00:01Z',
    })
    expect(parsed.practicing_now).toBe(3)
  })
})
