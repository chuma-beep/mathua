'use client'

import { useRouter } from 'next/navigation'
import { useTheme } from '../hooks/useTheme'
import { ensureGuestId, ensureGuestToken } from '../lib/auth'
import { useHomeRedirect } from '../hooks/useHomeRedirect'
import { useScrollReveal } from '../hooks/useScrollReveal'
import Header from '../components/Header'
import {
  CoverageSection,
  EditorialFooter,
  FaqSection,
  FeaturesSection,
  HeroSection,
  RankingSection,
  StatesSection,
  TrustSection,
} from './home/sections'

export default function HomePage() {
  const { theme, mounted } = useTheme()
  const { push } = useRouter()
  // Logged-in visits bounce to /profile; explicit Home clicks stay put.
  const homeChecked = useHomeRedirect(mounted)
  // Sections only exist once the redirect check resolves, so gate on it.
  useScrollReveal(mounted && homeChecked)

  if (!mounted || !homeChecked) {
    return <div style={{ background: 'var(--bg)', minHeight: '100vh' }} />
  }

  function handleGetStarted() {
    ensureGuestId()
    // Mint the guest token before landing on /profile so the first scores
    // fetch already carries a credential. Best-effort.
    ensureGuestToken().finally(() => push('/profile'))
  }

  return (
    <>
      <Header links={[{ label: 'Study', href: '/study' }, { label: 'How it works', href: '/how-it-works' }, { label: 'Docs', href: '/docs' }, { label: 'Note', href: '/note' }, { label: 'Leaderboard', href: '/leaderboard' }, { label: 'Login', href: '/login' }]} />
      <main className="max-w-container mx-auto px-4 sm:px-6 overflow-x-clip min-w-0">
        <HeroSection theme={theme} onGetStarted={handleGetStarted} />
        <TrustSection />
        <FeaturesSection />
        <StatesSection />
        <CoverageSection />
        <RankingSection />
        <FaqSection />
      </main>
      <EditorialFooter />
    </>
  )
}
