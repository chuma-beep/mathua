import { fireEvent, render, screen } from '@testing-library/react'
import { describe, expect, it } from 'vitest'
import InferenceLab from '../components/visual/templates/InferenceLab'
import PhasePortrait from '../components/visual/templates/PhasePortrait'
import GraphNetwork from '../components/visual/templates/GraphNetwork'
import OptimizationLab from '../components/visual/templates/OptimizationLab'

const readout = () => screen.getByTestId('visual-readout').textContent ?? ''

const choose = (label: string, value: string) =>
  fireEvent.change(screen.getByLabelText(label), { target: { value } })

const setSlider = (label: string, value: number) =>
  fireEvent.change(screen.getByLabelText(label), { target: { value: String(value) } })

describe('InferenceLab', () => {
  const params = { dist: 'z' as const, df: 10, alpha: 0.05, effect: 0.8, observed: 1.9 }

  it('draws the null and alternative curves and reports the test', () => {
    render(<InferenceLab title="Inference" params={params} />)
    const text = readout()
    expect(text).toContain('null distributionstandard normal')
    expect(text).toContain('critical value1.96')
    expect(text).toContain('p-value0.0574')
    expect(text).toContain('decisionfail to reject')
    expect(text).toContain('Type II error')
    expect(screen.getByRole('img').getAttribute('aria-label')).toContain('Two-sided z test')
  })

  it('moves the observation across the critical value', () => {
    render(<InferenceLab title="Inference" params={params} />)
    setSlider('Observed statistic', 2.5)
    expect(readout()).toContain('p-value0.0124')
    expect(readout()).toContain('decisionreject')
  })

  it('switches to the t distribution and exposes its degrees of freedom', () => {
    render(<InferenceLab title="Inference" params={params} />)
    choose('Distribution', 't')
    expect(readout()).toContain('null distributiont with 10 df')
    expect(screen.getByLabelText('Degrees of freedom')).toBeTruthy()
    setSlider('Degrees of freedom', 5)
    expect(readout()).toContain('t with 5 df')
  })

  it('changes alpha and the interval', () => {
    render(<InferenceLab title="Inference" params={params} />)
    setSlider('Alpha', 0.1)
    expect(readout()).toContain('critical value1.645')
    expect(readout()).toContain('alpha0.1')
  })
})

describe('PhasePortrait', () => {
  const params = {
    mode: 'line' as const,
    lineModel: 'logistic' as const,
    system: 'saddle' as const,
    showTrajectories: true,
  }

  it('shows the phase line and its equilibria', () => {
    render(<PhasePortrait title="Phase" params={params} />)
    expect(readout()).toContain("equationx' = x(1 − x)")
    expect(readout()).toContain('equilibria0 unstable, 1 stable')
    expect(screen.getByRole('img').getAttribute('aria-label')).toContain('Phase line')
  })

  it('switches to the saddle phase plane', () => {
    render(<PhasePortrait title="Phase" params={params} />)
    choose('View', 'plane')
    expect(readout()).toContain('systemSaddle')
    expect(readout()).toContain('eigenvalues2, -3')
    expect(readout()).toContain('classificationsaddle')
    expect(readout()).toContain('stabilityunstable')
  })

  it('switches the linear system and hides trajectories', () => {
    render(<PhasePortrait title="Phase" params={params} />)
    choose('View', 'plane')
    choose('System', 'spiral')
    expect(readout()).toContain('classificationstable spiral')
    expect(readout()).toContain('eigenvalues-0.5 ± 1i')
    fireEvent.click(screen.getByLabelText('Show trajectories'))
    expect(readout()).toContain('trajectorieshidden')
  })

  it('shows a threshold phase line with two equilibria', () => {
    render(<PhasePortrait title="Phase" params={params} />)
    choose('Equation', 'threshold')
    expect(readout()).toContain('equilibria-1 stable, 1 unstable')
  })
})

describe('GraphNetwork extensions', () => {
  const params = { algorithm: 'bfs' as const, start: 'A', showTree: true }

  it('orders a DAG by repeatedly removing a source', () => {
    render(<GraphNetwork title="Graph" params={params} />)
    choose('Algorithm', 'topological')
    expect(readout()).toContain('orderA, B, C, D, E')
    expect(readout()).toContain('statusok')
  })

  it('reports the maximum flow and the minimum cut', () => {
    render(<GraphNetwork title="Graph" params={params} />)
    choose('Algorithm', 'maxflow')
    expect(readout()).toContain('flow5')
    expect(readout()).toContain('minimum cutS | A, B, T')
    expect(screen.getByRole('img').getAttribute('aria-label')).toContain('flow network')
  })

  it('keeps the traversal algorithms working', () => {
    render(<GraphNetwork title="Graph" params={params} />)
    expect(readout()).toContain('orderA, B, D, C, E, F')
    choose('Algorithm', 'dfs')
    expect(readout()).toContain('orderA, B, C, F, E, D')
  })
})

describe('OptimizationLab nonlinear modes', () => {
  const params = {
    mode: 'lagrange' as const,
    startX: 1.5,
    startY: 2,
    rate: 0.15,
    steps: 4,
    candidateX: 0.5,
    candidateY: 1,
    rhs: 'none' as const,
    delta: 0,
  }

  it('solves the equality-constrained problem at the stationary point', () => {
    render(<OptimizationLab title="Optimization" params={params} />)
    setSlider('Candidate x', 1)
    const text = readout()
    expect(text).toContain('lambda2')
    expect(text).toContain('statusstationary on the constraint')
    expect(text).toContain('optimum(1, 1) = 2')
  })

  it('shows that a non-stationary candidate has a nonzero along-constraint derivative', () => {
    render(<OptimizationLab title="Optimization" params={params} />)
    expect(readout()).toContain('statusnot stationary')
    expect(readout()).toContain('along constraint-1.414')
  })

  it('takes a Newton step from the current guess', () => {
    render(<OptimizationLab title="Optimization" params={params} />)
    choose('Model', 'newton')
    setSlider('Steps', 1)
    const text = readout()
    expect(text).toContain('last x1.5')
    expect(text).toContain('next x1.2381')
    expect(text).toContain('f′(x)5.5')
    expect(text).toContain('f″(x)21')
  })

  it('stops Newton when the curvature is not positive', () => {
    render(<OptimizationLab title="Optimization" params={params} />)
    choose('Model', 'newton')
    setSlider('Start x', 0)
    expect(readout()).toContain('statusstationary')
    expect(readout()).toContain('note')
  })
})
