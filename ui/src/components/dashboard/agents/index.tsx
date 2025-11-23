import TransferAgentDialog from "./dialogs/transfer-agent-dialog";
import AgentContextsDialog from "./dialogs/agent-contexts-dialog";
import { internalAgentsColumns } from "./internal-agents-data-columns";
import { InternalAgentsDataTable } from "./internal-agents-data-table";
import { ExternalAgentsDataTable } from "./external-agents-data-table";
import { externalAgentsColumns } from "./external-agents-data-columns";
import UpdateAgentPromptDialog from "./dialogs/update-agent-prompt-dialog";
import CreateAgentContextDialog from "./dialogs/create-agent-context-dialog";
import AgentMCPConnectionsDialog from "./dialogs/agent-mcp-connections-dialog";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import UpsertAgentAttributeDialog from "./dialogs/upsert-agent-attribute-dialog";
import CreateAgentFirstStepDialog from "./dialogs/create-agent-first-step-dialog";
import { AgentDashboardContextProvider } from "./contexts/agent-dashboard-context";
import CreateAgentCollectionDialog from "./dialogs/create-agent-collection-dialog";
import SetupAgentLLMSettingsDialog from "./dialogs/setup-agent-llm-settings-dialog";
import CreateAgentMCPConnectionDialog from "./dialogs/create-agent-mcp-connection-dialog";
import AddMessageToAgentContextDialog from "./dialogs/add-message-to-agent-context-dialog";
import AgentImportationThirdStepDialog from "./dialogs/agent-importation-third-step-dialog";
import { TransferAgentDialogContextProvider } from "./contexts/transfer-agent-dialog-context";
import { AgentContextsDialogContextProvider } from "./contexts/agent-contexts-dialog-context";
import { AgentMCPConnectionsDialogContextProvider } from "./contexts/agent-mcp-connections-dialog-context";
import { UpdateAgentPromptDialogContextProvider } from "./contexts/update-agent-prompt-dialog-context";
import { CreateAgentContextDialogContextProvider } from "./contexts/create-agent-context-dialog-context";
import SetupAgentDatabaseSettingsDialog from "./dialogs/setup-agent-database-settings-first-step-dialog";
import { UpsertAgentAttributeDialogContextProvider } from "./contexts/upsert-agent-attribute-dialog-context";
import { CreateAgentFirstStepDialogContextProvider } from "./contexts/create-agent-first-step-dialog-context";
import { CreateAgentCollectionDialogContextProvider } from "./contexts/create-agent-collection-dialog-context";
import { SetupAgentLLMSettingsDialogContextProvider } from "./contexts/setup-agent-llm-settings-dialog-context";
import { CreateAgentMCPConnectionDialogContextProvider } from "./contexts/create-agent-mcp-connection-dialog-context";
import { AddMessageToAgentContextDialogContextProvider } from "./contexts/add-message-to-agent-context-dialog-context";
import { AgentImportationThirdStepDialogContextProvider } from "./contexts/agent-importation-third-step-dialog-context";
import { SetupAgentDatabaseSettingsFirstStepDialogContextProvider } from "./contexts/setup-agent-database-settings-first-step-dialog-context";
// import { AgentImportationSecondStepDialogContextProvider } from "./contexts/agent-importation-second-step-dialog-context";

export default function DashboardAgents() {
  return (
    <AgentDashboardContextProvider>
      <CreateAgentCollectionDialogContextProvider>
        <CreateAgentFirstStepDialogContextProvider>
          <TransferAgentDialogContextProvider>
            {/* <AgentImportationSecondStepDialogContextProvider> */}
            <AgentImportationThirdStepDialogContextProvider>
              <UpdateAgentPromptDialogContextProvider>
                <AgentContextsDialogContextProvider>
                  <CreateAgentContextDialogContextProvider>
                    <AddMessageToAgentContextDialogContextProvider>
                      <SetupAgentLLMSettingsDialogContextProvider>
                        <SetupAgentDatabaseSettingsFirstStepDialogContextProvider>
                          <AgentMCPConnectionsDialogContextProvider>
                            <CreateAgentMCPConnectionDialogContextProvider>
                              <UpsertAgentAttributeDialogContextProvider>
                                <div className="flex-1 space-y-4 p-8 pt-6">
                                  <div className="flex items-center justify-between space-y-2">
                                    <h2 className="text-3xl font-bold tracking-tight">
                                      Agents
                                    </h2>
                                    <div className="flex flex-row items-center gap-2">
                                      <CreateAgentCollectionDialog />
                                      <CreateAgentFirstStepDialog />
                                      <AgentImportationThirdStepDialog />
                                      <UpdateAgentPromptDialog />
                                      <CreateAgentContextDialog />
                                      <AddMessageToAgentContextDialog />
                                      <SetupAgentLLMSettingsDialog />
                                      <SetupAgentDatabaseSettingsDialog />
                                      <TransferAgentDialog />
                                      <AgentMCPConnectionsDialog />
                                      <CreateAgentMCPConnectionDialog />
                                      <AgentContextsDialog />
                                      <UpsertAgentAttributeDialog />
                                    </div>
                                  </div>
                                  <div>
                                    <Tabs
                                      defaultValue="internals"
                                      className="gap-6"
                                    >
                                      <TabsList className="gap-x-2">
                                        <TabsTrigger value="internals">
                                          Internal agents
                                        </TabsTrigger>
                                        <TabsTrigger value="externals">
                                          External agents
                                        </TabsTrigger>
                                      </TabsList>
                                      <TabsContent value="internals">
                                        <InternalAgentsDataTable
                                          columns={internalAgentsColumns}
                                        />
                                      </TabsContent>
                                      <TabsContent value="externals">
                                        <ExternalAgentsDataTable
                                          columns={externalAgentsColumns}
                                        />
                                      </TabsContent>
                                    </Tabs>
                                  </div>
                                </div>
                              </UpsertAgentAttributeDialogContextProvider>
                            </CreateAgentMCPConnectionDialogContextProvider>
                          </AgentMCPConnectionsDialogContextProvider>
                        </SetupAgentDatabaseSettingsFirstStepDialogContextProvider>
                      </SetupAgentLLMSettingsDialogContextProvider>
                    </AddMessageToAgentContextDialogContextProvider>
                  </CreateAgentContextDialogContextProvider>
                </AgentContextsDialogContextProvider>
              </UpdateAgentPromptDialogContextProvider>
            </AgentImportationThirdStepDialogContextProvider>
          </TransferAgentDialogContextProvider>
          {/* <AgentImportationSecondStepDialogContextProvider> */}
        </CreateAgentFirstStepDialogContextProvider>
      </CreateAgentCollectionDialogContextProvider>
    </AgentDashboardContextProvider>
  );
}
