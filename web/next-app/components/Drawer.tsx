'use client'

import Link from 'next/link'
import { usePathname } from 'next/navigation'
import { useAuthState } from '../hooks/useAuthState'

interface DrawerProps {
  collapsed: boolean
  onToggle: () => void
}

const NAV_ITEMS = [
  { label: 'Study', href: '/study', icon: '◐' },
  { label: 'Leaderboard', href: '/leaderboard', icon: '◎' },
  { label: 'Graph', href: '/graph', icon: '⬡' },
  { label: 'Practice', href: '/session', icon: 'π', auth: true },
  { label: 'Profile', href: '/profile', icon: '◑', auth: true },
]

export default function Drawer({ collapsed, onToggle }: DrawerProps) {
  const pathname = usePathname()
  const { loggedIn } = useAuthState()

  const items = NAV_ITEMS.filter((i) => !i.auth || loggedIn)

  return (
    <div className="flex flex-col h-full">
      <div className="flex items-center justify-between px-3 h-[53px] border-b border-mathua-border shrink-0">
        {!collapsed && (
          <span className="font-mono text-[11px] uppercase text-mathua-muted">Navigation</span>
        )}
        <button
          onClick={onToggle}
          aria-label={collapsed ? 'Expand drawer' : 'Collapse drawer'}
          className="font-mono text-xs text-mathua-muted hover:text-mathua-blue border border-mathua-border-strong px-1.5 py-0.5"
        >
          {collapsed ? '→' : '←'}
        </button>
      </div>
      <nav className="flex-1 p-2 space-y-0.5">
        {items.map((item) => {
          const active = pathname === item.href || pathname.startsWith(item.href + '/')
          return (
            <Link
              key={item.href}
              href={item.href}
              className={`flex items-center gap-2.5 px-2.5 py-2 border-l-2 text-xs font-mono transition-colors ${
                active
                  ? 'border-mathua-blue bg-mathua-surface text-mathua-blue'
                  : 'border-transparent text-mathua-muted hover:text-mathua-primary hover:bg-mathua-surface'
              }`}
              title={collapsed ? item.label : undefined}
            >
              <span className="text-[13px] w-4 text-center shrink-0">{item.icon}</span>
              {!collapsed && <span className="truncate">{item.label}</span>}
            </Link>
          )
        })}
      </nav>
      {!collapsed && (
        <div className="p-3 border-t border-mathua-border">
          <p className="font-mono text-[10px] text-mathua-muted leading-relaxed">
            Adaptive math — 437 concepts · mastery gated
          </p>
        </div>
      )}
    </div>
  )
}
