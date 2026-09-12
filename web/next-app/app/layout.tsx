import type { Metadata, Viewport } from 'next'
import Script from 'next/script'
import { IBM_Plex_Mono, IBM_Plex_Serif } from 'next/font/google'
import { Toaster } from '@/components/ui/sonner'
import ServiceWorkerRegistrar from '@/components/ServiceWorkerRegistrar'
import { SITE_URL } from '@/lib/site'
import './globals.css'

const ibmPlexMono = IBM_Plex_Mono({
  subsets: ['latin'],
  weight: ['400', '500'],
  style: ['normal', 'italic'],
  display: 'swap',
  // Only the faces a page actually renders are fetched; preloading all eight
  // (both families x 2 weights x 2 styles) cost ~105KB on the landing page.
  preload: false,
  variable: '--font-ibm-plex-mono',
})

const ibmPlexSerif = IBM_Plex_Serif({
  subsets: ['latin'],
  weight: ['400', '500'],
  style: ['normal', 'italic'],
  display: 'swap',
  preload: false,
  variable: '--font-ibm-plex-serif',
})

export const metadata: Metadata = {
  metadataBase: new URL(SITE_URL),
  title: {
    default: 'Mathua — Adaptive math learning platform',
    template: '%s · Mathua',
  },
  description: 'An open-source adaptive math learning platform. Mastery-gated, 630 concepts, generated problems, 150 XP quizzes.',
  manifest: '/manifest.webmanifest',
  openGraph: {
    title: 'Mathua — Adaptive math learning platform',
    description: 'Master prerequisites before you advance. Open-source, 630 concepts, spaced repetition, 150 XP mastery checks.',
    type: 'website',
    siteName: 'Mathua',
    images: [{ url: '/og.png', width: 1200, height: 630, alt: 'Mathua concept graph' }],
  },
  twitter: {
    card: 'summary_large_image',
    title: 'Mathua — Adaptive math learning platform',
    description: 'Master prerequisites before you advance. Open-source, 630 concepts, spaced repetition.',
    images: ['/og.png'],
  },
  icons: {
    icon: [
      { url: '/icons/icon.svg', type: 'image/svg+xml' },
      { url: '/icons/icon-192.png', sizes: '192x192', type: 'image/png' },
      { url: '/icons/icon-512.png', sizes: '512x512', type: 'image/png' },
    ],
    apple: [{ url: '/icons/icon-192.png', sizes: '192x192' }],
  },
}

export const viewport: Viewport = {
  themeColor: '#0c0d0f',
  width: 'device-width',
  initialScale: 1,
  viewportFit: 'cover',
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
        <a href="#main-content" className="skip-link">
          Skip to content
        </a>
        <div id="main-content" tabIndex={-1} className="relative z-10">
          {children}
        </div>
        <Toaster />
        <ServiceWorkerRegistrar />
      </body>
    </html>
  )
}
