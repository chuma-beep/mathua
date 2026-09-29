'use client'

import { useCallback, useEffect, useRef, useState } from 'react'
import Link from 'next/link'
import KatexContent from './KatexContent'
import ChoiceOptions from './ChoiceOptions'
import { Input } from '@/components/ui/input'
import { getLessonKPs, getLessonPractice, getLessonReadiness, submitStudyAnswer, type KpInfo, type PracticeQuestion, type ReadinessRes } from '../lib/api'
import { REQUIRED_IN_A_ROW, masteryEstimate, type Attempt } from '../lib/progression'
import { formatForGradingType } from '../lib/answerFormat'
import { concepts } from '../lib/conceptData'

// Append-only learning feed: intro once, then each submit locks its card and
// the next question opens below. Scrollback IS the performance history —
// there are no gating buttons between questions and no pre-answer hints.
// Post-mistake explanations render inline in the locked card.
interface QEntry {
  key: number
  kind: 'q'
  kpIndex: number
  q: PracticeQuestion
  answer: string
  checking: boolean
  servedAt: number
  feedback: { correct: boolean; text: string; xp: number } | null
}

type Entry =
  | { key: number; kind: 'intro'; kpIndex: number }
  | QEntry
  | { key: number; kind: 'kpdiv'; kpIndex: number }
  | { key: number; kind: 'halt' }
  | { key: number; kind: 'done' }

interface Props {
  conceptId: string
  returnTo?: string
}

let keySeq = 1
const nextKey = () => keySeq++

