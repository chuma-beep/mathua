'use client'

import { useEffect, useState } from 'react'
import Link from 'next/link'
import { useRouter } from 'next/navigation'
import { useTheme } from '../../hooks/useTheme'
import { getAuthHeaders, getUserInfo } from '../../lib/auth'
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
          // Guest: no scores, but still try activity/progress as guest (will be empty)
          const activityRes = await getActivity().catch(() => [])
          setActivity(activityRes as DailyActivity[])
          setScores(null)
          setProgress({})
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
        <div className="max-w-container mx-auto px-6 max-sm:px-4 py-20">
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
      </>
    )
  }

  if (error) {
    return (
      <>
        <Header />
        <div className="max-w-container mx-auto px-6 max-sm:px-4 py-20">
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
          <section className="mt-10">
            <div className="grid grid-cols-1 gap-4 lg:gap-6 min-w-0">
              <div className="border border-mathua-border p-4 bg-mathua-surface text-center min-w-0">
                <p className="font-mono text-xs text-mathua-secondary mb-3">Take a diagnostic to find your weak spots</p>
                <Link href="/onboard" className="font-mono text-xs text-mathua-blue hover:text-mathua-blue-hover">Start diagnostic →</Link>
              </div>
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
        <div className="max-w-container mx-auto px-6 max-sm:px-4 py-20">
          <div style={{ fontFamily: monoFont, fontSize: 13, color: 'var(--text-muted)', textAlign: 'center' }}>
            Loading scores…
          </div>
        </div>
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
