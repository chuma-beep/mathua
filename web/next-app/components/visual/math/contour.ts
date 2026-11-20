// Marching squares. One routine serves every figure that needs a level set: the
// conic sections and the contour plots are both "where does this field equal this
// number", and drawing them two ways would let the two disagree.
//
// A cell is skipped when any corner is not finite, because the alternative is
// interpolating toward an infinity and drawing a line that means nothing.

export interface Bounds {
  x0: number
  x1: number
  y0: number
  y1: number
}

export type Segment = [[number, number], [number, number]]

export function contourSegments(
  f: (x: number, y: number) => number,
  bounds: Bounds,
  level: number,
  nx = 60,
  ny = 60,
): Segment[] {
  const { x0, x1, y0, y1 } = bounds
  const dx = (x1 - x0) / nx
  const dy = (y1 - y0) / ny

  // One sample per grid node, reused by the four cells that share it.
  const values: number[][] = []
  for (let j = 0; j <= ny; j++) {
    const row: number[] = []
    for (let i = 0; i <= nx; i++) row.push(f(x0 + i * dx, y0 + j * dy))
    values.push(row)
  }

  const out: Segment[] = []

  for (let j = 0; j < ny; j++) {
    for (let i = 0; i < nx; i++) {
      const v = [values[j][i], values[j][i + 1], values[j + 1][i + 1], values[j + 1][i]]
      if (v.some((value) => !Number.isFinite(value))) continue

      const corners: Array<[number, number]> = [
        [x0 + i * dx, y0 + j * dy],
        [x0 + (i + 1) * dx, y0 + j * dy],
        [x0 + (i + 1) * dx, y0 + (j + 1) * dy],
        [x0 + i * dx, y0 + (j + 1) * dy],
      ]

      const crossings: Array<[number, number]> = []
      for (let e = 0; e < 4; e++) {
        const a = e
        const b = (e + 1) % 4
        if ((v[a] >= level) === (v[b] >= level)) continue
        const t = (level - v[a]) / (v[b] - v[a])
        crossings.push([
          corners[a][0] + t * (corners[b][0] - corners[a][0]),
          corners[a][1] + t * (corners[b][1] - corners[a][1]),
        ])
      }

      if (crossings.length === 2) out.push([crossings[0], crossings[1]])
      // A saddle cell crosses all four edges. Pairing in visit order is the
      // simple reading; it can pick the other topology on an exact saddle, which
      // only matters for a figure whose level set passes through a grid node.
      else if (crossings.length === 4) {
        out.push([crossings[0], crossings[1]])
        out.push([crossings[2], crossings[3]])
      }
    }
  }

  return out
}
