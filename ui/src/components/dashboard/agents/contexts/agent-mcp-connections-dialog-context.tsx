/* eslint-disable react-refresh/only-export-components */
import React from "react";
import { toast } from "sonner";
import { Agent } from "@/models/agent";
import { Pagination } from "@/models/pagination";
import { useAuth } from "@/contexts/auth-context";
import { MCPConnection } from "@/models/mcp-connection";
import {
  listAgentMCPConnections,
  ListAgentMCPConnectionsResponse,
} from "@/queries/api/agent.queries";

type AgentMCPConnectionsDialogContextType = {
  isOpen: boolean;
  setIsOpen: (isOpen: boolean) => void;
  agent: Agent | undefined;
  setAgent: (agent: Agent | undefined) => Promise<void>;
  mcpConnections: MCPConnection[];
  isLoading: boolean;
  fetchMCPConnections: (
    page?: number,
    search?: string,
    elementsPerPage?: number,
    sortBy?: string,
    sortDirection?: "asc" | "desc"
  ) => Promise<void>;
  mcpConnectionsPagination: Pagination;
};

export const AgentMCPConnectionsDialogContext =
  React.createContext<AgentMCPConnectionsDialogContextType>({
    isOpen: false,
    setIsOpen: () => {},
    agent: undefined,
    setAgent: async () => {},
    mcpConnections: [],
    isLoading: false,
    fetchMCPConnections: async () => {},
    mcpConnectionsPagination: {
      page: 1,
      elementsPerPage: 10,
      totalElements: 0,
    },
  });

export const AgentMCPConnectionsDialogContextProvider: React.FC<{
  children: React.ReactNode;
}> = ({ children }) => {
  const { jwt } = useAuth();

  const [isOpen, setIsOpen] = React.useState(false);
  const [agent, setAgent] = React.useState<Agent | undefined>(undefined);
  const [mcpConnections, setMCPConnections] = React.useState<MCPConnection[]>(
    []
  );
  const [isLoading, setIsLoading] = React.useState(false);
  const [mcpConnectionsPagination, setMCPConnectionsPagination] =
    React.useState<Pagination>({
      page: 1,
      elementsPerPage: 10,
      totalElements: 0,
    });

  const fetchMCPConnections = async (
    page?: number,
    search?: string,
    elementsPerPage?: number,
    sortBy?: string,
    sortDirection?: "asc" | "desc"
  ) => {
    if (!jwt || !agent) return;

    setIsLoading(true);

    try {
      // Fetch the latest MCP connections for the agent
      const response = await listAgentMCPConnections(
        jwt!,
        agent.collectionAddress,
        agent.collectionTokenId!,
        {
          page,
          search,
          elementsPerPage,
          sortBy,
          sortDirection,
        }
      );

      if (response.status > 299) {
        throw new Error(
          `Failed to fetch MCP connections for agent ${agent.id} with status ${response.status}`
        );
      }

      const data: ListAgentMCPConnectionsResponse = response.data;

      setMCPConnections(data.mcpConnections);
      setMCPConnectionsPagination({
        page: data.page,
        totalElements: data.totalElements,
        elementsPerPage: data.elementsPerPage,
      });
    } catch (error) {
      console.error("Failed to load MCP connections:", error);
      toast.error("Failed to load MCP connections");
      setMCPConnections([]);
    } finally {
      setIsLoading(false);
    }
  };

  const handleSetAgent = async (agent: Agent | undefined) => {
    if (!agent) {
      setAgent(undefined);
      setMCPConnections([]);
      return;
    }

    setAgent(agent);
    await fetchMCPConnections();
  };

  return (
    <AgentMCPConnectionsDialogContext.Provider
      value={{
        isOpen,
        setIsOpen,
        agent,
        setAgent: handleSetAgent,
        mcpConnections,
        isLoading,
        fetchMCPConnections,
        mcpConnectionsPagination,
      }}
    >
      {children}
    </AgentMCPConnectionsDialogContext.Provider>
  );
};

export const useAgentMCPConnectionsDialogContext = () => {
  const context = React.useContext(AgentMCPConnectionsDialogContext);
  if (!context) {
    throw new Error(
      "useAgentMCPConnectionsDialogContext must be used within AgentMCPConnectionsDialogContextProvider"
    );
  }
  return context;
};
