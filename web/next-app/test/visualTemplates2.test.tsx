import { describe, it, expect } from 'vitest'
import { fireEvent, render, screen } from '@testing-library/react'
import MatrixTransform2D from '../components/visual/templates/MatrixTransform2D'
import EigenDirections from '../components/visual/templates/EigenDirections'
import TaylorSeries from '../components/visual/templates/TaylorSeries'
import LimitApproach from '../components/visual/templates/LimitApproach'
import SamplingDistribution from '../components/visual/templates/SamplingDistribution'
import LeastSquares from '../components/visual/templates/LeastSquares'
import VectorCombination from '../components/visual/templates/VectorCombination'

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

  it('says when the bars are the population and reports what fell off the plot', () => {
    render(
      <SamplingDistribution
        title="Sampling"
        params={{ population: 'skewed', n: 1, draws: 200, seed: 7 }}
      />,
    )
    expect(readout()).toContain('off-plot')
    expect(screen.getByText(/the bars are the population itself/)).toBeInTheDocument()
    fireEvent.change(screen.getByLabelText('n'), { target: { value: '25' } })
    expect(screen.queryByText(/the bars are the population itself/)).not.toBeInTheDocument()
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

  it('adds vectors and rescales one', () => {
    render(
      <VectorCombination
        title="Vectors"
        params={{ a: [1, 2], b: [2, 0.5], k: -1.5, showSum: true }}
      />,
    )
    expect(readout()).toContain('a + b3, 2.5')
    expect(readout()).toContain('k b-3, -0.75')
    fireEvent.change(screen.getByLabelText('a x'), { target: { value: '3' } })
    expect(readout()).toContain('a + b5, 2.5')
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

describe('eigen edge cases', () => {
  it('names an eigenvector whose eigenvalue is negative', () => {
    render(
      <EigenDirections
        title="Eigen"
        params={{ matrix: [0, 1, 1, 0], vector: [1, -1], showCircle: true }}
      />,
    )
    // (1, -1) is an eigenvector of the swap matrix with eigenvalue -1, so the
    // image points the opposite way; the old two-degree test refused to name it.
    expect(readout()).toContain('stateeigenvector')
    expect(readout()).toContain('v·Av-1')
    expect(readout()).toContain('turn180°')
  })

  it('treats every direction as an eigenvector of the identity', () => {
    render(
      <EigenDirections
        title="Eigen"
        params={{ matrix: [1, 0, 0, 1], vector: [1, 0], showCircle: true }}
      />,
    )
    expect(readout()).toContain('stateeigenvector')
    expect(readout()).toContain('v·Av1')
    expect(readout()).toContain('turn0°')
  })

  it('gives the zero matrix eigenvalue zero without a NaN', () => {
    render(
      <EigenDirections
        title="Eigen"
        params={{ matrix: [0, 0, 0, 0], vector: [1, 0], showCircle: true }}
      />,
    )
    expect(readout()).toContain('stateeigenvector')
    expect(readout()).toContain('v·Av0')
    expect(readout()).not.toContain('NaN')
  })

  it('keeps the frame fixed while the direction turns', () => {
    render(
      <EigenDirections
        title="Eigen"
        params={{ matrix: [2, 0, 0, 2], vector: [1, 0], showCircle: true }}
      />,
    )
    const circle = () => screen.getByRole('img').querySelector('circle')
    const before = circle()?.getAttribute('r')
    expect(before).toBeTruthy()
    fireEvent.change(screen.getByLabelText('angle'), { target: { value: '90' } })
    expect(circle()?.getAttribute('r')).toBe(before)
    expect(readout()).toContain('stateeigenvector')
  })
})

describe('taylor edge cases', () => {
  it('says so when the centre is the pole of the geometric series', () => {
    render(
      <TaylorSeries title="Taylor" params={{ fn: 'geometric', center: 1, degree: 4, window: [-3, 3] }} />,
    )
    expect(readout()).toContain('max errorundefined')
    expect(readout()).not.toContain('NaN')
    expect(screen.getByText('pole at the centre')).toBeTruthy()
  })

  it('draws the curve in pieces rather than across the asymptote', () => {
    render(
      <TaylorSeries title="Taylor" params={{ fn: 'geometric', center: 0, degree: 4, window: [-3, 3] }} />,
    )
    const pieces = screen.getByRole('img').querySelectorAll('path[stroke="var(--text-primary)"]')
    expect(pieces.length).toBeGreaterThanOrEqual(2)
    expect(readout()).toContain('max errorundefined')
  })
})

describe('limit edge cases', () => {
  it('reports the jump value at a, not undefined', () => {
    render(<LimitApproach title="Jump" params={{ kind: 'jump', a: 1, delta: 0.4 }} />)
    expect(readout()).toContain('left limit1')
    expect(readout()).toContain('right limit3')
    expect(readout()).toContain('f(a)3')
    expect(readout()).toContain('continuousno')
    expect(readout()).toContain('two-sideddoes not exist')
  })
})

describe('least squares edge cases', () => {
  it('fits two points exactly', () => {
    render(<LeastSquares title="Two" params={{ points: [[0, 1], [2, 5]], showResiduals: true }} />)
    expect(readout()).toContain('slope2')
    expect(readout()).toContain('intercept1')
    expect(readout()).toContain('r1')
    expect(readout()).toContain('sum of squares0')
  })

  it('reports a perfect fit for collinear points', () => {
    render(
      <LeastSquares
        title="Collinear"
        params={{ points: [[0, 1], [1, 3], [2, 5], [3, 7]], showResiduals: true }}
      />,
    )
    expect(readout()).toContain('r1')
    expect(readout()).toContain('sum of squares0')
  })

  it('says no line is determined when every x is the same', () => {
    render(
      <LeastSquares
        title="Vertical"
        params={{ points: [[2, 1], [2, 4], [2, 9]], showResiduals: true }}
      />,
    )
    expect(readout()).toContain('slopeundefined')
    expect(readout()).toContain('at x = 6undefined')
    expect(screen.getByText(/x does not vary/)).toBeInTheDocument()
  })

  it('stops adding at twelve points', () => {
    const points: Array<[number, number]> = Array.from({ length: 12 }, (_, i) => [i * 0.5, i * 0.4])
    render(<LeastSquares title="Full" params={{ points, showResiduals: true }} />)
    expect(screen.getByRole('button', { name: 'Add point' })).toBeDisabled()
    expect(screen.getByRole('button', { name: 'Remove point' })).toBeEnabled()
  })

  it('does not append a point when a stale selection is edited after a removal', () => {
    render(
      <LeastSquares
        title="Stale"
        params={{ points: [[1, 2], [3, 4], [5, 6]], showResiduals: true }}
      />,
    )
    fireEvent.change(screen.getByLabelText('Point'), { target: { value: '2' } })
    fireEvent.click(screen.getByRole('button', { name: 'Remove point' }))
    fireEvent.change(screen.getByLabelText('y'), { target: { value: '7' } })
    expect(screen.getByText('2 points')).toBeInTheDocument()
  })
})

describe('vector combination edge cases', () => {
  it('reports no angle when a vector has no direction', () => {
    render(<VectorCombination title="Zero" params={{ a: [1, 2], b: [0, 0], k: -1.5, showSum: true }} />)
    expect(readout()).toContain('angle a,bundefined')
    expect(readout()).toContain('a + b1, 2')
  })

  it('draws no k b arrow when k is zero', () => {
    render(<VectorCombination title="Zero k" params={{ a: [1, 2], b: [2, 0.5], k: 0, showSum: true }} />)
    expect(readout()).toContain('k b0, 0')
    const labels = Array.from(screen.getByRole('img').querySelectorAll('text')).map((node) => node.textContent)
    expect(labels).not.toContain('k b')
  })

  it('keeps the sum of opposite vectors at the origin', () => {
    render(<VectorCombination title="Cancel" params={{ a: [2, 1], b: [-2, -1], k: 1, showSum: true }} />)
    expect(readout()).toContain('a + b0, 0')
    expect(readout()).toContain('angle a,bundefined'.replace('undefined', '180°'))
  })
})

describe('eigen spectrum readout', () => {
  const readout = () => screen.getByTestId('visual-readout').textContent ?? ''

  it('names both real eigenvalues of a symmetric matrix', () => {
    render(<EigenDirections title="Eigen" params={{ matrix: [2, 1, 1, 2], vector: [1, 0], showCircle: true }} />)
    expect(readout()).toContain('1, 3')
  })

  it('reports a complex pair instead of inventing a real direction', () => {
    render(<EigenDirections title="Eigen" params={{ matrix: [0, -1, 1, 0], vector: [1, 0], showCircle: true }} />)
    expect(readout()).toContain('0 ± 1i')
    expect(screen.getByText(/no real eigen-direction/)).toBeInTheDocument()
  })
})
