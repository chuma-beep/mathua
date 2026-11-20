import { fireEvent, render, screen } from '@testing-library/react'
import { describe, expect, it } from 'vitest'
import SinusoidExplorer from '../components/visual/templates/SinusoidExplorer'

const readout = () => screen.getByTestId('visual-readout').textContent ?? ''
const slider = (label: string) => screen.getByLabelText(label) as HTMLInputElement
const setSlider = (label: string, value: number) =>
  fireEvent.change(slider(label), { target: { value: String(value) } })
const setNumber = (label: string, value: number) =>
  fireEvent.change(screen.getByLabelText(label), { target: { value: String(value) } })

describe('sinusoid explorer', () => {
  it('starts from the sine wave and reports its four parameters', () => {
    render(<SinusoidExplorer title="The sinusoid" params={{}} />)
    const text = readout()
    expect(text).toContain('equationy = sin(x)')
    expect(text).toContain('amplitude1')
    expect(text).toContain('period2π')
    expect(text).toContain('frequency0.159 cycles/unit')
    expect(text).toContain('phase shift0')
    expect(text).toContain('maximum1')
    expect(text).toContain('minimum-1')
    expect(text).toContain('midline0')
    expect(screen.getByRole('img', { name: 'Graph of y = sin(x)' })).toBeTruthy()
  })

  it('changes the amplitude and the midline through A and D', () => {
    render(<SinusoidExplorer title="The sinusoid" params={{}} />)
    setSlider('Amplitude A', 2)
    setSlider('Vertical D', 1)
    const text = readout()
    expect(text).toContain('equationy = 2sin(x) + 1')
    expect(text).toContain('amplitude2')
    expect(text).toContain('maximum3')
    expect(text).toContain('minimum-1')
    expect(text).toContain('midline1')
  })

  it('turns the inside coefficient into a period', () => {
    render(<SinusoidExplorer title="The sinusoid" params={{}} />)
    setSlider('Coefficient B', 2)
    expect(readout()).toContain('periodπ')
    expect(readout()).toContain('frequency0.318 cycles/unit')
    expect(readout()).toContain('equationy = sin(2x)')
  })

  it('names the phase shift and its direction', () => {
    render(<SinusoidExplorer title="The sinusoid" params={{}} />)
    setNumber('Phase C value', Math.PI / 2)
    expect(readout()).toContain('phase shiftπ/2 right')
    expect(readout()).toContain('equationy = sin(x − π/2)')
    setNumber('Phase C value', -Math.PI / 2)
    expect(readout()).toContain('phase shiftπ/2 left')
    expect(readout()).toContain('equationy = sin(x + π/2)')
  })

  it('switches between sine and cosine', () => {
    render(<SinusoidExplorer title="The sinusoid" params={{}} />)
    fireEvent.change(screen.getByLabelText('Form'), { target: { value: 'cos' } })
    expect(readout()).toContain('formcosine')
    expect(readout()).toContain('equationy = cos(x)')
    expect(screen.getByRole('img', { name: 'Graph of y = cos(x)' })).toBeTruthy()
  })

  it('toggles the reference curve without changing the wave', () => {
    render(<SinusoidExplorer title="The sinusoid" params={{}} />)
    const toggle = screen.getByRole('switch', { name: 'Show reference' })
    expect(toggle.getAttribute('aria-checked')).toBe('true')
    fireEvent.click(toggle)
    expect(toggle.getAttribute('aria-checked')).toBe('false')
    expect(readout()).toContain('amplitude1')
  })

  it('keeps the amplitude positive for a negative A and reports the reflection', () => {
    render(<SinusoidExplorer title="The sinusoid" params={{}} />)
    setSlider('Amplitude A', -2)
    expect(readout()).toContain('amplitude2')
    expect(readout()).toContain('equationy = −2sin(x)')
  })

  it('uses the absolute value of B for the period', () => {
    render(<SinusoidExplorer title="The sinusoid" params={{}} />)
    setSlider('Coefficient B', 4)
    expect(readout()).toContain('periodπ/2')
  })

  it('keeps the keyboard path on the range input', () => {
    render(<SinusoidExplorer title="The sinusoid" params={{}} />)
    const amplitude = slider('Amplitude A')
    expect(amplitude.type).toBe('range')
    expect(amplitude.tabIndex).toBe(0)
    expect(amplitude.min).toBe('-3')
    expect(amplitude.max).toBe('3')
    expect(amplitude.step).toBe('0.1')
  })

  it('never prints NaN or Infinity at the slider extremes', () => {
    render(
      <SinusoidExplorer
        title="Extremes"
        params={{
          form: 'cos',
          amplitude: -3,
          frequency: 4,
          phase: 3.14,
          vertical: -2,
          showReference: true,
        }}
      />,
    )
    expect(readout()).not.toContain('NaN')
    expect(readout()).not.toContain('Infinity')
  })

  it('gives two mounted figures separate clip paths', () => {
    const { container } = render(
      <>
        <SinusoidExplorer title="First" params={{}} />
        <SinusoidExplorer title="Second" params={{ form: 'cos', amplitude: 2 }} />
      </>,
    )
    const ids = Array.from(container.querySelectorAll('clipPath')).map((node) => node.id)
    expect(ids.length).toBe(2)
    expect(new Set(ids).size).toBe(2)
  })
})
