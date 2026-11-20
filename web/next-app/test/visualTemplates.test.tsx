import { describe, it, expect } from 'vitest'
import { fireEvent, render, screen, within } from '@testing-library/react'
import MatrixTransform2D from '../components/visual/templates/MatrixTransform2D'
import TangentLine from '../components/visual/templates/TangentLine'
import RiemannSum from '../components/visual/templates/RiemannSum'
import ComplexMultiply from '../components/visual/templates/ComplexMultiply'
import NormalCurve from '../components/visual/templates/NormalCurve'
import Transform2D from '../components/visual/templates/Transform2D'
import VisualFigure from '../components/visual/VisualFigure'

function readout() {
  return screen.getByTestId('visual-readout').textContent ?? ''
}

describe('visual templates', () => {
  it('matrix transform updates the determinant when a coefficient changes', () => {
    render(
      <MatrixTransform2D
        title="Matrix"
        params={{ matrix: [1, 0.5, 0.5, 1], vector: [2, 1], showGrid: true }}
      />,
    )
    expect(readout()).toContain('0.75')
    fireEvent.change(screen.getByLabelText('a'), { target: { value: '2' } })
    expect(readout()).toContain('1.75')
  })

  it('tangent line moves its point and recomputes f(x)', () => {
    render(
      <TangentLine
        title="Tangent"
        params={{ coeffs: [0, 0, 1], x: 1, h: 0.8, window: [-2, 3] }}
      />,
    )
    expect(readout()).toContain('f(x)1')
    fireEvent.change(screen.getByLabelText('x'), { target: { value: '2' } })
    expect(readout()).toContain('f(x)4')
  })

  it('Riemann sum changes its partition count', () => {
    render(
      <RiemannSum
        title="Riemann"
        params={{ coeffs: [0, 0, 1], a: 0, b: 2, n: 8, rule: 'left' }}
      />,
    )
    fireEvent.change(screen.getByLabelText('n'), { target: { value: '0' } })
    const cells = within(screen.getByTestId('visual-readout')).getAllByRole('definition')
    expect(cells.map((c) => c.textContent)).toContain('1')
  })

  it('complex product changes when a real part changes', () => {
    render(
      <ComplexMultiply
        title="Product"
        params={{ z: [1.5, 0.8], w: [0.8, 1.2], showUnitCircle: true }}
      />,
    )
    fireEvent.change(screen.getByLabelText('z re'), { target: { value: '2' } })
    expect(readout()).toContain('2.15')
  })

  it('normal curve reports the probability at the moved point', () => {
    render(
      <NormalCurve
        title="Normal"
        params={{ mu: 0, sigma: 1, x: 1, bands: true }}
      />,
    )
    expect(readout()).toContain('0.84')
    fireEvent.change(screen.getByLabelText('x'), { target: { value: '0' } })
    expect(readout()).toContain('P(X ≤ x)0.5')
  })

  it('rotation leaves distances unchanged and moves the vertex', () => {
    render(
      <Transform2D
        title="Rotation"
        params={{
          kind: 'rotation',
          angle: 30,
          scale: 1.25,
          dx: 1,
          dy: 0.5,
          shape: [[0, 0], [2, 0], [1, 2]],
          showOrbit: true,
        }}
      />,
    )
    expect(readout()).toContain('1.73')
    fireEvent.change(screen.getByLabelText('angle'), { target: { value: '90' } })
    expect(readout()).toContain('(0, 2)')
    expect(readout()).toContain('side')
  })

  it('lets the learner build the polygon before transforming it', () => {
    render(
      <Transform2D
        title="Construction"
        params={{
          kind: 'rotation',
          angle: 30,
          scale: 1.25,
          dx: 1,
          dy: 0.5,
          shape: [[0, 0], [2, 0], [1, 2]],
          showOrbit: false,
        }}
      />,
    )
    expect(readout()).toContain('vertices3')

    // The first vertex is selected by default, and moving it moves its image:
    // a 30 degree turn sends (1, 0) to (0.87, 0.5).
    fireEvent.change(screen.getByLabelText('x'), { target: { value: '1' } })
    expect(readout()).toContain('(0.87, 0.5)')

    fireEvent.click(screen.getByRole('button', { name: 'Add vertex' }))
    expect(readout()).toContain('vertices4')
    fireEvent.click(screen.getByRole('button', { name: 'Remove vertex' }))
    expect(readout()).toContain('vertices3')
  })
})

