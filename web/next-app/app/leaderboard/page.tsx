'use client'

import { useTheme } from '../../hooks/useTheme'
import { useState, useEffect, useCallback } from 'react'
import Link from 'next/link'
import AppShell from '../../components/app-shell'
import SectionHeader from '../../components/SectionHeader'
import Footer from '../../components/Footer'
import Avatar from '../../components/Avatar'
import LeagueAvatar from '../../components/LeagueAvatar'
import { LeaderboardTableSkeleton } from '../../components/skeletons/LeaderboardSkeleton'
import { useAuthState } from '../../hooks/useAuthState'
import {
  getLeaderboard,
  getMyLeague,
  getSettings,
  type LeaderboardBoard,
  type LeaderboardEntry,
  type LeaderboardView,
  type MyLeague,
} from '../../lib/api'
import { resolveAvatar } from '../../lib/dicebear'
import { scrollToPageTop } from '../../lib/scroll'

type Tab = LeaderboardView | 'league'

const TABS: { id: Tab; label: string }[] = [
  { id: 'weekly', label: 'This week' },
  { id: 'all_time', label: 'All time' },
  { id: 'league', label: 'My league' },
]

const PAGE_SIZE = 25

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

// The server decides the avatar shape; the client only has to turn it into a
// URL. A generated avatar is seeded from the public id, never from the account.
function avatarURL(row: { id: string; avatar?: { kind: string; url?: string; style?: string; version?: number } }): string | undefined {
  const av = row.avatar
  if (!av) return undefined
  if (av.kind === 'custom' && av.url) {
    return av.version ? `${av.url}?v=${av.version}` : av.url
  }
  return resolveAvatar(
    { avatar_url: '' },
    { avatar_dicebear: av.style ? { style: av.style, seed: row.id } : null },
  ).url
}

// The server resolves a display name, but a blank cell is a silent failure and
// costs one line to prevent: fall back to the public id, which is never empty.
function nameOf(row: { name: string; id: string }): string {
  return row.name.trim() || row.id
}

function shortWeek(iso: string): string {
  const d = new Date(iso)
  if (Number.isNaN(d.getTime())) return ''
  return d.toLocaleDateString(undefined, { month: 'short', day: 'numeric', timeZone: 'UTC' })
}

function Podium({ entries }: { entries: LeaderboardEntry[] }) {
  if (entries.length < 3) return null
  return (
    <ol className="grid grid-cols-1 sm:grid-cols-3 gap-3 mb-6 list-none p-0">
      {entries.slice(0, 3).map((row, i) => (
        <li
          key={`podium-${row.id}`}
          className={`border border-mathua-border bg-mathua-surface px-4 py-3 flex items-center gap-3 min-w-0 ${i === 0 ? 'sm:order-2' : i === 1 ? 'sm:order-1' : 'sm:order-3'}`}
        >
          <span className="font-mono text-lg text-mathua-blue w-5 shrink-0">{row.rank}</span>
          <Avatar seed={row.id} name={nameOf(row)} size={30} url={avatarURL(row)} className="shrink-0" />
          <span className="min-w-0 flex-1">
            <span className="block text-sm text-mathua-primary truncate" title={nameOf(row)}>{nameOf(row)}</span>
            <span className="block font-mono text-[10px] uppercase tracking-wider text-mathua-muted">
              {row.xp} XP
            </span>
          </span>
        </li>
      ))}
    </ol>
  )
}

function PositionStrip({
  me,
  view,
}: {
  me: NonNullable<LeaderboardBoard['me']>
  view: LeaderboardView
}) {
  if (me.guest) {
    return (
      <p className="mb-4 text-sm text-mathua-secondary">
        You are browsing as a guest. <Link href="/login" className="text-mathua-blue hover:underline">Sign in</Link> to compete on the board.
      </p>
    )
  }
  if (!me.eligible || me.rank === 0) {
    return (
      <p className="mb-4 text-sm text-mathua-secondary">
        You have no place yet{view === 'weekly' ? ' this week' : ''}. Answer one question and your XP puts you on the board.
      </p>
    )
  }
  if (me.in_page) return null
  return (
    <p className="mb-4 font-mono text-xs text-mathua-primary border border-mathua-border bg-mathua-blue-faint px-3 py-2">
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
    <div className="flex flex-wrap items-center justify-between gap-3 mt-4">
      <p className="font-mono text-[11px] text-mathua-muted" aria-live="polite">
        {board.total === 0 ? 'No learners' : `${start}–${end} of ${board.total}`}
      </p>
      <div className="flex items-center gap-2">
        <button
          type="button"
          onClick={() => onPage(board.page - 1)}
          disabled={busy || board.page <= 1}
          className="font-mono text-xs px-3 py-1.5 border border-mathua-border text-mathua-primary disabled:opacity-40 disabled:cursor-not-allowed hover:border-mathua-blue"
        >
          Previous
        </button>
        <span className="font-mono text-[11px] text-mathua-muted">
          Page {board.page} of {board.total_pages}
        </span>
        <button
          type="button"
          onClick={() => onPage(board.page + 1)}
          disabled={busy || board.page >= board.total_pages}
          className="font-mono text-xs px-3 py-1.5 border border-mathua-border text-mathua-primary disabled:opacity-40 disabled:cursor-not-allowed hover:border-mathua-blue"
        >
          Next
        </button>
      </div>
    </div>
  )
}

