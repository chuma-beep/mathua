import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'

const pushMock = vi.fn()
const getSettingsMock = vi.fn()
const updateSettingsMock = vi.fn()
const updateProfileNameMock = vi.fn()
const changePasswordMock = vi.fn()
const uploadAvatarMock = vi.fn()
const deleteAvatarMock = vi.fn()
const getIdentitiesMock = vi.fn()
const deleteIdentityMock = vi.fn()
const createLinkTokenMock = vi.fn()
const requestEmailVerificationMock = vi.fn()
const startOAuthLoginMock = vi.fn()

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
  changePassword: (...a: unknown[]) => changePasswordMock(...a),
  getIdentities: (...a: unknown[]) => getIdentitiesMock(...a),
  deleteIdentity: (...a: unknown[]) => deleteIdentityMock(...a),
  createLinkToken: (...a: unknown[]) => createLinkTokenMock(...a),
  requestEmailVerification: (...a: unknown[]) => requestEmailVerificationMock(...a),
  uploadAvatar: (...a: unknown[]) => uploadAvatarMock(...a),
  deleteAvatar: (...a: unknown[]) => deleteAvatarMock(...a),
  getConfig: () => Promise.resolve({ providers: ['google', 'github'] }),
  startOAuthLogin: (...a: unknown[]) => startOAuthLoginMock(...a),
  OAUTH_LABELS: { google: 'Google', github: 'GitHub', facebook: 'Facebook', microsoft: 'Microsoft', apple: 'Apple' },
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
    getIdentitiesMock.mockReset().mockResolvedValue([])
    deleteIdentityMock.mockReset().mockResolvedValue(undefined)
    createLinkTokenMock.mockReset().mockResolvedValue('linktok')
    requestEmailVerificationMock.mockReset().mockResolvedValue(undefined)
    startOAuthLoginMock.mockClear()
    updateSettingsMock.mockReset().mockResolvedValue(undefined)
    updateProfileNameMock.mockReset().mockImplementation(async (name: string) => ({ student_id: 's1', name }))
    changePasswordMock.mockReset().mockResolvedValue(undefined)
    uploadAvatarMock.mockReset().mockResolvedValue(undefined)
    deleteAvatarMock.mockReset().mockResolvedValue(undefined)
    setUserInfoMock.mockClear()
  })

  // jsdom lacks implicit form submission: clicking a submit button does not
  // fire submit. Submit the enclosing form directly instead.
  const submitForm = (el: HTMLElement) => {
    const form = el.closest('form')
    if (!form) throw new Error('expected a form ancestor')
    fireEvent.submit(form)
  }

  it('prefills the display name and saves it', async () => {
    render(<SettingsPage />)
    const input = await screen.findByLabelText('Display name')
    expect((input as HTMLInputElement).value).toBe('Ada')
    fireEvent.change(input, { target: { value: '  Ada Lovelace  ' } })
    submitForm(screen.getByText('Save profile'))
    await waitFor(() => expect(updateProfileNameMock).toHaveBeenCalledWith('Ada Lovelace', ''))
    expect(setUserInfoMock).toHaveBeenCalledWith(expect.objectContaining({ name: 'Ada Lovelace', username: 'ada' }))
    expect(await screen.findByText('Profile saved.')).toBeInTheDocument()
  })

  it('saves recovery email alongside the name', async () => {
    render(<SettingsPage />)
    fireEvent.change(await screen.findByLabelText('Display name'), { target: { value: 'Ada' } })
    fireEvent.change(screen.getByLabelText(/Recovery email/, { exact: false }), { target: { value: 'ada@example.com' } })
    submitForm(screen.getByText('Save profile'))
    await waitFor(() => expect(updateProfileNameMock).toHaveBeenCalledWith('Ada', 'ada@example.com'))
  })

  it('surfaces a taken-email rejection instead of the generic error', async () => {
    updateProfileNameMock.mockRejectedValue(Object.assign(new Error('Update profile failed: 409: email already in use'), { status: 409, serverMessage: 'email already in use' }))
    render(<SettingsPage />)
    fireEvent.change(await screen.findByLabelText('Display name'), { target: { value: 'Ada' } })
    fireEvent.change(screen.getByLabelText(/Recovery email/, { exact: false }), { target: { value: 'taken@example.com' } })
    submitForm(screen.getByText('Save profile'))
    expect(await screen.findByText('email already in use')).toBeInTheDocument()
  })

  it('changes the password with current verification', async () => {
    changePasswordMock.mockReset().mockResolvedValue(undefined)
    render(<SettingsPage />)
    fireEvent.change(await screen.findByLabelText('Current password'), { target: { value: 'Engine!n1' } })
    fireEvent.change(screen.getByLabelText('New password'), { target: { value: 'N3w!passw' } })
    fireEvent.click(screen.getByRole('button', { name: 'Change password' }))
    await waitFor(() => expect(changePasswordMock).toHaveBeenCalledWith('Engine!n1', 'N3w!passw'))
    expect(await screen.findByText('Password changed.')).toBeInTheDocument()
  })

  it('renders one gallery button per curated style', async () => {
    render(<SettingsPage />)
    await screen.findByText('Pick a character')
    for (const style of DICEBEAR_STYLES) {
      expect(screen.getByRole('button', { name: `Pick ${style} character` })).toBeInTheDocument()
    }
  })

  it('Surprise me re-seeds the whole grid to the new character', async () => {
    render(<SettingsPage />)
    await screen.findByText('Pick a character')
    const before = screen.getByRole('button', { name: 'Pick bottts character' }).querySelector('img')!.getAttribute('src')!
    expect(before).toContain('seed=s1')
    fireEvent.click(screen.getByRole('button', { name: /surprise me/i }))
    await screen.findByText(/not saved yet/)
    const after = screen.getByRole('button', { name: 'Pick bottts character' }).querySelector('img')!.getAttribute('src')!
    expect(after).not.toBe(before)
    // Every tile shows the same new character in a different outfit.
    const seed = new URL(after).searchParams.get('seed')!
    expect(seed.length).toBeGreaterThan(0)
    for (const style of ['pixel-art', 'lorelei', 'micah']) {
      const src = screen.getByRole('button', { name: `Pick ${style} character` }).querySelector('img')!.getAttribute('src')!
      expect(new URL(src).searchParams.get('seed')).toBe(seed)
    }
  })

  it('picking an outfit after Surprise keeps the new character', async () => {
    render(<SettingsPage />)
    fireEvent.click(await screen.findByRole('button', { name: /surprise me/i }))
    await screen.findByText(/not saved yet/)
    const surpriseSeed = new URL(
      screen.getByRole('button', { name: 'Pick bottts character' }).querySelector('img')!.getAttribute('src')!,
    ).searchParams.get('seed')!
    fireEvent.click(screen.getByRole('button', { name: 'Pick pixel-art character' }))
    fireEvent.click(screen.getByRole('button', { name: 'Save image' }))
    await waitFor(() => expect(updateSettingsMock).toHaveBeenCalled())
    expect(updateSettingsMock.mock.calls[0][0].avatar_dicebear).toEqual({ style: 'pixel-art', seed: surpriseSeed })
  })

  it('Surprise me stages a preview without saving, Save image persists it', async () => {
    getSettingsMock.mockResolvedValue({ avatar_preset: 3 })
    render(<SettingsPage />)
    const surprise = await screen.findByRole('button', { name: /surprise me/i })
    expect(surprise.textContent).not.toMatch(/🎲/)
    expect(surprise.querySelector('svg')).not.toBeNull()
    fireEvent.click(surprise)
    expect(await screen.findByText(/not saved yet/)).toBeInTheDocument()
    expect(updateSettingsMock).not.toHaveBeenCalled()
    const save = screen.getByRole('button', { name: 'Save image' })
    expect(save).toBeEnabled()
    fireEvent.click(save)
    await waitFor(() => expect(updateSettingsMock).toHaveBeenCalled())
    const saved = updateSettingsMock.mock.calls[0][0]
    expect(saved.avatar_dicebear.style).toMatch(/^[a-z0-9-]+$/)
    expect(saved.avatar_dicebear.seed.length).toBeGreaterThan(0)
    expect(saved).not.toHaveProperty('avatar_preset')
    expect(await screen.findByText('Image saved.')).toBeInTheDocument()
  })

  it('gallery pick stages without saving until Save image', async () => {
    render(<SettingsPage />)
    fireEvent.click(await screen.findByRole('button', { name: 'Pick bottts character' }))
    expect(await screen.findByText(/not saved yet/)).toBeInTheDocument()
    expect(updateSettingsMock).not.toHaveBeenCalled()
    fireEvent.click(screen.getByRole('button', { name: 'Save image' }))
    await waitFor(() => expect(updateSettingsMock).toHaveBeenCalled())
    expect(updateSettingsMock.mock.calls[0][0].avatar_dicebear.style).toBe('bottts')
  })

  it('Save image is disabled with nothing staged', async () => {
    render(<SettingsPage />)
    await screen.findByText('Pick a character')
    expect(screen.getByRole('button', { name: 'Save image' })).toBeDisabled()
  })

  it('failed save keeps the staged pick and shows an error', async () => {
    updateSettingsMock.mockRejectedValue(new Error('offline'))
    render(<SettingsPage />)
    fireEvent.click(await screen.findByRole('button', { name: 'Pick bottts character' }))
    fireEvent.click(screen.getByRole('button', { name: 'Save image' }))
    expect(await screen.findByText(/check your connection/)).toBeInTheDocument()
    expect(screen.getByText(/not saved yet/)).toBeInTheDocument()
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

  it('lists connected accounts and disconnects', async () => {
    getIdentitiesMock.mockResolvedValue([{ provider: 'github', email: 'ada@example.com' }])
    render(<SettingsPage />)
    expect(await screen.findByText(/GitHub/)).toBeInTheDocument()
    fireEvent.click(screen.getByText('Disconnect'))
    await waitFor(() => expect(deleteIdentityMock).toHaveBeenCalledWith('github'))
  })

  it('connect starts a link-token OAuth dance', async () => {
    render(<SettingsPage />)
    fireEvent.click(await screen.findByRole('button', { name: 'Connect GitHub' }))
    await waitFor(() => expect(createLinkTokenMock).toHaveBeenCalled())
    expect(startOAuthLoginMock).toHaveBeenCalledWith('github', { intent: 'link', linkToken: 'linktok' })
  })

  it('prompts email verification for unverified addresses', async () => {
    mockUser = { student_id: 's1', name: 'Ada', username: 'ada', email: 'ada@example.com', email_verified: false } as never
    render(<SettingsPage />)
    fireEvent.click(await screen.findByRole('button', { name: 'Verify email' }))
    await waitFor(() => expect(requestEmailVerificationMock).toHaveBeenCalled())
    expect(await screen.findByText(/check your inbox/)).toBeInTheDocument()
  })
})
