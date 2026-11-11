'use client'

import { Suspense, useEffect, useState } from 'react'
import { useSearchParams } from 'next/navigation'
import Link from 'next/link'
import AppShell from '../../../components/app-shell'
import SectionHeader from '../../../components/SectionHeader'
import Footer from '../../../components/Footer'
import Loading from '../../../components/Loading'
import { unsubscribeNotifications } from '../../../lib/api'

function UnsubscribeInner() {
  const searchParams = useSearchParams()
  const [status, setStatus] = useState<'working' | 'done' | 'error'>('working')
  const [message, setMessage] = useState('')

  useEffect(() => {
    const token = searchParams.get('token')
    if (!token) {
      setStatus('error')
      setMessage('This unsubscribe link is missing its token.')
      return
    }
    unsubscribeNotifications(token)
      .then(() => setStatus('done'))
      .catch((e) => {
        setStatus('error')
        setMessage(e instanceof Error ? e.message : 'Unsubscribe failed, the link may have expired.')
      })
  }, [searchParams])

  return (
    <AppShell>
      <div className="max-w-container mx-auto px-4 sm:px-6 pb-[calc(80px+env(safe-area-inset-bottom))] lg:pb-0 overflow-x-hidden min-w-0">
        <section className="pt-8 max-w-md mx-auto mt-8 sm:mt-12 min-w-0 overflow-hidden text-center">
          <SectionHeader label="Email" title="Unsubscribe" />
          <div className="mt-6">
            {status === 'working' && <Loading label="UPDATING" />}
            {status === 'done' && (
              <>
                <p className="font-mono text-sm text-mathua-green">Email reminders are off.</p>
                <p className="text-mathua-secondary text-sm mt-2">
                  In-app notifications are unchanged. You can turn email back on in settings.
                </p>
                <Link href="/settings#notifications" className="inline-block mt-6 border border-mathua-blue text-mathua-blue hover:bg-mathua-blue-faint px-6 py-2 font-mono text-xs min-h-[36px]">
                  Notification settings →
                </Link>
              </>
            )}
            {status === 'error' && (
              <>
                <p className="font-mono text-sm text-mathua-red">{message}</p>
                <Link href="/settings#notifications" className="inline-block mt-6 border border-mathua-border text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue px-6 py-2 font-mono text-xs min-h-[36px]">
                  Open notification settings
                </Link>
              </>
            )}
          </div>
        </section>
      </div>
      <Footer />
    </AppShell>
  )
}

export default function UnsubscribePage() {
  return (
    <Suspense fallback={<div style={{ background: 'var(--bg)', minHeight: '100vh' }} />}>
      <UnsubscribeInner />
    </Suspense>
  )
}
