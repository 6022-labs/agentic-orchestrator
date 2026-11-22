const SESSION_KEY = "session"

// The gateway JWT carries { sub, iat, nbf, exp, sid }. There is no chain claim — the chain is fixed per deployment (see config.chainId).
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

// Fallback for blocked storage, so the refresh loop can still read back the token it was given.
let memorySession: StoredSession | null = null
let memoryOnly = false

export const getSession = (): StoredSession | null => {
  if (memoryOnly) return memorySession
  let raw: string | null
  try {
    raw = localStorage.getItem(SESSION_KEY)
  } catch {
    return memorySession
  }
  if (raw === null) return null
  try {
    const session = JSON.parse(raw) as Partial<StoredSession> | null
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
  memorySession = session
  try {
    localStorage.setItem(SESSION_KEY, JSON.stringify(session))
    memoryOnly = false
  } catch {
    // Blocked storage only loses persistence; the in-memory session still works.
    memoryOnly = true
    // A full quota can still hold the previous, now spent, refresh token.
    try {
      localStorage.removeItem(SESSION_KEY)
    } catch {
      // Unreadable storage holds nothing to clear.
    }
  }
}

export const clearStoredSession = () => {
  memorySession = null
  memoryOnly = false
  try {
    localStorage.removeItem(SESSION_KEY)
  } catch {
    // Blocked storage: nothing was persisted to clear.
  }
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
