import { describe, expect, it, vi, afterEach } from 'vitest'
import { fireEvent, render, screen } from '@testing-library/react'
import MatrixTransform2D from '../components/visual/templates/MatrixTransform2D'
import { PredictProvider } from '../components/visual/predict'

function readout() {
  return screen.getByTestId('visual-readout').textContent ?? ''
}

const params = (matrix: number[], extra: Record<string, unknown> = {}) => ({
  matrix,
  vector: [2, 1],
  showGrid: true,
  showBasis: true,
  ...extra,
})

afterEach(() => {
  vi.unstubAllGlobals()
  vi.restoreAllMocks()
})

describe('matrix figure: what the numbers say', () => {
  it('shows a finite inverse and a full image when A is invertible', () => {
    render(<MatrixTransform2D title="Invertible" params={params([2, 1, 1, 1])} />)
    expect(readout()).toContain('det A1')
    expect(readout()).toContain('A⁻¹[[1,-1],[-1,2]]')
    expect(readout()).toContain('unit squarefull parallelogram')
  })

  it('states that no inverse exists only when A is singular', () => {
    render(<MatrixTransform2D title="Singular" params={params([1, 2, 2, 4])} />)
    expect(readout()).toContain('det A0')
    expect(readout()).toContain('no inverse (singular)')
    expect(readout()).toContain('collapses to a line')
    expect(readout()).not.toContain('NaN')
  })

  it('calls a pixel-thin image thin without calling it singular', () => {
    // Determinant 0.005: at the figure's 41.9 px per unit the image is 0.21 px
    // high, so it is invisible, which is a fact about the drawing and not about
    // invertibility.
    render(<MatrixTransform2D title="Thin" params={params([1, 0, 0, 0.005])} />)
    expect(readout()).toContain('thinner than a pixel')
    expect(readout()).toContain('A⁻¹[[1,0],[0,200]]')
    expect(readout()).not.toContain('no inverse')
    expect(readout()).toContain('orientationpreserved')
    // The determinant is printed at four decimals so it does not read as zero
    // beside an inverse that exists.
    expect(readout()).toContain('det A0.005')
  })

  it('keeps the intermediate and composed results apart', () => {
    render(
      <MatrixTransform2D
        title="Composed"
        params={params([2, 0, 0, 1], { then: [2, 0, 0, 1], vector: [1, 2] })}
      />,
    )
    expect(readout()).toContain('|Av|2.83')
    expect(readout()).toContain('|BAv|4.47')
    expect(readout()).toContain('det BA4')
  })

  it('renders the authored two-option prediction before the drawing', () => {
    render(
      <PredictProvider
        predict={{
          question: 'Double only the first column. What happens to the area?',
          options: ['It doubles', 'It stays the same'],
          answer: 0,
          explain: 'Doubling one column doubles the determinant.',
        }}
      >
        <MatrixTransform2D title="Predict" params={params([1.2, 0.6, 0.4, 1.4])} />
      </PredictProvider>,
    )
    expect(screen.getByText('Predict first')).toBeInTheDocument()
    expect(screen.getAllByRole('radio')).toHaveLength(2)
  })
})

describe('matrix figure: controls', () => {
  it('names the active preset and falls back to Custom after an edit', () => {
    render(<MatrixTransform2D title="Presets" params={params([1.2, 0.6, 0.4, 1.4])} />)
    const preset = screen.getByLabelText('Preset matrix') as HTMLSelectElement
    expect(preset.value).toBe('')

    fireEvent.change(preset, { target: { value: 'Project' } })
    expect(preset.value).toBe('Project')
    expect(readout()).toContain('orientationcollapsed')
    expect(readout()).toContain('area0')

    fireEvent.change(screen.getByLabelText('a'), { target: { value: '0.5' } })
    expect(preset.value).toBe('')
    expect(screen.getByRole('option', { name: 'Custom' })).toBeInTheDocument()
  })

  it('walks the column combination one term at a time', () => {
    render(<MatrixTransform2D title="Combination" params={params([1.2, 0.6, 0.4, 1.4])} />)
    const step = screen.getByRole('button', { name: 'Step through the combination that builds Av' })
    // Finished picture on load.
    expect(screen.getByText('Av', { selector: 'span' })).toBeInTheDocument()
    fireEvent.click(step)
    expect(screen.getByText('v only')).toBeInTheDocument()
    fireEvent.click(step)
    expect(screen.getByText('x·A e₁')).toBeInTheDocument()
    fireEvent.click(step)
    expect(screen.getByText('+ y·A e₂')).toBeInTheDocument()
  })

  it('keeps the inverse visible while the matrix is edited', () => {
    render(<MatrixTransform2D title="Edit" params={params([2, 1, 1, 1])} />)
    expect(readout()).toContain('A⁻¹[[1,-1],[-1,2]]')
    fireEvent.change(screen.getByLabelText('d'), { target: { value: '0' } })
    // det = -1, still invertible, and the inverse follows.
    expect(readout()).toContain('det A-1')
    expect(readout()).toContain('A⁻¹[[0,1],[1,-2]]')
  })
})

