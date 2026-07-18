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

export const verifyAuth = async (address: string, signature: string) => {
  return await axios.post<{ token: string; expiresAt: string }>(
    `/auth/verify`,
    { address, signature },
    { headers: { "Content-Type": "application/json" } }
  )
}
