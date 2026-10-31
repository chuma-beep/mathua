'use client'

import { Suspense, useEffect, useState } from 'react'
import { useSearchParams } from 'next/navigation'
import Link from 'next/link'
import Footer from '../../components/Footer'
import AppShell from '../../components/app-shell'
import LearnStepper from '../../components/LearnStepper'
import LessonCredits from '../../components/LessonCredits'
import NextUpCard from '../../components/NextUpCard'
import type { Shelf } from '../../lib/nextUp'
import { fetchShelf } from '../../lib/recommendations'
import { parseReturnStack } from '../../lib/learnReturn'
import { ensureGuestToken } from '../../lib/auth'

/**
 * The learner's entry point: whatever the engine recommends, rendered as one card plus the
 * alternatives they may pick instead.
 *
 * This used to fetch five endpoints and rank the result in the browser. It now asks once. The
 * engine chooses; the learner disposes.
 */
function LearnEntry() {
  const [shelf, setShelf] = useState<Shelf | null>(null)
  const [failed, setFailed] = useState(false)
  const [attempt, setAttempt] = useState(0)
  const [hasDiag, setHasDiag] = useState(false)

  useEffect(() => {
    try {
      setHasDiag(!!sessionStorage.getItem('mathua_diag_session_onboard'))
    } catch {
      setHasDiag(false)
    }
  }, [])

  useEffect(() => {
    let cancelled = false
    setFailed(false)
    fetchShelf()
      .then(s => { if (!cancelled) setShelf(s) })
      .catch(() => { if (!cancelled) { setShelf(null); setFailed(true) } })
    return () => { cancelled = true }
  }, [attempt])

  const loading = shelf === null && !failed

  // The failure branch used to be unreachable: getNext answered an outage with
  // "All learned", and this catch left shelf null, which read as loading and
  // spun forever. Neither is true. The queue was not read, so say that.
  if (failed) {
    return (
      <div className="mt-6 mb-16 px-4">
        <div className="max-w-2xl mx-auto border border-mathua-border bg-mathua-surface p-4">
          <p className="text-mathua-primary text-sm">Could not load what to study next.</p>
          <p className="mt-1 font-mono text-xs text-mathua-muted">
            The engine could not be reached, so your queue was not read. Nothing was lost.
          </p>
          <div className="mt-3 flex flex-wrap items-center gap-3">
            <button
              type="button"
              onClick={() => setAttempt(n => n + 1)}
              className="font-mono text-xs px-3 py-1.5 border border-mathua-border text-mathua-primary hover:border-mathua-blue"
            >
              Try again
            </button>
            <Link href="/domains" className="font-mono text-xs text-mathua-blue hover:underline">
              Browse domains
            </Link>
          </div>
        </div>
      </div>
    )
  }

  if (loading) {
    return <div className="p-6 font-mono text-xs text-mathua-muted">Loading…</div>
  }

  return (
    <div className="mt-6 mb-16 px-4">
      <div className="max-w-2xl mx-auto">
        {hasDiag && (
          <div className="mb-4 border border-mathua-blue bg-mathua-surface p-4">
            <p className="font-mono text-xs text-mathua-blue">Diagnostic test in progress</p>
            <p className="mt-1 font-mono text-[11px] text-mathua-secondary">
              You have an active diagnostic session. Resume it to find your starting point.
            </p>
            <div className="mt-3">
              <Link
                href="/onboard"
                className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue-faint px-4 py-1.5 font-mono text-xs inline-flex items-center"
              >
                Resume diagnostic test →
              </Link>
            </div>
          </div>
        )}
        <NextUpCard shelf={shelf} />
        <Link href="/profile" className="mt-4 inline-block font-mono text-xs text-mathua-secondary hover:text-mathua-blue">← Profile</Link>
      </div>
    </div>
  )
}

function LearnContent() {
  const searchParams = useSearchParams()
  const concept = searchParams.get('concept') || ''
  const returnStack = parseReturnStack(searchParams)
  // A guest who lands here from a direct link never pressed "Get started", so
  // nothing minted their guest token. Without it the practice set is served
  // anonymously, the server anchors it to nobody, and every answer comes back
  // 400 "student identity required", which the stepper can only show as the
  // server being unreachable. Mint it before the first question is served.
  // A registered user or a stored guest token resolves at once. A mint that
  // fails, or stalls until its ten-second timeout, falls through to the page
  // as it was before, so an outage costs a wait and never a blank lesson.
  const [identityReady, setIdentityReady] = useState(false)
  useEffect(() => {
    let cancelled = false
    ensureGuestToken().finally(() => { if (!cancelled) setIdentityReady(true) })
    return () => { cancelled = true }
  }, [])

  if (!concept) {
    return <LearnEntry />
  }

  if (!identityReady) {
    return <div className="p-6 font-mono text-xs text-mathua-muted">Loading…</div>
  }

  return (
    <div className="mt-6 mb-16 px-4">
      <div className="max-w-2xl mx-auto mb-4">
        <Link href="/profile" className="font-mono text-xs text-mathua-secondary hover:text-mathua-blue">← Profile</Link>
      </div>
      {/* Key the stepper on the concept. A Continue link is a client-side
          navigation from /learn?concept=A to /learn?concept=B, and without a
          key React reuses the instance: the old lesson, readiness and mastered
          banner stay on screen under the new URL, so the link appears to do
          nothing. Remounting is also what lets the loaders run for B. */}
      <LearnStepper key={concept ?? 'none'} conceptId={concept} returnStack={returnStack} />
      <div className="max-w-2xl mx-auto">
        <LessonCredits conceptId={concept} />
      </div>
    </div>
  )
}

export default function LearnPage() {
  return (
    <AppShell>
      <Suspense fallback={<div className="p-6 font-mono text-xs text-mathua-muted">Loading…</div>}>
        <LearnContent />
      </Suspense>
      <Footer />
    </AppShell>
  )
}
