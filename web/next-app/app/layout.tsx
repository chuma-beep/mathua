import type { Metadata, Viewport } from 'next'
import Script from 'next/script'
import localFont from 'next/font/local'
import { Toaster } from '@/components/ui/sonner'
import ServiceWorkerRegistrar from '@/components/ServiceWorkerRegistrar'
import ChromeToggle from '@/components/ChromeToggle'
import { SITE_URL } from '@/lib/site'
// MathLive's own stylesheet for the <math-field> rendering (MathML). Loaded here,
// not in the lazily-loaded component, so the chunk has no CSS side-effect.
import 'mathlive/static.css'
import './globals.css'

// Self-hosted variable fonts (vendored from Google Fonts 2026-09-28:
// JetBrains Mono v24 + Space Grotesk v22, latin subsets). next/font/google
// fetches at build time and fails the build without network — vendoring
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

// Space Grotesk ships no italic face — only normal weights are requested.
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
  description: 'An open-source adaptive math learning platform. Mastery-gated, 657 concepts, generated problems, 50 XP quizzes.',
  manifest: '/manifest.webmanifest',
  openGraph: {
    title: 'Mathua',
    description: 'Master prerequisites before you advance. Open-source, 657 concepts, spaced repetition, 50 XP mastery checks.',
    type: 'website',
    siteName: 'Mathua',
    images: [{ url: '/og.png', width: 1200, height: 630, alt: 'Mathua concept graph' }],
  },
  twitter: {
    card: 'summary_large_image',
    title: 'Mathua',
    description: 'Master prerequisites before you advance. Open-source, 657 concepts, spaced repetition.',
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
        <div
          id="main-content"
          tabIndex={-1}
          // Scroll anchoring is off, deliberately. The chrome collapses the
          // header and its --chrome-top padding, which changes the height of the
          // document above the viewport. With anchoring on, the browser
          // compensates and reports that compensation as a scroll event -- and
          // since it moves in the same direction as the chrome just went, the
          // scroll logic read it as the user scrolling and undid the hide. The
          // content shifting when a toolbar collapses is the expected behaviour;
          // silently scrolling to compensate for it is not.
          className="relative z-10 [overflow-anchor:none]"
        >
          {children}
        </div>
        <Toaster />
        <ServiceWorkerRegistrar />
        {/* The chrome tap control, mounted here rather than inside Header: the
            header is `inert` while hidden, which would disable the one control
            meant to bring it back. It renders nothing on pages with no chrome
            (mobile /profile). */}
        <ChromeToggle />
      </body>
    </html>
  )
}
