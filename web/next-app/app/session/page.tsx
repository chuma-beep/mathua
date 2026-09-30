'use client'

import { Suspense, useEffect } from 'react'
import { useRouter, useSearchParams } from 'next/navigation'

// Deprecated route (PR6): the /session chooser pointed at /study instead of
// the Learn loop. /learn (no concept) loads the Next head directly, so this
// stub forwards old links and bookmarks, then deletes itself in a later
// release along with /concept.
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
