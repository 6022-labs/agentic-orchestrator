import axios from "axios"

import type { PromptTemplate } from "@/models/prompt-template"

// GET /agent/{collectionAddress}/{collectionTokenId}/prompt-template/{key} — 404 when unset.
export const getPromptTemplate = async (
  collectionAddress: string,
  collectionTokenId: string,
  key: string
) => {
  return await axios.get<PromptTemplate>(
    `/agent/${collectionAddress}/${collectionTokenId}/prompt-template/${key}`
  )
}

// POST /agent/{collectionAddress}/{collectionTokenId}/prompt-template/{key} — upsert the template.
export const savePromptTemplate = async (
  collectionAddress: string,
  collectionTokenId: string,
  key: string,
  content: string
) => {
  return await axios.post<PromptTemplate>(
    `/agent/${collectionAddress}/${collectionTokenId}/prompt-template/${key}`,
    { content }
  )
}
