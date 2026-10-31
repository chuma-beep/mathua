'use client'

import { useCallback, useEffect, useMemo, useRef, useState } from 'react'
import Link from 'next/link'
import KatexContent from './KatexContent'
import LessonDiagram from './LessonDiagram'
import LessonAssets from './LessonAssets'
import ChoiceOptions from './ChoiceOptions'
import { stemWithoutChoices } from '../lib/choices'
import ConceptReference from './ConceptReference'
import PrerequisitePanel from './PrerequisitePanel'
import { getLessonKPs, getLessonPractice, getLessonReadiness, submitStudyAnswer, getProgress, getErrorStatus, getTopicLockedPrerequisites, type EligibilityPrereq, type KpInfo, type PracticeQuestion, type ReadinessRes, type ConceptProgress } from '../lib/api'
import ReportButton from './ReportButton'
import { getUserInfo } from '../lib/auth'
import { upcomingLocked, hrefConceptId, type Shelf, type LockedSuccessor } from '../lib/nextUp'
import { backTarget, learnConceptPath, pushReturn, studyConceptPath, withReturnStack, type ReturnEntry } from '../lib/learnReturn'
import { fetchShelf } from '../lib/recommendations'
import { REQUIRED_IN_A_ROW } from '../lib/progression'
import { formatForGradingType, type AnswerFormat } from '../lib/answerFormat'
import { countsAsMastered } from '../lib/progress'
import { MathAnswerInput, type MathFocusHandle } from './math/MathInput'
import { concepts } from '../lib/conceptData'

// Append-only learning feed: intro once, then each submit locks its card and
// the next question opens below. Scrollback IS the performance history, 
// there are no gating buttons between questions and no pre-answer hints.
// Post-mistake explanations render inline in the locked card.
interface QEntry {
  key: number
  kind: 'q'
  kpIndex: number
  q: PracticeQuestion
  answer: string
  checking: boolean
  error?: string
  blockedPrerequisiteIds?: string[]
  servedAt: number
  feedback: { correct: boolean; text: string; expected: string; diagnosis?: string; xp: number } | null
}

type Entry =
  | { key: number; kind: 'intro'; kpIndex: number }
  | QEntry
  | { key: number; kind: 'kpdiv'; kpIndex: number }
  | { key: number; kind: 'halt' }
  | { key: number; kind: 'done'; note?: string }

interface Props {
  conceptId: string
  /**
   * The chain of pages the learner came through, most-recent-first (`/learn?concept=X`,
   * `/study?concept=Y`). Carried in the URL so a navigation does not strand it.
   */
  returnStack?: ReturnEntry[]
}

let keySeq = 1
const nextKey = () => keySeq++

// The scheduler's state word for a stored status, used only by the legacy fallback below. An
// unrecognised or missing status becomes '', which `PrerequisitePanel` renders as "Not started".
function prereqStateFromStatus(status: string | undefined): string {
  switch ((status ?? '').toUpperCase()) {
    case 'UNSEEN':
      return 'unlocked'
    case 'LEARNING':
    case 'PRACTICING':
      return 'in_progress'
    case 'MASTERED':
      return 'mastered'
    case 'DECAYING':
      return 'due_for_review'
    default:
      return ''
  }
}

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

