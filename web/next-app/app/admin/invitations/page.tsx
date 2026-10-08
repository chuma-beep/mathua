'use client'

import { Suspense, useEffect, useState } from 'react'
import { useRouter, useSearchParams } from 'next/navigation'
import Link from 'next/link'
import Header from '../../../components/Header'
import BottomTabs from '../../../components/BottomTabs'
import Footer from '../../../components/Footer'
import SectionHeader from '../../../components/SectionHeader'
import { acceptAdminInvitation, getErrorMessage } from '../../../lib/api'
import { isLoggedIn } from '../../../lib/auth'
import { useAuthState } from '../../../hooks/useAuthState'

// Accepting a staff invitation.
//
// This is the one admin surface a non-staff account reaches, because it is how a learner becomes
// staff. It requires the account to be signed in and — server-side — to hold the address the
// invitation was issued to, so possessing the token is not by itself enough. On success the role
// applies immediately; the account's existing session carries it because the role is read per
// request, so no re-login is needed.
function AcceptInner() {
  const params = useSearchParams()
  const router = useRouter()
  const { loggedIn } = useAuthState()
  const token = params?.get('token') ?? ''
  const [state, setState] = useState<'idle' | 'working' | 'done' | 'error'>('idle')
  const [message, setMessage] = useState('')

  useEffect(() => {
    if (!token) return
    if (!isLoggedIn()) return
    setState('working')
    acceptAdminInvitation(token)
      .then((res) => {
        setState('done')
        setMessage(`You are now a ${res.role}.`)
        setTimeout(() => router.push('/admin'), 1200)
      })
      .catch((e) => { setState('error'); setMessage(getErrorMessage(e)) })
  }, [token, router])

  return (
    <>
      <SectionHeader label="Administration" title="Invitation" />
      {!token && (
        <p role="alert" className="border border-mathua-red bg-mathua-surface p-4 font-mono text-xs text-mathua-red">
          This link is missing its invitation token.
        </p>
      )}
      {token && !loggedIn && (
        <p role="status" className="border border-mathua-border bg-mathua-surface p-4 font-mono text-xs text-mathua-secondary">
          Sign in with the account this invitation was sent to, then open the link again.
          {' '}<Link href={`/login?return=${encodeURIComponent(`/admin/invitations?token=${token}`)}`} className="text-mathua-blue hover:underline">Go to login →</Link>
        </p>
      )}
      {state === 'working' && (
        <p role="status" className="font-mono text-xs text-mathua-secondary">Applying your invitation…</p>
      )}
      {state === 'done' && (
        <p role="status" className="border border-mathua-green bg-mathua-surface p-4 font-mono text-xs text-mathua-green">{message}</p>
      )}
      {state === 'error' && (
        <p role="alert" className="border border-mathua-red bg-mathua-surface p-4 font-mono text-xs text-mathua-red">{message}</p>
      )}
    </>
  )
}

export default function AcceptInvitationPage() {
  return (
    <>
      <Header />
      <div className="pt-[var(--chrome-top)] lg:pt-0">
        <main className="mx-auto w-full max-w-[720px] min-w-0 px-4 sm:px-6 py-8 sm:py-12 overflow-x-hidden">
          <Suspense fallback={null}>
            <AcceptInner />
          </Suspense>
        </main>
      </div>
      <BottomTabs />
      <Footer />
    </>
  )
}
