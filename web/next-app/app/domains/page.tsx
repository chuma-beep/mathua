'use client'

import { useEffect, useMemo, useState } from 'react'
import Link from 'next/link'
import Header from '../../components/Header'
import BottomTabs from '../../components/BottomTabs'
import SectionHeader from '../../components/SectionHeader'
import Footer from '../../components/Footer'
import Loading from '../../components/Loading'
import DomainTable from '../../components/DomainTable'
import { concepts } from '../../lib/conceptData'
import { buildDomainRows, byNeedsAttention } from '../../lib/domainRows'
import { getProgress, type ConceptProgress } from '../../lib/api'
import { getGuestId, getUserInfo } from '../../lib/auth'

// Per-domain progress, off the dashboard.
//
// This block used to live on /profile as `DomainProgress`, and two things moved it here. The
// first is that the dashboard answers "what should I do now", and a table of subject percentages
// is not that: it competed with the recommendation for attention without being one. The second is
// that the component hardcoded its own 15-entry domain list, so `precalculus` and
// `machine_learning` never appeared at all — 50 concepts invisible, with nothing on screen to say
// so. Both fixes live in lib/domainRows and lib/graphDomains.
//
// The locked count went with it. It was the figure Phase 3 removed from PositionBlock as the one
// number a learner can neither act on nor change today, and it grows as they learn, so it reads
// as decline. Here it would have appeared 17 times instead of once.

export default function DomainsPage() {
  const [progress, setProgress] = useState<Record<string, ConceptProgress>>({})
  const [loading, setLoading] = useState(true)
  const [expanded, setExpanded] = useState<string | null>(null)
  const [needsAttentionFirst, setNeedsAttentionFirst] = useState(false)

  useEffect(() => {
    // Guests have progress too: the corpus is browsable before signing in, and a page that
    // showed zeroes to every guest is a page nobody signs up from.
    const id = getUserInfo()?.student_id || getGuestId() || ''
    if (!id) {
      setLoading(false)
      return
    }
    getProgress(id)
      .then(setProgress)
      .catch(() => setProgress({}))
      .finally(() => setLoading(false))
  }, [])

  const rows = useMemo(() => {
    const built = buildDomainRows(concepts, progress)
    return needsAttentionFirst ? [...built].sort(byNeedsAttention) : built
  }, [progress, needsAttentionFirst])

  if (loading) return <Loading />

  const totals = rows.reduce(
    (acc, r) => ({
      concepts: acc.concepts + r.total,
      mastered: acc.mastered + r.mastered,
      due: acc.due + r.dueForReview,
    }),
    { concepts: 0, mastered: 0, due: 0 },
  )

  return (
    <>
      <Header />
      <main className="mx-auto w-full max-w-[820px] min-w-0 px-4 sm:px-6 py-8 sm:py-12 overflow-x-hidden">
        <SectionHeader label="Progress" title="Domains" />

        <p className="font-mono text-[11px] text-mathua-secondary mb-6">
          {totals.mastered} of {totals.concepts} concepts mastered across {rows.length} domains
          {totals.due > 0 && <> · {totals.due} due for review</>}.{' '}
          <Link href="/graph" className="text-mathua-blue hover:underline">
            Open the graph
          </Link>
        </p>

        <div className="mb-3 flex justify-end">
          <button
            type="button"
            onClick={() => setNeedsAttentionFirst((v) => !v)}
            aria-pressed={needsAttentionFirst}
            className="font-mono text-[10px] uppercase px-3 min-h-[36px] py-1.5 border border-mathua-border text-mathua-muted hover:text-mathua-blue hover:border-mathua-secondary transition-colors"
          >
            {needsAttentionFirst ? 'Curriculum order' : 'Needs attention first'}
          </button>
        </div>

        <DomainTable
          rows={rows}
          onToggle={(d) => setExpanded((cur) => (cur === d ? null : d))}
          expanded={expanded}
        />
      </main>
      <Footer />
      <BottomTabs />
    </>
  )
}