describe('matrix figure: animation lifecycle', () => {
  function stubFrames() {
    const request = vi.fn(() => 1)
    const cancel = vi.fn()
    vi.stubGlobal('requestAnimationFrame', request)
    vi.stubGlobal('cancelAnimationFrame', cancel)
    return { request, cancel }
  }

  const time = () => (screen.getByLabelText('t') as HTMLInputElement).value

  it('schedules one loop on play and cancels it on pause', () => {
    const { request, cancel } = stubFrames()
    render(<MatrixTransform2D title="Play" params={params([1.2, 0.6, 0.4, 1.4])} />)
    expect(request).not.toHaveBeenCalled()
    expect(time()).toBe('1')

    fireEvent.click(screen.getByRole('button', { name: 'Play' }))
    expect(request).toHaveBeenCalledTimes(1)
    expect(time()).toBe('0')

    fireEvent.click(screen.getByRole('button', { name: 'Pause' }))
    expect(cancel).toHaveBeenCalled()
    expect(screen.getByRole('button', { name: 'Play' })).toBeInTheDocument()
  })

  it('cancels the frame on unmount', () => {
    const { cancel } = stubFrames()
    const view = render(<MatrixTransform2D title="Unmount" params={params([1.2, 0.6, 0.4, 1.4])} />)
    fireEvent.click(screen.getByRole('button', { name: 'Play' }))
    view.unmount()
    expect(cancel).toHaveBeenCalled()
  })

  it('stops playback and returns to A when the matrix is edited', () => {
    const { cancel } = stubFrames()
    render(<MatrixTransform2D title="Edit" params={params([1.2, 0.6, 0.4, 1.4])} />)
    fireEvent.click(screen.getByRole('button', { name: 'Play' }))
    expect(screen.getByRole('button', { name: 'Pause' })).toBeInTheDocument()

    fireEvent.change(screen.getByLabelText('a'), { target: { value: '2' } })
    expect(cancel).toHaveBeenCalled()
    expect(screen.getByRole('button', { name: 'Play' })).toBeInTheDocument()
    expect(time()).toBe('1')
  })

  it('stops playback when a preset is chosen', () => {
    stubFrames()
    render(<MatrixTransform2D title="Preset" params={params([1.2, 0.6, 0.4, 1.4])} />)
    fireEvent.click(screen.getByRole('button', { name: 'Play' }))
    fireEvent.change(screen.getByLabelText('Preset matrix'), { target: { value: 'Shear' } })
    expect(screen.getByRole('button', { name: 'Play' })).toBeInTheDocument()
    expect(time()).toBe('1')
  })

  it('keeps playing when the vector changes', () => {
    stubFrames()
    render(<MatrixTransform2D title="Vector" params={params([1.2, 0.6, 0.4, 1.4])} />)
    fireEvent.click(screen.getByRole('button', { name: 'Play' }))
    fireEvent.change(screen.getByLabelText('vx'), { target: { value: '1' } })
    expect(screen.getByRole('button', { name: 'Pause' })).toBeInTheDocument()
  })

  it('scrubbing cancels playback and sets the time directly', () => {
    const { cancel } = stubFrames()
    render(<MatrixTransform2D title="Scrub" params={params([1.2, 0.6, 0.4, 1.4])} />)
    fireEvent.click(screen.getByRole('button', { name: 'Play' }))
    fireEvent.change(screen.getByLabelText('t'), { target: { value: '0.5' } })
    expect(cancel).toHaveBeenCalled()
    expect(time()).toBe('0.5')
    expect(screen.getByRole('button', { name: 'Play' })).toBeInTheDocument()
  })

  it('Show A returns to the endpoint without playing', () => {
    stubFrames()
    render(<MatrixTransform2D title="Show A" params={params([1.2, 0.6, 0.4, 1.4])} />)
    fireEvent.change(screen.getByLabelText('t'), { target: { value: '0.3' } })
    fireEvent.click(screen.getByRole('button', { name: 'Show A' }))
    expect(time()).toBe('1')
  })

  it('never schedules a frame when the learner asks for reduced motion', () => {
    const { request } = stubFrames()
    Object.defineProperty(window, 'matchMedia', {
      configurable: true,
      value: (query: string) => ({
        matches: query.includes('reduce'),
        media: query,
        onchange: null,
        addListener: () => {},
        removeListener: () => {},
        addEventListener: () => {},
        removeEventListener: () => {},
        dispatchEvent: () => false,
      }),
    })
    render(<MatrixTransform2D title="Reduced" params={params([1.2, 0.6, 0.4, 1.4])} />)
    fireEvent.click(screen.getByRole('button', { name: 'Play' }))
    expect(request).not.toHaveBeenCalled()
    // The endpoint is shown instead, so the button is not a dead control.
    expect(time()).toBe('1')
    expect(screen.getByRole('button', { name: 'Play' })).toBeInTheDocument()
  })
})
