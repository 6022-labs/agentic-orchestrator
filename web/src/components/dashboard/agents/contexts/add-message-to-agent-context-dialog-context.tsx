/* eslint-disable react-refresh/only-export-components */
import React from "react";

type AddMessageToAgentContextDialogContextType = {
  isOpen: boolean;
  setIsOpen: (isOpen: boolean) => void;
  context: string;
  setContext: (context: string) => void;
};

export const AddMessageToAgentContextDialogContext =
  React.createContext<AddMessageToAgentContextDialogContextType>({
    isOpen: false,
    setIsOpen: async () => {},
    context: "",
    setContext: async () => {},
  });

export const AddMessageToAgentContextDialogContextProvider: React.FC<{
  children: React.ReactNode;
}> = ({ children }) => {
  const [isOpen, setIsOpen] = React.useState(false);
  const [context, setContext] = React.useState<string>("");

  return (
    <AddMessageToAgentContextDialogContext.Provider
      value={{
        isOpen,
        context,
        setIsOpen,
        setContext,
      }}
    >
      {children}
    </AddMessageToAgentContextDialogContext.Provider>
  );
};

export const useAddMessageToAgentContextDialogContext = () => {
  const context = React.useContext(AddMessageToAgentContextDialogContext);

  if (!context) {
    throw new Error(
      "useAddMessageToAgentContextDialogContext must be used within an AddMessageToAgentContextDialogContextProvider"
    );
  }

  return context;
};
