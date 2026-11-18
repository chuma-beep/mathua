import { describe, it, expect } from 'vitest'
import { fireEvent, render, screen } from '@testing-library/react'
import MatrixTransform2D from '../components/visual/templates/MatrixTransform2D'
import EigenDirections from '../components/visual/templates/EigenDirections'
import TaylorSeries from '../components/visual/templates/TaylorSeries'
import LimitApproach from '../components/visual/templates/LimitApproach'
import SamplingDistribution from '../components/visual/templates/SamplingDistribution'
import LeastSquares from '../components/visual/templates/LeastSquares'

function readout() {
  return screen.getByTestId('visual-readout').textContent ?? ''
}

describe('second-tranche templates', () => {
  it('composes two matrices and reports BA applied to the vector', () => {
    render(
      <MatrixTransform2D
        title="Composition"
        params={{
          matrix: [1, 0, 0, 1],
          then: [0, -1, 1, 0],
          vector: [2, 1],
          showGrid: true,
          showBasis: true,
        }}
      />,
    )
    // Identity, then a quarter turn: BA is the rotation, so (2, 1) lands on (-1, 2).
    expect(readout()).toContain('BAv(-1, 2)')
    expect(readout()).toContain('det BA1')
  })

  it('finds the eigen-direction of a symmetric matrix', () => {
    render(
      <EigenDirections
        title="Eigen"
        params={{ matrix: [2, 1, 1, 2], vector: [1, 0], showCircle: true }}
      />,
    )
    expect(readout()).toContain('turn')
    fireEvent.change(screen.getByLabelText('angle'), { target: { value: '45' } })
    // (1,1) is an eigenvector of [[2,1],[1,2]] with eigenvalue 3.
    expect(readout()).toMatch(/turn0°/)
    expect(readout()).toContain('v·Av3')
  })

  it('raises the degree and lowers the error', () => {
    render(
      <TaylorSeries title="Taylor" params={{ fn: 'sin', center: 0, degree: 3, window: [-6, 6] }} />,
    )
    expect(readout()).toContain('degree3')
    fireEvent.change(screen.getByLabelText('degree'), { target: { value: '9' } })
    expect(readout()).toContain('degree9')
    fireEvent.change(screen.getByLabelText('centre'), { target: { value: '1' } })
    expect(readout()).toContain('centre1')
  })

  it('separates a hole from a jump', () => {
    render(<LimitApproach title="Limit" params={{ kind: 'removable', a: 2, delta: 0.6 }} />)
    expect(readout()).toContain('two-sidedexists')
    expect(readout()).toContain('left limit4')
    fireEvent.change(screen.getByLabelText('a'), { target: { value: '1' } })
    expect(readout()).toContain('left limit2')
  })

  it('shrinks the standard error as the sample grows', () => {
    render(
      <SamplingDistribution
        title="Sampling"
        params={{ population: 'skewed', n: 9, draws: 200, seed: 7 }}
      />,
    )
    expect(readout()).toContain('standard error0.333')
    fireEvent.change(screen.getByLabelText('n'), { target: { value: '25' } })
    expect(readout()).toContain('standard error0.2')
  })

  it('moves the least-squares line when a point moves', () => {
    render(
      <LeastSquares
        title="Regression"
        params={{
          points: [
            [1, 2],
            [2, 3],
            [3, 3.5],
            [4, 6],
            [5, 6.5],
          ],
          showResiduals: true,
        }}
      />,
    )
    expect(readout()).toContain('slope1.2')
    const before = readout()
    fireEvent.change(screen.getByLabelText('y'), { target: { value: '8' } })
    expect(readout()).not.toBe(before)
  })
})
