import { describe, it, expect, beforeEach, beforeAll, vi } from 'vitest'
import { render, screen, fireEvent } from '@testing-library/react'

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

vi.mock('next/navigation', () => ({
  useRouter: () => ({ push: pushMock }),
  usePathname: () => '/login',
}))

vi.mock('../lib/api', () => ({
  signup: (...a: unknown[]) => signupMock(...a),
  login: (...a: unknown[]) => loginMock(...a),
  validateToken: (...a: unknown[]) => validateTokenMock(...a),
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
    fireEvent.change(screen.getByLabelText('Password'), { target: { value: 'Engine!n1' } })
    fireEvent.click(screen.getByRole('button', { name: 'Create Account' }))
    expect(signupMock).toHaveBeenCalledWith('Ada', 'ada', 'Engine!n1')
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
})
