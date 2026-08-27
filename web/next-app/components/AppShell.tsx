'use client'

import { useState, useEffect } from 'react'
import Header from './Header'
import Drawer from './Drawer'
import BottomTabs from './BottomTabs'

const STORAGE_KEY = 'mathua-drawer-open'

export default function AppShell({ children }: { children: React.ReactNode }) {
  const [drawerOpen, setDrawerOpen] = useState(true)
  const [mobileDrawer, setMobileDrawer] = useState(false)

  useEffect(() => {
    const stored = localStorage.getItem(STORAGE_KEY)
    if (stored === 'false') setDrawerOpen(false)
    if (stored === 'true') setDrawerOpen(true)
  }, [])

  const toggleDrawer = () => {
    const next = !drawerOpen
    setDrawerOpen(next)
    localStorage.setItem(STORAGE_KEY, String(next))
  }

  return (
    <>
      <Header drawerOpen={drawerOpen} onToggleDrawer={toggleDrawer} onOpenMobileDrawer={() => setMobileDrawer(true)} />
      <div className="flex max-w-container mx-auto">
        {/* Desktop drawer 240px */}
        <aside
          className={`${drawerOpen ? 'w-[240px]' : 'w-[48px]'} hidden lg:block shrink-0 sticky top-[52px] h-[calc(100vh-52px)] border-r-[0.5px] border-mathua-border bg-mathua-bg overflow-y-auto transition-all`}
        >
          <Drawer collapsed={!drawerOpen} onToggle={toggleDrawer} />
        </aside>
        <main className="flex-1 min-w-0 px-6 max-sm:px-4 pb-[56px] lg:pb-0">{children}</main>
      </div>
      <BottomTabs onMore={() => setMobileDrawer(true)} />
      {/* Mobile drawer sheet */}
      {mobileDrawer && (
        <>
          <div className="lg:hidden fixed inset-0 z-50 bg-black/40" onClick={() => setMobileDrawer(false)} />
          <div className="lg:hidden fixed left-0 top-0 bottom-0 w-[260px] z-50 bg-mathua-bg border-r border-mathua-border overflow-y-auto">
            <Drawer collapsed={false} onToggle={() => setMobileDrawer(false)} />
          </div>
        </>
      )}
    </>
  )
}
