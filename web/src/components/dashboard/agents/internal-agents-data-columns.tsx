import { Agent } from "@/models/agent";
import { toast } from "sonner";
import { mintAgent } from "@/lib/mint-agent";
import { Button } from "@/components/ui/button";
import { truncateEthAddress } from "@/lib/utils";
import { ColumnDef } from "@tanstack/react-table";
import { AuthContext } from "@/contexts/auth-context";
import ViewOnRaribleAction from "./row-actions/view-on-rarible-action";
import ViewOnOpenseaAction from "./row-actions/view-on-opensea-action";
import ViewOnScannerAction from "./row-actions/view-on-scanner-action";
import { UserModeratorContext } from "@/contexts/user-moderator-context";
import { AgentDashboardContext } from "./contexts/agent-dashboard-context";
import { WalletProviderContext } from "@/contexts/wallet-provider-context";
import { AgentLLMSetupStatus } from "@/models/enums/agent-llm-setup-status";
import { AgentCollectionsContext } from "@/contexts/agent-collections-context";
import { setupAgentDatabaseSettingsSecondStep } from "@/queries/api/agent.queries";
import CopyAddAgentCommandAction from "./row-actions/copy-add-agent-command-action";
import { AgentContextsDialogContext } from "./contexts/agent-contexts-dialog-context";
import { TransferAgentDialogContext } from "./contexts/transfer-agent-dialog-context";
import { AgentDatabaseSetupStatus } from "@/models/enums/agent-database-setup-status";
import CopyRemoveAgentCommandAction from "./row-actions/copy-remove-agent-command-action";
import { AgentMCPConnectionsDialogContext } from "./contexts/agent-mcp-connections-dialog-context";
import { UpdateAgentPromptDialogContext } from "./contexts/update-agent-prompt-dialog-context";
import { SetupAgentLLMSettingsDialogContext } from "./contexts/setup-agent-llm-settings-dialog-context";
import { SetupAgentDatabaseSettingsFirstStepDialogContext } from "./contexts/setup-agent-database-settings-first-step-dialog-context";
import {
  ArrowUp,
  ArrowDown,
  ArrowUpDown,
  EllipsisVertical,
} from "lucide-react";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { AgentRole, AgentRoleDisplayNames } from "@/models/enums/agent-role";
import CopyPairSlackUserToHumanAgentAction from "./row-actions/copy-pair-slack-user-to-human-agent-action";
import { CreateAgentFirstStepDialogContext } from "./contexts/create-agent-first-step-dialog-context";
import { UpsertAgentAttributeDialogContext } from "./contexts/upsert-agent-attribute-dialog-context";

