/* eslint-disable react-refresh/only-export-components */
import * as React from "react"
import { useLocation, useNavigate } from "react-router-dom"

// The gateway JWT carries { sub, iat, nbf, exp }. There is no chain claim — the
// chain is fixed per deployment (see config.chainId).
type UserData = {
  address: string
  exp?: number
}

type AuthContextType = {
  user: UserData | null
  jwt: string | null
  signIn: (token: string) => void
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
    return { address: payload.sub, exp: payload.exp }
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
  const expirationTimeout = React.useRef<number | null>(null)

  const clearExpirationTimeout = () => {
    if (expirationTimeout.current != null) {
      window.clearTimeout(expirationTimeout.current)
      expirationTimeout.current = null
    }
  }

  const signIn = React.useCallback((token: string) => {
    localStorage.setItem("jwt", token)
    setJwt(token)
  }, [])

  const signOut = React.useCallback(() => {
    clearExpirationTimeout()
    localStorage.removeItem("jwt")
    setJwt(null)
    setUser(null)
  }, [])

  React.useEffect(() => {
    clearExpirationTimeout()

    if (!jwt) {
      setUser(null)
      return undefined
    }

    const decoded = extractUserData(jwt)
    setUser(decoded)

    if (!decoded?.exp) {
      return undefined
    }

    const now = Date.now()
    const expiresAtMs = decoded.exp * 1000

    if (expiresAtMs <= now) {
      signOut()
      return undefined
    }

    expirationTimeout.current = window.setTimeout(signOut, expiresAtMs - now)
    return () => clearExpirationTimeout()
  }, [jwt, signOut])

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
