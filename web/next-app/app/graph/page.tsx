'use client'

import { useTheme } from '../../hooks/useTheme'
import { useState, useEffect, useMemo, useCallback, useRef, Suspense } from 'react'
import Link from 'next/link'
import { useRouter, useSearchParams } from 'next/navigation'
import AppShell from '../../components/app-shell'
import dynamic from 'next/dynamic'
import SectionHeader from '../../components/SectionHeader'
import ProgressSummary from '../../components/ProgressSummary'
import Footer from '../../components/Footer'
import GoalPicker from '../../components/GoalPicker'
import { concepts as conceptsData } from '../../lib/conceptData'
import { getScores, getGraph, getProgress, getWeaknesses, getGoal as fetchGoal, setGoal as saveGoal, clearGoal as deleteGoal, healthCheck, type ConceptProgress, type GraphRes, type Scores } from '../../lib/api'
import { getUserInfo } from '../../lib/auth'
import { useAuthState } from '../../hooks/useAuthState'
import { deriveStatuses } from '../../lib/graphStatus'
import { DOMAIN_ORDER, domainLabel } from '../../lib/graphDomains'
import { goalProgress as computeGoalProgress, nextGoalStep } from '../../lib/goalScope'
import { createInflightCache, createSessionCache } from '../../lib/requestCache'
import Loading from '../../components/Loading'

const graphLoadingStyle: React.CSSProperties = {
  height: 'clamp(480px, calc(100dvh - 360px), 1400px)',
  background: 'var(--surface)',
  borderRadius: '0',
  border: '0.5px solid var(--border)',
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  width: '100%',
  maxWidth: '100%',
}

const conceptGraphFlowChunk = import('../../components/ConceptGraphFlow')

const ConceptGraphFlow = dynamic(
  () => conceptGraphFlowChunk,
  {
    ssr: false,
    loading: () => (
      <div style={graphLoadingStyle}>
        <Loading label="LOADING GRAPH" />
      </div>
    ),
  }
)

const fallbackConcepts = conceptsData.map((c) => ({
  id: c.id, label: c.label, domain: c.domain, prerequisites: c.prerequisites,
}))


// Two caches, because the two kinds of data on this page have opposite lifetimes. See
// lib/requestCache.ts for why they were one helper and what that cost.
const oncePerSession = createSessionCache()
const dedupeInFlight = createInflightCache()

const GOAL_STORAGE_KEY = 'mathua_graph_goal'

function readLocalGoal(): string | null {
  if (typeof window === 'undefined') return null
  try {
    return window.localStorage.getItem(GOAL_STORAGE_KEY)
  } catch {
    return null
  }
}

function writeLocalGoal(id: string | null): void {
  if (typeof window === 'undefined') return
  try {
    if (id) window.localStorage.setItem(GOAL_STORAGE_KEY, id)
    else window.localStorage.removeItem(GOAL_STORAGE_KEY)
  } catch {
    // Storage can be unavailable in private modes. The account remains authoritative.
  }
}

