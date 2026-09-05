import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'

const pushMock = vi.fn()
const getSettingsMock = vi.fn()
const updateSettingsMock = vi.fn()
const updateProfileNameMock = vi.fn()
const uploadAvatarMock = vi.fn()
const deleteAvatarMock = vi.fn()

let mockUser: { student_id: string; name: string; username: string } | null = {
  student_id: 's1',
  name: 'Ada',
  username: 'ada',
}
const setUserInfoMock = vi.fn((info: unknown) => {
  mockUser = info as typeof mockUser
})

vi.mock('next/navigation', () => ({
  useRouter: () => ({ push: pushMock }),
}))

vi.mock('../lib/api', () => ({
  getSettings: (...a: unknown[]) => getSettingsMock(...a),
  updateSettings: (...a: unknown[]) => updateSettingsMock(...a),
  updateProfileName: (...a: unknown[]) => updateProfileNameMock(...a),
  uploadAvatar: (...a: unknown[]) => uploadAvatarMock(...a),
  deleteAvatar: (...a: unknown[]) => deleteAvatarMock(...a),
  avatarImageUrl: () => 'https://x/avatar',
  enableShare: vi.fn(),
  disableShare: vi.fn(),
}))

vi.mock('../lib/auth', () => ({
  isLoggedIn: () => true,
  getUserInfo: () => mockUser,
  setUserInfo: (info: unknown) => setUserInfoMock(info),
}))

vi.mock('../components/Header', () => ({ default: () => <div data-testid="header-stub" /> }))
vi.mock('../components/Footer', () => ({ default: () => <div data-testid="footer-stub" /> }))
vi.mock('../components/BottomTabs', () => ({ default: () => <div data-testid="tabs-stub" /> }))

import SettingsPage from '../app/settings/page'
import { DICEBEAR_STYLES } from '../lib/dicebear'

describe('Settings Profile section', () => {
  beforeEach(() => {
    mockUser = { student_id: 's1', name: 'Ada', username: 'ada' }
    getSettingsMock.mockReset().mockResolvedValue({})
    updateSettingsMock.mockReset().mockResolvedValue(undefined)
    updateProfileNameMock.mockReset().mockImplementation(async (name: string) => ({ student_id: 's1', name }))
    uploadAvatarMock.mockReset().mockResolvedValue(undefined)
    deleteAvatarMock.mockReset().mockResolvedValue(undefined)
    setUserInfoMock.mockClear()
  })

  it('prefills the display name and saves it', async () => {
    render(<SettingsPage />)
    const input = await screen.findByLabelText('Display name')
    expect((input as HTMLInputElement).value).toBe('Ada')
    fireEvent.change(input, { target: { value: '  Ada Lovelace  ' } })
    fireEvent.click(screen.getByText('Save name'))
    await waitFor(() => expect(updateProfileNameMock).toHaveBeenCalledWith('Ada Lovelace'))
    expect(setUserInfoMock).toHaveBeenCalledWith(expect.objectContaining({ name: 'Ada Lovelace', username: 'ada' }))
    expect(await screen.findByText('Name saved.')).toBeInTheDocument()
  })

  it('renders one gallery button per curated style', async () => {
    render(<SettingsPage />)
    await screen.findByText('Pick a character')
    for (const style of DICEBEAR_STYLES) {
      expect(screen.getByRole('button', { name: `Pick ${style} character` })).toBeInTheDocument()
    }
  })

  it('Surprise me persists a valid dicebear pick and drops the legacy preset', async () => {
    getSettingsMock.mockResolvedValue({ avatar_preset: 3 })
    render(<SettingsPage />)
    fireEvent.click(await screen.findByText(/Surprise me/))
    await waitFor(() => expect(updateSettingsMock).toHaveBeenCalled())
    const saved = updateSettingsMock.mock.calls[0][0]
    expect(saved.avatar_dicebear.style).toMatch(/^[a-z0-9-]+$/)
    expect(saved.avatar_dicebear.seed.length).toBeGreaterThan(0)
    expect(saved).not.toHaveProperty('avatar_preset')
  })

  it('rejects oversize uploads client-side without calling the API', async () => {
    render(<SettingsPage />)
    await screen.findByText('Pick a character')
    const input = document.querySelector('input[type="file"]') as HTMLInputElement
    const big = new File([new Uint8Array(600 * 1024)], 'big.png', { type: 'image/png' })
    fireEvent.change(input, { target: { files: [big] } })
    expect(await screen.findByText(/512KB or smaller/)).toBeInTheDocument()
    expect(uploadAvatarMock).not.toHaveBeenCalled()
  })

  it('shows Remove photo only when a custom photo is set', async () => {
    getSettingsMock.mockResolvedValue({ avatar_custom: true })
    render(<SettingsPage />)
    expect(await screen.findByText('Remove photo')).toBeInTheDocument()
  })
})
