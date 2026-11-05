/* eslint-disable react-refresh/only-export-components */
import * as React from "react"
import axios from "axios"
import { useLocation, useNavigate } from "react-router-dom"

import { logoutAuth, refreshAuth } from "@/queries/api/auth.queries"

const REFRESH_LEAD_MS = 60_000
// With the 5s request timeout, all 3 attempts start within the gateway's 30s reuse grace.
const REFRESH_RETRY_MS = 8_000
const REFRESH_MAX_ATTEMPTS = 3
const REFRESH_TOKEN_KEY = "refresh_token"

// The gateway JWT carries { sub, iat, nbf, exp }. There is no chain claim — the
// chain is fixed per deployment (see config.chainId).
type UserData = {
  address: string
  iat?: number
  exp?: number
}

type AuthContextType = {
  user: UserData | null
  jwt: string | null
  signIn: (token: string, refreshToken: string) => void
  signOut: () => void
}

const AuthContext = React.createContext<AuthContextType>({
  user: null,
  jwt: null,
  signIn: () => {},
  signOut: () => {},
})

const extractUserData = (token: string): UserData | null => {
  if (!token) return null

  try {
    const [, body] = token.split(".")
    const payload = JSON.parse(atob(body))
    return { address: payload.sub, iat: payload.iat, exp: payload.exp }
  } catch {
    return null
  }
}

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const navigate = useNavigate()
  const location = useLocation()

  const [user, setUser] = React.useState<UserData | null>(() =>
    extractUserData(localStorage.getItem("jwt") || "")
  )
  const [jwt, setJwt] = React.useState<string | null>(() =>
    localStorage.getItem("jwt")
  )
  const refreshTimeout = React.useRef<number | null>(null)
  // Client minus gateway clock, measured on fresh tokens, so a skewed clock can't loop refreshes.
  const clockSkewMs = React.useRef(0)

  const clearRefreshTimeout = () => {
    if (refreshTimeout.current != null) {
      window.clearTimeout(refreshTimeout.current)
      refreshTimeout.current = null
    }
  }

  const signIn = React.useCallback((token: string, refreshToken: string) => {
    const issuedAt = extractUserData(token)?.iat
    clockSkewMs.current = issuedAt ? Date.now() - issuedAt * 1000 : 0
    localStorage.setItem("jwt", token)
    localStorage.setItem(REFRESH_TOKEN_KEY, refreshToken)
    setJwt(token)
  }, [])

  const clearSession = React.useCallback(() => {
    clearRefreshTimeout()
    localStorage.removeItem("jwt")
    localStorage.removeItem(REFRESH_TOKEN_KEY)
    setJwt(null)
    setUser(null)
  }, [])

  const signOut = React.useCallback(() => {
    const refreshToken = localStorage.getItem(REFRESH_TOKEN_KEY)
    clearSession()
    if (refreshToken) logoutAuth(refreshToken).catch(() => undefined)
  }, [clearSession])

  React.useEffect(() => {
    clearRefreshTimeout()

    if (!jwt) {
      setUser(null)
      return undefined
    }

    const decoded = extractUserData(jwt)
    setUser(decoded)

    if (!decoded?.exp) {
      return undefined
    }

    const expiresAtMs = decoded.exp * 1000
    const lifetimeMs = decoded.iat ? expiresAtMs - decoded.iat * 1000 : 0
    const leadMs =
      lifetimeMs > 0
        ? Math.min(REFRESH_LEAD_MS, lifetimeMs / 2)
        : REFRESH_LEAD_MS
    const gatewayNow = () => Date.now() - clockSkewMs.current
    // Set on cleanup so a refresh still in flight can't revive a signed-out session.
    let cancelled = false
    let attempts = 0

    const refresh = async () => {
      attempts += 1
      // Read at call time: another tab may have rotated it.
      const refreshToken = localStorage.getItem(REFRESH_TOKEN_KEY)
      if (!refreshToken) {
        clearSession()
        return
      }
      try {
        const { data } = await refreshAuth(refreshToken)
        if (!cancelled) signIn(data.token, data.refreshToken)
      } catch (error) {
        if (cancelled) return
        const isUnauthorized =
          axios.isAxiosError(error) && error.response?.status === 401
        if (isUnauthorized || attempts >= REFRESH_MAX_ATTEMPTS) {
          clearSession()
          return
        }
        refreshTimeout.current = window.setTimeout(refresh, REFRESH_RETRY_MS)
      }
    }

    refreshTimeout.current = window.setTimeout(
      refresh,
      Math.max(expiresAtMs - leadMs - gatewayNow(), 0)
    )
    return () => {
      cancelled = true
      clearRefreshTimeout()
    }
  }, [jwt, signIn, clearSession])

  React.useEffect(() => {
    const isLoginRoute = location.pathname === "/login"
    // /verify-wallet is reachable without a session (bridge deep links land
    // there), so it is treated like /login for the unauthenticated redirect.
    const isPublicRoute = isLoginRoute || location.pathname === "/verify-wallet"

    if (jwt && isLoginRoute) {
      navigate("/dashboard")
    } else if (!jwt && !isPublicRoute) {
      navigate("/login")
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [jwt, location])

  return (
    <AuthContext.Provider value={{ user, jwt, signIn, signOut }}>
      {children}
    </AuthContext.Provider>
  )
}

export const useAuth = () => {
  const context = React.useContext(AuthContext)
  if (!context) throw new Error("useAuth must be used within an AuthProvider")
  return context
}
