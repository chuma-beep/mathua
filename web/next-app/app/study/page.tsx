'use client'

import { useState, useEffect, useMemo, Suspense, type MouseEvent as ReactMouseEvent } from 'react'
import { useSearchParams, useRouter } from 'next/navigation'
import Link from 'next/link'
import Header from '../../components/Header'
import BottomTabs from '../../components/BottomTabs'
import Footer from '../../components/Footer'
import StudySkeleton from '../../components/skeletons/StudySkeleton'
import { getLessons, getLessonBody, type LessonInfo, type LessonsRes } from '../../lib/api'
import { getUserInfo } from '../../lib/auth'
import { conceptLabels, domainOrder, lessonProgress } from './domains'
import { DomainDrillDown, DomainOverview, LessonDetail } from './components'
import { concepts as conceptCatalog } from '../../lib/conceptData'
import { topoRank } from '../../lib/topoRank'
import { planConceptNavigation } from '../../lib/conceptTarget'

let lessonsCache: { key: string; res: LessonsRes } | null = null
function getLessonsCached(studentId?: string): Promise<LessonsRes> {
  const key = studentId ?? ''
  if (lessonsCache && lessonsCache.key === key) return Promise.resolve(lessonsCache.res)
  return getLessons(studentId).then(res => {
    lessonsCache = { key, res }
    return res
  })
}

