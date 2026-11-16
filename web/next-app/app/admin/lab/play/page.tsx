'use client'

import AdminShell from '@/components/admin/AdminShell'
import { Suspense } from 'react'
import Link from 'next/link'
import { useSearchParams } from 'next/navigation'
import SectionHeader from '../../../../components/SectionHeader'
import LearnStepper from '../../../../components/LearnStepper'

// The learner route for one concept, mounted from the lab.
//
// This is the same component /learn?concept=... renders, not a mock-up of it, so
// an operator can see what a learner sees: the knowledge points, the served
// question, the field it gets and the feedback block. It is a real learning
// session, answering here writes evidence, which is why the header says so.
//
// A concept whose prerequisites are unmet still shows the lock the learner would
// see. That is not a limitation of the preview; it is the behaviour under test.
// For a locked concept the lab's generator is the place to see the questions.
function Playground() {
  const params = useSearchParams()
  const concept = params.get('concept') ?? ''

  return (
    <>
      {!concept ? (
        <p className="border border-mathua-border bg-mathua-surface p-4 font-mono text-xs text-mathua-secondary">
          No concept named. Pick one in the{' '}
          <Link href="/admin/lab" className="text-mathua-blue">
            lab
          </Link>
          .
        </p>
      ) : (
        <>
          <div className="mb-4 border border-mathua-border bg-mathua-surface p-4">
            <p className="font-mono text-[11px] leading-relaxed text-mathua-secondary">
              This is the live learner route for <span className="text-mathua-primary">{concept}</span>. Answers
              submitted here count as a real attempt for your account. To try an answer without writing evidence, use
              the{' '}
              <Link href="/admin/lab" className="text-mathua-blue">
                lab
              </Link>
              .
            </p>
          </div>
          <LearnStepper conceptId={concept} />
        </>
      )}
    </>
  )
}

export default function AdminLabPlayPage() {
  return (
    <AdminShell>
      <main className="mx-auto w-full max-w-[1200px] min-w-0 px-4 sm:px-6 py-8 sm:py-12 overflow-x-hidden">
        <SectionHeader label="Administration" title="Learner view" />
          <Suspense fallback={<p className="font-mono text-xs text-mathua-muted">Loading…</p>}>
          <Playground />
        </Suspense>
      </main>
    </AdminShell>
  )
}
