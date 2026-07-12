import { twMerge } from "tailwind-merge"
import { clsx, type ClassValue } from "clsx"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function truncateEthAddress(
  address: string | undefined,
  startLength: number = 6,
  endLength: number = 4
): string {
  if (!address) {
    return ""
  }

  if (address.length <= startLength + endLength) {
    return address
  }

  return `${address.slice(0, startLength)}...${address.slice(-endLength)}`
}

export function newId(): string {
  return typeof crypto !== "undefined" && "randomUUID" in crypto
    ? crypto.randomUUID()
    : Math.random().toString(36).slice(2)
}

const amazonPlatformRegionPatternRegex = /[a-z]{2}(-[a-z]+)+-\d+/i

// Pull the AWS region out of a Bedrock runtime endpoint URL, e.g.
// https://bedrock-runtime.us-east-1.amazonaws.com -> us-east-1.
export function extractAmazonPlatformRegionFromUrl(
  url: string
): string | undefined {
  const values = amazonPlatformRegionPatternRegex.exec(url)

  if (values?.length) {
    return values[0]
  }

  return undefined
}
