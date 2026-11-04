import { describe, expect, it } from 'vitest'
import { render, screen } from '@testing-library/react'
import KatexContent from '../components/KatexContent'

describe('KatexContent links', () => {
  it('keeps links as text by default', () => {
    render(<KatexContent>{'https://myhappr.com/wisddom'}</KatexContent>)
    expect(screen.queryByRole('link')).toBeNull()
  })

  it('renders a safe external link when links are allowed', () => {
    render(<KatexContent allowLinks>{'https://myhappr.com/wisddom'}</KatexContent>)
    const link = screen.getByRole('link', { name: 'https://myhappr.com/wisddom' })
    expect(link).toHaveAttribute('href', 'https://myhappr.com/wisddom')
    expect(link).toHaveAttribute('target', '_blank')
    expect(link).toHaveAttribute('rel', 'noreferrer noopener')
  })
})
