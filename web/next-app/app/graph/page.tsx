'use client'

import { useTheme } from '../../hooks/useTheme'
import { useState, useEffect, useMemo, useCallback, Suspense } from 'react'
import Link from 'next/link'
import { useRouter, useSearchParams } from 'next/navigation'
import Header from '../../components/Header'
import BottomTabs from '../../components/BottomTabs'
import dynamic from 'next/dynamic'
import SectionHeader from '../../components/SectionHeader'
import ProgressSummary from '../../components/ProgressSummary'
import Footer from '../../components/Footer'
import AsciiDivider from '../../components/AsciiDivider'
import { concepts as conceptsData } from '../../lib/conceptData'
import { getScores, getGraph, getProgress, getWeaknesses, healthCheck, type GraphRes, type Scores } from '../../lib/api'
import { getUserInfo } from '../../lib/auth'
import { useAuthState } from '../../hooks/useAuthState'
import { deriveStatuses } from '../../lib/graphStatus'
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

const domainLabels = {
  arithmetic: 'Arithmetic',
  fractions: 'Fractions',
  prealgebra: 'Pre-Algebra',
  algebra: 'Algebra',
  geometry: 'Geometry',
  trigonometry: 'Trigonometry',
  complex_numbers: 'Complex Numbers',
  precalculus: 'Precalculus',
  calculus: 'Calculus',
  linear_algebra: 'Linear Algebra',
  statistics: 'Statistics',
  discrete_math: 'Discrete Math',
  number_theory: 'Number Theory',
  differential_equations: 'Diff. Eqs.',
  abstract_algebra: 'Abstract Algebra',
  topology: 'Topology',
} satisfies Record<string, string>

const domainOrder = [
  'arithmetic', 'fractions', 'prealgebra', 'algebra', 'geometry',
  'trigonometry', 'complex_numbers', 'precalculus', 'calculus', 'linear_algebra',
  'statistics', 'discrete_math', 'number_theory', 'differential_equations',
  'abstract_algebra', 'topology',
]

// StrictMode double-invokes effects in dev; share one in-flight request per key
// so mounting twice never doubles network traffic.
const inflight = new Map<string, Promise<unknown>>()
function once<T>(key: string, run: () => Promise<T>): Promise<T> {
  if (!inflight.has(key)) {
    inflight.set(
      key,
      run().catch(err => {
        inflight.delete(key)
        throw err
      })
    )
  }
  return inflight.get(key) as Promise<T>
}

