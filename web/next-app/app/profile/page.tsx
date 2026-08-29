'use client'

import { useEffect, useState } from 'react'
import Link from 'next/link'
import { useRouter } from 'next/navigation'
import { useTheme } from '../../hooks/useTheme'
import { getAuthHeaders, getUserInfo, ensureGuestId } from '../../lib/auth'
import { getActivity, getProgress, getWeaknesses, getDueReviews } from '../../lib/api'
import type { DailyActivity, Scores, WeaknessRes, ConceptProgress } from '../../lib/api'
import Header from '../../components/Header'
import BottomTabs from '../../components/BottomTabs'
import ProfileStats from '../../components/ProfileStats'
import ActivityHeatmap from '../../components/ActivityHeatmap'
import DomainProgress from '../../components/DomainProgress'
import StrugglesSection from '../../components/StrugglesSection'
import Loading from '../../components/Loading'

const API_BASE = process.env.NEXT_PUBLIC_API_URL || ''

interface UserInfo {
  student_id: string
  name: string
  username: string
  concepts_mastered: number
  current_streak: number
  level: string
  diagnostic_completed: boolean
}

export default function ProfilePage() {
  const { theme, mounted } = useTheme()
  const router = useRouter()

  const [user, setUser] = useState<UserInfo | null>(null)
  const [scores, setScores] = useState<Scores | null>(null)
  const [activity, setActivity] = useState<DailyActivity[]>([])
  const [progress, setProgress] = useState<Record<string, ConceptProgress>>({})
  const [weaknesses, setWeaknesses] = useState<WeaknessRes | null>(null)
  const [dueReviews, setDueReviews] = useState(0)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')


  useEffect(() => {
    if (!mounted) return

    const info = getUserInfo()
    setUser(info)
    // Ensure guest has an ephemeral id so profile/diagnostic works without account
    if (!info) {
      ensureGuestId()
    }

    async function fetchData() {
      try {
        if (info) {
          const headers = { ...getAuthHeaders() }
          const [scoresRes, activityRes, progressRes, weaknessesRes] = await Promise.all([
            fetch(`${API_BASE}/api/scores/`, { headers }).then((r) => r.ok ? r.json() : null),
            getActivity(),
            getProgress(info!.student_id),
            getWeaknesses(),
          ])
          setScores(scoresRes as Scores | null)
          setActivity(activityRes)
          setProgress(progressRes)
          setWeaknesses(weaknessesRes)
          getDueReviews().then(r => setDueReviews(r.count)).catch(() => {})
        } else {
          // Guest: fetch progress via ephemeral guest_id so Study answers are visible
          const { getGuestId } = await import('../../lib/auth')
          const guestId = getGuestId() || ''
          const [activityRes, progressRes] = await Promise.all([
            getActivity().catch(() => [] as DailyActivity[]),
            guestId ? getProgress(guestId).catch(() => ({} as Record<string, ConceptProgress>)) : Promise.resolve({} as Record<string, ConceptProgress>),
          ])
          setActivity(activityRes as DailyActivity[])
          setScores(null)
          setProgress(progressRes as Record<string, ConceptProgress>)
          setWeaknesses(null)
        }
      } catch (err) {
        if (info) setError('Failed to load profile data')
      } finally {
        setLoading(false)
      }
    }

    fetchData()
  }, [mounted, router])

  const headingFont = "'IBM Plex Serif', serif"
  const monoFont = "'IBM Plex Mono', monospace"

  if (!mounted) {
    return <div style={{ background: 'var(--bg)', minHeight: '100vh' }} />
  }

  if (loading) {
    return (
      <>
        <Header />
        <div className="max-w-container mx-auto px-4 sm:px-6 py-20 pb-[calc(80px+env(safe-area-inset-bottom))] lg:pb-0 overflow-x-hidden min-w-0">
          <div
            style={{
              fontFamily: monoFont,
              fontSize: 13,
              color: 'var(--text-muted)',
              textAlign: 'center',
            }}
          >
            <Loading label="LOADING PROFILE" />
          </div>
        </div>
        <BottomTabs />
      </>
    )
  }

  if (error) {
    return (
      <>
        <Header />
        <div className="max-w-container mx-auto px-4 sm:px-6 py-20 pb-[calc(80px+env(safe-area-inset-bottom))] lg:pb-0 overflow-x-hidden min-w-0">
          <div
            style={{
              fontFamily: monoFont,
              fontSize: 13,
              color: 'var(--text-muted)',
              textAlign: 'center',
            }}
          >
            {error}
          </div>
        </div>
        <BottomTabs />
      </>
    )
  }

  if (!user) {
    return (
      <>
        <Header
          links={[
            { label: 'Study', href: '/study' },
            { label: 'Practice', href: '/session' },
            { label: 'Leaderboard', href: '/leaderboard' },
            { label: 'Graph', href: '/graph' },
          ]}
        />
        <div className="mx-auto w-full max-w-[820px] min-w-0 px-4 sm:px-6 py-8 sm:py-12 pb-[calc(80px+env(safe-area-inset-bottom))] lg:pb-12 overflow-x-hidden">
          <div className="border border-mathua-border p-6 text-center bg-mathua-surface min-w-0">
            <h2 className="font-serif text-[1.2rem] text-mathua-primary mb-2">Welcome to your profile</h2>
            <p className="font-mono text-xs text-mathua-secondary mb-4">Sign in to track XP, streaks, and mastery. Your activity heatmap will appear here once you start practicing.</p>
            <div className="flex flex-wrap gap-3 justify-center">
              <Link href="/login" className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white px-6 py-2 font-mono text-xs min-h-[36px] inline-flex items-center justify-center">Sign in</Link>
              <Link href="/session" className="border border-mathua-border text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue px-6 py-2 font-mono text-xs min-h-[36px] inline-flex items-center justify-center">Try as guest →</Link>
            </div>
          </div>
          <section className="mt-8 flex min-w-0 flex-col items-stretch">
            <h2 className="font-serif text-[1.05rem] font-normal text-mathua-primary mb-4 w-full">Activity</h2>
            <div className="w-full max-w-full min-w-0 flex justify-center overflow-hidden">
              <div className="w-full max-w-full min-w-0">
                <ActivityHeatmap data={activity} />
              </div>
            </div>
          </section>
          {/* Diagnostic CTA — guest */}
          <section className="mt-6 border border-mathua-blue bg-mathua-surface p-4 flex flex-col sm:flex-row items-center justify-between gap-3">
            <div className="min-w-0">
              <div className="flex items-center gap-2">
                <span className="bg-mathua-blue text-white px-2 py-0.5 font-mono text-[10px] uppercase tracking-wider">Recommended</span>
                <span className="font-mono text-xs text-mathua-primary truncate">Take a diagnostic to get a recommendation on where to start</span>
              </div>
              <p className="font-mono text-xs text-mathua-secondary mt-1">20–35 adaptive questions · finds your knowledge frontier</p>
            </div>
            <Link href="/onboard" className="shrink-0 border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white px-6 py-2 font-mono text-xs min-h-[36px] inline-flex items-center justify-center">Start diagnostic →</Link>
          </section>

          <section className="mt-10">
            <div className="grid grid-cols-1 gap-4 lg:gap-6 min-w-0">
              <DomainProgress progress={progress} />
            </div>
          </section>
        </div>
        <BottomTabs />
      </>
    )
  }

  if (!scores) {
    return (
      <>
        <Header />
        <div className="max-w-container mx-auto px-4 sm:px-6 py-20 pb-[calc(80px+env(safe-area-inset-bottom))] lg:pb-0 overflow-x-hidden min-w-0">
          <div style={{ fontFamily: monoFont, fontSize: 13, color: 'var(--text-muted)', textAlign: 'center' }}>
            Loading scores…
          </div>
        </div>
        <BottomTabs />
      </>
    )
  }

  return (
    <>
      <Header
        links={[
          { label: 'Study', href: '/study' },
          { label: 'Practice', href: '/session' },
          { label: 'Leaderboard', href: '/leaderboard' },
          { label: 'Graph', href: '/graph' },
          { label: 'Settings', href: '/settings' },
        ]}
      />

      <div className="mx-auto w-full max-w-[820px] min-w-0 px-4 sm:px-6 py-8 sm:py-12 pb-[calc(80px+env(safe-area-inset-bottom))] lg:pb-12 overflow-x-hidden">
        {/* Profile stats — mobile-first */}
        <ProfileStats name={user.name} scores={scores} />

        {scores.paused_until && (
          <div className="mt-6 border border-mathua-border bg-mathua-surface p-4 flex flex-col sm:flex-row items-center justify-between gap-3">
            <p className="font-mono text-xs text-mathua-primary min-w-0 truncate">
              ⏸ Paused until {scores.paused_until} — due reviews are hidden
            </p>
            <Link href="/settings" className="shrink-0 border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white px-4 py-2 font-mono text-xs min-h-[36px] inline-flex items-center justify-center">
              Resume
            </Link>
          </div>
        )}

        {/* Diagnostic CTA — both authed and guest via profile */}
        <section className="mt-6 border border-mathua-blue bg-mathua-surface p-4 flex flex-col sm:flex-row items-center justify-between gap-3">
          <div className="min-w-0">
            <div className="flex items-center gap-2">
              <span className="bg-mathua-blue text-white px-2 py-0.5 font-mono text-[10px] uppercase tracking-wider">
                {user.diagnostic_completed ? 'Retake' : 'Recommended'}
              </span>
              <span className="font-mono text-xs text-mathua-primary truncate">
                {user.diagnostic_completed ? 'Retake diagnostic to refresh recommendation' : 'Take a diagnostic to get a recommendation on where to start'}
              </span>
            </div>
            <p className="font-mono text-xs text-mathua-secondary mt-1">
              20–35 adaptive questions · finds your knowledge frontier
            </p>
          </div>
          <Link
            href="/onboard"
            className="shrink-0 border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white px-6 py-2 font-mono text-xs min-h-[36px] inline-flex items-center justify-center"
          >
            {user.diagnostic_completed ? 'Retake diagnostic →' : 'Start diagnostic →'}
          </Link>
        </section>

        {/* 150 XP Quiz gate (CONTEXT.md Quiz) */}
        {scores && scores.xp_total >= 150 && (
          <div className="mt-6 border border-mathua-blue bg-mathua-surface p-4 flex flex-col sm:flex-row items-center justify-between gap-3">
            <div className="min-w-0">
              <div className="flex items-center gap-2">
                <span className="bg-mathua-blue text-white px-2 py-0.5 font-mono text-[10px] uppercase tracking-wider">Quiz due</span>
                <span className="font-mono text-xs text-mathua-primary truncate">150 XP reached — mastery check recommended</span>
              </div>
              <div className="mt-2 h-1 bg-mathua-code overflow-hidden">
                <div className="h-full bg-mathua-blue" style={{ width: `${Math.min((scores.xp_total / 150) * 100, 100)}%` }} />
              </div>
            </div>
            <Link href="/goals?quiz=1" className="shrink-0 border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white px-6 py-2 font-mono text-xs min-h-[36px] inline-flex items-center justify-center">
              Take Test →
            </Link>
          </div>
        )}

        {dueReviews > 0 && (
          <Link
            href="/session"
            className="mt-6 flex w-full min-w-0 flex-col gap-2 bg-mathua-surface border border-yellow-500/40 px-4 py-3 hover:border-yellow-500 transition-colors sm:flex-row sm:items-center sm:justify-between"
          >
            <span className="font-mono text-xs text-yellow-400 min-w-0 truncate">
              ⏳ {dueReviews} concept{dueReviews !== 1 ? 's' : ''} due for review
            </span>
            <span className="font-mono text-[11px] text-yellow-400 border border-yellow-500/60 px-3 py-1.5 shrink-0 inline-flex items-center justify-center min-h-[36px] w-full sm:w-auto">
              Review Now →
            </span>
          </Link>
        )}

        {/* Activity heatmap — centered, GitHub-style, full-width on mobile */}
        <section className="mt-8 flex min-w-0 flex-col items-stretch">
          <h2 className="font-serif text-[1.05rem] font-normal text-mathua-primary mb-4 w-full">
            Activity
          </h2>
          <div className="w-full max-w-full min-w-0 flex justify-center overflow-hidden">
            <div className="w-full max-w-full min-w-0">
              <ActivityHeatmap data={activity} />
            </div>
          </div>
        </section>

        {/* Domain progress + Struggles — mobile-first: CTA on top, stacked */}
        <section className="mt-8 min-w-0">
          <div className="grid grid-cols-1 gap-4 lg:grid-cols-[3fr_2fr] lg:gap-6 min-w-0">
            <DomainProgress progress={progress} />
            <StrugglesSection weaknesses={weaknesses} />
          </div>
        </section>
      </div>
      <BottomTabs />
    </>
  )
}
