'use client'

import { useTheme } from '../../hooks/useTheme'
import { useState, useEffect, useCallback, useRef } from 'react'
import { useRouter } from 'next/navigation'
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
  type Question,
  type AnswerResult,
  type Scores,
} from '../../lib/api'
import { isLoggedIn, getUserInfo, clearToken, type UserInfo } from '../../lib/auth'

type Screen = 'name' | 'practice' | 'feedback'

export default function SessionPage() {
  const { mounted } = useTheme()
  const router = useRouter()

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
    xp_total: 0, xp_today: 0,
  })
  const [error, setError] = useState('')
  const [elapsed, setElapsed] = useState(0)
  const [loading, setLoading] = useState(true)
  const startRef = useRef(Date.now())
  const timerRef = useRef<ReturnType<typeof setInterval>>()

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

  if (!mounted) return <div style={{ background: 'var(--bg)', minHeight: '100vh' }} />

  return (
    <>
      <Header />
      <div className="max-w-container mx-auto px-6 max-sm:px-4">
      <section className="pt-8">
        <span className="flex justify-between items-center mb-4">
          <a href="/" className="text-mathua-secondary text-sm hover:text-mathua-primary">
            ← Back
          </a>
          <span className="flex gap-3 items-center">
            {user && <span className="font-mono text-[11px] text-mathua-muted">{user.name}</span>}
            <button onClick={() => { clearToken(); router.push('/login') }} className="text-mathua-secondary text-xs hover:text-mathua-red">
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
                autoFocus
                className="flex-1 bg-mathua-code border border-mathua-border rounded-md h-12 px-4 font-mono text-base text-mathua-primary placeholder:text-mathua-muted focus:outline-none focus:border-mathua-blue"
              />
              <button
                onClick={beginSessionName}
                className="bg-mathua-blue text-white hover:bg-mathua-blue-hover rounded-md h-12 px-8 font-medium text-sm"
              >
                Start
              </button>
            </div>
            <div className="mt-4 text-center">
              <a href="/login" className="text-mathua-secondary text-sm hover:text-mathua-blue">
                Have an account? Sign in
              </a>
            </div>
          </div>
        )}

        {loading && (
          <div className="max-w-md mx-auto mt-20 text-center">
            <p className="text-mathua-muted text-sm">Loading...</p>
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
                <div className="bg-mathua-surface border border-mathua-border rounded-lg p-5 space-y-4">
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
                    <div className="font-mono text-2xl text-mathua-gold mt-1">{scores.concepts_mastered}</div>
                  </div>
                  <div>
                    <span className="font-mono text-[10px] uppercase text-mathua-muted">XP</span>
                    <div className="flex items-baseline gap-2 mt-1">
                      <span className="font-mono text-2xl text-mathua-blue">{scores.xp_total ?? 0}</span>
                      <span className="font-mono text-[10px] text-mathua-muted">/ {Math.ceil(((scores.xp_total ?? 0) + 1) / 150) * 150}</span>
                    </div>
                    <div className="mt-2 h-1.5 bg-mathua-code rounded-full overflow-hidden">
                      <div
                        className="h-full bg-mathua-blue rounded-full transition-all duration-500"
                        style={{ width: `${(((scores.xp_total ?? 0) % 150) / 150) * 100}%` }}
                      />
                    </div>
                    {(scores.xp_total ?? 0) > 0 && (scores.xp_total ?? 0) % 150 === 0 && (
                      <div className="mt-1 font-mono text-[10px] text-mathua-gold uppercase">Quiz ready!</div>
                    )}
                  </div>
                </div>
              </div>
              <div className="flex-1">
                <div className="bg-mathua-surface border border-mathua-border rounded-lg p-8 max-md:p-5">
                  <div className="mb-6">
                    <span className="section-label">{question.concept_id}</span>
                    <h3 className="font-serif text-2xl font-medium text-mathua-primary mt-1">
                      {question.concept_name}
                    </h3>
                  </div>
                  <div className="bg-mathua-code border border-mathua-border rounded-md p-8 text-center mb-6">
                    <p className="text-mathua-primary text-2xl font-mono font-light whitespace-pre-wrap">
                      {question.question}
                    </p>
                  </div>
                  <div className="flex gap-3 mb-4">
                    <input
                      type="text"
                      value={answer}
                      onChange={(e) => setAnswer(e.target.value)}
                      onKeyDown={(e) => e.key === 'Enter' && handleSubmit()}
                      placeholder="Your answer"
                      autoFocus
                      className="flex-1 bg-mathua-code border border-mathua-border rounded-md h-12 px-4 font-mono text-base text-mathua-primary placeholder:text-mathua-muted focus:outline-none focus:border-mathua-blue"
                    />
                    <button
                      onClick={handleSubmit}
                      className="bg-mathua-blue text-white hover:bg-mathua-blue-hover rounded-md h-12 px-8 font-medium text-sm"
                    >
                      Submit
                    </button>
                  </div>
                  {question.lesson && (
                    <details className="mt-2">
                      <summary className="text-mathua-secondary text-sm cursor-pointer hover:text-mathua-blue">
                        Show lesson: {question.lesson.Title}
                      </summary>
                      <pre className="mt-2 bg-mathua-code p-4 rounded text-xs text-mathua-muted whitespace-pre-wrap max-h-60 overflow-auto">
                        {question.lesson.Body.slice(0, 2000)}
                      </pre>
                    </details>
                  )}
                </div>
              </div>
            </div>
          </>
        )}

        {screen === 'feedback' && lastResult && (
          <div className="max-w-xl mx-auto mt-8">
            <div className={`bg-mathua-surface border rounded-lg p-8 max-md:p-5 ${
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
                  className="w-full bg-mathua-blue text-white hover:bg-mathua-blue-hover rounded-md h-12 font-medium text-sm"
                >
                  Next Question
                </button>
              ) : (
                <p className="text-mathua-gold text-center">
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
