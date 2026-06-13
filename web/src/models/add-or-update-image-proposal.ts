import { KeyValue } from "./key-value";

export type AddOrUpdateImageProposal = {
  id: number;
  collectionAddress: string;
  tokenId: bigint;
  image: KeyValue;
};
