'use client'

import { useState, useEffect, useRef } from 'react'
import { toast } from 'sonner'
import KatexContent from '../../components/KatexContent'
import Loading from '../../components/Loading'
import SectionHeader from '../../components/SectionHeader'
import ProgressBar from '../../components/ProgressBar'
import SymbolPalette from '../../components/SymbolPalette'
import ReportButton from '../../components/ReportButton'
import {
  startGoalDiagnostic,
  submitGoalAnswer,
  resumeGoalDiagnostic,
  getGoalPlan,
  type GoalPlanRes,
  type DiagnosticProgress,
} from '../../lib/api'
import { GOALS_DIAG_KEY } from './constants'

const EMPTY_PLAN: GoalPlanRes = {
  readiness: 1,
  total_tested: 0,
  correct_count: 0,
  weak_areas: {},
  strong_areas: {},
}

// Goals diagnostic (CAT): select → answer loop → plan. Self-contained so the
// goals page only mounts it; owns its session + resume key.
export default function DiagnosticHost({
  startIds,
  resumeId,
  onComplete,
  onResumeExpired,
}: {
  startIds: string[] | null
  resumeId: string | null
  onComplete: (plan: GoalPlanRes) => void
  onResumeExpired: () => void
}) {
  const [loading, setLoading] = useState(false)
  const sessionId = useRef('')
  const conceptId = useRef('')
  const goalsInputRef = useRef<HTMLInputElement>(null)
  const startedRef = useRef(false)

  const [question, setQuestion] = useState('')
  const [conceptName, setConceptName] = useState('')
  const [questionCount, setQuestionCount] = useState(0)
  const [estimatedTotal, setEstimatedTotal] = useState(0)
  const [progress, setProgress] = useState<DiagnosticProgress | null>(null)
  const [answerInput, setAnswerInput] = useState('')
  const [lastResult, setLastResult] = useState<{ correct: boolean; feedback: string } | null>(null)
  const [accuracy, setAccuracy] = useState<{ correct: number; total: number }>({ correct: 0, total: 0 })

  async function start() {
    setLoading(true)
    try {
      const ids = startIds ?? []
      const res = await startGoalDiagnostic(ids)
      if (res.done) {
        onComplete(EMPTY_PLAN)
        return
      }
      sessionId.current = res.session_id
      try {
        sessionStorage.setItem(GOALS_DIAG_KEY, res.session_id)
      } catch { /* storage unavailable — session simply won't resume */ }
      setQuestion(res.question || '')
      conceptId.current = res.concept_id || ''
      setConceptName(res.concept_name || '')
      setQuestionCount(1)
      // Backend truth first; frontend estimate as fallback for older servers.
      if (res.progress && res.progress.cover_size > 0) {
        setProgress(res.progress)
        setEstimatedTotal(res.progress.cover_size)
      } else {
        // Fallback cover size: ~10 + log2(selected concepts)
        const est = Math.min(10 + Math.ceil(Math.log2(Math.max(ids.length, 1)) * 5), 50)
        setEstimatedTotal(est)
        setProgress(null)
      }
      setAccuracy({ correct: 0, total: 0 })
      setLastResult(null)
      setAnswerInput('')
    } catch {
      toast.error("Something went wrong, but we're working on it.")
      onResumeExpired()
    } finally {
      setLoading(false)
    }
  }

  async function resume() {
    const sid = resumeId ?? ''
    if (!sid) {
      onResumeExpired()
      return
    }
    setLoading(true)
    try {
      const data = await resumeGoalDiagnostic(sid)
      if (data.done) {
        try {
          sessionStorage.removeItem(GOALS_DIAG_KEY)
        } catch { /* ignore */ }
        toast.error('That diagnostic already finished — start a fresh one below.')
        onResumeExpired()
        return
      }
      sessionId.current = sid
      setQuestion(data.question || '')
      conceptId.current = data.concept_id || ''
      setConceptName(data.concept_name || '')
      if (data.progress) {
        setProgress(data.progress)
        setEstimatedTotal(data.progress.cover_size)
        setQuestionCount(data.progress.answered + 1)
      }
      setLastResult(null)
      setAnswerInput('')
    } catch {
      try {
        sessionStorage.removeItem(GOALS_DIAG_KEY)
      } catch { /* ignore */ }
      toast.error('Could not resume — that session expired. Start a fresh diagnostic.')
      onResumeExpired()
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    if (startedRef.current) return
    startedRef.current = true
    if (resumeId) void resume()
    else void start()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

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
      if (data.progress && data.progress.cover_size > 0) {
        setProgress(data.progress)
        setEstimatedTotal(data.progress.cover_size)
      }

      if (data.done) {
        setTimeout(async () => {
          try {
            const planRes = await getGoalPlan(sessionId.current)
            try {
              sessionStorage.removeItem(GOALS_DIAG_KEY)
            } catch { /* ignore */ }
            onComplete(planRes)
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

  return (
    <>
      <SectionHeader label={`Question ${progress ? progress.answered + 1 : questionCount}`} title={conceptName} />
      <div className="max-w-2xl mx-auto min-w-0 overflow-hidden px-2 sm:px-0">
        <ProgressBar
          answered={progress ? progress.answered + 1 : questionCount}
          coverDone={progress ? progress.cover_done : 0}
          coverSize={progress ? progress.cover_size : estimatedTotal}
        />

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

        <div className={`bg-mathua-surface border rounded-none p-4 sm:p-6 mb-6 transition-colors duration-200 w-full max-w-full min-w-0 overflow-hidden ${lastResult ? (lastResult.correct ? 'border-mathua-green' : 'border-mathua-red') : 'border-mathua-border'}`}>
          <div className="bg-mathua-code border border-mathua-border rounded-none p-4 sm:p-6 text-center mb-4 w-full max-w-full min-w-0 overflow-hidden">
            <div className="w-full max-w-full min-w-0 overflow-hidden">
              <KatexContent className="text-mathua-primary text-lg font-mono font-light whitespace-pre-wrap break-words">{question}</KatexContent>
            </div>
          </div>
          {!lastResult ? (
            <>
              <form
                onSubmit={e => { e.preventDefault(); submitAnswer() }}
                className="flex flex-col sm:flex-row gap-3 min-w-0"
              >
                <label htmlFor="goals-answer" className="sr-only">Your answer</label>
                <input
                  ref={goalsInputRef}
                  id="goals-answer"
                  type="text"
                  value={answerInput}
                  onChange={e => setAnswerInput(e.target.value)}
                  placeholder="Your answer..."
                  enterKeyHint="go"
                  disabled={loading}
                  className="flex-1 min-w-0 bg-mathua-code border border-mathua-border rounded-none h-24 sm:h-12 px-4 font-mono text-base text-mathua-primary placeholder:text-mathua-muted focus:outline-none focus:border-mathua-blue"
                />
                <button
                  type="submit"
                  disabled={!answerInput.trim() || loading}
                  className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white rounded-none h-12 px-8 font-medium text-sm disabled:opacity-50 shrink-0"
                >
                  Check Answer
                </button>
              </form>
              <SymbolPalette targetRef={goalsInputRef} onInsert={setAnswerInput} />
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
  )
}
