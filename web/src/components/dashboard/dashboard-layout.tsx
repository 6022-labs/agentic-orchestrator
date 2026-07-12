import { Outlet, useLocation } from "react-router-dom"

import { AppSidebar } from "@/components/dashboard/app-sidebar"
import { navTitleForPath } from "@/components/dashboard/nav-items"
import ThemeToggle from "@/components/dashboard/theme-toggle"
import { Separator } from "@/components/ui/separator"
import {
  SidebarInset,
  SidebarProvider,
  SidebarTrigger,
} from "@/components/ui/sidebar"

export default function DashboardLayout() {
  const location = useLocation()

  return (
    <SidebarProvider>
      <AppSidebar />
      <SidebarInset className="h-svh overflow-hidden">
        <header className="flex h-14 shrink-0 items-center gap-2 border-b px-4">
          <SidebarTrigger className="-ml-1" />
          <Separator orientation="vertical" className="mr-2 h-4" />
          <h1 className="text-base font-semibold">
            {navTitleForPath(location.pathname)}
          </h1>
          <div className="ml-auto">
            <ThemeToggle />
          </div>
        </header>
        <div className="flex min-h-0 flex-1 flex-col overflow-y-auto p-6">
          <Outlet />
        </div>
      </SidebarInset>
    </SidebarProvider>
  )
}
