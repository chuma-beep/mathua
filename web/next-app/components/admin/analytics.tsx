'use client'

import { useCallback, useEffect, useMemo, useRef, useState } from 'react'
import { adminAnalytics, adminAnalyticsLive, getAdminNotificationHealth, getErrorMessage } from '../../lib/api'
import type {
  AdminAnalytics,
  AdminAnalyticsLive,
  AdminNotificationProblem,
  AdminNotificationStat,
} from '../../lib/api'
import { adminMetric } from '../../lib/adminMetrics'
import Loading from '../Loading'

type Tab = 'overview' | 'engagement' | 'learning' | 'operations'
type Range = '15m' | '1h' | '24h' | '7d' | '30d' | '90d'
type ActivityPoint = NonNullable<AdminAnalytics['activity']>[number]

const RANGES: { value: Range; label: string }[] = [
  { value: '15m', label: '15m' },
  { value: '1h', label: '1h' },
  { value: '24h', label: '24h' },
  { value: '7d', label: '7d' },
  { value: '30d', label: '30d' },
  { value: '90d', label: '90d' },
]

const TABS: { value: Tab; label: string; sections: string[] }[] = [
  { value: 'overview', label: 'Overview', sections: ['accounts', 'totals', 'activity', 'concepts'] },
  { value: 'engagement', label: 'Engagement', sections: ['totals', 'activity', 'signups', 'quizzes', 'diagnostics'] },
  { value: 'learning', label: 'Learning', sections: ['totals', 'mastery', 'xp', 'concepts'] },
  { value: 'operations', label: 'Operations', sections: ['totals', 'reports', 'status'] },
]

function metric(id: string) {
  return adminMetric(id)
}

function pct(correct: number, answers: number): string {
  if (answers <= 0) return '—'
  return `${((correct / answers) * 100).toFixed(1)}%`
}

function number(value: number | undefined): string {
  return (value ?? 0).toLocaleString()
}

export default function AdminAnalyticsDashboard() {
  const [range, setRange] = useState<Range>('7d')
  const [tab, setTab] = useState<Tab>('overview')
  const [data, setData] = useState<AdminAnalytics | null>(null)
  const [live, setLive] = useState<AdminAnalyticsLive | null>(null)
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(true)
  const requestId = useRef(0)

  const sections = useMemo(() => TABS.find(item => item.value === tab)?.sections ?? [], [tab])

  const load = useCallback(async () => {
    const id = ++requestId.current
    setError('')
    setLoading(true)
    try {
      const [next, nextLive] = await Promise.all([
        adminAnalytics(range, { sections }),
        adminAnalyticsLive(),
      ])
      if (id !== requestId.current) return
      setData(next)
      setLive(nextLive)
    } catch (e) {
      if (id !== requestId.current) return
      setError(getErrorMessage(e))
    } finally {
      if (id === requestId.current) setLoading(false)
    }
  }, [range, sections])

  useEffect(() => {
    void load()
  }, [load])

  // The 15m and 1h views are close to live, so they refresh at the historical
  // cache TTL while visible. Daily ranges are not polled; the page is a report,
  // not a stream.
  useEffect(() => {
    if (range !== '15m' && range !== '1h') return
    const tick = () => {
      if (document.visibilityState === 'visible') void load()
    }
    const timer = window.setInterval(tick, 30_000)
    const onVisibility = () => {
      if (document.visibilityState === 'visible') void load()
    }
    document.addEventListener('visibilitychange', onVisibility)
    return () => {
      window.clearInterval(timer)
      document.removeEventListener('visibilitychange', onVisibility)
    }
  }, [range, load])

  return (
    <section className="mt-8" aria-labelledby="analytics-heading">
      <div className="flex flex-wrap items-end justify-between gap-3 border-b border-mathua-border pb-3">
        <div>
          <h2 id="analytics-heading" className="font-serif text-lg text-mathua-primary">Analytics</h2>
          <p className="mt-1 font-mono text-[10px] text-mathua-muted">
            Event metrics describe retained rows only; reset or delete removes history.
          </p>
        </div>
        <div className="flex flex-wrap gap-1" role="group" aria-label="Analytics range">
          {RANGES.map(item => (
            <button
              key={item.value}
              type="button"
              aria-pressed={range === item.value}
              onClick={() => {
                if (item.value === range) return
                requestId.current++
                setData(null)
                setRange(item.value)
              }}
              className={`border px-2.5 h-8 font-mono text-[11px] ${
                range === item.value
                  ? 'border-mathua-blue text-mathua-blue'
                  : 'border-mathua-border text-mathua-secondary hover:border-mathua-blue'
              }`}
            >
              {item.label}
            </button>
          ))}
        </div>
      </div>

      {error && (
        <p role="alert" className="mt-4 border border-mathua-red bg-mathua-surface p-3 font-mono text-[11px] text-mathua-red">
          {error}
        </p>
      )}

      {loading && !data && <div className="mt-4"><Loading label="LOADING ANALYTICS" /></div>}

      <LiveStrip live={live} />

      <div className="mt-6 flex flex-wrap gap-2" role="tablist" aria-label="Analytics section">
        {TABS.map(item => (
          <button
            key={item.value}
            type="button"
            role="tab"
            aria-selected={tab === item.value}
            onClick={() => {
              if (item.value === tab) return
              requestId.current++
              setData(null)
              setTab(item.value)
            }}
            className={`border px-3 h-9 font-mono text-[11px] ${
              tab === item.value
                ? 'border-mathua-blue text-mathua-blue'
                : 'border-mathua-border text-mathua-secondary hover:border-mathua-blue'
            }`}
          >
            {item.label}
          </button>
        ))}
      </div>

      {data && tab === 'overview' && <OverviewTab data={data} range={range} />}
      {data && tab === 'engagement' && <EngagementTab data={data} range={range} />}
      {data && tab === 'learning' && <LearningTab data={data} range={range} />}
      {data && tab === 'operations' && <OperationsTab data={data} range={range} />}
    </section>
  )
}

