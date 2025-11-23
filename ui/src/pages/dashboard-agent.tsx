import { useNavigate } from "react-router-dom";
import DashboardAgents from "@/components/dashboard/agents";
import { DashboardHeader } from "@/components/dashboard/header";

export default function DashboardAgentPage() {
  const navigate = useNavigate();

  return (
    <div>
      <header>
        <DashboardHeader
          onAgentsClick={() => navigate("/dashboard/agents")}
          onOverviewClick={() => navigate("/dashboard/overview")}
          onModerationClick={() => navigate("/dashboard/moderation")}
        />
      </header>
      <main>
        <DashboardAgents />
      </main>
    </div>
  );
}
