'use client'

import { useState, useEffect, useRef } from 'react'
import { useRouter } from 'next/navigation'
import { toast } from 'sonner'
import { useTheme } from '../../hooks/useTheme'
import AppShell from '../../components/app-shell'
import Footer from '../../components/Footer'
import {
  startGoalDiagnostic,
  startGoalDiagnosticName,
  submitGoalAnswer,
  getGoalPlan,
  resumeGoalDiagnostic,
  skipGoalQuestion,
  retryGoalQuestion,
  type GoalDiagStartRes,
  type GoalPlanRes,
  type DiagnosticProgress,
} from '../../lib/api'
import { getErrorMessage, getErrorStatus } from '../../lib/api'
import { formatForGradingType, type AnswerFormat } from '../../lib/answerFormat'
import {
  MAX_SKIPS,
  toSubmitError,
  type SubmitError,
} from '../../components/SubmitErrorBlock'
import { ensureGuestId, setUserInfo, getUserInfo } from '../../lib/auth'
import { concepts as conceptsData } from '../../lib/conceptData'
import { scopeSize } from '../../lib/diagnosticScope'
import type { MathFocusHandle } from '@/components/math/MathInput'
import { domainOrder, type DomainInfo } from './domains'
import { DiagnosticStep, ResultsStep, WelcomeStep } from './steps'

type Step = 'welcome' | 'diagnostic' | 'results'

// MA parity: the Diagnostic doesn't have to be completed at once. The session
// id is kept in sessionStorage for Back/refresh and in localStorage for a later
// visit, so leaving the site does not throw the run away.
const DIAG_KEY = 'mathua_diag_session_onboard'
const STEP_KEY = 'mathua_onboard_step'
const DIAG_AT_KEY = 'mathua_diag_session_at'
// The server keeps the run for the same seven days (diagnosticRetention in
// internal/server/server.go). Offering a resume the server has already swept
// ends in "the server restarted and forgot this run", so the two must agree.
const RESUME_TTL_MS = 7 * 24 * 60 * 60 * 1000

// The session id used to live only in sessionStorage, which a phone discards
// when the tab or PWA closes. The diagnostic is long enough that leaving for a
// while is normal, so keep a dated copy in localStorage and fall back to the
// per-tab value for older browsers or private mode.
function readStoredSession(): { sid: string; step: string } {
  let sid = ''
  let step = ''
  try {
    const at = Number(localStorage.getItem(DIAG_AT_KEY) || '0')
    if (at && Date.now() - at < RESUME_TTL_MS) {
      sid = localStorage.getItem(DIAG_KEY) || ''
      step = localStorage.getItem(STEP_KEY) || ''
    } else if (at) {
      localStorage.removeItem(DIAG_KEY)
      localStorage.removeItem(STEP_KEY)
      localStorage.removeItem(DIAG_AT_KEY)
    }
  } catch { /* storage unavailable */ }
  if (!sid) {
    try {
      sid = sessionStorage.getItem(DIAG_KEY) || ''
      step = sessionStorage.getItem(STEP_KEY) || ''
    } catch { /* storage unavailable */ }
  }
  return { sid, step }
}

function storeSession(sid: string, step: string): void {
  try {
    sessionStorage.setItem(DIAG_KEY, sid)
    sessionStorage.setItem(STEP_KEY, step)
  } catch { /* storage unavailable */ }
  try {
    localStorage.setItem(DIAG_KEY, sid)
    localStorage.setItem(STEP_KEY, step)
    localStorage.setItem(DIAG_AT_KEY, String(Date.now()))
  } catch { /* storage unavailable */ }
}

function storeStep(step: string): void {
  try { sessionStorage.setItem(STEP_KEY, step) } catch { /* storage unavailable */ }
  try {
    localStorage.setItem(STEP_KEY, step)
    localStorage.setItem(DIAG_AT_KEY, String(Date.now()))
  } catch { /* storage unavailable */ }
}