function GraphContent() {
  const { theme, mounted } = useTheme()
  const { replace } = useRouter()
  const searchParams = useSearchParams()
  const conceptParam = searchParams.get('concept')
  // ?domain= makes the filter linkable, so /domains can deep-link into one subject.
  const domainParam = searchParams.get('domain')
  const goalParam = searchParams.get('goal')
  const focusParam = searchParams.get('focus')
  const [graphData, setGraphData] = useState<GraphRes | null>(null)
  const [rawProgress, setRawProgress] = useState<Record<string, ConceptProgress>>({})
  const [weakByDomain, setWeakByDomain] = useState<Record<string, { id: string; label: string }[]> | undefined>(undefined)
  const [connected, setConnected] = useState(false)
  const [scores, setScores] = useState<Scores | null>(null)
  const [graphError, setGraphError] = useState(false)
  const [goalId, setGoalId] = useState<string | null>(goalParam)
  const [focusId, setFocusId] = useState<string | null>(focusParam)
  const [goalPathOnly, setGoalPathOnly] = useState(false)
  const [goalPickerOpen, setGoalPickerOpen] = useState(false)
  const [pendingLocalGoal, setPendingLocalGoal] = useState<string | null>(null)
  // Seeded from ?domain= so a deep link lands already filtered. The reset effect below
  // still runs: if the URL names a domain that does not exist, the graph shows nothing
  // filtered rather than silently ignoring the parameter.
  const [activeDomain, setActiveDomain] = useState<string | null>(domainParam)

  // Chip clicks write the filter into the URL so the view is shareable and survives a reload.
  // replace(), not push(): filtering is not navigation, and pushing would make Back walk through
  // every domain the learner glanced at.
  const setDomainFilter = (next: string | null) => {
    setActiveDomain(next)
    const q = new URLSearchParams(searchParams.toString())
    if (next) q.set('domain', next)
    else q.delete('domain')
    const qs = q.toString()
    replace(qs ? `/graph?${qs}` : '/graph', { scroll: false })
  }

  const replaceGraphQuery = useCallback((mutate: (q: URLSearchParams) => void) => {
    const q = new URLSearchParams(searchParams.toString())
    mutate(q)
    const qs = q.toString()
    replace(qs ? `/graph?${qs}` : '/graph', { scroll: false })
  }, [replace, searchParams])
  const [selectedId, setSelectedId] = useState<string | null>(conceptParam)
  const { loggedIn } = useAuthState()

  const loadGraph = useCallback(() => {
    setGraphError(false)
    oncePerSession('health', healthCheck).then(setConnected).catch(() => setConnected(false))
    oncePerSession('graph', getGraph).then(data => { setGraphData(data); setGraphError(false) }).catch(() => setGraphError(true))
  }, [])

  useEffect(() => {
    loadGraph()
  }, [loadGraph])

  useEffect(() => {
    if (!mounted) return
    const local = readLocalGoal()
    if (focusParam) {
      setFocusId(focusParam)
      setGoalId(null)
      return
    }
    if (goalParam) {
      setGoalId(goalParam)
      setFocusId(null)
      return
    }
    if (!loggedIn) {
      setGoalId(local)
      return
    }
    let cancelled = false
    fetchGoal()
      .then(account => {
        if (cancelled) return
        if (account) {
          setGoalId(account)
          setPendingLocalGoal(null)
          return
        }
        if (local) {
          setGoalId(null)
          setPendingLocalGoal(local)
          return
        }
        setGoalId(null)
        setPendingLocalGoal(null)
      })
      .catch(() => {
        if (!cancelled) setGoalId(local)
      })
    return () => {
      cancelled = true
    }
  }, [mounted, loggedIn, goalParam, focusParam])

  useEffect(() => {
    if (!loggedIn) return
    const user = getUserInfo()
    if (!user) return
    const id = user.student_id
    dedupeInFlight(`scores:${id}`, () => getScores(id)).then(setScores).catch(() => console.error('getScores failed'))
    dedupeInFlight(`progress:${id}`, () => getProgress(id)).then(setRawProgress).catch(() => console.error('getProgress failed'))
    dedupeInFlight('weaknesses', getWeaknesses).then(w => {
      const byDomain: Record<string, { id: string; label: string }[]> = {}
      for (const [domain, entries] of Object.entries(w.by_domain)) {
        byDomain[domain] = entries.map((e: any) => ({ id: e.id, label: e.label }))
      }
      setWeakByDomain(byDomain)
    }).catch(() => console.error('getWeaknesses failed'))
  }, [loggedIn])

  useEffect(() => {
    setSelectedId(conceptParam)
  }, [conceptParam])

  const concepts = useMemo(() => graphData
    ? graphData.nodes.map((n) => ({
        id: n.id,
        label: n.label,
        domain: n.domain,
        prerequisites: n.prerequisites,
      }))
    : [], [graphData])

  const statusSource = useMemo(() => (
    concepts.length > 0 ? concepts : fallbackConcepts
  ), [concepts])

  const conceptStatuses = useMemo(
    () => deriveStatuses(statusSource, rawProgress),
    [statusSource, rawProgress]
  )

  // Bar length per concept, straight from the server's `mastery_pct`.
  //
  // This used to be recomputed here as `streak / mastery_threshold.streak` against the
  // *bundled* corpus. Both halves of that were wrong: the streak resets to 1 on every
  // tier advance, so the bar dropped from 100% to ~10% three times while the learner was
  // doing everything right, and a bundled threshold is a second source of truth that
  // drifts from the server's whenever the corpus is rebuilt. The server now derives it
  // from the tier and the streak together, see `mastery.MasteryPct`.
  //
  // `decaying` still reads as a full bar, for the same reason `mastered` does: the
  // competence was demonstrated and what is outstanding is a retrieval check, not
  // attainment. Decay is carried by the node's colour alone, so "full bar" keeps meaning
  // "I learned this" and a learner never sees decay as losing it.
  const conceptProgress = useMemo(() => {
    const out: Record<string, number> = {}
    for (const [id, p] of Object.entries(rawProgress)) {
      const status = conceptStatuses[id]
      if (status === 'unseen' || status === 'locked') continue
      if (typeof p.mastery_pct !== 'number') continue
      out[id] = Math.min(1, Math.max(0, p.mastery_pct))
    }
    return out
  }, [rawProgress, conceptStatuses])

  const handleSelectionChange = useCallback((id: string | null) => {
    setSelectedId(id)
    replaceGraphQuery(q => {
      if (id) q.set('concept', id)
      else q.delete('concept')
    })
  }, [replaceGraphQuery])

  const handleSetGoal = useCallback(async (id: string) => {
    setGoalId(id)
    setFocusId(null)
    setGoalPathOnly(false)
    setActiveDomain(null)
    setPendingLocalGoal(null)
    writeLocalGoal(id)
    replaceGraphQuery(q => {
      q.set('goal', id)
      q.delete('focus')
      q.delete('domain')
    })
    if (loggedIn) {
      try {
        await saveGoal(id)
      } catch (error) {
        console.error('saveGoal failed', error)
      }
    }
  }, [loggedIn, replaceGraphQuery])

  const handleClearGoal = useCallback(async () => {
    setGoalId(null)
    setGoalPathOnly(false)
    writeLocalGoal(null)
    replaceGraphQuery(q => q.delete('goal'))
    if (loggedIn) {
      try {
        await deleteGoal()
      } catch (error) {
        console.error('deleteGoal failed', error)
      }
    }
  }, [loggedIn, replaceGraphQuery])

  const handleUseLocalGoal = useCallback(async () => {
    if (!pendingLocalGoal) return
    await handleSetGoal(pendingLocalGoal)
  }, [handleSetGoal, pendingLocalGoal])

  const handleDismissLocalGoal = useCallback(() => {
    setPendingLocalGoal(null)
  }, [])

  const handleFocus = useCallback((id: string) => {
    setFocusId(id)
    setGoalId(null)
    setGoalPathOnly(false)
    setActiveDomain(null)
    replaceGraphQuery(q => {
      q.set('focus', id)
      q.delete('goal')
      q.delete('domain')
    })
  }, [replaceGraphQuery])

  const handleClearFocus = useCallback(() => {
    setFocusId(null)
    replaceGraphQuery(q => q.delete('focus'))
  }, [replaceGraphQuery])

  const domains = useMemo(() => {
    const set = new Set<string>()
    for (const c of concepts) set.add(c.domain)
    return DOMAIN_ORDER.filter(d => set.has(d))
  }, [concepts])

  const onPathNodes = useMemo(() => {
    if (!activeDomain) return undefined
    // Single pass (replaces filter().map()).
    const ids: string[] = []
    for (const c of concepts) {
      if (c.domain === activeDomain) ids.push(c.id)
    }
    return ids
  }, [activeDomain, concepts])

  // O(1) node lookup for onNodeSelect (replaces find).
  const conceptById = useMemo(() => {
    const source = concepts.length > 0 ? concepts : fallbackConcepts
    return new Map(source.map(c => [c.id, c] as const))
  }, [concepts])

  // A deep link (?concept=) wins over the persisted/default domain filter, but
  // only once: afterwards the user's chip choice is authoritative.
  const deepLinkHandled = useRef(false)
  useEffect(() => {
    if (deepLinkHandled.current || concepts.length === 0) return
    // Not marked handled until the domain has been resolved against the canonical list, so a
    // ?domain= that has to be cleared still gets a chance to clear itself on a later pass.
    if (domainParam && !DOMAIN_ORDER.includes(domainParam as (typeof DOMAIN_ORDER)[number])) return
    deepLinkHandled.current = true
    // An unknown ?domain= is ignored rather than honoured: filtering on "bogus" matches no
    // concept, so the graph would render empty with no indication of why.
    if (activeDomain && !DOMAIN_ORDER.includes(activeDomain as (typeof DOMAIN_ORDER)[number])) {
      setActiveDomain(null)
    }
    if (!conceptParam || !activeDomain) return
    const c = conceptById.get(conceptParam)
    if (c && c.domain !== activeDomain) setActiveDomain(null)
  }, [conceptParam, domainParam, activeDomain, conceptById, concepts.length])

  const goalSource = concepts.length > 0 ? concepts : fallbackConcepts
  const goalConcept = goalId ? conceptById.get(goalId) ?? null : null
  const focusConcept = focusId ? conceptById.get(focusId) ?? null : null
  const goalSummary = useMemo(
    () => (goalId && goalConcept ? computeGoalProgress(goalId, goalSource, rawProgress) : null),
    [goalId, goalConcept, goalSource, rawProgress],
  )
  const nextStepId = useMemo(
    () => (goalId && goalConcept ? nextGoalStep(goalId, goalSource, rawProgress) : null),
    [goalId, goalConcept, goalSource, rawProgress],
  )
  const nextStepConcept = nextStepId ? conceptById.get(nextStepId) ?? null : null

  useEffect(() => {
    if (concepts.length === 0) return
    if (goalId && !goalConcept) {
      setGoalId(null)
      writeLocalGoal(null)
      replaceGraphQuery(q => q.delete('goal'))
    }
    if (focusId && !focusConcept) {
      setFocusId(null)
      replaceGraphQuery(q => q.delete('focus'))
    }
  }, [concepts.length, goalId, goalConcept, focusId, focusConcept, replaceGraphQuery])

  if (!mounted) {
    return (
      <AppShell>
        <div style={{ background: 'var(--bg)', minHeight: '100vh' }} />
        <Footer />
      </AppShell>
    )
  }

  return (
    <AppShell>
      <div className="w-full max-w-none px-4 sm:px-6 pb-[calc(80px+env(safe-area-inset-bottom))] lg:pb-0 overflow-x-hidden min-w-0">
      <section className="pt-8 min-w-0 overflow-hidden">
        <span className="flex justify-between mb-4">
          <button type="button" onClick={() => { if (window.history.length > 1) window.history.back() }} className="text-mathua-secondary text-sm hover:text-mathua-blue">
            ← Back
          </button>
        </span>
        <SectionHeader label="Your knowledge graph" title="Explore the concept map" />
        {connected && loggedIn && scores ? (
          <ProgressSummary scores={scores} weakByDomain={weakByDomain} />
        ) : connected && !loggedIn ? (
          <p className="text-mathua-secondary text-sm text-center max-w-[600px] mx-auto mt-4 mb-8">
            Sign in or{' '}
            <Link href="/learn" className="text-mathua-blue hover:underline">start learning</Link>
            {' '}to track your progress across the concept map.
          </p>
        ) : (
          <p className="text-mathua-secondary text-sm text-center max-w-[600px] mx-auto mt-4 mb-8">
            <Link href="/learn" className="text-mathua-blue hover:underline">Start learning</Link>
            {' '}to track your progress across the concept map.
          </p>
        )}
      </section>

      <div className="flex flex-col items-center gap-2 mb-6">
        <button
          type="button"
          onClick={() => setGoalPickerOpen(true)}
          className="w-full max-w-[320px] sm:w-auto font-mono text-[10px] uppercase tracking-wider px-4 min-h-[40px] border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white transition-colors"
        >
          {goalId ? 'Change goal' : 'Set a goal'}
        </button>
        {!goalId && (
          <p className="max-w-[420px] px-2 text-center text-[11px] text-mathua-secondary">
            Pick any concept and the map trims to the path that leads there.
          </p>
        )}
      </div>

      {!goalId && !focusId && domains.length > 0 && (
        <div className="flex flex-wrap gap-1.5 sm:gap-2 justify-center mb-6 min-w-0">
          <button
            type="button"
            onClick={() => setDomainFilter(null)}
            className={`font-mono text-[10px] uppercase px-3 min-h-[36px] py-1.5 border transition-colors min-w-0 truncate ${
              activeDomain === null
                ? 'bg-mathua-blue text-white border-mathua-blue'
                : 'border-mathua-border text-mathua-muted hover:text-mathua-blue hover:border-mathua-secondary'
            }`}
          >
            All
          </button>
          {domains.map(d => (
            <button
              type="button"
              key={d}
              onClick={() => setDomainFilter(activeDomain === d ? null : d)}
              className={`font-mono text-[10px] uppercase px-3 min-h-[36px] py-1.5 border transition-colors min-w-0 truncate ${
                activeDomain === d
                  ? 'bg-mathua-blue text-white border-mathua-blue'
                  : 'border-mathua-border text-mathua-muted hover:text-mathua-blue hover:border-mathua-secondary'
              }`}
            >
              {domainLabel(d)}
            </button>
          ))}
        </div>
      )}

      <div className="my-6 sm:my-8 w-full max-w-full min-w-0 overflow-hidden">
        {pendingLocalGoal && !goalId && (
          <div className="mb-6 border border-mathua-border bg-mathua-surface px-4 py-3">
            <p className="text-sm text-mathua-primary">
              This device has a saved goal: {conceptById.get(pendingLocalGoal)?.label ?? pendingLocalGoal}.
            </p>
            <div className="mt-3 flex flex-col gap-2 sm:flex-row">
              <button type="button" onClick={handleUseLocalGoal} className="w-full sm:w-auto font-mono text-[10px] uppercase tracking-wider px-3 min-h-[36px] border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white transition-colors">Use this goal</button>
              <button type="button" onClick={handleDismissLocalGoal} className="w-full sm:w-auto font-mono text-[10px] uppercase tracking-wider px-3 min-h-[36px] border border-mathua-border text-mathua-muted hover:text-mathua-primary transition-colors">Not now</button>
            </div>
          </div>
        )}

        {goalId && goalConcept && goalSummary && (
          <div className="mb-6 border border-mathua-blue bg-mathua-surface px-4 py-3">
            <div className="flex flex-col gap-2 sm:flex-row sm:items-start sm:justify-between sm:gap-3">
              <div className="min-w-0">
                <p className="font-mono text-[10px] uppercase tracking-wider text-mathua-muted">Goal</p>
                <p className="text-sm text-mathua-primary break-words">{goalConcept.label}</p>
              </div>
              <button type="button" onClick={handleClearGoal} className="self-start shrink-0 font-mono text-[10px] uppercase tracking-wider text-mathua-muted hover:text-mathua-red transition-colors sm:self-auto">Clear goal</button>
            </div>
            <div className="mt-3 h-1.5 w-full bg-mathua-border" aria-hidden="true">
              <div className="h-full bg-mathua-blue" style={{ width: `${goalSummary.pct}%` }} />
            </div>
            <p className="mt-2 font-mono text-[10px] uppercase tracking-wider text-mathua-muted">
              {goalSummary.learned} of {goalSummary.total} concepts learned
            </p>
            <div className="mt-3 flex flex-wrap items-center gap-2">
              <span className="font-mono text-[10px] uppercase tracking-wider text-mathua-muted">View</span>
              <button
                type="button"
                aria-pressed={!goalPathOnly}
                onClick={() => setGoalPathOnly(false)}
                className={`font-mono text-[10px] uppercase tracking-wider px-2.5 min-h-[32px] border transition-colors ${
                  !goalPathOnly
                    ? 'bg-mathua-blue text-white border-mathua-blue'
                    : 'border-mathua-border text-mathua-muted hover:text-mathua-blue'
                }`}
              >
                Full graph
              </button>
              <button
                type="button"
                aria-pressed={goalPathOnly}
                onClick={() => setGoalPathOnly(true)}
                className={`font-mono text-[10px] uppercase tracking-wider px-2.5 min-h-[32px] border transition-colors ${
                  goalPathOnly
                    ? 'bg-mathua-blue text-white border-mathua-blue'
                    : 'border-mathua-border text-mathua-muted hover:text-mathua-blue'
                }`}
              >
                Goal path
              </button>
            </div>
            {nextStepConcept ? (
              <Link href={`/learn?concept=${encodeURIComponent(nextStepConcept.id)}`} className="mt-3 inline-flex max-w-full items-start gap-1 break-words text-left font-mono text-xs uppercase tracking-wider text-mathua-blue leading-snug hover:text-mathua-blue-hover">
                Next step: {nextStepConcept.label} &#8594;
              </Link>
            ) : (
              <p className="mt-3 text-sm text-mathua-secondary">No ready step yet. Review a prerequisite to unlock the next one.</p>
            )}
          </div>
        )}

        {focusId && focusConcept && (
          <div className="mb-6 border border-mathua-border bg-mathua-surface px-4 py-3">
            <div className="flex flex-col gap-2 sm:flex-row sm:items-start sm:justify-between sm:gap-3">
              <div className="min-w-0">
                <p className="font-mono text-[10px] uppercase tracking-wider text-mathua-muted">Focus</p>
                <p className="text-sm text-mathua-primary break-words">{focusConcept.label}</p>
              </div>
              <button type="button" onClick={handleClearFocus} className="self-start shrink-0 font-mono text-[10px] uppercase tracking-wider text-mathua-muted hover:text-mathua-primary transition-colors sm:self-auto">Exit focus</button>
            </div>
            <p className="mt-2 text-sm text-mathua-secondary">Showing this concept, its prerequisite path and its direct unlocks.</p>
          </div>
        )}

        {graphError && (
          <div role="alert" className="mb-4 flex flex-wrap items-center gap-3 border border-mathua-red bg-mathua-surface px-4 py-3">
            <span className="font-mono text-xs text-mathua-red">Couldn&apos;t load the live graph: showing the bundled fallback.</span>
            <button
              type="button"
              onClick={loadGraph}
              className="font-mono text-xs text-mathua-blue hover:text-mathua-blue-hover uppercase tracking-wider"
            >
              Retry →
            </button>
          </div>
        )}
        <ConceptGraphFlow
          concepts={concepts.length > 0 ? concepts : fallbackConcepts}
          conceptStatuses={conceptStatuses}
          conceptProgress={conceptProgress}
          theme={theme}
          onPathNodes={onPathNodes}
          selectedId={selectedId}
          onSelectionChange={handleSelectionChange}
          focusDomain={activeDomain}
          goalId={goalId}
          trimToGoal={goalPathOnly}
          focusId={focusId}
          onSetGoal={handleSetGoal}
          onClearGoal={handleClearGoal}
          onFocus={handleFocus}
          onClearFocus={handleClearFocus}
        />
      </div>

    </div>

    <Footer />
    {goalPickerOpen && (
      <GoalPicker
        concepts={goalSource}
        currentGoalId={goalId}
        onSelect={id => {
          setGoalPickerOpen(false)
          void handleSetGoal(id)
        }}
        onClose={() => setGoalPickerOpen(false)}
      />
    )}
    </AppShell>
  )
}

export default function GraphPage() {
  return (
    <Suspense fallback={<div style={{ background: 'var(--bg)', minHeight: '100vh' }} />}>
      <GraphContent />
    </Suspense>
  )
}
