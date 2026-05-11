'use client'

import dynamic from 'next/dynamic'
import { useState, useEffect } from 'react'
import AsciiBanner from '../components/AsciiBanner'
import AsciiDivider from '../components/AsciiDivider'
import SectionHeader from '../components/SectionHeader'
import InfoCard from '../components/InfoCard'
import Pipeline from '../components/Pipeline'
import FormulaBlock from '../components/FormulaBlock'
import DomainTable from '../components/DomainTable'
import ProgressionLevels from '../components/ProgressionLevels'
import Button from '../components/Button'
import Footer from '../components/Footer'
import conceptsData from '../data/concepts.json'

const MathConceptGraph3D = dynamic(() => import('../components/MathConceptGraph3D'), {
  ssr: false,
  loading: () => (
    <div
      style={{
        height: '520px',
        background: 'var(--surface)',
        borderRadius: '8px',
        border: '0.5px solid var(--border)',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        color: 'var(--text-muted)',
        fontSize: '14px',
      }}
    >
      Loading graph...
    </div>
  ),
})

const ASCII_HERO = `      ┌─────────────────────────────────────────┐
      │     ╱╲  ╱╲  ╱╲    ╱╲  ╱╲  ╱╲          │
      │    ╱  ╲╱  ╲╱  ╲  ╱  ╲╱  ╲╱  ╲         │
      │   ╱           ╲╱             ╲        │
      │  ╱    ∂/∂x     ░░░░    ∫ f(x)  ╲       │
      │ ╱    lim Σ      ░░░░     ∇·F     ╲      │
      │╲    n→∞         ░░░░              ╱     │
      │ ╲               ░░░░             ╱      │
      │  ╲    e^{iπ}+1  ░░░░  det(A)   ╱       │
      │   ╲             ░░░░           ╱        │
      │    ╲           ░░░░           ╱         │
      │     ╲   MATHEMATICAL MASTERY  ╱          │
      │      └─────────────────────────┘         │`

const PIPELINE_STATES = [
  { label: 'UNSEEN', status: 'unseen' as const },
  { label: 'LEARNING', status: 'learning' as const },
  { label: 'PRACTICING', status: 'practicing' as const },
  { label: 'MASTERED', status: 'mastered' as const },
  { label: 'DECAYING', status: 'decaying' as const },
]

