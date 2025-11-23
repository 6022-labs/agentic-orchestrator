import axios from "axios";
// import { AgentSetupFormSchemaType } from "@/components/dashboard/agents/forms/agent-setup-form/agent-setup-form-schema";

export const importAgentFirstStep = async (
  jwt: string,
  agentCollectionAddress: string,
  agentCollectionTokenId: number
) => {
  return await axios.post(
    `agent/${agentCollectionAddress}/${agentCollectionTokenId}/import/first-step`,
    undefined,
    {
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${jwt}`,
      },
    }
  );
};

// export const importAgentSecondStep = async (
//   jwt: string,
//   agentCollectionAddress: string,
//   agentCollectionTokenId: number,
//   body: AgentSetupFormSchemaType
// ) => {
//   return await axios.post(
//     `/agent/${agentCollectionAddress}/${agentCollectionTokenId}/import/second-step`,
//     body,
//     {
//       headers: {
//         "Content-Type": "application/json",
//         Authorization: `Bearer ${jwt}`,
//       },
//     }
//   );
// };

export const importAgentThirdStep = async (
  jwt: string,
  agentCollectionAddress: string,
  agentCollectionTokenId: number
) => {
  return await axios.post(
    `agent/${agentCollectionAddress}/${agentCollectionTokenId}/import/third-step`,
    undefined,
    {
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${jwt}`,
      },
    }
  );
};