describe('visual dispatcher', () => {
  it('explains an unknown template without breaking the lesson', () => {
    render(<VisualFigure visual={{ type: 'not.a.template' }} />)
    expect(screen.getByText(/This figure could not be shown/)).toBeInTheDocument()
  })
})

describe('matrix edge cases', () => {
  it('reports the intermediate and composed lengths separately', () => {
    // A scales x by 2, then B scales x by 2 again: |Av| = 2.83 and |BAv| = 4.47.
    render(
      <MatrixTransform2D
        title="Composed"
        params={{
          matrix: [2, 0, 0, 1],
          then: [2, 0, 0, 1],
          vector: [1, 2],
          showGrid: true,
          showBasis: true,
        }}
      />,
    )
    expect(readout()).toContain('|Av|2.83')
    expect(readout()).toContain('|BAv|4.47')
  })

  it('collapses a singular matrix and keeps a zero vector finite', () => {
    render(
      <MatrixTransform2D
        title="Singular"
        params={{ matrix: [1, 2, 2, 4], vector: [0, 0], showGrid: true, showBasis: true }}
      />,
    )
    expect(readout()).toContain('orientationcollapsed')
    expect(readout()).toContain('area0')
    expect(readout()).toContain('v(0, 0)')
    expect(readout()).toContain('|Av|0')
    expect(readout()).not.toContain('NaN')
    expect(readout()).not.toContain('—')
  })

  it('stays finite with large entries', () => {
    render(
      <MatrixTransform2D
        title="Large"
        params={{ matrix: [9, -1, 0, 9], vector: [3, 3], showGrid: true, showBasis: true }}
      />,
    )
    expect(readout()).toContain('det A81')
    expect(readout()).toContain('area81')
    expect(readout()).not.toContain('NaN')
  })
})

describe('matrix determinant cases', () => {
  it('shows preserved, reversed and collapsed orientation as the determinant changes sign', () => {
    const positive = render(
      <MatrixTransform2D
        title="Positive"
        params={{ matrix: [1.2, 0.6, 0.4, 1.4], vector: [2, 1], showGrid: true, showBasis: true }}
      />,
    )
    expect(readout()).toContain('orientationpreserved')
    expect(readout()).toContain('area1.44')
    positive.unmount()

    const negative = render(
      <MatrixTransform2D
        title="Negative"
        params={{ matrix: [0.2, 1, 1, 0.2], vector: [2, 1], showGrid: true, showBasis: true }}
      />,
    )
    expect(readout()).toContain('det A-0.96')
    expect(readout()).toContain('orientationreversed')
    negative.unmount()

    render(
      <MatrixTransform2D
        title="Zero"
        params={{ matrix: [1, 1, 1, 1], vector: [2, 1], showGrid: true, showBasis: true }}
      />,
    )
    expect(readout()).toContain('orientationcollapsed')
    expect(readout()).toContain('area0')
  })
})

describe('tangent line edge cases', () => {
  const renderTangent = (params: Record<string, unknown>) =>
    render(<TangentLine title="Tangent" params={params} />)

  it('reports a zero difference when the curve is already a line', () => {
    // y = 1 + 2x: the secant between any two points has the tangent's slope.
    renderTangent({ coeffs: [1, 2], x: 0.5, h: 0.8, window: [-2, 3] })
    expect(readout()).toContain('secant2')
    expect(readout()).toContain("f'(x)2")
    expect(readout()).toContain('difference0')
  })

  it('handles a constant function, whose derivative has no terms', () => {
    renderTangent({ coeffs: [3], x: 1, h: 0.8, window: [-2, 3] })
    expect(readout()).toContain('f(x)3')
    expect(readout()).toContain("f'(x)0")
    expect(readout()).toContain('secant0')
    expect(readout()).not.toContain('NaN')
    expect(readout()).not.toContain('—')
  })

  it('keeps the secant point inside the frame when x + h passes the window', () => {
    const { container } = renderTangent({ coeffs: [0, 0, 1], x: 1, h: 2, window: [0, 1] })
    const label = Array.from(container.querySelectorAll('text')).find((node) => node.textContent === 'x + h')
    expect(label).toBeDefined()
    // The label sits 8 units right of the point; anything beyond the viewBox means the
    // point itself was clipped, which is the defect this regression covers.
    expect(Number(label?.getAttribute('x'))).toBeLessThan(310)
    expect(readout()).not.toContain('NaN')
  })

  it('stays finite at the minimum h', () => {
    renderTangent({ coeffs: [0, 0, 1], x: 1, h: 0.05, window: [-2, 3] })
    expect(readout()).not.toContain('NaN')
    expect(readout()).toContain('h0.05')
  })
})

