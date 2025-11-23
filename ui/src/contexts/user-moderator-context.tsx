/* eslint-disable react-refresh/only-export-components */
import React from "react";
import { toast } from "sonner";
import { Eip1193Provider } from "ethers";
import { useAuth } from "./auth-context";
import { useAppKitProvider } from "@reown/appkit/react";
import { AgentCollection } from "@/models/agent-collection";
import { useAgentCollections } from "./agent-collections-context";
import { moderatorOf } from "@/queries/contracts/agent-collections-manager";

type UserModeratorContextType = {
  collections: AgentCollection[];
};

export const UserModeratorContext =
  React.createContext<UserModeratorContextType>({
    collections: [],
  } as UserModeratorContextType);

export const UserModeratorProvider: React.FC<{
  children: React.ReactNode;
}> = ({ children }) => {
  const { jwt } = useAuth();
  const { collections } = useAgentCollections();
  const { walletProvider } = useAppKitProvider("eip155");
  const eip1193Provider = walletProvider as Eip1193Provider;

  const [userModeratorCollections, setModeratorCollections] = React.useState<
    AgentCollection[]
  >([]);

  const loadCollections = async () => {
    if (!jwt || !eip1193Provider || !collections || collections.length === 0)
      return;

    try {
      const addresses = await moderatorOf(eip1193Provider);

      if (!addresses || addresses.length === 0) {
        console.info("No moderator addresses found");
        setModeratorCollections([]);
        return;
      }

      const filteredCollections = collections.filter((collection) =>
        addresses
          .map((address) => address.toLowerCase())
          .includes(collection.address.toLowerCase())
      );

      if (filteredCollections.length === 0) {
        console.info("No collections found for the moderator addresses");
        setModeratorCollections([]);
        return;
      }

      setModeratorCollections(filteredCollections);
    } catch (error) {
      console.error("Failed to load moderator collections:", error);
      toast.error("Failed to load moderator collections");
      setModeratorCollections([]);
    }
  };

  React.useEffect(() => {
    if (jwt) {
      loadCollections();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [jwt, collections, walletProvider]);

  return (
    <UserModeratorContext.Provider
      value={{ collections: userModeratorCollections }}
    >
      {children}
    </UserModeratorContext.Provider>
  );
};

export const useUserModerator = () => {
  const context = React.useContext(UserModeratorContext);
  if (!context)
    throw new Error(
      "useUserModerator must be used within a UserModeratorProvider"
    );
  return context;
};
