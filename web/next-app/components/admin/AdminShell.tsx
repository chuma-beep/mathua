'use client'

import Link from 'next/link'
import { usePathname } from 'next/navigation'
import { useAdminMe } from '../../hooks/useAdminMe'
import { useTheme } from '../../hooks/useTheme'
import { ADMIN_ITEMS, isAdminNavActive } from './nav'
import { House, SunMoon } from '../icons'

// The admin surface is its own window, not a tab inside the learner shell. It
// has its own navigation, its own width and one explicit way back to the
// learner site, so a moderator never has to read the admin area through the
// same chrome as a lesson. The entry in the profile menu opens this in a new
// tab; nothing in the public header links here.
//
// Permission filtering is convenience only. Every admin endpoint re-reads the
// role on the server, so a hidden section is still refused if typed directly.
export default function AdminShell({ children }: { children: React.ReactNode }) {
  const pathname = usePathname()
  const { can, loading, me } = useAdminMe()
  const { mounted, toggleTheme } = useTheme()

  const items = loading ? [] : ADMIN_ITEMS.filter((item) => can(item.perm))

  return (
    <div className="min-h-screen bg-mathua-bg text-mathua-primary">
      <aside className="fixed inset-y-0 left-0 z-30 hidden w-60 flex-col border-r border-mathua-border bg-mathua-surface md:flex">
        <Link href="/admin" className="flex h-16 items-center gap-2 border-b border-mathua-border px-5">
          <span aria-hidden="true" className="font-mono text-lg text-mathua-blue">λ</span>
          <span className="min-w-0">
            <span className="block font-mono text-xs uppercase tracking-wider text-mathua-primary">Mathua</span>
            <span className="block font-mono text-[10px] uppercase tracking-wider text-mathua-muted">Administration</span>
          </span>
        </Link>

        <nav aria-label="Admin sections" className="min-h-0 flex-1 overflow-y-auto px-3 py-4">
          {loading ? (
            <p className="px-2 font-mono text-[10px] uppercase tracking-wider text-mathua-muted">Loading access</p>
          ) : items.length === 0 ? (
            <p className="px-2 font-mono text-[11px] text-mathua-muted">No staff sections for this account.</p>
          ) : (
            <ul className="space-y-1">
              {items.map((item) => {
                const active = isAdminNavActive(pathname, item.href)
                const Icon = item.icon
                return (
                  <li key={item.href}>
                    <Link
                      href={item.href}
                      aria-current={active ? 'page' : undefined}
                      className={
                        active
                          ? 'flex h-10 items-center gap-2.5 border-l-2 border-mathua-blue bg-mathua-blue-faint px-3 font-mono text-xs text-mathua-blue'
                          : 'flex h-10 items-center gap-2.5 border-l-2 border-transparent px-3 font-mono text-xs text-mathua-secondary hover:border-mathua-border-strong hover:text-mathua-primary'
                      }
                    >
                      <Icon size={15} aria-hidden="true" />
                      <span className="truncate">{item.label}</span>
                    </Link>
                  </li>
                )
              })}
            </ul>
          )}
        </nav>

        <div className="border-t border-mathua-border p-3">
          <Link
            href="/"
            className="flex h-10 items-center gap-2.5 px-3 font-mono text-xs text-mathua-secondary hover:text-mathua-blue"
          >
            <House size={15} aria-hidden="true" />
            <span>Learner site</span>
          </Link>
          {mounted && (
            <button
              type="button"
              onClick={toggleTheme}
              aria-label="Toggle theme"
              className="flex h-10 w-full items-center gap-2.5 px-3 font-mono text-xs text-mathua-secondary hover:text-mathua-blue"
            >
              <SunMoon size={15} isAnimated={false} aria-hidden="true" />
              <span>Toggle theme</span>
            </button>
          )}
          {me && (
            <p className="mt-2 truncate px-3 font-mono text-[10px] text-mathua-muted" title={me.id}>
              {me.role} · {me.id.slice(0, 8)}
            </p>
          )}
        </div>
      </aside>

      <div className="md:pl-60">
        <header className="sticky top-0 z-20 border-b border-mathua-border bg-mathua-surface md:hidden">
          <div className="flex h-14 items-center justify-between px-4">
            <Link href="/admin" className="flex items-center gap-2">
              <span aria-hidden="true" className="font-mono text-base text-mathua-blue">λ</span>
              <span className="font-mono text-[11px] uppercase tracking-wider text-mathua-primary">Mathua Admin</span>
            </Link>
            <Link href="/" className="flex h-11 items-center gap-1.5 font-mono text-[11px] text-mathua-secondary">
              <House size={14} aria-hidden="true" />
              Site
            </Link>
          </div>
          <nav aria-label="Admin sections mobile" className="overflow-x-auto border-t border-mathua-border px-2">
            <ul className="flex w-max gap-1 py-2">
              {items.map((item) => {
                const active = isAdminNavActive(pathname, item.href)
                return (
                  <li key={item.href}>
                    <Link
                      href={item.href}
                      aria-current={active ? 'page' : undefined}
                      className={
                        active
                          ? 'flex h-9 items-center border border-mathua-blue px-2.5 font-mono text-[11px] text-mathua-blue'
                          : 'flex h-9 items-center border border-mathua-border px-2.5 font-mono text-[11px] text-mathua-secondary'
                      }
                    >
                      {item.label}
                    </Link>
                  </li>
                )
              })}
            </ul>
          </nav>
        </header>

        <div className="flex min-h-screen flex-col">
          <div className="flex-1">{children}</div>
          <footer className="border-t border-mathua-border px-4 py-4 font-mono text-[10px] text-mathua-muted sm:px-6">
            Mathua administration. Actions here are audited.
          </footer>
        </div>
      </div>
    </div>
  )
}
