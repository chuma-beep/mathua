import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import StrugglesSection from '../components/StrugglesSection'
import type { WeaknessRes } from '../lib/api'

const weaknesses: WeaknessRes = {
  by_domain: {
    algebra: [
      { id: 'alg.quad.incomplete', label: 'Incomplete quadratics', weakness: 0.7 },
      { id: 'alg.log.properties', label: 'Log properties', weakness: 0.2 },
    ],
  },
}

describe('StrugglesSection', () => {
  it('sends weakness work to practice, not spaced review', () => {
    render(<StrugglesSection weaknesses={weaknesses} />)

    const link = screen.getByRole('link', { name: 'Practice now →' })
    expect(link.getAttribute('href')).toBe('/learn')
    expect(screen.queryByRole('link', { name: 'Review now →' })).toBeNull()
    expect(screen.getByText('1 struggling')).toBeTruthy()
  })

  it('renders nothing when there are no weak concepts', () => {
    const { container } = render(<StrugglesSection weaknesses={{ by_domain: {} }} />)
    expect(container.firstChild).toBeNull()
  })
})
