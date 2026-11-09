import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { createFileRoute, redirect } from "@tanstack/react-router"
import axios from "axios"
import { MoreVertical, Plus } from "lucide-react"
import { useMemo, useState } from "react"

import { UsersService } from "@/client"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Checkbox } from "@/components/ui/checkbox"
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
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
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

export const Route = createFileRoute("/_layout/perfils")({
  component: PerfilsPage,
  beforeLoad: async () => {
    const user = await UsersService.readUserMe()
    if (!user.is_superuser) {
      throw redirect({ to: "/" })
    }
  },
  head: () => ({
    meta: [{ title: "Perfils - Educly API" }],
  }),
})

type Perfil = {
  id: string
  user_id: string
  profile_bitmap: number
  personalized_trail: boolean
  avatar_url: string | null
  cover_url: string | null
  bio: string | null
}

type PerfilsResponse = {
  data: Perfil[]
  count: number
}

type PerfilFormData = {
  user_id: string
  profile_bitmap: number
  personalized_trail: boolean
  avatar_url: string
  cover_url: string
  bio: string
}

type BitmapOption = {
  value: number
  label: string
}

const API_BASE_URL = import.meta.env.VITE_API_URL
const PERFILES_QUERY_KEY = ["perfils"]
const USERS_QUERY_KEY = ["users", "for-perfils"]

const BITMAP_OPTIONS: BitmapOption[] = [
  { value: 1, label: "Free" },
  { value: 2, label: "Base" },
  { value: 4, label: "Freelancer" },
  { value: 8, label: "Hub IA" },
]

const DEFAULT_FORM: PerfilFormData = {
  user_id: "",
  profile_bitmap: 1,
  personalized_trail: false,
  avatar_url: "",
  cover_url: "",
  bio: "",
}

const getAuthHeaders = () => ({
  Authorization: `Bearer ${localStorage.getItem("access_token") || ""}`,
})

const labelsFromBitmap = (value: number) => {
  const labels = BITMAP_OPTIONS.filter((item) => (value & item.value) === item.value).map(
    (item) => item.label,
  )
  if (labels.length === 0) return ["Sem produto"]
  return labels
}

const bitmapFromSelectedOptions = (selected: number[]) => {
  return selected.reduce((acc, curr) => acc + curr, 0)
}

const selectedOptionsFromBitmap = (bitmap: number) => {
  return BITMAP_OPTIONS.filter((item) => (bitmap & item.value) === item.value).map(
    (item) => item.value,
  )
}

const getPerfils = async (): Promise<PerfilsResponse> => {
  const response = await axios.get<PerfilsResponse>(`${API_BASE_URL}/api/v1/perfils/`, {
    headers: getAuthHeaders(),
    params: { skip: 0, limit: 200 },
  })
  return response.data
}

const createPerfil = async (payload: PerfilFormData): Promise<Perfil> => {
  const response = await axios.post<Perfil>(
    `${API_BASE_URL}/api/v1/perfils/`,
    payload,
    { headers: getAuthHeaders() },
  )
  return response.data
}

const updatePerfil = async ({
  id,
  payload,
}: {
  id: string
  payload: Omit<PerfilFormData, "user_id">
}): Promise<Perfil> => {
  const response = await axios.patch<Perfil>(
    `${API_BASE_URL}/api/v1/perfils/${id}`,
    payload,
    { headers: getAuthHeaders() },
  )
  return response.data
}

const deletePerfil = async (id: string): Promise<void> => {
  await axios.delete(`${API_BASE_URL}/api/v1/perfils/${id}`, {
    headers: getAuthHeaders(),
  })
}

function validateCreateForm(data: PerfilFormData): string | null {
  if (!data.user_id) return "Selecione um usuario."
  if (data.profile_bitmap < 1 || data.profile_bitmap > 15) {
    return "Selecione ao menos um produto."
  }
  return null
}

function validateEditForm(data: Omit<PerfilFormData, "user_id">): string | null {
  if (data.profile_bitmap < 1 || data.profile_bitmap > 15) {
    return "Selecione ao menos um produto."
  }
  return null
}

