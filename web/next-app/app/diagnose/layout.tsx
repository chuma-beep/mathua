import type { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'Diagnostic Test \u2014 Mathua',
  description: 'A Computerised Adaptive Testing session that locates your position in the concept graph.',
}

export default function DiagnoseLayout({ children }: { children: React.ReactNode }) {
  return children
}
