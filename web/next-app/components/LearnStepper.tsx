'use client'

import { useCallback, useEffect, useRef, useState } from 'react'
import Link from 'next/link'
import KatexContent from './KatexContent'
import ChoiceOptions from './ChoiceOptions'
import { Input } from '@/components/ui/input'
import { getLessonKPs, getLessonPractice, getLessonReadiness, submitStudyAnswer, getActivity, getDueReviews, getProgress, getScores, getWeaknesses, getErrorStatus, type KpInfo, type PracticeQuestion, type ReadinessRes, type DailyActivity, type Scores, type WeaknessRes, type ConceptProgress } from '../lib/api'
import { getUserInfo } from '../lib/auth'
import { selectShelfHead, upcomingLocked, hrefConceptId, type Shelf, type LockedSuccessor } from '../lib/nextUp'
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
  feedback: { correct: boolean; text: string; diagnosis?: string; xp: number } | null
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

// Shelf re-fetch budget on the done card: past this the failure branch
// (Practice again + Back to Profile) renders instead of a stale spinner.
const SHELF_TIMEOUT_MS = 8000

function catalogEntries() {
  return concepts.map(c => ({ id: c.id, label: c.label, prerequisites: c.prerequisites ?? [], avgTimeSeconds: c.mastery_threshold?.avg_time_seconds }))
}

