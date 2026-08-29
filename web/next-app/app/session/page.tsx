'use client'

import { useTheme } from '../../hooks/useTheme'
import { useState, useEffect, useCallback, useRef } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import Image from 'next/image'
import { toast } from 'sonner'
import Header from '../../components/Header'
import BottomTabs from '../../components/BottomTabs'
import SectionHeader from '../../components/SectionHeader'
import Footer from '../../components/Footer'
import AsciiDivider from '../../components/AsciiDivider'
import KatexContent from '../../components/KatexContent'

import SymbolPalette from '../../components/SymbolPalette'
import {
  API_BASE,
  getConfig,
  startSession,
  startSessionName,
  submitAnswer,
  getScores,
  getSettings,
  updateSettings,
  startGoalDiagnosticName,
  submitGoalAnswer,
  getGoalPlan,
  setDailyXPGoal,
  getDueReviews,
  startReviewSession,
  submitReviewAnswer,
  type Question,
  type AnswerResult,
  type Scores,
  type GoalPlanRes,
} from '../../lib/api'
import { isLoggedIn, getUserInfo, clearToken, type UserInfo } from '../../lib/auth'
import conceptsData from '../../data/concepts.json'
import Loading from '../../components/Loading'

import DiagnosticResults from '../../components/DiagnosticResults'

type Screen = 'name' | 'diag_select' | 'diagnostic' | 'diagnostic_results' | 'practice' | 'review'

