import type { Vec2 } from './math/matrix'

// Label placement, not label guessing.
//
// A fixed offset per arrow is what put `Av` on top of `A e₁` whenever the two
// arrows pointed the same way, and what pushed a label through the frame edge at
// 320px. Each label now tries a short list of offsets and takes the first one
// that stays inside the frame and clear of every label already placed and every
// other arrow's knob. The search is driven only by the numbers being drawn, so it
// is pure and testable without a DOM.

export interface LabelFrame {
  width: number
  height: number
  pad: number
  toX: (value: number) => number
  toY: (value: number) => number
}

export interface LabelItem {
  id: string
  text: string
  /** The arrowhead the label names, in data coordinates. */
  tip: Vec2
  /** Which side of the shaft to prefer: 1 or -1. */
  side: number
}

export interface PlacedLabel {
  id: string
  text: string
  x: number
  y: number
}

export interface LabelBox {
  x: number
  y: number
  w: number
  h: number
}

/** Height of an 11px label with its halo, in user units. */
export const LABEL_HEIGHT = 13

/**
 * An estimate of rendered width for an 11px mono label. Deliberately generous:
 * an over-wide box moves a label slightly further than necessary, while an
 * under-wide box lets it collide.
 */
export function labelWidth(text: string): number {
  return text.length * 6.9 + 2
}

export function labelBox(text: string, x: number, y: number): LabelBox {
  return { x: x - labelWidth(text) / 2, y: y - LABEL_HEIGHT / 2, w: labelWidth(text), h: LABEL_HEIGHT }
}

export function boxesOverlap(a: LabelBox, b: LabelBox): boolean {
  const overlapX = Math.min(a.x + a.w, b.x + b.w) - Math.max(a.x, b.x)
  const overlapY = Math.min(a.y + a.h, b.y + b.h) - Math.max(a.y, b.y)
  return overlapX > 1 && overlapY > 1
}

/** The square a knob occupies, so a label is never placed on a grab handle. */
const KNOB = 26

export function placeLabels(items: LabelItem[], frame: LabelFrame): PlacedLabel[] {
  const taken: LabelBox[] = []
  const knobs: LabelBox[] = items.map((item) => ({
    x: frame.toX(item.tip[0]) - KNOB / 2,
    y: frame.toY(item.tip[1]) - KNOB / 2,
    w: KNOB,
    h: KNOB,
  }))

  return items.map((item, index) => {
    const tipX = frame.toX(item.tip[0])
    const tipY = frame.toY(item.tip[1])
    const dx = tipX - frame.width / 2
    const dy = frame.height / 2 - tipY
    const len = Math.hypot(dx, dy) || 1
    // [side, perpendicular offset, radial offset], in user units.
    const candidates: Array<[number, number, number]> = [
      [item.side, 15, 16],
      [item.side, 25, 18],
      [-item.side, 15, 16],
      [item.side, 15, 26],
      [-item.side, 25, 20],
      [item.side, 34, 18],
      [-item.side, 34, 22],
      // Near the frame edge the only room left is behind the arrowhead, so the
      // last candidates step back along the shaft instead of outward from it.
      [item.side, 15, -14],
      [-item.side, 15, -14],
      [item.side, 24, -22],
      [-item.side, 24, -22],
    ]

    for (const [side, perpendicular, radial] of candidates) {
      const x = tipX + (dx / len) * radial + (-dy / len) * perpendicular * side
      const y = tipY + (-dy / len) * radial + (-dx / len) * perpendicular * side
      const box = labelBox(item.text, x, y)
      const inside =
        box.x > frame.pad + 2 &&
        box.y > frame.pad + 2 &&
        box.x + box.w < frame.width - frame.pad - 2 &&
        box.y + box.h < frame.height - frame.pad - 2
      // Its own knob is skipped: the radial offset is chosen to clear that one.
      const clear =
        taken.every((other) => !boxesOverlap(box, other)) &&
        knobs.every((knob, knobIndex) => knobIndex === index || !boxesOverlap(box, knob))
      if (inside && clear) {
        taken.push(box)
        return { id: item.id, text: item.text, x, y }
      }
    }

    // No candidate cleared: keep the label beside the tip rather than dropping it,
    // and let the caller trust the overlap probe to report the collision.
    const x = tipX + (dx / len) * 12
    const y = tipY + (-dy / len) * 12
    taken.push(labelBox(item.text, x, y))
    return { id: item.id, text: item.text, x, y }
  })
}
