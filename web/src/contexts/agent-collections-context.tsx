/* eslint-disable react-refresh/only-export-components */
import React from "react";
import { toast } from "sonner";
import { useAuth } from "./auth-context";
import { AgentCollection } from "@/models/agent-collection";
import { getAgentCollections } from "@/queries/api/agent-collection.queries";

type AgentCollectionsContextType = {
  collections: AgentCollection[];
  loadCollections: () => Promise<void>;
};

export const AgentCollectionsContext =
  React.createContext<AgentCollectionsContextType>({
    collections: [],
    loadCollections: () => Promise.resolve(),
  } as AgentCollectionsContextType);

export const AgentCollectionsProvider: React.FC<{
  children: React.ReactNode;
}> = ({ children }) => {
  const { jwt } = useAuth();
  const [collections, setCollections] = React.useState<AgentCollection[]>([]);

  const loadCollections = async () => {
    if (!jwt) return;
    try {
      const response = await getAgentCollections();
      setCollections(response.data);
    } catch (error) {
      console.error("Failed to load agent collections:", error);
      toast.error("Failed to load agent collections");
      setCollections([]);
    }
  };

  React.useEffect(() => {
    if (jwt) {
      loadCollections();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [jwt]);

  return (
    <AgentCollectionsContext.Provider value={{ collections, loadCollections }}>
      {children}
    </AgentCollectionsContext.Provider>
  );
};

export const useAgentCollections = () => {
  const context = React.useContext(AgentCollectionsContext);
  if (!context)
    throw new Error(
      "useAgentCollections must be used within an AgentCollectionsProvider"
    );
  return context;
};
