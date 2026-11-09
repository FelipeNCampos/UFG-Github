import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { createFileRoute, redirect } from "@tanstack/react-router"
import axios from "axios"
import { Pencil, Plus } from "lucide-react"
import { useState } from "react"

import { UsersService } from "@/client"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"
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
import { Textarea } from "@/components/ui/textarea"
import useCustomToast from "@/hooks/useCustomToast"

export const Route = createFileRoute("/_layout/log-permissions")({
  component: LogPermissionsPage,
  beforeLoad: async () => {
    const user = await UsersService.readUserMe()
    if (!user.is_superuser) {
      throw redirect({ to: "/" })
    }
  },
  head: () => ({
    meta: [{ title: "LogPermissions - Educly API" }],
  }),
})

type LogPermission = {
  id: string
  action: string
  email: string
  reason: string | null
  status: string
  granted_plan: string | null
  created_at: string | null
}

type LogPermissionsResponse = {
  data: LogPermission[]
  count: number
}

type LogPermissionPayload = {
  action: string
  email: string
  reason: string
  status: string
  granted_plan: string
}

const API_BASE_URL = import.meta.env.VITE_API_URL

function getAuthHeaders() {
  return {
    Authorization: `Bearer ${localStorage.getItem(`access_token`) || ``}`,
  }
}

function statusBadgeClass(status: string): string {
  const normalized = status.trim().toLowerCase()
  if (normalized === "success" || normalized === "approved") {
    return "border-emerald-500/40 bg-emerald-500/15 text-emerald-300"
  }
  if (normalized === "pending") {
    return "border-amber-500/40 bg-amber-500/15 text-amber-300"
  }
  if (normalized === "error" || normalized === "rejected") {
    return "border-red-500/40 bg-red-500/15 text-red-300"
  }
  return "border-muted bg-secondary text-secondary-foreground"
}

