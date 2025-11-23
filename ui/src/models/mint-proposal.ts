import { KeyValue } from "./key-value";

export type MintProposal = {
  id: number;
  collectionAddress: string;
  name: string;
  role: string;
  owner: string;
  signer: string;
  mission: string;
  creator: string;
  wallets: string[];
  images: KeyValue[];
};
