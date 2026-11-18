import { describe, it, expect } from 'vitest'
import { fireEvent, render, screen, within } from '@testing-library/react'
import MatrixTransform2D from '../components/visual/templates/MatrixTransform2D'
import TangentLine from '../components/visual/templates/TangentLine'
import RiemannSum from '../components/visual/templates/RiemannSum'
import ComplexMultiply from '../components/visual/templates/ComplexMultiply'
import NormalCurve from '../components/visual/templates/NormalCurve'
import Transform2D from '../components/visual/templates/Transform2D'
import VisualFigure from '../components/visual/VisualFigure'

function readout() {
  return screen.getByTestId('visual-readout').textContent ?? ''
}

describe('visual templates', () => {
  it('matrix transform updates the determinant when a coefficient changes', () => {
    render(
      <MatrixTransform2D
        title="Matrix"
        params={{ matrix: [1, 0.5, 0.5, 1], vector: [2, 1], showGrid: true }}
      />,
    )
    expect(readout()).toContain('0.75')
    fireEvent.change(screen.getByLabelText('a'), { target: { value: '2' } })
    expect(readout()).toContain('1.75')
  })

  it('tangent line moves its point and recomputes f(x)', () => {
    render(
      <TangentLine
        title="Tangent"
        params={{ coeffs: [0, 0, 1], x: 1, h: 0.8, window: [-2, 3] }}
      />,
    )
    expect(readout()).toContain('f(x)1')
    fireEvent.change(screen.getByLabelText('x'), { target: { value: '2' } })
    expect(readout()).toContain('f(x)4')
  })

  it('Riemann sum changes its partition count', () => {
    render(
      <RiemannSum
        title="Riemann"
        params={{ coeffs: [0, 0, 1], a: 0, b: 2, n: 8, rule: 'left' }}
      />,
    )
    fireEvent.change(screen.getByLabelText('n'), { target: { value: '0' } })
    const cells = within(screen.getByTestId('visual-readout')).getAllByRole('definition')
    expect(cells.map((c) => c.textContent)).toContain('1')
  })

  it('complex product changes when a real part changes', () => {
    render(
      <ComplexMultiply
        title="Product"
        params={{ z: [1.5, 0.8], w: [0.8, 1.2], showUnitCircle: true }}
      />,
    )
    fireEvent.change(screen.getByLabelText('z re'), { target: { value: '2' } })
    expect(readout()).toContain('2.15')
  })

  it('normal curve reports the probability at the moved point', () => {
    render(
      <NormalCurve
        title="Normal"
        params={{ mu: 0, sigma: 1, x: 1, bands: true }}
      />,
    )
    expect(readout()).toContain('0.84')
    fireEvent.change(screen.getByLabelText('x'), { target: { value: '0' } })
    expect(readout()).toContain('P(X ≤ x)0.5')
  })

  it('rotation leaves distances unchanged and moves the vertex', () => {
    render(
      <Transform2D
        title="Rotation"
        params={{
          kind: 'rotation',
          angle: 30,
          scale: 1.25,
          dx: 1,
          dy: 0.5,
          shape: [[0, 0], [2, 0], [1, 2]],
          showOrbit: true,
        }}
      />,
    )
    expect(readout()).toContain('1.73')
    fireEvent.change(screen.getByLabelText('angle'), { target: { value: '90' } })
    expect(readout()).toContain('(0, 2)')
    expect(readout()).toContain('side')
  })
})

describe('visual dispatcher', () => {
  it('explains an unknown template without breaking the lesson', () => {
    render(<VisualFigure visual={{ type: 'not.a.template' }} />)
    expect(screen.getByText(/This figure could not be shown/)).toBeInTheDocument()
  })
})