function LogPermissionsPage() {
  const queryClient = useQueryClient()
  const { showErrorToast, showSuccessToast } = useCustomToast()

  const [search, setSearch] = useState("")
  const [statusFilter, setStatusFilter] = useState("all")

  const [createOpen, setCreateOpen] = useState(false)
  const [form, setForm] = useState<LogPermissionPayload>({
    action: "assign_access:webhook",
    email: "",
    reason: "",
    status: "pending",
    granted_plan: "",
  })

  const [editOpen, setEditOpen] = useState(false)
  const [editing, setEditing] = useState<LogPermission | null>(null)
  const [editForm, setEditForm] = useState<LogPermissionPayload>({
    action: "assign_access:webhook",
    email: "",
    reason: "",
    status: "pending",
    granted_plan: "",
  })

  const query = useQuery({
    queryKey: ["log-permissions", statusFilter, search],
    queryFn: async () => {
      const params: Record<string, string | number> = { skip: 0, limit: 500 }
      if (statusFilter !== "all") params.status = statusFilter
      if (search.trim()) params.q = search.trim()
      const response = await axios.get<LogPermissionsResponse>(
        `${API_BASE_URL}/api/v1/log-permissions/`,
        {
          headers: getAuthHeaders(),
          params,
        },
      )
      return response.data
    },
  })

  const createMutation = useMutation({
    mutationFn: async (payload: LogPermissionPayload) => {
      const response = await axios.post<LogPermission>(
        `${API_BASE_URL}/api/v1/log-permissions/`,
        payload,
        { headers: getAuthHeaders() },
      )
      return response.data
    },
    onSuccess: () => {
      showSuccessToast("LogPermission criado com sucesso")
      setCreateOpen(false)
      setForm({
        action: "assign_access:webhook",
        email: "",
        reason: "",
        status: "pending",
        granted_plan: "",
      })
      queryClient.invalidateQueries({ queryKey: ["log-permissions"] })
    },
    onError: () => showErrorToast("Erro ao criar LogPermission"),
  })

  const updateMutation = useMutation({
    mutationFn: async ({
      id,
      payload,
    }: {
      id: string
      payload: LogPermissionPayload
    }) => {
      const response = await axios.patch<LogPermission>(
        `${API_BASE_URL}/api/v1/log-permissions/${id}`,
        payload,
        { headers: getAuthHeaders() },
      )
      return response.data
    },
    onSuccess: () => {
      showSuccessToast("LogPermission atualizado com sucesso")
      setEditOpen(false)
      setEditing(null)
      queryClient.invalidateQueries({ queryKey: ["log-permissions"] })
    },
    onError: () => showErrorToast("Erro ao atualizar LogPermission"),
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
        Falha ao carregar LogPermissions.
      </div>
    )
  }

  const requests = query.data?.data ?? []
  const statuses = Array.from(new Set(requests.map((item) => item.status))).sort()

  return (
    <div className="flex flex-col gap-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">LogPermissions</h1>
          <p className="text-muted-foreground">
            Historico de tentativas de concessao de permissao.
          </p>
        </div>

        <Dialog open={createOpen} onOpenChange={setCreateOpen}>
          <DialogTrigger asChild>
            <Button>
              <Plus className="mr-2 h-4 w-4" />
              Add LogPermission
            </Button>
          </DialogTrigger>
          <DialogContent className="sm:max-w-lg">
            <DialogHeader>
              <DialogTitle>Novo LogPermission</DialogTitle>
              <DialogDescription>
                Preencha os dados para criar um log manual.
              </DialogDescription>
            </DialogHeader>

            <div className="grid gap-4 py-4">
              <Input
                placeholder="Action"
                value={form.action}
                onChange={(e) => setForm({ ...form, action: e.target.value })}
              />
              <Input
                placeholder="Email"
                value={form.email}
                onChange={(e) => setForm({ ...form, email: e.target.value })}
              />
              <Textarea
                placeholder="Reason"
                value={form.reason}
                onChange={(e) => setForm({ ...form, reason: e.target.value })}
                className="resize-none"
              />
              <Input
                placeholder="Granted plan"
                value={form.granted_plan}
                onChange={(e) => setForm({ ...form, granted_plan: e.target.value })}
              />
              <Select
                value={form.status}
                onValueChange={(value) => setForm({ ...form, status: value })}
              >
                <SelectTrigger>
                  <SelectValue placeholder="Status" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="success">success</SelectItem>
                  <SelectItem value="pending">pending</SelectItem>
                  <SelectItem value="error">error</SelectItem>
                </SelectContent>
              </Select>
            </div>

            <DialogFooter>
              <DialogClose asChild>
                <Button variant="outline" disabled={createMutation.isPending}>
                  Cancelar
                </Button>
              </DialogClose>
              <Button
                onClick={() => createMutation.mutate(form)}
                disabled={
                  createMutation.isPending ||
                  !form.action.trim() ||
                  !form.email.trim()
                }
              >
                {createMutation.isPending ? "Salvando..." : "Salvar"}
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>

      <div className="grid grid-cols-1 gap-3 md:grid-cols-2">
        <Input
          placeholder="Pesquisar por acao/email/motivo/plano"
          value={search}
          onChange={(e) => setSearch(e.target.value)}
        />
        <Select value={statusFilter} onValueChange={setStatusFilter}>
          <SelectTrigger>
            <SelectValue placeholder="Filtrar por status" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">Todos os status</SelectItem>
            {statuses.map((item) => (
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
              <TableHead>Acao</TableHead>
              <TableHead>Email</TableHead>
              <TableHead>Status</TableHead>
              <TableHead>Motivo</TableHead>
              <TableHead>Plano</TableHead>
              <TableHead>Criado em</TableHead>
              <TableHead>Acoes</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {requests.length === 0 ? (
              <TableRow>
                <TableCell className="text-muted-foreground" colSpan={7}>
                  Nenhum LogPermission encontrado.
                </TableCell>
              </TableRow>
            ) : (
              requests.map((request) => (
                <TableRow key={request.id}>
                  <TableCell className="font-medium">{request.action}</TableCell>
                  <TableCell>{request.email}</TableCell>
                  <TableCell>
                    <Badge variant="secondary" className={statusBadgeClass(request.status)}>
                      {request.status}
                    </Badge>
                  </TableCell>
                  <TableCell className="max-w-sm truncate">{request.reason || "N/A"}</TableCell>
                  <TableCell>{request.granted_plan || "N/A"}</TableCell>
                  <TableCell>
                    {request.created_at
                      ? new Date(request.created_at).toLocaleString()
                      : "N/A"}
                  </TableCell>
                  <TableCell>
                    <Button
                      size="sm"
                      variant="outline"
                      onClick={() => {
                        setEditing(request)
                        setEditForm({
                          action: request.action,
                          email: request.email,
                          reason: request.reason || "",
                          status: request.status,
                          granted_plan: request.granted_plan || "",
                        })
                        setEditOpen(true)
                      }}
                    >
                      <Pencil className="mr-1 h-4 w-4" />
                      Editar
                    </Button>
                  </TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </div>

      <Dialog
        open={editOpen}
        onOpenChange={(open) => {
          setEditOpen(open)
          if (!open) setEditing(null)
        }}
      >
        <DialogContent className="sm:max-w-lg">
          <DialogHeader>
            <DialogTitle>Editar LogPermission</DialogTitle>
            <DialogDescription>
              Atualize os dados do log selecionado.
            </DialogDescription>
          </DialogHeader>

          <div className="grid gap-4 py-4">
            <Input
              placeholder="Action"
              value={editForm.action}
              onChange={(e) =>
                setEditForm({ ...editForm, action: e.target.value })
              }
            />
            <Input
              placeholder="Email"
              value={editForm.email}
              onChange={(e) =>
                setEditForm({ ...editForm, email: e.target.value })
              }
            />
            <Textarea
              placeholder="Reason"
              value={editForm.reason}
              onChange={(e) => setEditForm({ ...editForm, reason: e.target.value })}
              className="resize-none"
            />
            <Input
              placeholder="Granted plan"
              value={editForm.granted_plan}
              onChange={(e) => setEditForm({ ...editForm, granted_plan: e.target.value })}
            />
            <Select
              value={editForm.status}
              onValueChange={(value) => setEditForm({ ...editForm, status: value })}
            >
              <SelectTrigger>
                <SelectValue placeholder="Status" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="success">success</SelectItem>
                <SelectItem value="pending">pending</SelectItem>
                <SelectItem value="error">error</SelectItem>
              </SelectContent>
            </Select>
          </div>

          <DialogFooter>
            <DialogClose asChild>
              <Button variant="outline" disabled={updateMutation.isPending}>
                Cancelar
              </Button>
            </DialogClose>
            <Button
              onClick={() => {
                if (!editing) return
                updateMutation.mutate({ id: editing.id, payload: editForm })
              }}
              disabled={
                updateMutation.isPending ||
                !editForm.action.trim() ||
                !editForm.email.trim()
              }
            >
              {updateMutation.isPending ? "Salvando..." : "Salvar"}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
