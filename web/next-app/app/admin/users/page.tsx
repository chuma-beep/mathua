'use client'

import AdminShell from '@/components/admin/AdminShell'
import { useCallback, useEffect, useState } from 'react'
import Link from 'next/link'
import SectionHeader from '../../../components/SectionHeader'
import ConfirmDialog from '../../../components/ConfirmDialog'
import Avatar from '../../../components/Avatar'
import { Skeleton } from '../../../components/ui/skeleton'
import { Table, TableBody, TableCaption, TableCell, TableHead, TableHeader, TableRow } from '../../../components/ui/table'
import { API_BASE, searchAdminUsers, setUserRole, getErrorMessage, type AdminUser, type Role } from '../../../lib/api'
import { dicebearUrl, isDicebearStyle } from '../../../lib/dicebear'
import { useAuthState } from '../../../hooks/useAuthState'
import { useAdminMe } from '../../../hooks/useAdminMe'

// Account administration, and the only place an administrator is made or un-made.
//
// The whole surface is: find an account, see its role, change the role. What is absent is the
// point. There is no editor for XP, mastery, attempts or concept progress, those are the
// learning system's records, not account attributes, and an administrator with a text box in
// front of them will eventually use it to mark a concept learned. One operation is offered,
// and it is the one an administrator is trusted to perform.
//
// Promotion goes through `PATCH /api/admin/users/{id}`, which is behind `requireAdmin` and
// delegates to one transactional primitive: read the before state, write the role, read the
// after state, append the audit event, commit. Nothing on this page decides anything, it asks.
//
// Both directions are behind a confirmation. See ConfirmDialog for why: an audited,
// immediately-effective privilege change should not be one click away in a list of rows, and the
// audit row is permanent even when the role is put back.

/** One row's pending change. Held as the whole target so the dialog can name the account. */
interface Pending {
  user: AdminUser
  next: Role
}

function labelFor(u: AdminUser): string {
  return u.name || u.username || u.email || u.id
}

function roleLabel(role: Role): string {
  return role === 'owner' ? 'owner' : role === 'admin' ? 'administrator' : role === 'moderator' ? 'moderator' : 'learner'
}

// The list is a scan, not a form. A joined date read as `2026-10-20` is precise and
// unreadable at a glance, so the row shows the day in the reader's locale and keeps
// the exact stamp in `dateTime`/`title`.
function fmtJoined(value: string): string {
  if (!value) return ''
  const d = new Date(value)
  if (Number.isNaN(d.getTime())) return value.slice(0, 10)
  return d.toLocaleDateString(undefined, { year: 'numeric', month: 'short', day: 'numeric' })
}

// The country arrives as an ISO code because that is what the edge sends. The
// list is read by people, so it is spelled out here; an unknown code is shown
// as-is rather than silently dropped.
function countryName(code: string): string {
  if (!code) return ''
  try {
    return new Intl.DisplayNames(undefined, { type: 'region' }).of(code) ?? code
  } catch {
    return code
  }
}

function roleTone(role: Role): string {
  switch (role) {
    case 'owner':
      return 'border-mathua-blue bg-mathua-blue-faint text-mathua-blue'
    case 'admin':
      return 'border-mathua-blue text-mathua-blue'
    case 'moderator':
      return 'border-mathua-green-faint text-mathua-green'
    default:
      return 'border-mathua-border text-mathua-muted'
  }
}

// The same precedence the sidebar, profile and account detail use: custom upload,
// Google photo, DiceBear pick, then the generated initial drawn by Avatar.
function avatarURL(u: AdminUser): string | undefined {
  const avatar = u.avatar
  if (!avatar) return undefined
  const stored = avatar.avatar_url
  if (stored) return stored.startsWith('/') ? `${API_BASE}${stored}` : stored
  const pick = avatar.avatar_dicebear
  return pick && isDicebearStyle(pick.style) ? dicebearUrl(pick.style, pick.seed) : undefined
}

function RoleBadge({ role }: { role: Role }) {
  return (
    <span className={`inline-flex items-center border px-1.5 py-0.5 font-mono text-[9px] uppercase tracking-wider ${roleTone(role)}`}>
      {roleLabel(role)}
    </span>
  )
}

