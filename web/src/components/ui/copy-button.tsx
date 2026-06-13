import { useState } from "react";
import { Copy, Check } from "lucide-react";
import { Button } from "@/components/ui/button";

interface CopyButtonProps {
  text: string;
  className?: string;
  resetDelay?: number; // Time (ms) before reverting to copy icon
}

export function CopyButton({
  text,
  className,
  resetDelay = 2000,
}: CopyButtonProps) {
  const [copied, setCopied] = useState(false);

  const copyToClipboard = () => {
    navigator.clipboard.writeText(text);
    setCopied(true);
    setTimeout(() => setCopied(false), resetDelay);
  };

  return (
    <Button
      size="icon"
      type="button"
      variant="ghost"
      className={className}
      onClick={copyToClipboard}
    >
      {copied ? <Check /> : <Copy />}
    </Button>
  );
}
