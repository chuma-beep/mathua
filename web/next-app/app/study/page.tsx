'use client'

import { useState, useEffect, useMemo } from 'react'
import Link from 'next/link'
import Header from '../../components/Header'
import SectionHeader from '../../components/SectionHeader'
import Footer from '../../components/Footer'
import AsciiDivider from '../../components/AsciiDivider'
import KatexContent from '../../components/KatexContent'
import SearchBar from '../../components/SearchBar'
import LessonQuiz from '../../components/LessonQuiz'
import MasteryBadge from '../../components/MasteryBadge'
import { getLessons, type LessonInfo } from '../../lib/api'
import { getUserInfo } from '../../lib/auth'

const domainOrder = [
  'counting', 'arith', 'fractions', 'prealgebra',
  'alg', 'trig', 'calc', 'calc.integral', 'calc.deriv', 'calc.limit',
  'linalg', 'linalg.vector', 'linalg.matrix', 'linalg.eigen', 'linalg.det',
  'linalg.transformations', 'linalg.span', 'linalg.basis',
  'complex', 'complex.concept',
  'alg.poly', 'alg.factor', 'alg.quad', 'alg.log', 'alg.systems',
  'arith.exp', 'arith.sqrt',
  'abstract', 'abstract.group', 'abstract.ring', 'abstract.homomorphism',
  'discrete', 'discrete.sets', 'discrete.combinatorics', 'nt',
]

const domainLabels: Record<string, string> = {
  'counting': 'Counting',
  'arith': 'Arithmetic',
  'fractions': 'Fractions',
  'prealgebra': 'Pre-Algebra',
  'alg': 'Algebra',
  'trig': 'Trigonometry',
  'calc': 'Calculus',
  'calc.integral': 'Calculus — Integrals',
  'calc.deriv': 'Calculus — Derivatives',
  'calc.limit': 'Calculus — Limits',
  'linalg': 'Linear Algebra',
  'linalg.vector': 'Linear Algebra — Vectors',
  'linalg.matrix': 'Linear Algebra — Matrices',
  'linalg.eigen': 'Linear Algebra — Eigen',
  'linalg.det': 'Linear Algebra — Determinants',
  'linalg.transformations': 'Linear Algebra — Transformations',
  'linalg.span': 'Linear Algebra — Span',
  'linalg.basis': 'Linear Algebra — Basis',
  'complex': 'Complex Numbers',
  'complex.concept': 'Complex Numbers',
  'alg.poly': 'Algebra — Polynomials',
  'alg.factor': 'Algebra — Factoring',
  'alg.quad': 'Algebra — Quadratics',
  'alg.log': 'Algebra — Logarithms',
  'alg.systems': 'Algebra — Linear Systems',
  'arith.exp': 'Arithmetic — Exponents',
  'arith.sqrt': 'Arithmetic — Radicals',
  'abstract': 'Abstract Algebra',
  'abstract.group': 'Abstract Algebra — Groups',
  'abstract.ring': 'Abstract Algebra — Rings',
  'abstract.homomorphism': 'Abstract Algebra — Homomorphisms',
  'discrete': 'Discrete Math',
  'discrete.sets': 'Discrete Math — Sets',
  'discrete.combinatorics': 'Discrete Math — Combinatorics',
  'nt': 'Number Theory',
}

function lessonProgress(lesson: LessonInfo): { mastered: number; total: number } {
  if (!lesson.progress) return { mastered: 0, total: lesson.concepts.length }
  let mastered = 0
  for (const cid of lesson.concepts) {
    const p = lesson.progress[cid]
    if (p && p.status === 'MASTERED') mastered++
  }
  return { mastered, total: lesson.concepts.length }
}

