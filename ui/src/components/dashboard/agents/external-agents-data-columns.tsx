import { Agent } from "@/models/agent";
import { Button } from "@/components/ui/button";
import { truncateEthAddress } from "@/lib/utils";
import { ColumnDef } from "@tanstack/react-table";
import { AuthContext } from "@/contexts/auth-context";
import ViewOnOpenseaAction from "./row-actions/view-on-opensea-action";
import ViewOnRaribleAction from "./row-actions/view-on-rarible-action";
import ViewOnScannerAction from "./row-actions/view-on-scanner-action";
import { AgentCollectionsContext } from "@/contexts/agent-collections-context";
import { AgentImportationStatus } from "@/models/enums/agent-importation-status";
import { AgentImportationThirdStepDialogContext } from "./contexts/agent-importation-third-step-dialog-context";
// import { AgentImportationSecondStepDialogContext } from "./contexts/agent-importation-second-step-dialog-context";
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

export const externalAgentsColumns: ColumnDef<Agent>[] = [
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
            row.original.owner
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
    accessorKey: "creationStatus",
    cell: ({ row }) => {
      return row.original.collectionTokenId;
    },
  },
  {
    accessorKey: "hasAccess",
    header: ({ column }) => {
      return (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
        >
          Access ?
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
    cell: () => {
      return "No";
    },
  },
  {
    id: "actions",
    cell: ({ row }) => {
      return (
        <div className="flex justify-end">
          <AuthContext.Consumer>
            {({ user }) => {
              return (
                <DropdownMenu>
                  <DropdownMenuTrigger asChild>
                    <Button variant="ghost" className="h-8 w-8 p-0">
                      <span className="sr-only">Open menu</span>
                      <EllipsisVertical className="ml-auto h-4 w-4" />
                    </Button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent align="end">
                    <DropdownMenuLabel>Actions</DropdownMenuLabel>
                    <ViewOnScannerAction address={row.original.walletAddress} />
                    <ViewOnOpenseaAction
                      collectionTokenId={row.original.collectionTokenId!}
                      collectionAddress={row.original.collectionAddress}
                    />
                    <ViewOnRaribleAction
                      collectionTokenId={row.original.collectionTokenId!}
                      collectionAddress={row.original.collectionAddress}
                    />
                    <DropdownMenuSeparator />
                    {user?.address.toLowerCase() ==
                      row.original.owner.toLowerCase() && (
                      <>
                        {/* {(!row.original.importationStatus ||
                        row.original.importationStatus ===
                          AgentImportationStatus.FirstStep) && (
                        <AgentImportationSecondStepDialogContext.Consumer>
                          {({ open }) => (
                            <DropdownMenuItem
                              onClick={() => {
                                open(row.original);
                              }}
                            >
                              Start import
                            </DropdownMenuItem>
                          )}
                        </AgentImportationSecondStepDialogContext.Consumer>
                      )} */}
                        {(row.original.importationStatus ===
                          AgentImportationStatus.SecondStep ||
                          row.original.importationStatus ===
                            AgentImportationStatus.ThirdStep) && (
                          <AgentImportationThirdStepDialogContext.Consumer>
                            {({ setAgent, setIsOpen }) => (
                              <DropdownMenuItem
                                onClick={() => {
                                  setAgent(row.original);
                                  setIsOpen(true);
                                }}
                              >
                                Finish import
                              </DropdownMenuItem>
                            )}
                          </AgentImportationThirdStepDialogContext.Consumer>
                        )}
                        <DropdownMenuItem>Transfer</DropdownMenuItem>
                      </>
                    )}
                  </DropdownMenuContent>
                </DropdownMenu>
              );
            }}
          </AuthContext.Consumer>
        </div>
      );
    },
  },
];
