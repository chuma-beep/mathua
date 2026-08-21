'use client'

import { useState, useEffect } from 'react'
import { isLoggedIn, getUserInfo, type UserInfo } from '../lib/auth'

export function useAuthState(): { loggedIn: boolean; user: UserInfo | null } {
  const [loggedIn, setLoggedIn] = useState(false)
  const [user, setUser] = useState<UserInfo | null>(null)

  useEffect(() => {
    const sync = () => {
      setLoggedIn(isLoggedIn())
      setUser(getUserInfo())
    }
    sync()
    window.addEventListener('auth-changed', sync)
    window.addEventListener('storage', sync)
    return () => {
      window.removeEventListener('auth-changed', sync)
      window.removeEventListener('storage', sync)
    }
  }, [])

  return { loggedIn, user }
}
