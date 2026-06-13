import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { ColumnDef } from "@tanstack/react-table";
import { MCPConnection } from "@/models/mcp-connection";
import { removeAgentMCPConnection } from "@/queries/api/agent.queries";
import { AgentMCPConnectionsDialogContext } from "../../contexts/agent-mcp-connections-dialog-context";
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

export const agentMCPConnectionsColumns: ColumnDef<MCPConnection>[] = [
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
    cell: ({ row }) => {
      return row.original.name;
    },
  },
  {
    id: "url",
    accessorKey: "url",
    header: ({ column }) => {
      return (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
        >
          Url
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
      return row.original.url;
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
              <DropdownMenuSeparator />
              <AgentMCPConnectionsDialogContext.Consumer>
                {({ agent, fetchMCPConnections }) => (
                  <DropdownMenuItem
                    onClick={async () => {
                      try {
                        await removeAgentMCPConnection(
                          agent!.collectionAddress,
                          agent!.collectionTokenId!,
                          row.original.id
                        );

                        toast.success("MCP connection removed successfully.");
                        await fetchMCPConnections();
                      } catch (error) {
                        toast.error("Error removing MCP connection.");
                        console.error(
                          "Error removing MCP connection:",
                          error
                        );
                      }
                    }}
                  >
                    Remove
                  </DropdownMenuItem>
                )}
              </AgentMCPConnectionsDialogContext.Consumer>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      );
    },
  },
];
