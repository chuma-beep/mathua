'use client'

import { Suspense, useCallback, useEffect, useMemo, useState } from 'react'
import { useSearchParams } from 'next/navigation'
import Link from 'next/link'
import AppShell from '../../components/app-shell'
import SectionHeader from '../../components/SectionHeader'
import Footer from '../../components/Footer'
import BoardTable from '../../components/leaderboard/board-table'
import LeagueBoard from '../../components/leaderboard/league-board'
import Podium from '../../components/leaderboard/podium'
import { MonthPicker, PeriodArrows, WeekPicker } from '../../components/leaderboard/period-picker'
import { useTheme } from '../../hooks/useTheme'
import { useAuthState } from '../../hooks/useAuthState'
import {
  getLeaderboard,
  getMyLeague,
  getSettings,
  type LeaderboardBoard,
  type LeaderboardView,
  type MyLeague,
} from '../../lib/api'
import {
  monthStartKey,
  monthsAvailable,
  parsePeriodKey,
  periodLabel,
  shiftMonthKey,
  shiftWeekKey,
  todayKey,
  weekStartKey,
} from '../../lib/leaderboard-period'
import { Button } from '../../components/ui/button'
import { scrollToPageTop } from '../../lib/scroll'
import { Tabs, TabsList, TabsTrigger } from '../../components/ui/tabs'

type Tab = LeaderboardView | 'league'

const PAGE_SIZE = 25

function parseTab(raw: string | null): Tab {
  if (raw === 'monthly' || raw === 'all_time' || raw === 'league') return raw
  return 'weekly'
}

function parsePage(raw: string | null): number {
  const value = Number(raw)
  return Number.isInteger(value) && value > 0 ? value : 1
}

function computeCountdown(): string {
  const now = new Date()
  const nextMonday = new Date(now)
  nextMonday.setUTCDate(now.getUTCDate() + ((7 - now.getUTCDay() + 1) % 7 || 7))
  nextMonday.setUTCHours(0, 0, 0, 0)
  const diff = nextMonday.getTime() - now.getTime()
  const d = Math.floor(diff / 86400000)
  const h = Math.floor((diff % 86400000) / 3600000)
  const m = Math.floor((diff % 3600000) / 60000)
  return `${d}d ${h}h ${m}m`
}

function PositionStrip({
  me,
  view,
  final,
}: {
  me: NonNullable<LeaderboardBoard['me']>
  view: LeaderboardView
  final?: boolean
}) {
  if (me.guest) {
    return (
      <p className="mb-4 text-sm text-mathua-secondary">
        You are browsing as a guest.{' '}
        <Link href="/login" className="text-mathua-blue hover:underline">
          Sign in
        </Link>{' '}
        to compete on the board.
      </p>
    )
  }
  if (!me.eligible || me.rank === 0) {
    const period =
      view === 'weekly' ? (final ? ' that week' : ' this week') : view === 'monthly' ? (final ? ' that month' : ' this month') : ''
    return (
      <p className="mb-4 text-sm text-mathua-secondary">
        You have no place yet{period}. Answer one question and your XP puts you on the board.
      </p>
    )
  }
  if (me.in_page) return null
  return (
    <p className="mb-4 border border-mathua-border bg-mathua-blue-faint px-3 py-2 font-mono text-xs text-mathua-primary">
      Your position: <span className="text-mathua-blue">#{me.rank}</span> with {me.xp} XP
    </p>
  )
}

function PageControls({
  board,
  onPage,
  busy,
}: {
  board: LeaderboardBoard
  onPage: (page: number) => void
  busy: boolean
}) {
  const start = board.total === 0 ? 0 : (board.page - 1) * board.page_size + 1
  const end = start === 0 ? 0 : start + board.entries.length - 1
  return (
    <div className="mt-4 flex flex-wrap items-center justify-between gap-3">
      <p className="font-mono text-[11px] text-mathua-muted" aria-live="polite">
        {board.total === 0 ? 'No learners' : `${start}–${end} of ${board.total}`}
      </p>
      <div className="flex items-center gap-2">
        <Button
          type="button"
          variant="outline"
          onClick={() => onPage(board.page - 1)}
          disabled={busy || board.page <= 1}
          className="h-8 rounded-none border-mathua-border font-mono text-xs"
        >
          Previous
        </Button>
        <span className="font-mono text-[11px] text-mathua-muted">
          Page {board.page} of {Math.max(board.total_pages, 1)}
        </span>
        <Button
          type="button"
          variant="outline"
          onClick={() => onPage(board.page + 1)}
          disabled={busy || board.page >= board.total_pages}
          className="h-8 rounded-none border-mathua-border font-mono text-xs"
        >
          Next
        </Button>
      </div>
    </div>
  )
}

