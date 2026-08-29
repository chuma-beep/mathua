'use client'

import { useState, useEffect, useCallback, useRef } from 'react'
import KatexContent from './KatexContent'
import { getLessonPractice, submitStudyAnswer, type PracticeQuestion } from '../lib/api'

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
  const loadTimes = useRef<Record<number, number>>({})

  const loadQuestions = useCallback(() => {
    setLoading(true)
    setError('')
    setAnswers({})
    setResults({})
    setXpMap({})
    setScore({ correct: 0, total: 0 })
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
    if (!userAnswer || results[i] !== undefined) return
    const q = questions[i]
    const elapsed = Math.max(0.5, (Date.now() - (loadTimes.current[i] ?? Date.now())) / 1000)
    try {
      const res = await submitStudyAnswer(conceptId, userAnswer, q.answer, elapsed)
      const key: 'correct' | 'incorrect' = res.correct ? 'correct' : 'incorrect'
      setResults(prev => ({ ...prev, [i]: key }))
      setXpMap(prev => ({ ...prev, [i]: res.xp ?? 0 }))
      setScore(prev => ({ correct: prev.correct + (res.correct ? 1 : 0), total: prev.total + 1 }))
    } catch {
      // Fallback to local grading if server unreachable
      const isCorrect = userAnswer.trim().toLowerCase() === q.answer.trim().toLowerCase()
      setResults(prev => ({ ...prev, [i]: isCorrect ? 'correct' : 'incorrect' }))
      setScore(prev => ({ correct: prev.correct + (isCorrect ? 1 : 0), total: prev.total + 1 }))
    }
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
        <p className="text-mathua-muted text-xs font-mono">No practice questions available for this concept.</p>
      </div>
    )
  }

  if (questions.length === 0) return null

  return (
    <div className="mt-8">
      <div className="flex items-center justify-between mb-4 border-b border-mathua-border pb-2">
        <h3 className="font-serif text-lg text-mathua-primary">
          Practice Questions ({questions.length})
        </h3>
        <div className="flex items-center gap-3">
          {score.total > 0 && (
            <span className="font-mono text-xs text-mathua-muted">
              {score.correct}/{score.total} correct
            </span>
          )}
          <button
            onClick={loadQuestions}
            className="font-mono text-[10px] text-mathua-blue hover:text-mathua-blue-hover transition-colors uppercase tracking-wider"
          >
            ↻ New
          </button>
        </div>
      </div>
      <div className="space-y-3">
        {questions.map((q, i) => {
          const result = results[i]
          const showAnswer = result !== undefined
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

                      <div className="mt-2 flex items-center gap-2">
                        <input
                          type="text"
                          value={answers[i] || ''}
                          onChange={e => setAnswers(prev => ({ ...prev, [i]: e.target.value }))}
                         onKeyDown={e => handleKeyDown(e, i)}
                         placeholder="Your answer…"
                        disabled={result !== undefined}
                        className={`flex-1 bg-mathua-bg border px-2.5 py-1.5 text-xs font-mono text-mathua-primary outline-none transition-colors rounded-none ${
                          result === 'correct'
                            ? 'border-green-500/60'
                            : result === 'incorrect'
                            ? 'border-red-500/60'
                            : 'border-mathua-border focus:border-mathua-blue'
                        }`}
                      />
                      {result === undefined && (
                        <button
                          onClick={() => handleCheck(i)}
                          className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white transition-colors px-3 py-1.5 text-xs font-mono rounded-none"
                        >
Submit
                        </button>
                      )}
                    </div>

                    {result === 'correct' && (
                      <p className="mt-2 text-xs font-mono text-green-400">
                        ✓ Correct!{xpMap[i] ? <span className="text-yellow-400"> +{xpMap[i]} XP</span> : null}
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
