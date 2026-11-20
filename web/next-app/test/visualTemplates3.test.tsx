import { describe, it, expect } from 'vitest'
import { fireEvent, render, screen } from '@testing-library/react'
import BasisChange from '../components/visual/templates/BasisChange'
import SpreadMeasure from '../components/visual/templates/SpreadMeasure'
import ConicSection from '../components/visual/templates/ConicSection'
import CoordinateWorkbench from '../components/visual/templates/CoordinateWorkbench'
import ComplexMapping from '../components/visual/templates/ComplexMapping'
import ContourMap from '../components/visual/templates/ContourMap'
import DistributionFamilies from '../components/visual/templates/DistributionFamilies'

function readout() {
  return screen.getByTestId('visual-readout').textContent ?? ''
}

describe('basis-change template', () => {
  it('reports the coordinates of the point in the drawn basis', () => {
    render(<BasisChange title="Basis" params={{ u: [1.5, 0.5], v: [0.5, 1.5], point: [2, 1] }} />)
    // P = [[1.5, 0.5], [0.5, 1.5]], det = 2, and (2, 1) is 1.25 u + 0.25 v.
    expect(readout()).toContain('det P2')
    expect(readout()).toContain('(α, β) = (1.25, 0.25)')
    expect(readout()).toContain('point(2, 1)')
  })

  it('recomputes the coordinates when a basis vector moves', () => {
    render(<BasisChange title="Basis" params={{ u: [1.5, 0.5], v: [0.5, 1.5], point: [2, 1] }} />)
    fireEvent.change(screen.getByLabelText('ux'), { target: { value: '3' } })
    expect(readout()).toContain('det P4')
    expect(readout()).toContain('(α, β) = (0.59, 0.47)')
  })
})

describe('spread template', () => {
  it('reports the mean, the sum of squares and both denominators', () => {
    render(<SpreadMeasure title="Spread" params={{ values: [1, 2, 2, 4, 6], squared: true, sample: true }} />)
    expect(readout()).toContain('mean3')
    expect(readout()).toContain('sum of squares16')
    expect(readout()).toContain('variance4')
    expect(readout()).toContain('std dev2')
    // [1,2,2,4,6]: gaps -2,-1,-1,1,3, squares 4,1,1,1,9 = 16; 16/4 = 4.
  })

  it('switches to the population denominator', () => {
    render(<SpreadMeasure title="Spread" params={{ values: [1, 2, 2, 4, 6], squared: true, sample: false }} />)
    expect(readout()).toContain('variance3.2')
    expect(readout()).toContain('divides byn = 5')
  })

  it('moves a point through its number field', () => {
    render(<SpreadMeasure title="Spread" params={{ values: [1, 2, 2, 4, 6], squared: true, sample: true }} />)
    fireEvent.change(screen.getByLabelText('x'), { target: { value: '10' } })
    // Moving P1 from 1 to 10 makes the spread much larger.
    expect(readout()).not.toContain('std dev2')
  })
})

describe('conic section', () => {
  it('calls a circle a circle and a zero discriminant two parallel lines', () => {
    render(<ConicSection title="Conic" params={{ A: 1, B: 0, C: 1 }} />)
    expect(readout()).toContain('typecircle')
    expect(readout()).toContain('discriminant-4')
    fireEvent.change(screen.getByLabelText('C'), { target: { value: '0' } })
    expect(readout()).toContain('typetwo parallel lines')
    expect(readout()).toContain('discriminant0')
  })
})

describe('contour map template', () => {
  it('reports the value and the gradient at the probe', () => {
    render(<ContourMap title="Contours" params={{ field: 'bowl', levels: 7, probe: [1, 1] }} />)
    // x^2 + y^2 at (1, 1) is 2, its gradient is (2, 2), and |grad| is sqrt 8.
    expect(readout()).toContain('2.828')
    fireEvent.change(screen.getByLabelText('x'), { target: { value: '0' } })
    // At (0, 1) the value is 1 and the gradient is vertical with length 2.
    expect(readout()).toContain('gradient angle90°')
  })
})

describe('complex mapping figure', () => {
  it('reports the same angle out as in, and changes when the opening changes', () => {
    render(
      <ComplexMapping
        title="Map"
        params={{ map: 'square', z: [0.8, 0.4], turn: 60 }}
      />,
    )
    expect(readout()).toContain('angle in60°')
    expect(readout()).toContain('angle out60°')
    fireEvent.change(screen.getByLabelText('opening'), { target: { value: '90' } })
    expect(readout()).toContain('angle in90°')
    expect(readout()).toContain('angle out90°')
  })
})

