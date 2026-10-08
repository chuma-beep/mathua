'use client'

import { useState, useEffect, useRef } from 'react'
import Link from 'next/link'
import { toast } from 'sonner'
import KatexContent from '../../components/KatexContent'
import Loading from '../../components/Loading'
import SectionHeader from '../../components/SectionHeader'
import BriefingCard, { QUIZ_BRIEFING } from '../../components/BriefingCard'
import ReportButton from '../../components/ReportButton'
import { startQuizSession, submitQuizAnswer, skipQuizQuestion } from '../../lib/api'
import { formatForGradingType, type AnswerFormat } from '../../lib/answerFormat'
import ChoiceOptions from '../../components/ChoiceOptions'
import SubmitErrorBlock, {
  MAX_SKIPS,
  toSubmitError,
  type SubmitError,
} from '../../components/SubmitErrorBlock'
import { concepts as conceptsData } from '../../lib/conceptData'
import { MathAnswerInput, type MathFocusHandle } from '@/components/math/MathInput'

type Phase = 'loading' | 'intro' | 'quiz' | 'done'

// Actionable quiz every 50 XP (engine.QuizGateXP): timed closed-book, own
// grading path, guest unlimited retake. Self-contained so the goals page only
// mounts it.
export default function QuizHost() {
  const [phase, setPhase] = useState<Phase>('intro')
  const [loading, setLoading] = useState(false)

  const quizSessionId = useRef('')
  const quizInputRef = useRef<MathFocusHandle | null>(null)
  const quizShownAt = useRef<number | null>(null)
  const quizConceptId = useRef('')
  const quizGradingType = useRef('')

  const [quizQuestion, setQuizQuestion] = useState('')
  const [quizConceptName, setQuizConceptName] = useState('')
  const [quizCount, setQuizCount] = useState(0)
  const [quizAnswerInput, setQuizAnswerInput] = useState('')
  // Verdict only. A test verifies knowledge and does not teach, so the response carries no
  // solution and this holds none: `conceptId` is here so a miss can hand the learner to Learn,
  // which is the one surface allowed to explain.
  //
  // `dontKnow` is what the learner did, not what they scored: an admitted unknown still counts as
  // a miss, but the verdict must not read as "you were close and wrong".
  const [quizLastResult, setQuizLastResult] = useState<{ correct: boolean; feedback: string; conceptId: string; xp?: number; dontKnow?: boolean } | null>(null)
  const [quizAccuracy, setQuizAccuracy] = useState({ correct: 0, total: 0 })
  const [quizClosedBook, setQuizClosedBook] = useState(false)
  const [quizQuestionsTotal, setQuizQuestionsTotal] = useState(0)
  const [quizRemedial, setQuizRemedial] = useState<string[]>([])
  const [submitError, setSubmitError] = useState<SubmitError | null>(null)
  const [finished, setFinished] = useState(false)
  const [skipCount, setSkipCount] = useState(0)
  const [answerFormat, setAnswerFormat] = useState<AnswerFormat>(() => formatForGradingType())
  // Next question staged from the submit response — revealed by goNextQuiz(),
  // never fetched. Cleared on advance, so double-press is a no-op.
  const pendingQuizNext = useRef<{
    question: string
    conceptId: string
    conceptName: string
    gradingType: string
  } | null>(null)

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
      setAnswerFormat(formatForGradingType(res.grading_type))
    quizGradingType.current = res.grading_type
      setQuizCount(1)
      setQuizAccuracy({ correct: 0, total: 0 })
      setQuizLastResult(null)
      setQuizAnswerInput('')
      setQuizClosedBook(!!res.closed_book)
      setQuizQuestionsTotal(res.questions_total ?? 0)
      setQuizRemedial([])
      pendingQuizNext.current = null
      setSubmitError(null)
      setSkipCount(0)
      setFinished(false)
      setPhase('quiz')
    } catch {
      toast.error('Quiz failed to start — try again.')
      setPhase('done')
    } finally {
      setLoading(false)
    }
  }

  // The quiz starts from the intro screen, not on mount — timing and
  // question reveal begin at Start, so the briefing costs nothing.

  // Per-question countdown for the timed closed-book quiz. Informational:
  // the engine grades over-time answers as slow, it never blocks submission.
  // Keyboard flow: the answer form unmounts while feedback shows, so the
  // cursor is restored here — whenever a fresh question is revealed.
  useEffect(() => {
    if (phase === 'quiz' && !quizLastResult && quizQuestion) quizInputRef.current?.focus()
  }, [phase, quizLastResult, quizQuestion])

  async function submitQuizAnswerFn(dontKnow = false) {
    if (!dontKnow && !quizAnswerInput.trim()) return
    setLoading(true)
    setSubmitError(null)
    try {
      const answer = dontKnow ? '' : quizAnswerInput.trim()
      const elapsed = Math.max(0.5, (Date.now() - (quizShownAt.current ?? Date.now())) / 1000)
      const data = await submitQuizAnswer(quizSessionId.current, quizConceptId.current, answer, elapsed, dontKnow)
      const correct = data.correct || false
      // An admitted unknown is not a wrong answer, so it does not get the wrong answer's words.
      const feedback = dontKnow ? '' : data.feedback || (correct ? 'Correct!' : 'Not quite.')
      setQuizAccuracy(prev => ({ correct: prev.correct + (correct ? 1 : 0), total: prev.total + 1 }))
      setQuizLastResult({ correct, feedback, conceptId: quizConceptId.current, xp: data.xp, dontKnow })
      if (data.remedial?.length) {
        setQuizRemedial(prev => Array.from(new Set([...prev, ...(data.remedial ?? [])])))
      }

      if (data.done) {
        setFinished(true)
        setTimeout(() => {
          setPhase('done')
          setLoading(false)
        }, 800)
        return
      }
      // Manual advance: stage the prefetched next question, reveal on Next.
      // The staged question stays hidden and the countdown starts at reveal,
      // so reading feedback never burns question time.
      pendingQuizNext.current = {
        question: data.question || '',
        conceptId: data.concept_id || '',
        conceptName: data.concept_name || '',
        gradingType: data.grading_type || '',
      }
      setLoading(false)
    } catch (e) {
      setSubmitError(toSubmitError(e, 'Failed to submit quiz answer.'))
      setLoading(false)
    }
  }

  function applyQuizQuestion(question: string, cid: string, name: string, gradingType: string) {
    setQuizQuestion(question)
    quizConceptId.current = cid
    quizShownAt.current = Date.now()
    setQuizConceptName(name)
    setQuizCount(prev => prev + 1)
    setQuizLastResult(null)
    setQuizAnswerInput('')
    setAnswerFormat(formatForGradingType(gradingType))
    quizGradingType.current = gradingType
    setSubmitError(null)
  }

  async function skipQuiz() {
    if (skipCount >= MAX_SKIPS || loading) return
    setLoading(true)
    setSubmitError(null)
    try {
      const data = await skipQuizQuestion(quizSessionId.current)
      setSkipCount(c => c + 1)
      if (data.done) {
        setFinished(true)
        setQuizLastResult({ correct: false, feedback: 'Skipped — no XP awarded.', conceptId: quizConceptId.current, xp: 0 })
        setTimeout(() => {
          setPhase('done')
          setLoading(false)
        }, 800)
        return
      }
      applyQuizQuestion(
        data.question || '',
        data.concept_id || '',
        data.concept_name || '',
        data.grading_type || '',
      )
    } catch (e) {
      setSubmitError(toSubmitError(e, 'Failed to skip quiz question.'))
    } finally {
      setLoading(false)
    }
  }

  function restartQuiz() {
    setSubmitError(null)
    setSkipCount(0)
    setFinished(false)
    void startQuiz()
  }

  function goNextQuiz() {
    const staged = pendingQuizNext.current
    if (!staged) return
    pendingQuizNext.current = null
    applyQuizQuestion(staged.question, staged.conceptId, staged.conceptName, staged.gradingType)
  }

  if (phase === 'intro') {
    return (
      <>
        <SectionHeader label="Quiz" title="Mastery check" />
        <div className="max-w-2xl mx-auto min-w-0 overflow-hidden px-2 sm:px-0">
          <BriefingCard eyebrow="Before you start" items={QUIZ_BRIEFING} />
          <div className="mt-6 text-center">
            <button
              type="button"
              onClick={() => { void startQuiz() }}
              disabled={loading}
              className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue-faint rounded-none h-12 px-8 font-medium text-sm disabled:opacity-50"
            >
              Start quiz →
            </button>
          </div>
        </div>
      </>
    )
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
        {quizAccuracy.total === 0 && skipCount === 0 ? (
          <>
            <SectionHeader label="Quiz" title="No questions available" />
            <p className="font-mono text-sm text-mathua-secondary mt-4">There are no quiz questions available right now, try again later.</p>
            <div className="mt-6 flex gap-3 justify-center">
              <button type="button" onClick={startQuiz} className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue-faint rounded-none h-12 px-8 text-sm">Try again →</button>
              <Link href="/profile" className="border border-mathua-border text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue rounded-none h-12 px-8 text-sm inline-flex items-center">Back to Profile →</Link>
            </div>
          </>
        ) : (
          <>
            <SectionHeader label="Quiz complete" title={`${quizAccuracy.correct}/${quizAccuracy.total} correct`} />
            <p className="font-mono text-sm text-mathua-secondary mt-4">XP awarded per correct answer, more for faster answers. Retake anytime.</p>
            {quizRemedial.length > 0 && (
              <div className="mt-6 border border-mathua-border bg-mathua-surface p-4 text-left">
                <p className="font-mono text-xs uppercase tracking-wider text-mathua-muted mb-2">Focus next in Study</p>
                <ul className="space-y-1">
                  {quizRemedial.map(id => {
                    const label = conceptsData.find(c => c.id === id)?.label ?? id
                    return (
                      <li key={id}>
                        <Link href={`/learn?concept=${encodeURIComponent(id)}`} className="font-mono text-sm text-mathua-blue hover:underline">
                          {label} →
                        </Link>
                      </li>
                    )
                  })}
                </ul>
              </div>
            )}
            <div className="mt-6 flex gap-3 justify-center">
              <button type="button" onClick={startQuiz} className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue-faint rounded-none h-12 px-8 text-sm">Retake Quiz →</button>
              <Link href="/profile" className="border border-mathua-border text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue rounded-none h-12 px-8 text-sm inline-flex items-center">Back to Profile →</Link>
            </div>
            <div className="mt-3 text-center">
              <Link href="/history" className="font-mono text-xs text-mathua-muted hover:text-mathua-blue">Every question you&apos;ve answered →</Link>
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
        </div>
        <div className={`bg-mathua-surface border rounded-none p-4 sm:p-6 mb-6 transition-colors w-full max-w-full min-w-0 overflow-hidden ${quizLastResult ? (quizLastResult.correct ? 'border-mathua-green-faint' : 'border-mathua-red-faint') : 'border-mathua-border'}`}>
          <div className="bg-mathua-code border border-mathua-border rounded-none p-4 sm:p-6 text-center mb-4">
            <KatexContent className="text-mathua-primary text-lg font-mono font-light whitespace-pre-wrap break-words">{quizQuestion}</KatexContent>
          </div>
          <ChoiceOptions question={quizQuestion} value={quizAnswerInput} onPick={setQuizAnswerInput} disabled={loading} />
          {!quizLastResult ? (
            <>
              <form onSubmit={e => { e.preventDefault(); void submitQuizAnswerFn(false) }} className="flex flex-col sm:flex-row gap-3 min-w-0">
                <label htmlFor="quiz-answer" className="sr-only">Your answer</label>
                <MathAnswerInput
                  id="quiz-answer"
                  value={quizAnswerInput}
                  onChange={setQuizAnswerInput}
                  onSubmit={() => void submitQuizAnswerFn(false)}
                  gradingType={quizGradingType.current}
                  conceptId={quizConceptId.current}
                  disabled={loading}
                  status={quizLastResult ? (quizLastResult.correct ? 'correct' : 'incorrect') : 'default'}
                  inputMode={answerFormat.inputMode}
              placeholder="Your answer..."
                  className="sm:flex-1"
                  focusRef={quizInputRef}
                />
                <button type="submit" disabled={!quizAnswerInput.trim() || loading} className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue-faint rounded-none h-12 px-6 font-medium text-sm disabled:opacity-50 shrink-0 w-auto self-end sm:self-auto">Check Answer</button>
              </form>
              <div className="mt-2 flex items-center justify-between gap-2">
                <p className="font-mono text-[11px] text-mathua-muted">{answerFormat.hint}</p>
                <button
                  type="button"
                  onClick={() => { void submitQuizAnswerFn(true) }}
                  disabled={loading}
                  className="shrink-0 font-mono text-[11px] text-mathua-muted hover:text-mathua-red underline underline-offset-2 disabled:opacity-50"
                >
                  I don&apos;t know
                </button>
              </div>
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
            <div className="animate-fadeIn">
              <p className={`text-base font-medium mb-2 ${quizLastResult.correct ? 'text-mathua-green' : 'text-mathua-red'}`}>{quizLastResult.dontKnow ? 'No answer' : quizLastResult.correct ? '✓ Correct!' : '✗ Not quite'}</p>
              {/* Assessment state, and nothing else. This used to render the served worked
                  solution and a one-sentence diagnosis on both verdicts; the diagnosis was
                  answer-revealing on its own ("you are exactly one away"), and the solution
                  handed over the answer mid-test, so the remaining questions measured whether
                  the learner had read it rather than what they knew. Instruction lives in Learn,
                  reached deliberately, one concept at a time. */}
              <KatexContent className="text-mathua-secondary text-sm">{quizLastResult.feedback}</KatexContent>
              {!quizLastResult.correct && quizLastResult.conceptId && (
                <Link
                  href={`/learn?concept=${encodeURIComponent(quizLastResult.conceptId)}`}
                  className="mt-3 inline-flex items-center border border-mathua-blue text-mathua-blue hover:bg-mathua-blue-faint px-4 py-2 font-mono text-[11px] uppercase tracking-[0.18em]"
                >
                  Learn this concept →
                </Link>
              )}
            </div>
          )}
        </div>

        {submitError ? (
          <div className="mb-6">
            <SubmitErrorBlock
              error={submitError}
              onRetry={() => { void submitQuizAnswerFn(false) }}
              onSkip={skipQuiz}
              skipsLeft={MAX_SKIPS - skipCount}
              onRestart={restartQuiz}
              restartLabel="Restart quiz"
              retrying={loading}
            />
          </div>
        ) : quizLastResult && !finished ? (
          <div className="mb-6 text-center">
            <button
              type="button"
              autoFocus
              onClick={goNextQuiz}
              className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue-faint rounded-none h-12 px-8 font-medium text-sm shrink-0"
            >
              Next →
            </button>
          </div>
        ) : quizLastResult && finished ? (
          <div className="mb-6 text-center text-mathua-muted text-xs font-mono">
            Wrapping up…
          </div>
        ) : null}
      </div>
    </>
  )
}
