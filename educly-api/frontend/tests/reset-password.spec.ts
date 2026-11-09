import { expect, test, type APIRequestContext } from "@playwright/test"
import { firstSuperuser, firstSuperuserPassword } from "./config"
import { findLastEmail, MAILCATCHER_HOST } from "./utils/mailcatcher"
import { createUser } from "./utils/privateApi"
import { randomEmail, randomPassword } from "./utils/random"
import { logInUser } from "./utils/user"

test.use({ storageState: { cookies: [], origins: [] } })

const API_BASE_URL = process.env.VITE_API_URL ?? "http://localhost:8000"

async function getResetLinkFromHtmlContent(
  request: APIRequestContext,
  email: string,
) {
  const tokenResponse = await request.post(
    `${API_BASE_URL}/api/v1/login/access-token`,
    {
    form: {
      username: firstSuperuser,
      password: firstSuperuserPassword,
    },
    },
  )
  if (!tokenResponse.ok()) {
    throw new Error(
      `Failed to authenticate superuser (${tokenResponse.status()})`,
    )
  }
  const tokenJson = (await tokenResponse.json()) as { access_token?: string }
  const accessToken = tokenJson.access_token
  if (!accessToken) {
    throw new Error("Failed to get superuser access token for reset-password test")
  }

  const htmlResponse = await request.post(
    `${API_BASE_URL}/api/v1/password-recovery-html-content/${encodeURIComponent(email)}`,
    {
      headers: { Authorization: `Bearer ${accessToken}` },
    },
  )
  if (htmlResponse.ok()) {
    const htmlContent = await htmlResponse.text()
    const linkMatch = htmlContent.match(
      /href="([^"]*\/reset-password\?token=[^"]+)"/,
    )
    if (linkMatch?.[1]) {
      const resetLink = linkMatch[1]
      if (resetLink.startsWith("http://localhost/")) {
        return resetLink.replace("http://localhost/", "http://localhost:5173/")
      }
      return resetLink
    }
  }

  const emailData = await findLastEmail({
    request,
    filter: (e) => e.recipients.includes(email),
    timeout: 8000,
  })
  const mailcatcherHtmlResponse = await request.get(
    `${MAILCATCHER_HOST}/messages/${emailData.id}.html`,
  )
  if (!mailcatcherHtmlResponse.ok()) {
    throw new Error(
      `Failed to fetch Mailcatcher message html (${mailcatcherHtmlResponse.status()})`,
    )
  }
  const mailcatcherHtml = await mailcatcherHtmlResponse.text()
  const linkMatch = mailcatcherHtml.match(
    /href="([^"]*\/reset-password\?token=[^"]+)"/,
  )

  if (!linkMatch?.[1]) {
    throw new Error("Could not find reset-password link in recovered email HTML")
  }

  const resetLink = linkMatch[1]
  if (resetLink.startsWith("http://localhost/")) {
    return resetLink.replace("http://localhost/", "http://localhost:5173/")
  }
  return resetLink
}

test("Password Recovery title is visible", async ({ page }) => {
  await page.goto("/recover-password")

  await expect(
    page.getByRole("heading", { name: "Password Recovery" }),
  ).toBeVisible()
})

test("Input is visible, empty and editable", async ({ page }) => {
  await page.goto("/recover-password")

  await expect(page.getByTestId("email-input")).toBeVisible()
  await expect(page.getByTestId("email-input")).toHaveText("")
  await expect(page.getByTestId("email-input")).toBeEditable()
})

test("Continue button is visible", async ({ page }) => {
  await page.goto("/recover-password")

  await expect(page.getByRole("button", { name: "Continue" })).toBeVisible()
})

test("User can reset password successfully using the link", async ({
  page,
  request,
}) => {
  const email = randomEmail()
  const password = randomPassword()
  const newPassword = randomPassword()

  await createUser({ email, password })

  await page.goto("/recover-password")
  await page.getByTestId("email-input").fill(email)

  await page.getByRole("button", { name: "Continue" }).click()

  const url = await getResetLinkFromHtmlContent(request, email)

  // Set the new password and confirm it
  await page.goto(url)

  await page.getByTestId("new-password-input").fill(newPassword)
  await page.getByTestId("confirm-password-input").fill(newPassword)
  await page.getByRole("button", { name: "Reset Password" }).click()
  await expect(page.getByText("Password updated successfully")).toBeVisible()

  // Check if the user is able to login with the new password
  await logInUser(page, email, newPassword)
})

test("Expired or invalid reset link", async ({ page }) => {
  const password = randomPassword()
  const invalidUrl = "/reset-password?token=invalidtoken"

  await page.goto(invalidUrl)

  await page.getByTestId("new-password-input").fill(password)
  await page.getByTestId("confirm-password-input").fill(password)
  await page.getByRole("button", { name: "Reset Password" }).click()

  await expect(page.getByText("Invalid token")).toBeVisible()
})

test("Weak new password validation", async ({ page, request }) => {
  const email = randomEmail()
  const password = randomPassword()
  const weakPassword = "123"

  await createUser({ email, password })

  await page.goto("/recover-password")
  await page.getByTestId("email-input").fill(email)
  await page.getByRole("button", { name: "Continue" }).click()

  const url = await getResetLinkFromHtmlContent(request, email)

  // Set a weak new password
  await page.goto(url)
  await page.getByTestId("new-password-input").fill(weakPassword)
  await page.getByTestId("confirm-password-input").fill(weakPassword)
  await page.getByRole("button", { name: "Reset Password" }).click()

  await expect(
    page.getByText("Password must be at least 8 characters"),
  ).toBeVisible()
})
