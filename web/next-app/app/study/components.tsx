'use client'

import { useState, useEffect, useRef } from 'react'
import Link from 'next/link'
import Image from 'next/image'
import SectionHeader from '../../components/SectionHeader'
import KatexContent from '../../components/KatexContent'
import SearchBar from '../../components/SearchBar'
import LessonQuiz from '../../components/LessonQuiz'
import ReportButton from '../../components/ReportButton'
import MasteryBadge from '../../components/MasteryBadge'
import { getLessonKPs, type LessonInfo, type Scores, type LessonKpsRes } from '../../lib/api'
import { conceptLabels, domainIcon, domainLabels, lessonProgress } from './domains'

export function QuizGateBanner({ scores }: { scores: Scores | null }) {
  // Batch 1: prefer backend gate (xp since last completion); fall back to
  // lifetime total so old mocks/e2e without the new fields still gate.
  const xp = scores?.xp_since_quiz ?? scores?.xp_total ?? 0
  const done = scores?.quiz_due ?? xp >= 150
  const goal = 150 // CONTEXT.md Quiz 150 XP gate MA verbatim
  const pct = Math.min((xp / goal) * 100, 100)
  return (
    <div className={`mt-6 border p-4 flex flex-col sm:flex-row items-center justify-between gap-3 ${done ? 'border-mathua-blue bg-mathua-surface' : 'border-mathua-border bg-mathua-surface'}`}>
      <div className="min-w-0 flex-1">
        <div className="flex items-center gap-2">
          <span className={`px-2 py-0.5 font-mono text-[10px] uppercase tracking-wider ${done ? 'bg-mathua-blue text-white' : 'bg-mathua-border text-mathua-muted'}`}>
            {done ? 'Quiz due' : '150 XP gate'}
          </span>
          <span className="font-mono text-xs text-mathua-primary truncate">
            {done ? '150 XP reached — take your mastery check' : `${xp} / ${goal} XP toward next quiz`}
          </span>
        </div>
        <div className="mt-2 h-1 bg-mathua-code overflow-hidden">
          <div className={`h-full transition-all ${done ? 'bg-mathua-blue' : 'bg-mathua-blue/60'}`} style={{ width: `${pct}%` }} />
        </div>
      </div>
      {done ? (
        <Link href="/goals?quiz=1" className="shrink-0 border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white px-6 py-2 font-mono text-xs min-h-[36px] inline-flex items-center justify-center">
          Take Test →
        </Link>
      ) : null}
    </div>
  )
}

// ── Domain Overview ─────────────────────────────────────

