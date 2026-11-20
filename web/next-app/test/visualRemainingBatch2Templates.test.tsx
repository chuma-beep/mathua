import { fireEvent, render, screen } from '@testing-library/react'
import { describe, expect, it } from 'vitest'
import ProbabilityTree from '../components/visual/templates/ProbabilityTree'
import FractionLine from '../components/visual/templates/FractionLine'

const readout = () => screen.getByTestId('visual-readout').textContent ?? ''
const setSlider = (label: string, value: number) =>
  fireEvent.change(screen.getByLabelText(label), { target: { value: String(value) } })

describe('probability tree figure', () => {
  it('shows the leaves and the posterior from the same probabilities', () => {
    render(<ProbabilityTree title="Tree" params={{}} />)
    expect(readout()).toContain('0.1085')
    expect(readout()).toContain('0.0876')
    expect(readout()).toContain('0.0095')
    expect(readout()).toContain('0.891')
    expect(screen.getByRole('img').getAttribute('aria-label')).toContain('Probability tree')
  })

  it('updates the posterior when the prior changes', () => {
    render(<ProbabilityTree title="Tree" params={{ prior: 0.5 }} />)
    const before = readout()
    setSlider('Prior', 0.9)
    expect(readout()).not.toBe(before)
    expect(readout()).toContain('P(D|+)')
  })

  it('keeps a zero prior at zero', () => {
    render(<ProbabilityTree title="Tree" params={{ prior: 0 }} />)
    expect(readout()).toContain('P(D|+)0')
    expect(readout()).toContain('P(+)0.1')
  })
})

describe('fraction line figure', () => {
  it('places both fractions and reports their comparison', () => {
    render(<FractionLine title="Fractions" params={{}} />)
    expect(readout()).toContain('3/4')
    expect(readout()).toContain('0.75')
    expect(readout()).toContain('1/2')
    expect(readout()).toContain('comparison>')
    expect(readout()).toContain('sum5/4')
    expect(screen.getByRole('img').getAttribute('aria-label')).toContain('3/4')
  })

  it('recognises equivalent fractions', () => {
    render(<FractionLine title="Fractions" params={{ numerator: 1, denominator: 2, otherNumerator: 2, otherDenominator: 4 }} />)
    expect(readout()).toContain('comparison=')
    expect(readout()).toContain('equivalentyes')
  })

  it('updates when the first numerator changes', () => {
    render(<FractionLine title="Fractions" params={{}} />)
    setSlider('Numerator', 1)
    expect(readout()).toContain('first fraction1/4')
    expect(readout()).toContain('comparison<')
  })

  it('keeps denominators positive and the readout finite', () => {
    render(<FractionLine title="Fractions" params={{ denominator: 8, otherDenominator: 8 }} />)
    expect(readout()).not.toContain('NaN')
    expect(readout()).not.toContain('Infinity')
  })
})
