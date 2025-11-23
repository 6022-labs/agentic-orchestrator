/* eslint-disable @typescript-eslint/no-explicit-any */
import { Agent } from "@/models/agent";
import axios, { AxiosResponse } from "axios";
import { AgentPromptTemplateKey } from "@/models/enums/agent-prompt-template-key";
import { AgentContextMessageRole } from "@/models/enums/agent-context-message-role";
import { CreateAgentFirstStepFormSchemaType } from "@/components/dashboard/agents/forms/create-agent-first-step-form/create-agent-first-step-form-schema";
import { SetupAgentLLMSettingsFormSchemaType } from "@/components/dashboard/agents/forms/setup-agent-llm-settings-form/setup-agent-llm-settings-form-schema";
import { CreateAgentMCPConnectionFormSchemaType } from "@/components/dashboard/agents/forms/create-agent-mcp-connection-form/create-agent-mcp-connection-form-schema";
import { SetupAgentDatabaseSettingsFirstStepFormSchemaType } from "@/components/dashboard/agents/forms/setup-agent-database-settings-first-step-form/setup-agent-database-settings-first-step-form-schema";
import { MCPConnection } from "@/models/mcp-connection";

export type ListAgentsResponse = {
  page: number;
  elementsPerPage: number;
  totalElements: number;
  agents: Agent[];
};

export const listAgents = async (
  jwt: string,
  params: {
    page?: number;
    elementsPerPage?: number;
    isInternal?: boolean;
    sortBy?: string;
    sortDirection?: "asc" | "desc";
  }
): Promise<AxiosResponse<ListAgentsResponse, any>> => {
  return await axios.get(`/agents`, {
    params,
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${jwt}`,
    },
  });
};

export const createAgentFirstStep = async (
  jwt: string,
  body: CreateAgentFirstStepFormSchemaType
) => {
  const formData = new FormData();
  Object.entries(body).forEach(([key, value]) => {
    formData.append(key, value as never);
  });

  return await axios.post(`/agent/create/first-step`, formData, {
    headers: {
      "Content-Type": "multipart/form-data",
      Authorization: `Bearer ${jwt}`,
    },
  });
};

export const createAgentSecondStep = async (jwt: string, agentId: string) => {
  return await axios.post(
    `/agent/create/second-step`,
    { agentId },
    {
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${jwt}`,
      },
    }
  );
};

export const getAgentPromptTemplate = async (
  jwt: string,
  agentCollectionAddress: string,
  agentCollectionTokenId: number,
  key: string
) => {
  return await axios.get(
    `/agent/${agentCollectionAddress}/${agentCollectionTokenId}/prompt-template/${key}`,
    {
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${jwt}`,
      },
    }
  );
};

export const saveAgentPromptTemplate = async (
  jwt: string,
  agentCollectionAddress: string,
  agentCollectionTokenId: number,
  key: AgentPromptTemplateKey,
  content: string
) => {
  return await axios.post(
    `/agent/${agentCollectionAddress}/${agentCollectionTokenId}/prompt-template/${key}`,
    { content },
    {
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${jwt}`,
      },
    }
  );
};

export type ListAgentContextKeysResponse = {
  page: number;
  elementsPerPage: number;
  totalElements: number;
  contextKeys: string[];
};

export const listAgentContextKeys = async (
  jwt: string,
  agentCollectionAddress: string,
  agentCollectionTokenId: number,
  params: {
    page?: number;
    search?: string;
    elementsPerPage?: number;
    sortBy?: string;
    sortDirection?: "asc" | "desc";
  }
): Promise<AxiosResponse<ListAgentContextKeysResponse, any>> => {
  return await axios.get(
    `/agent/${agentCollectionAddress}/${agentCollectionTokenId}/contexts/keys`,
    {
      params,
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${jwt}`,
      },
    }
  );
};

export const createAgentContext = async (
  jwt: string,
  agentCollectionAddress: string,
  agentCollectionTokenId: number,
  key: string
) => {
  return await axios.post(
    `/agent/${agentCollectionAddress}/${agentCollectionTokenId}/context`,
    { key },
    {
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${jwt}`,
      },
    }
  );
};

