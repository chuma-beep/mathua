import { useEffect } from 'react'
import type { RefObject } from 'react'
import { clientToSvg } from './scale'

/**
 * Pointer dragging for an SVG. The handler receives coordinates in the SVG's own
 * user units, which is what a template draws in, and pointer capture keeps a drag
 * alive when the pointer leaves the canvas.
 *
 * This is an addition to the keyboard and number-input controls, never their
 * replacement: every value a drag changes also has a native control.
 */
export function useDrag(
  svgRef: RefObject<SVGSVGElement | null>,
  onDrag: (point: { x: number; y: number }) => void,
  enabled = true,
) {
  useEffect(() => {
    const svg = svgRef.current
    if (!svg || !enabled) return

    let active = false
    const move = (event: PointerEvent) => {
      if (!active) return
      event.preventDefault()
      const p = clientToSvg(svg, event.clientX, event.clientY)
      onDrag(p)
    }
    const down = (event: PointerEvent) => {
      if (event.button !== 0) return
      active = true
      svg.setPointerCapture?.(event.pointerId)
      move(event)
    }
    const up = () => {
      active = false
    }

    svg.addEventListener('pointerdown', down)
    svg.addEventListener('pointermove', move)
    svg.addEventListener('pointerup', up)
    svg.addEventListener('pointercancel', up)
    return () => {
      svg.removeEventListener('pointerdown', down)
      svg.removeEventListener('pointermove', move)
      svg.removeEventListener('pointerup', up)
      svg.removeEventListener('pointercancel', up)
    }
  }, [svgRef, onDrag, enabled])
}
