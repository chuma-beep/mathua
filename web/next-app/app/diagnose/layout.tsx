import type { Metadata } from 'next'
import BottomTabs from '../../components/BottomTabs'

export const metadata: Metadata = {
  title: 'Diagnostic Test \u00b7 Mathua',
  description: 'A Computerised Adaptive Testing session that locates your position in the concept graph.',
}

export default function DiagnoseLayout({ children }: { children: React.ReactNode }) {
  return (
    <>
      {children}
      <BottomTabs />
    </>
  )
}
