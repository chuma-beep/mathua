'use client'

import { useEffect, useState, Suspense } from 'react'
import { useSearchParams } from 'next/navigation'
import Link from 'next/link'
import Header from '../../components/Header'
import BottomTabs from '../../components/BottomTabs'
import SectionHeader from '../../components/SectionHeader'
import Footer from '../../components/Footer'
import Loading from '../../components/Loading'
import ProfileStats from '../../components/ProfileStats'
import ActivityHeatmap from '../../components/ActivityHeatmap'
import DomainProgress from '../../components/DomainProgress'
import { getShareReport, type ShareReport } from '../../lib/api'

function ShareContent() {
  const searchParams = useSearchParams()
  const token = searchParams.get('token') || ''
  const [report, setReport] = useState<ShareReport | null>(null)
  const [error, setError] = useState('')

  useEffect(() => {
    if (!token) {
      setError('Missing share token.')
      return
    }
    getShareReport(token)
      .then(r => r ? setReport(r) : setError('This share link is invalid or has been disabled.'))
      .catch(() => setError('Failed to load the shared report.'))
  }, [token])

  if (error) {
    return (
      <>
        <Header />
        <div className="max-w-container mx-auto px-4 sm:px-6 pt-20 pb-[calc(80px+env(safe-area-inset-bottom))] lg:pb-0 text-center overflow-x-hidden min-w-0">
          <p className="text-mathua-muted text-sm">{error}</p>
          {!token && (
            <p className="text-mathua-muted text-xs font-mono mt-2">Ask the student to enable sharing in Settings → Share with parent / teacher.</p>
          )}
          <Link href="/" className="text-mathua-blue text-sm hover:underline mt-4 inline-block">
            ← Back home
          </Link>
        </div>
        <Footer />
        <BottomTabs />
      </>
    )
  }

  if (!report) {
    return (
      <>
        <Header />
        <div className="max-w-container mx-auto px-4 sm:px-6 pt-20 pb-[calc(80px+env(safe-area-inset-bottom))] lg:pb-0 text-center overflow-x-hidden min-w-0">
          <Loading label="LOADING REPORT" />
        </div>
        <Footer />
        <BottomTabs />
      </>
    )
  }

  return (
    <>
      <Header />
      <div className="max-w-container mx-auto px-4 sm:px-6 pb-[calc(80px+env(safe-area-inset-bottom))] lg:pb-0 overflow-x-hidden min-w-0">
        <section className="pt-8 min-w-0 overflow-hidden">
          <SectionHeader label="Read-only report" title={`${report.name}'s progress`} />
          <p className="text-mathua-muted text-xs font-mono text-center mt-3">
            Shared by the student for parent / teacher oversight
          </p>

          <div className="mx-auto w-full max-w-[820px] min-w-0 px-0 sm:px-4 py-8 overflow-hidden">
            <ProfileStats name={report.name} scores={report.scores} />

            <section className="mt-8 flex min-w-0 flex-col items-stretch">
              <h2 className="font-serif text-[1.05rem] font-normal text-mathua-primary mb-4 w-full">
                Activity
              </h2>
              <div className="w-full max-w-full min-w-0 flex justify-center overflow-hidden">
                <div className="w-full max-w-full min-w-0">
                  <ActivityHeatmap data={report.activity} />
                </div>
              </div>
            </section>

            <section className="mt-8 min-w-0">
              <DomainProgress progress={report.progress} />
            </section>
          </div>
        </section>
        <Footer />
      </div>
      <BottomTabs />
    </>
  )
}

export default function SharePage() {
  return (
    <Suspense fallback={
      <><Header /><div className="max-w-container mx-auto px-4 sm:px-6 pt-20 pb-[calc(80px+env(safe-area-inset-bottom))] lg:pb-0 text-center overflow-x-hidden min-w-0"><Loading label="LOADING REPORT" /></div></>
    }>
      <ShareContent />
    </Suspense>
  )
}
