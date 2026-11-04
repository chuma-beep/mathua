'use client'

import Link from 'next/link'
import { useEffect, useRef, useState } from 'react'
import { LOGIN_HERO } from '../../lib/loginHero'
import {
  PIPELINE_STATES,
  conceptCount,
  domainCounts,
  domainLabels,
  domainOrder,
  levels,
} from './data'

// The two hero actions. They are written out in full rather than composed, because
// the ghost button used to be the primary button plus overrides and the two
// background utilities (`bg-mathua-blue` and `bg-transparent`, `hover:bg-mathua-blue-hover`
// and `hover:bg-white/15`) collided: same property, same specificity, so the winner was
// whichever Tailwind emitted last and the GitHub button hovered to blue on the image.
// Separate constants with no overlapping property cannot collide. Both actions are
// content-width and left-aligned on phones: as full-width bars with only the page's
// 16px gutters they read as touching the screen edges.
const CTA_PRIMARY =
  'inline-flex min-h-[44px] items-center justify-center bg-mathua-blue px-6 py-3 text-sm font-medium text-white transition-colors hover:bg-mathua-blue-hover focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-mathua-blue'
// On the artwork the second action is a white hairline over a dark scrim, not the
// themed surface border: a `mathua-border` line is tuned for the page background,
// not for a photograph, and would disappear against the drawing.
// Wipe-fill ghost: on hover a blue layer sweeps left to right behind the label
// (origin-left scale-x). The fill lives on a child span, not the anchor, so the
// anchor never carries two background utilities at once, the collision that once
// turned this button blue. `overflow-hidden` clips the sweep; the label sits
// `relative` above it.
const CTA_GHOST_ON_ART =
  'group relative inline-flex min-h-[44px] items-center justify-center overflow-hidden border border-white/45 bg-transparent px-6 py-3 text-sm font-medium text-white transition-colors hover:border-white focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-white'

export function HeroSection({ onGetStarted }: { onGetStarted: () => void }) {
  return (
    // One responsive hero, not two stacked variants. Below lg it is the full-bleed
    // treatment (B): art behind the copy under a scrim. At lg and up it is the
    // side-by-side panel (A): copy left, art right, scrim gone and copy back on
    // theme tokens. A single h1 and a single img serve both, no duplicated
    // heading, no second download.
    //
    // The section sits in the page's capped, padded column, so the picture is a
    // banner aligned to the content rather than a window-wide bleed, every other
    // section shares that column, and a hero that ignored it would be the odd one
    // out. `isolate` keeps the background's stacking context local so it cannot
    // paint over the sticky header.
    <section className="group/hero relative isolate overflow-hidden border-b border-mathua-division lg:grid lg:grid-cols-2 lg:items-stretch lg:gap-10">
      {/* Art. Full-bleed backdrop below lg; the grid's right column at lg and up.
          `lg:relative` re-seats it as a grid item while the img keeps filling it
          with object-cover, and `lg:order-2` puts it right of the copy. DOM order
          stays art-first so that on phones, where both layers are positioned, 
          the copy paints above the picture.
          The picture never touches the header: `top-4` holds it off the section
          top on phones, `lg:mt-10` on desktop (`lg:top-auto lg:bottom-auto` cancel
          the phone offsets once it is a grid item). */}
      <div aria-hidden="true" className="hero-art-in absolute inset-x-0 bottom-0 top-4 overflow-hidden lg:relative lg:bottom-auto lg:top-auto lg:order-2 lg:mt-10 lg:min-h-[480px]">
        <img
          src={LOGIN_HERO.file}
          alt=""
          width={LOGIN_HERO.width}
          height={LOGIN_HERO.height}
          loading="eager"
          decoding="async"
          className="absolute inset-0 size-full object-cover object-center transition-transform duration-700 ease-out group-hover/hero:scale-[1.03] motion-reduce:transition-none motion-reduce:group-hover/hero:scale-100"
        />
        {/* Phone scrim only: it darkens the text side and fades toward the art, so
            the copy has contrast without flattening the drawing. Stronger than the
            earlier wash, because phone copy spans the full frame. */}
        <div className="absolute inset-0 bg-gradient-to-r from-black/70 via-black/50 to-black/35 lg:hidden" />
      </div>
      <div className="relative max-w-[72ch] px-4 py-16 [text-shadow:0_1px_14px_rgba(0,0,0,0.45)] md:py-24 lg:order-1 lg:px-0 lg:[text-shadow:none]">
        <div className="section-label hero-label reveal">( 01 · Title Block )</div>
        <h1
          className="reveal mt-5 max-w-[18ch] text-balance font-serif text-[clamp(2.5rem,6vw,4.5rem)] font-normal leading-[1.05] sm:leading-[0.95] tracking-tight text-white lg:text-mathua-primary"
          style={{ animationDelay: '80ms' }}
        >
          Adaptive math learning platform
        </h1>
        <p
          className="reveal mt-6 max-w-[46ch] text-pretty text-base sm:text-[15px] leading-relaxed text-white/85 lg:text-mathua-secondary"
          style={{ animationDelay: '160ms' }}
        >
          Mathua is an open-source adaptive math learning engine. You cannot advance until you
          have truly mastered the prerequisite, both speed and accuracy must be proven.
        </p>
        <div className="reveal mt-8 flex flex-col items-start gap-3 sm:flex-row sm:items-center" style={{ animationDelay: '240ms' }}>
          <button type="button" onClick={onGetStarted} className={CTA_PRIMARY}>
            Get started
          </button>
          <a
            href="https://github.com/chuma-beep/mathua"
            className={`${CTA_GHOST_ON_ART} lg:border-mathua-border lg:text-mathua-primary lg:hover:border-mathua-blue lg:focus-visible:outline-mathua-blue`}
          >
            <span
              aria-hidden="true"
              className="absolute inset-0 origin-left scale-x-0 bg-mathua-blue transition-transform duration-300 ease-out group-hover:scale-x-100 motion-reduce:transition-none"
            />
            <span className="relative group-hover:text-white">View on GitHub</span>
          </a>
        </div>
        <p className="reveal mt-4 font-mono text-[11px] text-mathua-blue" style={{ animationDelay: '320ms' }}>
          No account needed to start.
        </p>
      </div>
    </section>
  )
}

