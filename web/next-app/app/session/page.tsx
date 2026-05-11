'use client'

import { useState, useEffect } from 'react'
import SectionHeader from '../../components/SectionHeader'
import Footer from '../../components/Footer'
import AsciiDivider from '../../components/AsciiDivider'
import Pipeline from '../../components/Pipeline'

export default function SessionPage() {
  const [theme, setTheme] = useState<'dark' | 'light'>('dark')
  const [mounted, setMounted] = useState(false)

  useEffect(() => {
    const saved = localStorage.getItem('mathua-theme')
    setTheme(saved === 'light' || saved === 'dark' ? saved : 'dark')
    setMounted(true)
  }, [])

  useEffect(() => {
    if (!mounted) return
    document.documentElement.classList.toggle('dark', theme === 'dark')
  }, [theme, mounted])

  const toggleTheme = () => {
    const next = theme === 'dark' ? 'light' : 'dark'
    setTheme(next)
    localStorage.setItem('mathua-theme', next)
  }

  if (!mounted) return <div style={{ background: '#0a0f1a', minHeight: '100vh' }} />

  return (
    <div className="max-w-container mx-auto px-6 max-sm:px-4">
      <button
        onClick={toggleTheme}
        className="fixed top-[50px] right-5 z-[1000] bg-mathua-surface border border-mathua-border-strong text-mathua-primary px-3.5 py-2 rounded-md font-mono text-xs cursor-pointer transition-all duration-200 hover:border-mathua-blue hover:text-mathua-blue"
        aria-label="Toggle theme"
      >
        {theme === 'dark' ? '\u2600' : '\u263E'}
      </button>

      <section className="pt-8">
        <span className="flex justify-between items-center mb-4">
          <a href="/" className="text-mathua-secondary text-sm hover:text-mathua-primary">
            ← Back
          </a>
          <span className="font-mono text-[11px] text-mathua-muted">
            Session #1
          </span>
        </span>

        <SectionHeader label="Practice session" title="Build your mastery" />

        <div className="flex gap-8 items-start mt-10 max-md:flex-col">
          {/* Session Stats */}
          <div className="w-[240px] flex-shrink-0 max-md:w-full">
            <div className="bg-mathua-surface border border-mathua-border rounded-lg p-5 space-y-4">
              <div>
                <span className="font-mono text-[10px] uppercase tracking-[0.1em] text-mathua-muted">
                  Streak
                </span>
                <div className="font-mono text-2xl text-mathua-green mt-1">3</div>
              </div>
              <div>
                <span className="font-mono text-[10px] uppercase tracking-[0.1em] text-mathua-muted">
                  Time
                </span>
                <div className="font-mono text-2xl text-mathua-primary mt-1">12s</div>
              </div>
              <div>
                <span className="font-mono text-[10px] uppercase tracking-[0.1em] text-mathua-muted">
                  Mastered today
                </span>
                <div className="font-mono text-2xl text-mathua-gold mt-1">4</div>
              </div>
              <div>
                <span className="font-mono text-[10px] uppercase tracking-[0.1em] text-mathua-muted">
                  Progress
                </span>
                <pre className="font-mono text-xs text-mathua-blue mt-2 select-none whitespace-pre">
                  {'[████████░░░░░░░░░░░░]'}
                </pre>
              </div>
            </div>
          </div>

          {/* Main area */}
          <div className="flex-1">
            <div className="bg-mathua-surface border border-mathua-border rounded-lg p-8 max-md:p-5">
              <div className="mb-6">
                <span className="section-label">arith.add.single</span>
                <h3 className="font-serif text-2xl font-medium text-mathua-primary mt-1">
                  Single-digit addition
                </h3>
              </div>

              <div className="bg-mathua-code border border-mathua-border rounded-md p-8 text-center mb-6">
                <p className="text-mathua-primary text-4xl font-mono font-light mb-4">
                  7 + 5 = ?
                </p>
              </div>

              <div className="flex gap-3 mb-6">
                <input
                  type="text"
                  placeholder="Your answer"
                  className="flex-1 bg-mathua-code border border-mathua-border rounded-md h-12 px-4 font-mono text-base text-mathua-primary placeholder:text-mathua-muted focus:outline-none focus:border-mathua-blue transition-colors"
                />
                <button className="bg-mathua-blue text-white hover:bg-mathua-blue-hover rounded-md h-12 px-8 font-medium text-sm transition-colors">
                  Submit
                </button>
              </div>

              <div className="space-y-1">
                <button className="w-full bg-mathua-surface-elevated border border-mathua-border rounded-md h-10 text-sm text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue transition-colors">
                  I gave up — show solution
                </button>
                <button className="w-full bg-mathua-surface-elevated border border-mathua-border rounded-md h-10 text-sm text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue transition-colors">
                  Skip this concept
                </button>
              </div>
            </div>
          </div>
        </div>
      </section>

      <AsciiDivider pattern="wave" />

      <section className="py-20 max-sm:py-12">
        <SectionHeader label="Your state" title="Mastery pipeline" />
        <div className="flex justify-center mt-6">
          <Pipeline
            states={[
              { label: 'UNSEEN', status: 'unseen' },
              { label: 'LEARNING', status: 'learning' },
              { label: 'PRACTICING', status: 'practicing' },
              { label: 'MASTERED', status: 'mastered' },
              { label: 'DECAYING', status: 'decaying' },
            ]}
          />
        </div>
      </section>

      <Footer />
    </div>
  )
}
