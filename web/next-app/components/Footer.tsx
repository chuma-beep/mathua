import Link from 'next/link'
import { DICEBEAR_VERSION } from '../lib/dicebear'

interface FooterProps {
  className?: string
}

// Page-colored footer: blends with the background in both themes, separated
// by a hairline. Brand left, link columns grouped right, credits in a bottom
// bar. Full-width content.
const FOOTER_SHAPES = ['alpha', 'beta', 'gamma']
const shapeUrl = (seed: string) =>
  `https://api.dicebear.com/${DICEBEAR_VERSION}/shapes/svg?seed=${seed}&backgroundColor=1c3a5e`

const linkClass = 'block opacity-70 transition-opacity hover:opacity-100'
const headingClass = 'font-mono text-[10px] uppercase tracking-[0.2em] text-mathua-muted'

const CREDITS = ['MIT License', 'Avatars by DiceBear']

export default function Footer({ className = '' }: FooterProps) {
  return (
    <footer className={`mt-[10px] border-t border-mathua-border bg-mathua-bg text-mathua-primary ${className}`}>
      <div className="w-full px-4 py-10 sm:px-6 md:py-14 lg:px-10">
        <div className="flex flex-col gap-10 md:flex-row md:items-start md:justify-between md:gap-12">
          <div className="md:max-w-[34ch]">
            <div>
              <div className="font-sans text-lg">λ Mathua</div>
              <div className="mt-2 flex gap-2 md:mt-3">
                {FOOTER_SHAPES.map(seed => (
                  <img
                    key={seed}
                    src={shapeUrl(seed)}
                    alt=""
                    width={24}
                    height={24}
                    className="size-6 ring-1 ring-white/10 md:size-8"
                    loading="lazy"
                    referrerPolicy="no-referrer"
                  />
                ))}
              </div>
            </div>
            <p className="mt-3 text-sm leading-relaxed text-mathua-muted">
              Open-source adaptive math learning engine.
            </p>
          </div>

          <div className="grid grid-cols-2 gap-x-10 gap-y-8 sm:gap-x-16">
            <nav aria-labelledby="footer-learn">
              <h2 id="footer-learn" className={headingClass}>Learn</h2>
              <ul className="mt-3 space-y-2 font-mono text-xs">
                <li><Link href="/study" className={linkClass}>Study</Link></li>
                <li><Link href="/leaderboard" className={linkClass}>Leaderboard</Link></li>
                <li><Link href="/graph" className={linkClass}>Concept graph</Link></li>
              </ul>
            </nav>

            <nav aria-labelledby="footer-resources">
              <h2 id="footer-resources" className={headingClass}>Resources</h2>
              <ul className="mt-3 space-y-2 font-mono text-xs">
                <li><Link href="/docs" className={linkClass}>Docs</Link></li>
                <li><Link href="/docs/contributing" className={linkClass}>Contributing</Link></li>
                <li><Link href="/note" className={linkClass}>Creator&apos;s Note</Link></li>
                <li>
                  <a
                    href="https://github.com/chuma-beep/mathua"
                    target="_blank"
                    rel="noreferrer"
                    className={linkClass}
                  >
                    GitHub
                  </a>
                </li>
              </ul>
            </nav>
          </div>
        </div>

        <div className="mt-10 flex flex-wrap items-center gap-x-3 gap-y-1 border-t border-mathua-border pt-5 font-mono text-[10px] opacity-40">
          {CREDITS.map(credit => (
            <span key={credit} className="after:ml-3 after:content-['·'] last:after:content-none">
              {credit}
            </span>
          ))}
        </div>
      </div>
    </footer>
  )
}
