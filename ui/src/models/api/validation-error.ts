export type ApiValidationError<FieldName extends string = string> = {
  field: FieldName;
  message: string;
};

export type FormCompletionResult<
  FieldName extends string = string,
  TResult = void
> = {
  error?: ApiValidationError<FieldName>;
  result?: TResult;
};
