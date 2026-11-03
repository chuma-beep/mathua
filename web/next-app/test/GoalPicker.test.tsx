import { describe, it, expect, vi } from 'vitest'
import { render, screen, fireEvent } from '@testing-library/react'
import GoalPicker from '../components/GoalPicker'

const concepts = [
  { id: 'a', label: 'Alpha', domain: 'arithmetic' },
  { id: 'b', label: 'Beta', domain: 'algebra' },
  { id: 'c', label: 'Gamma', domain: 'geometry' },
]

describe('GoalPicker', () => {
  it('filters concepts and selects one', () => {
    const onSelect = vi.fn()
    const onClose = vi.fn()
    render(<GoalPicker concepts={concepts} onSelect={onSelect} onClose={onClose} />)

    expect(screen.getByRole('dialog', { name: 'Choose a goal' })).toBeInTheDocument()
    fireEvent.change(screen.getByLabelText('Search concepts for a goal'), {
      target: { value: 'Beta' },
    })

    expect(screen.getByText('Beta')).toBeInTheDocument()
    expect(screen.queryByText('Alpha')).not.toBeInTheDocument()
    expect(screen.queryByText('Gamma')).not.toBeInTheDocument()

    fireEvent.click(screen.getByRole('button', { name: /Beta/ }))
    expect(onSelect).toHaveBeenCalledWith('b')
  })

  it('closes on Escape and marks the current goal', () => {
    const onClose = vi.fn()
    render(
      <GoalPicker
        concepts={concepts}
        currentGoalId="b"
        onSelect={vi.fn()}
        onClose={onClose}
      />,
    )

    expect(screen.getByText('Current goal')).toBeInTheDocument()
    fireEvent.keyDown(window, { key: 'Escape' })
    expect(onClose).toHaveBeenCalled()
  })
})
