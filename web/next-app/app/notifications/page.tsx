'use client'

import { useCallback, useEffect, useState } from 'react'
import Link from 'next/link'
import AppShell from '@/components/app-shell'
import {
  listNotifications,
  markAllNotificationsRead,
  reportNotificationOpen,
  type AppNotification,
} from '@/lib/api'
import { useAuthState } from '@/hooks/useAuthState'

// The archive is the durable record of every reminder the pipeline produced.
// The badge and the dashboard only care about the newest unread items; this
// page is where a learner can see what arrived last month and what kind it was.
const KINDS: Array<{ value: string; label: string }> = [
  { value: '', label: 'All' },
  { value: 'review', label: 'Reviews' },
  { value: 'study', label: 'Study' },
  { value: 'milestone', label: 'Milestones' },
  { value: 'inactivity', label: 'Inactivity' },
]

function fmtDate(value?: string): string {
  if (!value) return ''
  const d = new Date(value)
  if (Number.isNaN(d.getTime())) return value
  return d.toLocaleString(undefined, {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

function kindLabel(kind: string): string {
  return KINDS.find(k => k.value === kind)?.label.replace(/s$/, '') || 'Notification'
}

export default function NotificationsArchivePage() {
  const { loggedIn } = useAuthState()
  const [notifications, setNotifications] = useState<AppNotification[]>([])
  const [unread, setUnread] = useState(0)
  const [cursor, setCursor] = useState('')
  const [hasMore, setHasMore] = useState(false)
  const [kind, setKind] = useState('')
  const [unreadOnly, setUnreadOnly] = useState(false)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [marking, setMarking] = useState(false)

  const load = useCallback(
    async (opts: { cursor?: string; append: boolean; kind: string; unreadOnly: boolean }) => {
      setLoading(true)
      setError('')
      try {
        const page = await listNotifications({
          limit: 20,
          cursor: opts.cursor ?? '',
          kind: opts.kind,
          unreadOnly: opts.unreadOnly,
        })
        setNotifications(prev => (opts.append ? [...prev, ...page.notifications] : page.notifications))
        setUnread(page.unread)
        setCursor(page.next_cursor)
        setHasMore(page.has_more)
      } catch (e) {
        setError(e instanceof Error ? e.message : 'Your notification history could not be loaded.')
      } finally {
        setLoading(false)
      }
    },
    [],
  )

  useEffect(() => {
    if (!loggedIn) return
    void load({ append: false, kind, unreadOnly })
  }, [loggedIn, kind, unreadOnly, load])

  function selectKind(next: string) {
    if (next === kind) return
    setCursor('')
    setHasMore(false)
    setKind(next)
  }

  function selectUnread(next: boolean) {
    if (next === unreadOnly) return
    setCursor('')
    setHasMore(false)
    setUnreadOnly(next)
  }

  async function handleOpen(notification: AppNotification) {
    setNotifications(prev =>
      prev.map(n => (n.id === notification.id ? { ...n, read: true } : n)),
    )
    try {
      const remaining = await reportNotificationOpen(notification.id)
      setUnread(remaining)
      window.dispatchEvent(new Event('notifications-changed'))
    } catch {
      // The link still works. A failed stamp only costs one engagement sample.
    }
    if (unreadOnly) {
      setNotifications(prev => prev.filter(n => n.id !== notification.id))
    }
  }

  async function handleMarkAllRead() {
    setMarking(true)
    setError('')
    try {
      await markAllNotificationsRead()
      setNotifications(prev => (unreadOnly ? [] : prev.map(n => ({ ...n, read: true }))))
      setUnread(0)
      window.dispatchEvent(new Event('notifications-changed'))
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Could not mark notifications read.')
    } finally {
      setMarking(false)
    }
  }

  function handleMore() {
    if (!cursor || loading) return
    void load({ cursor, append: true, kind, unreadOnly })
  }

  return (
    <AppShell>
      <main className="max-w-container mx-auto px-4 md:px-6 py-10 pb-28">
        <div className="flex flex-wrap items-end justify-between gap-3">
          <div>
            <h1 className="font-mono text-xl text-mathua-primary">Notification history</h1>
            <p className="mt-1 text-sm text-mathua-secondary">
              Every reminder, milestone and nudge Mathua sent you, newest first.
            </p>
          </div>
          <button
            type="button"
            onClick={() => void handleMarkAllRead()}
            disabled={marking || unread === 0}
            className="font-mono text-[11px] text-mathua-muted underline underline-offset-2 hover:text-mathua-blue disabled:opacity-40"
          >
            {marking ? 'Marking…' : 'Mark all read'}
          </button>
        </div>

        <div className="mt-6 flex flex-wrap items-center gap-2">
          {KINDS.map(option => (
            <button
              key={option.value || 'all'}
              type="button"
              aria-pressed={kind === option.value}
              onClick={() => selectKind(option.value)}
              className={`border px-3 h-8 font-mono text-[11px] ${kind === option.value ? 'border-mathua-blue text-mathua-blue bg-mathua-blue-faint' : 'border-mathua-border text-mathua-secondary'}`}
            >
              {option.label}
            </button>
          ))}
          <span className="mx-1 h-4 w-px bg-mathua-border" aria-hidden="true" />
          <button
            type="button"
            aria-pressed={!unreadOnly}
            onClick={() => selectUnread(false)}
            className={`border px-3 h-8 font-mono text-[11px] ${!unreadOnly ? 'border-mathua-blue text-mathua-blue bg-mathua-blue-faint' : 'border-mathua-border text-mathua-secondary'}`}
          >
            All
          </button>
          <button
            type="button"
            aria-pressed={unreadOnly}
            onClick={() => selectUnread(true)}
            className={`border px-3 h-8 font-mono text-[11px] ${unreadOnly ? 'border-mathua-blue text-mathua-blue bg-mathua-blue-faint' : 'border-mathua-border text-mathua-secondary'}`}
          >
            Unread{unread > 0 ? ` (${unread})` : ''}
          </button>
        </div>

        {!loggedIn && (
          <p className="mt-8 border border-mathua-border bg-mathua-surface p-4 text-sm text-mathua-secondary">
            <Link href="/login" className="text-mathua-blue underline underline-offset-2">
              Sign in
            </Link>{' '}
            to see your notification history.
          </p>
        )}

        {error && (
          <p className="mt-6 border border-[var(--accent-red)]/40 bg-mathua-surface p-4 font-mono text-xs text-mathua-secondary">
            {error}
          </p>
        )}

        {loggedIn && loading && notifications.length === 0 && (
          <p className="mt-6 font-mono text-xs text-mathua-muted">Loading notifications.</p>
        )}

        {loggedIn && !loading && notifications.length === 0 && !error && (
          <p className="mt-6 border border-mathua-border bg-mathua-surface p-4 text-sm text-mathua-secondary">
            {unreadOnly
              ? 'No unread notifications in this view.'
              : 'Nothing here yet. Reminders and updates will appear as they are sent.'}
          </p>
        )}

        {notifications.length > 0 && (
          <ul className="mt-6 space-y-3">
            {notifications.map(notification => (
              <li key={notification.id} className="border border-mathua-border bg-mathua-surface p-4">
                <div className="flex flex-wrap items-baseline justify-between gap-2">
                  <div className="flex min-w-0 items-baseline gap-2">
                    {!notification.read && (
                      <span
                        className="h-2 w-2 shrink-0 rounded-full bg-mathua-blue"
                        aria-label="Unread notification"
                      />
                    )}
                    {notification.href ? (
                      <Link
                        href={notification.href}
                        onClick={() => void handleOpen(notification)}
                        className="font-mono text-xs text-mathua-primary underline underline-offset-2 hover:text-mathua-blue"
                      >
                        {notification.title}
                      </Link>
                    ) : (
                      <span className="font-mono text-xs text-mathua-primary">{notification.title}</span>
                    )}
                  </div>
                  <span className="font-mono text-[10px] text-mathua-muted">
                    {fmtDate(notification.created_at)}
                  </span>
                </div>
                {notification.body && (
                  <p className="mt-2 text-sm text-mathua-secondary">{notification.body}</p>
                )}
                <p className="mt-2 font-mono text-[10px] uppercase tracking-wider text-mathua-muted">
                  {kindLabel(notification.kind)}
                </p>
              </li>
            ))}
          </ul>
        )}

        {hasMore && (
          <button
            type="button"
            onClick={handleMore}
            disabled={loading}
            className="mt-4 border border-mathua-border px-4 h-9 font-mono text-xs text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue disabled:opacity-50"
          >
            {loading ? 'Loading…' : 'Load more'}
          </button>
        )}
      </main>
    </AppShell>
  )
}
