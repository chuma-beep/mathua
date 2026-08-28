'use client'

import Link from 'next/link'
import { usePathname } from 'next/navigation'
import { useScrollDirection } from '../hooks/useScrollDirection'

const TABS = [
  { label: 'Home', href: '/', icon: '○' },
  { label: 'Study', href: '/study', icon: '◐' },
  { label: 'Practice', href: '/session', icon: 'π' },
  { label: 'Profile', href: '/profile', icon: '◑' },
  { label: 'Graph', href: '/graph', icon: '⬡' },
]

export default function BottomTabs() {
  const pathname = usePathname()
  const hidden = useScrollDirection({ hideThreshold: 12, topOffset: 40, idleMs: 300, bottomOffset: 24 })

  return (
    <nav
      className={`lg:hidden fixed inset-x-3 z-40 flex h-[56px] border-[0.5px] border-mathua-border bg-mathua-surface shadow-none transition-transform duration-300 ease-out will-change-transform ${
        hidden ? 'translate-y-[calc(100%+20px)] opacity-0' : 'translate-y-0 opacity-100'
      }`}
      style={{
        bottom: 'max(12px, env(safe-area-inset-bottom))',
        paddingBottom: 'env(safe-area-inset-bottom)',
        left: '12px',
        right: '12px',
        maxWidth: 'min(360px, calc(100% - 24px))',
        marginLeft: 'auto',
        marginRight: 'auto',
      }}
    >
      {TABS.map((t) => {
        const active = pathname === t.href || pathname.startsWith(t.href + '/')
        return (
          <Link
            key={t.href}
            href={t.href}
            className={`flex-1 flex flex-col items-center justify-center gap-0.5 font-mono border-t-[2px] min-h-[44px] transition-colors ${
              active ? 'border-mathua-blue text-mathua-blue' : 'border-transparent text-mathua-muted hover:text-mathua-primary'
            }`}
          >
            <span className="text-[14px] leading-none">{t.icon}</span>
            <span className="text-[9px] uppercase tracking-wide">{t.label}</span>
          </Link>
        )
      })}
    </nav>
  )
}
