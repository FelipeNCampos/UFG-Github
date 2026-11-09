import { useQuery } from "@tanstack/react-query"
import { createFileRoute, redirect } from "@tanstack/react-router"
import axios from "axios"
import { ChevronDown, ChevronUp } from "lucide-react"
import { useState } from "react"

import { UsersService } from "@/client"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { Skeleton } from "@/components/ui/skeleton"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"

export const Route = createFileRoute("/_layout/webhooks")({
  component: WebhooksPage,
  beforeLoad: async () => {
    const user = await UsersService.readUserMe()
    if (!user.is_superuser) {
      throw redirect({ to: "/" })
    }
  },
  head: () => ({
    meta: [{ title: "Webhooks - Educly API" }],
  }),
})

type WebhookEvent = {
  id: string
  provider: string
  url: string
  payload: Record<string, unknown>
  headers: Record<string, unknown>
  created_at: string | null
}

type WebhookResponse = {
  data: WebhookEvent[]
  count: number
}

const API_BASE_URL = import.meta.env.VITE_API_URL

function getAuthHeaders() {
  return {
    Authorization: `Bearer ${localStorage.getItem("access_token") || ""}`,
  }
}

function WebhooksPage() {
  const [search, setSearch] = useState("")
  const [provider, setProvider] = useState("all")
  const [expandedIds, setExpandedIds] = useState<Set<string>>(new Set())

  const query = useQuery({
    queryKey: ["webhooks", provider, search],
    queryFn: async () => {
      const params: Record<string, string | number> = { skip: 0, limit: 500 }
      if (provider !== "all") params.provider = provider
      if (search.trim()) params.q = search.trim()
      const response = await axios.get<WebhookResponse>(`${API_BASE_URL}/api/v1/webhooks/`, {
        headers: getAuthHeaders(),
        params,
      })
      return response.data
    },
  })

  if (query.isLoading) {
    return (
      <div className="space-y-4">
        <Skeleton className="h-10 w-64" />
        <Skeleton className="h-10 w-full" />
        <Skeleton className="h-72 w-full" />
      </div>
    )
  }

  if (query.isError) {
    return (
      <div className="rounded-md border border-destructive/30 bg-destructive/10 p-4 text-sm text-destructive">
        Falha ao carregar webhooks.
      </div>
    )
  }

  const events = query.data?.data ?? []
  const providers = Array.from(new Set(events.map((item) => item.provider))).sort()
  const toggleExpanded = (id: string) => {
    setExpandedIds((prev) => {
      const next = new Set(prev)
      if (next.has(id)) next.delete(id)
      else next.add(id)
      return next
    })
  }

  return (
    <div className="flex flex-col gap-6">
      <div>
        <h1 className="text-2xl font-bold tracking-tight">Webhooks</h1>
        <p className="text-muted-foreground">Visualize eventos recebidos por webhook.</p>
      </div>

      <div className="grid grid-cols-1 gap-3 md:grid-cols-2">
        <Input
          placeholder="Pesquisar por provider ou URL"
          value={search}
          onChange={(e) => setSearch(e.target.value)}
        />
        <Select value={provider} onValueChange={setProvider}>
          <SelectTrigger>
            <SelectValue placeholder="Filtrar por provider" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">Todos os providers</SelectItem>
            {providers.map((item) => (
              <SelectItem key={item} value={item}>
                {item}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>

      <div className="rounded-md border">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Provider</TableHead>
              <TableHead>URL</TableHead>
              <TableHead>Payload</TableHead>
              <TableHead>Criado em</TableHead>
              <TableHead>Acoes</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {events.length === 0 ? (
              <TableRow>
                <TableCell className="text-muted-foreground" colSpan={5}>
                  Nenhum webhook encontrado.
                </TableCell>
              </TableRow>
            ) : (
              events.flatMap((event) => {
                const isExpanded = expandedIds.has(event.id)
                return [
                  <TableRow key={event.id}>
                    <TableCell className="font-medium">{event.provider}</TableCell>
                    <TableCell className="max-w-sm truncate">{event.url}</TableCell>
                    <TableCell className="max-w-sm truncate font-mono text-xs">
                      {JSON.stringify(event.payload)}
                    </TableCell>
                    <TableCell>
                      {event.created_at ? new Date(event.created_at).toLocaleString() : "N/A"}
                    </TableCell>
                    <TableCell>
                      <Button
                        size="sm"
                        variant="outline"
                        onClick={() => toggleExpanded(event.id)}
                      >
                        {isExpanded ? (
                          <>
                            <ChevronUp className="mr-1 h-4 w-4" />
                            Recolher
                          </>
                        ) : (
                          <>
                            <ChevronDown className="mr-1 h-4 w-4" />
                            Expandir
                          </>
                        )}
                      </Button>
                    </TableCell>
                  </TableRow>,
                  ...(isExpanded
                    ? [
                        <TableRow key={`${event.id}-expanded`}>
                          <TableCell colSpan={5} className="whitespace-normal">
                            <div className="grid gap-4 md:grid-cols-2">
                              <div>
                                <p className="mb-2 text-xs font-semibold uppercase text-muted-foreground">
                                  Payload completo
                                </p>
                                <pre className="max-h-96 overflow-auto rounded-md border bg-muted/30 p-3 text-xs">
                                  {JSON.stringify(event.payload, null, 2)}
                                </pre>
                              </div>
                              <div>
                                <p className="mb-2 text-xs font-semibold uppercase text-muted-foreground">
                                  Headers
                                </p>
                                <pre className="max-h-96 overflow-auto rounded-md border bg-muted/30 p-3 text-xs">
                                  {JSON.stringify(event.headers, null, 2)}
                                </pre>
                              </div>
                            </div>
                          </TableCell>
                        </TableRow>,
                      ]
                    : []),
                ]
              })
            )}
          </TableBody>
        </Table>
      </div>
    </div>
  )
}
