import type { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'Architecture \u00b7 Mathua Docs',
  description: 'How the engine, scheduler, generators, graders, and storage layers compose into a single Go binary.',
}

export default function ArchitectureLayout({ children }: { children: React.ReactNode }) {
  return children
}
