'use client'

import { Suspense, useEffect } from 'react'
import { useRouter, useSearchParams } from 'next/navigation'

// Deprecated route (PR1): /concept?id=X now lives at /study?concept=X.
// Static export (`output: 'export'`) can't use next.config redirects, so
// this client stub forwards old links and bookmarks, then deletes itself
// in a later release along with /session.
function ConceptRedirect() {
  const router = useRouter()
  const searchParams = useSearchParams()
  useEffect(() => {
    const id = searchParams.get('id')
    router.replace(id ? `/study?concept=${encodeURIComponent(id)}` : '/study')
  }, [router, searchParams])
  return null
}

export default function ConceptPage() {
  return (
    <Suspense>
      <ConceptRedirect />
    </Suspense>
  )
}