function LiveStrip({ live }: { live: AdminAnalyticsLive | null }) {
  return (
    <dl className="mt-4 grid grid-cols-2 sm:grid-cols-5 gap-2">
      <Kpi label="Online now" value={number(live?.online_now)} hint="Signed in, five minutes" />
      <Kpi label="Practicing now (15m)" value={number(live?.practicing_now)} hint="Practice activity" />
      <Kpi label="Guests practicing" value={number(live?.guest_practice_now)} hint="Guest practice activity" />
      <Kpi label="Answers, 15m" value={number(live?.answers_15m)} hint="Submitted answers" />
      <Kpi label="Active sessions" value={number(live?.active_sessions)} hint="Touched in fifteen minutes" />
    </dl>
  )
}

function Kpi({ label, value, hint }: { label: string; value: string; hint?: string }) {
  return (
    <div className="border border-mathua-border bg-mathua-surface p-3 min-w-0">
      <dt className="font-mono text-[10px] uppercase tracking-wider text-mathua-muted">{label}</dt>
      <dd className="mt-1 font-mono text-xl text-mathua-primary tabular-nums">{value}</dd>
      {hint && <p className="mt-1 font-mono text-[9px] text-mathua-muted">{hint}</p>}
    </div>
  )
}

function PrimaryKpis({ data }: { data: AdminAnalytics }) {
  const totals = data.totals
  return (
    <dl className="mt-6 grid grid-cols-2 lg:grid-cols-4 gap-2">
      <Kpi label={metric('unique_learners')?.label ?? 'Unique learners'} value={number(totals?.unique_learners)} hint={metric('unique_learners')?.definition} />
      <Kpi label={metric('answers')?.label ?? 'Answers'} value={number(totals?.answers)} hint={metric('answers')?.definition} />
      <Kpi label={metric('accuracy')?.label ?? 'Accuracy'} value={pct(totals?.correct ?? 0, totals?.answers ?? 0)} hint={`${number(totals?.correct)} correct of ${number(totals?.answers)}`} />
      <Kpi label={metric('new_registrations')?.label ?? 'New registrations'} value={number(totals?.new_registrations)} hint={metric('new_registrations')?.definition} />
    </dl>
  )
}

