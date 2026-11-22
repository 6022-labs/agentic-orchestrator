import axios, { type InternalAxiosRequestConfig } from "axios"

import { requestFreshAccessToken } from "@/lib/fresh-access-token"
import { decodeJwt, getAccessToken } from "@/lib/session-tokens"

import { config as envConfig } from "./env"

// Codes the API returns for an access token it no longer accepts; other 401s are not about the session.
const SESSION_REJECTED_CODES = new Set(["unauthorized"])

export const initAxiosDefaults = () => {
  // Point at the orchestrator gateway through the env config so the /api
  // fallback applies when no .env is present.
  // Note: no blanket Content-Type default — axios picks application/json for
  // plain objects and multipart/form-data for FormData bodies; forcing json
  // here would make axios serialize FormData (drafts, image upload) to JSON.
  axios.defaults.baseURL = envConfig.gatewayApiBaseUrl

  axios.interceptors.request.use((config) => {
    const jwt = getAccessToken()
    if (jwt && !config.headers?.Authorization) {
      config.headers = config.headers ?? {}
      config.headers.Authorization = `Bearer ${jwt}`
    }
    return config
  })

  // A 401 on an authenticated call refreshes once and retries; a revoked session signs out through that refresh.
  axios.interceptors.response.use(undefined, async (error) => {
    if (!axios.isAxiosError(error) || error.response?.status !== 401)
      throw error
    const code = (error.response.data as { code?: unknown } | undefined)?.code
    const request = error.config as
      (InternalAxiosRequestConfig & { sessionRetried?: boolean }) | undefined
    const sent = request?.headers?.Authorization
    // /auth/* includes the refresh call itself, which waiting on would deadlock.
    if (
      typeof code !== "string" ||
      !SESSION_REJECTED_CODES.has(code) ||
      !request ||
      request.sessionRetried ||
      typeof sent !== "string" ||
      request.url?.startsWith("/auth/")
    )
      throw error
    request.sessionRetried = true
    const current = getAccessToken()
    const token =
      current && `Bearer ${current}` !== sent
        ? current
        : await requestFreshAccessToken()
    // Never replay a request under another wallet's identity.
    const sentAddress = decodeJwt(sent.slice("Bearer ".length))?.address
    if (
      !token ||
      decodeJwt(token)?.address?.toLowerCase() !== sentAddress?.toLowerCase()
    )
      throw error
    request.headers.Authorization = `Bearer ${token}`
    return axios(request)
  })
}
