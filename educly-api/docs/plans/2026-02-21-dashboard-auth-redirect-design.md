# Dashboard Redirect on Auth Failure

## Context
When the SPA is opened without a valid session the dashboard route still renders because `useAuth` has no guard: `currentUser` query errors with 404 but the UI keeps showing the dashboard header, leaving the logout button and admin controls inaccessible. We need to force a redirect back to `/login` and clear stored tokens whenever `/api/v1/users/me` responds with 404.

## Success criteria
- The dashboard immediately navigates to `/login` whenever the current-user fetch returns 404.
- Stored tokens (`localStorage` `access_token`) are removed before redirect.
- The behavior applies everywhere that uses the shared `useAuth` hook so no other route shows stale content.
- The fix is resilient to retries or intermittent failures (only act on HTTP 404; let other errors keep existing handling).

## Proposed implementation (recommended)
1. Keep the `currentUser` query inside `useAuth`. Observe its `error` state and check for `error?.status === 404`. When the check passes:
   - Remove `access_token` from `localStorage`.
   - Clear the cached `currentUser` query so stale errors don’t block retry.
   - Redirect via `navigate({ to: '/login' })`.
2. Wrap that logic in `useEffect` so it only runs once per 404 event and avoids triggering before React Router is ready.
3. Confirm there are no regressions by manually testing the app with and without a token.

## Risks & mitigations
- Clearing `localStorage` may also log out other tabs; this matches the expectation that a 404 means the session is invalid, so it is acceptable. If needed later we can target only the current tab via storage events.
- Redirecting on 404 must not run for other routes; keeping this inside `useAuth` scoped to `UsersService.readUserMe` ensures we only respond to that endpoint.
