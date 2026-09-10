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
  type GoalPlanRes,
  type DiagnosticProgress,
} from '../../lib/api'
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

  const [domains, setDomains] = useState<DomainInfo[]>([])

  const sessionId = useRef('')
  const onboardInputRef = useRef<HTMLInputElement>(null)
  const [question, setQuestion] = useState('')
  const conceptId = useRef('')
  const questionShownAt = useRef<number | null>(null)
  const [conceptName, setConceptName] = useState('')
  const [questionCount, setQuestionCount] = useState(0)
  const [progress, setProgress] = useState<DiagnosticProgress | null>(null)
  const [answerInput, setAnswerInput] = useState('')
  const [lastResult, setLastResult] = useState<{ correct: boolean; feedback: string } | null>(null)
  const [accuracy, setAccuracy] = useState<{ correct: number; total: number }>({ correct: 0, total: 0 })
  const [hasPaused, setHasPaused] = useState(false)

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
      setQuestionCount(1)
      setProgress(res.progress ?? null)
      setAccuracy({ correct: 0, total: 0 })
      setLastResult(null)
      setAnswerInput('')
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
      if (data.progress) {
        setProgress(data.progress)
        setQuestionCount(data.progress.answered + 1)
      }
      setLastResult(null)
      setAnswerInput('')
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

  async function submitAnswer() {
    if (!answerInput.trim()) return
    setLoading(true)
    try {
      const answer = answerInput.trim()
      const elapsed = Math.max(0.5, (Date.now() - (questionShownAt.current ?? Date.now())) / 1000)
      const data = await submitGoalAnswer(sessionId.current, conceptId.current, answer, elapsed)
      const correct = data.correct || false
      const feedback = data.feedback || (correct ? 'Correct!' : 'Not quite.')
      setAccuracy(prev => ({ correct: prev.correct + (correct ? 1 : 0), total: prev.total + 1 }))
      setLastResult({ correct, feedback })
      if (data.progress) setProgress(data.progress)

      if (data.done) {
        setTimeout(async () => {
          try {
            const planRes = await getGoalPlan(sessionId.current)
            setPlan(planRes)
            setStep('results')
            try {
              sessionStorage.removeItem(DIAG_KEY)
            } catch { /* ignore */ }
            setHasPaused(false)
          } catch {
            alert('Could not generate plan.')
          }
          setLoading(false)
        }, 800)
        return
      }

      setTimeout(() => {
        setQuestion(data.question || '')
        conceptId.current = data.concept_id || ''
        questionShownAt.current = Date.now()
        setConceptName(data.concept_name || '')
        setQuestionCount(prev => prev + 1)
        setLastResult(null)
        setAnswerInput('')
        setLoading(false)
      }, 1200)
    } catch {
      alert('Failed to submit answer.')
      setLoading(false)
    }
  }

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
              onToggle={toggleDomain}
              onSelectAll={selectAll}
              onStart={startDiagnostic}
              onResume={resumeDiagnostic}
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
              onSubmit={submitAnswer}
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