function LeaderboardFallback() {
  return (
    <AppShell>
      <div className="mx-auto max-w-container px-4 pb-[calc(80px+env(safe-area-inset-bottom))] sm:px-6 lg:pb-0">
        <section className="pt-8">
          <SectionHeader label="Competition" title="Leaderboard" />
        </section>
      </div>
      <Footer />
    </AppShell>
  )
}

function LeaderboardPageInner() {
  const params = useSearchParams()
  const { mounted } = useTheme()
  const { loggedIn } = useAuthState()
  const [countdown, setCountdown] = useState(computeCountdown)
  const [tab, setTab] = useState<Tab>(() => parseTab(params.get('view')))
  const [anchor, setAnchor] = useState<string | undefined>(() => params.get('anchor') ?? undefined)
  const [page, setPage] = useState(() => parsePage(params.get('page')))
  const [leagueWeek, setLeagueWeek] = useState<string | undefined>(() =>
    parseTab(params.get('view')) === 'league' ? params.get('anchor') ?? undefined : undefined,
  )
  const [reload, setReload] = useState(0)
  const [board, setBoard] = useState<LeaderboardBoard | null>(null)
  const [league, setLeague] = useState<MyLeague | null>(null)
  const [earliest, setEarliest] = useState<string | undefined>(undefined)
  const [loading, setLoading] = useState(true)
  const [failed, setFailed] = useState(false)
  const [showLeaderboards, setShowLeaderboards] = useState(true)

  const showMonth = useMemo(
    () => tab === 'monthly' || monthsAvailable(earliest) || monthsAvailable(league?.earliest_week),
    [earliest, league?.earliest_week, tab],
  )

  useEffect(() => {
    if (!loggedIn) {
      setShowLeaderboards(true)
      return
    }
    let live = true
    getSettings()
      .then((s) => {
        if (live) setShowLeaderboards(s.leaderboard_privacy !== 'hidden')
      })
      .catch(() => {
        if (live) setShowLeaderboards(true)
      })
    return () => {
      live = false
    }
  }, [loggedIn])

  useEffect(() => {
    const update = () => setCountdown(computeCountdown())
    update()
    const interval = setInterval(update, 60000)
    return () => clearInterval(interval)
  }, [])

  useEffect(() => {
    if (!mounted) return
    let live = true
    setLoading(true)
    setFailed(false)
    if (tab === 'league') {
      getMyLeague(leagueWeek)
        .then((data) => {
          if (!live) return
          setLeague(data)
          setLoading(false)
        })
        .catch(() => {
          if (live) {
            setFailed(true)
            setLoading(false)
          }
        })
    } else {
      getLeaderboard(tab, page, PAGE_SIZE, anchor)
        .then((data) => {
          if (!live) return
          setBoard(data)
          if (data.earliest_anchor) setEarliest(data.earliest_anchor)
          setLoading(false)
        })
        .catch(() => {
          if (live) {
            setFailed(true)
            setLoading(false)
          }
        })
    }
    return () => {
      live = false
    }
  }, [mounted, tab, anchor, page, leagueWeek, reload])

  useEffect(() => {
    const query = new URLSearchParams()
    if (tab !== 'weekly') query.set('view', tab)
    const activeAnchor = tab === 'league' ? leagueWeek : anchor
    if (activeAnchor) query.set('anchor', activeAnchor)
    if (page > 1) query.set('page', String(page))
    const suffix = query.toString()
    window.history.replaceState(null, '', suffix ? `${window.location.pathname}?${suffix}` : window.location.pathname)
  }, [tab, anchor, page, leagueWeek])

  const selectTab = useCallback((next: Tab) => {
    setTab(next)
    setPage(1)
    setAnchor(undefined)
    setLeagueWeek(undefined)
  }, [])

  const changePage = useCallback((next: number) => {
    setPage(next)
    scrollToPageTop()
  }, [])

  const retry = useCallback(() => setReload((n) => n + 1), [])

  const view = tab === 'league' ? null : tab
  const periodStart = board?.period_start ? board.period_start.slice(0, 10) : undefined
  const currentWeek = weekStartKey(new Date())
  const currentMonth = monthStartKey(new Date())
  const today = todayKey()
  const weekValue = periodStart ?? anchor ?? currentWeek
  const monthValue = periodStart ?? anchor ?? currentMonth

  const earliestWeek = useMemo(() => {
    const value = earliest ? parsePeriodKey(earliest) : undefined
    return value ? weekStartKey(value) : undefined
  }, [earliest])
  const earliestMonth = useMemo(() => {
    const value = earliest ? parsePeriodKey(earliest) : undefined
    return value ? monthStartKey(value) : undefined
  }, [earliest])

  if (!mounted) {
    return <LeaderboardFallback />
  }

  if (!showLeaderboards) {
    return (
      <AppShell>
        <div className="mx-auto max-w-container min-w-0 overflow-x-hidden px-4 pb-[calc(80px+env(safe-area-inset-bottom))] sm:px-6 lg:pb-0">
          <section className="min-w-0 overflow-hidden pt-8">
            <span className="mb-4 flex">
              <button
                type="button"
                onClick={() => {
                  if (window.history.length > 1) window.history.back()
                }}
                className="text-sm text-mathua-secondary hover:text-mathua-blue"
              >
                ← Back
              </button>
            </span>
            <SectionHeader label="Competition" title="Leaderboard" />
            <div className="mx-auto mt-6 max-w-[560px] border border-mathua-border bg-mathua-surface p-6 text-center">
              <p className="font-mono text-sm text-mathua-primary">Leaderboards are turned off</p>
              <p className="mt-2 text-sm text-mathua-secondary">
                You are anonymous on the public ranking, so your name never appears on it. Change this in{' '}
                <Link href="/settings" className="text-mathua-blue hover:underline">
                  Settings
                </Link>
                .
              </p>
            </div>
          </section>
        </div>
        <Footer />
      </AppShell>
    )
  }

  const tabs: { id: Tab; label: string }[] = [
    { id: 'weekly', label: 'Week' },
    ...(showMonth ? [{ id: 'monthly' as const, label: 'Month' }] : []),
    { id: 'all_time', label: 'All time' },
    { id: 'league', label: 'My league' },
  ]

  const atEarliest =
    view === 'weekly'
      ? !earliestWeek || weekValue <= earliestWeek
      : view === 'monthly'
        ? !earliestMonth || monthValue <= earliestMonth
        : true

  return (
    <AppShell>
      <div className="mx-auto max-w-container min-w-0 overflow-x-hidden px-4 pb-[calc(80px+env(safe-area-inset-bottom))] sm:px-6 lg:pb-0">
        <section className="min-w-0 overflow-hidden pt-8">
          <span className="mb-4 flex">
            <button
              type="button"
              onClick={() => {
                if (window.history.length > 1) window.history.back()
              }}
              className="text-sm text-mathua-secondary hover:text-mathua-blue"
            >
              ← Back
            </button>
          </span>

          <SectionHeader label="Competition" title="Leaderboard" />

          <p className="mx-auto mb-6 mt-4 max-w-[560px] px-2 text-center text-sm text-mathua-secondary">
            XP comes from answered questions, not from time spent. Weeks run Monday 00:00 to Monday
            00:00 UTC and months are calendar months in UTC. The current board is live; a closed
            period stays fixed.
          </p>

          <div className="mb-5 flex flex-wrap items-center justify-center gap-2">
            {view === 'weekly' && (
              <>
                <PeriodArrows
                  onPrev={() => {
                    setAnchor(shiftWeekKey(weekValue, -1))
                    setPage(1)
                  }}
                  onNext={() => {
                    setAnchor(shiftWeekKey(weekValue, 1))
                    setPage(1)
                  }}
                  prevDisabled={atEarliest}
                  nextDisabled={!board?.final}
                  label="week"
                />
                <WeekPicker
                  value={weekValue}
                  onChange={(next) => {
                    setAnchor(next)
                    setPage(1)
                  }}
                  min={earliestWeek}
                  max={today}
                />
                {board && !board.final && (
                  <span className="border border-mathua-border bg-mathua-surface px-3 py-2 font-mono text-[10px] uppercase text-mathua-muted">
                    Resets in <span className="text-mathua-blue">{countdown}</span>
                  </span>
                )}
                {board?.final && (
                  <span className="border border-mathua-border bg-mathua-surface px-3 py-2 font-mono text-[10px] uppercase text-mathua-muted">
                    Final
                  </span>
                )}
              </>
            )}
            {view === 'monthly' && (
              <>
                <PeriodArrows
                  onPrev={() => {
                    setAnchor(shiftMonthKey(monthValue, -1))
                    setPage(1)
                  }}
                  onNext={() => {
                    setAnchor(shiftMonthKey(monthValue, 1))
                    setPage(1)
                  }}
                  prevDisabled={atEarliest}
                  nextDisabled={!board?.final}
                  label="month"
                />
                <MonthPicker
                  value={monthValue}
                  onChange={(next) => {
                    setAnchor(next)
                    setPage(1)
                  }}
                  min={earliestMonth}
                  max={currentMonth}
                />
                {board?.final ? (
                  <span className="border border-mathua-border bg-mathua-surface px-3 py-2 font-mono text-[10px] uppercase text-mathua-muted">
                    Final
                  </span>
                ) : (
                  <span className="border border-mathua-border bg-mathua-surface px-3 py-2 font-mono text-[10px] uppercase text-mathua-muted">
                    Live
                  </span>
                )}
              </>
            )}
          </div>

          <Tabs value={tab} onValueChange={(value) => selectTab(value as Tab)}>
            <TabsList className="mb-6 flex h-auto w-full flex-wrap justify-center gap-1 rounded-none border-b border-mathua-border bg-transparent p-0">
              {tabs.map((item) => (
                <TabsTrigger
                  key={item.id}
                  value={item.id}
                  className="rounded-none border-b-2 border-transparent px-3 py-2 font-mono text-xs uppercase tracking-wider text-mathua-muted data-[state=active]:border-mathua-blue data-[state=active]:bg-transparent data-[state=active]:text-mathua-blue data-[state=active]:shadow-none sm:px-4"
                >
                  {item.label}
                </TabsTrigger>
              ))}
            </TabsList>
          </Tabs>

          {view && (
            <p className="mb-3 text-center font-mono text-[11px] uppercase tracking-wider text-mathua-muted">
              {periodLabel(view, board?.period_start)}
            </p>
          )}
        </section>

        <div className="mb-12 min-w-0">
          {view && (
            <>
              {failed && (
                <p className="py-8 text-center text-sm text-mathua-muted">
                  Could not load this board.{' '}
                  <Button type="button" variant="link" onClick={retry} className="h-auto p-0 font-mono text-xs">
                    Retry
                  </Button>
                </p>
              )}
              {!failed && (
                <>
                  {!loading && board && board.page === 1 && (
                    <Podium
                      entries={board.entries}
                      score={
                        view === 'all_time'
                          ? 'lifetime'
                          : view === 'weekly'
                            ? board?.final
                              ? 'that week'
                              : 'this week'
                            : board?.final
                              ? 'that month'
                              : 'this month'
                      }
                    />
                  )}
                  {!loading && board?.me && <PositionStrip me={board.me} view={view} final={board.final} />}
                  {!loading && board && board.entries.length === 0 && (
                    <p className="px-2 py-8 text-center text-sm text-mathua-muted">
                      {board.final
                        ? 'No one scored in this period.'
                        : view === 'weekly'
                          ? 'No one has earned XP this week yet. Answer a question and the board starts with you.'
                          : view === 'monthly'
                            ? 'No one has earned XP this month yet. Answer a question and the board starts with you.'
                            : 'No learner has earned XP yet. The first answer puts the first name here.'}{' '}
                      {!board.final && (
                        <Link href="/learn" className="text-mathua-blue hover:underline">
                          Start learning
                        </Link>
                      )}
                    </p>
                  )}
                  {(loading || (board && board.entries.length > 0)) && (
                    <BoardTable
                      board={board}
                      view={view}
                      loading={loading}
                      caption={`${periodLabel(view, board?.period_start)} XP standings`}
                    />
                  )}
                  {!loading && board && board.entries.length > 0 && (
                    <PageControls board={board} onPage={changePage} busy={loading} />
                  )}
                </>
              )}
            </>
          )}

          {tab === 'league' && (
            <>
              {failed && (
                <p className="py-8 text-center text-sm text-mathua-muted">
                  Could not load your league.{' '}
                  <Button type="button" variant="link" onClick={retry} className="h-auto p-0 font-mono text-xs">
                    Retry
                  </Button>
                </p>
              )}
              {!failed && (
                <LeagueBoard
                  league={league}
                  week={leagueWeek ?? league?.week ?? currentWeek}
                  earliest={league?.earliest_week}
                  loading={loading}
                  onWeek={(next) => setLeagueWeek(next)}
                />
              )}
            </>
          )}
        </div>
      </div>

      <Footer />
    </AppShell>
  )
}

export default function LeaderboardPage() {
  return (
    <Suspense fallback={<LeaderboardFallback />}>
      <LeaderboardPageInner />
    </Suspense>
  )
}
