import type { Metadata, Viewport } from 'next'
import Script from 'next/script'
import { JetBrains_Mono, Space_Grotesk } from 'next/font/google'
import { Toaster } from '@/components/ui/sonner'
import ServiceWorkerRegistrar from '@/components/ServiceWorkerRegistrar'
import { SITE_URL } from '@/lib/site'
import './globals.css'

const jetbrainsMono = JetBrains_Mono({
  subsets: ['latin'],
  weight: ['400', '500'],
  style: ['normal', 'italic'],
  display: 'swap',
  // Only the faces a page actually renders are fetched; preloading the
  // families costs ~105KB on the landing page.
  preload: false,
  variable: '--font-jetbrains-mono',
})

// Space Grotesk ships no italic face — only normal weights are requested.
const spaceGrotesk = Space_Grotesk({
  subsets: ['latin'],
  weight: ['400', '500', '700'],
  display: 'swap',
  preload: false,
  variable: '--font-space-grotesk',
})

export const metadata: Metadata = {
  metadataBase: new URL(SITE_URL),
  title: 'Mathua',
  applicationName: 'Mathua',
  description: 'An open-source adaptive math learning platform. Mastery-gated, 634 concepts, generated problems, 150 XP quizzes.',
  manifest: '/manifest.webmanifest',
  openGraph: {
    title: 'Mathua',
    description: 'Master prerequisites before you advance. Open-source, 634 concepts, spaced repetition, 150 XP mastery checks.',
    type: 'website',
    siteName: 'Mathua',
    images: [{ url: '/og.png', width: 1200, height: 630, alt: 'Mathua concept graph' }],
  },
  twitter: {
    card: 'summary_large_image',
    title: 'Mathua',
    description: 'Master prerequisites before you advance. Open-source, 634 concepts, spaced repetition.',
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
        {/* KaTeX builds some glyphs (≠, ≤, overlays) from private-use codepoints
            that render as tofu until KaTeX_Main arrives. katex.min.css uses
            font-display:block with no preload, so slow networks show broken
            math. Preload the Main faces up front. Filenames are content hashes
            from the katex package — refresh these hrefs when katex upgrades
            (a stale href 404s and degrades gracefully to status quo). */}
        <link rel="preload" href="/_next/static/media/KaTeX_Main-Regular.0462f03b.woff2" as="font" type="font/woff2" crossOrigin="anonymous" />
        <link rel="preload" href="/_next/static/media/KaTeX_Main-Bold.c3fb5ac2.woff2" as="font" type="font/woff2" crossOrigin="anonymous" />
        <link rel="preload" href="/_next/static/media/KaTeX_Main-Italic.8916142b.woff2" as="font" type="font/woff2" crossOrigin="anonymous" />
        <link rel="preload" href="/_next/static/media/KaTeX_Main-BoldItalic.6f2bb1df.woff2" as="font" type="font/woff2" crossOrigin="anonymous" />
      </head>
      <body className={`${jetbrainsMono.variable} ${spaceGrotesk.variable}`}>
        <Script
          src="/js/theme-init.js"
          strategy="beforeInteractive"
        />
        <div id="main-content" tabIndex={-1} className="relative z-10">
          {children}
        </div>
        <Toaster />
        <ServiceWorkerRegistrar />
      </body>
    </html>
  )
}
