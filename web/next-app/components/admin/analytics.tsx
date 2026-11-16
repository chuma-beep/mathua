'use client'

import { useCallback, useEffect, useMemo, useRef, useState } from 'react'
import {
  Area,
  AreaChart,
  CartesianGrid,
  ComposedChart,
  Line,
  LineChart,
  XAxis,
  YAxis,
} from 'recharts'
import { adminAnalytics, adminAnalyticsLive, getAdminNotificationHealth, getErrorMessage } from '../../lib/api'
import type {
  AdminAnalytics,
  AdminAnalyticsLive,
  AdminNotificationProblem,
  AdminNotificationStat,
} from '../../lib/api'
import { adminMetric } from '../../lib/adminMetrics'
import { Skeleton } from '../ui/skeleton'
import {
  ChartContainer,
  ChartLegend,
  ChartLegendContent,
  ChartTooltip,
  ChartTooltipContent,
  type ChartConfig,
} from '../ui/chart'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '../ui/tabs'
import { ToggleGroup, ToggleGroupItem } from '../ui/toggle-group'

type Tab = 'overview' | 'engagement' | 'learning' | 'operations'
type Range = '15m' | '1h' | '24h' | '7d' | '30d' | '90d'

const RANGES: { value: Range; label: string }[] = [
  { value: '15m', label: '15m' },
  { value: '1h', label: '1h' },
  { value: '24h', label: '24h' },
  { value: '7d', label: '7d' },
  { value: '30d', label: '30d' },
  { value: '90d', label: '90d' },
]

const TABS: { value: Tab; label: string; sections: string[] }[] = [
  { value: 'overview', label: 'Overview', sections: ['accounts', 'totals', 'activity', 'concepts', 'signups', 'mastery'] },
  { value: 'engagement', label: 'Engagement', sections: ['totals', 'activity', 'signups', 'quizzes', 'diagnostics'] },
  { value: 'learning', label: 'Learning', sections: ['totals', 'activity', 'mastery', 'xp', 'concepts'] },
  { value: 'operations', label: 'Operations', sections: ['totals', 'reports', 'status'] },
]

const GROWTH_CONFIG = {
  cumulativeRegistered: { label: 'Registered accounts', color: 'var(--chart-1)' },
  signups: { label: 'New registrations', color: 'var(--chart-2)' },
} satisfies ChartConfig

const ACTIVITY_CONFIG = {
  activeRegistered: { label: 'Registered', color: 'var(--chart-1)' },
  activeGuests: { label: 'Guests', color: 'var(--chart-2)' },
  activeUnclassified: { label: 'Unclassified', color: 'var(--chart-4)' },
  answers: { label: 'Answers', color: 'var(--chart-3)' },
} satisfies ChartConfig

const ACCURACY_CONFIG = {
  accuracy: { label: 'Accuracy', color: 'var(--chart-3)' },
} satisfies ChartConfig

const MASTERY_CONFIG = {
  mastery: { label: 'Mastery events', color: 'var(--chart-1)' },
} satisfies ChartConfig

const REGISTRATION_CONFIG = {
  signups: { label: 'New registrations', color: 'var(--chart-2)' },
} satisfies ChartConfig

const QUIZ_CONFIG = {
  quizzes: { label: 'Mastery checks', color: 'var(--chart-1)' },
  timedQuizzes: { label: 'Timed quizzes', color: 'var(--chart-3)' },
} satisfies ChartConfig

const XP_CONFIG = {
  xpAwarded: { label: 'Awarded', color: 'var(--chart-3)' },
  xpDeducted: { label: 'Deducted', color: 'var(--chart-4)' },
} satisfies ChartConfig

const REPORT_CONFIG = {
  reportsOpened: { label: 'Opened', color: 'var(--chart-4)' },
  reportsResolved: { label: 'Resolved', color: 'var(--chart-3)' },
} satisfies ChartConfig

const INTENSITY_CONFIG = {
  answersPerLearner: { label: 'Answers per active learner', color: 'var(--chart-2)' },
} satisfies ChartConfig

type Row = {
  start: number
  label: string
  signups: number
  cumulativeRegistered: number
  activeRegistered: number
  activeGuests: number
  activeUnclassified: number
  answers: number
  correct: number
  accuracy: number
  mastery: number
  xpAwarded: number
  xpDeducted: number
  quizzes: number
  timedQuizzes: number
  reportsOpened: number
  reportsResolved: number
  answersPerLearner: number
}

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

