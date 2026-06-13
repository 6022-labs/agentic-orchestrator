import DashboardOverview from "@/components/dashboard/overview";
import { DashboardHeader } from "@/components/dashboard/header";
import { useNavigate } from "react-router-dom";

export default function DashboardOverviewPage() {
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
        <DashboardOverview />
      </main>
    </div>
  );
}
