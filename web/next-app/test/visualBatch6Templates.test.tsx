import { fireEvent, render, screen } from '@testing-library/react'
import { describe, expect, it } from 'vitest'
import VectorFieldLab from '../components/visual/templates/VectorFieldLab'
import SurfaceLab from '../components/visual/templates/SurfaceLab'
import RegionIntegral from '../components/visual/templates/RegionIntegral'
import SeriesLab from '../components/visual/templates/SeriesLab'
import AccumulationLab from '../components/visual/templates/AccumulationLab'
import SolidRevolution from '../components/visual/templates/SolidRevolution'
import ConditionalArea from '../components/visual/templates/ConditionalArea'
import GraphNetwork from '../components/visual/templates/GraphNetwork'

const readout = () => screen.getByTestId('visual-readout').textContent ?? ''
const choose = (label: string, value: string) =>
  fireEvent.change(screen.getByLabelText(label), { target: { value } })
const setSlider = (label: string, value: number) =>
  fireEvent.change(screen.getByLabelText(label), { target: { value: String(value) } })

describe('VectorFieldLab', () => {
  const params = { preset: 'vortex' as const, radius: 1.2, showStreamlines: true }

  it('reports the circulation and checks Green', () => {
    render(<VectorFieldLab title="Field" params={params} />)
    const text = readout()
    expect(text).toContain('divergence0')
    expect(text).toContain('curl2')
    expect(text).toContain('circulation9.048')
    expect(text).toContain('Green checkholds')
    expect(screen.getByRole('img').getAttribute('aria-label')).toContain('Vector field')
  })

  it('changes to the source field and its flux', () => {
    render(<VectorFieldLab title="Field" params={params} />)
    choose('Field', 'source')
    const text = readout()
    expect(text).toContain('divergence2')
    expect(text).toContain('flux9.048')
    expect(text).toContain('circulation0')
  })

  it('grows the loop and the integral with it', () => {
    render(<VectorFieldLab title="Field" params={params} />)
    setSlider('Loop radius', 2)
    expect(readout()).toContain('circulation25.133')
  })
})

describe('SurfaceLab', () => {
  const params = { field: 'bowl' as const, x0: 1, y0: 0.6 }

  it('draws the tangent plane at the probe', () => {
    render(<SurfaceLab title="Surface" params={params} />)
    const text = readout()
    expect(text).toContain('z at probe1.36')
    expect(text).toContain('gradient(2, 1.2)')
    expect(text).toContain('classificationminimum')
    expect(screen.getByRole('img').getAttribute('aria-label')).toContain('tangent plane')
  })

  it('finds the saddle at the origin', () => {
    render(<SurfaceLab title="Surface" params={params} />)
    choose('Surface', 'saddle')
    setSlider('Probe x', 0)
    setSlider('Probe y', 0)
    const text = readout()
    expect(text).toContain('D = f_xx f_yy − f_xy²-4')
    expect(text).toContain('classificationsaddle')
  })
})

describe('RegionIntegral', () => {
  const params = { region: 'between' as const, density: 'one' as const, n: 10 }

  it('reports the estimate and the exact area', () => {
    render(<RegionIntegral title="Region" params={params} />)
    const text = readout()
    expect(text).toContain('exact0.1667')
    expect(text).toContain('cells inside')
    expect(text).toContain('area elementdx dy')
    expect(screen.getByRole('img').getAttribute('aria-label')).toContain('partitioned')
  })

  it('refines the grid toward the exact value', () => {
    render(<RegionIntegral title="Region" params={params} />)
    setSlider('Cells per side', 40)
    expect(readout()).toContain('cells inside')
    expect(readout()).toContain('exact0.1667')
  })

  it('switches to the polar region and names its area element', () => {
    render(<RegionIntegral title="Region" params={params} />)
    choose('Region', 'polar')
    const text = readout()
    expect(text).toContain('exact0.7854')
    expect(text).toContain('area elementr dr dθ')
  })
})

