'use client'

import { useRouter } from 'next/navigation'
import { useTheme } from '../hooks/useTheme'
import { ensureGuestId, ensureGuestToken } from '../lib/auth'
import { useHomeRedirect } from '../hooks/useHomeRedirect'
import Header from '../components/Header'
import AsciiDivider from '../components/AsciiDivider'
import Footer from '../components/Footer'
import {
  ContributingSection,
  CurriculumSection,
  FeaturesSection,
  HeroSection,
  PipelineSection,
  ProgressionSection,
} from './home/sections'

export default function HomePage() {
  const { theme, mounted } = useTheme()
  const { push } = useRouter()
  // Logged-in visits bounce to /profile; explicit Home clicks stay put.
  const homeChecked = useHomeRedirect(mounted)

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
      <div className="max-w-container mx-auto px-4 sm:px-6 pb-[calc(80px+env(safe-area-inset-bottom))] lg:pb-0 overflow-x-hidden min-w-0">
        <HeroSection theme={theme} onGetStarted={handleGetStarted} />
        <AsciiDivider pattern="wave" />
        <FeaturesSection />
        <PipelineSection />
        <AsciiDivider pattern="dash" />
        <CurriculumSection />
        <AsciiDivider pattern="wave" />
        <ProgressionSection />
        <AsciiDivider pattern="wave" />
        <ContributingSection />
        <Footer />
      </div>
    </>
  )
}
