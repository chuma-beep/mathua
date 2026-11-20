import { fireEvent, render, screen } from '@testing-library/react'
import { describe, expect, it } from 'vitest'
import SlopeField from '../components/visual/templates/SlopeField'
import RootIteration from '../components/visual/templates/RootIteration'
import GraphNetwork from '../components/visual/templates/GraphNetwork'

const readout = () => screen.getByTestId('visual-readout').textContent ?? ''

const choose = (label: string, value: string) => {
  fireEvent.change(screen.getByLabelText(label), { target: { value } })
}

const setSlider = (label: string, value: number) => {
  fireEvent.change(screen.getByLabelText(label), { target: { value: String(value) } })
}

describe('slope field figure', () => {
  it('draws the actual field and reports the integrated curve', () => {
    render(<SlopeField title="Slope field" params={{}} />)
    expect(readout()).toContain("y' = y")
    expect(readout()).toContain('(0, 1)')
    expect(readout()).toContain('RK4')
    expect(readout()).toContain('0.1')
    expect(readout()).not.toContain('NaN')
    expect(screen.getByRole('img').getAttribute('aria-label')).toContain("y' = y")
    expect(document.querySelectorAll('line').length).toBeGreaterThan(100)
  })

  it('shows an equilibrium slope at a constant solution', () => {
    render(<SlopeField title="Slope field" params={{ field: 'logistic', y0: 1 }} />)
    expect(readout()).toContain('y(1 − y)')
    expect(readout()).toContain('slope at start0')
    expect(readout()).toContain('ok')
  })

  it('switches integrator and step size', () => {
    render(<SlopeField title="Slope field" params={{ method: 'euler', step: 0.25 }} />)
    expect(readout()).toContain('Euler')
    expect(readout()).toContain('0.25')
    setSlider('step h', 0.05)
    expect(readout()).toContain('0.05')
  })

  it('keeps the curve finite when it leaves the window', () => {
    render(<SlopeField title="Slope field" params={{ field: 'exponential', y0: 3 }} />)
    expect(readout()).toContain('boundary')
    expect(readout()).not.toContain('Infinity')
    expect(readout()).not.toContain('NaN')
  })
})

describe('root iteration figure', () => {
  it('iterates Newton from the default guess toward √2', () => {
    render(<RootIteration title="Roots" params={{}} />)
    expect(readout()).toContain('x² − 2')
    expect(readout()).toContain('Newton')
    expect(readout()).toContain('√2')
    expect(readout()).toContain('1.414')
    expect(readout()).not.toContain('NaN')
  })

  it('switches to bisection and shows the bracket', () => {
    render(<RootIteration title="Roots" params={{}} />)
    choose('Method', 'bisection')
    expect(readout()).toContain('Bisection')
    expect(readout()).toContain('[0, 2]')
    expect(readout()).toContain('f(a), f(b)')
    expect(document.querySelectorAll('rect').length).toBeGreaterThan(0)
  })

  it('changes the function and reports its exact root', () => {
    render(<RootIteration title="Roots" params={{ fn: 'cubic' }} />)
    expect(readout()).toContain('x³ − x − 1')
    expect(readout()).toContain('1.32471795')
  })

  it('stops on a stationary Newton tangent', () => {
    render(<RootIteration title="Roots" params={{ fn: 'quadratic', start: 0 }} />)
    expect(readout()).toContain('stationary')
    expect(readout()).toContain('tangent is horizontal')
    expect(readout()).not.toContain('NaN')
  })

  it('respects the iteration count', () => {
    render(<RootIteration title="Roots" params={{ method: 'bisection', iterations: 2 }} />)
    expect(readout()).toContain('2')
    setSlider('Iterations', 8)
    expect(readout()).toContain('8')
  })
})

describe('graph network figure', () => {
  it('runs breadth-first search from A', () => {
    render(<GraphNetwork title="Graph" params={{}} />)
    expect(readout()).toContain('Breadth-first')
    expect(readout()).toContain('A, B, D, C, E, F')
    expect(screen.getByLabelText('Show tree')).toBeTruthy()
    expect(readout()).not.toContain('NaN')
  })

  it('reports Dijkstra distances and the start node', () => {
    render(<GraphNetwork title="Graph" params={{}} />)
    choose('Algorithm', 'dijkstra')
    choose('Start node', 'C')
    expect(readout()).toContain('C: 0')
    expect(readout()).toContain('A: 2')
    expect(readout()).toContain('F: 1')
  })

  it('runs depth-first search from A', () => {
    render(<GraphNetwork title="Graph" params={{ algorithm: 'dfs' }} />)
    expect(readout()).toContain('Depth-first')
    expect(readout()).toContain('A, B, C, F, E, D')
  })

  it('runs Kruskal without a start node', () => {
    render(<GraphNetwork title="Graph" params={{ algorithm: 'kruskal' }} />)
    expect(readout()).toContain('Kruskal')
    expect(readout()).toContain('not used')
    expect((screen.getByLabelText('Start node') as HTMLSelectElement).disabled).toBe(true)
  })

  it('toggles the tree overlay and stays finite', () => {
    render(<GraphNetwork title="Graph" params={{}} />)
    const toggle = screen.getByLabelText('Show tree')
    fireEvent.click(toggle)
    expect(readout()).not.toContain('NaN')
    expect(readout()).not.toContain('Infinity')
  })
})
