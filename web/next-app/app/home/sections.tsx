'use client'

import Link from 'next/link'
import { useEffect, useRef, useState } from 'react'
import { useRouter } from 'next/navigation'
import { LOGIN_HERO } from '../../lib/loginHero'
import { FEATURED_FIGURES } from '../../lib/featuredFigures'
import VisualFigure from '../../components/visual/VisualFigure'
import { API_BASE, getConfig, googleOneTap } from '../../lib/api'
import { establishSession } from '../../lib/auth'
import { conceptCount, domainCounts, domainLabels, domainOrder } from './data'

// The homepage keeps the drawing and the type of the rest of the app. The
// landing page is the one place with a painting, so the buttons on it carry a
// white hairline and every other control keeps the themed border.

const CTA_PRIMARY =
  'inline-flex min-h-[44px] items-center justify-center bg-mathua-blue px-6 py-3 text-sm font-medium text-white transition-colors hover:bg-mathua-blue-hover focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-mathua-blue'

const CTA_GHOST_ON_ART =
  'group relative inline-flex min-h-[44px] items-center justify-center overflow-hidden border border-white/45 bg-transparent px-6 py-3 text-sm font-medium text-white transition-colors hover:border-white focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-white lg:border-mathua-border lg:text-mathua-primary lg:hover:border-mathua-blue lg:focus-visible:outline-mathua-blue'

function GhostLabel({ children }: { children: React.ReactNode }) {
  return (
    <>
      <span
        aria-hidden="true"
        className="absolute inset-0 origin-left scale-x-0 bg-mathua-blue transition-transform duration-300 ease-out group-hover:scale-x-100 motion-reduce:transition-none"
      />
      <span className="relative group-hover:text-white">{children}</span>
    </>
  )
}

interface GoogleAccounts {
  initialize: (options: {
    client_id: string
    callback: (response: { credential?: string }) => void
    auto_select?: boolean
    cancel_on_tap_outside?: boolean
  }) => void
  prompt: () => void
}

interface GoogleWindow extends Window {
  google?: { accounts?: { id?: GoogleAccounts } }
}

/**
 * The Learn call to action opens this panel instead of navigating straight
 * away. Google One Tap is destroyed by a navigation and can be suppressed by
 * FedCM, so the prompt has to run on the page that shows it. The panel always
 * keeps a guest link, which is the path that cannot fail.
 */
function StartPanel({ onClose }: { onClose: () => void }) {
  const router = useRouter()
  const [prompted, setPrompted] = useState(false)
  const [busy, setBusy] = useState(false)
  const started = useRef(false)

  useEffect(() => {
    if (started.current) return
    started.current = true
    let cancelled = false

    async function start() {
      try {
        const cfg = await getConfig()
        const clientID = (cfg?.google_client_id || '').trim()
        if (!clientID || cancelled) return
        const src = 'https://accounts.google.com/gsi/client'
        if (!document.querySelector(`script[src="${src}"]`)) {
          await new Promise<void>((resolve, reject) => {
            const script = document.createElement('script')
            script.src = src
            script.async = true
            script.defer = true
            script.onload = () => resolve()
            script.onerror = () => reject(new Error('gis unavailable'))
            document.head.appendChild(script)
          })
        }
        if (cancelled) return
        const g = (window as GoogleWindow).google
        if (!g?.accounts?.id) return
        g.accounts.id.initialize({
          client_id: clientID,
          callback: async (resp) => {
            if (!resp?.credential) return
            setBusy(true)
            try {
              const data = await googleOneTap(resp.credential)
              await establishSession(data.token, {
                student_id: data.student_id,
                name: data.name,
                username: '',
                concepts_mastered: 0,
                current_streak: 0,
                level: 'Novice',
                diagnostic_completed: data.diagnostic_completed,
              })
              router.push('/learn')
            } catch {
              setBusy(false)
            }
          },
          auto_select: false,
          cancel_on_tap_outside: true,
        })
        g.accounts.id.prompt()
        setPrompted(true)
      } catch {
        // One Tap is optional. The two links below are the working paths.
      }
    }

    void start()
    return () => {
      cancelled = true
    }
  }, [router])

  return (
    <div className="reveal mt-6 border border-mathua-border bg-mathua-bg p-4 sm:p-5" data-testid="start-panel">
      <p className="font-mono text-[10px] uppercase tracking-wide text-mathua-muted">Start learning</p>
      <p className="mt-2 max-w-[52ch] text-[13px] leading-relaxed text-mathua-secondary">
        {busy
          ? 'Signing you in…'
          : prompted
            ? 'Choose an account in the Google prompt, or continue as a guest.'
            : 'Continue with Google to keep progress across devices, or start as a guest.'}
      </p>
      <div className="mt-3 flex flex-col gap-2 sm:flex-row sm:items-center">
        <a
          href={`${API_BASE}/api/auth/google/login?return=${encodeURIComponent('/learn')}`}
          className="inline-flex min-h-[44px] items-center justify-center border border-mathua-blue px-5 text-sm font-medium text-mathua-blue transition-colors hover:bg-mathua-blue-faint"
        >
          Continue with Google
        </a>
        <Link
          href="/learn"
          className="inline-flex min-h-[44px] items-center justify-center px-5 text-sm font-medium text-mathua-secondary underline-offset-4 hover:text-mathua-primary hover:underline"
        >
          Continue as a guest
        </Link>
        <button
          type="button"
          onClick={onClose}
          className="ml-auto min-h-[44px] px-2 font-mono text-[11px] text-mathua-muted hover:text-mathua-primary"
        >
          Close
        </button>
      </div>
      <p className="mt-3 font-mono text-[10px] text-mathua-muted">No account needed. Progress stays on this device.</p>
    </div>
  )
}

