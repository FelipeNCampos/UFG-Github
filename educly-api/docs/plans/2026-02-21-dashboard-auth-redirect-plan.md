# Dashboard Redirect Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Ensure that accessing `/` without an active session immediately clears the stored token and redirects to `/login` whenever `/api/v1/users/me` returns 404.

**Architecture:** Detect the 404 inside `useAuth` by watching the query error and use an effect to clear `access_token`, drop the cached `currentUser`, and navigate to `/login`, with the change validated by a Playwright spec that covers dashboard access without a valid token.

**Tech Stack:** React + TanStack Query/Router, Playwright (`bunx playwright test`), TypeScript.

---

### Task 1: Add regression test for dashboard without a session

**Files:**
- Modify: `frontend/tests/login.spec.ts`

**Step 1: Write the failing test**

Add a Playwright test such as:

```ts
test("Dashboard redirects to /login when no token", async ({ page }) => {
  await page.evaluate(() => localStorage.removeItem("access_token"))
  await page.goto("/")
  await page.waitForURL("/login")
  await expect(page).toHaveURL("/login")
})
```

**Step 2: Run test to confirm it fails**

Run `npm run test -- tests/login.spec.ts --grep "Dashboard redirects"` (or just the entire spec file). Expect it to fail because `/` currently stays on the dashboard.

**Step 3: Re-run the test suite later after implementing the fix**

After Task 2 is done rerun `npm run test -- tests/login.spec.ts` to see the new test pass along with the rest of the file.

**Step 4: Commit the test change**

```bash
git add frontend/tests/login.spec.ts
git commit -m "test: cover dashboard redirect when no token"
```

### Task 2: Make `useAuth` react to 404s from `/users/me`

**Files:**
- Modify: `frontend/src/hooks/useAuth.ts`

**Step 1: Implement the redirect effect**

Within `useAuth`, inspect the query state (e.g., `const { data: user, error, isError } = useQuery(...)`). Add an effect such as:

```ts
useEffect(() => {
  if (isError && error?.status === 404) {
    localStorage.removeItem("access_token")
    queryClient.removeQueries({ queryKey: ["currentUser"] })
    navigate({ to: "/login" })
  }
}, [isError, error?.status, navigate, queryClient])
```

Ensure `users/me` is the only endpoint triggering this effect by relying on the specific query key.

**Step 2: Run the Playwright suite to validate**

Run `npm run test -- tests/login.spec.ts` and confirm the new dashboard redirect test and the existing `invalid token` test both pass.

**Step 3: Commit the implementation**

```bash
git add frontend/src/hooks/useAuth.ts
git commit -m "feat: redirect dashboard when current-user 404s"
```

### Task 3: Manual sanity check

**Files:**
- (none)

**Step 1: Smoke test manually**

Start the frontend (`npm run dev`) and open `/` with and without a token stored in `localStorage`. Confirm it always redirects to `/login` when the backend returns 404 and the logout/menu options stay hidden until you log in.

**Step 2: Document the manual validation in the test run notes before closing the task.**

**Step 3: (Optional) Commit any adjustments discovered during manual QA**

If you needed to tweak styling or hooks because of timing, commit those changes with a descriptive message.


