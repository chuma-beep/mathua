import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen } from '@testing-library/react'
import FriendsPage from '../app/friends/page'
import { getConfig, getFriends, getFriendsLeaderboard } from '../lib/api'

vi.mock('../components/app-shell', () => ({
  default: ({ children }: { children: React.ReactNode }) => <div>{children}</div>,
}))

vi.mock('../lib/api', async (importOriginal) => {
  const actual = await importOriginal<typeof import('../lib/api')>()
  return {
    ...actual,
    getConfig: vi.fn(),
    getFriends: vi.fn(),
    getFriendsLeaderboard: vi.fn(),
    sendFriendRequest: vi.fn(),
    acceptFriendRequest: vi.fn(),
    declineFriendRequest: vi.fn(),
    cancelFriendRequest: vi.fn(),
    removeFriend: vi.fn(),
    blockFriend: vi.fn(),
    rotateFriendCode: vi.fn(),
  }
})

const feed = {
  identity: { profile_id: 'p_ABCDEFGHJKLMNPQR', friend_code: 'ABCD2345' },
  friends: [{
    profile_id: 'p_ZZZZZZZZZZZZZZZZ',
    name: 'Grace Hopper',
    username: 'grace',
    avatar: {},
    xp_total: 120,
    mastered: 4,
    friends_since: '2026-10-30T00:00:00Z',
  }],
  incoming: [],
  outgoing: [],
}

describe('FriendsPage', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    vi.mocked(getFriendsLeaderboard).mockResolvedValue([])
  })

  it('shows the friend code and accepted friends when the surface is enabled', async () => {
    vi.mocked(getConfig).mockResolvedValue({ auth_enabled: true, friends_enabled: true })
    vi.mocked(getFriends).mockResolvedValue(feed)
    render(<FriendsPage />)
    expect(await screen.findByText('ABCD2345')).toBeDefined()
    expect(await screen.findByText('Grace Hopper')).toBeDefined()
    expect(screen.getByText(/4 mastered/)).toBeDefined()
  })

  it('explains that the surface is off rather than showing a broken page', async () => {
    vi.mocked(getConfig).mockResolvedValue({ auth_enabled: true, friends_enabled: false })
    render(<FriendsPage />)
    expect(await screen.findByText('Friends are not enabled on this server.')).toBeDefined()
    expect(getFriends).not.toHaveBeenCalled()
  })
})
