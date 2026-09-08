'use client'

import Loading from '../../components/Loading'

import { useState, useEffect, useRef } from 'react'
import { useRouter } from 'next/navigation'
import { toast } from 'sonner'
import KatexContent from '../../components/KatexContent'
import { useTheme } from '../../hooks/useTheme'
import Header from '../../components/Header'
import BottomTabs from '../../components/BottomTabs'
import SectionHeader from '../../components/SectionHeader'
import Footer from '../../components/Footer'
import SymbolPalette from '../../components/SymbolPalette'
import ReportButton from '../../components/ReportButton'
import DiagnosticResults from '../../components/DiagnosticResults'
import ProgressBar from '../../components/ProgressBar'
import {
  startGoalDiagnostic,
  submitGoalAnswer,
  getGoalPlan,
  resumeGoalDiagnostic,
  type GoalPlanRes,
  type DiagnosticProgress,
} from '../../lib/api'
import { setUserInfo, getUserInfo } from '../../lib/auth'
import conceptsData from '../../data/concepts.json'

type Step = 'welcome' | 'diagnostic' | 'results'

interface DomainInfo {
  name: string
  count: number
  concepts: string[]
  selected: boolean
}

const domainLabels: Record<string, string> = {
  arithmetic: 'Arithmetic',
  fractions: 'Fractions',
  prealgebra: 'Pre-Algebra',
  algebra: 'Algebra',
  geometry: 'Geometry',
  trigonometry: 'Trigonometry',
  complex_numbers: 'Complex Numbers',
  precalculus: 'Precalculus',
  calculus: 'Calculus',
  linear_algebra: 'Linear Algebra',
  statistics: 'Statistics',
  discrete_math: 'Discrete Math',
  number_theory: 'Number Theory',
  differential_equations: 'Differential Equations',
  abstract_algebra: 'Abstract Algebra',
  topology: 'Topology',
}

