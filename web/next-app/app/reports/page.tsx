'use client'

import { useCallback, useEffect, useMemo, useState } from 'react'
import Link from 'next/link'
import AppShell from '@/components/app-shell'
import Footer from '@/components/Footer'
import SectionHeader from '@/components/SectionHeader'
import { getLearnerReport, getLearnerReportHistory } from '@/lib/api'
import type {
  LearnerAssessment,
  LearnerDomainSummary,
  LearnerHistoryEvent,
  LearnerMasteryRow,
  LearnerReport,
} from '@/lib/api'

const PERIODS = [
  { days: 7, label: 'Last 7 days' },
  { days: 30, label: 'Last 30 days' },
  { days: 90, label: 'Last 90 days' },
  { days: 0, label: 'All history' },
]

const TABS = [
  { id: 'overview', label: 'Overview' },
  { id: 'mastery', label: 'Mastery' },
  { id: 'assessments', label: 'Assessments' },
  { id: 'history', label: 'Learning history' },
] as const

type TabID = (typeof TABS)[number]['id']

function fmtDate(value?: string): string {
  if (!value) return 'Not recorded'
  const d = new Date(value)
  if (Number.isNaN(d.getTime())) return 'Not recorded'
  return d.toLocaleDateString(undefined, { year: 'numeric', month: 'short', day: 'numeric' })
}