export default function SessionPage() {
  const { mounted } = useTheme()
  const { push } = useRouter()

  const [screen, setScreen] = useState<Screen>('practice')
  const [name, setName] = useState('')
  const [user, setUser] = useState<UserInfo | null>(null)
  const [studentID, setStudentID] = useState('')
  const [sessionID, setSessionID] = useState('')
  const [question, setQuestion] = useState<Question | null>(null)
  const [attemptId, setAttemptId] = useState('')
  const [lastResult, setLastResult] = useState<AnswerResult | null>(null)
  const [submitted, setSubmitted] = useState(false)
  const [answer, setAnswer] = useState('')
  const [sessionStats, setSessionStats] = useState({ correct: 0, total: 0 })
  const [scores, setScores] = useState<Scores>({
    lifetime_points: 0, weekly_score: 0, speed_bonus: 0,
    concepts_mastered: 0, current_streak: 0, level: 'Novice',
    xp_total: 0, xp_today: 0, daily_xp_goal: 150,
  })
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(true)
  const [showTimer, setShowTimer] = useState(false)
  const [dueReviews, setDueReviews] = useState(0)
  const [elapsed, setElapsed] = useState(0)
  const startRef = useRef(Date.now())
  const timerRef = useRef<ReturnType<typeof setInterval>>()
  // Review mode state
  const [reviewSessionID, setReviewSessionID] = useState('')
  const [reviewStats, setReviewStats] = useState({ correct: 0, total: 0 })
  const [reviewDone, setReviewDone] = useState(false)
  const inputRef = useRef<HTMLInputElement>(null)
  const submittingRef = useRef(false)

   // Diagnostic state (guest mode)
  const [domains, setDomains] = useState<{ name: string; concepts: string[]; selected: boolean }[]>([])
  const diagSessionId = useRef('')
  const [diagQuestion, setDiagQuestion] = useState('')
  const diagConceptId = useRef('')
  const [diagConceptName, setDiagConceptName] = useState('')
  const [diagCount, setDiagCount] = useState(0)
  const [diagAnswer, setDiagAnswer] = useState('')
  const [diagLastResult, setDiagLastResult] = useState<{ correct: boolean; feedback: string } | null>(null)
  const [diagAccuracy, setDiagAccuracy] = useState<{ correct: number; total: number }>({ correct: 0, total: 0 })
  const [diagPlan, setDiagPlan] = useState<GoalPlanRes | null>(null)
  const guestStudentID = useRef('')
  const diagInputRef = useRef<HTMLInputElement>(null)

  const beginSessionAuth = useCallback(async () => {
    setError('')
    setLoading(true)
    try {
      const res = await startSession()
      setStudentID(res.student_id)
      setSessionID(res.session_id)
      setQuestion(res.question)
      setAttemptId(res.question?.attempt_id ?? '')
      setScreen('practice')
      setSubmitted(false)
      setSessionStats({ correct: 0, total: 0 })
      const s = await getScores(res.student_id).catch(() => null)
      if (!s) console.error('getScores failed')
      if (s) setScores(s)
      getDueReviews().then(r => setDueReviews(r.count)).catch(() => console.error('getDueReviews failed'))
    } catch {
      clearToken()
      setScreen('name')
      setError('Session expired. Please log in again.')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    getConfig()
      .then(config => {
        if (!config.auth_enabled) {
          setScreen('name')
          setLoading(false)
          return
        }
        const u = getUserInfo()
        if (u) {
          setUser(u)
          setScreen('practice')
          getSettings().then(s => setShowTimer(s.show_timer ?? false)).catch(() => console.error('getSettings failed'))
          beginSessionAuth()
        } else {
          setScreen('name')
          setLoading(false)
        }
      })
      .catch(() => {
        setScreen('name')
        setLoading(false)
      })
  }, [beginSessionAuth])

  useEffect(() => {
    if (screen === 'practice' && !submitted) {
      startRef.current = Date.now()
      inputRef.current?.focus()
    }
  }, [screen, question, submitted])

  useEffect(() => {
    if (screen === 'practice' && !submitted && showTimer) {
      timerRef.current = setInterval(() => {
        setElapsed((Date.now() - startRef.current) / 1000)
      }, 100)
    }
    return () => { if (timerRef.current) clearInterval(timerRef.current) }
  }, [screen, submitted, showTimer])

  const beginSessionName = useCallback(async () => {
    if (!name.trim()) { setError('Enter your name'); return }
    setError(''); setLoading(true)
    try {
      const res = await startSessionName(name.trim())
      setStudentID(res.student_id)
      setSessionID(res.session_id)
      setQuestion(res.question)
      setAttemptId(res.question?.attempt_id ?? '')
      setScreen('practice')
      setSubmitted(false)
      setSessionStats({ correct: 0, total: 0 })
      const s = await getScores(res.student_id).catch(() => null)
      if (s) setScores(s)
      if (!s) console.error('getScores failed')
    } catch {
      setError('Could not connect to server. Is the backend running?')
    } finally { setLoading(false) }
  }, [name])

  const handleSubmit = useCallback(async () => {
    if (!answer.trim() || !question || !attemptId || submittingRef.current) return
    submittingRef.current = true
    setError('')
    const e = (Date.now() - startRef.current) / 1000
    try {
      const res = await submitAnswer(sessionID, answer.trim(), e, attemptId)
      setLastResult(res.result)
      setSubmitted(true)
      if (res.result) {
        setSessionStats(prev => ({
          correct: prev.correct + (res.result!.correct ? 1 : 0),
          total: prev.total + 1,
        }))
      }
      if (res.next_question) {
        setQuestion(res.next_question)
        setAttemptId(res.next_question.attempt_id ?? '')
      } else {
        setQuestion(null)
        setAttemptId('')
      }
      const s = await getScores(studentID).catch(() => null)
      if (s) setScores(s)
    } catch {
      setError('Failed to submit answer. If you answered in another tab, refresh to continue.')
    } finally {
      submittingRef.current = false
    }
  }, [answer, question, attemptId, sessionID, studentID])

  const nextQuestion = useCallback(() => {
    setLastResult(null)
    setSubmitted(false)
    setAnswer('')
  }, [])

  const beginReviewSession = useCallback(async () => {
    setError('')
    setLoading(true)
    try {
      const res = await startReviewSession()
      setReviewSessionID(res.session_id)
      setQuestion(res.question)
      setAttemptId(res.question?.attempt_id ?? '')
      setScreen('review')
      setSubmitted(false)
      setReviewStats({ correct: 0, total: 0 })
      setReviewDone(false)
      startRef.current = Date.now()
    } catch {
      setError('Could not start review session.')
    } finally {
      setLoading(false)
    }
  }, [])

  const handleReviewSubmit = useCallback(async () => {
    if (!answer.trim() || !question || !attemptId || submittingRef.current) return
    submittingRef.current = true
    setError('')
    const e = (Date.now() - startRef.current) / 1000
    try {
      const res = await submitReviewAnswer(reviewSessionID, answer.trim(), e, attemptId)
      setLastResult(res.result)
      setSubmitted(true)
      if (res.result) {
        setReviewStats(prev => ({
          correct: prev.correct + (res.result!.correct ? 1 : 0),
          total: prev.total + 1,
        }))
      }
      if (res.next_question) {
        setQuestion(res.next_question)
        setAttemptId(res.next_question.attempt_id ?? '')
      } else {
        setQuestion(null)
        setAttemptId('')
        setReviewDone(true)
      }
      const s = await getScores(studentID).catch(() => null)
      if (s) setScores(s)
    } catch {
      setError('Failed to submit review answer. If you answered in another tab, refresh to continue.')
    } finally {
      submittingRef.current = false
    }
  }, [answer, question, attemptId, reviewSessionID, studentID])

  const nextReviewQuestion = useCallback(() => {
    if (reviewDone) {
      setScreen('practice')
      setReviewDone(false)
      return
    }
    setLastResult(null)
    setSubmitted(false)
    setAnswer('')
    startRef.current = Date.now()
  }, [reviewDone])

  // Guest diagnostic logic
  const domainOrder = ['arithmetic', 'fractions', 'prealgebra', 'algebra', 'geometry', 'trigonometry', 'complex_numbers', 'precalculus', 'calculus', 'linear_algebra', 'statistics', 'discrete_math', 'number_theory', 'differential_equations', 'abstract_algebra', 'topology']
  const domainLabels: Record<string, string> = {
    arithmetic: 'Arithmetic', fractions: 'Fractions', prealgebra: 'Pre-Algebra',
    algebra: 'Algebra', geometry: 'Geometry', trigonometry: 'Trigonometry',
    complex_numbers: 'Complex Numbers', precalculus: 'Precalculus', calculus: 'Calculus',
    linear_algebra: 'Linear Algebra', statistics: 'Statistics', discrete_math: 'Discrete Math',
    number_theory: 'Number Theory', differential_equations: 'Differential Equations',
    abstract_algebra: 'Abstract Algebra', topology: 'Topology',
  }

  function toggleDomain(name: string) {
    setDomains(prev => prev.map(d => ({ ...d, selected: d.name === name ? !d.selected : false })))
  }

  async function beginGuestDiagnostic() {
    const ids: string[] = []
    for (const d of domains) {
      if (d.selected) ids.push(...d.concepts)
    }
    if (ids.length === 0) return
    setLoading(true)
    try {
      const res = await startGoalDiagnosticName(name.trim(), ids)
      if (res.student_id) guestStudentID.current = res.student_id
      if (res.done) {
        setDiagPlan({ readiness: 1, total_tested: 0, correct_count: 0, weak_areas: {}, strong_areas: {} })
        setScreen('diagnostic_results')
        return
      }
      diagSessionId.current = res.session_id
      setDiagQuestion(res.question || '')
      diagConceptId.current = res.concept_id || ''
      setDiagConceptName(res.concept_name || '')
      setDiagCount(1)
      setDiagAccuracy({ correct: 0, total: 0 })
      setDiagLastResult(null)
      setDiagAnswer('')
      setScreen('diagnostic')
    } catch {
      toast.error("Something went wrong, but we're working on it.")
    } finally {
      setLoading(false)
    }
  }

  async function submitGuestDiagnostic() {
    if (!diagAnswer.trim()) return
    setLoading(true)
    try {
      const elapsed = 5.0
      const data = await submitGoalAnswer(diagSessionId.current, diagConceptId.current, diagAnswer.trim(), elapsed)
      const correct = data.correct || false
      const feedback = data.feedback || (correct ? 'Correct!' : 'Not quite.')
      setDiagAccuracy(prev => ({ correct: prev.correct + (correct ? 1 : 0), total: prev.total + 1 }))
      setDiagLastResult({ correct, feedback })

      if (data.done) {
        setTimeout(async () => {
          try {
            const planRes = await getGoalPlan(diagSessionId.current)
            setDiagPlan(planRes)
            setScreen('diagnostic_results')
          } catch {
            setError('Could not generate plan.')
          }
          setLoading(false)
        }, 800)
        return
      }

      setTimeout(() => {
        setDiagQuestion(data.question || '')
        diagConceptId.current = data.concept_id || ''
        setDiagConceptName(data.concept_name || '')
        setDiagCount(prev => prev + 1)
        setDiagLastResult(null)
        setDiagAnswer('')
        setLoading(false)
      }, 1200)
    } catch {
      setError('Failed to submit.')
      setLoading(false)
    }
  }

  function startGuestDiagnostic() {
    setError('')
    const raw = conceptsData as any[]
    const map = new Map<string, string[]>()
    for (const c of raw) {
      const list = map.get(c.domain) || []
      list.push(c.id)
      map.set(c.domain, list)
    }
    const result: { name: string; concepts: string[]; selected: boolean }[] = []
    map.forEach((concepts, name) => result.push({ name, concepts, selected: true }))
    result.sort((a, b) => domainOrder.indexOf(a.name) - domainOrder.indexOf(b.name))
    setDomains(result)
    setScreen('diag_select')
  }

  if (!mounted) return <div style={{ background: 'var(--bg)', minHeight: '100vh' }} />

  return (
    <>
      <Header />
      <div className="max-w-container mx-auto px-4 sm:px-6 pb-[calc(80px+env(safe-area-inset-bottom))] lg:pb-0 overflow-x-hidden min-w-0">
      <section className="pt-8 min-w-0 overflow-hidden">
        <span className="flex justify-between items-center mb-4">
          <Link href="/" className="text-mathua-secondary text-sm hover:text-mathua-primary">
            ← Back
          </Link>
          <span className="flex gap-3 items-center">
            {user && <span className="font-mono text-[11px] text-mathua-muted">{user.name}</span>}
            <button onClick={() => { clearToken(); push('/login') }} className="text-mathua-secondary text-xs hover:text-mathua-red">
              Logout
            </button>
          </span>
        </span>

        {screen === 'name' && (
          <div className="max-w-md mx-auto mt-20">
            <SectionHeader label="Practice session" title="What's your name?" />
            <p className="text-mathua-secondary text-sm text-center mt-2 mb-4">
              Sign up or login to save your progress permanently.
            </p>
            <div className="mt-6">
              <input
                type="text"
                value={name}
                onChange={(e) => setName(e.target.value)}
                onKeyDown={(e) => e.key === 'Enter' && beginSessionName()}
                placeholder="Your name"
                className="w-full bg-mathua-code border border-mathua-border rounded-none h-12 px-4 font-mono text-base text-mathua-primary placeholder:text-mathua-muted focus:outline-none focus:border-mathua-blue"
              />
              <div className="flex flex-col gap-3 mt-4">
                <button
                  onClick={beginSessionName}
                  className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white rounded-none h-12 px-4 font-medium text-sm"
                >
                  Learning Mode →
                </button>
                <button
                  onClick={startGuestDiagnostic}
                  className="border border-mathua-secondary text-mathua-secondary hover:bg-mathua-secondary hover:text-white rounded-none h-12 px-4 font-medium text-sm"
                >
                  Take a Diagnostic →
                </button>
              </div>
            </div>
            <div className="mt-6 text-center">
              <Link href="/study" className="block text-mathua-blue text-sm hover:text-mathua-blue-hover">
                Browse study lessons →
              </Link>
              <Link href="/login" className="block text-mathua-secondary text-sm hover:text-mathua-blue">
                Have an account? Sign in
              </Link>
            </div>
          </div>
        )}

        {screen === 'diag_select' && (
          <div className="max-w-4xl mx-auto mt-8">
            <SectionHeader label="Diagnostic" title="What do you want to learn?" />
            <p className="text-mathua-secondary text-sm text-center max-w-[500px] mx-auto mt-2 mb-6">
              We&apos;ll test your current knowledge and find the right starting point.
            </p>
            <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-3 mb-8">
              {domains.map(d => {
                const label = domainLabels[d.name] || d.name
                return (
                  <button
                    key={d.name}
                    onClick={() => toggleDomain(d.name)}
                    className={`rounded-none p-4 text-left transition-all text-sm ${
                      d.selected
                        ? 'bg-mathua-blue text-white'
                        : 'bg-mathua-surface border border-mathua-border text-mathua-secondary hover:bg-mathua-blue hover:text-white'
                    }`}
                  >
                    <div className="font-medium">{label}</div>
                    <div className={`font-mono text-[10px] mt-1 ${d.selected ? 'text-white/70' : 'text-mathua-muted'}`}>
                      {d.concepts.length} concepts
                    </div>
                  </button>
                )
              })}
            </div>
            <div className="flex gap-3 justify-center">
              <button onClick={() => setScreen('name')} className="bg-mathua-surface border border-mathua-border rounded-none h-12 px-8 text-sm text-mathua-secondary hover:text-mathua-blue">
                Back
              </button>
              <button onClick={beginGuestDiagnostic} disabled={loading} className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white rounded-none h-12 px-10 font-medium text-sm disabled:opacity-50">
                {loading ? (<><Loading inline size={13} /> Loading…</>) : 'Start diagnostic'}
              </button>
            </div>
          </div>
        )}

        {screen === 'diagnostic' && (
          <div className="max-w-2xl mx-auto mt-8">
            <SectionHeader label={`Question ${diagCount}`} title={diagConceptName} />
            <div className="bg-mathua-surface border border-mathua-border rounded-none p-6 mb-6">
              <div className="bg-mathua-code border border-mathua-border rounded-none p-6 text-center mb-4">
                <div className="">
                  <KatexContent className="text-mathua-primary text-lg font-mono font-light whitespace-pre-wrap">{diagQuestion}</KatexContent>
                </div>
               </div>
               <div className="flex flex-col sm:flex-row gap-3 min-w-0">
                  <input
                    ref={diagInputRef}
                    type="text"
                    value={diagAnswer}
                    onChange={(e) => setDiagAnswer(e.target.value)}
                    onKeyDown={(e) => e.key === 'Enter' && submitGuestDiagnostic()}
                    placeholder="Your answer..."
                    disabled={loading || diagLastResult !== null}
                    className="flex-1 min-w-0 bg-mathua-code border border-mathua-border rounded-none h-12 px-4 font-mono text-base text-mathua-primary placeholder:text-mathua-muted focus:outline-none focus:border-mathua-blue"
                />
                 <button onClick={submitGuestDiagnostic} disabled={!diagAnswer.trim() || loading || diagLastResult !== null} className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white rounded-none h-12 px-8 font-medium text-sm disabled:opacity-50 shrink-0 w-full sm:w-auto min-h-[36px]">
                   Check Answer
                 </button>
               </div>
               <SymbolPalette targetRef={diagInputRef} onInsert={setDiagAnswer} />
            </div>
            {diagLastResult && (
              <div className={`bg-mathua-surface border rounded-none p-4 mb-4 text-center ${diagLastResult.correct ? 'border-mathua-green' : 'border-mathua-red'}`}>
                <KatexContent className={diagLastResult.correct ? 'text-mathua-green' : 'text-mathua-red'}>{diagLastResult.feedback}</KatexContent>
              </div>
            )}
            <div className="text-center text-mathua-muted text-xs font-mono">
              {diagAccuracy.total > 0 && `${diagAccuracy.correct}/${diagAccuracy.total} correct`}
            </div>
          </div>
        )}

        {screen === 'diagnostic_results' && diagPlan && (
          <div className="mt-8">
            <SectionHeader label="Diagnostic complete" title="Your results" />
            <div className="mt-6">
              <DiagnosticResults
                plan={diagPlan}
                onStartPractice={async () => {
                  if (guestStudentID.current) {
                    try {
                      const sessRes = await fetch(`${API_BASE}/api/session`, {
                        method: 'POST',
                        headers: { 'Content-Type': 'application/json' },
                        body: JSON.stringify({ student_id: guestStudentID.current }),
                      })
                      const sessData = await sessRes.json()
                      setStudentID(sessData.student_id)
                      setSessionID(sessData.session_id)
                      setQuestion(sessData.question)
                      setScreen('practice')
                      getScores(guestStudentID.current).then(setScores).catch(() => {})
                    } catch { setError('Could not start practice') }
                  } else {
                    setScreen('practice')
                  }
                }}
              />
              <div className="mt-4 text-center">
                <Link href="/profile" className="font-mono text-xs text-mathua-muted hover:text-mathua-blue">View profile →</Link>
              </div>
            </div>
          </div>
        )}

        {loading && (
          <div className="max-w-md mx-auto mt-20 text-center">
            <Loading label="LOADING" />
          </div>
        )}
        {error && (
          <div className="max-w-md mx-auto mt-20 text-center">
            <p className="text-mathua-red text-sm">{error}</p>
          </div>
        )}

        {screen === 'practice' && (
          <>
            {!question && !submitted && (
              <div className="max-w-xl mx-auto mt-20 text-center">
                <p className="text-mathua-muted text-sm mb-4">No questions available.</p>
                <button
                  onClick={() => beginSessionAuth()}
                  className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white rounded-none h-12 px-8 font-medium text-sm"
                >
                  Start new session
                </button>
              </div>
            )}

            {dueReviews > 0 && (
              <div className="bg-mathua-surface border border-yellow-500/40 rounded-none p-3 mb-4 flex items-center justify-center gap-4">
                <p className="font-mono text-xs text-yellow-400">
                  {dueReviews} concept{dueReviews !== 1 ? 's' : ''} due for review
                </p>
                <button
                  onClick={beginReviewSession}
                  className="border border-yellow-500/60 text-yellow-400 hover:bg-yellow-500 hover:text-black rounded-none px-4 h-8 font-mono text-[11px] transition-colors"
                >
                  Review Now
                </button>
              </div>
            )}

            {question && (
              <div className="flex gap-8 items-start mt-4 max-md:flex-col">
                <div className="w-[200px] flex-shrink-0 max-md:w-full">
                  {/* Desktop sidebar */}
                  <div className="bg-mathua-surface border border-mathua-border rounded-none p-5 space-y-4 max-md:hidden">
                    <div
                      className="transition-opacity duration-200"
                      key={lastResult?.streak ?? 0}
                    >
                      <span className="font-mono text-[10px] uppercase text-mathua-muted">Streak</span>
                      <div className="font-mono text-2xl text-mathua-green mt-1">
                        {submitted && lastResult ? lastResult.streak : lastResult ? lastResult.streak : '--'}
                      </div>
                    </div>
                    <div>
                      <div className="flex items-center justify-between">
                        <span className="font-mono text-[10px] uppercase text-mathua-muted">Timer</span>
                        <button
                          onClick={() => {
                            const next = !showTimer
                            setShowTimer(next)
                            updateSettings({ show_timer: next }).catch(() => console.error('updateSettings failed'))
                          }}
                          className={`font-mono text-[10px] px-2 py-0.5 border transition-colors ${
                            showTimer
                              ? 'bg-mathua-blue text-white border-mathua-blue'
                              : 'text-mathua-muted border-mathua-border hover:text-mathua-primary'
                          }`}
                        >
                          {showTimer ? 'ON' : 'OFF'}
                        </button>
                      </div>
                      {showTimer && (
                        <div className="font-mono text-2xl text-mathua-primary mt-2">
                          {!submitted ? `${elapsed.toFixed(1)}s` : '·'}
                        </div>
                      )}
                    </div>
                    <div>
                      <span className="font-mono text-[10px] uppercase text-mathua-muted">This session</span>
                      <div className="font-mono text-lg text-mathua-primary mt-1">
                        <span className={sessionStats.total > 0 && sessionStats.correct / sessionStats.total >= 0.8 ? 'text-mathua-green' : sessionStats.total > 0 && sessionStats.correct / sessionStats.total < 0.5 ? 'text-mathua-red' : 'text-mathua-primary'}>
                          {sessionStats.correct}/{sessionStats.total}
                        </span>
                        <span className="text-mathua-muted text-xs ml-1">correct</span>
                      </div>
                      {sessionStats.total > 0 && (
                        <div className="mt-1 h-1 bg-mathua-code rounded-full overflow-hidden">
                          <div
                            className="h-full bg-mathua-blue rounded-full transition-all duration-500"
                            style={{ width: `${(sessionStats.correct / sessionStats.total) * 100}%` }}
                          />
                        </div>
                      )}
                    </div>
                    <div>
                      <span className="font-mono text-[10px] uppercase text-mathua-muted">Mastered total</span>
                      <div className="font-mono text-2xl text-mathua-blue mt-1">{scores.concepts_mastered}</div>
                    </div>
                    <div>
                      <span className="font-mono text-[10px] uppercase text-mathua-muted">Daily XP</span>
                      <div className="flex items-baseline gap-2 mt-1">
                        <span className="font-mono text-2xl text-mathua-blue">{scores.xp_today ?? 0}</span>
                        <span className="font-mono text-[10px] text-mathua-muted">/ {scores.daily_xp_goal ?? 150}</span>
                        <button
                          onClick={() => {
                            const g = prompt('Set daily XP goal:', String(scores.daily_xp_goal || 150))
                            if (g) {
                              const n = parseInt(g, 10)
                              if (n > 0 && n <= 10000) {
                                setDailyXPGoal(n).then(() => {
                                  setScores(prev => ({ ...prev, daily_xp_goal: n }))
                                }).catch(() => console.error('setDailyXPGoal failed'))
                              }
                            }
                          }}
                          className="font-mono text-[10px] text-mathua-blue hover:text-mathua-blue-hover ml-1"
                        >
                          edit
                        </button>
                      </div>
                      <div className="mt-2 h-1.5 bg-mathua-code rounded-full overflow-hidden">
                        <div
                          className="h-full bg-mathua-blue rounded-full transition-all duration-500"
                          style={{ width: `${Math.min(((scores.xp_today ?? 0) / (scores.daily_xp_goal || 150)) * 100, 100)}%` }}
                        />
                      </div>
                      {(scores.xp_today ?? 0) >= (scores.daily_xp_goal || 150) && (
                        <div className="mt-1 font-mono text-[10px] text-mathua-blue uppercase">Goal reached! ★</div>
                      )}
                    </div>
                  </div>
                  {/* Mobile stats bar */}
                  <div className="hidden max-md:flex bg-mathua-surface border border-mathua-border rounded-none p-3 mb-4 items-center justify-around text-center gap-1 sm:gap-2 min-w-0 overflow-hidden">
                    <div className="flex-1 min-w-0">
                      <span className="font-mono text-[8px] uppercase text-mathua-muted block truncate">Streak</span>
                      <span className="font-mono text-sm text-mathua-green truncate block">{submitted && lastResult ? lastResult.streak : lastResult ? lastResult.streak : '--'}</span>
                    </div>
                    <div className="w-px h-8 bg-mathua-border shrink-0" />
                    <div className="flex-1 min-w-0">
                      <span className="font-mono text-[8px] uppercase text-mathua-muted block truncate">Session</span>
                      <span className={`font-mono text-sm truncate block ${sessionStats.total > 0 && sessionStats.correct / sessionStats.total >= 0.8 ? 'text-mathua-green' : sessionStats.total > 0 && sessionStats.correct / sessionStats.total < 0.5 ? 'text-mathua-red' : 'text-mathua-primary'}`}>
                        {sessionStats.correct}/{sessionStats.total}
                      </span>
                    </div>
                    <div className="w-px h-8 bg-mathua-border shrink-0" />
                    <div className="flex-1 min-w-0">
                      <span className="font-mono text-[8px] uppercase text-mathua-muted block truncate">Mastered</span>
                      <span className="font-mono text-sm text-mathua-blue truncate block">{scores.concepts_mastered}</span>
                    </div>
                    <div className="w-px h-8 bg-mathua-border shrink-0" />
                    <div className="flex-1 min-w-0">
                      <span className="font-mono text-[8px] uppercase text-mathua-muted block truncate">XP today</span>
                      <span className="font-mono text-sm text-mathua-blue truncate block">{scores.xp_today ?? 0}</span>
                    </div>
                  </div>
                </div>
                <div className="flex-1 min-w-0">
                  <div className={`bg-mathua-surface border rounded-none p-8 max-md:p-5 transition-colors duration-300 ${
                    submitted && lastResult
                      ? lastResult.correct ? 'border-green-500/40' : 'border-red-500/40'
                      : 'border-mathua-border'
                  }`}>
                    <div className="mb-6">
                      <span className="section-label">{question.concept_id}</span>
                      <h3 className="font-serif text-2xl font-medium text-mathua-primary mt-1">
                        {question.concept_name}
                      </h3>
                    </div>

                    <div className={`bg-mathua-code border border-mathua-border rounded-none mb-6 ${question.diagram ? 'p-0' : 'p-8 text-center'}`}>
                      {question.diagram ? (
                         <div className="flex flex-col md:flex-row">
                           <div className="md:w-1/3 p-4 flex items-center justify-center bg-mathua-surface border-r border-mathua-border">
                             <Image src={question.diagram} alt="Diagram" width={200} height={180} className="max-w-full h-auto" style={{ maxHeight: '180px' }} unoptimized />
                           </div>
                          <div className="md:w-2/3 p-8 flex items-center justify-center">
                               <KatexContent className="text-mathua-primary text-2xl font-mono font-light whitespace-pre-wrap text-center">
                                 {question.question}
                               </KatexContent>
                             </div>
                           </div>
                          ) : (
                           <div className="">
                             <KatexContent className="text-mathua-primary text-2xl font-mono font-light whitespace-pre-wrap p-8">
                               {question.question}
                             </KatexContent>
                           </div>
                          )}
                       </div>

                       {!submitted ? (
                         <>
                           <div className="flex flex-col sm:flex-row gap-3 mb-4 min-w-0">
                             <input
                               ref={inputRef}
                              type="text"
                              value={answer}
                              onChange={(e) => setAnswer(e.target.value)}
                              onKeyDown={(e) => e.key === 'Enter' && handleSubmit()}
                             placeholder="Your answer"
                             className="flex-1 min-w-0 bg-mathua-code border border-mathua-border rounded-none h-12 px-4 font-mono text-base text-mathua-primary placeholder:text-mathua-muted focus:outline-none focus:border-mathua-blue"
                          />
                           <button
                             onClick={handleSubmit}
                             disabled={!answer.trim() || loading || submitted || submittingRef.current}
                             className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white rounded-none h-12 px-8 font-medium text-sm whitespace-nowrap shrink-0 w-full sm:w-auto min-h-[36px] disabled:opacity-50"
                           >
                             Check Answer
                           </button>
                         </div>
                         <SymbolPalette targetRef={inputRef} onInsert={setAnswer} />
                         {question.lesson && (
                          <details className="mt-2">
                            <summary className="text-mathua-secondary text-sm cursor-pointer hover:text-mathua-blue">
                              Show lesson: {question.lesson.Title}
                            </summary>
                            <div className="mt-2 bg-mathua-code border border-mathua-border rounded-none p-4 text-xs text-mathua-muted max-h-80 overflow-auto leading-relaxed">
                              <KatexContent>{question.lesson.Body.slice(0, 3000)}</KatexContent>
                            </div>
                          </details>
                        )}
                      </>
                    ) : (
                      <div className="space-y-4 animate-fadeIn">
                        <div className={`border-t pt-4 ${lastResult?.correct ? 'border-green-500/20' : 'border-red-500/20'}`}>
                          <div className="flex items-center gap-3 mb-3">
                            <span className={`text-2xl ${lastResult?.correct ? 'text-green-400' : 'text-red-400'}`}>
                              {lastResult?.correct ? '✓' : '✗'}
                            </span>
                            <div>
                              <p className={`font-serif text-lg font-medium ${lastResult?.correct ? 'text-green-400' : 'text-red-400'}`}>
                                {lastResult?.correct ? 'Correct!' : 'Incorrect'}
                              </p>
                              {lastResult?.explanation && (
                                <KatexContent className="text-mathua-secondary text-xs mt-1">{lastResult.explanation}</KatexContent>
                              )}
                            </div>
                          </div>
                          <div className="flex flex-wrap items-center gap-x-4 gap-y-1 text-xs font-mono text-mathua-muted mb-4">
                            <span>Streak: {lastResult?.streak ?? 0}/{lastResult?.required_streak ?? 0}</span>
                            <span>Status: {lastResult?.new_status ?? 'UNSEEN'}{lastResult?.new_status === 'MASTERED' ? ' ★' : ''}</span>
                            {lastResult && lastResult.xp > 0 && (
                              <span className="text-yellow-400">+{lastResult.xp} XP</span>
                            )}
                          </div>
                          {question ? (
                            <button
                              onClick={nextQuestion}
                              className="w-full border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white rounded-none h-12 font-medium text-sm transition-colors"
                            >
                              Next Question
                            </button>
                          ) : (
                            <p className="text-mathua-blue text-center text-sm">
                              All available concepts mastered! Come back tomorrow for reviews.
                            </p>
                          )}
                        </div>
                      </div>
                    )}
                  </div>
                </div>
              </div>
            )}

            {submitted && !question && lastResult && (
              <div className="max-w-xl mx-auto mt-8">
                <div className={`bg-mathua-surface border rounded-none p-8 max-md:p-5 ${
                  lastResult.correct ? 'border-mathua-green' : 'border-mathua-red'
                }`}>
                  <h3 className={`font-serif text-2xl font-medium mb-2 ${
                    lastResult.correct ? 'text-mathua-green' : 'text-mathua-red'
                  }`}>
                    {lastResult.correct ? 'Correct!' : 'Incorrect'}
                  </h3>
                  <p className="text-mathua-blue text-center mt-4">
                    All available concepts mastered! Come back tomorrow for reviews.
                  </p>
                </div>
              </div>
            )}
          </>
        )}

        {screen === 'review' && (
          <>
            {reviewDone ? (
              <div className="max-w-xl mx-auto mt-8">
                <div className="bg-mathua-surface border border-mathua-border rounded-none p-8 max-md:p-5 text-center">
                  <h3 className="font-serif text-2xl font-medium text-mathua-green mb-2">
                    Review Complete!
                  </h3>
                  <p className="font-mono text-sm text-mathua-secondary mb-4">
                    You reviewed {reviewStats.total} concept{reviewStats.total !== 1 ? 's' : ''}
                    {reviewStats.total > 0 && (
                      <> · {reviewStats.correct}/{reviewStats.total} correct</>
                    )}
                  </p>
                  <button
                    onClick={nextReviewQuestion}
                    className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white rounded-none h-12 px-8 font-medium text-sm transition-colors"
                  >
                    Back to Practice
                  </button>
                </div>
              </div>
            ) : !question ? (
              <div className="max-w-xl mx-auto mt-20 text-center">
                <p className="text-mathua-muted text-sm mb-4">No reviews due.</p>
                <button
                  onClick={() => setScreen('practice')}
                  className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white rounded-none h-12 px-8 font-medium text-sm"
                >
                  Back to Practice
                </button>
              </div>
            ) : (
              <div className="flex gap-8 items-start mt-4 max-md:flex-col">
                <div className="w-[200px] flex-shrink-0 max-md:w-full">
                  <div className="bg-mathua-surface border border-mathua-border rounded-none p-5 space-y-4 max-md:hidden">
                    <div>
                      <span className="font-mono text-[10px] uppercase text-mathua-muted">Review</span>
                      <div className="font-mono text-lg text-mathua-primary mt-1">
                        <span className={reviewStats.total > 0 && reviewStats.correct / reviewStats.total >= 0.8 ? 'text-mathua-green' : reviewStats.total > 0 && reviewStats.correct / reviewStats.total < 0.5 ? 'text-mathua-red' : 'text-mathua-primary'}>
                          {reviewStats.correct}/{reviewStats.total}
                        </span>
                        <span className="text-mathua-muted text-xs ml-1">correct</span>
                      </div>
                      {reviewStats.total > 0 && (
                        <div className="mt-1 h-1 bg-mathua-code rounded-full overflow-hidden">
                          <div
                            className="h-full bg-yellow-500 rounded-full transition-all duration-500"
                            style={{ width: `${(reviewStats.correct / reviewStats.total) * 100}%` }}
                          />
                        </div>
                      )}
                    </div>
                    <div>
                      <span className="font-mono text-[10px] uppercase text-mathua-muted">Streak</span>
                      <div className="font-mono text-2xl text-mathua-green mt-1">
                        {submitted && lastResult ? lastResult.streak : lastResult ? lastResult.streak : '--'}
                      </div>
                    </div>
                    <div>
                      <span className="font-mono text-[10px] uppercase text-mathua-muted">Mastered total</span>
                      <div className="font-mono text-2xl text-mathua-blue mt-1">{scores.concepts_mastered}</div>
                    </div>
                    <div>
                      <span className="font-mono text-[10px] uppercase text-mathua-muted">Daily XP</span>
                      <div className="font-mono text-2xl text-mathua-blue mt-1">{scores.xp_today ?? 0}</div>
                    </div>
                  </div>
                  <div className="hidden max-md:flex bg-mathua-surface border border-mathua-border rounded-none p-3 mb-4 items-center justify-around text-center gap-1 sm:gap-2 min-w-0 overflow-hidden">
                    <div className="flex-1 min-w-0">
                      <span className="font-mono text-[8px] uppercase text-mathua-muted block truncate">Review</span>
                      <span className={`font-mono text-sm truncate block ${reviewStats.total > 0 && reviewStats.correct / reviewStats.total >= 0.8 ? 'text-mathua-green' : reviewStats.total > 0 && reviewStats.correct / reviewStats.total < 0.5 ? 'text-mathua-red' : 'text-mathua-primary'}`}>
                        {reviewStats.correct}/{reviewStats.total}
                      </span>
                    </div>
                    <div className="w-px h-8 bg-mathua-border shrink-0" />
                    <div className="flex-1 min-w-0">
                      <span className="font-mono text-[8px] uppercase text-mathua-muted block truncate">Streak</span>
                      <span className="font-mono text-sm text-mathua-green truncate block">{submitted && lastResult ? lastResult.streak : lastResult ? lastResult.streak : '--'}</span>
                    </div>
                    <div className="w-px h-8 bg-mathua-border shrink-0" />
                    <div className="flex-1 min-w-0">
                      <span className="font-mono text-[8px] uppercase text-mathua-muted block truncate">Mastered</span>
                      <span className="font-mono text-sm text-mathua-blue truncate block">{scores.concepts_mastered}</span>
                    </div>
                  </div>
                </div>
                <div className="flex-1 min-w-0">
                  <div className={`bg-mathua-surface border rounded-none p-8 max-md:p-5 transition-colors duration-300 ${
                    submitted && lastResult
                      ? lastResult.correct ? 'border-green-500/40' : 'border-red-500/40'
                      : 'border-yellow-500/30'
                  }`}>
                    <div className="mb-4 flex items-center gap-2">
                      <span className="font-mono text-[10px] uppercase text-yellow-400">Review</span>
                      <span className="text-mathua-muted text-[10px]">·</span>
                      <span className="section-label">{question.concept_id}</span>
                    </div>
                      <h3 className="font-serif text-2xl font-medium text-mathua-primary mt-1 mb-6">
                       {question.concept_name}
                     </h3>
                     <div className={`bg-mathua-code border border-mathua-border rounded-none mb-6 ${question.diagram ? 'p-0' : 'p-8 text-center'}`}>
                       {question.diagram ? (
                         <div className="flex flex-col md:flex-row">
                           <div className="md:w-1/3 p-4 flex items-center justify-center bg-mathua-surface border-r border-mathua-border">
                             <Image src={question.diagram} alt="Diagram" width={200} height={180} className="max-w-full h-auto" style={{ maxHeight: '180px' }} unoptimized />
                           </div>
                             <div className="md:w-2/3 p-8 flex items-center justify-center">
                               <KatexContent className="text-mathua-primary text-2xl font-mono font-light whitespace-pre-wrap text-center">
                                 {question.question}
                               </KatexContent>
                             </div>
                           </div>
                          ) : (
                           <div className="">
                             <KatexContent className="text-mathua-primary text-2xl font-mono font-light whitespace-pre-wrap p-8">
                               {question.question}
                             </KatexContent>
                           </div>
                          )}
                       </div>
                        {!submitted ? (
                          <>
                            <div className="flex flex-col sm:flex-row gap-3 mb-4 min-w-0">
                              <input
                                ref={inputRef}
                               type="text"
                               value={answer}
                               onChange={(e) => setAnswer(e.target.value)}
                               onKeyDown={(e) => e.key === 'Enter' && handleReviewSubmit()}
                              placeholder="Your answer"
                              className="flex-1 min-w-0 bg-mathua-code border border-mathua-border rounded-none h-12 px-4 font-mono text-base text-mathua-primary placeholder:text-mathua-muted focus:outline-none focus:border-yellow-500"
                           />
                           <button
                             onClick={handleReviewSubmit}
                             disabled={!answer.trim() || loading || submitted || submittingRef.current}
                             className="border border-yellow-500/60 text-yellow-400 hover:bg-yellow-500 hover:text-black rounded-none h-12 px-8 font-medium text-sm whitespace-nowrap transition-colors shrink-0 w-full sm:w-auto min-h-[36px] disabled:opacity-50"
                           >
                             Check Answer
                           </button>
                         </div>
                         <SymbolPalette targetRef={inputRef} onInsert={setAnswer} />
                       </>
                    ) : (
                      <div className="space-y-4 animate-fadeIn">
                        <div className={`border-t pt-4 ${lastResult?.correct ? 'border-green-500/20' : 'border-red-500/20'}`}>
                          <div className="flex items-center gap-3 mb-3">
                            <span className={`text-2xl ${lastResult?.correct ? 'text-green-400' : 'text-red-400'}`}>
                              {lastResult?.correct ? '✓' : '✗'}
                            </span>
                            <div>
                              <p className={`font-serif text-lg font-medium ${lastResult?.correct ? 'text-green-400' : 'text-red-400'}`}>
                                {lastResult?.correct ? 'Correct!' : 'Incorrect'}
                              </p>
                              {lastResult?.explanation && (
                                <KatexContent className="text-mathua-secondary text-xs mt-1">{lastResult.explanation}</KatexContent>
                              )}
                            </div>
                          </div>
                          <div className="flex flex-wrap items-center gap-x-4 gap-y-1 text-xs font-mono text-mathua-muted mb-4">
                            <span>Streak: {lastResult?.streak ?? 0}/{lastResult?.required_streak ?? 0}</span>
                            <span>Status: {lastResult?.new_status ?? 'UNSEEN'}{lastResult?.new_status === 'MASTERED' ? ' ★' : ''}</span>
                            {lastResult && lastResult.xp > 0 && (
                              <span className="text-yellow-400">+{lastResult.xp} XP</span>
                            )}
                          </div>
                          {question ? (
                            <button
                              onClick={nextReviewQuestion}
                              className="w-full border border-yellow-500/60 text-yellow-400 hover:bg-yellow-500 hover:text-black rounded-none h-12 font-medium text-sm transition-colors"
                            >
                              Next Review Question
                            </button>
                          ) : (
                            <button
                              onClick={() => { setReviewDone(true); nextReviewQuestion() }}
                              className="w-full border border-mathua-green text-mathua-green hover:bg-mathua-green hover:text-white rounded-none h-12 font-medium text-sm transition-colors"
                            >
                              See Summary
                            </button>
                          )}
                        </div>
                      </div>
                    )}
                  </div>
                </div>
              </div>
            )}
          </>
        )}
      </section>

      <AsciiDivider pattern="wave" />
      <Footer />
    </div>
      <BottomTabs />
    </>
  )
}
