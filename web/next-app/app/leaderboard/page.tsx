'use client'

import { useTheme } from '../../hooks/useTheme'
import { useState, useEffect, useCallback } from 'react'
import Link from 'next/link'
import Header from '../../components/Header'
import BottomTabs from '../../components/BottomTabs'
import SectionHeader from '../../components/SectionHeader'
import Footer from '../../components/Footer'
import FormulaBlock from '../../components/FormulaBlock'
import ProgressionLevels from '../../components/ProgressionLevels'
import { getLeaderboard, getLeagues, API_BASE, type LeaderboardEntry, type LeagueBoard } from '../../lib/api'
import { resolveAvatar } from '../../lib/dicebear'
import Avatar from '../../components/Avatar'
import Loading from '../../components/Loading'

// Display chain: registered name → random/claimed username → Anonymous.
// Enforcement (required + trimmed at signup, service-level cap) should make
// the last link dead code, but the board must never render a blank cell.
function displayName(row: { name: string; username?: string }): string {
  if (row.name && row.name.trim()) return row.name
  if (row.username && row.username.trim()) return row.username
  return 'Anonymous'
}

interface BoardAvatar {
  student_id: string
  avatar_url?: string
  avatar_dicebear?: { style: string; seed: string } | null
  avatar_custom?: boolean
  avatar_version?: number
}

// Avatar precedence mirrors resolveAvatar, except a custom upload resolves
// to the public per-student photo endpoint (board-visible by design).
function boardAvatarUrl(row: BoardAvatar): string | undefined {
  if (row.avatar_custom && row.student_id) {
    const v = row.avatar_version ? `?v=${row.avatar_version}` : ''
    return `${API_BASE}/api/avatar/${row.student_id}${v}`
  }
  return resolveAvatar(
    { avatar_url: row.avatar_url },
    { avatar_dicebear: row.avatar_dicebear ?? null },
  ).url
}

const LEVELS = [  { num: '01', name: 'Novice', range: '0–31' },
  { num: '02', name: 'Apprentice', range: '32–63' },
  { num: '03', name: 'Student', range: '64–95' },
  { num: '04', name: 'Scholar', range: '96–127' },
  { num: '05', name: 'Adept', range: '128–159' },
  { num: '06', name: 'Expert', range: '160–191' },
  { num: '07', name: 'Master', range: '192–223' },
  { num: '08', name: 'Grandmaster', range: '224–255' },
  { num: '09', name: 'Math Architect', range: '256–284', elite: true },
]

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

