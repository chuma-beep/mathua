import Link from 'next/link'
import { DICEBEAR_VERSION } from '../lib/dicebear'

interface FooterProps {
  className?: string
}

// Prominent editorial footer: full-bleed ink block with the wordmark, the
// DiceBear shape marks, and the link columns. Rendered outside the page
// container so it spans the viewport.
const FOOTER_SHAPES = ['alpha', 'beta', 'gamma']
const shapeUrl = (seed: string) =>
  `https://api.dicebear.com/${DICEBEAR_VERSION}/shapes/svg?seed=${seed}&backgroundColor=1c3a5e`

const linkClass = 'block opacity-70 transition-opacity hover:opacity-100'

export default function Footer({ className = '' }: FooterProps) {
  return (
    <footer className={`bg-mathua-primary text-mathua-bg ${className}`}>
      <div className="mx-auto max-w-container px-4 py-14 sm:px-6">
        <div className="grid gap-8 md:grid-cols-4">
          <div>
            <div className="font-sans text-lg">λ Mathua</div>
            <div className="mt-2 flex gap-2">
              {FOOTER_SHAPES.map(seed => (
                <img
                  key={seed}
                  src={shapeUrl(seed)}
                  alt=""
                  width={32}
                  height={32}
                  className="size-8"
                  loading="lazy"
                  referrerPolicy="no-referrer"
                />
              ))}
            </div>
          </div>

          <nav className="space-y-2 font-mono text-xs" aria-label="Study">
            <Link href="/study" className={linkClass}>Study</Link>
            <Link href="/leaderboard" className={linkClass}>Leaderboard</Link>
            <Link href="/graph" className={linkClass}>Concept graph</Link>
          </nav>

          <nav className="space-y-2 font-mono text-xs" aria-label="Resources">
            <Link href="/docs" className={linkClass}>Docs</Link>
            <Link href="/docs/contributing" className={linkClass}>Contributing</Link>
            <Link href="/note" className={linkClass}>Creator&apos;s Note</Link>
            <a
              href="https://github.com/chuma-beep/mathua"
              target="_blank"
              rel="noreferrer"
              className={linkClass}
            >
              GitHub
            </a>
          </nav>

          <div className="space-y-2 font-mono text-[10px] opacity-40 md:text-right">
            <div>MIT License</div>
            <div>Avatars by DiceBear</div>
            <div>Draft · Sheet 01</div>
            <div>Scale 1:1 · Grid 22px</div>
          </div>
        </div>
      </div>
    </footer>
  )
}
