/* eslint-disable react-refresh/only-export-components */
import React from "react";
import { Agent } from "@/models/agent";

type UpsertAgentAttributeDialogContextType = {
  isOpen: boolean;
  setIsOpen: (isOpen: boolean) => void;
  agent: Agent | undefined;
  setAgent: (agent: Agent | undefined) => void;
};

export const UpsertAgentAttributeDialogContext =
  React.createContext<UpsertAgentAttributeDialogContextType>({
    isOpen: false,
    agent: undefined,
    setAgent: async () => {},
    setIsOpen: async () => {},
  });

export const UpsertAgentAttributeDialogContextProvider: React.FC<{
  children: React.ReactNode;
}> = ({ children }) => {
  const [isOpen, setIsOpen] = React.useState(false);
  const [agent, setAgent] = React.useState<Agent | undefined>(undefined);

  const handleSetAgent = async (agent: Agent | undefined) => {
    if (!agent) {
      setAgent(undefined);
      return;
    }

    setAgent(agent);
  };

  return (
    <UpsertAgentAttributeDialogContext.Provider
      value={{
        agent,
        isOpen,
        setIsOpen,
        setAgent: handleSetAgent,
      }}
    >
      {children}
    </UpsertAgentAttributeDialogContext.Provider>
  );
};

export const useUpsertAgentAttributeDialogContext = () => {
  const context = React.useContext(UpsertAgentAttributeDialogContext);

  if (!context) {
    throw new Error(
      "useUpsertAgentAttributeDialogContext must be used within a UpsertAgentAttributeDialogContextProvider"
    );
  }

  return context;
};
