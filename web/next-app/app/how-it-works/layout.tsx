import type { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'How Mathua Works \u00b7 Concept graph, scheduler, scoring',
  description: 'A technical explanation of the concept graph, student model, diagnostic algorithm, task selection, and scoring system.',
}

export default function HowItWorksLayout({ children }: { children: React.ReactNode }) {
  return children
}
