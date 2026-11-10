'use client'

import { useCallback, useEffect, useRef, useState } from 'react'
import Link from 'next/link'
import KatexContent from '../../components/KatexContent'
import SectionHeader from '../../components/SectionHeader'
import XPAdjustmentNote from '../../components/XPAdjustmentNote'
import XPTaskSummary from '../../components/XPTaskSummary'
import { MathAnswerInput, type MathFocusHandle } from '../../components/math/MathInput'
import { formatForQuestion } from '../../lib/answerFormat'
import { stemWithoutChoices } from '../../lib/choices'
import {
  finalizeTimedQuiz,
  getTimedQuiz,
  getTimedQuizResult,
  skipTimedQuizItem,
  startTimedQuizSession,
  submitTimedQuizAnswer,
  type TimedQuizAnswerRes,
  type TimedQuizQuestion,
  type TimedQuizResultRes,
  type TimedQuizSessionRes,
  type TimedQuizSummaryRes,
} from '../../lib/api'

type Phase = 'loading' | 'intro' | 'run' | 'done'

const PHONE_KEYPAD_STORAGE_KEY = 'mathua_timed_quiz_phone_keypad'

function newSubmissionId(): string {
  const c = globalThis.crypto
  if (c && typeof c.randomUUID === 'function') return c.randomUUID()
  return `timed-quiz-${Date.now()}-${Math.random().toString(36).slice(2)}`
}

function formatRemaining(ms: number): string {
  const total = Math.max(0, Math.ceil(ms / 1000))
  const minutes = Math.floor(total / 60)
  const seconds = total % 60
  return `${minutes}:${String(seconds).padStart(2, '0')}`
}

