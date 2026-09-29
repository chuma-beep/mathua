'use client'

import { Suspense } from 'react'
import { useSearchParams } from 'next/navigation'
import Link from 'next/link'
import Header from '../../components/Header'
import BottomTabs from '../../components/BottomTabs'
import Footer from '../../components/Footer'
import LearnStepper from '../../components/LearnStepper'

function LearnContent() {
  const searchParams = useSearchParams()
  const concept = searchParams.get('concept') || ''
  const returnTo = searchParams.get('return') || undefined

  if (!concept) {
    return (
      <div className="max-w-2xl mx-auto mt-8 p-6 border border-mathua-border bg-mathua-surface">
        <p className="font-mono text-xs text-mathua-muted">No concept selected.</p>
        <Link href="/profile" className="mt-3 inline-block font-mono text-xs text-mathua-blue">← Back to profile</Link>
      </div>
    )
  }

  return (
    <div className="mt-6 mb-16 px-4">
      <div className="max-w-2xl mx-auto mb-4 flex items-center justify-between">
        <Link href="/profile" className="font-mono text-xs text-mathua-secondary hover:text-mathua-blue">← Next up</Link>
        <Link href={`/study?concept=${encodeURIComponent(concept)}`} className="font-mono text-[11px] text-mathua-muted hover:text-mathua-blue">Reference</Link>
      </div>
      <LearnStepper conceptId={concept} returnTo={returnTo} />
    </div>
  )
}

export default function LearnPage() {
  return (
    <>
      <Header />
      <div className="pt-20 lg:pt-0">
        <Suspense fallback={<div className="p-6 font-mono text-xs text-mathua-muted">Loading…</div>}>
          <LearnContent />
        </Suspense>
      </div>
      <BottomTabs />
      <Footer />
    </>
  )
}
