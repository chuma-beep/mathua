'use client'

import { useCallback, useEffect, useState } from 'react'
import Header from '../../../components/Header'
import BottomTabs from '../../../components/BottomTabs'
import Footer from '../../../components/Footer'
import SectionHeader from '../../../components/SectionHeader'
import AdminNav from '../../../components/AdminNav'
import Loading from '../../../components/Loading'
import {
  listAdmins, listAdminInvitations, createAdminInvitation, revokeAdminInvitation,
  getErrorMessage, type AdminUser, type AdminInvitation, type Role,
} from '../../../lib/api'
import { useAdminMe } from '../../../hooks/useAdminMe'

// Contributors. Adding someone to the team is a deliberate workflow, not a config change: an
// owner or admin issues an invitation bound to an email, the invitation is a single-use,
// time-limited token, and accepting it is what assigns the role. The raw token is shown once and
// stored only as a hash, so a database read cannot be replayed as access.
//
// The role model is small on purpose. Owner is the only role that can mint another owner; admin
// runs the place but cannot crown a peer; moderator works the moderation queue and nothing else.
export default function ContributorsPage() {
  const { can, me } = useAdminMe()
  const [staff, setStaff] = useState<AdminUser[] | null>(null)
  const [invitations, setInvitations] = useState<AdminInvitation[] | null>(null)
  const [error, setError] = useState('')

  const canManage = can('admins.manage')

  const load = useCallback(() => {
    setError('')
    Promise.all([listAdmins(), listAdminInvitations()])
      .then(([a, i]) => { setStaff(a.admins); setInvitations(i.invitations) })
      .catch((e) => setError(getErrorMessage(e)))
  }, [])

  useEffect(() => { load() }, [load])

  return (
    <>
      <Header />
      <div className="pt-[var(--chrome-top)] lg:pt-0">
        <main className="mx-auto w-full max-w-[900px] min-w-0 px-4 sm:px-6 py-8 sm:py-12 overflow-x-hidden">
          <SectionHeader label="Administration" title="Contributors" />
          <AdminNav />

          {error && <p role="alert" className="mb-4 border border-mathua-red bg-mathua-surface p-4 font-mono text-xs text-mathua-red">{error}</p>}

          {canManage && <InviteForm assignable={me?.assignable_roles ?? []} onCreated={load} />}

          <section className="mt-8">
            <h2 className="mb-3 font-serif text-lg text-mathua-primary">Current staff</h2>
            {!staff ? <Loading label="LOADING STAFF" /> : (
              <ul className="space-y-1">
                {staff.map((u) => (
                  <li key={u.id} className="flex items-center justify-between border border-mathua-border bg-mathua-surface px-4 py-3">
                    <span className="min-w-0">
                      <span className="font-mono text-xs text-mathua-primary">{u.name}</span>
                      <span className="ml-2 font-mono text-[11px] text-mathua-muted">{u.email || u.username || u.id}</span>
                    </span>
                    <span className="font-mono text-[11px] uppercase tracking-wider text-mathua-blue">{u.role}</span>
                  </li>
                ))}
              </ul>
            )}
          </section>

          <section className="mt-8">
            <h2 className="mb-3 font-serif text-lg text-mathua-primary">Invitations</h2>
            {!invitations ? <Loading label="LOADING INVITATIONS" /> : invitations.length === 0 ? (
              <p className="font-mono text-xs text-mathua-muted">None yet.</p>
            ) : (
              <ul className="space-y-1">
                {invitations.map((inv) => {
                  const settled = inv.accepted_at || inv.revoked_at
                  const state = inv.accepted_at ? 'accepted' : inv.revoked_at ? 'revoked' : 'pending'
                  return (
                    <li key={inv.id} className="flex items-center justify-between gap-3 border border-mathua-border bg-mathua-surface px-4 py-3">
                      <span className="min-w-0">
                        <span className="font-mono text-xs text-mathua-primary">{inv.email}</span>
                        <span className="ml-2 font-mono text-[11px] uppercase tracking-wider text-mathua-blue">{inv.role}</span>
                        <span className="ml-2 font-mono text-[10px] text-mathua-muted">{state} · expires {inv.expires_at}</span>
                      </span>
                      {canManage && !settled && (
                        <button
                          type="button"
                          onClick={() => { revokeAdminInvitation(inv.id).then(load).catch((e) => setError(getErrorMessage(e))) }}
                          className="shrink-0 border border-mathua-border text-mathua-secondary hover:text-mathua-red px-3 h-8 font-mono text-[11px]"
                        >
                          Revoke
                        </button>
                      )}
                    </li>
                  )
                })}
              </ul>
            )}
          </section>
        </main>
      </div>
      <BottomTabs />
      <Footer />
    </>
  )
}

function InviteForm({ assignable, onCreated }: { assignable: Role[]; onCreated: () => void }) {
  const [email, setEmail] = useState('')
  const [role, setRole] = useState<Role>(assignable[0] ?? 'moderator')
  const [busy, setBusy] = useState(false)
  const [error, setError] = useState('')
  const [token, setToken] = useState<{ token: string; email: string } | null>(null)

  async function submit() {
    setBusy(true); setError('')
    try {
      const res = await createAdminInvitation(email.trim(), role)
      setToken({ token: res.token, email: res.email })
      setEmail('')
      onCreated()
    } catch (e) {
      setError(getErrorMessage(e))
    } finally {
      setBusy(false)
    }
  }

  return (
    <section className="border border-mathua-border bg-mathua-surface p-4">
      <h2 className="mb-3 font-serif text-lg text-mathua-primary">Invite a contributor</h2>
      <div className="flex flex-col gap-2 sm:flex-row">
        <input
          type="email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          placeholder="their account email"
          className="flex-1 border border-mathua-border bg-mathua-code px-3 h-11 font-mono text-xs text-mathua-primary"
        />
        <select
          value={role}
          onChange={(e) => setRole(e.target.value as Role)}
          className="border border-mathua-border bg-mathua-code px-3 h-11 font-mono text-xs text-mathua-primary"
        >
          {assignable.map((r) => <option key={r} value={r}>{r}</option>)}
        </select>
        <button
          type="button"
          disabled={busy || !email.trim()}
          onClick={submit}
          className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue-faint px-4 h-11 font-mono text-xs disabled:opacity-50"
        >
          Create invitation
        </button>
      </div>
      {error && <p role="alert" className="mt-2 font-mono text-[11px] text-mathua-red">{error}</p>}
      {token && (
        <div className="mt-4 border border-mathua-green bg-mathua-code p-3">
          <p className="font-mono text-[11px] text-mathua-green">
            Invitation for {token.email}. This link is shown once and never stored — send it to them now.
          </p>
          <code className="mt-2 block break-all font-mono text-[11px] text-mathua-primary">/admin/invitations?token={token.token}</code>
        </div>
      )}
    </section>
  )
}
