import axios from "axios"

// Gateway auth: GET /auth/challenge issues a single-use EIP-4361 message for
// the address; the wallet personal_signs the message verbatim and POST
// /auth/verify exchanges { address, signature } for a JWT.

export const getAuthChallenge = async (address: string, chainId: number) => {
  return await axios.get<{ message: string; expiresAt: string }>(
    `/auth/challenge`,
    { params: { address, chainId } }
  )
}

type AuthTokens = {
  token: string
  expiresAt: string
  refreshToken: string
  refreshExpiresAt: string
}

export const verifyAuth = async (address: string, signature: string) => {
  return await axios.post<AuthTokens>(
    `/auth/verify`,
    { address, signature },
    { headers: { "Content-Type": "application/json" } }
  )
}

export const refreshAuth = async (refreshToken: string) => {
  return await axios.post<AuthTokens>(
    `/auth/refresh`,
    { refreshToken },
    { timeout: 10_000 }
  )
}

export const logoutAuth = async (refreshToken: string) => {
  return await axios.post(`/auth/logout`, { refreshToken })
}
