'use client'

import { Suspense, useEffect, useState } from 'react'
import { useSearchParams } from 'next/navigation'
import Link from 'next/link'
import Header from '../../components/Header'
import BottomTabs from '../../components/BottomTabs'
import SectionHeader from '../../components/SectionHeader'
import Footer from '../../components/Footer'
import Loading from '../../components/Loading'
import { verifyEmail } from '../../lib/api'
import { getUserInfo, setUserInfo } from '../../lib/auth'

function VerifyEmailInner() {
  const searchParams = useSearchParams()
  const [status, setStatus] = useState<'working' | 'done' | 'error'>('working')
  const [message, setMessage] = useState('')

  useEffect(() => {
    const token = searchParams.get('token')
    if (!token) {
      setStatus('error')
      setMessage('This verification link is missing its token.')
      return
    }
    verifyEmail(token)
      .then(() => {
        const info = getUserInfo()
        if (info) setUserInfo({ ...info, email_verified: true })
        setStatus('done')
      })
      .catch((e: unknown) => {
        setStatus('error')
        setMessage(e instanceof Error ? e.message : 'Verification failed — the link may have expired.')
      })
  }, [searchParams])

  return (
    <>
      <Header />
      <div className="max-w-container mx-auto px-4 sm:px-6 pb-[calc(80px+env(safe-area-inset-bottom))] lg:pb-0 overflow-x-hidden min-w-0">
        <section className="pt-8 max-w-md mx-auto mt-8 sm:mt-12 min-w-0 overflow-hidden text-center">
          <SectionHeader label="Email" title="Verify email" />
          <div className="mt-6">
            {status === 'working' && <Loading label="VERIFYING" />}
            {status === 'done' && (
              <>
                <p className="font-mono text-sm text-mathua-green">Email verified.</p>
                <p className="text-mathua-secondary text-sm mt-2">Your logins now merge into one account.</p>
                <Link href="/profile" className="inline-block mt-6 border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white px-6 py-2 font-mono text-xs min-h-[36px]">
                  Continue to profile →
                </Link>
              </>
            )}
            {status === 'error' && <p className="font-mono text-sm text-mathua-red">{message}</p>}
          </div>
        </section>
      </div>
      <Footer />
      <BottomTabs />
    </>
  )
}

export default function VerifyEmailPage() {
  return (
    <Suspense fallback={<div style={{ background: 'var(--bg)', minHeight: '100vh' }} />}>
      <VerifyEmailInner />
    </Suspense>
  )
}
