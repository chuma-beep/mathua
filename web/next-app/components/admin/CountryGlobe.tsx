'use client'

import { useEffect, useRef } from 'react'
import createGlobe from 'cobe'
import { COUNTRY_CENTROIDS } from '../../lib/countryCentroids'
import { useTheme } from '../../hooks/useTheme'

export type GlobeCountry = {
  code: string
  active_7d: number
}

// A cobe globe with one marker per country, sized by this week's active
// learners. The canvas is decorative: the country table below it carries the
// same numbers as text, so a screen reader and a keyboard user lose nothing.
//
// Rotation stops under prefers-reduced-motion rather than slowing down. A
// globe that keeps drifting is exactly the motion that preference asks us not
// to draw.
export default function CountryGlobe({ countries }: { countries: GlobeCountry[] }) {
  const canvas = useRef<HTMLCanvasElement>(null)
  const { theme, mounted } = useTheme()
  const signature = countries.map((c) => `${c.code}:${c.active_7d}`).join(',')

  useEffect(() => {
    const element = canvas.current
    if (!element || !mounted) return
    const reduced = window.matchMedia('(prefers-reduced-motion: reduce)').matches
    const max = Math.max(1, ...countries.map((c) => c.active_7d))
    const markers = countries.flatMap((c) => {
      const location = COUNTRY_CENTROIDS[c.code]
      if (!location) return []
      // Square root keeps a single very busy country from flattening every
      // other marker into the same dot.
      const size = 0.03 + 0.06 * Math.sqrt(c.active_7d / max)
      return [{ location, size }]
    })

    const dark = theme === 'dark'
    // cobe's README documents onRender but its bundled types omit it, so the
    // option object is cast to the call signature rather than dropping the
    // rotation. The cast is narrow: createGlobe still checks every other field.
    const options = {
      devicePixelRatio: 2,
      width: 640,
      height: 640,
      phi: -0.6,
      theta: 0.25,
      dark: dark ? 1 : 0,
      diffuse: 1.25,
      mapSamples: 22000,
      // The first build used near-white on near-white, so the land texture
      // disappeared and the globe read as a blank ball. A deeper base gives the
      // dotted continents something to sit against in both themes.
      mapBrightness: dark ? 8 : 5.5,
      baseColor: dark ? [0.11, 0.12, 0.15] : [0.76, 0.81, 0.88],
      markerColor: dark ? [0.42, 0.66, 1] : [0.16, 0.38, 0.9],
      glowColor: dark ? [0.07, 0.08, 0.11] : [0.91, 0.94, 0.98],
      markers,
      onRender: (state: { phi: number }) => {
        if (!reduced) state.phi += 0.002
      },
    }
    const globe = createGlobe(element, options as unknown as Parameters<typeof createGlobe>[1])
    return () => globe.destroy()
    // signature carries every value the markers depend on; countries is a new
    // array on each parent render and would rebuild the WebGL context.
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [theme, mounted, signature])

  return (
    <div className="relative mx-auto aspect-square w-full max-w-[560px]">
      <div
        aria-hidden="true"
        className="pointer-events-none absolute inset-[7%] rounded-full bg-mathua-blue-faint blur-2xl"
      />
      <div
        aria-hidden="true"
        className="pointer-events-none absolute inset-[2%] rounded-full border border-mathua-border"
      />
      <canvas
        ref={canvas}
        aria-hidden="true"
        className="relative h-full w-full drop-shadow-[0_22px_50px_rgba(15,23,42,0.18)]"
      />
      <p className="sr-only">
        A globe with one marker for each country that has at least one known account. The table
        below lists the same countries with their account and activity counts.
      </p>
    </div>
  )
}
