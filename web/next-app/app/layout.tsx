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
    <html lang="en" className="dark" suppressHydrationWarning>
      <head>
        <meta name="color-scheme" content="dark light" />
        <script
          dangerouslySetInnerHTML={{
            __html: `(function(){try{var t=localStorage.getItem('mathua-theme');var d;if(t==='dark'||t==='light'){d=t==='dark'}else{d=window.matchMedia('(prefers-color-scheme:dark)').matches}if(d){document.documentElement.classList.add('dark')}else{document.documentElement.classList.remove('dark')}}catch(e){}})()`,
          }}
        />
        <link rel="preconnect" href="https://fonts.googleapis.com" />
        <link rel="preconnect" href="https://fonts.gstatic.com" crossOrigin="anonymous" />
        <link
          href="https://fonts.googleapis.com/css2?family=IBM+Plex+Mono:ital,wght@0,400;0,500;1,400&family=IBM+Plex+Serif:ital,wght@0,400;0,500;1,400&display=swap"
          rel="stylesheet"
        />
      </head>
      <body>
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
