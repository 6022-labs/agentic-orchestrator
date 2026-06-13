/* eslint-disable react-refresh/only-export-components */
import React from "react";
import { Eip1193Provider } from "ethers";
import { useAuth } from "@/contexts/auth-context";
import { useAppKitProvider } from "@reown/appkit/react";
import { canCreateCollection } from "@/queries/contracts/agent-collection-creator";

type CreateAgentCollectionDialogContextType = {
  isOpen: boolean;
  canCreateCollection: boolean;
  setIsOpen: (isOpen: boolean) => void;
};

export const CreateAgentCollectionDialogContext =
  React.createContext<CreateAgentCollectionDialogContextType>({
    isOpen: false,
    canCreateCollection: false,
    setIsOpen: async () => {},
  });

type CreateAgentCollectionDialogContextProviderProps = {
  children: React.ReactNode;
};

export const CreateAgentCollectionDialogContextProvider: React.FC<
  CreateAgentCollectionDialogContextProviderProps
> = ({ children }) => {
  const [isOpen, setIsOpen] = React.useState(false);
  const [authUserCanCreateCollection, setAuthUserCanCreateCollection] =
    React.useState(false);

  const { user } = useAuth();

  const { walletProvider } = useAppKitProvider("eip155");
  const eip1193Provider = walletProvider as Eip1193Provider;

  const checkCanCreateCollection = async () => {
    if (!user || !user.address || !eip1193Provider) {
      setAuthUserCanCreateCollection(false);
      return;
    }

    const canCreate = await canCreateCollection(user.address, eip1193Provider);
    setAuthUserCanCreateCollection(canCreate);
  };

  React.useEffect(() => {
    checkCanCreateCollection();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [user, eip1193Provider]);

  return (
    <CreateAgentCollectionDialogContext.Provider
      value={{
        isOpen,
        setIsOpen,
        canCreateCollection: authUserCanCreateCollection,
      }}
    >
      {children}
    </CreateAgentCollectionDialogContext.Provider>
  );
};

export const useCreateAgentCollectionDialog = () => {
  const context = React.useContext(CreateAgentCollectionDialogContext);

  if (!context) {
    throw new Error(
      "useCreateAgentCollectionDialog must be used within a CreateAgentCollectionDialogContextProvider"
    );
  }

  return context;
};