export default function HomePage() {
  const [theme, setTheme] = useState<'dark' | 'light'>('dark')
  const [mounted, setMounted] = useState(false)

  useEffect(() => {
    const saved = localStorage.getItem('mathua-theme')
    const initialTheme = saved === 'light' || saved === 'dark' ? saved : 'dark'
    setTheme(initialTheme)
    setMounted(true)
  }, [])

  const conceptCount = conceptsData.length
  const connectionCount = conceptsData.reduce(
    (sum: number, c: any) => sum + ((c.prerequisites as any[])?.length || 0),
    0
  )
  const domainCount = new Set(conceptsData.map((c: any) => c.domain)).size

  const domainOrder = [
    'counting',
    'arithmetic',
    'fractions',
    'prealgebra',
    'algebra',
    'geometry',
    'trigonometry',
    'calculus',
    'statistics',
    'linear_algebra',
    'discrete_math',
    'complex_numbers',
    'number_theory',
    'differential_equations',
    'abstract_algebra',
    'topology',
  ]
  const domainLabels: Record<string, string> = {
    counting: 'Counting',
    arithmetic: 'Arithmetic',
    fractions: 'Fractions',
    prealgebra: 'Pre-Algebra',
    algebra: 'Algebra',
    geometry: 'Geometry',
    trigonometry: 'Trigonometry',
    calculus: 'Calculus',
    statistics: 'Statistics',
    linear_algebra: 'Linear Algebra',
    discrete_math: 'Discrete Math',
    complex_numbers: 'Complex Numbers',
    number_theory: 'Number Theory',
    differential_equations: 'Differential Equations',
    abstract_algebra: 'Abstract Algebra',
    topology: 'Topology',
  }
  const domainCounts = conceptsData.reduce(
    (acc: Record<string, number>, c: any) => {
      acc[c.domain] = (acc[c.domain] || 0) + 1
      return acc
    },
    {} as Record<string, number>
  )

  const levels = [
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

  useEffect(() => {
    if (!mounted) return
    if (theme === 'dark') {
      document.documentElement.classList.add('dark')
    } else {
      document.documentElement.classList.remove('dark')
    }
  }, [theme, mounted])

  const toggleTheme = () => {
    const next = theme === 'dark' ? 'light' : 'dark'
    setTheme(next)
    localStorage.setItem('mathua-theme', next)
  }

  if (!mounted) {
    return <div style={{ background: '#0a0f1a', minHeight: '100vh' }} />
  }

  return (
    <div className="max-w-container mx-auto px-6 max-sm:px-4">
      {/* Theme Toggle */}
      <button
        onClick={toggleTheme}
        className="fixed top-[50px] right-5 z-[1000] bg-mathua-surface border border-mathua-border-strong text-mathua-primary px-3.5 py-2 rounded-md font-mono text-xs cursor-pointer transition-all duration-200 hover:border-mathua-blue hover:text-mathua-blue"
        aria-label="Toggle theme"
      >
        {theme === 'dark' ? '\u2600' : '\u263E'}
      </button>

      {/* ── Hero ── */}
      <section className="py-20 max-sm:py-12 text-center bg-mathua-bg">
        <div className="flex justify-end mb-6">
          <a href="/how-it-works" className="border border-mathua-border-strong text-mathua-primary rounded-md h-10 px-6 inline-flex items-center text-[13px] font-medium transition-all duration-200 hover:border-mathua-blue hover:text-mathua-blue">
            How it works
          </a>
        </div>

        <h1 className="font-serif font-semibold text-[56px] max-sm:text-[32px] leading-tight tracking-[-0.03em] text-mathua-primary mb-4">
          Master the foundation.
          <br />
          Earn the abstraction.
        </h1>
        <p className="text-mathua-secondary text-base max-w-[600px] mx-auto mb-8 leading-relaxed">
          Mathua is an open-source adaptive math learning engine. It never lets you advance until
          you have truly mastered the prerequisite — both speed and accuracy must be proven.
        </p>

        <div className="flex gap-3 justify-center items-center mb-8 max-sm:flex-col max-sm:[&_a]:w-full max-sm:[&_a]:max-w-[280px]">
          <a
            href="#"
            className="bg-mathua-blue text-white hover:bg-mathua-blue-hover rounded-md h-10 px-6 inline-flex items-center justify-center text-[13px] font-medium transition-all duration-200"
          >
            Open the web app
          </a>
          <a
            href="https://github.com/chuma-beep/mathua"
            className="border border-mathua-border-strong text-mathua-primary hover:border-mathua-blue hover:text-mathua-blue rounded-md h-10 px-6 inline-flex items-center justify-center text-[13px] font-medium transition-all duration-200"
          >
            View on GitHub
          </a>
        </div>

        <div className="flex gap-2 justify-center flex-wrap mb-10">
          <span className="font-mono text-[11px] font-medium px-3 py-1.5 rounded bg-mathua-surface-elevated text-mathua-secondary border border-mathua-border">
            {conceptCount} topics
          </span>
          <span className="font-mono text-[11px] font-medium px-3 py-1.5 rounded bg-mathua-surface-elevated text-mathua-secondary border border-mathua-border">
            {connectionCount} connections
          </span>
          <span className="font-mono text-[11px] font-medium px-3 py-1.5 rounded bg-mathua-surface-elevated text-mathua-secondary border border-mathua-border">
            {domainCount} domains
          </span>
          <span className="font-mono text-[11px] font-medium px-3 py-1.5 rounded bg-mathua-surface-elevated text-mathua-secondary border border-mathua-border">
            web + desktop
          </span>
          <span className="font-mono text-[11px] font-medium px-3 py-1.5 rounded bg-mathua-surface-elevated text-mathua-secondary border border-mathua-border">
            open source
          </span>
        </div>

        {/* ASCII Art Banner */}
        <AsciiBanner text={ASCII_HERO} className="mb-10 text-mathua-gold text-center" />

        {/* 3D Concept Graph */}
        <div className="bg-mathua-surface border border-mathua-border rounded-lg p-6">
          <MathConceptGraph3D
            theme={theme}
            concepts={conceptsData.map((c: any) => ({
              id: c.id,
              label: c.label,
              domain: c.domain,
              prerequisites: c.prerequisites,
            }))}
          />
        </div>
      </section>

      <AsciiDivider pattern="wave" />

      {/* ── How It Works ── */}
      <section className="py-20 max-sm:py-12">
        <SectionHeader label="How it works" title="A different kind of math system" />
        <div className="flex flex-wrap gap-6 mt-7 max-sm:flex-col">
          <InfoCard title="Mastery gating">
            You cannot advance until your streak and response time both meet the threshold. Knowing
            the answer is not enough — you must know it fast.
          </InfoCard>
          <InfoCard title="Concept graph">
            Every concept is a node with explicit prerequisites. The scheduler reads the graph and
            your progress to decide what you see next.
          </InfoCard>
          <InfoCard title="Spaced repetition">
            Concepts you master resurface automatically before they decay. Reviews are woven into
            your session — there is no separate review mode.
          </InfoCard>
        </div>
      </section>

      {/* ── Mastery Pipeline ── */}
      <section className="py-20 max-sm:py-12">
        <SectionHeader label="Mastery pipeline" title="Five states. One direction." />
        <Pipeline states={PIPELINE_STATES} className="my-7" />
        <div className="flex justify-center">
          <FormulaBlock
            code={`priority = (0.7 × days_since_last_seen) + (0.3 × (1 − mastery))
        + 5.0 if DECAYING  +  2.0 if newly_unlocked`}
          />
        </div>
        <p className="text-center text-mathua-secondary text-sm mt-5">
          The scheduler enforces three hard rules: prerequisites must be mastered before a concept
          unlocks, the same concept never appears twice in a row, and roughly 70% of each session is
          new material.
        </p>
      </section>

      <AsciiDivider pattern="dash" />

      {/* ── Curriculum ── */}
      <section className="py-20 max-sm:py-12">
        <SectionHeader label="Curriculum" title="What Mathua covers" />
        <DomainTable
          className="mt-5"
          rows={domainOrder.map((d) => ({
            domain: domainLabels[d],
            count: domainCounts[d],
          }))}
        />
        <p className="text-mathua-muted text-[13px] italic mt-3">
          Problems are generated on demand — never stored. There is nothing to memorise.
        </p>
      </section>

      <AsciiDivider pattern="wave" />

      {/* ── Progression ── */}
      <section className="py-20 max-sm:py-12">
        <SectionHeader label="Progression" title="Rank by mastery. Level by depth." />
        <div className="flex flex-wrap gap-8 items-start mt-7 max-sm:flex-col">
          <div className="flex-1 min-w-[280px]">
            <p className="text-mathua-secondary text-sm leading-relaxed text-center">
              The leaderboard resets every Monday at 00:00 UTC. Your score is calculated from three
              components:
            </p>
            <div className="flex justify-center mt-4">
              <FormulaBlock
                code={`score = (mastered_count × 100)
    + speed_bonus
     + (current_streak × 10)`}
              />
            </div>
          </div>
          <div className="flex-1 min-w-[280px]">
            <ProgressionLevels levels={levels} />
          </div>
        </div>
      </section>

      <AsciiDivider pattern="dash" />

      {/* ── Platforms ── */}
      <section className="py-20 max-sm:py-12">
        <SectionHeader label="Platforms" title="One engine. Two ways to run it." />
        <div className="flex flex-wrap gap-8 items-start mt-7 max-sm:flex-col">
          <div className="flex-1 min-w-[280px] bg-mathua-surface border border-mathua-border rounded-lg p-6">
            <div className="font-mono text-[10px] uppercase tracking-[0.12em] text-mathua-muted mb-2">
              Web
            </div>
            <h3 className="font-sans text-base font-semibold text-mathua-primary mb-3.5">
              Browser
            </h3>
            <ul className="list-none text-[13px] text-mathua-secondary leading-relaxed">
              <li className="before:content-['—'] before:text-mathua-muted before:mr-1">
                Runs in any modern browser
              </li>
              <li className="before:content-['—'] before:text-mathua-muted before:mr-1">
                React frontend with KaTeX math rendering
              </li>
              <li className="before:content-['—'] before:text-mathua-muted before:mr-1">
                Account required — progress syncs across devices
              </li>
              <li className="before:content-['—'] before:text-mathua-muted before:mr-1">
                Global weekly leaderboard
              </li>
              <li className="before:content-['—'] before:text-mathua-muted before:mr-1">
                Interactive concept graph view
              </li>
            </ul>
          </div>
          <div className="flex-1 min-w-[280px] bg-mathua-surface border border-mathua-border rounded-lg p-6">
            <div className="font-mono text-[10px] uppercase tracking-[0.12em] text-mathua-muted mb-2">
              Desktop TUI
            </div>
            <h3 className="font-sans text-base font-semibold text-mathua-primary mb-3.5">
              Terminal
            </h3>
            <ul className="list-none text-[13px] text-mathua-secondary leading-relaxed">
              <li className="before:content-['—'] before:text-mathua-muted before:mr-1">
                Single binary download — no runtime dependencies
              </li>
              <li className="before:content-['—'] before:text-mathua-muted before:mr-1">
                Bubble Tea terminal interface
              </li>
              <li className="before:content-['—'] before:text-mathua-muted before:mr-1">
                SQLite storage — all data stays on your machine
              </li>
              <li className="before:content-['—'] before:text-mathua-muted before:mr-1">
                No account needed
              </li>
              <li className="before:content-['—'] before:text-mathua-muted before:mr-1">
                Fully offline after download
              </li>
            </ul>
          </div>
        </div>
        <div className="flex justify-center mt-8">
          <pre className="bg-mathua-code border border-mathua-border rounded-md p-5 font-mono text-xs text-mathua-blue whitespace-pre overflow-x-auto text-center leading-relaxed">
            {`┌──────────────────┐    ┌──────────────────┐
│  Web (React)    │    │  TUI (Bubble Tea)│
└────────┬─────────┘    └────────┬─────────┘
         ▼                       ▼
┌──────────────────────────────────────┐
│            API Layer               │
└──────────────────┬───────────────────┘
                   ▼
┌──────────────────────────────────────┐
│        Scheduling Engine           │
└──────────────────┬────────────────┘
                   ▼
┌──────────────────────────────────────┐
│   Grading & Problem Generation    │
└──────────────────┬───────────────────┘
                   ▼
┌──────────────────────────────────────┐
│  Data Layer (PostgreSQL / SQLite) │
└──────────────────────────────────────┘`}
          </pre>
        </div>
      </section>

      <AsciiDivider pattern="wave" />

      {/* ── Contributing ── */}
      <section className="py-20 max-sm:py-12">
        <SectionHeader label="Contributing" title="Built to be extended." />
        <p className="text-mathua-secondary text-sm leading-relaxed text-center max-w-[600px] mx-auto">
          Every concept is a JSON node. Every problem is a Go generator function. Every contribution
          goes through a graph validator that rejects cycles and orphaned nodes automatically.
        </p>
        <div className="flex justify-center mt-4">
          <pre className="bg-mathua-code border border-mathua-border rounded-md p-5 font-mono text-xs text-mathua-blue whitespace-pre overflow-x-auto text-left leading-relaxed">
            {`{
  "id":                "arith.add.multi",
  "label":             "Multi-digit addition",
  "domain":            "arithmetic",
  "prerequisites":     ["arith.add.single", "arith.add.carry"],
  "mastery_threshold": {
    "streak":          5,
    "avg_time_seconds": 8
  }
}`}
          </pre>
        </div>
        <div className="flex justify-center mt-6">
          <a
            href="https://github.com/chuma-beep/mathua/blob/main/CONTRIBUTING.md"
            className="bg-mathua-blue text-white hover:bg-mathua-blue-hover rounded-md h-10 px-6 inline-flex items-center justify-center text-[13px] font-medium transition-all duration-200"
          >
            Read CONTRIBUTING.md →
          </a>
        </div>
      </section>

      <Footer />
    </div>
  )
}
