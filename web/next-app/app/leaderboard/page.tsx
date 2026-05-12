'use client'

import { useTheme } from '../../hooks/useTheme'
import { useState, useEffect } from 'react'
import Header from '../../components/Header'
import SectionHeader from '../../components/SectionHeader'
import Footer from '../../components/Footer'
import FormulaBlock from '../../components/FormulaBlock'
import ProgressionLevels from '../../components/ProgressionLevels'
import { getLeaderboard, type LeaderboardEntry } from '../../lib/api'

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
      .then((data) => { setEntries(data); setLoading(false) })
      .catch(() => setLoading(false))
  }, [])

  if (!mounted) return <div style={{ background: 'var(--bg)', minHeight: '100vh' }} />

  return (
    <>
      <Header />
      <div className="max-w-container mx-auto px-6 max-sm:px-4">
      <section className="pt-8">
        <span className="flex mb-4">
          <a href="/" className="text-mathua-secondary text-sm hover:text-mathua-primary">
            ← Back
          </a>
        </span>

        <SectionHeader label="Weekly Leaderboard" title="Compete. Improve. Rise." />

        <div className="flex justify-center mt-4 mb-10">
          <div className="bg-mathua-surface border border-mathua-border rounded-lg p-4 text-center min-w-[200px]">
            <span className="font-mono text-[10px] uppercase text-mathua-muted">
              Resets in
            </span>
            <div className="font-mono text-2xl text-mathua-blue mt-1">{countdown}</div>
          </div>
        </div>

        <p className="text-mathua-secondary text-sm text-center mb-6">
          Start a practice session to appear on the leaderboard. Scores reset every Monday.
        </p>

        <div className="flex justify-center mb-10">
          <FormulaBlock
            code={`weekly_score = (concepts_mastered_this_week * 100) + speed_bonus + (day_streak * 10)`}
          />
        </div>
      </section>

      <div className="mb-12">
        <table className="w-full border-collapse bg-mathua-surface rounded-lg overflow-hidden border border-mathua-border">
          <thead>
            <tr>
              <th className="font-mono text-[11px] uppercase text-mathua-muted text-left p-[14px_20px] max-sm:p-[10px_12px] bg-mathua-surface-elevated border-b border-mathua-border w-12">#</th>
              <th className="font-mono text-[11px] uppercase text-mathua-muted text-left p-[14px_20px] max-sm:p-[10px_12px] bg-mathua-surface-elevated border-b border-mathua-border">Name</th>
              <th className="font-mono text-[11px] uppercase text-mathua-muted text-right p-[14px_20px] max-sm:p-[10px_12px] bg-mathua-surface-elevated border-b border-mathua-border">Mastered</th>
              <th className="font-mono text-[11px] uppercase text-mathua-muted text-right p-[14px_20px] max-sm:p-[10px_12px] bg-mathua-surface-elevated border-b border-mathua-border">Score</th>
              <th className="font-mono text-[11px] uppercase text-mathua-muted text-right p-[14px_20px] max-sm:p-[10px_12px] bg-mathua-surface-elevated border-b border-mathua-border hidden md:table-cell">Level</th>
            </tr>
          </thead>
          <tbody>
            {loading && (
              <tr>
                <td colSpan={5} className="p-8 text-center text-mathua-muted text-sm">
                  Loading leaderboard...
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
                key={i}
                className={`border-b border-mathua-border last:border-b-0 ${
                  i < 3 ? 'bg-mathua-gold/5' : ''
                }`}
              >
                <td className="p-[14px_20px] max-sm:p-[10px_12px] font-mono text-sm text-mathua-muted">
                  {row.rank <= 3 ? (
                    <span className="text-mathua-gold">
                      {['\u2460', '\u2461', '\u2462'][row.rank - 1]}
                    </span>
                  ) : row.rank}
                </td>
                <td className="p-[14px_20px] max-sm:p-[10px_12px] text-mathua-primary text-sm font-medium">
                  {row.name}
                </td>
                <td className="p-[14px_20px] max-sm:p-[10px_12px] font-mono text-sm text-mathua-muted text-right">
                  {row.mastered}
                </td>
                <td className="p-[14px_20px] max-sm:p-[10px_12px] font-mono text-sm text-mathua-green text-right">
                  {row.score}
                </td>
                <td className="p-[14px_20px] max-sm:p-[10px_12px] text-mathua-secondary text-sm text-right hidden md:table-cell">
                  {row.level}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      <section className="py-20 max-sm:py-12">
        <SectionHeader label="Progression" title="Your permanent rank" />
        <div className="flex justify-center mt-6">
          <div className="w-full max-w-[400px]">
            <ProgressionLevels levels={LEVELS} />
          </div>
        </div>
        <p className="text-center text-mathua-muted text-[13px] italic mt-4">
          Your level is permanent and never decreases.
        </p>
      </section>

      <Footer />
    </div>
    </>
  )
}
