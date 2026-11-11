'use client'

import { useCallback, useEffect, useState } from 'react'
import AppShell from '../../../components/app-shell'
import Footer from '../../../components/Footer'
import SectionHeader from '../../../components/SectionHeader'
import AdminNav from '../../../components/AdminNav'
import Loading from '../../../components/Loading'
import {
  getAdminNotificationHealth,
  getAdminNotificationHistory,
  getErrorMessage,
  replayNotificationDelivery,
  type AdminNotificationDeliveryDetail,
  type AdminNotificationProblem,
  type AdminNotificationRetention,
  type AdminNotificationStat,
} from '../../../lib/api'

// Notification delivery health.
//
// Counts come from one grouped query and the problem list is capped server-side,
// so a bad provider cannot make this page heavy. Replay is the only write: it
// returns a failed or suppressed delivery to the queue. History answers the two
// questions an operator actually asks, "why did this one not arrive" and "what
// has this learner been sent".
function statusLabel(status: string): string {
  return status.charAt(0).toUpperCase() + status.slice(1)
}

export default function AdminNotificationsPage() {
  const [stats, setStats] = useState<AdminNotificationStat[] | null>(null)
  const [problems, setProblems] = useState<AdminNotificationProblem[]>([])
  const [retention, setRetention] = useState<AdminNotificationRetention | null>(null)
  const [historyFor, setHistoryFor] = useState('')
  const [history, setHistory] = useState<AdminNotificationDeliveryDetail[]>([])
  const [busy, setBusy] = useState('')
  const [error, setError] = useState('')

  const load = useCallback(() => {
    return getAdminNotificationHealth()
      .then((r) => {
        setStats(r.stats)
        setProblems(r.problems)
        setRetention(r.retention ?? null)
      })
      .catch((e) => {
        setError(getErrorMessage(e))
        setStats([])
      })
  }, [])

  useEffect(() => {
    void load()
  }, [load])

  const replay = (problem: AdminNotificationProblem) => {
    setBusy(`${problem.notification_id}-${problem.channel}`)
    setError('')
    replayNotificationDelivery(problem.notification_id, problem.channel)
      .then(() => load())
      .catch((e) => setError(getErrorMessage(e)))
      .finally(() => setBusy(''))
  }

  const showHistory = (problem: AdminNotificationProblem) => {
    if (historyFor === problem.notification_id) {
      setHistoryFor('')
      setHistory([])
      return
    }
    setBusy(`${problem.notification_id}-history`)
    setError('')
    getAdminNotificationHistory({ notificationId: problem.notification_id })
      .then((rows) => {
        setHistoryFor(problem.notification_id)
        setHistory(rows)
      })
      .catch((e) => setError(getErrorMessage(e)))
      .finally(() => setBusy(''))
  }

  return (
    <AppShell>
      <main className="mx-auto w-full max-w-[900px] min-w-0 px-4 sm:px-6 py-8 sm:py-12 overflow-x-hidden">
        <SectionHeader label="Administration" title="Notifications" />
        <AdminNav />

        {error && <p role="alert" className="mb-4 border border-mathua-red bg-mathua-surface p-4 font-mono text-xs text-mathua-red">{error}</p>}
        {!stats && !error && <Loading label="LOADING DELIVERY HEALTH" />}

        {stats && (
          <>
            {retention && (
              <p className="mb-4 font-mono text-[11px] text-mathua-muted">
                Retention: notifications and deliveries are pruned after {retention.days} days. A push device
                that fails {retention.push_failure_limit} times is removed after {retention.push_failure_days} days.
              </p>
            )}

            <h2 className="mb-2 font-mono text-xs text-mathua-primary">Deliveries by channel and status</h2>
            {stats.length === 0 ? (
              <p className="font-mono text-[11px] text-mathua-muted">No deliveries have been recorded yet.</p>
            ) : (
              <ul className="mb-8 space-y-1">
                {stats.map((stat) => (
                  <li key={`${stat.channel}-${stat.status}`} className="flex items-baseline justify-between gap-3 border border-mathua-border bg-mathua-surface px-4 py-2">
                    <span className="font-mono text-xs text-mathua-secondary">{stat.channel}</span>
                    <span className="font-mono text-xs text-mathua-primary">{statusLabel(stat.status)}</span>
                    <span className="font-mono text-xs text-mathua-primary">{stat.count}</span>
                  </li>
                ))}
              </ul>
            )}

            <h2 className="mb-2 font-mono text-xs text-mathua-primary">Failed and suppressed deliveries</h2>
            {problems.length === 0 ? (
              <p className="font-mono text-[11px] text-mathua-muted">Nothing is stuck. Retrying deliveries are not listed here.</p>
            ) : (
              <ul className="space-y-1">
                {problems.map((problem) => {
                  const key = `${problem.notification_id}-${problem.channel}`
                  const open = historyFor === problem.notification_id
                  return (
                    <li key={key} className="border border-mathua-border bg-mathua-surface px-4 py-3">
                      <div className="flex flex-wrap items-baseline justify-between gap-3">
                        <span className="font-mono text-xs text-mathua-primary">{problem.title || problem.notification_id}</span>
                        <span className="font-mono text-[10px] text-mathua-muted">{problem.channel} · {statusLabel(problem.status)} · {problem.attempts} attempt{problem.attempts === 1 ? '' : 's'}</span>
                      </div>
                      <div className="mt-1 font-mono text-[10px] text-mathua-secondary">
                        {problem.student_name || problem.student_id}
                        {problem.suppressed_reason ? ` · ${problem.suppressed_reason}` : ''}
                        {problem.last_error ? ` · ${problem.last_error}` : ''}
                      </div>
                      <div className="mt-2 flex flex-wrap gap-2">
                        <button
                          type="button"
                          disabled={busy === key}
                          onClick={() => replay(problem)}
                          className="border border-mathua-border px-2 h-7 font-mono text-[10px] uppercase tracking-wider text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue disabled:opacity-50"
                        >
                          {busy === key ? 'Replaying' : 'Replay'}
                        </button>
                        <button
                          type="button"
                          disabled={busy === `${problem.notification_id}-history`}
                          onClick={() => showHistory(problem)}
                          className="border border-mathua-border px-2 h-7 font-mono text-[10px] uppercase tracking-wider text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue disabled:opacity-50"
                        >
                          {open ? 'Hide history' : 'History'}
                        </button>
                      </div>
                      {open && (
                        <ul className="mt-2 space-y-1 border-t border-mathua-border pt-2">
                          {history.length === 0 && (
                            <li className="font-mono text-[10px] text-mathua-muted">No delivery rows for this notification.</li>
                          )}
                          {history.map((row) => (
                            <li key={`${row.notification_id}-${row.channel}`} className="font-mono text-[10px] text-mathua-secondary">
                              {row.channel} · {statusLabel(row.status)} · {row.attempts} attempt{row.attempts === 1 ? '' : 's'}
                              {row.provider_id ? ` · ${row.provider_id}` : ''}
                              {row.sent_at ? ` · sent ${row.sent_at}` : ''}
                              {row.opened_at ? ` · opened ${row.opened_at}` : ''}
                              {row.clicked_at ? ` · clicked ${row.clicked_at}` : ''}
                              {row.next_attempt_at ? ` · next ${row.next_attempt_at}` : ''}
                            </li>
                          ))}
                        </ul>
                      )}
                    </li>
                  )
                })}
              </ul>
            )}
          </>
        )}
      </main>
      <Footer />
    </AppShell>
  )
}
