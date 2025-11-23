/* eslint-disable react-refresh/only-export-components */
import React from "react";
import { Agent } from "@/models/agent";

type TransferAgentDialogContextType = {
  isOpen: boolean;
  setIsOpen: (isOpen: boolean) => void;
  agent: Agent | undefined;
  setAgent: (agent: Agent | undefined) => void;
};

export const TransferAgentDialogContext =
  React.createContext<TransferAgentDialogContextType>({
    isOpen: false,
    agent: undefined,
    setAgent: async () => {},
    setIsOpen: async () => {},
  });

export const TransferAgentDialogContextProvider: React.FC<{
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
    <TransferAgentDialogContext.Provider
      value={{
        agent,
        isOpen,
        setIsOpen,
        setAgent: handleSetAgent,
      }}
    >
      {children}
    </TransferAgentDialogContext.Provider>
  );
};

export const useTransferAgentDialogContext = () => {
  const context = React.useContext(TransferAgentDialogContext);

  if (!context) {
    throw new Error(
      "useTransferAgentDialogContext must be used within a TransferAgentDialogContextProvider"
    );
  }

  return context;
};
