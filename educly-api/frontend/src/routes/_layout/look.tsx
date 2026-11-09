import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { createFileRoute, redirect } from "@tanstack/react-router"
import axios from "axios"
import { MoreVertical } from "lucide-react"
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

export const Route = createFileRoute("/_layout/look")({
  component: LookPage,
  beforeLoad: async () => {
    const user = await UsersService.readUserMe()
    if (!user.is_superuser) {
      throw redirect({ to: "/" })
    }
  },
  head: () => ({
    meta: [{ title: "Look - Educly API" }],
  }),
})

type WebhookEvent = {
  id: string
  provider: string
  payload: Record<string, unknown>
  processed: boolean
}

type WebhookResponse = {
  data: WebhookEvent[]
  count: number
}

type User = {
  id: string
  email: string
  full_name?: string | null
  last_login?: string | null
}

type UsersResponse = {
  data: User[]
  count: number
}

type Product = {
  id: string
  name: string
  valor: number
  product_color?: string | null
  product_codes?: string[] | null
}

type ProductsResponse = {
  data: Product[]
  count: number
}

type LookRow = {
  id: string
  provider: string
  nome: string
  email: string
  ultimoLogin: string
  dataCompra: string
  produto: string
  status: string
  hasUser: boolean
  userId: string | null
  productId: string | null
  productBitmap: number
  productColor: string
}

type CreateAccountForm = {
  full_name: string
  password: string
  profile_bitmap: number
  personalized_trail: boolean
  avatar_url: string
  cover_url: string
  bio: string
}

const API_BASE_URL = import.meta.env.VITE_API_URL
const BITMAP_OPTIONS = [
  { value: 1, label: "Free" },
  { value: 2, label: "Base" },
  { value: 4, label: "Freelancer" },
  { value: 8, label: "Hub IA" },
]

const DEFAULT_CREATE_FORM: CreateAccountForm = {
  full_name: "",
  password: "",
  profile_bitmap: 1,
  personalized_trail: false,
  avatar_url: "",
  cover_url: "",
  bio: "",
}

function bitmapFromSelectedOptions(selected: number[]) {
  return selected.reduce((acc, curr) => acc + curr, 0)
}

function selectedOptionsFromBitmap(bitmap: number) {
  return BITMAP_OPTIONS.filter((item) => (bitmap & item.value) === item.value).map(
    (item) => item.value,
  )
}

function getAuthHeaders() {
  return {
    Authorization: `Bearer ${localStorage.getItem("access_token") || ""}`,
  }
}

function walkValues(value: unknown, out: string[]) {
  if (value === null || value === undefined) return
  if (typeof value === "string") {
    out.push(value)
    return
  }
  if (typeof value === "number") {
    out.push(String(value))
    return
  }
  if (Array.isArray(value)) {
    for (const item of value) walkValues(item, out)
    return
  }
  if (typeof value === "object") {
    for (const val of Object.values(value as Record<string, unknown>)) {
      walkValues(val, out)
    }
  }
}

function extractEmail(payload: Record<string, unknown>): string {
  const candidates: string[] = []
  walkValues(payload, candidates)
  const email = candidates.find((value) =>
    /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value.trim()),
  )
  return email?.trim().toLowerCase() ?? ""
}

function extractName(payload: Record<string, unknown>): string {
  const data = (payload.data as Record<string, unknown> | undefined) ?? {}
  const buyer = (data.buyer as Record<string, unknown> | undefined) ?? {}
  const directName = buyer.name
  if (typeof directName === "string" && directName.trim()) return directName.trim()
  const first = typeof buyer.first_name === "string" ? buyer.first_name.trim() : ""
  const last = typeof buyer.last_name === "string" ? buyer.last_name.trim() : ""
  return `${first} ${last}`.trim()
}