// Counts up when the stat scrolls into view. The target is rendered first so
// no-JS, reduced-motion, and e2e reads always see the real number; the
// observer restarts it from 0 only once the element is (barely) on screen.
function CountUp({ target, duration = 1200 }: { target: number; duration?: number }) {
  const ref = useRef<HTMLSpanElement>(null)
  useEffect(() => {
    const el = ref.current
    if (!el) return
    if (window.matchMedia('(prefers-reduced-motion: reduce)').matches) return
    let raf = 0
    const observer = new IntersectionObserver(
      entries => {
        if (!entries[0]?.isIntersecting) return
        observer.disconnect()
        const start = performance.now()
        const tick = (now: number) => {
          const t = Math.min(1, (now - start) / duration)
          const eased = 1 - Math.pow(1 - t, 3)
          el.textContent = String(Math.round(eased * target))
          if (t < 1) raf = requestAnimationFrame(tick)
        }
        raf = requestAnimationFrame(tick)
      },
      { threshold: 0 }
    )
    observer.observe(el)
    return () => {
      observer.disconnect()
      cancelAnimationFrame(raf)
    }
  }, [target, duration])
  return <span ref={ref}>{target}</span>
}

export function TrustSection() {
  return (
    <section data-reveal className="border-b border-mathua-division py-16">
      <div className="section-label">( 02 · Why it is trusted )</div>
      <div className="mt-5 grid gap-10 md:grid-cols-12 md:items-start">
        <p className="max-w-[52ch] text-sm leading-relaxed text-mathua-secondary md:col-span-7">
          Open-source under MIT. Sequencing inspired by Math Academy, with MIT OpenCourseWare and
          Art of Problem Solving consulted for curriculum structure and problem quality. Lessons
          build on OpenStax (CC BY-NC-SA 4.0) and Algebrica (CC BY-NC 4.0). No accounts to start,
          no paywall on the graph.
        </p>
        <div className="grid grid-cols-3 divide-x divide-mathua-division border border-mathua-border md:col-span-5">
          <div className="p-5">
            <div className="font-serif text-3xl tracking-tight text-mathua-primary"><CountUp target={conceptCount} /></div>
            <div className="mt-1 font-mono text-[10px] uppercase tracking-wider text-mathua-muted">worked concepts</div>
          </div>
          <div className="p-5">
            <div className="font-serif text-3xl tracking-tight text-mathua-primary">∞</div>
            <div className="mt-1 font-mono text-[10px] uppercase tracking-wider text-mathua-muted">generated problems</div>
          </div>
          <div className="p-5">
            <div className="font-serif text-3xl tracking-tight text-mathua-primary">5</div>
            <div className="mt-1 font-mono text-[10px] uppercase tracking-wider text-mathua-muted">questions per check</div>
          </div>
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
    body: 'Every concept is a node with explicit prerequisites. Mathua reads the graph and your progress together, and picks what you see next from the two.',
  },
  {
    numeral: 'III.',
    title: 'Spaced repetition',
    body: 'Concepts you master resurface automatically before they decay. Reviews are woven into your session; there is no separate review mode.',
  },
]

