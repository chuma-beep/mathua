'use client'

import { useTheme } from '../../hooks/useTheme'
import Header from '../../components/Header'
import dynamic from 'next/dynamic'
import SectionHeader from '../../components/SectionHeader'
import Footer from '../../components/Footer'
import conceptsData from '../../data/concepts.json'

const MathConceptGraph3D = dynamic(
  () => import('../../components/MathConceptGraph3D'),
  {
    ssr: false,
    loading: () => (
      <div
        style={{
          height: '520px',
          background: 'var(--surface)',
          borderRadius: '8px',
          border: '0.5px solid var(--border)',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          color: 'var(--text-muted)',
          fontSize: '14px',
        }}
      >
        Loading graph...
      </div>
    ),
  }
)

const concepts = conceptsData.map((c: any) => ({
  id: c.id,
  label: c.label,
  domain: c.domain,
  prerequisites: c.prerequisites,
}))

export default function GraphPage() {
  const { theme, mounted } = useTheme()

  if (!mounted) {
    return <div style={{ background: 'var(--bg)', minHeight: '100vh' }} />
  }

  return (
    <>
      <Header />
      <div className="max-w-container mx-auto px-6 max-sm:px-4">
      <section className="pt-8">
        <span className="flex justify-end mb-4">
          <a
            href="/"
            className="text-mathua-secondary text-sm hover:text-mathua-primary"
          >
            ← Back
          </a>
        </span>
        <SectionHeader label="Your knowledge graph" title="Explore the concept map" />
        <p className="text-mathua-secondary text-sm leading-relaxed text-center max-w-[600px] mx-auto mt-4">
          Each node is a math concept. Hover to see details. After taking the diagnostic test, your
          personal progress will be overlaid on this graph — mastered concepts glow green, and your
          learning path is highlighted in gold.
        </p>
      </section>

      <div className="my-8">
        <MathConceptGraph3D concepts={concepts} theme={theme} />
      </div>

      <Footer />
    </div>
    </>
  )
}
