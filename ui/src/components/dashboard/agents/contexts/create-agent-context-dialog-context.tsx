/* eslint-disable react-refresh/only-export-components */
import React from "react";

type CreateAgentContextDialogContextType = {
  isOpen: boolean;
  setIsOpen: (isOpen: boolean) => void;
};

export const CreateAgentContextDialogContext =
  React.createContext<CreateAgentContextDialogContextType>({
    isOpen: false,
    setIsOpen: async () => {},
  });

export const CreateAgentContextDialogContextProvider: React.FC<{
  children: React.ReactNode;
}> = ({ children }) => {
  const [isOpen, setIsOpen] = React.useState(false);

  return (
    <CreateAgentContextDialogContext.Provider
      value={{
        isOpen,
        setIsOpen,
      }}
    >
      {children}
    </CreateAgentContextDialogContext.Provider>
  );
};

export const useCreateAgentContextDialogContext = () => {
  const context = React.useContext(CreateAgentContextDialogContext);

  if (!context) {
    throw new Error(
      "useCreateAgentContextDialogContext must be used within a CreateAgentContextDialogContextProvider"
    );
  }

  return context;
};