function PerfilsPage() {
  const queryClient = useQueryClient()
  const { showErrorToast, showSuccessToast } = useCustomToast()

  const [search, setSearch] = useState("")
  const [bitmapFilters, setBitmapFilters] = useState<number[]>([])

  const [isCreateOpen, setIsCreateOpen] = useState(false)
  const [createForm, setCreateForm] = useState<PerfilFormData>(DEFAULT_FORM)
  const [createError, setCreateError] = useState<string | null>(null)

  const [editingPerfil, setEditingPerfil] = useState<Perfil | null>(null)
  const [editForm, setEditForm] = useState<Omit<PerfilFormData, "user_id">>({
    profile_bitmap: 1,
    personalized_trail: false,
    avatar_url: "",
    cover_url: "",
    bio: "",
  })
  const [editError, setEditError] = useState<string | null>(null)

  const [deletingPerfil, setDeletingPerfil] = useState<Perfil | null>(null)

  const perfilsQuery = useQuery({
    queryKey: PERFILES_QUERY_KEY,
    queryFn: getPerfils,
  })

  const usersQuery = useQuery({
    queryKey: USERS_QUERY_KEY,
    queryFn: () => UsersService.readUsers({ skip: 0, limit: 500 }),
  })

  const createMutation = useMutation({
    mutationFn: createPerfil,
    onSuccess: () => {
      showSuccessToast("Perfil criado com sucesso")
      setCreateForm(DEFAULT_FORM)
      setCreateError(null)
      setIsCreateOpen(false)
      queryClient.invalidateQueries({ queryKey: PERFILES_QUERY_KEY })
    },
    onError: () => showErrorToast("Erro ao criar perfil"),
  })

  const updateMutation = useMutation({
    mutationFn: updatePerfil,
    onSuccess: () => {
      showSuccessToast("Perfil atualizado com sucesso")
      setEditingPerfil(null)
      setEditError(null)
      queryClient.invalidateQueries({ queryKey: PERFILES_QUERY_KEY })
    },
    onError: () => showErrorToast("Erro ao atualizar perfil"),
  })

  const deleteMutation = useMutation({
    mutationFn: deletePerfil,
    onSuccess: () => {
      showSuccessToast("Perfil deletado com sucesso")
      setDeletingPerfil(null)
      queryClient.invalidateQueries({ queryKey: PERFILES_QUERY_KEY })
    },
    onError: () => showErrorToast("Erro ao deletar perfil"),
  })

  const users = usersQuery.data?.data ?? []
  const userById = useMemo(() => {
    return new Map(users.map((user) => [user.id, user]))
  }, [users])

  const availableUsersForCreate = useMemo(() => {
    const perfils = perfilsQuery.data?.data ?? []
    const usedIds = new Set(perfils.map((perfil) => perfil.user_id))
    return users.filter((user) => !usedIds.has(user.id))
  }, [users, perfilsQuery.data?.data])

  const filteredPerfils = useMemo(() => {
    const perfils = perfilsQuery.data?.data ?? []
    const query = search.trim().toLowerCase()

    return perfils.filter((perfil) => {
      const linkedUser = userById.get(perfil.user_id)
      const fullName = linkedUser?.full_name?.toLowerCase() ?? ""
      const email = linkedUser?.email?.toLowerCase() ?? ""
      const bitmapLabel = labelsFromBitmap(perfil.profile_bitmap).join(" ").toLowerCase()

      const matchesSearch =
        query.length === 0 ||
        fullName.includes(query) ||
        email.includes(query) ||
        bitmapLabel.includes(query)

      const matchesBitmap =
        bitmapFilters.length === 0 || bitmapFilters.includes(perfil.profile_bitmap)

      return matchesSearch && matchesBitmap
    })
  }, [perfilsQuery.data?.data, search, bitmapFilters, userById])

  const toggleBitmapFilter = (value: number, checked: boolean) => {
    setBitmapFilters((current) => {
      if (checked) return Array.from(new Set([...current, value]))
      return current.filter((item) => item !== value)
    })
  }

  const bitmapFilterLabel =
    bitmapFilters.length === 0
      ? "Todos os planos"
      : BITMAP_OPTIONS.filter((option) => bitmapFilters.includes(option.value))
          .map((option) => option.label)
          .join(", ")

  const openEdit = (perfil: Perfil) => {
    setEditingPerfil(perfil)
    setEditForm({
      profile_bitmap: perfil.profile_bitmap,
      personalized_trail: perfil.personalized_trail,
      avatar_url: perfil.avatar_url ?? "",
      cover_url: perfil.cover_url ?? "",
      bio: perfil.bio ?? "",
    })
    setEditError(null)
  }

  const handleCreate = () => {
    const validationError = validateCreateForm(createForm)
    if (validationError) {
      setCreateError(validationError)
      return
    }
    createMutation.mutate(createForm)
  }

  const handleEdit = () => {
    if (!editingPerfil) return
    const validationError = validateEditForm(editForm)
    if (validationError) {
      setEditError(validationError)
      return
    }
    updateMutation.mutate({ id: editingPerfil.id, payload: editForm })
  }

  const handleDelete = () => {
    if (!deletingPerfil) return
    deleteMutation.mutate(deletingPerfil.id)
  }

  const toggleCreateBitmapOption = (value: number, checked: boolean) => {
    const selected = selectedOptionsFromBitmap(createForm.profile_bitmap)
    const next = checked
      ? Array.from(new Set([...selected, value]))
      : selected.filter((item) => item !== value)
    setCreateForm({ ...createForm, profile_bitmap: bitmapFromSelectedOptions(next) })
  }

  const toggleEditBitmapOption = (value: number, checked: boolean) => {
    const selected = selectedOptionsFromBitmap(editForm.profile_bitmap)
    const next = checked
      ? Array.from(new Set([...selected, value]))
      : selected.filter((item) => item !== value)
    setEditForm({ ...editForm, profile_bitmap: bitmapFromSelectedOptions(next) })
  }

  if (perfilsQuery.isLoading || usersQuery.isLoading) {
    return (
      <div className="space-y-4">
        <Skeleton className="h-10 w-56" />
        <Skeleton className="h-10 w-full" />
        <Skeleton className="h-72 w-full" />
      </div>
    )
  }

  if (perfilsQuery.isError || usersQuery.isError) {
    return (
      <div className="rounded-md border border-destructive/30 bg-destructive/10 p-4 text-sm text-destructive">
        Falha ao carregar perfis.
      </div>
    )
  }

  return (
    <div className="flex flex-col gap-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Perfils</h1>
          <p className="text-muted-foreground">
            Dados sincronizados com o backend e usuario vinculado exibido.
          </p>
        </div>

        <Dialog
          open={isCreateOpen}
          onOpenChange={(open) => {
            setIsCreateOpen(open)
            if (!open) {
              setCreateForm(DEFAULT_FORM)
              setCreateError(null)
            }
          }}
        >
          <DialogTrigger asChild>
            <Button className="my-4">
              <Plus className="mr-2 h-4 w-4" />
              Create Perfil
            </Button>
          </DialogTrigger>
          <DialogContent className="sm:max-w-lg">
            <DialogHeader>
              <DialogTitle>Criar perfil</DialogTitle>
              <DialogDescription>
                Selecione o usuario vinculado e o plano (bitmap).
              </DialogDescription>
            </DialogHeader>

            <div className="grid gap-4 py-4">
              <div className="grid gap-2">
                <label className="text-sm font-medium">Usuario</label>
                <Select
                  disabled={availableUsersForCreate.length === 0}
                  value={createForm.user_id}
                  onValueChange={(value) =>
                    setCreateForm({ ...createForm, user_id: value })
                  }
                >
                  <SelectTrigger>
                    <SelectValue placeholder="Selecione o usuario" />
                  </SelectTrigger>
                  <SelectContent>
                    {availableUsersForCreate.length === 0 ? (
                      <SelectItem value="__no_users_available__" disabled>
                        Nenhum usuario disponivel (todos ja possuem perfil)
                      </SelectItem>
                    ) : (
                      availableUsersForCreate.map((user) => (
                        <SelectItem key={user.id} value={user.id}>
                          {user.full_name || user.email} ({user.email})
                        </SelectItem>
                      ))
                    )}
                  </SelectContent>
                </Select>
                {availableUsersForCreate.length === 0 ? (
                  <p className="text-xs text-muted-foreground">
                    Crie um novo usuario ou remova um perfil existente para liberar a selecao.
                  </p>
                ) : null}
              </div>

              <div className="grid gap-2">
                <label className="text-sm font-medium">Produtos</label>
                <div className="grid grid-cols-2 gap-2 rounded-md border p-3">
                  {BITMAP_OPTIONS.map((option) => {
                    const selected = selectedOptionsFromBitmap(createForm.profile_bitmap)
                    const checked = selected.includes(option.value)
                    return (
                      <label key={option.value} className="flex items-center gap-2 text-sm">
                        <Checkbox
                          checked={checked}
                          onCheckedChange={(value) =>
                            toggleCreateBitmapOption(option.value, value === true)
                          }
                        />
                        {option.label}
                      </label>
                    )
                  })}
                </div>
              </div>

              <label className="flex items-center gap-2 text-sm">
                <Checkbox
                  checked={createForm.personalized_trail}
                  onCheckedChange={(checked) =>
                    setCreateForm({
                      ...createForm,
                      personalized_trail: checked === true,
                    })
                  }
                />
                Trilha personalizada
              </label>

              <Input
                placeholder="Avatar URL"
                value={createForm.avatar_url}
                onChange={(e) =>
                  setCreateForm({ ...createForm, avatar_url: e.target.value })
                }
              />
              <Input
                placeholder="Cover URL"
                value={createForm.cover_url}
                onChange={(e) =>
                  setCreateForm({ ...createForm, cover_url: e.target.value })
                }
              />
              <Textarea
                placeholder="Bio"
                value={createForm.bio}
                onChange={(e) => setCreateForm({ ...createForm, bio: e.target.value })}
                className="resize-none"
              />
            </div>

            {createError ? (
              <p className="text-sm text-destructive">{createError}</p>
            ) : null}

            <DialogFooter>
              <DialogClose asChild>
                <Button variant="outline" disabled={createMutation.isPending}>
                  Cancelar
                </Button>
              </DialogClose>
              <Button
                onClick={handleCreate}
                disabled={
                  createMutation.isPending || availableUsersForCreate.length === 0
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
          placeholder="Buscar por usuario/email/plano"
          value={search}
          onChange={(e) => setSearch(e.target.value)}
        />
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="outline" className="justify-start">
              {bitmapFilterLabel}
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="start" className="w-64">
            {BITMAP_OPTIONS.map((option) => {
              const checked = bitmapFilters.includes(option.value)
              return (
                <DropdownMenuItem
                  key={option.value}
                  onSelect={(e) => e.preventDefault()}
                  onClick={() => toggleBitmapFilter(option.value, !checked)}
                  className="flex items-center gap-2"
                >
                  <Checkbox checked={checked} />
                  <span>{option.label}</span>
                </DropdownMenuItem>
              )
            })}
          </DropdownMenuContent>
        </DropdownMenu>
      </div>

      <p className="text-sm text-muted-foreground">
        {filteredPerfils.length} usuarios encontrados
      </p>

      <div className="rounded-md border">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Usuario</TableHead>
              <TableHead>Email</TableHead>
              <TableHead>Ultimo login</TableHead>
              <TableHead>Plano</TableHead>
              <TableHead>Trilha</TableHead>
              <TableHead>Bio</TableHead>
              <TableHead className="text-right">Acoes</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {filteredPerfils.length === 0 ? (
              <TableRow>
                <TableCell
                  colSpan={7}
                  className="h-24 text-center text-muted-foreground"
                >
                  Nenhum perfil encontrado.
                </TableCell>
              </TableRow>
            ) : (
              filteredPerfils.map((perfil) => {
                const user = userById.get(perfil.user_id)
                return (
                  <TableRow key={perfil.id}>
                    <TableCell className="font-medium">
                      {user?.full_name || "Usuario sem nome"}
                    </TableCell>
                    <TableCell>{user?.email || perfil.user_id}</TableCell>
                    <TableCell>
                      {(user as { last_login?: string | null } | undefined)?.last_login
                        ? new Date(
                            (user as { last_login?: string | null }).last_login as string,
                          ).toLocaleString()
                        : "N/A"}
                    </TableCell>
                    <TableCell>
                      <Badge
                        variant="secondary"
                        className={
                          perfil.profile_bitmap !== 1
                            ? "border-emerald-500/40 bg-emerald-500/15 text-emerald-300"
                            : undefined
                        }
                      >
                        {labelsFromBitmap(perfil.profile_bitmap).join(" + ")}
                      </Badge>
                    </TableCell>
                    <TableCell>
                      {perfil.personalized_trail ? "Personalizada" : "Padrao"}
                    </TableCell>
                    <TableCell className="max-w-[260px] truncate">
                      {perfil.bio || "N/A"}
                    </TableCell>
                    <TableCell className="text-right">
                      <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                          <Button variant="ghost" size="icon">
                            <MoreVertical className="h-4 w-4" />
                            <span className="sr-only">Abrir menu</span>
                          </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="end">
                          <DropdownMenuItem onClick={() => openEdit(perfil)}>
                            Editar
                          </DropdownMenuItem>
                          <DropdownMenuItem
                            className="text-destructive focus:text-destructive"
                            onClick={() => setDeletingPerfil(perfil)}
                          >
                            Deletar
                          </DropdownMenuItem>
                        </DropdownMenuContent>
                      </DropdownMenu>
                    </TableCell>
                  </TableRow>
                )
              })
            )}
          </TableBody>
        </Table>
      </div>

      <Dialog
        open={editingPerfil !== null}
        onOpenChange={(open) => {
          if (!open) {
            setEditingPerfil(null)
            setEditError(null)
          }
        }}
      >
        <DialogContent className="sm:max-w-lg">
          <DialogHeader>
            <DialogTitle>Editar perfil</DialogTitle>
            <DialogDescription>
              {editingPerfil
                ? `Usuario vinculado: ${
                    userById.get(editingPerfil.user_id)?.email || editingPerfil.user_id
                  }`
                : "Atualize os dados do perfil."}
            </DialogDescription>
          </DialogHeader>

          <div className="grid gap-4 py-4">
            <div className="grid gap-2">
              <label className="text-sm font-medium">Produtos</label>
              <div className="grid grid-cols-2 gap-2 rounded-md border p-3">
                {BITMAP_OPTIONS.map((option) => {
                  const selected = selectedOptionsFromBitmap(editForm.profile_bitmap)
                  const checked = selected.includes(option.value)
                  return (
                    <label key={option.value} className="flex items-center gap-2 text-sm">
                      <Checkbox
                        checked={checked}
                        onCheckedChange={(value) =>
                          toggleEditBitmapOption(option.value, value === true)
                        }
                      />
                      {option.label}
                    </label>
                  )
                })}
              </div>
            </div>

            <label className="flex items-center gap-2 text-sm">
              <Checkbox
                checked={editForm.personalized_trail}
                onCheckedChange={(checked) =>
                  setEditForm({
                    ...editForm,
                    personalized_trail: checked === true,
                  })
                }
              />
              Trilha personalizada
            </label>

            <Input
              placeholder="Avatar URL"
              value={editForm.avatar_url}
              onChange={(e) =>
                setEditForm({ ...editForm, avatar_url: e.target.value })
              }
            />
            <Input
              placeholder="Cover URL"
              value={editForm.cover_url}
              onChange={(e) =>
                setEditForm({ ...editForm, cover_url: e.target.value })
              }
            />
            <Textarea
              placeholder="Bio"
              value={editForm.bio}
              onChange={(e) => setEditForm({ ...editForm, bio: e.target.value })}
              className="resize-none"
            />
          </div>

          {editError ? <p className="text-sm text-destructive">{editError}</p> : null}

          <DialogFooter>
            <DialogClose asChild>
              <Button variant="outline" disabled={updateMutation.isPending}>
                Cancelar
              </Button>
            </DialogClose>
            <Button onClick={handleEdit} disabled={updateMutation.isPending}>
              {updateMutation.isPending ? "Salvando..." : "Salvar alteracoes"}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <Dialog
        open={deletingPerfil !== null}
        onOpenChange={(open) => {
          if (!open) setDeletingPerfil(null)
        }}
      >
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Deletar perfil</DialogTitle>
            <DialogDescription>
              {`Tem certeza que deseja deletar o perfil de ${
                deletingPerfil
                  ? userById.get(deletingPerfil.user_id)?.email || deletingPerfil.user_id
                  : "usuario"
              }? Esta acao nao pode ser desfeita.`}
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <DialogClose asChild>
              <Button variant="outline" disabled={deleteMutation.isPending}>
                Cancelar
              </Button>
            </DialogClose>
            <Button
              variant="destructive"
              onClick={handleDelete}
              disabled={deleteMutation.isPending}
            >
              {deleteMutation.isPending ? "Deletando..." : "Deletar"}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