function OverviewTab({ data, range }: { data: AdminAnalytics; range: Range }) {
  return (
    <>
      <PrimaryKpis data={data} />
      <div className="mt-6 grid items-start gap-4 lg:grid-cols-[1.4fr_1fr]">
        <ChartCard title="Learner activity" hint="Distinct learners per bucket, split by account class.">
          <ActivityBars points={data.activity ?? []} range={range} />
        </ChartCard>
        <ChartCard title="Top concepts" hint="Attempts in the selected range.">
          <ConceptTable rows={data.top_concepts ?? []} />
        </ChartCard>
      </div>
    </>
  )
}

function EngagementTab({ data, range }: { data: AdminAnalytics; range: Range }) {
  const totals = data.totals
  return (
    <>
      <PrimaryKpis data={data} />
      <div className="mt-6 grid items-start gap-4 lg:grid-cols-2">
        <ChartCard title="New registrations" hint="Registered accounts created in the range.">
          <TimeSeries points={data.signups ?? []} range={range} color="var(--accent-teal)" label="New registrations" />
        </ChartCard>
        <ChartCard title="Quiz completions" hint="Mastery checks and timed quizzes.">
          <DualSeries
            primary={data.quizzes ?? []}
            secondary={data.timed_quizzes ?? []}
            primaryLabel="Mastery checks"
            secondaryLabel="Timed quizzes"
            range={range}
          />
        </ChartCard>
        <ChartCard title="Diagnostic cohort" hint="Registered in the range, completed to date.">
          <div className="grid grid-cols-2 gap-2">
            <Kpi label="Registered" value={number(totals?.diagnostic_cohort_registered)} />
            <Kpi label="Completed" value={number(totals?.diagnostic_cohort_completed)} />
          </div>
          <p className="mt-3 font-mono text-[10px] text-mathua-muted">
            A cohort, not an event rate. Completion can happen after the registration range closes.
          </p>
        </ChartCard>
        <ChartCard title="Answers" hint="Answers submitted per bucket.">
          <TimeSeries points={activityAnswers(data.activity ?? [])} range={range} color="var(--accent-blue)" label="Answers" />
        </ChartCard>
      </div>
    </>
  )
}

function LearningTab({ data, range }: { data: AdminAnalytics; range: Range }) {
  const domainRows = Object.entries(data.mastered_by_concept ?? {})
    .sort((a, b) => b[1] - a[1])
    .slice(0, 12)
  return (
    <>
      <PrimaryKpis data={data} />
      <div className="mt-6 grid items-start gap-4 lg:grid-cols-2">
        <ChartCard title="Mastery" hint="Concepts mastered by the date recorded.">
          <TimeSeries points={data.mastery ?? []} range={range} color="var(--accent-green)" label="Mastery" />
        </ChartCard>
        <ChartCard title="XP" hint="Awarded and deducted ledger events.">
          <DualXp points={data.xp ?? []} range={range} />
        </ChartCard>
        <ChartCard title="Currently mastered" hint="Current state, not a range total.">
          <p className="font-mono text-3xl text-mathua-primary tabular-nums">{number(data.mastered_current)}</p>
          <p className="mt-2 font-mono text-[10px] text-mathua-muted">Re-mastery overwrites the date.</p>
        </ChartCard>
        <ChartCard title="Top mastered concepts" hint="Concept ids with current mastery counts.">
          {domainRows.length === 0 ? <Empty /> : (
            <ul className="space-y-1">
              {domainRows.map(([id, count]) => (
                <li key={id} className="flex justify-between gap-3 font-mono text-[11px]">
                  <span className="truncate text-mathua-secondary">{id}</span>
                  <span className="tabular-nums text-mathua-primary">{count}</span>
                </li>
              ))}
            </ul>
          )}
        </ChartCard>
      </div>
    </>
  )
}

function OperationsTab({ data, range }: { data: AdminAnalytics; range: Range }) {
  return (
    <>
      <PrimaryKpis data={data} />
      <div className="mt-6 grid items-start gap-4 lg:grid-cols-2">
        <ChartCard title="Question reports" hint="Opened and resolved in the range.">
          <DualReports points={data.reports ?? []} range={range} />
        </ChartCard>
        <ChartCard title="Report status" hint="Current state.">
          <SegmentedBar values={data.report_status ?? {}} />
        </ChartCard>
        <ChartCard title="Most reported concepts" hint="Reports in the range.">
          <ConceptTable rows={data.most_reported ?? []} />
        </ChartCard>
        <NotificationHealthCard />
      </div>
    </>
  )
}

