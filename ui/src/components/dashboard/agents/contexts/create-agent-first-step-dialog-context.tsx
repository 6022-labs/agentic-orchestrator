/* eslint-disable react-refresh/only-export-components */
import React from "react";
import { AgentRole, getAllAgentRoles } from "@/models/enums/agent-role";
import {
  CreateAgentFirstStepDisabledValues,
  CreateAgentFirstStepFormSchemaType,
} from "../forms/create-agent-first-step-form/create-agent-first-step-form-schema";

type CreateAgentFirstStepDialogContextType = {
  isOpen: boolean;
  accessibleRoles: AgentRole[];
  setIsOpen: (isOpen: boolean) => void;
  setAccessibleRoles: (roles: AgentRole[]) => void;
  defaultValues: Partial<CreateAgentFirstStepFormSchemaType>;
  disabledValues: CreateAgentFirstStepDisabledValues;
  setDisabledValues: (values: CreateAgentFirstStepDisabledValues) => void;
  setDefaultValues: (
    values: Partial<CreateAgentFirstStepFormSchemaType>
  ) => void;
};

export const CreateAgentFirstStepDialogContext =
  React.createContext<CreateAgentFirstStepDialogContextType>({
    isOpen: false,
    defaultValues: {},
    disabledValues: {},
    setIsOpen: async () => {},
    setDefaultValues: () => {},
    setDisabledValues: () => {},
    setAccessibleRoles: () => {},
    accessibleRoles: getAllAgentRoles(),
  });

type CreateAgentFirstStepDialogContextProviderProps = {
  children: React.ReactNode;
};

export const CreateAgentFirstStepDialogContextProvider: React.FC<
  CreateAgentFirstStepDialogContextProviderProps
> = ({ children }) => {
  const [isOpen, setIsOpen] = React.useState(false);
  const [accessibleRoles, setAccessibleRoles] = React.useState<AgentRole[]>(
    getAllAgentRoles()
  );
  const [defaultValues, setDefaultValues] = React.useState<
    Partial<CreateAgentFirstStepFormSchemaType>
  >({});
  const [disabledValues, setDisabledValues] =
    React.useState<CreateAgentFirstStepDisabledValues>({});

  return (
    <CreateAgentFirstStepDialogContext.Provider
      value={{
        isOpen,
        setIsOpen,
        accessibleRoles,
        setAccessibleRoles,
        defaultValues,
        setDefaultValues,
        disabledValues,
        setDisabledValues,
      }}
    >
      {children}
    </CreateAgentFirstStepDialogContext.Provider>
  );
};

export const useCreateAgentFirstStepDialog = () => {
  const context = React.useContext(CreateAgentFirstStepDialogContext);

  if (!context) {
    throw new Error(
      "useCreateAgentFirstStepDialog must be used within a CreateAgentFirstStepDialogContextProvider"
    );
  }

  return context;
};
