'use client'

import { useState, useEffect, useCallback } from 'react'
import { getLessonPractice, type PracticeQuestion } from '../lib/api'

interface LessonQuizProps {
  conceptId: string
  limit?: number
}

function normalize(s: string): string {
  return s.replace(/\s+/g, ' ').trim().toLowerCase()
}

function stripOuterParens(s: string): string {
  s = s.trim()
  if (s.startsWith('(') && s.endsWith(')')) return s.slice(1, -1).trim()
  return s
}

function answersMatch(userAnswer: string, expected: string): boolean {
  const a = normalize(userAnswer)
  const b = normalize(expected)
  if (a === b) return true
  if (stripOuterParens(a) === b) return true
  if (a === stripOuterParens(b)) return true
  const aSet = new Set(a.split(',').map(s => s.trim()))
  const bSet = new Set(b.split(',').map(s => s.trim()))
  if (aSet.size === bSet.size && Array.from(aSet).every(v => bSet.has(v))) return true
  return false
}

export default function LessonQuiz({ conceptId, limit = 5 }: LessonQuizProps) {
  const [questions, setQuestions] = useState<PracticeQuestion[]>([])
  const [answers, setAnswers] = useState<Record<number, string>>({})
  const [results, setResults] = useState<Record<number, 'correct' | 'incorrect'>>({})
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [score, setScore] = useState({ correct: 0, total: 0 })

  const loadQuestions = useCallback(() => {
    setLoading(true)
    setError('')
    setAnswers({})
    setResults({})
    setScore({ correct: 0, total: 0 })
    getLessonPractice(conceptId, limit)
      .then(res => setQuestions(res.questions))
      .catch(e => setError(e.message))
      .finally(() => setLoading(false))
  }, [conceptId, limit])

  useEffect(() => {
    loadQuestions()
  }, [loadQuestions])

  function handleCheck(i: number) {
    const userAnswer = (answers[i] || '').trim()
    if (!userAnswer) return
    const q = questions[i]
    const match = answersMatch(userAnswer, q.answer)
    if (match) {
      setResults(prev => ({ ...prev, [i]: 'correct' }))
      setScore(prev => ({ ...prev, correct: prev.correct + 1, total: prev.total + 1 }))
    } else {
      setResults(prev => ({ ...prev, [i]: 'incorrect' }))
      setScore(prev => ({ ...prev, total: prev.total + 1 }))
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
                    <p className="text-sm text-mathua-primary font-mono whitespace-pre-wrap">
                      {q.question}
                    </p>

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
                      <p className="mt-2 text-xs font-mono text-green-400">✓ Correct!</p>
                    )}
                    {result === 'incorrect' && (
                      <p className="mt-2 text-xs font-mono text-red-400">
                        ✗ Expected: {q.answer}
                      </p>
                    )}

                    {showAnswer && q.explanation && (
                      <p className="mt-1.5 text-xs font-mono text-mathua-secondary">
                        {q.explanation}
                      </p>
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
