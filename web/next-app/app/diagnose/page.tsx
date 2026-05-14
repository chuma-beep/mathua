import type { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'Diagnostic Test — Mathua',
  description: 'A computerized adaptive test that locates your position on the concept graph using as few questions as possible.',
}

'use client'

import { useTheme } from '../../hooks/useTheme'
import Link from 'next/link'
import Header from '../../components/Header'
import SectionHeader from '../../components/SectionHeader'
import D2Diagram from '../../components/D2Diagram'
import Footer from '../../components/Footer'
import AsciiDivider from '../../components/AsciiDivider'

export default function DiagnosePage() {
  const { theme, mounted } = useTheme()

  if (!mounted) return <div style={{ background: 'var(--bg)', minHeight: '100vh' }} />

  return (
    <>
      <Header />
      <div className="max-w-container mx-auto px-6 max-sm:px-4">
      <section className="pt-8">
        <span className="flex mb-4">
          <Link href="/" className="text-mathua-secondary text-sm hover:text-mathua-primary">
            ← Back
          </Link>
        </span>

        <SectionHeader label="Diagnostic Test" title="Find your knowledge frontier" />

        <p className="text-mathua-secondary text-sm leading-relaxed text-center max-w-[600px] mx-auto mt-4">
          A Computerised Adaptive Testing session that locates your position in the concept graph
           using as few questions as possible: typically 20–35 questions across the 284-concept
          space.
        </p>
      </section>

      <AsciiDivider pattern="wave" />

      {/* How it works */}
      <section className="py-20 max-sm:py-12">
        <SectionHeader label="Method" title="How the diagnostic works" />

        <div className="flex justify-center mt-8">
          <D2Diagram name="cat-diagnostic" theme={theme} />
        </div>

        <ul className="max-w-[600px] mx-auto mt-8 space-y-3 list-none">
          {[
            'The concept graph is sorted topologically: the diagnostic starts at the midpoint.',
            'Correct answers within the time limit move the probe forward toward harder concepts.',
            'Incorrect or slow answers move backward toward foundational material.',
            'After 3 consecutive correct answers in a region, the frontier is considered located.',
            'The diagnostic records a mastery estimate for every concept passed through.',
          ].map((step) => (
            <li
              key={step}
              className="text-mathua-secondary text-[0.95rem] leading-[1.7] pl-9 relative before:content-[counter(step)] before:absolute before:left-0 before:text-mathua-blue before:font-mono before:text-[13px]"
              style={{ counterIncrement: 'step-counter 1' }}
            >
              {step}
            </li>
          ))}
        </ul>
        <style>{`ul { counter-reset: step-counter; }`}</style>
      </section>

      <AsciiDivider pattern="dash" />

      {/* Actions */}
      <section className="py-20 max-sm:py-12 text-center">
        <SectionHeader label="Ready" title="Begin your diagnostic" />

        <p className="text-mathua-secondary text-sm leading-relaxed max-w-[500px] mx-auto mt-4 mb-8">
          The diagnostic takes 20–35 questions. You can retake it at any time: retaking never
          deletes progress, and Mathua always keeps the most optimistic estimate.
        </p>

        <div className="flex gap-3 justify-center max-sm:flex-col max-sm:items-center">
          <button className="bg-mathua-blue text-white hover:bg-mathua-blue-hover rounded-md h-10 px-6 text-[13px] font-medium transition-colors">
            Start Diagnostic
          </button>
          <button className="border border-mathua-border-strong text-mathua-primary hover:border-mathua-blue hover:text-mathua-blue rounded-md h-10 px-6 text-[13px] font-medium transition-colors">
            Read the full algorithm
          </button>
        </div>
      </section>

      <Footer />
    </div>
    </>
  )
}