function clearStoredSession(): void {
  try {
    sessionStorage.removeItem(DIAG_KEY)
    sessionStorage.removeItem(STEP_KEY)
  } catch { /* storage unavailable */ }
  try {
    localStorage.removeItem(DIAG_KEY)
    localStorage.removeItem(STEP_KEY)
    localStorage.removeItem(DIAG_AT_KEY)
  } catch { /* storage unavailable */ }
}

function clearStoredStep(): void {
  try { sessionStorage.removeItem(STEP_KEY) } catch { /* storage unavailable */ }
  try { localStorage.removeItem(STEP_KEY) } catch { /* storage unavailable */ }
}

function setStoredDemo(on: boolean): void {
  try {
    if (on) {
      sessionStorage.setItem(DEMO_KEY, '1')
      localStorage.setItem(DEMO_KEY, '1')
    } else {
      sessionStorage.removeItem(DEMO_KEY)
      localStorage.removeItem(DEMO_KEY)
    }
  } catch { /* storage unavailable */ }
}

function isStoredDemo(): boolean {
  try {
    return localStorage.getItem(DEMO_KEY) === '1' || sessionStorage.getItem(DEMO_KEY) === '1'
  } catch {
    return false
  }
}

// A guest run is saved under the browser's persistent guest id. Signup adopts
// that same row, so the diagnostic and its attempts survive account creation.
// The demo flag still means "not yet an account", not "will be discarded".
const DEMO_KEY = 'mathua_onboard_demo'

