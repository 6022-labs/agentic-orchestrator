import TransferAgentDialog from "./dialogs/transfer-agent-dialog";
import { ExternalAgentsDataTable } from "./external-agents-data-table";
import { externalAgentsColumns } from "./external-agents-data-columns";
import UpsertAgentAttributeDialog from "./dialogs/upsert-agent-attribute-dialog";
import { AgentDashboardContextProvider } from "./contexts/agent-dashboard-context";
import CreateAgentCollectionDialog from "./dialogs/create-agent-collection-dialog";
import { TransferAgentDialogContextProvider } from "./contexts/transfer-agent-dialog-context";
import { UpsertAgentAttributeDialogContextProvider } from "./contexts/upsert-agent-attribute-dialog-context";
import { CreateAgentCollectionDialogContextProvider } from "./contexts/create-agent-collection-dialog-context";

export default function DashboardAgents() {
  return (
    <AgentDashboardContextProvider>
      <CreateAgentCollectionDialogContextProvider>
        <TransferAgentDialogContextProvider>
          <UpsertAgentAttributeDialogContextProvider>
            <div className="flex-1 space-y-4 p-8 pt-6">
              <div className="flex items-center justify-between space-y-2">
                <h2 className="text-3xl font-bold tracking-tight">Agents</h2>
                <div className="flex flex-row items-center gap-2">
                  <CreateAgentCollectionDialog />
                  <TransferAgentDialog />
                  <UpsertAgentAttributeDialog />
                </div>
              </div>
              <div>
                <ExternalAgentsDataTable columns={externalAgentsColumns} />
              </div>
            </div>
          </UpsertAgentAttributeDialogContextProvider>
        </TransferAgentDialogContextProvider>
      </CreateAgentCollectionDialogContextProvider>
    </AgentDashboardContextProvider>
  );
}
