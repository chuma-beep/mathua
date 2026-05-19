'use client'

import { useState, useEffect } from 'react'
import Link from 'next/link'
import Header from '../../components/Header'
import SectionHeader from '../../components/SectionHeader'
import Footer from '../../components/Footer'
import AsciiDivider from '../../components/AsciiDivider'
import KatexContent from '../../components/KatexContent'
import { getLessons, type LessonInfo } from '../../lib/api'

export default function StudyPage() {
  const [lessonsByDomain, setLessonsByDomain] = useState<Record<string, LessonInfo[]>>({})
  const [selectedLesson, setSelectedLesson] = useState<LessonInfo | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    getLessons().then(res => {
      setLessonsByDomain(res.lessons)
      setLoading(false)
    }).catch(() => setLoading(false))
  }, [])

  const domainOrder = ['trig', 'calc', 'linalg', 'complex', 'alg', 'arith', 'abstract', 'discrete', 'nt']
  const domainLabels: Record<string, string> = {
    'trig': 'Trigonometry', 'calc': 'Calculus', 'linalg': 'Linear Algebra',
    'complex': 'Complex Numbers', 'alg': 'Algebra', 'arith': 'Arithmetic',
    'abstract': 'Abstract Algebra', 'discrete': 'Discrete Math', 'nt': 'Number Theory',
    'calc.integral': 'Calculus — Integrals', 'calc.deriv': 'Calculus — Derivatives',
    'calc.limit': 'Calculus — Limits',
    'linalg.vector': 'Linear Algebra — Vectors',
    'linalg.matrix': 'Linear Algebra — Matrices',
    'linalg.eigen': 'Linear Algebra — Eigen',
    'linalg.det': 'Linear Algebra — Determinants',
    'linalg.transformations': 'Linear Algebra — Transformations',
    'linalg.span': 'Linear Algebra — Span',
    'linalg.basis': 'Linear Algebra — Basis',
    'complex.concept': 'Complex Numbers',
    'alg.poly': 'Algebra — Polynomials',
    'alg.factor': 'Algebra — Factoring',
    'alg.quad': 'Algebra — Quadratics',
    'alg.log': 'Algebra — Logarithms',
    'alg.systems': 'Algebra — Linear Systems',
    'arith.exp': 'Arithmetic — Exponents',
    'arith.sqrt': 'Arithmetic — Radicals',
    'abstract.group': 'Abstract Algebra — Groups',
    'abstract.ring': 'Abstract Algebra — Rings',
    'abstract.homomorphism': 'Abstract Algebra — Homomorphisms',
    'discrete.sets': 'Discrete Math — Sets',
    'discrete.combinatorics': 'Discrete Math — Combinatorics',
  }

  const sortedDomains = Object.keys(lessonsByDomain).sort((a, b) => {
    const ai = domainOrder.indexOf(a)
    const bi = domainOrder.indexOf(b)
    return (ai === -1 ? 999 : ai) - (bi === -1 ? 999 : bi)
  })

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
              <p className="text-mathua-secondary text-sm text-center mt-2 mb-8">
              </p>

              {sortedDomains.length === 0 && (
                <p className="text-mathua-muted text-sm text-center">No lessons available.</p>
              )}

              {sortedDomains.map(domain => (
                <div key={domain} className="mb-8">
                  <h3 className="font-serif text-lg text-mathua-primary mb-3 border-b border-mathua-border pb-2">
                    {domainLabels[domain] || domain}
                  </h3>
                  <div className="grid grid-cols-1 sm:grid-cols-2 gap-2">
                    {lessonsByDomain[domain].map((lesson, i) => (
                      <button
                        key={i}
                        onClick={() => setSelectedLesson(lesson)}
                        className="text-left bg-mathua-surface border border-mathua-border rounded-none p-4 hover:border-mathua-blue transition-colors"
                      >
                        <div className="font-mono text-sm text-mathua-primary">{lesson.title}</div>
                        <div className="font-mono text-[10px] text-mathua-muted mt-1">
                          {lesson.concepts.length} concept{lesson.concepts.length !== 1 ? 's' : ''}
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