export default function StudyPage() {
  const [lessonsByDomain, setLessonsByDomain] = useState<Record<string, LessonInfo[]>>({})
  const [selectedLesson, setSelectedLesson] = useState<LessonInfo | null>(null)
  const [selectedDomain, setSelectedDomain] = useState<string | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const user = getUserInfo()
    const studentId = user?.student_id
    getLessons(studentId).then(res => {
      setLessonsByDomain(res.lessons)
      setLoading(false)
    }).catch(() => setLoading(false))
  }, [])

  const sortedDomains = useMemo(() => {
    return Object.keys(lessonsByDomain).sort((a, b) => {
      const ai = domainOrder.indexOf(a)
      const bi = domainOrder.indexOf(b)
      return (ai === -1 ? 999 : ai) - (bi === -1 ? 999 : bi)
    })
  }, [lessonsByDomain])

  const allLessons = useMemo(() => {
    const items: { title: string; body: string; domain: string; concepts: string[] }[] = []
    for (const [domain, lessons] of Object.entries(lessonsByDomain)) {
      for (const l of lessons) {
        items.push({ title: l.title, body: l.body, domain, concepts: l.concepts })
      }
    }
    return items
  }, [lessonsByDomain])

  const filteredDomains = useMemo(() => {
    if (!selectedDomain) return sortedDomains
    return sortedDomains.filter(d => d === selectedDomain)
  }, [sortedDomains, selectedDomain])

  if (loading) {
    return (
      <>
        <Header />
        <div className="max-w-container mx-auto px-6 max-sm:px-4 pt-20 text-center">
          <p className="text-mathua-muted text-sm">Loading lessons…</p>
        </div>
        <Footer />
      </>
    )
  }

  return (
    <>
      <Header />
      <div className="max-w-container mx-auto px-6 max-sm:px-4">
        <section className="pt-8">
          <span className="flex justify-between items-center mb-4">
            <Link href="/" className="text-mathua-secondary text-sm hover:text-mathua-primary">
              ← Back
            </Link>
          </span>

          {selectedLesson ? (
            <div className="max-w-4xl mx-auto mt-8 mb-16">
              <button
                onClick={() => setSelectedLesson(null)}
                className="text-mathua-secondary text-xs font-mono hover:text-mathua-blue mb-6"
              >
                ← All domains
              </button>
              <SectionHeader label="Lesson" title={selectedLesson.title} />

              <div className="bg-mathua-surface border border-mathua-border rounded-none p-6 mt-6 mb-6">
                <div className="flex items-center gap-4 flex-wrap">
                  <span className="text-mathua-muted text-xs font-mono">Concepts in this lesson:</span>
                  {selectedLesson.concepts.map((cid) => {
                    const p = selectedLesson.progress?.[cid]
                    return (
                      <Link
                        key={cid}
                        href={`/concept?id=${encodeURIComponent(cid)}`}
                        className="inline-flex items-center gap-1.5 border border-mathua-border px-2.5 py-1 text-xs font-mono text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue transition-colors rounded-none"
                      >
                        <MasteryBadge status={p?.status} size="sm" />
                        <span>{cid}</span>
                      </Link>
                    )
                  })}
                </div>
              </div>

              {selectedLesson.prerequisites && selectedLesson.prerequisites.length > 0 && (
                <div className="mb-8">
                  <h3 className="font-serif text-sm text-mathua-muted mb-3 font-mono border-b border-mathua-border pb-2">
                    Before you start ({selectedLesson.prerequisites.length})
                  </h3>
                  <div className="grid grid-cols-1 sm:grid-cols-2 gap-2">
                    {selectedLesson.prerequisites.map((p) => {
                      const pColor = p.status === 'MASTERED' ? 'text-green-400' : p.status === 'PRACTICING' ? 'text-yellow-400' : p.status === 'LEARNING' ? 'text-yellow-600' : 'text-mathua-muted'
                      return (
                        <Link
                          key={p.id}
                          href={`/concept?id=${encodeURIComponent(p.id)}`}
                          className="bg-mathua-surface border border-mathua-border rounded-none p-3 hover:border-mathua-blue transition-colors block"
                        >
                          <div className="flex items-center gap-2">
                            <MasteryBadge status={p.status} size="sm" />
                            <div className="font-mono text-xs text-mathua-primary">{p.label}</div>
                          </div>
                          <div className="flex items-center gap-2 mt-1.5">
                            <div className="flex-1 h-1 bg-mathua-bg rounded-full overflow-hidden">
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
                      )
                    })}
                  </div>
                </div>
              )}

              <div className="bg-mathua-surface border border-mathua-border rounded-none p-8">
                <KatexContent>{selectedLesson.body}</KatexContent>
              </div>

              {selectedLesson.concepts.slice(0, 3).map(cid => (
                <LessonQuiz key={cid} conceptId={cid} limit={4} />
              ))}

              {selectedLesson.dependents && selectedLesson.dependents.length > 0 && (
                <div className="mt-8 mb-8">
                  <h3 className="font-serif text-sm text-mathua-muted mb-3 font-mono border-b border-mathua-border pb-2">
                    What to study next ({selectedLesson.dependents.length})
                  </h3>
                  <div className="grid grid-cols-1 sm:grid-cols-2 gap-2">
                    {selectedLesson.dependents.map((d) => (
                      <Link
                        key={d.id}
                        href={`/concept?id=${encodeURIComponent(d.id)}`}
                        className="bg-mathua-surface border border-mathua-border rounded-none p-3 hover:border-mathua-blue transition-colors block"
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
                  className="inline-block border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white rounded-none h-12 px-8 font-medium text-sm leading-[48px]"
                >
                  Practice these concepts
                </Link>
              </div>
            </div>
          ) : (
            <div className="max-w-4xl mx-auto mt-8">
              <SectionHeader label="Study" title="Browse Lessons" />
              <p className="text-mathua-secondary text-sm text-center mt-2 mb-6">
              </p>

              <SearchBar
                items={allLessons}
                onSelect={(item) => {
                  const domain = item.domain
                  const lesson = lessonsByDomain[domain]?.find(l => l.title === item.title)
                  if (lesson) setSelectedLesson(lesson)
                }}
              />

              {sortedDomains.length > 1 && (
                <div className="flex flex-wrap gap-2 mb-6 justify-center">
                  <button
                    onClick={() => setSelectedDomain(null)}
                    className={`px-3 py-1 text-xs font-mono border rounded-none transition-colors ${
                      !selectedDomain
                        ? 'bg-mathua-blue text-white border-mathua-blue'
                        : 'border-mathua-border text-mathua-secondary hover:border-mathua-blue'
                    }`}
                  >
                    All
                  </button>
                  {sortedDomains.map(d => (
                    <button
                      key={d}
                      onClick={() => setSelectedDomain(d)}
                      className={`px-3 py-1 text-xs font-mono border rounded-none transition-colors ${
                        selectedDomain === d
                          ? 'bg-mathua-blue text-white border-mathua-blue'
                          : 'border-mathua-border text-mathua-secondary hover:border-mathua-blue'
                      }`}
                    >
                      {domainLabels[d] || d}
                    </button>
                  ))}
                </div>
              )}

              {filteredDomains.length === 0 && (
                <p className="text-mathua-muted text-sm text-center">No lessons available.</p>
              )}

              {filteredDomains.map(domain => (
                <div key={domain} className="mb-8">
                  <h3 className="font-serif text-lg text-mathua-primary mb-3 border-b border-mathua-border pb-2 flex items-center justify-between">
                    <span>{domainLabels[domain] || domain}</span>
                    <span className="text-xs font-mono text-mathua-muted font-normal">
                      {lessonsByDomain[domain].length} lesson{lessonsByDomain[domain].length !== 1 ? 's' : ''}
                    </span>
                  </h3>
                  <div className="grid grid-cols-1 sm:grid-cols-2 gap-2">
                    {lessonsByDomain[domain].map((lesson, i) => {
                      const { mastered, total } = lessonProgress(lesson)
                      const pct = total > 0 ? Math.round((mastered / total) * 100) : 0
                      return (
                        <button
                          key={i}
                          onClick={() => setSelectedLesson(lesson)}
                          className="text-left bg-mathua-surface border border-mathua-border rounded-none p-4 hover:border-mathua-blue transition-colors group"
                        >
                          <div className="font-mono text-sm text-mathua-primary group-hover:text-mathua-blue transition-colors">
                            {lesson.title}
                          </div>
                          <div className="flex items-center gap-2 mt-1.5">
                            <span className="font-mono text-[10px] text-mathua-muted">
                              {total} concept{total !== 1 ? 's' : ''}
                            </span>
                            {lesson.progress && (
                              <>
                                <div className="flex-1 max-w-[100px] h-1.5 bg-mathua-code rounded-full overflow-hidden">
                                  <div
                                    className="h-full bg-mathua-blue rounded-full transition-all"
                                    style={{ width: `${pct}%` }}
                                  />
                                </div>
                                <span className="font-mono text-[10px] text-mathua-muted">
                                  {mastered}/{total}
                                </span>
                              </>
                            )}
                            <span className="text-[10px] text-mathua-blue opacity-0 group-hover:opacity-100 transition-opacity ml-auto">
                              View →
                            </span>
                          </div>
                          {lesson.progress && (
                            <div className="flex items-center gap-1 mt-1.5">
                              {lesson.concepts.map((cid) => (
                                <MasteryBadge key={cid} status={lesson.progress?.[cid]?.status} size="sm" />
                              ))}
                            </div>
                          )}
                        </button>
                      )
                    })}
                  </div>
                </div>
              ))}
            </div>
          )}
        </section>

        <AsciiDivider pattern="wave" />
        <Footer />
      </div>
    </>
  )
}
