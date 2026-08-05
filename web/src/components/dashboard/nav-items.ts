import {
  Bot,
  Gavel,
  LayoutDashboard,
  MessagesSquare,
  type LucideIcon,
} from "lucide-react"

export type NavItem = {
  title: string
  path: string
  icon: LucideIcon
}

export const browseNavItems: NavItem[] = [
  { title: "Overview", path: "/dashboard/overview", icon: LayoutDashboard },
  { title: "Agents", path: "/dashboard/agents", icon: Bot },
  {
    title: "Conversations",
    path: "/dashboard/conversations",
    icon: MessagesSquare,
  },
]

export const governNavItems: NavItem[] = [
  { title: "Moderation", path: "/dashboard/moderation", icon: Gavel },
]

const allNavItems = [...browseNavItems, ...governNavItems]

// isNavItemActive is prefix-aware so detail routes (e.g. /dashboard/agents/…)
// still highlight and title their parent nav item.
export function isNavItemActive(item: NavItem, pathname: string): boolean {
  return pathname === item.path || pathname.startsWith(`${item.path}/`)
}

export function navTitleForPath(pathname: string): string {
  return (
    allNavItems.find((item) => isNavItemActive(item, pathname))?.title ?? ""
  )
}
