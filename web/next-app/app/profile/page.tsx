'use client'

import { useEffect, useState } from 'react'
import Link from 'next/link'
import { useRouter } from 'next/navigation'
import { useTheme } from '../../hooks/useTheme'
import { getAuthHeaders, getUserInfo } from '../../lib/auth'
import { getActivity, getProgress, getWeaknesses, getDueReviews } from '../../lib/api'
import type { DailyActivity, Scores, WeaknessRes, ConceptProgress } from '../../lib/api'
import Header from '../../components/Header'
import Drawer from '../../components/Drawer'
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
  const [drawerOpen, setDrawerOpen] = useState(true)
  const [mobileDrawer, setMobileDrawer] = useState(false)

  useEffect(() => {
    const stored = localStorage.getItem('mathua-drawer-open')
    if (stored === 'false') setDrawerOpen(false)
  }, [])

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
          drawerOpen={drawerOpen}
          onToggleDrawer={() => {
            const next = !drawerOpen
            setDrawerOpen(next)
            localStorage.setItem('mathua-drawer-open', String(next))
          }}
          onOpenMobileDrawer={() => setMobileDrawer(true)}
          hideNavWhenOpen
          links={[
            { label: 'Study', href: '/study' },
            { label: 'Practice', href: '/session' },
            { label: 'Leaderboard', href: '/leaderboard' },
            { label: 'Graph', href: '/graph' },
          ]}
        />
        <div className="flex max-w-container mx-auto">
          <aside className={`${drawerOpen ? 'w-[240px]' : 'w-[48px]'} hidden lg:block shrink-0 sticky top-[52px] h-[calc(100vh-52px)] border-r-[0.5px] border-mathua-border bg-mathua-bg overflow-y-auto transition-all`}>
            <Drawer collapsed={!drawerOpen} onToggle={() => {
              const next = !drawerOpen
              setDrawerOpen(next)
              localStorage.setItem('mathua-drawer-open', String(next))
            }} />
          </aside>
          <div className="flex-1 min-w-0 px-6 max-sm:px-4 py-12 pb-[56px] lg:pb-12">
          <div className="border border-mathua-border p-6 text-center bg-mathua-surface">
            <h2 style={{ fontFamily: headingFont, fontSize: '1.2rem', color: 'var(--text-primary)', marginBottom: 8 }}>Welcome to your profile</h2>
            <p style={{ fontFamily: monoFont, fontSize: 12, color: 'var(--text-secondary)', marginBottom: 16 }}>Sign in to track XP, streaks, and mastery. Your activity heatmap will appear here once you start practicing.</p>
            <div className="flex gap-3 justify-center">
              <Link href="/login" className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white px-6 py-2 font-mono text-xs">Sign in</Link>
              <Link href="/session" className="border border-mathua-border text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue px-6 py-2 font-mono text-xs">Try as guest →</Link>
            </div>
          </div>
          <section style={{ marginTop: 32, display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
            <h2 style={{ fontFamily: headingFont, fontSize: '1.05rem', fontWeight: 400, color: 'var(--text-primary)', marginBottom: 16, width: '100%', maxWidth: 820 }}>Activity</h2>
            <div style={{ width: '100%', maxWidth: 820, display: 'flex', justifyContent: 'center' }}>
              <ActivityHeatmap data={activity} />
            </div>
          </section>
          <section style={{ marginTop: 40 }}>
            <div style={{ display: 'grid', gridTemplateColumns: 'minmax(0, 3fr) minmax(0, 2fr)', gap: 24, alignItems: 'start' }} className="max-sm:block max-sm:[&>*+*]:mt-6">
              <DomainProgress progress={progress} />
              <div className="border border-mathua-border p-4 bg-mathua-surface text-center">
                <p className="font-mono text-xs text-mathua-secondary mb-3">Take a diagnostic to find your weak spots</p>
                <Link href="/onboard" className="font-mono text-xs text-mathua-blue hover:text-mathua-blue-hover">Start diagnostic →</Link>
              </div>
            </div>
          </section>
          </div>
        </div>
        <BottomTabs onMore={() => setMobileDrawer(true)} />
        {mobileDrawer && (
          <>
            <div className="lg:hidden fixed inset-0 z-50 bg-black/40" onClick={() => setMobileDrawer(false)} />
            <div className="lg:hidden fixed left-0 top-0 bottom-0 w-[260px] z-50 bg-mathua-bg border-r border-mathua-border overflow-y-auto">
              <Drawer collapsed={false} onToggle={() => setMobileDrawer(false)} />
            </div>
          </>
        )}
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
        drawerOpen={drawerOpen}
        onToggleDrawer={() => {
          const next = !drawerOpen
          setDrawerOpen(next)
          localStorage.setItem('mathua-drawer-open', String(next))
        }}
        onOpenMobileDrawer={() => setMobileDrawer(true)}
        hideNavWhenOpen
        links={[
          { label: 'Study', href: '/study' },
          { label: 'Practice', href: '/session' },
          { label: 'Leaderboard', href: '/leaderboard' },
          { label: 'Graph', href: '/graph' },
          { label: 'Settings', href: '/settings' },
        ]}
      />

      <div className="flex max-w-container mx-auto">
        <aside className={`${drawerOpen ? 'w-[240px]' : 'w-[48px]'} hidden lg:block shrink-0 sticky top-[52px] h-[calc(100vh-52px)] border-r-[0.5px] border-mathua-border bg-mathua-bg overflow-y-auto transition-all`}>
          <Drawer collapsed={!drawerOpen} onToggle={() => {
            const next = !drawerOpen
            setDrawerOpen(next)
            localStorage.setItem('mathua-drawer-open', String(next))
          }} />
        </aside>
        <div className="flex-1 min-w-0 px-6 max-sm:px-4 py-12 pb-[56px] lg:pb-12">
        {/* Profile stats */}
        <ProfileStats name={user.name} scores={scores} />

        {dueReviews > 0 && (
          <Link
            href="/session"
            className="block mt-6 bg-mathua-surface border border-yellow-500/40 rounded-none px-4 py-3 flex items-center justify-between hover:border-yellow-500 transition-colors"
          >
            <span className="font-mono text-xs text-yellow-400">
              ⏳ {dueReviews} concept{dueReviews !== 1 ? 's' : ''} due for review
            </span>
            <span className="font-mono text-[11px] text-yellow-400 border border-yellow-500/60 rounded-none px-3 py-1.5">
              Review Now →
            </span>
          </Link>
        )}

        {/* Activity heatmap — centered, GitHub-style */}
        <section style={{ marginTop: 32, display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
          <h2
            style={{
              fontFamily: headingFont,
              fontSize: '1.05rem',
              fontWeight: 400,
              color: 'var(--text-primary)',
              marginBottom: 16,
              width: '100%',
              maxWidth: 820,
            }}
          >
            Activity
          </h2>
          <div style={{ width: '100%', maxWidth: 820, display: 'flex', justifyContent: 'center' }}>
            <ActivityHeatmap data={activity} />
          </div>
        </section>

        {/* Domain progress + Struggles */}
        <section style={{ marginTop: 40 }}>
          <div
            style={{
              display: 'grid',
              gridTemplateColumns: 'minmax(0, 3fr) minmax(0, 2fr)',
              gap: 24,
              alignItems: 'start',
            }}
            className="max-sm:block max-sm:[&>*+*]:mt-6"
          >
            <DomainProgress progress={progress} />
            <StrugglesSection weaknesses={weaknesses} />
          </div>
        </section>
        </div>
      </div>
      <BottomTabs onMore={() => setMobileDrawer(true)} />
      {mobileDrawer && (
        <>
          <div className="lg:hidden fixed inset-0 z-50 bg-black/40" onClick={() => setMobileDrawer(false)} />
          <div className="lg:hidden fixed left-0 top-0 bottom-0 w-[260px] z-50 bg-mathua-bg border-r border-mathua-border overflow-y-auto">
            <Drawer collapsed={false} onToggle={() => setMobileDrawer(false)} />
          </div>
        </>
      )}
    </>
  )
}
