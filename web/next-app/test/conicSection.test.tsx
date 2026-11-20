import { describe, expect, it } from 'vitest'
import { fireEvent, render, screen } from '@testing-library/react'
import ConicSection from '../components/visual/templates/ConicSection'

const readout = () => screen.getByTestId('visual-readout').textContent ?? ''
const item = (label: string) => {
  const cells = Array.from(document.querySelectorAll('[data-testid="visual-readout"] dt'))
  const cell = cells.find((node) => node.textContent === label)
  return cell?.nextElementSibling?.textContent ?? ''
}
const preset = () => screen.getByLabelText('Preset equation') as HTMLSelectElement
const slider = (label: string) => screen.getByLabelText(label) as HTMLInputElement
const setSlider = (label: string, value: number) =>
  fireEvent.change(slider(label), { target: { value: String(value) } })

describe('conic section controls', () => {
  it('exposes all six coefficients and the preset control', () => {
    render(<ConicSection title="Conic" params={{ A: 1, B: 0, C: 0.5 }} />)
    for (const label of ['A', 'B', 'C', 'D', 'E', 'F']) {
      expect(slider(label)).toBeTruthy()
    }
    expect(preset()).toBeTruthy()
  })

  it('starts on the ellipse the old three-term params described', () => {
    render(<ConicSection title="Conic" params={{ A: 1, B: 0, C: 0.5 }} />)
    // F defaults to -1, so Ax² + Cy² = 1 is unchanged.
    expect(readout()).toContain('A, B, C1, 0, 0.5')
    expect(readout()).toContain('D, E, F0, 0, -1')
    expect(readout()).toContain('typeellipse')
    expect(readout()).toContain('centre(0, 0)')
    expect(readout()).toContain('axes1.414 × 1')
    expect(readout()).toContain('eccentricity0.707')
    expect(preset().value).toBe('Ellipse')
    expect(screen.getByRole('img', { name: 'ellipse from its general equation' })).toBeTruthy()
  })

  it('draws a parabola with the vertex, p and opening direction', () => {
    render(<ConicSection title="Conic" params={{ A: 1, B: 0, C: 0.5 }} />)
    fireEvent.change(preset(), { target: { value: 'Parabola' } })
    expect(readout()).toContain('typeparabola')
    expect(readout()).toContain('vertex(0, 0)')
    expect(readout()).toContain('opening(0, 1)')
    expect(readout()).toContain('p = 0.25')
    expect(readout()).toContain('eccentricity1')
    expect(readout()).toContain('squared terms1 written')
  })

  it('draws a hyperbola with its transverse orientation and eccentricity', () => {
    render(<ConicSection title="Conic" params={{ A: 1, B: 0, C: 0.5 }} />)
    fireEvent.change(preset(), { target: { value: 'Hyperbola' } })
    expect(readout()).toContain('typehyperbola')
    expect(readout()).toContain('axes1 × 1')
    expect(readout()).toContain('transverse-axis angle0°')
    expect(readout()).toContain('eccentricity1.414')
    expect(preset().value).toBe('Hyperbola')
  })

  it('translates a conic with the linear terms', () => {
    render(<ConicSection title="Conic" params={{ A: 1, B: 0, C: 0.5 }} />)
    fireEvent.change(preset(), { target: { value: 'Translated ellipse' } })
    expect(readout()).toContain('centre(1, -2)')
    expect(readout()).toContain('axes2 × 1')
    expect(readout()).toContain('eccentricity0.866')

    // The same circle moved by D and E: x² + y² - 2x + 2y - 2 = 0.
    fireEvent.change(preset(), { target: { value: 'Circle' } })
    setSlider('D', -2)
    setSlider('E', 2)
    expect(readout()).toContain('centre(1, -1)')
    expect(readout()).toContain('typecircle')
    expect(readout()).toContain('eccentricity0')
  })

  it('reports intersecting lines at the degenerate preset instead of a hyperbola', () => {
    render(<ConicSection title="Conic" params={{ A: 1, B: 0, C: 0.5 }} />)
    fireEvent.change(preset(), { target: { value: 'Degenerate' } })
    expect(readout()).toContain('typetwo intersecting lines')
    expect(readout()).toContain('discriminant4')
    expect(readout()).not.toContain('typehyperbola')
  })

  it('reports a point when both squares vanish at the origin', () => {
    render(<ConicSection title="Conic" params={{ A: 1, B: 0, C: 1 }} />)
    setSlider('F', 0)
    expect(readout()).toContain('typea single point')
    expect(readout()).toContain('centre(0, 0)')
  })

  it('reports no real points and says so inside the drawing', () => {
    render(<ConicSection title="Conic" params={{ A: 1, B: 0, C: 1 }} />)
    setSlider('F', 1)
    expect(readout()).toContain('typeno real points')
    expect(screen.getAllByText('no real points').length).toBeGreaterThan(0)
    expect(screen.getByRole('img', { name: 'no real points from its general equation' })).toBeTruthy()
  })

  it('distinguishes the all-plane equation from the empty one', () => {
    render(<ConicSection title="Conic" params={{ A: 1, B: 0, C: 1 }} />)
    setSlider('A', 0)
    setSlider('C', 0)
    setSlider('F', 0)
    expect(readout()).toContain('typeevery point')
    expect(screen.getAllByText('every point satisfies the equation').length).toBeGreaterThan(0)
    expect(screen.queryByText('no real points')).toBeNull()
  })

  it('flips the preset select to Custom after an edit and restores it on reset', () => {
    render(<ConicSection title="Conic" params={{ A: 1, B: 0, C: 1 }} />)
    expect(preset().value).toBe('Circle')
    setSlider('A', 0.8)
    expect(preset().value).toBe('')
    expect(screen.getByRole('option', { name: 'Custom' })).toBeTruthy()

    fireEvent.click(screen.getByText('Reset'))
    expect(preset().value).toBe('Circle')
    expect(readout()).toContain('typecircle')
  })

  it('moves from an ellipse to a hyperbola as B passes the balance point', () => {
    render(<ConicSection title="Conic" params={{ A: 1, B: 0, C: 0.5 }} />)
    setSlider('B', 1)
    expect(readout()).toContain('typeellipse')
    expect(readout()).toContain('discriminant-1')

    setSlider('B', 2)
    expect(readout()).toContain('typehyperbola')
    expect(readout()).toContain('discriminant2')
  })

  it('renders the rotated ellipse and turns it into a hyperbola past the balance point', () => {
    render(<ConicSection title="Conic" params={{ A: 1, B: 1, C: 1 }} />)
    expect(readout()).toContain('typeellipse')
    expect(readout()).toContain('minor-axis angle45°')

    setSlider('C', 0.5)
    setSlider('B', 2)
    expect(readout()).toContain('typehyperbola')
  })

  it('keeps a near-degenerate hyperbola a hyperbola and reports the uncertainty', () => {
    render(<ConicSection title="Conic" params={{ A: 1, B: 0, C: -1 }} />)
    setSlider('F', -0.0000001)
    expect(readout()).toContain('typehyperbola')
    expect(readout()).toContain('certaintynear-degenerate')
  })

  it('does not crash at the extremes of the controls', () => {
    render(<ConicSection title="Conic" params={{ A: 2, B: 2, C: 2, D: 4, E: 4, F: 4 }} />)
    expect(readout().length).toBeGreaterThan(0)
    expect(readout()).not.toContain('NaN')
    expect(readout()).not.toContain('Infinity')
  })

  it('gives each mounted figure its own clip path', () => {
    render(
      <>
        <ConicSection title="First" params={{ A: 1, B: 0, C: 1 }} />
        <ConicSection title="Second" params={{ A: 1, B: 0, C: 0.5 }} />
      </>,
    )
    const ids = Array.from(document.querySelectorAll('clipPath')).map((node) => node.id)
    expect(ids).toHaveLength(2)
    expect(new Set(ids).size).toBe(2)
  })
})

