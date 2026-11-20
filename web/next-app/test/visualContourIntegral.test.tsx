import { describe, expect, it } from 'vitest'
import { fireEvent, render, screen } from '@testing-library/react'
import ContourIntegral from '../components/visual/templates/ContourIntegral'

function readout() {
  return screen.getByTestId('visual-readout').textContent ?? ''
}

describe('contour integral template', () => {
  it('accumulates 2 pi i over the unit circle and names it in the readout', () => {
    render(
      <ContourIntegral
        title="Contour"
        params={{ path: 'circle', fn: 'inv', reverse: false, radius: 1, centerX: 0, t: 1 }}
      />,
    )
    expect(readout()).toContain('0 + 6.28i')
    expect(readout()).toContain('poleclear')
  })

  it('negates the integral when the path is traversed backwards', () => {
    render(
      <ContourIntegral
        title="Contour"
        params={{ path: 'circle', fn: 'inv', reverse: false, radius: 1, centerX: 0, t: 1 }}
      />,
    )
    fireEvent.click(screen.getByLabelText('Traverse backwards'))
    expect(readout()).toContain('0 − 6.28i')
  })

  it('gives zero for a circle that does not enclose the pole', () => {
    render(
      <ContourIntegral
        title="Contour"
        params={{ path: 'circle', fn: 'inv', reverse: false, radius: 1, centerX: 1.6, t: 1 }}
      />,
    )
    expect(readout()).toContain('total0 + 0i')
  })

  it('reports the crossing instead of a value when the path hits the pole', () => {
    render(
      <ContourIntegral
        title="Contour"
        params={{ path: 'line', fn: 'inv', reverse: false, radius: 1, centerX: 0, t: 1 }}
      />,
    )
    expect(readout()).toContain('undefined')
    expect(readout()).toContain('poleon the path')
    expect(screen.getByText(/path passes through 0/)).toBeInTheDocument()
  })
})
