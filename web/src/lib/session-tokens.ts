const SESSION_KEY = "session"

// The gateway JWT carries { sub, iat, nbf, exp }. There is no chain claim — the chain is fixed per deployment (see config.chainId).
export type UserData = {
  address: string
  iat?: number
  exp?: number
}

// Only the refresh token and its owner are stored; the access JWT stays in memory (the refresh token is still readable by XSS).
type StoredSession = {
  address: string
  refreshToken: string
}

let accessToken: string | null = null

export const getAccessToken = () => accessToken

export const setAccessToken = (token: string | null) => {
  accessToken = token
}

export const getSession = (): StoredSession | null => {
  try {
    const raw = localStorage.getItem(SESSION_KEY)
    const session = raw ? (JSON.parse(raw) as Partial<StoredSession>) : null
    if (
      typeof session?.address !== "string" ||
      typeof session.refreshToken !== "string"
    )
      return null
    return { address: session.address, refreshToken: session.refreshToken }
  } catch {
    return null
  }
}

export const saveSession = (session: StoredSession) => {
  localStorage.setItem(SESSION_KEY, JSON.stringify(session))
}

export const clearStoredSession = () => {
  localStorage.removeItem(SESSION_KEY)
}

// A null key means another tab cleared the whole storage.
export const isSessionStorageEvent = (event: StorageEvent) =>
  event.key === SESSION_KEY || event.key === null

export const decodeJwt = (
  token: string | null | undefined
): UserData | null => {
  if (!token) return null

  try {
    const [, body] = token.split(".")
    const payload = JSON.parse(atob(body))
    return { address: payload.sub, iat: payload.iat, exp: payload.exp }
  } catch {
    return null
  }
}
