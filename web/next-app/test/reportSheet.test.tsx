import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import ReportButton from '../components/ReportButton'
import { submitReport } from '../lib/api'

vi.mock('../lib/api', async (importOriginal) => ({
  ...(await importOriginal<typeof import('../lib/api')>()),
  submitReport: vi.fn(),
}))

// The report control is the one channel a learner has when a question is
// broken, so it must survive the diagnostic re-rendering underneath it. It used
// to be an inline block that moved when a verdict appeared, which dismissed the
// Android picker about half a second after it opened and left the learner with
// nothing to type in.
describe('the report sheet', () => {
  beforeEach(() => {
    vi.mocked(submitReport).mockReset()
    vi.mocked(submitReport).mockResolvedValue(undefined)
  })

  it('stays open when its parent re-renders', async () => {
    const { rerender } = render(<ReportButton kind="question" conceptId="a" question="q1" />)
    fireEvent.click(screen.getByRole('button', { name: /report a problem/i }))

    const dialog = screen.getByRole('dialog', { name: /report a problem/i })
    expect(dialog).toBeInTheDocument()

    rerender(<ReportButton kind="question" conceptId="a" question="q2" />)

    expect(screen.getByRole('dialog', { name: /report a problem/i })).toBeInTheDocument()
  })

  it('keeps the typed detail and sends it', async () => {
    render(<ReportButton kind="question" conceptId="a" question="q1" />)
    fireEvent.click(screen.getByRole('button', { name: /report a problem/i }))

    const box = screen.getByRole('textbox')
    fireEvent.change(box, { target: { value: 'the equation lost its exponent' } })
    expect(box).toHaveValue('the equation lost its exponent')

    fireEvent.click(screen.getByRole('button', { name: /send report/i }))

    await waitFor(() => expect(submitReport).toHaveBeenCalledTimes(1))
    expect(vi.mocked(submitReport).mock.calls[0][0]).toMatchObject({
      detail: 'the equation lost its exponent',
      concept_id: 'a',
      kind: 'question',
    })
  })

  it('resets the confirmation when the question changes', async () => {
    const { rerender } = render(<ReportButton kind="question" conceptId="a" question="q1" />)
    fireEvent.click(screen.getByRole('button', { name: /report a problem/i }))
    fireEvent.click(screen.getByRole('button', { name: /send report/i }))

    await waitFor(() => expect(screen.getByText(/thanks — reported/i)).toBeInTheDocument())

    rerender(<ReportButton kind="question" conceptId="a" question="q2" />)
    expect(screen.queryByText(/thanks — reported/i)).not.toBeInTheDocument()
    expect(screen.getByRole('button', { name: /report a problem/i })).toBeInTheDocument()
  })
})
