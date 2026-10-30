'use client'

import { useCallback, useEffect, useState } from 'react'
import Link from 'next/link'
import AppShell from '../../components/app-shell'
import SectionHeader from '../../components/SectionHeader'
import Footer from '../../components/Footer'
import Loading from '../../components/Loading'
import Avatar from '../../components/Avatar'
import {
  API_BASE,
  acceptFriendRequest,
  blockFriend,
  cancelFriendRequest,
  declineFriendRequest,
  getErrorMessage,
  getErrorStatus,
  getFriends,
  getFriendsLeaderboard,
  removeFriend,
  rotateFriendCode,
  sendFriendRequest,
  type FriendAvatar,
  type FriendLeaderboardRow,
  type FriendRequest,
  type FriendsFeed,
  type FriendSummary,
} from '../../lib/api'
import { useAuthState } from '../../hooks/useAuthState'
import { dicebearUrl, isDicebearStyle } from '../../lib/dicebear'

type LeaderboardView = 'weekly' | 'all_time'
// The page must distinguish four outcomes that used to collapse into one
// message: signed out, feature disabled, network failure and ready. Treating a
// failed config read as "disabled" told a signed-in learner on a working server
// that friends did not exist.
type Availability = 'checking' | 'ready' | 'guest' | 'disabled' | 'error'

// The friends surface is deliberately one page. Requests, the friend code and
// the friends leaderboard are all about the same small group, and splitting
// them into tabs on first release made an empty account look like three empty
// pages. The sidebar entry only appears when /api/config says the surface is
// enabled, so a deployment that has not cleared the rollout has no dead link.

function avatarImage(avatar: FriendAvatar): string | undefined {
  if (avatar.avatar_url) {
    return avatar.avatar_url.startsWith('/') ? `${API_BASE}${avatar.avatar_url}` : avatar.avatar_url
  }
  const pick = avatar.avatar_dicebear
  return pick && isDicebearStyle(pick.style) ? dicebearUrl(pick.style, pick.seed) : undefined
}

function PersonAvatar({ person, size = 40 }: { person: FriendSummary | FriendRequest; size?: number }) {
  return (
    <Avatar
      seed={person.profile_id}
      name={person.name}
      size={size}
      url={avatarImage(person.avatar)}
      preset={person.avatar.avatar_preset ?? null}
    />
  )
}

function dateLabel(value?: string): string {
  if (!value) return ''
  const date = new Date(value)
  return isNaN(date.getTime()) ? '' : date.toLocaleDateString(undefined, { year: 'numeric', month: 'short', day: 'numeric' })
}

function FriendRow({ friend, onRemove, onBlock, busy }: {
  friend: FriendSummary
  onRemove: () => void
  onBlock: () => void
  busy: boolean
}) {
  return (
    <li className="flex items-center gap-3 border border-mathua-border bg-mathua-surface px-3 py-3 min-w-0">
      <PersonAvatar person={friend} />
      <div className="min-w-0 flex-1">
        <Link href={`/friends/view?id=${encodeURIComponent(friend.profile_id)}`} className="block truncate text-sm text-mathua-primary hover:text-mathua-blue">
          {friend.name}
        </Link>
        <p className="truncate font-mono text-[11px] text-mathua-muted">
          {friend.mastered} mastered{friend.friends_since ? ` · since ${dateLabel(friend.friends_since)}` : ''}
        </p>
      </div>
      <div className="flex shrink-0 items-center gap-1">
        <button
          type="button"
          onClick={onRemove}
          disabled={busy}
          className="h-8 border border-mathua-border px-2 font-mono text-[10px] uppercase tracking-wider text-mathua-secondary hover:border-mathua-red hover:text-mathua-red disabled:opacity-50"
        >
          Remove
        </button>
        <button
          type="button"
          onClick={onBlock}
          disabled={busy}
          className="h-8 border border-mathua-border px-2 font-mono text-[10px] uppercase tracking-wider text-mathua-secondary hover:border-mathua-red hover:text-mathua-red disabled:opacity-50"
          aria-label={`Block ${friend.name}`}
        >
          Block
        </button>
      </div>
    </li>
  )
}

