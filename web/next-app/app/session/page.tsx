'use client'

import Link from 'next/link'
import { Suspense } from 'react'
import { useSearchParams } from 'next/navigation'
import Header from '../../components/Header'
import BottomTabs from '../../components/BottomTabs'
import SectionHeader from '../../components/SectionHeader'
import Footer from '../../components/Footer'
import AsciiDivider from '../../components/AsciiDivider'

function SessionContent() {
  const searchParams = useSearchParams()
  const concept = searchParams.get('concept')
  return (
    <>
      <Header />
      <div className="max-w-container mx-auto px-4 sm:px-6 pb-[calc(80px+env(safe-area-inset-bottom))] lg:pb-0 overflow-x-hidden min-w-0">
        <section className="pt-8 min-w-0 overflow-hidden">
          <span className="flex mb-4">
            <Link href="/" className="text-mathua-secondary text-sm hover:text-mathua-primary">
              ← Back
            </Link>
          </span>

          <SectionHeader label="Start" title="Where would you like to start?" />
          <p className="text-mathua-secondary text-sm leading-relaxed text-center max-w-[600px] mx-auto mt-4">
            Take a diagnostic test to get a recommendation on where to start, or jump straight into Study.
          </p>
          {concept && (
            <div className="max-w-2xl mx-auto mt-6 border border-mathua-blue bg-mathua-surface p-4 flex flex-col sm:flex-row items-center justify-between gap-3">
              <span className="font-mono text-xs text-mathua-primary min-w-0 truncate">
                Continue with {concept}
              </span>
              <Link
                href={`/concept?id=${encodeURIComponent(concept)}`}
                className="shrink-0 border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white px-6 py-2 font-mono text-xs min-h-[36px] inline-flex items-center justify-center"
              >
                Open concept →
              </Link>
            </div>
          )}

          <div className="max-w-2xl mx-auto mt-10 grid grid-cols-1 sm:grid-cols-2 gap-4">
            {/* Option 1 — Diagnostic (Recommended) */}
            <Link
              href="/onboard"
              className="group relative flex flex-col border-2 border-mathua-blue bg-mathua-surface p-6 text-left hover:bg-mathua-blue hover:text-white transition-colors"
            >
              <span className="absolute -top-3 left-4 bg-mathua-blue text-white px-2 py-0.5 font-mono text-[10px] uppercase tracking-wider">
                Recommended
              </span>
              <span className="font-mono text-[11px] uppercase tracking-wider text-mathua-blue group-hover:text-white/80">
                Option 1
              </span>
              <h3 className="font-serif text-xl font-medium text-mathua-primary group-hover:text-white mt-1">
                Take a diagnostic test
              </h3>
              <p className="font-mono text-xs text-mathua-secondary group-hover:text-white/80 mt-2 leading-relaxed">
                A diagnostic test to find your knowledge frontier and get a personalized starting recommendation.
              </p>
              <span className="mt-4 inline-flex items-center font-mono text-xs text-mathua-blue group-hover:text-white">
                Start diagnostic test →
              </span>
            </Link>

            {/* Option 2 — Study */}
            <Link
              href="/study"
              className="flex flex-col border border-mathua-border bg-mathua-surface p-6 text-left hover:border-mathua-blue hover:text-mathua-blue transition-colors"
            >
              <span className="font-mono text-[11px] uppercase tracking-wider text-mathua-muted">
                Option 2
              </span>
              <h3 className="font-serif text-xl font-medium text-mathua-primary mt-1">
                Go to Study
              </h3>
              <p className="font-mono text-xs text-mathua-secondary mt-2 leading-relaxed">
                Browse the corpus directly — lessons first, then practice. Progress and XP are still recorded.
              </p>
              <span className="mt-4 inline-flex items-center font-mono text-xs text-mathua-blue">
                Open Study →
              </span>
            </Link>
          </div>

          <p className="text-mathua-muted text-xs font-mono text-center mt-6">
            You can switch anytime. The diagnostic test never deletes progress.
          </p>
        </section>

        <AsciiDivider pattern="wave" />
        <Footer />
      </div>
      <BottomTabs />
    </>
  )
}

export default function SessionPage() {
  return (
    <Suspense>
      <SessionContent />
    </Suspense>
  )
}
