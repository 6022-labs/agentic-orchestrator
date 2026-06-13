/* eslint-disable react-refresh/only-export-components */
import React from "react";
import { toast } from "sonner";
import { Agent } from "@/models/agent";
import { Pagination } from "@/models/pagination";
import {
  listAgentContextKeys,
  ListAgentContextKeysResponse,
} from "@/queries/api/agent.queries";

type AgentContextsDialogContextType = {
  isOpen: boolean;
  setIsOpen: (isOpen: boolean) => void;
  agent: Agent | undefined;
  setAgent: (agent: Agent | undefined) => void;
  contextKeys: string[];
  isLoading: boolean;
  fetchContextKeys: (
    page?: number,
    search?: string,
    elementsPerPage?: number,
    sortBy?: string,
    sortDirection?: "asc" | "desc",
  ) => Promise<void>;
  contextKeysPagination: Pagination;
};

export const AgentContextsDialogContext =
  React.createContext<AgentContextsDialogContextType>({
    isOpen: false,
    agent: undefined,
    setAgent: async () => {},
    setIsOpen: async () => {},
    contextKeys: [],
    isLoading: false,
    fetchContextKeys: async () => {},
    contextKeysPagination: { page: 1, elementsPerPage: 10, totalElements: 0 },
  });

export const AgentContextsDialogContextProvider: React.FC<{
  children: React.ReactNode;
}> = ({ children }) => {
  const [isOpen, setIsOpen] = React.useState(false);
  const [agent, setAgent] = React.useState<Agent | undefined>(undefined);
  const [contextKeys, setContextKeys] = React.useState<string[]>([]);
  const [isLoading, setIsLoading] = React.useState(false);
  const [contextKeysPagination, setContextKeysPagination] =
    React.useState<Pagination>({
      page: 1,
      elementsPerPage: 10,
      totalElements: 0,
    });

  const fetchContextKeys = async (
    page?: number,
    search?: string,
    elementsPerPage?: number,
    sortBy?: string,
    sortDirection?: "asc" | "desc",
  ) => {
    if (!agent) return;

    setIsLoading(true);

    try {
      // Fetch the latest context keys for the agent
      const response = await listAgentContextKeys(
        agent.collectionAddress,
        agent.collectionTokenId!,
        {
          page,
          search,
          elementsPerPage,
          sortBy,
          sortDirection,
        },
      );

      if (response.status > 299) {
        throw new Error(
          `Failed to fetch context keys for agent ${agent.id} with status ${response.status}`,
        );
      }

      const data: ListAgentContextKeysResponse = response.data;

      setContextKeys(data.contextKeys);
      setContextKeysPagination({
        page: data.page,
        totalElements: data.totalElements,
        elementsPerPage: data.elementsPerPage,
      });
    } catch (error) {
      console.error("Failed to load context keys:", error);
      toast.error("Failed to load context keys");
    } finally {
      setIsLoading(false);
    }
  };

  const handleSetAgent = async (agent: Agent | undefined) => {
    if (!agent) {
      setAgent(undefined);
      setContextKeys([]);
      return;
    }

    setAgent(agent);
  };

  React.useEffect(() => {
    if (isOpen && agent) {
      fetchContextKeys();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [isOpen, agent]);

  return (
    <AgentContextsDialogContext.Provider
      value={{
        agent,
        isOpen,
        setAgent: handleSetAgent,
        setIsOpen,
        contextKeys,
        isLoading,
        fetchContextKeys,
        contextKeysPagination,
      }}
    >
      {children}
    </AgentContextsDialogContext.Provider>
  );
};

export const useAgentContextsDialogContext = () => {
  const context = React.useContext(AgentContextsDialogContext);

  if (!context) {
    throw new Error(
      "useAgentContextsDialogContext must be used within an AgentContextsDialogContextProvider",
    );
  }

  return context;
};
