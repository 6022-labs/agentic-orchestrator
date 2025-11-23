import { toast } from "sonner";
import { Plus } from "lucide-react";
import { Eip1193Provider } from "ethers";
import { Button } from "@/components/ui/button";
import { useAuth } from "@/contexts/auth-context";
import { useAppKitProvider } from "@reown/appkit/react";
import { FormCompletionResult } from "@/models/api/validation-error";
import CreateAgentCollectionForm from "../forms/create-agent-collection-form";
import { createCollection } from "@/queries/contracts/agent-collection-creator";
import { useCreateAgentCollectionDialog } from "../contexts/create-agent-collection-dialog-context";
import { CreateAgentCollectionFormSchemaType } from "../forms/create-agent-collection-form/create-agent-collection-form-schema";
import {
  Dialog,
  DialogTitle,
  DialogHeader,
  DialogTrigger,
  DialogContent,
  DialogDescription,
} from "@/components/ui/dialog";

export default function CreateAgentCollectionDialog() {
  const { isOpen, canCreateCollection, setIsOpen } =
    useCreateAgentCollectionDialog();

  const { user } = useAuth();
  const { walletProvider } = useAppKitProvider("eip155");
  const eip1193Provider = walletProvider as Eip1193Provider;

  const onFormComplete = async (
    values: CreateAgentCollectionFormSchemaType
  ): Promise<
    FormCompletionResult<keyof CreateAgentCollectionFormSchemaType>
  > => {
    try {
      await createCollection(
        values.name,
        values.symbol,
        values.description,
        user!.address!,
        user!.address!,
        eip1193Provider
      );
      toast.success("Agent collection created successfully");
      setIsOpen(false);
    } catch (error) {
      console.error("Error creating agent collection:", error);
      toast.error("Failed to create agent collection");

      return {};
    }

    return {};
  };

  return (
    <>
      {canCreateCollection && (
        <Dialog open={isOpen} onOpenChange={setIsOpen}>
          <DialogTrigger asChild>
            <Button type="button">
              <Plus className="mr-1" />
              Create Collection
            </Button>
          </DialogTrigger>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Create an agent collection</DialogTitle>
              <DialogDescription>
                You don't known how to create an agent collection? See our
                documentation.
              </DialogDescription>
            </DialogHeader>
            <CreateAgentCollectionForm onFormComplete={onFormComplete} />
          </DialogContent>
        </Dialog>
      )}
    </>
  );
}