describe('Riemann sum edge cases', () => {
  const rectCount = () => screen.getByRole('img').querySelectorAll('rect').length

  it('draws one rectangle when n is 1, and its width is the whole interval', () => {
    render(<RiemannSum title="One" params={{ coeffs: [0, 0, 1], a: 0, b: 2, n: 1, rule: 'left' }} />)
    expect(rectCount()).toBe(1)
    expect(readout()).toContain('sum0')
    expect(readout()).toContain('Δx2')
  })

  it('treats a reversed interval as the same interval rather than a mirrored one', () => {
    const forward = render(<RiemannSum title="Forward" params={{ coeffs: [0, 0, 1], a: 0, b: 2, n: 4, rule: 'left' }} />)
    expect(readout()).toContain('sum1.75')
    expect(rectCount()).toBe(4)
    forward.unmount()
    // a > b used to draw a mirrored curve with zero-width rectangles, because the
    // rectangle width was computed as X(x1) - X(x0) and clamped at zero.
    render(<RiemannSum title="Reversed" params={{ coeffs: [0, 0, 1], a: 2, b: 0, n: 4, rule: 'left' }} />)
    expect(readout()).toContain('sum1.75')
    expect(rectCount()).toBe(4)
    const widths = Array.from(screen.getByRole('img').querySelectorAll('rect')).map((r) => Number(r.getAttribute('width')))
    expect(Math.min(...widths)).toBeGreaterThan(0)
  })

  it('subtracts area below the axis instead of dropping it', () => {
    render(<RiemannSum title="Negative" params={{ coeffs: [-1], a: 0, b: 2, n: 4, rule: 'left' }} />)
    expect(readout()).toContain('sum-2')
    expect(readout()).toContain('exact-2')
    expect(readout()).toContain('error0')
    expect(rectCount()).toBe(4)
  })

  it('accepts the 256 rectangles the schema allows', () => {
    render(<RiemannSum title="Many" params={{ coeffs: [0, 0, 1], a: 0, b: 2, n: 100, rule: 'mid' }} />)
    expect(readout()).toContain('n100')
    expect(readout()).not.toContain('NaN')
    expect(rectCount()).toBe(100)
  })
})

describe('complex multiplication edge cases', () => {
  it('reports the wrapped argument so the product and the sum agree', () => {
    render(
      <ComplexMultiply
        title="Product"
        params={{ z: [-0.5, 0.866], w: [-0.5, 0.866], showUnitCircle: true }}
      />,
    )
    // 120 degrees plus 120 degrees is 240, which the product reports as -120.
    expect(readout()).toContain('arg zw-2.09')
    expect(readout()).toContain('arg sum-2.09')
  })

  it('says the argument is undefined at the origin instead of calling it zero', () => {
    render(
      <ComplexMultiply title="Product" params={{ z: [0, 0], w: [1, 0], showUnitCircle: true }} />,
    )
    expect(readout()).toContain('arg zundefined')
    expect(readout()).toContain('|zw|0')
    expect(readout()).toContain('arg zwundefined')
  })
})

describe('normal curve edge cases', () => {
  it('keeps a fixed window so sigma visibly widens the curve', () => {
    render(<NormalCurve title="Normal" params={{ mu: 0, sigma: 1, x: 1, bands: true }} />)
    // peak density is 1/sqrt(2 pi) at sigma 1 and a quarter of that at sigma 4
    expect(readout()).toContain('peak0.399')
    fireEvent.change(screen.getByLabelText('σ'), { target: { value: '4' } })
    expect(readout()).toContain('peak0.1')
    expect(readout()).toContain('σ4')
    // the band probability is the same curve area at every sigma
    expect(readout()).toContain('μ ± σ0.6827')
  })

  it('reports a value outside the window without pretending the marker is there', () => {
    render(<NormalCurve title="Normal" params={{ mu: 0, sigma: 1, x: 10, bands: true }} />)
    expect(readout()).toContain('z10')
    expect(readout()).toContain('P(X ≤ x)1')
  })

  it('survives the smallest sigma the schema allows', () => {
    render(<NormalCurve title="Normal" params={{ mu: 0, sigma: 0.1, x: 0.1, bands: true }} />)
    expect(readout()).toContain('σ0.1')
    expect(readout()).toContain('peak3.989')
    expect(readout()).not.toContain('NaN')
  })
})

