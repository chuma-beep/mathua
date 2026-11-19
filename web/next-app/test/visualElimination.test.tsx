import { describe, expect, it } from 'vitest'
import { fireEvent, render, screen } from '@testing-library/react'
import { eliminate, outcomeText, type Row } from '../components/visual/math/rowReduce'
import GaussianElimination from '../components/visual/templates/GaussianElimination'

const close = (a: number, b: number, tolerance = 1e-9) => expect(Math.abs(a - b)).toBeLessThan(tolerance)

function satisfies(rows: Row[], x: number, y: number): boolean {
  return rows.every((row) => Math.abs(row[0] * x + row[1] * y - row[2]) < 1e-9)
}

describe('row reduction', () => {
  it('solves a unique system and records every operation', () => {
    const rows: Row[] = [
      [2, 1, 5],
      [1, -1, 1],
    ]
    const outcome = eliminate(rows)
    expect(outcome.kind).toBe('unique')
    if (outcome.kind !== 'unique') return
    expect(outcome.values[0]).toBe(2)
    expect(outcome.values[1]).toBe(1)
    expect(satisfies(rows, 2, 1)).toBe(true)
    expect(outcome.steps.length).toBeGreaterThan(0)
    for (const step of outcome.steps) expect(step.description.length).toBeGreaterThan(10)
    expect(outcomeText(outcome)).toContain('x = 2')
  })

  it('keeps the solution set at every step, which is the claim the figure makes', () => {
    const rows: Row[] = [
      [3, -2, 4],
      [1, 5, 9],
    ]
    const outcome = eliminate(rows)
    expect(outcome.kind).toBe('unique')
    if (outcome.kind !== 'unique') return
    const [x, y] = outcome.values
    expect(satisfies(rows, x, y)).toBe(true)
    // Each intermediate matrix must still accept that same pair.
    for (const step of outcome.steps) expect(satisfies(step.matrix, x, y)).toBe(true)
  })

  it('reports infinite solutions and names the free variable', () => {
    const outcome = eliminate([
      [1, 1, 2],
      [2, 2, 4],
    ])
    expect(outcome.kind).toBe('infinite')
    if (outcome.kind !== 'infinite') return
    expect(outcome.free).toEqual([1])
    expect(outcomeText(outcome)).toBe('infinitely many solutions: y is free')
  })

  it('reports an inconsistent system as having no solution', () => {
    const outcome = eliminate([
      [1, 1, 2],
      [1, 1, 3],
    ])
    expect(outcome.kind).toBe('none')
    if (outcome.kind !== 'none') return
    expect(outcome.reason).toContain('0 =')
    expect(outcomeText(outcome)).toBe('no solution')
  })

  it('swaps to a pivot when the first column is empty', () => {
    const outcome = eliminate([
      [0, 1, 2],
      [1, 0, 3],
    ])
    expect(outcome.kind).toBe('unique')
    if (outcome.kind !== 'unique') return
    close(outcome.values[0], 3)
    close(outcome.values[1], 2)
  })
})

describe('elimination template', () => {
  function readout() {
    return screen.getByTestId('visual-readout').textContent ?? ''
  }

  it('shows the system, steps through it, and reads the answer', () => {
    render(<GaussianElimination title="Elimination" params={{ rows: [[2, 1, 5], [1, -1, 1]] }} />)
    expect(readout()).toContain('x = 2, y = 1')
    expect(readout()).toContain('of')
    const before = screen.getByTestId('elimination-step').textContent
    fireEvent.click(screen.getByRole('button', { name: 'Step' }))
    expect(screen.getByTestId('elimination-step').textContent).not.toBe(before)
    expect(readout()).toContain('1 of')
  })

  it('reports inconsistency when the constants contradict each other', () => {
    render(<GaussianElimination title="Elimination" params={{ rows: [[1, 1, 2], [1, 1, 3]] }} />)
    expect(readout()).toContain('no solution')
  })

  it('reports a free variable when the rows are multiples', () => {
    render(<GaussianElimination title="Elimination" params={{ rows: [[1, 1, 2], [2, 2, 4]] }} />)
    expect(readout()).toContain('infinitely many solutions')
    expect(readout()).toContain('y is free')
  })
})
