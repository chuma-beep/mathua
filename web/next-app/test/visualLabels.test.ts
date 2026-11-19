import { describe, expect, it } from 'vitest'
import { boxesOverlap, labelBox, placeLabels, type LabelFrame, type LabelItem } from '../components/visual/labels'
import type { Vec2 } from '../components/visual/math/matrix'

// The same geometry as the matrix figure, so the numbers here are the numbers the
// drawing uses rather than a convenient smaller frame.
const SIZE = 320
const PAD = 26
const UNIT = (SIZE - PAD * 2) / 4.8
const frame: LabelFrame = {
  width: SIZE,
  height: SIZE,
  pad: PAD,
  toX: (value) => SIZE / 2 + value * UNIT,
  toY: (value) => SIZE / 2 - value * UNIT,
}

const items = (c1: Vec2, c2: Vec2, v: Vec2, image: Vec2): LabelItem[] => [
  { id: 'image', text: 'Av', tip: image, side: 1 },
  { id: 'c1', text: 'A e₁', tip: c1, side: -1 },
  { id: 'c2', text: 'A e₂', tip: c2, side: 1 },
  { id: 'v', text: 'v', tip: v, side: -1 },
]

function expectClean(c1: Vec2, c2: Vec2, v: Vec2, image: Vec2) {
  const placed = placeLabels(items(c1, c2, v, image), frame)
  expect(placed).toHaveLength(4)
  for (const label of placed) {
    const box = labelBox(label.text, label.x, label.y)
    // Frame-edge clipping: every label is fully inside the plot rectangle.
    expect(box.x).toBeGreaterThanOrEqual(PAD)
    expect(box.y).toBeGreaterThanOrEqual(PAD)
    expect(box.x + box.w).toBeLessThanOrEqual(SIZE - PAD)
    expect(box.y + box.h).toBeLessThanOrEqual(SIZE - PAD)
  }
  for (let i = 0; i < placed.length; i++) {
    for (let j = i + 1; j < placed.length; j++) {
      expect(boxesOverlap(labelBox(placed[i].text, placed[i].x, placed[i].y), labelBox(placed[j].text, placed[j].x, placed[j].y))).toBe(false)
    }
  }
  return placed
}

describe('label placement', () => {
  it('never overlaps two labels in the configuration that used to fuse Av with A e₁', () => {
    // Both arrows point roughly the same way, which is what a fixed offset could
    // not survive: the old figure drew `Av` through `A e₁` here.
    expectClean([1.4, 0.5], [-0.3, 1.1], [1.6, 0.9], [1.97, 1.88])
  })

  it('keeps labels inside the frame when an arrow reaches the edge', () => {
    expectClean([2.2, 0.1], [0.1, 2.2], [-2.2, -0.1], [2.1, 2.0])
  })

  it('survives arrows crowding the origin', () => {
    expectClean([0.05, 0.05], [-0.05, 0.05], [0.05, -0.05], [0, 0.1])
  })

  it('places the same input deterministically', () => {
    const first = placeLabels(items([1.2, 0.4], [-0.3, 1.1], [1.2, -0.5], [1.71, -0.07]), frame)
    const second = placeLabels(items([1.2, 0.4], [-0.3, 1.1], [1.2, -0.5], [1.71, -0.07]), frame)
    expect(second).toEqual(first)
  })

  it('returns one label per item, in the order it was given', () => {
    const placed = placeLabels(items([1, 0], [0, 1], [0.5, 0.5], [0.5, 0.5]), frame)
    expect(placed.map((label) => label.id)).toEqual(['image', 'c1', 'c2', 'v'])
  })
})
