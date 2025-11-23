import "./index.css";
import App from "./App.tsx";
import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { BrowserRouter } from "react-router-dom";
import { Toaster } from "@/components/ui/sonner.tsx";
import { initAxiosDefaults } from "./configs/axios.ts";
import { AuthProvider } from "./contexts/auth-context.tsx";
import { WalletProvider } from "./contexts/wallet-provider-context.tsx";
import { UserModeratorProvider } from "./contexts/user-moderator-context.tsx";
import { AgentCollectionsProvider } from "./contexts/agent-collections-context.tsx";

initAxiosDefaults();

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <BrowserRouter>
      <AuthProvider>
        <WalletProvider>
          <AgentCollectionsProvider>
            <UserModeratorProvider>
              <Toaster />
              <App />
            </UserModeratorProvider>
          </AgentCollectionsProvider>
        </WalletProvider>
      </AuthProvider>
    </BrowserRouter>
  </StrictMode>
);
