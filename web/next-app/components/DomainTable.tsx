'use client'

import Link from 'next/link'
import type { DomainRow, DomainConceptStatus } from '../lib/domainRows'
import { DOMAIN_STATUS_LABEL, DOMAIN_STATUS_ORDER } from '../lib/domainRows'

interface Props {
  rows: DomainRow[]
  /** Expand control. Omitted on the read-only share report, which has nothing to expand into. */
  onToggle?: (domain: string) => void
  expanded?: string | null
  /** Links out. Off for /share: a shared report is not the learner, so /learn is not theirs. */
  links?: boolean
}

// The row itself, shared by /domains and /share.
//
// This was `DomainProgress` on /profile. It hardcoded a 15-entry domain list, so precalculus and
// machine_learning never rendered — 50 concepts invisible — and every row printed a locked count,
// the figure Phase 3 removed from PositionBlock as the one number a learner can neither act on nor
// change today. Both are fixed upstream in lib/domainRows; nothing here chooses a domain list.
export default function DomainTable({ rows, onToggle, expanded, links = true }: Props) {
  if (rows.length === 0) return null

  return (
    <div className="border-[0.5px] border-mathua-border">
      {rows.map((r, i) => {
        const pct = r.total > 0 ? Math.round((r.mastered / r.total) * 100) : 0
        const open = expanded === r.domain
        const counts = (
          <>
            {r.mastered}/{r.total} · {pct}%
            {r.dueForReview > 0 && ` · ${r.dueForReview} due`}
          </>
        )
        return (
          <div key={r.domain} className={i < rows.length - 1 ? 'border-b-[0.5px] border-mathua-border' : ''}>
            <div className="flex flex-wrap sm:flex-nowrap items-center gap-2 sm:gap-[10px] px-3 py-2 sm:px-[10px] sm:py-1.5 min-w-0">
              {links ? (
                <Link
                  href={`/graph?domain=${encodeURIComponent(r.domain)}`}
                  className="min-w-0 w-full sm:w-[130px] shrink-0 truncate text-mathua-primary text-[11px] hover:text-mathua-blue sm:shrink-0"
                >
                  {r.label}
                </Link>
              ) : (
                <div className="min-w-0 w-full sm:w-[130px] shrink-0 truncate text-mathua-primary text-[11px] sm:shrink-0">
                  {r.label}
                </div>
              )}
              <div className="flex-1 min-w-[60px] h-2 bg-mathua-border relative shrink">
                <div
                  className="absolute left-0 top-0 h-full bg-mathua-blue transition-all duration-300"
                  style={{ width: `${pct}%` }}
                />
              </div>
              <div className="hidden sm:block w-[190px] text-right shrink-0 text-mathua-muted text-[11px]">
                {counts}
              </div>
              <div className="sm:hidden w-full text-right text-mathua-muted text-[10px] leading-none">
                {counts}
              </div>
            </div>

            {(onToggle || (links && r.nextId)) && (
              <div className="flex flex-wrap items-center gap-2 px-3 pb-2 sm:px-[10px] min-w-0">
                {links && r.nextId && (
                  <Link
                    href={`/learn?concept=${encodeURIComponent(r.nextId)}`}
                    className="font-mono text-[10px] text-mathua-blue hover:underline min-w-0 truncate"
                  >
                    {r.inProgress > 0 || r.dueForReview > 0 ? 'Continue' : 'Start'}: {r.nextLabel}
                  </Link>
                )}
                {onToggle && (
                  <button
                    type="button"
                    onClick={() => onToggle(r.domain)}
                    aria-expanded={open}
                    className="font-mono text-[10px] text-mathua-muted hover:text-mathua-blue transition-colors"
                  >
                    {open ? 'Hide concepts' : `${r.concepts.length} concepts`}
                  </button>
                )}
              </div>
            )}

            {open &&
              DOMAIN_STATUS_ORDER.map((st: DomainConceptStatus) => {
                const group = r.concepts.filter((c) => c.status === st)
                if (group.length === 0) return null
                return (
                  <div key={st} className="px-3 pb-2 sm:px-[10px] min-w-0">
                    <div className="font-mono text-[10px] uppercase tracking-wider text-mathua-muted mb-1">
                      {DOMAIN_STATUS_LABEL[st]} ({group.length})
                    </div>
                    <ul className="space-y-0.5">
                      {group.map((c) =>
                        links ? (
                          <li key={c.id} className="min-w-0">
                            <Link
                              href={`/learn?concept=${encodeURIComponent(c.id)}`}
                              className="font-mono text-[11px] text-mathua-secondary hover:text-mathua-blue break-words"
                            >
                              {c.label}
                            </Link>
                          </li>
                        ) : (
                          <li key={c.id} className="font-mono text-[11px] text-mathua-secondary break-words">
                            {c.label}
                          </li>
                        ),
                      )}
                    </ul>
                  </div>
                )
              })}
          </div>
        )
      })}
    </div>
  )
}
