import * as React from "react"
import axios from "axios"

import {
  registerRefreshNow,
  settleFreshAccessToken,
} from "@/lib/fresh-access-token"
import { decodeJwt, getSession } from "@/lib/session-tokens"
import { refreshAuth } from "@/queries/api/auth.queries"

const REFRESH_LEAD_MS = 60_000
// With the 10s request timeout, the first 3 attempts start within the gateway's 30s reuse grace.
const REFRESH_RETRY_MS = 2_000
const REFRESH_FAST_ATTEMPTS = 3
// Past the grace a retry either succeeds or gets a 401, so failing offline never signs out.
const REFRESH_BACKOFF_MS = 60_000

type SessionRefreshOptions = {
  jwt: string | null
  // address owns the session this tab holds or is restoring.
  address: string | undefined
  // restoring has no access token yet, so it refreshes right away.
  restoring: boolean
  // Client minus gateway clock, measured on fresh tokens, so a skewed clock can't loop refreshes.
  clockSkewMs: React.RefObject<number>
  onRefreshed: (token: string, refreshToken: string) => void
  // onAccessToken keeps a fresh access token in memory when another tab already stored the rotated refresh token.
  onAccessToken: (token: string) => void
  // onEndedElsewhere ends this tab's session only: storage is empty or holds another wallet's session.
  onEndedElsewhere: () => void
  onExpired: () => void
  onSettled: () => void
}

// useSessionRefresh trades the stored refresh token for a new access token shortly before the current one expires.
export const useSessionRefresh = ({
  jwt,
  address,
  restoring,
  clockSkewMs,
  onRefreshed,
  onAccessToken,
  onEndedElsewhere,
  onExpired,
  onSettled,
}: SessionRefreshOptions) => {
  React.useEffect(() => {
    const exp = decodeJwt(jwt)?.exp
    if (!address || (!exp && !restoring)) return undefined

    const iat = decodeJwt(jwt)?.iat
    const expiresAtMs = exp ? exp * 1000 : 0
    const lifetimeMs = exp && iat ? expiresAtMs - iat * 1000 : 0
    const leadMs =
      lifetimeMs > 0
        ? Math.min(REFRESH_LEAD_MS, lifetimeMs / 2)
        : REFRESH_LEAD_MS
    const dueAtMs = exp ? expiresAtMs - leadMs : 0
    const gatewayNow = () => Date.now() - clockSkewMs.current
    // Set on cleanup so a refresh still in flight can't revive a signed-out session.
    let cancelled = false
    let inFlight = false
    let attempts = 0
    let timeout: number | undefined

    const refresh = async () => {
      attempts += 1
      inFlight = true
      let settled = true
      let fresh: string | null = null
      // Read at call time: another tab may have rotated it; never spend a token stored for another address.
      const session = getSession()
      try {
        if (!session || session.address !== address) {
          onEndedElsewhere()
          return
        }
        const { data } = await refreshAuth(session.refreshToken)
        const storedRefreshToken = getSession()?.refreshToken
        // Keep the successor while the spent token is still stored, even if cancelled: dropping it would leave a spent token behind.
        if (storedRefreshToken === session.refreshToken) {
          onRefreshed(data.token, data.refreshToken)
          fresh = data.token
        } else if (!storedRefreshToken) {
          if (!cancelled) onEndedElsewhere()
        } else if (getSession()?.address === address && !cancelled) {
          onAccessToken(data.token)
          fresh = data.token
        }
      } catch (error) {
        if (cancelled) return
        if (axios.isAxiosError(error) && error.response?.status === 401) {
          // Only expire what was spent: storage may already hold a newer session.
          if (getSession()?.refreshToken === session?.refreshToken) onExpired()
          else onEndedElsewhere()
          return
        }
        // Without an access token there is nothing to show yet: keep restoring while retrying.
        if (!jwt) settled = false
        timeout = window.setTimeout(
          refresh,
          attempts < REFRESH_FAST_ATTEMPTS
            ? REFRESH_RETRY_MS
            : REFRESH_BACKOFF_MS
        )
      } finally {
        inFlight = false
        settleFreshAccessToken(fresh)
        if (settled && !cancelled) onSettled()
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
    const unregisterRefreshNow = registerRefreshNow(() => {
      if (inFlight) return
      // A failed attempt already has its retry scheduled; spending more would push it past the reuse grace.
      if (attempts > 0) {
        settleFreshAccessToken(null)
        return
      }
      window.clearTimeout(timeout)
      refresh()
    })
    return () => {
      unregisterRefreshNow()
      cancelled = true
      window.clearTimeout(timeout)
      document.removeEventListener("visibilitychange", onVisibilityChange)
    }
  }, [
    jwt,
    address,
    restoring,
    clockSkewMs,
    onRefreshed,
    onAccessToken,
    onEndedElsewhere,
    onExpired,
    onSettled,
  ])
}