export function DomainOverview({
  domains,
  lessonsByDomain,
  domainAgg,
  allLessons,
  onSelectDomain,
  onSelectLesson,
}: {
  domains: string[]
  lessonsByDomain: Record<string, LessonInfo[]>
  domainAgg: Record<string, { mastered: number; total: number; pct: number }>
  allLessons: { title: string; body: string; domain: string; concepts: string[]; conceptLabels: string[] }[]
  onSelectDomain: (d: string) => void
  onSelectLesson: (lesson: LessonInfo) => void
}) {
  const DOMAINS_PER_PAGE = 20
  const [domainPage, setDomainPage] = useState(1)
  const listTopRef = useRef<HTMLDivElement>(null)
  const pageCount = Math.max(1, Math.ceil(domains.length / DOMAINS_PER_PAGE))
  const safePage = Math.min(domainPage, pageCount)
  const visibleDomains = domains.slice((safePage - 1) * DOMAINS_PER_PAGE, safePage * DOMAINS_PER_PAGE)
  const goToPage = (p: number) => {
    setDomainPage(Math.min(Math.max(1, p), pageCount))
    listTopRef.current?.scrollIntoView({ block: 'start' })
  }
  return (
    <div className="max-w-4xl mx-auto mt-8">
      <div ref={listTopRef} className="scroll-mt-20" />
      <SectionHeader label="Study" title="Browse Lessons" />

      <SearchBar
        items={allLessons}
        onSelect={(item) => {
          const lesson = lessonsByDomain[item.domain]?.find(l => l.title === item.title)
          if (lesson) onSelectLesson(lesson)
        }}
      />

      {domains.length === 0 && (
        <p className="text-mathua-muted text-sm text-center">No lessons available.</p>
      )}

      <div className="space-y-2">
        {visibleDomains.map((domain, idx) => {
          const lessons = lessonsByDomain[domain]
          const agg = domainAgg[domain]
          const label = domainLabels[domain] || domain
          const icon = domainIcon(domain)
          const allMastered = agg.total > 0 && agg.mastered === agg.total

          return (
            <button
              type="button"
              key={domain}
              onClick={() => onSelectDomain(domain)}
              className={`w-full text-left transition-all duration-200 group ${
                allMastered ? 'opacity-50 hover:opacity-70' : ''
              } animate-fadeIn`}
              style={{ animationDelay: `${idx * 30}ms` }}
            >
              <div className="flex items-stretch border border-mathua-border hover:shadow-card-hover transition-shadow bg-mathua-surface-elevated">
                <div
                  className={`w-[3px] shrink-0 transition-colors ${
                    allMastered
                      ? 'bg-mathua-green'
                      : agg.total > 0 && agg.mastered > 0
                      ? 'bg-mathua-blue'
                      : 'bg-mathua-border group-hover:bg-mathua-blue'
                  }`}
                />

                <div className="flex-1 min-w-0 p-3 sm:p-4 overflow-hidden">
                  <div className="flex items-center gap-2 sm:gap-3 min-w-0">
                    <span className="font-mono text-lg text-mathua-blue shrink-0 w-6 text-center select-none">
                      {icon}
                    </span>
                    <span className="font-mono text-sm text-mathua-primary group-hover:text-mathua-blue transition-colors min-w-0 flex-1 truncate">
                      {label}
                    </span>
                    <span className="font-mono text-[11px] text-mathua-muted shrink-0 ml-1">
                      {lessons.length} lesson{lessons.length !== 1 ? 's' : ''}
                    </span>
                  </div>

                  {agg.total > 0 && (
                    <div className="mt-2 pl-9">
                      <div className="h-[3px] bg-mathua-code">
                        <div
                          className={`h-full transition-all duration-500 ${
                            allMastered ? 'bg-mathua-green' : 'bg-mathua-blue'
                          }`}
                          style={{ width: `${agg.pct}%` }}
                        />
                      </div>
                      <div className="flex items-center gap-2 mt-1">
                        <span className="font-mono text-[10px] text-mathua-muted">
                          {agg.mastered} / {agg.total} concepts mastered
                        </span>
                        <span className={`font-mono text-[10px] ${
                          allMastered ? 'text-mathua-green' : 'text-mathua-blue'
                        }`}>
                          {agg.pct}%
                        </span>
                      </div>
                    </div>
                  )}
                </div>

                <div className="hidden sm:flex items-center pr-4 shrink-0">
                  <span className="font-mono text-[11px] text-mathua-blue opacity-0 group-hover:opacity-100 transition-opacity">
                    →
                  </span>
                </div>
              </div>
            </button>
          )
        })}
      </div>
      {pageCount > 1 && (
        <div className="flex items-center justify-center gap-3 mt-6">
          <button
            type="button"
            onClick={() => goToPage(safePage - 1)}
            disabled={safePage <= 1}
            aria-label="Previous page"
            className="border border-mathua-border text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue rounded-none h-11 px-5 font-mono text-xs disabled:opacity-40 disabled:cursor-not-allowed"
          >
            ← Prev
          </button>
          <span className="font-mono text-xs text-mathua-muted" aria-live="polite">
            Page {safePage} of {pageCount}
          </span>
          <button
            type="button"
            onClick={() => goToPage(safePage + 1)}
            disabled={safePage >= pageCount}
            aria-label="Next page"
            className="border border-mathua-border text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue rounded-none h-11 px-5 font-mono text-xs disabled:opacity-40 disabled:cursor-not-allowed"
          >
            Next →
          </button>
        </div>
      )}
      <Link href="/leaderboard" className="mt-4 flex items-center justify-between border border-mathua-border bg-mathua-surface p-3 hover:border-mathua-blue transition-colors">
        <span className="font-mono text-xs text-mathua-primary">Leaderboard</span>
        <span className="font-mono text-[11px] text-mathua-blue">See weekly ranking →</span>
      </Link>
    </div>
  )
}

// ── Domain Drill-Down ───────────────────────────────────

