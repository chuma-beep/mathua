'use client'

import Link from 'next/link'
import { useEffect, useLayoutEffect, useState } from 'react'
import { DICEBEAR_VERSION } from '../lib/dicebear'
import { CREATOR_LINKS, SUPPORT_URL } from '../lib/creatorLinks'
import { SITE_URL } from '../lib/site'

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
// Desktop ratios are 2.6:1 by design so the footer's height does not change from
// visit to visit. At the original 16:9 it was a full screen of artwork -- 810px at
// 1440 -- and the empty sky above the ridgeline was most of that. Mobile ratios
// deliberately are not uniform: the portrait set ranges from 9:16 to 2:3, and
// normalising it would mean cropping the two outliers.
//
// Three desktop pictures since 2026-10-18: the harbour, and two blue-and-cream
// coastal engravings added at the owner's request. They are dealt out over the
// six entries, two each, so every print is equally likely and every entry keeps
// its own portrait -- a phone still rotates through six pictures.
//
// All three are pale, so `tone` stays light on desktop and the ink does not change
// with the pick. The two newer prints were measured over the nav band after the
// scrim and came out better than the harbour print did (muted ink, 5th percentile:
// 7.35 and 5.42 against 4.73), so the two bands above carry them without a corner
// layer.
//
// PROVENANCE: every file named below is shipped, and none of them has a recorded
// author or licence. See the "Footer artwork" entry in THIRD_PARTY_NOTICES.md. The
// contract `test/footerArt.test.tsx` enforces is about the files -- they exist,
// declare their real dimensions, and carry a tone -- so replacing this set with
// known-provenance artwork needs no change here.
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
    desktop: { file: '/footer/coastal-temple.jpg', width: 1672, height: 643, tone: 'light' },
    mobile: { file: '/footer/mobile-field.jpg', width: 638, height: 1140, tone: 'dark' },
  },
  {
    desktop: { file: '/footer/coastal-temple.jpg', width: 1672, height: 643, tone: 'light' },
    mobile: { file: '/footer/mobile-lighthouse.jpg', width: 816, height: 1456, tone: 'dark' },
  },
  {
    desktop: { file: '/footer/coastal-villa.jpg', width: 1672, height: 643, tone: 'light' },
    mobile: { file: '/footer/mobile-clouds.jpg', width: 1290, height: 2580, tone: 'light' },
  },
  {
    desktop: { file: '/footer/coastal-villa.jpg', width: 1672, height: 643, tone: 'light' },
    mobile: { file: '/footer/mobile-library.jpg', width: 723, height: 1076, tone: 'dark' },
  },
]

const FOOTER_SHAPES = ['alpha', 'beta', 'gamma']
const shapeUrl = (seed: string) =>
  `https://api.dicebear.com/${DICEBEAR_VERSION}/shapes/svg?seed=${seed}&backgroundColor=1c3a5e`

// Every word here is drawn on the artwork, so the ink comes from the picture's
// tone rather than from the page theme.
// Full ink, not 80 percent. The mono face is 12px with hairline strokes, and
// fading it to 80 percent over an engraving is what made the links look washed
// out even where the scrim covered them. Hover adds an underline instead of
// fading the text, so the resting state keeps the full contrast. The floor fade
// below supplies the background; the links themselves carry no tint.
const linkClass = 'inline-block underline-offset-4 transition-colors hover:underline hover:decoration-1'
const headingClass = 'font-mono text-[10px] uppercase tracking-[0.2em] text-[var(--footer-ink-muted)]'

