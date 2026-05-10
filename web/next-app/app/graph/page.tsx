'use client'

import dynamic from 'next/dynamic'
import conceptsData from '../../data/concepts.json'

const MathConceptGraph3D = dynamic(() => import('../../components/MathConceptGraph3D'), {
  ssr: false,
  loading: () => (
    <div style={{
      height: '520px',
      background: 'var(--bg-card, #1a1a1a)',
      borderRadius: '8px',
      border: '0.5px solid var(--border, #333)',
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'center',
      color: 'var(--text-muted, #888)',
      fontSize: '14px',
    }}>Loading graph...</div>
  ),
})

const concepts = conceptsData.map((c: any) => ({
  id: c.id,
  label: c.label,
  domain: c.domain,
  prerequisites: c.prerequisites,
}))

export default function GraphPage() {
  return (
    <div className="container" style={{ paddingTop: '2rem' }}>
      <section>
        <span className="section-label">Your knowledge graph</span>
        <h2 style={{ textAlign: 'center' }}>Explore the concept map</h2>
        <p className="body-text" style={{ textAlign: 'center', maxWidth: '600px', margin: '0 auto 2rem' }}>
          Each node is a maths concept. Hover to see details. After taking the diagnostic test,
          your personal progress will be overlaid on this graph — mastered concepts glow green,
          and your learning path is highlighted in gold.
        </p>
      </section>
      <MathConceptGraph3D
        concepts={concepts}
      />
    </div>
  )
}