export default function LearnStepper({ conceptId, returnTo }: Props) {
  const [entries, setEntries] = useState<Entry[]>([])
  const [kps, setKps] = useState<KpInfo[]>([])
  const [kpIndex, setKpIndex] = useState(0)
  const [readiness, setReadiness] = useState<ReadinessRes | null>(null)
  const [bannerDismissed, setBannerDismissed] = useState(false)
  const [buffer, setBuffer] = useState<PracticeQuestion[]>([])
  const [history, setHistory] = useState<Attempt[]>([])
  const [difficulty, setDifficulty] = useState(0.4)
  const [consecutive, setConsecutive] = useState(0)
  const [misses, setMisses] = useState(0)
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(true)
  const seenRef = useRef<string[]>([])
  const bottomRef = useRef<HTMLDivElement>(null)
  const reduceMotion = useRef(false)
  const format = formatForGradingType(concepts.find(c => c.id === conceptId)?.grading_type)

  useEffect(() => {
    reduceMotion.current = typeof window !== 'undefined' && window.matchMedia?.('(prefers-reduced-motion: reduce)').matches
  }, [])

  const scrollToBottom = useCallback(() => {
    bottomRef.current?.scrollIntoView({ behavior: reduceMotion.current ? 'auto' : 'smooth', block: 'nearest' })
  }, [])

  const refill = useCallback(async (diff: number, exclude: string[]): Promise<PracticeQuestion[]> => {
    const pr = await getLessonPractice(conceptId, 3, { seed: Date.now() % 100000, exclude, difficulty: diff })
    const fresh = pr.questions.filter(q => !exclude.includes(q.question))
    return fresh.length > 0 ? fresh : pr.questions
  }, [conceptId])

  const load = useCallback(async () => {
    setLoading(true)
    setError('')
    try {
      const [kpRes, readyRes] = await Promise.all([getLessonKPs(conceptId), getLessonReadiness(conceptId)])
      setKps(kpRes.kps ?? [])
      setReadiness(readyRes)
      const pr = await getLessonPractice(conceptId, 3, { seed: Date.now() % 100000, difficulty: 0.4 })
      seenRef.current = pr.questions.map(q => q.question)
      setBuffer(pr.questions)
      setEntries([{ key: nextKey(), kind: 'intro', kpIndex: 0 }])
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Failed to load lesson')
    } finally {
      setLoading(false)
    }
  }, [conceptId])

  useEffect(() => {
    setEntries([])
    setKpIndex(0)
    setHistory([])
    setDifficulty(0.4)
    setConsecutive(0)
    setMisses(0)
    setBannerDismissed(false)
    load()
  }, [load])

  // Take the next question: buffer first, fetching when empty.
  async function takeNext(diff: number): Promise<PracticeQuestion | null> {
    if (buffer.length > 0) {
      const [head, ...rest] = buffer
      setBuffer(rest)
      if (rest.length < 2) {
        refill(diff, seenRef.current).then(fresh => {
          seenRef.current = [...seenRef.current, ...fresh.map(q => q.question)].slice(-20)
          setBuffer(prev => [...prev, ...fresh].slice(0, 6))
        }).catch(() => {})
      }
      return head
    }
    try {
      const fresh = await refill(diff, seenRef.current)
      seenRef.current = [...seenRef.current, ...fresh.map(q => q.question)].slice(-20)
      const [head, ...rest] = fresh
      setBuffer(rest)
      return head ?? null
    } catch {
      return null
    }
  }

  function appendAfter(ms: number, fn: () => void) {
    window.setTimeout(() => {
      fn()
      scrollToBottom()
    }, ms)
  }

  async function startPracticing() {
    const q = await takeNext(difficulty)
    if (!q) {
      setError('Could not load a question — check your connection and reload.')
      return
    }
    const entry: QEntry = { key: nextKey(), kind: 'q', kpIndex, q, answer: '', checking: false, servedAt: Date.now(), feedback: null }
    setEntries(prev => [...prev, entry])
    scrollToBottom()
  }

  function setEntry(key: number, patch: Partial<QEntry>) {
    setEntries(prev => prev.map(e => (e.key === key && e.kind === 'q' ? { ...e, ...patch } : e)))
  }

  async function handleCheck(key: number) {
    const entry = entries.find(e => e.key === key && e.kind === 'q') as QEntry | undefined
    if (!entry || !entry.answer.trim() || entry.checking || entry.feedback) return
    setEntry(key, { checking: true })
    const elapsed = Math.max(0.5, (Date.now() - entry.servedAt) / 1000)
    try {
      const res = await submitStudyAnswer(conceptId, entry.answer.trim(), entry.q.answer, elapsed, entry.q.question)
      if (res.ungraded) {
        setEntry(key, { checking: false })
        return
      }
      const h: Attempt = { correct: res.correct, difficulty, elapsed, variation: entry.q.question }
      const nextHistory = [...history, h]
      setHistory(nextHistory)
      if (res.correct) {
        const next = consecutive + 1
        const nextDiff = Math.min(1.0, difficulty + 0.15)
        setConsecutive(next)
        setMisses(0)
        setDifficulty(nextDiff)
        setEntry(key, { checking: false, feedback: { correct: true, text: '', xp: res.xp ?? 0 } })
        const est = masteryEstimate(nextHistory)
        const advance = est.decision === 'advance' || next >= REQUIRED_IN_A_ROW
        appendAfter(350, async () => {
          if (advance) {
            const total = Math.max(kps.length, 1)
            if (entry.kpIndex + 1 >= total) {
              setEntries(prev => [...prev, { key: nextKey(), kind: 'done' }])
            } else {
              const ni = entry.kpIndex + 1
              setKpIndex(ni)
              setConsecutive(0)
              setMisses(0)
              setHistory([])
              setDifficulty(0.4)
              const nq = await takeNext(0.4)
              if (!nq) return
              setEntries(prev => [
                ...prev,
                { key: nextKey(), kind: 'kpdiv', kpIndex: ni },
                { key: nextKey(), kind: 'q', kpIndex: ni, q: nq, answer: '', checking: false, servedAt: Date.now(), feedback: null },
              ])
            }
          } else {
            const nq = await takeNext(nextDiff)
            if (!nq) return
            setEntries(prev => [...prev, { key: nextKey(), kind: 'q', kpIndex: entry.kpIndex, q: nq, answer: '', checking: false, servedAt: Date.now(), feedback: null }])
          }
        })
      } else {
        const m = misses + 1
        const nextDiff = Math.max(0.3, difficulty - 0.15)
        setConsecutive(0)
        setMisses(m)
        setDifficulty(nextDiff)
        setEntry(key, { checking: false, feedback: { correct: false, text: res.explanation || entry.q.explanation, xp: 0 } })
        appendAfter(350, async () => {
          if (m >= 3) {
            // Safety net only: note + easier question below, prereq links.
            const nq = await takeNext(Math.max(0.3, nextDiff - 0.1))
            setEntries(prev => [...prev, { key: nextKey(), kind: 'halt' }, ...(nq ? [{ key: nextKey(), kind: 'q', kpIndex: entry.kpIndex, q: nq, answer: '', checking: false, servedAt: Date.now(), feedback: null } as QEntry] : [])])
            setMisses(0)
          } else {
            const nq = await takeNext(nextDiff)
            if (!nq) return
            setEntries(prev => [...prev, { key: nextKey(), kind: 'q', kpIndex: entry.kpIndex, q: nq, answer: '', checking: false, servedAt: Date.now(), feedback: null }])
          }
        })
      }
    } catch {
      const correct = entry.answer.trim().toLowerCase() === entry.q.answer.trim().toLowerCase()
      setHistory(prev => [...prev, { correct, difficulty, elapsed, variation: entry.q.question }])
      setEntry(key, { checking: false, feedback: { correct, text: entry.q.explanation, xp: 0 } })
      if (correct) {
        setConsecutive(c => c + 1)
        setDifficulty(d => Math.min(1.0, d + 0.15))
      } else {
        setConsecutive(0)
        setMisses(m => m + 1)
        setDifficulty(d => Math.max(0.3, d - 0.15))
      }
    }
  }

  if (loading && entries.length === 0) {
    return <div className="p-6 border border-mathua-border bg-mathua-surface"><p className="font-mono text-xs text-mathua-muted">Loading lesson…</p></div>
  }

  const prereqs = [...(readiness?.missing ?? []), ...(readiness?.weak ?? [])]
  const showBanner = readiness && !readiness.ready && !bannerDismissed && prereqs.length > 0
  const est = masteryEstimate(history)
  const totalAnswered = history.length
  const totalCorrect = history.filter(h => h.correct).length
  const totalXP = entries.reduce((s, e) => s + (e.kind === 'q' && e.feedback?.correct ? e.feedback.xp : 0), 0)

  return (
    <div className="max-w-2xl mx-auto">
      <div className="sticky top-0 z-10 bg-mathua-bg/95 backdrop-blur py-2 mb-4 flex items-center gap-2 font-mono text-[11px] text-mathua-muted border-b border-mathua-border">
        <span>KP {Math.min(kpIndex + 1, Math.max(kps.length, 1))}/{Math.max(kps.length, 1)}</span>
        {totalAnswered > 0 && <span>· {totalCorrect}/{totalAnswered} correct</span>}
        {totalXP > 0 && <span className="text-yellow-400">· +{totalXP} XP</span>}
        <span>· level {difficulty.toFixed(2)}</span>
        {est.score > 0 && <span title="Mastery estimate from accuracy, difficulty, variation and time">· {est.band}</span>}
      </div>

      {showBanner && (
        <div className="mb-4 border border-mathua-blue/50 bg-mathua-surface p-4" role="note" aria-label="Prerequisite suggestion">
          <div className="flex items-start justify-between gap-3">
            <div>
              <p className="font-mono text-[11px] uppercase tracking-wider text-mathua-blue">Before you start ({prereqs.length})</p>
              <p className="mt-1 font-mono text-xs text-mathua-secondary">This builds on {prereqs.slice(0, 3).map(p => p.label).join(', ')}. A quick review helps — or continue anyway.</p>
              <div className="mt-2 flex flex-wrap gap-2">
                {prereqs.slice(0, 3).map(p => (
                  <Link key={p.id} href={`/learn?concept=${encodeURIComponent(p.id)}&return=${encodeURIComponent(conceptId)}`} className="border border-mathua-border px-2.5 py-1.5 font-mono text-[11px] text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue">
                    Review: {p.label}
                  </Link>
                ))}
              </div>
            </div>
            <button type="button" onClick={() => setBannerDismissed(true)} className="font-mono text-[11px] text-mathua-muted hover:text-mathua-primary shrink-0" aria-label="Dismiss prerequisite suggestion">Continue anyway ✕</button>
          </div>
        </div>
      )}

      {error && <p className="mb-4 font-mono text-xs text-red-400">{error}</p>}

      <div className="space-y-4">
        {entries.map(e => {
          if (e.kind === 'intro') {
            const ikp = kps[e.kpIndex]
            return (
              <div key={e.key} className="border border-mathua-border bg-mathua-surface p-6">
                <p className="font-mono text-[11px] uppercase tracking-wider text-mathua-muted">{ikp?.label || 'Learn'}</p>
                <KatexContent className="mt-3 text-sm text-mathua-primary whitespace-pre-wrap">{ikp?.worked_example || 'Worked example unavailable — start practicing instead.'}</KatexContent>
                {ikp?.subgoals && ikp.subgoals.length > 0 && (
                  <details className="mt-3">
                    <summary className="font-mono text-[11px] text-mathua-muted cursor-pointer hover:text-mathua-blue">Steps ({ikp.subgoals.length})</summary>
                    <ol className="mt-2 space-y-1.5">
                      {ikp.subgoals.map((s, i) => <li key={i} className="font-mono text-xs text-mathua-secondary">Step {i + 1}: <KatexContent>{s}</KatexContent></li>)}
                    </ol>
                  </details>
                )}
                <button type="button" onClick={startPracticing} className="mt-5 border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white px-6 py-2.5 font-mono text-xs min-h-[44px]">Start practicing →</button>
              </div>
            )
          }
          if (e.kind === 'kpdiv') {
            const ikp = kps[e.kpIndex]
            return (
              <div key={e.key} className="border-t-2 border-mathua-blue/60 pt-4">
                <details>
                  <summary className="font-mono text-[11px] uppercase tracking-wider text-mathua-blue cursor-pointer">Next: {ikp?.label || `Knowledge point ${e.kpIndex + 1}`} (worked example)</summary>
                  <KatexContent className="mt-2 text-sm text-mathua-primary whitespace-pre-wrap">{ikp?.worked_example || ''}</KatexContent>
                </details>
              </div>
            )
          }
          if (e.kind === 'halt') {
            return (
              <div key={e.key} className="border border-red-500/40 bg-mathua-surface p-5">
                <p className="font-mono text-xs text-red-400">Stepping down a level — easier question below.</p>
                <div className="mt-3 flex flex-wrap gap-2">
                  {prereqs.slice(0, 2).map(p => (
                    <Link key={p.id} href={`/learn?concept=${encodeURIComponent(p.id)}&return=${encodeURIComponent(conceptId)}`} className="border border-mathua-border px-4 py-2 font-mono text-[11px] text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue inline-flex items-center">Review {p.label}</Link>
                  ))}
                  <Link href="/profile" className="border border-mathua-border px-4 py-2 font-mono text-[11px] text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue inline-flex items-center">Other task →</Link>
                </div>
              </div>
            )
          }
          if (e.kind === 'done') {
            return (
              <div key={e.key} className="border border-green-500/40 bg-mathua-surface p-6">
                <p className="font-mono text-xs text-green-400">✓ Complete — {totalCorrect}/{totalAnswered} correct · +{totalXP} XP · {est.band}</p>
                <p className="mt-2 font-mono text-[11px] text-mathua-secondary">Scroll up to review anything. Reviews are scheduled automatically.</p>
                <div className="mt-4 flex flex-wrap gap-2">
                  {returnTo && <Link href={`/learn?concept=${encodeURIComponent(returnTo)}`} className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white px-5 py-2 font-mono text-xs inline-flex items-center min-h-[40px]">← Back to {returnTo}</Link>}
                  <Link href="/profile" className="border border-mathua-border px-5 py-2 font-mono text-xs text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue inline-flex items-center min-h-[40px]">Next up →</Link>
                  <Link href={`/study?concept=${encodeURIComponent(conceptId)}`} className="border border-mathua-border px-5 py-2 font-mono text-xs text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue inline-flex items-center min-h-[40px]">Reference</Link>
                </div>
              </div>
            )
          }
          const qe = e as QEntry
          const locked = qe.feedback !== null
          return (
            <div key={qe.key} className={`border bg-mathua-surface p-6 ${qe.feedback?.correct ? 'border-green-500/30' : qe.feedback ? 'border-red-500/30' : 'border-mathua-border'}`}>
              <KatexContent className="text-sm text-mathua-primary font-mono whitespace-pre-wrap">{qe.q.question}</KatexContent>
              <ChoiceOptions question={qe.q.question} value={qe.answer} onPick={v => setEntry(qe.key, { answer: v })} disabled={locked || qe.checking} />
              {!locked && (
                <form onSubmit={ev => { ev.preventDefault(); handleCheck(qe.key) }} className="mt-3 flex flex-col sm:flex-row gap-2">
                  <label htmlFor={`learn-answer-${qe.key}`} className="sr-only">Your answer</label>
                  <Input id={`learn-answer-${qe.key}`} value={qe.answer} onChange={ev => setEntry(qe.key, { answer: ev.target.value })} placeholder="Your answer…" inputMode={format.inputMode} className="sm:flex-1 bg-mathua-bg" />
                  <button type="submit" disabled={qe.checking || !qe.answer.trim()} className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white px-6 h-12 text-sm font-mono disabled:opacity-50 w-auto self-end sm:self-auto shrink-0">Check</button>
                </form>
              )}
              {!locked && <p className="mt-2 font-mono text-[11px] text-mathua-muted">{format.hint}</p>}
              {locked && qe.feedback && (
                <div className="mt-3">
                  <p className={`font-mono text-xs ${qe.feedback.correct ? 'text-green-400' : 'text-red-400'}`}>
                    {qe.feedback.correct ? `✓ ${qe.feedback.xp ? `+${qe.feedback.xp} XP` : 'Correct'}` : '✗ Not quite — the method:'} <span className="text-mathua-muted">you answered “{qe.answer}”</span>
                  </p>
                  {!qe.feedback.correct && (
                    <KatexContent className="mt-1.5 text-xs font-mono text-mathua-secondary whitespace-pre-wrap">{qe.feedback.text}</KatexContent>
                  )}
                </div>
              )}
            </div>
          )
        })}
        <div ref={bottomRef} aria-hidden="true" />
      </div>
    </div>
  )
}
