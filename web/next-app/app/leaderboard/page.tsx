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
import { getLeaderboard, type LeaderboardEntry } from '../../lib/api'
import Loading from '../../components/Loading'

const LEVELS = [
  { num: '01', name: 'Novice', range: '0–31' },
  { num: '02', name: 'Apprentice', range: '32–63' },
  { num: '03', name: 'Student', range: '64–95' },
  { num: '04', name: 'Scholar', range: '96–127' },
  { num: '05', name: 'Adept', range: '128–159' },
  { num: '06', name: 'Expert', range: '160–191' },
  { num: '07', name: 'Master', range: '192–223' },
  { num: '08', name: 'Grandmaster', range: '224–255' },
  { num: '09', name: 'Math Architect', range: '256–284', elite: true },
]

export default function LeaderboardPage() {
  const { mounted } = useTheme()
  const [countdown, setCountdown] = useState('')
  const [entries, setEntries] = useState<LeaderboardEntry[]>([])
  const [loading, setLoading] = useState(true)

  const finishLoading = useCallback((data: LeaderboardEntry[]) => {
    setEntries(data)
    setLoading(false)
  }, [])

  useEffect(() => {
    const update = () => {
      const now = new Date()
      const nextMonday = new Date(now)
      nextMonday.setUTCDate(now.getUTCDate() + ((7 - now.getUTCDay() + 1) % 7 || 7))
      nextMonday.setUTCHours(0, 0, 0, 0)
      const diff = nextMonday.getTime() - now.getTime()
      const d = Math.floor(diff / 86400000)
      const h = Math.floor((diff % 86400000) / 3600000)
      const m = Math.floor((diff % 3600000) / 60000)
      setCountdown(`${d}d ${h}h ${m}m`)
    }
    update()
    const interval = setInterval(update, 60000)
    return () => clearInterval(interval)
  }, [])

  useEffect(() => {
    getLeaderboard()
      .then((data) => finishLoading(data))
      .catch((e) => { console.error('leaderboard fetch failed:', e); setLoading(false) })
  }, [finishLoading])

  if (!mounted) return <div style={{ background: 'var(--bg)', minHeight: '100vh' }} />

  return (
    <>
      <Header />
      <div className="max-w-container mx-auto px-4 sm:px-6 pb-[calc(80px+env(safe-area-inset-bottom))] lg:pb-0 overflow-x-hidden min-w-0">
      <section className="pt-8 min-w-0 overflow-hidden">
        <span className="flex mb-4">
          <Link href="/" className="text-mathua-secondary text-sm hover:text-mathua-primary">
            ← Back
          </Link>
        </span>

        <SectionHeader label="Weekly Leaderboard" title="Compete. Improve. Rise." />

        <div className="flex justify-center mt-4 mb-6 sm:mb-10 px-2">
          <div className="bg-mathua-surface border border-mathua-border rounded-none p-4 text-center w-full max-w-[260px] sm:min-w-[200px] sm:w-auto">
            <span className="font-mono text-[10px] uppercase text-mathua-muted">
              Resets in
            </span>
            <div className="font-mono text-2xl text-mathua-blue mt-1">{countdown}</div>
          </div>
        </div>

        <p className="text-mathua-secondary text-sm text-center mb-6 px-2">
          Start a practice session to appear on the leaderboard. Scores reset every Monday.
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
        <div className="w-full max-w-full overflow-x-auto overscroll-x-contain -mx-4 px-4 sm:mx-0 sm:px-0">
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
              {!loading && entries.length === 0 && (
                <tr>
                  <td colSpan={5} className="p-8 text-center text-mathua-muted text-sm">
                    No data yet. Start a practice session to appear here!
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
                  <td className="p-2.5 sm:p-[14px_20px] text-mathua-primary text-sm font-medium max-w-[110px] sm:max-w-none truncate">
                    <span className="block truncate" title={row.name}>{row.name}</span>
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

      <section className="py-12 sm:py-20 min-w-0 overflow-hidden">
        <SectionHeader label="Progression" title="Your permanent rank" />
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
