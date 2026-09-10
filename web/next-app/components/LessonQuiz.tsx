'use client'

import { useEffect, useCallback, useRef, useReducer, useState } from 'react'
import KatexContent from './KatexContent'
import ReportButton from './ReportButton'
import { getLessonPractice, submitStudyAnswer, type PracticeQuestion } from '../lib/api'
import { applyResult, initialState, type StreakState } from '../lib/progression'

interface LessonQuizProps {
  conceptId: string
  limit?: number
}

interface QuizAnswerState {
  answers: Record<number, string>
  results: Record<number, 'correct' | 'incorrect'>
  xpMap: Record<number, number>
  checking: Record<number, boolean>
  score: { correct: number; total: number }
  streak: StreakState
}

type QuizAnswerAction =
  | { type: 'reset' }
  | { type: 'setAnswer'; index: number; value: string }
  | { type: 'checkStart'; index: number }
  | { type: 'checkEnd'; index: number }
  | { type: 'gradeServer'; index: number; correct: boolean; xp: number }
  | { type: 'gradeLocal'; index: number; correct: boolean }

function initialQuizAnswerState(): QuizAnswerState {
  return {
    answers: {},
    results: {},
    xpMap: {},
    checking: {},
    score: { correct: 0, total: 0 },
    streak: initialState(),
  }
}

function quizAnswerReducer(state: QuizAnswerState, action: QuizAnswerAction): QuizAnswerState {
  switch (action.type) {
    case 'reset':
      return initialQuizAnswerState()
    case 'setAnswer':
      return {
        ...state,
        answers: { ...state.answers, [action.index]: action.value },
      }
    case 'checkStart':
      return {
        ...state,
        checking: { ...state.checking, [action.index]: true },
      }
    case 'checkEnd':
      return {
        ...state,
        checking: { ...state.checking, [action.index]: false },
      }
    case 'gradeServer': {
      const key: 'correct' | 'incorrect' = action.correct ? 'correct' : 'incorrect'
      return {
        ...state,
        results: { ...state.results, [action.index]: key },
        xpMap: { ...state.xpMap, [action.index]: action.xp },
        score: {
          correct: state.score.correct + (action.correct ? 1 : 0),
          total: state.score.total + 1,
        },
        streak: applyResult(state.streak, action.correct),
      }
    }
    case 'gradeLocal': {
      const key: 'correct' | 'incorrect' = action.correct ? 'correct' : 'incorrect'
      return {
        ...state,
        results: { ...state.results, [action.index]: key },
        score: {
          correct: state.score.correct + (action.correct ? 1 : 0),
          total: state.score.total + 1,
        },
        streak: applyResult(state.streak, action.correct),
      }
    }
  }
}

export default function LessonQuiz({ conceptId, limit = 5 }: LessonQuizProps) {
  const [questions, setQuestions] = useState<PracticeQuestion[] | null>(null)
  const [error, setError] = useState('')
  const [quiz, dispatch] = useReducer(quizAnswerReducer, undefined, initialQuizAnswerState)
  const { answers, results, xpMap, checking, score, streak } = quiz
  const checkingRef = useRef<Record<number, boolean>>({})
  const loadTimes = useRef<Record<number, number>>({})

  // Derived during render: null questions means a fetch is in flight.
  const loading = questions === null

  const loadQuestions = useCallback(() => {
    setQuestions(null)
    setError('')
    dispatch({ type: 'reset' })
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
      .catch(e => {
        setError(e.message)
        setQuestions([])
      })
  }, [conceptId, limit])

  useEffect(() => {
    loadQuestions()
  }, [loadQuestions])

  async function handleCheck(i: number) {
    const userAnswer = (answers[i] || '').trim()
    if (!userAnswer || results[i] !== undefined || checkingRef.current[i]) return
    checkingRef.current[i] = true
    const q = questions?.[i]
    if (!q) {
      checkingRef.current[i] = false
      return
    }
    const elapsed = Math.max(0.5, (Date.now() - (loadTimes.current[i] ?? Date.now())) / 1000)
    dispatch({ type: 'checkStart', index: i })
    try {
      const res = await submitStudyAnswer(conceptId, userAnswer, q.answer, elapsed)
      dispatch({ type: 'gradeServer', index: i, correct: res.correct, xp: res.xp ?? 0 })
    } catch {
      // Fallback to local grading if server unreachable
      const correct = userAnswer.trim().toLowerCase() === q.answer.trim().toLowerCase()
      dispatch({ type: 'gradeLocal', index: i, correct })
    } finally {
      checkingRef.current[i] = false
      dispatch({ type: 'checkEnd', index: i })
    }
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
        <button type="button" onClick={loadQuestions} className="mt-3 font-mono text-[11px] text-mathua-blue hover:text-mathua-blue-hover uppercase tracking-wider">
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
            type="button"
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
              key={`${conceptId}::${q.question}::${q.answer}`}
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

                      <form
                        onSubmit={e => { e.preventDefault(); handleCheck(i) }}
                        className="mt-2 flex flex-col sm:flex-row sm:items-center gap-2"
                      >
                        <label htmlFor={`lesson-quiz-${conceptId}-${i}`} className="sr-only">Your answer</label>
                        <input
                          id={`lesson-quiz-${conceptId}-${i}`}
                          type="text"
                          value={answers[i] || ''}
                          onChange={e => dispatch({ type: 'setAnswer', index: i, value: e.target.value })}
                          placeholder="Your answer…"
                          aria-label={`Your answer for question ${i + 1}`}
                          enterKeyHint="go"
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
                          type="submit"
                          disabled={!!checking[i]}
                          aria-label={`Check Answer for question ${i + 1}`}
                          className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white transition-colors px-3 h-12 text-sm font-mono rounded-none disabled:opacity-50 disabled:cursor-not-allowed w-full sm:w-auto shrink-0"
                        >
                        {checking[i] ? 'Checking…' : 'Check Answer'}
                        </button>
                      )}
                      </form>

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

                    <div className="mt-2 flex justify-end">
                      <ReportButton
                        conceptId={conceptId}
                        kind={showAnswer && q.explanation ? 'explanation' : 'question'}
                        question={q.question}
                        expected={q.answer}
                        explanation={q.explanation}
                      />
                    </div>
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
