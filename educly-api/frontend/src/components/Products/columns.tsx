import type { ColumnDef } from "@tanstack/react-table"

import type { ProductPublic } from "@/client"
import { Badge } from "@/components/ui/badge"
import { ProductActionsMenu } from "./ProductActionsMenu"

export const productColumns: ColumnDef<ProductPublic>[] = [
  {
    accessorKey: "name",
    header: "Name",
    cell: ({ row }) => <span className="font-medium">{row.original.name}</span>,
  },
  {
    accessorKey: "product_id",
    header: "Product ID",
    cell: ({ row }) => (
      <span className="text-sm font-mono text-muted-foreground">
        {row.original.product_id}
      </span>
    ),
  },
  {
    accessorKey: "product_codes",
    header: "Product Codes",
    cell: ({ row }) => {
      const product = row.original as unknown as { product_codes?: string[] | string | null }
      const codes = Array.isArray(product.product_codes)
        ? product.product_codes
        : typeof product.product_codes === "string"
          ? [product.product_codes]
          : []
      if (codes.length === 0) return <span className="text-muted-foreground">N/A</span>
      return (
        <div className="flex flex-wrap gap-1">
          {codes.map((code) => (
            <Badge key={code} variant="secondary">
              {code}
            </Badge>
          ))}
        </div>
      )
    },
  },
  {
    accessorKey: "valor",
    header: "Valor",
    cell: ({ row }) => (
      <span className="text-sm">
        {(row.original as unknown as { valor?: number | null }).valor ?? "N/A"}
      </span>
    ),
  },
  {
    accessorKey: "description",
    header: "Description",
    cell: ({ row }) => (
      <span className="text-sm text-muted-foreground line-clamp-1">
        {row.original.description || "N/A"}
      </span>
    ),
  },
  {
    accessorKey: "created_at",
    header: "Created",
    cell: ({ row }) => {
      const date = row.original.created_at
      if (!date) return <span className="text-muted-foreground">N/A</span>
      return (
        <span className="text-sm text-muted-foreground">
          {new Date(date).toLocaleDateString()}
        </span>
      )
    },
  },
  {
    id: "actions",
    header: () => <span className="sr-only">Actions</span>,
    cell: ({ row }) => (
      <div className="flex justify-end">
        <ProductActionsMenu product={row.original} />
      </div>
    ),
  },
]
