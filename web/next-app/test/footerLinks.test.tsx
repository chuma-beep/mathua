import { render, screen } from '@testing-library/react'
import { describe, expect, it } from 'vitest'
import Footer from '../components/Footer'
import { CREATOR_LINKS, SUPPORT_URL } from '../lib/creatorLinks'

// The footer is the one surface every page shares. A support link that points at
// the private dashboard, or a social profile that still has the old handle, is a
// site-wide mistake rather than a one-page one, so it is worth asserting here.
describe('the footer creator links', () => {
  it('opens the public support page, not the creator dashboard', () => {
    render(<Footer />)

    const support = screen.getByRole('link', { name: 'Support Mathua' })
    expect(support).toHaveAttribute('href', SUPPORT_URL)
    expect(support).toHaveAttribute('target', '_blank')
    expect(support).toHaveAttribute('rel', 'noreferrer')
    expect(support.getAttribute('href')).not.toContain('dashboard')
  })

  it('publishes every public profile beside the project links', () => {
    render(<Footer />)

    for (const link of CREATOR_LINKS) {
      const anchor = screen.getByRole('link', { name: link.label })
      expect(anchor).toHaveAttribute('href', link.href)
      expect(anchor).toHaveAttribute('target', '_blank')
      expect(anchor).toHaveAttribute('rel', 'noreferrer')
    }
  })
})
