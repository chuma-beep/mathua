'use client'

import { useCallback, useEffect, useState } from 'react'
import Link from 'next/link'
import AppShell from '../../../components/app-shell'
import Footer from '../../../components/Footer'
import SectionHeader from '../../../components/SectionHeader'
import AdminNav from '../../../components/AdminNav'
import { useAdminMe } from '../../../hooks/useAdminMe'
import { listAdminAudit, getErrorMessage, type AdminAuditEvent } from '../../../lib/api'

// The audit trail.
//
// Read-only, and not by convention: the server registers GET for this path and answers 405 to
// everything else, and the table has no update or delete method on either store. So there is no
// "fix a typo in an event" affordance here, which is the correct absence, an operator who
// notices a wrong row must not be able to correct it, because the value of the trail is that it
// records what the system actually did.
//
// No filter and no search yet. Both are future work and neither is needed at this size; a
// filter control that hides events invites an administrator to curate the view of their own
// trail, which is the wrong direction.
const PAGE_SIZE = 50

export default function AdminAuditPage() {
  const [events, setEvents] = useState<AdminAuditEvent[] | null>(null)
  const [error, setError] = useState('')
  const [openId, setOpenId] = useState<number | null>(null)
  const [query, setQuery] = useState('')
  const [appliedQuery, setAppliedQuery] = useState('')
  const [offset, setOffset] = useState(0)
  const [total, setTotal] = useState(0)

  const load = useCallback(() => {
    setError('')
    listAdminAudit(appliedQuery, PAGE_SIZE, offset)
      .then((res) => { setEvents(res.events); setTotal(res.total ?? res.events.length) })
      .catch((e) => { setError(getErrorMessage(e)); setEvents([]); setTotal(0) })
  }, [appliedQuery, offset])

  useEffect(() => { load() }, [load])

  const { isStaff } = useAdminMe()

  return (
    <AppShell>
      <main className="mx-auto w-full max-w-[820px] min-w-0 px-4 sm:px-6 py-8 sm:py-12 overflow-x-hidden">
        <SectionHeader label="Administration" title="Audit" />
        <AdminNav />

        <form
          onSubmit={(event) => { event.preventDefault(); setAppliedQuery(query.trim()); setOffset(0) }}
          className="mb-4 flex flex-wrap gap-2"
        >
          <label htmlFor="audit-search" className="sr-only">Search the audit log</label>
          <input
            id="audit-search"
            value={query}
            onChange={(event) => setQuery(event.target.value)}
            placeholder="Search action, actor, entity or payload"
            className="h-9 min-w-0 flex-1 border border-mathua-border bg-mathua-code px-3 font-mono text-xs text-mathua-primary outline-none focus:border-mathua-blue"
          />
          <button type="submit" className="h-9 border border-mathua-blue px-3 font-mono text-xs text-mathua-blue hover:bg-mathua-blue hover:text-white">
            Search
          </button>
          {(appliedQuery || query) && (
            <button
              type="button"
              onClick={() => { setQuery(''); setAppliedQuery(''); setOffset(0) }}
              className="h-9 border border-mathua-border px-3 font-mono text-xs text-mathua-secondary hover:text-mathua-blue"
            >
              Clear
            </button>
          )}
        </form>

        {!isStaff && (
          <p role="status" className="border border-mathua-border bg-mathua-surface p-4 font-mono text-xs text-mathua-secondary">
            This area needs an administrator role.
          </p>
        )}

        {error && (
          <p role="alert" className="mb-4 border border-mathua-red bg-mathua-surface p-4 font-mono text-xs text-mathua-red">
            {error}
          </p>
        )}

        {events === null && !error && (
          <p className="font-mono text-xs text-mathua-muted">Loading the log…</p>
        )}

        {events !== null && events.length === 0 && (
          <p className="font-mono text-xs text-mathua-muted">
            No administrative activity has been recorded yet.
          </p>
        )}

        {events !== null && events.length > 0 && (
          <ol className="space-y-2">
            {events.map((e) => {
              const open = openId === e.id
              return (
                <li key={e.id} className="border border-mathua-border bg-mathua-surface min-w-0">
                  <button
                    type="button"
                    onClick={() => setOpenId(open ? null : e.id)}
                    aria-expanded={open}
                    className="w-full text-left px-3 py-3 font-mono text-[11px] min-w-0 hover:bg-mathua-surface-highlight"
                  >
                    <span className="block text-mathua-muted">{e.created_at}</span>
                    <span className="block text-mathua-primary">{e.action}</span>
                    <span className="block text-mathua-muted truncate">
                      by {e.actor_id || 'unknown'} → {e.entity_type} {e.entity_id}
                    </span>
                    {!open && <span className="block text-mathua-blue">Before / after →</span>}
                  </button>
                  {open && (
                    <div className="px-3 pb-3 min-w-0">
                      <div className="grid grid-cols-1 sm:grid-cols-2 gap-2">
                        <Snapshot label="Before" json={e.before_json} />
                        <Snapshot label="After" json={e.after_json} />
                      </div>
                      <p className="mt-2 font-mono text-[10px] text-mathua-muted">
                        Event {e.id}. This record is written once and cannot be edited or deleted.
                      </p>
                    </div>
                  )}
                </li>
              )
            })}
          </ol>
        )}

        {total > PAGE_SIZE && (
          <div className="mt-6 flex flex-wrap items-center justify-between gap-3 border-t border-mathua-border pt-4">
            <span className="font-mono text-[11px] text-mathua-muted">
              {offset + 1}–{Math.min(offset + (events?.length ?? 0), total)} of {total}
            </span>
            <div className="flex gap-2">
              <button
                type="button"
                disabled={offset === 0}
                onClick={() => setOffset(Math.max(0, offset - PAGE_SIZE))}
                className="h-8 border border-mathua-border px-3 font-mono text-[10px] uppercase tracking-wider text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue disabled:opacity-40"
              >
                Previous
              </button>
              <button
                type="button"
                disabled={offset + PAGE_SIZE >= total}
                onClick={() => setOffset(offset + PAGE_SIZE)}
                className="h-8 border border-mathua-border px-3 font-mono text-[10px] uppercase tracking-wider text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue disabled:opacity-40"
              >
                Next
              </button>
            </div>
          </div>
        )}

        <Link href="/admin" className="mt-10 inline-block font-mono text-xs text-mathua-secondary hover:text-mathua-blue">
          ← Overview
        </Link>
      </main>
      <Footer />
    </AppShell>
  )
}

// Snapshot pretty-prints one side of a change.
//
// The JSON is rendered as text, not parsed into fields, on purpose: the schema is the server's
// to change, and a page that renders whatever it is given keeps showing the record when the
// server starts sending one more field. A long single line also overflows on a phone.
function Snapshot({ label, json }: { label: string; json: string }) {
  let body = json
  try {
    body = JSON.stringify(JSON.parse(json), null, 2)
  } catch {
    // Not JSON, or truncated. Shown as-is rather than hidden: an unparseable record is itself
    // information, and silently rendering nothing would make it look like the row is empty.
  }
  return (
    <div className="border border-mathua-border min-w-0">
      <p className="px-2 py-1 font-mono text-[10px] uppercase tracking-wider text-mathua-muted border-b border-mathua-border">
        {label}
      </p>
      <pre className="p-2 font-mono text-[10px] text-mathua-secondary whitespace-pre-wrap break-all overflow-x-hidden">
        {body || ', '}
      </pre>
    </div>
  )
}