function StudyContent() {
  const searchParams = useSearchParams()
  const { push, replace } = useRouter()
  const lessonParam = searchParams.get('lesson')
  const domainParam = searchParams.get('domain')
  const conceptParam = searchParams.get('concept')
  // Set by the Reference link in Learn, so the page can offer a way back to the
  // question the learner was answering instead of only onwards into Learn.
  const fromParam = searchParams.get('from') || ''

  const [lessonsByDomain, setLessonsByDomain] = useState<Record<string, LessonInfo[]>>({})
  const [selectedBody, setSelectedBody] = useState<string | null>(null)
  const [bodyError, setBodyError] = useState(false)
  const [bodyRetry, setBodyRetry] = useState(0)
  const [selectedLesson, setSelectedLesson] = useState<LessonInfo | null>(null)
  const [selectedDomain, setSelectedDomain] = useState<string | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const user = getUserInfo()
    const studentId = user?.student_id
    getLessonsCached(studentId).then(res => {
      // PR7: within-domain lessons follow topo (prereq) order, not title
      // order. Copy — never mutate the cached response.
      const rank = topoRank(conceptCatalog)
      const lessonRank = (l: LessonInfo) =>
        Math.min(...l.concepts.map(c => rank.get(c) ?? Number.MAX_SAFE_INTEGER))
      const ordered: Record<string, LessonInfo[]> = {}
      for (const [domain, lessons] of Object.entries(res.lessons)) {
        ordered[domain] = [...lessons].sort(
          (a, b) => lessonRank(a) - lessonRank(b) || (a.title < b.title ? -1 : 1),
        )
      }
      setLessonsByDomain(ordered)
      setLoading(false)
    }).catch((e) => { console.error('getLessons failed:', e); setLoading(false) })
  }, [])

  // O(1) lookups for URL → state sync (replaces find-in-loop).
  const lessonByTitle = useMemo(() => {
    const m = new Map<string, LessonInfo>()
    for (const lessons of Object.values(lessonsByDomain)) {
      for (const l of lessons) m.set(l.title, l)
    }
    return m
  }, [lessonsByDomain])
  const lessonByConcept = useMemo(() => {
    const m = new Map<string, LessonInfo>()
    for (const lessons of Object.values(lessonsByDomain)) {
      for (const l of lessons) {
        for (const c of l.concepts) {
          if (!m.has(c)) m.set(c, l)
        }
      }
    }
    return m
  }, [lessonsByDomain])

  // Concept reveal: the URL is the state. Chips navigate (push across
  // lessons, replace for in-lesson jumps) and bump this request; the
  // LessonDetail effect opens + scrolls to the matching block. Direct loads
  // and /concept redirects arrive via the URL-sync effect below.
  const [revealReq, setRevealReq] = useState<{ cid: string; n: number } | null>(null)
  const bumpReveal = (cid: string) =>
    setRevealReq((r) => ({ cid, n: r?.cid === cid ? r.n + 1 : 1 }))

  const handleConceptSelect = (cid: string, e: ReactMouseEvent<HTMLAnchorElement>): void => {
    // New-tab/middle clicks keep default link behavior.
    if (e.metaKey || e.ctrlKey || e.shiftKey || e.altKey || e.button !== 0) return
    e.preventDefault()
    const nav = planConceptNavigation(cid, lessonByConcept, selectedLesson?.title ?? null, conceptParam)
    if (nav.action === 'reveal') {
      // Re-click: URL unchanged, so the URL-sync effect won't fire.
      bumpReveal(cid)
    } else if (nav.action === 'replace') {
      replace(nav.href, { scroll: false })
    } else {
      push(nav.href)
    }
    // Cross-URL reveal comes from the URL-sync effect below.
  }
  useEffect(() => {
    if (lessonParam) {
      setSelectedLesson(lessonByTitle.get(lessonParam) ?? null)
      return
    }
    if (conceptParam && Object.keys(lessonsByDomain).length > 0) {
      const found = lessonByConcept.get(conceptParam)
      if (found) {
        setSelectedLesson(found)
        // Direct loads, redirects, and chip navigations land here: reveal
        // the concept's block once LessonDetail has its KP shards.
        bumpReveal(conceptParam)
        return
      }
    }

    if (domainParam && lessonsByDomain[domainParam]) {
      setSelectedDomain(domainParam)
    } else {
      setSelectedDomain(null)
    }
  }, [lessonParam, domainParam, conceptParam, lessonsByDomain, lessonByTitle, lessonByConcept])

  // Lazily fetch the selected lesson's markdown body.
  useEffect(() => {
    if (!selectedLesson || selectedLesson.body) {
      setSelectedBody(selectedLesson?.body ?? null)
      return
    }
    let cancelled = false
    setSelectedBody(null)
    setBodyError(false)
    getLessonBody(selectedLesson.title)
      .then(body => { if (!cancelled) setSelectedBody(body) })
      .catch((e) => { console.error('lesson body failed:', e); if (!cancelled) setBodyError(true) })
    return () => { cancelled = true }
  }, [selectedLesson, bodyRetry])

  const hydratedLesson = useMemo(
    () => (selectedLesson ? { ...selectedLesson, body: selectedBody ?? '' } : null),
    [selectedLesson, selectedBody]
  )

  // Popstate: browser back/forward
  useEffect(() => {
    const onPop = () => {
      const params = window.location.search
      if (!params.includes('lesson=') && !params.includes('domain=') && !params.includes('concept=')) {
        setSelectedLesson(null)
        setSelectedDomain(null)
      }
    }
    window.addEventListener('popstate', onPop)
    return () => window.removeEventListener('popstate', onPop)
  }, [])

  const sortedDomains = useMemo(() => {
    return Object.keys(lessonsByDomain).sort((a, b) => {
      const ai = domainOrder.indexOf(a)
      const bi = domainOrder.indexOf(b)
      return (ai === -1 ? 999 : ai) - (bi === -1 ? 999 : bi)
    })
  }, [lessonsByDomain])

  // Aggregate progress per domain. Not rendered anywhere on Study — it survives
  // only to decide whether to offer the "Start here" orientation to a learner
  // who has not mastered anything anywhere. Progression itself lives on Profile.
  const domainAgg = useMemo(() => {
    const agg: Record<string, { mastered: number; total: number; pct: number }> = {}
    for (const [domain, lessons] of Object.entries(lessonsByDomain)) {
      let mastered = 0, total = 0
      for (const lesson of lessons) {
        const p = lessonProgress(lesson)
        mastered += p.mastered
        total += p.total
      }
      agg[domain] = { mastered, total, pct: total > 0 ? Math.round((mastered / total) * 100) : 0 }
    }
    return agg
  }, [lessonsByDomain])

  const allLessons = useMemo(() => {
    const items: { title: string; body: string; domain: string; concepts: string[]; conceptLabels: string[] }[] = []
    for (const [domain, lessons] of Object.entries(lessonsByDomain)) {
      for (const l of lessons) {
        items.push({
          title: l.title,
          body: l.body ?? '',
          domain,
          concepts: l.concepts,
          conceptLabels: l.concepts.map(cid => conceptLabels.get(cid) || cid),
        })
      }
    }
    return items
  }, [lessonsByDomain])

  if (loading) {
    return (
      <>
        <Header />
        <div className="pt-[var(--chrome-top)] lg:pt-0">
          <StudySkeleton />
        </div>
        <BottomTabs />
        <Footer />
      </>
    )
  }

  const backHref = selectedLesson
    ? (selectedDomain ? `/study?domain=${encodeURIComponent(selectedDomain)}` : '/study')
    : selectedDomain
    ? '/study'
    : '/'

  return (
    <>
      <Header />
      <div className="max-w-container mx-auto px-4 sm:px-6 pb-[calc(80px+env(safe-area-inset-bottom))] lg:pb-0 overflow-x-hidden min-w-0">
        <section className="pt-8 min-w-0 overflow-hidden">
          {!selectedLesson && (
            <span className="flex justify-between items-center mb-4">
              <Link href={backHref} className="text-mathua-secondary text-sm hover:text-mathua-primary">
                ← Back
              </Link>
            </span>
          )}
          {selectedLesson ? (
            // ── Lesson Detail ──
            <>
              {bodyError && (
                <div role="alert" className="mb-4 flex flex-wrap items-center gap-3 border border-mathua-red bg-mathua-surface px-4 py-3">
                  <span className="font-mono text-xs text-mathua-red">Couldn&apos;t load the lesson text.</span>
                  <button
                    type="button"
                    onClick={() => setBodyRetry(n => n + 1)}
                    className="font-mono text-xs text-mathua-blue hover:text-mathua-blue-hover uppercase tracking-wider"
                  >
                    Retry →
                  </button>
                </div>
              )}
              <LessonDetail
              lesson={hydratedLesson ?? selectedLesson}
              domain={selectedDomain}
              fromConcept={fromParam}
              revealReq={revealReq}
              onConceptSelect={handleConceptSelect}
              onBack={() => {
                setSelectedLesson(null)
                const url = selectedDomain
                  ? '/study?domain=' + encodeURIComponent(selectedDomain)
                  : '/study'
                push(url)
              }}
            />
            </>
          ) : selectedDomain ? (
            // ── Domain Drill-Down ──
            <DomainDrillDown
              domain={selectedDomain}
              lessons={lessonsByDomain[selectedDomain] || []}
              onBack={() => {
                setSelectedDomain(null)
                push('/study')
              }}
              onSelectLesson={(lesson) => {
                setSelectedLesson(lesson)
                push('/study?domain=' + encodeURIComponent(selectedDomain) + '&lesson=' + encodeURIComponent(lesson.title))
              }}
            />
          ) : (
            // ── Domain Overview ──
            <>
              {Object.values(domainAgg).every(a => a.mastered === 0) && (
                <div className="max-w-4xl mx-auto mt-2 mb-4 border border-mathua-border bg-mathua-surface p-4">
                  <h3 className="font-mono text-[11px] text-mathua-muted uppercase tracking-wider mb-2">Start here</h3>
                  <p className="font-mono text-xs text-mathua-secondary mb-3">New here? Follow the order, it respects prerequisites.</p>
                  <div className="flex flex-wrap gap-2">
                    {[
                      { prefix: 'arith', label: 'Arithmetic' },
                      { prefix: 'frac', label: 'Fractions' },
                      { prefix: 'prealg', label: 'Pre-Algebra' },
                    ].map(d => (
                      <button
                        type="button"
                        key={d.prefix}
                        onClick={() => {
                          const match = sortedDomains.find(s => s === d.prefix || s.startsWith(d.prefix + '.') || s.startsWith(d.prefix))
                          const target = match ?? sortedDomains.find(s => s.startsWith(d.prefix.slice(0, 4))) ?? d.prefix
                          setSelectedDomain(target)
                          push('/study?domain=' + encodeURIComponent(target))
                        }}
                        className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue-faint px-4 py-2 font-mono text-xs min-h-[36px] inline-flex items-center justify-center"
                      >
                        {d.label} →
                      </button>
                    ))}
                  </div>
                  <p className="font-mono text-[10px] text-mathua-muted mt-3">Study is the reference library: browse any lesson, read it in full, nothing is tracked. Learn is where you answer questions, and the scheduler picks what comes next from the prerequisite graph.</p>
                </div>
              )}
              <DomainOverview
                domains={sortedDomains}
                lessonsByDomain={lessonsByDomain}
                allLessons={allLessons}
                onSelectDomain={(d) => {
                  setSelectedDomain(d)
                  push('/study?domain=' + encodeURIComponent(d))
                }}
                onSelectLesson={(lesson) => {
                  setSelectedLesson(lesson)
                  push('/study?lesson=' + encodeURIComponent(lesson.title))
                }}
              />
            </>
          )}
        </section>
      </div>

      <BottomTabs />
      <Footer />
    </>
  )
}

export default function StudyPage() {
  return (
    <Suspense fallback={
      <><Header /><div className="pt-[var(--chrome-top)] lg:pt-0"><StudySkeleton /></div></>
    }>
      <StudyContent />
    </Suspense>
  )
}
