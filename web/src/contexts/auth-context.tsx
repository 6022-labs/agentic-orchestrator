/* eslint-disable react-refresh/only-export-components */
import * as React from "react"
import { Loader2 } from "lucide-react"
import { useLocation, useNavigate } from "react-router-dom"

import { Button } from "@/components/ui/button"
import { useSessionRefresh } from "@/hooks/use-session-refresh"
import {
  clearStoredSession,
  decodeJwt,
  getSession,
  isSessionStorageEvent,
  saveSession,
  setAccessToken,
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

  const [initialSession] = React.useState(getSession)
  const [jwt, setJwt] = React.useState<string | null>(null)
  const user = React.useMemo(() => decodeJwt(jwt), [jwt])
  // A stored session has no access token yet: hold the app until the first refresh settles.
  const [restoring, setRestoring] = React.useState(!!initialSession)
  const [sessionAddress, setSessionAddress] = React.useState(
    initialSession?.address
  )
  const [storedAddress, setStoredAddress] = React.useState(
    initialSession?.address
  )
  const [endedElsewhere, setEndedElsewhere] = React.useState(false)
  // Mirrors sessionAddress synchronously so a storage event never compares against a stale owner.
  const owner = React.useRef(initialSession?.address)
  const clockSkewMs = React.useRef(0)

  const setOwner = React.useCallback((address: string | undefined) => {
    owner.current = address
    setSessionAddress(address)
  }, [])

  const setToken = React.useCallback((token: string | null) => {
    setAccessToken(token)
    setJwt(token)
  }, [])

  const signIn = React.useCallback(
    (token: string, refreshToken: string) => {
      const decoded = decodeJwt(token)
      if (!decoded) return
      clockSkewMs.current = decoded.iat ? Date.now() - decoded.iat * 1000 : 0
      saveSession({ address: decoded.address, refreshToken })
      setOwner(decoded.address)
      setStoredAddress(decoded.address)
      setEndedElsewhere(false)
      setToken(token)
    },
    [setOwner, setToken]
  )

  const keepAccessToken = React.useCallback(
    (token: string) => {
      const issuedAt = decodeJwt(token)?.iat
      clockSkewMs.current = issuedAt ? Date.now() - issuedAt * 1000 : 0
      setToken(token)
    },
    [setToken]
  )

  const clearSession = React.useCallback(() => {
    clearStoredSession()
    setOwner(undefined)
    setStoredAddress(undefined)
    setRestoring(false)
    setToken(null)
  }, [setOwner, setToken])

  const endSessionHere = React.useCallback(() => {
    setOwner(undefined)
    setEndedElsewhere(true)
    setRestoring(false)
    setToken(null)
  }, [setOwner, setToken])

  const signOut = React.useCallback(() => {
    const session = getSession()
    // Storage may already hold another wallet's session whose event hasn't arrived: never revoke it from here.
    if (session?.address !== owner.current) {
      endSessionHere()
      return
    }
    clearSession()
    if (session) logoutAuth(session.refreshToken).catch(() => undefined)
  }, [clearSession, endSessionHere])

  const endRestoring = React.useCallback(() => setRestoring(false), [])

  const adoptStoredSession = React.useCallback(
    (address: string) => {
      const session = getSession()
      if (session?.address.toLowerCase() !== address.toLowerCase()) return false
      setOwner(session.address)
      setRestoring(true)
      return true
    },
    [setOwner]
  )

  useSessionRefresh({
    jwt,
    address: sessionAddress,
    restoring,
    clockSkewMs,
    onRefreshed: signIn,
    onAccessToken: keepAccessToken,
    onEndedElsewhere: endSessionHere,
    onExpired: clearSession,
    onSettled: endRestoring,
  })

  // Another tab signed out or signed in another wallet: end this tab's session without revoking it. Same-address rotations need nothing: the next refresh reads the stored token.
  React.useEffect(() => {
    const onStorage = (event: StorageEvent) => {
      if (!isSessionStorageEvent(event)) return
      const next = getSession()
      setStoredAddress(next?.address)
      if (!owner.current) return
      if (next?.address !== owner.current) {
        endSessionHere()
        return
      }
      setEndedElsewhere(false)
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
      {/* The public wallet-verification page doesn't need the session, so it isn't held behind the restore. */}
      {restoring && location.pathname !== "/verify-wallet" ? (
        <div className="flex h-svh flex-col items-center justify-center gap-3">
          <Loader2 className="size-6 animate-spin text-muted-foreground" />
          <p className="text-sm text-muted-foreground">
            Restoring your session…
          </p>
          {/* Retries continue while the gateway is down; this is the way out. */}
          <Button variant="ghost" size="sm" onClick={signOut}>
            Sign out
          </Button>
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