export function HeroSection() {
  const [panelOpen, setPanelOpen] = useState(false)
  return (
    <section className="group/hero relative isolate overflow-hidden border-b border-mathua-division lg:grid lg:grid-cols-2 lg:items-stretch lg:gap-10">
      <div
        aria-hidden="true"
        className="hero-art-in absolute inset-x-0 bottom-0 top-4 overflow-hidden lg:relative lg:bottom-auto lg:top-auto lg:order-2 lg:mt-10 lg:min-h-[480px]"
      >
        <img
          src={LOGIN_HERO.file}
          width={LOGIN_HERO.width}
          height={LOGIN_HERO.height}
          loading="eager"
          decoding="async"
          className="absolute inset-0 size-full object-cover object-center transition-transform duration-700 ease-out group-hover/hero:scale-[1.03] motion-reduce:transition-none motion-reduce:group-hover/hero:scale-100"
        />
        <div className="absolute inset-0 bg-gradient-to-b from-black/80 via-black/65 to-black/55 lg:hidden" />
      </div>
      <div className="relative max-w-[72ch] px-4 py-16 [text-shadow:0_1px_14px_rgba(0,0,0,0.45)] md:py-24 lg:order-1 lg:px-0 lg:[text-shadow:none]">
        <div className="section-label">( 01 · Title Block )</div>
        <h1
          className="reveal mt-5 max-w-[18ch] text-balance font-serif text-[clamp(2.5rem,6vw,4.5rem)] font-normal leading-[1.05] sm:leading-[0.95] tracking-tight text-white lg:text-mathua-primary"
          style={{ transitionDelay: '80ms' }}
        >
          Don’t just finish mathematics. Understand it.
        </h1>
        <p
          className="reveal mt-6 max-w-[46ch] text-pretty text-base sm:text-[15px] leading-relaxed text-white/85 lg:text-mathua-secondary"
          style={{ transitionDelay: '160ms' }}
        >
          Mathua adapts to what you know, strengthens what you don’t, and makes sure you master the foundations
          before moving forward.
        </p>
        <div className="reveal mt-8 flex flex-col items-start gap-3 sm:flex-row sm:items-center" style={{ transitionDelay: '240ms' }}>
          <button
            type="button"
            onClick={() => setPanelOpen((open) => !open)}
            aria-expanded={panelOpen}
            aria-controls="start-panel"
            className={CTA_PRIMARY}
          >
            Start learning for free
          </button>
          <a href="#how-learning-works" className={CTA_GHOST_ON_ART}>
            <GhostLabel>See how it works</GhostLabel>
          </a>
        </div>
        <p className="reveal mt-4 font-mono text-[11px] text-mathua-blue" style={{ transitionDelay: '320ms' }}>
          No account needed. Progress stays on this device.
        </p>
        <div id="start-panel">
          {panelOpen && <StartPanel onClose={() => setPanelOpen(false)} />}
        </div>
      </div>
    </section>
  )
}