function bucketRows(data: AdminAnalytics, range: Range): Row[] {
  const map = new Map<number, Row>()
  const row = (start: number): Row => {
    const existing = map.get(start)
    if (existing) return existing
    const next: Row = {
      start,
      label: rangeLabel(start, range),
      signups: 0,
      cumulativeRegistered: 0,
      activeRegistered: 0,
      activeGuests: 0,
      activeUnclassified: 0,
      answers: 0,
      correct: 0,
      accuracy: 0,
      mastery: 0,
      xpAwarded: 0,
      xpDeducted: 0,
      quizzes: 0,
      timedQuizzes: 0,
      reportsOpened: 0,
      reportsResolved: 0,
      answersPerLearner: 0,
    }
    map.set(start, next)
    return next
  }

  for (const item of data.activity ?? []) {
    const target = row(item.start)
    target.activeRegistered = item.registered
    target.activeGuests = item.guests
    target.activeUnclassified = item.unclassified
    target.answers = item.answers
    target.correct = item.correct
    target.accuracy = item.answers > 0 ? (item.correct / item.answers) * 100 : 0
    const learners = item.registered + item.guests + item.unclassified
    target.answersPerLearner = learners > 0 ? item.answers / learners : 0
  }
  for (const item of data.signups ?? []) row(item.start).signups = item.value
  for (const item of data.mastery ?? []) row(item.start).mastery = item.value
  for (const item of data.xp ?? []) {
    const target = row(item.start)
    target.xpAwarded = item.awarded
    target.xpDeducted = item.deducted
  }
  for (const item of data.quizzes ?? []) row(item.start).quizzes = item.value
  for (const item of data.timed_quizzes ?? []) row(item.start).timedQuizzes = item.value
  for (const item of data.reports ?? []) {
    const target = row(item.start)
    target.reportsOpened = item.opened
    target.reportsResolved = item.resolved
  }

  const rows = [...map.values()].sort((a, b) => a.start - b.start)
  // Registered accounts are a stock, not a flow. The current total minus the
  // registrations created in this window gives the opening balance, and each
  // bucket's signups roll it forward. A deletion makes the line an estimate,
  // which the card hint says.
  if (data.accounts?.registered && data.totals) {
    let running = data.accounts.registered - data.totals.new_registrations
    for (const item of rows) {
      running += item.signups
      item.cumulativeRegistered = running
    }
  }
  return rows
}

const RANGE_VALUES = new Set(RANGES.map(item => item.value))
const TAB_VALUES = new Set(TABS.map(item => item.value))

function initialRange(): Range {
  if (typeof window === 'undefined') return '7d'
  const value = new URLSearchParams(window.location.search).get('range')
  return value && RANGE_VALUES.has(value as Range) ? (value as Range) : '7d'
}

function initialTab(): Tab {
  if (typeof window === 'undefined') return 'overview'
  const value = new URLSearchParams(window.location.search).get('tab')
  return value && TAB_VALUES.has(value as Tab) ? (value as Tab) : 'overview'
}

// The range and tab are part of the page URL so an operator can share exactly
// what they are looking at. replaceState keeps the back button on the page
// rather than stepping through every control change.
function syncAnalyticsUrl(range: Range, tab: Tab) {
  if (typeof window === 'undefined') return
  const url = new URL(window.location.href)
  if (range === '7d') url.searchParams.delete('range')
  else url.searchParams.set('range', range)
  if (tab === 'overview') url.searchParams.delete('tab')
  else url.searchParams.set('tab', tab)
  window.history.replaceState({}, '', url.toString())
}

