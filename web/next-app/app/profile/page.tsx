'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import { useTheme } from '../../hooks/useTheme'
import { getAuthHeaders, getUserInfo } from '../../lib/auth'
import { getActivity, getProgress, getWeaknesses } from '../../lib/api'
import type { DailyActivity, Scores, WeaknessRes, ConceptProgress } from '../../lib/api'
import Header from '../../components/Header'
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
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    if (!mounted) return

    const info = getUserInfo()
    if (!info) {
      router.push('/login')
      return
    }
    setUser(info)

    async function fetchData() {
      try {
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
      } catch (err) {
        setError('Failed to load profile data')
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

  if (error || !user || !scores) {
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
            {error || 'Unable to load profile. Are you logged in?'}
          </div>
        </div>
      </>
    )
  }

  return (
    <>
      <Header links={[
        { label: 'Study', href: '/study' },
        { label: 'Practice', href: '/session' },
        { label: 'Leaderboard', href: '/leaderboard' },
        { label: 'Graph', href: '/graph' },
        { label: 'Settings', href: '/settings' },
      ]} />

      <div className="max-w-container mx-auto px-6 max-sm:px-4 py-12">
        {/* Profile stats */}
        <ProfileStats name={user.name} scores={scores} />

        {/* Activity heatmap */}
        <section style={{ marginTop: 32 }}>
          <h2
            style={{
              fontFamily: headingFont,
              fontSize: '1.05rem',
              fontWeight: 400,
              color: 'var(--text-primary)',
              marginBottom: 16,
            }}
          >
            Activity
          </h2>
          <ActivityHeatmap data={activity} />
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
    </>
  )
}
