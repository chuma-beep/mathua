'use client'

import { useEffect, useMemo, useState } from 'react'
import Link from 'next/link'
import { useTheme } from '../../hooks/useTheme'
import { getUserInfo, ensureGuestId, ensureGuestToken, getGuestId, isLoggedIn } from '../../lib/auth'
import { ensureDicebearAvatar, resolveAvatar } from '../../lib/dicebear'
import { getActivity, getProgress, getWeaknesses, getDueReviews, getEfficacy, getScores, getSettings } from '../../lib/api'
import type { DailyActivity, Scores, WeaknessRes, ConceptProgress, EfficacyReport } from '../../lib/api'
import ProfileStats from '../../components/ProfileStats'
import ActivityHeatmap from '../../components/ActivityHeatmap'
import DomainProgress from '../../components/DomainProgress'
import StrugglesSection from '../../components/StrugglesSection'
import ProfileSkeleton from '../../components/skeletons/ProfileSkeleton'
import { AppSidebar } from '../../components/app-sidebar'
import { SidebarInset, SidebarProvider, SidebarTrigger } from '../../components/ui/sidebar'
import NextUpSummary from '../../components/NextUpSummary'
import PositionBlock from '../../components/PositionBlock'
import DailyGoalControl, { getGuestGoal } from '../../components/DailyGoalControl'
import { isNewUser, recentlyUnlocked, hrefConceptId, RECENT_UNLOCK_DAYS, type Shelf } from '../../lib/nextUp'
import { fetchShelf } from '../../lib/recommendations'
import { concepts as conceptCatalog } from '../../lib/conceptData'

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
  const { mounted } = useTheme()

  const [user, setUser] = useState<UserInfo | null>(null)
  const [scores, setScores] = useState<Scores | null>(null)
  const [activity, setActivity] = useState<DailyActivity[]>([])
  const [progress, setProgress] = useState<Record<string, ConceptProgress>>({})
  const [weaknesses, setWeaknesses] = useState<WeaknessRes | null>(null)
  const [efficacy, setEfficacy] = useState<EfficacyReport | null>(null)
  const [dueReviews, setDueReviews] = useState(0)
  const [shelf, setShelf] = useState<Shelf | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [goalOverride, setGoalOverride] = useState<number | null>(null)
  const [avatarUrl, setAvatarUrl] = useState<string | undefined>(undefined)
  const [avatarPreset, setAvatarPreset] = useState<number | null>(null)

  // Guest daily-goal override (server default applies otherwise).
  useEffect(() => {
    if (mounted && !getUserInfo()) {
      setGoalOverride(getGuestGoal())
    }
  }, [mounted])

  // Null when an older server omits it. The quiz block hides rather than substituting a
  // literal, because a bar computed against a gate the engine does not use is worse than no bar.
  const quizGate = scores?.quiz_gate_xp && scores.quiz_gate_xp > 0 ? scores.quiz_gate_xp : null

  const effScores = useMemo(
    () => (scores && goalOverride ? { ...scores, daily_xp_goal: goalOverride } : scores),
    [scores, goalOverride],
  )

  // New = nothing mastered and no answered questions anywhere.
  const isNew = scores != null && isNewUser({ conceptsMastered: scores.concepts_mastered, activity })

  // The head, as the engine ranked it.
  //
  // This used to be re-derived here from `progress`, `weaknesses`, `activity`, `dueReviews`
  // and `scores` — five inputs the browser combined with its own copy of the eligibility and
  // ranking rules, which had already drifted from internal/scheduler: a duplicated XP award
  // that still halved reviews, and a different ranking order. Now it arrives decided.
  //
  // A null shelf is the pre-fetch state rather than an empty one: rendering "nothing
  // recommended" while the request is in flight would misreport the learner.
  const head = shelf

  // What the head names, for "where am I". Read off the same head as the recommendation, so
  // the two cannot disagree — which is the property that matters, and the reason this block
  // keeps the label even though it no longer offers a button: the engine owns both the name
  // and the destination, so a learner is never shown one concept's name and sent to another.
  const frontierLabel = head?.next.title ?? null

  // Recently unlocked: successors unlocked by recently-active concepts
  // (recency derived client-side from activity — no endpoint needed). Head
  // and queue destinations are filtered so no two items share an href and
  // nothing re-links the head concept.
  const unlockRows = useMemo(() => {
    // head is null until the recommendation arrives. Deriving rows from it unconditionally
    // would throw on the first render — the memo ran before with a synchronously computed head,
    // which is exactly what moving the decision to the engine changed.
    if (!head) return []
    const rows = recentlyUnlocked({
      catalog: conceptCatalog.map(c => ({ id: c.id, label: c.label, prerequisites: c.prerequisites ?? [] })),
      progress,
      activity,
      recentDays: RECENT_UNLOCK_DAYS,
    })
    const headCid = hrefConceptId(head.next.href)
    const seen = new Set([head.next.href, ...head.alternatives.map(a => a.href)])
    return rows.filter(r => {
      const href = `/learn?concept=${encodeURIComponent(r.id)}`
      if (seen.has(href)) return false
      const cid = hrefConceptId(href)
      if (cid !== null && cid === headCid) return false
      seen.add(href)
      return true
    }).slice(0, 5)
  }, [head, progress, activity])


  useEffect(() => {
    if (!mounted) return

    const info = getUserInfo()
    setUser(info)
    // Ensure guest has an ephemeral id so profile/diagnostic works without account
    if (!info) {
      ensureGuestId()
      // Best-effort: mint a guest bearer token so the guest reads their own
      // progress/scores as themselves (no credential = 401 for UUID rows).
      ensureGuestToken()
    }

      // Avatar precedence lives in resolveAvatar (shared with Header):
      // custom upload > Google photo > DiceBear pick > legacy preset >
      // automatic initial.
      const hydrateAvatar = (s: {
        avatar_custom?: boolean
        avatar_dicebear?: { style: string; seed: string } | null
        avatar_preset?: number | null
      }) => {
        const resolved = resolveAvatar(info, s)
        setAvatarPreset(resolved.preset ?? null)
        setAvatarUrl(resolved.url)
      }
      getSettings().then(hydrateAvatar).catch(() => {
        if (info?.avatar_url) setAvatarUrl(info.avatar_url)
      })
      // Surprise DiceBear combo for photo-less users (all auth flows land here).
      ensureDicebearAvatar().then(pick => {
        if (pick) hydrateAvatar({ avatar_dicebear: { style: pick.style, seed: pick.seed } })
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
          // The recommendation is the engine's answer, not a re-ranking of the fetches above.
          // Fetched beside them rather than derived from them: the report's numbers and the
          // recommendation come from one authority each, and neither is computed from the other.
          fetchShelf().then(setShelf).catch(() => setShelf(null))
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
      } catch {
        if (info) setError('Failed to load profile data')
      } finally {
        setLoading(false)
      }
    }

    fetchData()
  }, [mounted])

  const monoFont = "var(--font-jetbrains-mono), 'JetBrains Mono', monospace"

  if (!mounted) {
    return <div style={{ background: 'var(--bg)', minHeight: '100vh' }} />
  }

  if (loading) {
    return (
      <SidebarProvider>
        <AppSidebar name="…" studentId="…" />
        <SidebarInset>
          <div className="mb-2 flex justify-start md:hidden px-4 sm:px-6 pt-8">
            <SidebarTrigger variant="ghost" />
          </div>
          <ProfileSkeleton />
        </SidebarInset>
      </SidebarProvider>
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
          <div id="profile-main" className="mx-auto w-full max-w-[820px] min-w-0 px-4 sm:px-6 py-8 sm:py-12 overflow-x-hidden">
          {/* Mobile only: sidebar is an overlay sheet, so content needs an opener.
              Desktop toggles from the sidebar header + edge rail. */}
          <div className="mb-2 flex justify-start md:hidden">
            <SidebarTrigger variant="ghost" />
          </div>
          <div className="border border-mathua-border p-6 text-center bg-mathua-surface min-w-0">
            <h2 className="font-serif text-[1.2rem] text-mathua-primary mb-2">Welcome to your profile</h2>
            <p className="font-mono text-xs text-mathua-secondary mb-4">Sign in to track XP, streaks and mastery. Your activity heatmap will appear here once you start practicing.</p>
            <div className="flex flex-col sm:flex-row flex-wrap gap-3 justify-center">
              <Link href="/login" className="w-full sm:w-auto border border-mathua-blue text-mathua-blue hover:bg-mathua-blue-faint px-6 py-2 font-mono text-xs min-h-[36px] inline-flex items-center justify-center text-center whitespace-nowrap">Sign in</Link>
              <Link href="/learn" className="w-full sm:w-auto border border-mathua-border text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue px-6 py-2 font-mono text-xs min-h-[36px] inline-flex items-center justify-center text-center whitespace-nowrap">Try as guest →</Link>
            </div>
          </div>
          {(Object.keys(progress).length === 0 && !activity.some(d => d.questions > 0)) && (
            <section className="mt-6 border border-mathua-border bg-mathua-surface p-4">
              <h3 className="font-mono text-[11px] text-mathua-muted uppercase tracking-wider mb-3">What to do first</h3>
              <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
                <div className="border border-mathua-border p-3 bg-mathua-surface-elevated">
                  <div className="font-mono text-xs text-mathua-blue mb-1">1. Take diagnostic test</div>
                  <p className="font-mono text-[11px] text-mathua-secondary">Finds your starting point</p>
                  <Link href="/onboard" className="font-mono text-[10px] text-mathua-blue hover:text-mathua-blue-hover mt-2 inline-block">Start diagnostic test →</Link>
                </div>
                <div className="border border-mathua-border p-3 bg-mathua-surface-elevated">
                  <div className="font-mono text-xs text-mathua-blue mb-1">2. Start answering in Learn</div>
                  <p className="font-mono text-[11px] text-mathua-secondary">Questions are generated; nothing to memorise</p>
                  <Link href="/learn" className="font-mono text-[10px] text-mathua-blue hover:text-mathua-blue-hover mt-2 inline-block">Go to Learn →</Link>
                </div>
                <div className="border border-mathua-border p-3 bg-mathua-surface-elevated">
                  <div className="font-mono text-xs text-mathua-blue mb-1">3. Practice → see XP</div>
                  <p className="font-mono text-[11px] text-mathua-secondary">Answer questions · XP shows on Profile</p>
                  <Link href="/learn" className="font-mono text-[10px] text-mathua-blue hover:text-mathua-blue-hover mt-2 inline-block">Continue learning →</Link>
                </div>
              </div>
            </section>
          )}
          {/* Diagnostic CTA — guest */}
          <section className="mt-6 w-full max-w-full min-w-0 overflow-hidden border border-mathua-blue bg-mathua-surface p-4 flex flex-col sm:flex-row items-stretch sm:items-center justify-between gap-3">
            <div className="w-full sm:flex-1 min-w-0 overflow-hidden">
              <div className="flex flex-col items-start gap-1.5 sm:flex-row sm:items-center sm:gap-2 min-w-0">
                <span className="shrink-0 bg-mathua-blue text-white px-2 py-0.5 font-mono text-[10px] uppercase tracking-wider">Recommended</span>
                <span className="min-w-0 break-words [overflow-wrap:anywhere] leading-snug font-mono text-[11px] sm:text-xs text-mathua-primary">Take a diagnostic to get a recommendation on where to start</span>
              </div>
              <p className="font-mono text-xs text-mathua-secondary mt-1 break-words [overflow-wrap:anywhere]">Diagnostic test · finds your starting point</p>
            </div>
            <Link href="/onboard" className="w-full sm:w-auto sm:shrink-0 border border-mathua-blue text-mathua-blue hover:bg-mathua-blue-faint px-6 py-2 font-mono text-xs min-h-[36px] inline-flex items-center justify-center text-center whitespace-nowrap">Start diagnostic test →</Link>
          </section>

          <section id="activity" aria-label="Activity" className="mt-10 flex min-w-0 flex-col items-stretch scroll-mt-28">
            <h2 className="font-serif text-[1.05rem] font-normal text-mathua-primary mb-4 w-full">
              Activity
            </h2>
            <div className="w-full max-w-full min-w-0 flex justify-center overflow-hidden">
              <div className="w-full max-w-full min-w-0">
                <ActivityHeatmap data={activity} />
              </div>
            </div>
          </section>

          <section id="domains" aria-label="By domain" className="mt-8 min-w-0 scroll-mt-28">
            <div className="grid grid-cols-1 gap-4 min-w-0">
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
              <Link href="/login" className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue-faint px-6 py-2 font-mono text-xs min-h-[36px] inline-flex items-center justify-center">Sign in</Link>
              <button type="button" onClick={() => window.location.reload()} className="border border-mathua-border text-mathua-secondary px-6 py-2 font-mono text-xs min-h-[36px] inline-flex items-center justify-center">Retry</button>
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
        <a href="#profile-main" className="sr-only focus:not-sr-only focus:absolute focus:z-50 focus:bg-mathua-surface focus:px-4 focus:py-2 focus:font-mono focus:text-xs focus:text-mathua-blue">
          Skip to profile content
        </a>
        <div id="profile-main" className="mx-auto w-full max-w-[820px] min-w-0 px-4 sm:px-6 py-8 sm:py-12 overflow-x-hidden">
          {/* Mobile only: sidebar is an overlay sheet, so content needs an opener.
              Desktop toggles from the sidebar header + edge rail. */}
          <div className="mb-2 flex justify-start md:hidden">
            <SidebarTrigger variant="ghost" />
          </div>
          <div className="min-w-0">
        {/* Profile stats — mobile-first */}
        <ProfileStats
          name={user.name}
          scores={effScores ?? scores}
          avatarSeed={user.student_id}
          avatarUrl={avatarPreset !== null ? undefined : avatarUrl}
          avatarPreset={avatarPreset}
        />

        <PositionBlock
          catalogue={conceptCatalog}
          progress={progress}
          frontierLabel={frontierLabel}
        />

        {effScores && !isLoggedIn() && (
          <DailyGoalControl
            current={effScores.daily_xp_goal}
            onSaved={g => {
              setGoalOverride(g)
              const sid = user.student_id || getGuestId() || ''
              if (sid) {
                getScores(sid).then(setScores).catch(() => {})
              }
            }}
          />
        )}
        <div className="mt-3 flex flex-wrap gap-x-4 gap-y-1 font-mono text-xs">
          <Link href="/plan" className="text-mathua-blue hover:text-mathua-blue-hover">
            Plan your learning → finish-date estimates
          </Link>
        </div>

        <NextUpSummary shelf={head} />

        {(isNew && !user.diagnostic_completed) && (
          <section className="mt-6 border border-mathua-border bg-mathua-surface p-4">
            <h3 className="font-mono text-[11px] text-mathua-muted uppercase tracking-wider mb-3">What to do first</h3>
            <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
              <div className="border border-mathua-border p-3 bg-mathua-surface-elevated">
                <div className="font-mono text-xs text-mathua-blue mb-1">1. Take diagnostic test</div>
                <p className="font-mono text-[11px] text-mathua-secondary">Finds your starting point</p>
                <Link href="/onboard" className="font-mono text-[10px] text-mathua-blue hover:text-mathua-blue-hover mt-2 inline-block">Start diagnostic test →</Link>
              </div>
              <div className="border border-mathua-border p-3 bg-mathua-surface-elevated">
                <div className="font-mono text-xs text-mathua-blue mb-1">2. Start answering in Learn</div>
                <p className="font-mono text-[11px] text-mathua-secondary">Questions are generated; nothing to memorise</p>
                <Link href="/learn" className="font-mono text-[10px] text-mathua-blue hover:text-mathua-blue-hover mt-2 inline-block">Go to Learn →</Link>
              </div>
              <div className="border border-mathua-border p-3 bg-mathua-surface-elevated">
                <div className="font-mono text-xs text-mathua-blue mb-1">3. Practice → see XP</div>
                <p className="font-mono text-[11px] text-mathua-secondary">Answer questions · XP shows below</p>
                <Link href="/learn" className="font-mono text-[10px] text-mathua-blue hover:text-mathua-blue-hover mt-2 inline-block">Continue learning →</Link>
              </div>
            </div>
          </section>
        )}

        {scores.paused_until && (
          <div className="mt-6 w-full max-w-full min-w-0 overflow-hidden border border-mathua-border bg-mathua-surface p-4 flex flex-col sm:flex-row items-stretch sm:items-center justify-between gap-3">
            <p className="font-mono text-xs text-mathua-primary min-w-0 break-words [overflow-wrap:anywhere] leading-snug">
              ⏸ Paused until {scores.paused_until}: due reviews are hidden
            </p>
            <Link href="/settings" className="w-full sm:w-auto sm:shrink-0 border border-mathua-blue text-mathua-blue hover:bg-mathua-blue-faint px-4 py-2 font-mono text-xs min-h-[36px] inline-flex items-center justify-center text-center whitespace-nowrap">
              Resume
            </Link>
          </div>
        )}

        {dueReviews > 0 && (
          <Link
            href="/review"
            className="mt-6 flex w-full min-w-0 flex-col gap-2 bg-mathua-surface border border-mathua-blue-faint px-4 py-3 hover:border-mathua-blue transition-colors sm:flex-row sm:items-center sm:justify-between"
          >
            <span className="font-mono text-xs text-mathua-blue min-w-0 break-words [overflow-wrap:anywhere] leading-snug">
              ⏳ {dueReviews} concept{dueReviews !== 1 ? 's' : ''} due for review
            </span>
            <span className="font-mono text-[11px] text-mathua-blue border border-mathua-blue-faint px-3 py-1.5 shrink-0 inline-flex items-center justify-center min-h-[36px] w-full sm:w-auto text-center whitespace-nowrap">
              Review Now →
            </span>
          </Link>
        )}

        {/* Quiz gate. Every number here is the server's: whether the gate is open
            (quiz_due), how much it is (quiz_gate_xp) and how far along the learner is
            (xp_since_quiz). The block had a literal 50 in the copy, in the bar and in an
            `xp_total >= 50` fallback for a field the server does send, so ADR-020's rescale
            could move the gate and leave all three describing the old one. */}
        {scores && scores.quiz_due && quizGate && (
          <div className="mt-6 w-full min-w-0 overflow-hidden border border-mathua-blue bg-mathua-surface p-4 flex flex-col sm:flex-row items-stretch sm:items-center justify-between gap-3">
            <div className="w-full sm:flex-1 min-w-0">
              <div className="flex flex-col items-start gap-1.5 sm:flex-row sm:items-center sm:gap-2 flex-wrap">
                <span className="shrink-0 bg-mathua-blue text-white px-2 py-0.5 font-mono text-[10px] uppercase tracking-wider">Quiz due</span>
                <span className="break-words font-mono text-[11px] sm:text-xs text-mathua-primary">{quizGate} XP reached: mastery check recommended</span>
              </div>
              <div className="mt-2 h-1 bg-mathua-code overflow-hidden">
                <div className="h-full bg-mathua-blue" style={{ width: `${Math.min(((scores.xp_since_quiz ?? scores.xp_total) / quizGate) * 100, 100)}%` }} />
              </div>
            </div>
            <Link href="/goals?quiz=1" className="w-full sm:w-auto sm:shrink-0 border border-mathua-blue text-mathua-blue hover:bg-mathua-blue-faint px-6 py-2 font-mono text-xs min-h-[36px] inline-flex items-center justify-center text-center whitespace-nowrap">
              Start the quiz →
            </Link>
          </div>
        )}

        {/* Diagnostic CTA — new users only: hidden once completed or once
            the learner is no longer new. Retake stays URL-reachable. */}
        {!user.diagnostic_completed && isNew && head?.next.kind !== 'diagnostic' && (
        <section className="mt-6 w-full max-w-full min-w-0 overflow-hidden border border-mathua-blue bg-mathua-surface p-4 flex flex-col sm:flex-row items-stretch sm:items-center justify-between gap-3">
          <div className="w-full sm:flex-1 min-w-0 overflow-hidden">
            <div className="flex flex-col items-start gap-1.5 sm:flex-row sm:items-center sm:gap-2 min-w-0">
              <span className="shrink-0 bg-mathua-blue text-white px-2 py-0.5 font-mono text-[10px] uppercase tracking-wider">
                Recommended
              </span>
              <span className="min-w-0 break-words [overflow-wrap:anywhere] leading-snug font-mono text-[11px] sm:text-xs text-mathua-primary">
                Take a diagnostic to get a recommendation on where to start
              </span>
            </div>
            <p className="font-mono text-xs text-mathua-secondary mt-1 break-words [overflow-wrap:anywhere]">
              Diagnostic test · finds your starting point
            </p>
          </div>
          <Link
            href="/onboard"
            className="w-full sm:w-auto sm:shrink-0 border border-mathua-blue text-mathua-blue hover:bg-mathua-blue-faint px-6 py-2 font-mono text-xs min-h-[36px] inline-flex items-center justify-center text-center whitespace-nowrap"
          >
            Start diagnostic test →
          </Link>
        </section>
        )}

        {/* ── The report ────────────────────────────────────────────────
            Restored to the hub: a learner expects the whole picture on their
            own page, not behind a separate destination. Activity leads, then
            the durable picture, then the detail. Cohort figures stay on
            /docs/efficacy — these are the learner's own numbers. */}
        <section id="activity" aria-label="Activity" className="mt-10 flex min-w-0 flex-col items-stretch scroll-mt-28">
          <h2 className="font-serif text-[1.05rem] font-normal text-mathua-primary mb-4 w-full">
            Activity
          </h2>
          <div className="w-full max-w-full min-w-0 flex justify-center overflow-hidden">
            <div className="w-full max-w-full min-w-0">
              <ActivityHeatmap data={activity} />
            </div>
          </div>
        </section>

        <section id="domains" aria-label="By domain" className="mt-8 min-w-0 scroll-mt-28">
          <h2 className="font-serif text-[1.05rem] font-normal text-mathua-primary mb-4">
            By domain
          </h2>
          <DomainProgress progress={progress} />
        </section>

        {efficacy && efficacy.concepts_touched > 0 && (
          <section id="how-doing" aria-label="How you are doing" className="mt-8 min-w-0 scroll-mt-28">
            <h2 className="font-serif text-[1.05rem] font-normal text-mathua-primary mb-4">
              How you&apos;re doing
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

        <section id="struggles" aria-label="Struggles" className="mt-8 min-w-0 scroll-mt-28">
          <StrugglesSection weaknesses={weaknesses} />
          <div className="mt-3 text-center">
            <Link href="/history" className="font-mono text-xs text-mathua-blue hover:text-mathua-blue-hover">
              Every question you&apos;ve answered →
            </Link>
          </div>
        </section>

        {unlockRows.length > 0 && (
          <section aria-label="Recently unlocked" className="mt-8 w-full max-w-full min-w-0 overflow-hidden border border-mathua-border bg-mathua-surface p-4">
            <h3 className="font-mono text-[11px] text-mathua-muted uppercase tracking-wider mb-3">
              Recently unlocked
            </h3>
            <ul className="space-y-1.5">
              {unlockRows.map(r => (
                <li key={r.id}>
                  <Link
                    href={`/learn?concept=${encodeURIComponent(r.id)}`}
                    className="flex items-baseline gap-2 font-mono text-[11px] text-mathua-secondary hover:text-mathua-blue min-w-0"
                  >
                    <span className="shrink-0 uppercase tracking-wider text-mathua-muted">Unlocked</span>
                    <span className="truncate">{r.label}</span>
                    <span className="ml-auto shrink-0 text-mathua-muted">via {r.via}</span>
                  </Link>
                </li>
              ))}
            </ul>
          </section>
        )}
          </div>
        </div>
      </SidebarInset>
    </SidebarProvider>
  )
}
