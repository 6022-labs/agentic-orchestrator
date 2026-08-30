import axios from "axios"

import type { ApiValidationError } from "@/models/api/validation-error"

// The gateway returns 400 validation errors as { field, message } and other
// errors as { code, message }. These helpers narrow each shape for the UI.

export const extractValidationError = <FieldName extends string = string>(
  error: unknown
): ApiValidationError<FieldName> | undefined => {
  if (!axios.isAxiosError(error)) {
    return undefined
  }

  if (error.response?.status !== 400) {
    return undefined
  }

  const data = error.response.data as { field?: string; message?: string }

  if (typeof data?.field !== "string" || typeof data?.message !== "string") {
    return undefined
  }

  return { field: data.field as FieldName, message: data.message }
}

// Toasts are for humans: anything longer than this, or that looks like a raw
// payload dump (ethers errors embed whole JSON-RPC bodies), is replaced by the
// caller's fallback. Full diagnostics stay in the console.
const MAX_TOAST_LENGTH = 160

const sanitize = (candidate: string | undefined, fallback: string): string => {
  if (!candidate) return fallback
  const message = candidate.trim()
  if (
    message.length === 0 ||
    message.length > MAX_TOAST_LENGTH ||
    message.includes('{"') ||
    message.includes("\n")
  ) {
    return fallback
  }
  return message
}

type ChainErrorShape = {
  code?: string | number
  shortMessage?: string
  reason?: string
  revert?: { name?: string }
  parsedRevertName?: string
}

// Unified gateway/axios + ethers error to human-message mapping.
export const extractErrorMessage = (
  error: unknown,
  fallback = "Something went wrong"
): string => {
  if (axios.isAxiosError(error)) {
    const data = error.response?.data as
      { code?: string; message?: string; field?: string } | undefined
    if (data?.message) {
      return sanitize(data.message, fallback)
    }
    return sanitize(error.message, fallback)
  }

  const chainError = (error ?? {}) as ChainErrorShape
  switch (chainError.code) {
    case "ACTION_REJECTED":
      return "Transaction rejected in your wallet"
    case "INSUFFICIENT_FUNDS":
      return "Not enough funds to pay for gas"
    case "CALL_EXCEPTION": {
      const revertName = chainError.parsedRevertName ?? chainError.revert?.name
      return revertName
        ? `Contract reverted: ${revertName}`
        : "The transaction would revert on chain"
    }
    case "NETWORK_ERROR":
    case "TIMEOUT":
      return "Network error — could not reach the chain RPC"
  }

  if (chainError.parsedRevertName) {
    return `Contract reverted: ${chainError.parsedRevertName}`
  }

  const candidate =
    chainError.shortMessage ??
    chainError.reason ??
    (error instanceof Error ? error.message : undefined)
  return sanitize(candidate, fallback)
}
