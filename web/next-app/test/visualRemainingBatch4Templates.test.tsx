import { describe, expect, it } from 'vitest'
import { fireEvent, render, screen } from '@testing-library/react'
import NumberLineLab from '../components/visual/templates/NumberLineLab'
import ModularLab from '../components/visual/templates/ModularLab'
import GroupTable from '../components/visual/templates/GroupTable'
import DiscreteLab from '../components/visual/templates/DiscreteLab'
import OptimizationLab from '../components/visual/templates/OptimizationLab'
import AnalysisLab from '../components/visual/templates/AnalysisLab'
import ModelFit from '../components/visual/templates/ModelFit'

const readout = () => screen.getByTestId('visual-readout').textContent ?? ''

describe('arithmetic and algebra lab', () => {
  it('shows signed order and a balanced equation', () => {
    const { unmount } = render(<NumberLineLab params={{ mode: 'signed', a: -4, b: 2 }} />)
    expect(readout()).toContain('-4 is left of 2')
    unmount()
    render(<NumberLineLab params={{ mode: 'balance', guess: 2, constant: 3, total: 8 }} />)
    expect(readout()).toContain('solution5')
  })
})

describe('number theory lab', () => {
  it('shows a residue and a Euclidean gcd', () => {
    const { unmount } = render(<ModularLab params={{ mode: 'residue', a: 17, n: 5 }} />)
    expect(readout()).toContain('2')
    unmount()
    render(<ModularLab params={{ mode: 'euclid', a: 30, b: 17 }} />)
    expect(readout()).toContain('gcd')
    expect(readout()).toContain('1')
  })
})

describe('group table', () => {
  it('shows Z4 order and generated subgroup', () => {
    render(<GroupTable params={{ group: 'z4', element: 1 }} />)
    expect(readout()).toContain('Z4')
    expect(readout()).toContain('4')
    expect(readout()).toContain('abelian')
    expect(screen.getByRole('img').getAttribute('aria-label')).toContain('Z4')
  })
})

describe('discrete lab', () => {
  it('shows a truth table and a binomial count', () => {
    const { unmount } = render(<DiscreteLab params={{ mode: 'truth', op: 'implies' }} />)
    expect(readout()).toContain('truth rows4')
    expect(readout()).toContain('true rows3')
    unmount()
    render(<DiscreteLab params={{ mode: 'combinations', n: 5, k: 2 }} />)
    expect(readout()).toContain('10')
  })
})

describe('optimization lab', () => {
  it('shows a descent path and a feasible optimum', () => {
    const { unmount } = render(
      <OptimizationLab params={{ mode: 'descent', startX: 2.5, startY: 2, rate: 0.15, steps: 6 }} />,
    )
    expect(readout()).toContain('x² + 2y²')
    expect(readout()).toContain('running')
    unmount()
    render(<OptimizationLab params={{ mode: 'feasible', candidateX: 1, candidateY: 1 }} />)
    expect(readout()).toContain('feasible')
    expect(readout()).toContain('8')
  })
})

describe('analysis lab', () => {
  it('shows an epsilon threshold and a metric ball', () => {
    const { unmount } = render(<AnalysisLab params={{ mode: 'sequence', sequence: 'reciprocal', epsilon: 0.2 }} />)
    expect(readout()).toContain('1/n')
    expect(readout()).toContain('threshold')
    unmount()
    render(<AnalysisLab params={{ mode: 'metric', metric: 'manhattan', radius: 2 }} />)
    expect(readout()).toContain('manhattan ball')
    expect(readout()).toContain('lattice points')
  })
})

describe('machine-learning fit', () => {
  it('shows a fitted line and a logistic probability', () => {
    const { unmount } = render(<ModelFit params={{ mode: 'regression', rate: 0.1, steps: 8 }} />)
    expect(readout()).toContain('linear regression')
    expect(readout()).toContain('MSE')
    unmount()
    render(<ModelFit params={{ mode: 'logistic', threshold: 0 }} />)
    expect(readout()).toContain('probability at x = 1')
    expect(readout()).toContain('0.88')
  })

  it('keeps a keyboard-operable rate slider', () => {
    render(<ModelFit params={{ mode: 'regression', rate: 0.1, steps: 8 }} />)
    const slider = screen.getByLabelText('Learning rate') as HTMLInputElement
    fireEvent.change(slider, { target: { value: '0.2' } })
    expect(readout()).toContain('linear regression')
  })
})
