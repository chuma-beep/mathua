import { describe, it, expect } from 'vitest'
import { fireEvent, render, screen } from '@testing-library/react'
import ComplexLevels from '../components/visual/templates/ComplexLevels'

function readout() {
  return screen.getByTestId('visual-readout').textContent ?? ''
}

describe('complex level curves', () => {
  it('reports the right angle between the families at a regular point', () => {
    render(<ComplexLevels title="Levels" params={{ fn: 'square', probe: [1, 0.6], levels: 5, showImag: true }} />)
    expect(readout()).toContain('f(z)0.64 + 1.2i')
    expect(readout()).toContain('Cauchy-Riemannsatisfied')
    expect(readout()).toContain('angle90°')
  })

  it('calls the angle undefined at a critical point and the argument undefined at a zero', () => {
    render(<ComplexLevels title="Levels" params={{ fn: 'square', probe: [1, 0.6], levels: 5, showImag: true }} />)
    fireEvent.change(screen.getByLabelText('x'), { target: { value: '0' } })
    fireEvent.change(screen.getByLabelText('y'), { target: { value: '0' } })
    expect(readout()).toContain('arg fundefined')
    expect(readout()).toContain('angleundefined')
  })

  it('separates a pole from a zero and from a critical point', () => {
    const { unmount } = render(
      <ComplexLevels title="Levels" params={{ fn: 'inverse', probe: [0, 0], levels: 5, showImag: true }} />,
    )
    expect(readout()).toContain('f(z)undefined')
    // A pole is not a failed analyticity test; there is nothing to test there.
    expect(readout()).toContain('Cauchy-Riemannundefined')
    const svg = screen.getByRole('img')
    expect(svg.textContent).toContain('pole: no value here')
    expect(svg.textContent).not.toContain('critical point')
    unmount()

    // z squared at the origin is a zero and a critical point, not a pole.
    render(<ComplexLevels title="Levels" params={{ fn: 'square', probe: [0, 0], levels: 5, showImag: true }} />)
    expect(readout()).toContain('Cauchy-Riemannsatisfied')
    expect(readout()).toContain('angleundefined')
    const svg2 = screen.getByRole('img')
    expect(svg2.textContent).toContain('critical point: no tangent')
    expect(svg2.textContent).not.toContain('pole')
  })

  it('finds the critical points of sin away from its zeros', () => {
    render(
      <ComplexLevels
        title="Levels"
        params={{ fn: 'sine', probe: [Math.PI / 2, 0], levels: 5, showImag: true }}
      />,
    )
    expect(readout()).toContain('Cauchy-Riemannsatisfied')
    expect(readout()).toContain('angleundefined')
    // sin(pi/2 + 0i) = 1, so the argument exists (and is 0) even though the
    // derivative vanishes there.
    expect(readout()).toContain('arg f0')
  })
})
