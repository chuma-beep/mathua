'use client'

import { Skeleton } from '../ui/skeleton'

// Row-level skeletons for leaderboard/page.tsx: table body + leagues.
// Page keeps its static header/countdown/formula rendered — only the
// dynamic zones below are placeholder rows, so no layout shift on swap.
export function LeaderboardTableSkeleton({ rows = 8 }: { rows?: number }) {
  return (
    <>
      {Array.from({ length: rows }).map((_, i) => (
        <tr key={`sk-${i}`} className="border-b border-mathua-border last:border-b-0" aria-hidden="true">
          <td className="p-2.5 sm:p-[14px_20px]">
            <Skeleton className="h-4 w-6" />
          </td>
          <td className="p-2.5 sm:p-[14px_20px]">
            <span className="flex items-center gap-2 min-w-0">
              <Skeleton className="size-7 rounded-full shrink-0" />
              <Skeleton className="h-4 flex-1 min-w-0" />
            </span>
          </td>
          <td className="p-2.5 sm:p-[14px_20px]">
            <Skeleton className="h-4 w-10 ml-auto" />
          </td>
          <td className="p-2.5 sm:p-[14px_20px]">
            <Skeleton className="h-4 w-10 ml-auto" />
          </td>
          <td className="p-2.5 sm:p-[14px_20px] hidden md:table-cell">
            <Skeleton className="h-4 w-16 ml-auto" />
          </td>
        </tr>
      ))}
    </>
  )
}

export function LeaguesSkeleton() {
  return (
    <div
      role="status"
      aria-live="polite"
      aria-busy="true"
      aria-label="Loading leagues"
      className="grid grid-cols-1 md:grid-cols-2 gap-4 w-full max-w-full min-w-0 overflow-hidden"
    >
      {Array.from({ length: 2 }).map((_, c) => (
        <div key={c} className="border border-mathua-border bg-mathua-surface w-full min-w-0 overflow-hidden">
          <div className="flex items-center justify-between px-4 py-3 bg-mathua-surface-elevated border-b border-mathua-border">
            <Skeleton className="h-4 w-24" />
            <Skeleton className="h-3 w-16" />
          </div>
          <div className="divide-y divide-mathua-border">
            {Array.from({ length: 4 }).map((_, i) => (
              <div key={i} className="flex items-center gap-3 px-4 py-2 min-w-0" aria-hidden="true">
                <Skeleton className="h-3 w-6 shrink-0" />
                <Skeleton className="size-6 rounded-full shrink-0" />
                <Skeleton className="h-3 flex-1 min-w-0" />
                <Skeleton className="h-3 w-16 shrink-0" />
              </div>
            ))}
          </div>
        </div>
      ))}
    </div>
  )
}
