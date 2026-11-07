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
  // endedElsewhere is true once another tab ended this tab's session, so the login page doesn't auto-prompt.
  endedElsewhere: boolean
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
  // Kept in step with every setJwt so a storage event never compares against a stale address.
  const userAddress = React.useRef(user?.address)
  const [endedElsewhere, setEndedElsewhere] = React.useState(false)

  const signIn = React.useCallback((token: string, refreshToken: string) => {
    const decoded = decodeJwt(token)
    clockSkewMs.current = decoded?.iat ? Date.now() - decoded.iat * 1000 : 0
    saveTokens(token, refreshToken)
    userAddress.current = decoded?.address
    setEndedElsewhere(false)
    setJwt(token)
  }, [])

  const clearSession = React.useCallback(() => {
    clearTokens()
    userAddress.current = undefined
    setJwt(null)
  }, [])

  const signOut = React.useCallback(() => {
    const refreshToken = getRefreshToken()
    clearSession()
    if (refreshToken) logoutAuth(refreshToken).catch(() => undefined)
  }, [clearSession])

  const endRestoring = React.useCallback(() => setRestoring(false), [])

  const endSessionHere = React.useCallback(() => {
    userAddress.current = undefined
    setEndedElsewhere(true)
    setRestoring(false)
    setJwt(null)
  }, [])

  useSessionRefresh({
    jwt,
    clockSkewMs,
    onRefreshed: signIn,
    onCleared: clearSession,
    onSuperseded: endSessionHere,
    onExpired: clearSession,
    onSettled: endRestoring,
  })

  // Another tab rotated this session or signed out; another wallet's sign-in there ends this tab's session without revoking it.
  React.useEffect(() => {
    const onStorage = (event: StorageEvent) => {
      if (!isJwtStorageEvent(event) || !userAddress.current) return
      const next = getJwt()
      if (decodeJwt(next)?.address !== userAddress.current) {
        endSessionHere()
        return
      }
      setRestoring(false)
      setJwt(next)
    }
    window.addEventListener("storage", onStorage)
    return () => window.removeEventListener("storage", onStorage)
  }, [endSessionHere])

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
      value={{ user, jwt, signIn, signOut, endedElsewhere }}
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
