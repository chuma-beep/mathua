import { describe, it, expect } from 'vitest'
import { render, screen, fireEvent } from '@testing-library/react'
import Avatar from '../components/Avatar'

describe('Avatar', () => {
  it('renders the photo url when provided', () => {
    render(<Avatar seed="s1" name="Ada" url="https://x/a.png" />)
    expect(screen.getByAltText('Ada avatar')).toHaveAttribute('src', 'https://x/a.png')
  })

  it('falls back to the initial when the photo fails', () => {
    render(<Avatar seed="s1" name="Ada" url="https://x/broken.png" />)
    fireEvent.error(screen.getByAltText('Ada avatar'))
    expect(screen.getByLabelText('Ada avatar')).toHaveTextContent('A')
  })

  it('retries each new url instead of latching the failure', () => {
    const { rerender } = render(<Avatar seed="s1" name="Ada" url="https://x/one.png" />)
    fireEvent.error(screen.getByAltText('Ada avatar'))
    expect(screen.getByLabelText('Ada avatar')).toHaveTextContent('A')
    rerender(<Avatar seed="s1" name="Ada" url="https://x/two.png" />)
    expect(screen.getByAltText('Ada avatar')).toHaveAttribute('src', 'https://x/two.png')
  })
})
