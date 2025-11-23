import { useNavigate } from "react-router-dom";
import { DashboardHeader } from "@/components/dashboard/header";
import DashboardModeration from "@/components/dashboard/moderation";

export default function DashboardModerationPage() {
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
        <DashboardModeration />
      </main>
    </div>
  );
}
