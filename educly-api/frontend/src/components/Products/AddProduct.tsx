import { zodResolver } from "@hookform/resolvers/zod"
import { useMutation, useQueryClient } from "@tanstack/react-query"
import { Plus, X } from "lucide-react"
import { useState } from "react"
import { useForm } from "react-hook-form"
import { z } from "zod"

import { ProductsService } from "@/client"
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
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form"
import { Input } from "@/components/ui/input"
import { LoadingButton } from "@/components/ui/loading-button"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { Textarea } from "@/components/ui/textarea"
import useCustomToast from "@/hooks/useCustomToast"
import { handleError } from "@/utils"
import { getReferenceById, PRODUCT_REFERENCES } from "./productCodeOptions"

const formSchema = z.object({
  name: z.string().min(1, { message: "Name is required" }).max(255),
  product_id: z.string().min(1, { message: "Product ID is required" }).max(255),
  product_codes: z.array(z.string().min(1).max(50)).min(1, {
    message: "Add at least one product code",
  }),
  description: z.string().max(255).optional(),
  valor: z.number().int().min(0),
})

type FormData = z.infer<typeof formSchema>

const AddProduct = () => {
  const [isOpen, setIsOpen] = useState(false)
  const [newCode, setNewCode] = useState("")
  const queryClient = useQueryClient()
  const { showSuccessToast, showErrorToast } = useCustomToast()

  const form = useForm<FormData>({
    resolver: zodResolver(formSchema),
    mode: "onBlur",
    criteriaMode: "all",
    defaultValues: {
      name: "",
      product_id: "",
      product_codes: [],
      description: "",
      valor: 0,
    },
  })

  const mutation = useMutation({
    mutationFn: (data: FormData) =>
      ProductsService.createProduct({ requestBody: data as never }),
    onSuccess: () => {
      showSuccessToast("Product created successfully")
      form.reset()
      setNewCode("")
      setIsOpen(false)
    },
    onError: handleError.bind(showErrorToast),
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: ["products"] })
    },
  })

  const addCode = () => {
    const value = newCode.trim()
    if (!value) return
    const current = form.getValues("product_codes")
    if (current.includes(value)) {
      setNewCode("")
      return
    }
    form.setValue("product_codes", [...current, value], { shouldValidate: true })
    setNewCode("")
  }

  const removeCode = (code: string) => {
    const current = form.getValues("product_codes")
    form.setValue(
      "product_codes",
      current.filter((item) => item !== code),
      { shouldValidate: true },
    )
  }

  const onSubmit = (data: FormData) => {
    mutation.mutate(data)
  }

  return (
    <Dialog open={isOpen} onOpenChange={setIsOpen}>
      <DialogTrigger asChild>
        <Button className="my-4">
          <Plus className="mr-2" />
          Add Product
        </Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>Add Product</DialogTitle>
          <DialogDescription>
            Fill in the form below to add a new product to the system.
          </DialogDescription>
        </DialogHeader>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)}>
            <div className="grid gap-4 py-4">
              <FormField
                control={form.control}
                name="name"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>
                      Name <span className="text-destructive">*</span>
                    </FormLabel>
                    <FormControl>
                      <Input placeholder="Product name" {...field} required />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="product_id"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>
                      Product ID <span className="text-destructive">*</span>
                    </FormLabel>
                    <Select
                      value={field.value}
                      onValueChange={(value) => {
                        const ref = getReferenceById(value)
                        field.onChange(value)
                        form.setValue("valor", ref?.valor ?? 0, {
                          shouldValidate: true,
                        })
                        if (!form.getValues("name")) {
                          form.setValue("name", ref?.label ?? "", {
                            shouldValidate: true,
                          })
                        }
                      }}
                    >
                      <SelectTrigger>
                        <SelectValue placeholder="Select a product reference" />
                      </SelectTrigger>
                      <SelectContent>
                        {PRODUCT_REFERENCES.map((reference) => (
                          <SelectItem key={reference.id} value={reference.id}>
                            {reference.label} ({reference.id})
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="product_codes"
                render={() => (
                  <FormItem>
                    <FormLabel>
                      Product Codes <span className="text-destructive">*</span>
                    </FormLabel>
                    <div className="flex gap-2">
                      <Input
                        placeholder="Add code and press Enter"
                        value={newCode}
                        onChange={(e) => setNewCode(e.target.value)}
                        onKeyDown={(e) => {
                          if (e.key === "Enter") {
                            e.preventDefault()
                            addCode()
                          }
                        }}
                      />
                      <Button type="button" variant="secondary" onClick={addCode}>
                        Add
                      </Button>
                    </div>
                    <div className="flex min-h-8 flex-wrap gap-2">
                      {form.watch("product_codes").map((code) => (
                        <Badge
                          key={code}
                          variant="secondary"
                          className="flex items-center gap-1"
                        >
                          {code}
                          <button
                            type="button"
                            className="cursor-pointer"
                            onClick={() => removeCode(code)}
                            aria-label={`Remove ${code}`}
                          >
                            <X className="h-3 w-3" />
                          </button>
                        </Badge>
                      ))}
                    </div>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="valor"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Valor</FormLabel>
                    <Input value={String(field.value ?? 0)} readOnly disabled />
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="description"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Description</FormLabel>
                    <FormControl>
                      <Textarea
                        placeholder="Product description"
                        className="resize-none"
                        {...field}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

            </div>

            <DialogFooter>
              <DialogClose asChild>
                <Button variant="outline" disabled={mutation.isPending}>
                  Cancel
                </Button>
              </DialogClose>
              <LoadingButton type="submit" loading={mutation.isPending}>
                Save
              </LoadingButton>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  )
}

export default AddProduct
