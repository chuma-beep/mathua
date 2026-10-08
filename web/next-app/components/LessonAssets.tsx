'use client'

import type { LessonAsset } from '../lib/api'
import LessonDiagram from './LessonDiagram'
import { resolveLessonImageUrl } from './KatexContent'

interface LessonAssetsProps {
  assets?: LessonAsset[]
  /** Taller for a full-page reference, shorter inside a Learn card. */
  maxHeight?: number
}

/**
 * A lesson's figures, rendered from the asset list rather than from the markdown.
 *
 * The lesson body remains canonical and the reference panel renders it verbatim — this exists for the
 * surface that shows only part of a lesson. `/learn` presents one knowledge point at a time,
 * and without this the figures in that knowledge point's section were dropped: the worked
 * example arrived as text with its diagram missing, and the concept-level diagram that the
 * API returns was discarded by the client entirely.
 *
 * Figures go through `LessonDiagram` rather than a bare `<img>` so they get the same
 * treatment as everywhere else: SVG inlined and themeable, PNGs optimized, `meta.json`
 * titles and CC BY-NC attribution applied, and a fetch failure falling back to a plain
 * image so a missing asset never blanks the card.
 */
export default function LessonAssets({ assets, maxHeight }: LessonAssetsProps) {
  const figures = (assets ?? []).filter((a) => a.kind === 'image' && a.src)
  if (figures.length === 0) return null
  return (
    <div className="mt-4 space-y-4">
      {figures.map((a) => (
        <figure key={a.id}>
          <LessonDiagram
            src={resolveLessonImageUrl(a.src)}
            alt={a.alt || a.id}
            maxHeight={maxHeight}
          />
          {/* Rendered only when the lesson body carries one. Measured at 0 across the
              current corpus — these figures caption themselves inside the SVG — so this is
              the path for authored assets, not the vendored ones. If a body ever gains
              markdown captions *and* the SVG already draws its own, the caption will read
              twice, and `TestLessonCaptionsDoNotDuplicate` is where that should be caught. */}
          {a.caption && (
            <figcaption className="mt-2 font-mono text-[11px] text-mathua-secondary">
              {a.caption}
            </figcaption>
          )}
        </figure>
      ))}
    </div>
  )
}