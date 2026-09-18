import { describe, it, expect, vi } from 'vitest'
import { render } from '@testing-library/react'
import KatexContent from '../components/KatexContent'

describe('KatexContent KaTeX fallback', () => {
  it('replaces KaTeX errors with a labeled code fallback and keeps valid math', () => {
    const err = vi.spyOn(console, 'error').mockImplementation(() => {})
    try {
      const { container } = render(
        <KatexContent>{'ok $x^2$ and broken $\\boguscommandxyz$ here'}</KatexContent>
      )
      // No red raw-TeX error nodes leak to the user.
      expect(container.querySelector('.katex-error')).toBeNull()
      expect(container.innerHTML).not.toMatch(/204,\s*0,\s*0|#cc0000/i)
      const fb = container.querySelector('.math-fallback')
      expect(fb).not.toBeNull()
      expect(fb?.textContent).toContain('\\boguscommandxyz')
      expect(fb?.getAttribute('title')).toMatch(/failed to render/i)
      // Valid math on the same page still renders.
      expect(container.querySelector('.katex')).not.toBeNull()
      // Failure is logged for telemetry instead of silent.
      expect(err).toHaveBeenCalled()
    } finally {
      err.mockRestore()
    }
  })

  it('renders clean pages with no fallback and no error logging', () => {
    const err = vi.spyOn(console, 'error').mockImplementation(() => {})
    try {
      const { container } = render(
        <KatexContent>{'The integral $\\int_0^1 x\\,dx = 1/2$ converges.'}</KatexContent>
      )
      expect(container.querySelector('.katex-error')).toBeNull()
      expect(container.querySelector('.math-fallback')).toBeNull()
      expect(container.querySelector('.katex')).not.toBeNull()
      expect(err).not.toHaveBeenCalled()
    } finally {
      err.mockRestore()
    }
  })
})