const PRACTICE = {
  question: 'What is 7 × 8?',
  options: [54, 56, 58, 64],
  answer: 56,
  working: '7 × 8 = 56. The multiplication table row for 7 lands on 56 at the eighth step.',
}

/** A real question, answered locally, in the same markup as a lesson card. */
export function PracticePreview() {
  const [picked, setPicked] = useState<number | null>(null)
  const correct = picked === PRACTICE.answer
  return (
    <section data-reveal className="border-b border-mathua-division py-16" aria-labelledby="practice-heading">
      <div className="section-label">( 02 · A session, not a screenshot )</div>
      <div className="mt-4 grid items-start gap-6 lg:grid-cols-[minmax(0,1fr)_260px] lg:gap-10">
        <div className="border border-mathua-border bg-mathua-surface p-5 sm:p-6">
          <div className="flex items-center justify-between gap-3 font-mono text-[10px] uppercase text-mathua-muted">
            <span>Practice · arithmetic</span>
            <span>Example</span>
          </div>
          <h2 id="practice-heading" className="mt-3 font-serif text-xl text-mathua-primary sm:text-2xl">
            {PRACTICE.question}
          </h2>
          <div className="mt-4 grid grid-cols-2 gap-2">
            {PRACTICE.options.map((option) => {
              const chosen = picked === option
              const isAnswer = option === PRACTICE.answer
              return (
                <button
                  key={option}
                  type="button"
                  onClick={() => setPicked(option)}
                  className={
                    'min-h-[44px] border px-4 font-mono text-sm transition-colors ' +
                    (picked === null
                      ? 'border-mathua-border text-mathua-primary hover:border-mathua-blue hover:text-mathua-blue'
                      : chosen && isAnswer
                        ? 'border-mathua-blue bg-mathua-blue-faint text-mathua-blue'
                        : chosen
                          ? 'border-mathua-red text-mathua-red'
                          : isAnswer
                            ? 'border-mathua-blue text-mathua-blue'
                            : 'border-mathua-border text-mathua-muted')
                  }
                >
                  {option}
                </button>
              )
            })}
          </div>
          <div aria-live="polite" className="mt-4 min-h-[48px]">
            {picked !== null && (
              <p className={'text-[13px] leading-relaxed ' + (correct ? 'text-mathua-primary' : 'text-mathua-secondary')}>
                <span className="font-mono text-[10px] uppercase text-mathua-muted">{correct ? 'Correct · ' : 'Not yet · '}</span>
                {PRACTICE.working}
              </p>
            )}
          </div>
        </div>
        <div className="border border-mathua-border bg-mathua-surface p-5">
          <p className="font-mono text-[10px] uppercase text-mathua-muted">Concept progress</p>
          <p className="mt-2 text-sm text-mathua-primary">Multiplication facts</p>
          <div className="mt-4 h-1.5 w-full bg-mathua-division" role="img" aria-label="Mastery 4 of 5 checks">
            <div
              className={'h-full bg-mathua-blue transition-[width] duration-300 ' + (correct ? 'w-full' : 'w-4/5')}
            />
          </div>
          <p className="mt-2 font-mono text-[10px] text-mathua-muted">{correct ? '5 of 5 checks' : '4 of 5 checks'}</p>
          <p className="mt-4 text-[12px] leading-relaxed text-mathua-secondary">
            A concept is mastered from recent evidence: accuracy, speed, variety and time since the last attempt.
          </p>
        </div>
      </div>
    </section>
  )
}

const STEPS = [
  { title: 'Learn', body: 'A short lesson introduces the idea and works one example. The figure and the algebra stay on the same screen.' },
  { title: 'Practice', body: 'Problems are generated for the concept, so you cannot memorise one set of answers and move on.' },
  { title: 'Master', body: 'Mastery needs a run of correct answers at the target speed, across different forms of the same idea.' },
  { title: 'Review', body: 'What you mastered comes back on a schedule, before it fades, rather than all at once.' },
]