function GraphContent() {
  const { theme, mounted } = useTheme()
  const { push, replace } = useRouter()
  const searchParams = useSearchParams()
  const conceptParam = searchParams.get('concept')
  const [graphData, setGraphData] = useState<GraphRes | null>(null)
  const [rawProgress, setRawProgress] = useState<Record<string, { status?: string; streak?: number }>>({})
  const [weakByDomain, setWeakByDomain] = useState<Record<string, { id: string; label: string }[]> | undefined>(undefined)
  const [connected, setConnected] = useState(false)
  const [scores, setScores] = useState<Scores | null>(null)
  const [graphError, setGraphError] = useState(false)
  const [activeDomain, setActiveDomain] = useState<string | null>(null)
  const [selectedId, setSelectedId] = useState<string | null>(conceptParam)
  const { loggedIn } = useAuthState()

  const loadGraph = useCallback(() => {
    setGraphError(false)
    once('health', healthCheck).then(setConnected).catch(() => setConnected(false))
    once('graph', getGraph).then(data => { setGraphData(data); setGraphError(false) }).catch(() => setGraphError(true))
  }, [])

  useEffect(() => {
    loadGraph()
  }, [loadGraph])

  useEffect(() => {
    if (!loggedIn) return
    const user = getUserInfo()
    if (!user) return
    const id = user.student_id
    once(`scores:${id}`, () => getScores(id)).then(setScores).catch(() => console.error('getScores failed'))
    once(`progress:${id}`, () => getProgress(id)).then(setRawProgress).catch(() => console.error('getProgress failed'))
    once('weaknesses', getWeaknesses).then(w => {
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

  // Progress toward mastery per concept, mirroring the engine's rule:
  // streak / mastery_threshold.streak, capped at 1. Mastered is full;
  // locked/unseen get no bar.
  const conceptProgress = useMemo(() => {
    const thresholds = new Map<string, number>()
    for (const c of conceptsData) {
      thresholds.set(c.id, c.mastery_threshold?.streak ?? 10)
    }
    const out: Record<string, number> = {}
    for (const [id, p] of Object.entries(rawProgress)) {
      const status = conceptStatuses[id]
      if (status === 'mastered') {
        out[id] = 1
        continue
      }
      if (status === 'unseen' || status === 'locked') continue
      out[id] = Math.min(1, (p.streak ?? 0) / (thresholds.get(id) ?? 10))
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
    return domainOrder.filter(d => set.has(d))
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

  if (!mounted) return <div style={{ background: 'var(--bg)', minHeight: '100vh' }} />

  return (
    <>
      <Header />
      <div className="max-w-container mx-auto px-4 sm:px-6 pb-[calc(80px+env(safe-area-inset-bottom))] lg:pb-0 overflow-x-hidden min-w-0">
      <section className="pt-8 min-w-0 overflow-hidden">
        <span className="flex justify-between mb-4">
          <button type="button" onClick={() => { if (window.history.length > 1) window.history.back() }} className="text-mathua-secondary text-sm hover:text-mathua-primary">
            ← Back
          </button>
        </span>
        <SectionHeader label="Your knowledge graph" title="Explore the concept map" />
        {connected && loggedIn && scores ? (
          <ProgressSummary scores={scores} weakByDomain={weakByDomain} />
        ) : connected && !loggedIn ? (
          <p className="text-mathua-secondary text-sm text-center max-w-[600px] mx-auto mt-4 mb-8">
            Sign in or{' '}
            <Link href="/session" className="text-mathua-blue hover:underline">start learning</Link>
            {' '}to track your progress across the concept map.
          </p>
        ) : (
          <p className="text-mathua-secondary text-sm text-center max-w-[600px] mx-auto mt-4 mb-8">
            <Link href="/session" className="text-mathua-blue hover:underline">Start learning</Link>
            {' '}to track your progress across the concept map.
          </p>
        )}
      </section>

      {domains.length > 0 && (
        <div className="flex flex-wrap gap-1.5 sm:gap-2 justify-center mb-6 min-w-0">
          <button
            type="button"
            onClick={() => setActiveDomain(null)}
            className={`font-mono text-[10px] uppercase px-3 min-h-[36px] py-1.5 border transition-colors min-w-0 truncate ${
              activeDomain === null
                ? 'bg-mathua-blue text-white border-mathua-blue'
                : 'border-mathua-border text-mathua-muted hover:text-mathua-primary hover:border-mathua-secondary'
            }`}
          >
            All
          </button>
          {domains.map(d => (
            <button
              type="button"
              key={d}
              onClick={() => setActiveDomain(activeDomain === d ? null : d)}
              className={`font-mono text-[10px] uppercase px-3 min-h-[36px] py-1.5 border transition-colors min-w-0 truncate ${
                activeDomain === d
                  ? 'bg-mathua-blue text-white border-mathua-blue'
                  : 'border-mathua-border text-mathua-muted hover:text-mathua-primary hover:border-mathua-secondary'
              }`}
            >
              {domainLabels[d] || d}
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
          onNodeSelect={(nodeId) => {
            const c = conceptById.get(nodeId)
            if (c) push(`/concept?id=${encodeURIComponent(c.id)}`)
          }}
        />
      </div>

      <AsciiDivider pattern="wave" />
      <Footer />
    </div>
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
