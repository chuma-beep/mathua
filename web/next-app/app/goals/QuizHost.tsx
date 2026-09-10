'use client'

import { useState, useEffect, useRef } from 'react'
import Link from 'next/link'
import { toast } from 'sonner'
import KatexContent from '../../components/KatexContent'
import Loading from '../../components/Loading'
import SectionHeader from '../../components/SectionHeader'
import SymbolPalette from '../../components/SymbolPalette'
import ReportButton from '../../components/ReportButton'
import { startQuizSession, submitQuizAnswer } from '../../lib/api'
import { concepts as conceptsData } from '../../lib/conceptData'

type Phase = 'loading' | 'quiz' | 'done'

// Actionable quiz every 150 XP: timed closed-book, own grading path, guest
// unlimited retake. Self-contained so the goals page only mounts it.
export default function QuizHost() {
  const [phase, setPhase] = useState<Phase>('loading')
  const [loading, setLoading] = useState(false)

  const quizSessionId = useRef('')
  const quizInputRef = useRef<HTMLInputElement>(null)
  const quizShownAt = useRef<number | null>(null)
  const quizConceptId = useRef('')
  const startedRef = useRef(false)

  const [quizQuestion, setQuizQuestion] = useState('')
  const [quizConceptName, setQuizConceptName] = useState('')
  const [quizCount, setQuizCount] = useState(0)
  const [quizAnswerInput, setQuizAnswerInput] = useState('')
  const [quizLastResult, setQuizLastResult] = useState<{ correct: boolean; feedback: string; xp?: number } | null>(null)
  const [quizAccuracy, setQuizAccuracy] = useState({ correct: 0, total: 0 })
  const [quizTimeLimit, setQuizTimeLimit] = useState(0)
  const [quizRemaining, setQuizRemaining] = useState(0)
  const [quizClosedBook, setQuizClosedBook] = useState(false)
  const [quizQuestionsTotal, setQuizQuestionsTotal] = useState(0)
  const [quizRemedial, setQuizRemedial] = useState<string[]>([])

  async function startQuiz() {
    setPhase('loading')
    setLoading(true)
    try {
      const res = await startQuizSession()
      if (res.done) {
        setPhase('done')
        return
      }
      quizSessionId.current = res.session_id
      setQuizQuestion(res.question || '')
      quizConceptId.current = res.concept_id || ''
      setQuizConceptName(res.concept_name || '')
      quizShownAt.current = Date.now()
      setQuizCount(1)
      setQuizAccuracy({ correct: 0, total: 0 })
      setQuizLastResult(null)
      setQuizAnswerInput('')
      setQuizTimeLimit(res.time_limit_seconds ?? 0)
      setQuizRemaining(res.time_limit_seconds ?? 0)
      setQuizClosedBook(!!res.closed_book)
      setQuizQuestionsTotal(res.questions_total ?? 0)
      setQuizRemedial([])
      setPhase('quiz')
    } catch {
      toast.error('Quiz failed to start — try again.')
      setPhase('done')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    if (startedRef.current) return
    startedRef.current = true
    void startQuiz()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  // Per-question countdown for the timed closed-book quiz. Informational:
  // the engine grades over-time answers as slow, it never blocks submission.
  useEffect(() => {
    if (phase !== 'quiz' || quizLastResult || quizTimeLimit <= 0) return
    const started = quizShownAt.current ?? Date.now()
    const tick = () => setQuizRemaining(Math.max(0, quizTimeLimit - (Date.now() - started) / 1000))
    tick()
    const id = window.setInterval(tick, 250)
    return () => window.clearInterval(id)
  }, [phase, quizLastResult, quizTimeLimit, quizQuestion])

  async function submitQuizAnswerFn() {
    if (!quizAnswerInput.trim()) return
    setLoading(true)
    try {
      const answer = quizAnswerInput.trim()
      const elapsed = Math.max(0.5, (Date.now() - (quizShownAt.current ?? Date.now())) / 1000)
      const data = await submitQuizAnswer(quizSessionId.current, quizConceptId.current, answer, elapsed)
      const correct = data.correct || false
      const feedback = data.feedback || (correct ? 'Correct!' : 'Not quite.')
      setQuizAccuracy(prev => ({ correct: prev.correct + (correct ? 1 : 0), total: prev.total + 1 }))
      setQuizLastResult({ correct, feedback, xp: data.xp })
      if (data.remedial?.length) {
        setQuizRemedial(prev => Array.from(new Set([...prev, ...(data.remedial ?? [])])))
      }

      if (data.done) {
        setTimeout(() => {
          setPhase('done')
          setLoading(false)
        }, 800)
        return
      }
      setTimeout(() => {
        setQuizQuestion(data.question || '')
        quizConceptId.current = data.concept_id || ''
        quizShownAt.current = Date.now()
        setQuizConceptName(data.concept_name || '')
        setQuizCount(prev => prev + 1)
        setQuizLastResult(null)
        setQuizAnswerInput('')
        setQuizTimeLimit(data.time_limit_seconds ?? 0)
        setQuizRemaining(data.time_limit_seconds ?? 0)
        setLoading(false)
      }, 1200)
    } catch {
      toast.error('Failed to submit quiz answer.')
      setLoading(false)
    }
  }

  if (phase === 'loading') {
    return (
      <div className="py-20 text-center">
        <Loading label="LOADING QUIZ" />
      </div>
    )
  }

  if (phase === 'done') {
    return (
      <div className="max-w-2xl mx-auto text-center">
        {quizAccuracy.total === 0 ? (
          <>
            <SectionHeader label="Quiz" title="No questions available" />
            <p className="font-mono text-sm text-mathua-secondary mt-4">There are no quiz questions available right now, try again later.</p>
            <div className="mt-6 flex gap-3 justify-center">
              <button type="button" onClick={startQuiz} className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white rounded-none h-12 px-8 text-sm">Try again →</button>
              <Link href="/profile" className="border border-mathua-border text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue rounded-none h-12 px-8 text-sm inline-flex items-center">Back to Profile →</Link>
            </div>
          </>
        ) : (
          <>
            <SectionHeader label="Quiz complete" title={`${quizAccuracy.correct}/${quizAccuracy.total} correct`} />
            <p className="font-mono text-sm text-mathua-secondary mt-4">TaskQuiz 20 XP awarded per correct, retake anytime.</p>
            {quizRemedial.length > 0 && (
              <div className="mt-6 border border-mathua-border bg-mathua-surface p-4 text-left">
                <p className="font-mono text-xs uppercase tracking-wider text-mathua-muted mb-2">Focus next in Study</p>
                <ul className="space-y-1">
                  {quizRemedial.map(id => {
                    const label = conceptsData.find(c => c.id === id)?.label ?? id
                    return (
                      <li key={id}>
                        <Link href={`/study?concept=${encodeURIComponent(id)}`} className="font-mono text-sm text-mathua-blue hover:underline">
                          {label} →
                        </Link>
                      </li>
                    )
                  })}
                </ul>
              </div>
            )}
            <div className="mt-6 flex gap-3 justify-center">
              <button type="button" onClick={startQuiz} className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white rounded-none h-12 px-8 text-sm">Retake Quiz →</button>
              <Link href="/profile" className="border border-mathua-border text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue rounded-none h-12 px-8 text-sm inline-flex items-center">Back to Profile →</Link>
            </div>
          </>
        )}
      </div>
    )
  }

  return (
    <>
      <SectionHeader label={`Quiz question ${quizCount}${quizQuestionsTotal > 0 ? ` of ${quizQuestionsTotal}` : ''}`} title={quizConceptName} />
      <div className="max-w-2xl mx-auto min-w-0 overflow-hidden px-2 sm:px-0">
        <div className="mb-4 flex flex-wrap items-center gap-x-3 gap-y-1 text-xs font-mono text-mathua-muted justify-center">
          <span className={quizAccuracy.correct / Math.max(quizAccuracy.total, 1) >= 0.7 ? 'text-mathua-green' : ''}>{quizAccuracy.correct}/{quizAccuracy.total} correct</span>
          {quizLastResult?.xp ? <span className="text-yellow-400">+{quizLastResult.xp} XP (TaskQuiz 20)</span> : null}
          {quizClosedBook ? <span className="text-mathua-muted">Closed book</span> : null}
          {quizTimeLimit > 0 && !quizLastResult ? (
            <span className={quizRemaining <= 0 ? 'text-mathua-red' : quizRemaining <= 3 ? 'text-yellow-400' : 'text-mathua-muted'}>
              {quizRemaining <= 0 ? 'Time up (counts as slow)' : `${Math.ceil(quizRemaining)}s`}
            </span>
          ) : null}
        </div>
        <div className={`bg-mathua-surface border rounded-none p-4 sm:p-6 mb-6 transition-colors w-full max-w-full min-w-0 overflow-hidden ${quizLastResult ? (quizLastResult.correct ? 'border-green-500/40' : 'border-red-500/40') : 'border-mathua-border'}`}>
          <div className="bg-mathua-code border border-mathua-border rounded-none p-4 sm:p-6 text-center mb-4">
            <KatexContent className="text-mathua-primary text-lg font-mono font-light whitespace-pre-wrap break-words">{quizQuestion}</KatexContent>
          </div>
          {!quizLastResult ? (
            <>
              <form onSubmit={e => { e.preventDefault(); submitQuizAnswerFn() }} className="flex flex-col sm:flex-row gap-3 min-w-0">
                <label htmlFor="quiz-answer" className="sr-only">Your answer</label>
                <input ref={quizInputRef} id="quiz-answer" type="text" value={quizAnswerInput} onChange={e => setQuizAnswerInput(e.target.value)} placeholder="Your answer..." enterKeyHint="go" disabled={loading} className="flex-1 min-w-0 bg-mathua-code border border-mathua-border rounded-none h-24 sm:h-12 px-4 font-mono text-base text-mathua-primary placeholder:text-mathua-muted focus:outline-none focus:border-mathua-blue" />
                <button type="submit" disabled={!quizAnswerInput.trim() || loading} className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white rounded-none h-12 px-8 font-medium text-sm disabled:opacity-50">Check Answer</button>
              </form>
              <SymbolPalette targetRef={quizInputRef} onInsert={setQuizAnswerInput} />
              <div className="mt-2 flex justify-end">
                <ReportButton
                  key={quizQuestion}
                  conceptId={quizConceptId.current}
                  kind="question"
                  question={quizQuestion}
                  source="quiz"
                  sessionId={quizSessionId.current}
                />
              </div>
            </>
          ) : (
            <div className="animate-fadeIn text-center">
              <p className={`text-base font-medium mb-2 ${quizLastResult.correct ? 'text-mathua-green' : 'text-mathua-red'}`}>{quizLastResult.correct ? '✓ Correct!' : '✗ Not quite'}</p>
              <KatexContent className="text-mathua-secondary text-sm">{quizLastResult.feedback}</KatexContent>
              {loading && <p className="text-mathua-muted text-xs mt-2"><Loading inline size={11} /> Loading next…</p>}
            </div>
          )}
        </div>
      </div>
    </>
  )
}