describe('coordinate workbench', () => {
  const params = { mode: 'measure', a: [-3, 1], b: [2, 4], c: [1, -2], slope: 1 }

  it('shows the distance and midpoint working, and follows a typed coordinate', () => {
    render(<CoordinateWorkbench title="Coordinates" params={params} />)
    expect(readout()).toContain('distance5.831')
    expect(readout()).toContain('midpoint(-0.5, 2.5)')
    expect(screen.getByTestId('visual-working').textContent).toContain('√34')
    fireEvent.change(screen.getByLabelText('x'), { target: { value: '2' } })
    expect(readout()).toContain('distance3')
    expect(readout()).toContain('Δx0')
  })

  it('names a parallel pair instead of inventing an intersection', () => {
    render(<CoordinateWorkbench title="Coordinates" params={{ ...params, mode: 'intersect' }} />)
    expect(readout()).toContain('statusmeets')
    fireEvent.change(screen.getByLabelText('slope'), { target: { value: '0.6' } })
    expect(readout()).toContain('statusparallel')
    expect(readout()).toContain('intersectionnone')
  })
})

describe('distribution families template', () => {
  it('puts the same interval at different probabilities for each family', () => {
    render(<DistributionFamilies title="Families" params={{ emphasis: 'uniform', x: 1 }} />)
    expect(readout()).toContain('P uniform0.7887')
    expect(readout()).toContain('P exp0.6321')
    expect(readout()).toContain('mean0')
    fireEvent.change(screen.getByLabelText('x'), { target: { value: '0' } })
    expect(readout()).toContain('P uniform0.5')
    expect(readout()).toContain('P normal0.5')
  })
})

describe('basis change edge cases', () => {
  it('reads coordinates back out of the identity basis', () => {
    render(<BasisChange title="Basis" params={{ u: [1, 0], v: [0, 1], point: [2, 1] }} />)
    expect(screen.getByTestId('visual-readout').textContent).toContain('(α, β) = (2, 1)')
  })

  it('says a dependent pair gives no coordinates at all', () => {
    render(<BasisChange title="Basis" params={{ u: [1, 0], v: [2, 0], point: [1, 1] }} />)
    const readout = screen.getByTestId('visual-readout').textContent ?? ''
    expect(readout).toContain('coordsundefined')
    expect(readout).toContain('αu + βvno basis')
    expect(screen.getByRole('img').textContent).toContain('u and v are dependent')
  })

  it('does not print a thin basis as singular, and warns that its coordinates are huge', () => {
    render(<BasisChange title="Basis" params={{ u: [1, 0], v: [1.0001, 0.0001], point: [2, 1] }} />)
    const readout = screen.getByTestId('visual-readout').textContent ?? ''
    // det is 1e-4: it must not round to a bare "0" beside coordinates that exist.
    expect(readout).not.toContain('det P0')
    expect(readout).toContain('e-4')
    // The coordinates are genuine, and enormous, so the figure says so rather than
    // hiding them or drawing a decomposition that leaves the plot.
    expect(readout).toContain('(-9999, 10000)')
    expect(screen.getByRole('img').textContent).toContain('nearly dependent')
  })
})

describe('spread measure edge cases', () => {
  it('draws no bars and says so when every value equals the mean', () => {
    render(<SpreadMeasure title="Spread" params={{ values: [5, 5, 5, 5], squared: true, sample: true }} />)
    expect(readout()).toContain('std dev0')
    expect(readout()).toContain('sum of squares0')
    const heights = Array.from(screen.getByRole('img').querySelectorAll('rect')).map((r) => r.getAttribute('height'))
    expect(heights.every((h) => h === '0')).toBe(true)
    expect(screen.getByRole('img').textContent).toContain('no spread')
  })

  it('reports the sample and population variance for two points', () => {
    render(<SpreadMeasure title="Spread" params={{ values: [2, 6], squared: true, sample: true }} />)
    expect(readout()).toContain('mean4')
    expect(readout()).toContain('sum of squares8')
    expect(readout()).toContain('variance8')
    expect(readout()).toContain('divides byn − 1 = 1')
    fireEvent.click(screen.getByLabelText('Sample (n − 1)'))
    expect(readout()).toContain('variance4')
    expect(readout()).toContain('divides byn = 2')
  })
})

describe('contour edge cases', () => {
  it('reports no gradient at a critical point and hides the level direction', () => {
    const { container } = render(
      <ContourMap title="Contour" params={{ field: 'bowl', levels: 7, probe: [0, 0] }} />,
    )
    expect(readout()).toContain('|∇f|0')
    expect(readout()).toContain('gradient angleundefined')
    // A zero gradient has no perpendicular, so the dashed level direction must not
    // be drawn: every direction is tangent at a critical point.
    expect(container.querySelector('line[stroke-dasharray="4 3"]')).toBeNull()
    expect(container.textContent).toContain('critical point')
  })

  it('reports the gradient of x times y and its angle off the axes', () => {
    render(<ContourMap title="Contour" params={{ field: 'product', levels: 5, probe: [1, 1] }} />)
    expect(readout()).toContain('|∇f|1.414')
    expect(readout()).toContain('gradient angle45°')
  })

  it('accepts the lowest and highest level counts', () => {
    render(<ContourMap title="Contour" params={{ field: 'saddle', levels: 3, probe: [1, 0.5] }} />)
    expect(readout()).toContain('levels3')
    fireEvent.change(screen.getByLabelText('levels'), { target: { value: '12' } })
    expect(readout()).toContain('levels12')
  })
})

