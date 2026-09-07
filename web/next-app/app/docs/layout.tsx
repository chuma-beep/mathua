'use client'

import Header from '../../components/Header'
import BottomTabs from '../../components/BottomTabs'

export default function DocsLayout({ children }: { children: React.ReactNode }) {
  return (
    <>
      <Header links={[
        { label: 'Study', href: '/study' },
        { label: 'Docs', href: '/docs' },
        { label: 'Architecture', href: '/docs/architecture' },
        { label: 'Contributing', href: '/docs/contributing' },
        { label: 'Note', href: '/note' },
      ]} />
      <main>{children}</main>
      <footer style={{
        borderTop: '0.5px solid var(--border)',
        padding: '32px 24px',
        textAlign: 'center',
        fontFamily: "'IBM Plex Mono', monospace",
        fontSize: '12px',
        color: 'var(--text-muted)',
        lineHeight: 1.6,
      }}>
        Mathua · MIT License
      </footer>
      <BottomTabs />
    </>
  )
}
