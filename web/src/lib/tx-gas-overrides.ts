import { ethers } from "ethers"

import { config } from "@/configs/env"

const TX_GAS_MULTIPLIER_SCALE = 1_000n

export type TxGasOverrides = {
  gasPrice?: bigint
  maxFeePerGas?: bigint
  maxPriorityFeePerGas?: bigint
}

// Fixed-point ceiling so fractional multipliers (e.g. 1.5) round up, never down.
function scaleWithMultiplier(value: bigint, multiplier: number): bigint {
  if (multiplier <= 1) return value

  const scaledMultiplier = BigInt(
    Math.ceil(multiplier * Number(TX_GAS_MULTIPLIER_SCALE))
  )

  return (
    (value * scaledMultiplier + TX_GAS_MULTIPLIER_SCALE - 1n) /
    TX_GAS_MULTIPLIER_SCALE
  )
}

// Scales the provider's suggested fees by VITE_TX_GAS_MULTIPLIER; returns {}
// when the multiplier is 1 so callers fall back to ethers defaults.
export async function getScaledTxGasOverrides(
  signer: ethers.Signer
): Promise<TxGasOverrides> {
  const multiplier = config.txGasMultiplier

  if (multiplier <= 1 || !signer.provider) {
    return {}
  }

  const feeData = await signer.provider.getFeeData()
  const gasOverrides: TxGasOverrides = {}

  if (feeData.maxFeePerGas !== null || feeData.maxPriorityFeePerGas !== null) {
    if (feeData.maxFeePerGas !== null) {
      gasOverrides.maxFeePerGas = scaleWithMultiplier(
        feeData.maxFeePerGas,
        multiplier
      )
    }

    if (feeData.maxPriorityFeePerGas !== null) {
      gasOverrides.maxPriorityFeePerGas = scaleWithMultiplier(
        feeData.maxPriorityFeePerGas,
        multiplier
      )
    }

    return gasOverrides
  }

  if (feeData.gasPrice !== null) {
    gasOverrides.gasPrice = scaleWithMultiplier(feeData.gasPrice, multiplier)
  }

  return gasOverrides
}
