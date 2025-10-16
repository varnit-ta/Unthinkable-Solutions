/**
 * Authentication Module
 * 
 * Provides authentication context and hooks for managing user sessions.
 * Handles token persistence in localStorage and provides authentication state
 * throughout the application.
 * 
 * @module auth
 */

import { createContext, useContext, useEffect, useMemo, useState } from 'react'

/**
 * Authentication state structure
 */
type AuthState = {
  token: string | null
  setToken: (t: string | null) => void
}

/**
 * React context for authentication state
 */
const AuthContext = createContext<AuthState | undefined>(undefined)

/**
 * Authentication Provider Component
 * 
 * Wraps the application to provide authentication state and token management.
 * Automatically syncs authentication token with localStorage for persistence.
 * 
 * @param {Object} props - Component props
 * @param {React.ReactNode} props.children - Child components
 * @returns {JSX.Element} Provider component
 */
export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [token, setToken] = useState<string | null>(() => 
    localStorage.getItem('token')
  )

  useEffect(() => {
    if (token) {
      localStorage.setItem('token', token)
    } else {
      localStorage.removeItem('token')
    }
  }, [token])

  const value = useMemo(() => ({ token, setToken }), [token])

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  )
}

/**
 * Custom hook to access authentication state
 * 
 * Provides access to the current authentication token and setter function.
 * Must be used within an AuthProvider component.
 * 
 * @returns {AuthState} Authentication state object
 * @throws {Error} If used outside of AuthProvider
 * 
 * @example
 * const { token, setToken } = useAuth()
 * if (token) {
 *   // User is authenticated
 * }
 */
export function useAuth() {
  const ctx = useContext(AuthContext)

  if (!ctx) {
    throw new Error('useAuth must be used within AuthProvider')
  }

  return ctx
}


