const SESSION_KEY = "session"

// The gateway JWT carries { sub, iat, nbf, exp }. There is no chain claim — the chain is fixed per deployment (see config.chainId).
export type UserData = {
  address: string
  iat?: number
  exp?: number
}

type StoredSession = {
  jwt: string
  refreshToken: string
}

// getSession reads both tokens at once: one key means no tab ever sees half a sign-in.
export const getSession = (): StoredSession | null => {
  try {
    const raw = localStorage.getItem(SESSION_KEY)
    return raw ? (JSON.parse(raw) as StoredSession) : null
  } catch {
    return null
  }
}

export const getJwt = () => getSession()?.jwt ?? null

export const saveTokens = (jwt: string, refreshToken: string) => {
  localStorage.setItem(SESSION_KEY, JSON.stringify({ jwt, refreshToken }))
}

export const clearTokens = () => {
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
