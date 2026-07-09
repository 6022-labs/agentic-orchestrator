/* eslint-disable react-refresh/only-export-components */
import React from "react";
import { Agent } from "@/models/agent";
import { Pagination } from "@/models/pagination";
import { listAgents, ListAgentsResponse } from "@/queries/api/agent.queries";

type AgentDashboardContextType = {
  fetchExternalAgents: (
    page?: number,
    elementsPerPage?: number,
    sortBy?: string,
    sortDirection?: "asc" | "desc"
  ) => Promise<void>;
  externalAgents: Agent[];
  externalPagination: Pagination;
};

export const AgentDashboardContext =
  React.createContext<AgentDashboardContextType>({
    fetchExternalAgents: async () => {},
    externalAgents: [],
    externalPagination: { page: 1, elementsPerPage: 10, totalElements: 0 },
  });

export const AgentDashboardContextProvider: React.FC<{
  children: React.ReactNode;
}> = ({ children }) => {
  const [externalAgents, setExternalAgents] = React.useState<Agent[]>([]);

  const [externalPagination, setExternalPagination] =
    React.useState<Pagination>({
      page: 1,
      elementsPerPage: 10,
      totalElements: 0,
    });

  const fetchExternalAgents = async (
    page = externalPagination.page,
    elementsPerPage = externalPagination.elementsPerPage,
    sortBy?: string,
    sortDirection?: "asc" | "desc"
  ) => {
    const response = await listAgents({
      page,
      elementsPerPage,
      sortBy,
      sortDirection,
    });

    if (response.status > 299) {
      throw new Error(
        `Failed to fetch external agents with status ${response.status}`
      );
    }

    const data: ListAgentsResponse = response.data;
    setExternalAgents(data.agents as Agent[]);
    setExternalPagination({
      page: data.page,
      elementsPerPage: data.elementsPerPage,
      totalElements: data.totalElements,
    });
  };

  React.useEffect(() => {
    setExternalAgents([]);
    setExternalPagination({ page: 1, elementsPerPage: 10, totalElements: 0 });
  }, []);

  return (
    <AgentDashboardContext.Provider
      value={{
        fetchExternalAgents,
        externalAgents,
        externalPagination,
      }}
    >
      {children}
    </AgentDashboardContext.Provider>
  );
};

export const useAgentDashboard = () => {
  const context = React.useContext(AgentDashboardContext);

  if (!context) {
    throw new Error(
      "useAgentDashboard must be used within an AgentDashboardContextProvider"
    );
  }

  return context;
};
