import type { APIRequestContext } from "@playwright/test"

type Email = {
  id: number
  recipients: string[]
  subject: string
}

export const MAILCATCHER_HOST =
  process.env.MAILCATCHER_HOST?.replace("localhost", "127.0.0.1") ??
  "http://127.0.0.1:1080"

function normalizeRecipient(recipient: unknown): string[] {
  if (typeof recipient === "string") {
    const trimmed = recipient.trim()
    const unwrapped = trimmed.replace(/^<(.+)>$/, "$1")
    return Array.from(new Set([trimmed, unwrapped]))
  }

  if (
    recipient &&
    typeof recipient === "object" &&
    "mailbox" in recipient &&
    "domain" in recipient
  ) {
    const mailbox = (recipient as { mailbox?: unknown }).mailbox
    const domain = (recipient as { domain?: unknown }).domain
    if (typeof mailbox === "string" && typeof domain === "string") {
      const email = `${mailbox}@${domain}`
      return [email, `<${email}>`]
    }
  }

  return []
}

function normalizeEmail(raw: unknown): Email | null {
  if (!raw || typeof raw !== "object") {
    return null
  }

  const candidate = raw as {
    id?: unknown
    subject?: unknown
    recipients?: unknown
  }

  if (typeof candidate.id !== "number") {
    return null
  }

  const recipientsRaw = Array.isArray(candidate.recipients)
    ? candidate.recipients
    : []
  const recipients = Array.from(
    new Set(recipientsRaw.flatMap((recipient) => normalizeRecipient(recipient))),
  )

  return {
    id: candidate.id,
    subject: typeof candidate.subject === "string" ? candidate.subject : "",
    recipients,
  }
}

async function assertMailcatcherUp(request: APIRequestContext): Promise<void> {
  let response: Awaited<ReturnType<APIRequestContext["get"]>>
  try {
    response = await request.get(`${MAILCATCHER_HOST}/messages`, {
      timeout: 3000,
    })
  } catch (error) {
    const err = new Error(
      `Mailcatcher unavailable at ${MAILCATCHER_HOST}. Start Mailcatcher before running this test.`,
    )
    ;(err as Error & { cause?: unknown }).cause = error
    throw err
  }

  if (!response.ok()) {
    throw new Error(
      `Mailcatcher unavailable at ${MAILCATCHER_HOST} (status ${response.status()}).`,
    )
  }
}

async function findEmail({
  request,
  filter,
}: {
  request: APIRequestContext
  filter?: (email: Email) => boolean
}) {
  const response = await request.get(`${MAILCATCHER_HOST}/messages`)
  if (!response.ok()) {
    throw new Error(
      `Mailcatcher request failed at ${MAILCATCHER_HOST}/messages (status ${response.status()}).`,
    )
  }

  const payload = await response.json()
  const emails = Array.isArray(payload)
    ? payload
        .map((entry) => normalizeEmail(entry))
        .filter((entry): entry is Email => entry !== null)
    : []

  if (filter) {
    const filteredEmails = emails.filter(filter)
    return filteredEmails[filteredEmails.length - 1] ?? null
  }

  return emails[emails.length - 1] ?? null
}

export async function findLastEmail({
  request,
  filter,
  timeout = 5000,
}: {
  request: APIRequestContext
  filter?: (email: Email) => boolean
  timeout?: number
}) {
  await assertMailcatcherUp(request)
  const deadline = Date.now() + timeout

  while (Date.now() < deadline) {
    const emailData = await findEmail({ request, filter })
    if (emailData) {
      return emailData
    }

    // Wait briefly before polling Mailcatcher again.
    await new Promise((resolve) => setTimeout(resolve, 100))
  }

  const smtpHost = process.env.SMTP_HOST ?? "undefined"
  const smtpPort = process.env.SMTP_PORT ?? "undefined"
  const smtpUser = process.env.SMTP_USER ?? "undefined"
  throw new Error(
    `Timeout while trying to get latest email from ${MAILCATCHER_HOST}. SMTP_HOST=${smtpHost}, SMTP_PORT=${smtpPort}, SMTP_USER=${smtpUser}. For Mailcatcher-based tests, use SMTP_HOST=127.0.0.1 and SMTP_PORT=1025 (or SMTP_HOST=mailcatcher in Docker), and leave SMTP_USER/SMTP_PASSWORD unset.`,
  )
}
