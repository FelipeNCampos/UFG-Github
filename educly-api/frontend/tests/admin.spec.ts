import { expect, test } from "@playwright/test"
import { firstSuperuser, firstSuperuserPassword } from "./config.ts"
import { createUser } from "./utils/privateApi"
import { randomEmail, randomPassword } from "./utils/random"
import { logInUser } from "./utils/user"

test("Admin page is accessible and shows correct title", async ({ page }) => {
  await page.goto("/admin")
  await expect(page.getByRole("heading", { name: "Users" })).toBeVisible()
  await expect(
    page.getByText("Manage user accounts and permissions"),
  ).toBeVisible()
})

test("Add User button is visible", async ({ page }) => {
  await page.goto("/admin")
  await expect(page.getByRole("button", { name: "Add User" })).toBeVisible()
})

test.describe("Admin user management", () => {
  test("Create a new user successfully", async ({ page }) => {
    await page.goto("/admin")

    const email = randomEmail()
    const password = randomPassword()
    const fullName = "Test User Admin"

    await page.getByRole("button", { name: "Add User" }).click()

    await page.getByPlaceholder("Email").fill(email)
    await page.getByPlaceholder("Full name").fill(fullName)
    await page.getByPlaceholder("Password").first().fill(password)
    await page.getByPlaceholder("Password").last().fill(password)

    await page.getByRole("button", { name: "Save" }).click()

    await expect(page.getByText("User created successfully")).toBeVisible({
      timeout: 20000,
    })

    await expect(page.getByRole("dialog")).not.toBeVisible()

    const userRow = page.getByRole("row").filter({ hasText: email })
    await expect(userRow).toBeVisible()
  })

  test("Create a superuser", async ({ page }) => {
    await page.goto("/admin")

    const email = randomEmail()
    const password = randomPassword()
    const fullName = "Test Superuser"

    await page.getByRole("button", { name: "Add User" }).click()

    await page.getByPlaceholder("Email").fill(email)
    await page.getByPlaceholder("Full name").fill(fullName)
    await page.getByPlaceholder("Password").first().fill(password)
    await page.getByPlaceholder("Password").last().fill(password)
    await page.getByLabel("Is superuser?").check()
    await page.getByLabel("Is active?").check()

    await page.getByRole("button", { name: "Save" }).click()

    await expect(page.getByText("User created successfully")).toBeVisible({
      timeout: 20000,
    })

    await expect(page.getByRole("dialog")).not.toBeVisible()

    const userRow = page.getByRole("row").filter({ hasText: email })
    await expect(
      userRow.getByRole("cell").filter({ hasText: "Superuser" }).first(),
    ).toBeVisible()
  })

  test("Edit a user successfully", async ({ page }) => {
    await page.goto("/admin")

    const email = randomEmail()
    const password = randomPassword()
    const originalName = "Original Name"
    const updatedName = "Updated Name"

    await page.getByRole("button", { name: "Add User" }).click()
    await page.getByPlaceholder("Email").fill(email)
    await page.getByPlaceholder("Full name").fill(originalName)
    await page.getByPlaceholder("Password").first().fill(password)
    await page.getByPlaceholder("Password").last().fill(password)
    await page.getByRole("button", { name: "Save" }).click()

    await expect(page.getByText("User created successfully")).toBeVisible()
    await expect(page.getByRole("dialog")).not.toBeVisible()

    const userRow = page.getByRole("row").filter({ hasText: email })
    await userRow.getByRole("button").click()

    await page.getByRole("menuitem", { name: "Edit User" }).click()

    await page.getByPlaceholder("Full name").fill(updatedName)
    await page.getByRole("button", { name: "Save" }).click()

    await expect(page.getByText("User updated successfully")).toBeVisible()
    await expect(
      page.getByRole("row").filter({ hasText: email }).getByText(updatedName),
    ).toBeVisible()
  })

  test("Delete a user successfully", async ({ page }) => {
    await page.goto("/admin")

    const email = randomEmail()
    const password = randomPassword()
    const fullName = "Test User Delete"

    await page.getByRole("button", { name: "Add User" }).click()
    await page.getByPlaceholder("Email").fill(email)
    await page.getByPlaceholder("Full name").fill(fullName)
    await page.getByPlaceholder("Password").first().fill(password)
    await page.getByPlaceholder("Password").last().fill(password)
    await page.getByRole("button", { name: "Save" }).click()

    await expect(page.getByText("User created successfully")).toBeVisible({
      timeout: 10000,
    })

    await expect(page.getByRole("dialog")).not.toBeVisible()

    const userRow = page.getByRole("row").filter({ hasText: email })
    await userRow.getByRole("button").click()

    await page.getByRole("menuitem", { name: "Delete User" }).click()

    await page.getByRole("button", { name: "Delete" }).click()

    await expect(
      page.getByText("The user was deleted successfully"),
    ).toBeVisible({ timeout: 10000 })

    await expect(
      page.getByRole("row").filter({ hasText: email }),
    ).not.toBeVisible()
  })

  test("Cancel user creation", async ({ page }) => {
    await page.goto("/admin")

    await page.getByRole("button", { name: "Add User" }).click()
    await page.getByPlaceholder("Email").fill("test@example.com")

    await page.getByRole("button", { name: "Cancel" }).click()

    await expect(page.getByRole("dialog")).not.toBeVisible()
  })

  test("Email is required and must be valid", async ({ page }) => {
    await page.goto("/admin")

    await page.getByRole("button", { name: "Add User" }).click()

    await page.getByPlaceholder("Email").fill("invalid-email")
    await page.getByPlaceholder("Email").blur()

    await expect(page.getByText("Invalid email address")).toBeVisible()
  })

  test("Password must be at least 8 characters", async ({ page }) => {
    await page.goto("/admin")

    await page.getByRole("button", { name: "Add User" }).click()

    await page.getByPlaceholder("Email").fill(randomEmail())
    await page.getByPlaceholder("Password").first().fill("short")
    await page.getByPlaceholder("Password").last().fill("short")
    await page.getByRole("button", { name: "Save" }).click()

    await expect(
      page.getByText("Password must be at least 8 characters"),
    ).toBeVisible()
  })

  test("Passwords must match", async ({ page }) => {
    await page.goto("/admin")

    await page.getByRole("button", { name: "Add User" }).click()

    await page.getByPlaceholder("Email").fill(randomEmail())
    await page.getByPlaceholder("Password").first().fill(randomPassword())
    await page.getByPlaceholder("Password").last().fill("different12345")
    await page.getByPlaceholder("Password").last().blur()

    await expect(page.getByText("The passwords don't match")).toBeVisible()
  })
})