function extractProductCode(
  payload: Record<string, unknown>,
  knownCodes: Set<string>,
): string {
  const candidates: string[] = []
  walkValues(payload, candidates)
  return candidates.find((value) => knownCodes.has(value.trim()))?.trim() ?? ""
}

function formatDateValue(value: unknown): string | null {
  if (typeof value === "number" && Number.isFinite(value)) {
    const normalized = value > 1e12 ? value : value * 1000
    const date = new Date(normalized)
    return Number.isNaN(date.getTime()) ? null : date.toLocaleString()
  }
  if (typeof value === "string" && value.trim()) {
    const trimmed = value.trim()
    const asNumber = Number(trimmed)
    if (!Number.isNaN(asNumber) && trimmed.length >= 10) {
      return formatDateValue(asNumber)
    }
    const date = new Date(trimmed)
    return Number.isNaN(date.getTime()) ? null : date.toLocaleString()
  }
  return null
}

function extractPurchaseDate(payload: Record<string, unknown>): string {
  const data = (payload.data as Record<string, unknown> | undefined) ?? {}
  const purchase = (data.purchase as Record<string, unknown> | undefined) ?? {}

  const candidates: unknown[] = [
    purchase.approved_date,
    purchase.order_date,
    purchase.date_next_charge,
    data.creation_date,
    payload.creation_date,
  ]

  for (const candidate of candidates) {
    const formatted = formatDateValue(candidate)
    if (formatted) return formatted
  }

  return "N/A"
}

function productColorClass(color: string): string {
  if (color === "green") return "border-emerald-500/40 bg-emerald-500/15 text-emerald-300"
  if (color === "blue") return "border-sky-500/40 bg-sky-500/15 text-sky-300"
  if (color === "purple") return "border-violet-500/40 bg-violet-500/15 text-violet-300"
  return "border-muted bg-secondary text-secondary-foreground"
}

function statusBadgeClass(status: string): string {
  const normalized = status.trim().toLowerCase()
  if (normalized === "processado") {
    return "border-emerald-500/40 bg-emerald-500/15 text-emerald-300"
  }
  if (normalized.includes("erro")) {
    return "border-red-500/40 bg-red-500/15 text-red-300"
  }
  return "border-muted bg-secondary text-secondary-foreground"
}

