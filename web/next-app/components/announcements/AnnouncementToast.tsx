'use client'

import { useCallback, useEffect, useRef } from 'react'
import { toast } from 'sonner'
import { markAnnouncementRead, type Announcement } from '@/lib/api'

// Routes where a toast must never appear. A graded question or a review session
// is the one place an interruption costs more than the message is worth.
const SUPPRESSED = ['/session', '/timed-quiz', '/review', '/announcements']

function firstLine(body: string): string {
  const line = body.split(/\r?\n/).map((part) => part.trim()).find(Boolean) ?? ''
  return line.length > 150 ? `${line.slice(0, 147)}…` : line
}

function suppressedHere(pathname: string): boolean {
  return SUPPRESSED.some((path) => pathname === path || pathname.startsWith(`${path}/`))
}

/**
 * The announcement arrives as a toast and collapses into a bell when the reader
 * touches the page without opening it. The bell keeps a red dot and stays until
 * the announcement is opened, so an ignored toast is never the same as a read
 * one.
 *
 * The collapsed bell sits where the toast left from: bottom centre on a phone,
 * lifted above the tab bar, and bottom right on a desktop.
 */
export default function AnnouncementToast({ announcement }: { announcement: Announcement | null }) {
  const shown = useRef<string | null>(null)
  const toastId = announcement ? `announcement-${announcement.id}` : null
  const unread = Boolean(announcement && !announcement.read)

  const record = useCallback(() => {
    if (!announcement) return
    void markAnnouncementRead(announcement.id).catch(() => {
      // The shell refreshes the count; a failed write must not break the page.
    })
  }, [announcement])

  useEffect(() => {
    if (!announcement || !unread || toastId === null) return
    if (typeof window === 'undefined' || suppressedHere(window.location.pathname)) return
    if (shown.current === announcement.id) return

    // The Toaster subscribes in its own effect, and this effect runs first.
    // Emitting on the next tick means the subscriber exists. The guard moves
    // inside the timer because React runs effects twice in development.
    const timer = window.setTimeout(() => {
      if (shown.current === announcement.id) return
      shown.current = announcement.id
      toast.custom(
        () => (
          <div className="w-[min(90vw,360px)] border border-mathua-border bg-mathua-bg p-3 font-mono shadow-lg">
            <p className="text-[10px] tracking-wider text-mathua-muted">Announcement</p>
            <p className="mt-1 text-[13px] leading-snug text-mathua-primary">{announcement.title}</p>
            <p className="mt-1 line-clamp-2 text-[11px] leading-snug text-mathua-secondary">
              {firstLine(announcement.body)}
            </p>
            <div className="mt-3 flex items-center gap-3">
              <button
                type="button"
                onClick={() => {
                  record()
                  toast.dismiss(toastId)
                }}
                className="border border-mathua-blue px-2 py-0.5 text-[11px] text-mathua-blue hover:bg-mathua-blue-faint"
              >
                Read
              </button>
              <button
                type="button"
                onClick={() => {
                  // Dismissing is not reading: the bell in the assistant dock
                  // keeps its dot until the announcement is opened.
                  toast.dismiss(toastId)
                }}
                className="text-[11px] text-mathua-muted hover:text-mathua-primary"
              >
                Dismiss
              </button>
            </div>
          </div>
        ),
        { id: toastId, duration: 9000 },
      )
    }, 0)

    return () => window.clearTimeout(timer)
  }, [announcement, unread, toastId, record])

  return null
}
