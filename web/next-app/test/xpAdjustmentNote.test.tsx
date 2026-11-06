import { describe, expect, it } from 'vitest'
import { render, screen } from '@testing-library/react'
import XPAdjustmentNote from '../components/XPAdjustmentNote'

describe('XPAdjustmentNote', () => {
  it('renders an actual negative performance adjustment and its explanation', () => {
    render(
      <XPAdjustmentNote
        adjustment={-2}
        note="-2 XP: Task performance was extremely poor across this task."
        source="performance"
      />,
    )
    const line = screen.getByRole('status')
    expect(line.textContent).toContain('Performance')
    expect(line.textContent).toContain('-2 XP')
    expect(line.textContent).toContain('extremely poor')
  })

  it('renders a zero adjustment explanation rather than hiding it', () => {
    render(
      <XPAdjustmentNote
        adjustment={0}
        note="0 XP: Task performance was poor."
        source="performance"
      />,
    )
    expect(screen.getByRole('status').textContent).toContain('0 XP')
  })

  it('renders nothing when there is neither an amount nor a note', () => {
    const { container } = render(<XPAdjustmentNote adjustment={0} />)
    expect(container.textContent).toBe('')
  })
})
