export const PRODUCT_REFERENCES = [
  { id: "FREE", label: "Free", valor: 1 },
  { id: "BASE", label: "Base", valor: 2 },
  { id: "FREELANCER", label: "Freelancer", valor: 4 },
  { id: "HUB_IA", label: "Hub IA", valor: 8 },
  { id: "COMBO_FREELANCER_AI", label: "Combo Freelancer + Hub IA", valor: 12 },
] as const

export function getReferenceById(productId?: string | null) {
  if (!productId) return null
  return PRODUCT_REFERENCES.find((item) => item.id === productId) ?? null
}
