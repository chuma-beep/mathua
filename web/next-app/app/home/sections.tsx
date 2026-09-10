'use client'

import dynamic from 'next/dynamic'
import { useEffect, useRef, useState } from 'react'
import SectionHeader from '../../components/SectionHeader'
import Pipeline from '../../components/Pipeline'
import FormulaBlock from '../../components/FormulaBlock'
import DomainTable from '../../components/DomainTable'
import ProgressionLevels from '../../components/ProgressionLevels'
import Loading from '../../components/Loading'
import {
  PIPELINE_STATES,
  conceptCount,
  connectionCount,
  domainCount,
  domainCounts,
  domainLabels,
  domainOrder,
  heroConcepts,
  levels,
} from './data'
import {
  bodyFont,
  codeQuoteStyle,
  ctaPrimaryStyle,
  ctaSecondaryStyle,
  headingFont,
  loadingGraphStyle,
  statsRowStyle,
} from './styles'

const MathConceptGraph3D = dynamic(() => import('../../components/MathConceptGraph3D'), {
  ssr: false,
  loading: () => (
    <div style={loadingGraphStyle}>
      <Loading label="LOADING GRAPH" />
    </div>
  ),
})

export function LazyGraphMount({ children }: { children: React.ReactNode }) {
  const ref = useRef<HTMLDivElement>(null)
  const [visible, setVisible] = useState(false)

  useEffect(() => {
    const el = ref.current
    if (!el || typeof IntersectionObserver === 'undefined') {
      setVisible(true)
      return
    }
    let idleId: number | undefined
    let timeoutId: ReturnType<typeof setTimeout> | undefined
    const obs = new IntersectionObserver(
      entries => {
        if (!entries[0]?.isIntersecting) return
        obs.disconnect()
        // Mount after the page is interactive: the WebGL context + 630-node
        // scene must not compete with hydration/first paint.
        if (typeof window.requestIdleCallback === 'function') {
          idleId = window.requestIdleCallback(() => setVisible(true), { timeout: 1500 })
        } else {
          timeoutId = setTimeout(() => setVisible(true), 200)
        }
      },
      { rootMargin: '200px' }
    )
    obs.observe(el)
    return () => {
      obs.disconnect()
      if (idleId !== undefined) window.cancelIdleCallback?.(idleId)
      if (timeoutId !== undefined) clearTimeout(timeoutId)
    }
  }, [])

  return <div ref={ref}>{visible ? children : <div style={loadingGraphStyle}><Loading label="PREPARING GRAPH" /></div>}</div>
}

export function HeroSection({ theme, onGetStarted }: { theme: 'dark' | 'light'; onGetStarted: () => void }) {
  return (
    <section className="py-20 max-sm:py-12 text-center" style={{ background: 'transparent' }}>
      <h1
        style={{
          fontFamily: headingFont,
          fontWeight: 400,
          fontSize: 'clamp(2.2rem, 5vw, 3.8rem)',
          lineHeight: 1.2,
          color: 'var(--text-primary)',
          marginBottom: '1rem',
          letterSpacing: '-0.01em',
        }}
      >
        Adaptive math learning platform
      </h1>
      <p
        style={{
          fontFamily: bodyFont,
          fontSize: 'clamp(0.95rem, 2.5vw, 1.1rem)',
          color: 'var(--text-secondary)',
          maxWidth: '600px',
          margin: '0 auto 2rem',
          lineHeight: 1.85,
        }}
      >
        Mathua is an open-source adaptive math learning engine. It never lets you advance until
        you have truly mastered the prerequisite, both speed and accuracy must be proven.
      </p>

      <div className="flex gap-3 justify-center items-center mb-10 max-sm:flex-col max-sm:[&_a]:w-full max-sm:[&_a]:max-w-[280px] max-sm:px-2 min-w-0">
        <button
          type="button"
          onClick={onGetStarted}
          style={ctaPrimaryStyle}
          className="inline-flex items-center justify-center min-h-[44px] max-sm:w-full max-sm:max-w-[280px]"
        >
          Get started
        </button>
        <a href="https://github.com/chuma-beep/mathua" style={ctaSecondaryStyle} className="inline-flex items-center justify-center min-h-[44px] max-sm:w-full max-sm:max-w-[280px]">
          View on GitHub
        </a>
      </div>

      <div style={statsRowStyle} className="flex flex-wrap justify-center gap-x-2 gap-y-1 px-2 text-center">
        <span>{conceptCount} topics</span>
        <span style={{ color: 'var(--border-strong)' }}>·</span>
        <span>{connectionCount} connections</span>
        <span style={{ color: 'var(--border-strong)' }}>·</span>
        <span>{domainCount} domains</span>
      </div>

      {/* 3D Concept Graph */}
      <div
        className="p-3 sm:p-6 w-full max-w-full min-w-0 overflow-hidden"
        style={{
          border: '0.5px solid var(--border)',
          borderRadius: 0,
          background: 'var(--graph-surface)',
        }}
      >
        <div className="w-full max-w-full min-w-0 overflow-hidden">
          <LazyGraphMount>
            <MathConceptGraph3D theme={theme} concepts={heroConcepts} />
          </LazyGraphMount>
        </div>
      </div>
    </section>
  )
}

