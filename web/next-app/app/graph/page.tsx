'use client'

import { useTheme } from '../../hooks/useTheme'
import { useState, useEffect, useMemo, useCallback, Suspense } from 'react'
import Link from 'next/link'
import { useRouter, useSearchParams } from 'next/navigation'
import Header from '../../components/Header'
import dynamic from 'next/dynamic'
import SectionHeader from '../../components/SectionHeader'
import ProgressSummary from '../../components/ProgressSummary'
import Footer from '../../components/Footer'
import AsciiDivider from '../../components/AsciiDivider'
import conceptsData from '../../data/concepts.json'
import { getScores, getGraph, getProgress, getWeaknesses, healthCheck, type GraphRes, type Scores } from '../../lib/api'
import { isLoggedIn, getUserInfo } from '../../lib/auth'
import { useAuthState } from '../../hooks/useAuthState'
import { deriveStatuses, type MasteryStatus } from '../../lib/graphStatus'
import Loading from '../../components/Loading'

const graphLoadingStyle: React.CSSProperties = {
  height: 'clamp(320px, 50vh, 520px)',
  background: 'var(--surface)',
  borderRadius: '8px',
  border: '0.5px solid var(--border)',
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
}

const ConceptGraphFlow = dynamic(
  () => import('../../components/ConceptGraphFlow'),
  {
    ssr: false,
    loading: () => (
      <div style={graphLoadingStyle}>
        <Loading label="LOADING GRAPH" />
      </div>
    ),
  }
)

const fallbackConcepts = (conceptsData as any[]).map((c: any) => ({
  id: c.id, label: c.label, domain: c.domain, prerequisites: c.prerequisites,
}))

const domainLabels: Record<string, string> = {
  counting: 'Counting',
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
}

const domainOrder = [
  'counting', 'arithmetic', 'fractions', 'prealgebra', 'algebra', 'geometry',
  'trigonometry', 'complex_numbers', 'precalculus', 'calculus', 'linear_algebra',
  'statistics', 'discrete_math', 'number_theory', 'differential_equations',
  'abstract_algebra', 'topology',
]

function GraphContent() {
  const { theme, mounted } = useTheme()
  const { push, replace } = useRouter()
  const searchParams = useSearchParams()
  const conceptParam = searchParams.get('concept')
  const [graphData, setGraphData] = useState<GraphRes | null>(null)
  const [rawProgress, setRawProgress] = useState<Record<string, { status?: string }>>({})
  const [weakByDomain, setWeakByDomain] = useState<Record<string, { id: string; label: string }[]> | undefined>(undefined)
  const [connected, setConnected] = useState(false)
  const [scores, setScores] = useState<Scores | null>(null)
  const [activeDomain, setActiveDomain] = useState<string | null>(null)
  const [selectedId, setSelectedId] = useState<string | null>(conceptParam)
  const { loggedIn } = useAuthState()

  useEffect(() => {
    healthCheck().then(setConnected).catch(() => setConnected(false))
    getGraph().then(setGraphData).catch(() => console.error('getGraph failed'))
  }, [])

  useEffect(() => {
    if (!loggedIn) return
    const user = getUserInfo()
    if (!user) return
    getScores(user.student_id).then(setScores).catch(() => console.error('getScores failed'))
    getProgress(user.student_id).then(setRawProgress).catch(() => console.error('getProgress failed'))
    getWeaknesses().then(w => {
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
    return concepts.filter(c => c.domain === activeDomain).map(c => c.id)
  }, [activeDomain, concepts])

  if (!mounted) return <div style={{ background: 'var(--bg)', minHeight: '100vh' }} />

  return (
    <>
      <Header />
      <div className="max-w-container mx-auto px-6 max-sm:px-4">
      <section className="pt-8">
        <span className="flex justify-between mb-4">
          <Link href="/" className="text-mathua-secondary text-sm hover:text-mathua-primary">
            ← Back
          </Link>
        </span>
        <SectionHeader label="Your knowledge graph" title="Explore the concept map" />
        {connected && loggedIn && scores ? (
          <ProgressSummary scores={scores} weakByDomain={weakByDomain} />
        ) : connected && !loggedIn ? (
          <p className="text-mathua-secondary text-sm text-center max-w-[600px] mx-auto mt-4 mb-8">
            Sign in or{' '}
            <Link href="/session" className="text-mathua-blue hover:underline">start a practice session</Link>
            {' '}to track your progress across the concept map.
          </p>
        ) : (
          <p className="text-mathua-secondary text-sm text-center max-w-[600px] mx-auto mt-4 mb-8">
            <Link href="/session" className="text-mathua-blue hover:underline">Start practicing</Link>
            {' '}to track your progress across the concept map.
          </p>
        )}
      </section>

      {domains.length > 0 && (
        <div className="flex flex-wrap gap-2 justify-center mb-6">
          <button
            onClick={() => setActiveDomain(null)}
            className={`font-mono text-[10px] uppercase px-3 h-7 border transition-colors ${
              activeDomain === null
                ? 'bg-mathua-blue text-white border-mathua-blue'
                : 'border-mathua-border text-mathua-muted hover:text-mathua-primary hover:border-mathua-secondary'
            }`}
          >
            All
          </button>
          {domains.map(d => (
            <button
              key={d}
              onClick={() => setActiveDomain(activeDomain === d ? null : d)}
              className={`font-mono text-[10px] uppercase px-3 h-7 border transition-colors ${
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

      <div className="my-8">
        <ConceptGraphFlow
          concepts={concepts.length > 0 ? concepts : fallbackConcepts}
          conceptStatuses={conceptStatuses}
          theme={theme}
          onPathNodes={onPathNodes}
          selectedId={selectedId}
          onSelectionChange={handleSelectionChange}
          onNodeSelect={(nodeId) => {
            const c = (concepts.length > 0 ? concepts : fallbackConcepts).find(n => n.id === nodeId)
            if (c) push(`/concept?id=${encodeURIComponent(c.id)}`)
          }}
        />
      </div>

      <AsciiDivider pattern="wave" />
      <Footer />
    </div>
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
