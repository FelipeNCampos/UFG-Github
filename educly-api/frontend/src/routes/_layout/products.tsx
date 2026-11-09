import { useQuery } from "@tanstack/react-query"
import { createFileRoute, redirect } from "@tanstack/react-router"

import { ProductsService, UsersService } from "@/client"
import { DataTable } from "@/components/Common/DataTable"
import AddProduct from "@/components/Products/AddProduct"
import { productColumns } from "@/components/Products/columns"
import { Skeleton } from "@/components/ui/skeleton"

function getProductsQueryOptions() {
  return {
    queryFn: () => ProductsService.readProducts({ skip: 0, limit: 100 }),
    queryKey: ["products"],
  }
}

export const Route = createFileRoute("/_layout/products")({
  component: Products,
  beforeLoad: async () => {
    try {
      const user = await UsersService.readUserMe()
      if (!user.is_superuser) {
        throw redirect({
          to: "/",
        })
      }
    } catch (error) {
      const status = (error as { status?: number } | undefined)?.status
      if (status === 401 || status === 403 || status === 404) {
        localStorage.removeItem("access_token")
        throw redirect({ to: "/login" })
      }
      throw error
    }
  },
  head: () => ({
    meta: [
      {
        title: "Products - FastAPI Template",
      },
    ],
  }),
})

function ProductsTableContent() {
  const {
    data: products,
    isLoading,
    isError,
  } = useQuery({
    ...getProductsQueryOptions(),
    retry: 1,
  })

  if (isLoading) {
    return <ProductsTableSkeleton />
  }

  if (isError || !products) {
    return (
      <div className="rounded-md border border-destructive/30 bg-destructive/10 p-4 text-sm text-destructive">
        Falha ao carregar produtos. Verifique se o backend esta online e com migracoes aplicadas.
      </div>
    )
  }

  return <DataTable columns={productColumns} data={products.data} />
}

function ProductsTableSkeleton() {
  return (
    <div className="space-y-4">
      <Skeleton className="h-10 w-full" />
      <Skeleton className="h-64 w-full" />
    </div>
  )
}

function ProductsTable() {
  return <ProductsTableContent />
}

function Products() {
  return (
    <div className="flex flex-col gap-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Products</h1>
          <p className="text-muted-foreground">Manage products catalog</p>
        </div>
        <AddProduct />
      </div>
      <ProductsTable />
    </div>
  )
}
