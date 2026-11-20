'use client'

import { useState } from 'react'
import AppShell from '@/components/app-shell'
import VisualFigure from '@/components/visual/VisualFigure'
import { useAnnouncementsFeed } from '@/components/announcements/useAnnouncementsFeed'
import { FEATURED_FIGURES } from '@/lib/featuredFigures'
import { markAnnouncementsRead } from '@/lib/api'
import { useAuthState } from '@/hooks/useAuthState'

function fmtDate(value?: string): string {
  if (!value) return ''
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return ''
  return date.toLocaleDateString(undefined, { year: 'numeric', month: 'short', day: 'numeric' })
}

/**
 * Where the toast and the card both land. The three figures at the top are the
 * reason the announcement exists, so they are part of the message rather than
 * behind a lesson the reader has not unlocked yet.
 */
export default function AnnouncementsPage() {
  const { loggedIn } = useAuthState()
  const { items, unread, loading, error, reload, markRead } = useAnnouncementsFeed(loggedIn)
  const [marking, setMarking] = useState(false)

  async function markAll() {
    setMarking(true)
    try {
      await markAnnouncementsRead()
      await reload()
    } finally {
      setMarking(false)
    }
  }

  return (
    <AppShell>
      <div className="mt-6 mb-16 px-4">
        <div className="mx-auto max-w-2xl">
          <div className="flex flex-wrap items-end justify-between gap-3">
            <div>
              <h1 className="text-lg text-mathua-primary">Announcements</h1>
              <p className="mt-1 text-[13px] text-mathua-secondary">
                Product news and what changed. {unread > 0 ? `${unread} unread.` : 'Nothing unread.'}
              </p>
            </div>
            {unread > 0 ? (
              <button
                type="button"
                onClick={markAll}
                disabled={marking}
                className="border border-mathua-border px-3 py-1 font-mono text-[11px] text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue disabled:opacity-40"
              >
                {marking ? 'Marking…' : 'Mark all read'}
              </button>
            ) : null}
          </div>

          <section className="mt-8">
            {!loggedIn ? (
              <p className="font-mono text-[11px] text-mathua-muted">
                <a href="/login" className="text-mathua-blue hover:underline">
                  Sign in
                </a>{' '}
                to see announcements for your account.
              </p>
            ) : loading ? (
              <p className="font-mono text-[11px] text-mathua-muted">Loading…</p>
            ) : error ? (
              <div className="border border-mathua-border p-4">
                <p className="font-mono text-[11px] text-mathua-secondary">{error}</p>
                <button
                  type="button"
                  onClick={() => void reload()}
                  className="mt-3 border border-mathua-border px-3 py-1 font-mono text-[11px] text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue"
                >
                  Try again
                </button>
              </div>
            ) : items.length === 0 ? (
              <p className="font-mono text-[11px] text-mathua-muted">No announcements yet.</p>
            ) : (
              <ol className="space-y-4">
                {items.map((item) => (
                  <li
                    key={item.id}
                    className={`border p-4 ${item.read ? 'border-mathua-border' : 'border-mathua-blue bg-mathua-surface'}`}
                  >
                    <div className="flex flex-wrap items-center justify-between gap-2">
                      <p className="text-sm text-mathua-primary">{item.title}</p>
                      <p className="font-mono text-[10px] text-mathua-muted">{fmtDate(item.published_at)}</p>
                    </div>
                    <p className="mt-2 whitespace-pre-wrap text-[13px] text-mathua-secondary">{item.body}</p>
                    {!item.read ? (
                      <button
                        type="button"
                        onClick={() => void markRead(item.id)}
                        className="mt-3 font-mono text-[11px] text-mathua-blue hover:underline"
                      >
                        Mark read
                      </button>
                    ) : null}
                  </li>
                ))}
              </ol>
            )}
          </section>

          <section className="mt-6">
            <p className="font-mono text-[10px] tracking-wider text-mathua-muted">From the figures</p>
            <div className="mt-2 grid grid-cols-1 gap-4">
              {FEATURED_FIGURES.map((figure) => (
                <VisualFigure key={figure.type} visual={figure} />
              ))}
            </div>
          </section>
        </div>
      </div>
    </AppShell>
  )
}
