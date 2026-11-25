'use client'

import type { ReactNode } from 'react'
import { PredictBlock, usePredict } from './predict'
import type { ReadoutItem } from './types'

interface FigureFrameProps {
  title: string
  caption?: string
  readout: ReadoutItem[]
  /** One line telling the reader how the drawing responds. Kept out of the SVG:
   *  text inside the drawing scales with it and collided with the plot. */
  hint?: string
  onReset?: () => void
  /** Control row rendered outside the bordered figure, so the drawing box stays
   *  a drawing and the buttons can take a full touch target on a phone. */
  controls?: ReactNode
  children: ReactNode
}

/**
 * The chrome every template shares: a title, the drawing, a live readout, a
 * caption and a reset. Templates draw only the SVG and controls; keeping the
 * frame here is what makes six figures look like one system.
 *
 * The readout is a <dl aria-live="polite">, so a learner who changes a parameter
 * with the keyboard hears the new value. It is never the only place a value
 * appears: every control carries its own visible output.
 */
export default function FigureFrame({ title, caption, readout, hint, onReset, controls, children }: FigureFrameProps) {
  const predict = usePredict()
  return (
    <>
    <figure className="mt-5 border border-mathua-border bg-mathua-bg">
      <div className="flex items-center justify-between gap-3 border-b border-mathua-border px-3 py-2">
        {/* No `uppercase` on a title or a label: CSS case folding rewrites the mathematics.
            The readout once showed "P(X ≤ X)" for P(X ≤ x), "Μ ± Σ" for μ ± σ and "ΔΧ"
            for Δx, because the Greek lower-case letters have upper-case codepoints. */}
        <p className="font-mono text-[11px] tracking-wider text-mathua-muted">{title}</p>
        {onReset && (
          <button
            type="button"
            onClick={onReset}
            className="inline-flex h-9 items-center border border-mathua-border px-3 font-mono text-xs text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue"
          >
            Reset
          </button>
        )}
      </div>
      <div className="p-3">
        {/* Drawing and controls share one measure. An SVG capped at 420px with controls stretched
            across a 700px card reads as two unrelated columns, which is what the complex figure
            looked like before this wrapper existed. */}
        <div className="mx-auto w-full max-w-[420px]">
          {hint && <p className="mb-2 font-mono text-[10px] text-mathua-muted">{hint}</p>}
          {predict && <PredictBlock predict={predict} />}
          {children}
        </div>
      </div>
      {readout.length > 0 && (
        <div className="border-t border-mathua-border px-3 py-2">
          <dl
            aria-live="polite"
            data-testid="visual-readout"
            className="mx-auto grid w-full max-w-[420px] grid-cols-2 gap-x-4 gap-y-1 sm:grid-cols-3"
          >
            {readout.map((item) => (
              <div key={item.label}>
                <dt className="font-mono text-[9px] tracking-wide text-mathua-muted">{item.label}</dt>
                <dd className="font-mono text-[11px] tabular-nums text-mathua-primary">{item.value}</dd>
              </div>
            ))}
          </dl>
        </div>
      )}
      {caption && (
        <figcaption className="border-t border-mathua-border px-3 py-2 font-mono text-[11px] leading-relaxed text-mathua-secondary">
          {caption}
        </figcaption>
      )}
    </figure>
    {controls && (
      <div data-testid="figure-controls" className="mt-2 border border-mathua-border bg-mathua-surface p-3">
        {controls}
      </div>
    )}
    </>
  )
}
