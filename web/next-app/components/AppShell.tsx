'use client'

import Header from './Header'
import BottomTabs from './BottomTabs'

export default function AppShell({ children }: { children: React.ReactNode }) {
  return (
    <>
      <Header />
      <div className="max-w-container mx-auto px-6 max-sm:px-4 pb-[80px] lg:pb-0">{children}</div>
      <BottomTabs />
    </>
  )
}
