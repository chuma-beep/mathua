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
          href="https://fonts.googleapis.com/css2?family=Inter:ital,opsz,wght@0,14..32,400;0,14..32,600&family=Lora:ital,wght@0,400;0,500;1,400&family=JetBrains+Mono:ital,wght@0,400;0,500;1,400&display=swap"
          rel="stylesheet"
        />
      </head>
      <body>{children}</body>
    </html>
  )
}