export default function AdminAnalyticsDashboard() {
  const [range, setRange] = useState<Range>(initialRange)
  const [tab, setTab] = useState<Tab>(initialTab)
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
    <section className="mt-10" aria-labelledby="analytics-heading">
      <div className="border-b border-mathua-border pb-4">
        <div className="flex flex-wrap items-start justify-between gap-3">
          <div>
            <p className="font-mono text-[10px] uppercase tracking-[0.16em] text-mathua-muted">Analytics</p>
            <h2 id="analytics-heading" className="mt-1 font-serif text-2xl text-mathua-primary">Growth and learning</h2>
            <p className="mt-1 max-w-2xl font-mono text-[10px] leading-5 text-mathua-muted">
              User growth, practice and learning outcomes for the selected window. Event metrics describe retained rows only. All times are UTC.
            </p>
          </div>
          <div className="flex items-center gap-2">
            <span className="font-mono text-[10px] text-mathua-muted">
              {data ? `Updated ${updatedLabel(data.generated_at)}` : 'Loading window'}
            </span>
            <button
              type="button"
              onClick={() => void load()}
              className="h-8 border border-mathua-border px-2.5 font-mono text-[11px] text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue"
            >
              Refresh
            </button>
          </div>
        </div>
        <ToggleGroup
          type="single"
          value={range}
          onValueChange={(value) => {
            if (!value || value === range) return
            requestId.current++
            setRange(value as Range)
            syncAnalyticsUrl(value as Range, tab)
          }}
          aria-label="Analytics range"
          className="mt-3 flex justify-start gap-1 overflow-x-auto"
        >
          {RANGES.map(item => (
            <ToggleGroupItem
              key={item.value}
              value={item.value}
              aria-label={item.label}
              className="h-8 shrink-0 rounded-none border border-mathua-border px-2.5 font-mono text-[11px] text-mathua-secondary data-[state=on]:border-mathua-blue data-[state=on]:bg-mathua-blue-faint data-[state=on]:text-mathua-blue"
            >
              {item.label}
            </ToggleGroupItem>
          ))}
        </ToggleGroup>
      </div>

      {error && (
        <div role="alert" className="mt-4 flex flex-wrap items-center justify-between gap-3 border border-mathua-red bg-mathua-surface p-3">
          <p className="font-mono text-[11px] text-mathua-red">{error}</p>
          <button
            type="button"
            onClick={() => void load()}
            className="h-8 border border-mathua-red px-2.5 font-mono text-[11px] text-mathua-red hover:bg-mathua-red-faint"
          >
            Retry
          </button>
        </div>
      )}

      {loading && !data && <AnalyticsSkeleton />}

      <LiveStrip live={live} />

      <Tabs
        value={tab}
        onValueChange={(value) => {
          if (value === tab) return
          requestId.current++
          setTab(value as Tab)
          syncAnalyticsUrl(range, value as Tab)
        }}
        className="mt-6"
      >
        <TabsList className="flex h-auto w-full justify-start gap-5 overflow-x-auto rounded-none border-b border-mathua-border bg-transparent p-0">
          {TABS.map(item => (
            <TabsTrigger
              key={item.value}
              value={item.value}
              className="-mb-px h-10 shrink-0 rounded-none border-b-2 border-transparent bg-transparent px-0 font-mono text-[11px] text-mathua-secondary shadow-none data-[state=active]:border-mathua-blue data-[state=active]:bg-transparent data-[state=active]:text-mathua-blue data-[state=active]:shadow-none"
            >
              {item.label}
            </TabsTrigger>
          ))}
        </TabsList>

        <TabsContent value="overview">{data && <OverviewTab data={data} range={range} />}</TabsContent>
        <TabsContent value="engagement">{data && <EngagementTab data={data} range={range} />}</TabsContent>
        <TabsContent value="learning">{data && <LearningTab data={data} range={range} />}</TabsContent>
        <TabsContent value="operations">{data && <OperationsTab data={data} range={range} />}</TabsContent>
      </Tabs>
    </section>
  )
}

function AnalyticsSkeleton() {
  return (
    <div className="mt-4 space-y-4" aria-busy="true">
      <span className="sr-only">Loading analytics</span>
      <div className="grid grid-cols-2 gap-2 sm:grid-cols-5">
        {Array.from({ length: 5 }, (_, index) => (
          <Skeleton key={index} className="h-24 w-full rounded-none" />
        ))}
      </div>
      <div className="grid gap-4 lg:grid-cols-3">
        <Skeleton className="h-[290px] w-full rounded-none lg:col-span-2" />
        <Skeleton className="h-[290px] w-full rounded-none" />
      </div>
    </div>
  )
}

function LiveStrip({ live }: { live: AdminAnalyticsLive | null }) {
  return (
    <dl className="mt-4 grid grid-cols-2 gap-2 sm:grid-cols-5">
      <Kpi label="Online now" value={number(live?.online_now)} hint="Signed in, five minutes" />
      <Kpi label="Practicing now (15m)" value={number(live?.practicing_now)} hint="Practice activity" />
      <Kpi label="Guests practicing" value={number(live?.guest_practice_now)} hint="Guest practice activity" />
      <Kpi label="Answers, 15m" value={number(live?.answers_15m)} hint="Submitted answers" />
      <Kpi label="Active sessions" value={number(live?.active_sessions)} hint="Touched in fifteen minutes" />
    </dl>
  )
}

function Kpi({ label, value, hint, delta }: { label: string; value: string; hint?: string; delta?: string }) {
  return (
    <div className="border border-mathua-border bg-mathua-surface p-3 min-w-0">
      <dt className="font-mono text-[10px] uppercase tracking-wider text-mathua-muted">{label}</dt>
      <dd className="mt-1 flex flex-wrap items-baseline gap-2">
        <span className="font-mono text-2xl text-mathua-primary tabular-nums">{value}</span>
        {delta && <span className="font-mono text-[10px] text-mathua-muted">{delta}</span>}
      </dd>
      {hint && <p className="mt-1 font-mono text-[9px] leading-4 text-mathua-muted">{hint}</p>}
    </div>
  )
}

