'use client'

import type { ReactNode } from 'react'
import type { ReadoutItem } from './types'

interface FigureFrameProps {
  title: string
  caption?: string
  readout: ReadoutItem[]
  onReset?: () => void
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
export default function FigureFrame({ title, caption, readout, onReset, children }: FigureFrameProps) {
  return (
    <figure className="mt-5 border border-mathua-border bg-mathua-bg">
      <div className="flex items-center justify-between gap-3 border-b border-mathua-border px-3 py-2">
        <p className="font-mono text-[11px] uppercase tracking-wider text-mathua-muted">{title}</p>
        {onReset && (
          <button
            type="button"
            onClick={onReset}
            className="border border-mathua-border px-2 py-0.5 font-mono text-[10px] text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue"
          >
            Reset
          </button>
        )}
      </div>
      <div className="p-3">{children}</div>
      {readout.length > 0 && (
        <dl
          aria-live="polite"
          data-testid="visual-readout"
          className="grid grid-cols-2 gap-x-4 gap-y-1 border-t border-mathua-border px-3 py-2 sm:grid-cols-3"
        >
          {readout.map((item) => (
            <div key={item.label}>
              <dt className="font-mono text-[9px] uppercase tracking-wide text-mathua-muted">{item.label}</dt>
              <dd className="font-mono text-[11px] tabular-nums text-mathua-primary">{item.value}</dd>
            </div>
          ))}
        </dl>
      )}
      {caption && (
        <figcaption className="border-t border-mathua-border px-3 py-2 font-mono text-[10px] leading-relaxed text-mathua-secondary">
          {caption}
        </figcaption>
      )}
    </figure>
  )
}