function fmtDateTime(value?: string): string {
  if (!value) return 'Not recorded'
  const d = new Date(value)
  if (Number.isNaN(d.getTime())) return 'Not recorded'
  return d.toLocaleString(undefined, { year: 'numeric', month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' })
}

function stateLabel(row: LearnerMasteryRow): string {
  if (row.decaying) return 'Due for review'
  switch (row.state) {
    case 'mastered':
      return 'Mastered'
    case 'due_for_review':
      return 'Due for review'
    case 'in_progress':
      return 'In progress'
    case 'unlocked':
      return 'Accessible'
    case 'upcoming':
      return 'Not yet accessible'
    case 'locked':
      return 'Locked'
    default:
      return row.state || 'Unknown'
  }
}

function stateTone(row: LearnerMasteryRow): string {
  if (row.decaying || row.state === 'due_for_review') return 'border-[var(--accent-red)]/40 text-[var(--accent-red)]'
  switch (row.state) {
    case 'mastered':
      return 'border-[var(--accent-green)]/40 text-[var(--accent-green)]'
    case 'in_progress':
      return 'border-[var(--accent-blue)]/40 text-[var(--accent-blue)]'
    default:
      return 'border-[var(--border)] text-mathua-secondary'
  }
}

function Metric({ label, value, note }: { label: string; value: string; note: string }) {
  return (
    <div className="border border-[var(--border)] rounded-none p-4 min-w-0">
      <div className="font-mono text-[10px] uppercase tracking-wider text-mathua-muted">{label}</div>
      <div className="text-2xl font-semibold mt-1 tabular-nums">{value}</div>
      <p className="text-xs text-mathua-secondary mt-1 leading-snug">{note}</p>
    </div>
  )
}

export default function ReportsPage() {
  const [tab, setTab] = useState<TabID>('overview')
  const [days, setDays] = useState(30)
  const [domain, setDomain] = useState('')
  const [domainOptions, setDomainOptions] = useState<LearnerDomainSummary[]>([])
  const [report, setReport] = useState<LearnerReport | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [visibleMastery, setVisibleMastery] = useState(100)
  const [history, setHistory] = useState<LearnerHistoryEvent[]>([])
  const [historyTotal, setHistoryTotal] = useState(0)
  const [historyLoading, setHistoryLoading] = useState(false)

  const load = useCallback(async (nextDays: number, nextDomain: string) => {
    setLoading(true)
    setError('')
    try {
      const data = await getLearnerReport({ days: nextDays, domain: nextDomain || undefined })
      setReport(data)
      setHistory(data.history ?? [])
      setHistoryTotal(data.history_total)
      setVisibleMastery(100)
      if (!nextDomain) setDomainOptions(data.domains ?? [])
    } catch (err) {
      setError(err instanceof Error ? err.message : 'The report could not be loaded.')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    void load(days, domain)
  }, [days, domain, load])

  const loadMoreHistory = useCallback(async () => {
    if (!report || historyLoading || history.length >= historyTotal) return
    setHistoryLoading(true)
    try {
      const page = await getLearnerReportHistory({ limit: 25, offset: history.length })
      setHistory((current) => {
        const seen = new Set(current.map((e) => `${e.kind}:${e.id}`))
        return [...current, ...page.events.filter((e) => !seen.has(`${e.kind}:${e.id}`))]
      })
      setHistoryTotal(page.total)
    } catch {
      setError('More history could not be loaded.')
    } finally {
      setHistoryLoading(false)
    }
  }, [history.length, historyLoading, historyTotal, report])

  const activity = useMemo(() => (report?.activity ?? []).slice(-90), [report])
  const maxQuestions = useMemo(() => activity.reduce((max, day) => Math.max(max, day.questions), 1), [activity])
  const mastery = report?.mastery ?? []
  const assessments = report?.assessments ?? []

  return (
    <AppShell>
      <main className="max-w-container mx-auto px-4 sm:px-6 pb-[calc(80px+env(safe-area-inset-bottom))] lg:pb-0 overflow-x-hidden min-w-0">
        <div className="print-report">
          <section className="pt-8 min-w-0 overflow-hidden">
            <div className="flex flex-wrap items-start justify-between gap-4 report-no-print">
              <div className="min-w-0">
                <h1 className="text-2xl font-semibold">Reports</h1>
                <p className="text-sm text-mathua-secondary mt-1 max-w-2xl">
                  A record of what you have learned, the evidence behind it and what you have attempted. It is built
                  from stored progress, answers and completions, not from estimates.
                </p>
              </div>
              <button
                type="button"
                onClick={() => window.print()}
                disabled={!report || loading}
                className="h-9 px-4 font-mono text-[11px] uppercase tracking-wider border border-[var(--border)] rounded-none hover:bg-[var(--surface)] disabled:opacity-50"
              >
                Print report
              </button>
            </div>

            {report && (
              <p className="text-xs text-mathua-muted mt-3">
                {report.learner.name || 'Learner'} · generated {fmtDateTime(report.generated_at)} · {report.period.label}
              </p>
            )}

            <div className="flex flex-wrap gap-2 mt-5 report-no-print" role="group" aria-label="Reporting period">
              {PERIODS.map((p) => (
                <button
                  key={p.days}
                  type="button"
                  onClick={() => setDays(p.days)}
                  aria-pressed={days === p.days}
                  className={`px-3 h-8 font-mono text-[11px] uppercase tracking-wider border rounded-none transition-colors ${
                    days === p.days
                      ? 'border-mathua-blue bg-mathua-blue text-white'
                      : 'border-mathua-border text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue'
                  }`}
                >
                  {p.label}
                </button>
              ))}
              <label className="sr-only" htmlFor="report-domain">
                Domain
              </label>
              <select
                id="report-domain"
                value={domain}
                onChange={(event) => setDomain(event.target.value)}
                className="h-8 px-2 border border-[var(--border)] rounded-none bg-[var(--bg)] font-mono text-[11px] uppercase tracking-wider"
              >
                <option value="">All domains</option>
                {domainOptions.map((d) => (
                  <option key={d.id} value={d.id}>
                    {d.id}
                  </option>
                ))}
              </select>
            </div>

            {error && (
              <div className="mt-5 border border-[var(--accent-red)]/40 rounded-none p-4" role="alert">
                <p className="text-sm">{error}</p>
                <button
                  type="button"
                  onClick={() => void load(days, domain)}
                  className="mt-2 text-sm underline underline-offset-4"
                >
                  Try again
                </button>
              </div>
            )}

            {loading && !report && (
              <div className="mt-6 grid gap-3 sm:grid-cols-2 lg:grid-cols-4" aria-busy="true" aria-label="Loading report">
                {[0, 1, 2, 3].map((i) => (
                  <div key={i} className="h-28 rounded-none border border-[var(--border)] bg-[var(--surface)] animate-pulse" />
                ))}
              </div>
            )}

            {report && (
              <>
                <div className="flex flex-wrap gap-2 mt-6 report-no-print" role="tablist" aria-label="Report sections">
                  {TABS.map((t) => (
                    <button
                      key={t.id}
                      type="button"
                      role="tab"
                      aria-selected={tab === t.id}
                      aria-controls={`panel-${t.id}`}
                      id={`tab-${t.id}`}
                      onClick={() => setTab(t.id)}
                      className={`px-3 py-2 text-sm border rounded-none transition-colors ${
                        tab === t.id
                          ? 'border-mathua-blue bg-mathua-blue text-white font-medium'
                          : 'border-mathua-border text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue'
                      }`}
                    >
                      {t.label}
                    </button>
                  ))}
                </div>

                <div
                  id="panel-overview"
                  role="tabpanel"
                  aria-labelledby="tab-overview"
                  data-report-panel
                  className={`pt-6 ${tab === 'overview' ? '' : 'hidden'}`}
                >
                  <SectionHeader title="Overview" />
                  <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-4">
                    <Metric
                      label="Mastered"
                      value={String(report.overview.mastered)}
                      note={`of ${report.overview.total} concepts in view (${report.overview.total ? Math.round((report.overview.mastered / report.overview.total) * 100) : 0}%). Mastered only, not viewed or unlocked.`}
                    />
                    <Metric
                      label="In progress"
                      value={String(report.overview.in_progress)}
                      note="Concepts started with evidence below the mastery threshold."
                    />
                    <Metric
                      label="Accessible"
                      value={String(report.overview.unlocked)}
                      note="Unlocked by prerequisites and not yet mastered or in progress."
                    />
                    <Metric
                      label="Locked"
                      value={String(report.overview.locked)}
                      note="Prerequisites are not yet met. Locked is reachability, not a judgement of ability."
                    />
                    <Metric
                      label="XP lifetime"
                      value={String(report.overview.xp_lifetime)}
                      note="The lifetime ledger total. It does not change with the period filter."
                    />
                    <Metric
                      label="XP this period"
                      value={String(report.overview.xp_period)}
                      note={`Ledger entries from ${report.period.label.toLowerCase()}.`}
                    />
                    <Metric
                      label="Active days"
                      value={String(report.overview.activity_days)}
                      note="Days with at least one recorded answer in the period. Missing days are not counted as zero activity."
                    />
                    <Metric
                      label="Answers"
                      value={`${report.overview.correct}/${report.overview.questions}`}
                      note="Correct answers out of recorded answers in the period, across all activity."
                    />
                  </div>

                  <h3 className="font-mono text-[11px] uppercase tracking-wider text-mathua-muted mt-8 mb-2">Progress by domain</h3>
                  {(report.domains ?? []).length === 0 ? (
                    <p className="text-sm text-mathua-secondary">No domains are available for this filter.</p>
                  ) : (
                    <div className="overflow-x-auto border border-[var(--border)] rounded-none">
                      <table className="w-full text-sm">
                        <thead>
                          <tr className="text-left font-mono text-[10px] uppercase tracking-wider text-mathua-muted border-b border-[var(--border)]">
                            <th className="px-3 py-2">Domain</th>
                            <th className="px-3 py-2 text-right">Mastered</th>
                            <th className="px-3 py-2 text-right">In progress</th>
                            <th className="px-3 py-2 text-right">Accessible</th>
                            <th className="px-3 py-2 text-right">Concepts</th>
                          </tr>
                        </thead>
                        <tbody>
                          {(report.domains ?? []).map((d) => (
                            <tr key={d.id} className="border-b border-[var(--border)] last:border-0">
                              <td className="px-3 py-2">{d.id}</td>
                              <td className="px-3 py-2 text-right tabular-nums">{d.masteredCount}</td>
                              <td className="px-3 py-2 text-right tabular-nums">{d.inProgressCount}</td>
                              <td className="px-3 py-2 text-right tabular-nums">{d.unlockedCount}</td>
                              <td className="px-3 py-2 text-right tabular-nums">{d.conceptCount}</td>
                            </tr>
                          ))}
                        </tbody>
                      </table>
                    </div>
                  )}

                  <h3 className="font-mono text-[11px] uppercase tracking-wider text-mathua-muted mt-8 mb-2">Learning activity</h3>
                  {activity.length === 0 ? (
                    <p className="text-sm text-mathua-secondary">
                      No recorded activity in this period. This is an absence of records, not proof that no learning
                      happened.
                    </p>
                  ) : (
                    <figure className="border border-[var(--border)] rounded-none p-4">
                      <figcaption className="text-xs text-mathua-secondary mb-3">
                        Recorded answers per active day. Showing the most recent {activity.length} active days in the
                        period, not a continuous calendar.
                      </figcaption>
                      <div className="flex items-end gap-1 h-32" aria-hidden="true">
                        {activity.map((day) => (
                          <div
                            key={day.date}
                            className="flex-1 min-w-[3px] bg-[var(--accent-blue)]/60 rounded-none"
                            style={{ height: `${Math.max(4, (day.questions / maxQuestions) * 100)}%` }}
                            title={`${day.date}: ${day.questions} answers, ${day.correct} correct`}
                          />
                        ))}
                      </div>
                      <table className="sr-only">
                        <caption>Answers per active day</caption>
                        <thead>
                          <tr><th>Date</th><th>Answers</th><th>Correct</th></tr>
                        </thead>
                        <tbody>
                          {activity.map((day) => (
                            <tr key={day.date}><td>{day.date}</td><td>{day.questions}</td><td>{day.correct}</td></tr>
                          ))}
                        </tbody>
                      </table>
                    </figure>
                  )}

                  <div className="mt-8 border border-[var(--border)] rounded-none p-4">
                    <h3 className="font-mono text-[11px] uppercase tracking-wider text-mathua-muted">What this report can and cannot say</h3>
                    <ul className="mt-2 text-xs text-mathua-secondary space-y-1 list-disc pl-5">
                      {(report.limitations ?? []).map((item) => <li key={item}>{item}</li>)}
                    </ul>
                  </div>
                </div>

                <div
                  id="panel-mastery"
                  role="tabpanel"
                  aria-labelledby="tab-mastery"
                  data-report-panel
                  className={`pt-6 ${tab === 'mastery' ? '' : 'hidden'}`}
                >
                  <SectionHeader title="Mastery" />
                  <p className="text-sm text-mathua-secondary mb-4 max-w-3xl">
                    The state below comes from the same eligibility and mastery rules the rest of Mathua uses. A concept
                    is mastered only when its stored progress says so, not because it was opened, unlocked or assessed
                    by placement.
                  </p>
                  {mastery.length === 0 ? (
                    <p className="text-sm text-mathua-secondary">No concepts match this domain.</p>
                  ) : (
                    <>
                      <div className="hidden md:block overflow-x-auto border border-[var(--border)] rounded-none">
                        <table className="w-full text-sm">
                          <thead>
                            <tr className="text-left font-mono text-[10px] uppercase tracking-wider text-mathua-muted border-b border-[var(--border)]">
                              <th className="px-3 py-2">Concept</th>
                              <th className="px-3 py-2">Domain</th>
                              <th className="px-3 py-2">State</th>
                              <th className="px-3 py-2 text-right">Mastery</th>
                              <th className="px-3 py-2">Evidence</th>
                            </tr>
                          </thead>
                          <tbody>
                            {mastery.slice(0, visibleMastery).map((row) => (
                              <tr key={row.concept_id} className="border-b border-[var(--border)] last:border-0 align-top">
                                <td className="px-3 py-2">
                                  <Link href={row.href} className="hover:underline underline-offset-4">{row.label}</Link>
                                  <div className="text-[11px] text-mathua-muted font-mono">{row.concept_id}</div>
                                </td>
                                <td className="px-3 py-2 text-mathua-secondary">{row.domain}</td>
                                <td className="px-3 py-2">
                                  <span className={`inline-block px-2 py-0.5 rounded-none border text-[11px] ${stateTone(row)}`}>
                                    {stateLabel(row)}
                                  </span>
                                  {row.placement_seeded && (
                                    <div className="text-[11px] text-mathua-muted mt-1">Placement seeded</div>
                                  )}
                                </td>
                                <td className="px-3 py-2 text-right tabular-nums">{Math.round(row.mastery_pct * 100)}%</td>
                                <td className="px-3 py-2 text-xs text-mathua-secondary">
                                  <div>{row.mastered_at ? `Mastered ${fmtDate(row.mastered_at)}` : 'No mastery date stored'}</div>
                                  <div>{row.attempts} attempts, streak {row.streak}</div>
                                  {row.next_review && <div>Review due {fmtDate(row.next_review)}</div>}
                                  {(row.prerequisite_gaps ?? []).length > 0 && (
                                    <div className="mt-1">
                                      Gaps: {(row.prerequisite_gaps ?? []).map((g) => g.label).join(', ')}
                                    </div>
                                  )}
                                </td>
                              </tr>
                            ))}
                          </tbody>
                        </table>
                      </div>
                      <div className="md:hidden space-y-3">
                        {mastery.slice(0, visibleMastery).map((row) => (
                          <div key={row.concept_id} className="border border-[var(--border)] rounded-none p-3">
                            <div className="flex items-start justify-between gap-3">
                              <Link href={row.href} className="font-medium hover:underline underline-offset-4">{row.label}</Link>
                              <span className={`shrink-0 px-2 py-0.5 rounded-none border text-[11px] ${stateTone(row)}`}>{stateLabel(row)}</span>
                            </div>
                            <div className="text-[11px] text-mathua-muted font-mono mt-0.5">{row.domain} · {row.concept_id}</div>
                            <div className="text-xs text-mathua-secondary mt-2">
                              <div>{Math.round(row.mastery_pct * 100)}% mastery · {row.attempts} attempts · streak {row.streak}</div>
                              <div>{row.mastered_at ? `Mastered ${fmtDate(row.mastered_at)}` : 'No mastery date stored'}</div>
                              {row.next_review && <div>Review due {fmtDate(row.next_review)}</div>}
                              {(row.prerequisite_gaps ?? []).length > 0 && (
                                <div className="mt-1">Gaps: {(row.prerequisite_gaps ?? []).map((g) => g.label).join(', ')}</div>
                              )}
                            </div>
                          </div>
                        ))}
                      </div>
                      {visibleMastery < mastery.length && (
                        <button
                          type="button"
                          onClick={() => setVisibleMastery((n) => n + 100)}
                          className="mt-4 h-8 px-3 font-mono text-[11px] uppercase tracking-wider border border-[var(--border)] rounded-none report-no-print"
                        >
                          Show 100 more ({mastery.length - visibleMastery} left)
                        </button>
                      )}
                    </>
                  )}
                </div>

                <div
                  id="panel-assessments"
                  role="tabpanel"
                  aria-labelledby="tab-assessments"
                  data-report-panel
                  className={`pt-6 ${tab === 'assessments' ? '' : 'hidden'}`}
                >
                  <SectionHeader title="Assessments" />
                  <p className="text-sm text-mathua-secondary mb-4 max-w-3xl">
                    Quizzes and the diagnostic, as recorded. A quiz does not store a per-question score, so its status
                    comes from the recorded session and completion rather than an invented percentage.
                  </p>
                  {assessments.length === 0 ? (
                    <p className="text-sm text-mathua-secondary">No assessments have been recorded.</p>
                  ) : (
                    <div className="space-y-3">
                      {assessments.map((a: LearnerAssessment) => (
                        <div key={`${a.type}:${a.id}`} className="border border-[var(--border)] rounded-none p-4">
                          <div className="flex flex-wrap items-start justify-between gap-3">
                            <div className="min-w-0">
                              <div className="text-[11px] font-mono uppercase tracking-wider text-mathua-muted">{a.type === 'quiz' ? 'Mastery check' : 'Diagnostic'}</div>
                              <h3 className="font-medium">{a.title}</h3>
                            </div>
                            <span className="shrink-0 px-2 py-0.5 rounded-none border border-[var(--border)] text-[11px]">{a.status}</span>
                          </div>
                          <div className="text-xs text-mathua-secondary mt-2 space-y-1">
                            <div>
                              {a.completed_at ? `Completed ${fmtDateTime(a.completed_at)}` : a.started_at ? `Started ${fmtDateTime(a.started_at)}` : 'Time not recorded'}
                            </div>
                            <div>
                              {a.questions} answers, {a.correct} correct
                              {typeof a.accuracy === 'number' ? ` (${Math.round(a.accuracy * 100)}%)` : ''}
                            </div>
                            {(a.findings ?? []).length > 0 && (
                              <div>Concepts answered incorrectly: {(a.findings ?? []).map((f) => f.label).join(', ')}</div>
                            )}
                            {a.note && <div>{a.note}</div>}
                          </div>
                          {a.href && (
                            <Link href={a.href} className="inline-block mt-3 text-sm underline underline-offset-4">Open assessment</Link>
                          )}
                        </div>
                      ))}
                    </div>
                  )}
                </div>

                <div
                  id="panel-history"
                  role="tabpanel"
                  aria-labelledby="tab-history"
                  data-report-panel
                  className={`pt-6 ${tab === 'history' ? '' : 'hidden'}`}
                >
                  <SectionHeader title="Learning history" />
                  <p className="text-sm text-mathua-secondary mb-4 max-w-3xl">
                    Answers, mastery transitions, quiz completions and the diagnostic completion. This list is
                    paginated and newest first. It contains only events the system stored.
                  </p>
                  {history.length === 0 ? (
                    <p className="text-sm text-mathua-secondary">Nothing has been recorded yet.</p>
                  ) : (
                    <ol className="border border-[var(--border)] rounded-none divide-y divide-[var(--border)]">
                      {history.map((event) => (
                        <li key={`${event.kind}:${event.id}`} className="p-3 flex flex-wrap items-start justify-between gap-2">
                          <div className="min-w-0">
                            <div className="text-sm">
                              {event.href ? <Link href={event.href} className="hover:underline underline-offset-4">{event.title}</Link> : event.title}
                            </div>
                            <div className="text-[11px] text-mathua-muted font-mono uppercase tracking-wider">{event.kind}{event.source ? ` · ${event.source}` : ''}</div>
                            {event.detail && <div className="text-xs text-mathua-secondary mt-0.5">{event.detail}</div>}
                          </div>
                          <div className="text-right text-xs text-mathua-secondary shrink-0">
                            <div>{fmtDateTime(event.at)}</div>
                            {typeof event.xp === 'number' && event.xp !== 0 && <div>{event.xp > 0 ? '+' : ''}{event.xp} XP</div>}
                          </div>
                        </li>
                      ))}
                    </ol>
                  )}
                  {history.length < historyTotal && (
                    <button
                      type="button"
                      onClick={() => void loadMoreHistory()}
                      disabled={historyLoading}
                      className="mt-4 h-8 px-3 font-mono text-[11px] uppercase tracking-wider border border-[var(--border)] rounded-none disabled:opacity-50 report-no-print"
                    >
                      {historyLoading ? 'Loading' : `Load more (${historyTotal - history.length} left)`}
                    </button>
                  )}
                </div>
              </>
            )}
          </section>
        </div>
      </main>
      <Footer />
    </AppShell>
  )
}
