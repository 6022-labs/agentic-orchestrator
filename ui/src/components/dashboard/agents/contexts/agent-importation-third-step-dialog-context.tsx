/* eslint-disable react-refresh/only-export-components */
import React from "react";
import { Agent } from "@/models/agent";

type AgentImportationThirdStepDialogContextType = {
  isOpen: boolean;
  setIsOpen: (isOpen: boolean) => void;
  agent: Agent | undefined;
  setAgent: (agent: Agent) => void;
};

export const AgentImportationThirdStepDialogContext =
  React.createContext<AgentImportationThirdStepDialogContextType>({
    isOpen: false,
    agent: undefined,
    setAgent: async () => {},
    setIsOpen: async () => {},
  });

export const AgentImportationThirdStepDialogContextProvider: React.FC<{
  children: React.ReactNode;
}> = ({ children }) => {
  const [isOpen, setIsOpen] = React.useState(false);
  const [agent, setAgent] = React.useState<Agent | undefined>(undefined);

  return (
    <AgentImportationThirdStepDialogContext.Provider
      value={{
        agent,
        isOpen,
        setAgent,
        setIsOpen,
      }}
    >
      {children}
    </AgentImportationThirdStepDialogContext.Provider>
  );
};

export const useAgentImportationThirdStepDialog = () => {
  const context = React.useContext(AgentImportationThirdStepDialogContext);

  if (!context) {
    throw new Error(
      "useAgentImportationThirdStepDialog must be used within an AgentImportationThirdStepDialogContextProvider"
    );
  }

  return context;
};
