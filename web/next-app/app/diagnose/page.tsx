'use client'

import { useState, useEffect } from 'react'
import SectionHeader from '../../components/SectionHeader'
import Footer from '../../components/Footer'
import AsciiDivider from '../../components/AsciiDivider'

export default function DiagnosePage() {
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
        <span className="flex mb-4">
          <a href="/" className="text-mathua-secondary text-sm hover:text-mathua-primary">
            ← Back
          </a>
        </span>

        <SectionHeader label="Diagnostic Test" title="Find your knowledge frontier" />

        <p className="text-mathua-secondary text-sm leading-relaxed text-center max-w-[600px] mx-auto mt-4">
          A Computerised Adaptive Testing session that locates your position in the concept graph
          using as few questions as possible — typically 20–35 questions across the 284-concept
          space.
        </p>
      </section>

      <AsciiDivider pattern="wave" />

      {/* How it works */}
      <section className="py-20 max-sm:py-12">
        <SectionHeader label="Method" title="How the diagnostic works" />

        <div className="flex justify-center mt-8">
          <pre className="bg-mathua-code border border-mathua-border rounded-md p-6 font-mono text-xs text-mathua-blue whitespace-pre overflow-x-auto leading-relaxed inline-block">
            {`┌──────────────────────────────────────────────┐
│          TOPOLOGICALLY SORTED GRAPH           │
│   Node 0 ──→ Node 1 ──→ ... ──→ Node N     │
│   ^                         ^                 │
│   │      Binary search       │                │
│   │       finds frontier      │                │
│   └───────────────────────────┘               │
│                                                │
│   CORRECT under time  →  move forward          │
│   INCORRECT / slow    →  move backward         │
│   3 consecutive hits   →  FRONTIER LOCKED      │
└──────────────────────────────────────────────┘`}
          </pre>
        </div>

        <ul className="max-w-[600px] mx-auto mt-8 space-y-3 list-none">
          {[
            'The concept graph is sorted topologically — the diagnostic starts at the midpoint.',
            'Correct answers within the time limit move the probe forward toward harder concepts.',
            'Incorrect or slow answers move backward toward foundational material.',
            'After 3 consecutive correct answers in a region, the frontier is considered located.',
            'The diagnostic records a mastery estimate for every concept passed through.',
          ].map((step, i) => (
            <li
              key={i}
              className="text-mathua-secondary text-[0.95rem] leading-[1.7] pl-9 relative before:content-[counter(step)] before:absolute before:left-0 before:text-mathua-blue before:font-mono before:text-[13px]"
              style={{ counterIncrement: 'step-counter 1' }}
            >
              {step}
            </li>
          ))}
        </ul>
        <style>{`ul { counter-reset: step-counter; }`}</style>
      </section>

      <AsciiDivider pattern="dash" />

      {/* Actions */}
      <section className="py-20 max-sm:py-12 text-center">
        <SectionHeader label="Ready" title="Begin your diagnostic" />

        <p className="text-mathua-secondary text-sm leading-relaxed max-w-[500px] mx-auto mt-4 mb-8">
          The diagnostic takes 20–35 questions. You can retake it at any time — retaking never
          deletes progress, and Mathua always keeps the most optimistic estimate.
        </p>

        <div className="flex gap-3 justify-center max-sm:flex-col max-sm:items-center">
          <button className="bg-mathua-blue text-white hover:bg-mathua-blue-hover rounded-md h-10 px-6 text-[13px] font-medium transition-colors">
            Start Diagnostic
          </button>
          <button className="border border-mathua-border-strong text-mathua-primary hover:border-mathua-blue hover:text-mathua-blue rounded-md h-10 px-6 text-[13px] font-medium transition-colors">
            Read the full algorithm
          </button>
        </div>
      </section>

      <Footer />
    </div>
  )
}
