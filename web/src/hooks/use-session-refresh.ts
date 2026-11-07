import * as React from "react"
import axios from "axios"

import { decodeJwt, getJwt, getRefreshToken } from "@/lib/session-tokens"
import { refreshAuth } from "@/queries/api/auth.queries"

const REFRESH_LEAD_MS = 60_000
// With the 10s request timeout, the first 3 attempts start within the gateway's 30s reuse grace.
const REFRESH_RETRY_MS = 2_000
const REFRESH_FAST_ATTEMPTS = 3
// Past the grace a retry either succeeds or gets a 401, so failing offline never signs out.
const REFRESH_BACKOFF_MS = 60_000

type SessionRefreshOptions = {
  jwt: string | null
  // Client minus gateway clock, measured on fresh tokens, so a skewed clock can't loop refreshes.
  clockSkewMs: React.RefObject<number>
  onRefreshed: (token: string, refreshToken: string) => void
  onCleared: () => void
  // onSuperseded ends this tab's session only: storage holds another wallet's session.
  onSuperseded: () => void
  onExpired: () => void
  onSettled: () => void
}

// useSessionRefresh rotates the session shortly before the access token expires.
export const useSessionRefresh = ({
  jwt,
  clockSkewMs,
  onRefreshed,
  onCleared,
  onSuperseded,
  onExpired,
  onSettled,
}: SessionRefreshOptions) => {
  React.useEffect(() => {
    const decoded = decodeJwt(jwt)
    if (!decoded?.exp) return undefined

    const expiresAtMs = decoded.exp * 1000
    const lifetimeMs = decoded.iat ? expiresAtMs - decoded.iat * 1000 : 0
    const leadMs =
      lifetimeMs > 0
        ? Math.min(REFRESH_LEAD_MS, lifetimeMs / 2)
        : REFRESH_LEAD_MS
    const dueAtMs = expiresAtMs - leadMs
    const gatewayNow = () => Date.now() - clockSkewMs.current
    // Set on cleanup so a refresh still in flight can't revive a signed-out session.
    let cancelled = false
    let inFlight = false
    let attempts = 0
    let timeout: number | undefined

    const refresh = async () => {
      attempts += 1
      inFlight = true
      try {
        // Read at call time: another tab may have rotated it; never spend a token stored for another address.
        const refreshToken = getRefreshToken()
        const storedAddress = decodeJwt(getJwt())?.address
        if (!refreshToken || !storedAddress) {
          onCleared()
          return
        }
        if (storedAddress !== decoded.address) {
          onSuperseded()
          return
        }
        const { data } = await refreshAuth(refreshToken)
        const storedRefreshToken = getRefreshToken()
        // Keep the successor while the spent token is still stored, even if cancelled: dropping it would leave a spent token behind.
        if (storedRefreshToken === refreshToken)
          onRefreshed(data.token, data.refreshToken)
        // Storage emptied mid-flight (sign-out, cleared site data): end quietly, nothing expired.
        else if (!storedRefreshToken && !cancelled) onCleared()
      } catch (error) {
        if (cancelled) return
        if (axios.isAxiosError(error) && error.response?.status === 401) {
          onExpired()
          return
        }
        timeout = window.setTimeout(
          refresh,
          attempts < REFRESH_FAST_ATTEMPTS
            ? REFRESH_RETRY_MS
            : REFRESH_BACKOFF_MS
        )
      } finally {
        inFlight = false
        if (!cancelled) onSettled()
      }
    }

    // Hidden tabs and sleeping laptops delay timers; refresh as soon as the page is back if it is due.
    const onVisibilityChange = () => {
      if (document.visibilityState !== "visible" || inFlight) return
      if (gatewayNow() < dueAtMs) return
      window.clearTimeout(timeout)
      refresh()
    }

    timeout = window.setTimeout(refresh, Math.max(dueAtMs - gatewayNow(), 0))
    document.addEventListener("visibilitychange", onVisibilityChange)
    return () => {
      cancelled = true
      window.clearTimeout(timeout)
      document.removeEventListener("visibilitychange", onVisibilityChange)
    }
  }, [
    jwt,
    clockSkewMs,
    onRefreshed,
    onCleared,
    onSuperseded,
    onExpired,
    onSettled,
  ])
}
