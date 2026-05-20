'use client'

import { useState, useEffect } from 'react'
import { getLessonPractice, type PracticeQuestion } from '../lib/api'

interface LessonQuizProps {
  conceptId: string
  limit?: number
}

export default function LessonQuiz({ conceptId, limit = 5 }: LessonQuizProps) {
  const [questions, setQuestions] = useState<PracticeQuestion[]>([])
  const [revealed, setRevealed] = useState<Set<number>>(new Set())
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    setLoading(true)
    setError('')
    setRevealed(new Set())
    getLessonPractice(conceptId, limit)
      .then(res => setQuestions(res.questions))
      .catch(e => setError(e.message))
      .finally(() => setLoading(false))
  }, [conceptId, limit])

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

  function toggleReveal(i: number) {
    setRevealed(prev => {
      const next = new Set(prev)
      if (next.has(i)) next.delete(i)
      else next.add(i)
      return next
    })
  }

  return (
    <div className="mt-8">
      <h3 className="font-serif text-lg text-mathua-primary mb-4 border-b border-mathua-border pb-2">
        Practice Questions ({questions.length})
      </h3>
      <div className="space-y-3">
        {questions.map((q, i) => {
          const isRevealed = revealed.has(i)
          return (
            <div
              key={i}
              className="border border-mathua-border bg-mathua-surface rounded-none"
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
                    <button
                      onClick={() => toggleReveal(i)}
                      className={`mt-2 text-xs font-mono transition-colors ${
                        isRevealed ? 'text-mathua-green' : 'text-mathua-blue hover:text-mathua-blue-hover'
                      }`}
                    >
                      {isRevealed ? '▲ Hide answer' : '▼ Reveal answer'}
                    </button>
                    {isRevealed && (
                      <div className="mt-3 pt-3 border-t border-mathua-border space-y-1">
                        <p className="text-xs font-mono text-mathua-green">
                          Answer: {q.answer}
                        </p>
                        {q.explanation && (
                          <p className="text-xs font-mono text-mathua-secondary">
                            {q.explanation}
                          </p>
                        )}
                      </div>
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
