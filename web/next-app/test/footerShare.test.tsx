import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import { afterEach, describe, expect, it, vi } from 'vitest'
import Footer from '../components/Footer'
import { SITE_URL } from '../lib/site'

const originalShare = navigator.share
const originalClipboard = navigator.clipboard

afterEach(() => {
  Object.defineProperty(navigator, 'share', { configurable: true, value: originalShare })
  Object.defineProperty(navigator, 'clipboard', { configurable: true, value: originalClipboard })
  vi.restoreAllMocks()
})

describe('the footer share action', () => {
  it('shows the support-by-sharing message', () => {
    render(<Footer />)
    expect(
      screen.getByText('Know someone who wants to learn? Share Mathua with them.')
    ).toBeInTheDocument()
  })

  it('uses the native share sheet when the browser provides one', async () => {
    const share = vi.fn().mockResolvedValue(undefined)
    Object.defineProperty(navigator, 'share', { configurable: true, value: share })

    render(<Footer />)
    fireEvent.click(screen.getByRole('button', { name: 'Share Mathua' }))

    await waitFor(() =>
      expect(share).toHaveBeenCalledWith({
        title: 'Mathua',
        text: 'Learn math with Mathua.',
        url: SITE_URL,
      })
    )
  })

  it('copies the site URL and confirms when native sharing is unavailable', async () => {
    const writeText = vi.fn().mockResolvedValue(undefined)
    Object.defineProperty(navigator, 'share', { configurable: true, value: undefined })
    Object.defineProperty(navigator, 'clipboard', {
      configurable: true,
      value: { writeText },
    })

    render(<Footer />)
    fireEvent.click(screen.getByRole('button', { name: 'Share Mathua' }))

    await waitFor(() => expect(writeText).toHaveBeenCalledWith(SITE_URL))
    expect(await screen.findByRole('button', { name: 'Link copied' })).toBeInTheDocument()
  })
})
