import { describe, it, expect, beforeAll, beforeEach, vi } from 'vitest'
import { render, screen, fireEvent, within } from '@testing-library/react'
import { deriveStatuses } from '../lib/graphStatus'

const reactFlowMocks = vi.hoisted(() => ({
  zoomIn: vi.fn(),
  zoomOut: vi.fn(),
  fitView: vi.fn(),
  getViewport: vi.fn(() => ({ x: 0, y: 0, zoom: 1 })),
  setViewport: vi.fn(),
  setCenter: vi.fn(),
}))

vi.mock('@xyflow/react', async importOriginal => {
  const actual = await importOriginal<typeof import('@xyflow/react')>()
  return {
    ...actual,
    useReactFlow: () => reactFlowMocks,
  }
})

import ConceptGraphFlow from '../components/ConceptGraphFlow'

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

  it('shows a progress bar and mastery percentage on nodes with progress', () => {
    const { container } = render(
      <ConceptGraphFlow concepts={concepts} conceptProgress={{ count: 0.4 }} />
    )
    expect(container.querySelector('[data-progress="40"]')).toBeInTheDocument()
    // getByRole would compute an empty accessible name here: jsdom gives the
    // React Flow node wrapper no dimensions, so RF hides it with display:none.
    // Match the aria-label directly instead.
    const node = container.querySelector('.concept-node[aria-label*="40% toward mastery"]')
    expect(node).not.toBeNull()
    expect(node!.querySelector('.node-progress-tip')!.textContent).toBe('40% toward mastery')
  })

  it('shows a not-started tooltip on nodes without progress', () => {
    const { container } = render(
      <ConceptGraphFlow concepts={concepts} conceptProgress={{ count: 0.4 }} />
    )
    const node = container.querySelector('.concept-node[aria-label*="Addition"]')
    expect(node).not.toBeNull()
    expect(node!.querySelector('.node-progress-tip')!.textContent).toBe('Not started')
  })

  it('derives locked statuses that flow into node colors', () => {
    const statuses = deriveStatuses(concepts)
    expect(statuses['count']).toBe('unseen')
    expect(statuses['add']).toBe('locked')
    expect(statuses['mul']).toBe('locked')
  })

  it('opens the list view and selects a concept via its buttons', () => {
    const onSelectionChange = vi.fn()
    render(<ConceptGraphFlow concepts={concepts} onSelectionChange={onSelectionChange} />)
    fireEvent.click(screen.getByRole('button', { name: 'Open concept list' }))
    const dialog = screen.getByRole('dialog', { name: 'Concept list' })
    expect(dialog).toBeInTheDocument()
    fireEvent.click(within(dialog).getByRole('button', { name: /Addition/ }))
    expect(onSelectionChange).toHaveBeenCalledWith('add')
    expect(screen.queryByRole('dialog')).not.toBeInTheDocument()
  })

  it('filters the list view and reports no matches for junk queries', () => {
    render(<ConceptGraphFlow concepts={concepts} />)
    fireEvent.click(screen.getByRole('button', { name: 'Open concept list' }))
    const dialog = screen.getByRole('dialog', { name: 'Concept list' })
    fireEvent.change(within(dialog).getByLabelText('Filter concepts'), {
      target: { value: 'zzzz' },
    })
    expect(within(dialog).getByText('No matching concepts')).toBeInTheDocument()
    fireEvent.keyDown(within(dialog).getByLabelText('Filter concepts'), {
      key: 'Escape',
    })
    expect(screen.queryByRole('dialog')).not.toBeInTheDocument()
  })
})

describe('ConceptGraphFlow keyboard controls', () => {
  beforeEach(() => {
    reactFlowMocks.zoomIn.mockClear()
    reactFlowMocks.zoomOut.mockClear()
    reactFlowMocks.fitView.mockClear()
    reactFlowMocks.setViewport.mockClear()
    reactFlowMocks.getViewport.mockClear()
    reactFlowMocks.getViewport.mockImplementation(() => ({ x: 0, y: 0, zoom: 1 }))
  })

  function renderEngaged(props = {}) {
    const utils = render(<ConceptGraphFlow concepts={concepts} {...props} />)
    fireEvent.mouseEnter(screen.getByTestId('graph-wrapper'))
    return utils
  }

  it('zooms with plus and minus and fits with 0 when the graph is engaged', () => {
    renderEngaged()
    const wrapper = screen.getByTestId('graph-wrapper')
    fireEvent.keyDown(wrapper, { key: '+' })
    expect(reactFlowMocks.zoomIn).toHaveBeenCalledWith({ duration: 200 })
    fireEvent.keyDown(wrapper, { key: '-' })
    expect(reactFlowMocks.zoomOut).toHaveBeenCalledWith({ duration: 200 })
    fireEvent.keyDown(wrapper, { key: '0' })
    expect(reactFlowMocks.fitView).toHaveBeenCalledWith({
      padding: 0.15,
      maxZoom: 1,
      duration: 300,
    })
  })

  it('pans with arrow keys using the map camera convention', () => {
    renderEngaged()
    const wrapper = screen.getByTestId('graph-wrapper')
    fireEvent.keyDown(wrapper, { key: 'ArrowLeft' })
    expect(reactFlowMocks.setViewport).toHaveBeenCalledWith(
      { x: 90, y: 0, zoom: 1 },
      { duration: 200 }
    )
    fireEvent.keyDown(wrapper, { key: 'ArrowRight' })
    expect(reactFlowMocks.setViewport).toHaveBeenCalledWith(
      { x: -90, y: 0, zoom: 1 },
      { duration: 200 }
    )
  })

  it('ignores shortcuts while typing in the search input', () => {
    renderEngaged()
    fireEvent.keyDown(screen.getByPlaceholderText('Search concepts…'), { key: '+' })
    expect(reactFlowMocks.zoomIn).not.toHaveBeenCalled()
  })

  it('ignores shortcuts when the graph is neither focused nor hovered', () => {
    render(<ConceptGraphFlow concepts={concepts} />)
    fireEvent.keyDown(window, { key: '+' })
    expect(reactFlowMocks.zoomIn).not.toHaveBeenCalled()
  })

  it('Escape clears the selection when the list view is closed', () => {
    const onSelectionChange = vi.fn()
    renderEngaged({ selectedId: 'mul', onSelectionChange })
    fireEvent.keyDown(screen.getByTestId('graph-wrapper'), { key: 'Escape' })
    expect(onSelectionChange).toHaveBeenCalledWith(null)
  })

  it('Escape closes the list view before clearing the selection', () => {
    const onSelectionChange = vi.fn()
    renderEngaged({ selectedId: 'mul', onSelectionChange })
    fireEvent.click(screen.getByRole('button', { name: 'Open concept list' }))
    expect(screen.getByRole('dialog', { name: 'Concept list' })).toBeInTheDocument()
    // isolate from React Flow's mount-time empty-selection callback
    onSelectionChange.mockClear()
    fireEvent.keyDown(window, { key: 'Escape' })
    expect(screen.queryByRole('dialog')).not.toBeInTheDocument()
    expect(onSelectionChange).not.toHaveBeenCalledWith(null)
  })
})
