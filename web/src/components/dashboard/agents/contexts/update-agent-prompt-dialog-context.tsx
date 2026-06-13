/* eslint-disable react-refresh/only-export-components */
import React from "react";
import { toast } from "sonner";
import { Agent } from "@/models/agent";
import { getAgentPromptTemplate } from "@/queries/api/agent.queries";
import {
  AgentPromptTemplateKey,
  GetAgentPromptTemplateKeysFromRole,
} from "@/models/enums/agent-prompt-template-key";

type UpdateAgentPromptDialogContextType = {
  isOpen: boolean;
  setIsOpen: (isOpen: boolean) => void;
  agent: Agent | undefined;
  isPromptTemplateLoading: boolean;
  setAgent: (agent: Agent | undefined) => void;
  existingPromptsTemplates: Map<string, string>;
  agentPromptTemplateKeys: AgentPromptTemplateKey[];
  reloadExistingPromptTemplate: (key: AgentPromptTemplateKey) => Promise<void>;
};

export const UpdateAgentPromptDialogContext =
  React.createContext<UpdateAgentPromptDialogContextType>({
    isOpen: false,
    agent: undefined,
    setAgent: async () => {},
    setIsOpen: async () => {},
    agentPromptTemplateKeys: [],
    isPromptTemplateLoading: false,
    reloadExistingPromptTemplate: async () => {},
    existingPromptsTemplates: new Map<string, string>(),
  });

export const UpdateAgentPromptDialogContextProvider: React.FC<{
  children: React.ReactNode;
}> = ({ children }) => {

  const [isOpen, setIsOpen] = React.useState(false);
  const [agent, setAgent] = React.useState<Agent | undefined>(undefined);

  const [isPromptTemplateLoading, setIsPromptTemplateLoading] =
    React.useState(false);

  const [agentPromptTemplateKeys, setAgentPromptTemplateKeys] = React.useState<
    AgentPromptTemplateKey[]
  >([]);
  const [existingPromptsTemplates, setExistingPromptsTemplates] =
    React.useState<Map<string, string>>(new Map<string, string>());

  const handleSetAgent = async (agent: Agent | undefined) => {
    if (!agent) {
      setAgent(undefined);
      setAgentPromptTemplateKeys([]);
      return;
    }

    setAgent(agent);
    setIsPromptTemplateLoading(true);

    const agentPromptTemplateKeys = GetAgentPromptTemplateKeysFromRole(
      agent.role
    );

    setAgentPromptTemplateKeys(agentPromptTemplateKeys);

    await loadExistingPromptTemplates(agent, agentPromptTemplateKeys);
  };

  const loadExistingPromptTemplates = async (
    agent: Agent | undefined,
    keys: AgentPromptTemplateKey[]
  ) => {
    if (!agent) return;

    try {
      const loadExistingPromptPromises = [];

      for (const key of keys) {
        loadExistingPromptPromises.push(loadExistingPromptTemplate(agent, key));
      }

      await Promise.all(loadExistingPromptPromises);
    } catch (error) {
      console.error("Error fetching agent prompt templates:", error);
      toast.error("Failed to fetch agent prompt templates");
    } finally {
      setIsPromptTemplateLoading(false);
    }
  };

  const loadExistingPromptTemplate = async (
    agent: Agent | undefined,
    key: AgentPromptTemplateKey
  ) => {
    if (!agent) return;
    try {
      const response = await getAgentPromptTemplate(
        agent.collectionAddress,
        agent.collectionTokenId!,
        key
      );

      if (response.status > 299) {
        throw new Error(
          `Failed to fetch agent prompt template for key: ${key}, status: ${response.status}`
        );
      }

      setExistingPromptsTemplates((oldExistingPrompts) => {
        const newExistingPrompts = new Map(oldExistingPrompts);
        newExistingPrompts.set(key, response.data);
        return newExistingPrompts;
      });
    } catch (error) {
      console.error(
        "Error fetching agent prompt template key",
        key,
        ":",
        error
      );
      toast.error("Failed to fetch agent prompt template key: " + key);
    }
  };

  const reloadExistingPromptTemplate = async (key: AgentPromptTemplateKey) => {
    try {
      setIsPromptTemplateLoading(true);
      await loadExistingPromptTemplate(agent, key);
    } catch (error) {
      console.error(
        "Error reloading agent prompt template key:",
        key,
        ":",
        error
      );
      toast.error("Failed to reload agent prompt template key: " + key);
    } finally {
      setIsPromptTemplateLoading(false);
    }
  };

  return (
    <UpdateAgentPromptDialogContext.Provider
      value={{
        agent,
        isOpen,
        setIsOpen,
        isPromptTemplateLoading,
        agentPromptTemplateKeys,
        existingPromptsTemplates,
        setAgent: handleSetAgent,
        reloadExistingPromptTemplate,
      }}
    >
      {children}
    </UpdateAgentPromptDialogContext.Provider>
  );
};

export const useUpdateAgentPromptDialogContext = () => {
  const context = React.useContext(UpdateAgentPromptDialogContext);

  if (!context) {
    throw new Error(
      "useUpdateAgentPromptDialogContext must be used within an UpdateAgentPromptDialogContextProvider"
    );
  }

  return context;
};
