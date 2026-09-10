'use client'

import { useState, useEffect, useMemo, Suspense } from 'react'
import { useSearchParams, useRouter } from 'next/navigation'
import Link from 'next/link'
import Header from '../../components/Header'
import BottomTabs from '../../components/BottomTabs'
import Footer from '../../components/Footer'
import AsciiDivider from '../../components/AsciiDivider'
import Loading from '../../components/Loading'
import { getLessons, getLessonBody, getScores, type LessonInfo, type LessonsRes, type Scores } from '../../lib/api'
import { getUserInfo, getGuestId } from '../../lib/auth'
import { conceptLabels, domainOrder, lessonProgress } from './domains'
import { DomainDrillDown, DomainOverview, LessonDetail, QuizGateBanner } from './components'

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
  const { push } = useRouter()
  const lessonParam = searchParams.get('lesson')
  const domainParam = searchParams.get('domain')
  const conceptParam = searchParams.get('concept')

  const [lessonsByDomain, setLessonsByDomain] = useState<Record<string, LessonInfo[]>>({})
  const [selectedBody, setSelectedBody] = useState<string | null>(null)
  const [bodyError, setBodyError] = useState(false)
  const [bodyRetry, setBodyRetry] = useState(0)
  const [selectedLesson, setSelectedLesson] = useState<LessonInfo | null>(null)
  const [selectedDomain, setSelectedDomain] = useState<string | null>(null)
  const [loading, setLoading] = useState(true)
  const [scores, setScores] = useState<Scores | null>(null)

  useEffect(() => {
    const user = getUserInfo()
    const studentId = user?.student_id
    getLessonsCached(studentId).then(res => {
      setLessonsByDomain(res.lessons)
      setLoading(false)
    }).catch((e) => { console.error('getLessons failed:', e); setLoading(false) })
    const sid = studentId || getGuestId() || ''
    if (sid) {
      getScores(sid).then(setScores).catch(() => {})
    }
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

  // URL → state sync
  useEffect(() => {
    if (lessonParam) {
      setSelectedLesson(lessonByTitle.get(lessonParam) ?? null)
      return
    }
    if (conceptParam && Object.keys(lessonsByDomain).length > 0) {
      const found = lessonByConcept.get(conceptParam)
      if (found) { setSelectedLesson(found); return }
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

  // Aggregate progress per domain
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
        <div className="max-w-container mx-auto px-4 sm:px-6 pt-20 pb-[calc(80px+env(safe-area-inset-bottom))] lg:pb-0 text-center overflow-x-hidden min-w-0">
          <Loading label="LOADING LESSONS" />
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

          <QuizGateBanner scores={scores} />

          {selectedLesson ? (
            // ── Lesson Detail ──
            <>
              {bodyError && (
                <div role="alert" className="mb-4 flex flex-wrap items-center gap-3 border border-mathua-red bg-mathua-surface px-4 py-3">
                  <span className="font-mono text-xs text-mathua-red">Couldn&apos;t load the lesson text, practice below still works.</span>
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
              agg={domainAgg[selectedDomain]}
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
                        className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white px-4 py-2 font-mono text-xs min-h-[36px] inline-flex items-center justify-center"
                      >
                        {d.label} →
                      </button>
                    ))}
                  </div>
                  <p className="font-mono text-[10px] text-mathua-muted mt-3">Study is a library (Lesson = corpus): browse any lesson, but practice respects the DAG. Each lesson shows 2 in a row to advance before you practice.</p>
                </div>
              )}
              <DomainOverview
                domains={sortedDomains}
                lessonsByDomain={lessonsByDomain}
                domainAgg={domainAgg}
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

        <AsciiDivider pattern="wave" />
        <Footer />
      </div>
      <BottomTabs />
    </>
  )
}

export default function StudyPage() {
  return (
    <Suspense fallback={
      <><Header /><div className="max-w-container mx-auto px-4 sm:px-6 pt-20 pb-[calc(80px+env(safe-area-inset-bottom))] lg:pb-0 text-center overflow-x-hidden min-w-0"><Loading label="LOADING LESSONS" /></div><Footer /></>
    }>
      <StudyContent />
    </Suspense>
  )
}
