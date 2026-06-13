/* eslint-disable react-refresh/only-export-components */
import React from "react";

type CreateAgentMCPConnectionDialogContextType = {
  isOpen: boolean;
  setIsOpen: (isOpen: boolean) => void;
};

export const CreateAgentMCPConnectionDialogContext =
  React.createContext<CreateAgentMCPConnectionDialogContextType>({
    isOpen: false,
    setIsOpen: async () => {},
  });

export const CreateAgentMCPConnectionDialogContextProvider: React.FC<{
  children: React.ReactNode;
}> = ({ children }) => {
  const [isOpen, setIsOpen] = React.useState(false);

  return (
    <CreateAgentMCPConnectionDialogContext.Provider
      value={{ isOpen, setIsOpen }}
    >
      {children}
    </CreateAgentMCPConnectionDialogContext.Provider>
  );
};

export const useCreateAgentMCPConnectionDialogContext = () => {
  const context = React.useContext(CreateAgentMCPConnectionDialogContext);
  if (!context) {
    throw new Error(
      "useCreateAgentMCPConnectionDialogContext must be used within CreateAgentMCPConnectionDialogContextProvider"
    );
  }
  return context;
};