describe('complex mapping edge cases', () => {
  it('reports a doubled angle at the critical point of z squared', () => {
    render(<ComplexMapping title="Map" params={{ map: 'square', z: [0, 0], turn: 60 }} />)
    expect(readout()).toContain('angle out120°')
    expect(readout()).toContain('conformalno')
    expect(screen.getByRole('img').textContent).toContain('the angle doubles')
  })

  it('reports no value at a pole and never puts NaN in the drawing', () => {
    render(<ComplexMapping title="Map" params={{ map: 'inverse', z: [0, 0], turn: 60 }} />)
    expect(readout()).toContain('f(z)undefined')
    expect(readout()).toContain('angle outundefined')
    expect(readout()).toContain('conformalno')
    expect(screen.getByRole('img').textContent).toContain('no value here')
    expect(screen.getByRole('img').innerHTML).not.toContain('NaN')
  })
})

describe('coordinate workbench edge cases', () => {
  const base = { mode: 'intersect' as const, a: [-2, 0] as [number, number], b: [2, 2] as [number, number], c: [0, 1] as [number, number], slope: 0.5 }

  it('refuses to invent a line through a single point', () => {
    render(<CoordinateWorkbench title="Coordinate" params={{ ...base, b: [-2, 0] }} />)
    expect(readout()).toContain('line 1undefined')
    expect(readout()).toContain('intersectionundefined')
    expect(readout()).toContain('statusundefined')
    expect(screen.getByText('A and B coincide: no line')).toBeInTheDocument()
  })

  it('meets a vertical line at the shared x', () => {
    render(<CoordinateWorkbench title="Coordinate" params={{ ...base, a: [2, -2], b: [2, 2], c: [0, -2], slope: 2 }} />)
    expect(readout()).toContain('line 1x = 2')
    expect(readout()).toContain('intersection(2, 2)')
    expect(readout()).toContain('statusmeets')
  })

  it('reports coincident lines as every shared point rather than a coordinate', () => {
    render(<CoordinateWorkbench title="Coordinate" params={base} />)
    expect(readout()).toContain('statuscoincident')
    expect(readout()).toContain('intersectionall shared')
    expect(screen.getByText('coincident: every point shared')).toBeInTheDocument()
  })

  it('reports parallel lines with no shared point', () => {
    render(<CoordinateWorkbench title="Coordinate" params={{ ...base, c: [0, 3] }} />)
    expect(readout()).toContain('statusparallel')
    expect(readout()).toContain('intersectionnone')
  })

  it('says when the meeting point lies outside the frame', () => {
    render(<CoordinateWorkbench title="Coordinate" params={{ ...base, c: [1, -2], slope: 1 }} />)
    expect(readout()).toContain('statusmeets')
    expect(screen.getByText('meets outside this window')).toBeInTheDocument()
  })

  it('gives a coincident pair distance zero and its own midpoint', () => {
    render(
      <CoordinateWorkbench
        title="Coordinate"
        params={{ mode: 'measure', a: [1, 1], b: [1, 1], c: [0, 0], slope: 1 }}
      />,
    )
    expect(readout()).toContain('distance0')
    expect(readout()).toContain('midpoint(1, 1)')
    expect(screen.getByTestId('visual-working').textContent).toContain('√0')
  })
})

describe('distribution family edge cases', () => {
  it('reports nothing below zero for the exponential, and says which tail it is', () => {
    render(<DistributionFamilies title="Families" params={{ emphasis: 'exponential', x: -4.5 }} />)
    expect(readout()).toContain('f(x)0')
    expect(readout()).toContain('P(X ≤ x)0')
    expect(readout()).toContain('P(X > x)1')
    expect(readout()).toContain('mean1')
    expect(readout()).toContain('sd1')
  })

  it('reaches certainty at the right edge of the window', () => {
    render(<DistributionFamilies title="Families" params={{ emphasis: 'normal', x: 5.5 }} />)
    expect(readout()).toContain('P normal1')
    expect(readout()).toContain('P uniform1')
    expect(readout()).toContain('P exp0.9959')
    expect(readout()).toContain('P(X > x)0')
  })

  it('keeps the three probabilities distinct at the same x and updates with the slider', () => {
    render(<DistributionFamilies title="Families" params={{ emphasis: 'normal', x: 1 }} />)
    expect(readout()).toContain('P normal0.8413')
    expect(readout()).toContain('P uniform0.7887')
    expect(readout()).toContain('P exp0.6321')
    fireEvent.change(screen.getByLabelText('x'), { target: { value: '3' } })
    expect(readout()).toContain('P exp0.9502')
    expect(readout()).toContain('P(X > x)0.0013')
  })
})
