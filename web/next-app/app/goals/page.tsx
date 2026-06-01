'use client'

import { useState, useEffect, useRef } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import { toast } from 'sonner'
import KatexContent from '../../components/KatexContent'
import { useTheme } from '../../hooks/useTheme'
import Header from '../../components/Header'
import SectionHeader from '../../components/SectionHeader'
import ProgressSummary from '../../components/ProgressSummary'
import Footer from '../../components/Footer'
import {
  startGoalDiagnostic,
  getGoalPlan,
  getScores,
  type GoalPlanRes,
  type Scores,
} from '../../lib/api'
import { isLoggedIn, getUserInfo } from '../../lib/auth'
import conceptsData from '../../data/concepts.json'

type Step = 'select' | 'diagnostic' | 'results'

interface DomainInfo {
  name: string
  count: number
  concepts: string[]
  selected: boolean
}

const domainLabels: Record<string, string> = {
  counting: 'Counting',
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

export default function GoalsPage() {
  const { mounted } = useTheme()
  const { push } = useRouter()

  const [step, setStep] = useState<Step>('select')
  const [scores, setScores] = useState<Scores | null>(null)
  const [loading, setLoading] = useState(false)

  // Step 1: goal selection
  const [domains, setDomains] = useState<DomainInfo[]>([])
  const [customConcepts] = useState<string[]>([])

  // Step 2: diagnostic
  const sessionId = useRef('')
  const tokenRef = useRef('')
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

  // Load scores and domains on mount
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
  }, [mounted])

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
      const order = ['counting', 'arithmetic', 'fractions', 'prealgebra', 'algebra', 'geometry', 'trigonometry', 'complex_numbers', 'precalculus', 'calculus', 'linear_algebra', 'statistics', 'discrete_math', 'number_theory', 'differential_equations', 'abstract_algebra', 'topology']
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
      const token = tokenRef.current
      const res = await fetch(`/api/goal/diagnostic/answer`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token || ''}` },
        body: JSON.stringify({ session_id: sessionId.current, concept_id: conceptId.current, answer, elapsed }),
      })
      const data = await res.json()
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

  function startPractice() {
    push('/session')
  }

  if (!mounted) return <div style={{ background: 'var(--bg)', minHeight: '100vh' }} />

  return (
    <>
      <Header />
      <div className="max-w-container mx-auto px-6 max-sm:px-4">
        <section className="pt-8">
          <span className="flex mb-4">
            <Link href="/" className="text-mathua-secondary text-sm hover:text-mathua-primary">Back</Link>
          </span>

          {/* === STEP 1: Goal Selection === */}
          {step === 'select' && (
            <>
              {scores && <ProgressSummary scores={scores} weakByDomain={weakByDomain || undefined} />}
              <SectionHeader label="Step 1" title="What do you want to learn?" />
              <p className="text-mathua-secondary text-sm text-center max-w-[600px] mx-auto mt-2 mb-8">
                Pick one or more topics. The system will test your prerequisite knowledge and build a personalized study plan.
              </p>

              <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-3 max-w-4xl mx-auto mb-8">
                {domains.map(d => {
                  const label = domainLabels[d.name] || d.name
                  return (
                    <button
                      key={d.name}
                      onClick={() => toggleDomain(d.name)}
                    className={`rounded-lg p-4 text-left transition-all text-sm ${
                      d.selected
                        ? 'bg-mathua-blue text-white'
                        : 'bg-mathua-surface border border-mathua-border text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue'
                    }`}
                    >
                      <div className="font-medium">{label}</div>
                      <div className={`font-mono text-[10px] mt-1 ${d.selected ? 'text-white/70' : 'text-mathua-muted'}`}>
                        {d.count} concepts
                      </div>
                    </button>
                  )
                })}
              </div>

              <div className="text-center">
                <button
                  onClick={startDiagnostic}
                  disabled={selectedConceptIds().length === 0 || loading}
                  className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white rounded-md h-12 px-10 font-medium text-sm disabled:opacity-50"
                >
                  {loading ? 'Loading…' : `Start Diagnostic (${selectedConceptIds().length} concepts selected)`}
                </button>
              </div>
            </>
          )}

          {/* === STEP 2: Diagnostic === */}
          {step === 'diagnostic' && (
            <>
              <SectionHeader label={`Question ${questionCount} of ~${estimatedTotal}`} title={conceptName} />
              <div className="max-w-2xl mx-auto">
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

                <div className={`bg-mathua-surface border rounded-lg p-6 mb-6 transition-colors duration-200 ${
                  lastResult
                    ? lastResult.correct ? 'border-green-500/40' : 'border-red-500/40'
                    : 'border-mathua-border'
                }`}>
                  <div className="bg-mathua-code border border-mathua-border rounded-md p-6 text-center mb-4">
                    <div className="select-none" onCopy={(e) => e.preventDefault()}>
                       <KatexContent className="text-mathua-primary text-lg font-mono font-light whitespace-pre-wrap">
                         {question}
                       </KatexContent>
                     </div>
                   </div>

                   {!lastResult ? (
                     <div className="flex gap-3">
                       <input
                         type="text"
                         value={answerInput}
                         onChange={(e) => setAnswerInput(e.target.value)}
                         onPaste={(e) => e.preventDefault()}
                         onKeyDown={(e) => e.key === 'Enter' && submitAnswer()}
                         placeholder="Your answer..."
                        disabled={loading}
                        className="flex-1 bg-mathua-code border border-mathua-border rounded-md h-12 px-4 font-mono text-base text-mathua-primary placeholder:text-mathua-muted focus:outline-none focus:border-mathua-blue"
                      />
                      <button
                        onClick={submitAnswer}
                        disabled={!answerInput.trim() || loading}
                        className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white rounded-md h-12 px-8 font-medium text-sm disabled:opacity-50 whitespace-nowrap"
                      >
                        Check Answer
                      </button>
                    </div>
                  ) : (
                    <div className="animate-fadeIn text-center">
                      <p className={`text-base font-medium mb-2 ${lastResult.correct ? 'text-mathua-green' : 'text-mathua-red'}`}>
                        {lastResult.correct ? '✓ Correct!' : '✗ Not quite'}
                      </p>
                      <KatexContent className="text-mathua-secondary text-sm">{lastResult.feedback}</KatexContent>
                      {loading && <p className="text-mathua-muted text-xs mt-2">Loading next question…</p>}
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

              <div className="max-w-2xl mx-auto">
                {/* Readiness meter */}
                <div className="text-center mb-10">
                  <div className="relative inline-flex items-center justify-center">
                    <svg width="140" height="140" className="-rotate-90">
                      <circle cx="70" cy="70" r="60" fill="none" stroke="var(--code-bg)" strokeWidth="10" />
                      <circle
                        cx="70" cy="70" r="60" fill="none" stroke="var(--accent-blue)"
                        strokeWidth="10" strokeLinecap="round"
                        strokeDasharray={`${2 * Math.PI * 60}`}
                        strokeDashoffset={`${2 * Math.PI * 60 * (1 - plan.readiness)}`}
                        className="transition-all duration-1000"
                      />
                    </svg>
                    <span className="absolute text-3xl font-mono font-light text-mathua-primary">
                      {Math.round(plan.readiness * 100)}%
                    </span>
                  </div>
                  <p className="text-mathua-secondary text-sm mt-2">
                    readiness for your selected topics
                  </p>
                  {plan.total_tested > 0 && (
                    <p className="text-mathua-muted text-xs font-mono mt-1">
                      {plan.correct_count}/{plan.total_tested} correct
                    </p>
                  )}
                </div>

                <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-8">
                  {/* Weak areas */}
                  {Object.keys(plan.weak_areas).length > 0 && (
                    <div>
                      <h3 className="font-serif text-sm text-mathua-red mb-3 font-mono border-b border-mathua-border pb-2">
                        Needs attention ({Object.values(plan.weak_areas).flat().length})
                      </h3>
                      <div className="space-y-2">
                        {Object.entries(plan.weak_areas).map(([domain, concepts]) => {
                          const label = domainLabels[domain] || domain
                          return (
                            <div key={domain} className="bg-mathua-surface border border-mathua-border rounded-none p-3">
                              <div className="flex items-center justify-between">
                                <span className="font-mono text-xs text-mathua-primary">{label}</span>
                                <span className="font-mono text-xs text-mathua-red">{concepts.length}</span>
                              </div>
                            </div>
                          )
                        })}
                      </div>
                    </div>
                  )}

                  {/* Strong areas */}
                  {Object.keys(plan.strong_areas).length > 0 && (
                    <div>
                      <h3 className="font-serif text-sm text-mathua-green mb-3 font-mono border-b border-mathua-border pb-2">
                        Strong areas ({Object.values(plan.strong_areas).flat().length})
                      </h3>
                      <div className="space-y-2">
                        {Object.entries(plan.strong_areas).map(([domain, concepts]) => {
                          const label = domainLabels[domain] || domain
                          return (
                            <div key={domain} className="bg-mathua-surface border border-mathua-border rounded-none p-3">
                              <div className="flex items-center justify-between">
                                <span className="font-mono text-xs text-mathua-primary">{label}</span>
                                <span className="font-mono text-xs text-mathua-green">{concepts.length}</span>
                              </div>
                            </div>
                          )
                        })}
                      </div>
                    </div>
                  )}
                </div>

                <div className="text-center">
                  <button
                    onClick={startPractice}
                    className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white rounded-none h-12 px-10 font-medium text-sm"
                  >
                    Start practicing →
                  </button>
                </div>
              </div>
            </>
          )}
        </section>
      </div>
      <Footer />
    </>
  )
}
