import React from "react";
import { Button } from "@/components/ui/button";
import { useAuth } from "@/contexts/auth-context";
import { useAgentMCPConnectionsDialogContext } from "../../contexts/agent-mcp-connections-dialog-context";
import {
  ColumnDef,
  flexRender,
  SortingState,
  useReactTable,
  getCoreRowModel,
} from "@tanstack/react-table";
import {
  Table,
  TableRow,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
} from "@/components/ui/table";
import {
  Select,
  SelectItem,
  SelectValue,
  SelectTrigger,
  SelectContent,
} from "@/components/ui/select";

interface DataTableProps<TData, TValue> {
  columns: ColumnDef<TData, TValue>[];
  searchQuery?: string;
}

export function AgentMCPConnectionsDataTable<TData, TValue>({
  columns,
  searchQuery,
}: DataTableProps<TData, TValue>) {
  const { mcpConnections, mcpConnectionsPagination, fetchMCPConnections } =
    useAgentMCPConnectionsDialogContext();
  const { jwt } = useAuth();

  const [sorting, setSorting] = React.useState<SortingState>([]);

  const table = useReactTable({
    data: mcpConnections as TData[],
    columns,
    onSortingChange: setSorting,
    getCoreRowModel: getCoreRowModel(),
    state: {
      sorting,
    },
  });

  React.useEffect(() => {
    const sortBy = sorting[0]?.id as string | undefined;
    const sortDirection = sorting[0]?.desc ? "desc" : "asc";

    fetchMCPConnections(
      mcpConnectionsPagination.page,
      searchQuery,
      mcpConnectionsPagination.elementsPerPage,
      sortBy,
      sortBy ? sortDirection : undefined
    );
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [jwt, sorting, searchQuery]);

  const onNextPage = () => {
    const maxPage = Math.ceil(
      mcpConnectionsPagination.totalElements /
        mcpConnectionsPagination.elementsPerPage
    );
    if (mcpConnectionsPagination.page < maxPage) {
      const sortBy = sorting[0]?.id as string | undefined;
      const sortDirection = sorting[0]?.desc ? "desc" : "asc";
      fetchMCPConnections(
        mcpConnectionsPagination.page + 1,
        searchQuery,
        mcpConnectionsPagination.elementsPerPage,
        sortBy,
        sortBy ? sortDirection : undefined
      );
    }
  };

  const onPreviousPage = () => {
    if (mcpConnectionsPagination.page > 1) {
      const sortBy = sorting[0]?.id as string | undefined;
      const sortDirection = sorting[0]?.desc ? "desc" : "asc";
      fetchMCPConnections(
        mcpConnectionsPagination.page - 1,
        searchQuery,
        mcpConnectionsPagination.elementsPerPage,
        sortBy,
        sortBy ? sortDirection : undefined
      );
    }
  };

  const onElementsPerPageChange = (value: string) => {
    const sortBy = sorting[0]?.id as string | undefined;
    const sortDirection = sorting[0]?.desc ? "desc" : "asc";
    fetchMCPConnections(
      1,
      searchQuery,
      parseInt(value, 10),
      sortBy,
      sortBy ? sortDirection : undefined
    );
  };

  return (
    <div>
      <div className="rounded-md border grid overflow-hidden">
        <Table className="text-center">
          <TableHeader>
            {table.getHeaderGroups().map((headerGroup) => (
              <TableRow key={headerGroup.id}>
                {headerGroup.headers.map((header) => {
                  return (
                    <TableHead className="text-center" key={header.id}>
                      {header.isPlaceholder
                        ? null
                        : flexRender(
                            header.column.columnDef.header,
                            header.getContext()
                          )}
                    </TableHead>
                  );
                })}
              </TableRow>
            ))}
          </TableHeader>
          <TableBody>
            {table.getRowModel().rows?.length ? (
              table.getRowModel().rows.map((row) => (
                <TableRow
                  key={row.id}
                  data-state={row.getIsSelected() && "selected"}
                >
                  {row.getVisibleCells().map((cell) => (
                    <TableCell key={cell.id}>
                      {flexRender(
                        cell.column.columnDef.cell,
                        cell.getContext()
                      )}
                    </TableCell>
                  ))}
                </TableRow>
              ))
            ) : (
              <TableRow>
                <TableCell
                  colSpan={columns.length}
                  className="h-24 text-center"
                >
                  No results.
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </div>
      <div className="flex items-center justify-between w-full py-4 px-2">
        <div className="flex items-center space-x-2">
          <Button
            variant="outline"
            size="sm"
            onClick={onPreviousPage}
            disabled={mcpConnectionsPagination.page <= 1}
          >
            Previous
          </Button>
          <Button
            variant="outline"
            size="sm"
            onClick={onNextPage}
            disabled={
              mcpConnectionsPagination.page >=
              Math.ceil(
                mcpConnectionsPagination.totalElements /
                  mcpConnectionsPagination.elementsPerPage
              )
            }
          >
            Next
          </Button>
        </div>
        <div>
          <Select
            value={mcpConnectionsPagination.elementsPerPage.toString()}
            onValueChange={onElementsPerPageChange}
          >
            <SelectTrigger className="w-20">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              {[10, 25, 50, 100].map((value) => (
                <SelectItem key={value} value={value.toString()}>
                  {value}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
      </div>
    </div>
  );
}
