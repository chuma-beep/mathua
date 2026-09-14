'use client'

import { useState, useEffect, useCallback } from 'react'
import Link from 'next/link'
import Header from '../../components/Header'
import BottomTabs from '../../components/BottomTabs'
import SectionHeader from '../../components/SectionHeader'
import Footer from '../../components/Footer'
import Loading from '../../components/Loading'
import ProgressCardList from '../../components/ProgressCardList'
import { getAttempts, type AttemptRecord } from '../../lib/api'

const PAGE_SIZE = 50

const SOURCES = ['all', 'diagnostic', 'quiz', 'practice', 'review'] as const

export default function ProgressCardPage() {
  const [incorrectOnly, setIncorrectOnly] = useState(true)
  const [source, setSource] = useState<(typeof SOURCES)[number]>('all')
  const [attempts, setAttempts] = useState<AttemptRecord[]>([])
  const [total, setTotal] = useState(0)
  const [loading, setLoading] = useState(true)
  const [loadingMore, setLoadingMore] = useState(false)
  const [error, setError] = useState('')

  const load = useCallback(async (offset: number, append: boolean) => {
    if (append) setLoadingMore(true)
    else setLoading(true)
    setError('')
    try {
      const res = await getAttempts({
        source: source === 'all' ? undefined : source,
        incorrect_only: incorrectOnly || undefined,
        limit: PAGE_SIZE,
        offset: offset || undefined,
      })
      setTotal(res.total)
      setAttempts(prev => (append ? [...prev, ...res.attempts] : res.attempts))
    } catch {
      setError('Couldn\u2019t load your history. Sign in to see answers tied to your account.')
    } finally {
      setLoading(false)
      setLoadingMore(false)
    }
  }, [source, incorrectOnly])

  useEffect(() => {
    void load(0, false)
  }, [load])

  return (
    <>
      <Header />
      <div className="max-w-container mx-auto px-4 sm:px-6 pb-[calc(80px+env(safe-area-inset-bottom))] lg:pb-0 overflow-x-hidden min-w-0">
        <section className="pt-8 min-w-0 overflow-hidden">
          <span className="flex mb-4">
            <Link href="/profile" className="text-mathua-secondary text-sm hover:text-mathua-primary">
              ← Back
            </Link>
          </span>
          <SectionHeader label="Progress card" title="Answered questions" />
          <p className="text-mathua-secondary text-sm text-center max-w-[600px] mx-auto mt-2 mb-6 px-2">
            The questions are listed on each row, and you can see your answers and compare them with the correct answer.
          </p>

          <div className="max-w-2xl mx-auto min-w-0">
            <div className="flex gap-2 flex-wrap mb-2">
              {(['incorrect', 'all'] as const).map(m => {
                const active = (m === 'incorrect') === incorrectOnly
                return (
                  <button
                    key={m}
                    type="button"
                    onClick={() => setIncorrectOnly(m === 'incorrect')}
                    className={`px-3 h-8 font-mono text-[11px] uppercase tracking-wider border transition-colors ${
                      active
                        ? 'border-mathua-blue text-mathua-blue'
                        : 'border-mathua-border text-mathua-muted hover:text-mathua-primary'
                    }`}
                  >
                    {m === 'incorrect' ? 'Missed' : 'All'}
                  </button>
                )
              })}
            </div>
            <div className="flex gap-2 flex-wrap mb-6 items-center">
              <span className="font-mono text-[10px] text-mathua-muted uppercase tracking-wider">Source:</span>
              {SOURCES.map(s => (
                <button
                  key={s}
                  type="button"
                  onClick={() => setSource(s)}
                  className={`px-3 h-8 font-mono text-[11px] uppercase tracking-wider border transition-colors ${
                    source === s
                      ? 'border-mathua-blue text-mathua-blue'
                      : 'border-mathua-border text-mathua-muted hover:text-mathua-primary'
                  }`}
                >
                  {s}
                </button>
              ))}
            </div>

            {loading ? (
              <div className="text-center py-8">
                <Loading label="LOADING PROGRESS" />
              </div>
            ) : error ? (
              <div className="text-center py-8">
                <p className="text-mathua-muted text-sm">{error}</p>
                <Link href="/login" className="text-mathua-blue text-sm hover:underline mt-4 inline-block">
                  Sign in →
                </Link>
              </div>
            ) : (
              <>
                <p className="font-mono text-[11px] text-mathua-muted mb-3">
                  Showing {attempts.length} of {total}
                </p>
                <ProgressCardList attempts={attempts} />
                {attempts.length < total && (
                  <div className="text-center mt-6">
                    <button
                      type="button"
                      onClick={() => { void load(attempts.length, true) }}
                      disabled={loadingMore}
                      className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white px-6 h-11 font-mono text-xs disabled:opacity-50"
                    >
                      {loadingMore ? 'Loading…' : 'Load more'}
                    </button>
                  </div>
                )}
              </>
            )}
          </div>
        </section>
      </div>
      <Footer />
      <BottomTabs />
    </>
  )
}