export function FeaturesSection() {
  return (
    <section data-reveal className="border-b border-mathua-division py-16">
      <div className="section-label">( 03 · 3 Main Features )</div>
      <div className="mt-6 grid gap-px border border-mathua-division bg-mathua-division md:grid-cols-3">
        {FEATURES.map(f => (
          <div key={f.title} className="bg-mathua-bg p-6">
            <div className="font-mono text-xs text-mathua-blue">{f.numeral}</div>
            <h4 className="mt-3 font-serif text-lg text-mathua-primary">{f.title}</h4>
            <p className="mt-2 text-[13px] leading-relaxed text-mathua-secondary">{f.body}</p>
          </div>
        ))}
      </div>
    </section>
  )
}

const chipClass = (status: string) => {
  if (status === 'mastered') return 'border border-mathua-blue bg-mathua-blue text-mathua-bg'
  if (status === 'decaying') return 'border border-dashed border-mathua-border text-mathua-primary'
  return 'border border-mathua-border text-mathua-primary'
}

export function StatesSection() {
  return (
    <section data-reveal className="border-b border-mathua-division py-16">
      <div className="section-label">( 04 · The States )</div>
      <div className="mt-6 flex flex-wrap items-center gap-2 font-mono text-xs">
        {PIPELINE_STATES.map((s, i) => (
          <span key={s.label} className="flex items-center gap-2">
            {i > 0 && <span className="text-mathua-muted">→</span>}
            <span className={`px-3 py-1.5 ${chipClass(s.status)}`}>{s.label}</span>
          </span>
        ))}
      </div>
      <p className="mt-4 font-mono text-[11px] text-mathua-muted">
        example · <span className="text-mathua-blue">arith.add.single</span> · 10 correct in a row at ≤8s → MASTERED
      </p>
      <div className="mt-6 overflow-x-auto border border-mathua-border bg-mathua-primary p-8 font-mono text-[13px] leading-loose text-mathua-bg dark:bg-mathua-code dark:text-mathua-primary">
        <span className="opacity-40">priority</span> = (0.7 × days_since_last_seen)
        <br />
        <span className="pl-20">+ (0.3 × (1 − mastery))</span>
        <br />
        <span className="pl-20">
          + 5.0{' '}
          <span className="underline decoration-mathua-blue decoration-2 underline-offset-2">if DECAYING</span>
        </span>
        <br />
        <span className="pl-20">
          + 2.0{' '}
          <span className="underline decoration-mathua-blue decoration-2 underline-offset-2">if newly_unlocked</span>
        </span>
      </div>
      <p className="mt-4 max-w-[60ch] text-[13px] leading-relaxed text-mathua-secondary">
        The scheduler enforces three hard rules: prerequisites must be mastered before a concept
        unlocks, the same concept never appears twice in a row and roughly 70% of each session
        is new material.
      </p>
    </section>
  )
}

function DomainTable({ rows }: { rows: [string, number][] }) {
  return (
    <table className="w-full font-mono text-xs">
      <tbody className="[&_td]:border-b [&_td]:border-mathua-division [&_td]:py-1.5">
        {rows.map(([domain, count]) => (
          <tr key={domain}>
            <td className="text-mathua-secondary">{domain}</td>
            <td className="text-right text-mathua-primary">{count}</td>
          </tr>
        ))}
      </tbody>
    </table>
  )
}

export function CoverageSection() {
  const rows: [string, number][] = domainOrder.map(d => [domainLabels[d], domainCounts[d]])
  const half = Math.ceil(rows.length / 2)
  return (
    <section data-reveal className="border-b border-mathua-division py-16">
      <div className="flex flex-wrap items-end justify-between gap-2">
        <div className="section-label">( 05 · What Mathua covers )</div>
        <Link
          href="/domains"
          className="font-mono text-[11px] text-mathua-muted transition-colors hover:text-mathua-blue"
        >
          All {conceptCount} concepts by subject →
        </Link>
      </div>
      <div className="mt-6 grid gap-x-12 gap-y-6 md:grid-cols-2">
        <DomainTable rows={rows.slice(0, half)} />
        <DomainTable rows={rows.slice(half)} />
      </div>
      <p className="mt-4 font-mono text-[11px] text-mathua-muted">
        Problems are generated on demand, never stored. There is nothing to memorise.
      </p>
    </section>
  )
}

