import { describe, expect, it } from 'vitest'
import { fireEvent, render, screen } from '@testing-library/react'
import UnitCircle from '../components/visual/templates/UnitCircle'

const readout = () => screen.getByTestId('visual-readout').textContent ?? ''
const angleSlider = () => screen.getByLabelText('Angle') as HTMLInputElement
const commonAngle = () => screen.getByLabelText('Common angle') as HTMLSelectElement
const setAngle = (value: number) => fireEvent.change(angleSlider(), { target: { value: String(value) } })
const svg = () => screen.getByRole('img') as unknown as SVGSVGElement

describe('unit circle and sinusoid', () => {
  it('links the default angle to the point and the two values', () => {
    render(<UnitCircle title="Unit circle" params={{}} />)
    expect(readout()).toContain('angle45°')
    expect(readout()).toContain('angle on circle45°')
    expect(readout()).toContain('radians0.785 rad')
    expect(readout()).toContain('point(0.707, 0.707)')
    expect(readout()).toContain('cos θ0.7071')
    expect(readout()).toContain('sin θ0.7071')
    expect(readout()).toContain('quadrantI')
    expect(readout()).toContain('reference angle45°')
    expect(readout()).toContain('sin² + cos²1')
    expect(commonAngle().value).toBe('45')
    expect(screen.getByRole('img', { name: 'Unit circle at 45° with its sine and cosine graph' })).toBeTruthy()
  })

  it('reads the cardinal point at 90 degrees', () => {
    render(<UnitCircle title="Unit circle" params={{}} />)
    setAngle(90)
    expect(readout()).toContain('point(0, 1)')
    expect(readout()).toContain('cos θ0')
    expect(readout()).toContain('sin θ1')
    expect(readout()).toContain('tan θundefined')
    expect(readout()).toContain('quadranton an axis')
    expect(readout()).toContain('reference angle90°')
  })

  it('selects a common angle and updates the point', () => {
    render(<UnitCircle title="Unit circle" params={{}} />)
    fireEvent.change(commonAngle(), { target: { value: '180' } })
    expect(angleSlider().value).toBe('180')
    expect(readout()).toContain('point(-1, 0)')
    expect(readout()).toContain('cos θ-1')
    expect(readout()).toContain('sin θ0')
    expect(readout()).toContain('quadranton an axis')
  })

  it('keeps a negative angle and shows its positive coterminal angle', () => {
    render(<UnitCircle title="Unit circle" params={{}} />)
    setAngle(-45)
    expect(readout()).toContain('angle-45°')
    expect(readout()).toContain('angle on circle315°')
    expect(readout()).toContain('point(0.707, -0.707)')
    expect(readout()).toContain('quadrantIV')
    expect(commonAngle().value).toBe('315')
  })

  it('shows a custom angle when no common angle matches', () => {
    render(<UnitCircle title="Unit circle" params={{}} />)
    setAngle(37)
    expect(commonAngle().value).toBe('')
    expect(screen.getByRole('option', { name: 'Custom' })).toBeTruthy()
  })

  it('switches the graph curves without hiding the readout values', () => {
    render(<UnitCircle title="Unit circle" params={{}} />)
    expect(document.querySelector('[data-part="sine-curve"]')).not.toBeNull()
    expect(document.querySelector('[data-part="cosine-curve"]')).not.toBeNull()
    fireEvent.click(screen.getByLabelText('Show sine'))
    expect(document.querySelector('[data-part="sine-curve"]')).toBeNull()
    expect(document.querySelector('[data-part="cosine-curve"]')).not.toBeNull()
    expect(readout()).toContain('sin θ0.7071')

    fireEvent.click(screen.getByLabelText('Show cosine'))
    expect(document.querySelector('[data-part="cosine-curve"]')).toBeNull()
    expect(readout()).toContain('cos θ0.7071')
  })

  it('marks the current angle on both functions', () => {
    render(<UnitCircle title="Unit circle" params={{}} />)
    setAngle(30)
    const sineMarker = document.querySelector('[data-part="sine-marker"]')
    const cosineMarker = document.querySelector('[data-part="cosine-marker"]')
    expect(sineMarker).not.toBeNull()
    expect(cosineMarker).not.toBeNull()
    expect(sineMarker?.getAttribute('cy')).not.toBe(cosineMarker?.getAttribute('cy'))
  })

  it('drags the point around the circle and updates the angle', () => {
    render(<UnitCircle title="Unit circle" params={{}} />)
    const drawing = svg()
    fireEvent.pointerDown(drawing, { button: 0, pointerId: 1, clientX: 170, clientY: 56 })
    fireEvent.pointerMove(drawing, { pointerId: 1, clientX: 170, clientY: 56 })
    fireEvent.pointerUp(drawing, { pointerId: 1 })
    expect(readout()).toContain('angle90°')
    expect(readout()).toContain('point(0, 1)')
  })

  it('projects a drag outside the circle back onto it', () => {
    render(<UnitCircle title="Unit circle" params={{}} />)
    const drawing = svg()
    fireEvent.pointerDown(drawing, { button: 0, pointerId: 2, clientX: 300, clientY: 220 })
    fireEvent.pointerMove(drawing, { pointerId: 2, clientX: 300, clientY: 220 })
    fireEvent.pointerUp(drawing, { pointerId: 2 })
    expect(readout()).toContain('sin² + cos²1')
    expect(readout()).not.toContain('NaN')
  })

  it('keeps the control keyboard reachable', () => {
    render(<UnitCircle title="Unit circle" params={{}} />)
    const slider = angleSlider()
    expect(slider.type).toBe('range')
    expect(slider.tabIndex).toBe(0)
    slider.focus()
    fireEvent.keyDown(slider, { key: 'ArrowRight' })
    expect(readout()).not.toContain('NaN')
  })

  it('resets the angle and both graph toggles', () => {
    render(<UnitCircle title="Unit circle" params={{}} />)
    setAngle(210)
    fireEvent.click(screen.getByLabelText('Show sine'))
    fireEvent.click(screen.getByText('Reset'))
    expect(readout()).toContain('angle45°')
    expect(readout()).toContain('point(0.707, 0.707)')
    expect(document.querySelector('[data-part="sine-curve"]')).not.toBeNull()
  })

  it('does not produce non-finite values at the angle limits', () => {
    render(<UnitCircle title="Unit circle" params={{ angle: 360 }} />)
    expect(readout()).not.toContain('NaN')
    expect(readout()).not.toContain('Infinity')
    setAngle(-360)
    expect(readout()).not.toContain('NaN')
    expect(readout()).not.toContain('Infinity')
    setAngle(270)
    expect(readout()).toContain('tan θundefined')
    expect(readout()).toContain('point(0, -1)')
  })
})
