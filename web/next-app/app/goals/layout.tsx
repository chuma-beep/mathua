import type { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'Goals \u00b7 Mathua',
  description: 'Set learning goals and run a diagnostic to build a personalized study plan.',
}

export default function GoalsLayout({ children }: { children: React.ReactNode }) {
  return children
}
