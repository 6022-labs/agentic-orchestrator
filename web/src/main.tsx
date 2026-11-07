import "./index.css"
import App from "./App.tsx"
import { StrictMode } from "react"
import { createRoot } from "react-dom/client"
import { BrowserRouter } from "react-router-dom"
import { ThemeProvider } from "next-themes"
import { TooltipProvider } from "@/components/ui/tooltip.tsx"
import { Toaster } from "@/components/ui/sonner.tsx"
import ConfigGate from "@/components/config-gate.tsx"
import { initAxiosDefaults } from "@/configs/axios.ts"
import { AuthProvider } from "@/contexts/auth-context.tsx"
import { WalletProvider } from "@/contexts/wallet-provider-context.tsx"

initAxiosDefaults()

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <ThemeProvider
      attribute="class"
      defaultTheme="system"
      enableSystem
      enableColorScheme
      disableTransitionOnChange
      storageKey="agentic-theme"
    >
      <TooltipProvider>
        <ConfigGate>
          <BrowserRouter>
            <Toaster />
            <AuthProvider>
              <WalletProvider>
                <App />
              </WalletProvider>
            </AuthProvider>
          </BrowserRouter>
        </ConfigGate>
      </TooltipProvider>
    </ThemeProvider>
  </StrictMode>
)
