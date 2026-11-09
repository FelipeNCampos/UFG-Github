# Localhost 3000 CORS Design

## Goal
Ensure `http://localhost:3000` is permitted by default in both code and env configuration so tooling (storybook or frontend dev server) can interact with the API without extra per-machine steps.

## Context
Currently the backend restricts `BACKEND_CORS_ORIGINS` to the values in `.env` and the default frontend host. The goal is to relax this to include port 3000 for local dev scenarios while keeping existing origin handling intact.

## Proposed Approach
1. **Settings-level default**: Extend `Settings.BACKEND_CORS_ORIGINS` default to include `http://localhost:3000`. That keeps automated deployments/tests within this repo aware of the origin without relying on derived env files.
2. **`.env` alignment**: Add `http://localhost:3000` to the `.env`’s `BACKEND_CORS_ORIGINS` string so local developers loading that file see the same origins as the runtime defaults.
3. **No extra logic needed**: The existing `parse_cors` hook already handles comma-separated strings and `all_cors_origins` continues building the list, so no extra code adjustments are required.

Once these values are synchronized, the backend accepts requests from Storybook or other frontends running on 3000 without manual env tweaks.
