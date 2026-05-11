'use client'

import { useState, useEffect } from 'react'
import SectionHeader from '../../components/SectionHeader'
import Footer from '../../components/Footer'
import FormulaBlock from '../../components/FormulaBlock'
import ProgressionLevels from '../../components/ProgressionLevels'

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

const MOCK_LEADERBOARD = [
  { rank: 1, name: 'alice', mastered: 247, streak: 14, level: 'Grandmaster' },
  { rank: 2, name: 'bob', mastered: 238, streak: 9, level: 'Grandmaster' },
  { rank: 3, name: 'carol', mastered: 221, streak: 21, level: 'Master' },
  { rank: 4, name: 'dave', mastered: 198, streak: 7, level: 'Master' },
  { rank: 5, name: 'eve', mastered: 185, streak: 12, level: 'Expert' },
  { rank: 6, name: 'frank', mastered: 172, streak: 5, level: 'Expert' },
  { rank: 7, name: 'grace', mastered: 164, streak: 18, level: 'Expert' },
  { rank: 8, name: 'heidi', mastered: 151, streak: 3, level: 'Adept' },
  { rank: 9, name: 'ivan', mastered: 143, streak: 6, level: 'Adept' },
  { rank: 10, name: 'judy', mastered: 128, streak: 1, level: 'Adept' },
]

export default function LeaderboardPage() {
  const [theme, setTheme] = useState<'dark' | 'light'>('dark')
  const [mounted, setMounted] = useState(false)
  const [countdown, setCountdown] = useState('')

  useEffect(() => {
    const saved = localStorage.getItem('mathua-theme')
    setTheme(saved === 'light' || saved === 'dark' ? saved : 'dark')
    setMounted(true)
  }, [])

  useEffect(() => {
    if (!mounted) return
    document.documentElement.classList.toggle('dark', theme === 'dark')
  }, [theme, mounted])

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

  const toggleTheme = () => {
    const next = theme === 'dark' ? 'light' : 'dark'
    setTheme(next)
    localStorage.setItem('mathua-theme', next)
  }

  if (!mounted) return <div style={{ background: 'var(--bg)', minHeight: '100vh' }} />

  return (
    <div className="max-w-container mx-auto px-6 max-sm:px-4">
      <button
        onClick={toggleTheme}
        className="fixed top-[50px] right-5 z-[1000] border border-[var(--border-strong)] text-[var(--text-muted)] px-3 py-1.5 font-mono text-xs cursor-pointer transition-all duration-200 hover:text-[var(--accent-gold)] hover:border-[var(--accent-gold)]"
        style={{ borderRadius: 0, background: 'var(--bg)' }}
        aria-label="Toggle theme"
      >
        {theme === 'dark' ? '\u2600' : '\u263E'}
      </button>

      <section className="pt-8">
        <span className="flex mb-4">
          <a href="/" className="text-mathua-secondary text-sm hover:text-mathua-primary">
            ← Back
          </a>
        </span>

        <SectionHeader label="Weekly Leaderboard" title="Compete. Improve. Rise." />

        <div className="flex justify-center mt-4 mb-10">
          <div className="bg-mathua-surface border border-mathua-border rounded-lg p-4 text-center min-w-[200px]">
            <span className="font-mono text-[10px] uppercase tracking-[0.1em] text-mathua-muted">
              Resets in
            </span>
            <div className="font-mono text-2xl text-mathua-blue mt-1">{countdown}</div>
          </div>
        </div>

        <p className="text-mathua-secondary text-sm leading-relaxed text-center mb-6">
          The leaderboard resets every Monday at 00:00 UTC. Your score is calculated from three
          components:
        </p>

        <div className="flex justify-center mb-10">
          <FormulaBlock
            code={`weekly_score = (concepts_mastered_this_week × 100)
            + speed_bonus
            + (current_day_streak × 10)`}
          />
        </div>
      </section>

      {/* Leaderboard Table */}
      <div className="mb-12">
        <table className="w-full border-collapse bg-mathua-surface rounded-lg overflow-hidden border border-mathua-border">
          <thead>
            <tr>
              <th className="font-mono text-[11px] uppercase tracking-[0.1em] text-mathua-muted text-left p-[14px_20px] bg-mathua-surface-elevated border-b border-mathua-border w-12">
                #
              </th>
              <th className="font-mono text-[11px] uppercase tracking-[0.1em] text-mathua-muted text-left p-[14px_20px] bg-mathua-surface-elevated border-b border-mathua-border">
                Name
              </th>
              <th className="font-mono text-[11px] uppercase tracking-[0.1em] text-mathua-muted text-right p-[14px_20px] bg-mathua-surface-elevated border-b border-mathua-border">
                Mastered
              </th>
              <th className="font-mono text-[11px] uppercase tracking-[0.1em] text-mathua-muted text-right p-[14px_20px] bg-mathua-surface-elevated border-b border-mathua-border">
                Streak
              </th>
              <th className="font-mono text-[11px] uppercase tracking-[0.1em] text-mathua-muted text-right p-[14px_20px] bg-mathua-surface-elevated border-b border-mathua-border hidden md:table-cell">
                Level
              </th>
            </tr>
          </thead>
          <tbody>
            {MOCK_LEADERBOARD.map((row, i) => (
              <tr
                key={row.rank}
                className={`border-b border-mathua-border last:border-b-0 ${
                  i < 3 ? 'bg-mathua-gold/5' : ''
                }`}
              >
                <td className="p-[14px_20px] font-mono text-sm text-mathua-muted">
                  {row.rank <= 3 ? (
                    <span className="text-mathua-gold">
                      {['\u2460', '\u2461', '\u2462'][row.rank - 1]}
                    </span>
                  ) : (
                    row.rank
                  )}
                </td>
                <td className="p-[14px_20px] text-mathua-primary text-sm font-medium">
                  {row.name}
                </td>
                <td className="p-[14px_20px] font-mono text-sm text-mathua-muted text-right">
                  {row.mastered}
                </td>
                <td className="p-[14px_20px] font-mono text-sm text-mathua-green text-right">
                  {row.streak}
                </td>
                <td className="p-[14px_20px] text-mathua-secondary text-sm text-right hidden md:table-cell">
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
          Your level is permanent and never decreases. A bad week does not undo months of learning.
        </p>
      </section>

      <Footer />
    </div>
  )
}
