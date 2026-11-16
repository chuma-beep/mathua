import { describe, expect, it } from 'vitest'
import { render, screen } from '@testing-library/react'
import LeagueAvatar, { LEAGUE_TIER_ORDER, leagueFrameGeometry } from '../components/LeagueAvatar'

const base = { seed: 's1', name: 'Ada' }

describe('LeagueAvatar', () => {
  it.each(LEAGUE_TIER_ORDER)('draws a %s frame', (tier) => {
    render(<LeagueAvatar {...base} tier={tier} weeklyXp={12} size={72} />)
    const frame = screen.getByText(`${tier.charAt(0).toUpperCase()}${tier.slice(1)} league`).closest('[data-league-frame]')
    expect(frame).toHaveAttribute('data-league-frame', tier)
    expect(frame?.querySelector('svg[aria-hidden="true"]')).not.toBeNull()
  })

  it('stays plain without a tier', () => {
    const { container } = render(<LeagueAvatar {...base} weeklyXp={12} />)
    expect(container.querySelector('[data-league-frame]')).toBeNull()
    expect(screen.queryByText(/league$/)).toBeNull()
  })

  it.each([0, null, undefined])('stays plain when weekly XP is %s', (xp) => {
    const { container } = render(<LeagueAvatar {...base} tier="gold" weeklyXp={xp} />)
    expect(container.querySelector('[data-league-frame]')).toBeNull()
    expect(container.querySelector('svg[aria-hidden="true"]')).toBeNull()
  })

  it('never puts frame geometry inside the image radius', () => {
    for (const size of [24, 26, 32, 72]) {
      const geometry = leagueFrameGeometry(size)
      expect(geometry.frameInnerRadius).toBeGreaterThan(geometry.imageRadius)
      expect(geometry.outer).toBeGreaterThan(size)
    }
  })

  it('collapses the detailed geometry below 20 pixels', () => {
    const { container, unmount } = render(<LeagueAvatar {...base} tier="diamond" weeklyXp={12} size={18} />)
    expect(container.querySelector('[data-league-frame]')).toHaveAttribute('data-collapsed', 'true')
    expect(container.querySelectorAll('path')).toHaveLength(0)
    unmount()

    const large = render(<LeagueAvatar {...base} tier="diamond" weeklyXp={12} size={72} />)
    expect(large.container.querySelector('[data-league-frame]')).toHaveAttribute('data-collapsed', 'false')
    expect(large.container.querySelectorAll('path').length).toBeGreaterThan(0)
  })

})
