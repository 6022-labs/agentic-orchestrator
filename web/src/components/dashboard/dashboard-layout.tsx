import { Outlet, useLocation } from "react-router-dom"

import { AppSidebar } from "@/components/dashboard/app-sidebar"
import { navTitleForPath } from "@/components/dashboard/nav-items"
import OrchestratorAddressChip from "@/components/dashboard/orchestrator-address-chip"
import ThemeToggle from "@/components/dashboard/theme-toggle"
import { Separator } from "@/components/ui/separator"
import {
  MAIN_CONTENT_ID,
  PAGE_TITLE_ID,
  useRouteFocus,
} from "@/hooks/use-route-focus"
import {
  SidebarInset,
  SidebarProvider,
  SidebarTrigger,
} from "@/components/ui/sidebar"

export default function DashboardLayout() {
  const location = useLocation()
  useRouteFocus()

  return (
    <>
      <a
        href={`#${MAIN_CONTENT_ID}`}
        onClick={(event) => {
          // A fragment navigation would add the hash to the route.
          event.preventDefault()
          document.getElementById(MAIN_CONTENT_ID)?.focus()
        }}
        className="sr-only focus:not-sr-only focus:fixed focus:top-3 focus:left-3 focus:z-50 focus:rounded-md focus:bg-background focus:px-3 focus:py-2 focus:text-sm focus:font-medium focus:shadow-md focus:ring-2 focus:ring-ring focus:outline-none"
      >
        Skip to content
      </a>
      <SidebarProvider>
        <AppSidebar />
        <SidebarInset className="h-svh overflow-hidden">
          <header className="flex h-14 shrink-0 items-center gap-2 border-b px-4">
            <SidebarTrigger className="-ml-1" />
            <Separator orientation="vertical" className="mr-2 h-4" />
            <h1
              id={PAGE_TITLE_ID}
              tabIndex={-1}
              className="text-base font-semibold"
            >
              {navTitleForPath(location.pathname)}
            </h1>
            <div className="ml-auto flex items-center gap-1">
              <OrchestratorAddressChip />
              <ThemeToggle />
            </div>
          </header>
          <div
            id={MAIN_CONTENT_ID}
            tabIndex={-1}
            className="flex min-h-0 flex-1 flex-col overflow-y-auto p-6 outline-none"
          >
            <Outlet />
          </div>
        </SidebarInset>
      </SidebarProvider>
    </>
  )
}
