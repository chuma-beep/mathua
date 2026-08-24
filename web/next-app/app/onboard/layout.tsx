import type { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'Onboarding \u00b7 Mathua',
  description: 'Select what you want to learn and let Mathua build a personalized study plan.',
}

export default function OnboardLayout({ children }: { children: React.ReactNode }) {
  return children
}
