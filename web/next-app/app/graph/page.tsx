'use client'

import { useTheme } from '../../hooks/useTheme'
import { useState, useEffect, useMemo, useCallback, useRef, Suspense } from 'react'
import Link from 'next/link'
import { useRouter, useSearchParams } from 'next/navigation'
import Header from '../../components/Header'
import BottomTabs from '../../components/BottomTabs'
import dynamic from 'next/dynamic'
import SectionHeader from '../../components/SectionHeader'
import ProgressSummary from '../../components/ProgressSummary'
import Footer from '../../components/Footer'
import { concepts as conceptsData } from '../../lib/conceptData'
import { getScores, getGraph, getProgress, getWeaknesses, healthCheck, type GraphRes, type Scores } from '../../lib/api'
import { getUserInfo } from '../../lib/auth'
import { useAuthState } from '../../hooks/useAuthState'
import { deriveStatuses } from '../../lib/graphStatus'
import { DOMAIN_ORDER, domainLabel } from '../../lib/graphDomains'
import { createInflightCache, createSessionCache } from '../../lib/requestCache'
import Loading from '../../components/Loading'

const graphLoadingStyle: React.CSSProperties = {
  height: 'clamp(320px, 60dvh, 520px)',
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

function GraphContent() {
  const { theme, mounted } = useTheme()
  const { push, replace } = useRouter()
  const searchParams = useSearchParams()
  const conceptParam = searchParams.get('concept')
  // ?domain= makes the filter linkable, so /domains can deep-link into one subject.
  const domainParam = searchParams.get('domain')
  const [graphData, setGraphData] = useState<GraphRes | null>(null)
  const [rawProgress, setRawProgress] = useState<Record<string, { status?: string; streak?: number; mastery_pct?: number }>>({})
  const [weakByDomain, setWeakByDomain] = useState<Record<string, { id: string; label: string }[]> | undefined>(undefined)
  const [connected, setConnected] = useState(false)
  const [scores, setScores] = useState<Scores | null>(null)
  const [graphError, setGraphError] = useState(false)
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
  // from the tier and the streak together — see `mastery.MasteryPct`.
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
    if (id) {
      replace(`/graph?concept=${encodeURIComponent(id)}`, { scroll: false })
    } else {
      replace('/graph', { scroll: false })
    }
  }, [replace])

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

  if (!mounted) return <div style={{ background: 'var(--bg)', minHeight: '100vh' }} />

  return (
    <>
      <Header />
      <div className="max-w-container mx-auto px-4 sm:px-6 pb-[calc(80px+env(safe-area-inset-bottom))] lg:pb-0 overflow-x-hidden min-w-0">
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

      {domains.length > 0 && (
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
          onNodeSelect={(nodeId) => {
            const c = conceptById.get(nodeId)
            if (c) push(`/study?concept=${encodeURIComponent(c.id)}`)
          }}
        />
      </div>

    </div>

    <Footer />
    <BottomTabs />
    </>
  )
}

export default function GraphPage() {
  return (
    <Suspense fallback={<div style={{ background: 'var(--bg)', minHeight: '100vh' }} />}>
      <GraphContent />
    </Suspense>
  )
}