function BoardTable({ board, view }: { board: LeaderboardBoard; view: LeaderboardView }) {
  const scoreLabel = view === 'weekly' ? 'This week' : 'Lifetime'
  return (
    <div className="w-[calc(100%+2rem)] -mx-4 px-4 sm:w-full sm:mx-0 sm:px-0 overflow-x-auto overscroll-x-contain">
      <table className="w-full min-w-[320px] border-collapse bg-mathua-surface border border-mathua-border">
        <caption className="sr-only">
          {view === 'weekly' ? 'This week’s XP standings' : 'Lifetime XP standings'}
        </caption>
        <thead>
          <tr>
            <th scope="col" className="font-mono text-[11px] uppercase text-mathua-muted text-left p-2.5 sm:p-[14px_20px] bg-mathua-surface-elevated border-b border-mathua-border w-10 sm:w-12">#</th>
            <th scope="col" className="font-mono text-[11px] uppercase text-mathua-muted text-left p-2.5 sm:p-[14px_20px] bg-mathua-surface-elevated border-b border-mathua-border">Learner</th>
            <th scope="col" className="font-mono text-[11px] uppercase text-mathua-muted text-right p-2.5 sm:p-[14px_20px] bg-mathua-surface-elevated border-b border-mathua-border">{scoreLabel}</th>
            <th scope="col" className="font-mono text-[11px] uppercase text-mathua-muted text-right p-2.5 sm:p-[14px_20px] bg-mathua-surface-elevated border-b border-mathua-border hidden md:table-cell">Level</th>
          </tr>
        </thead>
        <tbody>
          {board.entries.map((row) => (
            <tr
              key={row.id}
              className={`border-b border-mathua-border last:border-b-0 ${row.is_me ? 'bg-mathua-blue-faint' : ''}`}
            >
              <td className="p-2.5 sm:p-[14px_20px] font-mono text-sm text-mathua-muted">{row.rank}</td>
              <td className="p-2.5 sm:p-[14px_20px] text-mathua-primary text-sm font-medium max-w-[160px] sm:max-w-none">
                <span className="flex items-center gap-2 min-w-0">
                  <Avatar seed={row.id} name={nameOf(row)} size={28} url={avatarURL(row)} className="shrink-0" />
                  <span className="block truncate" title={nameOf(row)}>{nameOf(row)}</span>
                  {row.is_me && (
                    <span className="font-mono text-[10px] uppercase text-mathua-blue shrink-0">you</span>
                  )}
                </span>
              </td>
              <td className="p-2.5 sm:p-[14px_20px] font-mono text-sm text-mathua-green text-right whitespace-nowrap">{row.xp}</td>
              <td className="p-2.5 sm:p-[14px_20px] text-mathua-secondary text-sm text-right hidden md:table-cell truncate max-w-[130px]">{row.level}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

export default function LeaderboardPage() {
  const { mounted } = useTheme()
  const { loggedIn } = useAuthState()
  const [countdown, setCountdown] = useState(computeCountdown)
  const [tab, setTab] = useState<Tab>('weekly')
  const [page, setPage] = useState(1)
  const [reload, setReload] = useState(0)
  const [board, setBoard] = useState<LeaderboardBoard | null>(null)
  const [league, setLeague] = useState<MyLeague | null>(null)
  const [loading, setLoading] = useState(true)
  const [failed, setFailed] = useState(false)
  // Visible until the account's own settings say otherwise. A signed-out visitor
  // has no settings to read, and a failed read cannot be evidence that a learner
  // turned leaderboards off.
  const [showLeaderboards, setShowLeaderboards] = useState(true)

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
    let live = true
    setLoading(true)
    setFailed(false)
    if (tab === 'league') {
      getMyLeague()
        .then((data) => { if (live) { setLeague(data); setLoading(false) } })
        .catch(() => { if (live) { setFailed(true); setLoading(false) } })
    } else {
      getLeaderboard(tab, page, PAGE_SIZE)
        .then((data) => { if (live) { setBoard(data); setLoading(false) } })
        .catch(() => { if (live) { setFailed(true); setLoading(false) } })
    }
    return () => { live = false }
  }, [tab, page, reload])

  const selectTab = useCallback((next: Tab) => {
    setTab(next)
    setPage(1)
  }, [])

  const changePage = useCallback((next: number) => {
    setPage(next)
    scrollToPageTop()
  }, [])

  const retry = useCallback(() => setReload((n) => n + 1), [])

  if (!mounted) {
    return (
      <AppShell>
        <div style={{ background: 'var(--bg)', minHeight: '100vh' }} />
        <Footer />
      </AppShell>
    )
  }

  // The rail and the header already drop the destination; this states why, so a
  // learner who reaches the URL directly is not looking at a page they believe
  // they turned off. The server has anonymised the rows either way.
  if (!showLeaderboards) {
    return (
      <AppShell>
        <div className="max-w-container mx-auto px-4 sm:px-6 pb-[calc(80px+env(safe-area-inset-bottom))] lg:pb-0 overflow-x-hidden min-w-0">
          <section className="pt-8 min-w-0 overflow-hidden">
            <span className="flex mb-4">
              <button type="button" onClick={() => { if (window.history.length > 1) window.history.back() }} className="text-mathua-secondary text-sm hover:text-mathua-blue">
                ← Back
              </button>
            </span>

            <SectionHeader label="Competition" title="Leaderboard" />

            <div className="border border-mathua-border bg-mathua-surface p-6 mt-6 max-w-[560px] mx-auto text-center">
              <p className="font-mono text-sm text-mathua-primary">Leaderboards are turned off</p>
              <p className="text-mathua-secondary text-sm mt-2">
                You are anonymous on the public ranking, so your name never appears on it. Change this in <Link href="/settings" className="text-mathua-blue hover:underline">Settings</Link>.
              </p>
            </div>
          </section>
        </div>

        <Footer />
      </AppShell>
    )
  }

  return (
    <AppShell>
      <div className="max-w-container mx-auto px-4 sm:px-6 pb-[calc(80px+env(safe-area-inset-bottom))] lg:pb-0 overflow-x-hidden min-w-0">
        <section className="pt-8 min-w-0 overflow-hidden">
          <span className="flex mb-4">
            <button type="button" onClick={() => { if (window.history.length > 1) window.history.back() }} className="text-mathua-secondary text-sm hover:text-mathua-blue">
              ← Back
            </button>
          </span>

          <SectionHeader label="Competition" title="Leaderboard" />

          <p className="text-mathua-secondary text-sm text-center mt-4 mb-6 px-2 max-w-[560px] mx-auto">
            XP comes from answered questions, not from time spent. The week runs Monday 00:00 to
            Monday 00:00 UTC, and the board resets to zero each week while lifetime XP keeps counting.
          </p>

          {tab === 'weekly' && (
            <div className="flex justify-center mb-6 px-2">
              <div className="bg-mathua-surface border border-mathua-border p-4 text-center w-full max-w-[260px] sm:min-w-[200px] sm:w-auto">
                <span className="font-mono text-[10px] uppercase text-mathua-muted">Resets in</span>
                <div className="font-mono text-2xl text-mathua-blue mt-1">{countdown}</div>
              </div>
            </div>
          )}

          <div role="tablist" aria-label="Leaderboard view" className="flex flex-wrap justify-center gap-1 border-b border-mathua-border mb-6">
            {TABS.map((t) => (
              <button
                key={t.id}
                type="button"
                role="tab"
                id={`tab-${t.id}`}
                aria-selected={tab === t.id}
                aria-controls={`panel-${t.id}`}
                onClick={() => selectTab(t.id)}
                className={`font-mono text-xs uppercase tracking-wider px-3 sm:px-4 py-2 border-b-2 -mb-px whitespace-nowrap ${
                  tab === t.id
                    ? 'border-mathua-blue text-mathua-blue'
                    : 'border-transparent text-mathua-muted hover:text-mathua-primary'
                }`}
              >
                {t.label}
              </button>
            ))}
          </div>
        </section>

        <div
          role="tabpanel"
          id={`panel-${tab}`}
          aria-labelledby={`tab-${tab}`}
          className="mb-12 min-w-0"
        >
          {loading && (
            <div className="w-[calc(100%+2rem)] -mx-4 px-4 sm:w-full sm:mx-0 sm:px-0 overflow-x-auto">
              <span role="status" aria-live="polite" className="sr-only">Loading leaderboard</span>
              <table className="w-full min-w-[320px] border-collapse bg-mathua-surface border border-mathua-border">
                <tbody><LeaderboardTableSkeleton rows={8} /></tbody>
              </table>
            </div>
          )}

          {!loading && failed && (
            <p className="text-center text-mathua-muted text-sm py-8">
              Could not load this board.{' '}
              <button type="button" onClick={retry} className="text-mathua-blue hover:underline">Retry</button>
            </p>
          )}

          {!loading && !failed && tab !== 'league' && board && (
            <>
              {board.page === 1 && <Podium entries={board.entries} />}
              {board.me && <PositionStrip me={board.me} view={tab} />}
              {board.entries.length === 0 ? (
                <p className="text-center text-mathua-muted text-sm py-8 px-2">
                  {tab === 'weekly'
                    ? 'No one has earned XP this week yet. Answer a question and the board starts with you.'
                    : 'No learner has earned XP yet. The first answer puts the first name here.'}
                  {' '}<Link href="/learn" className="text-mathua-blue hover:underline">Start learning</Link>
                </p>
              ) : (
                <>
                  <BoardTable board={board} view={tab} />
                  <PageControls board={board} onPage={changePage} busy={loading} />
                </>
              )}
            </>
          )}

          {!loading && !failed && tab === 'league' && (
            league === null ? (
              <p className="text-center text-mathua-muted text-sm py-8 px-2">
                <Link href="/login" className="text-mathua-blue hover:underline">Sign in</Link> to join a league.
              </p>
            ) : league.tier === '' ? (
              <p className="text-center text-mathua-muted text-sm py-8 px-2">
                You are not in a league yet. Answer one question this week and you join a Bronze
                league with room. The top five in each league move up a tier every Monday.
              </p>
            ) : (
              <>
                <div className="flex flex-wrap items-baseline justify-between gap-2 mb-4">
                  <h2 className="font-mono text-xs uppercase tracking-wider text-mathua-primary">
                    {league.tier} league {league.league_no}
                  </h2>
                  <p className="font-mono text-[11px] text-mathua-muted">
                    {shortWeek(league.week_start)} to {shortWeek(league.week_end)} · {league.size} learners
                  </p>
                </div>
                <p className="text-mathua-secondary text-sm mb-4">
                  Top {league.promote} move up a tier next Monday, bottom {league.demote} move down.
                  Ties on XP are ordered by name so positions are stable.
                </p>
                <div className="w-[calc(100%+2rem)] -mx-4 px-4 sm:w-full sm:mx-0 sm:px-0 overflow-x-auto">
                  <table className="w-full min-w-[320px] border-collapse bg-mathua-surface border border-mathua-border">
                    <caption className="sr-only">Your league standings this week</caption>
                    <thead>
                      <tr>
                        <th scope="col" className="font-mono text-[11px] uppercase text-mathua-muted text-left p-2.5 sm:p-[14px_20px] bg-mathua-surface-elevated border-b border-mathua-border w-10 sm:w-12">#</th>
                        <th scope="col" className="font-mono text-[11px] uppercase text-mathua-muted text-left p-2.5 sm:p-[14px_20px] bg-mathua-surface-elevated border-b border-mathua-border">Learner</th>
                        <th scope="col" className="font-mono text-[11px] uppercase text-mathua-muted text-right p-2.5 sm:p-[14px_20px] bg-mathua-surface-elevated border-b border-mathua-border">XP</th>
                      </tr>
                    </thead>
                    <tbody>
                      {league.members.map((m) => (
                        <tr key={m.id} className={`border-b border-mathua-border last:border-b-0 ${m.is_me ? 'bg-mathua-blue-faint' : ''}`}>
                          <td className="p-2.5 sm:p-[14px_20px] font-mono text-sm text-mathua-muted">{m.rank}</td>
                          <td className="p-2.5 sm:p-[14px_20px] text-mathua-primary text-sm font-medium max-w-[160px] sm:max-w-none">
                            <span className="flex items-center gap-2 min-w-0">
                              <LeagueAvatar tier={league.tier} weeklyXp={m.xp} seed={m.id} name={nameOf(m)} size={26} url={avatarURL(m)} className="shrink-0" />
                              <span className="block truncate" title={nameOf(m)}>{nameOf(m)}</span>
                              {m.is_me && <span className="font-mono text-[10px] uppercase text-mathua-blue shrink-0">you</span>}
                            </span>
                          </td>
                          <td className="p-2.5 sm:p-[14px_20px] font-mono text-sm text-mathua-green text-right whitespace-nowrap">{m.xp}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </>
            )
          )}
        </div>
      </div>

      <Footer />
    </AppShell>
  )
}
