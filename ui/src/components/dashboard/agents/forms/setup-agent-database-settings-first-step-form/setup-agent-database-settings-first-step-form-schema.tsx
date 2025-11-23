import { z } from "zod";
import { DatabaseConnectionDriver } from "@/models/enums/database-connection-driver";

export type SetupAgentDatabaseSettingsFirstStepFormSchemaType = {
  host: string;
  port: number;
  driver: DatabaseConnectionDriver;
};

export const createSetupAgentDatabaseSettingsFirstStepFormSchema = () => {
  return z.object({
    host: z.string().superRefine((value, ctx) => {
      if (!value) {
        ctx.addIssue({
          code: "custom",
          message: "Database host is required.",
        });
      }
    }),
    port: z
      .number()
      .min(1, "Port must be a positive number.")
      .max(65535, "Port must be less than or equal to 65535."),
    driver: z.enum(DatabaseConnectionDriver),
  });
};
