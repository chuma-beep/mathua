import { describe, it, expect, beforeEach, beforeAll, vi } from 'vitest'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'

beforeAll(() => {
  Object.defineProperty(window, 'matchMedia', {
    writable: true,
    value: (query: string) => ({
      matches: false,
      media: query,
      onchange: null,
      addListener: () => {},
      removeListener: () => {},
      addEventListener: () => {},
      removeEventListener: () => {},
      dispatchEvent: () => false,
    }),
  })
})

const pushMock = vi.fn()
const signupMock = vi.fn()
const loginMock = vi.fn()
const validateTokenMock = vi.fn()
const getConfigMock = vi.fn()
const requestResetMock = vi.fn()
const completeResetMock = vi.fn()

vi.mock('next/navigation', () => ({
  useRouter: () => ({ push: pushMock }),
  usePathname: () => '/login',
  useSearchParams: () => new URLSearchParams(),
}))

vi.mock('../lib/api', () => ({
  signup: (...a: unknown[]) => signupMock(...a),
  login: (...a: unknown[]) => loginMock(...a),
  validateToken: (...a: unknown[]) => validateTokenMock(...a),
  requestPasswordReset: (...a: unknown[]) => requestResetMock(...a),
  completePasswordReset: (...a: unknown[]) => completeResetMock(...a),
  startOAuthLogin: vi.fn(),
  OAUTH_LABELS: { google: 'Google', github: 'GitHub', facebook: 'Facebook', microsoft: 'Microsoft', apple: 'Apple' },
  getConfig: (...a: unknown[]) => getConfigMock(...a),
  API_BASE: '',
}))

vi.mock('../components/Header', () => ({ default: () => <div data-testid="header-stub" /> }))
vi.mock('../components/Footer', () => ({ default: () => <div data-testid="footer-stub" /> }))

import LoginPage from '../app/login/page'

