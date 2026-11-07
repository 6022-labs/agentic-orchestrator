/* eslint-disable react-refresh/only-export-components */
import * as React from "react"
import { Loader2 } from "lucide-react"
import { useLocation, useNavigate } from "react-router-dom"

import { useSessionRefresh } from "@/hooks/use-session-refresh"
import {
  clearTokens,
  decodeJwt,
  getJwt,
  getSession,
  isSessionStorageEvent,
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
  // storedAddress is the address of the session in storage, which may belong to another tab.
  storedAddress: string | undefined
  // adoptStoredSession joins the stored session if it belongs to address; returns whether it did.
  adoptStoredSession: (address: string) => boolean
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
    return !!exp && exp * 1000 <= Date.now() && !!getSession()
  })
  const clockSkewMs = React.useRef(0)
  // Kept in step with every setJwt so a storage event never compares against a stale address.
  const userAddress = React.useRef(user?.address)
  const [endedElsewhere, setEndedElsewhere] = React.useState(false)
  const [storedAddress, setStoredAddress] = React.useState(user?.address)

  const signIn = React.useCallback((token: string, refreshToken: string) => {
    const decoded = decodeJwt(token)
    clockSkewMs.current = decoded?.iat ? Date.now() - decoded.iat * 1000 : 0
    saveTokens(token, refreshToken)
    userAddress.current = decoded?.address
    setStoredAddress(decoded?.address)
    setEndedElsewhere(false)
    setJwt(token)
  }, [])

  const clearSession = React.useCallback(() => {
    clearTokens()
    userAddress.current = undefined
    setStoredAddress(undefined)
    setJwt(null)
  }, [])

  const endSessionHere = React.useCallback(() => {
    userAddress.current = undefined
    setEndedElsewhere(true)
    setRestoring(false)
    setJwt(null)
  }, [])

  const signOut = React.useCallback(() => {
    const session = getSession()
    // Storage may already hold another wallet's session whose event hasn't arrived: never revoke it from here.
    if (decodeJwt(session?.jwt)?.address !== userAddress.current) {
      endSessionHere()
      return
    }
    clearSession()
    if (session) logoutAuth(session.refreshToken).catch(() => undefined)
  }, [clearSession, endSessionHere])

  const endRestoring = React.useCallback(() => setRestoring(false), [])

  const adoptStoredSession = React.useCallback((address: string) => {
    const session = getSession()
    const decoded = decodeJwt(session?.jwt)
    if (!session || decoded?.address.toLowerCase() !== address.toLowerCase())
      return false
    userAddress.current = decoded.address
    setRestoring(!!decoded.exp && decoded.exp * 1000 <= Date.now())
    setJwt(session.jwt)
    return true
  }, [])

  useSessionRefresh({
    jwt,
    clockSkewMs,
    onRefreshed: signIn,
    onEndedElsewhere: endSessionHere,
    onExpired: clearSession,
    onSettled: endRestoring,
  })

  // Another tab rotated this session or signed out; another wallet's sign-in there ends this tab's session without revoking it.
  React.useEffect(() => {
    const onStorage = (event: StorageEvent) => {
      if (!isSessionStorageEvent(event)) return
      const next = getJwt()
      setStoredAddress(decodeJwt(next)?.address)
      if (!userAddress.current) return
      if (decodeJwt(next)?.address !== userAddress.current) {
        endSessionHere()
        return
      }
      setEndedElsewhere(false)
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
      value={{
        user,
        jwt,
        signIn,
        signOut,
        endedElsewhere,
        storedAddress,
        adoptStoredSession,
      }}
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