const domainOrder = ['arithmetic', 'fractions', 'prealgebra', 'algebra', 'geometry', 'trigonometry', 'complex_numbers', 'precalculus', 'calculus', 'linear_algebra', 'statistics', 'discrete_math', 'number_theory', 'differential_equations', 'abstract_algebra', 'topology']

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
    const raw = conceptsData as any[]
    const map = new Map<string, string[]>()
    for (const c of raw) {
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
              <button onClick={() => setStep('welcome')} className="text-mathua-secondary text-sm hover:text-mathua-primary">← Back</button>
            </span>
          )}

          {/* === WELCOME === */}
          {step === 'welcome' && (
            <div className="max-w-4xl mx-auto">
              <div className="text-center mb-4 mt-6 sm:mt-8 px-2 min-w-0">
                <div className="font-mono text-[11px] uppercase text-mathua-muted mb-3">Welcome to Mathua</div>
                <h1 className="font-serif text-2xl sm:text-4xl font-medium text-mathua-primary px-2">
                  Let&apos;s find your starting point
                </h1>
                <p className="text-mathua-secondary text-sm mt-3 max-w-[500px] mx-auto px-2">
                  Select what you want to learn. We&apos;ll test your current knowledge and build a personalized plan.
                </p>
              </div>

              <div className="flex gap-3 justify-center mb-6 px-2">
                <button
                  onClick={selectAll}
                  className="bg-mathua-surface border border-mathua-border rounded-none h-10 min-h-[36px] px-6 text-sm text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue"
                >
                  Select everything
                </button>
              </div>

              <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-2 sm:gap-3 mb-8 min-w-0">
                {domains.map(d => {
                  const label = domainLabels[d.name] || d.name
                  return (
                    <button
                      key={d.name}
                      onClick={() => toggleDomain(d.name)}
                    className={`rounded-none p-3 sm:p-4 text-left transition-all text-sm min-h-[60px] min-w-0 overflow-hidden ${
                      d.selected
                        ? 'bg-mathua-blue text-white'
                        : 'bg-mathua-surface border border-mathua-border text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue'
                    }`}
                    >
                      <div className="font-medium truncate sm:whitespace-normal sm:line-clamp-2 break-words text-[13px] sm:text-sm">{label}</div>
                      <div className={`font-mono text-[10px] mt-1 truncate ${d.selected ? 'text-white/70' : 'text-mathua-muted'}`}>
                        {d.count} concepts
                      </div>
                    </button>
                  )
                })}
              </div>

              <div className="text-center px-4">
                {hasPaused && (
                  <button
                    onClick={resumeDiagnostic}
                    disabled={loading}
                    className="border border-mathua-blue bg-mathua-blue text-white hover:opacity-90 rounded-none h-12 min-h-[44px] px-6 sm:px-10 font-medium text-sm disabled:opacity-50 max-w-full mb-3"
                  >
                    {loading ? (<><Loading inline size={13} /> Loading…</>) : 'Continue diagnostic →'}
                  </button>
                )}
                <button
                  onClick={startDiagnostic}
                  disabled={selectedConceptIds().length === 0 || loading}
                  className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white rounded-none h-12 min-h-[44px] px-6 sm:px-10 font-medium text-sm disabled:opacity-50 max-w-full"
                >
                  {loading ? (<><Loading inline size={13} /> Loading…</>) : `Start diagnostic test (${selectedConceptIds().length} concepts)`}
                </button>
              </div>
            </div>
          )}

          {/* === DIAGNOSTIC === */}
          {step === 'diagnostic' && (
            <>
              <SectionHeader label={`Question ${progress ? progress.answered + 1 : questionCount}`} title={conceptName} />
              <div className="max-w-2xl mx-auto px-2 sm:px-0 min-w-0 overflow-hidden">
                <ProgressBar
                  answered={progress ? progress.answered + 1 : questionCount}
                  coverDone={progress ? progress.cover_done : 0}
                  coverSize={progress ? progress.cover_size : 0}
                />
                <div className="bg-mathua-surface border border-mathua-border rounded-none p-4 sm:p-6 mb-6 w-full max-w-full min-w-0 overflow-hidden">
                  <div className="bg-mathua-code border border-mathua-border rounded-none p-4 sm:p-6 text-center mb-4 w-full max-w-full min-w-0 overflow-hidden">
                    <div className="w-full max-w-full min-w-0 overflow-hidden">
                       <KatexContent className="text-mathua-primary text-lg font-mono font-light whitespace-pre-wrap break-words">{question}</KatexContent>
                     </div>
                   </div>
                   <div className="flex flex-col sm:flex-row gap-3 min-w-0">
                      <input
                        ref={onboardInputRef}
                        type="text"
                        value={answerInput}
                        onChange={(e) => setAnswerInput(e.target.value)}
                        onKeyDown={(e) => e.key === 'Enter' && submitAnswer()}
                        placeholder="Your answer..."
                       disabled={loading || lastResult !== null}
                        className="flex-1 min-w-0 bg-mathua-code border border-mathua-border rounded-none h-24 sm:h-12 px-4 font-mono text-base text-mathua-primary placeholder:text-mathua-muted focus:outline-none focus:border-mathua-blue"
                    />
                     <button
                       onClick={submitAnswer}
                       disabled={!answerInput.trim() || loading || lastResult !== null}
                       className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white rounded-none h-12 min-h-[36px] px-8 font-medium text-sm disabled:opacity-50 shrink-0 w-full sm:w-auto"
                     >
                       Check Answer
                     </button>
                   </div>
                    <SymbolPalette targetRef={onboardInputRef} onInsert={setAnswerInput} />
                    <div className="mt-2 flex justify-end">
                      <ReportButton
                        key={question}
                        conceptId={conceptId.current}
                        kind="question"
                        question={question}
                        source="diagnostic"
                        sessionId={sessionId.current}
                      />
                    </div>
                </div>

                {lastResult && (
                  <div className={`bg-mathua-surface border rounded-none p-4 mb-4 text-center ${lastResult.correct ? 'border-mathua-green' : 'border-mathua-red'}`}>
                    <KatexContent className={lastResult.correct ? 'text-mathua-green' : 'text-mathua-red'}>{lastResult.feedback}</KatexContent>
                  </div>
                )}

                <div className="text-center text-mathua-muted text-xs font-mono">
                  {accuracy.total > 0 && `${accuracy.correct}/${accuracy.total} correct`}
                </div>
              </div>
            </>
          )}

          {/* === RESULTS === */}
          {step === 'results' && plan && (
            <div className="max-w-3xl mx-auto min-w-0 overflow-hidden px-2 sm:px-0">
              <SectionHeader label="Your results" title="Here's what we found" />
              <div className="mt-6 min-w-0 overflow-hidden">
                <DiagnosticResults plan={plan} onStartPractice={finishOnboarding} />
              </div>
            </div>
          )}
        </section>
      </div>
      <Footer />
      <BottomTabs />
    </>
  )
}