export default function LeaderboardPage() {
  const { mounted } = useTheme()
  // Lazy initializer derives the first value during render (no init effect).
  const [countdown, setCountdown] = useState(computeCountdown)
  const [entries, setEntries] = useState<LeaderboardEntry[]>([])
  const [leagues, setLeagues] = useState<LeagueBoard | null>(null)
  const [leaguesFailed, setLeaguesFailed] = useState(false)
  const [loadError, setLoadError] = useState(false)
  const [loading, setLoading] = useState(true)

  const finishLoading = useCallback((data: LeaderboardEntry[]) => {
    setEntries(data)
    setLoading(false)
  }, [])

  useEffect(() => {
    const update = () => {
      setCountdown(computeCountdown())
    }
    update()
    const interval = setInterval(update, 60000)
    return () => clearInterval(interval)
  }, [])

  useEffect(() => {
    getLeaderboard()
      .then((data) => { setLoadError(false); finishLoading(data) })
      .catch((e) => { console.error('leaderboard fetch failed:', e); setLoadError(true); setLoading(false) })
    getLeagues()
      .then(setLeagues)
      .catch((e) => { console.error('leagues fetch failed:', e); setLeaguesFailed(true) })
  }, [finishLoading])

  if (!mounted) return <div style={{ background: 'var(--bg)', minHeight: '100vh' }} />

  return (
    <>
      <Header />
      <div className="max-w-container mx-auto px-4 sm:px-6 pb-[calc(80px+env(safe-area-inset-bottom))] lg:pb-0 overflow-x-hidden min-w-0">
      <section className="pt-8 min-w-0 overflow-hidden">
        <span className="flex mb-4">
          <button type="button" onClick={() => { if (window.history.length > 1) window.history.back() }} className="text-mathua-secondary text-sm hover:text-mathua-primary">
            ← Back
          </button>
        </span>

        <SectionHeader label="Weekly Leaderboard" title="Weekly rankings" />

        <div className="flex justify-center mt-4 mb-6 sm:mb-10 px-2">
          <div className="bg-mathua-surface border border-mathua-border rounded-none p-4 text-center w-full max-w-[260px] sm:min-w-[200px] sm:w-auto">
            <span className="font-mono text-[10px] uppercase text-mathua-muted">
              Resets in
            </span>
            <div className="font-mono text-2xl text-mathua-blue mt-1">{countdown}</div>
          </div>
        </div>

        <p className="text-mathua-secondary text-sm text-center mb-6 px-2">
          Scores reset every Monday at 00:00 UTC.
        </p>

        <div className="flex justify-center mb-6 sm:mb-10 w-full max-w-full min-w-0 overflow-hidden px-2">
          <div className="w-full max-w-full min-w-0 overflow-hidden flex justify-center">
            <FormulaBlock
              code={`weekly_score = (concepts_mastered_this_week * 100) + speed_bonus + (day_streak * 10)`}
            />
          </div>
        </div>
      </section>

      <div className="mb-12 w-full max-w-full min-w-0 overflow-hidden">
        <div className="w-[calc(100%+2rem)] -mx-4 px-4 sm:w-full sm:mx-0 sm:px-0 overflow-x-auto overscroll-x-contain">
          <table className="w-full min-w-[320px] border-collapse bg-mathua-surface rounded-none overflow-hidden border border-mathua-border">
            <thead>
              <tr>
                <th className="font-mono text-[11px] uppercase text-mathua-muted text-left p-2.5 sm:p-[14px_20px] bg-mathua-surface-elevated border-b border-mathua-border w-10 sm:w-12">#</th>
                <th className="font-mono text-[11px] uppercase text-mathua-muted text-left p-2.5 sm:p-[14px_20px] bg-mathua-surface-elevated border-b border-mathua-border">Name</th>
                <th className="font-mono text-[11px] uppercase text-mathua-muted text-right p-2.5 sm:p-[14px_20px] bg-mathua-surface-elevated border-b border-mathua-border">Mastered</th>
                <th className="font-mono text-[11px] uppercase text-mathua-muted text-right p-2.5 sm:p-[14px_20px] bg-mathua-surface-elevated border-b border-mathua-border">Score</th>
                <th className="font-mono text-[11px] uppercase text-mathua-muted text-right p-2.5 sm:p-[14px_20px] bg-mathua-surface-elevated border-b border-mathua-border hidden md:table-cell">Level</th>
              </tr>
            </thead>
            <tbody>
              {loading && (
                <tr>
                  <td colSpan={5} className="p-8 text-center text-mathua-muted text-sm">
                    <Loading label="LOADING LEADERBOARD" />
                  </td>
                </tr>
              )}
              {!loading && loadError && (
                <tr>
                  <td colSpan={5} className="p-8 text-center text-mathua-muted text-sm">
                    Couldn&apos;t load the leaderboard. <button type="button" onClick={() => window.location.reload()} className="text-mathua-blue hover:underline">Retry</button>
                  </td>
                </tr>
              )}
              {!loading && !loadError && entries.length === 0 && (
                <tr>
                  <td colSpan={5} className="p-8 text-center text-mathua-muted text-sm">
                    No data yet.
                  </td>
                </tr>
              )}
              {entries.map((row, i) => (
                <tr
                  key={`entry-${row.rank}`}
                  className={`border-b border-mathua-border last:border-b-0 ${
                    i < 3 ? 'bg-mathua-blue/5' : ''
                  }`}
                >
                  <td className="p-2.5 sm:p-[14px_20px] font-mono text-sm text-mathua-muted">
                    {row.rank <= 3 ? (
                      <span className="text-mathua-blue">
                        {['\u2460', '\u2461', '\u2462'][row.rank - 1]}
                      </span>
                    ) : row.rank}
                  </td>
                  <td className="p-2.5 sm:p-[14px_20px] text-mathua-primary text-sm font-medium max-w-[110px] sm:max-w-none">
                    <span className="flex items-center gap-2 min-w-0">
                      <Avatar seed={row.student_id || displayName(row)} name={displayName(row)} size={28} url={boardAvatarUrl(row)} className="shrink-0" />
                      <span className="block truncate" title={displayName(row)}>{displayName(row)}</span>
                    </span>
                  </td>
                  <td className="p-2.5 sm:p-[14px_20px] font-mono text-sm text-mathua-muted text-right whitespace-nowrap">
                    {row.mastered}
                  </td>
                  <td className="p-2.5 sm:p-[14px_20px] font-mono text-sm text-mathua-green text-right whitespace-nowrap">
                    {row.score}
                  </td>
                  <td className="p-2.5 sm:p-[14px_20px] text-mathua-secondary text-sm text-right hidden md:table-cell truncate max-w-[110px]">
                    {row.level}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>

      <section className="mt-12 mb-12 min-w-0 overflow-hidden">
        <SectionHeader label="Leagues" title="Weekly promotion & demotion" />
        <p className="text-mathua-secondary text-sm text-center mb-6 px-2">
          Top 2 in each league promote each Monday; the bottom 2 demote.
        </p>
        {!leagues && !leaguesFailed ? (
          <div className="text-center py-8">
            <Loading label="LOADING LEAGUES" />
          </div>
        ) : leaguesFailed ? (
          <p className="text-center text-mathua-muted text-sm">Couldn&apos;t load leagues. <button type="button" onClick={() => window.location.reload()} className="text-mathua-blue hover:underline">Retry</button></p>
        ) : leagues!.leagues.length === 0 ? (
          <p className="text-center text-mathua-muted text-sm"><Link href="/login" className="text-mathua-blue hover:text-mathua-blue-hover">Sign in</Link> to join a league.</p>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4 w-full max-w-full min-w-0 overflow-hidden">
            {leagues!.leagues.map((lg) => (
              <div key={lg.tier} className="border border-mathua-border bg-mathua-surface w-full max-w-full min-w-0 overflow-hidden">
                <div className="flex items-center justify-between px-4 py-3 bg-mathua-surface-elevated border-b border-mathua-border">
                  <span className="font-mono text-xs uppercase tracking-wider text-mathua-primary">
                    {lg.tier}
                  </span>
                  <span className="font-mono text-[10px] text-mathua-muted">
                    {lg.members.length} members
                  </span>
                </div>
                <div className="divide-y divide-mathua-border">
                  {lg.members.map((m, i) => (
                    <div key={m.student_id} className="flex items-center gap-3 px-4 py-2 min-w-0">
                      <span className="font-mono text-xs text-mathua-muted w-6 shrink-0">{i + 1}</span>
                      <Avatar seed={m.student_id || displayName(m)} name={displayName(m)} size={24} url={boardAvatarUrl(m)} className="shrink-0" />
                      <span className="font-mono text-xs text-mathua-primary truncate flex-1 min-w-0" title={displayName(m)}>
                        {displayName(m)}
                      </span>
                      {m.moved === 1 && (
                        <span className="font-mono text-[10px] text-mathua-green shrink-0">▲ promoted</span>
                      )}
                      {m.moved === -1 && (
                        <span className="font-mono text-[10px] text-mathua-red shrink-0">▼ demoted</span>
                      )}
                      <span className="font-mono text-[10px] text-mathua-muted shrink-0">
                        {m.weekly_mastered} this week
                      </span>
                    </div>
                  ))}
                </div>
              </div>
            ))}
          </div>
        )}
      </section>

      <section className="py-12 sm:py-20 min-w-0 overflow-hidden">
        <SectionHeader label="Progression" title="Lifetime rank" />
        <div className="flex justify-center mt-6 w-full max-w-full min-w-0 overflow-hidden px-2">
          <div className="w-full max-w-[400px] min-w-0">
            <ProgressionLevels levels={LEVELS} />
          </div>
        </div>
      </section>

      <Footer />
    </div>
      <BottomTabs />
    </>
  )
}