test.describe("Admin page access control", () => {
  test.use({ storageState: { cookies: [], origins: [] } })

  test("Non-superuser cannot access admin page", async ({ page }) => {
    const email = randomEmail()
    const password = randomPassword()

    await createUser({ email, password })
    await logInUser(page, email, password)

    await page.goto("/admin")

    await expect(page.getByRole("heading", { name: "Users" })).not.toBeVisible()
    await expect(page).not.toHaveURL(/\/admin/)
  })

  test("Superuser can access admin page", async ({ page }) => {
    await logInUser(page, firstSuperuser, firstSuperuserPassword)

    await page.goto("/admin")

    await expect(page.getByRole("heading", { name: "Users" })).toBeVisible()
  })
})

test.describe("Admin content management", () => {
  test("CRUD completo de conteudo e topicos com filtros", async ({ page }) => {
    const contents: Array<{
      id: string
      plan: string
      language: string
      day_id: number
      title: string
      content: string
      created_at: string
    }> = []

    const topicsByContent: Record<
      string,
      Array<{
        id: string
        content_id: string
        title: string
        body: string
        behavior: "atividade" | "explicacao"
        activity_type?: "MatchWords" | "FillBlanks" | "Quizz" | "PromptTreiner" | null
        payload: Record<string, unknown>
        order: number
        created_at: string
      }>
    > = {}

    await page.route("**/api/v1/content/**", async (route) => {
      const request = route.request()
      const method = request.method()
      const url = new URL(request.url())
      const pathname = url.pathname

      const jsonResponse = (body: unknown, status = 200) =>
        route.fulfill({
          status,
          contentType: "application/json",
          body: JSON.stringify(body),
        })

      if (method === "GET" && pathname.endsWith("/api/v1/content/items")) {
        const plan = url.searchParams.get("plan")
        const language = url.searchParams.get("language")
        const search = url.searchParams.get("search")?.toLowerCase().trim()

        const filtered = contents.filter((item) => {
          const matchPlan = !plan || item.plan === plan
          const matchLanguage = !language || item.language === language
          const matchSearch =
            !search ||
            item.title.toLowerCase().includes(search) ||
            item.content.toLowerCase().includes(search)
          return matchPlan && matchLanguage && matchSearch
        })

        await jsonResponse({ data: filtered, count: filtered.length })
        return
      }

      if (method === "POST" && pathname.endsWith("/api/v1/content/items")) {
        const payload = (request.postDataJSON() as {
          plan: string
          language: string
          day_id: number
          title: string
          content: string
        }) ?? {
          plan: "base",
          language: "pt-br",
          day_id: 1,
          title: "",
          content: "",
        }

        const created = {
          id: `content-${Date.now()}-${Math.random().toString(36).slice(2, 8)}`,
          ...payload,
          created_at: new Date().toISOString(),
        }
        contents.push(created)
        topicsByContent[created.id] = []
        await jsonResponse(created)
        return
      }

      if (method === "PUT" && pathname.includes("/api/v1/content/items/")) {
        const contentId = pathname.split("/").pop() || ""
        const payload = (request.postDataJSON() as Partial<(typeof contents)[number]>) ?? {}
        const index = contents.findIndex((item) => item.id === contentId)
        if (index === -1) {
          await jsonResponse({ detail: "Content not found" }, 404)
          return
        }
        contents[index] = { ...contents[index], ...payload }
        await jsonResponse(contents[index])
        return
      }

      if (method === "DELETE" && pathname.includes("/api/v1/content/items/")) {
        const contentId = pathname.split("/").pop() || ""
        const index = contents.findIndex((item) => item.id === contentId)
        if (index >= 0) {
          contents.splice(index, 1)
        }
        delete topicsByContent[contentId]
        await jsonResponse({ message: "Content deleted successfully" })
        return
      }

      if (method === "GET" && pathname.includes("/api/v1/content/items/") && pathname.endsWith("/topics")) {
        const parts = pathname.split("/")
        const contentId = parts[parts.length - 2]
        const behavior = url.searchParams.get("behavior")
        const search = url.searchParams.get("search")?.toLowerCase().trim()
        const topics = topicsByContent[contentId] ?? []
        const filtered = topics.filter((item) => {
          const matchBehavior = !behavior || item.behavior === behavior
          const matchSearch =
            !search ||
            item.title.toLowerCase().includes(search) ||
            item.body.toLowerCase().includes(search)
          return matchBehavior && matchSearch
        })
        await jsonResponse({ data: filtered, count: filtered.length })
        return
      }

      if (method === "POST" && pathname.includes("/api/v1/content/items/") && pathname.endsWith("/topics")) {
        const parts = pathname.split("/")
        const contentId = parts[parts.length - 2]
        const payload = (request.postDataJSON() as {
          content_id: string
          behavior: "atividade" | "explicacao"
          activity_type?: "MatchWords" | "FillBlanks" | "Quizz" | "PromptTreiner"
          payload: Record<string, unknown>
          order: number
        }) ?? {
          content_id: contentId,
          behavior: "explicacao",
          payload: { title: "", content: "" },
          order: 1,
        }
        const payloadTitle = payload.payload?.title
        const payloadBody =
          payload.behavior === "explicacao"
            ? payload.payload?.content
            : payload.payload?.description ?? payload.payload?.sentence ?? payload.payload?.question
        const created = {
          id: `topic-${Date.now()}-${Math.random().toString(36).slice(2, 8)}`,
          content_id: contentId,
          title: typeof payloadTitle === "string" ? payloadTitle : "",
          body: typeof payloadBody === "string" ? payloadBody : "",
          behavior: payload.behavior,
          activity_type: payload.activity_type ?? null,
          payload: payload.payload ?? {},
          order: payload.order,
          created_at: new Date().toISOString(),
        }
        if (!topicsByContent[contentId]) {
          topicsByContent[contentId] = []
        }
        topicsByContent[contentId].push(created)
        await jsonResponse(created)
        return
      }

      if (method === "PUT" && pathname.includes("/api/v1/content/topics/")) {
        const topicId = pathname.split("/").pop() || ""
        const payload = (request.postDataJSON() as {
          behavior?: "atividade" | "explicacao"
          activity_type?: "MatchWords" | "FillBlanks" | "Quizz" | "PromptTreiner"
          payload?: Record<string, unknown>
          order?: number
        }) ?? {}

        for (const [contentId, topics] of Object.entries(topicsByContent)) {
          const index = topics.findIndex((topic) => topic.id === topicId)
          if (index !== -1) {
            const nextPayload = payload.payload ?? topics[index].payload
            const payloadTitle = nextPayload?.title
            const payloadBody =
              (payload.behavior ?? topics[index].behavior) === "explicacao"
                ? nextPayload?.content
                : nextPayload?.description ?? nextPayload?.sentence ?? nextPayload?.question
            topics[index] = {
              ...topics[index],
              ...payload,
              content_id: contentId,
              payload: nextPayload,
              title: typeof payloadTitle === "string" ? payloadTitle : topics[index].title,
              body: typeof payloadBody === "string" ? payloadBody : topics[index].body,
            }
            await jsonResponse(topics[index])
            return
          }
        }

        await jsonResponse({ detail: "Topic not found" }, 404)
        return
      }

      if (method === "DELETE" && pathname.includes("/api/v1/content/topics/")) {
        const topicId = pathname.split("/").pop() || ""
        for (const topics of Object.values(topicsByContent)) {
          const index = topics.findIndex((topic) => topic.id === topicId)
          if (index !== -1) {
            topics.splice(index, 1)
            break
          }
        }
        await jsonResponse({ message: "Topic deleted successfully" })
        return
      }

      await route.continue()
    })

    const id = Math.floor(Math.random() * 100000)
    const dayId = 500000 + id
    const contentTitle = `Conteudo Playwright ${id}`
    const updatedContentTitle = `${contentTitle} Atualizado`
    const topicTitle = `Topico Playwright ${id}`
    const updatedTopicTitle = `${topicTitle} Editado`

    await page.goto("/content")

    await expect(
      page.getByRole("heading", { name: "Conteudo" }),
    ).toBeVisible()

    await page.getByRole("button", { name: "Novo Conteudo" }).click()

    const contentDialog = page.getByRole("dialog")
    await contentDialog.getByRole("combobox").first().click()
    await page.getByRole("option", { name: "base" }).click()
    await contentDialog.getByPlaceholder("Idioma").fill("pt-br")
    await contentDialog.getByPlaceholder("Dia").fill(String(dayId))
    await contentDialog.getByPlaceholder("Titulo").fill(contentTitle)
    await contentDialog
      .getByPlaceholder("Conteudo")
      .fill("Conteudo criado via Playwright")
    await contentDialog.getByRole("button", { name: "Salvar" }).click()

    await expect(contentDialog).not.toBeVisible({ timeout: 15000 })
    await expect(page.getByText(contentTitle)).toBeVisible()
    await expect(
      page.getByRole("tab", { name: new RegExp(`Dia ${dayId}`) }),
    ).toBeVisible()

    await page.getByRole("button", { name: "Editar" }).first().click()
    const editContentDialog = page.getByRole("dialog")
    await editContentDialog.getByPlaceholder("Titulo").fill(updatedContentTitle)
    await editContentDialog.getByRole("button", { name: "Salvar" }).click()

    await expect(editContentDialog).not.toBeVisible({ timeout: 15000 })
    await expect(page.getByText(updatedContentTitle)).toBeVisible()

    await page.getByRole("button", { name: "Novo Topico" }).click()
    const topicDialog = page.getByRole("dialog")
    await topicDialog.getByPlaceholder("Titulo da explicacao").fill(topicTitle)
    await topicDialog.getByPlaceholder("Ordem").fill("1")
    await topicDialog
      .getByPlaceholder("Conteudo da explicacao")
      .fill("Corpo do topico de teste")
    await topicDialog.getByRole("button", { name: "Salvar" }).click()

    await expect(topicDialog).not.toBeVisible({ timeout: 15000 })
    await expect(page.getByRole("cell", { name: topicTitle })).toBeVisible()

    const topicRow = page.getByRole("row").filter({ hasText: topicTitle })
    await topicRow.getByRole("button", { name: "Editar" }).click()

    const editTopicDialog = page.getByRole("dialog")
    await editTopicDialog.getByPlaceholder("Titulo da explicacao").fill(updatedTopicTitle)
    await editTopicDialog.getByRole("button", { name: "Salvar" }).click()

    await expect(editTopicDialog).not.toBeVisible({ timeout: 15000 })
    await expect(page.getByRole("cell", { name: updatedTopicTitle })).toBeVisible()

    await page.getByPlaceholder("Pesquisar conteudo...").fill(updatedContentTitle)
    await expect(page.getByText(updatedContentTitle)).toBeVisible()

    await page.getByPlaceholder("Pesquisar topico...").fill(updatedTopicTitle)
    await expect(page.getByRole("cell", { name: updatedTopicTitle })).toBeVisible()

    const updatedTopicRow = page
      .getByRole("row")
      .filter({ hasText: updatedTopicTitle })
    await updatedTopicRow.getByRole("button", { name: "Remover" }).click()

    const deleteTopicDialog = page.getByRole("alertdialog")
    await deleteTopicDialog.getByRole("button", { name: "Remover" }).click()

    await expect(deleteTopicDialog).not.toBeVisible({ timeout: 15000 })
    await expect(page.getByRole("cell", { name: updatedTopicTitle })).not.toBeVisible()

    await page.getByRole("button", { name: "Remover" }).first().click()
    const deleteContentDialog = page.getByRole("alertdialog")
    await deleteContentDialog.getByRole("button", { name: "Remover" }).click()

    await expect(deleteContentDialog).not.toBeVisible({ timeout: 15000 })
    await expect(page.getByText(updatedContentTitle)).not.toBeVisible()
  })
})
