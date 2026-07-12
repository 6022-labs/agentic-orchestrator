import type { KeyValue } from "@/models/key-value"

export type AddOrUpdateImageProposal = {
  id: number
  collectionAddress: string
  tokenId: bigint
  image: KeyValue
}
