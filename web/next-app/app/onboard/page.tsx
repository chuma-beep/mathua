'use client'

import { useState, useEffect, useRef } from 'react'
import { useRouter } from 'next/navigation'
import { toast } from 'sonner'
import { useTheme } from '../../hooks/useTheme'
import Header from '../../components/Header'
import BottomTabs from '../../components/BottomTabs'
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
import { setUserInfo, getUserInfo } from '../../lib/auth'
import { concepts as conceptsData } from '../../lib/conceptData'
import { scopeSize } from '../../lib/diagnosticScope'
import type { MathFocusHandle } from '@/components/math/MathInput'
import { domainOrder, type DomainInfo } from './domains'
import { DiagnosticStep, ResultsStep, WelcomeStep } from './steps'

type Step = 'welcome' | 'diagnostic' | 'results'

// MA parity: the Diagnostic doesn't have to be completed at once — the
// session id persists in sessionStorage so Back/refresh resumes it.
const DIAG_KEY = 'mathua_diag_session_onboard'

// A guest run is a demo: the server opens a session against a student row
// created from a typed name, and nothing can sign back into that row. Session
// storage keeps the fact, so a refresh mid-test still ends on a result that
// says it will not be here tomorrow.
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
  // Next question staged from the submit response — revealed by goNext(),
  // never fetched. Cleared on advance, so double-press is a no-op.
  const pendingNext = useRef<{ question: string; conceptId: string; conceptName: string; gradingType: string } | null>(null)
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
  // Concept the visible feedback belongs to — the retry call must name it,
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
    try {
      setHasPaused(!!sessionStorage.getItem(DIAG_KEY))
    } catch {
      setHasPaused(false)
    }
    try {
      setDemo(sessionStorage.getItem(DEMO_KEY) === '1')
    } catch {
      setDemo(false)
    }
  }, [mounted])

  // A paused run covers the domains it was started on. Changing the scope makes
  // that run no longer the one being asked for, so it is forgotten rather than
  // offered above the Start button, where Continue used to resume the old
  // subject after the learner had just picked a different one.
  function dropPausedRun() {
    if (!hasPaused) return
    sessionId.current = ''
    try { sessionStorage.removeItem(DIAG_KEY) } catch { /* ignore */ }
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
    try {
      sessionStorage.setItem(DIAG_KEY, res.session_id)
      setHasPaused(true)
    } catch { /* storage unavailable — session simply won't resume */ }
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
      try { sessionStorage.removeItem(DEMO_KEY) } catch { /* ignore */ }
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
      const res = await startGoalDiagnosticName(name, ids)
      if (res.done) {
        toast.error('No diagnostic questions are available for the selected areas. Try selecting more domains.')
        return
      }
      setDemo(true)
      setNeedsName(false)
      try { sessionStorage.setItem(DEMO_KEY, '1') } catch { /* ignore */ }
      beginSession(res)
    } catch {
      toast.error("Something went wrong, but we're working on it.")
    } finally {
      setLoading(false)
    }
  }

  async function resumeDiagnostic() {
    let sid = ''
    try {
      sid = sessionStorage.getItem(DIAG_KEY) || ''
    } catch {
      sid = ''
    }
    if (!sid) return
    setLoading(true)
    try {
      const data = await resumeGoalDiagnostic(sid)
      if (data.done) {
        try {
          sessionStorage.removeItem(DIAG_KEY)
        } catch { /* ignore */ }
        setHasPaused(false)
        toast.error('That diagnostic already finished — start a fresh one below.')
        setLoading(false)
        return
      }
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
      try {
        sessionStorage.removeItem(DIAG_KEY)
      } catch { /* ignore */ }
      setHasPaused(false)
      toast.error('Could not resume — that session expired. Start a fresh diagnostic.')
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
    try {
      const answer = dontKnow ? '' : answerInput.trim()
      const elapsed = Math.max(0.5, (Date.now() - (questionShownAt.current ?? Date.now())) / 1000)
      const data = await submitGoalAnswer(sessionId.current, conceptId.current, answer, elapsed, dontKnow)
      const correct = data.correct || false
      // An admitted unknown is not a wrong answer, so it does not get the wrong answer's words.
      // This step renders the feedback and no separate verdict line, so the label goes here.
      const feedback = dontKnow ? 'No answer' : data.feedback || (correct ? 'Correct!' : 'Not quite.')
      setAccuracy(prev => ({ correct: prev.correct + (correct ? 1 : 0), total: prev.total + 1 }))
      setLastResult({ correct, feedback })
      setRetryAvailable(data.done ? false : data.retry_available === true)
      setRetryConcept(data.done ? '' : conceptId.current)
      if (data.progress) setProgress(data.progress)

      if (data.done) {
        setFinished(true)
        setTimeout(() => {
          void fetchPlan()
        }, 800)
        return
      }

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
      try {
        sessionStorage.removeItem(DIAG_KEY)
      } catch { /* ignore */ }
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
        setLastResult({ correct: false, feedback: 'Skipped — no evidence recorded.' })
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

  function finishOnboarding() {
    const user = getUserInfo()
    if (user) {
      setUserInfo({ ...user, diagnostic_completed: true })
    }
    push('/profile')
  }

  if (!mounted) return <div style={{ background: 'var(--bg)', minHeight: '100vh' }} />

  return (
    <>
      <Header />
      <div className="max-w-container mx-auto px-4 sm:px-6 pb-[calc(80px+env(safe-area-inset-bottom))] lg:pb-0 overflow-x-hidden min-w-0">
        <section className="pt-8 min-w-0 overflow-hidden">
          {step !== 'welcome' && (
            <span className="flex mb-4">
              <button type="button" onClick={() => setStep('welcome')} className="text-mathua-secondary text-sm hover:text-mathua-primary">← Back</button>
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
      <BottomTabs />
    </>
  )
}
