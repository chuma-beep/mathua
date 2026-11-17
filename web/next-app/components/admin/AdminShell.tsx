'use client'

import { useEffect, useState } from 'react'
import Link from 'next/link'
import { usePathname } from 'next/navigation'
import { useAdminMe } from '../../hooks/useAdminMe'
import { useTheme } from '../../hooks/useTheme'
import { ADMIN_GROUPS, ADMIN_ITEMS, type AdminNavItem, isAdminNavActive } from './nav'
import { ChevronsLeft, ChevronsRight, Close, House, PanelLeft, SunMoon } from '../icons'
import { SheetClose } from '../ui/sheet'
import { Sheet, SheetContent, SheetTitle, SheetTrigger } from '../ui/sheet'

// The admin surface is its own window, not a tab inside the learner shell. It
// has its own navigation, its own width and one explicit way back to the
// learner site, so a moderator never has to read the admin area through the
// same chrome as a lesson. The entry in the profile menu opens this in a new
// tab; nothing in the public header links here.
//
// Permission filtering is convenience only. Every admin endpoint re-reads the
// role on the server, so a hidden section is still refused if typed directly.

type SidebarProps = {
  items: AdminNavItem[]
  pathname: string | null
  loading: boolean
  me: { id: string; role: string } | null
  mounted: boolean
  toggleTheme: () => void
  collapsed: boolean
  onToggle?: () => void
  onNavigate?: () => void
}

function AdminSidebar({ items, pathname, loading, me, mounted, toggleTheme, collapsed, onToggle, onNavigate }: SidebarProps) {
  return (
    <>
      <div className={`flex h-16 items-center border-b border-mathua-border ${collapsed ? 'justify-center px-1' : 'justify-between px-5'}`}>
        <Link href="/admin" onClick={onNavigate} className="flex min-w-0 items-center gap-2">
          <span aria-hidden="true" className="font-mono text-lg text-mathua-blue">λ</span>
          {!collapsed && (
            <span className="min-w-0">
              <span className="block font-mono text-xs uppercase tracking-wider text-mathua-primary">Mathua</span>
              <span className="block font-mono text-[10px] uppercase tracking-wider text-mathua-muted">Administration</span>
            </span>
          )}
        </Link>
        {onToggle && (
          <button
            type="button"
            onClick={onToggle}
            aria-label={collapsed ? 'Expand admin navigation' : 'Collapse admin navigation'}
            aria-expanded={!collapsed}
            title={collapsed ? 'Expand admin navigation' : 'Collapse admin navigation'}
            className="flex h-8 w-8 shrink-0 items-center justify-center text-mathua-secondary hover:text-mathua-blue"
          >
            {collapsed ? <ChevronsRight size={15} aria-hidden="true" /> : <ChevronsLeft size={15} aria-hidden="true" />}
          </button>
        )}
      </div>

      <nav aria-label="Admin sections" className="min-h-0 flex-1 overflow-y-auto px-3 py-4">
        {loading ? (
          <p className="px-2 font-mono text-[10px] uppercase tracking-wider text-mathua-muted">Loading access</p>
        ) : items.length === 0 ? (
          <p className="px-2 font-mono text-[11px] text-mathua-muted">No staff sections for this account.</p>
        ) : (
          <div className="space-y-5">
            {ADMIN_GROUPS.map((group) => {
              const groupItems = items.filter((item) => item.group === group)
              if (groupItems.length === 0) return null
              return (
                <div key={group}>
                  {!collapsed && <p className="px-3 pb-1 font-mono text-[9px] uppercase tracking-[0.18em] text-mathua-muted">{group}</p>}
                  <ul className="space-y-1">
                    {groupItems.map((item) => {
                      const active = isAdminNavActive(pathname ?? '', item.href)
                      const Icon = item.icon
                      return (
                        <li key={item.href}>
                          <Link
                            href={item.href}
                            onClick={onNavigate}
                            aria-current={active ? 'page' : undefined}
                            title={collapsed ? item.label : undefined}
                            className={
                              active
                                ? `flex h-10 items-center gap-2.5 bg-mathua-blue-faint font-mono text-xs text-mathua-blue ${collapsed ? 'justify-center px-0' : 'px-3'}`
                                : `flex h-10 items-center gap-2.5 font-mono text-xs text-mathua-secondary hover:bg-mathua-surface-elevated hover:text-mathua-primary ${collapsed ? 'justify-center px-0' : 'px-3'}`
                            }
                          >
                            <Icon size={15} aria-hidden="true" />
                            {!collapsed && <span className="truncate">{item.label}</span>}
                          </Link>
                        </li>
                      )
                    })}
                  </ul>
                </div>
              )
            })}
          </div>
        )}
      </nav>

      <div className="border-t border-mathua-border p-3">
        <Link
          href="/"
          onClick={onNavigate}
          className={`flex h-10 items-center gap-2.5 font-mono text-xs text-mathua-secondary hover:text-mathua-blue ${collapsed ? 'justify-center px-0' : 'px-3'}`}
        >
          <House size={15} aria-hidden="true" />
          {!collapsed && <span>Learner site</span>}
        </Link>
        {mounted && (
          <button
            type="button"
            onClick={toggleTheme}
            aria-label="Toggle theme"
            className={`flex h-10 w-full items-center gap-2.5 font-mono text-xs text-mathua-secondary hover:text-mathua-blue ${collapsed ? 'justify-center px-0' : 'px-3'}`}
          >
            <SunMoon size={15} isAnimated={false} aria-hidden="true" />
            {!collapsed && <span>Toggle theme</span>}
          </button>
        )}
        {me && !collapsed && (
          <p className="mt-2 truncate px-3 font-mono text-[10px] text-mathua-muted" title={me.id}>
            {me.role} · {me.id.slice(0, 8)}
          </p>
        )}
      </div>
    </>
  )
}