function PrimaryKpis({ data }: { data: AdminAnalytics }) {
  const totals = data.totals
  const previous = data.previous
  return (
    <dl className="mt-5 grid grid-cols-2 gap-2 lg:grid-cols-4">
      <Kpi
        label={metric('unique_learners')?.label ?? 'Unique learners'}
        value={number(totals?.unique_learners)}
        hint={metric('unique_learners')?.definition}
        delta={deltaLabel(totals?.unique_learners, previous?.unique_learners)}
      />
      <Kpi
        label={metric('answers')?.label ?? 'Answers'}
        value={number(totals?.answers)}
        hint={metric('answers')?.definition}
        delta={deltaLabel(totals?.answers, previous?.answers)}
      />
      <Kpi
        label={metric('accuracy')?.label ?? 'Accuracy'}
        value={pct(totals?.correct ?? 0, totals?.answers ?? 0)}
        hint={`${number(totals?.correct)} correct of ${number(totals?.answers)}`}
        delta={ppDelta(totals?.correct, totals?.answers, previous?.correct, previous?.answers)}
      />
      <Kpi
        label={metric('new_registrations')?.label ?? 'New registrations'}
        value={number(totals?.new_registrations)}
        hint={metric('new_registrations')?.definition}
        delta={deltaLabel(totals?.new_registrations, previous?.new_registrations)}
      />
    </dl>
  )
}

function OverviewTab({ data, range }: { data: AdminAnalytics; range: Range }) {
  const rows = bucketRows(data, range)
  return (
    <>
      <PrimaryKpis data={data} />
      <div className="mt-5 grid items-start gap-4 lg:grid-cols-3">
        <ChartCard className="lg:col-span-2" title="User growth and participation" hint="Cumulative registered accounts and new registrations. The opening balance is derived from the current total, so a deletion makes it an estimate.">
          <ComposedTrend rows={rows} config={GROWTH_CONFIG} leftDomain={['dataMin - 20', 'dataMax + 20']}>
            <Line yAxisId="left" type="monotone" dataKey="cumulativeRegistered" stroke="var(--color-cumulativeRegistered)" strokeWidth={2} dot={false} />
            <Line yAxisId="right" type="monotone" dataKey="signups" stroke="var(--color-signups)" strokeWidth={2} dot={false} />
          </ComposedTrend>
        </ChartCard>
        <ChartCard title="Account mix" hint="Current classification of every student row.">
          <AccountMix data={data} />
        </ChartCard>
      </div>
      <div className="mt-4 grid items-start gap-4 lg:grid-cols-3">
        <ChartCard title="Accuracy" hint="Correct answers divided by answers in each bucket.">
          <AccuracyTrend rows={rows} />
        </ChartCard>
        <ChartCard title="Mastery events" hint="Concepts recorded as mastered in each bucket.">
          <MasteryTrend rows={rows} />
        </ChartCard>
        <ChartCard title="Practice intensity" hint="Answers divided by active learners in each bucket.">
          <IntensityTrend rows={rows} />
        </ChartCard>
      </div>
      <ChartCard className="mt-4" title="Top concepts" hint="Attempts in the selected range.">
        <ConceptTable rows={data.top_concepts ?? []} />
      </ChartCard>
    </>
  )
}

function EngagementTab({ data, range }: { data: AdminAnalytics; range: Range }) {
  const totals = data.totals
  const rows = bucketRows(data, range)
  return (
    <>
      <PrimaryKpis data={data} />
      <div className="mt-5 grid items-start gap-4 lg:grid-cols-2">
        <ChartCard title="Learner activity" hint="Distinct learners per bucket, split by account class.">
          <ActivityTrend rows={rows} />
        </ChartCard>
        <ChartCard title="New registrations" hint="Registered accounts created in the range.">
          <RegistrationTrend rows={rows} />
        </ChartCard>
        <ChartCard title="Quiz completions" hint="Mastery checks and timed quizzes.">
          <QuizTrend rows={rows} />
        </ChartCard>
        <ChartCard title="Diagnostic cohort" hint="Registered in the range, completed to date.">
          <div className="grid grid-cols-2 gap-2">
            <Kpi label="Registered" value={number(totals?.diagnostic_cohort_registered)} />
            <Kpi label="Completed" value={number(totals?.diagnostic_cohort_completed)} />
          </div>
          <p className="mt-3 font-mono text-[10px] leading-5 text-mathua-muted">
            A cohort, not an event rate. Completion can happen after the registration range closes.
          </p>
        </ChartCard>
      </div>
    </>
  )
}

