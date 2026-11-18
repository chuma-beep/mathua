'use client'

import type { LeaderboardBoard, LeaderboardView } from '@/lib/api'
import { avatarURL, nameOf } from '@/lib/leaderboard-display'
import { scoreLabel } from '@/lib/leaderboard-period'
import LeagueAvatar from '@/components/LeagueAvatar'
import { LeaderboardTableSkeleton } from '@/components/skeletons/LeaderboardSkeleton'
import { Table, TableBody, TableCaption, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'

export default function BoardTable({
  board,
  view,
  loading,
  caption,
}: {
  board: LeaderboardBoard | null
  view: LeaderboardView
  loading: boolean
  caption: string
}) {
  return (
    <div className="w-[calc(100%+2rem)] -mx-4 sm:w-full sm:mx-0">
      <Table className="min-w-[320px] border border-mathua-border bg-mathua-surface">
        <TableCaption className="sr-only">{caption}</TableCaption>
        <TableHeader>
          <TableRow className="hover:bg-transparent">
            <TableHead className="w-10 bg-mathua-surface-elevated p-2.5 text-left sm:w-12 sm:p-[14px_20px]">#</TableHead>
            <TableHead className="bg-mathua-surface-elevated p-2.5 sm:p-[14px_20px]">Learner</TableHead>
            <TableHead className="bg-mathua-surface-elevated p-2.5 text-right sm:p-[14px_20px]">
              {scoreLabel(view, board?.final)}
            </TableHead>
            <TableHead className="hidden bg-mathua-surface-elevated p-2.5 text-right md:table-cell sm:p-[14px_20px]">
              Level
            </TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {loading && <LeaderboardTableSkeleton rows={8} />}
          {!loading &&
            board?.entries.map((row) => (
              <TableRow
                key={row.id}
                className={row.is_me ? 'bg-mathua-blue-faint hover:bg-mathua-blue-faint' : undefined}
              >
                <TableCell className="p-2.5 text-sm text-mathua-muted sm:p-[14px_20px]">{row.rank}</TableCell>
                <TableCell className="max-w-[160px] p-2.5 text-sm font-medium text-mathua-primary sm:max-w-none sm:p-[14px_20px]">
                  <span className="flex min-w-0 items-center gap-2">
                    <LeagueAvatar
                      tier={row.league}
                      weeklyXp={row.league_xp}
                      seed={row.id}
                      name={nameOf(row)}
                      size={28}
                      url={avatarURL(row)}
                      className="shrink-0"
                    />
                    <span className="block truncate" title={nameOf(row)}>
                      {nameOf(row)}
                    </span>
                    {row.is_me && (
                      <span className="shrink-0 font-mono text-[10px] uppercase text-mathua-blue">you</span>
                    )}
                  </span>
                </TableCell>
                <TableCell className="whitespace-nowrap p-2.5 text-right text-sm text-mathua-green sm:p-[14px_20px]">
                  {row.xp}
                </TableCell>
                <TableCell className="hidden whitespace-nowrap p-2.5 text-right text-sm text-mathua-secondary md:table-cell sm:p-[14px_20px]">
                  {row.level}
                </TableCell>
              </TableRow>
            ))}
        </TableBody>
      </Table>
    </div>
  )
}
