import Blockies from "react-blockies";
import { useAppKit } from "@reown/appkit/react";
import { useAuth } from "@/contexts/auth-context";
import { cn, truncateEthAddress } from "@/lib/utils";
import { useUserModerator } from "@/contexts/user-moderator-context";

type DashboardHeaderProps = {
  className?: string;
  onAgentsClick?: () => void;
  onOverviewClick?: () => void;
  onModerationClick?: () => void;
};

export function DashboardHeader(props: DashboardHeaderProps) {
  const { open } = useAppKit();
  const { user } = useAuth();
  const { collections } = useUserModerator();

  const Headers = [
    {
      title: "Overview",
      onClick: () => {
        if (props.onOverviewClick) {
          props.onOverviewClick();
        }
      },
    },
    {
      title: "Agents",
      onClick: () => {
        if (props.onAgentsClick) {
          props.onAgentsClick();
        }
      },
    },
  ];

  if (collections.length !== 0) {
    Headers.push({
      title: "Moderation",
      onClick: () => {
        if (props.onModerationClick) {
          props.onModerationClick();
        }
      },
    });
  }

  return (
    <div className={cn("border-b", props.className)}>
      <div className="flex h-16 items-center px-4">
        <nav className="flex items-center space-x-4 lg:space-x-6 mx-6">
          {Headers.map((header, index) => (
            <button
              key={index}
              className="text-base font-medium text-muted-foreground cursor-pointer hover:text-accent"
              onClick={header.onClick}
            >
              {header.title}
            </button>
          ))}
        </nav>
        <div className="ml-auto flex items-center space-x-4">
          <div className="text-muted-foreground">
            {user?.address ? truncateEthAddress(user.address, 8, 6) : ""}
          </div>
          <button
            className="inline-flex items-center justify-center gap-2 whitespace-nowrap text-sm font-medium transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50 [&amp;_svg]:pointer-events-none [&amp;_svg]:size-4 [&amp;_svg]:shrink-0 hover:bg-accent hover:text-accent-foreground px-4 py-2 relative h-8 w-8 rounded-full"
            type="button"
            onClick={() => open()}
          >
            <span className="relative flex shrink-0 overflow-hidden rounded-full h-8 w-8">
              <Blockies seed={user?.address ?? ""} />
            </span>
          </button>
        </div>
      </div>
    </div>
  );
}
