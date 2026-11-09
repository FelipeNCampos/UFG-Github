import { createFileRoute, redirect } from "@tanstack/react-router"

export const Route = createFileRoute("/_layout/request-acess")({
  beforeLoad: () => {
    throw redirect({ to: "/log-permissions" })
  },
  component: () => null,
})
