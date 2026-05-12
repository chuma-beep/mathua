'use client'

import { useTheme } from '../../hooks/useTheme'
import { useState, useEffect } from 'react'
import Header from '../../components/Header'
import dynamic from 'next/dynamic'
import SectionHeader from '../../components/SectionHeader'
import Footer from '../../components/Footer'
import { getGraph, healthCheck, type GraphRes } from '../../lib/api'
import type { MasteryStatus } from '../../components/MathConceptGraph3D'
import conceptsData from '../../data/concepts.json'

const MathConceptGraph3D = dynamic(
  () => import('../../components/MathConceptGraph3D'),
  {
    ssr: false,
    loading: () => (
      <div style={{
        height: 'clamp(320px, 50vh, 520px)',
        background: 'var(--surface)',
        borderRadius: '8px',
        border: '0.5px solid var(--border)',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        color: 'var(--text-muted)',
        fontSize: '14px',
      }}>
        Loading graph...
      </div>
    ),
  }
)

const fallbackConcepts = (conceptsData as any[]).map((c: any) => ({
  id: c.id, label: c.label, domain: c.domain, prerequisites: c.prerequisites,
}))

export default function GraphPage() {
  const { theme, mounted } = useTheme()
  const [graphData, setGraphData] = useState<GraphRes | null>(null)
  const [conceptStatuses, setConceptStatuses] = useState<Record<string, MasteryStatus>>({})
  const [connected, setConnected] = useState(false)

  useEffect(() => {
    healthCheck().then((ok) => {
      setConnected(ok)
      if (ok) {
        getGraph().then(setGraphData)
      }
    })
  }, [])

  const concepts = graphData
    ? graphData.nodes.map((n) => ({
        id: n.id,
        label: n.label,
        domain: n.domain,
        prerequisites: n.prerequisites,
      }))
    : []

  if (!mounted) return <div style={{ background: 'var(--bg)', minHeight: '100vh' }} />

  return (
    <>
      <Header />
      <div className="max-w-container mx-auto px-6 max-sm:px-4">
      <section className="pt-8">
        <span className="flex justify-between mb-4">
          <a href="/" className="text-mathua-secondary text-sm hover:text-mathua-primary">
            ← Back
          </a>
          {connected && (
            <span className="font-mono text-[11px] text-mathua-green">● Connected</span>
          )}
        </span>
        <SectionHeader label="Your knowledge graph" title="Explore the concept map" />
        <p className="text-mathua-secondary text-sm text-center max-w-[600px] mx-auto mt-4">
          {connected
            ? 'Each node is a math concept. Hover to see details. Mastered concepts glow green, and your learning path is highlighted in gold.'
            : 'Start the server (./mathua --serve) and run a diagnostic to overlay your personal progress on this graph.'}
        </p>
      </section>

      <div className="my-8">
        <MathConceptGraph3D
          concepts={concepts.length > 0 ? concepts : fallbackConcepts}
          conceptStatuses={conceptStatuses}
          theme={theme}
        />
      </div>

      <Footer />
    </div>
    </>
  )
}
