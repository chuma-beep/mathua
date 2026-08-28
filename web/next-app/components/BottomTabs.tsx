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
      className={`lg:hidden fixed bottom-3 inset-x-3 sm:left-1/2 sm:-translate-x-1/2 sm:w-[min(360px,calc(100%-24px))] sm:max-w-none z-40 h-[56px] rounded-none border-[0.5px] border-mathua-border bg-mathua-surface shadow-none flex transition-transform duration-300 ease-out will-change-transform ${
        hidden ? 'translate-y-[calc(100%+20px)] opacity-0' : 'translate-y-0 opacity-100'
      }`}
    >
      {TABS.map((t) => {
        const active = pathname === t.href || pathname.startsWith(t.href + '/')
        return (
          <Link
            key={t.href}
            href={t.href}
            className={`flex-1 flex flex-col items-center justify-center gap-0.5 font-mono border-t-[2px] transition-colors ${
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
