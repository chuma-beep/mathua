'use client'

import { useCallback, useEffect, useMemo, useState } from 'react'
import Link from 'next/link'
import AppShell from '../../../components/app-shell'
import Footer from '../../../components/Footer'
import SectionHeader from '../../../components/SectionHeader'
import AdminNav from '../../../components/AdminNav'
import Loading from '../../../components/Loading'
import KatexContent from '../../../components/KatexContent'
import ReportReplay from '../../../components/ReportReplay'
import {
  listAdminReports, moderateReport, bulkModerateReports, getErrorMessage, getAdminReportImage,
  type QuestionReport, type ModerationStatus,
} from '../../../lib/api'
import { useAdminMe } from '../../../hooks/useAdminMe'

// The moderation queue. Reported Content is a first-class surface, not an afterthought: a report
// is a person telling us something is wrong, and the queue exists so an operator can see it,
// judge it, and record the decision without leaving the deployed app.
//
// A report is never deleted. Resolving stamps who decided and why, and the decision is written to
// the audit trail; the row keeps the complaint so the next person can see what was decided before.

const FILTERS: { key: string; label: string }[] = [
  { key: 'open', label: 'Open' },
  { key: 'reviewing', label: 'Reviewing' },
  { key: 'resolved', label: 'Resolved' },
  { key: 'dismissed', label: 'Dismissed' },
  { key: 'all', label: 'All' },
]

const PAGE_SIZE = 20

