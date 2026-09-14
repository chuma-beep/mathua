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
  submitGoalAnswer,
  getGoalPlan,
  resumeGoalDiagnostic,
  skipGoalQuestion,
  type GoalPlanRes,
  type DiagnosticProgress,
} from '../../lib/api'
import { getErrorMessage } from '../../lib/api'
import { formatForGradingType, type AnswerFormat } from '../../lib/answerFormat'
import {
  MAX_SKIPS,
  toSubmitError,
  type SubmitError,
} from '../../components/SubmitErrorBlock'
import { setUserInfo, getUserInfo } from '../../lib/auth'
import { concepts as conceptsData } from '../../lib/conceptData'
import { domainOrder, type DomainInfo } from './domains'
import { DiagnosticStep, ResultsStep, WelcomeStep } from './steps'

type Step = 'welcome' | 'diagnostic' | 'results'

// MA parity: the Diagnostic doesn't have to be completed at once — the
// session id persists in sessionStorage so Back/refresh resumes it.
const DIAG_KEY = 'mathua_diag_session_onboard'

export default function OnboardPage() {
  const { mounted } = useTheme()
  const { push } = useRouter()

  const [step, setStep] = useState<Step>('welcome')
  const [loading, setLoading] = useState(false)
  const [confirming, setConfirming] = useState(false)

  const [domains, setDomains] = useState<DomainInfo[]>([])

  const sessionId = useRef('')
  const onboardInputRef = useRef<HTMLInputElement>(null)
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
  const [submitError, setSubmitError] = useState<SubmitError | null>(null)
  const [planError, setPlanError] = useState('')
  const [finished, setFinished] = useState(false)
  const [skipCount, setSkipCount] = useState(0)
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
    setDomains(result)
    try {
      setHasPaused(!!sessionStorage.getItem(DIAG_KEY))
    } catch {
      setHasPaused(false)
    }
  }, [mounted])

  function toggleDomain(name: string) {
    setDomains(prev => prev.map(d => d.name === name ? { ...d, selected: !d.selected } : d))
  }

  function selectAll() {
    setDomains(prev => prev.map(d => ({ ...d, selected: true })))
  }

  function selectedConceptIds(): string[] {
    const ids: string[] = []
    for (const d of domains) {
      if (d.selected) ids.push(...d.concepts)
    }
    return ids
  }

  async function startDiagnostic() {
    const ids = selectedConceptIds()
    if (ids.length === 0) return
    setLoading(true)
    try {
      const res = await startGoalDiagnostic(ids)
      if (res.done) {
        toast.error('No diagnostic questions are available for the selected areas. Try selecting more domains.')
        setLoading(false)
        return
      }
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
      setStep('diagnostic')
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
    if (!dontKnow && !answerInput.trim()) return
    setLoading(true)
    setSubmitError(null)
    try {
      const answer = dontKnow ? '' : answerInput.trim()
      const elapsed = Math.max(0.5, (Date.now() - (questionShownAt.current ?? Date.now())) / 1000)
      const data = await submitGoalAnswer(sessionId.current, conceptId.current, answer, elapsed, dontKnow)
      const correct = data.correct || false
      const feedback = data.feedback || (correct ? 'Correct!' : 'Not quite.')
      setAccuracy(prev => ({ correct: prev.correct + (correct ? 1 : 0), total: prev.total + 1 }))
      setLastResult({ correct, feedback })
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
    setFinished(false)
    void startDiagnostic()
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
              selectedCount={selectedConceptIds().length}
              confirming={confirming}
              onToggle={toggleDomain}
              onSelectAll={selectAll}
              onStart={() => setConfirming(true)}
              onResume={resumeDiagnostic}
              onBegin={() => { setConfirming(false); void startDiagnostic() }}
              onCancel={() => setConfirming(false)}
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
              onRetryPlan={() => { setLoading(true); void fetchPlan() }}
            />
          )}

          {step === 'results' && plan && (
            <ResultsStep plan={plan} onStartPractice={finishOnboarding} />
          )}
        </section>
      </div>
      <Footer />
      <BottomTabs />
    </>
  )
}
