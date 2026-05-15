import type { Metadata } from 'next'
import Script from 'next/script'
import { IBM_Plex_Mono, IBM_Plex_Serif } from 'next/font/google'
import './globals.css'

const ibmPlexMono = IBM_Plex_Mono({
  subsets: ['latin'],
  weight: ['400', '500'],
  style: ['normal', 'italic'],
  display: 'swap',
  variable: '--font-ibm-plex-mono',
})

const ibmPlexSerif = IBM_Plex_Serif({
  subsets: ['latin'],
  weight: ['400', '500'],
  style: ['normal', 'italic'],
  display: 'swap',
  variable: '--font-ibm-plex-serif',
})

export const metadata: Metadata = {
  title: 'Mathua | Math Understanding Agent',
  description: 'An open-source adaptive math learning engine. Master the foundation. Earn the abstraction.',
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en" className="dark" suppressHydrationWarning>
      <head>
        <meta name="color-scheme" content="dark light" />
      </head>
      <body className={`${ibmPlexMono.variable} ${ibmPlexSerif.variable}`}>
        <Script
          src="/js/theme-init.js"
          strategy="beforeInteractive"
        />
        <div
          aria-hidden="true"
          className="pointer-events-none fixed inset-0 opacity-[0.25]"
          style={{
            backgroundImage:
              "linear-gradient(var(--border-strong) 1px, transparent 1px), linear-gradient(90deg, var(--border-strong) 1px, transparent 1px)",
            backgroundSize: "48px 48px",
            maskImage: "radial-gradient(ellipse at center, black 40%, transparent 78%)",
          }}
        />
        <div className="relative z-10">{children}</div>
      </body>
    </html>
  )
}