export function DomainDrillDown({
  domain,
  lessons,
  agg,
  onBack,
  onSelectLesson,
}: {
  domain: string
  lessons: LessonInfo[]
  agg?: { mastered: number; total: number; pct: number }
  onBack: () => void
  onSelectLesson: (lesson: LessonInfo) => void
}) {
  const label = domainLabels[domain] || domain
  const icon = domainIcon(domain)

  return (
    <div className="max-w-7xl mx-auto mt-8 mb-16">
      <button
        type="button"
        onClick={onBack}
        className="text-mathua-secondary text-xs font-mono hover:text-mathua-blue mb-6"
      >
        ← All domains
      </button>

      <div className="flex items-center gap-3 mb-4">
        <span className="font-mono text-2xl text-mathua-blue select-none">{icon}</span>
        <SectionHeader label="Domain" title={label} className="flex-1" />
      </div>

      {agg && agg.total > 0 && (
        <div className="border border-mathua-border p-4 mb-6 bg-mathua-surface">
          <div className="flex items-center gap-3 mb-2">
            <span className="font-mono text-[11px] text-mathua-muted">Domain progress</span>
            <span className="font-mono text-[11px] text-mathua-blue">
              {agg.mastered} / {agg.total} concepts
            </span>
            <span className={`font-mono text-[11px] ${
              agg.mastered === agg.total ? 'text-mathua-green' : 'text-mathua-blue'
            }`}>
              {agg.pct}%
            </span>
          </div>
          <div className="h-[3px] bg-mathua-code">
            <div
              className={`h-full transition-all duration-500 ${
                agg.mastered === agg.total ? 'bg-mathua-green' : 'bg-mathua-blue'
              }`}
              style={{ width: `${agg.pct}%` }}
            />
          </div>
        </div>
      )}

      <div className="grid grid-cols-1 sm:grid-cols-2 gap-2">
        {lessons.map((lesson, i) => {
          const { mastered, total } = lessonProgress(lesson)
          const pct = total > 0 ? Math.round((mastered / total) * 100) : 0
          const allDone = total > 0 && mastered === total

          return (
            <button
              type="button"
              key={lesson.title}
              onClick={() => onSelectLesson(lesson)}
              className={`text-left transition-all duration-200 group animate-fadeIn ${
                allDone ? 'opacity-50 hover:opacity-70' : ''
              }`}
              style={{ animationDelay: `${i * 30}ms` }}
            >
              <div className="border border-mathua-border hover:shadow-card-hover transition-shadow bg-mathua-surface-elevated p-3 sm:p-4 min-w-0 overflow-hidden">
                <div className="flex items-center justify-between gap-2 min-w-0">
                  <div className="font-mono text-sm text-mathua-primary group-hover:text-mathua-blue transition-colors min-w-0 flex-1 line-clamp-2 break-words">
                    {lesson.title}
                  </div>
                  <span className="hidden sm:inline font-mono text-[11px] text-mathua-blue opacity-0 group-hover:opacity-100 transition-opacity shrink-0 ml-2">
                    View →
                  </span>
                </div>

                <div className="flex items-center gap-2 mt-2">
                  <span className="font-mono text-[10px] text-mathua-muted">
                    {total} concept{total !== 1 ? 's' : ''}
                  </span>
                  {lesson.progress && (
                    <>
                      <div className="flex-1 h-[3px] bg-mathua-code max-w-[120px]">
                        <div
                          className={`h-full transition-all ${
                            allDone ? 'bg-mathua-green' : 'bg-mathua-blue'
                          }`}
                          style={{ width: `${pct}%` }}
                        />
                      </div>
                      <span className="font-mono text-[10px] text-mathua-muted">
                        {mastered}/{total}
                      </span>
                    </>
                  )}
                </div>

                {lesson.progress && (
                  <div className="flex items-center gap-1 mt-2">
                    {lesson.concepts.map((cid) => (
                      <MasteryBadge key={cid} status={lesson.progress?.[cid]?.status} size="sm" />
                    ))}
                  </div>
                )}
              </div>
            </button>
          )
        })}
      </div>
    </div>
  )
}

// ── Lesson Detail ───────────────────────────────────────

