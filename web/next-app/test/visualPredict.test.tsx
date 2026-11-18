import { describe, expect, it } from 'vitest'
import { fireEvent, render, screen } from '@testing-library/react'
import FigureFrame from '../components/visual/FigureFrame'
import { PredictProvider } from '../components/visual/predict'
import type { PredictSpec } from '../lib/api'

const predict: PredictSpec = {
  question: 'Which is steeper at the point?',
  options: ['the tangent', 'the secant'],
  answer: 0,
  explain: 'The tangent is the limit of the secant as h shrinks.',
}

function frame() {
  return render(
    <PredictProvider predict={predict}>
      <FigureFrame title="Tangent" readout={[]}>
        <div data-testid="drawing" />
      </FigureFrame>
    </PredictProvider>,
  )
}

describe('predict before you reveal', () => {
  it('asks the question before the figure is explored, and hides the answer until chosen', () => {
    frame()
    expect(screen.getByText('Which is steeper at the point?')).toBeTruthy()
    expect(screen.getByText('the tangent')).toBeTruthy()
    expect(screen.queryByText(/limit of the secant/)).toBeNull()
    expect(screen.getByTestId('drawing')).toBeTruthy()
  })

  it('tells a wrong answer it is wrong and still explains why', () => {
    frame()
    fireEvent.click(screen.getByLabelText('the secant'))
    expect(screen.getByText(/Not quite/)).toBeTruthy()
    expect(screen.getByText(/limit of the secant/)).toBeTruthy()
    // The figure stays usable: the prediction never blocks the drawing.
    expect(screen.getByTestId('drawing')).toBeTruthy()
  })

  it('confirms a correct answer', () => {
    frame()
    fireEvent.click(screen.getByLabelText('the tangent'))
    expect(screen.getByText(/Yes\./)).toBeTruthy()
  })

  it('renders no question when the visual carries none', () => {
    render(
      <PredictProvider predict={null}>
        <FigureFrame title="Plain" readout={[]}>
          <div />
        </FigureFrame>
      </PredictProvider>,
    )
    expect(screen.queryByText('Predict first')).toBeNull()
  })
})
