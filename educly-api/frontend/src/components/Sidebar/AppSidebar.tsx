import { BookOpenText, Eye, Home, Package, ShieldCheck, Users, Webhook } from "lucide-react"

import { SidebarAppearance } from "@/components/Common/Appearance"
import { Logo } from "@/components/Common/Logo"
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
} from "@/components/ui/sidebar"
import useAuth from "@/hooks/useAuth"
import { type Item, Main } from "./Main"
import { User } from "./User"

const baseItems: Item[] = [{ icon: Home, title: "Dashboard", path: "/" }]

export function AppSidebar() {
  const { user: currentUser } = useAuth()

  const sections = currentUser?.is_superuser
    ? [
        {
          title: "Geral",
          items: [
            ...baseItems,
            { icon: Package, title: "Perfils", path: "/perfils" },
            { icon: Eye, title: "Look", path: "/look" },
          ],
        },
        {
          title: "Tecnica",
          items: [
            { icon: Users, title: "Admin", path: "/admin" },
            { icon: BookOpenText, title: "Conteudo", path: "/content" },
            { icon: Package, title: "Products", path: "/products" },
            { icon: Webhook, title: "Webhooks", path: "/webhooks" },
            { icon: ShieldCheck, title: "LogPermissions", path: "/log-permissions" },
          ],
        },
      ]
    : [{ title: "Geral", items: baseItems }]

  return (
    <Sidebar collapsible="icon">
      <SidebarHeader className="px-4 py-6 group-data-[collapsible=icon]:px-0 group-data-[collapsible=icon]:items-center">
        <Logo variant="responsive" />
      </SidebarHeader>
      <SidebarContent>
        <Main sections={sections} />
      </SidebarContent>
      <SidebarFooter>
        <SidebarAppearance />
        <User user={currentUser} />
      </SidebarFooter>
    </Sidebar>
  )
}

export default AppSidebar