export default function LearnStepper({ conceptId, returnStack = [] }: Props) {
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
  // sitting next to "+3 XP", the XP sums every card in the feed while the count summed only
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
  const [startingPractice, setStartingPractice] = useState(false)
  // Next head for the done state: re-fetched on entering done (after the
  // scheduler ingests the last answer), so Continue names what will load.
  // No navigation until the click, learner disposes via the alternatives.
  const [nextShelf, setNextShelf] = useState<Shelf | null>(null)
  const [nextLoading, setNextLoading] = useState(false)
  const [shelfProgress, setShelfProgress] = useState<Record<string, ConceptProgress>>({})
  const nextFetchedRef = useRef(false)
  const seenRef = useRef<string[]>([])
  const bottomRef = useRef<HTMLDivElement>(null)
  const questionCardRef = useRef<HTMLDivElement>(null)
  const answerFormRef = useRef<HTMLFormElement>(null)
  const answerFocusRef = useRef<MathFocusHandle | null>(null)
  const practiceStartedRef = useRef(false)
  const reduceMotion = useRef(false)
  const lastEntry = entries[entries.length - 1]
  const lastEntryKey = lastEntry?.key
  const lastEntryKind = lastEntry?.kind
  const practiceStarted = entries.some(e => e.kind === 'q')
  const practiceBlocked = entries.some(e => e.kind === 'q' && e.blockedPrerequisiteIds !== undefined)
  const conceptGradingType = concepts.find(c => c.id === conceptId)?.grading_type
  // The served question carries the format its answer actually needs. The
  // concept record only names the grading type, which is why a numeric concept
  // whose answer is "5 R 2" used to get a digits-only keypad.
  const [servedFormat, setServedFormat] = useState<AnswerFormat | null>(null)
  const format = servedFormat ?? formatForGradingType(conceptGradingType)
  // Where "back" goes, and the label a learner reads. The label is resolved from the concept the
  // return path names, never printed as the raw id a storage key would show.
  const back = backTarget(returnStack)
  const returnTo = returnStack.length > 0 ? hrefConceptId(returnStack[0]) : null
  // A learner reading "← Back to frac.add.diff" is reading a storage key.
  const returnToLabel = returnTo ? (concepts.find(c => c.id === returnTo)?.label ?? returnTo) : null
  // The stack to put on any link that leaves this page, so this page is reachable again.
  const currentPath = learnConceptPath(conceptId)
  const prereqStack = pushReturn(returnStack, currentPath)

  useEffect(() => {
    reduceMotion.current = typeof window !== 'undefined' && window.matchMedia?.('(prefers-reduced-motion: reduce)').matches
  }, [])

  // Follow new cards after they mount, never answer edits or feedback updates.
  useEffect(() => {
    const options: ScrollIntoViewOptions = { behavior: reduceMotion.current ? 'auto' : 'smooth', block: 'nearest' }
    if (lastEntryKind === 'done') {
      bottomRef.current?.scrollIntoView(options)
      return
    }
    const card = questionCardRef.current
    if (lastEntryKind !== 'q' || !card) return
    answerFocusRef.current?.focus()
    const topPadding = parseFloat(getComputedStyle(document.documentElement).scrollPaddingTop) || 0
    const bottomMargin = parseFloat(getComputedStyle(card).scrollMarginBottom) || 0
    const availableHeight = (window.visualViewport?.height ?? window.innerHeight) - topPadding - bottomMargin
    const target = card.getBoundingClientRect().height <= availableHeight ? card : answerFormRef.current
    target?.scrollIntoView(options)
  }, [lastEntryKey, lastEntryKind])

  const refill = useCallback(async (diff: number, exclude: string[]): Promise<PracticeQuestion[]> => {
    const pr = await getLessonPractice(conceptId, 3, { seed: Date.now() % 100000, exclude, difficulty: diff, kp_ordinal: kpIndex + 1 })
    const fresh = pr.questions.filter(q => !exclude.includes(q.question))
    const answerable = fresh.length > 0 ? fresh : pr.questions
    const first = answerable[0]
    if (first?.answer_hint) {
      setServedFormat({ hint: first.answer_hint, inputMode: first.answer_input_mode === 'numeric' ? 'numeric' : 'text' })
    }
    return answerable
  }, [conceptId, kpIndex])

  const load = useCallback(async () => {
    practiceStartedRef.current = false
    setStartingPractice(false)
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
      const pr = await getLessonPractice(conceptId, 3, { seed: Date.now() % 100000, difficulty: 0.4, kp_ordinal: 1 })
      seenRef.current = pr.questions.map(q => q.question)
      setBuffer(pr.questions)
      const first = pr.questions[0]
      if (first?.answer_hint) {
        setServedFormat({ hint: first.answer_hint, inputMode: first.answer_input_mode === 'numeric' ? 'numeric' : 'text' })
      }
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
  // secondary action on the done card, never a silent fallback.
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

  // Re-fetch the ranked head on entering done or when re-entering a mastered concept:
  // navigation waits for the Continue click (deterministic for e2e).
  // A slow network must not strand the card: past the timeout the failure
  // branch (Practice again + Back to Profile) renders instead of a stale
  // spinner, and never a self-link.
  const isMastered = readiness?.state === 'mastered'
  const hasDone = entries.some(e => e.kind === 'done')
  useEffect(() => {
    if ((!hasDone && !isMastered) || nextFetchedRef.current) return
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
      // current page, a scroll-to-top dead end. That argument moved to the engine with the
      // rest of the decision.
      setNextShelf(next)
    }).catch(() => {
      nextFetchedRef.current = false
    }).finally(() => setNextLoading(false))
  }, [hasDone, isMastered, conceptId])

  // Only a real concept change may render a navigating <Link>. A same-concept
  // href (however many query params it carries) would leave the learner on the
  // page they just finished, so every Continue surface shares this guard
  // instead of trusting the server to exclude the current concept.
  const nextHeadId = nextShelf ? hrefConceptId(nextShelf.next.href) : null
  const canContinue = !nextLoading && !!nextShelf && nextHeadId !== null && nextHeadId !== conceptId

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
    window.setTimeout(fn, ms)
  }

  async function startPracticing() {
    if (practiceStartedRef.current) return
    practiceStartedRef.current = true
    setStartingPractice(true)
    setError('')
    const q = await takeNext(difficulty)
    setStartingPractice(false)
    if (!q) {
      practiceStartedRef.current = false
      setError('Could not load a question. Check your connection and try again.')
      return
    }
    const entry: QEntry = { key: nextKey(), kind: 'q', kpIndex, q, answer: '', checking: false, servedAt: Date.now(), feedback: null }
    setEntries(prev => [...prev, entry])
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
   * this is the second line of defence rather than the only one, but a client that can be
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
    if (!entry || !entry.answer.trim() || entry.checking || entry.feedback || entry.blockedPrerequisiteIds !== undefined) return
    setEntry(key, { checking: true, error: '' })
    const elapsed = Math.max(0.5, (Date.now() - entry.servedAt) / 1000)
    try {
      const res = await submitStudyAnswer(conceptId, entry.answer.trim(), elapsed, entry.q.question)
      if (res.ungraded) {
        setEntry(key, { error: 'The grader is temporarily unavailable. Your answer is kept. Press Check to retry.' })
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
            text: res.explanation ?? '',
            expected: res.expected_answer ?? '',
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
            text: res.explanation || '',
            expected: res.expected_answer ?? '',
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
      const blockedPrerequisiteIds = getTopicLockedPrerequisites(e)
      if (blockedPrerequisiteIds !== null) {
        setEntry(key, {
          blockedPrerequisiteIds,
          error: 'This topic now needs prerequisite practice before you can continue. This answer was not graded.',
        })
        return
      }
      // The server is the only grader. This never decides whether the answer
      // was right, a network fault or an expired server-side record is not a
      // miss, and grading locally would have made it one.
      //
      // Re-serve a rejected question on 409, including an expired server-side
      // record. If that fetch also fails, keep the answer and release Check.
      // A transport failure alone cannot tell us whether the server recorded
      // the answer.
      if (getErrorStatus(e) === 409) {
        try {
          const fresh = await refill(difficulty, seenRef.current)
          if (fresh.length > 0) {
            seenRef.current = [...seenRef.current, fresh[0].question].slice(-20)
            setEntry(key, { q: fresh[0], answer: '', servedAt: Date.now(), error: 'That question could not be graded. Here is a fresh one. Enter your answer and press Check.' })
            setBuffer(prev => [...prev, ...fresh.slice(1)].slice(0, 6))
            return
          }
        } catch {
          // Replacement failures must also release the submit control.
        }
        setEntry(key, { error: 'That question could not be graded and no replacement could be loaded. Press Check to try again or reload the lesson.' })
        return
      }
      const status = getErrorStatus(e)
      const message = status === 401
        ? 'Your session ended. Sign in again to submit your answer.'
        : status === 429
          ? 'Too many requests. Wait a moment and press Check again. Your answer is kept.'
          : status !== undefined
            ? `The server could not grade this answer (HTTP ${status}). Your answer is kept. Press Check to retry.`
            : 'Could not reach the server to grade that answer. Check your connection and press Check to retry. Your answer is kept.'
      setEntry(key, { error: message })
    } finally {
      setEntry(key, { checking: false })
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
  // this renders its reason and a path into each prerequisite's own reference page.
  if (ineligible && entries.length === 0) {
    const title = concepts.find(c => c.id === conceptId)?.label ?? conceptId
    return (
      <div className="max-w-2xl mx-auto">
        <PrerequisitePanel
          topicTitle={title}
          prerequisites={ineligible.prerequisites ?? []}
          returnStack={prereqStack}
        />
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

  // The server's authoritative prerequisite list, with a fallback to the older `missing`/`weak`
  // buckets for a payload that predates `prerequisites`. Every entry is openable: a prerequisite's
  // mastery never gates access to its reference material (ADR-054). Review routing for a mastered
  // concept happens on the study page's Practise control, not by withholding a link here.
  const prerequisites: EligibilityPrereq[] = (() => {
    if (!readiness) return []
    if (readiness.prerequisites && readiness.prerequisites.length > 0) return readiness.prerequisites
    return [...(readiness.missing ?? []), ...(readiness.weak ?? [])].map(p => ({
      conceptId: p.id,
      title: p.label,
      state: prereqStateFromStatus(p.status),
      met: countsAsMastered({ status: p.status }),
    }))
  })()
  const unmetPrereqs = prerequisites.filter(p => !p.met)
  // `met` is eligibility, not learning: the server marks a prerequisite met when
  // the learner may proceed past it, which includes one that is unlocked but
  // untouched. The panel prints the state per row, so the intro must not claim
  // the learner has "already learned" a row that reads Not started.
  const unstartedPrereqs = prerequisites.filter(p => p.met && (p.state === 'unlocked' || p.state === 'upcoming'))
  const dueReviewPrereqs = prerequisites.filter(p => p.met && p.state === 'due_for_review')
  const prereqIntro = (() => {
    const names = (unmetPrereqs.length > 0 ? unmetPrereqs : dueReviewPrereqs)
      .slice(0, 3).map(p => p.title).join(', ')
    if (unmetPrereqs.length === 0 && unstartedPrereqs.length > 0) {
      const unstarted = unstartedPrereqs.slice(0, 3).map(p => p.title).join(', ')
      return `These prerequisites are not started yet: ${unstarted}. Study them first or continue anyway.`
    }
    if (unmetPrereqs.length === 0 && dueReviewPrereqs.length > 0) {
      return `You have already learned ${names}, a short check is due or continue anyway.`
    }
    if (unmetPrereqs.length === 0) {
      return 'You have already learned what this builds on. You can revisit it any time.'
    }
    return `This builds on ${names}. A quick study helps or continue anyway.`
  })()
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
          <span title="How you're doing on this concept, from your recent answers. Progress only, mastery is decided separately.">
            · {band}
          </span>
        )}
      </div>

      {prerequisites.length > 0 && !bannerDismissed && !practiceBlocked && (
        <PrerequisitePanel
          prerequisites={prerequisites}
          returnStack={prereqStack}
          heading={`Prerequisites (${prerequisites.length})`}
          intro={prereqIntro}
          onDismiss={() => setBannerDismissed(true)}
        />
      )}

      {error && <p className="mb-4 font-mono text-xs text-red-400">{error}</p>}

      {isMastered && (
        <div className="border border-mathua-green-faint bg-mathua-surface p-5 mb-4">
          <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3">
            <div>
              <p className="font-mono text-xs text-green-400">✓ You have mastered this concept</p>
              <p className="mt-1 font-mono text-[11px] text-mathua-secondary">
                You can move on to the next recommended concept or practice this lesson again below.
              </p>
            </div>
            <div className="flex flex-wrap items-center gap-2 shrink-0">
              {canContinue && nextShelf && (
                <Link
                  href={nextShelf.next.href}
                  className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue-faint px-4 py-2 font-mono text-xs inline-flex items-center min-h-[40px]"
                >
                  Continue to Next Concept ({nextShelf.next.title}) →
                </Link>
              )}
              {nextLoading && (
                <span className="font-mono text-xs text-mathua-muted px-4 py-2 inline-flex items-center min-h-[40px]">
                  Finding next concept…
                </span>
              )}
              {!nextLoading && !canContinue && (
                <Link
                  href="/learn"
                  className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue-faint px-4 py-2 font-mono text-xs inline-flex items-center min-h-[40px]"
                >
                  Continue to Next Concept →
                </Link>
              )}
              <Link
                href="/profile"
                className="border border-mathua-border px-4 py-2 font-mono text-xs text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue inline-flex items-center min-h-[40px]"
              >
                Profile
              </Link>
            </div>
          </div>
        </div>
      )}

      {/* The reference layer, in the flow.
          This replaces the "Reference" link the done card used to carry to `/study`: the same
          material, reachable without leaving the questions, and impossible to wander off in, 
          there is nothing in it that navigates. Collapsed by default, so it never stands between
          a learner and the first question; the worked example above it is the intended first
          read. It renders no mastery state and awards nothing, which is the point: reading
          explains, and only an answer moves the concept. */}
      <ConceptReference conceptId={conceptId} />

      <div className="space-y-4">
        {entries.map(e => {
          if (e.kind === 'intro') {
            const ikp = kps[e.kpIndex]
            return (
              <div key={e.key} className="border border-mathua-border bg-mathua-surface p-6">
                <p className="font-mono text-[11px] uppercase tracking-wider text-mathua-muted">{ikp?.label || 'Learn'}</p>
                <KatexContent className="mt-3 text-sm text-mathua-primary whitespace-pre-wrap">{ikp?.worked_example || 'Worked example unavailable, start practicing instead.'}</KatexContent>
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
                    example stays on screen as scrollback, nothing is skipped.
                    Matches the kpdiv "Next: … (worked example)" label below, and
                    MA has no skip affordance at all: the worked example is the
                    scaffolding for the question that follows (example-problem
                    pair), so the label must not invite bypassing it. */}
                <div className="mt-5 flex flex-wrap items-center gap-3">
                  {!practiceStarted && (
                    <button type="button" onClick={startPracticing} disabled={startingPractice} className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue-faint px-6 py-2.5 font-mono text-xs min-h-[44px] disabled:opacity-50">
                      {isMastered ? 'Practice again →' : 'Next →'}
                    </button>
                  )}
                  {isMastered && canContinue && nextShelf && (
                    <Link
                      href={nextShelf.next.href}
                      className="border border-mathua-border text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue px-6 py-2.5 font-mono text-xs min-h-[44px] inline-flex items-center"
                    >
                      Continue to Next Concept →
                    </Link>
                  )}
                </div>
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
                <p className="font-mono text-xs text-red-400">Stepping down a level, easier question below.</p>
                <div className="mt-3 flex flex-wrap gap-2">
                  {/* A prerequisite is reference here too: study the material rather than being
                      sent to a tutorial under the name of a check. Practising it is a separate,
                      eligibility-gated step on that concept's own page. */}
                  {prerequisites.slice(0, 2).map(p => (
                    <Link key={p.conceptId} href={withReturnStack(studyConceptPath(p.conceptId), prereqStack)} className="border border-mathua-border px-4 py-2 font-mono text-[11px] text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue inline-flex items-center">Study {p.title}</Link>
                  ))}
                  <Link href="/profile" className="border border-mathua-border px-4 py-2 font-mono text-[11px] text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue inline-flex items-center">Other task →</Link>
                </div>
              </div>
            )
          }
          if (e.kind === 'done') {
            const upcoming: LockedSuccessor[] = upcomingLocked(catalogEntries(), shelfProgress, conceptId)
            return (
              <div key={e.key} className="border border-mathua-green-faint bg-mathua-surface p-6">
                <p className="font-mono text-xs text-green-400">✓ Complete, {totalCorrect}/{totalAnswered} correct · +{totalXP} XP · {band}</p>
                {e.note && <p role="status" className="mt-2 font-mono text-[11px] text-mathua-secondary">{e.note}</p>}
                <p className="mt-2 font-mono text-[11px] text-mathua-secondary">Scroll up to review anything. Reviews are scheduled automatically.</p>
                <div className="mt-4 flex flex-wrap gap-2">
                  {back && <Link href={back.href} className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue-faint px-5 py-2 font-mono text-xs inline-flex items-center min-h-[40px]">← Back to {returnToLabel}</Link>}
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
                      Coming up, unlocks once this concept is mastered
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
          const blocked = qe.blockedPrerequisiteIds !== undefined
          return (
            <div key={qe.key} ref={qe.key === lastEntryKey ? questionCardRef : undefined} className={`clear-of-tabs border bg-mathua-surface p-6 ${qe.feedback?.correct ? 'border-mathua-green-faint' : qe.feedback ? 'border-mathua-red-faint' : 'border-mathua-border'}`}>
              <KatexContent className="text-sm text-mathua-primary font-mono whitespace-pre-wrap">{stemWithoutChoices(qe.q.question)}</KatexContent>
              <ChoiceOptions question={qe.q.question} value={qe.answer} onPick={v => setEntry(qe.key, { answer: v })} disabled={locked || blocked || qe.checking} />
              {!locked && !blocked && (
                <form ref={qe.key === lastEntryKey ? answerFormRef : undefined} onSubmit={ev => { ev.preventDefault(); handleCheck(qe.key) }} className="clear-of-tabs mt-3 flex flex-col sm:flex-row gap-2">
                  <label htmlFor={`learn-answer-${qe.key}`} className="sr-only">Your answer</label>
                  <MathAnswerInput
                    id={`learn-answer-${qe.key}`}
                    focusRef={qe.key === lastEntryKey ? answerFocusRef : undefined}
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
              {!locked && !blocked && <p className="mt-2 font-mono text-[11px] text-mathua-muted">{format.hint}</p>}
              {!locked && qe.error && <p role="alert" className="mt-2 font-mono text-xs text-red-400">{qe.error}</p>}
              {blocked && (
                <div>
                  <p className="mt-2 font-mono text-xs text-mathua-secondary">Your answer: {qe.answer}</p>
                  <PrerequisitePanel
                    heading="Review prerequisites to continue"
                    intro="Study the prerequisites below, then use their practice action. Return here when you are ready."
                    prerequisites={(qe.blockedPrerequisiteIds?.length ? qe.blockedPrerequisiteIds : prerequisites.flatMap(p => p.met ? [] : [p.conceptId])).map(id => ({
                      conceptId: id,
                      title: prerequisites.find(p => p.conceptId === id)?.title ?? concepts.find(c => c.id === id)?.label ?? id,
                      state: prerequisites.find(p => p.conceptId === id)?.state ?? '',
                      met: false,
                    }))}
                    returnStack={prereqStack}
                  />
                  <button type="button" onClick={practiceAgain} className="mt-3 border border-mathua-border px-4 py-2 font-mono text-xs text-mathua-secondary hover:text-mathua-blue">Restart lesson</button>
                </div>
              )}
              <div className="mt-2 flex justify-end">
                <ReportButton
                  key={`learn-${qe.key}`}
                  conceptId={conceptId}
                  kind="question"
                  question={qe.q.question}
                  source="practice"
                  answerGiven={qe.answer}
                  gradingType={conceptGradingType}
                  answerHint={format.hint}
                  answerInputMode={format.inputMode}
                />
              </div>
              {locked && qe.feedback && (
                <div className="mt-3">
                  <p className={`font-mono text-xs ${qe.feedback.correct ? 'text-green-400' : 'text-red-400'}`}>
                    {qe.feedback.correct ? `✓ ${qe.feedback.xp ? `+${qe.feedback.xp} XP` : 'Correct'}` : '✗ Not quite.'}{' '}
                    <span className="text-mathua-muted">you answered “{qe.answer}”</span>
                  </p>
                  <p className="mt-1.5 font-mono text-xs text-mathua-primary">
                    Answer: <KatexContent>{qe.feedback.expected}</KatexContent>
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
