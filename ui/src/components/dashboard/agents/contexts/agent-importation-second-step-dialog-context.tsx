// /* eslint-disable react-refresh/only-export-components */
// import React from "react";
// import { toast } from "sonner";
// import { Agent } from "@/models/agent";
// import { useAuth } from "@/contexts/auth-context";
// import { useAgentDashboard } from "./agent-dashboard-context";
// import { DatabaseCredentials } from "@/models/database-credentials";
// import { AgentBaseInformation } from "@/models/agent-base-information";
// import { AgentSetupFormSchemaType } from "../forms/agent-setup-form/agent-setup-form-schema";
// import {
//   importAgentFirstStep,
//   importAgentSecondStep,
// } from "@/queries/api/import-agent.queries";

// type AgentImportationSecondStepDialogContextType = {
//   isOpen: boolean;
//   close: () => void;
//   externalAgent?: Agent;
//   hideCloseButton: () => void;
//   isCloseButtonHidden: boolean;
//   displayCloseButton: () => void;
//   open: (agent: Agent) => Promise<void>;
//   agentBaseInformation?: AgentBaseInformation;
//   onFormComplete: (
//     values: AgentSetupFormSchemaType
//   ) => Promise<DatabaseCredentials | undefined>;
// };

// export const AgentImportationSecondStepDialogContext =
//   React.createContext<AgentImportationSecondStepDialogContextType>({
//     isOpen: false,
//     close: () => {},
//     open: async () => {},
//     agentBaseInformation: undefined,
//     isCloseButtonHidden: false,
//     onFormComplete: async () => {
//       return {} as DatabaseCredentials;
//     },
//     externalAgent: undefined,
//     hideCloseButton: async () => {},
//     displayCloseButton: async () => {},
//   });

// type AgentImportationSecondStepDialogContextProviderProps = {
//   children: React.ReactNode;
// };

// export const AgentImportationSecondStepDialogContextProvider: React.FC<
//   AgentImportationSecondStepDialogContextProviderProps
// > = ({ children }) => {
//   const { jwt } = useAuth();
//   const { fetchAgents } = useAgentDashboard();

//   const [isOpen, setIsOpen] = React.useState(false);
//   const [isCloseButtonHidden, setIsCloseButtonHidden] = React.useState(false);
//   const [agentBaseInformation, setAgentBaseInformation] = React.useState<
//     AgentBaseInformation | undefined
//   >(undefined);
//   const [externalAgent, setExternalAgent] = React.useState<Agent | undefined>(
//     undefined
//   );

//   const hideCloseButton = () => {
//     setIsCloseButtonHidden(true);
//   };

//   const displayCloseButton = () => {
//     setIsCloseButtonHidden(false);
//   };

//   const open = async (agent: Agent) => {
//     try {
//       if (!jwt) {
//         throw new Error("JWT token is null");
//       }
//       if (agent.collectionTokenId === undefined) {
//         throw new Error("Agent does not have a collection token ID");
//       }

//       setExternalAgent(agent);
//       const agentBaseInformation = await importAgentFirstStep(
//         jwt!,
//         agent.collectionAddress,
//         agent.collectionTokenId
//       );
//       if (agentBaseInformation.status !== 200) {
//         throw new Error("Failed to fetch agent metadata");
//       }

//       setAgentBaseInformation(agentBaseInformation.data);
//       setIsOpen(true);
//     } catch (error) {
//       console.error("Error opening agent importation dialog:", error);
//       toast("Failed to open agent importation dialog");
//     }
//   };

//   const close = () => {
//     setIsOpen(false);
//     setExternalAgent(undefined);
//     setAgentBaseInformation(undefined);
//   };

//   const onFormComplete = async (values: AgentSetupFormSchemaType) => {
//     try {
//       const response = await importAgentSecondStep(
//         jwt!,
//         externalAgent!.collectionAddress!,
//         externalAgent!.collectionTokenId!,
//         values
//       );

//       if (response.status !== 200) {
//         throw new Error("Failed to register agent information");
//       }

//       toast("Agent information registered successfully");

//       fetchAgents();

//       return response.data as DatabaseCredentials;
//     } catch (error) {
//       console.error("Error completing form:", error);
//       toast("Failed to register agent information");
//       return undefined;
//     }
//   };

//   return (
//     <AgentImportationSecondStepDialogContext.Provider
//       value={{
//         open,
//         close,
//         isOpen,
//         externalAgent,
//         onFormComplete,
//         hideCloseButton,
//         displayCloseButton,
//         isCloseButtonHidden,
//         agentBaseInformation,
//       }}
//     >
//       {children}
//     </AgentImportationSecondStepDialogContext.Provider>
//   );
// };

// export const useAgentImportationSecondStepDialog = () => {
//   const context = React.useContext(AgentImportationSecondStepDialogContext);

//   if (!context) {
//     throw new Error(
//       "useAgentImportationSecondStepDialog must be used within an AgentImportationSecondStepDialogContextProvider"
//     );
//   }

//   return context;
// };
