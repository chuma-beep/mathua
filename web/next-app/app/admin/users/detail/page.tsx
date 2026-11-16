'use client'

import AdminShell from '@/components/admin/AdminShell'
import { Suspense, useEffect, useState } from 'react'
import Link from 'next/link'
import { useSearchParams } from 'next/navigation'
import SectionHeader from '../../../../components/SectionHeader'
import Avatar from '../../../../components/Avatar'
import { dicebearUrl, isDicebearStyle } from '../../../../lib/dicebear'
import {
  API_BASE,
  getAdminUser,
  getErrorMessage,
  type AdminUser,
  type AdminUserAvatar,
  type AdminUserDetail,
  type Role,
} from '../../../../lib/api'
import { useAdminMe } from '../../../../hooks/useAdminMe'

// One account, as an administrator needs to see it while answering a report or a support
// message: who they are, what they have done in aggregate, what they sent us and what happened
// to their role.
//
// The address is `?id=`, not a path segment, because this site is a static export. A dynamic
// segment would need `generateStaticParams`, which can only list ids that existed at build time,
// so every account created afterwards would 404.
//
// The page shows counts, never answers. `AdminUserReportSchema` in lib/api.ts omits the
// learner's answer and the expected answer for the same reason the schema there omits a
// password hash: a field the client does not name cannot reach the page.

function roleLabel(role: Role): string {
  return role === 'owner' ? 'owner'
    : role === 'admin' ? 'administrator'
    : role === 'moderator' ? 'moderator'
    : 'learner'
}

function labelFor(user: { name?: string; username?: string; email?: string; id: string }): string {
  return user.name || user.username || user.email || user.id
}

// The same precedence the sidebar and profile use: custom upload, Google
// photo, DiceBear pick, then the generated initial. The server resolves the
// custom-upload path because only it knows the account's photo version.
function UserAvatar({ user, avatar }: { user: AdminUser; avatar: AdminUserAvatar }) {
  const pick = avatar.avatar_dicebear
  const storedURL = avatar.avatar_url
  const url = storedURL
    ? (storedURL.startsWith('/') ? `${API_BASE}${storedURL}` : storedURL)
    : (pick && isDicebearStyle(pick.style) ? dicebearUrl(pick.style, pick.seed) : undefined)
  return (
    <Avatar
      seed={user.id}
      name={labelFor(user)}
      size={56}
      url={url}
      preset={avatar.avatar_preset ?? null}
    />
  )
}

function Count({ label, value }: { label: string; value: number | string }) {
  return (
    <div className="border border-mathua-border bg-mathua-surface p-3">
      <dt className="font-mono text-[10px] uppercase tracking-wider text-mathua-muted">{label}</dt>
      <dd className="mt-1 font-mono text-sm text-mathua-primary">{value}</dd>
    </div>
  )
}

