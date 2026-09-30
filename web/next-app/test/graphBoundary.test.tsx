import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import { GraphErrorBoundary } from '../app/home/sections'

function Boom(): React.ReactNode {
  throw new Error('chunk failed')
}

describe('GraphErrorBoundary', () => {
  it('renders children when the graph loads', () => {
    render(
      <GraphErrorBoundary fallback={<div>poster fallback</div>}>
        <div>webgl graph</div>
      </GraphErrorBoundary>,
    )
    expect(screen.getByText('webgl graph')).toBeTruthy()
  })

  it('falls back to the poster instead of crashing the page', () => {
    vi.spyOn(console, 'error').mockImplementation(() => {})
    render(
      <GraphErrorBoundary fallback={<div>poster fallback</div>}>
        <Boom />
      </GraphErrorBoundary>,
    )
    expect(screen.getByText('poster fallback')).toBeTruthy()
    expect(screen.queryByText('Something broke on this page.')).toBeNull()
    vi.restoreAllMocks()
  })
})
