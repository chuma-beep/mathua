'use client'

import { useTheme } from '../../hooks/useTheme'
import { useState, useEffect, useCallback, useRef } from 'react'
import Header from '../../components/Header'
import SectionHeader from '../../components/SectionHeader'
import Footer from '../../components/Footer'
import AsciiDivider from '../../components/AsciiDivider'
import Pipeline from '../../components/Pipeline'
import {
  startSession,
  submitAnswer,
  getScores,
  type Question,
  type AnswerResult,
  type Scores,
} from '../../lib/api'

type Screen = 'name' | 'practice' | 'feedback'

export default function SessionPage() {
  const { mounted } = useTheme()

  const [screen, setScreen] = useState<Screen>('name')
  const [name, setName] = useState('')
  const [studentID, setStudentID] = useState('')
  const [sessionID, setSessionID] = useState('')
  const [question, setQuestion] = useState<Question | null>(null)
  const [lastResult, setLastResult] = useState<AnswerResult | null>(null)
  const [answer, setAnswer] = useState('')
  const [scores, setScores] = useState<Scores>({
    lifetime_points: 0, weekly_score: 0, speed_bonus: 0,
    concepts_mastered: 0, current_streak: 0, level: 'Novice',
  })
  const [error, setError] = useState('')
  const [connected, setConnected] = useState<boolean | null>(null)
  const [elapsed, setElapsed] = useState(0)
  const startRef = useRef(Date.now())
  const timerRef = useRef<ReturnType<typeof setInterval>>()

  useEffect(() => {
    if (screen === 'practice') {
      startRef.current = Date.now()
      timerRef.current = setInterval(() => {
        setElapsed((Date.now() - startRef.current) / 1000)
      }, 100)
    }
    return () => { if (timerRef.current) clearInterval(timerRef.current) }
  }, [screen, question])

  const beginSession = useCallback(async () => {
    if (!name.trim()) {
      setError('Enter your name')
      return
    }
    setError('')
    try {
      const res = await startSession(name.trim())
      setStudentID(res.student_id)
      setSessionID(res.session_id)
      setQuestion(res.question)
      setScreen('practice')
      const s = await getScores(res.student_id).catch(() => null)
      if (s) setScores(s)
    } catch {
      setError('Could not connect to server. Is the backend running?')
    }
  }, [name])

  const handleSubmit = useCallback(async () => {
    if (!answer.trim() || !question) return
    const e = (Date.now() - startRef.current) / 1000
    try {
      const res = await submitAnswer(sessionID, studentID, answer.trim(), e)
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
          <span className="font-mono text-[11px] text-mathua-muted">
            {screen === 'name' ? 'New session' : `Session`}
          </span>
        </span>

        {screen === 'name' && (
          <div className="max-w-md mx-auto mt-20">
            <SectionHeader label="Practice session" title="What's your name?" />
            <div className="flex gap-3 mt-6">
              <input
                type="text"
                value={name}
                onChange={(e) => setName(e.target.value)}
                onKeyDown={(e) => e.key === 'Enter' && beginSession()}
                placeholder="Your name"
                autoFocus
                className="flex-1 bg-mathua-code border border-mathua-border rounded-md h-12 px-4 font-mono text-base text-mathua-primary placeholder:text-mathua-muted focus:outline-none focus:border-mathua-blue"
              />
              <button
                onClick={beginSession}
                className="bg-mathua-blue text-white hover:bg-mathua-blue-hover rounded-md h-12 px-8 font-medium text-sm"
              >
                Start
              </button>
            </div>
            {error && <p className="text-mathua-red text-sm mt-2">{error}</p>}
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
