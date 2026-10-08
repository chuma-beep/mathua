'use client'

import Link from 'next/link'
import { usePathname } from 'next/navigation'
import { useAdminMe } from '../hooks/useAdminMe'

// The admin sub-navigation.
//
// Each entry names the permission it needs, and the list is filtered against /api/admin/me, so a
// moderator sees Reports and Content while an owner or admin sees everything. This is a
// convenience and it is labelled as one in the code: hiding a link is not authorization. Every
// section's API is gated server-side, so a caller who types a hidden URL reaches the page and is
// refused by the request.
//
// The permission string is the same one the route is registered under on the server, which is
// the point: the two lists are the same vocabulary and a section cannot be added to one without
// naming it in the other.
const ITEMS: { label: string; href: string; perm: string }[] = [
  { label: 'Overview', href: '/admin', perm: 'admin.access' },
  { label: 'Reports', href: '/admin/reports', perm: 'reports.read' },
  { label: 'Users', href: '/admin/users', perm: 'users.read' },
  { label: 'Contributors', href: '/admin/contributors', perm: 'admins.read' },
  { label: 'Content', href: '/admin/content', perm: 'content.read' },
  { label: 'Audit Log', href: '/admin/audit', perm: 'audit.read' },
]

export default function AdminNav() {
  const pathname = usePathname()
  const { can, loading } = useAdminMe()
  if (loading) return null

  const items = ITEMS.filter((i) => can(i.perm))
  if (items.length === 0) return null

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
