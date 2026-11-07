/* eslint-disable react-refresh/only-export-components */
import * as React from "react"
import { Loader2 } from "lucide-react"
import { useLocation, useNavigate } from "react-router-dom"

import { useSessionRefresh } from "@/hooks/use-session-refresh"
import {
  clearTokens,
  decodeJwt,
  getJwt,
  getRefreshToken,
  isJwtStorageEvent,
  saveTokens,
  type UserData,
} from "@/lib/session-tokens"
import { logoutAuth } from "@/queries/api/auth.queries"

type AuthContextType = {
  user: UserData | null
  jwt: string | null
  signIn: (token: string, refreshToken: string) => void
  signOut: () => void
  // releaseSession ends the session when the wallet no longer matches it; returns true when it signed out.
  releaseSession: () => boolean
}

const AuthContext = React.createContext<AuthContextType | null>(null)

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const navigate = useNavigate()
  const location = useLocation()

  const [jwt, setJwt] = React.useState<string | null>(getJwt)
  const user = React.useMemo(() => decodeJwt(jwt), [jwt])
  // An expired stored token would render the dashboard into a burst of 401s; hold it until the first refresh settles.
  const [restoring, setRestoring] = React.useState(() => {
    const exp = user?.exp
    return !!exp && exp * 1000 <= Date.now() && !!getRefreshToken()
  })
  const clockSkewMs = React.useRef(0)
  const userAddress = React.useRef(user?.address)
  // Address of a session another tab signed in, and of one this tab left: another wallet's session is never revoked from here.
  const foreignAddress = React.useRef<string | null>(null)
  const ignoredAddress = React.useRef<string | null>(null)

  React.useEffect(() => {
    userAddress.current = user?.address
  }, [user])

  const storeSession = React.useCallback(
    (token: string, refreshToken: string) => {
      const issuedAt = decodeJwt(token)?.iat
      clockSkewMs.current = issuedAt ? Date.now() - issuedAt * 1000 : 0
      saveTokens(token, refreshToken)
      setJwt(token)
    },
    []
  )

  const signIn = React.useCallback(
    (token: string, refreshToken: string) => {
      foreignAddress.current = null
      ignoredAddress.current = null
      storeSession(token, refreshToken)
    },
    [storeSession]
  )

  const clearSession = React.useCallback(() => {
    clearTokens()
    setJwt(null)
  }, [])

  const signOut = React.useCallback(() => {
    const refreshToken = getRefreshToken()
    clearSession()
    if (refreshToken) logoutAuth(refreshToken).catch(() => undefined)
  }, [clearSession])

  const releaseSession = React.useCallback(() => {
    if (user && user.address === foreignAddress.current) {
      ignoredAddress.current = user.address
      setJwt(null)
      setRestoring(false)
      return false
    }
    signOut()
    return true
  }, [user, signOut])

  const endRestoring = React.useCallback(() => setRestoring(false), [])

  useSessionRefresh({
    jwt,
    clockSkewMs,
    onRefreshed: storeSession,
    onExpired: clearSession,
    onSettled: endRestoring,
  })

  // Another tab signed in, rotated the token or signed out: follow it.
  React.useEffect(() => {
    const onStorage = (event: StorageEvent) => {
      if (!isJwtStorageEvent(event)) return
      const next = getJwt()
      const nextAddress = decodeJwt(next)?.address
      if (nextAddress && nextAddress === ignoredAddress.current) return
      if (nextAddress && nextAddress !== userAddress.current)
        foreignAddress.current = nextAddress
      setRestoring(false)
      setJwt(next)
    }
    window.addEventListener("storage", onStorage)
    return () => window.removeEventListener("storage", onStorage)
  }, [])

  React.useEffect(() => {
    if (restoring) return
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
  }, [jwt, location, restoring])

  return (
    <AuthContext.Provider
      value={{ user, jwt, signIn, signOut, releaseSession }}
    >
      {restoring ? (
        <div className="flex h-svh flex-col items-center justify-center gap-3">
          <Loader2 className="size-6 animate-spin text-muted-foreground" />
          <p className="text-sm text-muted-foreground">
            Restoring your session…
          </p>
        </div>
      ) : (
        children
      )}
    </AuthContext.Provider>
  )
}

export const useAuth = () => {
  const context = React.useContext(AuthContext)
  if (!context) throw new Error("useAuth must be used within an AuthProvider")
  return context
}