function ChartCard({ title, hint, children }: { title: string; hint?: string; children: React.ReactNode }) {
  return (
    <section className="border border-mathua-border bg-mathua-surface p-4 min-w-0">
      <h3 className="font-mono text-[11px] uppercase tracking-wider text-mathua-secondary">{title}</h3>
      {hint && <p className="mt-1 font-mono text-[10px] text-mathua-muted">{hint}</p>}
      <div className="mt-3">{children}</div>
    </section>
  )
}

function Empty() {
  return <p className="font-mono text-[11px] text-mathua-muted">No rows in this window.</p>
}

function ActivityBars({ points, range }: { points: ActivityPoint[]; range: Range }) {
  if (points.length === 0) return <Empty />
  const max = Math.max(1, ...points.map(p => p.registered + p.guests + p.unclassified))
  const width = 600
  const height = 150
  const bar = width / points.length
  return (
    <div>
      <svg viewBox={`0 0 ${width} ${height}`} className="h-auto w-full" role="img" aria-label="Learner activity by account class">
        <text x="0" y="10" fill="var(--text-muted)" fontSize="9" fontFamily="monospace">{max}</text>
        {[0, 0.5, 1].map(t => (
          <line key={t} x1={0} x2={width} y1={height - t * height} y2={height - t * height} stroke="var(--border)" strokeWidth="1" />
        ))}
        {points.map((p, i) => {
          const total = p.registered + p.guests + p.unclassified
          const h = (total / max) * height
          const guestH = (p.guests / Math.max(1, total)) * h
          const unclassifiedH = (p.unclassified / Math.max(1, total)) * h
          const registeredH = h - guestH - unclassifiedH
          const x = i * bar + bar * 0.15
          const w = Math.max(1, bar * 0.7)
          return (
            <g key={p.start}>
              <rect x={x} y={height - registeredH} width={w} height={registeredH} fill="var(--accent-blue)" opacity="0.85" />
              <rect x={x} y={height - registeredH - guestH} width={w} height={guestH} fill="var(--accent-teal)" opacity="0.75" />
              <rect x={x} y={height - h} width={w} height={unclassifiedH} fill="var(--accent-red)" opacity="0.55" />
            </g>
          )
        })}
      </svg>
      <Legend items={[
        ['Registered', 'var(--accent-blue)'],
        ['Guests', 'var(--accent-teal)'],
        ['Unclassified', 'var(--accent-red)'],
      ]} />
      <p className="mt-1 font-mono text-[9px] text-mathua-muted">{rangeLabel(points[0].start, range)} to {rangeLabel(points[points.length - 1].start, range)}</p>
    </div>
  )
}

function TimeSeries({ points, range, color, label }: { points: { start: number; value: number }[]; range: Range; color: string; label: string }) {
  if (points.length === 0) return <Empty />
  const width = 600
  const height = 140
  const max = Math.max(1, ...points.map(p => p.value))
  const step = points.length > 1 ? width / (points.length - 1) : width
  const coords = points.map((p, i) => [i * step, height - (p.value / max) * (height - 12) - 2] as const)
  const path = coords.map(([x, y], i) => `${i === 0 ? 'M' : 'L'}${x.toFixed(1)},${y.toFixed(1)}`).join(' ')
  const area = `${path} L${width},${height} L0,${height} Z`
  return (
    <div>
      <svg viewBox={`0 0 ${width} ${height}`} className="h-auto w-full" role="img" aria-label={label}>
        <text x="0" y="10" fill="var(--text-muted)" fontSize="9" fontFamily="monospace">{max}</text>
        {[0, 0.5, 1].map(t => (
          <line key={t} x1={0} x2={width} y1={height - t * (height - 12) - 2} y2={height - t * (height - 12) - 2} stroke="var(--border)" strokeWidth="1" />
        ))}
        <path d={area} fill={color} opacity="0.12" />
        <path d={path} fill="none" stroke={color} strokeWidth="2" vectorEffect="non-scaling-stroke" />
        <circle cx={coords[coords.length - 1][0]} cy={coords[coords.length - 1][1]} r="3" fill={color} />
      </svg>
      <p className="mt-1 font-mono text-[9px] text-mathua-muted">{rangeLabel(points[0].start, range)} to {rangeLabel(points[points.length - 1].start, range)}</p>
    </div>
  )
}

