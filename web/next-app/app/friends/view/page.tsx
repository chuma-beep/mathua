'use client'

import { Suspense, useEffect, useState } from 'react'
import Link from 'next/link'
import { useSearchParams } from 'next/navigation'
import AppShell from '../../../components/app-shell'
import SectionHeader from '../../../components/SectionHeader'
import Footer from '../../../components/Footer'
import Loading from '../../../components/Loading'
import Avatar from '../../../components/Avatar'
import DomainTable from '../../../components/DomainTable'
import { concepts } from '../../../lib/conceptData'
import { buildDomainRows } from '../../../lib/domainRows'
import { API_BASE, getErrorMessage, getFriendGraph, getFriendProfile, type ConceptProgress, type FriendGraphNode, type FriendSummary } from '../../../lib/api'
import { dicebearUrl, isDicebearStyle } from '../../../lib/dicebear'

// The friend profile is a read-only shared graph. It uses the same domain
// builder as /domains and /share so a friend sees the same shape the owner
// sees, but with no links, no weakness, no attempts and no due dates.
//
// The address is ?id=<profile_id> because this site is a static export; a
// dynamic segment would only prerender ids that existed at build time.

function avatarImage(person: FriendSummary): string | undefined {
  if (person.avatar.avatar_url) {
    return person.avatar.avatar_url.startsWith('/') ? `${API_BASE}${person.avatar.avatar_url}` : person.avatar.avatar_url
  }
  const pick = person.avatar.avatar_dicebear
  return pick && isDicebearStyle(pick.style) ? dicebearUrl(pick.style, pick.seed) : undefined
}

function FriendView() {
  const searchParams = useSearchParams()
  const profileID = searchParams.get('id') || ''
  const [profile, setProfile] = useState<FriendSummary | null>(null)
  const [nodes, setNodes] = useState<FriendGraphNode[]>([])
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (!profileID) {
      setError('Missing friend.')
      setLoading(false)
      return
    }
    let cancelled = false
    Promise.all([getFriendProfile(profileID), getFriendGraph(profileID)])
      .then(([nextProfile, nextNodes]) => {
        if (cancelled) return
        setProfile(nextProfile)
        setNodes(nextNodes)
      })
      .catch((e) => {
        if (!cancelled) setError(getErrorMessage(e))
      })
      .finally(() => {
        if (!cancelled) setLoading(false)
      })
    return () => { cancelled = true }
  }, [profileID])

  if (loading) {
    return (
      <AppShell>
        <div className="mx-auto w-full max-w-[920px] px-4 py-16 text-center">
          <Loading label="LOADING FRIEND" />
        </div>
        <Footer />
      </AppShell>
    )
  }

  if (error || !profile) {
    return (
      <AppShell>
        <div className="mx-auto w-full max-w-[920px] px-4 sm:px-6 py-8 sm:py-12 min-w-0">
          <SectionHeader label="Your circle" title="Friend" />
          <div className="border border-mathua-border bg-mathua-surface p-6">
            <p className="text-sm text-mathua-primary">{error || 'This profile is not available.'}</p>
            <p className="mt-2 font-mono text-[11px] text-mathua-muted">
              A friend may have turned off friends or removed the friendship.
            </p>
            <Link href="/friends" className="mt-4 inline-block font-mono text-[11px] uppercase tracking-wider text-mathua-blue hover:underline">
              Back to friends
            </Link>
          </div>
        </div>
        <Footer />
      </AppShell>
    )
  }

  const progress: Record<string, ConceptProgress> = {}
  for (const node of nodes) {
    progress[node.concept_id] = { status: node.status, streak: 0 }
  }
  const rows = buildDomainRows(concepts, progress)
  const mastered = rows.reduce((total, row) => total + row.mastered, 0)

  return (
    <AppShell>
      <div className="mx-auto w-full max-w-[920px] min-w-0 px-4 sm:px-6 py-8 sm:py-12 overflow-x-hidden">
        <Link href="/friends" className="mb-4 inline-block font-mono text-[11px] uppercase tracking-wider text-mathua-muted hover:text-mathua-blue">
          Back to friends
        </Link>
        <SectionHeader label="Shared graph" title={profile.name} />

        <section className="mb-8 flex flex-wrap items-center gap-4 border border-mathua-border bg-mathua-surface p-4">
          <Avatar
            seed={profile.profile_id}
            name={profile.name}
            size={56}
            url={avatarImage(profile)}
            preset={profile.avatar.avatar_preset ?? null}
          />
          <div className="min-w-0 flex-1">
            {profile.username && <p className="font-mono text-[11px] text-mathua-muted">@{profile.username}</p>}
            <p className="mt-1 text-sm text-mathua-primary">
              {mastered} concept{mastered === 1 ? '' : 's'} mastered
            </p>
            <p className="font-mono text-[11px] text-mathua-muted">
              {profile.xp_total} XP total
              {profile.friends_since ? ` · friends since ${new Date(profile.friends_since).toLocaleDateString()}` : ''}
            </p>
          </div>
        </section>

        <p className="mb-4 max-w-[62ch] text-sm text-mathua-secondary">
          Only concept states are shared. Answers, attempts, weaknesses and review dates stay private.
        </p>

        {rows.length === 0 ? (
          <div className="border border-mathua-border bg-mathua-surface p-6">
            <p className="text-sm text-mathua-primary">Nothing to show yet.</p>
            <p className="mt-1 font-mono text-[11px] text-mathua-muted">Concepts appear as they are studied.</p>
          </div>
        ) : (
          <DomainTable rows={rows} links={false} />
        )}
      </div>
      <Footer />
    </AppShell>
  )
}

export default function FriendViewPage() {
  return (
    <Suspense fallback={
      <AppShell>
        <div className="mx-auto w-full max-w-[920px] px-4 py-16 text-center">
          <Loading label="LOADING FRIEND" />
        </div>
      </AppShell>
    }>
      <FriendView />
    </Suspense>
  )
}