// The credits strip names each license the shipped content carries, not only the software.
// Lesson content is CC BY-NC 4.0 from Algebrica and CC BY-NC-SA 4.0 from OpenStax: both require
// attribution, and neither is MIT, so each gets its own line. The glossary draws on OpenStax,
// DLMF and others without reusing their text. A short line cannot carry title, author, source
// and license for every source, so the line names them and "Licenses" links to /licenses, which
// renders THIRD_PARTY_NOTICES.md in full.
const CREDITS = [
  'MIT License',
  'Avatars by DiceBear',
  'Lesson content: Algebrica (CC BY-NC 4.0)',
  'Lesson content: OpenStax (CC BY-NC-SA 4.0)',
  'Glossary: OpenStax, DLMF and others',
]

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
  const [shareStatus, setShareStatus] = useState<'idle' | 'copied' | 'error'>('idle')

  useEffect(() => {
    if (shareStatus === 'idle') return
    const timer = window.setTimeout(() => setShareStatus('idle'), 2200)
    return () => window.clearTimeout(timer)
  }, [shareStatus])

  async function shareMathua() {
    if (typeof navigator !== 'undefined' && typeof navigator.share === 'function') {
      try {
        await navigator.share({
          title: 'Mathua',
          text: 'Learn math with Mathua.',
          url: SITE_URL,
        })
        return
      } catch (error) {
        if (error instanceof DOMException && error.name === 'AbortError') return
      }
    }

    try {
      await navigator.clipboard.writeText(SITE_URL)
      setShareStatus('copied')
    } catch {
      setShareStatus('error')
    }
  }

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
            <a
              href={SUPPORT_URL}
              target="_blank"
              rel="noreferrer"
              className="mt-4 inline-block border border-current px-2.5 py-1 font-mono text-[11px] uppercase tracking-[0.12em] transition-colors hover:bg-white/10"
            >
              Support Mathua
            </a>
            <a
              href="https://github.com/chuma-beep/mathua"
              target="_blank"
              rel="noreferrer"
              className="mt-2 block font-mono text-[11px] uppercase tracking-[0.12em] underline underline-offset-4 transition-colors hover:no-underline"
            >
              Star Mathua on GitHub
            </a>
            <div className="mt-4 max-w-[250px]">
              <p className="text-sm leading-relaxed text-[var(--footer-ink-muted)]">
                Know someone who wants to learn? Share Mathua with them.
              </p>
              <button
                type="button"
                onClick={shareMathua}
                className="mt-2 inline-flex min-h-[32px] items-center justify-center bg-mathua-blue px-2.5 py-1 font-mono text-[11px] uppercase tracking-[0.12em] text-white transition-colors hover:bg-mathua-blue-hover"
              >
                {shareStatus === 'copied'
                  ? 'Link copied'
                  : shareStatus === 'error'
                    ? 'Copy failed'
                    : 'Share Mathua'}
              </button>
              <span className="sr-only" aria-live="polite">
                {shareStatus === 'copied'
                  ? 'Link copied'
                  : shareStatus === 'error'
                    ? 'Copy failed'
                    : ''}
              </span>
            </div>
          </div>

          <div className="grid grid-cols-2 gap-x-10 gap-y-8 sm:grid-cols-3 sm:gap-x-12">
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
                <li><Link href="/how-it-works" className={linkClass}>How it works</Link></li>
                <li><Link href="/docs" className={linkClass}>Docs</Link></li>
                <li><Link href="/docs/contributing" className={linkClass}>Contributing</Link></li>
                <li><Link href="/note" className={linkClass}>Creator&apos;s Note</Link></li>
                <li><Link href="/contact" className={linkClass}>Contact admins</Link></li>
                <li>
                  <a
                    href="https://discord.gg/qXPeQ2jys"
                    target="_blank"
                    rel="noreferrer"
                    className={linkClass}
                  >
                    Discord
                  </a>
                </li>
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

            <nav aria-labelledby="footer-elsewhere">
              <h2 id="footer-elsewhere" className={headingClass}>Elsewhere</h2>
              <ul className="mt-3 space-y-2 font-mono text-xs">
                {CREATOR_LINKS.map(link => (
                  <li key={link.href}>
                    <a
                      href={link.href}
                      target="_blank"
                      rel="noreferrer"
                      className={linkClass}
                    >
                      {link.label}
                    </a>
                  </li>
                ))}
              </ul>
            </nav>
          </div>
        </div>

        <div className="footer-credits border-t pt-5 font-mono text-[10px] opacity-80">
          {CREDITS.map(credit => (
            <span key={credit} className="after:ml-3 after:content-['·']">
              {credit}
            </span>
          ))}
          <Link href="/licenses" className="underline hover:no-underline">
            Licenses
          </Link>
          {' · '}
          <Link href="/privacy" className="underline hover:no-underline">
            Privacy
          </Link>
        </div>
      </div>
    </footer>
  )
}
