'use client'

import { Suspense, useEffect, useMemo, useState } from 'react'
import { useSearchParams } from 'next/navigation'
import Link from 'next/link'
import Header from '../../components/Header'
import BottomTabs from '../../components/BottomTabs'
import Footer from '../../components/Footer'
import LearnStepper from '../../components/LearnStepper'
import NextUpCard from '../../components/NextUpCard'
import { getActivity, getDueReviews, getProgress, getScores, getWeaknesses, type DailyActivity, type Scores, type WeaknessRes, type ConceptProgress } from '../../lib/api'
import { getUserInfo, getGuestId } from '../../lib/auth'
import { selectShelfHead } from '../../lib/nextUp'
import { concepts as conceptCatalog } from '../../lib/conceptData'

function LearnEntry() {
  const [activity, setActivity] = useState<DailyActivity[]>([])
  const [progress, setProgress] = useState<Record<string, ConceptProgress>>({})
  const [weaknesses, setWeaknesses] = useState<WeaknessRes | null>(null)
  const [dueReviews, setDueReviews] = useState(0)
  const [scores, setScores] = useState<Scores | null>(null)
  const [diagnosticCompleted, setDiagnosticCompleted] = useState(false)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const info = getUserInfo()
    setDiagnosticCompleted(info?.diagnostic_completed ?? false)
    const sid = info?.student_id || getGuestId() || ''
    Promise.all([
      getActivity().catch(() => [] as DailyActivity[]),
      sid ? getProgress(sid).catch(() => ({} as Record<string, ConceptProgress>)) : Promise.resolve({} as Record<string, ConceptProgress>),
      getWeaknesses().catch(() => ({ by_domain: {} } as WeaknessRes)),
      getDueReviews().catch(() => ({ count: 0 })),
      sid ? getScores(sid).catch(() => null) : Promise.resolve(null),
    ]).then(([a, p, w, r, s]) => {
      setActivity(a)
      setProgress(p)
      setWeaknesses(w)
      setDueReviews(r.count ?? 0)
      setScores(s as Scores | null)
      setLoading(false)
    }).catch(() => setLoading(false))
  }, [])

  const shelf = useMemo(
    () =>
      selectShelfHead({
        dueReviews,
        weaknesses,
        progress,
        activity,
        diagnosticCompleted,
        conceptsMastered: scores?.concepts_mastered ?? 0,
        catalog: conceptCatalog.map(c => ({ id: c.id, label: c.label, prerequisites: c.prerequisites ?? [], avgTimeSeconds: c.mastery_threshold?.avg_time_seconds })),
      }),
    [dueReviews, weaknesses, progress, activity, diagnosticCompleted, scores?.concepts_mastered],
  )

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
