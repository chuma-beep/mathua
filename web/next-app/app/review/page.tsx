'use client'

import Link from 'next/link'
import Header from '../../components/Header'
import BottomTabs from '../../components/BottomTabs'
import Footer from '../../components/Footer'
import ReviewHost from './ReviewHost'

export default function ReviewPage() {
  return (
    <>
      <Header />
      <div className="max-w-container mx-auto px-4 sm:px-6 pb-[calc(80px+env(safe-area-inset-bottom))] lg:pb-0 overflow-x-hidden min-w-0">
        <section className="pt-8 min-w-0 overflow-hidden">
          <span className="flex mb-4">
            <Link href="/profile" className="text-mathua-secondary text-sm hover:text-mathua-primary">
              ← Back
            </Link>
          </span>
          <ReviewHost />
          <p className="text-mathua-muted text-xs font-mono text-center mt-6">
            Reviews record progress and XP like practice.
          </p>
        </section>
      </div>
      <Footer />
      <BottomTabs />
    </>
  )
}
