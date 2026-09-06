'use client'

import { useState, useEffect, useCallback, useRef } from 'react'
import KatexContent from './KatexContent'
import { getLessonPractice, submitStudyAnswer, type PracticeQuestion } from '../lib/api'
import { applyResult, initialState, type StreakState } from '../lib/progression'

interface LessonQuizProps {
  conceptId: string
  limit?: number
}

export default function LessonQuiz({ conceptId, limit = 5 }: LessonQuizProps) {
  const [questions, setQuestions] = useState<PracticeQuestion[]>([])
  const [answers, setAnswers] = useState<Record<number, string>>({})
  const [results, setResults] = useState<Record<number, 'correct' | 'incorrect'>>({})
  const [xpMap, setXpMap] = useState<Record<number, number>>({})
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [score, setScore] = useState({ correct: 0, total: 0 })
  const [streak, setStreak] = useState<StreakState>(initialState())
  const [checking, setChecking] = useState<Record<number, boolean>>({})
  const checkingRef = useRef<Record<number, boolean>>({})
  const loadTimes = useRef<Record<number, number>>({})

  const loadQuestions = useCallback(() => {
    setLoading(true)
    setError('')
    setAnswers({})
    setResults({})
    setXpMap({})
    setScore({ correct: 0, total: 0 })
    setStreak(initialState())
    setChecking({})
    checkingRef.current = {}
    loadTimes.current = {}
    getLessonPractice(conceptId, limit)
      .then(res => {
        setQuestions(res.questions)
        const now = Date.now()
        const map: Record<number, number> = {}
        res.questions.forEach((_, i) => { map[i] = now })
        loadTimes.current = map
      })
      .catch(e => setError(e.message))
      .finally(() => setLoading(false))
  }, [conceptId, limit])

  useEffect(() => {
    loadQuestions()
  }, [loadQuestions])

  async function handleCheck(i: number) {
    const userAnswer = (answers[i] || '').trim()
    if (!userAnswer || results[i] !== undefined || checkingRef.current[i]) return
    checkingRef.current[i] = true
    const q = questions[i]
    const elapsed = Math.max(0.5, (Date.now() - (loadTimes.current[i] ?? Date.now())) / 1000)
    setChecking(prev => ({ ...prev, [i]: true }))
    let correct = false
    try {
      const res = await submitStudyAnswer(conceptId, userAnswer, q.answer, elapsed)
      correct = res.correct
      const key: 'correct' | 'incorrect' = res.correct ? 'correct' : 'incorrect'
      setResults(prev => ({ ...prev, [i]: key }))
      setXpMap(prev => ({ ...prev, [i]: res.xp ?? 0 }))
      setScore(prev => ({ correct: prev.correct + (res.correct ? 1 : 0), total: prev.total + 1 }))
    } catch {
      // Fallback to local grading if server unreachable
      correct = userAnswer.trim().toLowerCase() === q.answer.trim().toLowerCase()
      setResults(prev => ({ ...prev, [i]: correct ? 'correct' : 'incorrect' }))
      setScore(prev => ({ correct: prev.correct + (correct ? 1 : 0), total: prev.total + 1 }))
    } finally {
      checkingRef.current[i] = false
      setChecking(prev => ({ ...prev, [i]: false }))
    }
    setStreak(prev => applyResult(prev, correct))
  }

  function handleKeyDown(e: React.KeyboardEvent, i: number) {
    if (e.key === 'Enter') handleCheck(i)
  }

  if (loading) {
    return (
      <div className="mt-8 p-6 border border-mathua-border bg-mathua-surface">
        <p className="text-mathua-muted text-xs font-mono">Loading practice questions…</p>
      </div>
    )
  }

  if (error) {
    return (
      <div className="mt-8 p-6 border border-mathua-border bg-mathua-surface">
        <p className="text-mathua-muted text-xs font-mono">Couldn&apos;t load practice questions — check your connection.</p>
        <button onClick={loadQuestions} className="mt-3 font-mono text-[11px] text-mathua-blue hover:text-mathua-blue-hover uppercase tracking-wider">
          Retry →
        </button>
      </div>
    )
  }

  if (questions.length === 0) {
    return (
      <div className="mt-8 p-6 border border-mathua-border bg-mathua-surface">
        <p className="text-mathua-muted text-xs font-mono">No practice questions available for this concept yet.</p>
      </div>
    )
  }

  return (
    <div className="mt-8">
      <div className="flex items-center justify-between mb-4 border-b border-mathua-border pb-2">
        <h3 className="font-serif text-lg text-mathua-primary">
          Practice Questions ({questions.length})
        </h3>
        <div className="flex items-center gap-3">
          {!streak.advanced && streak.consecutive > 0 && (
            <span title="Answer 2 in a row to advance to the next concept" className="font-mono text-[10px] text-mathua-blue uppercase tracking-wider">
              streak {streak.consecutive}/2
            </span>
          )}
          {score.total > 0 && (
            <span className="font-mono text-xs text-mathua-muted">
              {score.correct}/{score.total} correct
            </span>
          )}
          <button
            onClick={loadQuestions}
            aria-label="Load new questions"
            title="Load new questions"
            className="font-mono text-[10px] text-mathua-blue hover:text-mathua-blue-hover transition-colors uppercase tracking-wider"
          >
            New questions
          </button>
        </div>
      </div>

      {streak.advanced && (
        <div className="mb-4 border border-mathua-green bg-mathua-surface p-4 flex items-center gap-3">
          <span className="text-green-400 text-lg">✓</span>
          <div>
            <p className="font-mono text-xs text-mathua-green">2 in a row — concept advanced</p>
            <p className="font-mono text-[10px] text-mathua-secondary mt-0.5">
              {score.correct}/{score.total} correct · next concept unlocked
            </p>
          </div>
        </div>
      )}

      <div className="space-y-3">
        {questions.map((q, i) => {
          const result = results[i]
          const showAnswer = result !== undefined
          const locked = streak.advanced
          return (
            <div
              key={i}
              className={`border rounded-none transition-colors ${
                result === 'correct'
                  ? 'border-green-500/40 bg-mathua-surface'
                  : result === 'incorrect'
                  ? 'border-red-500/40 bg-mathua-surface'
                  : 'border-mathua-border bg-mathua-surface'
              }`}
            >
              <div className="p-4">
                <div className="flex items-start gap-3">
                  <span className="font-mono text-[10px] text-mathua-muted mt-0.5 shrink-0 w-5">
                    {i + 1}.
                  </span>
                  <div className="flex-1 min-w-0">
                    <div>
                        <KatexContent className="text-sm text-mathua-primary font-mono whitespace-pre-wrap">
                          {q.question}
                        </KatexContent>
                      </div>

                      <div className="mt-2 flex flex-col sm:flex-row sm:items-center gap-2">
                        <input
                          type="text"
                          value={answers[i] || ''}
                          onChange={e => setAnswers(prev => ({ ...prev, [i]: e.target.value }))}
                         onKeyDown={e => handleKeyDown(e, i)}
                          placeholder="Your answer…"
                        disabled={result !== undefined || locked}
                        className={`flex-1 min-w-0 bg-mathua-bg border px-4 h-24 sm:h-12 text-base font-mono text-mathua-primary placeholder:text-mathua-muted outline-none transition-colors rounded-none ${
                          result === 'correct'
                            ? 'border-green-500/60'
                            : result === 'incorrect'
                            ? 'border-red-500/60'
                            : 'border-mathua-border focus:border-mathua-blue'
                        }`}
                      />
                      {result === undefined && !locked && (
                        <button
                          onClick={() => handleCheck(i)}
                          disabled={!!checking[i]}
                          className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white transition-colors px-3 h-12 text-sm font-mono rounded-none disabled:opacity-50 disabled:cursor-not-allowed w-full sm:w-auto shrink-0"
                        >
                        {checking[i] ? 'Checking…' : 'Submit'}
                        </button>
                      )}
                    </div>

                    {result === 'correct' && (
                      <p className="mt-2 text-xs font-mono text-green-400">
                        ✓ Correct!{xpMap[i] ? <span title="Experience points — progress toward your daily goal" className="text-yellow-400"> +{xpMap[i]} XP</span> : null}
                      </p>
                    )}
                    {result === 'incorrect' && (
                      <p className="mt-2 text-xs font-mono text-red-400">
                        ✗ Expected: <KatexContent>{q.answer}</KatexContent>
                      </p>
                    )}

                    {showAnswer && q.explanation && (
                      <KatexContent className="mt-1.5 text-xs font-mono text-mathua-secondary">
                        {q.explanation}
                      </KatexContent>
                    )}
                  </div>
                </div>
              </div>
            </div>
          )
        })}
      </div>
    </div>
  )
}
