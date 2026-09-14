import { describe, it, expect, vi } from 'vitest'
import { render, screen, fireEvent } from '@testing-library/react'
import ProgressCardList from '../components/ProgressCardList'
import ProgressCardPage from '../app/progress-card/page'
import type { AttemptRecord } from '../lib/api'

// next/navigation stub for tests without a router
vi.mock('next/navigation', () => ({
  usePathname: () => '/progress-card',
  useRouter: () => ({ push: () => {} }),
}))

vi.mock('../lib/api', async (importOriginal) => {
  const mod = await importOriginal<typeof import('../lib/api')>()
  return {
    ...mod,
    getAttempts: vi.fn(async () => ({ attempts: [], total: 0 })),
  }
})

const rows: AttemptRecord[] = [
  {
    session_id: 's1',
    student_id: 'u1',
    concept_id: 'arith.add.single',
    concept_name: 'Single-digit addition',
    answer: '5',
    expected: '4',
    correct: false,
    elapsed_seconds: 3,
    timestamp: '2026-09-10T10:00:00Z',
    question: '2 + 2 = ?',
    source: 'diagnostic',
    explanation: '2 + 2 = 4',
  },
  {
    session_id: 's1',
    student_id: 'u1',
    concept_id: 'arith.add.single',
    concept_name: 'Single-digit addition',
    answer: '4',
    expected: '4',
    correct: true,
    elapsed_seconds: 2,
    timestamp: '2026-09-11T10:00:00Z',
    question: '2 + 2 = ?',
    source: 'quiz',
    explanation: '',
  },
]

describe('ProgressCardList', () => {
  it('renders verdicts, corrections, and concept links', () => {
    render(<ProgressCardList attempts={rows} />)
    expect(screen.getAllByText('2 + 2 = ?').length).toBe(2)
    expect(screen.getByText(/You: 5/)).toBeInTheDocument()
    expect(screen.getByText(/Correct: 4/)).toBeInTheDocument()
    expect(screen.getAllByText('Single-digit addition →').length).toBeGreaterThan(0)
    expect(screen.getByText('2 + 2 = 4')).toBeInTheDocument()
  })

  it('marks legacy rows without question text', () => {
    render(<ProgressCardList attempts={[{ ...rows[0], question: '' }]} />)
    expect(screen.getByText(/question unavailable/i)).toBeInTheDocument()
  })

  it('explains the empty state', () => {
    render(<ProgressCardList attempts={[]} />)
    expect(screen.getByText(/no answered questions yet/i)).toBeInTheDocument()
  })
})

describe('ProgressCardPage', () => {
  it('defaults to missed filter and switches to all', async () => {
    const { getAttempts } = await import('../lib/api')
    vi.mocked(getAttempts).mockResolvedValue({ attempts: rows, total: 2 })
    render(<ProgressCardPage />)
    await screen.findByText('Showing 2 of 2')
    expect(vi.mocked(getAttempts)).toHaveBeenCalledWith(
      expect.objectContaining({ incorrect_only: true }),
    )
    fireEvent.click(screen.getByRole('button', { name: 'All' }))
    await screen.findByText('Showing 2 of 2')
    expect(vi.mocked(getAttempts)).toHaveBeenLastCalledWith(
      expect.not.objectContaining({ incorrect_only: true }),
    )
  })

  it('filters by source', async () => {
    const { getAttempts } = await import('../lib/api')
    vi.mocked(getAttempts).mockResolvedValue({ attempts: [rows[1]], total: 1 })
    render(<ProgressCardPage />)
    await screen.findByText('Showing 1 of 1')
    fireEvent.click(screen.getByRole('button', { name: 'quiz' }))
    await screen.findByText('Showing 1 of 1')
    expect(vi.mocked(getAttempts)).toHaveBeenLastCalledWith(
      expect.objectContaining({ source: 'quiz' }),
    )
  })
})
