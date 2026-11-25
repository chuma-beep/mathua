'use client'

import Footer from '../components/Footer'
import Header from '../components/Header'
import { useHomeRedirect } from '../hooks/useHomeRedirect'
import { useScrollReveal } from '../hooks/useScrollReveal'
import { useTheme } from '../hooks/useTheme'
import {
  AdaptiveSection,
  ClosingCta,
  CoverageSection,
  FaqSection,
  FiguresSection,
  HeroSection,
  OpenSourceSection,
  PracticePreview,
  StepsSection,
} from './home/sections'

const LINKS = [
  { label: 'Learn', href: '/learn' },
  { label: 'How it works', href: '/how-it-works' },
  { label: 'Docs', href: '/docs' },
  { label: 'Note', href: '/note' },
  { label: 'Leaderboard', href: '/leaderboard' },
  { label: 'Login', href: '/login' },
]

export default function HomePage() {
  const { mounted } = useTheme()
  const homeChecked = useHomeRedirect(mounted)
  useScrollReveal(mounted && homeChecked)

  if (!mounted || !homeChecked) {
    return <div style={{ background: 'var(--bg)', minHeight: '100vh' }} />
  }

  return (
    <>
      <Header links={LINKS} />
      <main className="max-w-container mx-auto px-4 sm:px-6 overflow-x-clip min-w-0">
        <HeroSection />
        <PracticePreview />
        <StepsSection />
        <FiguresSection />
        <AdaptiveSection />
        <CoverageSection />
        <OpenSourceSection />
        <ClosingCta />
        <FaqSection />
      </main>
      <Footer />
    </>
  )
}
