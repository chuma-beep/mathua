'use client'

import { useTheme } from '../../hooks/useTheme'
import { useState, useEffect, useCallback, useRef } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import Image from 'next/image'
import Header from '../../components/Header'
import SectionHeader from '../../components/SectionHeader'
import Footer from '../../components/Footer'
import AsciiDivider from '../../components/AsciiDivider'
import Pipeline from '../../components/Pipeline'
import {
  startSession,
  startSessionName,
  submitAnswer,
  getScores,
  startGoalDiagnosticName,
  getGoalPlan,
  setDailyXPGoal,
  type Question,
  type AnswerResult,
  type Scores,
  type GoalPlanRes,
} from '../../lib/api'
import { isLoggedIn, getUserInfo, clearToken, type UserInfo } from '../../lib/auth'
import conceptsData from '../../data/concepts.json'

type Screen = 'name' | 'diag_select' | 'diagnostic' | 'practice' | 'feedback'

export default function SessionPage() {
  const { mounted } = useTheme()
  const { push } = useRouter()

  const [screen, setScreen] = useState<Screen>('practice')
  const [name, setName] = useState('')
  const [user, setUser] = useState<UserInfo | null>(null)
  const [studentID, setStudentID] = useState('')
  const [sessionID, setSessionID] = useState('')
  const [question, setQuestion] = useState<Question | null>(null)
  const [lastResult, setLastResult] = useState<AnswerResult | null>(null)
  const [answer, setAnswer] = useState('')
  const [scores, setScores] = useState<Scores>({
    lifetime_points: 0, weekly_score: 0, speed_bonus: 0,
    concepts_mastered: 0, current_streak: 0, level: 'Novice',
    xp_total: 0, xp_today: 0, daily_xp_goal: 150,
  })
  const [error, setError] = useState('')
  const [elapsed, setElapsed] = useState(0)
  const [loading, setLoading] = useState(true)
  const startRef = useRef(Date.now())
  const timerRef = useRef<ReturnType<typeof setInterval>>()

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
  const diagPlan = useRef<GoalPlanRes | null>(null)
  const guestStudentID = useRef('')

  useEffect(() => {
    const u = getUserInfo()
    if (u) {
      setUser(u)
      setScreen('practice')
      beginSessionAuth()
    } else {
      setScreen('name')
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    if (screen === 'practice') {
      startRef.current = Date.now()
      timerRef.current = setInterval(() => {
        setElapsed((Date.now() - startRef.current) / 1000)
      }, 100)
    }
    return () => { if (timerRef.current) clearInterval(timerRef.current) }
  }, [screen, question])

  const beginSessionAuth = useCallback(async () => {
    setError('')
    setLoading(true)
    try {
      const res = await startSession()
      setStudentID(res.student_id)
      setSessionID(res.session_id)
      setQuestion(res.question)
      setScreen('practice')
      const s = await getScores(res.student_id).catch(() => null)
      if (s) setScores(s)
    } catch {
      setError('Could not connect to server. Is the backend running?')
    } finally {
      setLoading(false)
    }
  }, [])

  const beginSessionName = useCallback(async () => {
    if (!name.trim()) { setError('Enter your name'); return }
    setError(''); setLoading(true)
    try {
      const res = await startSessionName(name.trim())
      setStudentID(res.student_id)
      setSessionID(res.session_id)
      setQuestion(res.question)
      setScreen('practice')
      const s = await getScores(res.student_id).catch(() => null)
      if (s) setScores(s)
    } catch {
      setError('Could not connect to server. Is the backend running?')
    } finally { setLoading(false) }
  }, [name])

  const handleSubmit = useCallback(async () => {
    if (!answer.trim() || !question) return
    const e = (Date.now() - startRef.current) / 1000
    try {
      const res = await submitAnswer(sessionID, answer.trim(), e)
      setLastResult(res.result)
      setAnswer('')
      if (timerRef.current) clearInterval(timerRef.current)
      setScreen('feedback')
      // Pre-load next question
      if (res.next_question) {
        setQuestion(res.next_question)
      } else {
        setQuestion(null)
      }
      const s = await getScores(studentID).catch(() => null)
      if (s) setScores(s)
    } catch {
      setError('Failed to submit answer')
    }
  }, [answer, question, sessionID, studentID])

  const nextQuestion = useCallback(() => {
    setLastResult(null)
    setScreen('practice')
  }, [])

  // Guest diagnostic logic
  const domainOrder = ['counting', 'arithmetic', 'fractions', 'prealgebra', 'algebra', 'geometry', 'trigonometry', 'complex_numbers', 'precalculus', 'calculus', 'linear_algebra', 'statistics', 'discrete_math', 'number_theory', 'differential_equations', 'abstract_algebra', 'topology']
  const domainLabels: Record<string, string> = {
    counting: 'Counting', arithmetic: 'Arithmetic', fractions: 'Fractions', prealgebra: 'Pre-Algebra',
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
        diagPlan.current = { readiness: 1, total_tested: 0, correct_count: 0, weak_areas: {}, strong_areas: {} }
        setScreen('practice')
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
      setError('Could not start diagnostic.')
    } finally {
      setLoading(false)
    }
  }

  async function submitGuestDiagnostic() {
    if (!diagAnswer.trim()) return
    setLoading(true)
    try {
      const elapsed = 5.0
      const res = await fetch(`/api/goal/diagnostic/answer`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ session_id: diagSessionId.current, concept_id: diagConceptId.current, answer: diagAnswer.trim(), elapsed }),
      })
      const data = await res.json()
      const correct = data.correct || false
      const feedback = data.feedback || (correct ? 'Correct!' : 'Not quite.')
      setDiagAccuracy(prev => ({ correct: prev.correct + (correct ? 1 : 0), total: prev.total + 1 }))
      setDiagLastResult({ correct, feedback })

      if (data.done) {
        setTimeout(async () => {
          try {
            const planRes = await getGoalPlan(diagSessionId.current)
            diagPlan.current = planRes
            if (guestStudentID.current) {
              const sessRes = await fetch('/api/session', {
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
            }
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
      <div className="max-w-container mx-auto px-6 max-sm:px-4">
      <section className="pt-8">
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
            <div className="flex gap-3 mt-6">
              <input
                type="text"
                value={name}
                onChange={(e) => setName(e.target.value)}
                onKeyDown={(e) => e.key === 'Enter' && beginSessionName()}
                placeholder="Your name"
                className="flex-1 bg-mathua-code border border-mathua-border rounded-none h-12 px-4 font-mono text-base text-mathua-primary placeholder:text-mathua-muted focus:outline-none focus:border-mathua-blue"
              />
              <button
                onClick={startGuestDiagnostic}
                className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white rounded-none h-12 px-8 font-medium text-sm"
              >
                Start
              </button>
            </div>
            <div className="mt-4 text-center">
              <Link href="/login" className="text-mathua-secondary text-sm hover:text-mathua-blue">
                Have an account? Sign in
              </Link>
            </div>
          </div>
        )}

        {screen === 'diag_select' && (
          <div className="max-w-4xl mx-auto mt-8">
            <SectionHeader label="Diagnostic" title="What do you want to learn?" />
            <p className="text-mathua-secondary text-sm text-center max-w-[500px] mx-auto mt-2 mb-6">
              We'll test your current knowledge and find the right starting point.
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
                {loading ? 'Loading…' : 'Start diagnostic'}
              </button>
            </div>
          </div>
        )}

        {screen === 'diagnostic' && (
          <div className="max-w-2xl mx-auto mt-8">
            <SectionHeader label={`Question ${diagCount}`} title={diagConceptName} />
            <div className="bg-mathua-surface border border-mathua-border rounded-none p-6 mb-6">
              <div className="bg-mathua-code border border-mathua-border rounded-none p-6 text-center mb-4">
                <p className="text-mathua-primary text-lg font-mono font-light whitespace-pre-wrap">{diagQuestion}</p>
              </div>
              <div className="flex gap-3">
                <input
                  type="text"
                  value={diagAnswer}
                  onChange={(e) => setDiagAnswer(e.target.value)}
                  onKeyDown={(e) => e.key === 'Enter' && submitGuestDiagnostic()}
                  placeholder="Your answer..."
                  disabled={loading || diagLastResult !== null}
                  className="flex-1 bg-mathua-code border border-mathua-border rounded-none h-12 px-4 font-mono text-base text-mathua-primary placeholder:text-mathua-muted focus:outline-none focus:border-mathua-blue"
                />
                <button onClick={submitGuestDiagnostic} disabled={!diagAnswer.trim() || loading || diagLastResult !== null} className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white rounded-none h-12 px-8 font-medium text-sm disabled:opacity-50">
                  Check Answer
                </button>
              </div>
            </div>
            {diagLastResult && (
              <div className={`bg-mathua-surface border rounded-none p-4 mb-4 text-center ${diagLastResult.correct ? 'border-mathua-green' : 'border-mathua-red'}`}>
                <p className={diagLastResult.correct ? 'text-mathua-green' : 'text-mathua-red'}>{diagLastResult.feedback}</p>
              </div>
            )}
            <div className="text-center text-mathua-muted text-xs font-mono">
              {diagAccuracy.total > 0 && `${diagAccuracy.correct}/${diagAccuracy.total} correct`}
            </div>
          </div>
        )}

        {loading && (
          <div className="max-w-md mx-auto mt-20 text-center">
            <p className="text-mathua-muted text-sm">Loading…</p>
          </div>
        )}
        {error && (
          <div className="max-w-md mx-auto mt-20 text-center">
            <p className="text-mathua-red text-sm">{error}</p>
          </div>
        )}

        {screen === 'practice' && question && (
          <>
            <div className="flex gap-8 items-start mt-4 max-md:flex-col">
              <div className="w-[200px] flex-shrink-0 max-md:w-full">
                <div className="bg-mathua-surface border border-mathua-border rounded-none p-5 space-y-4">
                  <div>
                    <span className="font-mono text-[10px] uppercase text-mathua-muted">Streak</span>
                    <div className="font-mono text-2xl text-mathua-green mt-1">
                      {lastResult ? lastResult.streak : '--'}
                    </div>
                  </div>
                  <div>
                    <span className="font-mono text-[10px] uppercase text-mathua-muted">Time</span>
                    <div className="font-mono text-2xl text-mathua-primary mt-1">{elapsed.toFixed(1)}s</div>
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
                              }).catch(() => {})
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
              </div>
              <div className="flex-1">
                <div className="bg-mathua-surface border border-mathua-border rounded-none p-8 max-md:p-5">
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
                          <p className="text-mathua-primary text-2xl font-mono font-light whitespace-pre-wrap text-center">
                            {question.question}
                          </p>
                        </div>
                      </div>
                    ) : (
                      <p className="text-mathua-primary text-2xl font-mono font-light whitespace-pre-wrap p-8">
                        {question.question}
                      </p>
                    )}
                  </div>
                  <div className="flex gap-3 mb-4">
                    <input
                      type="text"
                      value={answer}
                      onChange={(e) => setAnswer(e.target.value)}
                      onKeyDown={(e) => e.key === 'Enter' && handleSubmit()}
                      placeholder="Your answer"
                      className="flex-1 bg-mathua-code border border-mathua-border rounded-none h-12 px-4 font-mono text-base text-mathua-primary placeholder:text-mathua-muted focus:outline-none focus:border-mathua-blue"
                    />
                    <button
                      onClick={handleSubmit}
                      className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white rounded-none h-12 px-8 font-medium text-sm"
                    >
                      Check Answer
                    </button>
                  </div>
                  {question.lesson && (
                    <details className="mt-2">
                      <summary className="text-mathua-secondary text-sm cursor-pointer hover:text-mathua-blue">
                        Show lesson: {question.lesson.Title}
                      </summary>
                      <div className="mt-2 bg-mathua-code border border-mathua-border rounded-none p-4 text-xs text-mathua-muted whitespace-pre-wrap max-h-80 overflow-auto leading-relaxed">
                        {question.lesson.Body.slice(0, 3000)}
                      </div>
                    </details>
                  )}
                </div>
              </div>
            </div>
          </>
        )}

        {screen === 'feedback' && lastResult && (
          <div className="max-w-xl mx-auto mt-8">
            <div className={`bg-mathua-surface border rounded-none p-8 max-md:p-5 ${
              lastResult.correct ? 'border-mathua-green' : 'border-mathua-red'
            }`}>
              <h3 className={`font-serif text-2xl font-medium mb-2 ${
                lastResult.correct ? 'text-mathua-green' : 'text-mathua-red'
              }`}>
                {lastResult.correct ? 'Correct!' : 'Incorrect'}
              </h3>
              {lastResult.explanation && (
                <p className="text-mathua-secondary text-sm mb-2">
                  {lastResult.explanation}
                </p>
              )}
              <p className="text-mathua-muted text-xs mb-6">
                Streak: {lastResult.streak}/{lastResult.required_streak} ·
                Status: {lastResult.new_status}
                {lastResult.new_status === 'MASTERED' && ' ★'}
                {lastResult.xp > 0 && ` ·  +${lastResult.xp} XP`}
              </p>
              {question ? (
                <button
                  onClick={nextQuestion}
                  className="w-full border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white rounded-none h-12 font-medium text-sm"
                >
                  Next Question
                </button>
              ) : (
                <p className="text-mathua-blue text-center">
                  All available concepts mastered! Come back tomorrow for reviews.
                </p>
              )}
            </div>
          </div>
        )}
      </section>

      <AsciiDivider pattern="wave" />

      <section className="py-20 max-sm:py-12">
        <SectionHeader label="Your state" title="Mastery pipeline" />
        <div className="flex justify-center mt-6">
          <Pipeline
            states={[
              { label: 'UNSEEN', status: 'unseen' },
              { label: 'LEARNING', status: 'learning' },
              { label: 'PRACTICING', status: 'practicing' },
              { label: 'MASTERED', status: 'mastered' },
              { label: 'DECAYING', status: 'decaying' },
            ]}
          />
        </div>
      </section>

      <Footer />
    </div>
    </>
  )
}