function DualSeries({ primary, secondary, primaryLabel, secondaryLabel, range }: {
  primary: { start: number; value: number }[]
  secondary: { start: number; value: number }[]
  primaryLabel: string
  secondaryLabel: string
  range: Range
}) {
  if (primary.length === 0 && secondary.length === 0) return <Empty />
  const width = 600
  const height = 140
  const max = Math.max(1, ...primary.map(p => p.value), ...secondary.map(p => p.value))
  const draw = (points: { start: number; value: number }[], color: string) => {
    if (points.length === 0) return null
    const step = points.length > 1 ? width / (points.length - 1) : width
    const path = points.map((p, i) => `${i === 0 ? 'M' : 'L'}${(i * step).toFixed(1)},${(height - (p.value / max) * (height - 12) - 2).toFixed(1)}`).join(' ')
    return <path d={path} fill="none" stroke={color} strokeWidth="2" vectorEffect="non-scaling-stroke" />
  }
  return (
    <div>
      <svg viewBox={`0 0 ${width} ${height}`} className="h-auto w-full" role="img" aria-label={`${primaryLabel} and ${secondaryLabel}`}>
        <text x="0" y="10" fill="var(--text-muted)" fontSize="9" fontFamily="monospace">{max}</text>
        {[0, 0.5, 1].map(t => (
          <line key={t} x1={0} x2={width} y1={height - t * (height - 12) - 2} y2={height - t * (height - 12) - 2} stroke="var(--border)" strokeWidth="1" />
        ))}
        {draw(primary, 'var(--accent-blue)')}
        {draw(secondary, 'var(--accent-green)')}
      </svg>
      <Legend items={[[primaryLabel, 'var(--accent-blue)'], [secondaryLabel, 'var(--accent-green)']]} />
      <p className="mt-1 font-mono text-[9px] text-mathua-muted">{rangeLabel(primary[0]?.start ?? 0, range)}</p>
    </div>
  )
}

function DualXp({ points, range }: { points: { start: number; awarded: number; deducted: number }[]; range: Range }) {
  if (points.length === 0) return <Empty />
  const width = 600
  const height = 140
  const max = Math.max(1, ...points.map(p => Math.max(p.awarded, p.deducted)))
  const step = points.length > 1 ? width / (points.length - 1) : width
  return (
    <div>
      <svg viewBox={`0 0 ${width} ${height}`} className="h-auto w-full" role="img" aria-label="XP awarded and deducted">
        <text x="0" y="10" fill="var(--text-muted)" fontSize="9" fontFamily="monospace">{max}</text>
        {[0, 0.5, 1].map(t => (
          <line key={t} x1={0} x2={width} y1={height - t * (height - 12) - 2} y2={height - t * (height - 12) - 2} stroke="var(--border)" strokeWidth="1" />
        ))}
        {points.map((p, i) => {
          const x = i * step
          const ah = (p.awarded / max) * (height - 12)
          const dh = (p.deducted / max) * (height - 12)
          return (
            <g key={p.start}>
              <rect x={x - 2} y={height - ah - 2} width="4" height={ah} fill="var(--accent-green)" opacity="0.8" />
              <rect x={x + 1} y={height - dh - 2} width="4" height={dh} fill="var(--accent-red)" opacity="0.65" />
            </g>
          )
        })}
      </svg>
      <Legend items={[['Awarded', 'var(--accent-green)'], ['Deducted', 'var(--accent-red)']]} />
      <p className="mt-1 font-mono text-[9px] text-mathua-muted">{rangeLabel(points[0].start, range)}</p>
    </div>
  )
}

function DualReports({ points, range }: { points: { start: number; opened: number; resolved: number }[]; range: Range }) {
  return (
    <DualSeries
      primary={points.map(p => ({ start: p.start, value: p.opened }))}
      secondary={points.map(p => ({ start: p.start, value: p.resolved }))}
      primaryLabel="Opened"
      secondaryLabel="Resolved"
      range={range}
    />
  )
}

