import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import { useAppKit } from "@reown/appkit/react";

export function LoginForm({
  className,
  ...props
}: React.ComponentProps<"form">) {
  const { open } = useAppKit();

  return (
    <form className={cn("flex flex-col gap-6", className)} {...props}>
      <div className="flex flex-col items-center gap-2 text-center">
        <h1 className="text-2xl font-bold">Login to your account</h1>
        <p className="text-muted-foreground text-sm text-balance">
          Use your web3 wallet to login
        </p>
      </div>
      <div className="grid gap-6">
        <Button type="button" className="w-full" onClick={() => open()}>
          Login
        </Button>
      </div>
    </form>
  );
}