export default function AdminShell({ children }: { children: React.ReactNode }) {
  const pathname = usePathname()
  const { can, loading, me } = useAdminMe()
  const { mounted, toggleTheme } = useTheme()
  const [navOpen, setNavOpen] = useState(false)
  const [collapsed, setCollapsed] = useState(false)

  useEffect(() => {
    try {
      setCollapsed(localStorage.getItem('mathua_admin_sidebar') === 'collapsed')
    } catch {
      // Storage can be disabled; expanded is the safe default.
    }
  }, [])

  const toggleCollapsed = () => {
    setCollapsed((value) => {
      const next = !value
      try {
        localStorage.setItem('mathua_admin_sidebar', next ? 'collapsed' : 'expanded')
      } catch {
        // The visual state still changes for this session.
      }
      return next
    })
  }

  const items = loading ? [] : ADMIN_ITEMS.filter((item) => can(item.perm))

  return (
    <div className="min-h-screen bg-mathua-bg text-mathua-primary">
      <aside className={`fixed inset-y-0 left-0 z-30 hidden flex-col border-r border-mathua-border bg-mathua-surface transition-[width] duration-150 ease-out motion-reduce:transition-none md:flex ${collapsed ? 'w-14' : 'w-60'}`}>
        <AdminSidebar
          items={items}
          pathname={pathname}
          loading={loading}
          me={me}
          mounted={mounted}
          toggleTheme={toggleTheme}
          collapsed={collapsed}
          onToggle={toggleCollapsed}
        />
      </aside>

      <div className={`transition-[padding] duration-150 ease-out motion-reduce:transition-none ${collapsed ? 'md:pl-14' : 'md:pl-60'}`}>
        <header className="sticky top-0 z-20 border-b border-mathua-border bg-mathua-surface md:hidden">
          <div className="flex h-14 items-center justify-between px-4">
            <div className="flex items-center gap-3">
              <Sheet open={navOpen} onOpenChange={setNavOpen}>
                <SheetTrigger asChild>
                  <button
                    type="button"
                    aria-label="Open admin navigation"
                    className="flex h-9 items-center gap-1.5 border border-mathua-border bg-mathua-surface px-2.5 font-mono text-[11px] text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue"
                  >
                    <PanelLeft size={15} aria-hidden="true" />
                    Menu
                  </button>
                </SheetTrigger>
                <SheetContent side="left" className="w-60 border-r border-mathua-border bg-mathua-surface p-0">
                  <SheetTitle className="sr-only">Admin navigation</SheetTitle>
                  <SheetClose
                    aria-label="Close admin navigation"
                    className="absolute right-2 top-2 z-10 flex h-9 w-9 items-center justify-center border border-mathua-border bg-mathua-surface text-mathua-muted hover:border-mathua-blue hover:text-mathua-blue"
                  >
                    <Close size={16} aria-hidden="true" />
                  </SheetClose>
                  <div className="flex h-full flex-col">
                    <AdminSidebar
                      items={items}
                      pathname={pathname}
                      loading={loading}
                      me={me}
                      mounted={mounted}
                      toggleTheme={toggleTheme}
                      collapsed={false}
                      onNavigate={() => setNavOpen(false)}
                    />
                  </div>
                </SheetContent>
              </Sheet>
              <Link href="/admin" className="flex items-center gap-2">
                <span aria-hidden="true" className="font-mono text-base text-mathua-blue">λ</span>
                <span className="font-mono text-[11px] uppercase tracking-wider text-mathua-primary">Mathua Admin</span>
              </Link>
            </div>
            <Link href="/" className="flex h-11 items-center gap-1.5 font-mono text-[11px] text-mathua-secondary">
              <House size={14} aria-hidden="true" />
              Site
            </Link>
          </div>
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
