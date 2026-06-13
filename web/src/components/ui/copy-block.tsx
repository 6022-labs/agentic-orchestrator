import { CopyButton } from "./copy-button";
import { Card, CardContent } from "@/components/ui/card";

interface CopyBlockProps {
  children: string;
  resetDelay?: number; // Time (ms) before reverting to copy icon
}

export function CopyBlock({ children, resetDelay = 2000 }: CopyBlockProps) {
  return (
    <Card className="flex flex-row items-center justify-between py-1 px-2 bg-muted rounded-lg">
      <CardContent className="p-0 text-sm font-mono truncate">
        {children}
      </CardContent>
      <CopyButton text={children} className="ml-2" resetDelay={resetDelay} />
    </Card>
  );
}
