'use client'

import { useState, useEffect, Suspense } from 'react'
import { useRouter } from 'next/navigation'
import { useTheme } from '../../hooks/useTheme'
import ProgressSummary from '../../components/ProgressSummary'
import Header from '../../components/Header'
import BottomTabs from '../../components/BottomTabs'
import SectionHeader from '../../components/SectionHeader'
import Footer from '../../components/Footer'
import DiagnosticResults from '../../components/DiagnosticResults'
import Loading from '../../components/Loading'
import { getScores, getWeaknesses, type GoalPlanRes, type Scores } from '../../lib/api'
import { isLoggedIn, getUserInfo } from '../../lib/auth'
import { concepts as conceptsData } from '../../lib/conceptData'
import QuizHost from './QuizHost'
import DiagnosticHost from './DiagnosticHost'
import { GOALS_DIAG_KEY } from './constants'

type Step = 'select' | 'diagnostic' | 'results' | 'quiz'

interface DomainInfo {
  name: string
  count: number
  concepts: string[]
  selected: boolean
}

const domainLabels = {
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
} satisfies Record<string, string>

function GoalsContent() {
  const { mounted } = useTheme()
  const { push } = useRouter()

  const [step, setStep] = useState<Step>(() => {
    if (typeof window === 'undefined') return 'select'
    return new URLSearchParams(window.location.search).get('quiz') === '1' ? 'quiz' : 'select'
  })
  const [scores, setScores] = useState<Scores | null>(null)

  // Step 1: goal selection
  const [domains, setDomains] = useState<DomainInfo[]>([])
  const [customConcepts] = useState<string[]>([])
  const [hasPaused, setHasPaused] = useState(false)

  // Diagnostic handoff (the host owns the session once started)
  const [diagStart, setDiagStart] = useState<string[] | null>(null)
  const [diagResume, setDiagResume] = useState<string | null>(null)

  // Step 3: results
  const [plan, setPlan] = useState<GoalPlanRes | null>(null)
  const [weakByDomain, setWeakByDomain] = useState<Record<string, { id: string; label: string }[]> | null>(null)

  // Load scores and domains on mount — QuizHost auto-starts when ?quiz=1.
  useEffect(() => {
    if (!mounted) return
    const loggedIn = isLoggedIn()
    if (loggedIn) {
      const user = getUserInfo()
      if (user) {
        getScores(user.student_id).then(setScores).catch(e => console.error('scores fetch failed:', e))
      }
      getWeaknesses().then(d => {
        if (d.by_domain) setWeakByDomain(d.by_domain)
      }).catch(e => console.error('weaknesses fetch failed:', e))
    }
    buildDomains()
    try {
      setHasPaused(!!sessionStorage.getItem(GOALS_DIAG_KEY))
    } catch {
      setHasPaused(false)
    }
  }, [mounted])

  function buildDomains() {
    const map = new Map<string, string[]>()
    for (const c of conceptsData) {
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

  function startDiagnostic() {
    const ids = selectedConceptIds()
    if (ids.length === 0) return
    setDiagStart(ids)
    setDiagResume(null)
    setStep('diagnostic')
  }

  function resumeDiagnostic() {
    let sid = ''
    try {
      sid = sessionStorage.getItem(GOALS_DIAG_KEY) || ''
    } catch {
      sid = ''
    }
    if (!sid) return
    setDiagResume(sid)
    setDiagStart(null)
    setStep('diagnostic')
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
            <button type="button" onClick={() => { if (window.history.length > 1) window.history.back(); else push('/profile') }} className="text-mathua-secondary text-sm hover:text-mathua-primary">← Back</button>
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
                      type="button"
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
                    type="button"
                    onClick={resumeDiagnostic}
                    className="border border-mathua-blue bg-mathua-blue text-white hover:opacity-90 rounded-none h-12 min-h-[44px] px-6 sm:px-10 font-medium text-sm disabled:opacity-50 max-w-full mb-3"
                  >
                    Continue diagnostic →
                  </button>
                )}
                <button
                  type="button"
                  onClick={startDiagnostic}
                  disabled={selectedConceptIds().length === 0}
                  className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white rounded-none h-12 min-h-[44px] px-6 sm:px-10 font-medium text-sm disabled:opacity-50 max-w-full"
                >
                  {`Start diagnostic test (${selectedConceptIds().length} concepts selected)`}
                </button>
              </div>
            </>
          )}

          {/* === STEP 2: Diagnostic === */}
          {step === 'diagnostic' && (
            <DiagnosticHost
              startIds={diagStart}
              resumeId={diagResume}
              onComplete={(p) => {
                setPlan(p)
                setHasPaused(false)
                setStep('results')
              }}
              onResumeExpired={() => {
                setHasPaused(false)
                setStep('select')
              }}
            />
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
                <button type="button" onClick={() => setStep('quiz')} className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white rounded-none h-12 px-8 font-medium text-sm">Take Quiz (150 XP gate) →</button>
                <p className="font-mono text-xs text-mathua-muted mt-2">Guest allowed, unlimited retake</p>
              </div>
            </>
          )}

          {/* === QUIZ — actionable every 150 XP, own grading path, guest unlimited === */}
          {step === 'quiz' && <QuizHost />}
        </section>
      </div>
      <Footer />
      <BottomTabs />
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
