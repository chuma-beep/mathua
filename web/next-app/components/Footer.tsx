import Link from 'next/link'
import { DICEBEAR_VERSION } from '../lib/dicebear'

interface FooterProps {
  className?: string
}

// Page-colored footer: blends with the background in both themes, separated
// by a hairline. Full-width content; compact two-column layout on phones.
const FOOTER_SHAPES = ['alpha', 'beta', 'gamma']
const shapeUrl = (seed: string) =>
  `https://api.dicebear.com/${DICEBEAR_VERSION}/shapes/svg?seed=${seed}&backgroundColor=1c3a5e`

const linkClass = 'block opacity-70 transition-opacity hover:opacity-100'

export default function Footer({ className = '' }: FooterProps) {
  return (
    <footer className={`mt-[10px] border-t border-mathua-border bg-mathua-bg text-mathua-primary ${className}`}>
      <div className="w-full px-4 py-10 sm:px-6 md:py-14 lg:px-10">
        <div className="grid grid-cols-2 gap-x-6 gap-y-8 md:grid-cols-4 md:gap-8 lg:gap-12">
          <div className="col-span-2 flex items-center gap-3 md:col-span-1 md:block">
            <div className="font-sans text-lg">λ Mathua</div>
            <div className="flex gap-2 md:mt-2">
              {FOOTER_SHAPES.map(seed => (
                <img
                  key={seed}
                  src={shapeUrl(seed)}
                  alt=""
                  width={32}
                  height={32}
                  className="size-8 ring-1 ring-white/10"
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

          <div className="col-span-2 flex flex-wrap gap-x-3 gap-y-1 font-mono text-[10px] opacity-40 md:col-span-1 md:block md:space-y-2 md:text-right">
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
