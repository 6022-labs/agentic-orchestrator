import React from "react";
import { useForm } from "react-hook-form";
import { Form } from "@/components/ui/form";
import { zodResolver } from "@hookform/resolvers/zod";
import { SetupAgentDatabaseFirstStepInformationForm } from "./setup-agent-database-first-step-information";
import {
  SetupAgentDatabaseSettingsFirstStepFormSchemaType,
  createSetupAgentDatabaseSettingsFirstStepFormSchema,
} from "./setup-agent-database-settings-first-step-form-schema";
import {
  SetupAgentDatabaseFirstStepResult,
  SetupAgentDatabaseFirstStepResultInformation,
} from "./setup-agent-database-first-step-result-information";
import {
  Carousel,
  CarouselApi,
  CarouselItem,
  CarouselContent,
} from "@/components/ui/carousel";
import type { FormCompletionResult } from "@/models/api/validation-error";

type SetupAgentDatabaseSettingsFirstStepFormProps = {
  hideCloseButton: () => void;
  onFormComplete: (
    values: SetupAgentDatabaseSettingsFirstStepFormSchemaType
  ) => Promise<
    FormCompletionResult<
      keyof SetupAgentDatabaseSettingsFirstStepFormSchemaType,
      SetupAgentDatabaseFirstStepResult
    >
  >;
  onFormFinished: () => void;
  disabledValues?: Map<string, boolean>;
  defaultValues?: SetupAgentDatabaseSettingsFirstStepFormSchemaType;
};

export default function SetupAgentDatabaseSettingsFirstStepForm(
  props: Readonly<SetupAgentDatabaseSettingsFirstStepFormProps>
) {
  const [carouselApi, setCarouselApi] = React.useState<CarouselApi>();
  const [currentStep, setCurrentStep] = React.useState(0);

  const [result, setResult] = React.useState<
    SetupAgentDatabaseFirstStepResult | undefined
  >(undefined);

  const setupAgentDatabaseSettingsFormSchema =
    createSetupAgentDatabaseSettingsFirstStepFormSchema();

  const form = useForm<SetupAgentDatabaseSettingsFirstStepFormSchemaType>({
    resolver: zodResolver(setupAgentDatabaseSettingsFormSchema),
    defaultValues: props.defaultValues,
  });

  const onDatabaseInformationFormNext = async () => {
    const values = form.getValues();

    const submission = await props.onFormComplete(values);
    if (submission?.error) {
      form.setError(
        submission.error
          .field as keyof SetupAgentDatabaseSettingsFirstStepFormSchemaType,
        {
          type: "server",
          message: submission.error.message,
        }
      );
      return;
    }

    const data = submission?.result;
    if (!data) {
      return;
    }

    setResult(data);
    props.hideCloseButton();

    setCurrentStep(1);
    carouselApi?.scrollNext();
  };

  const onResultInformationNext = async () => {
    props.onFormFinished();
  };

  return (
    <Form {...form}>
      <form className="h-fit">
        <Carousel
          opts={{
            align: "start",
            dragFree: true,
            watchDrag: false,
          }}
          setApi={setCarouselApi}
        >
          <CarouselContent>
            <CarouselItem>
              <SetupAgentDatabaseFirstStepInformationForm
                form={form}
                nextButtonLabel="Submit"
                onNextValidStep={onDatabaseInformationFormNext}
                className={currentStep !== 0 ? "hidden" : ""}
              />
            </CarouselItem>
            <CarouselItem>
              <SetupAgentDatabaseFirstStepResultInformation
                result={result}
                driver={form.getValues("driver")}
                onNextValidStep={onResultInformationNext}
                className={currentStep !== 1 ? "hidden" : ""}
              />
            </CarouselItem>
          </CarouselContent>
        </Carousel>
      </form>
    </Form>
  );
}
