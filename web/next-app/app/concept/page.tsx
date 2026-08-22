'use client'

import { useState, useEffect, Suspense } from 'react'
import { useSearchParams } from 'next/navigation'
import Link from 'next/link'
import Header from '../../components/Header'
import SectionHeader from '../../components/SectionHeader'
import Footer from '../../components/Footer'
import AsciiDivider from '../../components/AsciiDivider'
import KatexContent from '../../components/KatexContent'
import LessonQuiz from '../../components/LessonQuiz'
import { getConceptDetail, type ConceptDetailRes } from '../../lib/api'
import Loading from '../../components/Loading'

function slugify(text: string): string {
  return text
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/^-|-$/g, '')
}

function extractToc(body: string): { id: string; label: string; level: number }[] {
  const headings: { id: string; label: string; level: number }[] = []
  const lines = body.split('\n')
  for (const line of lines) {
    const match = line.match(/^(#{2,3})\s+(.+)/)
    if (match) {
      const level = match[1].length
      const label = match[2].trim()
      const id = slugify(label)
      headings.push({ id, label, level })
    }
  }
  return headings
}

function ConceptContent() {
  const searchParams = useSearchParams()
  const conceptId = searchParams.get('id') || ''
  const [detail, setDetail] = useState<ConceptDetailRes | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    if (!conceptId) {
      setError('No concept specified')
      setLoading(false)
      return
    }
    getConceptDetail(conceptId)
      .then(setDetail)
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false))
  }, [conceptId])

  if (loading) {
    return (
      <>
        <Header />
        <div className="max-w-container mx-auto px-6 max-sm:px-4 pt-20 text-center">
          <Loading label="LOADING CONCEPT" />
        </div>
        <Footer />
      </>
    )
  }

  if (error || !detail) {
    return (
      <>
        <Header />
        <div className="max-w-container mx-auto px-6 max-sm:px-4 pt-20 text-center">
          <p className="text-mathua-muted text-sm">{error || 'Concept not found'}</p>
          <Link href="/" className="text-mathua-blue text-sm hover:underline mt-4 inline-block">
            ← Back home
          </Link>
        </div>
        <Footer />
      </>
    )
  }

  const masteryPct = detail.progress?.mastery_pct ?? 0
  const statusColors: Record<string, string> = {
    MASTERED: 'text-green-400',
    PRACTICING: 'text-yellow-400',
    LEARNING: 'text-yellow-600',
  }
  const statusColor = statusColors[detail.progress?.status ?? ''] || 'text-mathua-muted'
  const toc = extractToc(detail.lesson?.body || '')

  return (
    <>
      <Header />
      <div className="max-w-container mx-auto px-6 max-sm:px-4">
        <section className="pt-8">
          <span className="flex items-center gap-2 mb-4 text-xs font-mono">
            <Link href="/" className="text-mathua-secondary hover:text-mathua-primary">
              Home
            </Link>
            <span className="text-mathua-muted">/</span>
            <Link href="/study" className="text-mathua-secondary hover:text-mathua-primary">
              Study
            </Link>
            {detail.concept && (
              <>
                <span className="text-mathua-muted">/</span>
                <span className="text-mathua-muted">{detail.concept.domain}</span>
                <span className="text-mathua-muted">/</span>
                <span className="text-mathua-primary">{detail.concept.label}</span>
              </>
            )}
          </span>

          <div className="max-w-7xl mx-auto mt-8 mb-16">
            <SectionHeader label={detail.concept.domain} title={detail.concept.label} />
            <p className="text-mathua-muted text-xs font-mono text-center -mt-4 mb-8">
              {detail.concept.domain}.{detail.concept.subdomain} &middot; {detail.concept.id}
            </p>

            <div className="bg-mathua-surface border border-mathua-border rounded-none p-6 mb-8">
              <div className="flex items-center justify-between flex-wrap gap-4">
                <div>
                  <span className="text-mathua-muted text-xs font-mono">status</span>
                  <p className={`font-mono text-sm mt-0.5 ${statusColor}`}>
                    {detail.progress?.status || 'unseen'}
                  </p>
                </div>
                {detail.progress && (
                  <>
                    <div className="text-right">
                      <span className="text-mathua-muted text-xs font-mono">mastery</span>
                      <p className="font-mono text-sm text-mathua-primary mt-0.5">
                        {Math.round(masteryPct * 100)}%
                      </p>
                    </div>
                    <div className="text-right">
                      <span className="text-mathua-muted text-xs font-mono">streak</span>
                      <p className="font-mono text-sm text-mathua-primary mt-0.5">
                        {detail.progress.streak} / {detail.progress.required_streak}
                      </p>
                    </div>
                  </>
                )}
                <div>
                  <span className="text-mathua-muted text-xs font-mono">unlocked</span>
                  <p className={`font-mono text-sm mt-0.5 ${detail.unlocked ? 'text-green-400' : 'text-red-400'}`}>
                    {detail.unlocked ? 'yes' : 'no'}
                  </p>
                </div>
              </div>
            </div>

            {detail.prerequisites.length > 0 && (
              <div className="mb-8">
                <h3 className="font-serif text-sm text-mathua-muted mb-3 font-mono">
                  prerequisites ({detail.prerequisites.length})
                </h3>
                <div className="grid grid-cols-1 sm:grid-cols-2 gap-2">
                  {detail.prerequisites.map((p) => {
                    const pColor = statusColors[p.status] || 'text-mathua-muted'
                    return (
                      <Link
                        key={p.id}
                        href={`/concept?id=${encodeURIComponent(p.id)}`}
                        className="bg-mathua-surface border border-mathua-border rounded-none p-3 hover:border-mathua-blue transition-colors block"
                      >
                        <div className="font-mono text-xs text-mathua-primary">{p.label}</div>
                        <div className="flex items-center gap-2 mt-1">
                          <span className={`font-mono text-[10px] ${pColor}`}>{p.status}</span>
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

            {detail.dependents && detail.dependents.length > 0 && (
              <div className="mb-8">
                <h3 className="font-serif text-sm text-mathua-muted mb-3 font-mono">
                  builds toward ({detail.dependents.length})
                </h3>
                <div className="grid grid-cols-1 sm:grid-cols-2 gap-2">
                  {detail.dependents.map((d) => {
                    const dColor = d.status === 'MASTERED' ? 'text-green-400' : d.status === 'PRACTICING' ? 'text-yellow-400' : d.status === 'LEARNING' ? 'text-yellow-600' : 'text-mathua-muted'
                    return (
                      <Link
                        key={d.id}
                        href={`/concept?id=${encodeURIComponent(d.id)}`}
                        className="bg-mathua-surface border border-mathua-border rounded-none p-3 hover:border-mathua-blue transition-colors block"
                      >
                        <div className="font-mono text-xs text-mathua-primary">{d.label}</div>
                        <span className={`font-mono text-[10px] ${dColor}`}>{d.status}</span>
                      </Link>
                    )
                  })}
                </div>
              </div>
            )}

            {detail.prerequisites.length > 0 && (
              <div className="mb-8">
                <p className="font-mono text-xs text-mathua-muted border-b border-mathua-border pb-2 mb-3">
                  Requires: {detail.prerequisites.map((p, i) => (
                    <span key={p.id}>
                      {i > 0 && <span className="mx-1 text-mathua-border">·</span>}
                      <Link href={`/concept?id=${encodeURIComponent(p.id)}`} className="text-mathua-blue hover:text-mathua-blue-hover transition-colors">
                        {p.label}
                      </Link>
                    </span>
                  ))}
                </p>
              </div>
            )}

            {detail.lesson && (
              <div className="mb-8">
                <div className="flex gap-6">
                  {toc.length > 0 && (
                    <aside className="hidden lg:block w-48 shrink-0">
                      <div className="sticky top-24">
                        <h4 className="font-mono text-[10px] uppercase text-mathua-muted mb-3 tracking-wider">
                          In this lesson
                        </h4>
                        <nav className="space-y-1">
                          {toc.map((h) => (
                            <a
                              key={h.id}
                              href={`#${h.id}`}
                              className={`block font-mono text-xs text-mathua-secondary hover:text-mathua-blue transition-colors ${
                                h.level === 3 ? 'pl-3' : ''
                              }`}
                            >
                              {h.label}
                            </a>
                          ))}
                        </nav>
                      </div>
                    </aside>
                  )}
                  <div className="flex-1 min-w-0">
                    <div className="bg-mathua-surface border border-mathua-border rounded-none p-6 md:p-8 lg:p-10">
                      <div className="text-mathua-secondary text-xs font-mono mb-4">
                        {detail.lesson.title}
                      </div>
                      {detail.lesson.concepts && detail.lesson.concepts.length > 0 && (
                        <div className="text-mathua-muted text-[10px] font-mono mb-4">
                          concepts: {detail.lesson.concepts.join(', ')}
                        </div>
                      )}
                      <div id="lesson-body">
                        <KatexContent>{detail.lesson.body}</KatexContent>
                      </div>
                    </div>
                    <LessonQuiz conceptId={conceptId} />
                  </div>
                </div>
              </div>
            )}

            <div className="text-center">
              <Link
                href={`/session?concept=${encodeURIComponent(conceptId)}`}
                className="inline-block border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white rounded-none h-12 px-8 font-medium text-sm leading-[48px]"
              >
                Practice {detail.concept.label}
              </Link>
            </div>
          </div>
        </section>

        <AsciiDivider pattern="wave" />
        <Footer />
      </div>
    </>
  )
}

export default function ConceptPage() {
  return (
    <Suspense fallback={
      <><Header /><div className="max-w-container mx-auto px-6 max-sm:px-4 pt-20 text-center"><Loading label="LOADING CONCEPT" /></div><Footer /></>
    }>
      <ConceptContent />
    </Suspense>
  )
}
