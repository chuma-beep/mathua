'use client'

import Loading from '../../components/Loading'
import { useState, useEffect, useRef, Suspense } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import { toast } from 'sonner'
import KatexContent from '../../components/KatexContent'
import { useTheme } from '../../hooks/useTheme'
import Header from '../../components/Header'
import SectionHeader from '../../components/SectionHeader'
import ProgressSummary from '../../components/ProgressSummary'
import Footer from '../../components/Footer'
import SymbolPalette from '../../components/SymbolPalette'
import DiagnosticResults from '../../components/DiagnosticResults'
import {
  startGoalDiagnostic,
  submitGoalAnswer,
  getGoalPlan,
  getScores,
  startQuizSession,
  submitQuizAnswer,
  type GoalPlanRes,
  type Scores,
} from '../../lib/api'
import { isLoggedIn, getUserInfo } from '../../lib/auth'
import { useSearchParams } from 'next/navigation'
import conceptsData from '../../data/concepts.json'

type Step = 'select' | 'diagnostic' | 'results' | 'quiz' | 'quiz_done'

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

function GoalsContent() {
  const { mounted } = useTheme()
  const { push } = useRouter()
  const searchParams = useSearchParams()

  const [step, setStep] = useState<Step>('select')
  const [scores, setScores] = useState<Scores | null>(null)
  const [loading, setLoading] = useState(false)

  // Step 1: goal selection
  const [domains, setDomains] = useState<DomainInfo[]>([])
  const [customConcepts] = useState<string[]>([])

  // Step 2: diagnostic
  const sessionId = useRef('')
  const tokenRef = useRef('')
  const goalsInputRef = useRef<HTMLInputElement>(null)
  const [question, setQuestion] = useState('')
  const conceptId = useRef('')
  const [conceptName, setConceptName] = useState('')
  const [questionCount, setQuestionCount] = useState(0)
  const [estimatedTotal, setEstimatedTotal] = useState(0)
  const [answerInput, setAnswerInput] = useState('')
  const [lastResult, setLastResult] = useState<{ correct: boolean; feedback: string } | null>(null)
  const [accuracy, setAccuracy] = useState<{ correct: number; total: number }>({ correct: 0, total: 0 })

  // Step 3: results
  const [plan, setPlan] = useState<GoalPlanRes | null>(null)
  const [weakByDomain, setWeakByDomain] = useState<Record<string, { id: string; label: string }[]> | null>(null)

  // Quiz (actionable, 150 XP gate, guest unlimited, own grading path)
  const quizSessionId = useRef('')
  const quizInputRef = useRef<HTMLInputElement>(null)
  const quizShownAt = useRef<number | null>(null)
  const [quizQuestion, setQuizQuestion] = useState('')
  const quizConceptId = useRef('')
  const [quizConceptName, setQuizConceptName] = useState('')
  const [quizCount, setQuizCount] = useState(0)
  const [quizAnswerInput, setQuizAnswerInput] = useState('')
  const [quizLastResult, setQuizLastResult] = useState<{ correct: boolean; feedback: string; xp?: number } | null>(null)
  const [quizAccuracy, setQuizAccuracy] = useState({ correct: 0, total: 0 })
  const [quizDone, setQuizDone] = useState(false)

  // Load scores and domains on mount — auto-start quiz if ?quiz=1
  useEffect(() => {
    if (!mounted) return
    const token = localStorage.getItem('mathua_token')
    tokenRef.current = token
    const loggedIn = isLoggedIn()
    if (loggedIn) {
      const user = getUserInfo()
      if (user) {
        getScores(user.student_id).then(setScores).catch(e => console.error('scores fetch failed:', e))
      }
      fetch('/api/weaknesses', {
        headers: { Authorization: `Bearer ${token || ''}` },
      }).then(r => r.json()).then(d => {
        if (d.by_domain) setWeakByDomain(d.by_domain)
      }).catch(e => console.error('weaknesses fetch failed:', e))
    }
    buildDomains()
    if (searchParams.get('quiz') === '1') {
      // Reuse: quiz host — auto-start actionable quiz (150 XP, 80% own grading, guest allowed)
      setTimeout(() => { startQuiz() }, 300)
    }
  }, [mounted, searchParams])

  function buildDomains() {
    const raw = conceptsData as any[]
    const map = new Map<string, string[]>()
    for (const c of raw) {
      const list = map.get(c.domain) || []
      list.push(c.id)
      map.set(c.domain, list)
    }
    const result: DomainInfo[] = []
    map.forEach((concepts, name) => {
      result.push({
        name,
        count: concepts.length,
        concepts,
        selected: false,
      })
    })
    result.sort((a, b) => {
      const order = ['arithmetic', 'fractions', 'prealgebra', 'algebra', 'geometry', 'trigonometry', 'complex_numbers', 'precalculus', 'calculus', 'linear_algebra', 'statistics', 'discrete_math', 'number_theory', 'differential_equations', 'abstract_algebra', 'topology']
      return order.indexOf(a.name) - order.indexOf(b.name)
    })
    setDomains(result)
  }

  function toggleDomain(name: string) {
    setDomains(prev => prev.map(d => d.name === name ? { ...d, selected: !d.selected } : d))
  }

  function selectedConceptIds(): string[] {
    const ids: string[] = []
    for (const d of domains) {
      if (d.selected) {
        ids.push(...d.concepts)
      }
    }
    ids.push(...customConcepts)
    return ids
  }

  async function startDiagnostic() {
    const ids = selectedConceptIds()
    if (ids.length === 0) return
    setLoading(true)
    try {
      const res = await startGoalDiagnostic(ids)
      if (res.done) {
        setPlan({
          readiness: 1,
          total_tested: 0,
          correct_count: 0,
          weak_areas: {},
          strong_areas: {},
        })
        setStep('results')
        return
      }
      sessionId.current = res.session_id
      setQuestion(res.question || '')
      conceptId.current = res.concept_id || ''
      setConceptName(res.concept_name || '')
      setQuestionCount(1)
      // Estimate total: ~10 + log2(selected concepts)
      const est = Math.min(10 + Math.ceil(Math.log2(ids.length) * 5), 50)
      setEstimatedTotal(est)
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

  async function submitAnswer() {
    if (!answerInput.trim()) return
    setLoading(true)
    try {
      const answer = answerInput.trim()
      const elapsed = 5.0
      const data = await submitGoalAnswer(sessionId.current, conceptId.current, answer, elapsed)
      const correct = data.correct || false
      const feedback = data.feedback || (correct ? 'Correct!' : 'Not quite.')
      setAccuracy(prev => ({
        correct: prev.correct + (correct ? 1 : 0),
        total: prev.total + 1,
      }))
      setLastResult({ correct, feedback })

      if (data.done) {
        setTimeout(async () => {
          try {
            const planRes = await getGoalPlan(sessionId.current)
            setPlan(planRes)
            setStep('results')
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

  async function startQuiz() {
    setLoading(true)
    try {
      const res = await startQuizSession()
      if (res.done) {
        setQuizDone(true)
        setStep('quiz_done')
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
      setQuizDone(false)
      setStep('quiz')
    } catch {
      toast.error("Quiz failed to start — try again.")
    } finally {
      setLoading(false)
    }
  }

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

      if (data.done) {
        setTimeout(() => {
          setQuizDone(true)
          setStep('quiz_done')
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
        setLoading(false)
      }, 1200)
    } catch {
      toast.error("Failed to submit quiz answer.")
      setLoading(false)
    }
  }

  function startPractice() {
    push('/profile')
  }

  if (!mounted) return <div style={{ background: 'var(--bg)', minHeight: '100vh' }} />

  return (
    <>
      <Header />
      <div className="max-w-container mx-auto px-4 sm:px-6 pb-[calc(80px+env(safe-area-inset-bottom))] lg:pb-0 overflow-x-hidden min-w-0">
        <section className="pt-8 min-w-0 overflow-hidden">
          <span className="flex mb-4">
            <Link href="/" className="text-mathua-secondary text-sm hover:text-mathua-primary">Back</Link>
          </span>

          {/* === STEP 1: Goal Selection === */}
          {step === 'select' && (
            <>
              {scores && <div className="min-w-0 overflow-hidden"><ProgressSummary scores={scores} weakByDomain={weakByDomain || undefined} /></div>}
              <SectionHeader label="Step 1" title="What do you want to learn?" />
              <p className="text-mathua-secondary text-sm text-center max-w-[600px] mx-auto mt-2 mb-8 px-2">
                Pick one or more topics. The system will test your prerequisite knowledge and build a personalized study plan.
              </p>

              <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-2 sm:gap-3 max-w-4xl mx-auto mb-8 min-w-0">
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
                <button
                  onClick={startDiagnostic}
                  disabled={selectedConceptIds().length === 0 || loading}
                  className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white rounded-none h-12 min-h-[44px] px-6 sm:px-10 font-medium text-sm disabled:opacity-50 max-w-full"
                >
                  {loading ? (<><Loading inline size={13} /> Loading…</>) : `Start Diagnostic (${selectedConceptIds().length} concepts selected)`}
                </button>
              </div>
            </>
          )}

          {/* === STEP 2: Diagnostic === */}
          {step === 'diagnostic' && (
            <>
              <SectionHeader label={`Question ${questionCount} of ~${estimatedTotal}`} title={conceptName} />
              <div className="max-w-2xl mx-auto min-w-0 overflow-hidden px-2 sm:px-0">
                {/* Progress bar */}
                <div className="mb-4">
                  <div className="flex justify-between text-[10px] font-mono text-mathua-muted mb-1">
                    <span>Progress</span>
                    <span>{questionCount > estimatedTotal ? estimatedTotal : questionCount} / {estimatedTotal}</span>
                  </div>
                  <div className="h-1.5 bg-mathua-code rounded-full overflow-hidden">
                    <div
                      className="h-full bg-mathua-blue rounded-full transition-all duration-500"
                      style={{ width: `${Math.min((questionCount / estimatedTotal) * 100, 100)}%` }}
                    />
                  </div>
                </div>

                {/* Accuracy display */}
                {accuracy.total > 0 && (
                  <div className="mb-4 flex items-center gap-2 text-xs font-mono text-mathua-muted justify-center">
                    <span className={accuracy.correct / accuracy.total >= 0.7 ? 'text-mathua-green' : accuracy.correct / accuracy.total < 0.4 ? 'text-mathua-red' : ''}>
                      {accuracy.correct}/{accuracy.total}
                    </span>
                    <span>correct</span>
                    <div className="w-20 h-1 bg-mathua-code rounded-full overflow-hidden">
                      <div
                        className="h-full bg-mathua-blue rounded-full transition-all"
                        style={{ width: `${(accuracy.correct / Math.max(accuracy.total, 1)) * 100}%` }}
                      />
                    </div>
                  </div>
                )}

                <div className={`bg-mathua-surface border rounded-none p-4 sm:p-6 mb-6 transition-colors duration-200 w-full max-w-full min-w-0 overflow-hidden ${
                  lastResult
                    ? lastResult.correct ? 'border-green-500/40' : 'border-red-500/40'
                    : 'border-mathua-border'
                }`}>
                  <div className="bg-mathua-code border border-mathua-border rounded-none p-4 sm:p-6 text-center mb-4 w-full max-w-full min-w-0 overflow-hidden">
                    <div className="w-full max-w-full min-w-0 overflow-hidden">
                       <KatexContent className="text-mathua-primary text-lg font-mono font-light whitespace-pre-wrap break-words">
                         {question}
                       </KatexContent>
                     </div>
                   </div>

                    {!lastResult ? (
                      <>
                      <div className="flex flex-col sm:flex-row gap-3 min-w-0">
                        <input
                          ref={goalsInputRef}
                          type="text"
                          value={answerInput}
                          onChange={(e) => setAnswerInput(e.target.value)}
                          onKeyDown={(e) => e.key === 'Enter' && submitAnswer()}
                          placeholder="Your answer..."
                         disabled={loading}
                         className="flex-1 min-w-0 bg-mathua-code border border-mathua-border rounded-none h-12 px-4 font-mono text-base text-mathua-primary placeholder:text-mathua-muted focus:outline-none focus:border-mathua-blue"
                      />
                       <button
                         onClick={submitAnswer}
                         disabled={!answerInput.trim() || loading}
                         className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white rounded-none h-12 min-h-[36px] px-8 font-medium text-sm disabled:opacity-50 whitespace-nowrap shrink-0 w-full sm:w-auto"
                       >
                         Check Answer
                       </button>
                     </div>
                      <SymbolPalette targetRef={goalsInputRef} onInsert={setAnswerInput} />
                      </>
                    ) : (
                    <div className="animate-fadeIn text-center">
                      <p className={`text-base font-medium mb-2 ${lastResult.correct ? 'text-mathua-green' : 'text-mathua-red'}`}>
                        {lastResult.correct ? '✓ Correct!' : '✗ Not quite'}
                      </p>
                      <KatexContent className="text-mathua-secondary text-sm">{lastResult.feedback}</KatexContent>
                      {loading && <p className="text-mathua-muted text-xs mt-2"><Loading inline size={11} /> Loading next question…</p>}
                    </div>
                  )}
                </div>
              </div>
            </>
          )}

          {/* === STEP 3: Results === */}
          {step === 'results' && plan && (
            <>
              <SectionHeader label="Your results" title="Here's what we found" />
              <div className="mt-6 min-w-0 overflow-hidden">
                <DiagnosticResults plan={plan} onStartPractice={startPractice} />
              </div>
              {/* Quiz CTA — actionable after diagnostic, also reachable via ?quiz=1 */}
              <div className="mt-6 text-center">
                <button onClick={startQuiz} className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white rounded-none h-12 px-8 font-medium text-sm">Take Quiz (150 XP gate) →</button>
                <p className="font-mono text-xs text-mathua-muted mt-2">Guest allowed — unlimited retake</p>
              </div>
            </>
          )}

          {/* === QUIZ (reuse) — actionable every 150 XP, own grading path, guest unlimited === */}
          {step === 'quiz' && (
            <>
              <SectionHeader label={`Quiz ${quizCount} of 5`} title={quizConceptName} />
              <div className="max-w-2xl mx-auto min-w-0 overflow-hidden px-2 sm:px-0">
                <div className="mb-4 flex items-center gap-2 text-xs font-mono text-mathua-muted justify-center">
                  <span className={quizAccuracy.correct / Math.max(quizAccuracy.total, 1) >= 0.7 ? 'text-mathua-green' : ''}>{quizAccuracy.correct}/{quizAccuracy.total} correct</span>
                  {quizLastResult?.xp ? <span className="text-yellow-400">+{quizLastResult.xp} XP (TaskQuiz 20)</span> : null}
                </div>
                <div className={`bg-mathua-surface border rounded-none p-4 sm:p-6 mb-6 transition-colors w-full max-w-full min-w-0 overflow-hidden ${quizLastResult ? (quizLastResult.correct ? 'border-green-500/40' : 'border-red-500/40') : 'border-mathua-border'}`}>
                  <div className="bg-mathua-code border border-mathua-border rounded-none p-4 sm:p-6 text-center mb-4">
                    <KatexContent className="text-mathua-primary text-lg font-mono font-light whitespace-pre-wrap break-words">{quizQuestion}</KatexContent>
                  </div>
                  {!quizLastResult ? (
                    <>
                      <div className="flex flex-col sm:flex-row gap-3 min-w-0">
                        <input ref={quizInputRef} type="text" value={quizAnswerInput} onChange={e => setQuizAnswerInput(e.target.value)} onKeyDown={e => e.key === 'Enter' && submitQuizAnswerFn()} placeholder="Your answer..." disabled={loading} className="flex-1 min-w-0 bg-mathua-code border border-mathua-border rounded-none h-12 px-4 font-mono text-base text-mathua-primary placeholder:text-mathua-muted focus:outline-none focus:border-mathua-blue" />
                        <button onClick={submitQuizAnswerFn} disabled={!quizAnswerInput.trim() || loading} className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white rounded-none h-12 px-8 font-medium text-sm disabled:opacity-50">Check Answer</button>
                      </div>
                      <SymbolPalette targetRef={quizInputRef} onInsert={setQuizAnswerInput} />
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
          )}
          {step === 'quiz_done' && (
            <div className="max-w-2xl mx-auto text-center">
              <SectionHeader label="Quiz complete" title={`${quizAccuracy.correct}/${quizAccuracy.total} correct`} />
              <p className="font-mono text-sm text-mathua-secondary mt-4">TaskQuiz 20 XP awarded per correct — retake anytime.</p>
              <div className="mt-6 flex gap-3 justify-center">
                <button onClick={startQuiz} className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white rounded-none h-12 px-8 text-sm">Retake Quiz →</button>
                <Link href="/profile" className="border border-mathua-border text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue rounded-none h-12 px-8 text-sm inline-flex items-center">Back to Profile →</Link>
              </div>
            </div>
          )}
        </section>
      </div>
      <Footer />
    </>
  )
}

export default function GoalsPage() {
  return (
    <Suspense fallback={<><div style={{ background: 'var(--bg)', minHeight: '100vh' }}><Loading label="LOADING" /></div></>}>
      <GoalsContent />
    </Suspense>
  )
}
