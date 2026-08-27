'use client'

import { useState, useEffect } from 'react'
import Header from './Header'
import Drawer from './Drawer'
import BottomTabs from './BottomTabs'

const STORAGE_KEY = 'mathua-drawer-open'

export default function AppShell({ children }: { children: React.ReactNode }) {
  const [mobileDrawer, setMobileDrawer] = useState(false)

  return (
    <>
      <Header onOpenMobileDrawer={() => setMobileDrawer(true)} />
      <div className="max-w-container mx-auto px-6 max-sm:px-4 pb-[56px] lg:pb-0">{children}</div>
      <BottomTabs onMore={() => setMobileDrawer(true)} />
      {/* Mobile drawer sheet — only mobile/tablet */}
      {mobileDrawer && (
        <>
          <div className="fixed inset-0 z-50 bg-black/40 lg:hidden" onClick={() => setMobileDrawer(false)} />
          <div className="fixed left-0 top-0 bottom-0 w-[260px] z-50 bg-mathua-bg border-r border-mathua-border overflow-y-auto lg:hidden">
            <Drawer collapsed={false} onToggle={() => setMobileDrawer(false)} />
          </div>
        </>
      )}
    </>
  )
}
