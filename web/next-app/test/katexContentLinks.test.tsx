import { describe, expect, it } from 'vitest'
import { render, screen } from '@testing-library/react'
import KatexContent from '../components/KatexContent'

describe('KatexContent links', () => {
  it('keeps links as text by default', () => {
    render(<KatexContent>{'https://selar.com/showlove/wisddom'}</KatexContent>)
    expect(screen.queryByRole('link')).toBeNull()
  })

  it('renders a safe external link when links are allowed', () => {
    render(<KatexContent allowLinks>{'https://selar.com/showlove/wisddom'}</KatexContent>)
    const link = screen.getByRole('link', { name: 'https://selar.com/showlove/wisddom' })
    expect(link).toHaveAttribute('href', 'https://selar.com/showlove/wisddom')
    expect(link).toHaveAttribute('target', '_blank')
    expect(link).toHaveAttribute('rel', 'noreferrer noopener')
  })
})
