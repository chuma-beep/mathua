'use client'

import { useEffect, useRef, useState } from 'react'
import { usePathname, useRouter } from 'next/navigation'
import { toast } from '@/components/ui/sonner'
import {
  listAnnouncements,
  listMyReports,
  listNotifications,
  markAnnouncementRead,
  markMyReportsRead,
  markNotificationRead,
} from '@/lib/api'
import { getToken } from '@/lib/auth'

// A graded question or a review session is the one place an interruption costs
// more than the message is worth. The inboxes are excluded because opening one
// is already the reader's answer to the message.
const SUPPRESSED = [
  '/session',
  '/timed-quiz',
  '/review',
  '/announcements',
  '/notifications',
  '/my-reports',
]

const KIND_LABELS: Record<string, string> = {
  review: 'Review reminder',
  study: 'Study reminder',
  milestone: 'Milestone',
  inactivity: 'Reminder',
}

const REPORT_STATUS: Record<string, string> = {
  open: 'open',
  reviewing: 'being reviewed',
  resolved: 'resolved',
  dismissed: 'dismissed',
}

type Candidate = {
  key: string
  label: string
  title: string
  description: string
  href: string
  at: number
  markRead: () => Promise<void>
}

function time(value: string | undefined): number {
  const parsed = value ? Date.parse(value) : NaN
  return Number.isFinite(parsed) ? parsed : 0
}

function firstLine(body: string): string {
  const line =
    body
      .split(/\r?\n/)
      .map((part) => part.trim())
      .find(Boolean) ?? ''
  return line.length > 160 ? `${line.slice(0, 157)}…` : line
}

function suppressedHere(pathname: string | null): boolean {
  if (!pathname) return false
  return SUPPRESSED.some((path) => pathname === path || pathname.startsWith(`${path}/`))
}

/**
 * Every unread item from the three inboxes arrives as one sonner toast at the
 * bottom of the screen. Reading it records the item and opens its page.
 * Dismissing or letting it time out records nothing, so the bell keeps its red
 * dot and the item stays in the inbox. There is no custom toast markup: the
 * shadcn Toaster already carries the app's border, type and position.
 */
export default function NotificationToasts() {
  const pathname = usePathname()
  const router = useRouter()
  const [candidates, setCandidates] = useState<Candidate[]>([])
  const shown = useRef(new Set<string>())

  useEffect(() => {
    let live = true
    const load = async () => {
      // Guests without a token have no inbox, and asking anyway only produces
      // a 401 on pages whose shell has not minted a guest token yet.
      if (!getToken()) return
      const [announcements, notifications, reports] = await Promise.allSettled([
        listAnnouncements(),
        listNotifications({ unreadOnly: true, limit: 10 }),
        listMyReports(),
      ])
      if (!live) return

      const next: Candidate[] = []
      if (announcements.status === 'fulfilled') {
        for (const item of announcements.value.announcements) {
          if (item.read) continue
          next.push({
            key: `announcement:${item.id}`,
            label: 'Announcement',
            title: item.title,
            description: firstLine(item.body),
            href: '/announcements',
            at: time(item.published_at || item.created_at),
            markRead: async () => {
              await markAnnouncementRead(item.id)
              window.dispatchEvent(new Event('announcements-changed'))
            },
          })
        }
      }
      if (notifications.status === 'fulfilled') {
        for (const item of notifications.value.notifications) {
          if (item.read) continue
          next.push({
            key: `notification:${item.id}`,
            label: KIND_LABELS[item.kind] ?? 'Notification',
            title: item.title,
            description: firstLine(item.body),
            href: item.href || '/notifications',
            at: time(item.created_at),
            markRead: async () => {
              await markNotificationRead(item.id)
              window.dispatchEvent(new Event('notifications-changed'))
            },
          })
        }
      }
      if (reports.status === 'fulfilled') {
        for (const report of reports.value.reports) {
          if (!report.unread) continue
          next.push({
            key: `report:${report.id}`,
            label: 'Report update',
            title: `Report #${report.id} is ${REPORT_STATUS[report.status] ?? report.status}`,
            description: firstLine(report.resolution) || 'Open your reports to see the update.',
            href: '/my-reports',
            at: time(report.resolved_at || report.created_at),
            markRead: async () => {
              await markMyReportsRead()
              window.dispatchEvent(new Event('report-updates-changed'))
            },
          })
        }
      }
      setCandidates(next.sort((a, b) => b.at - a.at))
    }

    void load()
    const onVisible = () => {
      if (document.visibilityState === 'visible') void load()
    }
    window.addEventListener('focus', load)
    window.addEventListener('auth-changed', load)
    window.addEventListener('announcements-changed', load)
    window.addEventListener('notifications-changed', load)
    window.addEventListener('report-updates-changed', load)
    document.addEventListener('visibilitychange', onVisible)
    return () => {
      live = false
      window.removeEventListener('focus', load)
      window.removeEventListener('auth-changed', load)
      window.removeEventListener('announcements-changed', load)
      window.removeEventListener('notifications-changed', load)
      window.removeEventListener('report-updates-changed', load)
      document.removeEventListener('visibilitychange', onVisible)
    }
  }, [])

  useEffect(() => {
    if (suppressedHere(pathname)) return
    if (typeof document !== 'undefined' && document.visibilityState === 'hidden') return
    const next = candidates.find((candidate) => !shown.current.has(candidate.key))
    if (!next) return
    shown.current.add(next.key)

    const id = `mathua-${next.key}`
    toast(next.title, {
      id,
      description: `${next.label} · ${next.description}`,
      duration: 9000,
      action: {
        label: 'Read',
        onClick: () => {
          toast.dismiss(id)
          void next
            .markRead()
            .catch(() => {
              // The inbox retries the read when it opens; the click still travels.
            })
            .finally(() => router.push(next.href))
        },
      },
      cancel: {
        label: 'Dismiss',
        onClick: () => {
          // Dismissing is not reading: the bell keeps its dot.
        },
      },
    })
  }, [candidates, pathname, router])

  return null
}