export function RankingSection() {
  return (
    <section data-reveal className="border-b border-mathua-division py-16">
      <div className="section-label">( 06 · Ranking System )</div>
      <p className="mt-4 max-w-[52ch] text-sm text-mathua-secondary">
        The leaderboard resets every Monday at 00:00 UTC. Your score is calculated from three
        components:
      </p>
      <div className="mt-4 overflow-x-auto border border-mathua-border bg-mathua-primary p-6 font-mono text-[13px] leading-loose text-mathua-bg dark:bg-mathua-code dark:text-mathua-primary">
        <span className="opacity-40">score</span> = (mastered_count × 100)
        <br />
        <span className="pl-6">+ speed_bonus</span>
        <br />
        <span className="pl-6">+ (current_streak × 10)</span>
      </div>
      <div className="mt-8 font-mono text-[10px] uppercase tracking-[0.2em] text-mathua-muted">
        Levels of mastery
      </div>
      <div className="mt-3 grid grid-cols-3 gap-px border border-mathua-division bg-mathua-division md:grid-cols-9">
        {levels.map(l => {
          const elite = 'elite' in l && l.elite
          return (
            <div key={l.name} className={elite ? 'bg-mathua-primary p-3 text-mathua-bg dark:bg-mathua-code dark:text-mathua-primary' : 'bg-mathua-bg p-3'}>
              <div className={elite ? 'font-mono text-[10px] opacity-50' : 'font-mono text-[10px] text-mathua-muted'}>
                {l.num}
              </div>
              <div className={`mt-1 font-serif text-[13px] ${elite ? '' : 'text-mathua-primary'}`}>{l.name}</div>
              <div className={elite ? 'mt-1 font-mono text-[10px] opacity-50' : 'mt-1 font-mono text-[10px] text-mathua-muted'}>
                {l.range}
              </div>
            </div>
          )
        })}
      </div>
    </section>
  )
}

const FAQS = [
  {
    q: 'Do I need an account to start?',
    a: 'No. Begin studying immediately, accounts are optional and only needed to sync progress across devices.',
  },
  {
    q: 'Why does speed matter, not just accuracy?',
    a: 'Fluency requires retrieval under time pressure. Knowing the answer is not enough; you must know it fast, or the foundation will not hold the next concept.',
  },
  {
    q: 'When does a mastery check appear?',
    a: 'After enough learning XP. The short check verifies what you have learned; concepts you have already mastered return through spaced review instead.',
  },
  {
    q: 'Is Mathua free and open-source?',
    a: 'Yes, MIT licensed, with no paywall on the graph. Every concept and generator is in the repository.',
  },
]

export function FaqSection() {
  return (
    <section data-reveal className="grid gap-10 py-16 md:grid-cols-2">
      <div className="min-w-0">
        <div className="section-label">( 07 · Questions, answered. )</div>
        <div className="mt-5 divide-y divide-mathua-division border-y border-mathua-division">
          {FAQS.map(f => (
            <details key={f.q} className="group py-3">
              <summary className="flex cursor-pointer list-none items-center justify-between gap-3 font-serif text-[15px] text-mathua-primary">
                <span>{f.q}</span>
                <span className="font-mono text-mathua-muted transition-transform group-open:rotate-45">+</span>
              </summary>
              <p className="mt-2 text-[13px] leading-relaxed text-mathua-secondary">{f.a}</p>
            </details>
          ))}
        </div>
      </div>
      <div className="min-w-0">
        <div className="section-label">( 08 · How it is extended. )</div>
        <p className="mt-4 text-[13px] leading-relaxed text-mathua-secondary">
          Every concept is a JSON node. Every problem is a Go generator function. Every
          contribution goes through a graph validator that rejects cycles and orphaned nodes
          automatically.
        </p>
        <pre className="mt-4 overflow-x-auto overscroll-x-contain whitespace-pre border border-mathua-border bg-mathua-primary p-4 font-mono text-[11px] leading-relaxed text-mathua-bg sm:p-5 dark:bg-mathua-code dark:text-mathua-primary">
{`{
  "id":                "arith.add.multi",
  "label":             "Multi-digit addition",
  "domain":            "arithmetic",
  "prerequisites":     ["arith.add.single",
                        "arith.add.carry"],
  "mastery_threshold": {
    "streak":          5,
    "avg_time_seconds": 8
  }
}`}
        </pre>
        <Link
          href="/docs/contributing"
          className="mt-4 inline-block font-mono text-[11px] text-mathua-muted transition-colors hover:text-mathua-blue"
        >
          Read the contributing guide →
        </Link>
      </div>
    </section>
  )
}
