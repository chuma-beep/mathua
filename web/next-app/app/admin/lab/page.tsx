'use client'

import AdminShell from '@/components/admin/AdminShell'
import { useCallback, useEffect, useState } from 'react'
import Link from 'next/link'
import SectionHeader from '../../../components/SectionHeader'
import Loading from '../../../components/Loading'
import LabQuestionCard from '../../../components/LabQuestionCard'
import {
  generateLabQuestions, getErrorMessage, listAdminContent,
  type AdminContentConcept, type LabGenerateRes,
} from '../../../lib/api'

// The reproduction bench.
//
// A report says "this question was broken". The queue can replay that one
// report; this page answers the question behind it for any concept, including
// the ones nobody has reported yet: what does the generator produce, what
// format is the field given, and does a correct answer grade as correct.
//
// Both endpoints behind it are read-only, and the page says so where it matters.
// Generating and grading write no attempt, XP, streak or mastery, so an operator
// can try every question in the corpus without touching a learner's record.

const DIFFICULTIES = [0.3, 0.5, 0.8, 1] as const
const COUNTS = [1, 3, 5, 10] as const

export default function AdminLabPage() {
  const [view, setView] = useState<'generate' | 'formats'>('generate')

  const [q, setQ] = useState('')
  const [results, setResults] = useState<AdminContentConcept[] | null>(null)
  const [total, setTotal] = useState(0)
  const [searchError, setSearchError] = useState('')

  const [selected, setSelected] = useState<AdminContentConcept | null>(null)
  const [difficulty, setDifficulty] = useState<number>(0.5)
  const [count, setCount] = useState<number>(3)
  const [generated, setGenerated] = useState<LabGenerateRes | null>(null)
  const [genBusy, setGenBusy] = useState(false)
  const [genError, setGenError] = useState('')

  const [overrides, setOverrides] = useState<AdminContentConcept[] | null>(null)
  const [overrideError, setOverrideError] = useState('')

  useEffect(() => {
    let cancelled = false
    const t = setTimeout(() => {
      listAdminContent(q)
        .then((r) => { if (!cancelled) { setResults(r.concepts); setTotal(r.total); setSearchError('') } })
        .catch((e) => { if (!cancelled) { setSearchError(getErrorMessage(e)); setResults([]) } })
    }, 200)
    return () => { cancelled = true; clearTimeout(t) }
  }, [q])

  useEffect(() => {
    if (view !== 'formats' || overrides !== null) return
    let cancelled = false
    listAdminContent('', { overridden: true })
      .then((r) => { if (!cancelled) setOverrides(r.concepts) })
      .catch((e) => { if (!cancelled) { setOverrideError(getErrorMessage(e)); setOverrides([]) } })
    return () => { cancelled = true }
  }, [view, overrides])

  const generate = useCallback(async (concept: AdminContentConcept, diff: number, n: number) => {
    setGenBusy(true)
    setGenError('')
    try {
      setGenerated(await generateLabQuestions(concept.id, n, diff))
    } catch (e) {
      setGenerated(null)
      setGenError(getErrorMessage(e))
    } finally {
      setGenBusy(false)
    }
  }, [])

  function pick(concept: AdminContentConcept) {
    setSelected(concept)
    setGenerated(null)
    void generate(concept, difficulty, count)
  }

  function tryFromTable(concept: AdminContentConcept) {
    setView('generate')
    pick(concept)
  }

  return (
    <AdminShell>
      <main className="mx-auto w-full max-w-[1100px] min-w-0 px-4 sm:px-6 py-8 sm:py-12 overflow-x-hidden">
        <SectionHeader label="Administration" title="Lab" />
  
        <p className="mb-4 max-w-3xl font-mono text-[11px] leading-relaxed text-mathua-secondary">
          Generate any concept&rsquo;s questions, answer them with the field a learner gets, and grade them with the
          same engine. Nothing here writes an attempt, XP, streak or mastery.
        </p>

        <div className="mb-4 flex flex-wrap gap-2">
          <button
            type="button"
            onClick={() => setView('generate')}
            className={
              view === 'generate'
                ? 'border border-mathua-blue text-mathua-blue px-3 h-9 font-mono text-xs'
                : 'border border-mathua-border text-mathua-secondary hover:text-mathua-blue px-3 h-9 font-mono text-xs'
            }
          >
            Generate
          </button>
          <button
            type="button"
            onClick={() => setView('formats')}
            className={
              view === 'formats'
                ? 'border border-mathua-blue text-mathua-blue px-3 h-9 font-mono text-xs'
                : 'border border-mathua-border text-mathua-secondary hover:text-mathua-blue px-3 h-9 font-mono text-xs'
            }
          >
            Formats that differ
          </button>
        </div>

        {view === 'generate' ? (
          <>
            <input
              type="search"
              value={q}
              onChange={(e) => setQ(e.target.value)}
              placeholder="Search concepts by id or label"
              className="mb-3 w-full border border-mathua-border bg-mathua-code px-3 h-11 font-mono text-xs text-mathua-primary"
            />

            {searchError && (
              <p role="alert" className="mb-4 border border-mathua-red bg-mathua-surface p-4 font-mono text-xs text-mathua-red">
                {searchError}
              </p>
            )}
            {!results && !searchError && <Loading label="LOADING CONCEPTS" />}

            {results && (
              <ul className="mb-6 max-h-72 overflow-y-auto border border-mathua-border">
                {results.map((c) => (
                  <li key={c.id} className="border-b border-mathua-border last:border-b-0">
                    <button
                      type="button"
                      onClick={() => pick(c)}
                      className={
                        selected?.id === c.id
                          ? 'w-full bg-mathua-hover px-4 py-2 text-left'
                          : 'w-full px-4 py-2 text-left hover:bg-mathua-hover'
                      }
                    >
                      <span className="flex items-baseline justify-between gap-3">
                        <span className="font-mono text-xs text-mathua-primary">{c.label}</span>
                        <span className="font-mono text-[10px] text-mathua-muted">{c.id}</span>
                      </span>
                      <span className="mt-1 block font-mono text-[10px] text-mathua-secondary">
                        {c.domain} · {c.grading_type} · field {c.answer_input_mode}
                        {c.answer_overridden ? ' (declared)' : ''}
                      </span>
                    </button>
                  </li>
                ))}
              </ul>
            )}

            {results && results.length === 0 && (
              <p className="mb-6 border border-mathua-border bg-mathua-surface p-4 font-mono text-xs text-mathua-secondary">
                No concept matches that search.
              </p>
            )}
            {results && results.length > 0 && (
              <p className="mb-4 font-mono text-[11px] text-mathua-muted">{total} matching</p>
            )}

            {selected && (
              <section className="border border-mathua-border bg-mathua-surface">
                <header className="border-b border-mathua-border px-4 py-3">
                  <div className="flex flex-wrap items-baseline justify-between gap-2">
                    <h2 className="font-mono text-xs text-mathua-primary">{selected.label}</h2>
                    <span className="font-mono text-[10px] text-mathua-muted">{selected.id}</span>
                  </div>
                  <p className="mt-1 font-mono text-[10px] text-mathua-secondary">
                    grading type {selected.grading_type} · field {selected.answer_input_mode} ·{' '}
                    {selected.answer_overridden ? 'declared by the generator' : 'derived from the grading type'}
                  </p>
                  <p className="mt-1 font-mono text-[11px] text-mathua-secondary">{selected.answer_hint}</p>
                </header>

                <div className="border-b border-mathua-border px-4 py-3">
                  <div className="flex flex-wrap items-center gap-4">
                    <div>
                      <span className="mr-2 font-mono text-[10px] uppercase tracking-wider text-mathua-muted">Difficulty</span>
                      {DIFFICULTIES.map((d) => (
                        <button
                          key={d}
                          type="button"
                          onClick={() => { setDifficulty(d); void generate(selected, d, count) }}
                          className={
                            difficulty === d
                              ? 'ml-1 border border-mathua-blue px-2 h-8 font-mono text-xs text-mathua-blue'
                              : 'ml-1 border border-mathua-border px-2 h-8 font-mono text-xs text-mathua-secondary hover:text-mathua-blue'
                          }
                        >
                          {d}
                        </button>
                      ))}
                    </div>
                    <div>
                      <span className="mr-2 font-mono text-[10px] uppercase tracking-wider text-mathua-muted">Count</span>
                      {COUNTS.map((n) => (
                        <button
                          key={n}
                          type="button"
                          onClick={() => { setCount(n); void generate(selected, difficulty, n) }}
                          className={
                            count === n
                              ? 'ml-1 border border-mathua-blue px-2 h-8 font-mono text-xs text-mathua-blue'
                              : 'ml-1 border border-mathua-border px-2 h-8 font-mono text-xs text-mathua-secondary hover:text-mathua-blue'
                          }
                        >
                          {n}
                        </button>
                      ))}
                    </div>
                    <button
                      type="button"
                      onClick={() => void generate(selected, difficulty, count)}
                      disabled={genBusy}
                      className="border border-mathua-primary px-3 h-9 font-mono text-xs text-mathua-primary hover:bg-mathua-hover disabled:opacity-50"
                    >
                      {genBusy ? 'Generating…' : 'Generate again'}
                    </button>
                    <Link
                      href={`/admin/lab/play?concept=${encodeURIComponent(selected.id)}`}
                      className="border border-mathua-border px-3 h-9 inline-flex items-center font-mono text-xs text-mathua-secondary hover:text-mathua-blue"
                    >
                      Open the learner view
                    </Link>
                  </div>
                </div>

                {genError && (
                  <p role="alert" className="border-b border-mathua-border px-4 py-3 font-mono text-xs text-mathua-red">
                    {genError}
                  </p>
                )}
                {genBusy && !generated && <Loading label="GENERATING QUESTIONS" />}

                {generated && (
                  <ul className="space-y-3 p-4">
                    {generated.questions.map((question, i) => (
                      <LabQuestionCard
                        key={`${generated.concept_id}-${i}-${question.question}`}
                        question={question}
                        conceptId={generated.concept_id}
                        gradingType={generated.grading_type}
                        index={i}
                      />
                    ))}
                  </ul>
                )}
              </section>
            )}
          </>
        ) : (
          <section>
            <p className="mb-3 max-w-3xl font-mono text-[11px] leading-relaxed text-mathua-secondary">
              Every concept whose served field or hint differs from what its grading type would give. These are the
              concepts where the grader and the input could disagree, which is how a numeric concept ended up with a
              digits-only keypad in front of an answer like &ldquo;5 R 2&rdquo;.
            </p>

            {overrideError && (
              <p role="alert" className="mb-4 border border-mathua-red bg-mathua-surface p-4 font-mono text-xs text-mathua-red">
                {overrideError}
              </p>
            )}
            {!overrides && !overrideError && <Loading label="LOADING FORMATS" />}

            {overrides && overrides.length === 0 && (
              <p className="border border-mathua-border bg-mathua-surface p-4 font-mono text-xs text-mathua-secondary">
                No concept overrides its grading type. That would mean the guard has nothing to protect.
              </p>
            )}

            {overrides && overrides.length > 0 && (
              <div className="overflow-x-auto border border-mathua-border">
                <table className="w-full border-collapse text-left">
                  <thead>
                    <tr className="border-b border-mathua-border bg-mathua-code">
                      <th className="px-3 py-2 font-mono text-[10px] uppercase tracking-wider text-mathua-muted">Concept</th>
                      <th className="px-3 py-2 font-mono text-[10px] uppercase tracking-wider text-mathua-muted">Grading type</th>
                      <th className="px-3 py-2 font-mono text-[10px] uppercase tracking-wider text-mathua-muted">Field</th>
                      <th className="px-3 py-2 font-mono text-[10px] uppercase tracking-wider text-mathua-muted">Hint</th>
                      <th className="px-3 py-2 font-mono text-[10px] uppercase tracking-wider text-mathua-muted">Try</th>
                    </tr>
                  </thead>
                  <tbody>
                    {overrides.map((c) => (
                      <tr key={c.id} className="border-b border-mathua-border last:border-b-0">
                        <td className="px-3 py-2">
                          <span className="block font-mono text-xs text-mathua-primary">{c.label}</span>
                          <span className="block font-mono text-[10px] text-mathua-muted">{c.id}</span>
                        </td>
                        <td className="px-3 py-2 font-mono text-[11px] text-mathua-secondary">{c.grading_type}</td>
                        <td className="px-3 py-2 font-mono text-[11px] text-mathua-primary">{c.answer_input_mode}</td>
                        <td className="px-3 py-2 font-mono text-[11px] text-mathua-secondary">{c.answer_hint}</td>
                        <td className="px-3 py-2">
                          <button
                            type="button"
                            onClick={() => tryFromTable(c)}
                            className="border border-mathua-border px-2 h-8 font-mono text-[11px] text-mathua-secondary hover:text-mathua-blue"
                          >
                            Generate
                          </button>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
            {overrides && overrides.length > 0 && (
              <p className="mt-3 font-mono text-[11px] text-mathua-muted">{overrides.length} concepts</p>
            )}
          </section>
        )}
      </main>
    </AdminShell>
  )
}
