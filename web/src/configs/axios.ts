import axios from "axios"

import { getJwt } from "@/lib/session-tokens"

import { config as envConfig } from "./env"

export const initAxiosDefaults = () => {
  // Point at the orchestrator gateway through the env config so the /api
  // fallback applies when no .env is present.
  // Note: no blanket Content-Type default — axios picks application/json for
  // plain objects and multipart/form-data for FormData bodies; forcing json
  // here would make axios serialize FormData (drafts, image upload) to JSON.
  axios.defaults.baseURL = envConfig.gatewayApiBaseUrl

  axios.interceptors.request.use((config) => {
    const jwt = getJwt()
    if (jwt && !config.headers?.Authorization) {
      config.headers = config.headers ?? {}
      config.headers.Authorization = `Bearer ${jwt}`
    }
    return config
  })
}
