import { twMerge } from "tailwind-merge";
import { clsx, type ClassValue } from "clsx";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function truncateEthAddress(
  address: string | undefined,
  startLength: number = 6,
  endLength: number = 4
): string {
  if (!address) {
    return "";
  }

  // Ensure the address is not shorter than the sum of the start and end lengths
  if (address.length <= startLength + endLength) {
    return address; // Return the original address if it's too short to truncate
  }

  return `${address.slice(0, startLength)}...${address.slice(-endLength)}`;
}

export function strToBytes(hex: string): Uint8Array {
  if (hex.length % 2 !== 0) {
    throw new Error("Invalid hex string");
  }
  const bytes = new Uint8Array(hex.length);
  for (let i = 0; i < hex.length; i += 1) {
    bytes[i] = hex.substring(i, i + 1).charCodeAt(0);
  }
  return bytes;
}

const amazonPlatformRegionPatternRegex = /[a-z]{2}(-[a-z]+)+-\d+/i;

export function extractAmazonPlatformRegionFromUrl(
  url: string
): string | undefined {
  const values = amazonPlatformRegionPatternRegex.exec(url);

  if (values?.length) {
    return values[0];
  }

  return undefined;
}
