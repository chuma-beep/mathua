'use client'

import { useEffect, useMemo, useState } from 'react'
import Link from 'next/link'
import Header from '../../components/Header'
import BottomTabs from '../../components/BottomTabs'
import Footer from '../../components/Footer'
import ProfileSkeleton from '../../components/skeletons/ProfileSkeleton'
import PositionBlock from '../../components/PositionBlock'
import DomainProgress from '../../components/DomainProgress'
import ActivityHeatmap from '../../components/ActivityHeatmap'
import StrugglesSection from '../../components/StrugglesSection'
import {
  getActivity,
  getEfficacy,
  getProgress,
  getScores,
  getWeaknesses,
} from '../../lib/api'
import type {
  ConceptProgress,
  DailyActivity,
  EfficacyReport,
  Scores,
  WeaknessRes,
} from '../../lib/api'
import { getUserInfo, getGuestId, ensureGuestId, ensureGuestToken } from '../../lib/auth'
import { selectShelfHead, hrefConceptId, recentlyUnlocked, RECENT_UNLOCK_DAYS } from '../../lib/nextUp'
import { countOverall } from '../../lib/progress'
import catalogue from '../../data/concepts.json'

type Concept = {
  id: string
  label: string
  domain: string
  prerequisites?: string[]
}

// The report. It lives on its own page because ADR-001 locks /profile as the
// hub and ADR-023 as a compact one — the heatmap, the struggles, the efficacy
// numbers and the per-domain picture are all report content, and none of it is
// "what should I do next".
//
// Ordered by the questions a learner actually asks: where am I, how is it
// distributed, how am I doing, what have I been doing, what went wrong, what
// opened up.
export default function ProgressPage() {
  const [scores, setScores] = useState<Scores | null>(null)
  const [activity, setActivity] = useState<DailyActivity[]>([])
  const [progress, setProgress] = useState<Record<string, ConceptProgress>>({})
  const [weaknesses, setWeaknesses] = useState<WeaknessRes | null>(null)
  const [efficacy, setEfficacy] = useState<EfficacyReport | null>(null)
  const [loading, setLoading] = useState(true)
  // Read in the effect, never during render: getUserInfo touches
  // localStorage, and these pages are statically prerendered.
  const [diagnosticCompleted, setDiagnosticCompleted] = useState(false)

  useEffect(() => {
    const info = getUserInfo()
    setDiagnosticCompleted(info?.diagnostic_completed ?? false)
    if (!info) {
      ensureGuestId()
      ensureGuestToken()
    }
    const sid = info?.student_id ?? getGuestId() ?? ''
    const load = async () => {
      try {
        const [scoresRes, progressRes, activityRes, weaknessRes] = await Promise.all([
          sid ? getScores(sid).catch(() => null as Scores | null) : Promise.resolve(null as Scores | null),
          sid ? getProgress(sid).catch(() => ({} as Record<string, ConceptProgress>)) : Promise.resolve({} as Record<string, ConceptProgress>),
          getActivity().catch(() => [] as DailyActivity[]),
          getWeaknesses().catch(() => ({ by_domain: {} }) as WeaknessRes),
        ])
        setScores(scoresRes)
        setProgress(progressRes)
        setActivity(activityRes)
        setWeaknesses(weaknessRes)
        getEfficacy().then(setEfficacy).catch(() => {})
      } finally {
        setLoading(false)
      }
    }
    void load()
  }, [])

  // concepts.json already carries label and domain, so one local type serves
  // both the shelf (which wants CatalogEntry) and the counts (which want the
  // structural subset CatalogueConcept describes).
  const concepts = catalogue as Concept[]

  const catalog = useMemo(
    () =>
      concepts.map(c => ({
        id: c.id,
        label: c.label,
        prerequisites: c.prerequisites ?? [],
        avgTimeSeconds: 0,
      })),
    [concepts],
  )

  // The frontier is the concept the scheduler would serve next, which is
  // already computed from the same inputs the report displays — so "where am
  // I" and "what happens next" cannot disagree.
  const head = useMemo(
    () =>
      selectShelfHead({
        dueReviews: 0,
        weaknesses,
        progress,
        activity,
        diagnosticCompleted,
        conceptsMastered: scores?.concepts_mastered ?? 0,
        catalog,
      }),
    [weaknesses, progress, activity, scores, catalog, diagnosticCompleted],
  )
  const frontierCid = hrefConceptId(head.next.href)
  const frontierLabel = frontierCid
    ? (concepts.find(c => c.id === frontierCid)?.label ?? frontierCid)
    : null

  const unlockRows = useMemo(() => {
    const headCid = hrefConceptId(head.next.href)
    const seen = new Set([head.next.href, ...head.alternatives.map(a => a.href)])
    return recentlyUnlocked({
      catalog,
      progress,
      activity,
      recentDays: RECENT_UNLOCK_DAYS,
    })
      .filter(r => {
        const href = `/learn?concept=${encodeURIComponent(r.id)}`
        if (seen.has(href)) return false
        if (r.id === headCid) return false
        seen.add(href)
        return true
      })
      .slice(0, 5)
  }, [head, progress, activity, catalog])

  if (loading) {
    return (
      <>
        <Header />
        <div className="pt-20 lg:pt-0 px-4">
          <ProfileSkeleton />
        </div>
        <Footer />
        <BottomTabs />
      </>
    )
  }

  const counts = countOverall(concepts, progress)

  return (
    <>
      <Header />
      <div className="max-w-container mx-auto px-4 sm:px-6 pb-[calc(80px+env(safe-area-inset-bottom))] lg:pb-0 overflow-x-hidden min-w-0">
        <section className="pt-8 min-w-0">
          <span className="flex mb-4">
            <Link href="/profile" className="text-mathua-secondary text-sm hover:text-mathua-primary">
              ← Profile
            </Link>
          </span>

          <div className="space-y-8 min-w-0">
            <PositionBlock
              catalogue={concepts}
              progress={progress}
              frontierLabel={frontierLabel}
              frontierHref={frontierCid ? `/learn?concept=${encodeURIComponent(frontierCid)}` : null}
            />

            <section id="domains" aria-label="By domain" className="min-w-0 scroll-mt-28">
              <h2 className="font-serif text-[1.05rem] font-normal text-mathua-primary mb-3">
                By domain
              </h2>
              <DomainProgress progress={progress} />
            </section>

            {efficacy && efficacy.concepts_touched > 0 && (
              <section id="how-doing" aria-label="How you are doing" className="min-w-0 scroll-mt-28">
                <h2 className="font-serif text-[1.05rem] font-normal text-mathua-primary mb-3">
                  How you&apos;re doing
                </h2>
                <div className="grid grid-cols-2 sm:grid-cols-4 gap-3 min-w-0">
                  {[
                    { label: 'First-pass', value: `${Math.round(efficacy.first_pass_rate * 100)}%`, hint: 'correct on attempt 1' },
                    { label: 'Second-pass', value: `${Math.round(efficacy.second_pass_rate * 100)}%`, hint: 'correct within 2 tries' },
                    { label: 'Avg attempts', value: efficacy.avg_attempts_per_concept.toFixed(2), hint: 'per concept' },
                    { label: 'Concepts', value: String(efficacy.concepts_touched), hint: `${efficacy.total_attempts} attempts` },
                  ].map(m => (
                    <div key={m.label} className="border border-mathua-border bg-mathua-surface p-3 min-w-0">
                      <div className="font-mono text-[10px] uppercase text-mathua-muted">{m.label}</div>
                      <div className="font-mono text-xl text-mathua-blue mt-1 truncate">{m.value}</div>
                      <div className="font-mono text-[10px] text-mathua-secondary mt-0.5 truncate">{m.hint}</div>
                    </div>
                  ))}
                </div>
              </section>
            )}

            <section id="activity" aria-label="Activity" className="min-w-0 scroll-mt-28">
              <h2 className="font-serif text-[1.05rem] font-normal text-mathua-primary mb-3">
                Activity
              </h2>
              <ActivityHeatmap data={activity} />
            </section>

            <StrugglesSection weaknesses={weaknesses} />

            {unlockRows.length > 0 && (
              <section aria-label="Recently unlocked" className="min-w-0">
                <h2 className="font-serif text-[1.05rem] font-normal text-mathua-primary mb-3">
                  Recently unlocked
                </h2>
                <ul className="border border-mathua-border bg-mathua-surface p-4 space-y-1.5">
                  {unlockRows.map(r => (
                    <li key={r.id}>
                      <Link
                        href={`/learn?concept=${encodeURIComponent(r.id)}`}
                        className="flex items-baseline gap-2 font-mono text-[11px] text-mathua-secondary hover:text-mathua-blue min-w-0"
                      >
                        <span className="shrink-0 uppercase tracking-wider text-mathua-muted">Unlocked</span>
                        <span className="truncate">{r.label}</span>
                        <span className="ml-auto shrink-0 text-mathua-muted">via {r.via}</span>
                      </Link>
                    </li>
                  ))}
                </ul>
              </section>
            )}

            <div className="pt-2 text-center">
              <Link href="/history" className="font-mono text-xs text-mathua-blue hover:text-mathua-blue-hover">
                Every question you&apos;ve answered →
              </Link>
            </div>
          </div>

          {counts.total > 0 && (
            <p className="mt-8 font-mono text-[10px] text-mathua-muted text-center">
              {counts.mastered} mastered · {counts.unlocked} unlocked · {counts.locked} locked of {counts.total}
            </p>
          )}
        </section>
      </div>
      <Footer />
      <BottomTabs />
    </>
  )
}