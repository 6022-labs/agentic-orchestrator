import { ExternalLink } from "lucide-react";
import { Button } from "@/components/ui/button";

interface ViewOnScannerLinkProps {
  href: string;
  className?: string;
}

export function ViewOnScannerLink({ href, className }: ViewOnScannerLinkProps) {
  return (
    <Button
      size="icon"
      type="button"
      variant="ghost"
      className={className}
      onClick={() => window.open(href, "_blank")}
    >
      <ExternalLink />
    </Button>
  );
}