export function StepsSection() {
  return (
    <section id="how-learning-works" data-reveal className="border-b border-mathua-division py-16">
      <div className="section-label">( 03 · How learning works )</div>
      <div className="mt-4 grid gap-px border border-mathua-division bg-mathua-division sm:grid-cols-2 lg:grid-cols-4">
        {STEPS.map((step, index) => (
          <div key={step.title} className="bg-mathua-bg p-6">
            <span className="font-mono text-xs text-mathua-blue">{String(index + 1).padStart(2, '0')}</span>
            <h2 className="mt-3 font-serif text-lg text-mathua-primary">{step.title}</h2>
            <p className="mt-2 text-[13px] leading-relaxed text-mathua-secondary">{step.body}</p>
          </div>
        ))}
      </div>
      <Link href="/how-it-works" className="mt-4 inline-block font-mono text-[11px] text-mathua-muted hover:text-mathua-blue">
        How the engine decides →
      </Link>
    </section>
  )
}

export function FiguresSection() {
  return (
    <section data-reveal className="border-b border-mathua-division py-16">
      <div className="section-label">( 04 · Interactive mathematics )</div>
      <p className="mt-3 max-w-[62ch] text-[15px] leading-relaxed text-mathua-secondary">
        Some ideas are easier to see than to read. Change the angle, the field or the slice and the result follows.
      </p>
      <div className="mt-6 grid gap-4 lg:grid-cols-3">
        {FEATURED_FIGURES.map((figure) => (
          <article key={figure.type} className="flex flex-col border border-mathua-border bg-mathua-surface p-4">
            <h2 className="font-serif text-lg text-mathua-primary">{figure.title}</h2>
            <p className="mt-1 text-[12px] leading-relaxed text-mathua-secondary">{figure.caption}</p>
            <div className="mt-3">
              <VisualFigure visual={figure} />
            </div>
          </article>
        ))}
      </div>
      <Link href="/announcements" className="mt-4 inline-block font-mono text-[11px] text-mathua-muted hover:text-mathua-blue">
        More about the interactive figures →
      </Link>
    </section>
  )
}

const ADAPTIVE = [
  {
    title: 'Prerequisites are enforced',
    body: 'A concept unlocks only when the ideas under it are mastered. The map is a check, not a suggestion.',
  },
  {
    title: 'Mastery is evidence',
    body: 'Correctness alone is not enough. Speed, variety across question forms and the time since the last attempt all count.',
  },
  {
    title: 'Review is scheduled',
    body: 'Mastered ideas return on a spacing schedule built from how quickly each one was answered.',
  },
]

export function AdaptiveSection() {
  return (
    <section data-reveal className="border-b border-mathua-division py-16">
      <div className="section-label">( 05 · Adaptive learning )</div>
      <div className="mt-4 grid gap-px border border-mathua-division bg-mathua-division md:grid-cols-3">
        {ADAPTIVE.map((item) => (
          <div key={item.title} className="bg-mathua-bg p-6">
            <h2 className="font-serif text-lg text-mathua-primary">{item.title}</h2>
            <p className="mt-2 text-[13px] leading-relaxed text-mathua-secondary">{item.body}</p>
          </div>
        ))}
      </div>
    </section>
  )
}

export function CoverageSection() {
  return (
    <section data-reveal className="border-b border-mathua-division py-16">
      <div className="flex flex-wrap items-baseline justify-between gap-3">
        <div className="section-label">( 06 · What Mathua covers )</div>
        <Link href="/domains" className="font-mono text-[11px] text-mathua-muted hover:text-mathua-blue">
          All {conceptCount} concepts by subject →
        </Link>
      </div>
      <div className="mt-4 grid gap-px border border-mathua-division bg-mathua-division sm:grid-cols-2 lg:grid-cols-3">
        {domainOrder.map((domain) => (
          <Link
            key={domain}
            href="/domains"
            className="group flex items-baseline justify-between gap-3 bg-mathua-bg px-4 py-3 transition-colors hover:bg-mathua-surface"
          >
            <span className="text-sm text-mathua-secondary group-hover:text-mathua-primary">{domainLabels[domain]}</span>
            <span className="font-mono text-xs text-mathua-muted">{domainCounts[domain]}</span>
          </Link>
        ))}
      </div>
      <p className="mt-4 max-w-[62ch] text-[12px] leading-relaxed text-mathua-muted">
        Every concept has a lesson and generated problems. Nothing is a fixed answer key.
      </p>
    </section>
  )
}