// The done guard compares concept ids, never raw href strings: live URLs may
// carry &seed=/&difficulty=/&exclude= while the stepper still won't reset,
// so only a real concept change may render a navigating <Link>. Shared with
// Profile dedupe via lib/nextUp; re-exported under the guard's name for tests.
export { hrefConceptId as headConceptId }

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
  // Next head for the done state: re-fetched on entering done (after the
  // scheduler ingests the last answer), so Continue names what will load.
  // No navigation until the click — learner disposes via the alternatives.
  const [nextShelf, setNextShelf] = useState<Shelf | null>(null)
  const [nextLoading, setNextLoading] = useState(false)
  const [shelfProgress, setShelfProgress] = useState<Record<string, ConceptProgress>>({})
  const nextFetchedRef = useRef(false)
  const seenRef = useRef<string[]>([])
  const bottomRef = useRef<HTMLDivElement>(null)
  const reduceMotion = useRef(false)
  const format = formatForGradingType(concepts.find(c => c.id === conceptId)?.grading_type)
  // A learner reading "← Back to frac.add.diff" is reading a storage key.
  const returnToLabel = concepts.find(c => c.id === returnTo)?.label ?? returnTo

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
    setNextShelf(null)
    setNextLoading(false)
    nextFetchedRef.current = false
    load()
  }, [load])

  // Deliberate same-concept restart with fresh variants. Rendered as the
  // secondary action on the done card — never a silent fallback.
  const practiceAgain = useCallback(() => {
    setEntries([])
    setKpIndex(0)
    setHistory([])
    setDifficulty(0.4)
    setConsecutive(0)
    setMisses(0)
    setBannerDismissed(false)
    setError('')
    seenRef.current = []
    setNextShelf(null)
    setNextLoading(false)
    nextFetchedRef.current = false
    load()
  }, [load])

  // Re-fetch the ranked head on entering done: the scheduler has just
  // ingested the last answer, so the head may have moved. Fetch in place —
  // navigation waits for the Continue click (deterministic for e2e).
  // A slow network must not strand the card: past the timeout the failure
  // branch (Practice again + Back to Profile) renders instead of a stale
  // spinner, and never a self-link.
  const hasDone = entries.some(e => e.kind === 'done')
  useEffect(() => {
    if (!hasDone || nextFetchedRef.current) return
    nextFetchedRef.current = true
    setNextLoading(true)
    const info = getUserInfo()
    const sid = info?.student_id ?? ''
    const timeout = new Promise<never>((_, reject) => {
      window.setTimeout(() => reject(new Error('shelf-timeout')), SHELF_TIMEOUT_MS)
    })
    Promise.race([
      Promise.all([
        getActivity().catch(() => [] as DailyActivity[]),
        sid ? getProgress(sid).catch(() => ({} as Record<string, ConceptProgress>)) : Promise.resolve({} as Record<string, ConceptProgress>),
        getWeaknesses().catch(() => ({ by_domain: {} } as WeaknessRes)),
        getDueReviews().catch(() => ({ count: 0 })),
        sid ? getScores(sid).catch(() => null) : Promise.resolve(null),
      ]),
      timeout,
    ]).then(([a, p, w, r, s]) => {
      const progress = p as Record<string, ConceptProgress>
      setShelfProgress(progress)
      setNextShelf(selectShelfHead({
        dueReviews: r.count ?? 0,
        weaknesses: w,
        progress,
        activity: a,
        diagnosticCompleted: info?.diagnostic_completed ?? false,
        conceptsMastered: (s as Scores | null)?.concepts_mastered ?? 0,
        catalog: catalogEntries(),
        // Never head the concept just finished: it is unmastered by
        // definition here, so without this it would resume itself and the
        // Continue link would point at the current page (scroll-to-top dead
        // end). Passed only by this done re-fetch.
        excludeConceptIds: [conceptId],
      }))
    }).catch(() => {
      nextFetchedRef.current = false
    }).finally(() => setNextLoading(false))
  }, [hasDone])

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
      const res = await submitStudyAnswer(conceptId, entry.answer.trim(), elapsed, entry.q.question)
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
        setEntry(key, {
          checking: false,
          feedback: {
            correct: true,
            text: res.explanation ?? entry.q.explanation,
            diagnosis: res.diagnosis,
            xp: res.xp ?? 0,
          },
        })
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
        setEntry(key, {
          checking: false,
          feedback: {
            correct: false,
            text: res.explanation || entry.q.explanation,
            diagnosis: res.diagnosis,
            xp: 0,
          },
        })
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
    } catch (e) {
      // The server is the only grader. This never decides whether the answer
      // was right — a network fault or an expired server-side record is not a
      // miss, and grading locally would have made it one.
      //
      // 409 means the server has no record of this question, so resubmitting
      // it could never succeed: re-serve instead, which writes a fresh record
      // and leaves the learner with a question they can actually answer. Any
      // other failure keeps the question and the learner's answer so a retry
      // can succeed. Nothing was recorded either way — no attempt, streak,
      // weakness or XP moved.
      if (getErrorStatus(e) === 409) {
        const fresh = await refill(difficulty, seenRef.current)
        if (fresh.length > 0) {
          seenRef.current = [...seenRef.current, fresh[0].question].slice(-20)
          setEntry(key, { checking: false, q: fresh[0], answer: '', servedAt: Date.now() })
          setBuffer(prev => [...prev, ...fresh.slice(1)].slice(0, 6))
          setError('That question had expired on the server, so it could not be graded. Here is a fresh one — nothing was recorded.')
          return
        }
        setError('That question had expired and no replacement could be loaded. Reload to continue.')
        setEntry(key, { checking: false })
        return
      }
      setEntry(key, { checking: false })
      setError('Could not reach the server to grade that answer. Check your connection and submit again — nothing was recorded.')
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
      <div className="sticky top-0 z-10 bg-mathua-bg backdrop-blur py-2 mb-4 flex items-center gap-2 font-mono text-[11px] text-mathua-muted border-b border-mathua-border">
        <span>KP {Math.min(kpIndex + 1, Math.max(kps.length, 1))}/{Math.max(kps.length, 1)}</span>
        {totalAnswered > 0 && <span>· {totalCorrect}/{totalAnswered} correct</span>}
        {totalXP > 0 && <span className="text-yellow-400">· +{totalXP} XP</span>}
        <span>· level {difficulty.toFixed(2)}</span>
        {est.score > 0 && <span title="Mastery estimate from accuracy, difficulty, variation and time">· {est.band}</span>}
      </div>

      {showBanner && (
        <div className="mb-4 border border-mathua-blue bg-mathua-surface p-4" role="note" aria-label="Prerequisite suggestion">
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
                {/* Progression, not a waiver: the feed is append-only, so the
                    example stays on screen as scrollback — nothing is skipped.
                    Matches the kpdiv "Next: … (worked example)" label below, and
                    MA has no skip affordance at all: the worked example is the
                    scaffolding for the question that follows (example-problem
                    pair), so the label must not invite bypassing it. */}
                <button type="button" onClick={startPracticing} className="mt-5 border border-mathua-blue text-mathua-blue hover:bg-mathua-blue-faint px-6 py-2.5 font-mono text-xs min-h-[44px]">Next →</button>
                <p className="mt-2 font-mono text-[11px] text-mathua-muted">Work it through on paper first — the question below is the same move.</p>
              </div>
            )
          }
          if (e.kind === 'kpdiv') {
            const ikp = kps[e.kpIndex]
            return (
              <div key={e.key} className="border-t-2 border-mathua-blue pt-4">
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
            const headId = nextShelf ? hrefConceptId(nextShelf.next.href) : null
            // Only a real concept change may render a navigating <Link>: the
            // stepper keys everything off conceptId, so a same-concept href
            // (however many query params it carries) would scroll to top and
            // strand the learner.
            const canContinue = !nextLoading && nextShelf && headId !== null && headId !== conceptId
            const upcoming: LockedSuccessor[] = upcomingLocked(catalogEntries(), shelfProgress, conceptId)
            return (
              <div key={e.key} className="border border-green-500/40 bg-mathua-surface p-6">
                <p className="font-mono text-xs text-green-400">✓ Complete — {totalCorrect}/{totalAnswered} correct · +{totalXP} XP · {est.band}</p>
                <p className="mt-2 font-mono text-[11px] text-mathua-secondary">Scroll up to review anything. Reviews are scheduled automatically.</p>
                <div className="mt-4 flex flex-wrap gap-2">
                  {returnTo && <Link href={`/learn?concept=${encodeURIComponent(returnTo)}`} className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue-faint px-5 py-2 font-mono text-xs inline-flex items-center min-h-[40px]">← Back to {returnToLabel}</Link>}
                  {nextLoading && <span className="font-mono text-xs text-mathua-muted inline-flex items-center min-h-[40px]">Finding what&apos;s next…</span>}
                  {canContinue && nextShelf && (
                    <Link href={nextShelf.next.href} className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue-faint px-5 py-2 font-mono text-xs inline-flex items-center min-h-[40px]">{nextShelf.next.badge}: {nextShelf.next.title} →</Link>
                  )}
                  {!nextLoading && !nextShelf && (
                    <Link href="/profile" className="border border-mathua-border px-5 py-2 font-mono text-xs text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue inline-flex items-center min-h-[40px]">Back to Profile</Link>
                  )}
                  <button type="button" onClick={practiceAgain} className="border border-mathua-border px-5 py-2 font-mono text-xs text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue inline-flex items-center min-h-[40px]">Practice again</button>
                  <Link href={`/study?concept=${encodeURIComponent(conceptId)}&from=${encodeURIComponent(conceptId)}`} className="border border-mathua-border px-5 py-2 font-mono text-xs text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue inline-flex items-center min-h-[40px]">Reference</Link>
                </div>
                {!nextLoading && nextShelf && nextShelf.alternatives.length > 0 && (
                  <details className="mt-3 border border-mathua-border">
                    <summary className="font-mono text-[11px] text-mathua-secondary cursor-pointer px-4 py-2.5">
                      Or pick something else ({nextShelf.alternatives.length})
                    </summary>
                    <div className="px-4 pb-4 grid grid-cols-1 sm:grid-cols-2 gap-2">
                      {nextShelf.alternatives.map((it, i) => (
                        <Link
                          key={`${it.kind}:${it.href}:${i}`}
                          href={it.href}
                          className="border border-mathua-border p-3 hover:border-mathua-blue transition-colors block min-w-0"
                        >
                          <div className="font-mono text-xs text-mathua-primary truncate">{it.title}</div>
                          <div className="mt-1 font-mono text-[11px] text-mathua-blue">{it.cta}</div>
                        </Link>
                      ))}
                    </div>
                  </details>
                )}
                {upcoming.length > 0 && (
                  <div className="mt-3 border border-mathua-border">
                    <p className="font-mono text-[11px] text-mathua-secondary px-4 py-2.5">
                      Coming up — unlocks once this concept is mastered
                    </p>
                    <div className="px-4 pb-4 grid grid-cols-1 sm:grid-cols-2 gap-2">
                      {upcoming.map(u => (
                        <Link
                          key={u.id}
                          href={`/learn?concept=${encodeURIComponent(u.id)}`}
                          className="border border-mathua-border p-3 hover:border-mathua-blue transition-colors block min-w-0"
                        >
                          <div className="font-mono text-xs text-mathua-primary truncate">{u.label}</div>
                          <div className="mt-1 font-mono text-[11px] text-mathua-muted">Needs: {u.missing.join(', ')}</div>
                        </Link>
                      ))}
                    </div>
                  </div>
                )}
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
                  <button type="submit" disabled={qe.checking || !qe.answer.trim()} className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue-faint px-6 h-12 text-sm font-mono disabled:opacity-50 w-auto self-end sm:self-auto shrink-0">Check</button>
                </form>
              )}
              {!locked && <p className="mt-2 font-mono text-[11px] text-mathua-muted">{format.hint}</p>}
              {locked && qe.feedback && (
                <div className="mt-3">
                  <p className={`font-mono text-xs ${qe.feedback.correct ? 'text-green-400' : 'text-red-400'}`}>
                    {qe.feedback.correct ? `✓ ${qe.feedback.xp ? `+${qe.feedback.xp} XP` : 'Correct'}` : '✗ Not quite.'}{' '}
                    <span className="text-mathua-muted">you answered “{qe.answer}”</span>
                  </p>
                  <p className="mt-1.5 font-mono text-xs text-mathua-primary">
                    Answer: <KatexContent>{qe.q.answer}</KatexContent>
                  </p>
                  {/* The explanation is the solution for *this* instance, shown
                      on both verdicts: a correct answer still needs the
                      reasoning that made it right. Guarded on a non-empty
                      string so a server that returns nothing leaves no empty
                      section behind. */}
                  {/* The mistake, when it can be named with certainty. It is
                      a description of what happened, never a hint about what
                      to do, and it is empty far more often than not. */}
                  {qe.feedback.diagnosis && (
                    <p data-diagnosis className="mt-2 font-mono text-[11px] text-mathua-muted">{qe.feedback.diagnosis}</p>
                  )}
                  {qe.feedback.text && (
                    <div className="mt-1.5">
                      <p className="font-mono text-[11px] uppercase tracking-wider text-mathua-muted">
                        {qe.feedback.correct ? 'Why' : 'How'}
                      </p>
                      <KatexContent className="mt-1 text-xs font-mono text-mathua-secondary whitespace-pre-wrap">{qe.feedback.text}</KatexContent>
                    </div>
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
