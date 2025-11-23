// import AgentSetupForm from "../forms/agent-setup-form";
// import { DatabaseConnectionDriver } from "@/models/enums/database-connection-driver";
// import { LLMConnectionPlatformType } from "@/models/enums/llm-connection-platform-type";
// import { useAgentImportationSecondStepDialog } from "../contexts/agent-importation-second-step-dialog-context";
// import {
//   getLLMConnectionPlatformType,
//   LLMConnectionPlatform,
// } from "@/models/enums/llm-connection-platform";
// import {
//   Dialog,
//   DialogTitle,
//   DialogHeader,
//   DialogContent,
//   DialogDescription,
// } from "@/components/ui/dialog";
// import { AgentRole } from "@/models/enums/agent-role";

// export default function AgentImportationSecondStepDialog() {
//   const {
//     close,
//     isOpen,
//     onFormComplete,
//     hideCloseButton,
//     isCloseButtonHidden,
//     agentBaseInformation,
//   } = useAgentImportationSecondStepDialog();

//   return (
//     <Dialog
//       open={isOpen}
//       onOpenChange={(isOpen) => {
//         if (!isOpen) {
//           close();
//         }
//       }}
//     >
//       <DialogContent hideCloseButton={isCloseButtonHidden}>
//         <DialogHeader>
//           <DialogTitle>Import an agent</DialogTitle>
//           <DialogDescription>
//             You don't known how to import an agent? See our documentation.
//           </DialogDescription>
//         </DialogHeader>
//         <AgentSetupForm
//           onFinishSetup={close}
//           onFormComplete={onFormComplete}
//           hideCloseButton={hideCloseButton}
//           disabledValues={
//             new Map([
//               ["name", true],
//               ["role", true],
//               ["image", true],
//             ])
//           }
//           defaultValues={{
//             image: new File([], ""),
//             name: agentBaseInformation?.name ?? "",
//             role: agentBaseInformation?.role ?? AgentRole.None,
//             llmPlatformType: agentBaseInformation?.llmPlatform
//               ? getLLMConnectionPlatformType(agentBaseInformation.llmPlatform)
//               : LLMConnectionPlatformType.None,
//             llmModel: agentBaseInformation?.llmModel ?? "",
//             llmPlatform:
//               agentBaseInformation?.llmPlatform ?? LLMConnectionPlatform.None,
//             databaseHost: "",
//             databasePort: "",
//             databaseDriver: DatabaseConnectionDriver.None,
//           }}
//         />
//       </DialogContent>
//     </Dialog>
//   );
// }
