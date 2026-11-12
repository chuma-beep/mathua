import type { Metadata, Viewport } from 'next'
import Script from 'next/script'
import localFont from 'next/font/local'
import { Toaster } from '@/components/ui/sonner'
import ServiceWorkerRegistrar from '@/components/ServiceWorkerRegistrar'
import ChromeToggle from '@/components/ChromeToggle'
import AssistantDock from '@/components/AssistantDock'
import ScrollToTopOnNavigate from '@/components/ScrollToTopOnNavigate'
import { SITE_URL } from '@/lib/site'
import { Analytics } from '@vercel/analytics/react'
// MathLive's own stylesheet for the <math-field> rendering (MathML). Loaded here,
// not in the lazily-loaded component, so the chunk has no CSS side-effect.
import 'mathlive/static.css'
import './globals.css'

// Self-hosted variable fonts (vendored from Google Fonts 2026-09-28:
// JetBrains Mono v24 + Space Grotesk v22, latin subsets). next/font/google
// fetches at build time and fails the build without network, vendoring
// keeps Vercel/Docker builds hermetic.
const jetbrainsMono = localFont({
  src: [
    { path: './fonts/jetbrains-mono-latin-normal.woff2', weight: '400 500', style: 'normal' },
    { path: './fonts/jetbrains-mono-latin-italic.woff2', weight: '400 500', style: 'italic' },
  ],
  display: 'swap',
  // Only the faces a page actually renders are fetched; preloading the
  // families costs ~105KB on the landing page.
  preload: false,
  variable: '--font-jetbrains-mono',
})

// Space Grotesk ships no italic face, only normal weights are requested.
const spaceGrotesk = localFont({
  src: [{ path: './fonts/space-grotesk-latin.woff2', weight: '400 500 700', style: 'normal' }],
  display: 'swap',
  preload: false,
  variable: '--font-space-grotesk',
})

export const metadata: Metadata = {
  metadataBase: new URL(SITE_URL),
  title: 'Mathua',
  applicationName: 'Mathua',
  description: 'An open-source adaptive math learning platform. Mastery-gated, 2500 concepts, generated problems, mastery checks.',
  manifest: '/manifest.webmanifest',
  openGraph: {
    title: 'Mathua',
    description: 'Master prerequisites before you advance. Open-source, 2500 concepts, spaced repetition, mastery checks.',
    type: 'website',
    siteName: 'Mathua',
    images: [{ url: '/og.png', width: 1200, height: 630, alt: 'Mathua: an adaptive math learning platform' }],
  },
  twitter: {
    card: 'summary_large_image',
    site: '@wisdomsglow',
    creator: '@wisdomsglow',
    title: 'Mathua',
    description: 'Master prerequisites before you advance. Open-source, 2500 concepts, spaced repetition.',
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
            from the katex package, refresh these hrefs when katex upgrades
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
        <div
          id="main-content"
          tabIndex={-1}
          // Scroll anchoring used to be off here, because the chrome collapsed the
          // header and its --chrome-top padding and the browser's compensation for
          // that shift arrived as a scroll event that undid the hide. The header
          // now slides on a transform and changes no layout, so there is nothing
          // left to compensate for, and the browser's default is the better one:
          // content that loads above the viewport should not jerk the page.
          className="relative z-10"
        >
          {children}
        </div>
        <Toaster />
        {/* Vercel serves the analytics script at /_vercel/insights. The static
            export that the e2e suite serves has no such route, and its 404 is a
            console error, which the suite correctly fails on. */}
        {process.env.VERCEL === '1' && <Analytics />}
        <ServiceWorkerRegistrar />
        {/* The chrome tap control, mounted here rather than inside Header: the
            header is `inert` while hidden, which would disable the one control
            meant to bring it back. It renders nothing on pages with no chrome
            (mobile /profile). */}
        <ChromeToggle />
        {/* Mounted once here rather than per page so a conversation survives
            navigation; it renders nothing when the server reports AI is off. */}
        <AssistantDock />
        {/* A route change starts at the top, including on iOS Safari, where the
            browser can otherwise leave the old offset in place until the new
            page has laid out. Back and forward keep their restored position. */}
        <ScrollToTopOnNavigate />
      </body>
    </html>
  )
}
