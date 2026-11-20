import { describe, expect, it } from 'vitest'
import { render } from '@testing-library/react'
import SolidWireframe from '../components/visual/templates/solid/SolidWireframe'

describe('wireframe fallback', () => {
  it('draws the same solid and marks itself as the fallback', () => {
    const { container } = render(
      <SolidWireframe theta={0.6} phi={0.35} cut={2} showCut onOrbit={() => {}} />,
    )
    const svg = container.querySelector('svg[aria-label="A solid of revolution drawn as a wireframe"]')
    expect(svg).not.toBeNull()
    expect(container.querySelectorAll('path').length).toBeGreaterThan(20)
    expect(container.textContent).toContain('wireframe fallback')
  })

  it('omits the cross-section when the cut is hidden', () => {
    const { container } = render(
      <SolidWireframe theta={0} phi={0} cut={1} showCut={false} onOrbit={() => {}} />,
    )
    const teal = container.querySelector('path[stroke="var(--accent-teal)"]')
    expect(teal).toBeNull()
  })
})
