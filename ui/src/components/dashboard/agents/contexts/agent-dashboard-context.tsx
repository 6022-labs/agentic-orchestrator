/* eslint-disable react-refresh/only-export-components */
import React from "react";
import { Agent } from "@/models/agent";
import { Pagination } from "@/models/pagination";
import { useAuth } from "@/contexts/auth-context";
import { listAgents, ListAgentsResponse } from "@/queries/api/agent.queries";

type AgentDashboardContextType = {
  fetchInternalAgents: (
    page?: number,
    elementsPerPage?: number,
    sortBy?: string,
    sortDirection?: "asc" | "desc"
  ) => Promise<void>;
  fetchExternalAgents: (
    page?: number,
    elementsPerPage?: number,
    sortBy?: string,
    sortDirection?: "asc" | "desc"
  ) => Promise<void>;
  internalAgents: Agent[];
  externalAgents: Agent[];
  internalPagination: Pagination;
  externalPagination: Pagination;
};

export const AgentDashboardContext =
  React.createContext<AgentDashboardContextType>({
    fetchInternalAgents: async () => {},
    fetchExternalAgents: async () => {},
    internalAgents: [],
    externalAgents: [],
    internalPagination: { page: 1, elementsPerPage: 10, totalElements: 0 },
    externalPagination: { page: 1, elementsPerPage: 10, totalElements: 0 },
  });

export const AgentDashboardContextProvider: React.FC<{
  children: React.ReactNode;
}> = ({ children }) => {
  const { jwt } = useAuth();
  const [internalAgents, setInternalAgents] = React.useState<Agent[]>([]);
  const [externalAgents, setExternalAgents] = React.useState<Agent[]>([]);

  const [internalPagination, setInternalPagination] =
    React.useState<Pagination>({
      page: 1,
      elementsPerPage: 10,
      totalElements: 0,
    });

  const [externalPagination, setExternalPagination] =
    React.useState<Pagination>({
      page: 1,
      elementsPerPage: 10,
      totalElements: 0,
    });

  const fetchInternalAgents = async (
    page = internalPagination.page,
    elementsPerPage = internalPagination.elementsPerPage,
    sortBy?: string,
    sortDirection?: "asc" | "desc"
  ) => {
    if (!jwt) return;

    const response = await listAgents(jwt, {
      page,
      elementsPerPage,
      isInternal: true,
      sortBy,
      sortDirection,
    });

    if (response.status > 299) {
      throw new Error(
        `Failed to fetch internal agents with status ${response.status}`
      );
    }

    const data: ListAgentsResponse = response.data;
    setInternalAgents(data.agents as Agent[]);
    setInternalPagination({
      page: data.page,
      elementsPerPage: data.elementsPerPage,
      totalElements: data.totalElements,
    });
  };

  const fetchExternalAgents = async (
    page = externalPagination.page,
    elementsPerPage = externalPagination.elementsPerPage,
    sortBy?: string,
    sortDirection?: "asc" | "desc"
  ) => {
    if (!jwt) return;

    const response = await listAgents(jwt, {
      page,
      elementsPerPage,
      isInternal: false,
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
    setInternalAgents([]);
    setExternalAgents([]);
    setInternalPagination({ page: 1, elementsPerPage: 10, totalElements: 0 });
    setExternalPagination({ page: 1, elementsPerPage: 10, totalElements: 0 });
  }, [jwt]);

  return (
    <AgentDashboardContext.Provider
      value={{
        fetchInternalAgents,
        fetchExternalAgents,
        internalAgents,
        externalAgents,
        internalPagination,
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
