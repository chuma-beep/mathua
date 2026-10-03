'use client'

import { useState, useEffect, useRef } from 'react'
import Link from 'next/link'
import { toast } from 'sonner'
import KatexContent from '../../components/KatexContent'
import Loading from '../../components/Loading'
import SectionHeader from '../../components/SectionHeader'
import BriefingCard, { REVIEW_BRIEFING } from '../../components/BriefingCard'
import ReportButton from '../../components/ReportButton'
import { startReviewSession, submitReviewAnswer } from '../../lib/api'
import { formatForGradingType, type AnswerFormat } from '../../lib/answerFormat'
import ChoiceOptions from '../../components/ChoiceOptions'
import SubmitErrorBlock, { toSubmitError, type SubmitError } from '../../components/SubmitErrorBlock'
import { MathAnswerInput, type MathFocusHandle } from '@/components/math/MathInput'

type Phase = 'loading' | 'intro' | 'run' | 'done'

// Due-review runner: each run pulls only concepts whose memory is fading
// (POST /api/reviews/session → /api/reviews/answer). Self-contained so the
// review page only mounts it.
export default function ReviewHost() {
  const [phase, setPhase] = useState<Phase>('intro')
  const [loading, setLoading] = useState(false)

  const reviewSessionId = useRef('')
  const reviewInputRef = useRef<MathFocusHandle | null>(null)
  const reviewShownAt = useRef<number | null>(null)
  const reviewConceptId = useRef('')
  const reviewGradingType = useRef('')
  const reviewAttemptId = useRef('')

  const [reviewQuestion, setReviewQuestion] = useState('')
  const [reviewConceptName, setReviewConceptName] = useState('')
  const [reviewCount, setReviewCount] = useState(0)
  const [reviewAnswerInput, setReviewAnswerInput] = useState('')
  const [reviewLastResult, setReviewLastResult] = useState<{ correct: boolean; feedback: string; explanation?: string; diagnosis?: string; xp?: number } | null>(null)
  const [reviewAccuracy, setReviewAccuracy] = useState({ correct: 0, total: 0 })
  const [reviewXp, setReviewXp] = useState(0)
  const [submitError, setSubmitError] = useState<SubmitError | null>(null)
  const [finished, setFinished] = useState(false)
  const [empty, setEmpty] = useState(false)
  const [answerFormat, setAnswerFormat] = useState<AnswerFormat>(() => formatForGradingType())
  // Next question staged from the submit response — revealed by goNext(),
  // never fetched. Cleared on advance, so double-press is a no-op.
  const pendingReviewNext = useRef<{
    question: string
    conceptId: string
    conceptName: string
    attemptId: string
    gradingType: string
  } | null>(null)

  async function startReview() {
    setPhase('loading')
    setLoading(true)
    try {
      const res = await startReviewSession()
      if (!res.question) {
        setEmpty(true)
        setPhase('done')
        return
      }
      reviewSessionId.current = res.session_id
      setReviewQuestion(res.question.question)
      reviewConceptId.current = res.question.concept_id
      setReviewConceptName(res.question.concept_name)
      reviewAttemptId.current = res.question.attempt_id ?? ''
      reviewShownAt.current = Date.now()
      setAnswerFormat(formatForGradingType(res.question.grading_type))
    reviewGradingType.current = res.question.grading_type
      setReviewCount(1)
      setReviewAccuracy({ correct: 0, total: 0 })
      setReviewXp(0)
      setReviewLastResult(null)
      setReviewAnswerInput('')
      setEmpty(false)
      pendingReviewNext.current = null
      setSubmitError(null)
      setFinished(false)
      setPhase('run')
    } catch {
      toast.error('Review failed to start — try again.')
      setPhase('intro')
    } finally {
      setLoading(false)
    }
  }

  // The review starts from the intro screen, not on mount — timing and
  // question reveal begin at Start, so the briefing costs nothing.
  useEffect(() => {
    if (phase === 'run' && !reviewLastResult && reviewQuestion) reviewInputRef.current?.focus()
  }, [phase, reviewLastResult, reviewQuestion])

  async function submitReviewFn() {
    if (!reviewAnswerInput.trim()) return
    setLoading(true)
    setSubmitError(null)
    try {
      const answer = reviewAnswerInput.trim()
      const elapsed = Math.max(0.5, (Date.now() - (reviewShownAt.current ?? Date.now())) / 1000)
      const data = await submitReviewAnswer(reviewSessionId.current, answer, elapsed, reviewAttemptId.current)
      // Grader unavailable: never counts as a miss; ask the student to retry.
      if (data.result?.ungraded) {
        setSubmitError({ message: 'Grading is temporarily unavailable — please try again.' })
        setLoading(false)
        return
      }
      const correct = data.result?.correct || false
      const feedback = data.result?.feedback || (correct ? 'Correct!' : 'Not quite.')
      const xp = data.result?.xp ?? 0
      setReviewAccuracy(prev => ({ correct: prev.correct + (correct ? 1 : 0), total: prev.total + 1 }))
      setReviewXp(prev => prev + xp)
      setReviewLastResult({ correct, feedback, explanation: data.result?.explanation, diagnosis: data.result?.diagnosis, xp })
      if (data.done || !data.next_question) {
        setFinished(true)
        setTimeout(() => {
          setPhase('done')
          setLoading(false)
        }, 800)
        return
      }
      // Manual advance: stage the prefetched next question, reveal on Next.
      pendingReviewNext.current = {
        question: data.next_question.question,
        conceptId: data.next_question.concept_id,
        conceptName: data.next_question.concept_name,
        attemptId: data.next_question.attempt_id ?? '',
        gradingType: data.next_question.grading_type ?? '',
      }
      setLoading(false)
    } catch (e) {
      setSubmitError(toSubmitError(e, 'Failed to submit review answer.'))
      setLoading(false)
    }
  }

  function applyReviewQuestion(question: string, cid: string, name: string, attemptId: string, gradingType: string) {
    setReviewQuestion(question)
    reviewConceptId.current = cid
    reviewShownAt.current = Date.now()
    setReviewConceptName(name)
    reviewAttemptId.current = attemptId
    setAnswerFormat(formatForGradingType(gradingType))
    reviewGradingType.current = gradingType
    setReviewCount(prev => prev + 1)
    setReviewLastResult(null)
    setReviewAnswerInput('')
    setSubmitError(null)
  }

  function restartReview() {
    setSubmitError(null)
    setFinished(false)
    setEmpty(false)
    void startReview()
  }

  function goNext() {
    const staged = pendingReviewNext.current
    if (!staged) return
    pendingReviewNext.current = null
    applyReviewQuestion(staged.question, staged.conceptId, staged.conceptName, staged.attemptId, staged.gradingType)
  }

  if (phase === 'intro') {
    return (
      <>
        <SectionHeader label="Review" title="Due for review" />
        <div className="max-w-2xl mx-auto min-w-0 overflow-hidden px-2 sm:px-0">
          <BriefingCard eyebrow="Before you start" items={REVIEW_BRIEFING} />
          <div className="mt-6 text-center">
            <button
              type="button"
              onClick={() => { void startReview() }}
              disabled={loading}
              className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue-faint rounded-none h-12 px-8 font-medium text-sm disabled:opacity-50"
            >
              Start review →
            </button>
          </div>
        </div>
      </>
    )
  }

  if (phase === 'loading') {
    return (
      <div className="py-20 text-center">
        <Loading label="LOADING REVIEW" />
      </div>
    )
  }

  if (phase === 'done') {
    return (
      <div className="max-w-2xl mx-auto text-center">
        {empty ? (
          <>
            <SectionHeader label="Review" title="All caught up" />
            <p className="font-mono text-sm text-mathua-secondary mt-4">Nothing is due for review right now. New reviews appear here as memories fade.</p>
            <div className="mt-6 flex gap-3 justify-center">
              {/* A review run is a Learn sub-mode that has finished. Sending
                  the learner into the reference library at that point is the
                  exact confusion this app is trying to avoid, so both exits go
                  back to doing. */}
              <Link href="/learn" className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue-faint rounded-none h-12 px-8 text-sm inline-flex items-center">Answer more questions →</Link>
              <Link href="/profile" className="border border-mathua-border text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue rounded-none h-12 px-8 text-sm inline-flex items-center">Back to Profile →</Link>
            </div>
          </>
        ) : (
          <>
            <SectionHeader label="Review complete" title={`${reviewAccuracy.correct}/${reviewAccuracy.total} correct`} />
            <p className="font-mono text-sm text-mathua-secondary mt-4">Reviews record progress like practice{reviewXp > 0 ? ` · +${reviewXp} XP this run` : ''}.</p>
            <div className="mt-6 flex gap-3 justify-center">
              <button type="button" onClick={restartReview} className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue-faint rounded-none h-12 px-8 text-sm">Review again →</button>
              <Link href="/profile" className="border border-mathua-border text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue rounded-none h-12 px-8 text-sm inline-flex items-center">Back to Profile →</Link>
            </div>
            <div className="mt-3 text-center">
              <Link href="/history" className="font-mono text-xs text-mathua-muted hover:text-mathua-primary">Every question you&apos;ve answered →</Link>
            </div>
          </>
        )}
      </div>
    )
  }

  return (
    <>
      <SectionHeader label={`Review question ${reviewCount}`} title={reviewConceptName} />
      <div className="max-w-2xl mx-auto min-w-0 overflow-hidden px-2 sm:px-0">
        <div className="mb-4 flex flex-wrap items-center gap-x-3 gap-y-1 text-xs font-mono text-mathua-muted justify-center">
          <span className={reviewAccuracy.correct / Math.max(reviewAccuracy.total, 1) >= 0.7 ? 'text-mathua-green' : ''}>{reviewAccuracy.correct}/{reviewAccuracy.total} correct</span>
          {reviewXp > 0 ? <span className="text-yellow-400">+{reviewXp} XP</span> : null}
        </div>
        <div className={`bg-mathua-surface border rounded-none p-4 sm:p-6 mb-6 transition-colors w-full max-w-full min-w-0 overflow-hidden ${reviewLastResult ? (reviewLastResult.correct ? 'border-green-500/40' : 'border-red-500/40') : 'border-mathua-border'}`}>
          <div className="bg-mathua-code border border-mathua-border rounded-none p-4 sm:p-6 text-center mb-4">
            <KatexContent className="text-mathua-primary text-lg font-mono font-light whitespace-pre-wrap break-words">{reviewQuestion}</KatexContent>
          </div>
          <ChoiceOptions question={reviewQuestion} value={reviewAnswerInput} onPick={setReviewAnswerInput} disabled={loading} />
          {!reviewLastResult ? (
            <>
              <form onSubmit={e => { e.preventDefault(); void submitReviewFn() }} className="flex flex-col sm:flex-row gap-3 min-w-0">
                <label htmlFor="review-answer" className="sr-only">Your answer</label>
                <MathAnswerInput
                  id="review-answer"
                  value={reviewAnswerInput}
                  onChange={setReviewAnswerInput}
                  onSubmit={() => void submitReviewFn()}
                  gradingType={reviewGradingType.current}
                  conceptId={reviewConceptId.current}
                  disabled={loading}
                  status={reviewLastResult ? (reviewLastResult.correct ? 'correct' : 'incorrect') : 'default'}
                  inputMode={answerFormat.inputMode}
              placeholder="Your answer..."
                  className="sm:flex-1"
                  focusRef={reviewInputRef}
                />
                <button type="submit" disabled={!reviewAnswerInput.trim() || loading} className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue-faint rounded-none h-12 px-6 font-medium text-sm disabled:opacity-50 shrink-0 w-auto self-end sm:self-auto">Check Answer</button>
              </form>
              <div className="mt-2 flex items-center justify-between gap-2">
                <p className="font-mono text-[11px] text-mathua-muted">{answerFormat.hint}</p>
              </div>
              <div className="mt-2 flex justify-end">
                <ReportButton
                  key={reviewQuestion}
                  conceptId={reviewConceptId.current}
                  kind="question"
                  question={reviewQuestion}
                  source="review"
                  sessionId={reviewSessionId.current}
                  attemptId={reviewAttemptId.current}
                />
              </div>
            </>
          ) : (
            <div className="animate-fadeIn">
              <p className={`text-base font-medium mb-2 text-center ${reviewLastResult.correct ? 'text-mathua-green' : 'text-mathua-red'}`}>{reviewLastResult.correct ? '✓ Correct!' : '✗ Not quite'}</p>
              {/* The explanation is the solution for this instance and is the
                  body of the result; `feedback` is the grader's verdict token
                  ("Incorrect"), so showing it above the working is noise. Fall
                  back to it only when there is no explanation at all. */}
              {reviewLastResult.explanation ? (
                <KatexContent className="text-mathua-secondary text-sm text-left whitespace-pre-wrap">{reviewLastResult.explanation}</KatexContent>
              ) : (
                <KatexContent className="text-mathua-secondary text-sm">{reviewLastResult.feedback}</KatexContent>
              )}
              {/* Names the mistake when it is certain; empty otherwise. */}
              {reviewLastResult.diagnosis && (
                <p className="mt-2 text-center font-mono text-[11px] text-mathua-muted">{reviewLastResult.diagnosis}</p>
              )}
            </div>
          )}
        </div>

        {submitError ? (
          <div className="mb-6">
            <SubmitErrorBlock
              error={submitError}
              onRetry={() => { void submitReviewFn() }}
              onRestart={restartReview}
              restartLabel="Restart review"
              retrying={loading}
            />
          </div>
        ) : reviewLastResult && !finished ? (
          <div className="mb-6 text-center">
            <button
              type="button"
              autoFocus
              onClick={goNext}
              className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue-faint rounded-none h-12 px-8 font-medium text-sm shrink-0"
            >
              Next →
            </button>
          </div>
        ) : reviewLastResult && finished ? (
          <div className="mb-6 text-center text-mathua-muted text-xs font-mono">
            Wrapping up…
          </div>
        ) : null}
      </div>
    </>
  )
}
