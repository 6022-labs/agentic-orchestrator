const relativeTime = new Intl.RelativeTimeFormat(undefined, {
  numeric: "auto",
})

const UNITS: { unit: Intl.RelativeTimeFormatUnit; seconds: number }[] = [
  { unit: "second", seconds: 1 },
  { unit: "minute", seconds: 60 },
  { unit: "hour", seconds: 3600 },
  { unit: "day", seconds: 86400 },
]

// "5 minutes ago" / "yesterday"; past ~7 days falls back to the locale date.
export function formatRelativeTime(iso: string): string {
  const date = new Date(iso)
  if (Number.isNaN(date.getTime())) return iso

  const elapsedSeconds = (date.getTime() - Date.now()) / 1000
  if (Math.abs(elapsedSeconds) > 7 * 86400) {
    return date.toLocaleDateString()
  }

  for (let i = UNITS.length - 1; i >= 0; i--) {
    const { unit, seconds } = UNITS[i]
    if (Math.abs(elapsedSeconds) >= seconds || unit === "second") {
      return relativeTime.format(Math.round(elapsedSeconds / seconds), unit)
    }
  }
  return relativeTime.format(0, "second")
}
