import React from "react";
import LoginPage from "./pages/login";
import { useSiweConfig } from "./hooks/useSiwe";
import { createAppKit } from "@reown/appkit/react";
import VerifyWalletPage from "./pages/verify-wallet";
import DashboardAgentPage from "./pages/dashboard-agent";
import { Route, Routes, Navigate } from "react-router-dom";
import { EthersAdapter } from "@reown/appkit-adapter-ethers";
import DashboardOverviewPage from "./pages/dashboard-overview";
import DashboardModerationPage from "./pages/dashboard-moderation";
import {
  hardhat,
  polygon,
  polygonAmoy,
  AppKitNetwork,
  citreaTestnet,
} from "@reown/appkit/networks";

const projectId = import.meta.env.VITE_APP_KIT_PROJECT_ID;

let networks: [AppKitNetwork, ...AppKitNetwork[]];

if (import.meta.env.VITE_CHAIN_ID == "137") {
  networks = [polygon];
}

if (import.meta.env.VITE_CHAIN_ID == "80002") {
  networks = [polygonAmoy];
}

if (import.meta.env.VITE_CHAIN_ID == "31337") {
  networks = [hardhat];
}

if (import.meta.env.VITE_CHAIN_ID == "5115") {
  networks = [citreaTestnet];
}

const metadata = {
  name: "Agentic AI Orchestrator",
  description: "Agentic AI Orchestrator is a platform for AI orchestration.",
  url: "https://reown.com/appkit", // origin must match your domain & subdomain
  debug: import.meta.env.VITE_ENVIRONMENT != "production",
  icons: ["/logo-192x192.png", "/logo-512x512.png"],
};

function App() {
  const siweConfig = useSiweConfig();

  React.useMemo(() => {
    createAppKit({
      adapters: [new EthersAdapter()],
      networks,
      metadata,
      projectId,
      features: {
        send: false,
        swaps: false,
        email: false,
        onramp: false,
        socials: false,
        analytics: false,
        connectMethodsOrder: ["wallet"],
      },
      siweConfig,
    });
  }, [siweConfig]);

  return (
    <Routes>
      <Route path="/" element={<Navigate to="/login" />} />
      <Route path="/login" element={<LoginPage />} />
      <Route path="/verify-wallet" element={<VerifyWalletPage />} />
      <Route
        path="/dashboard"
        element={<Navigate to="/dashboard/overview" />}
      />
      <Route path="/dashboard/overview" element={<DashboardOverviewPage />} />
      <Route path="/dashboard/agents" element={<DashboardAgentPage />} />
      <Route
        path="/dashboard/moderation"
        element={<DashboardModerationPage />}
      />
    </Routes>
  );
}

export default App;
