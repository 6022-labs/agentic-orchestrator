const JWT_KEY = "jwt"
const REFRESH_TOKEN_KEY = "refresh_token"

// The gateway JWT carries { sub, iat, nbf, exp }. There is no chain claim — the chain is fixed per deployment (see config.chainId).
export type UserData = {
  address: string
  iat?: number
  exp?: number
}

export const getJwt = () => localStorage.getItem(JWT_KEY)

export const getRefreshToken = () => localStorage.getItem(REFRESH_TOKEN_KEY)

export const saveTokens = (jwt: string, refreshToken: string) => {
  localStorage.setItem(JWT_KEY, jwt)
  localStorage.setItem(REFRESH_TOKEN_KEY, refreshToken)
}

export const clearTokens = () => {
  localStorage.removeItem(JWT_KEY)
  localStorage.removeItem(REFRESH_TOKEN_KEY)
}

// A null key means another tab cleared the whole storage.
export const isJwtStorageEvent = (event: StorageEvent) =>
  event.key === JWT_KEY || event.key === null

export const decodeJwt = (token: string | null): UserData | null => {
  if (!token) return null

  try {
    const [, body] = token.split(".")
    const payload = JSON.parse(atob(body))
    return { address: payload.sub, iat: payload.iat, exp: payload.exp }
  } catch {
    return null
  }
}