function SegmentedBar({ values }: { values: Record<string, number> }) {
  const entries = Object.entries(values).filter(([, value]) => value > 0)
  const total = entries.reduce((sum, [, value]) => sum + value, 0)
  if (total === 0) return <Empty />
  const colors: Record<string, string> = {
    open: 'var(--accent-red)',
    reviewing: 'var(--accent-teal)',
    resolved: 'var(--accent-green)',
    dismissed: 'var(--text-muted)',
  }
  return (
    <div>
      <div className="flex h-3 w-full overflow-hidden border border-mathua-border">
        {entries.map(([key, value]) => (
          <div key={key} style={{ width: `${(value / total) * 100}%`, background: colors[key] ?? 'var(--accent-blue)' }} title={`${key}: ${value}`} />
        ))}
      </div>
      <ul className="mt-2 flex flex-wrap gap-x-4 gap-y-1">
        {entries.map(([key, value]) => (
          <li key={key} className="font-mono text-[10px] text-mathua-secondary">
            {key} <span className="tabular-nums text-mathua-primary">{value}</span>
          </li>
        ))}
      </ul>
    </div>
  )
}

function ConceptTable({ rows }: { rows: { concept_id: string; attempts: number; correct: number; reports: number }[] }) {
  if (rows.length === 0) return <Empty />
  return (
    <table className="w-full table-fixed border-collapse font-mono text-[11px]">
      <thead>
        <tr className="text-left text-mathua-muted">
          <th className="w-[58%] py-1 font-normal">Concept</th>
          <th className="w-[22%] py-1 text-right font-normal">Attempts</th>
          <th className="w-[20%] py-1 text-right font-normal">Correct</th>
        </tr>
      </thead>
      <tbody>
        {rows.slice(0, 10).map(row => (
          <tr key={row.concept_id} className="border-t border-mathua-border">
            <td className="overflow-hidden text-ellipsis whitespace-nowrap py-1 pr-2 text-mathua-secondary" title={row.concept_id}>{row.concept_id}</td>
            <td className="py-1 text-right tabular-nums text-mathua-primary">{row.attempts}</td>
            <td className="py-1 text-right tabular-nums text-mathua-secondary">{pct(row.correct, row.attempts)}</td>
          </tr>
        ))}
      </tbody>
    </table>
  )
}

function NotificationHealthCard() {
  const [stats, setStats] = useState<AdminNotificationStat[]>([])
  const [problems, setProblems] = useState<AdminNotificationProblem[]>([])
  const [error, setError] = useState('')

  useEffect(() => {
    let cancelled = false
    getAdminNotificationHealth(10)
      .then(result => {
        if (cancelled) return
        setStats(result.stats)
        setProblems(result.problems)
      })
      .catch(e => { if (!cancelled) setError(getErrorMessage(e)) })
    return () => { cancelled = true }
  }, [])

  return (
    <ChartCard title="Notification delivery" hint="Current delivery state, all time.">
      {error ? <p className="font-mono text-[11px] text-mathua-red">{error}</p> : (
        <>
          <div className="grid grid-cols-2 gap-2">
            {stats.slice(0, 4).map(stat => (
              <Kpi key={`${stat.channel}-${stat.status}`} label={`${stat.channel} ${stat.status}`} value={number(stat.count)} />
            ))}
          </div>
          <p className="mt-3 font-mono text-[10px] text-mathua-muted">
            {problems.length === 0 ? 'No failed or suppressed deliveries.' : `${problems.length} failed or suppressed deliveries.`}
          </p>
        </>
      )}
    </ChartCard>
  )
}

function Legend({ items }: { items: [string, string][] }) {
  return (
    <ul className="mt-2 flex flex-wrap gap-x-4 gap-y-1">
      {items.map(([label, color]) => (
        <li key={label} className="flex items-center gap-1.5 font-mono text-[9px] text-mathua-muted">
          <span className="inline-block h-2 w-2" style={{ background: color }} />
          {label}
        </li>
      ))}
    </ul>
  )
}

function activityAnswers(points: ActivityPoint[]) {
  return points.map(p => ({ start: p.start, value: p.answers }))
}

function rangeLabel(start: number, range: Range): string {
  const date = new Date(start * 1000)
  if (range === '15m' || range === '1h' || range === '24h') {
    return date.toLocaleTimeString(undefined, { hour: '2-digit', minute: '2-digit' })
  }
  return date.toLocaleDateString(undefined, { month: 'short', day: 'numeric' })
}
