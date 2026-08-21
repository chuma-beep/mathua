import { describe, it, expect, beforeAll, vi } from 'vitest'
import { render, screen, fireEvent } from '@testing-library/react'
import ConceptGraphFlow from '../components/ConceptGraphFlow'
import { deriveStatuses } from '../lib/graphStatus'

class ResizeObserverMock {
  observe() {}
  unobserve() {}
  disconnect() {}
}

beforeAll(() => {
  Object.defineProperty(globalThis, 'ResizeObserver', {
    writable: true,
    value: ResizeObserverMock,
  })
  Object.defineProperty(globalThis, 'DOMMatrixReadOnly', {
    writable: true,
    value: class {
      m22 = 1
      constructor(_transform?: string) {}
    },
  })
})

const concepts = [
  { id: 'count', label: 'Counting', domain: 'counting', prerequisites: [] },
  { id: 'add', label: 'Addition', domain: 'arithmetic', prerequisites: ['count'] },
  { id: 'mul', label: 'Multiplication', domain: 'arithmetic', prerequisites: ['add'] },
]

describe('ConceptGraphFlow', () => {
  it('renders without crashing and shows the search box', () => {
    render(<ConceptGraphFlow concepts={concepts} />)
    expect(screen.getByPlaceholderText('Search concepts…')).toBeInTheDocument()
  })

  it('renders the empty state when there are no concepts', () => {
    render(<ConceptGraphFlow concepts={[]} />)
    expect(screen.getByText('No concepts to display')).toBeInTheDocument()
  })

  it('search finds a concept and selecting it calls onSelectionChange', async () => {
    const onSelectionChange = vi.fn()
    render(<ConceptGraphFlow concepts={concepts} onSelectionChange={onSelectionChange} />)
    const input = screen.getByPlaceholderText('Search concepts…')
    fireEvent.change(input, { target: { value: 'multipl' } })
    const hits = await screen.findAllByText('Multiplication')
    fireEvent.click(hits[0])
    expect(onSelectionChange).toHaveBeenCalledWith('mul')
  })

  it('derives locked statuses that flow into node colors', () => {
    const statuses = deriveStatuses(concepts)
    expect(statuses['count']).toBe('unseen')
    expect(statuses['add']).toBe('locked')
    expect(statuses['mul']).toBe('locked')
  })
})
