'use client'

import { Suspense, useEffect } from 'react'
import { useRouter, useSearchParams } from 'next/navigation'

// Deprecated route (PR1): /concept?id=X now lives at /study?concept=X.
// Static export (`output: 'export'`) can't use next.config redirects, so
// this client stub forwards old links and bookmarks.
//
// Whether it should still exist is a question about external traffic, not
// about this codebase: nothing here links to it, it is not in nav, it is not
// in INDEXABLE_ROUTES, and no page requests it. The only thing it serves is
// an inbound link from outside. Deleting it turns an old bookmark into a 404;
// keeping it costs one build slot and a blank frame before the redirect.
//
// The behaviour is pinned by test/deprecatedRoutes.test.tsx, so whichever way
// that is decided, the change is deliberate and the loss is visible.
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
