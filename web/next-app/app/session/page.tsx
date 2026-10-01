'use client'

import { Suspense, useEffect } from 'react'
import { useRouter, useSearchParams } from 'next/navigation'

// Deprecated route (PR6): the /session chooser pointed at /study instead of
// the Learn loop. /learn (no concept) loads the Next head directly, so this
// stub forwards old links and bookmarks.
//
// Same open question as /concept, and for the same reasons — nothing internal
// links here, and the only caller is an inbound link from outside. See
// test/deprecatedRoutes.test.tsx, which pins what this forwards so the decision
// to keep or delete is made on evidence rather than by accident.
function SessionRedirect() {
  const router = useRouter()
  const searchParams = useSearchParams()
  useEffect(() => {
    const concept = searchParams.get('concept')
    router.replace(concept ? `/learn?concept=${encodeURIComponent(concept)}` : '/learn')
  }, [router, searchParams])
  return null
}

export default function SessionPage() {
  return (
    <Suspense>
      <SessionRedirect />
    </Suspense>
  )
}
