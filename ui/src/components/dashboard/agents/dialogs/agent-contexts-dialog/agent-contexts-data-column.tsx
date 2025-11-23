import { Button } from "@/components/ui/button";
import { ColumnDef } from "@tanstack/react-table";
import { AgentContextsDialogContext } from "../../contexts/agent-contexts-dialog-context";
import { AddMessageToAgentContextDialogContext } from "../../contexts/add-message-to-agent-context-dialog-context";
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

export const agentContextsColumns: ColumnDef<string>[] = [
  {
    id: "key",
    header: ({ column }) => {
      return (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
        >
          Key
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
      return row.original;
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
              <>
                <DropdownMenuSeparator />
                <AgentContextsDialogContext.Consumer>
                  {({ setIsOpen: setAgentContextIsOpen }) => (
                    <AddMessageToAgentContextDialogContext.Consumer>
                      {({ setContext, setIsOpen: setAddMessageIsOpen }) => (
                        <DropdownMenuItem
                          onClick={async () => {
                            setAgentContextIsOpen(false);

                            setContext(row.original);
                            setAddMessageIsOpen(true);
                          }}
                        >
                          Add message
                        </DropdownMenuItem>
                      )}
                    </AddMessageToAgentContextDialogContext.Consumer>
                  )}
                </AgentContextsDialogContext.Consumer>
              </>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      );
    },
  },
];