describe('transform edge cases', () => {
  const base = { angle: 0, scale: 1, dx: 1, dy: 0.5, shape: [[0, 0], [2, 0], [1, 2]], showOrbit: true }

  it('keeps the matrix, offset and area scale honest for translation and reflection', () => {
    const first = render(<Transform2D title="T" params={{ ...base, kind: 'translation' }} />)
    expect(readout()).toContain('matrix[[1, 0], [0, 1]]')
    expect(readout()).toContain('offset(1, 0.5)')
    expect(readout()).toContain('area scale1')
    first.unmount()

    render(<Transform2D title="T" params={{ ...base, kind: 'reflection' }} />)
    expect(readout()).toContain('matrix[[-1, 0], [0, 1]]')
    // A reflection is an isometry: no side changes, no area change.
    expect(readout()).toContain('max side change0')
    expect(readout()).toContain('area scale1')
  })

  it('handles a degenerate collinear polygon without NaN', () => {
    render(<Transform2D title="T" params={{ ...base, kind: 'dilation', scale: 2, shape: [[0, 0], [1, 0], [2, 0]] }} />)
    expect(readout()).toContain('area scale4')
    expect(readout()).toContain('max side change2')
    expect(readout()).not.toContain('NaN')
  })

  it('stops adding at eight vertices and removing at three', () => {
    const eight = [[0, 0], [1, 0], [2, 1], [3, 0], [4, 1], [3, 2], [2, 1], [1, 2]]
    const first = render(<Transform2D title="T" params={{ ...base, kind: 'rotation', shape: eight }} />)
    expect(readout()).toContain('vertices8')
    expect(screen.getByRole('button', { name: 'Add vertex' })).toBeDisabled()
    first.unmount()

    render(<Transform2D title="T" params={{ ...base, kind: 'rotation', shape: [[0, 0], [2, 0], [1, 2]] }} />)
    expect(readout()).toContain('vertices3')
    expect(screen.getByRole('button', { name: 'Remove vertex' })).toBeDisabled()
  })

  it('does not rescale the frame when a vertex is edited', () => {
    render(
      <Transform2D
        title="T"
        params={{ kind: 'rotation', angle: 30, scale: 1, dx: 1, dy: 0.5, shape: [[0, 0], [2, 0], [1, 2]], showOrbit: true }}
      />,
    )
    const svg = screen.getByRole('img') as unknown as SVGSVGElement
    const xs = () => Array.from(svg.querySelectorAll('circle')).map((circle) => circle.getAttribute('cx'))
    const before = xs()
    fireEvent.change(screen.getByLabelText('x'), { target: { value: '5' } })
    const after = xs()
    // Circle order is the three orbit rings, then the learner's vertices, then
    // their images. Moving vertex 0 must not move vertex 1: the frame comes from
    // the authored polygon, not from the live shape, so the view never rescales
    // underneath a drag.
    expect(after[4]).toBe(before[4])
  })
})

describe('matrix inverse readout', () => {
  const readout = () => screen.getByTestId('visual-readout').textContent ?? ''

  it('reports the inverse when one exists and says when it does not', () => {
    const first = render(
      <MatrixTransform2D title="Inverse" params={{ matrix: [2, 1, 1, 1], vector: [2, 1], showGrid: true, showBasis: true }} />,
    )
    // Compact on purpose: the spaced form wrapped to three lines in a 320px cell.
    expect(readout()).toContain('[[1,-1],[-1,2]]')
    first.unmount()
    render(
      <MatrixTransform2D title="Singular" params={{ matrix: [1, 2, 2, 4], vector: [1, 1], showGrid: true, showBasis: true }} />,
    )
    expect(readout()).toContain('no inverse (singular)')
    expect(screen.getByText(/no inverse/)).toBeInTheDocument()
  })
})

describe('matrix column combination', () => {
  it('states Av as the sum of the scaled columns', () => {
    render(
      <MatrixTransform2D
        title="Columns"
        params={{ matrix: [1.2, 0.6, 0.4, 1.4], vector: [2, 1], showGrid: true, showBasis: true }}
      />,
    )
    const text = screen.getByTestId('visual-readout').textContent ?? ''
    // The arrows draw the two scaled columns; the readout has to say that their sum
    // is Av, because a reader should not have to infer the equation from the picture.
    expect(text).toContain('Av = 2·(1.2, 0.4) + 1·(0.6, 1.4)')
    expect(text).toContain('(3, 2.2)')
  })
})
