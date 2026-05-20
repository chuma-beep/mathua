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
import { getLessons, type LessonInfo } from '../../lib/api'

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

export default function StudyPage() {
  const [lessonsByDomain, setLessonsByDomain] = useState<Record<string, LessonInfo[]>>({})
  const [selectedLesson, setSelectedLesson] = useState<LessonInfo | null>(null)
  const [selectedDomain, setSelectedDomain] = useState<string | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    getLessons().then(res => {
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
              <div className="bg-mathua-surface border border-mathua-border rounded-none p-8 mt-6">
                <div className="text-mathua-secondary text-xs font-mono mb-4">
                  Concepts: {selectedLesson.concepts.join(', ')}
                </div>
                <KatexContent>{selectedLesson.body}</KatexContent>
              </div>
              <div className="mt-6">
                {selectedLesson.concepts.map((cid) => (
                  <Link
                    key={cid}
                    href={`/concept?id=${encodeURIComponent(cid)}`}
                    className="inline-block mr-2 mb-2 border border-mathua-border text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue px-3 py-1 text-xs font-mono rounded-none transition-colors"
                  >
                    Practice: {cid}
                  </Link>
                ))}
              </div>
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
                    {lessonsByDomain[domain].map((lesson, i) => (
                      <button
                        key={i}
                        onClick={() => setSelectedLesson(lesson)}
                        className="text-left bg-mathua-surface border border-mathua-border rounded-none p-4 hover:border-mathua-blue transition-colors group"
                      >
                        <div className="font-mono text-sm text-mathua-primary group-hover:text-mathua-blue transition-colors">
                          {lesson.title}
                        </div>
                        <div className="flex items-center gap-2 mt-1">
                          <span className="font-mono text-[10px] text-mathua-muted">
                            {lesson.concepts.length} concept{lesson.concepts.length !== 1 ? 's' : ''}
                          </span>
                          <span className="text-[10px] text-mathua-blue opacity-0 group-hover:opacity-100 transition-opacity">
                            View →
                          </span>
                        </div>
                      </button>
                    ))}
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
