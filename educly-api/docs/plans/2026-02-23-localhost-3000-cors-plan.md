# Localhost 3000 CORS Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Allow `http://localhost:3000` to be accepted by the backend CORS stack without manual environment edits.

**Architecture:** Update the base settings defaults so any process that instantiates `Settings` already permits port 3000, and keep the `.env` defaults aligned for developers who source that file.

**Tech Stack:** Python, FastAPI, Pydantic settings, dotenv-based environment configuration.

---

### Task 1: Permit localhost:3000 in code and env

**Files:**
- Modify: `backend/app/core/config.py`
- Modify: `.env`

**Step 1: Update defaults**
Extend the `BACKEND_CORS_ORIGINS` default list in `Settings` to include `http://localhost:3000` and append the same origin to the `BACKEND_CORS_ORIGINS` line in `.env` so both runtime defaults and local env files match.

**Step 2: Run backend tests**
Run `pytest backend -k cors -q` (or `pytest backend -q` if the targeted filter is unavailable) to confirm nothing breaks when the origin list changes.

**Step 3: Smoke-launch the API**
Start the API via `uvicorn backend.app.main:app --reload --port 8000` to ensure runtime loading of the new defaults works.

**Step 4: Manual CORS verification**
While the server runs, execute `curl -I -H "Origin: http://localhost:3000" http://localhost:8000/health` and check `access-control-allow-origin: http://localhost:3000` is present.

**Step 5: Commit**
Stage the two modified files and commit with a message like `feat: allow localhost 3000 cors origin`.