describe('conic focal interaction', () => {
  it('shows foci and the constant focal sum for an ellipse', () => {
    render(<ConicSection title="Conic" params={{ A: 1, B: 0, C: 0.5 }} />)
    expect(item('foci')).toContain('F1 (0, -1)')
    expect(item('foci')).toContain('F2 (0, 1)')
    expect(item('focal property')).toContain('PF1 + PF2')
    expect(item('focal property')).toContain('2a = 2.8284')
    expect(item('eccentricity')).toBe('0.707 (constant)')
    const probe = screen.getByLabelText('Probe position')
    expect(probe.getAttribute('type')).toBe('range')
    expect(probe.getAttribute('min')).toBe('0')
  })

  it('moves the probe along the curve without changing the conic eccentricity or focal sum', () => {
    render(<ConicSection title="Conic" params={{ A: 1, B: 0, C: 0.5 }} />)
    const probe = screen.getByLabelText('Probe position') as HTMLInputElement
    const before = item('probe')
    const eccentricity = item('eccentricity')
    const focalProperty = item('focal property')
    fireEvent.change(probe, { target: { value: '1.9' } })
    expect(item('probe')).not.toBe(before)
    expect(item('probe')).toContain('(0.95, -0.46)')
    expect(item('eccentricity')).toBe(eccentricity)
    expect(item('focal property')).toBe(focalProperty)
  })

  it('switches hyperbola branches and keeps the focal difference on both', () => {
    render(<ConicSection title="Conic" params={{ A: 1, B: 0, C: 0.5 }} />)
    fireEvent.change(preset(), { target: { value: 'Hyperbola' } })
    expect(item('probe')).toMatch(/^\(1\./)
    expect(item('focal property')).toContain('|PF1 − PF2|')
    expect(item('focal property')).toContain('2a = 2')

    const first = screen.getByRole('button', { name: 'First branch' })
    const second = screen.getByRole('button', { name: 'Second branch' })
    expect(first.getAttribute('aria-pressed')).toBe('true')
    fireEvent.click(second)
    expect(second.getAttribute('aria-pressed')).toBe('true')
    expect(item('probe')).toMatch(/^\(-1\./)
    expect(item('focal property')).toContain('|PF1 − PF2|')
    expect(item('focal property')).toContain('2a = 2')
  })

  it('draws the parabola focus, directrix and latus rectum', () => {
    const { container } = render(<ConicSection title="Conic" params={{ A: 1, B: 0, C: 0.5 }} />)
    fireEvent.change(preset(), { target: { value: 'Parabola' } })
    expect(item('foci')).toContain('(0, 0.25)')
    expect(item('directrix')).toContain('through (0, -0.25), direction (1, 0)')
    expect(item('latus rectum')).toBe('1')
    expect(item('focal property')).toContain('PF = ')
    expect(item('distance to directrix')).toBe('0.74')
    const directrix = container.querySelector('line[stroke="var(--chart-3)"]')
    expect(directrix).not.toBeNull()
    expect(container.querySelector('[data-part="latus-rectum"]')).not.toBeNull()
  })

  it('states the coincident-focus circle case without inventing two foci', () => {
    const { container } = render(<ConicSection title="Conic" params={{ A: 1, B: 0, C: 0.5 }} />)
    fireEvent.change(preset(), { target: { value: 'Circle' } })
    expect(item('foci')).toBe('both at (0, 0)')
    expect(item('focal property')).toContain('PF = 1 = r')
    expect(item('eccentricity')).toBe('0 (constant)')
    expect(container.querySelector('svg')?.textContent).toContain('F')
    expect(container.querySelector('svg')?.textContent).not.toContain('F1')
  })

  it('disables focal interaction for a degenerate equation', () => {
    render(<ConicSection title="Conic" params={{ A: 1, B: 0, C: 0.5 }} />)
    fireEvent.change(preset(), { target: { value: 'Degenerate' } })
    expect(screen.queryByLabelText('Probe position')).toBeNull()
    expect(screen.getByText('Focal probe unavailable for this equation.')).toBeTruthy()
    expect(item('foci')).toBe('—')
    expect(item('focal property')).toBe('—')
  })

  it('projects a pointer drag onto the nearest point of the curve', () => {
    const { container } = render(<ConicSection title="Conic" params={{ A: 1, B: 0, C: 0.5 }} />)
    const svg = container.querySelector('svg') as SVGSVGElement
    const before = item('probe')
    fireEvent.pointerDown(svg, { clientX: 280, clientY: 60, button: 0, pointerId: 1 })
    expect(item('probe')).not.toBe(before)
    expect(item('probe')).toMatch(/^\(/)
    expect(readout()).not.toContain('NaN')
  })

  it('keeps focus labels and the probe control keyboard reachable', () => {
    const { container } = render(<ConicSection title="Conic" params={{ A: 1, B: 0, C: 0.5 }} />)
    const svg = container.querySelector('svg') as SVGSVGElement
    const labels = Array.from(svg.querySelectorAll('text')).map((node) => node.textContent)
    expect(labels).toContain('F1')
    expect(labels).toContain('F2')

    const probe = screen.getByLabelText('Probe position') as HTMLInputElement
    probe.focus()
    expect(document.activeElement).toBe(probe)
    expect(probe.getAttribute('step')).toBe(String(Math.PI / 90))
  })

  it('does not leak non-finite values through the focal readouts', () => {
    render(<ConicSection title="Conic" params={{ A: 2, B: 2, C: 2, D: 4, E: 4, F: 4 }} />)
    expect(readout()).not.toContain('NaN')
    expect(readout()).not.toContain('Infinity')
    expect(item('focal property')).not.toContain('NaN')
  })
})
