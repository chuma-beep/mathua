'use client'

import { Suspense } from 'react'
import Link from 'next/link'
import Header from '../../components/Header'
import BottomTabs from '../../components/BottomTabs'
import Footer from '../../components/Footer'
import PlanEditor from '../../components/PlanEditor'

// Learning planner: destination + (daily effort XOR deadline) → estimate.
// Dates render as months (estimates, never promises); quiz eligibility and
// mastery are independent of every number on this page. The editor is shared
// with Settings (compact there); this page keeps the full estimate display.
function PlanContent() {
  return (
    <div className="max-w-2xl mx-auto mt-6 mb-16 px-4">
      <p className="font-mono text-[11px] uppercase tracking-wider text-mathua-muted">Learning planner</p>
      <h1 className="font-serif text-2xl text-mathua-primary mt-1">Plan your learning</h1>
      <PlanEditor />
      <Link href="/profile" className="mt-4 inline-block font-mono text-xs text-mathua-blue">← Back to profile</Link>
    </div>
  )
}

export default function PlanPage() {
  return (
    <>
      <Header />
      <div className="pt-20 lg:pt-0">
        <Suspense fallback={<div className="p-6 font-mono text-xs text-mathua-muted">Loading…</div>}>
          <PlanContent />
        </Suspense>
      </div>
      <BottomTabs />
      <Footer />
    </>
  )
}
