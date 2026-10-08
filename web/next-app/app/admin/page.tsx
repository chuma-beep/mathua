'use client'

import { useEffect, useState } from 'react'
import Link from 'next/link'
import Header from '../../components/Header'
import BottomTabs from '../../components/BottomTabs'
import Footer from '../../components/Footer'
import SectionHeader from '../../components/SectionHeader'
import AdminNav from '../../components/AdminNav'
import Loading from '../../components/Loading'
import { adminOverview, getErrorMessage, type AdminOverview } from '../../lib/api'
import { useAdminMe } from '../../hooks/useAdminMe'

// The admin landing page answers one question: what needs attention right now.
//
// So the top of it is the moderation queue — outstanding reports — because that is the only
// figure here that is a task rather than a fact. The rest are counts an operator can act on and
// the recent trail, which is the answer to "what changed while I was away". Everything is read
// live on each visit; nothing is precomputed, because a number that must be refreshed is a number
// that will be stale in the screenshot someone argues from.
export default function AdminOverviewPage() {
  const { isStaff } = useAdminMe()
  const [data, setData] = useState<AdminOverview | null>(null)
  const [error, setError] = useState('')

  useEffect(() => {
    let cancelled = false
    adminOverview()
      .then((d) => { if (!cancelled) setData(d) })
      .catch((e) => { if (!cancelled) setError(getErrorMessage(e)) })
    return () => { cancelled = true }
  }, [])

  return (
    <>
      <Header />
      <div className="pt-[var(--chrome-top)] lg:pt-0">
        <main className="mx-auto w-full max-w-[900px] min-w-0 px-4 sm:px-6 py-8 sm:py-12 overflow-x-hidden">
          <SectionHeader label="Administration" title="Overview" />
          <AdminNav />

          {!isStaff && (
            <p role="status" className="border border-mathua-border bg-mathua-surface p-4 font-mono text-xs text-mathua-secondary">
              This area needs a staff role. If you hold one, sign in with that account.
            </p>
          )}

          {error && (
            <p role="alert" className="border border-mathua-red bg-mathua-surface p-4 font-mono text-xs text-mathua-red">
              {error}
            </p>
          )}

          {!data && !error && <Loading label="LOADING OVERVIEW" />}

          {data && (
            <>
              <Link
                href="/admin/reports?status=open"
                className="mb-8 flex items-center justify-between border border-mathua-border bg-mathua-surface p-4 hover:border-mathua-blue"
              >
                <span className="font-mono text-xs uppercase tracking-wider text-mathua-muted">Reports needing attention</span>
                <span className={`font-mono text-2xl tabular-nums ${data.reports.outstanding > 0 ? 'text-mathua-red' : 'text-mathua-primary'}`}>
                  {data.reports.outstanding}
                </span>
              </Link>

              <dl className="grid grid-cols-2 sm:grid-cols-3 gap-2">
                <Stat label="Open" value={data.reports.open} />
                <Stat label="Reviewing" value={data.reports.reviewing} />
                <Stat label="Resolved" value={data.reports.resolved} />
                <Stat label="Staff" value={data.staff} />
                <Stat label="Owners" value={data.owners} />
                <Stat label="Accounts" value={data.users} />
                <Stat label="Concepts" value={data.concepts} />
                <Stat label="Domains" value={data.domains} />
                <Stat label="Stored questions" value={data.questions} />
              </dl>

              {/* Zero owners is a broken deployment: no owner means nobody can grant owner, and
                  every administrative route is a 403 including the one that would fix it. The
                  recovery is the operator command, which is why it is named here. */}
              {data.owners === 0 && (
                <p role="alert" className="mt-4 border border-mathua-red bg-mathua-surface p-4 font-mono text-xs text-mathua-red">
                  No owners exist. Every administrative route answers 403 until one is promoted
                  from the machine: <code>mathua admin promote &lt;email&gt;</code>
                </p>
              )}

              <section className="mt-10">
                <div className="flex items-baseline justify-between gap-3">
                  <h2 className="font-serif text-lg text-mathua-primary">Recent admin activity</h2>
                  <Link href="/admin/audit" className="font-mono text-[11px] text-mathua-blue hover:underline">
                    Full log →
                  </Link>
                </div>
                {data.recent_activity.length === 0 ? (
                  <p className="mt-3 font-mono text-xs text-mathua-muted">
                    Nothing yet. Role changes and moderation decisions are recorded here.
                  </p>
                ) : (
                  <ul className="mt-3 space-y-1">
                    {data.recent_activity.map((e) => (
                      <li key={e.id} className="font-mono text-[11px] text-mathua-secondary">
                        <span className="text-mathua-muted">{e.created_at}</span>{' '}
                        <span className="text-mathua-primary">{e.action}</span>{' '}
                        {e.entity_type}#{e.entity_id}
                      </li>
                    ))}
                  </ul>
                )}
              </section>
            </>
          )}

          <Link href="/profile" className="mt-10 inline-block font-mono text-xs text-mathua-secondary hover:text-mathua-blue">
            ← Profile
          </Link>
        </main>
      </div>
      <BottomTabs />
      <Footer />
    </>
  )
}

function Stat({ label, value }: { label: string; value: number }) {
  return (
    <div className="border border-mathua-border bg-mathua-surface p-4 min-w-0">
      <dt className="font-mono text-[10px] uppercase tracking-wider text-mathua-muted">{label}</dt>
      <dd className="mt-1 font-mono text-2xl text-mathua-primary tabular-nums">{value.toLocaleString()}</dd>
    </div>
  )
}
