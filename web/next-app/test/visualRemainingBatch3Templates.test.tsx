import { fireEvent, render, screen } from '@testing-library/react'
import { describe, expect, it } from 'vitest'
import FunctionGraph from '../components/visual/templates/FunctionGraph'

const readout = () => screen.getByTestId('visual-readout').textContent ?? ''
const slider = (label: string) => screen.getByLabelText(label) as HTMLInputElement
const setSlider = (label: string, value: number) =>
  fireEvent.change(slider(label), { target: { value: String(value) } })
const choose = (label: string, value: string) =>
  fireEvent.change(screen.getByLabelText(label), { target: { value } })

describe('function graph template', () => {
  it('starts with the quadratic and reports its vertex and ends', () => {
    render(
      <FunctionGraph
        title="Function graph"
        params={{ family: 'quadratic', a: 1, h: 0, k: 0, showInverse: false }}
      />,
    )
    const text = readout()
    expect(text).toContain('Quadratic')
    expect(text).toContain('x²')
    expect(text).toContain('left +∞, right +∞')
    expect(text).toContain('[0, ∞)')
    expect(text).toContain('(0, 0)')
    expect(text).toContain('not a function on this domain')
    expect(screen.getByRole('img')).toHaveAccessibleName(/Graph of x²/)
  })

  it('switches to a negative cubic and reverses the ends', () => {
    render(<FunctionGraph title="Function graph" params={{ family: 'quadratic' }} />)
    choose('Family', 'cubic')
    setSlider('Scale a', -1)
    const text = readout()
    expect(text).toContain('−x³')
    expect(text).toContain('left +∞, right −∞')
    expect(text).toContain('all real y')
  })

  it('shows the exponential inverse as a function', () => {
    render(
      <FunctionGraph
        title="Function graph"
        params={{ family: 'exponential', a: 1, h: 0, k: 0, showInverse: true }}
      />,
    )
    const text = readout()
    expect(text).toContain('2^x')
    expect(text).toContain('inversea function')
    expect(text).toContain('log₂')
    expect(screen.getAllByText('inverse').length).toBeGreaterThan(0)
  })

  it('keeps the quadratic inverse labelled as not a function', () => {
    render(
      <FunctionGraph
        title="Function graph"
        params={{ family: 'quadratic', a: 1, h: 0, k: 0, showInverse: true }}
      />,
    )
    expect(readout()).toContain('not a function on this domain')
  })

  it('moves the vertex with h and k', () => {
    render(<FunctionGraph title="Function graph" params={{ family: 'quadratic' }} />)
    setSlider('Shift h', 1)
    setSlider('Lift k', -2)
    const text = readout()
    expect(text).toContain('(x − 1)² − 2')
    expect(text).toContain('(1, -2)')
    expect(text).toContain('[-2, ∞)')
  })

  it('toggles the reflected curve and keeps the readout finite', () => {
    render(<FunctionGraph title="Function graph" params={{ family: 'cubic' }} />)
    fireEvent.click(screen.getByLabelText('Show inverse'))
    expect(screen.getAllByText('inverse').length).toBeGreaterThan(0)
    expect(readout()).not.toContain('NaN')
    expect(readout()).not.toContain('Infinity')
  })

  it('resets every control', () => {
    render(
      <FunctionGraph
        title="Function graph"
        params={{ family: 'cubic', a: -1, h: 1, k: -2, showInverse: true }}
      />,
    )
    setSlider('Scale a', 0.5)
    fireEvent.click(screen.getByRole('button', { name: 'Reset' }))
    const text = readout()
    expect(text).toContain('Cubic')
    expect(text).toContain('−(x − 1)³ − 2')
    expect(text).toContain('inversea function')
  })

  it('is keyboard reachable', () => {
    render(<FunctionGraph title="Function graph" params={{ family: 'quadratic' }} />)
    const input = slider('Scale a')
    expect(input.type).toBe('range')
    input.focus()
    fireEvent.keyDown(input, { key: 'ArrowRight' })
    expect(readout()).not.toContain('NaN')
  })

  it('stays finite at extreme allowed coefficients', () => {
    render(
      <FunctionGraph
        title="Function graph"
        params={{ family: 'exponential', a: 2, h: -2, k: 2, showInverse: true }}
      />,
    )
    const text = readout()
    expect(text).not.toContain('NaN')
    expect(text).not.toContain('Infinity')
    expect(text).toContain('2·2^(x + 2) + 2')
  })
})
