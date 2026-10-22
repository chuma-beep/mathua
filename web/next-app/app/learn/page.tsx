'use client'

import { Suspense, useEffect, useState } from 'react'
import { useSearchParams } from 'next/navigation'
import Link from 'next/link'
import Header from '../../components/Header'
import BottomTabs from '../../components/BottomTabs'
import Footer from '../../components/Footer'
import LearnStepper from '../../components/LearnStepper'
import LessonCredits from '../../components/LessonCredits'
import NextUpCard from '../../components/NextUpCard'
import type { Shelf } from '../../lib/nextUp'
import { fetchShelf } from '../../lib/recommendations'
import { parseReturnStack } from '../../lib/learnReturn'

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

  if (!concept) {
    return <LearnEntry />
  }

  return (
    <div className="mt-6 mb-16 px-4">
      <div className="max-w-2xl mx-auto mb-4">
        <Link href="/profile" className="font-mono text-xs text-mathua-secondary hover:text-mathua-blue">← Profile</Link>
      </div>
      <LearnStepper conceptId={concept} returnStack={returnStack} />
      <div className="max-w-2xl mx-auto">
        <LessonCredits conceptId={concept} />
      </div>
    </div>
  )
}

export default function LearnPage() {
  return (
    <>
      <Header />
      <div className="pt-[var(--chrome-top)] lg:pt-0">
        <Suspense fallback={<div className="p-6 font-mono text-xs text-mathua-muted">Loading…</div>}>
          <LearnContent />
        </Suspense>
      </div>
      <BottomTabs />
      <Footer />
    </>
  )
}
