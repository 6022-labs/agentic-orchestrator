/* eslint-disable @typescript-eslint/no-explicit-any */
import { Agent } from "@/models/agent";
import axios, { AxiosResponse } from "axios";

export type ListAgentsResponse = {
  page: number;
  elementsPerPage: number;
  totalElements: number;
  agents: Agent[];
};

export const listAgents = async (params: {
  page?: number;
  elementsPerPage?: number;
  sortBy?: string;
  sortDirection?: "asc" | "desc";
}): Promise<AxiosResponse<ListAgentsResponse, any>> => {
  return await axios.get(`/agents`, { params });
};