function AccountDetail() {
  const params = useSearchParams()
  const id = params.get('id') ?? ''
  const { isStaff } = useAdminMe()
  const [data, setData] = useState<AdminUserDetail | null>(null)
  const [error, setError] = useState('')

  useEffect(() => {
    if (!id) {
      setError('Open this page from a row on the users list or from a report.')
      return
    }
    let live = true
    getAdminUser(id)
      .then((res) => { if (live) { setData(res); setError('') } })
      .catch((e) => { if (live) setError(getErrorMessage(e)) })
    return () => { live = false }
  }, [id])

  const activity = data?.activity

  return (
    <AdminShell>
      <main className="mx-auto w-full max-w-[1200px] min-w-0 px-4 sm:px-6 py-8 sm:py-12 overflow-x-hidden">
        <SectionHeader label="Administration" title="Account" />
  
        <p className="mb-4">
          <Link href="/admin/users" className="font-mono text-xs text-mathua-blue hover:underline">
            ← All accounts
          </Link>
        </p>

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

        {!data && !error && (
          <p className="font-mono text-xs text-mathua-muted">Loading account…</p>
        )}

        {data && (
          <>
            <section className="border border-mathua-border bg-mathua-surface p-4">
              <div className="flex items-center gap-3">
                <UserAvatar user={data.user} avatar={data.avatar} />
                <h2 className="font-mono text-sm text-mathua-primary">{labelFor(data.user)}</h2>
              </div>
              <dl className="mt-3 grid grid-cols-1 sm:grid-cols-2 gap-2 font-mono text-[11px] text-mathua-secondary">
                <div><dt className="text-mathua-muted">Account id</dt><dd className="break-all" title={data.user.id}>{data.user.id}</dd></div>
                <div><dt className="text-mathua-muted">Role</dt><dd className={data.user.role === 'student' ? 'text-mathua-muted' : 'text-mathua-blue'}>{roleLabel(data.user.role)}</dd></div>
                <div><dt className="text-mathua-muted">Email</dt><dd className="break-all">{data.user.email || ', '}</dd></div>
                <div><dt className="text-mathua-muted">Username</dt><dd>{data.user.username ? `@${data.user.username}` : ', '}</dd></div>
                <div><dt className="text-mathua-muted">Created</dt><dd>{data.user.created_at || ', '}</dd></div>
                <div>
                  <dt className="text-mathua-muted">Sign-in</dt>
                  <dd>{data.user.has_password ? 'Password' : 'OAuth provider, no password'}</dd>
                </div>
              </dl>
              {!data.user.email && !data.user.has_password && (
                <p className="mt-3 font-mono text-[11px] text-mathua-muted">
                  No address and no password: this is an unclaimed guest. Their reports and attempts
                  are real, but the account cannot be signed into or recovered.
                </p>
              )}
            </section>

            {activity && (
              <section className="mt-4">
                <h2 className="mb-2 font-mono text-xs uppercase tracking-wider text-mathua-muted">Activity</h2>
                <dl className="grid grid-cols-2 sm:grid-cols-4 gap-2">
                  <Count label="XP total" value={activity.xp_total} />
                  <Count label="XP today" value={activity.xp_today} />
                  <Count label="Concepts touched" value={activity.concepts_touched} />
                  <Count label="Concepts mastered" value={activity.concepts_mastered} />
                  <Count label="Attempts" value={activity.attempts} />
                  <Count label="Reports filed" value={activity.reports_filed} />
                  <Count label="Last active day" value={activity.last_active_day || ', '} />
                  <Count label="Last attempt" value={activity.last_attempt_at ? activity.last_attempt_at.slice(0, 10) : ', '} />
                </dl>
                <p className="mt-2 font-mono text-[10px] text-mathua-muted">
                  Diagnostic {activity.diagnostic_completed_at ? `completed ${activity.diagnostic_completed_at.slice(0, 10)}` : 'not completed'}.
                  Counts only: answers and attempt rows are not shown here.
                </p>
              </section>
            )}

            <section className="mt-6">
              <h2 className="mb-2 font-mono text-xs uppercase tracking-wider text-mathua-muted">
                Reports from this account ({activity?.reports_filed ?? 0})
              </h2>
              {data.reports.length === 0 ? (
                <p className="font-mono text-xs text-mathua-muted">They have not reported anything.</p>
              ) : (
                <ul className="border border-mathua-border">
                  {data.reports.map((r) => (
                    <li key={r.id} className="border-b border-mathua-border last:border-b-0 p-3">
                      <div className="flex flex-wrap items-baseline justify-between gap-2">
                        <span className="font-mono text-[11px] text-mathua-primary">
                          {r.concept_id || '(no concept)'}
                          <span className="ml-2 text-mathua-muted">{r.kind} · {r.reason}</span>
                        </span>
                        <span className="font-mono text-[10px] text-mathua-muted">#{r.id} · {r.status} · {r.created_at.slice(0, 10)}</span>
                      </div>
                      {r.detail && (
                        <p className="mt-1 font-mono text-[11px] text-mathua-secondary whitespace-pre-wrap break-words">{r.detail}</p>
                      )}
                    </li>
                  ))}
                </ul>
              )}
            </section>

            <section className="mt-6">
              <h2 className="mb-2 font-mono text-xs uppercase tracking-wider text-mathua-muted">Role history</h2>
              {data.history.length === 0 ? (
                <p className="font-mono text-xs text-mathua-muted">No role changes recorded.</p>
              ) : (
                <ul className="border border-mathua-border">
                  {data.history.map((h) => (
                    <li key={h.id} className="border-b border-mathua-border last:border-b-0 p-3 font-mono text-[11px] text-mathua-secondary">
                      <span className="text-mathua-primary">{h.action}</span>
                      {' by '}
                      <span className="break-all">{h.actor_id || '(system)'}</span>
                      <span className="text-mathua-muted"> · {h.created_at.slice(0, 10)}</span>
                    </li>
                  ))}
                </ul>
              )}
            </section>
          </>
        )}
      </main>
    </AdminShell>
  )
}

export default function AdminUserDetailPage() {
  return (
    <Suspense fallback={<div style={{ background: 'var(--bg)', minHeight: '100vh' }} />}>
      <AccountDetail />
    </Suspense>
  )
}
