import type { KeyValue } from "@/models/key-value"

export type MintProposal = {
  id: number
  collectionAddress: string
  name: string
  role: string
  owner: string
  mission: string
  creator: string
  wallets: string[]
  images: KeyValue[]
}
