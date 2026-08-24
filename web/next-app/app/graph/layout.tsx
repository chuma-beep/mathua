import type { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'Knowledge Graph \u00b7 Mathua',
  description: 'Explore the Mathua concept map · view your progress across the full curriculum.',
}

export default function GraphLayout({ children }: { children: React.ReactNode }) {
  return children
}