function RequestRow({ request, onAccept, onDecline, onCancel, busy }: {
  request: FriendRequest
  onAccept?: () => void
  onDecline?: () => void
  onCancel?: () => void
  busy: boolean
}) {
  return (
    <li className="flex items-center gap-3 border border-mathua-border px-3 py-3 min-w-0">
      <PersonAvatar person={request} />
      <div className="min-w-0 flex-1">
        <p className="truncate text-sm text-mathua-primary">{request.name}</p>
        <p className="font-mono text-[11px] text-mathua-muted">
          {request.direction === 'incoming' ? 'Wants to be friends' : 'Request sent'} · {dateLabel(request.created_at)}
        </p>
      </div>
      <div className="flex shrink-0 gap-1">
        {request.direction === 'incoming' ? (
          <>
            <button type="button" onClick={onAccept} disabled={busy} className="h-8 bg-mathua-blue px-3 font-mono text-[10px] uppercase tracking-wider text-white disabled:opacity-50">
              Accept
            </button>
            <button type="button" onClick={onDecline} disabled={busy} className="h-8 border border-mathua-border px-3 font-mono text-[10px] uppercase tracking-wider text-mathua-secondary hover:border-mathua-red hover:text-mathua-red disabled:opacity-50">
              Decline
            </button>
          </>
        ) : (
          <button type="button" onClick={onCancel} disabled={busy} className="h-8 border border-mathua-border px-3 font-mono text-[10px] uppercase tracking-wider text-mathua-secondary hover:border-mathua-red hover:text-mathua-red disabled:opacity-50">
            Cancel
          </button>
        )}
      </div>
    </li>
  )
}