export const internalAgentsColumns: ColumnDef<Agent>[] = [
  {
    id: "name",
    accessorKey: "name",
    header: ({ column }) => {
      return (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
        >
          Name
          {column.getIsSorted() === false ? (
            <ArrowUpDown className="ml-2 h-4 w-4" />
          ) : column.getIsSorted() === "asc" ? (
            <ArrowUp className="ml-2 h-4 w-4" />
          ) : (
            <ArrowDown className="ml-2 h-4 w-4" />
          )}
        </Button>
      );
    },
  },
  {
    id: "walletAddress",
    accessorKey: "walletAddress",
    header: ({ column }) => {
      return (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
        >
          Address
          {column.getIsSorted() === false ? (
            <ArrowUpDown className="ml-2 h-4 w-4" />
          ) : column.getIsSorted() === "asc" ? (
            <ArrowUp className="ml-2 h-4 w-4" />
          ) : (
            <ArrowDown className="ml-2 h-4 w-4" />
          )}
        </Button>
      );
    },
    cell: ({ row }) => {
      return (
        <a
          className="underline cursor-pointer"
          href={
            import.meta.env.VITE_SCANNER_BASE_URL +
            "/address/" +
            row.original.walletAddress
          }
          target="_blank"
          rel="noreferrer"
        >
          {truncateEthAddress(row.original.walletAddress, 8, 6)}
        </a>
      );
    },
  },
  {
    id: "owner",
    accessorKey: "owner",
    header: ({ column }) => {
      return (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
        >
          Owner
          {column.getIsSorted() === false ? (
            <ArrowUpDown className="ml-2 h-4 w-4" />
          ) : column.getIsSorted() === "asc" ? (
            <ArrowUp className="ml-2 h-4 w-4" />
          ) : (
            <ArrowDown className="ml-2 h-4 w-4" />
          )}
        </Button>
      );
    },
    cell: ({ row }) => {
      return (
        <a
          className="underline cursor-pointer"
          href={
            import.meta.env.VITE_SCANNER_BASE_URL +
            "/address/" +
            row.original.owner
          }
          target="_blank"
          rel="noreferrer"
        >
          {truncateEthAddress(row.original.owner, 8, 6)}
        </a>
      );
    },
  },
  {
    id: "role",
    header: "Role",
    accessorKey: "role",
    cell: ({ row }) => {
      return AgentRoleDisplayNames[row.original.role];
    },
  },
  {
    id: "collectionAddress",
    accessorKey: "collectionAddress",
    header: ({ column }) => {
      return (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
        >
          Collection
          {column.getIsSorted() === false ? (
            <ArrowUpDown className="ml-2 h-4 w-4" />
          ) : column.getIsSorted() === "asc" ? (
            <ArrowUp className="ml-2 h-4 w-4" />
          ) : (
            <ArrowDown className="ml-2 h-4 w-4" />
          )}
        </Button>
      );
    },
    cell: ({ row }) => {
      return (
        <AgentCollectionsContext.Consumer>
          {({ collections }) => {
            const collection = collections.find(
              (c) => c.address === row.original.collectionAddress
            );

            return (
              <a
                className="underline cursor-pointer"
                href={
                  import.meta.env.VITE_SCANNER_BASE_URL +
                  "/address/" +
                  row.original.collectionAddress
                }
                target="_blank"
                rel="noreferrer"
              >
                {collection?.name
                  ? collection.name +
                    " (" +
                    truncateEthAddress(row.original.collectionAddress, 8, 6) +
                    ")"
                  : truncateEthAddress(row.original.collectionAddress, 8, 6)}
              </a>
            );
          }}
        </AgentCollectionsContext.Consumer>
      );
    },
  },
  {
    header: "Token ID",
    accessorKey: "collectionTokenId",
    cell: ({ row }) => {
      return row.original.collectionTokenId;
    },
  },
  {
    header: "Proposal ID",
    accessorKey: "proposalId",
    cell: ({ row }) => {
      return row.original.proposalId;
    },
  },
  {
    header: "Setup",
    cell: ({ row }) => {
      return (
        <>
          {row.original.isSetupCompleted ? (
            <span className="text-green-500">Completed</span>
          ) : (
            <span className="text-red-500">In Progress</span>
          )}
        </>
      );
    },
  },
  {
    id: "actions",
    cell: ({ row }) => {
      return (
        <div className="flex justify-end">
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" className="h-8 w-8 p-0">
                <span className="sr-only">Open menu</span>
                <EllipsisVertical className="ml-auto h-4 w-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuLabel>Actions</DropdownMenuLabel>
              <AuthContext.Consumer>
                {({ user }) => {
                  return (
                    <>
                      {row.original.collectionTokenId ? (
                        <>
                          <DropdownMenuSeparator />
                          <ViewOnScannerAction
                            address={row.original.walletAddress}
                          />
                          <ViewOnOpenseaAction
                            collectionTokenId={row.original.collectionTokenId}
                            collectionAddress={row.original.collectionAddress}
                          />
                          <ViewOnRaribleAction
                            collectionTokenId={row.original.collectionTokenId}
                            collectionAddress={row.original.collectionAddress}
                          />
                          <DropdownMenuSeparator />
                          <CopyAddAgentCommandAction
                            collectionTokenId={row.original.collectionTokenId}
                            collectionAddress={row.original.collectionAddress}
                          />
                          <CopyRemoveAgentCommandAction
                            collectionTokenId={row.original.collectionTokenId}
                            collectionAddress={row.original.collectionAddress}
                          />
                          <CopyPairSlackUserToHumanAgentAction
                            collectionTokenId={row.original.collectionTokenId}
                            collectionAddress={row.original.collectionAddress}
                          />
                          {row.original.owner?.toLowerCase() ===
                            user?.address.toLowerCase() && (
                            <>
                              <DropdownMenuSeparator />
                              <TransferAgentDialogContext.Consumer>
                                {({ setAgent, setIsOpen }) => (
                                  <DropdownMenuItem
                                    onClick={() => {
                                      setAgent(row.original);
                                      setIsOpen(true);
                                    }}
                                  >
                                    Transfer
                                  </DropdownMenuItem>
                                )}
                              </TransferAgentDialogContext.Consumer>
                              {row.original.role === AgentRole.Human && (
                                <CreateAgentFirstStepDialogContext.Consumer>
                                  {({
                                    setIsOpen,
                                    setDefaultValues,
                                    setDisabledValues,
                                    setAccessibleRoles,
                                  }) => (
                                    <DropdownMenuItem
                                      onClick={() => {
                                        setDefaultValues({
                                          collectionAddress:
                                            row.original.collectionAddress,
                                          cloneOfAgentId:
                                            row.original.collectionTokenId,
                                          role: AgentRole.Clone,
                                        });
                                        setDisabledValues({
                                          role: true,
                                          cloneOfAgentId: true,
                                          collectionAddress: true,
                                        });
                                        setAccessibleRoles([AgentRole.Clone]);
                                        setIsOpen(true);
                                      }}
                                    >
                                      Create digital twin
                                    </DropdownMenuItem>
                                  )}
                                </CreateAgentFirstStepDialogContext.Consumer>
                              )}
                              {row.original.role !== AgentRole.Human && (
                                <>
                                  <SetupAgentLLMSettingsDialogContext.Consumer>
                                    {({ setAgent, setIsOpen }) => (
                                      <DropdownMenuItem
                                        onClick={() => {
                                          setAgent(row.original);
                                          setIsOpen(true);
                                        }}
                                      >
                                        {row.original.llmSetupStatus ===
                                        AgentLLMSetupStatus.Established
                                          ? "Manage LLM settings"
                                          : "Setup LLM settings"}
                                      </DropdownMenuItem>
                                    )}
                                  </SetupAgentLLMSettingsDialogContext.Consumer>
                                  <AgentMCPConnectionsDialogContext.Consumer>
                                    {({ setAgent, setIsOpen }) => (
                                      <DropdownMenuItem
                                        onClick={() => {
                                          setAgent(row.original);
                                          setIsOpen(true);
                                        }}
                                      >
                                        Manage MCPs
                                      </DropdownMenuItem>
                                    )}
                                  </AgentMCPConnectionsDialogContext.Consumer>
                                  <UpsertAgentAttributeDialogContext.Consumer>
                                    {({ setAgent, setIsOpen }) => (
                                      <DropdownMenuItem
                                        onClick={() => {
                                          setAgent(row.original);
                                          setIsOpen(true);
                                        }}
                                      >
                                        Manage attributes
                                      </DropdownMenuItem>
                                    )}
                                  </UpsertAgentAttributeDialogContext.Consumer>
                                </>
                              )}
                              {!row.original.databaseSetupStatus && (
                                <SetupAgentDatabaseSettingsFirstStepDialogContext.Consumer>
                                  {({ setAgent, setIsOpen }) => (
                                    <DropdownMenuItem
                                      onClick={() => {
                                        setAgent(row.original);
                                        setIsOpen(true);
                                      }}
                                    >
                                      Setup Database settings
                                    </DropdownMenuItem>
                                  )}
                                </SetupAgentDatabaseSettingsFirstStepDialogContext.Consumer>
                              )}
                              {row.original.databaseSetupStatus ===
                                AgentDatabaseSetupStatus.Configured && (
                                <AgentDashboardContext.Consumer>
                                  {({ fetchInternalAgents }) => (
                                    <DropdownMenuItem
                                      onClick={async () => {
                                        try {
                                          await setupAgentDatabaseSettingsSecondStep(
                                            row.original.collectionAddress,
                                            row.original.collectionTokenId!
                                          );
                                          toast.success(
                                            "Agent database setup completed successfully."
                                          );
                                        } catch (error) {
                                          console.error(
                                            "Error setting up agent database:",
                                            error
                                          );
                                          toast.error(
                                            "Failed to setup agent database. Please try again."
                                          );
                                        } finally {
                                          fetchInternalAgents();
                                        }
                                      }}
                                    >
                                      Setup Database
                                    </DropdownMenuItem>
                                  )}
                                </AgentDashboardContext.Consumer>
                              )}
                              {row.original.isSetupCompleted && (
                                <>
                                  {row.original.role !== AgentRole.Human && (
                                    <>
                                      <UpdateAgentPromptDialogContext.Consumer>
                                        {({ setAgent, setIsOpen }) => (
                                          <DropdownMenuItem
                                            onClick={() => {
                                              setAgent(row.original);
                                              setIsOpen(true);
                                            }}
                                          >
                                            Manage prompts
                                          </DropdownMenuItem>
                                        )}
                                      </UpdateAgentPromptDialogContext.Consumer>
                                      <AgentContextsDialogContext.Consumer>
                                        {({ setAgent, setIsOpen }) => (
                                          <DropdownMenuItem
                                            onClick={() => {
                                              setAgent(row.original);
                                              setIsOpen(true);
                                            }}
                                          >
                                            Manage contexts
                                          </DropdownMenuItem>
                                        )}
                                      </AgentContextsDialogContext.Consumer>
                                    </>
                                  )}
                                </>
                              )}
                            </>
                          )}
                        </>
                      ) : (
                        <>
                          {row.original.creator.toLowerCase() ===
                            user?.address.toLowerCase() &&
                            !row.original.proposalId && (
                              <>
                                <DropdownMenuSeparator />
                                <WalletProviderContext.Consumer>
                                  {({ wallet }) => (
                                    <AgentDashboardContext.Consumer>
                                      {({ fetchInternalAgents }) => (
                                        <DropdownMenuItem
                                          onClick={async () => {
                                            try {
                                              const result = await mintAgent(
                                                row.original,
                                                wallet!
                                              );

                                              if (result.type === "proposal") {
                                                toast.success(
                                                  "Agent mint proposal created successfully."
                                                );

                                                return;
                                              }

                                              toast.success(
                                                "Agent minted successfully."
                                              );
                                            } catch (error) {
                                              console.error(
                                                "Failed to mint agent:",
                                                error
                                              );
                                              toast.error(
                                                "Failed to mint agent. Please try again."
                                              );
                                            } finally {
                                              fetchInternalAgents();
                                            }
                                          }}
                                        >
                                          <UserModeratorContext.Consumer>
                                            {({ collections }) => (
                                              <>
                                                {collections.some(
                                                  (c) =>
                                                    c.address.toLowerCase() ===
                                                    row.original.collectionAddress.toLowerCase()
                                                )
                                                  ? "Mint"
                                                  : "Propose mint"}
                                              </>
                                            )}
                                          </UserModeratorContext.Consumer>
                                        </DropdownMenuItem>
                                      )}
                                    </AgentDashboardContext.Consumer>
                                  )}
                                </WalletProviderContext.Consumer>
                              </>
                            )}
                        </>
                      )}
                    </>
                  );
                }}
              </AuthContext.Consumer>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      );
    },
  },
];
