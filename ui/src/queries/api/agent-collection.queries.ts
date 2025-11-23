/* eslint-disable @typescript-eslint/no-explicit-any */
import axios, { AxiosResponse } from "axios";
import { AgentCollection } from "@/models/agent-collection";

export const getAgentCollections = async (
  jwt: string
): Promise<AxiosResponse<AgentCollection[], any>> => {
  return await axios.get(`/agent-collections`, {
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${jwt}`,
    },
  });
};
