'use client'

import { useEffect, useState } from 'react'
import AppShell from '../../../components/app-shell'
import Footer from '../../../components/Footer'
import SectionHeader from '../../../components/SectionHeader'
import AdminNav from '../../../components/AdminNav'
import Loading from '../../../components/Loading'
import { getAdminNotificationHealth, getErrorMessage, type AdminNotificationProblem, type AdminNotificationStat } from '../../../lib/api'

// Notification delivery health.
//
// Read-only for the same reason the audit log is: the useful question is "is
// anything stuck", and every repair is a provider or a preference, not a row an
// operator should edit here. Counts come from one grouped query; the problem
// list is capped server-side so a bad provider cannot make this page heavy.
function statusLabel(status: string): string {
  return status.charAt(0).toUpperCase() + status.slice(1)
}

export default function AdminNotificationsPage() {
  const [stats, setStats] = useState<AdminNotificationStat[] | null>(null)
  const [problems, setProblems] = useState<AdminNotificationProblem[]>([])
  const [error, setError] = useState('')

  useEffect(() => {
    let cancelled = false
    getAdminNotificationHealth()
      .then((r) => {
        if (cancelled) return
        setStats(r.stats)
        setProblems(r.problems)
      })
      .catch((e) => {
        if (cancelled) return
        setError(getErrorMessage(e))
        setStats([])
      })
    return () => { cancelled = true }
  }, [])

  return (
    <AppShell>
      <main className="mx-auto w-full max-w-[900px] min-w-0 px-4 sm:px-6 py-8 sm:py-12 overflow-x-hidden">
        <SectionHeader label="Administration" title="Notifications" />
        <AdminNav />

        {error && <p role="alert" className="mb-4 border border-mathua-red bg-mathua-surface p-4 font-mono text-xs text-mathua-red">{error}</p>}
        {!stats && !error && <Loading label="LOADING DELIVERY HEALTH" />}

        {stats && (
          <>
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
                {problems.map((problem) => (
                  <li key={`${problem.notification_id}-${problem.channel}`} className="border border-mathua-border bg-mathua-surface px-4 py-3">
                    <div className="flex flex-wrap items-baseline justify-between gap-3">
                      <span className="font-mono text-xs text-mathua-primary">{problem.title || problem.notification_id}</span>
                      <span className="font-mono text-[10px] text-mathua-muted">{problem.channel} · {statusLabel(problem.status)} · {problem.attempts} attempt{problem.attempts === 1 ? '' : 's'}</span>
                    </div>
                    <div className="mt-1 font-mono text-[10px] text-mathua-secondary">
                      {problem.student_name || problem.student_id}
                      {problem.suppressed_reason ? ` · ${problem.suppressed_reason}` : ''}
                      {problem.last_error ? ` · ${problem.last_error}` : ''}
                    </div>
                  </li>
                ))}
              </ul>
            )}
          </>
        )}
      </main>
      <Footer />
    </AppShell>
  )
}