function LookPage() {
  const queryClient = useQueryClient()
  const { showErrorToast, showSuccessToast } = useCustomToast()

  const [search, setSearch] = useState("")
  const [providerFilter, setProviderFilter] = useState("all")
  const [statusFilter, setStatusFilter] = useState("all")

  const [createRow, setCreateRow] = useState<LookRow | null>(null)
  const [createForm, setCreateForm] = useState<CreateAccountForm>(DEFAULT_CREATE_FORM)
  const [createError, setCreateError] = useState<string | null>(null)

  const [editProductRow, setEditProductRow] = useState<LookRow | null>(null)
  const [selectedProductId, setSelectedProductId] = useState<string>("")
  const [retryNoAccountRow, setRetryNoAccountRow] = useState<LookRow | null>(null)

  const webhooksQuery = useQuery({
    queryKey: ["look", "webhooks"],
    queryFn: async () => {
      const response = await axios.get<WebhookResponse>(`${API_BASE_URL}/api/v1/webhooks/`, {
        headers: getAuthHeaders(),
        params: { skip: 0, limit: 1000 },
      })
      return response.data
    },
  })

  const usersQuery = useQuery({
    queryKey: ["look", "users"],
    queryFn: async () => {
      const response = await axios.get<UsersResponse>(`${API_BASE_URL}/api/v1/users/`, {
        headers: getAuthHeaders(),
        params: { skip: 0, limit: 1000 },
      })
      return response.data
    },
  })

  const productsQuery = useQuery({
    queryKey: ["look", "products"],
    queryFn: async () => {
      const response = await axios.get<ProductsResponse>(`${API_BASE_URL}/api/v1/products/`, {
        headers: getAuthHeaders(),
        params: { skip: 0, limit: 1000 },
      })
      return response.data
    },
  })

  const createAccountMutation = useMutation({
    mutationFn: async ({ row, form }: { row: LookRow; form: CreateAccountForm }) => {
      await axios.post(
        `${API_BASE_URL}/api/v1/users/`,
        {
          email: row.email,
          password: form.password,
          full_name: form.full_name.trim(),
          profile_bitmap: form.profile_bitmap,
          personalized_trail: form.personalized_trail,
          avatar_url: form.avatar_url.trim() || null,
          cover_url: form.cover_url.trim() || null,
          bio: form.bio.trim() || null,
        },
        { headers: getAuthHeaders() },
      )
    },
    onSuccess: () => {
      showSuccessToast("Conta criada com sucesso")
      setCreateRow(null)
      setCreateForm(DEFAULT_CREATE_FORM)
      setCreateError(null)
      queryClient.invalidateQueries({ queryKey: ["look", "users"] })
      queryClient.invalidateQueries({ queryKey: ["look", "webhooks"] })
    },
    onError: () => showErrorToast("Erro ao criar conta"),
  })

  const updatePurchaseProductMutation = useMutation({
    mutationFn: async ({ webhookId, productId }: { webhookId: string; productId: string }) => {
      await axios.patch(
        `${API_BASE_URL}/api/v1/webhooks/${webhookId}/purchase-product`,
        { product_id: productId },
        { headers: getAuthHeaders() },
      )
    },
    onSuccess: () => {
      showSuccessToast("Produto da compra atualizado")
      setEditProductRow(null)
      setSelectedProductId("")
      queryClient.invalidateQueries({ queryKey: ["look", "webhooks"] })
      queryClient.invalidateQueries({ queryKey: ["look", "users"] })
    },
    onError: () => showErrorToast("Erro ao atualizar produto da compra"),
  })

  const retryAssignMutation = useMutation({
    mutationFn: async (webhookId: string) => {
      await axios.post(
        `${API_BASE_URL}/api/v1/webhooks/${webhookId}/retry-assign-access`,
        {},
        { headers: getAuthHeaders() },
      )
    },
    onSuccess: () => {
      showSuccessToast("Atribuicao de acesso refeita com sucesso")
      queryClient.invalidateQueries({ queryKey: ["look", "webhooks"] })
      queryClient.invalidateQueries({ queryKey: ["look", "users"] })
    },
    onError: () => showErrorToast("Nao foi possivel atribuir acesso novamente"),
  })

  const webhooks = webhooksQuery.data?.data ?? []
  const users = usersQuery.data?.data ?? []
  const products = productsQuery.data?.data ?? []

  const userByEmail = useMemo(
    () => new Map(users.map((user) => [user.email.toLowerCase(), user])),
    [users],
  )

  const codeToProduct = useMemo(() => {
    const map = new Map<string, Product>()
    for (const product of products) {
      for (const code of product.product_codes ?? []) {
        const normalized = String(code).trim()
        if (normalized) map.set(normalized, product)
      }
    }
    return map
  }, [products])

  const knownCodes = useMemo(() => new Set(codeToProduct.keys()), [codeToProduct])
  const providers = Array.from(new Set(webhooks.map((item) => item.provider))).sort()

  const rows = useMemo<LookRow[]>(() => {
    return webhooks.map((event) => {
      const email = extractEmail(event.payload)
      const user = userByEmail.get(email)
      const payloadName = extractName(event.payload)
      const nome = user ? user.full_name?.trim() || payloadName || "sem nome" : "conta nao criada"

      const code = extractProductCode(event.payload, knownCodes)
      const product = code ? codeToProduct.get(code) : undefined
      const produto = product?.name ?? "produto nao encontrado"

      return {
        id: event.id,
        provider: event.provider,
        nome,
        email: email || "email nao encontrado",
        ultimoLogin: user?.last_login ? new Date(user.last_login).toLocaleString() : "nao logou",
        dataCompra: extractPurchaseDate(event.payload),
        produto,
        status: event.processed ? "processado" : "pendente",
        hasUser: Boolean(user),
        userId: user?.id ?? null,
        productId: product?.id ?? null,
        productBitmap: product?.valor ?? 1,
        productColor: product?.product_color?.trim().toLowerCase() || "green",
      }
    })
  }, [webhooks, userByEmail, knownCodes, codeToProduct])

  const filteredRows = useMemo(() => {
    const query = search.trim().toLowerCase()
    return rows.filter((row) => {
      const providerOk = providerFilter === "all" || row.provider === providerFilter
      const statusOk = statusFilter === "all" || row.status === statusFilter
      const searchOk =
        query.length === 0 ||
        row.provider.toLowerCase().includes(query) ||
        row.nome.toLowerCase().includes(query) ||
        row.email.toLowerCase().includes(query) ||
        row.dataCompra.toLowerCase().includes(query) ||
        row.produto.toLowerCase().includes(query) ||
        row.status.toLowerCase().includes(query)
      return providerOk && statusOk && searchOk
    })
  }, [rows, providerFilter, statusFilter, search])

  const openCreateAccount = (row: LookRow) => {
    const suggestedName =
      row.nome === "conta nao criada" || row.nome === "sem nome" ? "" : row.nome
    setCreateForm({
      ...DEFAULT_CREATE_FORM,
      full_name: suggestedName,
      profile_bitmap: row.productBitmap || 1,
    })
    setCreateError(null)
    setCreateRow(row)
  }

  const openEditProduct = (row: LookRow) => {
    setEditProductRow(row)
    setSelectedProductId(row.productId || "")
  }

  const submitCreateAccount = () => {
    if (!createRow) return
    if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(createRow.email)) {
      setCreateError("Email invalido para criar conta.")
      return
    }
    if (!createForm.password || createForm.password.length < 8) {
      setCreateError("Senha precisa ter pelo menos 8 caracteres.")
      return
    }
    if (!createForm.full_name.trim()) {
      setCreateError("Nome completo e obrigatorio.")
      return
    }
    if (createForm.profile_bitmap < 1 || createForm.profile_bitmap > 15) {
      setCreateError("Selecione ao menos um plano.")
      return
    }
    setCreateError(null)
    createAccountMutation.mutate({ row: createRow, form: createForm })
  }

  const submitEditProduct = () => {
    if (!editProductRow || !selectedProductId) return
    updatePurchaseProductMutation.mutate({
      webhookId: editProductRow.id,
      productId: selectedProductId,
    })
  }

  const retryAssignAccess = (row: LookRow) => {
    if (!row.hasUser) {
      setRetryNoAccountRow(row)
      return
    }
    retryAssignMutation.mutate(row.id)
  }

  if (webhooksQuery.isLoading || usersQuery.isLoading || productsQuery.isLoading) {
    return (
      <div className="space-y-4">
        <Skeleton className="h-10 w-64" />
        <Skeleton className="h-10 w-full" />
        <Skeleton className="h-72 w-full" />
      </div>
    )
  }

  if (webhooksQuery.isError || usersQuery.isError || productsQuery.isError) {
    return (
      <div className="rounded-md border border-destructive/30 bg-destructive/10 p-4 text-sm text-destructive">
        Falha ao carregar dados da aba Look.
      </div>
    )
  }

  return (
    <div className="flex flex-col gap-6">
      <div>
        <h1 className="text-2xl font-bold tracking-tight">Look</h1>
        <p className="text-muted-foreground">
          Visao consolidada de webhooks com usuario e produto.
        </p>
      </div>

      <div className="grid grid-cols-1 gap-3 md:grid-cols-3">
        <Input
          placeholder="Pesquisar por provider, nome, email, produto ou status"
          value={search}
          onChange={(e) => setSearch(e.target.value)}
        />
        <Select value={providerFilter} onValueChange={setProviderFilter}>
          <SelectTrigger>
            <SelectValue placeholder="Filtrar por provider" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">Todos os providers</SelectItem>
            {providers.map((provider) => (
              <SelectItem key={provider} value={provider}>
                {provider}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
        <Select value={statusFilter} onValueChange={setStatusFilter}>
          <SelectTrigger>
            <SelectValue placeholder="Filtrar por status" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">Todos os status</SelectItem>
            <SelectItem value="pendente">pendente</SelectItem>
            <SelectItem value="processado">processado</SelectItem>
          </SelectContent>
        </Select>
      </div>

      <div className="rounded-md border">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Provider</TableHead>
              <TableHead>Nome</TableHead>
              <TableHead>Email</TableHead>
              <TableHead>Ultimo login</TableHead>
              <TableHead>Data da compra</TableHead>
              <TableHead>Produto</TableHead>
              <TableHead>Status</TableHead>
              <TableHead className="text-right">Acoes</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {filteredRows.length === 0 ? (
              <TableRow>
                <TableCell colSpan={8} className="text-muted-foreground">
                  Nenhum registro encontrado.
                </TableCell>
              </TableRow>
            ) : (
              filteredRows.map((row) => (
                <TableRow key={row.id}>
                  <TableCell>{row.provider}</TableCell>
                  <TableCell>{row.nome}</TableCell>
                  <TableCell>{row.email}</TableCell>
                  <TableCell>{row.ultimoLogin}</TableCell>
                  <TableCell>{row.dataCompra}</TableCell>
                  <TableCell>
                    <Badge variant="secondary" className={productColorClass(row.productColor)}>
                      {row.produto}
                    </Badge>
                  </TableCell>
                  <TableCell>
                    <Badge variant="secondary" className={statusBadgeClass(row.status)}>
                      {row.status}
                    </Badge>
                  </TableCell>
                  <TableCell className="text-right">
                    <DropdownMenu>
                      <DropdownMenuTrigger asChild>
                        <Button variant="ghost" size="icon">
                          <MoreVertical className="h-4 w-4" />
                          <span className="sr-only">Abrir acoes</span>
                        </Button>
                      </DropdownMenuTrigger>
                      <DropdownMenuContent align="end">
                        <DropdownMenuItem
                          disabled={row.hasUser || !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(row.email)}
                          onClick={() => openCreateAccount(row)}
                        >
                          Criar conta
                        </DropdownMenuItem>
                        <DropdownMenuItem disabled={!row.hasUser} onClick={() => openEditProduct(row)}>
                          Editar produto
                        </DropdownMenuItem>
                        <DropdownMenuItem onClick={() => retryAssignAccess(row)}>
                          Tentar atribuir cargo novamente
                        </DropdownMenuItem>
                      </DropdownMenuContent>
                    </DropdownMenu>
                  </TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </div>

      <Dialog
        open={createRow !== null}
        onOpenChange={(open) => {
          if (!open) {
            setCreateRow(null)
            setCreateForm(DEFAULT_CREATE_FORM)
            setCreateError(null)
          }
        }}
      >
        <DialogContent className="sm:max-w-lg">
          <DialogHeader>
            <DialogTitle>Criar conta</DialogTitle>
            <DialogDescription>
              Criar usuario e perfil vinculados ao email do webhook.
            </DialogDescription>
          </DialogHeader>

          <div className="grid gap-3 py-3">
            <Input value={createRow?.email ?? ""} disabled />
            <Input
              placeholder="Nome completo"
              value={createForm.full_name}
              onChange={(e) => setCreateForm({ ...createForm, full_name: e.target.value })}
            />
            <Input
              type="password"
              placeholder="Senha (minimo 8 caracteres)"
              value={createForm.password}
              onChange={(e) => setCreateForm({ ...createForm, password: e.target.value })}
            />
            <div className="grid gap-2">
              <label className="text-sm font-medium">Plano</label>
              <div className="grid grid-cols-2 gap-2 rounded-md border p-3">
                {BITMAP_OPTIONS.map((option) => {
                  const selected = selectedOptionsFromBitmap(createForm.profile_bitmap)
                  const checked = selected.includes(option.value)
                  return (
                    <label key={option.value} className="flex items-center gap-2 text-sm">
                      <Checkbox
                        checked={checked}
                        onCheckedChange={(value) => {
                          const next = value === true
                            ? Array.from(new Set([...selected, option.value]))
                            : selected.filter((item) => item !== option.value)
                          setCreateForm({
                            ...createForm,
                            profile_bitmap: bitmapFromSelectedOptions(next),
                          })
                        }}
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
              onChange={(e) => setCreateForm({ ...createForm, avatar_url: e.target.value })}
            />
            <Input
              placeholder="Cover URL"
              value={createForm.cover_url}
              onChange={(e) => setCreateForm({ ...createForm, cover_url: e.target.value })}
            />
            <Textarea
              placeholder="Bio"
              value={createForm.bio}
              onChange={(e) => setCreateForm({ ...createForm, bio: e.target.value })}
              className="resize-none"
            />
          </div>

          {createError ? <p className="text-sm text-destructive">{createError}</p> : null}

          <DialogFooter>
            <DialogClose asChild>
              <Button variant="outline" disabled={createAccountMutation.isPending}>
                Cancelar
              </Button>
            </DialogClose>
            <Button onClick={submitCreateAccount} disabled={createAccountMutation.isPending}>
              {createAccountMutation.isPending ? "Criando..." : "Criar conta"}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <Dialog
        open={editProductRow !== null}
        onOpenChange={(open) => {
          if (!open) {
            setEditProductRow(null)
            setSelectedProductId("")
          }
        }}
      >
        <DialogContent className="sm:max-w-lg">
          <DialogHeader>
            <DialogTitle>Editar produto da compra</DialogTitle>
            <DialogDescription>
              Remove o valor antigo do bitmap e aplica o novo produto escolhido.
            </DialogDescription>
          </DialogHeader>

          <div className="grid gap-3 py-3">
            <Input value={editProductRow?.email ?? ""} disabled />
            <Select value={selectedProductId} onValueChange={setSelectedProductId}>
              <SelectTrigger>
                <SelectValue placeholder="Selecione o novo produto" />
              </SelectTrigger>
              <SelectContent>
                {products.map((product) => (
                  <SelectItem key={product.id} value={product.id}>
                    {product.name} (bitmap {product.valor})
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          <DialogFooter>
            <DialogClose asChild>
              <Button variant="outline" disabled={updatePurchaseProductMutation.isPending}>
                Cancelar
              </Button>
            </DialogClose>
            <Button
              onClick={submitEditProduct}
              disabled={updatePurchaseProductMutation.isPending || !selectedProductId}
            >
              {updatePurchaseProductMutation.isPending ? "Salvando..." : "Salvar produto"}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <Dialog
        open={retryNoAccountRow !== null}
        onOpenChange={(open) => {
          if (!open) setRetryNoAccountRow(null)
        }}
      >
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle>Conta nao encontrada</DialogTitle>
            <DialogDescription>
              Nao existe conta para o email deste webhook. Crie a conta para depois tentar atribuir o cargo.
            </DialogDescription>
          </DialogHeader>
          <div className="py-2">
            <Input value={retryNoAccountRow?.email ?? ""} disabled />
          </div>
          <DialogFooter>
            <DialogClose asChild>
              <Button variant="outline">Fechar</Button>
            </DialogClose>
            <Button
              onClick={() => {
                if (!retryNoAccountRow) return
                const targetRow = retryNoAccountRow
                setRetryNoAccountRow(null)
                openCreateAccount(targetRow)
              }}
            >
              Criar conta
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
