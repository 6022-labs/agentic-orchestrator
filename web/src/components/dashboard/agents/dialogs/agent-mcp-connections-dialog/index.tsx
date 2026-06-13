import React from "react";
import { Plus } from "lucide-react";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { agentMCPConnectionsColumns } from "./agent-mcp-connections-data-column";
import { AgentMCPConnectionsDataTable } from "./agent-mcp-connections-data-table";
import { useAgentMCPConnectionsDialogContext } from "../../contexts/agent-mcp-connections-dialog-context";
import { useCreateAgentMCPConnectionDialogContext } from "../../contexts/create-agent-mcp-connection-dialog-context";
import {
  Dialog,
  DialogTitle,
  DialogHeader,
  DialogContent,
  DialogDescription,
} from "@/components/ui/dialog";

export default function MCPConnectionsDialog() {
  const { isOpen, setIsOpen, agent, setAgent, isLoading } =
    useAgentMCPConnectionsDialogContext();
  const { setIsOpen: setCreateModalOpen } =
    useCreateAgentMCPConnectionDialogContext();

  const [searchQuery, setSearchQuery] = React.useState<string | undefined>(
    undefined
  );

  return (
    <Dialog
      open={isOpen}
      onOpenChange={(open) => {
        setIsOpen(open);
        if (!open) {
          setAgent(undefined);
        }
      }}
    >
      <DialogContent loading={isLoading}>
        <DialogHeader>
          <DialogTitle>MCP Connections</DialogTitle>
          <DialogDescription>
            <span>Manage MCP connections for agent </span>
            <span className="font-bold">{agent?.name}</span>.
          </DialogDescription>
        </DialogHeader>
        <div className="flex-1 space-y-4 pt-2">
          <div className="flex w-full items-center space-y-2">
            <div className="flex w-full flex-row items-center justify-between gap-4">
              <Input
                placeholder="Search keys..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
              />
              <Button
                onClick={() => {
                  setIsOpen(false);
                  setCreateModalOpen(true);
                }}
              >
                <Plus className="mr-1" />
                Add MCP
              </Button>
            </div>
          </div>
          <div>
            <AgentMCPConnectionsDataTable
              searchQuery={searchQuery}
              columns={agentMCPConnectionsColumns}
            />
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
}
