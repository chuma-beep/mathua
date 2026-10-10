'use client'

import Link from 'next/link'
import { useEffect, useLayoutEffect, useState } from 'react'
import { DICEBEAR_VERSION } from '../lib/dicebear'

interface FooterProps {
  className?: string
}

// The pool the footer draws its backdrop from.
//
// Real artwork, and an <img> rather than a CSS background. That is the whole
// technique: `width: 100%; height: auto` lets the picture set its own height from
// its bytes, so nothing is cropped, nothing is stretched, and nothing has to know
// the aspect ratio in advance. A background image cannot do this -- it has no
// intrinsic size, so the box must be given a height, and the two available
// choices are `cover`, which crops, and a percentage, which distorts.
//
// Each entry carries two pictures, because a landscape print shown across a
// phone is a letterbox slot 90px tall, and a portrait print shown across a
// desktop is a tall thin strip with dead space either side. `<picture>` with a
// `max-width: 639px` source swaps between them. The pairing is arbitrary and can
// be reshuffled freely; what matters is that both halves are usable at the
// breakpoint that selects them.
//
// `tone` is a property of the picture, not of the page. Every word of the footer
// is drawn on top of the artwork in both themes, so the ink has to suit the
// artwork -- a pale print takes dark text, a deep one takes light -- and the page
// theme has no say in it. It is declared per variant because the two pictures of
// one entry can be opposites: the coastal town is cream, the blue field is not.
//
// Desktop ratios are 16:9 by design so the footer's height does not change from
// visit to visit. Mobile ratios deliberately are not -- the portrait set ranges
// from 9:16 to 2:3, and normalising it would mean cropping the two outliers.
export interface FooterArtVariant {
  file: string
  width: number
  height: number
  tone: 'light' | 'dark'
}

export interface FooterArt {
  desktop: FooterArtVariant
  mobile: FooterArtVariant
}

export const FOOTER_ART: FooterArt[] = [
  {
    desktop: { file: '/footer/coastal-town.jpg', width: 1672, height: 643, tone: 'light' },
    mobile: { file: '/footer/mobile-reader.jpg', width: 1080, height: 1920, tone: 'light' },
  },
  {
    desktop: { file: '/footer/coastal-town.jpg', width: 1672, height: 643, tone: 'light' },
    mobile: { file: '/footer/mobile-coast.jpg', width: 675, height: 1200, tone: 'light' },
  },
  {
    desktop: { file: '/footer/grey-mountains.jpg', width: 1920, height: 738, tone: 'light' },
    mobile: { file: '/footer/mobile-field.jpg', width: 638, height: 1140, tone: 'dark' },
  },
  {
    desktop: { file: '/footer/coastal-town.jpg', width: 1672, height: 643, tone: 'light' },
    mobile: { file: '/footer/mobile-lighthouse.jpg', width: 816, height: 1456, tone: 'dark' },
  },
  {
    desktop: { file: '/footer/grey-mountains.jpg', width: 1920, height: 738, tone: 'light' },
    mobile: { file: '/footer/mobile-clouds.jpg', width: 1290, height: 2580, tone: 'light' },
  },
  {
    desktop: { file: '/footer/grey-mountains.jpg', width: 1920, height: 738, tone: 'light' },
    mobile: { file: '/footer/mobile-library.jpg', width: 723, height: 1076, tone: 'dark' },
  },
]

const FOOTER_SHAPES = ['alpha', 'beta', 'gamma']
const shapeUrl = (seed: string) =>
  `https://api.dicebear.com/${DICEBEAR_VERSION}/shapes/svg?seed=${seed}&backgroundColor=1c3a5e`

// Every word here is drawn on the artwork, so the ink comes from the picture's
// tone rather than from the page theme.
const linkClass = 'block opacity-80 transition-opacity hover:opacity-100'
const headingClass = 'font-mono text-[10px] uppercase tracking-[0.2em] text-[var(--footer-ink-muted)]'

const CREDITS = ['MIT License', 'Avatars by DiceBear']

const useIsomorphicLayoutEffect = typeof window === 'undefined' ? useEffect : useLayoutEffect

// Picks one entry from FOOTER_ART, once per app load.
//
// A layout effect on the client, a passive one on the server, which has no layout
// to read. useLayoutEffect rather than useEffect so the swap lands before the
// browser paints and the footer never flashes the wrong artwork on the way to the
// right one. The pick happens in the effect and never during render:
// Math.random() in render would disagree with the server's markup and fail
// hydration.
function useFooterArt(): FooterArt {
  const [art, setArt] = useState<FooterArt>(FOOTER_ART[0])

  useIsomorphicLayoutEffect(() => {
    setArt(FOOTER_ART[Math.floor(Math.random() * FOOTER_ART.length)])
  }, [])

  return art
}

export default function Footer({ className = '' }: FooterProps) {
  const art = useFooterArt()

  return (
    <footer
      // Two tones, not one. The desktop and mobile pictures are chosen by the
      // same entry and can be opposites -- cream above, deep blue below -- so the
      // ink has to change with the breakpoint. globals.css reads these in the
      // matching media query.
      data-tone-desktop={art.desktop.tone}
      data-tone-mobile={art.mobile.tone}
      className={`footer relative mt-[10px] border-t border-mathua-border ${className}`}
    >
      {/* Artwork and text share one grid cell. The row is max(image, content):
          the picture sets the height when it is the taller of the two, and the
          content sets it when the text is -- so text can never be clipped or
          overflow invisibly, whatever the viewport or the type size. A
          position:absolute overlay would have no such floor. */}
      {/* aria-hidden sits on the <img>, not the <picture>: the container has no
          ARIA role of its own, and the empty alt is what actually marks it
          decorative. */}
      <picture className="footer-art pointer-events-none">
        <source
          media="(max-width: 639px)"
          srcSet={art.mobile.file}
          width={art.mobile.width}
          height={art.mobile.height}
        />
        <img
          src={art.desktop.file}
          alt=""
          aria-hidden="true"
          width={art.desktop.width}
          height={art.desktop.height}
          loading="lazy"
          decoding="async"
          className="footer-art-img"
        />
      </picture>

      {/* Light for dark ink, dark for light ink. Only the top and bottom bands
          get it: every word of the footer sits in one of those two, so the middle
          of the print -- the part with no text over it -- is left alone. A scrim
          across the whole picture would tint the artwork, which is the thing
          being avoided by not filtering it. */}
      <div aria-hidden="true" className="footer-scrim pointer-events-none" />

      <div className="footer-content">
        <div className="footer-content-row">
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
            <p className="mt-3 text-sm leading-relaxed text-[var(--footer-ink-muted)]">
              Open-source adaptive math learning engine.
            </p>
          </div>

          <div className="grid grid-cols-2 gap-x-10 gap-y-8 sm:gap-x-16">
            <nav aria-labelledby="footer-learn">
              <h2 id="footer-learn" className={headingClass}>Learn</h2>
              <ul className="mt-3 space-y-2 font-mono text-xs">
                <li><Link href="/domains" className={linkClass}>Domains</Link></li>
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

        <div className="footer-credits border-t pt-5 font-mono text-[10px] opacity-80">
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
