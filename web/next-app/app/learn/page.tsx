'use client'

import { Suspense, useEffect, useState } from 'react'
import { useSearchParams } from 'next/navigation'
import Link from 'next/link'
import Header from '../../components/Header'
import BottomTabs from '../../components/BottomTabs'
import Footer from '../../components/Footer'
import LearnStepper from '../../components/LearnStepper'
import NextUpCard from '../../components/NextUpCard'
import type { Shelf } from '../../lib/nextUp'
import { fetchShelf } from '../../lib/recommendations'

/**
 * The learner's entry point: whatever the engine recommends, rendered as one card plus the
 * alternatives they may pick instead.
 *
 * This used to fetch five endpoints and rank the result in the browser. It now asks once. The
 * engine chooses; the learner disposes.
 */
function LearnEntry() {
  const [shelf, setShelf] = useState<Shelf | null>(null)

  useEffect(() => {
    let cancelled = false
    fetchShelf()
      .then(s => { if (!cancelled) setShelf(s) })
      .catch(() => { if (!cancelled) setShelf(null) })
    return () => { cancelled = true }
  }, [])

  const loading = shelf === null

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
  const returnTo = searchParams.get('return') || undefined

  if (!concept) {
    return <LearnEntry />
  }

  return (
    <div className="mt-6 mb-16 px-4">
      <div className="max-w-2xl mx-auto mb-4 flex items-center justify-between">
        <Link href="/profile" className="font-mono text-xs text-mathua-secondary hover:text-mathua-blue">← Profile</Link>
        {/* `from` lets the reference page offer a way back to this question
            instead of only sending the learner onwards into Learn. */}
        <Link href={`/study?concept=${encodeURIComponent(concept)}&from=${encodeURIComponent(concept)}`} className="font-mono text-[11px] text-mathua-muted hover:text-mathua-blue">Reference</Link>
      </div>
      <LearnStepper conceptId={concept} returnTo={returnTo} />
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
