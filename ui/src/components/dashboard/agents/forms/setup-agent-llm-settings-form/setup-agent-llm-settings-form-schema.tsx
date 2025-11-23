import { z } from "zod";
import { LLMConnectionPlatform } from "@/models/enums/llm-connection-platform";
import { LLMConnectionPlatformType } from "@/models/enums/llm-connection-platform-type";

export type SetupAgentLLMSettingsFormSchemaType = {
  platformType: LLMConnectionPlatformType;
  model: string;
  platform: LLMConnectionPlatform;
  apiKey?: string;
  baseUrl?: string;
  amazonPlatformRegion?: string;
};

export const createSetupAgentLLMSettingsFormSchema = () => {
  return z
    .object({
      platformType: z.enum(LLMConnectionPlatformType),
      platform: z.enum(LLMConnectionPlatform),
      model: z.string().min(2, {
        message: "Model must be at least 2 characters.",
      }),
      baseUrl: z
        .string()
        .url({
          protocol: /^https?$/,
          hostname: z.regexes.domain,
        })
        .optional(),
      apiKey: z.string().optional(),
      amazonPlatformRegion: z
        .string()
        .regex(/^[a-z]{2}(-[a-z]+)+-\d+$/i, {
          message:
            "Region must be a valid AWS region like 'us-west-2' or 'ap-southeast-1'.",
        })
        .optional(),
    })
    .superRefine((data, ctx) => {
      if (
        data.platform == LLMConnectionPlatform.AmazonBedrock &&
        (!data.amazonPlatformRegion || data.amazonPlatformRegion.trim() === "")
      ) {
        ctx.addIssue({
          code: "custom",
          message: "Region is required for Amazon platform.",
          path: ["amazonPlatformRegion"],
        });
      }
    });
};
