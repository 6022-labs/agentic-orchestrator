import React from "react"
import { Navigate, Route, Routes } from "react-router-dom"
import { createAppKit } from "@reown/appkit/react"
import { EthersAdapter } from "@reown/appkit-adapter-ethers"
import { type AppKitNetwork, polygonAmoy } from "@reown/appkit/networks"

import { config } from "@/configs/env"
import { getAppConfig } from "@/configs/app-config"
import { networkByChainId } from "@/configs/networks"
import AppKitThemeSync from "@/components/appkit-theme-sync"
import LoginPage from "@/pages/login"
import VerifyWalletPage from "@/pages/verify-wallet"
import DashboardLayout from "@/components/dashboard/dashboard-layout"
import DashboardOverviewPage from "@/pages/dashboard-overview"
import DashboardAgentsPage from "@/pages/dashboard-agents"
import DashboardAgentDetailPage from "@/pages/dashboard-agent-detail"
import DashboardConversationsPage from "@/pages/dashboard-conversations"
import DashboardConversationDetailPage from "@/pages/dashboard-conversation-detail"
import DashboardModerationPage from "@/pages/dashboard-moderation"

const metadata = {
  name: "Agentic Orchestrator",
  description: "Console for the agentic orchestrator.",
  url: window.location.origin,
  icons: ["/logo-192x192.png"],
}

function App() {
  React.useMemo(() => {
    // App renders under ConfigGate, so the runtime config is loaded by now.
    // AppKit is initialized with every chain the orchestrator serves.
    const configured = Object.keys(getAppConfig().evmChains ?? {})
      .map((id) => networkByChainId[Number(id)])
      .filter((n): n is AppKitNetwork => Boolean(n))
    const networks: [AppKitNetwork, ...AppKitNetwork[]] =
      configured.length > 0
        ? (configured as [AppKitNetwork, ...AppKitNetwork[]])
        : [polygonAmoy]
    createAppKit({
      adapters: [new EthersAdapter()],
      networks,
      metadata,
      projectId: config.appKitProjectId,
      features: {
        email: false,
        socials: false,
        analytics: false,
        connectMethodsOrder: ["wallet"],
      },
    })
  }, [])

  return (
    <>
      <AppKitThemeSync />
      <Routes>
        <Route path="/" element={<Navigate to="/login" replace />} />
        <Route path="/login" element={<LoginPage />} />
        <Route path="/verify-wallet" element={<VerifyWalletPage />} />
        <Route path="/dashboard" element={<DashboardLayout />}>
          <Route
            index
            element={<Navigate to="/dashboard/overview" replace />}
          />
          <Route path="overview" element={<DashboardOverviewPage />} />
          <Route path="agents" element={<DashboardAgentsPage />} />
          <Route
            path="agents/:collectionAddress/:collectionTokenId"
            element={<DashboardAgentDetailPage />}
          />
          <Route
            path="conversations"
            element={<DashboardConversationsPage />}
          />
          <Route
            path="conversations/:conversationId"
            element={<DashboardConversationDetailPage />}
          />
          <Route path="moderation" element={<DashboardModerationPage />} />
        </Route>
        <Route path="*" element={<Navigate to="/login" replace />} />
      </Routes>
    </>
  )
}

export default App
