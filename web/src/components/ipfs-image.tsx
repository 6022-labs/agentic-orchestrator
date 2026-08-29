import React from "react"

import { ipfsImageUrls } from "@/lib/ipfs"

// A rate-limited gateway stalls instead of erroring, so onError never fires.
const CANDIDATE_TIMEOUT_MS = 8000

// IpfsImage walks the configured gateways until one loads.
export default function IpfsImage({
  value,
  alt,
  className,
  fallback = null,
}: {
  value?: string | null
  alt: string
  className?: string
  fallback?: React.ReactNode
}) {
  const candidates = React.useMemo(() => ipfsImageUrls(value), [value])
  const [renderedValue, setRenderedValue] = React.useState(value)
  const [index, setIndex] = React.useState(0)
  const [loaded, setLoaded] = React.useState(false)

  // Reset during render: an effect would pair new candidates with the old index.
  if (renderedValue !== value) {
    setRenderedValue(value)
    setIndex(0)
    setLoaded(false)
  }

  const next = React.useCallback(() => {
    setIndex((current) => current + 1)
    setLoaded(false)
  }, [])

  React.useEffect(() => {
    if (loaded || index >= candidates.length) return
    const timer = window.setTimeout(next, CANDIDATE_TIMEOUT_MS)
    return () => window.clearTimeout(timer)
  }, [candidates.length, index, loaded, next])

  const src = candidates[index]
  if (!src) return <>{fallback}</>

  return (
    <img
      key={src}
      src={src}
      alt={alt}
      className={className}
      onLoad={() => setLoaded(true)}
      onError={next}
    />
  )
}