describe('LoginPage', () => {
  beforeEach(() => {
    pushMock.mockClear()
    signupMock.mockClear()
    loginMock.mockClear()
    validateTokenMock.mockClear()
    requestResetMock.mockReset().mockResolvedValue(undefined)
    completeResetMock.mockReset()
    getConfigMock.mockReset().mockResolvedValue({ auth_enabled: true })
  })

  it('requires username and password on login and does not call the API', () => {
    render(<LoginPage />)
    fireEvent.click(screen.getByTestId("auth-submit"))
    expect(screen.getByText('Username is required')).toBeInTheDocument()
    expect(screen.getByText('Password is required')).toBeInTheDocument()
    expect(loginMock).not.toHaveBeenCalled()
  })

  it('rejects a weak signup password with the policy message', () => {
    render(<LoginPage />)
    fireEvent.click(screen.getByRole('button', { name: 'Sign Up' }))
    fireEvent.change(screen.getByLabelText('Name'), { target: { value: 'Ada' } })
    fireEvent.change(screen.getByLabelText('Username'), { target: { value: 'ada' } })
    fireEvent.change(screen.getByLabelText('Email'), { target: { value: 'ada@example.com' } })
    fireEvent.change(screen.getByLabelText('Password'), { target: { value: 'abc' } })
    fireEvent.click(screen.getByRole('button', { name: 'Create Account' }))
    expect(
      screen.getByText('Password needs at least 8 characters, a number, a special character')
    ).toBeInTheDocument()
    expect(signupMock).not.toHaveBeenCalled()
  })

  it('clears field errors live after a failed attempt', () => {
    render(<LoginPage />)
    fireEvent.click(screen.getByRole('button', { name: 'Sign Up' }))
    fireEvent.change(screen.getByLabelText('Name'), { target: { value: 'Ada' } })
    fireEvent.change(screen.getByLabelText('Username'), { target: { value: 'ada' } })
    fireEvent.change(screen.getByLabelText('Email'), { target: { value: 'ada@example.com' } })
    fireEvent.change(screen.getByLabelText('Password'), { target: { value: 'abc' } })
    fireEvent.click(screen.getByRole('button', { name: 'Create Account' }))
    expect(screen.getByText('Password needs at least 8 characters, a number, a special character')).toBeInTheDocument()
    fireEvent.change(screen.getByLabelText('Password'), { target: { value: 'Engine!n1' } })
    expect(screen.queryByText(/Password needs/)).not.toBeInTheDocument()
  })

  it('submits a valid signup with trimmed values', () => {
    signupMock.mockResolvedValue({ token: 't', student_id: 's1', name: 'Ada', diagnostic_completed: true })
    validateTokenMock.mockResolvedValue({ valid: true })
    render(<LoginPage />)
    fireEvent.click(screen.getByRole('button', { name: 'Sign Up' }))
    fireEvent.change(screen.getByLabelText('Name'), { target: { value: '  Ada  ' } })
    fireEvent.change(screen.getByLabelText('Username'), { target: { value: '  ada  ' } })
    fireEvent.change(screen.getByLabelText('Email'), { target: { value: '  Ada@Example.com  ' } })
    fireEvent.change(screen.getByLabelText('Password'), { target: { value: 'Engine!n1' } })
    fireEvent.click(screen.getByRole('button', { name: 'Create Account' }))
    expect(signupMock).toHaveBeenCalledWith('Ada', 'ada', 'Engine!n1', 'Ada@Example.com')
  })

  it('requires an email on signup', () => {
    render(<LoginPage />)
    fireEvent.click(screen.getByRole('button', { name: 'Sign Up' }))
    fireEvent.change(screen.getByLabelText('Name'), { target: { value: 'Ada' } })
    fireEvent.change(screen.getByLabelText('Username'), { target: { value: 'ada' } })
    fireEvent.change(screen.getByLabelText('Password'), { target: { value: 'Engine!n1' } })
    fireEvent.click(screen.getByRole('button', { name: 'Create Account' }))
    expect(screen.getByText('Email is required')).toBeInTheDocument()
    expect(signupMock).not.toHaveBeenCalled()
  })

  it('rejects a malformed signup email', () => {
    render(<LoginPage />)
    fireEvent.click(screen.getByRole('button', { name: 'Sign Up' }))
    fireEvent.change(screen.getByLabelText('Name'), { target: { value: 'Ada' } })
    fireEvent.change(screen.getByLabelText('Username'), { target: { value: 'ada' } })
    fireEvent.change(screen.getByLabelText('Email'), { target: { value: 'not-an-email' } })
    fireEvent.change(screen.getByLabelText('Password'), { target: { value: 'Engine!n1' } })
    fireEvent.click(screen.getByRole('button', { name: 'Create Account' }))
    expect(screen.getByText('Enter a valid email address')).toBeInTheDocument()
    expect(signupMock).not.toHaveBeenCalled()
  })

  it('toggles password visibility with the eye button', () => {
    render(<LoginPage />)
    const pw = screen.getByLabelText('Password') as HTMLInputElement
    expect(pw.type).toBe('password')
    fireEvent.click(screen.getByRole('button', { name: 'Show password' }))
    expect(pw.type).toBe('text')
    fireEvent.click(screen.getByRole('button', { name: 'Hide password' }))
    expect(pw.type).toBe('password')
  })

  it('rejects passwords whose only special is a non-ASCII letter (server parity)', () => {
    render(<LoginPage />)
    fireEvent.click(screen.getByRole('button', { name: 'Sign Up' }))
    fireEvent.change(screen.getByLabelText('Name'), { target: { value: 'Ada' } })
    fireEvent.change(screen.getByLabelText('Username'), { target: { value: 'ada' } })
    fireEvent.change(screen.getByLabelText('Email'), { target: { value: 'ada@example.com' } })
    fireEvent.change(screen.getByLabelText('Password'), { target: { value: 'abcdé123' } })
    fireEvent.click(screen.getByRole('button', { name: 'Create Account' }))
    expect(screen.getByText(/a special character/)).toBeInTheDocument()
    expect(signupMock).not.toHaveBeenCalled()
  })

  it('rejects overlong passwords at the bcrypt boundary', () => {
    render(<LoginPage />)
    fireEvent.click(screen.getByRole('button', { name: 'Sign Up' }))
    fireEvent.change(screen.getByLabelText('Name'), { target: { value: 'Ada' } })
    fireEvent.change(screen.getByLabelText('Username'), { target: { value: 'ada' } })
    fireEvent.change(screen.getByLabelText('Email'), { target: { value: 'ada@example.com' } })
    fireEvent.change(screen.getByLabelText('Password'), { target: { value: `${'a1!'.repeat(25)}` } })
    fireEvent.click(screen.getByRole('button', { name: 'Create Account' }))
    expect(screen.getByText(/at most 72 characters/)).toBeInTheDocument()
    expect(signupMock).not.toHaveBeenCalled()
  })

  it('rejects invalid usernames on signup', () => {
    render(<LoginPage />)
    fireEvent.click(screen.getByRole('button', { name: 'Sign Up' }))
    fireEvent.change(screen.getByLabelText('Name'), { target: { value: 'Ada' } })
    fireEvent.change(screen.getByLabelText('Username'), { target: { value: 'ab' } })
    fireEvent.change(screen.getByLabelText('Email'), { target: { value: 'ada@example.com' } })
    fireEvent.change(screen.getByLabelText('Password'), { target: { value: 'Engine!n1' } })
    fireEvent.click(screen.getByRole('button', { name: 'Create Account' }))
    expect(screen.getByText('Username must be 3-20 characters')).toBeInTheDocument()
    expect(signupMock).not.toHaveBeenCalled()
    fireEvent.change(screen.getByLabelText('Username'), { target: { value: 'bad name!' } })
    fireEvent.click(screen.getByRole('button', { name: 'Create Account' }))
    expect(screen.getByText(/may only contain letters/)).toBeInTheDocument()
    expect(signupMock).not.toHaveBeenCalled()
  })

  it('routes to profile after a valid login', async () => {
    loginMock.mockResolvedValue({ token: 't', student_id: 's1', name: 'Ada', diagnostic_completed: true })
    validateTokenMock.mockResolvedValue({ valid: true })
    render(<LoginPage />)
    fireEvent.change(screen.getByLabelText('Username'), { target: { value: 'ada' } })
    fireEvent.change(screen.getByLabelText('Password'), { target: { value: 'Engine!n1' } })
    fireEvent.click(screen.getByTestId('auth-submit'))
    await waitFor(() => expect(pushMock).toHaveBeenCalledWith('/profile'))
  })

  it('shows a disabled notice when the server has auth disabled', async () => {
    getConfigMock.mockResolvedValue({ auth_enabled: false })
    render(<LoginPage />)
    expect(await screen.findByText(/Accounts are disabled on this server/)).toBeInTheDocument()
    expect(screen.getByTestId('auth-submit')).toBeDisabled()
  })

  it('forgot flow sends a reset request with anti-enumeration copy', async () => {
    render(<LoginPage />)
    fireEvent.click(await screen.findByText('Forgot password?'))
    fireEvent.change(screen.getByLabelText('Username or email'), { target: { value: 'ada' } })
    fireEvent.click(screen.getByText('Send reset link'))
    await waitFor(() => expect(requestResetMock).toHaveBeenCalledWith('ada'))
    expect(await screen.findByText(/reset link is on its way/)).toBeInTheDocument()
  })

  it('renders one button per configured non-google provider', async () => {
    getConfigMock.mockResolvedValue({ auth_enabled: true, providers: ['google', 'github', 'apple'] })
    render(<LoginPage />)
    expect(await screen.findByRole('button', { name: 'Continue with GitHub' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Continue with Apple' })).toBeInTheDocument()
    expect(screen.queryByRole('button', { name: 'Continue with Facebook' })).toBeNull()
  })
})
