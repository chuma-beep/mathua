import {
  Bell,
  ChartNoAxesColumn,
  Dices,
  FileText,
  Handshake,
  Layers,
  MessagesSquare,
  Shield,
  Trophy,
  Users,
  type IconComponent,
} from '../icons'

// One list for the admin shell and the legacy horizontal nav. Each entry names
// the permission its route is registered under, so hiding a link and refusing
// the API use the same vocabulary. Hiding is convenience only: every admin
// endpoint re-reads the role on the server.
export interface AdminNavItem {
  label: string
  href: string
  perm: string
  icon: IconComponent
  group: 'Operate' | 'Publish' | 'Review'
}

export const ADMIN_ITEMS: AdminNavItem[] = [
  { label: 'Overview', href: '/admin', perm: 'admin.access', icon: ChartNoAxesColumn, group: 'Operate' },
  { label: 'Reports', href: '/admin/reports', perm: 'reports.read', icon: FileText, group: 'Operate' },
  { label: 'Users', href: '/admin/users', perm: 'users.read', icon: Users, group: 'Operate' },
  { label: 'Leagues', href: '/admin/leagues', perm: 'users.read', icon: Trophy, group: 'Operate' },
  { label: 'Notifications', href: '/admin/notifications', perm: 'users.read', icon: Bell, group: 'Operate' },
  { label: 'Contributors', href: '/admin/contributors', perm: 'admins.read', icon: Handshake, group: 'Publish' },
  { label: 'Content', href: '/admin/content', perm: 'content.read', icon: Layers, group: 'Publish' },
  { label: 'Announcements', href: '/admin/announcements', perm: 'content.read', icon: MessagesSquare, group: 'Publish' },
  { label: 'Lab', href: '/admin/lab', perm: 'content.read', icon: Dices, group: 'Publish' },
  { label: 'Audit Log', href: '/admin/audit', perm: 'audit.read', icon: Shield, group: 'Review' },
]

export function isAdminNavActive(pathname: string, href: string): boolean {
  if (!pathname) return false
  // Overview is the section root, not a prefix of every admin route.
  if (href === '/admin') return pathname === '/admin'
  return pathname === href || pathname.startsWith(href + '/')
}

export const ADMIN_GROUPS: AdminNavItem['group'][] = ['Operate', 'Publish', 'Review']

// Subpages that are not sections of their own. The mapping keeps a breadcrumb
// readable without adding a route table the nav already owns.
const SUBPAGE_LABELS: Record<string, string> = {
  detail: 'Account',
  play: 'Playground',
}

export interface AdminCrumb {
  label: string
  href: string
}

// adminBreadcrumb returns the section and any subpage trail for a path. The
// section root has no crumb because the page heading already names it.
export function adminBreadcrumb(pathname: string): AdminCrumb[] {
  if (!pathname) return []
  const section = [...ADMIN_ITEMS]
    .sort((a, b) => b.href.length - a.href.length)
    .find((item) => item.href !== '/admin' && isAdminNavActive(pathname, item.href))
  if (!section) return []
  const rest = pathname.slice(section.href.length).split('/').filter(Boolean)
  if (rest.length === 0) return []
  return [
    { label: section.label, href: section.href },
    ...rest.map((segment, index) => ({
      label: SUBPAGE_LABELS[segment] ?? segment.replace(/-/g, ' '),
      href: `${section.href}/${rest.slice(0, index + 1).join('/')}`,
    })),
  ]
}
