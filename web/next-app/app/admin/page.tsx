'use client'

import { useEffect, useState } from 'react'
import Link from 'next/link'
import Header from '../../components/Header'
import BottomTabs from '../../components/BottomTabs'
import Footer from '../../components/Footer'
import SectionHeader from '../../components/SectionHeader'
import AdminNav from '../../components/AdminNav'
import Loading from '../../components/Loading'
import { adminOverview, type AdminOverview } from '../../lib/api'
import { getErrorMessage } from '../../lib/api'
import { useAuthState } from '../../hooks/useAuthState'

// A landing page for administration.
//
// Every figure is read live from /api/admin/overview on each visit, and the five of them are
// the whole thing: accounts, administrators, concepts, domains and stored questions. That is not
// a placeholder list — it is the set of numbers an operator can act on without a separate tool,
// and a dashboard is primarily a place to navigate from rather than a place to read.
//
// No charts, no retention curves, no cohort comparison, and nothing about how learners are
// doing. Learning performance is analytics, it belongs on the data model rather than bolted onto
// an authorization layer, and the audit trail below is the more useful thing to have here
// anyway: it says who changed what, which is the question an administrator actually arrives
// with.
export default function AdminOverviewPage() {
  const { user } = useAuthState()
  const [data, setData] = useState<AdminOverview | null>(null)
  const [error, setError] = useState('')

  useEffect(() => {
    let cancelled = false
    adminOverview()
      .then((d) => { if (!cancelled) setData(d) })
      .catch((e) => { if (!cancelled) setError(getErrorMessage(e)) })
    return () => { cancelled = true }
  }, [])

  const isAdmin = user?.role === 'admin'

  return (
    <>
      <Header />
      <div className="pt-[var(--chrome-top)] lg:pt-0">
        <main className="mx-auto w-full max-w-[820px] min-w-0 px-4 sm:px-6 py-8 sm:py-12 overflow-x-hidden">
          <SectionHeader label="Administration" title="Overview" />
          <AdminNav isAdmin={isAdmin} />

          {/* The page renders for anyone who navigates here, and the refusal comes from the
              API. Saying so here keeps the two honest: the nav is a convenience, and a learner
              who typed the URL sees this sentence rather than a page that looks broken. */}
          {!isAdmin && (
            <p role="status" className="border border-mathua-border bg-mathua-surface p-4 font-mono text-xs text-mathua-secondary">
              This area needs an administrator role. If you are one, sign in with that account.
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
              <dl className="grid grid-cols-2 sm:grid-cols-3 gap-2">
                <Stat label="Accounts" value={data.users} />
                <Stat label="Administrators" value={data.admins} />
                <Stat label="Concepts" value={data.concepts} />
                <Stat label="Domains" value={data.domains} />
                <Stat label="Stored questions" value={data.questions} />
              </dl>

              {/* An administrator count of zero is a broken deployment, and it is worth saying
                  out loud rather than leaving as a number. The server refuses every
                  administrative route at that point, including the one that would fix it, so the
                  recovery is the operator command — which is why it is named here. */}
              {data.admins === 0 && (
                <p role="alert" className="mt-4 border border-mathua-red bg-mathua-surface p-4 font-mono text-xs text-mathua-red">
                  No administrators exist. Every administrative route will answer 403 until one is
                  promoted from the machine: <code>mathua admin promote &lt;email&gt;</code>
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
                    Nothing yet. Role changes are recorded here.
                  </p>
                ) : (
                  <ul className="mt-3 space-y-1">
                    {data.recent_activity.map((e) => (
                      <li key={e.id} className="font-mono text-[11px] text-mathua-secondary">
                        <span className="text-mathua-muted">{e.created_at}</span>{' '}
                        <span className="text-mathua-primary">{e.action}</span>{' '}
                        {e.entity_id}
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
      {/* tabular-nums: a column of counts is a column to compare down, and proportional figures
          make the third digit shift. */}
      <dd className="mt-1 font-mono text-2xl text-mathua-primary tabular-nums">{value.toLocaleString()}</dd>
    </div>
  )
}