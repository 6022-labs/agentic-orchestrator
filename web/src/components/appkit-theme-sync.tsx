import React from "react"
import { useTheme } from "next-themes"
import { useAppKitTheme } from "@reown/appkit/react"

// The Reown modal ignores the .dark class; mirror the resolved theme into it.
export default function AppKitThemeSync() {
  const { resolvedTheme } = useTheme()
  const { setThemeMode } = useAppKitTheme()

  React.useEffect(() => {
    if (resolvedTheme === "light" || resolvedTheme === "dark") {
      setThemeMode(resolvedTheme)
    }
  }, [resolvedTheme, setThemeMode])

  return null
}
