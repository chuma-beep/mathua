'use client'

import { useCallback, useEffect, useMemo, useState } from 'react'
import { listAnnouncements, markAnnouncementRead, type Announcement } from '@/lib/api'
import { ensureGuestToken } from '@/lib/auth'

interface Feed {
  items: Announcement[]
  unread: number
}

const EMPTY: Feed = { items: [], unread: 0 }

/**
 * The announcements feed and its read state. One hook serves the card and the
 * page, so both show the same order and the same unread count.
 *
 * The list and the count live in one state object because marking read changes
 * both. Splitting them meant the decrement depended on a state updater that
 * React may not have run yet, and the count drifted by one.
 */
export function useAnnouncementsFeed(enabled = true) {
  const [feed, setFeed] = useState<Feed>(EMPTY)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  const reload = useCallback(async () => {
    setLoading(true)
    setError('')
    try {
      // A guest has no account, so the feed answers 401 until a guest token
      // exists. Mint one and try once more. A signed-in learner never reaches
      // the retry.
      const response = await listAnnouncements().catch(async (cause) => {
        const minted = await ensureGuestToken().catch(() => false)
        if (!minted) throw cause
        return listAnnouncements()
      })
      setFeed({ items: response.announcements, unread: response.unread })
    } catch (cause) {
      setError(cause instanceof Error ? cause.message : 'Announcements could not be loaded.')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    if (!enabled) {
      setLoading(false)
      return
    }
    void reload()
  }, [enabled, reload])

  const markRead = useCallback(
    async (id: string) => {
      setFeed((previous) => {
        const target = previous.items.find((item) => item.id === id)
        if (!target || target.read) return previous
        return {
          items: previous.items.map((item) => (item.id === id ? { ...item, read: true } : item)),
          unread: Math.max(0, previous.unread - 1),
        }
      })
      try {
        await markAnnouncementRead(id)
        window.dispatchEvent(new Event('announcements-changed'))
      } catch {
        // A failed write must not leave the list disagreeing with the server.
        void reload()
      }
    },
    [reload],
  )

  const newestUnread = useMemo(() => feed.items.find((item) => !item.read) ?? null, [feed.items])

  return { items: feed.items, unread: feed.unread, loading, error, newestUnread, reload, markRead }
}
