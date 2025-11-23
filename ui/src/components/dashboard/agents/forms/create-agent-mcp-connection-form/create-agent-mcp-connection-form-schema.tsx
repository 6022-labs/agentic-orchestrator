import { z } from "zod";

export enum MCPAuthenticationMethod {
  NONE = "none",
  API_KEY = "api-key",
  CUSTOM_HEADER = "custom-header",
}

export const MCPAuthenticationMethodDisplayNames: Record<
  MCPAuthenticationMethod,
  string
> = {
  [MCPAuthenticationMethod.NONE]: "No Authentication",
  [MCPAuthenticationMethod.API_KEY]: "API Key",
  [MCPAuthenticationMethod.CUSTOM_HEADER]: "Custom Header",
};

// Regex copied from backend (Go) for https URL validation
const hostnameLabelPattern =
  "(?:[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)";
const tldLabelPattern = "[A-Za-z]{2,63}";
const hostnamePattern = "(?:" + hostnameLabelPattern + ".)+" + tldLabelPattern;
const ipv4SegmentPattern = "(?:25[0-5]|2[0-4]\\d|1\\d\\d|[1-9]?\\d)";
const ipv4Pattern = "(?:" + ipv4SegmentPattern + ".){3}" + ipv4SegmentPattern;
const httpsURLRegex = new RegExp(
  "^https://(?:" +
    ipv4Pattern +
    "|" +
    hostnamePattern +
    ")(?::\\d{1,5})?(?:[/?#].*)?$",
  "i"
);

export const createAgentMCPConnectionFormSchema = z
  .object({
    name: z
      .string()
      .min(3, "Name must be at least 3 characters")
      .regex(
        /^[a-z][a-z0-9_-]*$/,
        "Name must start with a letter and contain only letters, digits, '-' and '_'"
      ),
    url: z
      .string()
      .min(2, "Url is required")
      .refine((val) => httpsURLRegex.test(val.trim()), {
        message:
          "Url must be a valid https URL with a valid hostname or IPv4 address",
      }),
    authenticationMethod: z.enum(MCPAuthenticationMethod, {
      error: () => ({ message: "Authentication method is required" }),
    }),
    headerAuthKey: z.string().optional(),
    headerAuthValue: z.string().optional(),
  })
  .superRefine((data, ctx) => {
    if (data.authenticationMethod === MCPAuthenticationMethod.NONE) {
      if (data.headerAuthKey || data.headerAuthValue) {
        ctx.addIssue({
          code: "custom",
          message:
            "Header auth key and value must be empty when authentication method is 'none'",
          path: ["headerAuthValue"],
        });
      }
    } else {
      if (!data.headerAuthKey) {
        ctx.addIssue({
          code: "custom",
          message:
            "Header auth key must be provided if authentication is not 'none'",
          path: ["headerAuthKey"],
        });
      }
      if (!data.headerAuthValue) {
        ctx.addIssue({
          code: "custom",
          message:
            "Header auth value must be provided if authentication is not 'none'",
          path: ["headerAuthValue"],
        });
      }
    }
  });

export type CreateAgentMCPConnectionFormSchemaType = z.infer<
  typeof createAgentMCPConnectionFormSchema
>;
