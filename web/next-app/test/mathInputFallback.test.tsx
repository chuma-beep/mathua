// The fallback path: what happens when MathLive itself cannot load.
//
// Its own file because the only way to reach the boundary is to make the module
// fail, and a module mock is file-scoped — mocking it here would stop every other
// MathInput test from ever rendering the real editor.
//
// The scenario is not hypothetical. MathLive is ~800 kB and arrives as a lazily
// fetched chunk: a blocked asset, an aggressive proxy, a flaky connection, or a
// browser the build does not target. A learner in that situation must still be
// able to answer, and the plain input they fall back to grades identically for
// every concept in the corpus.

import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, waitFor, fireEvent } from '@testing-library/react'

vi.mock('../components/math/MathLiveField', () => ({
  default: function BrokenMathLiveField(): never {
    throw new Error('Failed to fetch dynamically imported module: MathLiveField')
  },
}))

const { MathAnswerInput } = await import('../components/math/MathInput')

describe('MathAnswerInput fallback', () => {
  beforeEach(() => {
    vi.spyOn(console, 'error').mockImplementation(() => {})
  })

  it('renders a working plain input when MathLive throws on mount', async () => {
    const seen: string[] = []
    render(
      <MathAnswerInput
        value=""
        onChange={v => seen.push(v)}
        gradingType="numeric"
        conceptId="arith.add.single"
        placeholder="Your answer"
      />,
    )

    // The fallback is a real input, not a dead placeholder.
    const input = await screen.findByPlaceholderText(/your answer/i)
    expect(input.tagName.toLowerCase()).toBe('input')

    // And it is wired: typing reaches the host exactly as the editor's would.
    fireEvent.change(input, { target: { value: '1/2' } })
    expect(seen[seen.length - 1]).toBe('1/2')
  })

  it('still submits on Enter after falling back', async () => {
    const onSubmit = vi.fn()
    render(
      <MathAnswerInput
        value=""
        onChange={() => {}}
        onSubmit={onSubmit}
        gradingType="numeric"
        conceptId="arith.add.single"
        placeholder="Your answer"
      />,
    )
    const input = await screen.findByPlaceholderText(/your answer/i)
    fireEvent.keyDown(input, { key: 'Enter' })
    expect(onSubmit).toHaveBeenCalledTimes(1)
  })

  it('logs the failure rather than swallowing it', async () => {
    render(<MathAnswerInput value="" onChange={() => {}} gradingType="numeric" conceptId="arith.add.single" placeholder="Your answer" />)
    await waitFor(() => expect(screen.getByPlaceholderText(/your answer/i)).toBeTruthy())
    // A silent 800 kB failure is indistinguishable from a broken deploy.
    expect(console.error).toHaveBeenCalledWith(
      expect.stringContaining('mathlive failed to load'),
      expect.anything(),
    )
  })
})
