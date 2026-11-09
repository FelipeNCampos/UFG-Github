import { createFileRoute, redirect } from "@tanstack/react-router"

import { UsersService } from "@/client"
import ContentManagement from "@/components/Admin/ContentManagement"

export const Route = createFileRoute("/_layout/content")({
  component: ContentPage,
  beforeLoad: async () => {
    const user = await UsersService.readUserMe()
    if (!user.is_superuser) {
      throw redirect({
        to: "/",
      })
    }
  },
  head: () => ({
    meta: [
      {
        title: "Conteudo - FastAPI Template",
      },
    ],
  }),
})

function ContentPage() {
  return <ContentManagement />
}
