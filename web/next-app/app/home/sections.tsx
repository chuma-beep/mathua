'use client'

import dynamic from 'next/dynamic'
import { useEffect, useRef, useState } from 'react'
import SectionHeader from '../../components/SectionHeader'
import Pipeline from '../../components/Pipeline'
import FormulaBlock from '../../components/FormulaBlock'
import DomainTable from '../../components/DomainTable'
import ProgressionLevels from '../../components/ProgressionLevels'
import Loading from '../../components/Loading'
import { loadPositionEntries } from '../../lib/graphPositions'
import { domainColor } from '../../lib/graphDomains'
import {
  PIPELINE_STATES,
  conceptCount,
  connectionCount,
  domainCount,
  domainCounts,
  domainLabels,
  domainOrder,
  levels,
} from './data'
import {
  bodyFont,
  codeQuoteStyle,
  ctaPrimaryStyle,
  ctaSecondaryStyle,
  headingFont,
  loadingGraphStyle,
  monoFont,
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

export function LazyGraphMount({ children, fallback }: { children: React.ReactNode; fallback?: React.ReactNode }) {
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

  return <div ref={ref}>{visible ? children : fallback ?? <div style={loadingGraphStyle}><Loading label="PREPARING GRAPH" /></div>}</div>
}

// Static SVG snapshot of the DAG (orthographic x/y scatter from the same
// precomputed positions the WebGL scene uses). Positions arrive via fetch
// (public/ copy) so the 26KB payload stays out of the page bundle; the box
// paints with first paint and dots pop in when the fetch resolves.
export function GraphPoster() {
  const [entries, setEntries] = useState<Array<[string, [number, number, number]]>>([])
  useEffect(() => {
    let live = true
    loadPositionEntries().then(list => {
      if (live) setEntries(list)
    })
    return () => {
      live = false
    }
  }, [])
  const W = 400
  const H = 400
  const PAD = 16
  let minX = Infinity
  let maxX = -Infinity
  let minY = Infinity
  let maxY = -Infinity
  for (const [, p] of entries) {
    if (p[0] < minX) minX = p[0]
    if (p[0] > maxX) maxX = p[0]
    if (p[1] < minY) minY = p[1]
    if (p[1] > maxY) maxY = p[1]
  }
  const spanX = Math.max(maxX - minX, 1e-6)
  const spanY = Math.max(maxY - minY, 1e-6)
  // Fit data bounds exactly (no letterbox); flip y so root layers read top-down.
  const dot = (p: [number, number, number]): [number, number] => [
    PAD + ((p[0] - minX) / spanX) * (W - 2 * PAD),
    PAD + (1 - (p[1] - minY) / spanY) * (H - 2 * PAD),
  ]
  return (
    <div style={loadingGraphStyle} aria-hidden="true">
      <svg viewBox={`0 0 ${W} ${H}`} width="100%" height="100%" preserveAspectRatio="xMidYMid meet" role="presentation">
        {entries.map(([id, p], i) => {
          const [cx, cy] = dot(p)
          return (
            <circle
              key={id}
              cx={cx.toFixed(1)}
              cy={cy.toFixed(1)}
              r={i % 12 === 0 ? 2.6 : 1.6}
              fill={i % 12 === 0 ? 'var(--accent-blue)' : 'var(--text-muted)'}
              opacity={i % 12 === 0 ? 0.9 : 0.7}
            />
          )
        })}
      </svg>
    </div>
  )
}

// Static cluster key (no canvas cost): one dot per domain so the node
// cloud reads as grouped clusters. Counts come from the build-generated
// meta file, so the concept corpus stays out of the page bundle.
function DomainLegend({ theme }: { theme: 'dark' | 'light' }) {
  return (
    <div
      aria-label="Domains in the graph"
      style={{
        display: 'flex',
        flexWrap: 'wrap',
        justifyContent: 'center',
        gap: '4px 12px',
        marginTop: '0.75rem',
        fontFamily: monoFont,
        fontSize: '11px',
        color: 'var(--text-muted)',
      }}
    >
      {domainOrder
        .filter(domain => domainCounts[domain] > 0)
        .map(domain => (
          <span key={domain} title={`${domainCounts[domain]} topics`} style={{ whiteSpace: 'nowrap' }}>
            <span style={{ color: domainColor(domain, theme) }}>●</span>{' '}
            {domain.replace(/_/g, ' ')}
          </span>
        ))}
    </div>
  )
}

export function HeroSection({ theme, onGetStarted }: { theme: 'dark' | 'light'; onGetStarted: () => void }) {
  return (
    <section className="py-20 max-sm:py-8 text-center" style={{ background: 'transparent' }}>
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

      <div className="flex gap-3 justify-center items-center mb-10 max-sm:mb-6 max-sm:flex-col max-sm:[&_a]:w-full max-sm:[&_a]:max-w-[280px] max-sm:px-2 min-w-0">
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

      <div style={{ ...statsRowStyle, marginBottom: '1.25rem' }} className="flex flex-wrap justify-center gap-x-2 gap-y-1 px-2 text-center">
        <span>{conceptCount} topics</span>
        <span style={{ color: 'var(--border-strong)' }}>·</span>
        <span>{connectionCount} connections</span>
        <span style={{ color: 'var(--border-strong)' }}>·</span>
        <span>{domainCount} domains</span>
      </div>
      <DomainLegend theme={theme} />
      <div
        className="sm:hidden"
        style={{ fontFamily: monoFont, fontSize: '11px', color: 'var(--text-muted)', marginTop: '0.5rem' }}
      >
        Pinch to zoom · tap a node to explore
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
          <LazyGraphMount fallback={<GraphPoster />}>
            <MathConceptGraph3D theme={theme} />
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

export function SocialProofSection() {
  return (
    <section className="py-20 max-sm:py-12">
      <SectionHeader title="Why it is trusted" />
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
        Open-source under MIT. Sequencing inspired by Math Academy, content shaped by
        OpenStax, MIT OpenCourseWare, and Art of Problem Solving. Lessons build on
        Algebrica under CC BY-NC 4.0. No accounts to start, no paywall on the graph.
      </p>
      <div style={statsRowStyle} className="flex flex-wrap justify-center gap-x-2 gap-y-1 px-2 text-center mt-6">
        <span>{conceptCount} worked concepts</span>
        <span style={{ color: 'var(--border-strong)' }}>·</span>
        <span>generated problems, never stored</span>
        <span style={{ color: 'var(--border-strong)' }}>·</span>
        <span>150 XP mastery-check quiz</span>
      </div>
    </section>
  )
}

const FAQS = [
  {
    q: 'Do I need an account to start?',
    a: 'No. Get started creates a guest profile instantly. Sign up later only if you want your progress on another device.',
  },
  {
    q: 'Why does speed matter, not just accuracy?',
    a: 'Each concept has a streak and a time threshold. A slow correct answer counts as weak mastery, so fluency is proven before you advance.',
  },
  {
    q: 'What happens every 150 XP?',
    a: 'A timed, closed-book mastery-check quiz at 80% difficulty over recent material, with immediate remedial work and a retake path.',
  },
  {
    q: 'Is Mathua free and open-source?',
    a: 'Yes, MIT licensed. The concept graph, generators, and lessons are in the repo and validated by automated checks on every change.',
  },
]

export function FaqSection() {
  return (
    <section className="py-20 max-sm:py-12">
      <SectionHeader title="Questions, answered." />
      <div style={{ marginTop: '2rem' }}>
        {FAQS.map(f => (
          <details key={f.q} style={{ borderTop: '0.5px solid var(--border)', padding: '1rem 0' }}>
            <summary
              style={{
                fontFamily: headingFont,
                fontSize: '1rem',
                color: 'var(--text-primary)',
                cursor: 'pointer',
              }}
            >
              {f.q}
            </summary>
            <p
              style={{
                fontFamily: bodyFont,
                fontSize: '0.95rem',
                color: 'var(--text-secondary)',
                lineHeight: 1.85,
                marginTop: '0.5rem',
                marginLeft: 'clamp(0.5rem, 3vw, 1.5rem)',
              }}
            >
              {f.a}
            </p>
          </details>
        ))}
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
