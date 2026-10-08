'use client'

import { useCallback, useEffect, useMemo, useRef, useState } from 'react'
import Link from 'next/link'
import KatexContent from './KatexContent'
import LessonDiagram from './LessonDiagram'
import LessonAssets from './LessonAssets'
import ChoiceOptions from './ChoiceOptions'
import ConceptReference from './ConceptReference'
import PrerequisitePanel from './PrerequisitePanel'
import { getLessonKPs, getLessonPractice, getLessonReadiness, submitStudyAnswer, getProgress, getErrorStatus, type KpInfo, type PracticeQuestion, type ReadinessRes, type ConceptProgress } from '../lib/api'
import { getUserInfo } from '../lib/auth'
import { upcomingLocked, hrefConceptId, type Shelf, type LockedSuccessor } from '../lib/nextUp'
import { fetchShelf } from '../lib/recommendations'
import { REQUIRED_IN_A_ROW } from '../lib/progression'
import { formatForGradingType } from '../lib/answerFormat'
import { countsAsMastered } from '../lib/progress'
import { MathAnswerInput } from './math/MathInput'
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
  | { key: number; kind: 'done'; note?: string }

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
  // The concept's own diagram. The KP endpoint has always returned it and the client has
  // always dropped it, so the one figure guaranteed to be relevant to the concept being
  // studied was the one figure /learn never showed.
  const [diagram, setDiagram] = useState<string | null>(null)
  const [kpIndex, setKpIndex] = useState(0)
  const [readiness, setReadiness] = useState<ReadinessRes | null>(null)
  // The server refused this topic because its prerequisites are unmet. The
  // practice endpoint enforces this; the UI renders the explanation so a locked
  // topic never looks like a broken page or an empty question set.
  const [ineligible, setIneligible] = useState<ReadinessRes | null>(null)
  const [bannerDismissed, setBannerDismissed] = useState(false)
  const [buffer, setBuffer] = useState<PracticeQuestion[]>([])
  // Counts for the header, per *concept*. Not per knowledge point: this used to reset on
  // every KP advance, so a learner three questions in with two correct saw "1/2 correct"
  // sitting next to "+3 XP" — the XP sums every card in the feed while the count summed only
  // the current section. Two scopes on one line, and the count was the one that looked wrong.
  // Two numbers rather than an array because nothing else needed the answers themselves.
  const [answeredCount, setAnsweredCount] = useState(0)
  const [correctCount, setCorrectCount] = useState(0)
  // Plain-language progress reading, computed server-side by mastery.EvidenceBand from
  // the same evidence the ladder decides on. The client used to compute its own score
  // and its own bands here; see test/masteryAuthority.test.ts.
  const [band, setBand] = useState('')
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
  const conceptGradingType = concepts.find(c => c.id === conceptId)?.grading_type
  const format = formatForGradingType(conceptGradingType)
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
    setIneligible(null)
    try {
      const [kpRes, readyRes] = await Promise.all([getLessonKPs(conceptId), getLessonReadiness(conceptId)])
      setKps(kpRes.kps ?? [])
      setDiagram(kpRes.diagram ?? null)
      setReadiness(readyRes)
      // The server is the authority: if it says the topic's prerequisites are
      // unmet, do not fetch questions it would refuse. Show the path instead.
      if (readyRes.eligible === false) {
        setIneligible(readyRes)
        return
      }
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
    setAnsweredCount(0)
    setCorrectCount(0)
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
    setAnsweredCount(0)
    setCorrectCount(0)
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
        fetchShelf([conceptId]),
        // Progress is still needed here, for "Coming up": the successors this concept will
        // unlock once it is mastered.
        sid ? getProgress(sid).catch(() => ({} as Record<string, ConceptProgress>)) : Promise.resolve({} as Record<string, ConceptProgress>),
      ]),
      timeout,
    ]).then(([next, progress]) => {
      setShelfProgress(progress)
      // The concept just finished is excluded server-side. It is unmastered by definition
      // here, so without that it would head itself and the Continue link would point at the
      // current page — a scroll-to-top dead end. That argument moved to the engine with the
      // rest of the decision.
      setNextShelf(next)
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

  /**
   * Move to the next knowledge point, or finish the concept.
   *
   * Shared by the 2-in-a-row rule and by question exhaustion, because "this section is
   * over" means the same thing whichever caused it.
   */
  async function advanceSection(from: QEntry, why: 'complete' | 'exhausted') {
    const total = Math.max(kps.length, 1)
    const note = why === 'exhausted'
      ? 'This topic ran out of new questions, so the section ended here.'
      : undefined
    if (from.kpIndex + 1 >= total) {
      setEntries(prev => [...prev, { key: nextKey(), kind: 'done', note }])
      return
    }
    const ni = from.kpIndex + 1
    setKpIndex(ni)
    setConsecutive(0)
    setMisses(0)
    setDifficulty(0.4)
    // Ask for the next section's question before announcing the section, so a concept with
    // nothing left anywhere goes straight to the completion card instead of showing a
    // heading and then nothing under it.
    const nq = await takeNext(0.4)
    setEntries(prev => [...prev, { key: nextKey(), kind: 'kpdiv', kpIndex: ni }])
    if (nq) {
      setEntries(prev => [
        ...prev,
        { key: nextKey(), kind: 'q', kpIndex: ni, q: nq, answer: '', checking: false, servedAt: Date.now(), feedback: null },
      ])
    } else {
      setEntries(prev => [
        ...prev,
        {
          key: nextKey(),
          kind: 'done',
          note: note ?? 'This topic ran out of new questions, so there was nothing left to work through.',
        },
      ])
    }
  }

  /**
   * Append the next question in this section, or move on if there isn't one.
   *
   * This is the fix for the dead end. `takeNext` returns null when the buffer is empty and
   * the refill yields nothing, which used to hit `if (!nq) return` at three sites: no card,
   * no message, no way forward, and the last verdict became the end of the page. A learner
   * hit this on `geo.basic.points_lines`, a concept whose generator has exactly four
   * distinct question texts, so it was guaranteed on the fourth answer rather than an edge
   * case.
   *
   * The practice endpoint now falls back to repeats rather than returning an empty 200, so
   * this is the second line of defence rather than the only one — but a client that can be
   * stranded by a 200 with an empty array has no business relying on the server not to.
   */
  async function continueInSection(from: QEntry, diff: number) {
    const nq = await takeNext(diff)
    if (nq) {
      setEntries(prev => [
        ...prev,
        { key: nextKey(), kind: 'q', kpIndex: from.kpIndex, q: nq, answer: '', checking: false, servedAt: Date.now(), feedback: null },
      ])
      return
    }
    await advanceSection(from, 'exhausted')
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
      setAnsweredCount(n => n + 1)
      if (res.correct) setCorrectCount(n => n + 1)
      if (res.evidence_band) setBand(res.evidence_band)
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
        // 2-in-a-row, and nothing else. This used to read
        //   `est.decision === 'advance' || next >= REQUIRED_IN_A_ROW`
        // where `est` came from a second, client-side copy of the mastery evidence model.
        // The first operand was dead: `decision === 'advance'` required the last two
        // attempts correct, which is the same fact as `next >= REQUIRED_IN_A_ROW`, so it
        // could never be the deciding half. Proven exhaustively in
        // test/masteryAuthority.test.ts. The model is gone; this rule is the whole gate.
        const advance = next >= REQUIRED_IN_A_ROW
        appendAfter(350, async () => {
          if (advance) {
            await advanceSection(entry, 'complete')
          } else {
            await continueInSection(entry, nextDiff)
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
            if (nq) {
              setEntries(prev => [
                ...prev,
                { key: nextKey(), kind: 'halt' },
                { key: nextKey(), kind: 'q', kpIndex: entry.kpIndex, q: nq, answer: '', checking: false, servedAt: Date.now(), feedback: null } as QEntry,
              ])
            } else {
              // The halt note still belongs even with no question to put under it, and the
              // section still has to end somewhere the learner can move on from.
              setEntries(prev => [...prev, { key: nextKey(), kind: 'halt' }])
              await advanceSection(entry, 'exhausted')
            }
            setMisses(0)
          } else {
            await continueInSection(entry, nextDiff)
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

  // Which sections have been attempted, meaning answered rather than served. A 'q' entry
  // exists from the moment the question is put on screen, so keying on the entry alone would
  // reveal the steps at the very moment the learner is about to try the question unaided.
  // Derived from the feed, which is append-only, so no state to keep in sync.
  const answeredSections = useMemo(() => {
    const seen = new Set<number>()
    for (const e of entries) if (e.kind === 'q' && e.feedback) seen.add(e.kpIndex)
    return seen
  }, [entries])

  if (loading && entries.length === 0) {
    return <div className="p-6 border border-mathua-border bg-mathua-surface"><p className="font-mono text-xs text-mathua-muted">Loading lesson…</p></div>
  }

  // A locked topic is explained, never started. The server refused the questions;
  // this renders its reason and a path into the prerequisite's own learning flow.
  if (ineligible && entries.length === 0) {
    const title = concepts.find(c => c.id === conceptId)?.label ?? conceptId
    return (
      <div className="max-w-2xl mx-auto">
        <PrerequisitePanel topicTitle={title} prerequisites={ineligible.prerequisites ?? []} />
        <div className="mt-4 flex flex-wrap gap-3">
          <Link href="/domains" className="border border-mathua-border px-4 py-2 font-mono text-[11px] text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue">
            ← Choose another topic
          </Link>
          <Link href="/profile" className="border border-mathua-border px-4 py-2 font-mono text-[11px] text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue">
            Back to Profile
          </Link>
        </div>
      </div>
    )
  }

  const prereqs = [...(readiness?.missing ?? []), ...(readiness?.weak ?? [])]
  // A prerequisite the learner has already demonstrated is a *retrieval* question, not a
  // teaching one, and the two must not be collapsed.
  //
  // The server already draws the distinction — `ready` buckets a prerequisite into
  // `missing` only when UNSEEN and into `weak` otherwise, so a decayed one arrives in
  // `weak` (`server.go` handleReadiness) — and then the client concatenated the buckets
  // identically and gave every entry the same `/learn` link. That turned "you have
  // mastered this and it is due a check" into "here is the tutorial again", which is the
  // purgatory arriving through the back door after `buildCandidates` was fixed.
  //
  // So the split is kept and each half gets the destination that matches it: still being
  // learned goes to `/learn`, already demonstrated goes to `/review`, which is
  // problems-first and fetches no lesson.
  const toLearn = prereqs.filter(p => !countsAsMastered({ status: p.status }))
  const toReview = prereqs.filter(p => countsAsMastered({ status: p.status }))
  const onlyReview = toLearn.length === 0 && toReview.length > 0
  const showBanner = readiness && !readiness.ready && !bannerDismissed && prereqs.length > 0
  const totalAnswered = answeredCount
  const totalCorrect = correctCount
  const totalXP = entries.reduce((s, e) => s + (e.kind === 'q' && e.feedback?.correct ? e.feedback.xp : 0), 0)

  return (
    <div className="max-w-2xl mx-auto">
      <div className="sticky top-0 z-10 bg-mathua-bg backdrop-blur py-2 mb-4 flex items-center gap-2 font-mono text-[11px] text-mathua-muted border-b border-mathua-border">
        <span>Section {Math.min(kpIndex + 1, Math.max(kps.length, 1))} of {Math.max(kps.length, 1)}</span>
        {totalAnswered > 0 && <span>· {totalCorrect}/{totalAnswered} correct</span>}
        {totalXP > 0 && <span className="text-yellow-400">· +{totalXP} XP</span>}
        {band && (
          <span title="How you're doing on this concept, from your recent answers. Progress only — mastery is decided separately.">
            · {band}
          </span>
        )}
      </div>

      {showBanner && (
        <div className="mb-4 border border-mathua-blue bg-mathua-surface p-4" role="note" aria-label="Prerequisite suggestion">
          <div className="flex items-start justify-between gap-3">
            <div>
              <p className="font-mono text-[11px] uppercase tracking-wider text-mathua-blue">
                {onlyReview ? 'Due a retrieval check' : `Before you start (${prereqs.length})`}
              </p>
              <p className="mt-1 font-mono text-xs text-mathua-secondary">
                {onlyReview
                  ? `You have already learned ${toReview.slice(0, 3).map(p => p.label).join(', ')} — a short check is due, or continue anyway.`
                  : `This builds on ${toLearn.slice(0, 3).map(p => p.label).join(', ')}. A quick review helps — or continue anyway.`}
              </p>
              <div className="mt-2 flex flex-wrap gap-2">
                {toLearn.slice(0, 3).map(p => (
                  <Link key={p.id} href={`/learn?concept=${encodeURIComponent(p.id)}&return=${encodeURIComponent(conceptId)}`} className="border border-mathua-border px-2.5 py-1.5 font-mono text-[11px] text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue">
                    Learn: {p.label}
                  </Link>
                ))}
                {toReview.slice(0, 3).map(p => (
                  <Link key={p.id} href="/review" className="border border-mathua-border px-2.5 py-1.5 font-mono text-[11px] text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue">
                    Check: {p.label}
                  </Link>
                ))}
              </div>
            </div>
            <button type="button" onClick={() => setBannerDismissed(true)} className="font-mono text-[11px] text-mathua-muted hover:text-mathua-primary shrink-0" aria-label="Dismiss prerequisite suggestion">Continue anyway ✕</button>
          </div>
        </div>
      )}

      {error && <p className="mb-4 font-mono text-xs text-red-400">{error}</p>}

      {/* The reference layer, in the flow.
          This replaces the "Reference" link the done card used to carry to `/study`: the same
          material, reachable without leaving the questions, and impossible to wander off in —
          there is nothing in it that navigates. Collapsed by default, so it never stands between
          a learner and the first question; the worked example above it is the intended first
          read. It renders no mastery state and awards nothing, which is the point: reading
          explains, and only an answer moves the concept. */}
      <ConceptReference conceptId={conceptId} prerequisites={prereqs} />

      <div className="space-y-4">
        {entries.map(e => {
          if (e.kind === 'intro') {
            const ikp = kps[e.kpIndex]
            return (
              <div key={e.key} className="border border-mathua-border bg-mathua-surface p-6">
                <p className="font-mono text-[11px] uppercase tracking-wider text-mathua-muted">{ikp?.label || 'Learn'}</p>
                <KatexContent className="mt-3 text-sm text-mathua-primary whitespace-pre-wrap">{ikp?.worked_example || 'Worked example unavailable — start practicing instead.'}</KatexContent>
                {/* Figures for this knowledge point, plus the concept's own diagram. The
                    worked example is text lifted out of the lesson body, so its figures do
                    not come with it; selecting them per knowledge point is what keeps the
                    shorter surface shorter without discarding instruction. */}
                <LessonAssets assets={ikp?.assets} />
                {diagram && (
                  <figure className="mt-4">
                    <LessonDiagram src={diagram} alt={`Diagram for ${ikp?.label || 'this topic'}`} />
                  </figure>
                )}
                {/* Steps appear once this section has been attempted. Listed before the
                    question they are a syllabus the learner has not chosen yet, and the
                    effect was to read the steps and skip the work they were meant to
                    support. Derived from the feed, which is append-only, so no state to keep
                    in sync and no way for the two to disagree. */}
                {ikp?.subgoals && ikp.subgoals.length > 0 && answeredSections.has(e.kpIndex) && (
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
                  <LessonAssets assets={ikp?.assets} />
                </details>
              </div>
            )
          }
          if (e.kind === 'halt') {
            return (
              <div key={e.key} className="border border-mathua-red-faint bg-mathua-surface p-5">
                <p className="font-mono text-xs text-red-400">Stepping down a level — easier question below.</p>
                <div className="mt-3 flex flex-wrap gap-2">
                  {/* Same split as the readiness banner: a decayed prerequisite is owed
                      a retrieval check, not another tutorial. */}
                  {toLearn.slice(0, 2).map(p => (
                    <Link key={p.id} href={`/learn?concept=${encodeURIComponent(p.id)}&return=${encodeURIComponent(conceptId)}`} className="border border-mathua-border px-4 py-2 font-mono text-[11px] text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue inline-flex items-center">Review {p.label}</Link>
                  ))}
                  {toReview.slice(0, 2).map(p => (
                    <Link key={p.id} href="/review" className="border border-mathua-border px-4 py-2 font-mono text-[11px] text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue inline-flex items-center">Check {p.label}</Link>
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
              <div key={e.key} className="border border-mathua-green-faint bg-mathua-surface p-6">
                <p className="font-mono text-xs text-green-400">✓ Complete — {totalCorrect}/{totalAnswered} correct · +{totalXP} XP · {band}</p>
                {e.note && <p role="status" className="mt-2 font-mono text-[11px] text-mathua-secondary">{e.note}</p>}
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
            <div key={qe.key} className={`border bg-mathua-surface p-6 ${qe.feedback?.correct ? 'border-mathua-green-faint' : qe.feedback ? 'border-mathua-red-faint' : 'border-mathua-border'}`}>
              <KatexContent className="text-sm text-mathua-primary font-mono whitespace-pre-wrap">{qe.q.question}</KatexContent>
              <ChoiceOptions question={qe.q.question} value={qe.answer} onPick={v => setEntry(qe.key, { answer: v })} disabled={locked || qe.checking} />
              {!locked && (
                <form onSubmit={ev => { ev.preventDefault(); handleCheck(qe.key) }} className="mt-3 flex flex-col sm:flex-row gap-2">
                  <label htmlFor={`learn-answer-${qe.key}`} className="sr-only">Your answer</label>
                  <MathAnswerInput
                    id={`learn-answer-${qe.key}`}
                    value={qe.answer}
                    onChange={v => setEntry(qe.key, { answer: v })}
                    onSubmit={() => handleCheck(qe.key)}
                    gradingType={conceptGradingType}
                    conceptId={conceptId}
                    disabled={qe.checking}
                    status={qe.feedback ? (qe.feedback.correct ? 'correct' : 'incorrect') : 'default'}
                    inputMode={format.inputMode}
                  placeholder="Your answer"
                    className="sm:flex-1"
                  />
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
