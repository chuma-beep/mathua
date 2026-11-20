import { describe, expect, it } from 'vitest'
import { fireEvent, render, screen } from '@testing-library/react'
import OptimizationLab from '../components/visual/templates/OptimizationLab'

const readout = () => screen.getByTestId('visual-readout').textContent ?? ''

const choose = (label: string, value: string) => {
  fireEvent.change(screen.getByLabelText(label), { target: { value } })
}

const setSlider = (label: string, value: string) => {
  fireEvent.change(screen.getByLabelText(label), { target: { value } })
}

describe('constrained optimization figure', () => {
  it('shows the exact optimum, active set and dual multipliers', () => {
    render(<OptimizationLab params={{ mode: 'feasible', candidateX: 2, candidateY: 1 }} />)
    expect(readout()).toContain('maximise 3x + 2y')
    expect(readout()).toContain('KKT at candidateholds')
    expect(readout()).toContain('optimum(2, 1) = 8')
    expect(readout()).toContain('x + 2y ≤ 4: 0.333')
    expect(readout()).toContain('2x + y ≤ 5: 1.333')
    expect(readout()).toContain('dual value8')
    expect(readout()).toContain('primal-dual gap0')
    expect(readout()).toContain('shadow price—')
    expect(screen.getByRole('img').getAttribute('aria-label')).toContain('candidate at (2, 1)')
  })

  it('rejects an interior candidate', () => {
    render(<OptimizationLab params={{ mode: 'feasible', candidateX: 1, candidateY: 1 }} />)
    expect(readout()).toContain('KKT at candidatefails')
    expect(readout()).toContain('active at candidatenone')
    expect(readout()).toContain('feasible')
  })

  it('shows a violation for a candidate outside the polygon', () => {
    render(<OptimizationLab params={{ mode: 'feasible', candidateX: 2, candidateY: 1 }} />)
    setSlider('Candidate x', '3')
    setSlider('Candidate y', '3')
    expect(readout()).toContain('violates a constraint')
  })

  it('prices the first bound after a right-hand-side shift', () => {
    render(<OptimizationLab params={{ mode: 'feasible', candidateX: 2, candidateY: 1 }} />)
    choose('Shift bound', 'first')
    setSlider('Change δ', '0.2')
    expect(readout()).toContain('optimum(1.93, 1.13) = 8.067')
    expect(readout()).toContain('shadow pricex + 2y: 0.333 per unit')
    expect(readout()).toContain('allowable rangex + 2y: [2.5, 10]')
    expect(readout()).toContain('dual value8.067')
    expect(readout()).toContain('primal-dual gap0')
  })

  it('prices the second bound independently', () => {
    render(<OptimizationLab params={{ mode: 'feasible', candidateX: 2, candidateY: 1 }} />)
    choose('Shift bound', 'second')
    setSlider('Change δ', '0.2')
    expect(readout()).toContain('shadow price2x + y: 1.333 per unit')
    expect(readout()).toContain('allowable range2x + y: [2, 8]')
    expect(readout()).toContain('optimum(2.13, 0.93) = 8.267')
  })

  it('moves the active set when the candidate reaches a vertex', () => {
    render(<OptimizationLab params={{ mode: 'feasible', candidateX: 1, candidateY: 1 }} />)
    setSlider('Candidate x', '2')
    setSlider('Candidate y', '1')
    expect(readout()).toContain('active at candidatex + 2y')
    expect(readout()).toContain('KKT at candidateholds')
  })

  it('accepts a keyboard-operated shift and stays finite at both ends', () => {
    render(<OptimizationLab params={{ mode: 'feasible', candidateX: 2, candidateY: 1 }} />)
    const select = screen.getByLabelText('Shift bound') as HTMLSelectElement
    expect(select.tagName).toBe('SELECT')
    choose('Shift bound', 'first')
    setSlider('Change δ', '-1')
    expect(readout()).not.toContain('NaN')
    expect(readout()).not.toContain('Infinity')
    setSlider('Change δ', '1')
    expect(readout()).not.toContain('NaN')
    expect(readout()).not.toContain('Infinity')
  })

  it('resets the candidate and the shift', () => {
    render(<OptimizationLab params={{ mode: 'feasible', candidateX: 2, candidateY: 1 }} />)
    choose('Shift bound', 'first')
    setSlider('Change δ', '0.5')
    setSlider('Candidate x', '0')
    fireEvent.click(screen.getByText('Reset'))
    expect(readout()).toContain('candidate(2, 1)')
    expect(readout()).toContain('shadow price—')
  })

  it('keeps the gradient-descent mode intact', () => {
    render(<OptimizationLab params={{ mode: 'descent', startX: 2.5, startY: 2, rate: 0.15, steps: 6 }} />)
    expect(readout()).toContain('x² + 2y²')
    expect(readout()).toContain('running')
    expect(screen.getByLabelText('Learning rate')).toBeTruthy()
  })
})
