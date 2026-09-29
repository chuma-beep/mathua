import { describe, it, expect, vi } from 'vitest'
import { render, screen, fireEvent } from '@testing-library/react'
import DangerZone, { attemptsToCSV } from '../components/DangerZone'

vi.mock('next/navigation', () => ({
  useRouter: () => ({ push: vi.fn() }),
}))

describe('DangerZone', () => {
  it('keeps reset disabled until the exact phrase is typed', () => {
    render(<DangerZone />)
    const button = screen.getByRole('button', { name: /Reset everything above/ })
    expect(button).toBeDisabled()
    fireEvent.change(screen.getByLabelText(/Type.*to confirm/i), { target: { value: 'reset' } })
    expect(button).toBeDisabled()
    fireEvent.change(screen.getByLabelText(/Type.*to confirm/i), { target: { value: 'reset my progress' } })
    expect(button).not.toBeDisabled()
  })

  it('states the keeps list including destination, deadline and pace', () => {
    render(<DangerZone />)
    expect(screen.getByText(/Destination, deadline and pace/i)).toBeTruthy()
    expect(screen.getByText(/profile and leaderboard/i)).toBeTruthy()
  })
})

describe('attemptsToCSV', () => {
  it('quotes and headers rows', () => {
    const csv = attemptsToCSV([
      { session_id: 's', student_id: 'u', concept_id: 'a', answer: '4', expected: '4', correct: true, elapsed_seconds: 3, timestamp: '2026-01-01', question: '2+2?', source: 'practice', explanation: '' },
      { session_id: 's', student_id: 'u', concept_id: 'b', answer: 'x "y"', expected: 'z', correct: false, elapsed_seconds: 9, timestamp: '2026-01-02', question: 'q', source: 'quiz', explanation: '' },
    ])
    const lines = csv.split('\n')
    expect(lines[0]).toBe('timestamp,concept_id,question,answer,expected,correct,elapsed_seconds,source')
    expect(lines).toHaveLength(3)
    expect(lines[2]).toContain('"x ""y"""')
  })
})