export default function FriendsPage() {
  const { loggedIn } = useAuthState()
  const [availability, setAvailability] = useState<Availability>('checking')
  const [feed, setFeed] = useState<FriendsFeed | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [notice, setNotice] = useState('')
  const [query, setQuery] = useState('')
  const [busy, setBusy] = useState(false)
  const [copied, setCopied] = useState(false)
  const [confirmRotate, setConfirmRotate] = useState(false)
  const [view, setView] = useState<LeaderboardView>('weekly')
  const [rows, setRows] = useState<FriendLeaderboardRow[]>([])
  const [rowsLoading, setRowsLoading] = useState(false)
  const [attempt, setAttempt] = useState(0)

  const load = useCallback(async () => {
    const next = await getFriends()
    setFeed(next)
  }, [])

  useEffect(() => {
    if (!loggedIn) {
      setAvailability('guest')
      setLoading(false)
      return
    }
    let cancelled = false
    setLoading(true)
    setError('')
    getFriends()
      .then((next) => {
        if (cancelled) return
        setFeed(next)
        setAvailability('ready')
      })
      .catch((e) => {
        if (cancelled) return
        const status = getErrorStatus(e)
        if (status === 404) {
          setAvailability('disabled')
        } else if (status === 401) {
          setAvailability('guest')
        } else {
          setAvailability('error')
          setError(getErrorMessage(e))
        }
      })
      .finally(() => {
        if (!cancelled) setLoading(false)
      })
    return () => { cancelled = true }
  }, [loggedIn, attempt])

  useEffect(() => {
    if (availability !== 'ready') return
    let cancelled = false
    setRowsLoading(true)
    getFriendsLeaderboard(view)
      .then((next) => { if (!cancelled) setRows(next) })
      .catch(() => { if (!cancelled) setRows([]) })
      .finally(() => { if (!cancelled) setRowsLoading(false) })
    return () => { cancelled = true }
  }, [availability, view, feed])

  const run = async (action: () => Promise<void>, message?: string) => {
    setBusy(true)
    setError('')
    setNotice('')
    try {
      await action()
      await load()
      if (message) setNotice(message)
    } catch (e) {
      setError(getErrorMessage(e))
    } finally {
      setBusy(false)
    }
  }

  const handleSend = async (event: React.FormEvent) => {
    event.preventDefault()
    const value = query.trim()
    if (!value || busy) return
    setBusy(true)
    setError('')
    setNotice('')
    try {
      const status = await sendFriendRequest(value)
      setQuery('')
      await load()
      setNotice(status === 'accepted' ? 'You are now friends.' : 'Friend request sent.')
    } catch (e) {
      setError(getErrorMessage(e))
    } finally {
      setBusy(false)
    }
  }

  const copyCode = async () => {
    if (!feed) return
    try {
      await navigator.clipboard.writeText(feed.identity.friend_code)
      setCopied(true)
      setTimeout(() => setCopied(false), 2000)
    } catch {
      setError('Could not copy the code. Select it and copy manually.')
    }
  }

  const handleRotate = async () => {
    setConfirmRotate(false)
    setBusy(true)
    setError('')
    try {
      const code = await rotateFriendCode()
      setFeed((current) => current ? { ...current, identity: { ...current.identity, friend_code: code } } : current)
    } catch (e) {
      setError(getErrorMessage(e))
    } finally {
      setBusy(false)
    }
  }

  if (loading) {
    return (
      <AppShell>
        <div className="mx-auto w-full max-w-[920px] px-4 py-16 text-center">
          <Loading label="LOADING FRIENDS" />
        </div>
        <Footer />
      </AppShell>
    )
  }

  if (availability === 'guest') {
    return (
      <AppShell>
        <div className="mx-auto w-full max-w-[920px] px-4 sm:px-6 py-8 sm:py-12 min-w-0">
          <SectionHeader label="Your circle" title="Friends" />
          <div className="border border-mathua-border bg-mathua-surface p-6">
            <p className="text-sm text-mathua-primary">Sign in to use friends.</p>
            <p className="mt-2 font-mono text-[11px] text-mathua-muted">
              Friends are per account, so this page needs a signed-in learner.
            </p>
            <Link href="/login" className="mt-4 inline-block font-mono text-[11px] uppercase tracking-wider text-mathua-blue hover:underline">
              Sign in
            </Link>
          </div>
        </div>
        <Footer />
      </AppShell>
    )
  }

  if (availability === 'disabled') {
    return (
      <AppShell>
        <div className="mx-auto w-full max-w-[920px] px-4 sm:px-6 py-8 sm:py-12 min-w-0">
          <SectionHeader label="Your circle" title="Friends" />
          <div className="border border-mathua-border bg-mathua-surface p-6">
            <p className="text-sm text-mathua-primary">Friends are not enabled on this server.</p>
            <p className="mt-2 font-mono text-[11px] text-mathua-muted">
              Start the API with MATHUA_FRIENDS=on to preview this page locally.
            </p>
          </div>
        </div>
        <Footer />
      </AppShell>
    )
  }

  if (availability === 'error') {
    return (
      <AppShell>
        <div className="mx-auto w-full max-w-[920px] px-4 sm:px-6 py-8 sm:py-12 min-w-0">
          <SectionHeader label="Your circle" title="Friends" />
          <div className="border border-mathua-border bg-mathua-surface p-6">
            <p role="alert" className="text-sm text-mathua-primary">
              {error || 'Could not load friends.'}
            </p>
            <button
              type="button"
              onClick={() => setAttempt((value) => value + 1)}
              className="mt-4 h-9 border border-mathua-border px-3 font-mono text-[10px] uppercase tracking-wider text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue"
            >
              Try again
            </button>
          </div>
        </div>
        <Footer />
      </AppShell>
    )
  }

  const friendCount = feed?.friends.length ?? 0
  const incomingCount = feed?.incoming.length ?? 0

  return (
    <AppShell>
      <div className="mx-auto w-full max-w-[920px] min-w-0 px-4 sm:px-6 py-8 sm:py-12 overflow-x-hidden">
        <SectionHeader label="Your circle" title="Friends" />
        <p className="mb-6 max-w-[62ch] text-sm text-mathua-secondary">
          Add someone by username or friend code. Friends see each other&apos;s names and a shared view of mastered concepts.
          Nothing about answers, attempts or weaknesses is shared.
        </p>

        {error && (
          <div role="alert" className="mb-4 border border-mathua-red bg-mathua-surface p-3 font-mono text-xs text-mathua-red">
            {error}
          </div>
        )}
        {notice && (
          <div role="status" className="mb-4 border border-mathua-border bg-mathua-surface p-3 font-mono text-xs text-mathua-secondary">
            {notice}
          </div>
        )}

        <section className="mb-8 border border-mathua-border bg-mathua-surface p-4">
          <div className="flex flex-wrap items-start justify-between gap-4">
            <div className="min-w-0">
              <p className="font-mono text-[10px] uppercase tracking-[0.15em] text-mathua-muted">Your friend code</p>
              <p className="mt-1 font-mono text-2xl tracking-[0.2em] text-mathua-primary">{feed?.identity.friend_code ?? '-'}</p>
              <p className="mt-1 break-all font-mono text-[10px] text-mathua-muted">Profile {feed?.identity.profile_id ?? '-'}</p>
            </div>
            <div className="flex flex-wrap gap-2">
              <button type="button" onClick={copyCode} className="h-9 border border-mathua-border px-3 font-mono text-[10px] uppercase tracking-wider text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue">
                {copied ? 'Copied' : 'Copy code'}
              </button>
              {confirmRotate ? (
                <>
                  <button type="button" onClick={handleRotate} disabled={busy} className="h-9 bg-mathua-red px-3 font-mono text-[10px] uppercase tracking-wider text-white disabled:opacity-50">
                    Confirm rotate
                  </button>
                  <button type="button" onClick={() => setConfirmRotate(false)} className="h-9 border border-mathua-border px-3 font-mono text-[10px] uppercase tracking-wider text-mathua-secondary">
                    Cancel
                  </button>
                </>
              ) : (
                <button type="button" onClick={() => setConfirmRotate(true)} className="h-9 border border-mathua-border px-3 font-mono text-[10px] uppercase tracking-wider text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue">
                  Rotate code
                </button>
              )}
            </div>
          </div>
          {confirmRotate && (
            <p className="mt-3 font-mono text-[11px] text-mathua-red">
              Rotating invalidates the old code. Anyone who saved it cannot use it to find you again.
            </p>
          )}
        </section>

        <section className="mb-8">
          <h2 className="mb-3 font-serif text-[1.05rem] text-mathua-primary">Add a friend</h2>
          <form onSubmit={handleSend} className="flex flex-wrap gap-2">
            <label htmlFor="friend-query" className="sr-only">Username or friend code</label>
            <input
              id="friend-query"
              value={query}
              onChange={(event) => setQuery(event.target.value)}
              placeholder="username or friend code"
              autoComplete="off"
              className="h-11 min-w-0 flex-1 border border-mathua-border bg-mathua-code px-3 font-mono text-xs text-mathua-primary outline-none focus:border-mathua-blue"
            />
            <button
              type="submit"
              disabled={busy || !query.trim()}
              className="h-11 bg-mathua-blue px-4 font-mono text-[11px] uppercase tracking-wider text-white disabled:opacity-50"
            >
              Send request
            </button>
          </form>
          <p className="mt-2 font-mono text-[10px] text-mathua-muted">
            A request is accepted when the other learner sends one back.
          </p>
        </section>

        {incomingCount > 0 && (
          <section className="mb-8">
            <h2 className="mb-3 font-serif text-[1.05rem] text-mathua-primary">
              Requests <span className="font-mono text-xs text-mathua-muted">({incomingCount} incoming)</span>
            </h2>
            <ul className="space-y-2">
              {feed?.incoming.map((request) => (
                <RequestRow
                  key={`in-${request.profile_id}`}
                  request={request}
                  busy={busy}
                  onAccept={() => run(() => acceptFriendRequest(request.profile_id), 'Friend request accepted.')}
                  onDecline={() => run(() => declineFriendRequest(request.profile_id), 'Friend request declined.')}
                />
              ))}
              {feed?.outgoing.map((request) => (
                <RequestRow
                  key={`out-${request.profile_id}`}
                  request={request}
                  busy={busy}
                  onCancel={() => run(() => cancelFriendRequest(request.profile_id), 'Friend request cancelled.')}
                />
              ))}
            </ul>
          </section>
        )}

        <section className="mb-8">
          <h2 className="mb-3 font-serif text-[1.05rem] text-mathua-primary">
            Friends <span className="font-mono text-xs text-mathua-muted">({friendCount})</span>
          </h2>
          {friendCount === 0 ? (
            <div className="border border-mathua-border bg-mathua-surface p-6">
              <p className="text-sm text-mathua-primary">No friends yet.</p>
              <p className="mt-1 font-mono text-[11px] text-mathua-muted">
                Share your friend code, or add someone with theirs.
              </p>
            </div>
          ) : (
            <ul className="space-y-2">
              {feed?.friends.map((friend) => (
                <FriendRow
                  key={friend.profile_id}
                  friend={friend}
                  busy={busy}
                  onRemove={() => run(() => removeFriend(friend.profile_id), `${friend.name} removed.`)}
                  onBlock={() => run(() => blockFriend(friend.profile_id), `${friend.name} blocked.`)}
                />
              ))}
            </ul>
          )}
        </section>

        <section>
          <div className="mb-3 flex flex-wrap items-center justify-between gap-2">
            <h2 className="font-serif text-[1.05rem] text-mathua-primary">Friends leaderboard</h2>
            <div className="flex gap-2" role="tablist" aria-label="Friends leaderboard period">
              {(['weekly', 'all_time'] as LeaderboardView[]).map((id) => (
                <button
                  key={id}
                  type="button"
                  role="tab"
                  aria-selected={view === id}
                  onClick={() => setView(id)}
                  className={`h-8 border px-3 font-mono text-[10px] uppercase tracking-wider ${
                    view === id
                      ? 'border-mathua-blue bg-mathua-blue text-white'
                      : 'border-mathua-border text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue'
                  }`}
                >
                  {id === 'weekly' ? 'This week' : 'All time'}
                </button>
              ))}
            </div>
          </div>
          {rowsLoading ? (
            <p className="font-mono text-[11px] text-mathua-muted">Loading standings</p>
          ) : rows.length === 0 ? (
            <div className="border border-mathua-border bg-mathua-surface p-6">
              <p className="text-sm text-mathua-primary">No standings yet.</p>
              <p className="mt-1 font-mono text-[11px] text-mathua-muted">Add a friend or earn XP to appear here.</p>
            </div>
          ) : (
            <ol className="divide-y divide-mathua-border border border-mathua-border">
              {rows.map((row, index) => (
                <li key={row.profile_id} className="flex items-center gap-3 px-3 py-2.5">
                  <span className="w-6 shrink-0 font-mono text-[11px] text-mathua-muted">{index + 1}</span>
                  <Avatar
                    seed={row.profile_id}
                    name={row.name}
                    size={28}
                    url={avatarImage(row.avatar)}
                    preset={row.avatar.avatar_preset ?? null}
                  />
                  <span className="min-w-0 flex-1 truncate text-sm text-mathua-primary">
                    {row.name}{row.is_me ? ' (you)' : ''}
                  </span>
                  <span className="shrink-0 font-mono text-[11px] tabular-nums text-mathua-secondary">{row.xp} XP</span>
                </li>
              ))}
            </ol>
          )}
        </section>
      </div>
      <Footer />
    </AppShell>
  )
}