export function LessonDetail({
  lesson,
  domain,
  onBack,
}: {
  lesson: LessonInfo
  domain: string | null
  onBack: () => void
}) {
  // KP-aware: fetch knowledge-point shards for the lesson's concepts.
  const [kpMap, setKpMap] = useState<Record<string, LessonKpsRes>>({})
  useEffect(() => {
    if (!lesson.concepts?.length) return
    let cancelled = false
    lesson.concepts.slice(0, 3).forEach(cid => {
      getLessonKPs(cid).then(res => {
        if (!cancelled && res.kps?.length) setKpMap(prev => ({ ...prev, [cid]: res }))
      }).catch(() => {})
    })
    return () => { cancelled = true }
  }, [lesson])

  // Single pass: collect only concepts that actually have KP shards.
  const kpSections: { cid: string; kps: NonNullable<LessonKpsRes['kps']>; diagram: LessonKpsRes['diagram'] }[] = []
  for (const cid of lesson.concepts.slice(0, 3)) {
    const kps = kpMap[cid]?.kps ?? []
    if (kps.length > 0) kpSections.push({ cid, kps, diagram: kpMap[cid]?.diagram })
  }
  const hasKps = kpSections.length > 0

  return (
    <div className="max-w-7xl mx-auto mt-8 mb-16">
      <button
        type="button"
        onClick={onBack}
        className="text-mathua-secondary text-xs font-mono hover:text-mathua-blue mb-6"
      >
        ← {domain ? domainLabels[domain] || domain : 'All domains'}
      </button>
      <SectionHeader label="Lesson" title={lesson.title} />

      <div className="bg-mathua-surface border border-mathua-border p-4 sm:p-6 mt-6 mb-6 min-w-0 overflow-hidden">
        <div className="flex items-center gap-2 sm:gap-3 flex-wrap min-w-0">
          <span className="text-mathua-muted text-xs font-mono shrink-0">Concepts:</span>
          {lesson.concepts.map((cid) => {
            const p = lesson.progress?.[cid]
            return (
              <Link
                key={cid}
                href={`/concept?id=${encodeURIComponent(cid)}`}
                className="inline-flex items-center gap-1.5 border border-mathua-border px-2.5 py-1.5 min-h-[36px] text-xs font-mono text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue transition-colors max-w-full"
              >
                <MasteryBadge status={p?.status} size="sm" />
                <span className="truncate">{conceptLabels.get(cid) || cid}</span>
              </Link>
            )
          })}
        </div>
      </div>

      {lesson.prerequisites && lesson.prerequisites.length > 0 && (
        <div className="mb-8">
          <h3 className="font-mono text-[11px] text-mathua-muted mb-3 border-b border-mathua-border pb-2 uppercase tracking-wider">
            Before you start ({lesson.prerequisites.length})
          </h3>
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-2">
            {lesson.prerequisites.map((p) => (
              <Link
                key={p.id}
                href={`/concept?id=${encodeURIComponent(p.id)}`}
                className="bg-mathua-surface border border-mathua-border p-3 hover:border-mathua-blue transition-colors block"
              >
                <div className="flex items-center gap-2">
                  <MasteryBadge status={p.status} size="sm" />
                  <div className="font-mono text-xs text-mathua-primary">{p.label}</div>
                </div>
                <div className="flex items-center gap-2 mt-1.5">
                  <div className="flex-1 h-1 bg-mathua-bg overflow-hidden">
                    <div
                      className="h-full bg-mathua-blue transition-all"
                      style={{ width: `${Math.round(p.mastery_pct * 100)}%` }}
                    />
                  </div>
                  <span className="font-mono text-[10px] text-mathua-muted">
                    {Math.round(p.mastery_pct * 100)}%
                  </span>
                </div>
              </Link>
            ))}
          </div>
        </div>
      )}

      {hasKps ? (
        <div className="space-y-5">
          {kpSections.map(({ cid, kps, diagram }) => {
            return (
              <div key={cid} className="w-full max-w-full min-w-0 overflow-hidden">
                <div className="flex items-center gap-2 mb-3">
                  <span className="bg-mathua-blue text-white px-2 py-0.5 font-mono text-[10px] uppercase tracking-wider">
                    Worked example
                  </span>
                  <span className="font-mono text-[10px] text-mathua-muted">{cid}</span>
                </div>
                {kps.map((kp, k) => (
                  <div key={`${cid}-${kp.label}`} className="border border-mathua-border bg-mathua-surface p-4 mb-3 w-full max-w-full min-w-0 overflow-hidden">
                    <p className="font-mono text-xs text-mathua-primary">
                      {k + 1}. {kp.label}
                    </p>
                    {kp.subgoals.length > 0 && (
                      <ul className="mt-2 space-y-1">
                        {kp.subgoals.map((sg) => (
                          <li
                            key={sg}
                            className="font-mono text-[11px] text-mathua-secondary pl-3 relative before:content-['–'] before:absolute before:left-0"
                          >
                            {sg}
                          </li>
                        ))}
                      </ul>
                    )}
                    <details className="mt-2">
                      <summary className="font-mono text-[10px] text-mathua-blue uppercase tracking-wider cursor-pointer">
                        Worked example
                      </summary>
                      <div className="mt-2 flex flex-col md:flex-row gap-4 items-start">
                        {diagram && (
                          <div className="shrink-0 bg-mathua-code border border-mathua-border p-2 flex items-center justify-center">
                            <Image src={diagram} alt={`Worked diagram for ${conceptLabels.get(cid) || cid}`} width={220} height={180} className="max-w-full h-auto" style={{ maxHeight: '180px' }} unoptimized />
                          </div>
                        )}
                        <div className="bg-mathua-code border border-mathua-border p-3 text-sm flex-1 min-w-0 overflow-hidden">
                          <KatexContent>{kp.worked_example}</KatexContent>
                        </div>
                      </div>
                    </details>
                    <div className="mt-2 flex justify-end gap-3">
                      <ReportButton
                        conceptId={cid}
                        lessonId={lesson.title}
                        kind="worked_example"
                        question={`${kp.label}: ${kp.worked_example}`}
                      />
                      {diagram && (
                        <ReportButton
                          conceptId={cid}
                          lessonId={lesson.title}
                          kind="diagram"
                          question={diagram}
                        />
                      )}
                    </div>
                  </div>
                ))}
              </div>
            )
          })}
        </div>
      ) : (
        <div className="bg-mathua-surface border border-mathua-border p-4 sm:p-6 md:p-8 lg:p-10 w-full max-w-full min-w-0 overflow-hidden">
          <div className="flex items-center gap-2 mb-4 border-b border-mathua-border pb-3">
            <span className="bg-mathua-blue text-white px-2 py-0.5 font-mono text-[10px] uppercase tracking-wider">
              Worked example
            </span>
            <span className="font-mono text-[10px] text-mathua-muted">
              study this first, then practice below
            </span>
          </div>
          <div className="w-full max-w-full min-w-0 overflow-hidden">
            <KatexContent>{lesson.body}</KatexContent>
          </div>
          <div className="mt-3 flex justify-end">
            <ReportButton
              conceptId={lesson.concepts[0]}
              lessonId={lesson.title}
              kind="lesson_body"
              question={lesson.body?.slice(0, 2000)}
            />
          </div>
        </div>
      )}

      <div className="mt-2 flex items-center gap-2">
        <span className="bg-mathua-border text-mathua-primary px-2 py-0.5 font-mono text-[10px] uppercase tracking-wider">
          Practice
        </span>
        <span className="font-mono text-[10px] text-mathua-muted">
          2 in a row to advance
        </span>
      </div>
      {lesson.concepts.slice(0, 3).map(cid => (
        <LessonQuiz key={cid} conceptId={cid} limit={4} />
      ))}

      {lesson.dependents && lesson.dependents.length > 0 && (
        <div className="mt-6 mb-8">
          <h3 className="font-mono text-[11px] text-mathua-muted mb-3 border-b border-mathua-border pb-2 uppercase tracking-wider">
            What to study next ({lesson.dependents.length})
          </h3>
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-2">
            {lesson.dependents.map((d) => (
              <Link
                key={d.id}
                href={`/concept?id=${encodeURIComponent(d.id)}`}
                className="bg-mathua-surface border border-mathua-border p-3 hover:border-mathua-blue transition-colors block"
              >
                <div className="flex items-center gap-2">
                  <MasteryBadge status={d.status} size="sm" />
                  <div className="font-mono text-xs text-mathua-primary">{d.label}</div>
                </div>
              </Link>
            ))}
          </div>
        </div>
      )}

      <div className="mt-8 text-center">
        <Link
          href={`/session`}
          className="inline-block border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white h-12 px-8 font-medium text-sm leading-[48px]"
        >
          Start practicing these concepts
        </Link>
      </div>
    </div>
  )
}
