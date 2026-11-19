import { describe, expect, it } from 'vitest'
import { X_MAX, X_MIN, crossSectionArea, lateralArea, profile, radius, slope, surface, volume } from '../components/visual/math/solid'
import { pickRenderer } from '../components/visual/templates/Solid3D'

const close = (a: number, b: number, tolerance = 1e-9) => expect(Math.abs(a - b)).toBeLessThan(tolerance)

describe('solid of revolution', () => {
  it('has the authored profile and its slope', () => {
    close(radius(0), 1)
    close(slope(0), 0.56)
    expect(radius(1)).toBeGreaterThan(1.3)
    expect(radius(3)).toBeLessThan(0.7)
  })

  it('measures volume by disks against an independent Riemann sum', () => {
    const n = 20000
    const h = (X_MAX - X_MIN) / n
    let riemann = 0
    for (let i = 0; i < n; i++) {
      const x = X_MIN + (i + 0.5) * h
      riemann += Math.PI * radius(x) ** 2 * h
    }
    close(volume(), riemann, 1e-3)
    expect(volume()).toBeGreaterThan(0)
  })

  it('reports the disk area and a finite lateral area', () => {
    close(crossSectionArea(0), Math.PI)
    expect(lateralArea()).toBeGreaterThan(0)
    expect(Number.isFinite(lateralArea())).toBe(true)
  })

  it('builds a closed mesh with unit normals', () => {
    const mesh = surface(8, 12)
    expect(mesh.positions.length).toBe((8 + 1) * (12 + 1) * 3)
    expect(mesh.normals.length).toBe(mesh.positions.length)
    expect(mesh.indices.length).toBe(8 * 12 * 6)
    for (let i = 0; i < mesh.normals.length; i += 3) {
      const len = Math.hypot(mesh.normals[i], mesh.normals[i + 1], mesh.normals[i + 2])
      close(len, 1, 1e-9)
    }
    const xs = []
    for (let i = 0; i < mesh.positions.length; i += 3) xs.push(mesh.positions[i])
    close(Math.min(...xs), -2, 1e-9)
    close(Math.max(...xs), 2, 1e-9)
  })

  it('samples the profile from end to end', () => {
    const points = profile(10)
    expect(points.length).toBe(11)
    close(points[0][0], X_MIN)
    close(points[points.length - 1][0], X_MAX)
  })
})

describe('renderer dispatch', () => {
  it('uses WebGL when it is available and the wireframe otherwise', () => {
    expect(pickRenderer(true)).toBe('webgl')
    expect(pickRenderer(false)).toBe('wireframe')
    expect(pickRenderer(true, true)).toBe('wireframe')
    expect(pickRenderer(false, true)).toBe('wireframe')
  })
})
