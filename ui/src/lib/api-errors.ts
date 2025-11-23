import axios from "axios";
import type { ApiValidationError } from "@/models/api/validation-error";

export const extractValidationError = <FieldName extends string = string>(
  error: unknown
): ApiValidationError<FieldName> | undefined => {
  if (!axios.isAxiosError(error)) {
    return undefined;
  }

  if (error.response?.status !== 400) {
    return undefined;
  }

  const data = error.response.data as {
    field?: string;
    message?: string;
  };

  if (typeof data?.field !== "string" || typeof data?.message !== "string") {
    return undefined;
  }

  return {
    field: data.field as FieldName,
    message: data.message,
  };
};
