'use client'

import { useCallback, useEffect, useState } from 'react'
import Link from 'next/link'
import Header from '../../../components/Header'
import BottomTabs from '../../../components/BottomTabs'
import Footer from '../../../components/Footer'
import SectionHeader from '../../../components/SectionHeader'
import AdminNav from '../../../components/AdminNav'
import ConfirmDialog from '../../../components/ConfirmDialog'
import { searchAdminUsers, setUserRole, type AdminUser, type Role } from '../../../lib/api'
import { getErrorMessage } from '../../../lib/api'
import { useAuthState } from '../../../hooks/useAuthState'

// Account administration, and the only place an administrator is made or un-made.
//
// The whole surface is: find an account, see its role, change the role. What is absent is the
// point. There is no editor for XP, mastery, attempts or concept progress — those are the
// learning system's records, not account attributes, and an administrator with a text box in
// front of them will eventually use it to mark a concept learned. One operation is offered,
// and it is the one an administrator is trusted to perform.
//
// Promotion goes through `PATCH /api/admin/users/{id}`, which is behind `requireAdmin` and
// delegates to one transactional primitive: read the before state, write the role, read the
// after state, append the audit event, commit. Nothing on this page decides anything — it asks.
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

export default function AdminUsersPage() {
  const { user } = useAuthState()
  const [query, setQuery] = useState('')
  const [users, setUsers] = useState<AdminUser[] | null>(null)
  // Two failures, because they are two different facts and one slot lost one of them. Refreshing
  // the list after a refused role change used to clear the shared error, so the message
  // explaining *why* the change did not happen was erased by the retry that confirmed nothing
  // changed — and it vanished faster than a reader could see it.
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

  const isAdmin = user?.role === 'admin'

  async function confirmChange() {
    if (!pending) return
    setBusy(true)
    setActionError('')
    setNotice('')
    try {
      const res = await setUserRole(pending.user.id, pending.next)
      // Re-read from the server rather than patching local state optimistically. The server may
      // have refused — the last administrator is the case that matters — and this way the row
      // always shows what is actually stored.
      load(query)
      setNotice(
        res.changed
          ? `${labelFor(pending.user)} is now ${pending.next === 'admin' ? 'an administrator' : 'a learner'}.`
          : `${labelFor(pending.user)} was already ${pending.next === 'admin' ? 'an administrator' : 'a learner'}.`,
      )
      setPending(null)
    } catch (e) {
      // The dialog closes and the error is shown inline. Staying open would hide the message
      // behind the modal and leave the row's state unverified. The list is re-read because the
      // server may have applied part of the request — but that refresh writes to listError, not
      // here, so the refusal survives it.
      setActionError(getErrorMessage(e))
      setPending(null)
      load(query)
    } finally {
      setBusy(false)
    }
  }

  return (
    <>
      <Header />
      <div className="pt-[var(--chrome-top)] lg:pt-0">
        <main className="mx-auto w-full max-w-[900px] min-w-0 px-4 sm:px-6 py-8 sm:py-12 overflow-x-hidden">
          <SectionHeader label="Administration" title="Users" />
          <AdminNav isAdmin={isAdmin} />

          {!isAdmin && (
            <p role="status" className="border border-mathua-border bg-mathua-surface p-4 font-mono text-xs text-mathua-secondary">
              This area needs an administrator role. If you are one, sign in with that account.
            </p>
          )}

          <form
            className="mb-4 flex gap-2"
            onSubmit={(e) => { e.preventDefault(); setActionError(''); load(query) }}
          >
            <label htmlFor="admin-user-search" className="sr-only">Search accounts</label>
            <input
              id="admin-user-search"
              type="search"
              value={query}
              onChange={(e) => setQuery(e.target.value)}
              placeholder="Search by email, id, username or name…"
              className="flex-1 min-w-0 border border-mathua-border bg-mathua-surface px-3 h-11 font-mono text-xs text-mathua-primary placeholder:text-mathua-muted focus:border-mathua-blue focus:outline-none"
            />
            <button type="submit" className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue-faint px-4 h-11 font-mono text-xs">
              Search
            </button>
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
              "Loading accounts…" underneath it — a page claiming to be busy while it had already
              given up, and would keep claiming indefinitely. */}
          {users === null && !listError && !actionError && (
            <p className="font-mono text-xs text-mathua-muted">Loading accounts…</p>
          )}

          {users !== null && users.length === 0 && (
            <p className="font-mono text-xs text-mathua-muted">No accounts match that search.</p>
          )}

          {users !== null && users.length > 0 && (
            <div className="border border-mathua-border overflow-x-hidden">
              <table className="w-full font-mono text-xs">
                <caption className="sr-only">Accounts, their roles, and the action available on each</caption>
                <thead>
                  <tr className="border-b border-mathua-border text-mathua-muted text-left">
                    <th scope="col" className="px-3 py-2 font-normal">Name</th>
                    <th scope="col" className="px-3 py-2 font-normal">Email</th>
                    <th scope="col" className="px-3 py-2 font-normal">Role</th>
                    <th scope="col" className="px-3 py-2 font-normal text-right">Actions</th>
                  </tr>
                </thead>
                <tbody>
                  {users.map((u) => {
                    const isAdminRow = u.role === 'admin'
                    const self = u.id === user?.student_id
                    return (
                      <tr key={u.id} className="border-b border-mathua-border last:border-b-0 align-top">
                        <td className="px-3 py-2 min-w-0">
                          <div className="text-mathua-primary truncate">{labelFor(u)}</div>
                          {/* The id is how an administrator names a person over chat, and the
                              provider flag answers "why can't they sign in with a password" —
                              both are facts about the account, not credentials. */}
                          <div className="text-mathua-muted truncate" title={u.id}>{u.id}</div>
                          {u.username && <div className="text-mathua-muted truncate">@{u.username}</div>}
                          {u.created_at && (
                            <div className="text-mathua-muted">Joined {u.created_at.slice(0, 10)}</div>
                          )}
                          {!u.has_password && (
                            <div className="text-mathua-muted">Signed in with a provider, not a password</div>
                          )}
                        </td>
                        <td className="px-3 py-2 min-w-0 text-mathua-secondary break-all">{u.email || '—'}</td>
                        <td className="px-3 py-2 whitespace-nowrap">
                          {isAdminRow ? (
                            <span className="text-mathua-blue">admin{self ? ' (you)' : ''}</span>
                          ) : (
                            <span className="text-mathua-muted">student</span>
                          )}
                        </td>
                        <td className="px-3 py-2 text-right whitespace-nowrap">
                          <button
                            type="button"
                            onClick={() => {
                              setActionError('')
                              setNotice('')
                              setPending({ user: u, next: isAdminRow ? 'student' : 'admin' })
                            }}
                            aria-haspopup="dialog"
                            className="border border-mathua-border px-3 h-11 font-mono text-[11px] text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue"
                          >
                            {isAdminRow ? 'Remove administrator' : 'Make administrator'}
                          </button>
                        </td>
                      </tr>
                    )
                  })}
                </tbody>
              </table>
            </div>
          )}

          <p className="mt-4 font-mono text-[11px] text-mathua-muted">
            Learning state — XP, mastery, attempts, concept progress — is not editable here. It is
            the learning system&apos;s record, and every change is visible at /admin/audit.
          </p>

          <Link href="/admin" className="mt-10 inline-block font-mono text-xs text-mathua-secondary hover:text-mathua-blue">
            ← Overview
          </Link>
        </main>
      </div>
      <BottomTabs />
      <Footer />

      <ConfirmDialog
        open={pending !== null}
        title={pending?.next === 'admin' ? 'Make administrator?' : 'Remove administrator access?'}
        confirmLabel={pending?.next === 'admin' ? 'Make administrator' : 'Remove access'}
        tone={pending?.next === 'admin' ? 'default' : 'destructive'}
        busy={busy}
        onCancel={() => setPending(null)}
        onConfirm={confirmChange}
      >
        {pending && (
          <>
            <p>
              {pending.next === 'admin'
                ? `${labelFor(pending.user)} will be able to access the Mathua administration interface and manage administrator access.`
                : `${labelFor(pending.user)} will immediately lose administrative access.`}
            </p>
            <p className="text-mathua-muted">Recorded in the audit log.</p>
          </>
        )}
      </ConfirmDialog>
    </>
  )
}