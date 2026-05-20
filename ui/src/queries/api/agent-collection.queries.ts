/* eslint-disable @typescript-eslint/no-explicit-any */
import axios, { AxiosResponse } from "axios";
import { AgentCollection } from "@/models/agent-collection";

export const getAgentCollections = async (): Promise<
  AxiosResponse<AgentCollection[], any>
> => {
  return await axios.get(`/agent-collections`);
};
