/* eslint-disable react-refresh/only-export-components */
import React from "react";
import { Agent } from "@/models/agent";
import { getAgentLLMSettings } from "@/queries/api/agent.queries";
import { AgentBaseLLMSettings } from "@/models/agent-base-llm-settings";

type SetupAgentLLMSettingsDialogContextType = {
  isOpen: boolean;
  setIsOpen: (isOpen: boolean) => void;
  agent: Agent | undefined;
  setAgent: (agent: Agent | undefined) => void;
  isExistingLLMSettingsLoading: boolean;
  existingLLMSettings: AgentBaseLLMSettings | undefined;
};

export const SetupAgentLLMSettingsDialogContext =
  React.createContext<SetupAgentLLMSettingsDialogContextType>({
    isOpen: false,
    agent: undefined,
    setAgent: async () => {},
    setIsOpen: async () => {},
    existingLLMSettings: undefined,
    isExistingLLMSettingsLoading: false,
  });

export const SetupAgentLLMSettingsDialogContextProvider: React.FC<{
  children: React.ReactNode;
}> = ({ children }) => {
  const [isOpen, setIsOpen] = React.useState(false);
  const [agent, setAgent] = React.useState<Agent | undefined>(undefined);

  const [isExistingLLMSettingsLoading, setIsExistingLLMSettingsLoading] =
    React.useState(false);

  const [existingLLMSettings, setExistingLLMSettings] = React.useState<
    AgentBaseLLMSettings | undefined
  >(undefined);

  const handleSetAgent = async (agent: Agent | undefined) => {
    if (!agent) {
      setAgent(undefined);
      setExistingLLMSettings(undefined);
      return;
    }

    setAgent(agent);
    setIsExistingLLMSettingsLoading(true);

    try {
      const response = await getAgentLLMSettings(
        agent.collectionAddress,
        agent.collectionTokenId!
      );

      if (response.status !== 200) {
        throw new Error("Failed to fetch agent LLM settings");
      }

      setExistingLLMSettings(response.data);
    } catch (error) {
      console.error("Error fetching agent LLM settings:", error);
    } finally {
      setIsExistingLLMSettingsLoading(false);
    }
  };

  return (
    <SetupAgentLLMSettingsDialogContext.Provider
      value={{
        agent,
        isOpen,
        setIsOpen,
        existingLLMSettings,
        setAgent: handleSetAgent,
        isExistingLLMSettingsLoading,
      }}
    >
      {children}
    </SetupAgentLLMSettingsDialogContext.Provider>
  );
};

export const useSetupAgentLLMSettingsDialogContext = () => {
  const context = React.useContext(SetupAgentLLMSettingsDialogContext);

  if (!context) {
    throw new Error(
      "useSetupAgentLLMSettingsDialogContext must be used within an SetupAgentLLMSettingsDialogContextProvider"
    );
  }

  return context;
};
