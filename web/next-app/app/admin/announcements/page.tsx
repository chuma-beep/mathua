'use client'

import { useCallback, useEffect, useState } from 'react'
import AppShell from '../../../components/app-shell'
import AdminNav from '../../../components/AdminNav'
import Footer from '../../../components/Footer'
import KatexContent from '../../../components/KatexContent'
import Loading from '../../../components/Loading'
import SectionHeader from '../../../components/SectionHeader'
import {
  createAdminAnnouncement,
  getErrorMessage,
  listAdminAnnouncements,
  updateAdminAnnouncement,
  type Announcement,
} from '../../../lib/api'
import { useAdminMe } from '../../../hooks/useAdminMe'

function fmtDate(value: string): string {
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

export default function AdminAnnouncementsPage() {
  const { can, loading: meLoading } = useAdminMe()
  const [items, setItems] = useState<Announcement[] | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [title, setTitle] = useState('')
  const [body, setBody] = useState('')
  const [editingId, setEditingId] = useState<string | null>(null)
  const [editingPublished, setEditingPublished] = useState(false)
  const [saving, setSaving] = useState(false)

  const load = useCallback(async () => {
    setLoading(true)
    setError('')
    try {
      const data = await listAdminAnnouncements()
      setItems(data.announcements)
    } catch (e) {
      setError(getErrorMessage(e))
      setItems([])
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    if (!meLoading) void load()
  }, [load, meLoading])

  const resetForm = () => {
    setTitle('')
    setBody('')
    setEditingId(null)
    setEditingPublished(false)
  }

  const save = async (publish: boolean) => {
    if (!title.trim() || !body.trim()) {
      setError('A title and a body are required.')
      return
    }
    setSaving(true)
    setError('')
    try {
      if (editingId) {
        await updateAdminAnnouncement(editingId, {
          title,
          body,
          published: publish ? true : editingPublished,
        })
      } else {
        await createAdminAnnouncement({ title, body, published: publish })
      }
      resetForm()
      await load()
      window.dispatchEvent(new Event('announcements-changed'))
    } catch (e) {
      setError(getErrorMessage(e))
    } finally {
      setSaving(false)
    }
  }

  const togglePublished = async (announcement: Announcement) => {
    setSaving(true)
    setError('')
    try {
      await updateAdminAnnouncement(announcement.id, { published: !announcement.published })
      await load()
      window.dispatchEvent(new Event('announcements-changed'))
    } catch (e) {
      setError(getErrorMessage(e))
    } finally {
      setSaving(false)
    }
  }

  const edit = (announcement: Announcement) => {
    setEditingId(announcement.id)
    setEditingPublished(announcement.published)
    setTitle(announcement.title)
    setBody(announcement.body)
    setError('')
  }

  const canManage = can('content.manage')

  return (
    <AppShell>
      <main className="mx-auto w-full max-w-[900px] min-w-0 px-4 sm:px-6 py-8 sm:py-12 overflow-x-hidden">
        <SectionHeader label="Administration" title="Announcements" />
        <AdminNav />

        {error && (
          <p role="alert" className="mb-4 border border-mathua-red bg-mathua-surface p-4 font-mono text-xs text-mathua-red">
            {error}
          </p>
        )}

        {canManage && (
          <section className="mb-8 border border-mathua-border bg-mathua-surface p-4" aria-labelledby="compose-heading">
            <h2 id="compose-heading" className="font-mono text-sm uppercase tracking-wider text-mathua-muted">
              {editingId ? 'Edit announcement' : 'New announcement'}
            </h2>
            <label htmlFor="announcement-title" className="mt-4 block font-mono text-[11px] text-mathua-secondary">
              Title
            </label>
            <input
              id="announcement-title"
              value={title}
              maxLength={120}
              onChange={(e) => setTitle(e.target.value)}
              className="mt-1 w-full border border-mathua-border bg-mathua-code px-3 h-11 font-mono text-xs text-mathua-primary"
            />
            <label htmlFor="announcement-body" className="mt-4 block font-mono text-[11px] text-mathua-secondary">
              Body, Markdown and LaTeX supported
            </label>
            <textarea
              id="announcement-body"
              value={body}
              maxLength={5000}
              rows={10}
              onChange={(e) => setBody(e.target.value)}
              className="mt-1 w-full border border-mathua-border bg-mathua-code px-3 py-2 font-mono text-xs text-mathua-primary"
            />
            {body.trim() && (
              <div className="mt-4 border border-mathua-border bg-mathua-bg p-4" aria-label="Announcement preview">
                <p className="font-mono text-[10px] uppercase tracking-wider text-mathua-muted">Preview</p>
                <h3 className="mt-2 font-mono text-sm text-mathua-primary">{title || 'Untitled'}</h3>
                <KatexContent className="mt-2 text-sm text-mathua-secondary">{body}</KatexContent>
              </div>
            )}
            <div className="mt-4 flex flex-wrap gap-2">
              <button
                type="button"
                onClick={() => void save(false)}
                disabled={saving}
                className="border border-mathua-border px-3 h-10 font-mono text-xs text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue disabled:opacity-50"
              >
                {saving ? 'Saving...' : editingId ? 'Save changes' : 'Save draft'}
              </button>
              <button
                type="button"
                onClick={() => void save(true)}
                disabled={saving}
                className="border border-mathua-blue bg-mathua-blue px-3 h-10 font-mono text-xs text-white disabled:opacity-50"
              >
                {saving ? 'Saving...' : editingId ? 'Save and publish' : 'Publish'}
              </button>
              {editingId && (
                <button
                  type="button"
                  onClick={resetForm}
                  className="border border-mathua-border px-3 h-10 font-mono text-xs text-mathua-secondary"
                >
                  Cancel
                </button>
              )}
            </div>
          </section>
        )}

        {loading && <Loading label="LOADING ANNOUNCEMENTS" />}

        {!loading && items && items.length === 0 && (
          <p className="border border-mathua-border bg-mathua-surface p-4 text-sm text-mathua-secondary">
            No announcements yet.
          </p>
        )}

        {!loading && items && items.length > 0 && (
          <ul className="space-y-3">
            {items.map((announcement) => (
              <li key={announcement.id} className="border border-mathua-border bg-mathua-surface p-4">
                <div className="flex flex-wrap items-start justify-between gap-2">
                  <div className="min-w-0">
                    <h2 className="font-mono text-sm text-mathua-primary">{announcement.title}</h2>
                    <p className="mt-1 font-mono text-[10px] text-mathua-muted">
                      {announcement.published ? 'Published' : 'Draft'}
                      {announcement.published_at ? ` ${fmtDate(announcement.published_at)}` : ''}
                      {` - created ${fmtDate(announcement.created_at)}`}
                    </p>
                  </div>
                  {canManage && (
                    <div className="flex flex-wrap gap-2">
                      <button
                        type="button"
                        onClick={() => edit(announcement)}
                        className="border border-mathua-border px-3 h-9 font-mono text-[11px] text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue"
                      >
                        Edit
                      </button>
                      <button
                        type="button"
                        onClick={() => void togglePublished(announcement)}
                        disabled={saving}
                        className="border border-mathua-border px-3 h-9 font-mono text-[11px] text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue disabled:opacity-50"
                      >
                        {announcement.published ? 'Unpublish' : 'Publish'}
                      </button>
                    </div>
                  )}
                </div>
                <div className="mt-3 border-t border-mathua-border pt-3">
                  <KatexContent className="text-sm text-mathua-secondary">{announcement.body}</KatexContent>
                </div>
              </li>
            ))}
          </ul>
        )}
      </main>
      <Footer />
    </AppShell>
  )
}
