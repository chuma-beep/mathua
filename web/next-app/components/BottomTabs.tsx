'use client'

import Link from 'next/link'
import { usePathname } from 'next/navigation'
import { useAuthState } from '../hooks/useAuthState'

const TABS = [
  { label: 'Home', href: '/', icon: '○' },
  { label: 'Study', href: '/study', icon: '◐' },
  { label: 'Practice', href: '/session', icon: 'π' },
  { label: 'Profile', href: '/profile', icon: '◑' },
  { label: 'More', href: '#more', icon: '⋯' },
]

interface Props {
  onMore: () => void
}

export default function BottomTabs({ onMore }: Props) {
  const pathname = usePathname()
  const { loggedIn } = useAuthState()

  return (
    <nav className="lg:hidden fixed bottom-0 left-0 right-0 z-40 h-[56px] border-t-[0.5px] border-mathua-border bg-mathua-bg/95 backdrop-blur flex">
      {TABS.map((t) => {
        const isMore = t.href === '#more'
        const active = !isMore && (pathname === t.href || pathname.startsWith(t.href + '/'))
        const content = (
          <>
            <span className="text-[14px] leading-none">{t.icon}</span>
            <span className="text-[9px] uppercase tracking-wide">{t.label}</span>
          </>
        )
        if (isMore) {
          return (
            <button
              key={t.label}
              onClick={onMore}
              className="flex-1 flex flex-col items-center justify-center gap-0.5 font-mono text-mathua-muted hover:text-mathua-primary"
            >
              {content}
            </button>
          )
        }
        // Hide Profile tab for guests? Keep visible but will show guest profile
        return (
          <Link
            key={t.href}
            href={t.href}
            className={`flex-1 flex flex-col items-center justify-center gap-0.5 font-mono border-t-[2px] transition-colors ${
              active ? 'border-mathua-blue text-mathua-blue' : 'border-transparent text-mathua-muted hover:text-mathua-primary'
            }`}
          >
            {content}
          </Link>
        )
      })}
    </nav>
  )
}