export function OpenSourceSection() {
  return (
    <section data-reveal className="border-b border-mathua-division py-16">
      <div className="section-label">( 07 · Open source )</div>
      <div className="mt-4 grid gap-6 md:grid-cols-2">
        <p className="max-w-[52ch] text-[14px] leading-relaxed text-mathua-secondary">
          Mathua is MIT licensed. Lessons draw on OpenStax (CC BY-NC-SA 4.0) and Algebrica (CC BY-NC 4.0), and the
          engine, generators and interface are in one repository. You can read it, run it, or change it.
        </p>
        <div className="flex flex-col items-start gap-3 sm:flex-row sm:items-center">
          <a
            href="https://github.com/chuma-beep/mathua"
            target="_blank"
            rel="noreferrer"
            className="inline-flex min-h-[44px] items-center border border-mathua-border px-5 text-sm text-mathua-primary transition-colors hover:border-mathua-blue hover:text-mathua-blue"
          >
            View on GitHub
          </a>
          <Link href="/docs/contributing" className="font-mono text-[11px] text-mathua-muted hover:text-mathua-blue">
            Read the contributing guide →
          </Link>
        </div>
      </div>
    </section>
  )
}

export function ClosingCta() {
  return (
    <section data-reveal className="border-b border-mathua-division py-16">
      <div className="border border-mathua-border bg-mathua-primary p-8 text-mathua-bg dark:bg-mathua-code dark:text-mathua-primary">
        <div className="section-label text-mathua-bg/60 dark:text-mathua-muted">( 08 · Start here )</div>
        <h2 className="mt-3 font-serif text-2xl sm:text-3xl">Start with what you know.</h2>
        <p className="mt-2 max-w-[52ch] text-[14px] leading-relaxed opacity-80">
          A short diagnostic finds the edge of what you can do, and the plan starts there.
        </p>
        <div className="mt-6 flex flex-col gap-3 sm:flex-row">
          <Link
            href="/learn"
            className="inline-flex min-h-[44px] items-center justify-center bg-mathua-blue px-6 text-sm font-medium text-white transition-colors hover:bg-mathua-blue-hover"
          >
            Start learning for free
          </Link>
          <Link
            href="/login"
            className="inline-flex min-h-[44px] items-center justify-center border border-mathua-bg/40 px-6 text-sm font-medium transition-colors hover:border-mathua-bg dark:border-mathua-border dark:hover:border-mathua-primary"
          >
            Create an account
          </Link>
        </div>
      </div>
    </section>
  )
}

const FAQS = [
  {
    q: 'Do I need an account?',
    a: 'No. Guest practice works in the browser. An account keeps progress across devices and lets reminders reach you.',
  },
  {
    q: 'Why does speed count?',
    a: 'A mastered fact is one you can use without stopping to rebuild it. The check measures a run of correct answers at a target pace, not a single lucky answer.',
  },
  {
    q: 'Is it free?',
    a: 'Yes. Mathua is MIT licensed and has no paywall. You can also self-host it.',
  },
]

export function FaqSection() {
  return (
    <section data-reveal className="py-16">
      <div className="section-label">( 09 · Questions, answered. )</div>
      <div className="mt-4 max-w-[72ch]">
        {FAQS.map((item) => (
          <details key={item.q} className="group border-b border-mathua-division py-3">
            <summary className="flex cursor-pointer list-none items-center justify-between gap-3 font-serif text-[15px] text-mathua-primary">
              {item.q}
              <span aria-hidden="true" className="font-mono text-lg text-mathua-muted transition-transform group-open:rotate-45">
                +
              </span>
            </summary>
            <p className="mt-2 text-[13px] leading-relaxed text-mathua-secondary">{item.a}</p>
          </details>
        ))}
      </div>
    </section>
  )
}
