'use client'

import { useState, useEffect, useRef } from 'react'
import Header from '../../components/Header'
import Footer from '../../components/Footer'
import BottomTabs from '../../components/BottomTabs'
import AsciiDivider from '../../components/AsciiDivider'
import { bodyStyle, h1Style, sections } from './styles'
import {
  ConceptGraphSection,
  DiagnosticSection,
  GeneratorsSection,
  SchedulerSection,
  ScoringSection,
  SpacedRepetitionSection,
  StudentModelSection,
  SymbolicGradingSection,
} from './sections'

function NavSidebar({ activeSection }: { activeSection: string }) {
  return (
    <nav className="w-[220px] flex-shrink-0 sticky top-[100px] h-fit max-md:fixed max-md:top-[57px] max-md:left-0 max-md:right-0 max-md:w-full max-md:z-[99] max-md:flex max-md:overflow-x-auto max-md:p-[8px_16px] max-md:gap-2"
      style={{ background: 'var(--bg)' }}
    >
      {sections.map((section) => (
        <a
          key={section.id}
          href={`#${section.id}`}
          className="block text-[12px] py-2 px-3 mb-1 transition-colors max-md:mb-0 max-md:whitespace-nowrap max-md:min-h-[36px] max-md:flex max-md:items-center max-md:shrink-0"
          style={{
            fontFamily: "'IBM Plex Mono', monospace",
            color: activeSection === section.id ? 'var(--accent-blue)' : 'var(--text-muted)',
            borderLeft: activeSection === section.id ? '2px solid var(--accent-blue)' : '2px solid transparent',
            textDecoration: 'none',
            borderRadius: 0,
            background: 'transparent',
          }}
        >
          {section.label}
        </a>
      ))}
    </nav>
  )
}

export default function HowItWorksPage() {
  const [activeSection, setActiveSection] = useState('concept-graph')
  const observerRef = useRef<IntersectionObserver | null>(null)

  useEffect(() => {
    observerRef.current = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting) {
            setActiveSection(entry.target.id)
          }
        })
      },
      { rootMargin: '-20% 0px -70% 0px' }
    )

    sections.forEach(({ id }) => {
      const el = document.getElementById(id)
      if (el) {
        observerRef.current?.observe(el)
      }
    })

    return () => observerRef.current?.disconnect()
  }, [])

  return (
    <>
      <Header links={[{ label: 'Study', href: '/study' }, { label: 'Docs', href: '/docs' }, { label: 'Note', href: '/note' }]} />
      <div className="min-h-screen w-full max-w-full min-w-0 overflow-x-hidden">
      <div className="flex max-w-[960px] mx-auto p-4 sm:p-[32px_24px] gap-4 sm:gap-10 max-md:flex-col max-md:pb-[calc(80px+env(safe-area-inset-bottom))] min-w-0">
        <NavSidebar activeSection={activeSection} />

        <main className="max-w-[720px] flex-1 min-w-0 overflow-hidden max-md:mt-20">
          <section id="intro" className="mb-12 pb-8" style={{ borderBottom: '0.5px solid var(--border)' }}>
            <h1 style={h1Style}>
              How Mathua Works
            </h1>
            <p style={bodyStyle}>
              The engine behind the learning: a technical explanation of the concept graph, student
              model, diagnostic algorithm, task selection, and scoring system.
            </p>
          </section>

          <ConceptGraphSection />
          <AsciiDivider pattern="dash" />
          <StudentModelSection />
          <AsciiDivider pattern="dash" />
          <SpacedRepetitionSection />
          <AsciiDivider pattern="dash" />
          <DiagnosticSection />
          <AsciiDivider pattern="dash" />
          <SchedulerSection />
          <AsciiDivider pattern="dash" />
          <ScoringSection />
          <AsciiDivider pattern="dash" />
          <GeneratorsSection />
          <AsciiDivider pattern="dash" />
          <SymbolicGradingSection />

          <section className="mt-12 pt-6 text-center" style={{ borderTop: '0.5px solid var(--border)' }}>
            <p style={{ ...bodyStyle, textAlign: 'center' }}>
              Ready to find your starting point?
            </p>
            <div className="flex gap-3 justify-center mt-4 max-sm:flex-col max-sm:items-center">
              <a href="/onboard" style={{ fontFamily: "'IBM Plex Mono', monospace", fontSize: '13px', color: 'var(--bg)', background: 'var(--accent-blue)', padding: '10px 22px', textDecoration: 'none' }}>
                Start diagnostic test →
              </a>
              <a href="/study" style={{ fontFamily: "'IBM Plex Mono', monospace", fontSize: '13px', color: 'var(--text-secondary)', border: '0.5px solid var(--border-strong)', padding: '10px 22px', textDecoration: 'none' }}>
                Open Study →
              </a>
            </div>
          </section>
        </main>
      </div>
    </div>
      <Footer />
      <BottomTabs />
    </>
  )
}
