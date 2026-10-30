'use client'

import { Fragment, useCallback, useEffect, useMemo, useRef, useState } from 'react'
import Link from 'next/link'
import { useRouter } from 'next/navigation'
import AppShell from '../../components/app-shell'
import SectionHeader from '../../components/SectionHeader'
import Footer from '../../components/Footer'
import Loading from '../../components/Loading'
import PrerequisitePathDialog, { type PrereqDialogStatus } from '../../components/PrerequisitePathDialog'
import { pushReturn, studyConceptPath } from '../../lib/learnReturn'
import { useTheme } from '../../hooks/useTheme'
import {
  getCurriculumDomains,
  getCurriculumDomain,
  getLessonEligibility,
  type DomainSummary,
  type DomainDetail,
  type TopicSummary,
  type EligibilityPrereq,
} from '../../lib/api'
import { DOMAIN_ORDER, domainLabel } from '../../lib/graphDomains'

/**
 * The curriculum browser: domain -> topic, driven entirely by the server's
 * curriculum endpoints. The client never derives a state from raw progress; it
 * renders `state` and asks the server before entering the learning loop.
 *
 * Domains enter from `/profile`'s "Choose what to learn", and they are also a
 * bottom tab, so this is reachable without going through Settings or Study.
 */

function stateLabel(state: string): string {
  switch (state) {
    case 'mastered':
      return 'Mastered'
    case 'due_for_review':
      return 'Due for review'
    case 'in_progress':
      return 'In progress'
    case 'unlocked':
      return 'Available'
    case 'upcoming':
      return 'Upcoming'
    case 'locked':
      return 'Locked'
    default:
      return state
  }
}

function available(row: DomainSummary): number {
  return row.unlockedCount + row.inProgressCount
}

