import axios from "axios"

// Agent-node gateway auth: GET /auth/challenge issues a single-use nonce; the
// wallet personal_signs the raw nonce and POST /auth/verify exchanges
// { address, signature, nonce } for a JWT.

export const getAuthChallenge = async () => {
  return await axios.get<{ nonce: string; expiresAt: string }>(
    `/auth/challenge`
  )
}

export const verifyAuth = async (
  address: string,
  signature: string,
  nonce: string
) => {
  return await axios.post<{ token: string; expiresAt: string }>(
    `/auth/verify`,
    { address, signature, nonce },
    { headers: { "Content-Type": "application/json" } }
  )
}
