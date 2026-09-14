'use client'

import { Skeleton } from '../ui/skeleton'

// Mirrors study/page.tsx DomainOverview: banner strip + domain card grid.
export default function StudySkeleton() {
  return (
    <div
      role="status"
      aria-live="polite"
      aria-busy="true"
      aria-label="Loading lessons"
      className="max-w-container mx-auto px-4 sm:px-6 pb-[calc(80px+env(safe-area-inset-bottom))] lg:pb-0 overflow-x-hidden min-w-0"
    >
      <section className="pt-8 min-w-0 overflow-hidden">
        <Skeleton className="h-10 w-full mb-4" />
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
          {Array.from({ length: 6 }).map((_, i) => (
            <div key={i} className="border border-mathua-border bg-mathua-surface p-4">
              <Skeleton className="h-4 w-2/3 mb-3" />
              <Skeleton className="h-3 w-full mb-2" />
              <Skeleton className="h-3 w-4/5 mb-4" />
              <Skeleton className="h-2 w-1/2" />
            </div>
          ))}
        </div>
      </section>
    </div>
  )
}