export default function TimedQuizHost() {
  const [phase, setPhase] = useState<Phase>('loading')
  const [summary, setSummary] = useState<TimedQuizSummaryRes | null>(null)
  const [session, setSession] = useState<TimedQuizSessionRes | null>(null)
  const [result, setResult] = useState<TimedQuizResultRes | null>(null)
  const [index, setIndex] = useState(0)
  const [answer, setAnswer] = useState('')
  const [feedback, setFeedback] = useState<TimedQuizAnswerRes | null>(null)
  const [error, setError] = useState('')
  const [busy, setBusy] = useState(false)
  const [remaining, setRemaining] = useState(0)
  // Mathua's keypad is the default because it carries the symbols a phone
  // keypad cannot. The learner can hand the answer field back to the device
  // keypad from the control above it, and that choice follows them to the next
  // quiz. Above, not below: a long answer-format hint below the field can push
  // the control under the open keypad, where it cannot be tapped.
  const [phoneKeypad, setPhoneKeypad] = useState(false)
  const mathFocusRef = useRef<MathFocusHandle | null>(null)
  const phoneInputRef = useRef<HTMLInputElement | null>(null)
  const focusAfterKeypadSwitch = useRef(false)
  const offsetRef = useRef(0)
  const submissionIds = useRef<Record<string, string>>({})
  const questionShownAt = useRef(0)
  const expiredRef = useRef(false)

  const loadSummary = useCallback(async () => {
    try {
      const data = await getTimedQuiz()
      setSummary(data)
      setPhase('intro')
    } catch {
      setError('Could not load the timed quiz. Try again.')
      setPhase('intro')
    }
  }, [])

  useEffect(() => {
    void loadSummary()
  }, [loadSummary])

  useEffect(() => {
    try {
      if (window.localStorage.getItem(PHONE_KEYPAD_STORAGE_KEY) === '1') setPhoneKeypad(true)
    } catch {
      // A blocked storage read leaves the math keypad as the default.
    }
  }, [])

  useEffect(() => {
    if (!focusAfterKeypadSwitch.current) return
    focusAfterKeypadSwitch.current = false
    if (phoneKeypad) phoneInputRef.current?.focus()
    else mathFocusRef.current?.focus()
  }, [phoneKeypad])

  const toggleKeypad = useCallback(() => {
    const next = !phoneKeypad
    focusAfterKeypadSwitch.current = true
    // The math keypad is a page-wide singleton, so hiding it explicitly keeps
    // it from covering the phone-keypad input after the editor unmounts.
    if (next) window.mathVirtualKeyboard?.hide()
    try {
      window.localStorage.setItem(PHONE_KEYPAD_STORAGE_KEY, next ? '1' : '0')
    } catch {
      // The choice still applies to this run when storage is unavailable.
    }
    setPhoneKeypad(next)
  }, [phoneKeypad])

  const finish = useCallback((data: TimedQuizResultRes) => {
    setResult(data)
    setSession(null)
    setFeedback(null)
    setPhase('done')
  }, [])

  const expire = useCallback(async (runId: string) => {
    if (expiredRef.current) return
    expiredRef.current = true
    try {
      finish(await finalizeTimedQuiz(runId))
    } catch {
      setError('The quiz time ran out, but the result could not be loaded. Try again.')
      setPhase('intro')
    }
  }, [finish])

  useEffect(() => {
    if (phase !== 'run' || !session) return
    const tick = () => {
      const left = session.deadline_at - (Date.now() + offsetRef.current)
      setRemaining(left)
      if (left <= 0) void expire(session.run_id)
    }
    tick()
    const id = window.setInterval(tick, 500)
    return () => window.clearInterval(id)
  }, [phase, session, expire])

  async function start() {
    setBusy(true)
    setError('')
    expiredRef.current = false
    try {
      const data = await startTimedQuizSession()
      if ('expired' in data) {
        finish(data.result)
        return
      }
      offsetRef.current = data.server_now - Date.now()
      sessionStorage.setItem('timed_quiz_run', data.run_id)
      setSession(data)
      setIndex(0)
      setAnswer('')
      setFeedback(null)
      questionShownAt.current = Date.now()
      setPhase('run')
    } catch {
      setError('The timed quiz could not start. Try again.')
    } finally {
      setBusy(false)
    }
  }

  async function viewLatest() {
    if (!summary?.latest) return
    setBusy(true)
    setError('')
    try {
      finish(await getTimedQuizResult(summary.latest.run_id))
    } catch {
      setError('The last timed quiz result could not be loaded.')
    } finally {
      setBusy(false)
    }
  }

  async function submit(dontKnow: boolean) {
    if (!session) return
    const item = session.items[index]
    if (!item) return
    if (!dontKnow && !answer.trim()) return
    setBusy(true)
    setError('')
    const submissionId = submissionIds.current[item.item_id] ?? newSubmissionId()
    submissionIds.current[item.item_id] = submissionId
    try {
      const data = await submitTimedQuizAnswer(
        session.run_id,
        item.item_id,
        dontKnow ? '' : answer.trim(),
        Math.max(0.5, (Date.now() - questionShownAt.current) / 1000),
        submissionId,
        dontKnow,
      )
      if (data.expired && data.result) {
        finish(data.result)
        return
      }
      if (data.ungraded) {
        setError('Grading is temporarily unavailable. Try again.')
        return
      }
      setFeedback(data)
      if (data.completed && data.result) finish(data.result)
    } catch {
      setError('The answer could not be submitted. Try again.')
    } finally {
      setBusy(false)
    }
  }

  async function skip() {
    if (!session) return
    const item = session.items[index]
    if (!item) return
    setBusy(true)
    setError('')
    try {
      const data = await skipTimedQuizItem(session.run_id, item.item_id)
      if (data.expired && data.result) {
        finish(data.result)
        return
      }
      if (data.completed && data.result) {
        finish(data.result)
        return
      }
      advance()
    } catch {
      setError('The question could not be skipped. Try again.')
    } finally {
      setBusy(false)
    }
  }

  function advance() {
    if (!session) return
    const next = index + 1
    setFeedback(null)
    setAnswer('')
    questionShownAt.current = Date.now()
    if (next >= session.items.length) {
      void finalizeTimedQuiz(session.run_id).then(finish).catch(() => {
        setError('The quiz could not be finalized. Try again.')
      })
      return
    }
    setIndex(next)
  }

  if (phase === 'loading') {
    return <p className="text-mathua-secondary text-sm font-mono text-center py-16">Loading timed quiz...</p>
  }

  if (phase === 'done' && result) {
    const resultAdjustment = result.items.reduce((sum, item) => sum + (item.xp_adjustment ?? 0), 0)
    const lastAdjustmentItem = result.items.slice().reverse().find(item => (item.xp_adjustment ?? 0) !== 0 || item.xp_note)
    return (
      <div className="max-w-2xl mx-auto min-w-0">
        <SectionHeader label="Timed quiz result" title={`${result.score_percent}% correct`} />
        <div className="border border-mathua-border rounded-none p-5 mb-6">
          <div className="grid grid-cols-2 sm:grid-cols-4 gap-4 text-center">
            <Stat label="Correct" value={`${result.correct_count}/${result.question_count}`} />
            <Stat label="Unanswered" value={String(result.unanswered_count)} />
            <Stat label="XP earned" value={`+${result.xp_earned}`} />
            <Stat label="Score" value={`${result.score_percent}%`} />
          </div>
        </div>

        <XPTaskSummary
          answerXP={result.xp_earned}
          adjustment={resultAdjustment}
          note={lastAdjustmentItem?.xp_note}
          source={lastAdjustmentItem?.xp_source}
          className="mb-6 text-left"
        />

        <div className="space-y-3">
          {result.items.map((item) => (
            <div key={item.item_id} className={`border rounded-none p-4 ${item.correct ? 'border-mathua-green-faint' : 'border-mathua-red-faint'}`}>
              <div className="flex flex-wrap items-baseline justify-between gap-2 mb-2">
                <span className="font-mono text-[11px] uppercase tracking-wider text-mathua-muted">{item.concept_name}</span>
                <span className={`font-mono text-[11px] uppercase tracking-wider ${item.correct ? 'text-mathua-green' : 'text-mathua-red'}`}>
                  {item.status === 'skipped' ? 'Unanswered' : item.correct ? 'Correct' : 'Incorrect'}
                </span>
              </div>
              <KatexContent className="text-mathua-primary text-sm whitespace-pre-wrap break-words">{stemWithoutChoices(item.question)}</KatexContent>
              {item.status !== 'skipped' && (
                <p className="mt-2 font-mono text-xs text-mathua-secondary">
                  Your answer: {item.answer || 'None'} · Expected: {item.expected_answer}
                </p>
              )}
              {item.explanation && (
                <KatexContent className="mt-2 text-mathua-secondary text-sm whitespace-pre-wrap">{item.explanation}</KatexContent>
              )}
            </div>
          ))}
        </div>

        {result.remediation.length > 0 && (
          <div className="mt-6 border border-mathua-border rounded-none p-4">
            <p className="font-mono text-[11px] uppercase tracking-wider text-mathua-muted mb-3">Review these concepts</p>
            <div className="flex flex-wrap gap-2">
              {result.remediation.map((entry) => (
                <Link key={entry.concept_id} href={entry.href} className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue-faint px-3 py-1.5 font-mono text-xs rounded-none">
                  {entry.concept_name}
                </Link>
              ))}
            </div>
          </div>
        )}

        <div className="mt-6 flex flex-wrap gap-3 justify-center">
          <button type="button" onClick={() => void loadSummary()} className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue-faint rounded-none h-11 px-6 text-sm">
            Back to timed quiz
          </button>
          <Link href="/profile" className="border border-mathua-border text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue rounded-none h-11 px-6 text-sm inline-flex items-center">
            Back to profile
          </Link>
        </div>
      </div>
    )
  }

  if (phase === 'run' && session) {
    const item = session.items[index]
    if (!item) return null
    const isLast = index >= session.items.length - 1
    const answerFormat = formatForQuestion(item)
    return (
      <div className="max-w-2xl mx-auto min-w-0">
        <div className="flex flex-wrap items-center justify-between gap-3 mb-4">
          <SectionHeader label={`Timed question ${index + 1} of ${session.items.length}`} title={item.concept_name} />
          <span className={`font-mono text-sm tabular-nums ${remaining <= 60_000 ? 'text-mathua-red' : 'text-mathua-secondary'}`}>
            {formatRemaining(remaining)}
          </span>
        </div>
        <div className="border border-mathua-border rounded-none p-4 sm:p-6 mb-6">
          <div className="bg-mathua-code border border-mathua-border rounded-none p-4 sm:p-6 text-center mb-4">
            <KatexContent className="text-mathua-primary text-lg font-mono font-light whitespace-pre-wrap break-words">{stemWithoutChoices(item.question)}</KatexContent>
          </div>
          {!feedback && (
            <div className="mb-2 flex justify-end">
              <button
                type="button"
                onClick={toggleKeypad}
                className="border border-mathua-border rounded-none px-2 py-1 font-mono text-[11px] text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue"
              >
                {phoneKeypad ? 'Use math keypad' : 'Use phone keypad'}
              </button>
            </div>
          )}
          {!feedback ? (
            <form onSubmit={(event) => { event.preventDefault(); void submit(false) }} className="flex flex-col sm:flex-row gap-3">
              <label htmlFor="timed-quiz-answer" className="sr-only">Your answer</label>
              {phoneKeypad ? (
                <input
                  ref={phoneInputRef}
                  id="timed-quiz-answer"
                  value={answer}
                  onChange={(event) => setAnswer(event.target.value)}
                  inputMode={answerFormat.inputMode}
                  enterKeyHint="go"
                  autoComplete="off"
                  placeholder="Your answer..."
                  disabled={busy}
                  className="flex-1 min-w-0 border border-mathua-border rounded-none bg-mathua-bg px-3 py-2 font-mono text-sm text-mathua-primary focus:border-mathua-blue focus:outline-none disabled:opacity-50"
                />
              ) : (
                <MathAnswerInput
                  id="timed-quiz-answer"
                  value={answer}
                  onChange={setAnswer}
                  onSubmit={() => void submit(false)}
                  gradingType={item.grading_type}
                  conceptId={item.concept_id}
                  disabled={busy}
                  placeholder="Your answer..."
                  inputMode={answerFormat.inputMode}
                  focusRef={mathFocusRef}
                  className="flex-1 min-w-0"
                />
              )}
              <button type="submit" disabled={!answer.trim() || busy} className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue-faint rounded-none h-11 px-6 text-sm disabled:opacity-50">
                Check answer
              </button>
            </form>
          ) : (
            <div className={`border rounded-none p-4 ${feedback.correct ? 'border-mathua-green-faint' : 'border-mathua-red-faint'}`}>
              <p className={`font-mono text-[11px] uppercase tracking-wider mb-2 ${feedback.correct ? 'text-mathua-green' : 'text-mathua-red'}`}>
                {feedback.correct ? 'Correct' : 'Incorrect'} {feedback.xp ? `· +${feedback.xp} XP` : ''}
              </p>
              {feedback.explanation && (
                <KatexContent className="text-mathua-secondary text-sm whitespace-pre-wrap">{feedback.explanation}</KatexContent>
              )}
              <XPAdjustmentNote
                adjustment={feedback.xp_adjustment}
                note={feedback.xp_note}
                source={feedback.xp_source}
                className="mt-3"
              />
              <button type="button" onClick={advance} className="mt-4 border border-mathua-blue text-mathua-blue hover:bg-mathua-blue-faint rounded-none h-10 px-5 text-sm">
                {isLast ? 'Finish' : 'Next question'}
              </button>
            </div>
          )}
          {!feedback && (
            <div className="mt-3 flex flex-wrap items-center justify-between gap-3">
              <p className="font-mono text-[11px] text-mathua-muted">{answerFormat.hint}</p>
              <div className="flex gap-3">
                <button type="button" onClick={() => void submit(true)} disabled={busy} className="font-mono text-[11px] text-mathua-muted hover:text-mathua-red underline underline-offset-2">
                  I don&apos;t know
                </button>
                <button type="button" onClick={() => void skip()} disabled={busy} className="font-mono text-[11px] text-mathua-muted hover:text-mathua-blue underline underline-offset-2">
                  Skip
                </button>
              </div>
            </div>
          )}
        </div>
        {error && <p className="text-mathua-red text-sm font-mono text-center">{error}</p>}
      </div>
    )
  }

  return (
    <div className="max-w-2xl mx-auto">
      <SectionHeader label="Timed quiz" title="Measure fluency under a deadline" />
      <div className="border border-mathua-border rounded-none p-5">
        <p className="text-mathua-secondary text-sm mb-4">
          A timed quiz samples concepts you have already studied. It uses the same grading and XP rules as practice, but it never moves your spaced-review schedule.
        </p>
        {summary && (
          <div className="grid grid-cols-2 gap-4 mb-5">
            <Stat label="Eligible XP" value={`${summary.eligible_xp}/${summary.gate_xp}`} />
            <Stat label="Status" value={summary.active ? 'In progress' : summary.due ? 'Due' : 'Not yet due'} />
          </div>
        )}
        {summary?.active && (
          <button type="button" onClick={() => void start()} disabled={busy} className="w-full border border-mathua-blue text-mathua-blue hover:bg-mathua-blue-faint rounded-none h-11 px-6 text-sm mb-3">
            Resume timed quiz
          </button>
        )}
        {!summary?.active && summary?.due && (
          <button type="button" onClick={() => void start()} disabled={busy} className="w-full border border-mathua-blue text-mathua-blue hover:bg-mathua-blue-faint rounded-none h-11 px-6 text-sm mb-3">
            Start timed quiz
          </button>
        )}
        {!summary?.active && summary && !summary.due && (
          <p className="font-mono text-xs text-mathua-muted mb-3">
            Earn {Math.max(0, summary.gate_xp - summary.eligible_xp)} more eligible XP to unlock the next timed quiz.
          </p>
        )}
        {summary?.latest && (
          <button type="button" onClick={() => void viewLatest()} disabled={busy} className="w-full border border-mathua-border text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue rounded-none h-11 px-6 text-sm">
            View last result ({summary.latest.correct_count}/{summary.latest.question_count})
          </button>
        )}
        {error && <p className="mt-3 text-mathua-red text-sm font-mono">{error}</p>}
        <div className="mt-4 text-center">
          <Link href="/learn" className="font-mono text-xs text-mathua-muted hover:text-mathua-blue">Go to Learn</Link>
        </div>
      </div>
    </div>
  )
}

function Stat({ label, value }: { label: string; value: string }) {
  return (
    <div>
      <p className="font-mono text-[10px] uppercase tracking-wider text-mathua-muted">{label}</p>
      <p className="text-lg font-semibold tabular-nums text-mathua-primary">{value}</p>
    </div>
  )
}