export default function AdminUsersPage() {
  const { user } = useAuthState()
  const [query, setQuery] = useState('')
  const [users, setUsers] = useState<AdminUser[] | null>(null)
  // Two failures, because they are two different facts and one slot lost one of them. Refreshing
  // the list after a refused role change used to clear the shared error, so the message
  // explaining *why* the change did not happen was erased by the retry that confirmed nothing
  // changed and it vanished faster than a reader could see it.
  const [listError, setListError] = useState('')
  const [actionError, setActionError] = useState('')
  const [notice, setNotice] = useState('')
  const [pending, setPending] = useState<Pending | null>(null)
  const [busy, setBusy] = useState(false)

  const load = useCallback((q: string) => {
    searchAdminUsers(q)
      .then((res) => { setUsers(res.users); setListError('') })
      .catch((e) => setListError(getErrorMessage(e)))
  }, [])

  useEffect(() => { load('') }, [load])

  const { can, isStaff, me } = useAdminMe()
  const canManage = can('admins.manage')

  async function confirmChange() {
    if (!pending) return
    setBusy(true)
    setActionError('')
    setNotice('')
    try {
      const res = await setUserRole(pending.user.id, pending.next)
      // Re-read from the server rather than patching local state optimistically. The server may
      // have refused, the last administrator is the case that matters and this way the row
      // always shows what is actually stored.
      load(query)
      setNotice(
        res.changed
          ? `${labelFor(pending.user)} is now ${roleLabel(pending.next)}.`
          : `${labelFor(pending.user)} was already ${roleLabel(pending.next)}.`,
      )
      setPending(null)
    } catch (e) {
      // The dialog closes and the error is shown inline. Staying open would hide the message
      // behind the modal and leave the row's state unverified. The list is re-read because the
      // server may have applied part of the request, but that refresh writes to listError, not
      // here, so the refusal survives it.
      setActionError(getErrorMessage(e))
      setPending(null)
      load(query)
    } finally {
      setBusy(false)
    }
  }

  return (
    <AdminShell>
      <main className="mx-auto w-full max-w-[1200px] min-w-0 px-4 sm:px-6 py-8 sm:py-12 overflow-x-hidden">
        <SectionHeader label="Administration" title="Users" />
        <p className="-mt-3 mb-6 font-mono text-[11px] text-mathua-muted">
          Find an account, check its role and change it. Every change is recorded at{' '}
          <Link href="/admin/audit" className="text-mathua-secondary hover:text-mathua-blue hover:underline">/admin/audit</Link>.
        </p>

        {!isStaff && (
          <p role="status" className="border border-mathua-border bg-mathua-surface p-4 font-mono text-xs text-mathua-secondary">
            This area needs a staff role. If you hold one, sign in with that account.
          </p>
        )}

        <form
          className="mb-3 flex flex-wrap items-center gap-2"
          onSubmit={(e) => { e.preventDefault(); setActionError(''); load(query) }}
        >
          <label htmlFor="admin-user-search" className="sr-only">Search accounts</label>
          <input
            id="admin-user-search"
            type="search"
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder="Search by email, id, username or name…"
            className="h-11 w-full min-w-0 border border-mathua-border bg-mathua-surface px-3 font-mono text-xs text-mathua-primary placeholder:text-mathua-muted focus:border-mathua-blue focus:outline-none sm:h-10 sm:w-auto sm:flex-1"
          />
          <button type="submit" className="h-11 flex-1 border border-mathua-blue px-4 font-mono text-xs text-mathua-blue hover:bg-mathua-blue-faint sm:h-10 sm:flex-none">
            Search
          </button>
          {query && (
            <button
              type="button"
              onClick={() => { setQuery(''); setActionError(''); load('') }}
              className="h-11 flex-1 border border-mathua-border px-3 font-mono text-xs text-mathua-secondary hover:text-mathua-primary sm:h-10 sm:flex-none"
            >
              Clear
            </button>
          )}
        </form>

        {actionError && (
          <p role="alert" className="mb-4 border border-mathua-red bg-mathua-surface p-4 font-mono text-xs text-mathua-red">
            {actionError}
          </p>
        )}
        {listError && (
          <p role="alert" className="mb-4 border border-mathua-red bg-mathua-surface p-4 font-mono text-xs text-mathua-red">
            {listError}
          </p>
        )}
        {notice && (
          <p role="status" className="mb-4 border border-mathua-green-faint bg-mathua-surface p-4 font-mono text-xs text-mathua-green">
            {notice}
          </p>
        )}

        {/* Guarded on the error as well as the null. The other admin pages do this and this one
            did not, which meant a learner who navigated here saw an error box from the 403 and
            "Loading accounts…" underneath it, a page claiming to be busy while it had already
            given up, and would keep claiming indefinitely. */}
        {users === null && !listError && !actionError && (
          <>
            <p role="status" className="sr-only">Loading accounts…</p>
            <div className="border border-mathua-border" aria-hidden="true">
              {[0, 1, 2].map((i) => (
                <div key={i} className="flex items-center gap-3 border-b border-mathua-border p-3 last:border-b-0">
                  <Skeleton className="h-9 w-9 rounded-full" />
                  <div className="flex-1 space-y-2">
                    <Skeleton className="h-2.5 w-40" />
                    <Skeleton className="h-2 w-64 max-w-full" />
                  </div>
                </div>
              ))}
            </div>
          </>
        )}

        {users !== null && users.length === 0 && (
          <p className="font-mono text-xs text-mathua-muted">
            {query ? `No accounts match “${query}”.` : 'No accounts have been created yet.'}
          </p>
        )}

        {users !== null && users.length > 0 && (
          <div className="border border-mathua-border bg-mathua-surface">
            <Table className="block md:table" role="table">
              <TableCaption className="sr-only">Accounts, their roles, and the action available on each</TableCaption>
              <TableHeader className="hidden md:table-header-group" role="rowgroup">
                <TableRow className="hover:bg-transparent" role="row">
                  <TableHead className="px-4 py-2" role="columnheader">Account</TableHead>
                  <TableHead className="hidden px-4 py-2 md:table-cell" role="columnheader">Role</TableHead>
                  <TableHead className="px-4 py-2 text-right" role="columnheader">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody className="block md:table-row-group" role="rowgroup">
                {users.map((u) => {
                  const self = u.id === user?.student_id
                  return (
                    <TableRow key={u.id} className="block px-4 py-3 align-top md:table-row md:px-0 md:py-0" role="row">
                      <TableCell className="block p-0 md:table-cell md:px-4 md:py-3" role="cell">
                        <div className="flex items-start gap-3">
                          <Avatar
                            seed={u.id}
                            name={labelFor(u)}
                            size={38}
                            url={avatarURL(u)}
                            preset={u.avatar?.avatar_preset ?? null}
                          />
                          <div className="min-w-0">
                            <div className="flex flex-wrap items-center gap-2">
                              <Link
                                href={`/admin/users/detail?id=${encodeURIComponent(u.id)}`}
                                className="truncate font-mono text-xs text-mathua-primary hover:text-mathua-blue hover:underline"
                              >
                                {labelFor(u)}
                              </Link>
                              <span className="md:hidden"><RoleBadge role={u.role} /></span>
                              {self && (
                                <span className="border border-mathua-blue px-1.5 py-0.5 font-mono text-[9px] uppercase tracking-wider text-mathua-blue">
                                  you
                                </span>
                              )}
                            </div>
                            <div className="mt-0.5 truncate font-mono text-[11px] text-mathua-secondary" title={u.email}>
                              {u.email || 'no email'}
                            </div>
                            {/* The id is how an administrator names a person over chat, and the
                                provider flag answers "why can't they sign in with a password",
                                both are facts about the account, not credentials. */}
                            <div className="mt-0.5 truncate font-mono text-[10px] text-mathua-muted" title={u.id}>
                              @{u.username || 'no username'} · {u.id}
                            </div>
                            <div className="mt-1 flex flex-wrap items-center gap-x-2 gap-y-1 font-mono text-[10px] text-mathua-muted">
                              {u.created_at && (
                                <time dateTime={u.created_at} title={u.created_at}>
                                  Joined {fmtJoined(u.created_at)}
                                </time>
                              )}
                              {!u.has_password && (
                                <span title="Signed in with a provider, not a password">provider sign-in</span>
                              )}
                            </div>
                            {/* The same once-a-day capture that records the country records
                                this. Both are coarse on purpose: a date, a device class and a
                                country, never a time series, a user agent or an address. */}
                            {u.last_seen_at && (
                              <div className="mt-0.5 flex flex-wrap items-center gap-x-2 font-mono text-[10px] text-mathua-muted">
                                <time dateTime={u.last_seen_at} title={u.last_seen_at}>
                                  Last seen {fmtJoined(u.last_seen_at)}
                                </time>
                                {u.device && <span title="Device class from the sign-in request">{u.device}</span>}
                                {u.country && <span title="Country from the last signed-in visit">{countryName(u.country)}</span>}
                              </div>
                            )}
                          </div>
                        </div>
                      </TableCell>
                      <TableCell className="hidden p-0 md:table-cell md:px-4 md:py-3" role="cell">
                        <RoleBadge role={u.role} />
                      </TableCell>
                      <TableCell className="mt-3 block p-0 md:mt-0 md:table-cell md:px-4 md:py-3 md:text-right md:whitespace-nowrap" role="cell">
                        {canManage && !self ? (
                          <label className="sr-only" htmlFor={`role-${u.id}`}>Change role for {labelFor(u)}</label>
                        ) : null}
                        {canManage && !self ? (
                          <select
                            id={`role-${u.id}`}
                            value=""
                            onChange={(e) => {
                              const next = e.target.value as Role
                              if (!next) return
                              setActionError('')
                              setNotice('')
                              setPending({ user: u, next })
                            }}
                            className="mt-2 h-11 w-full border border-mathua-border bg-mathua-surface px-2 font-mono text-[11px] text-mathua-secondary md:mt-0 md:h-9 md:w-auto"
                          >
                            <option value="">Change role…</option>
                            {(me?.assignable_roles ?? []).filter((r) => r !== u.role).map((r) => (
                              <option key={r} value={r}>{roleLabel(r)}</option>
                            ))}
                          </select>
                        ) : (
                          <span className="font-mono text-[10px] text-mathua-muted">
                            {self ? 'your account' : '—'}
                          </span>
                        )}
                      </TableCell>
                    </TableRow>
                  )
                })}
              </TableBody>
            </Table>
          </div>
        )}

        {users !== null && users.length > 0 && (
          <p className="mt-3 font-mono text-[10px] text-mathua-muted">
            {users.length} {users.length === 1 ? 'account' : 'accounts'} shown
            {users.length >= 50 ? '. Refine the search to see more.' : '.'}
          </p>
        )}

        <p className="mt-6 font-mono text-[11px] text-mathua-muted">
          Learning state, XP, mastery, attempts, concept progress, is not editable here. It is
          the learning system&apos;s record, and every change is visible at /admin/audit.
        </p>

        <Link href="/admin" className="mt-10 inline-block font-mono text-xs text-mathua-secondary hover:text-mathua-blue">
          ← Overview
        </Link>
      </main>

      <ConfirmDialog
        open={pending !== null}
        title={pending ? `Change role to ${roleLabel(pending.next)}?` : 'Change role?'}
        confirmLabel={pending ? `Make ${roleLabel(pending.next)}` : 'Confirm'}
        tone={pending && pending.next === 'student' ? 'destructive' : 'default'}
        busy={busy}
        onCancel={() => setPending(null)}
        onConfirm={confirmChange}
      >
        {pending && (
          <>
            <p>
              {labelFor(pending.user)} will become {roleLabel(pending.next)}. The change takes
              effect on their next request, their existing session is re-read, so no re-login is
              needed.
            </p>
            <p className="text-mathua-muted">Recorded in the audit log.</p>
          </>
        )}
      </ConfirmDialog>
    </AdminShell>
  )
}