const FEATURES = [
  {
    numeral: 'I.',
    title: 'Mastery gating',
    body: 'You cannot advance until your streak and response time both meet the threshold. Knowing the answer is not enough; you must know it fast.',
  },
  {
    numeral: 'II.',
    title: 'Concept graph',
    body: 'Every concept is a node with explicit prerequisites. The scheduler reads the graph and your progress to decide what you see next.',
  },
  {
    numeral: 'III.',
    title: 'Spaced repetition',
    body: 'Concepts you master resurface automatically before they decay. Reviews are woven into your session; there is no separate review mode.',
  },
]

export function FeaturesSection() {
  return (
    <section className="py-20 max-sm:py-12">
      <SectionHeader title="3 main features" />
      <div style={{ marginTop: '2rem' }}>
        {FEATURES.map(f => (
          <div key={f.title} style={{ borderTop: '0.5px solid var(--border)', padding: '1.2rem 0' }}>
            <div style={{ display: 'flex', gap: '0.5rem', marginBottom: '0.3rem' }}>
              <span style={{ fontFamily: headingFont, fontSize: '1.1rem', color: 'var(--accent-blue)' }}>{f.numeral}</span>
              <span style={{ fontFamily: headingFont, fontSize: '1.1rem', color: 'var(--text-primary)' }}>{f.title}</span>
            </div>
            <p style={{ fontFamily: bodyFont, fontSize: '1rem', color: 'var(--text-secondary)', lineHeight: 1.85, marginLeft: 'clamp(0.5rem, 3vw, 1.5rem)' }}>
              {f.body}
            </p>
          </div>
        ))}
      </div>
    </section>
  )
}

export function PipelineSection() {
  return (
    <section className="py-20 max-sm:py-12">
      <SectionHeader title="The states" />
      <Pipeline states={PIPELINE_STATES} className="my-7" />
      <div className="flex justify-center">
        <FormulaBlock
          code={`priority = (0.7 × days_since_last_seen)
         + (0.3 × (1 − mastery))
         + 5.0  if DECAYING
         + 2.0  if newly_unlocked`}
        />
      </div>
      <p
        style={{
          fontFamily: bodyFont,
          fontSize: '0.95rem',
          color: 'var(--text-secondary)',
          textAlign: 'center',
          marginTop: '1.25rem',
          lineHeight: 1.85,
        }}
      >
        The scheduler enforces three hard rules: prerequisites must be mastered before a concept
        unlocks, the same concept never appears twice in a row, and roughly 70% of each session is
        new material.
      </p>
    </section>
  )
}

export function CurriculumSection() {
  return (
    <section className="py-20 max-sm:py-12">
      <SectionHeader title="What Mathua covers" />
      <DomainTable
        className="mt-5"
        rows={domainOrder.map((d) => ({
          domain: domainLabels[d],
          count: domainCounts[d],
        }))}
      />
      <p
        style={{
          fontFamily: bodyFont,
          fontStyle: 'italic',
          fontSize: '13px',
          color: 'var(--text-muted)',
          marginTop: '0.75rem',
        }}
      >
        Problems are generated on demand, never stored. There is nothing to memorise.
      </p>
    </section>
  )
}

export function ProgressionSection() {
  return (
    <section className="py-20 max-sm:py-12">
      <SectionHeader title="Ranking system" />
      <div className="flex flex-wrap gap-8 items-start mt-7 max-sm:flex-col max-sm:gap-4">
        <div className="flex-1 min-w-[280px] max-sm:min-w-0">
          <p
            style={{
              fontFamily: bodyFont,
              fontSize: '0.95rem',
              color: 'var(--text-secondary)',
              lineHeight: 1.85,
              textAlign: 'center',
            }}
          >
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
        <div className="flex-1 min-w-[280px] max-sm:min-w-0">
          <ProgressionLevels levels={levels} />
        </div>
      </div>
    </section>
  )
}

export function ContributingSection() {
  return (
    <section className="py-20 max-sm:py-12">
      <SectionHeader title="How it is extended." />
      <p
        style={{
          fontFamily: bodyFont,
          fontSize: '0.95rem',
          color: 'var(--text-secondary)',
          lineHeight: 1.85,
          textAlign: 'center',
          maxWidth: '600px',
          margin: '0 auto',
        }}
      >
        Every concept is a JSON node. Every problem is a Go generator function. Every contribution
        goes through a graph validator that rejects cycles and orphaned nodes automatically.
      </p>
      <div className="flex justify-center mt-4">
        <pre style={codeQuoteStyle}>
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
          href="/docs/contributing"
          className="link-underline"
          style={{
            fontFamily: headingFont,
            fontStyle: 'italic',
            fontSize: '1rem',
            color: 'var(--accent-blue)',
          }}
        >
          Read the contributing guide →
        </a>
      </div>
    </section>
  )
}
