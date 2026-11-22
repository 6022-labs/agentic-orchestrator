// The mounted session refresh registers here, so a 401 reuses its single refresh path instead of racing it.
type RefreshNow = () => void

let refreshNow: RefreshNow | undefined
let waiters: Array<(token: string | null) => void> = []

export const registerRefreshNow = (fn: RefreshNow) => {
  refreshNow = fn
  return () => {
    if (refreshNow === fn) refreshNow = undefined
  }
}

// requestFreshAccessToken resolves with the next access token, or null when the session could not be refreshed.
export const requestFreshAccessToken = (): Promise<string | null> => {
  if (!refreshNow) return Promise.resolve(null)
  const fresh = new Promise<string | null>((resolve) => waiters.push(resolve))
  refreshNow()
  return fresh
}

export const settleFreshAccessToken = (token: string | null) => {
  const settled = waiters
  waiters = []
  settled.forEach((resolve) => resolve(token))
}
