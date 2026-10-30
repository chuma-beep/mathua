import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, fireEvent } from '@testing-library/react'
import FriendsPage from '../app/friends/page'
import { getFriends, getFriendsLeaderboard } from '../lib/api'

vi.mock('../components/app-shell', () => ({
  default: ({ children }: { children: React.ReactNode }) => <div>{children}</div>,
}))

vi.mock('../hooks/useAuthState', () => ({
  useAuthState: () => ({ loggedIn: true, user: null }),
}))

vi.mock('../lib/api', async (importOriginal) => {
  const actual = await importOriginal<typeof import('../lib/api')>()
  return {
    ...actual,
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
    vi.mocked(getFriends).mockResolvedValue(feed)
    render(<FriendsPage />)
    expect(await screen.findByText('ABCD2345')).toBeDefined()
    expect(await screen.findByText('Grace Hopper')).toBeDefined()
    expect(screen.getByText(/4 mastered/)).toBeDefined()
    // One list, not a friends list plus a separate standings block.
    expect(screen.queryByText('Friends leaderboard')).toBeNull()
    expect(screen.getAllByText('—').length).toBeGreaterThan(0)
  })

  it('explains that the surface is off when the API answers 404', async () => {
    vi.mocked(getFriends).mockRejectedValue(Object.assign(new Error('not found'), { status: 404 }))
    render(<FriendsPage />)
    expect(await screen.findByText('Friends are not enabled on this server.')).toBeDefined()
  })

  it('asks a signed-out learner to sign in instead of claiming the feature is off', async () => {
    vi.mocked(getFriends).mockRejectedValue(Object.assign(new Error('not authenticated'), { status: 401 }))
    render(<FriendsPage />)
    expect(await screen.findByText('Sign in to use friends.')).toBeDefined()
  })
})

describe('friend row actions', () => {
  it('keeps Block and Remove behind one menu', async () => {
    vi.mocked(getFriends).mockResolvedValue(feed)
    render(<FriendsPage />)
    fireEvent.click(await screen.findByRole('button', { name: 'More actions for Grace Hopper' }))
    expect(screen.getByRole('menuitem', { name: 'Block' })).toBeDefined()
    expect(screen.getByRole('menuitem', { name: 'Remove friend' })).toBeDefined()
  })
})
