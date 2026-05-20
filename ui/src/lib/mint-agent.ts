import { strToBytes } from "./utils";
import { Agent } from "@/models/agent";
import { Eip1193Provider, ethers } from "ethers";
import { createAgentSecondStep } from "@/queries/api/agent.queries";
import {
  mint,
  moderatorCount,
  hasModeratorRole,
  createMintProposal,
} from "@/queries/contracts/agent-collection";

type MintAgentResult = {
  type: "normal" | "proposal";
};

export async function mintAgent(
  agent: Agent,
  eip1193Provider: Eip1193Provider
): Promise<MintAgentResult> {
  const count = await moderatorCount(agent.collectionAddress, eip1193Provider);
  const isModerator = await hasModeratorRole(
    agent.collectionAddress,
    agent.creator, // Only the creator can mint the agent
    eip1193Provider
  );

  const response = await createAgentSecondStep(agent!.id);

  if (response.status !== 200) {
    throw new Error(
      "Received non-200 response from agent creation second step"
    );
  }

  const nonce = response.data.nonce.replace(/-/g, "");
  const byte32Nonce = ethers.keccak256(strToBytes(nonce));

  if (count !== 0 && !isModerator) {
    await createMintProposal(
      agent.collectionAddress,
      response.data.signer,
      agent.name!,
      agent.role,
      [agent.walletAddress],
      agent.images,
      agent.attributes,
      agent.cloneOf,
      byte32Nonce,
      response.data.signature,
      eip1193Provider
    );

    return { type: "proposal" };
  }

  await mint(
    agent.collectionAddress,
    response.data.signer,
    agent.name!,
    agent.role,
    [agent.walletAddress],
    agent.images,
    agent.attributes,
    agent.cloneOf,
    byte32Nonce,
    response.data.signature,
    eip1193Provider
  );

  return { type: "normal" };
}
