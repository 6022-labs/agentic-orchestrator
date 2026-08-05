import BlockiesRaw from "react-blockies"
import { ChevronsUpDown, LogOut, Wallet } from "lucide-react"
import { Link, useLocation } from "react-router-dom"
import { useAppKit } from "@reown/appkit/react"

import { truncateEthAddress } from "@/lib/utils"
import { useAuth } from "@/contexts/auth-context"
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  useSidebar,
} from "@/components/ui/sidebar"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import {
  browseNavItems,
  governNavItems,
  isNavItemActive,
  type NavItem,
} from "./nav-items"

// react-blockies ships CJS; under this bundler the default import can arrive
// wrapped as { default }, so unwrap it to get the actual component.
const Blockies = ((BlockiesRaw as unknown as { default?: typeof BlockiesRaw })
  .default ?? BlockiesRaw) as typeof BlockiesRaw

function NavGroup({ label, items }: { label: string; items: NavItem[] }) {
  const location = useLocation()

  return (
    <SidebarGroup>
      <SidebarGroupLabel>{label}</SidebarGroupLabel>
      <SidebarGroupContent>
        <SidebarMenu>
          {items.map((item) => (
            <SidebarMenuItem key={item.path}>
              <SidebarMenuButton
                asChild
                isActive={isNavItemActive(item, location.pathname)}
                tooltip={item.title}
              >
                <Link to={item.path}>
                  <item.icon />
                  <span>{item.title}</span>
                </Link>
              </SidebarMenuButton>
            </SidebarMenuItem>
          ))}
        </SidebarMenu>
      </SidebarGroupContent>
    </SidebarGroup>
  )
}

export function AppSidebar() {
  const { open } = useAppKit()
  const { user, signOut } = useAuth()
  const { setOpenMobile } = useSidebar()

  // Close the mobile sheet first so the AppKit modal doesn't stack under it.
  const openWalletModal = () => {
    setOpenMobile(false)
    open()
  }

  return (
    <Sidebar collapsible="icon">
      <SidebarHeader>
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton size="lg" asChild tooltip="Overview">
              <Link to="/dashboard/overview">
                <img
                  src="/logo-192x192.png"
                  alt="Logo"
                  className="size-8 shrink-0 rounded-md"
                />
                <div className="grid flex-1 text-left leading-tight">
                  <span className="truncate font-semibold">
                    Agentic Orchestrator
                  </span>
                  <span className="truncate text-xs text-muted-foreground">
                    Orchestrator console
                  </span>
                </div>
              </Link>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarHeader>
      <SidebarContent>
        <NavGroup label="Browse" items={browseNavItems} />
        <NavGroup label="Govern" items={governNavItems} />
      </SidebarContent>
      <SidebarFooter>
        <SidebarMenu>
          <SidebarMenuItem>
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <SidebarMenuButton size="lg" tooltip="Account">
                  <span className="flex size-8 shrink-0 items-center justify-center overflow-hidden rounded-full">
                    <Blockies seed={user?.address ?? ""} />
                  </span>
                  <span className="flex-1 truncate text-left">
                    {user?.address
                      ? truncateEthAddress(user.address, 8, 6)
                      : ""}
                  </span>
                  <ChevronsUpDown className="ml-auto size-4" />
                </SidebarMenuButton>
              </DropdownMenuTrigger>
              <DropdownMenuContent
                side="top"
                align="start"
                className="w-(--radix-dropdown-menu-trigger-width) min-w-48"
              >
                <DropdownMenuItem onSelect={openWalletModal}>
                  <Wallet />
                  Wallet
                </DropdownMenuItem>
                <DropdownMenuItem onSelect={signOut}>
                  <LogOut />
                  Sign out
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarFooter>
    </Sidebar>
  )
}
