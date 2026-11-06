import { describe, expect, it } from 'vitest'
import { render, screen } from '@testing-library/react'
import XPTaskSummary from '../components/XPTaskSummary'

describe('XPTaskSummary', () => {
  it('shows answer XP, the actual adjustment and net XP', () => {
    render(
      <XPTaskSummary
        answerXP={3}
        adjustment={-5}
        note="-5 XP: Task performance was extremely poor across this task."
        source="performance"
      />,
    )
    const summary = screen.getByRole('group', { name: 'Task XP summary' })
    expect(summary.textContent).toContain('XP from answers')
    expect(summary.textContent).toContain('+3 XP')
    expect(summary.textContent).toContain('Performance adjustment')
    expect(summary.textContent).toContain('-5 XP')
    expect(summary.textContent).toContain('Net XP')
    expect(summary.textContent).toContain('-2 XP')
    expect(summary.textContent).toContain('extremely poor')
  })

  it('shows a zero adjustment and its explanation', () => {
    render(
      <XPTaskSummary
        answerXP={0}
        adjustment={0}
        note="0 XP: Task performance was poor."
        source="performance"
      />,
    )
    const summary = screen.getByRole('group', { name: 'Task XP summary' })
    expect(summary.textContent).toContain('0 XP')
    expect(summary.textContent).toContain('Task performance was poor')
    expect(summary.textContent).not.toContain('No accountability adjustment')
  })

  it('shows no-adjustment copy when there is no adjustment or note', () => {
    render(<XPTaskSummary answerXP={4} adjustment={0} />)
    const summary = screen.getByRole('group', { name: 'Task XP summary' })
    expect(summary.textContent).toContain('+4 XP')
    expect(summary.textContent).toContain('No accountability adjustment on this task.')
  })

  it('uses the actual capped adjustment for net XP', () => {
    render(
      <XPTaskSummary
        answerXP={3}
        adjustment={-1}
        note="-1 XP: Task performance was extremely poor across this task."
        source="performance"
      />,
    )
    const summary = screen.getByRole('group', { name: 'Task XP summary' })
    expect(summary.textContent).toContain('-1 XP')
    expect(summary.textContent).toContain('+2 XP')
    expect(summary.textContent).not.toContain('-2 XP')
  })
})
