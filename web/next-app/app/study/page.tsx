'use client'

import { Suspense, useEffect } from 'react'
import { useRouter, useSearchParams } from 'next/navigation'

// Closed route: Study is no longer a place a learner goes.
//
// It used to be two surfaces at once — the learning loop's answer seam is still called
// `/api/study/answer`, and the page was a browsable library of whole lessons with its own
// search, domain index and a link in every header, tab bar, sidebar and footer. The second one
// is what is closed. A learner was sent to browse a reference library before being offered
// anything to solve, and could read their way to a concept instead of being taught it, which
// makes reading look equivalent to learning when only demonstrated answers count.
//
// The material did not go with the route. It is now the reference panel inside `/learn`, scoped
// to the concept being learned, and it can explain but cannot hand over a curriculum.
//
// This stub forwards rather than 404s, because old links and bookmarks exist. Static export
// (`output: 'export'`) cannot use next.config redirects, so it is a client stub — the same
// pattern as `/concept` and `/session`, and the cost is one blank frame before the redirect.
//
// The concept parameter is carried, because `/study?concept=X` was the deep link for one
// concept's reference and `/learn?concept=X` is the learning loop for the same concept: the
// learner lands on the teaching slice for it rather than on the article. `?domain=` and
// `?lesson=` have no `/learn` equivalent and are dropped rather than guessed at — a domain
// index is orientation, and that is `/domains`.
//
// Pinned by test/deprecatedRoutes.test.tsx, so whichever way the decision to delete this
// eventually goes, the change is deliberate.
function StudyRedirect() {
  const router = useRouter()
  const searchParams = useSearchParams()
  useEffect(() => {
    const concept = searchParams.get('concept')
    router.replace(concept ? `/learn?concept=${encodeURIComponent(concept)}` : '/learn')
  }, [router, searchParams])
  return null
}

export default function StudyPage() {
  return (
    <Suspense>
      <StudyRedirect />
    </Suspense>
  )
}