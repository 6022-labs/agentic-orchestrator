import React from "react";
import { Plus } from "lucide-react";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { agentContextsColumns } from "./agent-contexts-data-column";
import { AgentContextsDataTable } from "./agent-contexts-data-table";
import { useAgentContextsDialogContext } from "../../contexts/agent-contexts-dialog-context";
import { useCreateAgentContextDialogContext } from "../../contexts/create-agent-context-dialog-context";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";

export default function AgentContextsDialog() {
  const { isOpen, setIsOpen, agent, setAgent, isLoading } =
    useAgentContextsDialogContext();
  const { setIsOpen: setCreateModalOpen } =
    useCreateAgentContextDialogContext();

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
          setSearchQuery(undefined);
        }
      }}
    >
      <DialogContent loading={isLoading}>
        <DialogHeader>
          <DialogTitle>Agent Contexts</DialogTitle>
          <DialogDescription>
            <span>Manage agent contexts for </span>
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
                Add Agent Context
              </Button>
            </div>
          </div>
          <div>
            <AgentContextsDataTable
              searchQuery={searchQuery}
              columns={agentContextsColumns}
            />
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
}