export const addMessageToContext = async (
  jwt: string,
  agentCollectionAddress: string,
  agentCollectionTokenId: number,
  contextKey: string,
  message: string,
  role: AgentContextMessageRole
) => {
  return await axios.post(
    `/agent/${agentCollectionAddress}/${agentCollectionTokenId}/context/${contextKey}/message`,
    { message, role },
    {
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${jwt}`,
      },
    }
  );
};

export const setupAgentLLMSettings = async (
  jwt: string,
  agentCollectionAddress: string,
  agentCollectionTokenId: number,
  body: SetupAgentLLMSettingsFormSchemaType
) => {
  return await axios.post(
    `/agent/${agentCollectionAddress}/${agentCollectionTokenId}/llm/settings`,
    body,
    {
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${jwt}`,
      },
    }
  );
};

export const setupAgentDatabaseSettingsFirstStep = async (
  jwt: string,
  agentCollectionAddress: string,
  agentCollectionTokenId: number,
  body: SetupAgentDatabaseSettingsFirstStepFormSchemaType
) => {
  return await axios.post(
    `/agent/${agentCollectionAddress}/${agentCollectionTokenId}/database/settings/first-step`,
    body,
    {
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${jwt}`,
      },
    }
  );
};

export const setupAgentDatabaseSettingsSecondStep = async (
  jwt: string,
  agentCollectionAddress: string,
  agentCollectionTokenId: number
) => {
  return await axios.post(
    `/agent/${agentCollectionAddress}/${agentCollectionTokenId}/database/settings/second-step`,
    null,
    {
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${jwt}`,
      },
    }
  );
};

export const getAgentLLMSettings = async (
  jwt: string,
  agentCollectionAddress: string,
  agentCollectionTokenId: number
) => {
  return await axios.get(
    `/agent/${agentCollectionAddress}/${agentCollectionTokenId}/llm/settings`,
    {
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${jwt}`,
      },
    }
  );
};

export const createAgentMCPConnection = async (
  jwt: string,
  agentCollectionAddress: string,
  agentCollectionTokenId: number,
  params: CreateAgentMCPConnectionFormSchemaType
) => {
  // Destructure to remove headerAuthKey and headerAuthValue from params
  const { headerAuthKey, headerAuthValue, ...rest } = params;

  const body: typeof rest & {
    headerAuth: Record<string, string>;
  } = {
    ...rest,
    headerAuth: {},
  };

  if (headerAuthKey && headerAuthValue) {
    body.headerAuth[headerAuthKey] = headerAuthValue;
  }

  return await axios.post(
    `/agent/${agentCollectionAddress}/${agentCollectionTokenId}/mcp/connection`,
    body,
    {
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${jwt}`,
      },
    }
  );
};

export const removeAgentMCPConnection = async (
  jwt: string,
  agentCollectionAddress: string,
  agentCollectionTokenId: number,
  mcpConnectionId: string
) => {
  return await axios.delete(
    `/agent/${agentCollectionAddress}/${agentCollectionTokenId}/mcp/connection/${mcpConnectionId}`,
    {
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${jwt}`,
      },
    }
  );
};

export type ListAgentMCPConnectionsResponse = {
  page: number;
  elementsPerPage: number;
  totalElements: number;
  mcpConnections: MCPConnection[];
};

export const listAgentMCPConnections = async (
  jwt: string,
  agentCollectionAddress: string,
  agentCollectionTokenId: number,
  params: {
    page?: number;
    search?: string;
    elementsPerPage?: number;
    sortBy?: string;
    sortDirection?: "asc" | "desc";
  }
): Promise<AxiosResponse<ListAgentMCPConnectionsResponse, any>> => {
  return await axios.get(
    `/agent/${agentCollectionAddress}/${agentCollectionTokenId}/mcp/connections`,
    {
      params,
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${jwt}`,
      },
    }
  );
};
