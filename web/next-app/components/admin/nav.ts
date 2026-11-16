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
}

export const ADMIN_ITEMS: AdminNavItem[] = [
  { label: 'Overview', href: '/admin', perm: 'admin.access', icon: ChartNoAxesColumn },
  { label: 'Reports', href: '/admin/reports', perm: 'reports.read', icon: FileText },
  { label: 'Users', href: '/admin/users', perm: 'users.read', icon: Users },
  { label: 'Leagues', href: '/admin/leagues', perm: 'users.read', icon: Trophy },
  { label: 'Notifications', href: '/admin/notifications', perm: 'users.read', icon: Bell },
  { label: 'Contributors', href: '/admin/contributors', perm: 'admins.read', icon: Handshake },
  { label: 'Content', href: '/admin/content', perm: 'content.read', icon: Layers },
  { label: 'Announcements', href: '/admin/announcements', perm: 'content.read', icon: MessagesSquare },
  { label: 'Lab', href: '/admin/lab', perm: 'content.read', icon: Dices },
  { label: 'Audit Log', href: '/admin/audit', perm: 'audit.read', icon: Shield },
]

export function isAdminNavActive(pathname: string, href: string): boolean {
  if (!pathname) return false
  return pathname === href || pathname.startsWith(href + '/')
}
