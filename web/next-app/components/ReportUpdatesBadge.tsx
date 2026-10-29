'use client'

import Link from 'next/link'
import { useEffect, useState } from 'react'
import { listAnnouncements, listMyReports } from '../lib/api'

// A small count beside the avatar. The page that clears it dispatches
// report-updates-changed, so the badge clears without a reload.
export default function ReportUpdatesBadge({ enabled }: { enabled: boolean }) {
  const [unread, setUnread] = useState(0)

  useEffect(() => {
    if (!enabled) {
      setUnread(0)
      return
    }
    let cancelled = false
    const load = async () => {
      const [reports, announcements] = await Promise.allSettled([listMyReports(), listAnnouncements()])
      if (cancelled) return
      let total = 0
      let loaded = false
      if (reports.status === 'fulfilled') {
        total += reports.value.unread
        loaded = true
      }
      if (announcements.status === 'fulfilled') {
        total += announcements.value.unread
        loaded = true
      }
      // A transient failure must not put an error in the header. The page
      // itself reports the failure when the learner opens it.
      if (loaded) setUnread(total)
    }
    void load()
    const onVisible = () => {
      if (document.visibilityState === 'visible') void load()
    }
    window.addEventListener('report-updates-changed', load)
    window.addEventListener('announcements-changed', load)
    window.addEventListener('focus', load)
    document.addEventListener('visibilitychange', onVisible)
    return () => {
      cancelled = true
      window.removeEventListener('report-updates-changed', load)
      window.removeEventListener('announcements-changed', load)
      window.removeEventListener('focus', load)
      document.removeEventListener('visibilitychange', onVisible)
    }
  }, [enabled])

  if (!enabled || unread === 0) return null

  return (
    <Link
      href="/my-reports"
      aria-label={unread === 1 ? '1 unread notification' : `${unread} unread notifications`}
      className="flex min-h-[44px] min-w-[24px] items-center justify-center font-mono text-[10px] text-mathua-blue"
      title="Open notifications"
    >
      <span className="border border-mathua-blue px-1.5 py-0.5 leading-none">{unread}</span>
    </Link>
  )
}