export default function DomainsPage() {
  const { mounted } = useTheme()
  const router = useRouter()

  const [domains, setDomains] = useState<DomainSummary[] | null>(null)
  const [selected, setSelected] = useState<string | null>(null)
  const [detail, setDetail] = useState<DomainDetail | null>(null)
  const [detailLoading, setDetailLoading] = useState(false)
  const [busy, setBusy] = useState<string | null>(null)
  const [dialog, setDialog] = useState<{ conceptId: string; title: string } | null>(null)
  const [prereqs, setPrereqs] = useState<EligibilityPrereq[]>([])
  const [dlgStatus, setDlgStatus] = useState<PrereqDialogStatus>('loading')
  const [dlgError, setDlgError] = useState('')
  const [error, setError] = useState('')
  // Monotonic request id: a late response for concept A must never overwrite
  // the dialog now showing concept B.
  const seqRef = useRef(0)
  // Element to return focus to when the dialog closes.
  const triggerRef = useRef<HTMLButtonElement | null>(null)

  const loadDomains = useCallback(() => {
    setError('')
    setDomains(null)
    getCurriculumDomains()
      .then(setDomains)
      .catch(() => {
        setDomains([])
        setError("Couldn't load the curriculum.")
      })
  }, [])

  useEffect(() => {
    if (!mounted) return
    loadDomains()
  }, [mounted, loadDomains])

  const ordered = useMemo(() => {
    if (!domains) return []
    const byId = new Map(domains.map((d) => [d.id, d]))
    const out: DomainSummary[] = []
    for (const id of DOMAIN_ORDER) {
      const d = byId.get(id)
      if (d) out.push(d)
    }
    // Any domain the corpus grows that DOMAIN_ORDER does not yet name still
    // appears, rather than silently vanishing from the picker.
    for (const d of domains) if (!out.includes(d)) out.push(d)
    return out
  }, [domains])

  // A group heading earns its space only when the subject has more than one
  // row. With a single calculus domain the row says "Calculus I" on its own.
  const calcCount = ordered.filter((d) => d.id.startsWith('calc')).length

  const openDomain = useCallback(async (id: string) => {
    setSelected(id)
    setDialog(null)
    setDetail(null)
    setDetailLoading(true)
    setError('')
    try {
      setDetail(await getCurriculumDomain(id))
    } catch {
      setError("Couldn't load domain details.")
    } finally {
      setDetailLoading(false)
    }
  }, [])

  async function baseEligibility(conceptId: string): Promise<boolean | EligibilityPrereq[] | null> {
    const el = await getLessonEligibility(conceptId)
    if (el.eligible === false) return el.prerequisites ?? []
    return true
  }

  async function chooseTopic(topic: TopicSummary) {
    setError('')
    setDialog(null)
    if (topic.state === 'mastered' || topic.state === 'due_for_review') {
      router.push('/review')
      return
    }
    setBusy(topic.conceptId)
    try {
      const result = await baseEligibility(topic.conceptId)
      if (result === true) {
        router.push(`/learn?concept=${encodeURIComponent(topic.conceptId)}`)
      } else if (Array.isArray(result)) {
        seqRef.current++
        setDialog({ conceptId: topic.conceptId, title: topic.title })
        setPrereqs(result)
        setDlgStatus(result.length === 0 ? 'empty' : 'ready')
        setDlgError('')
      }
    } catch {
      setError("Couldn't start this topic.")
    } finally {
      setBusy(null)
    }
  }

  // Selecting a locked topic explains the prerequisite path instead of starting.
  // The dialog opens immediately in a loading state so the action is visible
  // regardless of where the row sits in a long list; the fetch then fills it.
  // Only the latest request may write, so rapid clicks on different rows can
  // never leave stale prerequisites on screen.
  async function explainTopic(topic: TopicSummary, trigger?: HTMLButtonElement | null) {
    setError('')
    setDlgError('')
    triggerRef.current = trigger ?? null
    const seq = ++seqRef.current
    setDialog({ conceptId: topic.conceptId, title: topic.title })
    setPrereqs([])
    setDlgStatus('loading')
    setBusy(topic.conceptId)
    try {
      const el = await getLessonEligibility(topic.conceptId)
      if (seqRef.current !== seq) return
      const list = el.prerequisites ?? []
      setPrereqs(list)
      setDlgStatus(list.length === 0 ? 'empty' : 'ready')
    } catch {
      if (seqRef.current !== seq) return
      setDlgStatus('error')
      setDlgError("Couldn't load this topic's prerequisites.")
    } finally {
      if (seqRef.current === seq) setBusy(null)
    }
  }

  function retryDialog() {
    if (!dialog) return
    void explainTopic(
      { conceptId: dialog.conceptId, title: dialog.title, state: 'locked', depth: 0 },
      triggerRef.current,
    )
  }

  function closeDialog() {
    seqRef.current++
    setDialog(null)
    setBusy(null)
    triggerRef.current?.focus()
    triggerRef.current = null
  }

  if (!mounted || domains === null) {
    return (
      <AppShell>
        <Loading label="Loading curriculum" full />
        <Footer />
      </AppShell>
    )
  }

  const selectedRow = ordered.find((d) => d.id === selected) ?? null

  return (
    <AppShell>
      <main className="mx-auto w-full max-w-[820px] min-w-0 px-4 sm:px-6 py-8 sm:py-12 overflow-x-hidden">
        {selected === null ? (
          <>
            <SectionHeader label="Curriculum" title="Choose what to learn" />
            <p className="font-mono text-[11px] text-mathua-secondary mb-6">
              Browse a domain and pick a topic. Eligibility is checked by Mathua; you cannot
              start a topic whose prerequisites are unmet.
            </p>

            {error && (
              <div className="mb-6 p-4 border border-mathua-border bg-mathua-surface">
                <p role="alert" className="font-mono text-xs text-mathua-secondary mb-3">
                  {error}
                </p>
                <button
                  type="button"
                  onClick={loadDomains}
                  className="font-mono text-xs text-mathua-primary border border-mathua-border px-3 py-1 hover:bg-mathua-blue-faint transition-colors"
                >
                  Retry loading curriculum
                </button>
              </div>
            )}

            <ul className="border border-mathua-border divide-y divide-mathua-border bg-mathua-surface">
              {ordered.map((d, i) => {
                // The three calculus domains are one subject broken into stages.
                // The heading keeps the top-level picker reading as a list of
                // subjects without hiding the per-course rows that carry the
                // counts and the diagnostic scope.
                const startsCalculus =
                  calcCount > 1 &&
                  d.id.startsWith('calc') &&
                  (i === 0 || !ordered[i - 1].id.startsWith('calc'))
                return (
                  <Fragment key={d.id}>
                    {startsCalculus && (
                      <li className="px-4 py-2 font-mono text-[10px] uppercase tracking-wider text-mathua-muted">
                        Calculus
                      </li>
                    )}
                    <li>
                      <button
                        type="button"
                        onClick={() => void openDomain(d.id)}
                        className="w-full text-left px-4 py-3 hover:bg-mathua-blue-faint transition-colors"
                      >
                        <span className="flex items-baseline justify-between gap-3 min-w-0">
                          <span className="font-mono text-xs text-mathua-primary truncate">
                            {domainLabel(d.id)}
                          </span>
                          <span className="shrink-0 font-mono text-[11px] text-mathua-secondary">
                            {available(d)} / {d.conceptCount} concepts available
                          </span>
                        </span>
                      </button>
                    </li>
                  </Fragment>
                )
              })}
            </ul>
          </>
        ) : (
          <>
            <button
              type="button"
              onClick={() => { setSelected(null); setDetail(null); setDialog(null); setError('') }}
              className="font-mono text-xs text-mathua-secondary hover:text-mathua-blue"
            >
              ← Domains
            </button>
            <div className="mt-4">
              <SectionHeader
                label="Curriculum"
                title={selectedRow ? domainLabel(selectedRow.id) : selected}
              />
            </div>

            {selectedRow && (
              <p className="font-mono text-[11px] text-mathua-secondary mb-6">
                Your progress: {available(selectedRow)} / {selectedRow.conceptCount} concepts available
                {selectedRow.masteredCount > 0 && ` · ${selectedRow.masteredCount} mastered`}
              </p>
            )}

            {error && !detailLoading && !detail && (
              <div className="mb-6 p-4 border border-mathua-border bg-mathua-surface">
                <p role="alert" className="font-mono text-xs text-mathua-secondary mb-3">
                  {error}
                </p>
                <button
                  type="button"
                  onClick={() => void openDomain(selected)}
                  className="font-mono text-xs text-mathua-primary border border-mathua-border px-3 py-1 hover:bg-mathua-blue-faint transition-colors"
                >
                  Retry loading domain
                </button>
              </div>
            )}

            {detailLoading && <p className="font-mono text-xs text-mathua-muted">Loading topics…</p>}

            {!detailLoading && detail && detail.topics.length === 0 && (
              <p className="font-mono text-xs text-mathua-secondary">
                No topics are currently available here. Your prerequisite path will unlock
                topics as you demonstrate mastery.
              </p>
            )}

            {!detailLoading && detail && detail.topics.length > 0 && (
              <>
                {!detail.topics.some((t) =>
                  ['unlocked', 'in_progress', 'mastered', 'due_for_review'].includes(t.state),
                ) && (
                  <p className="mb-4 border border-mathua-border bg-mathua-surface p-3 font-mono text-[11px] text-mathua-secondary">
                    No topics are currently available here. Your prerequisite path will unlock
                    topics as you demonstrate mastery.
                  </p>
                )}
                <ul className="border border-mathua-border divide-y divide-mathua-border bg-mathua-surface">
                  {detail.topics.map((topic) => {
                  const learnable = topic.state === 'unlocked' || topic.state === 'in_progress'
                  const review = topic.state === 'mastered' || topic.state === 'due_for_review'
                  return (
                    <li key={topic.conceptId} className="px-4 py-3 min-w-0">
                      <div className="flex flex-col sm:flex-row sm:items-center gap-2 sm:gap-4 min-w-0">
                        <div className="min-w-0 flex-1">
                          <div className="font-mono text-xs text-mathua-primary break-words [overflow-wrap:anywhere]">
                            {topic.title}
                          </div>
                          {/* State is text first, so it survives without color. */}
                          <div className="font-mono text-[10px] uppercase tracking-wider text-mathua-muted">
                            {stateLabel(topic.state)}
                          </div>
                        </div>
                        {learnable && (
                          <button
                            type="button"
                            disabled={busy === topic.conceptId}
                            onClick={() => void chooseTopic(topic)}
                            className="w-full sm:w-auto shrink-0 border border-mathua-blue text-mathua-blue hover:bg-mathua-blue-faint px-4 py-2 font-mono text-[11px] min-h-[36px] inline-flex items-center justify-center disabled:opacity-60"
                          >
                            {busy === topic.conceptId ? 'Checking…' : 'Learn'}
                          </button>
                        )}
                        {review && (
                          <Link
                            href="/review"
                            className="w-full sm:w-auto shrink-0 border border-mathua-blue text-mathua-blue hover:bg-mathua-blue-faint px-4 py-2 font-mono text-[11px] min-h-[36px] inline-flex items-center justify-center"
                          >
                            Review
                          </Link>
                        )}
                        {!learnable && !review && (
                          <button
                            type="button"
                            onClick={(e) => void explainTopic(topic, e.currentTarget)}
                            className="w-full sm:w-auto shrink-0 border border-mathua-border text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue px-4 py-2 font-mono text-[11px] min-h-[36px] inline-flex items-center justify-center"
                          >
                            View prerequisite path
                          </button>
                        )}
                      </div>
                    </li>
                  )
                })}
              </ul>
              </>
            )}

            <PrerequisitePathDialog
              open={dialog !== null}
              targetTitle={dialog?.title ?? ''}
              prerequisites={prereqs}
              status={dlgStatus}
              error={dlgError}
              returnStack={dialog ? pushReturn([], studyConceptPath(dialog.conceptId)) : []}
              onClose={closeDialog}
              onRetry={retryDialog}
            />

            {error && detail && (
              <p role="alert" className="mt-4 font-mono text-[11px] text-mathua-secondary">
                {error}
              </p>
            )}
          </>
        )}
      </main>
      <Footer />
    </AppShell>
  )
}
