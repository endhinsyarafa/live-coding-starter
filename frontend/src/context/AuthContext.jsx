import { createContext, useContext, useEffect, useMemo, useState } from 'react'
import { getMe, login as loginRequest, logout as logoutRequest } from '../services/authService'

const AuthContext = createContext(null)

export function AuthProvider({ children }) {
  const [token, setToken] = useState(() => localStorage.getItem('auth_token'))
  const [user, setUser] = useState(null)
  const [loading, setLoading] = useState(Boolean(token))
  const [authMessage, setAuthMessage] = useState('')

  useEffect(() => {
    if (!token) {
      setUser(null)
      setLoading(false)
      return
    }

    async function loadUser() {
      try {
        setLoading(true)
        const currentUser = await getMe()
        setUser(currentUser)
      } catch (err) {
        clearSession(err?.response?.data?.message || 'Sesi berakhir.')
      } finally {
        setLoading(false)
      }
    }

    loadUser()
  }, [token])

  useEffect(() => {
    function handleUnauthorized(event) {
      clearSession(event.detail?.message || 'Sesi berakhir.')
    }

    window.addEventListener('auth:unauthorized', handleUnauthorized)
    return () => window.removeEventListener('auth:unauthorized', handleUnauthorized)
  }, [])

  async function login(credentials) {
    const result = await loginRequest(credentials)
    localStorage.setItem('auth_token', result.token)
    setToken(result.token)
    setUser(result.user)
    setAuthMessage('')
    return result.user
  }

  async function logout() {
    try {
      if (token) {
        await logoutRequest()
      }
    } finally {
      clearSession('')
    }
  }

  function clearSession(message) {
    localStorage.removeItem('auth_token')
    setToken(null)
    setUser(null)
    setAuthMessage(message)
  }

  const value = useMemo(
    () => ({
      authMessage,
      isAuthenticated: Boolean(token && user),
      isAdmin: user?.role === 'admin',
      loading,
      login,
      logout,
      token,
      user,
    }),
    [authMessage, loading, token, user],
  )

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (!context) {
    throw new Error('useAuth must be used inside AuthProvider')
  }

  return context
}
