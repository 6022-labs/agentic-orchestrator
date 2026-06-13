/* eslint-disable react-refresh/only-export-components */
import React from "react";
import { Agent } from "@/models/agent";

type SetupAgentDatabaseSettingsFirstStepDialogContextType = {
  isOpen: boolean;
  setIsOpen: (isOpen: boolean) => void;
  agent: Agent | undefined;
  setAgent: (agent: Agent | undefined) => void;
  isCloseButtonHidden: boolean;
  hideCloseButton: () => void;
  displayCloseButton: () => void;
};

export const SetupAgentDatabaseSettingsFirstStepDialogContext =
  React.createContext<SetupAgentDatabaseSettingsFirstStepDialogContextType>({
    isOpen: false,
    agent: undefined,
    setAgent: async () => {},
    setIsOpen: async () => {},
    isCloseButtonHidden: false,
    hideCloseButton: async () => {},
    displayCloseButton: async () => {},
  });

export const SetupAgentDatabaseSettingsFirstStepDialogContextProvider: React.FC<{
  children: React.ReactNode;
}> = ({ children }) => {
  const [isOpen, setIsOpen] = React.useState(false);
  const [agent, setAgent] = React.useState<Agent | undefined>(undefined);
  const [isCloseButtonHidden, setIsCloseButtonHidden] = React.useState(false);

  const hideCloseButton = () => {
    setIsCloseButtonHidden(true);
  };

  const displayCloseButton = () => {
    setIsCloseButtonHidden(false);
  };

  const handleSetAgent = async (agent: Agent | undefined) => {
    if (!agent) {
      setAgent(undefined);
      return;
    }

    setAgent(agent);
  };

  return (
    <SetupAgentDatabaseSettingsFirstStepDialogContext.Provider
      value={{
        agent,
        isOpen,
        setIsOpen,
        setAgent: handleSetAgent,
        hideCloseButton,
        displayCloseButton,
        isCloseButtonHidden,
      }}
    >
      {children}
    </SetupAgentDatabaseSettingsFirstStepDialogContext.Provider>
  );
};

export const useSetupAgentDatabaseSettingsFirstStepDialogContext = () => {
  const context = React.useContext(
    SetupAgentDatabaseSettingsFirstStepDialogContext
  );

  if (!context) {
    throw new Error(
      "useSetupAgentDatabaseSettingsFirstStepDialogContext must be used within an SetupAgentDatabaseSettingsFirstStepDialogContextProvider"
    );
  }

  return context;
};
