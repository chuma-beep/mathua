'use client'

import { useState, useEffect, useRef, type MouseEvent as ReactMouseEvent } from 'react'
import Link from 'next/link'
import LessonDiagram from '../../components/LessonDiagram'
import SectionHeader from '../../components/SectionHeader'
import KatexContent from '../../components/KatexContent'
import SearchBar from '../../components/SearchBar'
import ReportMenu from '../../components/ReportMenu'
import MasteryBadge from '../../components/MasteryBadge'
import { getLessonKPs, type LessonInfo, type LessonKpsRes } from '../../lib/api'
import { stripMathDelimiters } from '../../lib/lessonMath'
import { conceptLabels, domainIcon, domainLabels, lessonProgress } from './domains'

// Normalize a KP section the way Go normSectionKey does ( mirrored by
// norm_section in scripts/audit_lessons.py): strip math delimiters,
// collapse whitespace. Shared-lesson concepts often point at the same
// section under byte-identical names; normalization keeps those deduped.
function normSection(s: string): string {
  return s.replace(/\\\\[()[\]]|\\[()[\]]|\$\$?/g, '').split(/\s+/).join(' ').trim()
}

// Short stable hash for body text (djb2). The dedupe key is
// (normalized section + body hash): same name and same body is by
// definition a true duplicate, so unlike either half alone this key
// cannot hide distinct content while still catching renames.
function hashBody(s: string): string {
  let h = 5381
  for (let i = 0; i < s.length; i++) h = ((h << 5) + h + s.charCodeAt(i)) >>> 0
  return h.toString(36)
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
  revealReq,
  onConceptSelect,
}: {
  lesson: LessonInfo
  domain: string | null
  onBack: () => void
  revealReq?: { cid: string; n: number } | null
  onConceptSelect?: (cid: string, e: ReactMouseEvent<HTMLAnchorElement>) => void
}) {
  // KP-aware: fetch knowledge-point shards for the lesson's concepts.
  const [kpMap, setKpMap] = useState<Record<string, LessonKpsRes>>({})
  // Concepts whose shard fetch settled (success or failure): the reveal
  // effect waits for these so it never falls back while blocks are still
  // loading. Only the first 3 concepts can ever produce blocks.
  const [settled, setSettled] = useState<Record<string, boolean>>({})
  useEffect(() => {
    if (!lesson.concepts?.length) return
    let cancelled = false
    lesson.concepts.slice(0, 3).forEach(cid => {
      getLessonKPs(cid).then(res => {
        if (cancelled) return
        if (res.kps?.length) setKpMap(prev => ({ ...prev, [cid]: res }))
        setSettled(prev => ({ ...prev, [cid]: true }))
      }).catch(() => {
        if (!cancelled) setSettled(prev => ({ ...prev, [cid]: true }))
      })
    })
    return () => { cancelled = true }
  }, [lesson])

  // Concept reveal state: controlled <details> openness (seeded with the
  // first block, as before), a short-lived highlight on the revealed block
  // or header, and a guard so each reveal request runs once.
  const [openKeys, setOpenKeys] = useState<Set<string> | null>(null)
  const [highlight, setHighlight] = useState<string | null>(null)
  const revealedRef = useRef<string | null>(null)
  const highlightTimer = useRef<ReturnType<typeof setTimeout> | null>(null)
  const topRef = useRef<HTMLDivElement>(null)
  useEffect(() => {
    setOpenKeys(null)
    setSettled({})
    setHighlight(null)
    revealedRef.current = null
    if (highlightTimer.current) clearTimeout(highlightTimer.current)
  }, [lesson.title])

  // Collapse pass: multi-concept lessons share sources, so identical
  // worked bodies (and diagram assets) recur. Each unique body renders
  // once, with combined labels — repeats are dropped silently, never
  // pointed at. Empty bodies never dedupe.
  interface ExampleBlock {
    key: string
    label: string
    also: string[]
    subgoals: string[]
    body: string
    cid: string
    k: number
    diagram: string | null
  }
  const blocks: ExampleBlock[] = []
  const seenBodies = new Map<string, ExampleBlock>()
  const seenDiagrams = new Set<string>()
  for (const cid of lesson.concepts.slice(0, 3)) {
    const kps = kpMap[cid]?.kps ?? []
    if (kps.length === 0) continue
    const diagram = kpMap[cid]?.diagram ?? null
    kps.forEach((kp, k) => {
      const body = kp.worked_example ?? ''
      const anchor = `lesson-kp-${cid}-${k}`
      if (body) {
        const key = `${normSection(kp.section ?? '')}::${hashBody(body)}`
        const first = seenBodies.get(key)
        if (first) {
          first.also.push(kp.label)
          return
        }
      }
      const block: ExampleBlock = {
        key: anchor,
        label: kp.label,
        also: [],
        subgoals: kp.subgoals ?? [],
        body,
        cid,
        k,
        diagram: null,
      }
      if (body) seenBodies.set(`${normSection(kp.section ?? '')}::${hashBody(body)}`, block)
      if (diagram && !seenDiagrams.has(diagram)) {
        seenDiagrams.add(diagram)
        block.diagram = diagram
      }
      blocks.push(block)
    })
  }

  // Seed: first block open, matching the previous uncontrolled default.
  useEffect(() => {
    if (openKeys === null && blocks.length > 0) {
      setOpenKeys(new Set([blocks[0].key]))
    }
  }, [blocks, openKeys])

  const flashHighlight = (key: string): void => {
    setHighlight(key)
    if (highlightTimer.current) clearTimeout(highlightTimer.current)
    highlightTimer.current = setTimeout(() => setHighlight(null), 1800)
  }

  const scrollToKey = (key: string | null): void => {
    const reduce = typeof window !== 'undefined' && window.matchMedia?.('(prefers-reduced-motion: reduce)').matches
    const target = key ? document.getElementById(key) : topRef.current
    target?.scrollIntoView({ behavior: reduce ? 'auto' : 'smooth', block: 'start' })
  }

  // rAF is unavailable in some test DOMs; setTimeout keeps the same
  // paint-then-scroll ordering there.
  const afterPaint = (fn: () => void): void => {
    if (typeof requestAnimationFrame !== 'undefined') requestAnimationFrame(fn)
    else setTimeout(fn, 0)
  }

  // One effect keyed on the reveal request (URL state): resolve is already
  // done by the page — here open the concept's blocks and scroll to the
  // first. Waits for the concept's shard fetch so a slow load never falls
  // back spuriously. Concepts beyond the first 3 can never have blocks.
  useEffect(() => {
    if (!revealReq) return
    const key = `${revealReq.cid}:${revealReq.n}`
    if (revealedRef.current === key) return
    const inScope = lesson.concepts.slice(0, 3).includes(revealReq.cid)
    if (inScope && !settled[revealReq.cid]) return
    revealedRef.current = key
    const matches = blocks.filter((b) => b.cid === revealReq.cid)
    if (matches.length > 0) {
      const keys = matches.map((m) => m.key)
      setOpenKeys((prev) => new Set([...(prev ?? []), ...keys]))
      // The details element exists regardless of open state; paint first.
      afterPaint(() => {
        scrollToKey(matches[0].key)
        flashHighlight(matches[0].key)
      })
    } else {
      // No block for this concept: lesson top + header highlight so the
      // chip never silently does nothing.
      afterPaint(() => {
        scrollToKey(null)
        flashHighlight('header')
      })
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  })

  return (
    <div className="max-w-7xl mx-auto mt-8 mb-16">
      <button
        type="button"
        onClick={onBack}
        className="text-mathua-secondary text-xs font-mono hover:text-mathua-blue mb-6"
      >
        ← {domain ? domainLabels[domain] || domain : 'All domains'}
      </button>
      <div ref={topRef} className={highlight === 'header' ? 'ring-1 ring-mathua-blue rounded-none transition-shadow' : 'transition-shadow'}>
      <SectionHeader label="Lesson" title={lesson.title} />
      </div>

      <div className="bg-mathua-surface border border-mathua-border p-4 sm:p-6 mt-6 mb-6 min-w-0 overflow-hidden">
        <div className="flex items-center gap-2 sm:gap-3 flex-wrap min-w-0">
          <span className="text-mathua-muted text-xs font-mono shrink-0">Concepts:</span>
          {lesson.concepts.map((cid) => {
            const p = lesson.progress?.[cid]
            return (
              <Link
                key={cid}
                href={`/study?concept=${encodeURIComponent(cid)}`}
                onClick={(e) => onConceptSelect?.(cid, e)}
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
                href={`/study?concept=${encodeURIComponent(p.id)}`}
                className="bg-mathua-surface border border-mathua-border p-3 hover:border-mathua-blue transition-colors flex items-center gap-2"
              >
                <MasteryBadge status={p.status} size="sm" />
                <div className="font-mono text-xs text-mathua-primary truncate">{p.label}</div>
                <div className="ml-auto font-mono text-[10px] text-mathua-muted shrink-0">
                  {p.status === 'MASTERED' ? 'mastered' : p.status === 'UNSEEN' || !p.status ? 'not started' : 'in progress'}
                </div>
              </Link>
            ))}
          </div>
        </div>
      )}

      {blocks.length > 0 ? (
        <div className="mb-8">
          <h3 className="font-mono text-[11px] text-mathua-muted mb-3 border-b border-mathua-border pb-2 uppercase tracking-wider">
            Worked examples ({blocks.length})
          </h3>
          <div className="space-y-3">
            {blocks.map((b, i) => (
              <div key={b.key} id={b.key} className={`border border-mathua-border bg-mathua-surface p-4 w-full max-w-full min-w-0 overflow-hidden scroll-mt-20 transition-shadow ${highlight === b.key ? 'ring-1 ring-mathua-blue' : ''}`}>
                <div className="flex items-start gap-2 min-w-0">
                  <div className="flex-1 min-w-0">
                    {b.body ? (
                      <details
                        open={openKeys?.has(b.key) ?? i === 0}
                        onToggle={(e) => {
                          const open = e.currentTarget.open
                          setOpenKeys((prev) => {
                            const next = new Set(prev ?? (blocks[0] ? [blocks[0].key] : []))
                            if (open) next.add(b.key)
                            else next.delete(b.key)
                            return next
                          })
                        }}
                      >
                        <summary className="font-mono text-xs text-mathua-primary cursor-pointer">
                          {i + 1}. {stripMathDelimiters(b.label)}
                          {b.also.length > 0 && (
                            <span className="text-mathua-muted"> · also {b.also.map(stripMathDelimiters).join(', ')}</span>
                          )}
                        </summary>
                        {b.diagram && (
                          <div className="mt-3">
                            <LessonDiagram src={b.diagram} alt={`Worked diagram for ${conceptLabels.get(b.cid) || b.cid}`} />
                          </div>
                        )}
                        {b.subgoals.length > 0 && (
                          <ul className="mt-2 space-y-1">
                            {b.subgoals.map((sg) => (
                              <li
                                key={sg}
                                className="font-mono text-[11px] text-mathua-secondary pl-3 relative before:content-['–'] before:absolute before:left-0"
                              >
                                {stripMathDelimiters(sg)}
                              </li>
                            ))}
                          </ul>
                        )}
                        <div className="mt-2 bg-mathua-code border border-mathua-border p-3 text-sm overflow-hidden">
                          <KatexContent>{b.body}</KatexContent>
                        </div>
                      </details>
                    ) : (
                      <>
                        <p className="font-mono text-xs text-mathua-primary">
                          {i + 1}. {stripMathDelimiters(b.label)}
                          {b.also.length > 0 && (
                            <span className="text-mathua-muted"> · also {b.also.map(stripMathDelimiters).join(', ')}</span>
                          )}
                        </p>
                        {b.subgoals.length > 0 && (
                          <ul className="mt-2 space-y-1">
                            {b.subgoals.map((sg) => (
                              <li
                                key={sg}
                                className="font-mono text-[11px] text-mathua-secondary pl-3 relative before:content-['–'] before:absolute before:left-0"
                              >
                                {stripMathDelimiters(sg)}
                              </li>
                            ))}
                          </ul>
                        )}
                      </>
                    )}
                  </div>
                  <ReportMenu
                    conceptId={b.cid}
                    lessonId={lesson.title}
                    kind={b.diagram ? 'diagram' : 'worked_example'}
                    question={b.diagram ?? `${b.label}: ${b.body}`.slice(0, 2000)}
                    blockId={`${b.cid}/${b.k}`}
                  />
                </div>
              </div>
            ))}
          </div>
        </div>
      ) : (
        <div className="bg-mathua-surface border border-mathua-border p-4 sm:p-6 md:p-8 lg:p-10 w-full max-w-full min-w-0 overflow-hidden">
          <div className="flex items-center gap-2 mb-4 border-b border-mathua-border pb-3">
            <span className="bg-mathua-blue text-white px-2 py-0.5 font-mono text-[10px] uppercase tracking-wider">
              Worked example
            </span>
            <span className="font-mono text-[10px] text-mathua-muted">
              study this first, then start learning below
            </span>
          </div>
          <div className="w-full max-w-full min-w-0 overflow-hidden">
            <KatexContent>{lesson.body}</KatexContent>
          </div>
          <div className="mt-3 flex justify-end">
            <ReportMenu
              conceptId={lesson.concepts[0]}
              lessonId={lesson.title}
              kind="lesson_body"
              question={lesson.body?.slice(0, 2000)}
              blockId={`${lesson.concepts[0] ?? 'lesson'}/body`}
            />
          </div>
        </div>
      )}

      <div id="practice" className="mt-2 flex items-center gap-2 scroll-mt-24">
        <span className="bg-mathua-border text-mathua-primary px-2 py-0.5 font-mono text-[10px] uppercase tracking-wider">
          Learn
        </span>
        <span className="font-mono text-[10px] text-mathua-muted">
          interactive loop, one concept at a time
        </span>
      </div>
      {lesson.concepts.slice(0, 3).map(cid => (
        <Link
          key={cid}
          href={`/learn?concept=${encodeURIComponent(cid)}`}
          className="mt-3 flex items-center justify-between gap-3 border border-mathua-border bg-mathua-surface p-4 hover:border-mathua-blue transition-colors"
        >
          <span className="font-mono text-xs text-mathua-primary truncate">
            {conceptLabels.get(cid) || cid}
          </span>
          <span className="shrink-0 font-mono text-xs text-mathua-blue">
            Start learning →
          </span>
        </Link>
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
                href={`/study?concept=${encodeURIComponent(d.id)}`}
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

    </div>
  )
}
