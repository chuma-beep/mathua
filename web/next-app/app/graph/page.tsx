'use client'

import { useTheme } from '../../hooks/useTheme'
import { useState, useEffect } from 'react'
import Link from 'next/link'
import Header from '../../components/Header'
import dynamic from 'next/dynamic'
import SectionHeader from '../../components/SectionHeader'
import ProgressSummary from '../../components/ProgressSummary'
import Footer from '../../components/Footer'
import conceptsData from '../../data/concepts.json'
import { getScores, getGraph, healthCheck, type GraphRes, type Scores } from '../../lib/api'
import { isLoggedIn, getUserInfo } from '../../lib/auth'
import type { MasteryStatus } from '../../components/MathConceptGraph3D'

const graphLoadingStyle: React.CSSProperties = {
  height: 'clamp(320px, 50vh, 520px)',
  background: 'var(--surface)',
  borderRadius: '8px',
  border: '0.5px solid var(--border)',
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  color: 'var(--text-muted)',
  fontSize: '14px',
}

const MathConceptGraph3D = dynamic(
  () => import('../../components/MathConceptGraph3D'),
  {
    ssr: false,
    loading: () => (
      <div style={graphLoadingStyle}>
        Loading graph&hellip;
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
  const [conceptStatuses] = useState<Record<string, MasteryStatus>>({})
  const [connected, setConnected] = useState(false)
  const [loggedIn, setLoggedIn] = useState(false)
  const [scores, setScores] = useState<Scores | null>(null)

  useEffect(() => {
    const loggedInVal = isLoggedIn()
    setLoggedIn(loggedInVal)
    healthCheck().then((ok) => {
      setConnected(ok)
      if (ok) {
        getGraph().then(setGraphData)
      }
      if (ok && loggedInVal) {
        const user = getUserInfo()
        if (user) {
          getScores(user.student_id).then(setScores).catch(() => {})
        }
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
          <Link href="/" className="text-mathua-secondary text-sm hover:text-mathua-primary">
            ← Back
          </Link>
        </span>
        <SectionHeader label="Your knowledge graph" title="Explore the concept map" />
        {connected && loggedIn && scores ? (
          <ProgressSummary scores={scores} />
        ) : connected && !loggedIn ? (
          <p className="text-mathua-secondary text-sm text-center max-w-[600px] mx-auto mt-4 mb-8">
            Sign in or{' '}
            <Link href="/session" className="text-mathua-gold hover:underline">start a practice session</Link>
            {' '}to track your progress across the concept map.
          </p>
        ) : (
          <p className="text-mathua-secondary text-sm text-center max-w-[600px] mx-auto mt-4 mb-8">
            <Link href="/session" className="text-mathua-gold hover:underline">Start practicing</Link>
            {' '}to track your progress across the concept map.
          </p>
        )}
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
