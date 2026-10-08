'use client'

import { useCallback, useEffect, useMemo, useState } from 'react'
import Link from 'next/link'
import { useRouter } from 'next/navigation'
import Header from '../../components/Header'
import BottomTabs from '../../components/BottomTabs'
import SectionHeader from '../../components/SectionHeader'
import Footer from '../../components/Footer'
import Loading from '../../components/Loading'
import PrerequisitePanel from '../../components/PrerequisitePanel'
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
  const [blocked, setBlocked] = useState<{ title: string; prerequisites: EligibilityPrereq[] } | null>(null)
  const [error, setError] = useState('')

  useEffect(() => {
    if (!mounted) return
    getCurriculumDomains()
      .then(setDomains)
      .catch(() => {
        setDomains([])
        setError("Couldn't load the curriculum.")
      })
  }, [mounted])

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

  const openDomain = useCallback(async (id: string) => {
    setSelected(id)
    setBlocked(null)
    setDetail(null)
    setDetailLoading(true)
    try {
      setDetail(await getCurriculumDomain(id))
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
    setBlocked(null)
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
        setBlocked({ title: topic.title, prerequisites: result })
      }
    } catch {
      setError("Couldn't start this topic.")
    } finally {
      setBusy(null)
    }
  }

  // Selecting a locked topic explains the prerequisite path instead of starting.
  async function explainTopic(topic: TopicSummary) {
    setError('')
    setBlocked(null)
    setBusy(topic.conceptId)
    try {
      const el = await getLessonEligibility(topic.conceptId)
      setBlocked({ title: topic.title, prerequisites: el.prerequisites ?? [] })
    } catch {
      setError("Couldn't load this topic's prerequisites.")
    } finally {
      setBusy(null)
    }
  }

  if (!mounted || domains === null) return <Loading label="Loading curriculum" full />

  const selectedRow = ordered.find((d) => d.id === selected) ?? null

  return (
    <>
      <Header />
      <main className="mx-auto w-full max-w-[820px] min-w-0 px-4 sm:px-6 py-8 sm:py-12 overflow-x-hidden">
        {selected === null ? (
          <>
            <SectionHeader label="Curriculum" title="Choose what to learn" />
            <p className="font-mono text-[11px] text-mathua-secondary mb-6">
              Browse a domain and pick a topic. Eligibility is checked by Mathua; you cannot
              start a topic whose prerequisites are unmet.
            </p>

            {error && (
              <p role="alert" className="mb-4 font-mono text-[11px] text-mathua-secondary">
                {error}
              </p>
            )}

            <ul className="border border-mathua-border divide-y divide-mathua-border bg-mathua-surface">
              {ordered.map((d) => (
                <li key={d.id}>
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
              ))}
            </ul>
          </>
        ) : (
          <>
            <button
              type="button"
              onClick={() => { setSelected(null); setDetail(null); setBlocked(null) }}
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
                            disabled={busy === topic.conceptId}
                            onClick={() => void explainTopic(topic)}
                            className="w-full sm:w-auto shrink-0 border border-mathua-border text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue px-4 py-2 font-mono text-[11px] min-h-[36px] inline-flex items-center justify-center disabled:opacity-60"
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

            {blocked && (
              <PrerequisitePanel
                topicTitle={blocked.title}
                prerequisites={blocked.prerequisites}
                onDismiss={() => setBlocked(null)}
              />
            )}

            {error && (
              <p role="alert" className="mt-4 font-mono text-[11px] text-mathua-secondary">
                {error}
              </p>
            )}
          </>
        )}
      </main>
      <Footer />
      <BottomTabs />
    </>
  )
}