export default function OnboardPage() {
  const { mounted } = useTheme()
  const { push } = useRouter()

  const [step, setStep] = useState<Step>('welcome')
  const [loading, setLoading] = useState(false)
  const [confirming, setConfirming] = useState(false)

  const [domains, setDomains] = useState<DomainInfo[]>([])

  const sessionId = useRef('')
  const onboardInputRef = useRef<MathFocusHandle | null>(null)
  const [question, setQuestion] = useState('')
  const conceptId = useRef('')
  // Next question staged from the submit response, revealed by goNext(),
  // never fetched. Cleared on advance, so double-press is a no-op.
  const pendingNext = useRef<{ question: string; conceptId: string; conceptName: string; gradingType: string } | null>(null)
  // Whether the last attempt was an admitted unknown. A retry that forgets this
  // posts an empty graded answer, so the learner is marked wrong for a question
  // they said they could not answer.
  const lastAttemptDontKnow = useRef(false)
  const questionShownAt = useRef<number | null>(null)
  const [conceptName, setConceptName] = useState('')
  const [questionCount, setQuestionCount] = useState(0)
  const [progress, setProgress] = useState<DiagnosticProgress | null>(null)
  const [answerInput, setAnswerInput] = useState('')
  const [lastResult, setLastResult] = useState<{ correct: boolean; feedback: string } | null>(null)
  const [accuracy, setAccuracy] = useState<{ correct: number; total: number }>({ correct: 0, total: 0 })
  const [hasPaused, setHasPaused] = useState(false)
  // A guest has no token until something mints one, and /onboard was one of the
  // entry points that never did: pressing Begin sent an unauthenticated start,
  // the server answered 401, and the catch showed "Something went wrong" for a
  // screen that was working as designed. The 401 is now the signal to offer a
  // named demo run instead of an error.
  const [needsName, setNeedsName] = useState(false)
  const [demoName, setDemoName] = useState('')
  const [demo, setDemo] = useState(false)
  const [submitError, setSubmitError] = useState<SubmitError | null>(null)
  const [planError, setPlanError] = useState('')
  const [finished, setFinished] = useState(false)
  const [skipCount, setSkipCount] = useState(0)
  const [retryAvailable, setRetryAvailable] = useState(false)
  // Concept the visible feedback belongs to, the retry call must name it,
  // since conceptId.current already advanced to the staged next question.
  const [retryConcept, setRetryConcept] = useState('')
  const [answerFormat, setAnswerFormat] = useState<AnswerFormat>(() => formatForGradingType())

  const [plan, setPlan] = useState<GoalPlanRes | null>(null)

  useEffect(() => {
    if (!mounted) return
    const map = new Map<string, string[]>()
    for (const c of conceptsData) {
      const list = map.get(c.domain) || []
      list.push(c.id)
      map.set(c.domain, list)
    }
    const result: DomainInfo[] = []
    map.forEach((concepts, name) => {
      result.push({ name, count: concepts.length, concepts, selected: false })
    })
    result.sort((a, b) => domainOrder.indexOf(a.name) - domainOrder.indexOf(b.name))
    // The selection is this visit's, not the browser's. Persisting it meant a
    // learner who came back to pick algebra found machine learning still ticked
    // and started a run covering both. A fresh entry starts with nothing chosen;
    // the pause is still offered separately, and only for a run that has not had
    // its scope changed underneath it.
    setDomains(result)
    const stored = readStoredSession()
    const saved = stored.sid
    const activeStep = stored.step
    setHasPaused(!!saved)
    setDemo(isStoredDemo())
    if (activeStep === 'diagnostic' && saved) {
      void resumeDiagnostic(saved)
    }
  }, [mounted])

  // A paused run covers the domains it was started on. Changing the scope makes
  // that run no longer the one being asked for, so it is forgotten rather than
  // offered above the Start button, where Continue used to resume the old
  // subject after the learner had just picked a different one.
  function dropPausedRun() {
    if (!hasPaused) return
    sessionId.current = ''
    clearStoredSession()
    setHasPaused(false)
  }

  function toggleDomain(name: string) {
    dropPausedRun()
    setDomains(prev => prev.map(d => d.name === name ? { ...d, selected: !d.selected } : d))
  }

  // One control for both directions: when every domain is already chosen the
  // same button clears them, because a full grid gives the learner nothing left
  // to add and a second trip through 21 cards to undo it is busywork.
  function selectAll() {
    dropPausedRun()
    setDomains(prev => {
      const allSelected = prev.length > 0 && prev.every(d => d.selected)
      return prev.map(d => ({ ...d, selected: !allSelected }))
    })
  }

  function selectedConceptIds(): string[] {
    const ids: string[] = []
    for (const d of domains) {
      if (d.selected) ids.push(...d.concepts)
    }
    return ids
  }

  // What the assessment can actually reach. Not the number of ticked concepts:
  // the server walks the same closure before it builds anything, and the flat sum
  // both overstates the questions that will be asked and omits the prerequisites
  // that come along for free.
  function scopeCount(): number {
    return scopeSize(selectedConceptIds())
  }

  // Everything both a signed-in start and a demo start must do once the server
  // has handed back a session. Kept in one place because the two paths differ
  // only in which call opened it, and a field added to one but not the other is
  // a bug that would only appear for guests.
  function beginSession(res: GoalDiagStartRes) {
    sessionId.current = res.session_id
    storeSession(res.session_id, 'diagnostic')
    setHasPaused(true)
    setQuestion(res.question || '')
    conceptId.current = res.concept_id || ''
    questionShownAt.current = Date.now()
    setConceptName(res.concept_name || '')
    setAnswerFormat(formatForGradingType(res.grading_type))
    setQuestionCount(1)
    setProgress(res.progress ?? null)
    setAccuracy({ correct: 0, total: 0 })
    setLastResult(null)
    setAnswerInput('')
    pendingNext.current = null
    setSubmitError(null)
    setPlanError('')
    setFinished(false)
    setSkipCount(0)
    setRetryAvailable(false)
    setRetryConcept('')
    setStep('diagnostic')
  }

  async function startDiagnostic() {
    const ids = selectedConceptIds()
    if (ids.length === 0) return
    setLoading(true)
    try {
      const res = await startGoalDiagnostic(ids)
      if (res.done) {
        toast.error('No diagnostic questions are available for the selected areas. Try selecting more domains.')
        return
      }
      setDemo(false)
      setNeedsName(false)
      setStoredDemo(false)
      beginSession(res)
    } catch (err) {
      if (getErrorStatus(err) === 401) {
        setNeedsName(true)
        return
      }
      toast.error("Something went wrong, but we're working on it.")
    } finally {
      setLoading(false)
    }
  }

  async function startDemo() {
    const ids = selectedConceptIds()
    const name = demoName.trim()
    if (ids.length === 0 || name === '') return
    setLoading(true)
    try {
      const res = await startGoalDiagnosticName(name, ids, ensureGuestId() ?? undefined)
      if (res.done) {
        toast.error('No diagnostic questions are available for the selected areas. Try selecting more domains.')
        return
      }
      setDemo(true)
      setNeedsName(false)
      setStoredDemo(true)
      beginSession(res)
    } catch {
      toast.error("Something went wrong, but we're working on it.")
    } finally {
      setLoading(false)
    }
  }

  async function resumeDiagnostic(savedSid?: string) {
    let sid = savedSid || ''
    if (!sid) sid = readStoredSession().sid
    if (!sid) return
    setLoading(true)
    try {
      const data = await resumeGoalDiagnostic(sid)
      if (data.done) {
        clearStoredSession()
        setHasPaused(false)
        setStep('welcome')
        toast.error('That diagnostic already finished, start a fresh one below.')
        setLoading(false)
        return
      }
      storeStep('diagnostic')
      sessionId.current = sid
      setQuestion(data.question || '')
      conceptId.current = data.concept_id || ''
      questionShownAt.current = Date.now()
      setConceptName(data.concept_name || '')
      setAnswerFormat(formatForGradingType(data.grading_type))
      if (data.progress) {
        setProgress(data.progress)
        setQuestionCount(data.progress.answered + 1)
      }
      setLastResult(null)
      setAnswerInput('')
      pendingNext.current = null
      setSubmitError(null)
      setPlanError('')
      setFinished(false)
      setSkipCount(0)
      setRetryAvailable(false)
      setRetryConcept('')
      setStep('diagnostic')
    } catch {
      clearStoredSession()
      setHasPaused(false)
      setStep('welcome')
      toast.error('Could not resume, that session expired. Start a fresh diagnostic.')
    } finally {
      setLoading(false)
    }
  }

  async function submitAnswer(dontKnow = false) {
    // One answer per question. The field's Enter listener can call this while the
    // previous submit is in flight or its feedback is showing.
    if (loading || lastResult) return
    if (!dontKnow && !answerInput.trim()) return
    setLoading(true)
    setSubmitError(null)
    lastAttemptDontKnow.current = dontKnow
    try {
      const answer = dontKnow ? '' : answerInput.trim()
      const elapsed = Math.max(0.5, (Date.now() - (questionShownAt.current ?? Date.now())) / 1000)
      const data = await submitGoalAnswer(sessionId.current, conceptId.current, answer, elapsed, dontKnow)
      const correct = data.correct || false
      // An admitted unknown is not a wrong answer, so it does not get the wrong answer's words.
      // This step renders the feedback and no separate verdict line, so the label goes here.
      if (data.progress) setProgress(data.progress)

      if (data.done) {
        setFinished(true)
        // The plan request can take its whole timeout, and leaving the loading
        // flag set makes every control look dead once the last answer is in.
        setLoading(false)
        setTimeout(() => {
          void fetchPlan()
        }, 800)
        return
      }

      // The server promises a question whenever done is false, so an absent one
      // is a broken response. Throwing reaches the error block, which offers a
      // retry, instead of staging '' and leaving Next looking inert.
      if (!data.question) {
        throw new Error('The diagnostic returned no next question. Try again.')
      }

      if (dontKnow) {
        // Cleanly advance to the next probe without holding on a dead-end feedback card
        setAccuracy(prev => ({ correct: prev.correct, total: prev.total + 1 }))
        applyQuestion(data.question, data.concept_id || '', data.concept_name || '', data.grading_type || '')
        setLoading(false)
        return
      }

      const feedback = data.feedback || (correct ? 'Correct!' : 'Not quite.')
      setAccuracy(prev => ({ correct: prev.correct + (correct ? 1 : 0), total: prev.total + 1 }))
      setLastResult({ correct, feedback })
      setRetryAvailable(data.retry_available === true)
      setRetryConcept(conceptId.current)

      // Manual advance: stage the prefetched next question, reveal on Next.
      // The staged question stays hidden, so answer timing starts at reveal.
      pendingNext.current = {
        question: data.question || '',
        conceptId: data.concept_id || '',
        conceptName: data.concept_name || '',
        gradingType: data.grading_type || '',
      }
      setLoading(false)
    } catch (e) {
      // A 404 means no session under this id any more, which is what happens when
      // the final answer committed but its response was lost: the server deleted
      // the finished session and the retry has nothing to answer. Ask for the
      // plan before calling a completed run a failure.
      if ((e as { status?: number } | null)?.status === 404) {
        try {
          const planRes = await getGoalPlan(sessionId.current)
          if (!planRes) throw new Error('no plan to recover')
          setPlan(planRes)
          setFinished(true)
          setStep('results')
          try {
            clearStoredSession()
          } catch { /* ignore */ }
          setHasPaused(false)
          setLoading(false)
          return
        } catch {
          // No plan either, so the run really is gone and the error block below
          // offers the restart, which is the only honest option left.
        }
      }
      setSubmitError(toSubmitError(e, 'Failed to submit answer.'))
      setLoading(false)
    }
  }

  async function fetchPlan() {
    setPlanError('')
    try {
      const planRes = await getGoalPlan(sessionId.current)
      setPlan(planRes)
      setStep('results')
      clearStoredSession()
      setHasPaused(false)
    } catch (e) {
      setPlanError(getErrorMessage(e) || 'Could not generate plan.')
    } finally {
      setLoading(false)
    }
  }

  function applyQuestion(question: string, cid: string, name: string, gradingType: string) {
    setQuestion(question)
    conceptId.current = cid
    questionShownAt.current = Date.now()
    setConceptName(name)
    setAnswerFormat(formatForGradingType(gradingType))
    setQuestionCount(prev => prev + 1)
    setLastResult(null)
    setAnswerInput('')
    setSubmitError(null)
    setRetryAvailable(false)
    setRetryConcept('')
  }

  async function retryQuestion() {
    if (loading || !retryConcept) return
    setLoading(true)
    setSubmitError(null)
    try {
      const data = await retryGoalQuestion(sessionId.current, retryConcept)
      setRetryAvailable(false)
      setRetryConcept('')
      pendingNext.current = null
      setQuestion(data.question || '')
      conceptId.current = data.concept_id || ''
      questionShownAt.current = Date.now()
      setConceptName(data.concept_name || '')
      setAnswerFormat(formatForGradingType(data.grading_type || ''))
      setLastResult(null)
      setAnswerInput('')
    } catch (e) {
      setSubmitError(toSubmitError(e, 'Failed to retry question.'))
    } finally {
      setLoading(false)
    }
  }

  async function skipAnswer() {
    if (skipCount >= MAX_SKIPS || loading) return
    setLoading(true)
    setSubmitError(null)
    try {
      const data = await skipGoalQuestion(sessionId.current)
      setSkipCount(c => c + 1)
      if (data.done) {
        setFinished(true)
        setLastResult({ correct: false, feedback: 'Skipped, no evidence recorded.' })
        setTimeout(() => {
          void fetchPlan()
        }, 800)
        return
      }
      applyQuestion(data.question || '', data.concept_id || '', data.concept_name || '', data.grading_type || '')
    } catch (e) {
      setSubmitError(toSubmitError(e, 'Failed to skip question.'))
    } finally {
      setLoading(false)
    }
  }

  function restartDiagnostic() {
    setSubmitError(null)
    setPlanError('')
    setSkipCount(0)
    setRetryAvailable(false)
    setRetryConcept('')
    setFinished(false)
    // A demo restart must not fall back to the signed-in call: it would 401,
    // and the learner would be asked for their name again after already giving
    // it once for this run.
    if (demo) void startDemo()
    else void startDiagnostic()
  }

  function goNext() {
    const staged = pendingNext.current
    if (!staged) return
    pendingNext.current = null
    applyQuestion(staged.question, staged.conceptId, staged.conceptName, staged.gradingType)
  }

  // Keyboard flow: put the cursor back in the answer box whenever a fresh
  // question is revealed (start, resume, or manual Next).
  useEffect(() => {
    if (step === 'diagnostic' && !lastResult && question) onboardInputRef.current?.focus()
  }, [step, lastResult, question])

  // A step replaces a long welcome screen with a short question screen. On a
  // phone the old scroll offset survives that change and lands on the footer,
  // so reset the viewport whenever the step changes. The focus effect above
  // runs first; this then wins the scroll position.
  useEffect(() => {
    if (mounted) window.scrollTo({ top: 0, left: 0, behavior: 'auto' })
  }, [mounted, step])

  function finishOnboarding() {
    const user = getUserInfo()
    if (user) {
      setUserInfo({ ...user, diagnostic_completed: true })
    }
    push('/profile')
  }

  if (!mounted) return <div style={{ background: 'var(--bg)', minHeight: '100vh' }} />

  return (
    <AppShell>
      <div className="max-w-container mx-auto px-4 sm:px-6 pb-[calc(80px+env(safe-area-inset-bottom))] lg:pb-0 overflow-x-hidden min-w-0">
        <section className="pt-8 min-w-0 overflow-hidden">
          {step !== 'welcome' && (
            <span className="flex mb-4">
              <button
                type="button"
                onClick={() => {
                  clearStoredStep()
                  setStep('welcome')
                }}
                className="text-mathua-secondary text-sm hover:text-mathua-primary"
              >
                ← Back
              </button>
            </span>
          )}

          {step === 'welcome' && (
            <WelcomeStep
              domains={domains}
              loading={loading}
              hasPaused={hasPaused}
              selectedCount={scopeCount()}
              confirming={confirming}
              onToggle={toggleDomain}
              onSelectAll={selectAll}
              onStart={() => setConfirming(true)}
              onResume={resumeDiagnostic}
              onBegin={() => { setConfirming(false); void startDiagnostic() }}
              onCancel={() => setConfirming(false)}
              needsName={needsName}
              demoName={demoName}
              onDemoName={setDemoName}
              onStartDemo={() => { void startDemo() }}
              onCancelDemo={() => setNeedsName(false)}
            />
          )}

          {step === 'diagnostic' && (
            <DiagnosticStep
              question={question}
              conceptName={conceptName}
              progress={progress}
              questionCount={questionCount}
              answerInput={answerInput}
              lastResult={lastResult}
              accuracy={accuracy}
              loading={loading}
              inputRef={onboardInputRef}
              conceptId={conceptId.current}
              sessionId={sessionId.current}
              onInputChange={setAnswerInput}
              onSubmit={() => { void submitAnswer(false) }}
              onDontKnow={() => { void submitAnswer(true) }}
              onRetrySubmit={() => { void submitAnswer(lastAttemptDontKnow.current) }}
              onNext={goNext}
              done={finished}
              answerFormat={answerFormat}
              submitError={submitError}
              planError={planError}
              onSkip={skipAnswer}
              skipsLeft={MAX_SKIPS - skipCount}
              onRestart={restartDiagnostic}
              retryAvailable={retryAvailable}
              onRetryQuestion={() => { void retryQuestion() }}
              onRetryPlan={() => { setLoading(true); void fetchPlan() }}
            />
          )}

          {step === 'results' && plan && (
            <ResultsStep plan={plan} onStartPractice={finishOnboarding} demo={demo} />
          )}
        </section>
      </div>
      <Footer />
    </AppShell>
  )
}
