'use client'

import Link from 'next/link'
import { usePathname } from 'next/navigation'

// The admin sub-navigation.
//
// Three entries, all implemented. There is deliberately no "Curriculum", "Questions" or
// "Publishing" here: those are real future phases and a navigation entry that leads nowhere is
// worse than its absence, because it advertises a capability the deployment does not have.
//
// Note what is NOT here either: the existing content-triage page at /admin/reports. That one is
// gated by the operator's shared password rather than by a role, so it belongs to a different
// trust model and putting it in this list would imply the two authorize each other. It stays
// reachable at its own URL.

// AdminNav renders nothing for a caller who is not an administrator.
//
// This is a convenience and it is labelled as one in the code: hiding a link is not
// authorization. The backend reads the role on every request, so a learner who types /admin/users
// directly reaches the page and is refused by the API with a 403 — which this nav is careful not
// to pretend otherwise about.
export default function AdminNav({ isAdmin }: { isAdmin?: boolean }) {
  const pathname = usePathname()
  if (!isAdmin) return null

  const items = [
    { label: 'Overview', href: '/admin' },
    { label: 'Users', href: '/admin/users' },
    { label: 'Audit', href: '/admin/audit' },
  ]

  return (
    <nav aria-label="Admin sections" className="mb-6 flex flex-wrap gap-2">
      {items.map((item) => {
        const active = pathname === item.href
        return (
          <Link
            key={item.href}
            href={item.href}
            aria-current={active ? 'page' : undefined}
            className={
              active
                ? 'border border-mathua-blue text-mathua-blue px-3 h-11 inline-flex items-center font-mono text-xs'
                : 'border border-mathua-border text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue px-3 h-11 inline-flex items-center font-mono text-xs'
            }
          >
            {item.label}
          </Link>
        )
      })}
    </nav>
  )
}