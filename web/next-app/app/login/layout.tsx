import type { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'Login \u2014 Mathua',
  description: 'Sign in or create an account to track your progress across the concept map.',
}

export default function LoginLayout({ children }: { children: React.ReactNode }) {
  return children
}
