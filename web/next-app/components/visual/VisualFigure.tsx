'use client'

import { Component, Suspense, lazy, useMemo, type ComponentType, type ReactNode } from 'react'
import FigureFrame from './FigureFrame'
import { VISUAL_TEMPLATES } from './registry'
import { PredictProvider } from './predict'
import { useReducedMotion } from './useReducedMotion'
import type { TemplateProps, VisualSpec } from './types'

// React requires a lazy component's type to be identical on every render, and
// creating one inside the component breaks that rule (and the `static-components`
// lint rule): a new type remounts the figure on every keystroke. Wrapping every
// template once at module load keeps the types stable, and the chunk is still
// fetched only when that template first renders, so a lesson with no visual
// downloads none of this.
const LAZY_TEMPLATES: Map<string, ComponentType<TemplateProps>> = new Map(
  Object.entries(VISUAL_TEMPLATES).map(([type, entry]) => [type, lazy(entry.load)] as const),
)

class VisualErrorBoundary extends Component<{ fallback: ReactNode; children: ReactNode }, { failed: boolean }> {
  state = { failed: false }

  static getDerivedStateFromError() {
    return { failed: true }
  }

  render() {
    return this.state.failed ? this.props.fallback : this.props.children
  }
}

interface VisualFigureProps {
  visual: VisualSpec
}

/**
 * The dispatcher. One place decides whether a spec becomes a figure, a loading
 * placeholder, or a quiet fallback, and no path here can take down the lesson:
 * an unknown type, a schema failure or a render error all leave the surrounding
 * card intact.
 */
export default function VisualFigure({ visual }: VisualFigureProps) {
  const reducedMotion = useReducedMotion()
  const entry = VISUAL_TEMPLATES[visual.type]
  const parsed = useMemo(() => (entry ? entry.schema.safeParse(visual.params ?? {}) : null), [entry, visual.params])
  const Template = LAZY_TEMPLATES.get(visual.type) ?? null

  if (!entry || !parsed || !parsed.success || !Template) {
    return (
      <FigureFrame title={visual.title || 'Interactive figure'} caption={visual.caption} readout={[]}>
        <p className="font-mono text-[11px] text-mathua-muted">
          This figure could not be shown. The lesson continues below.
        </p>
      </FigureFrame>
    )
  }

  return (
    <VisualErrorBoundary
      fallback={
        <FigureFrame title={visual.title || 'Interactive figure'} caption={visual.caption} readout={[]}>
          <p className="font-mono text-[11px] text-mathua-muted">
            This figure could not be shown. The lesson continues below.
          </p>
        </FigureFrame>
      }
    >
      <Suspense
        fallback={
          <FigureFrame title={visual.title || 'Interactive figure'} readout={[]}>
            <div
              className={`h-[260px] w-full bg-mathua-surface ${reducedMotion ? '' : 'animate-pulse'}`}
              aria-hidden="true"
            />
          </FigureFrame>
        }
      >
        <PredictProvider predict={visual.predict}>
          <Template params={parsed.data} title={visual.title} caption={visual.caption} />
        </PredictProvider>
      </Suspense>
    </VisualErrorBoundary>
  )
}
