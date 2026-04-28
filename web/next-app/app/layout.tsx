import type { Metadata } from 'next'
import './globals.css'

export const metadata: Metadata = {
  title: 'Mathua — Math Understanding Agent',
  description: 'An open-source adaptive math learning engine. Master the foundation. Earn the abstraction.',
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en" className="dark">
      <body>{children}</body>
    </html>
  )
}