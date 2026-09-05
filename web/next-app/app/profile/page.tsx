'use client'

import { useEffect, useMemo, useState } from 'react'
import Link from 'next/link'
import { useRouter } from 'next/navigation'
import { useTheme } from '../../hooks/useTheme'
import { getUserInfo, ensureGuestId, getGuestId } from '../../lib/auth'
import { getActivity, getProgress, getWeaknesses, getDueReviews, getEfficacy, getScores, getSettings } from '../../lib/api'
import type { DailyActivity, Scores, WeaknessRes, ConceptProgress, EfficacyReport } from '../../lib/api'
import ProfileStats from '../../components/ProfileStats'
import ActivityHeatmap from '../../components/ActivityHeatmap'
import DomainProgress from '../../components/DomainProgress'
import StrugglesSection from '../../components/StrugglesSection'
import Loading from '../../components/Loading'
import { AppSidebar } from '../../components/app-sidebar'
import { SidebarInset, SidebarProvider, SidebarTrigger } from '../../components/ui/sidebar'
import NextUpCard from '../../components/NextUpCard'
import { selectNextUp } from '../../lib/nextUp'

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
  const [efficacy, setEfficacy] = useState<EfficacyReport | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [avatarUrl, setAvatarUrl] = useState<string | undefined>(undefined)
  const [avatarPreset, setAvatarPreset] = useState<number | null>(null)

  const nextUp = useMemo(
    () =>
      selectNextUp({
        dueReviews,
        weaknesses,
        progress,
        activity,
        diagnosticCompleted: user?.diagnostic_completed ?? false,
        conceptsMastered: scores?.concepts_mastered ?? 0,
      }),
    [dueReviews, weaknesses, progress, activity, user?.diagnostic_completed, scores?.concepts_mastered],
  )


  useEffect(() => {
    if (!mounted) return

    const info = getUserInfo()
    setUser(info)
    // Ensure guest has an ephemeral id so profile/diagnostic works without account
    if (!info) {
      ensureGuestId()
    }

     // hydrate avatar from /api/me and settings preset
     if (info?.avatar_url) setAvatarUrl(info.avatar_url)
     getSettings().then(s => {
       if (typeof (s as unknown as { avatar_preset?: number }).avatar_preset === 'number') {
         const p = (s as unknown as { avatar_preset?: number }).avatar_preset!
         setAvatarPreset(p)
         setAvatarUrl(undefined)
       } else if (info?.avatar_url) {
         setAvatarUrl(info.avatar_url)
       }
     }).catch(() => {})

     async function fetchData() {
      try {
        if (info) {
          const [scoresRes, activityRes, progressRes, weaknessesRes] = await Promise.all([
            getScores(info!.student_id).catch(() => null as Scores | null),
            getActivity().catch(() => [] as DailyActivity[]),
            getProgress(info!.student_id).catch(() => ({} as Record<string, ConceptProgress>)),
            getWeaknesses().catch(() => ({ by_domain: {} } as WeaknessRes)),
          ])
          if (!scoresRes) {
            // 401 or fetch failure — token likely expired (JWT_SECRET rotated) or network.
            // Keep guest-like state so page doesn't hang on "Loading scores…"
            setError('Session expired or scores unavailable — please sign in again')
          }
          setScores(scoresRes as Scores | null)
          setActivity(activityRes)
          setProgress(progressRes)
          setWeaknesses(weaknessesRes)
          getDueReviews().then(r => setDueReviews(r.count)).catch(() => {})
          getEfficacy().then(setEfficacy).catch(() => {})
        } else {
          // Guest: fetch progress via ephemeral guest_id so Study answers are visible
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
      <div className="mx-auto px-4 sm:px-6 py-20 overflow-x-hidden min-w-0">
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
    )
  }

  if (error) {
    return (
      <div className="mx-auto px-4 sm:px-6 py-20 overflow-x-hidden min-w-0">
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
    )
  }

  if (!user) {
    return (
      <SidebarProvider>
        <AppSidebar
          name="Guest"
          studentId={getGuestId() || 'guest'}
        />
        <SidebarInset>
          <div className="flex h-[53px] shrink-0 items-center gap-2 border-b border-mathua-border px-4">
            <SidebarTrigger />
            <span className="font-mono text-[11px] uppercase tracking-wider text-mathua-muted">Profile</span>
          </div>
          <div id="profile-main" className="mx-auto w-full max-w-[820px] min-w-0 px-4 sm:px-6 py-8 sm:py-12 overflow-x-hidden">
          <div className="border border-mathua-border p-6 text-center bg-mathua-surface min-w-0">
            <h2 className="font-serif text-[1.2rem] text-mathua-primary mb-2">Welcome to your profile</h2>
            <p className="font-mono text-xs text-mathua-secondary mb-4">Sign in to track XP, streaks, and mastery. Your activity heatmap will appear here once you start practicing.</p>
            <div className="flex flex-col sm:flex-row flex-wrap gap-3 justify-center">
              <Link href="/login" className="w-full sm:w-auto border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white px-6 py-2 font-mono text-xs min-h-[36px] inline-flex items-center justify-center text-center whitespace-nowrap">Sign in</Link>
              <Link href="/session" className="w-full sm:w-auto border border-mathua-border text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue px-6 py-2 font-mono text-xs min-h-[36px] inline-flex items-center justify-center text-center whitespace-nowrap">Try as guest →</Link>
            </div>
          </div>
          {(Object.keys(progress).length === 0 && activity.length === 0) && (
            <section className="mt-6 border border-mathua-border bg-mathua-surface p-4">
              <h3 className="font-mono text-[11px] text-mathua-muted uppercase tracking-wider mb-3">What to do first</h3>
              <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
                <div className="border border-mathua-border p-3 bg-mathua-surface-elevated">
                  <div className="font-mono text-xs text-mathua-blue mb-1">1. Take diagnostic test</div>
                  <p className="font-mono text-[11px] text-mathua-secondary">Finds your knowledge frontier</p>
                  <Link href="/onboard" className="font-mono text-[10px] text-mathua-blue hover:text-mathua-blue-hover mt-2 inline-block">Start diagnostic test →</Link>
                </div>
                <div className="border border-mathua-border p-3 bg-mathua-surface-elevated">
                  <div className="font-mono text-xs text-mathua-blue mb-1">2. Pick a lesson in Study</div>
                  <p className="font-mono text-[11px] text-mathua-secondary">Start with Arithmetic → Fractions → Pre-Algebra</p>
                  <Link href="/study" className="font-mono text-[10px] text-mathua-blue hover:text-mathua-blue-hover mt-2 inline-block">Browse Study →</Link>
                </div>
                <div className="border border-mathua-border p-3 bg-mathua-surface-elevated">
                  <div className="font-mono text-xs text-mathua-blue mb-1">3. Practice → see XP</div>
                  <p className="font-mono text-[11px] text-mathua-secondary">2 in a row to advance · XP shows on Profile</p>
                  <Link href="/study" className="font-mono text-[10px] text-mathua-blue hover:text-mathua-blue-hover mt-2 inline-block">Start Study →</Link>
                </div>
              </div>
            </section>
          )}
          <section className="mt-8 flex min-w-0 flex-col items-stretch">
            <h2 className="font-serif text-[1.05rem] font-normal text-mathua-primary mb-4 w-full">Activity</h2>
            <div className="w-full max-w-full min-w-0 flex justify-center overflow-hidden">
              <div className="w-full max-w-full min-w-0">
                <ActivityHeatmap data={activity} />
              </div>
            </div>
          </section>
          {/* Diagnostic CTA — guest */}
          <section className="mt-6 w-full max-w-full min-w-0 overflow-hidden border border-mathua-blue bg-mathua-surface p-4 flex flex-col sm:flex-row items-stretch sm:items-center justify-between gap-3">
            <div className="w-full sm:flex-1 min-w-0 overflow-hidden">
              <div className="flex flex-col items-start gap-1.5 sm:flex-row sm:items-center sm:gap-2 min-w-0">
                <span className="shrink-0 bg-mathua-blue text-white px-2 py-0.5 font-mono text-[10px] uppercase tracking-wider">Recommended</span>
                <span className="min-w-0 break-words [overflow-wrap:anywhere] leading-snug font-mono text-[11px] sm:text-xs text-mathua-primary">Take a diagnostic to get a recommendation on where to start</span>
              </div>
              <p className="font-mono text-xs text-mathua-secondary mt-1 break-words [overflow-wrap:anywhere]">Diagnostic test · finds your knowledge frontier</p>
            </div>
            <Link href="/onboard" className="w-full sm:w-auto sm:shrink-0 border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white px-6 py-2 font-mono text-xs min-h-[36px] inline-flex items-center justify-center text-center whitespace-nowrap">Start diagnostic test →</Link>
          </section>

          <section className="mt-10">
            <div className="grid grid-cols-1 gap-4 lg:gap-6 min-w-0">
              <DomainProgress progress={progress} />
            </div>
          </section>
          </div>
        </SidebarInset>
      </SidebarProvider>
    )
  }

  if (user && !scores) {
    // Avoid infinite "Loading scores…" when token is expired or backend returns 401.
    // Show actionable error with sign-in CTA instead of hanging.
    if (error) {
      return (
        <div className="mx-auto px-4 sm:px-6 py-20 overflow-x-hidden min-w-0">
          <div style={{ fontFamily: monoFont, fontSize: 13, color: 'var(--text-muted)', textAlign: 'center' }}>
            {error}
            <div className="mt-4 flex gap-3 justify-center">
              <Link href="/login" className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white px-6 py-2 font-mono text-xs min-h-[36px] inline-flex items-center justify-center">Sign in</Link>
              <button onClick={() => window.location.reload()} className="border border-mathua-border text-mathua-secondary px-6 py-2 font-mono text-xs min-h-[36px] inline-flex items-center justify-center">Retry</button>
            </div>
          </div>
        </div>
      )
    }
    return (
      <div className="mx-auto px-4 sm:px-6 py-20 overflow-x-hidden min-w-0">
        <div style={{ fontFamily: monoFont, fontSize: 13, color: 'var(--text-muted)', textAlign: 'center' }}>
          Loading scores…
        </div>
      </div>
    )
  }

  return (
    <SidebarProvider>
      <AppSidebar
        name={user.name}
        studentId={user.student_id}
        username={user.username}
        level={scores?.level}
        streak={scores?.current_streak}
        avatarUrl={avatarUrl}
        avatarPreset={avatarPreset}
        dueReviews={dueReviews}
      />
      <SidebarInset>
        <div className="flex h-[53px] shrink-0 items-center gap-2 border-b border-mathua-border px-4">
          <SidebarTrigger />
          <span className="font-mono text-[11px] uppercase tracking-wider text-mathua-muted">Profile</span>
        </div>

        <a href="#profile-main" className="sr-only focus:not-sr-only focus:absolute focus:z-50 focus:bg-mathua-surface focus:px-4 focus:py-2 focus:font-mono focus:text-xs focus:text-mathua-blue">
          Skip to profile content
        </a>
        <div id="profile-main" className="mx-auto w-full max-w-[820px] min-w-0 px-4 sm:px-6 py-8 sm:py-12 overflow-x-hidden">
          <div className="min-w-0">
        {/* Profile stats — mobile-first */}
        <ProfileStats
          name={user.name}
          scores={scores}
          avatarSeed={user.student_id}
          avatarUrl={avatarPreset !== null ? undefined : avatarUrl}
          avatarPreset={avatarPreset}
        />

        <NextUpCard next={nextUp} />

        {(scores.concepts_mastered === 0 && !user.diagnostic_completed) && (
          <section className="mt-6 border border-mathua-border bg-mathua-surface p-4">
            <h3 className="font-mono text-[11px] text-mathua-muted uppercase tracking-wider mb-3">What to do first</h3>
            <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
              <div className="border border-mathua-border p-3 bg-mathua-surface-elevated">
                <div className="font-mono text-xs text-mathua-blue mb-1">1. Take diagnostic test</div>
                <p className="font-mono text-[11px] text-mathua-secondary">Finds your knowledge frontier</p>
                <Link href="/onboard" className="font-mono text-[10px] text-mathua-blue hover:text-mathua-blue-hover mt-2 inline-block">Start diagnostic test →</Link>
              </div>
              <div className="border border-mathua-border p-3 bg-mathua-surface-elevated">
                <div className="font-mono text-xs text-mathua-blue mb-1">2. Pick a lesson in Study</div>
                <p className="font-mono text-[11px] text-mathua-secondary">Start with Arithmetic → Fractions → Pre-Algebra</p>
                <Link href="/study" className="font-mono text-[10px] text-mathua-blue hover:text-mathua-blue-hover mt-2 inline-block">Browse Study →</Link>
              </div>
              <div className="border border-mathua-border p-3 bg-mathua-surface-elevated">
                <div className="font-mono text-xs text-mathua-blue mb-1">3. Practice → see XP</div>
                <p className="font-mono text-[11px] text-mathua-secondary">2 in a row to advance · XP shows below</p>
                <Link href="/study" className="font-mono text-[10px] text-mathua-blue hover:text-mathua-blue-hover mt-2 inline-block">Start Study →</Link>
              </div>
            </div>
          </section>
        )}

        {scores.paused_until && (
          <div className="mt-6 w-full max-w-full min-w-0 overflow-hidden border border-mathua-border bg-mathua-surface p-4 flex flex-col sm:flex-row items-stretch sm:items-center justify-between gap-3">
            <p className="font-mono text-xs text-mathua-primary min-w-0 break-words [overflow-wrap:anywhere] leading-snug">
              ⏸ Paused until {scores.paused_until} — due reviews are hidden
            </p>
            <Link href="/settings" className="w-full sm:w-auto sm:shrink-0 border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white px-4 py-2 font-mono text-xs min-h-[36px] inline-flex items-center justify-center text-center whitespace-nowrap">
              Resume
            </Link>
          </div>
        )}

        {/* Diagnostic CTA — suppressed when NextUp hero already covers diagnostic */}
        {nextUp.kind !== 'diagnostic' && (
        <section className="mt-6 w-full max-w-full min-w-0 overflow-hidden border border-mathua-blue bg-mathua-surface p-4 flex flex-col sm:flex-row items-stretch sm:items-center justify-between gap-3">
          <div className="w-full sm:flex-1 min-w-0 overflow-hidden">
            <div className="flex flex-col items-start gap-1.5 sm:flex-row sm:items-center sm:gap-2 min-w-0">
              <span className="shrink-0 bg-mathua-blue text-white px-2 py-0.5 font-mono text-[10px] uppercase tracking-wider">
                {user.diagnostic_completed ? 'Retake' : 'Recommended'}
              </span>
              <span className="min-w-0 break-words [overflow-wrap:anywhere] leading-snug font-mono text-[11px] sm:text-xs text-mathua-primary">
                {user.diagnostic_completed ? 'Retake diagnostic to refresh recommendation' : 'Take a diagnostic to get a recommendation on where to start'}
              </span>
            </div>
            <p className="font-mono text-xs text-mathua-secondary mt-1 break-words [overflow-wrap:anywhere]">
              Diagnostic test · finds your knowledge frontier
            </p>
          </div>
          <Link
            href="/onboard"
            className="w-full sm:w-auto sm:shrink-0 border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white px-6 py-2 font-mono text-xs min-h-[36px] inline-flex items-center justify-center text-center whitespace-nowrap"
          >
            {user.diagnostic_completed ? 'Retake diagnostic test →' : 'Start diagnostic test →'}
          </Link>
        </section>
        )}

        {/* 150 XP Quiz gate (CONTEXT.md Quiz) */}
        {scores && scores.xp_total >= 150 && (
          <div className="mt-6 w-full max-w-full min-w-0 overflow-hidden border border-mathua-blue bg-mathua-surface p-4 flex flex-col sm:flex-row items-stretch sm:items-center justify-between gap-3">
            <div className="w-full sm:flex-1 min-w-0 overflow-hidden">
              <div className="flex flex-col items-start gap-1.5 sm:flex-row sm:items-center sm:gap-2 min-w-0">
                <span className="shrink-0 bg-mathua-blue text-white px-2 py-0.5 font-mono text-[10px] uppercase tracking-wider">Quiz due</span>
                <span className="min-w-0 break-words [overflow-wrap:anywhere] leading-snug font-mono text-[11px] sm:text-xs text-mathua-primary">150 XP reached — mastery check recommended</span>
              </div>
              <div className="mt-2 h-1 bg-mathua-code overflow-hidden">
                <div className="h-full bg-mathua-blue" style={{ width: `${Math.min((scores.xp_total / 150) * 100, 100)}%` }} />
              </div>
            </div>
            <Link href="/goals?quiz=1" className="w-full sm:w-auto sm:shrink-0 border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white px-6 py-2 font-mono text-xs min-h-[36px] inline-flex items-center justify-center text-center whitespace-nowrap">
              Take Test →
            </Link>
          </div>
        )}

        {dueReviews > 0 && (
          <Link
            href="/session"
            className="mt-6 flex w-full min-w-0 flex-col gap-2 bg-mathua-surface border border-yellow-500/40 px-4 py-3 hover:border-yellow-500 transition-colors sm:flex-row sm:items-center sm:justify-between"
          >
            <span className="font-mono text-xs text-yellow-400 min-w-0 break-words [overflow-wrap:anywhere] leading-snug">
              ⏳ {dueReviews} concept{dueReviews !== 1 ? 's' : ''} due for review
            </span>
            <span className="font-mono text-[11px] text-yellow-400 border border-yellow-500/60 px-3 py-1.5 shrink-0 inline-flex items-center justify-center min-h-[36px] w-full sm:w-auto text-center whitespace-nowrap">
              Review Now →
            </span>
          </Link>
        )}

        {/* Activity heatmap — centered, GitHub-style, full-width on mobile */}
        <section id="activity" className="mt-8 flex min-w-0 flex-col items-stretch scroll-mt-28">
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
        <section id="domains" className="mt-8 min-w-0 scroll-mt-28">
          <div className="grid grid-cols-1 gap-4 min-w-0">
            <DomainProgress progress={progress} />
          </div>
        </section>

        <section id="struggles" className="mt-8 min-w-0 scroll-mt-28">
          <StrugglesSection weaknesses={weaknesses} />
        </section>

        {/* Efficacy — first-pass / second-pass instrumentation */}
        {efficacy && efficacy.concepts_touched > 0 && (
          <section id="efficacy" className="mt-10 min-w-0 scroll-mt-28">
            <h2 className="font-serif text-[1.05rem] font-normal text-mathua-primary mb-4">
              Efficacy
            </h2>
            <div className="grid grid-cols-2 sm:grid-cols-4 gap-3 min-w-0">
              {[
                { label: 'First-pass', value: `${Math.round(efficacy.first_pass_rate * 100)}%`, hint: 'correct on attempt 1' },
                { label: 'Second-pass', value: `${Math.round(efficacy.second_pass_rate * 100)}%`, hint: 'correct within 2 tries' },
                { label: 'Avg attempts', value: efficacy.avg_attempts_per_concept.toFixed(2), hint: 'per concept' },
                { label: 'Concepts', value: String(efficacy.concepts_touched), hint: `${efficacy.total_attempts} attempts` },
              ].map(m => (
                <div key={m.label} className="border border-mathua-border bg-mathua-surface p-3 min-w-0">
                  <div className="font-mono text-[10px] uppercase text-mathua-muted">{m.label}</div>
                  <div className="font-mono text-xl text-mathua-blue mt-1 truncate">{m.value}</div>
                  <div className="font-mono text-[10px] text-mathua-secondary mt-0.5 truncate">{m.hint}</div>
                </div>
              ))}
            </div>
          </section>
        )}


          </div>
        </div>
      </SidebarInset>
    </SidebarProvider>
  )
}
