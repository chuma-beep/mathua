'use client'

import { useState, useEffect } from 'react'
import Header from './Header'
import Drawer from './Drawer'
import BottomTabs from './BottomTabs'

const STORAGE_KEY = 'mathua-drawer-open'

export default function AppShell({ children }: { children: React.ReactNode }) {
  return (
    <>
      <Header />
      <div className="max-w-container mx-auto px-6 max-sm:px-4 pb-[56px] lg:pb-0">{children}</div>
      <BottomTabs />
    </>
  )
}