export default function AdminReportsPage() {
  const { can } = useAdminMe()
  const [filter, setFilter] = useState('open')
  const [query, setQuery] = useState('')
  const [appliedQuery, setAppliedQuery] = useState('')
  const [offset, setOffset] = useState(0)
  const [total, setTotal] = useState(0)
  const [reports, setReports] = useState<QuestionReport[] | null>(null)
  const [error, setError] = useState('')

  // One bug produces many reports. Grouping the queue by concept is what makes
  // that visible, and what gives the bulk decision a subject: closing "these
  // eleven" is closing the bug, not eleven unrelated complaints. Order is the
  // first appearance in the server's newest-first list.
  const groups = useMemo(() => {
    const byConcept = new Map<string, QuestionReport[]>()
    for (const r of reports ?? []) {
      const key = r.concept_id || ''
      const bucket = byConcept.get(key)
      if (bucket) bucket.push(r)
      else byConcept.set(key, [r])
    }
    return [...byConcept.entries()]
  }, [reports])

  const load = useCallback(() => {
    setError('')
    listAdminReports(filter, appliedQuery, PAGE_SIZE, offset)
      .then((r) => { setReports(r.reports); setTotal(r.total ?? r.reports.length) })
      .catch((e) => { setError(getErrorMessage(e)); setReports([]); setTotal(0) })
  }, [filter, appliedQuery, offset])

  useEffect(() => { load() }, [load])

  return (
    <AppShell>
      <main className="mx-auto w-full max-w-[900px] min-w-0 px-4 sm:px-6 py-8 sm:py-12 overflow-x-hidden">
        <SectionHeader label="Administration" title="Reports" />
        <AdminNav />

        <form
          onSubmit={(event) => { event.preventDefault(); setAppliedQuery(query.trim()); setOffset(0) }}
          className="mb-4 flex flex-wrap gap-2"
        >
          <label htmlFor="report-search" className="sr-only">Search reports</label>
          <input
            id="report-search"
            value={query}
            onChange={(event) => setQuery(event.target.value)}
            placeholder="Search concept, question, reporter or resolution"
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

        <div className="mb-4 flex flex-wrap gap-2">
          {FILTERS.map((f) => (
            <button
              key={f.key}
              type="button"
              onClick={() => { setFilter(f.key); setOffset(0) }}
              className={
                filter === f.key
                  ? 'border border-mathua-blue text-mathua-blue px-3 h-9 font-mono text-xs'
                  : 'border border-mathua-border text-mathua-secondary hover:text-mathua-blue px-3 h-9 font-mono text-xs'
              }
            >
              {f.label}
            </button>
          ))}
        </div>

        {error && (
          <p role="alert" className="mb-4 border border-mathua-red bg-mathua-surface p-4 font-mono text-xs text-mathua-red">{error}</p>
        )}

        {!reports && !error && <Loading label="LOADING REPORTS" />}
        {reports && reports.length === 0 && !error && (
          <p role="status" className="border border-mathua-border bg-mathua-surface p-4 font-mono text-xs text-mathua-secondary">
            Nothing here.
          </p>
        )}

        <div className="space-y-6">
          {groups.map(([concept, items]) => {
            const unsettled = items.filter((r) => r.status === 'open' || r.status === 'reviewing')
            return (
              <section key={concept} className="space-y-3">
                <div className="flex flex-wrap items-baseline justify-between gap-2 border-b border-mathua-border pb-2">
                  <h2 className="font-mono text-xs text-mathua-primary">
                    {concept || (items.every((r) => r.kind === 'message') ? 'Messages' : '(no concept)')}
                  </h2>
                  <span className="font-mono text-[10px] text-mathua-muted">
                    {items.length} report{items.length === 1 ? '' : 's'}
                    {unsettled.length !== items.length ? `, ${unsettled.length} unsettled` : ''}
                  </span>
                </div>
                {can('reports.moderate') && concept && unsettled.length > 1 && (
                  <BulkClose conceptId={concept} count={unsettled.length} onDone={load} />
                )}
                <ul className="space-y-3">
                  {items.map((r) => (
                    <ReportCard key={r.id} report={r} canModerate={can('reports.moderate')} onChanged={load} />
                  ))}
                </ul>
              </section>
            )
          })}
        </div>

        {total > PAGE_SIZE && (
          <div className="mt-6 flex flex-wrap items-center justify-between gap-3 border-t border-mathua-border pt-4">
            <span className="font-mono text-[11px] text-mathua-muted">
              {offset + 1}–{Math.min(offset + (reports?.length ?? 0), total)} of {total}
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
      </main>
      <Footer />
    </AppShell>
  )
}

function StatusPill({ status }: { status: string }) {
  const tone =
    status === 'open' ? 'text-mathua-red border-mathua-red'
      : status === 'reviewing' ? 'text-mathua-blue border-mathua-blue'
        : 'text-mathua-muted border-mathua-border'
  return <span className={`border ${tone} px-2 py-0.5 font-mono text-[10px] uppercase tracking-wider`}>{status}</span>
}

// The image route is behind the admin token, and an <img> tag cannot send an
// Authorization header. Fetch the bytes with the normal API helper and hand the
// browser a short-lived object URL instead of weakening the route to public.
function ReportImage({ reportID }: { reportID: number }) {
  const [url, setURL] = useState('')
  const [failed, setFailed] = useState(false)

  useEffect(() => {
    let live = true
    let objectURL = ''
    setURL('')
    setFailed(false)
    getAdminReportImage(reportID)
      .then((blob) => {
        objectURL = URL.createObjectURL(blob)
        if (live) setURL(objectURL)
        else URL.revokeObjectURL(objectURL)
      })
      .catch(() => {
        if (live) setFailed(true)
      })
    return () => {
      live = false
      if (objectURL) URL.revokeObjectURL(objectURL)
    }
  }, [reportID])

  if (failed) {
    return <p className="mt-2 font-mono text-[11px] text-mathua-red">The screenshot could not be loaded.</p>
  }
  if (!url) {
    return <p className="mt-2 font-mono text-[11px] text-mathua-muted">Loading screenshot…</p>
  }
  return (
    <a href={url} target="_blank" rel="noreferrer" className="mt-2 inline-block">
      {/* eslint-disable-next-line @next/next/no-img-element */}
      <img
        src={url}
        alt={`Screenshot attached to report ${reportID}`}
        className="max-h-[32rem] max-w-full border border-mathua-border"
      />
      <span className="mt-1 block font-mono text-[10px] text-mathua-blue">Open full size →</span>
    </a>
  )
}

function ReportCard({ report, canModerate, onChanged }: { report: QuestionReport; canModerate: boolean; onChanged: () => void }) {
  const [open, setOpen] = useState(false)
  const [resolution, setResolution] = useState('')
  const [busy, setBusy] = useState(false)
  const [error, setError] = useState('')

  async function decide(status: ModerationStatus) {
    if ((status === 'resolved' || status === 'dismissed') && !resolution.trim()) {
      setError('A reason is required to resolve or dismiss.')
      return
    }
    setBusy(true)
    setError('')
    try {
      await moderateReport(report.id, status, resolution.trim())
      onChanged()
    } catch (e) {
      setError(getErrorMessage(e))
    } finally {
      setBusy(false)
    }
  }

  return (
    <li className="border border-mathua-border bg-mathua-surface">
      <button
        type="button"
        onClick={() => setOpen((o) => !o)}
        className="flex w-full items-center justify-between gap-3 p-4 text-left"
      >
        <span className="min-w-0">
          <span className="font-mono text-xs text-mathua-primary">{report.concept_id || '(no concept)'}</span>
          <span className="ml-2 font-mono text-[11px] text-mathua-muted">{report.kind} · {report.reason}</span>
          <span className="ml-2 font-mono text-[11px] text-mathua-blue">
            {report.reporter_username
              ? `@${report.reporter_username}`
              : report.reporter_name || report.reporter_id || '(anon)'}
          </span>
        </span>
        <span className="flex shrink-0 items-center gap-2">
          <StatusPill status={report.status} />
          <span className="font-mono text-[10px] text-mathua-muted">#{report.id}</span>
        </span>
      </button>

      {open && (
        <div className="border-t border-mathua-border p-4">
          <dl className="grid grid-cols-1 sm:grid-cols-2 gap-2 font-mono text-[11px] text-mathua-secondary">
            <div>
              <dt className="text-mathua-muted">Reported by</dt>
              <dd>
                {report.reporter_id ? (
                  <Link
                    href={`/admin/users/detail?id=${encodeURIComponent(report.reporter_id)}`}
                    className="break-all text-mathua-blue hover:underline"
                    title={`Open this account (${report.reporter_id})`}
                  >
                    {report.reporter_username
                      ? `@${report.reporter_username}`
                      : report.reporter_name || report.reporter_id}
                  </Link>
                ) : '(anon)'}
                {report.reporter_username && report.reporter_name && (
                  <span className="ml-1 text-mathua-muted">({report.reporter_name})</span>
                )}
              </dd>
            </div>
            <div><dt className="text-mathua-muted">When</dt><dd>{report.created_at}</dd></div>
            <div><dt className="text-mathua-muted">Source</dt><dd>{report.source || ', '}</dd></div>
            <div><dt className="text-mathua-muted">Lesson</dt><dd>{report.lesson_id || ', '}</dd></div>
          </dl>

          {report.question && (
            <div className="mt-3 border border-mathua-border bg-mathua-code p-3">
              <p className="mb-1 font-mono text-[10px] uppercase tracking-wider text-mathua-muted">Question</p>
              <KatexContent className="text-mathua-primary text-sm">{report.question}</KatexContent>
            </div>
          )}
          {report.expected && (
            <p className="mt-3 font-mono text-[11px] text-mathua-secondary">
              Reporter&apos;s expected answer: <span className="text-mathua-primary">{report.expected}</span>
            </p>
          )}
          {report.explanation && (
            <p className="mt-2 font-mono text-[11px] text-mathua-secondary">{report.explanation}</p>
          )}
          {report.detail && (
            <p className="mt-3 border-l-2 border-mathua-border pl-3 font-mono text-[11px] text-mathua-secondary">{report.detail}</p>
          )}
          {report.has_image && (
            <div className="mt-3">
              <p className="font-mono text-[10px] uppercase tracking-wider text-mathua-muted">Screenshot</p>
              <ReportImage reportID={report.id} />
            </div>
          )}

          <ReportReplay report={report} />

          {report.status === 'resolved' || report.status === 'dismissed' ? (
            <p className="mt-3 font-mono text-[11px] text-mathua-muted">
              {report.status} by {report.resolved_by || '(unknown)'} · {report.resolved_at}
              {report.resolution ? `, ${report.resolution}` : ''}
            </p>
          ) : null}

          {canModerate && (report.status === 'open' || report.status === 'reviewing') && (
            <div className="mt-4">
              <label htmlFor={`res-${report.id}`} className="mb-1 block font-mono text-[10px] uppercase tracking-wider text-mathua-muted">
                Decision reason
              </label>
              <textarea
                id={`res-${report.id}`}
                value={resolution}
                onChange={(e) => setResolution(e.target.value)}
                rows={2}
                className="w-full border border-mathua-border bg-mathua-code p-2 font-mono text-xs text-mathua-primary"
                placeholder="What did you decide, and why?"
              />
              {error && <p role="alert" className="mt-2 font-mono text-[11px] text-mathua-red">{error}</p>}
              <div className="mt-2 flex flex-wrap gap-2">
                {report.status === 'open' && (
                  <button type="button" disabled={busy} onClick={() => decide('reviewing')}
                    className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue-faint px-3 h-9 font-mono text-xs disabled:opacity-50">
                    Mark reviewing
                  </button>
                )}
                <button type="button" disabled={busy} onClick={() => decide('resolved')}
                  className="border border-mathua-green text-mathua-green hover:bg-mathua-green-faint px-3 h-9 font-mono text-xs disabled:opacity-50">
                  Resolve
                </button>
                <button type="button" disabled={busy} onClick={() => decide('dismissed')}
                  className="border border-mathua-border text-mathua-secondary hover:text-mathua-red px-3 h-9 font-mono text-xs disabled:opacity-50">
                  Dismiss
                </button>
              </div>
            </div>
          )}
        </div>
      )}
    </li>
  )
}

// Closing a concept's reports after the fix ships. The action names the concept
// rather than the ids on screen, so a report that arrived while the operator was
// reading is closed too, it is the same bug. The server still writes one
// decision record per report, and never revisits one that already has a verdict.
function BulkClose({ conceptId, count, onDone }: { conceptId: string; count: number; onDone: () => void }) {
  const [armed, setArmed] = useState<ModerationStatus | null>(null)
  const [resolution, setResolution] = useState('')
  const [busy, setBusy] = useState(false)
  const [error, setError] = useState('')
  const [notice, setNotice] = useState('')

  async function close(status: ModerationStatus) {
    if (!resolution.trim()) {
      setError('A reason is required to resolve or dismiss.')
      return
    }
    setBusy(true)
    setError('')
    setNotice('')
    try {
      const res = await bulkModerateReports({ concept_id: conceptId }, status, resolution.trim())
      const failed = res.failed.length ? ` ${res.failed.length} could not be updated (${res.failed.map((f) => '#' + f.id).join(', ')}).` : ''
      setNotice(`Closed ${res.updated} report${res.updated === 1 ? '' : 's'}${res.skipped ? `, ${res.skipped} already settled` : ''}.${failed}`)
      setArmed(null)
      setResolution('')
      onDone()
    } catch (e) {
      setError(getErrorMessage(e))
    } finally {
      setBusy(false)
    }
  }

  return (
    <div className="border border-mathua-border bg-mathua-code p-3">
      <div className="flex flex-wrap items-center justify-between gap-2">
        <p className="font-mono text-[11px] text-mathua-secondary">
          {count} unsettled report{count === 1 ? '' : 's'} for this concept. Closing here settles every open or reviewing one.
        </p>
        {!armed && (
          <div className="flex flex-wrap gap-2">
            <button type="button" onClick={() => setArmed('resolved')}
              className="border border-mathua-green text-mathua-green hover:bg-mathua-green-faint px-3 h-8 font-mono text-[11px]">
              Resolve all
            </button>
            <button type="button" onClick={() => setArmed('dismissed')}
              className="border border-mathua-border text-mathua-secondary hover:text-mathua-red px-3 h-8 font-mono text-[11px]">
              Dismiss all
            </button>
          </div>
        )}
      </div>

      {armed && (
        <div className="mt-3">
          <label htmlFor={`bulk-res-${conceptId}`} className="mb-1 block font-mono text-[10px] uppercase tracking-wider text-mathua-muted">
            Reason for closing {count} report{count === 1 ? '' : 's'} as {armed}
          </label>
          <textarea
            id={`bulk-res-${conceptId}`}
            value={resolution}
            onChange={(e) => setResolution(e.target.value)}
            rows={2}
            className="w-full border border-mathua-border bg-mathua-surface p-2 font-mono text-xs text-mathua-primary"
            placeholder="What changed, and where?"
          />
          <div className="mt-2 flex flex-wrap gap-2">
            <button type="button" disabled={busy} onClick={() => close(armed)}
              className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue-faint px-3 h-8 font-mono text-[11px] disabled:opacity-50">
              {busy ? 'Closing...' : `Confirm ${armed}`}
            </button>
            <button type="button" disabled={busy} onClick={() => { setArmed(null); setError('') }}
              className="border border-mathua-border text-mathua-secondary px-3 h-8 font-mono text-[11px] disabled:opacity-50">
              Cancel
            </button>
          </div>
        </div>
      )}

      {error && <p role="alert" className="mt-2 font-mono text-[11px] text-mathua-red">{error}</p>}
      {notice && <p role="status" className="mt-2 font-mono text-[11px] text-mathua-green">{notice}</p>}
    </div>
  )
}