function LearningTab({ data, range }: { data: AdminAnalytics; range: Range }) {
  const rows = bucketRows(data, range)
  const masteryEvents = rows.reduce((sum, item) => sum + item.mastery, 0)
  const domainRows = Object.entries(data.mastered_by_concept ?? {})
    .sort((a, b) => b[1] - a[1])
    .slice(0, 12)
  return (
    <>
      <PrimaryKpis data={data} />
      <div className="mt-5 grid items-start gap-4 lg:grid-cols-2">
        <ChartCard title="Accuracy" hint="Correct answers divided by answers in each bucket.">
          <AccuracyTrend rows={rows} />
        </ChartCard>
        <ChartCard title="Mastery events" hint="Concepts recorded as mastered in each bucket.">
          <MasteryTrend rows={rows} />
        </ChartCard>
        <ChartCard title="XP ledger" hint="Awarded and deducted events.">
          <XPTrend rows={rows} />
        </ChartCard>
        <ChartCard title="Currently mastered" hint="Current state, not a range total.">
          <div className="flex flex-wrap items-baseline gap-3">
            <p className="font-mono text-3xl text-mathua-primary tabular-nums">{number(data.mastered_current)}</p>
            <p className="font-mono text-[10px] text-mathua-muted">
              {data.accounts?.registered ? `Current masteries across ${number(data.accounts.registered)} registered accounts` : 'Re-mastery overwrites the date.'}
            </p>
          </div>
          <div className="mt-4 grid grid-cols-2 gap-2">
            <Kpi label="Mastery events" value={number(masteryEvents)} hint="Concepts recorded as mastered in range" />
            <Kpi label="Timed quizzes" value={number(data.totals?.timed_quiz_completions)} hint="Separate assessment" />
          </div>
        </ChartCard>
        <ChartCard title="Top mastered concepts" hint="Concept ids with current mastery counts.">
          {domainRows.length === 0 ? <Empty /> : (
            <ul className="space-y-1">
              {domainRows.map(([id, count]) => (
                <li key={id} className="flex justify-between gap-3 font-mono text-[11px]">
                  <span className="truncate text-mathua-secondary" title={id}>{id}</span>
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
  const rows = bucketRows(data, range)
  return (
    <>
      <PrimaryKpis data={data} />
      <div className="mt-5 grid items-start gap-4 lg:grid-cols-2">
        <ChartCard title="Question reports" hint="Opened and resolved in the range.">
          <ReportTrend rows={rows} />
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

function ChartCard({ title, hint, className = '', children }: { title: string; hint?: string; className?: string; children: React.ReactNode }) {
  return (
    <section className={`min-w-0 border border-mathua-border bg-mathua-surface p-4 ${className}`}>
      <h3 className="font-mono text-[11px] uppercase tracking-wider text-mathua-secondary">{title}</h3>
      {hint && <p className="mt-1 font-mono text-[10px] leading-5 text-mathua-muted">{hint}</p>}
      <div className="mt-3">{children}</div>
    </section>
  )
}

function Empty() {
  return <p className="py-8 text-center font-mono text-[11px] text-mathua-muted">No data in this window.</p>
}

function ComposedTrend({
  rows,
  config,
  rightAxis,
  leftDomain,
  children,
}: {
  rows: Row[]
  config: ChartConfig
  rightAxis?: 'count' | 'percent'
  leftDomain?: [string, string]
  children: React.ReactNode
}) {
  if (rows.length === 0) return <Empty />
  return (
    <ChartContainer config={config} className="h-[220px] w-full">
      <ComposedChart data={rows} margin={{ top: 8, right: 4, left: -18, bottom: 0 }}>
        <CartesianGrid vertical={false} strokeDasharray="3 3" />
        <XAxis dataKey="label" tickLine={false} axisLine={false} tickMargin={8} minTickGap={22} />
        <YAxis yAxisId="left" allowDecimals={false} tickLine={false} axisLine={false} width={52} domain={leftDomain} tickFormatter={value => Number(value).toLocaleString()} />
        <YAxis
          yAxisId="right"
          orientation="right"
          allowDecimals={rightAxis === 'percent'}
          tickLine={false}
          axisLine={false}
          width={40}
          domain={rightAxis === 'percent' ? [0, 100] : undefined}
          tickFormatter={rightAxis === 'percent' ? value => `${value}%` : undefined}
        />
        <ChartTooltip content={<ChartTooltipContent indicator="line" />} />
        <ChartLegend content={<ChartLegendContent />} />
        {children}
      </ComposedChart>
    </ChartContainer>
  )
}

function ActivityTrend({ rows }: { rows: Row[] }) {
  if (rows.length === 0) return <Empty />
  return (
    <ChartContainer config={ACTIVITY_CONFIG} className="h-[220px] w-full">
      <ComposedChart data={rows} margin={{ top: 8, right: 4, left: -18, bottom: 0 }}>
        <CartesianGrid vertical={false} strokeDasharray="3 3" />
        <XAxis dataKey="label" tickLine={false} axisLine={false} tickMargin={8} minTickGap={22} />
        <YAxis yAxisId="left" allowDecimals={false} tickLine={false} axisLine={false} width={44} />
        <YAxis yAxisId="right" orientation="right" allowDecimals={false} tickLine={false} axisLine={false} width={40} />
        <ChartTooltip content={<ChartTooltipContent indicator="dot" />} />
        <ChartLegend content={<ChartLegendContent />} />
        <Line yAxisId="left" type="monotone" dataKey="activeRegistered" stroke="var(--color-activeRegistered)" strokeWidth={2} dot={false} />
        <Line yAxisId="left" type="monotone" dataKey="activeGuests" stroke="var(--color-activeGuests)" strokeWidth={2} dot={false} />
        <Line yAxisId="left" type="monotone" dataKey="activeUnclassified" stroke="var(--color-activeUnclassified)" strokeWidth={2} dot={false} />
        <Line yAxisId="right" type="monotone" dataKey="answers" stroke="var(--color-answers)" strokeWidth={2} dot={false} />
      </ComposedChart>
    </ChartContainer>
  )
}

function RegistrationTrend({ rows }: { rows: Row[] }) {
  if (rows.length === 0) return <Empty />
  return (
    <ChartContainer config={REGISTRATION_CONFIG} className="h-[220px] w-full">
      <AreaChart data={rows} margin={{ top: 8, right: 4, left: -18, bottom: 0 }}>
        <defs>
          <linearGradient id="signup-fill" x1="0" y1="0" x2="0" y2="1">
            <stop offset="0%" stopColor="var(--color-signups)" stopOpacity={0.32} />
            <stop offset="100%" stopColor="var(--color-signups)" stopOpacity={0.02} />
          </linearGradient>
        </defs>
        <CartesianGrid vertical={false} strokeDasharray="3 3" />
        <XAxis dataKey="label" tickLine={false} axisLine={false} tickMargin={8} minTickGap={22} />
        <YAxis allowDecimals={false} tickLine={false} axisLine={false} width={44} />
        <ChartTooltip content={<ChartTooltipContent indicator="line" />} />
        <ChartLegend content={<ChartLegendContent />} />
        <Area type="monotone" dataKey="signups" stroke="var(--color-signups)" strokeWidth={2} fill="url(#signup-fill)" />
      </AreaChart>
    </ChartContainer>
  )
}

function QuizTrend({ rows }: { rows: Row[] }) {
  if (rows.length === 0) return <Empty />
  return (
    <ChartContainer config={QUIZ_CONFIG} className="h-[220px] w-full">
      <ComposedChart data={rows} margin={{ top: 8, right: 4, left: -18, bottom: 0 }}>
        <CartesianGrid vertical={false} strokeDasharray="3 3" />
        <XAxis dataKey="label" tickLine={false} axisLine={false} tickMargin={8} minTickGap={22} />
        <YAxis allowDecimals={false} tickLine={false} axisLine={false} width={44} />
        <ChartTooltip content={<ChartTooltipContent indicator="line" />} />
        <ChartLegend content={<ChartLegendContent />} />
        <Line type="monotone" dataKey="quizzes" stroke="var(--color-quizzes)" strokeWidth={2} dot={false} />
        <Line type="monotone" dataKey="timedQuizzes" stroke="var(--color-timedQuizzes)" strokeWidth={2} dot={false} />
      </ComposedChart>
    </ChartContainer>
  )
}

function XPTrend({ rows }: { rows: Row[] }) {
  if (rows.length === 0) return <Empty />
  return (
    <ChartContainer config={XP_CONFIG} className="h-[220px] w-full">
      <ComposedChart data={rows} margin={{ top: 8, right: 4, left: -18, bottom: 0 }}>
        <CartesianGrid vertical={false} strokeDasharray="3 3" />
        <XAxis dataKey="label" tickLine={false} axisLine={false} tickMargin={8} minTickGap={22} />
        <YAxis allowDecimals={false} tickLine={false} axisLine={false} width={48} />
        <ChartTooltip content={<ChartTooltipContent indicator="line" />} />
        <ChartLegend content={<ChartLegendContent />} />
        <Line type="monotone" dataKey="xpAwarded" stroke="var(--color-xpAwarded)" strokeWidth={2} dot={false} />
        <Line type="monotone" dataKey="xpDeducted" stroke="var(--color-xpDeducted)" strokeWidth={2} dot={false} />
      </ComposedChart>
    </ChartContainer>
  )
}

function ReportTrend({ rows }: { rows: Row[] }) {
  if (rows.length === 0) return <Empty />
  return (
    <ChartContainer config={REPORT_CONFIG} className="h-[220px] w-full">
      <ComposedChart data={rows} margin={{ top: 8, right: 4, left: -18, bottom: 0 }}>
        <CartesianGrid vertical={false} strokeDasharray="3 3" />
        <XAxis dataKey="label" tickLine={false} axisLine={false} tickMargin={8} minTickGap={22} />
        <YAxis allowDecimals={false} tickLine={false} axisLine={false} width={44} />
        <ChartTooltip content={<ChartTooltipContent indicator="line" />} />
        <ChartLegend content={<ChartLegendContent />} />
        <Line type="monotone" dataKey="reportsOpened" stroke="var(--color-reportsOpened)" strokeWidth={2} dot={false} />
        <Line type="monotone" dataKey="reportsResolved" stroke="var(--color-reportsResolved)" strokeWidth={2} dot={false} />
      </ComposedChart>
    </ChartContainer>
  )
}

function IntensityTrend({ rows }: { rows: Row[] }) {
  if (rows.length === 0) return <Empty />
  return (
    <ChartContainer config={INTENSITY_CONFIG} className="h-[220px] w-full">
      <AreaChart data={rows} margin={{ top: 8, right: 4, left: -18, bottom: 0 }}>
        <defs>
          <linearGradient id="intensity-fill" x1="0" y1="0" x2="0" y2="1">
            <stop offset="0%" stopColor="var(--color-answersPerLearner)" stopOpacity={0.28} />
            <stop offset="100%" stopColor="var(--color-answersPerLearner)" stopOpacity={0.02} />
          </linearGradient>
        </defs>
        <CartesianGrid vertical={false} strokeDasharray="3 3" />
        <XAxis dataKey="label" tickLine={false} axisLine={false} tickMargin={8} minTickGap={22} />
        <YAxis allowDecimals tickLine={false} axisLine={false} width={44} tickFormatter={value => Number(value).toFixed(1)} />
        <ChartTooltip content={<ChartTooltipContent indicator="line" />} />
        <ChartLegend content={<ChartLegendContent />} />
        <Area type="monotone" dataKey="answersPerLearner" stroke="var(--color-answersPerLearner)" strokeWidth={2} fill="url(#intensity-fill)" />
      </AreaChart>
    </ChartContainer>
  )
}

function AccuracyTrend({ rows }: { rows: Row[] }) {
  if (rows.length === 0) return <Empty />
  return (
    <ChartContainer config={ACCURACY_CONFIG} className="h-[220px] w-full">
      <LineChart data={rows} margin={{ top: 8, right: 4, left: -18, bottom: 0 }}>
        <CartesianGrid vertical={false} strokeDasharray="3 3" />
        <XAxis dataKey="label" tickLine={false} axisLine={false} tickMargin={8} minTickGap={22} />
        <YAxis domain={[0, 100]} allowDecimals={false} tickLine={false} axisLine={false} width={44} tickFormatter={value => `${value}%`} />
        <ChartTooltip content={<ChartTooltipContent indicator="line" />} />
        <ChartLegend content={<ChartLegendContent />} />
        <Line type="monotone" dataKey="accuracy" stroke="var(--color-accuracy)" strokeWidth={2} dot={false} />
      </LineChart>
    </ChartContainer>
  )
}

function MasteryTrend({ rows }: { rows: Row[] }) {
  if (rows.length === 0) return <Empty />
  return (
    <ChartContainer config={MASTERY_CONFIG} className="h-[220px] w-full">
      <LineChart data={rows} margin={{ top: 8, right: 4, left: -18, bottom: 0 }}>
        <CartesianGrid vertical={false} strokeDasharray="3 3" />
        <XAxis dataKey="label" tickLine={false} axisLine={false} tickMargin={8} minTickGap={22} />
        <YAxis allowDecimals={false} tickLine={false} axisLine={false} width={44} />
        <ChartTooltip content={<ChartTooltipContent indicator="line" />} />
        <ChartLegend content={<ChartLegendContent />} />
        <Line type="monotone" dataKey="mastery" stroke="var(--color-mastery)" strokeWidth={2} dot={false} />
      </LineChart>
    </ChartContainer>
  )
}

function AccountMix({ data }: { data: AdminAnalytics }) {
  const accounts = data.accounts
  if (!accounts || accounts.total === 0) return <Empty />
  const entries = [
    ['registered', 'Registered', 'var(--chart-1)', accounts.registered],
    ['guest', 'Guests', 'var(--chart-2)', accounts.guest],
    ['unclassified', 'Unclassified', 'var(--chart-4)', accounts.unclassified],
  ] as const
  return (
    <div>
      <p className="font-mono text-3xl text-mathua-primary tabular-nums">{number(accounts.total)}</p>
      <p className="mt-1 font-mono text-[10px] text-mathua-muted">student rows, all time</p>
      <div className="mt-4 flex h-3 w-full overflow-hidden border border-mathua-border">
        {entries.map(([key, , color, value]) => (
          <div key={key} style={{ width: `${(value / accounts.total) * 100}%`, background: color }} title={`${key}: ${value}`} />
        ))}
      </div>
      <ul className="mt-3 space-y-2">
        {entries.map(([key, label, color, value]) => (
          <li key={key} className="flex items-center justify-between gap-3 font-mono text-[11px]">
            <span className="flex items-center gap-2 text-mathua-secondary">
              <span className="h-2 w-2" style={{ background: color }} />
              {label}
            </span>
            <span className="tabular-nums text-mathua-primary">{number(value)}</span>
          </li>
        ))}
      </ul>
      <p className="mt-3 font-mono text-[9px] leading-4 text-mathua-muted">
        Unclassified rows predate account classification and are kept separate.
      </p>
    </div>
  )
}

function SegmentedBar({ values }: { values: Record<string, number> }) {
  const entries = Object.entries(values).filter(([, value]) => value > 0)
  const total = entries.reduce((sum, [, value]) => sum + value, 0)
  if (total === 0) return <Empty />
  const colors: Record<string, string> = {
    open: 'var(--chart-4)',
    reviewing: 'var(--chart-2)',
    resolved: 'var(--chart-3)',
    dismissed: 'var(--text-muted)',
  }
  return (
    <div>
      <div className="flex h-3 w-full overflow-hidden border border-mathua-border">
        {entries.map(([key, value]) => (
          <div key={key} style={{ width: `${(value / total) * 100}%`, background: colors[key] ?? 'var(--chart-1)' }} title={`${key}: ${value}`} />
        ))}
      </div>
      <ul className="mt-3 flex flex-wrap gap-x-4 gap-y-2">
        {entries.map(([key, value]) => (
          <li key={key} className="flex items-center gap-1.5 font-mono text-[10px] text-mathua-secondary">
            <span className="h-2 w-2" style={{ background: colors[key] ?? 'var(--chart-1)' }} />
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
    <ChartCard title="Notification delivery" hint="All-time delivery rows by channel and status.">
      {error ? <p className="font-mono text-[11px] text-mathua-red">{error}</p> : (
        <>
          <div className="grid grid-cols-2 gap-2">
            {stats.slice(0, 4).map(stat => (
              <Kpi key={`${stat.channel}-${stat.status}`} label={`${stat.channel} ${stat.status}`} value={number(stat.count)} />
            ))}
          </div>
          <p className="mt-3 font-mono text-[10px] leading-5 text-mathua-muted">
            {problems.length === 0 ? 'No failed or suppressed deliveries.' : `${problems.length} failed or suppressed deliveries.`}
          </p>
        </>
      )}
    </ChartCard>
  )
}

function deltaLabel(current: number | undefined, previous: number | undefined): string | undefined {
  if (current === undefined || previous === undefined) return undefined
  if (previous === 0) return current === 0 ? 'no change' : 'new in this window'
  const change = ((current - previous) / previous) * 100
  const sign = change >= 0 ? '+' : ''
  return `${sign}${change.toFixed(1)}% vs previous`
}

function ppDelta(correct: number | undefined, answers: number | undefined, previousCorrect: number | undefined, previousAnswers: number | undefined): string | undefined {
  if (!answers || !previousAnswers || correct === undefined || previousCorrect === undefined) return undefined
  const change = (correct / answers - previousCorrect / previousAnswers) * 100
  const sign = change >= 0 ? '+' : ''
  return `${sign}${change.toFixed(1)} pp vs previous`
}

function updatedLabel(value: string | undefined): string {
  if (!value) return 'unknown'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return 'unknown'
  return `${date.toLocaleTimeString('en-GB', { hour: '2-digit', minute: '2-digit', timeZone: 'UTC' })} UTC`
}

function rangeLabel(start: number, range: Range): string {
  const date = new Date(start * 1000)
  if (range === '15m' || range === '1h' || range === '24h') {
    return date.toLocaleTimeString('en-GB', { hour: '2-digit', minute: '2-digit', timeZone: 'UTC' })
  }
  return date.toLocaleDateString('en-GB', { month: 'short', day: 'numeric', timeZone: 'UTC' })
}