describe('SeriesLab', () => {
  const params = { kind: 'geometric' as const, terms: 12 }

  it('shows the partial sum approaching one', () => {
    render(<SeriesLab title="Series" params={params} />)
    const text = readout()
    expect(text).toContain('statusconverges')
    expect(text).toContain('limit1')
    expect(text).toContain('partial sum0.99976')
    expect(text).toContain('ratio |aₙ₊₁/aₙ|0.5')
  })

  it('shows the harmonic series diverging', () => {
    render(<SeriesLab title="Series" params={params} />)
    choose('Series', 'harmonic')
    const text = readout()
    expect(text).toContain('statusdiverges')
    expect(text).toContain('limitno finite limit')
  })

  it('adds terms on demand', () => {
    render(<SeriesLab title="Series" params={params} />)
    setSlider('Terms', 2)
    expect(readout()).toContain('partial sum0.75')
  })
})

describe('AccumulationLab', () => {
  const params = { fn: 'square' as const, a: 0.5, x: 2 }

  it('differentiates the area back to f(x)', () => {
    render(<AccumulationLab title="Area" params={params} />)
    const text = readout()
    expect(text).toContain('area A(x)2.625')
    expect(text).toContain('f(x)4')
    expect(text).toContain('theorem holds')
    expect(screen.getByRole('img').getAttribute('aria-label')).toContain('Area under')
  })

  it('moves the interval', () => {
    render(<AccumulationLab title="Area" params={params} />)
    setSlider('Right end x', 3)
    expect(readout()).toContain('area A(x)')
    setSlider('Left end a', 0)
    expect(readout()).toContain('area A(x)9')
  })
})

describe('SolidRevolution', () => {
  const params = { mode: 'washer' as const, slice: 1.2 }

  it('reports the washer slice and the total volume', () => {
    render(<SolidRevolution title="Solid" params={params} />)
    const text = readout()
    expect(text).toContain('radius there1.24')
    expect(text).toContain('slice area4.3279')
    expect(text).toContain('total volume8.1263')
    expect(screen.getByRole('img').getAttribute('aria-label')).toContain('Washer body')
  })

  it('switches to shells and keeps a different volume', () => {
    render(<SolidRevolution title="Solid" params={params} />)
    choose('Method', 'shell')
    const text = readout()
    expect(text).toContain('methodCylindrical shell')
    expect(text).toContain('formulaV = 2π ∫ x f(x) dx')
  })

  it('moves the slice', () => {
    render(<SolidRevolution title="Solid" params={params} />)
    setSlider('Slice position', 2)
    expect(readout()).toContain('radius there1.4')
  })
})

describe('ConditionalArea', () => {
  const params = { pB: 0.4, pAgivenB: 0.8, pAgivenBc: 0.3 }

  it('draws the multiplication rule as an area', () => {
    render(<ConditionalArea title="Conditional" params={params} />)
    const text = readout()
    expect(text).toContain('P(A ∩ B)0.32')
    expect(text).toContain('P(A)0.5')
    expect(text).toContain('independencefails')
    expect(screen.getByRole('img').getAttribute('aria-label')).toContain('Unit square split')
  })

  it('detects independence when the bands level', () => {
    render(<ConditionalArea title="Conditional" params={params} />)
    setSlider('P(A | B)', 0.3)
    expect(readout()).toContain('independenceholds')
  })
})

describe('GraphNetwork bipartition and colouring', () => {
  const params = { algorithm: 'bfs' as const, start: 'A', showTree: true }

  it('splits the even cycle into two sides', () => {
    render(<GraphNetwork title="Graph" params={params} />)
    choose('Algorithm', 'bipartite')
    const text = readout()
    expect(text).toContain('bipartitionA, C, E | B, D, F')
    expect(text).toContain('statusok')
  })

  it('colours the odd cycle with three colours', () => {
    render(<GraphNetwork title="Graph" params={params} />)
    choose('Algorithm', 'colouring')
    const text = readout()
    expect(text).toContain('colour count3')
    expect(text).toContain('B: 3')
    expect(text).toContain('C: 2')
  })
})
