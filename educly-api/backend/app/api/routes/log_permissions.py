import uuid
from typing import Any

from fastapi import APIRouter, HTTPException, Query
from sqlmodel import col, func, select

from app.api.deps import CurrentUser, SessionDep
from app.models import (
    LogPermission,
    LogPermissionCreate,
    LogPermissionPublic,
    LogPermissionsPublic,
    LogPermissionUpdate,
)

router = APIRouter(prefix="/log-permissions", tags=["log-permissions"])


@router.get("/", response_model=LogPermissionsPublic)
def read_log_permissions(
    session: SessionDep,
    current_user: CurrentUser,
    skip: int = 0,
    limit: int = 100,
    status: str | None = Query(default=None, max_length=30),
    q: str | None = Query(default=None, max_length=255),
) -> Any:
    if not current_user.is_superuser:
        raise HTTPException(status_code=403, detail="Not enough permissions")

    statement = select(LogPermission)
    count_statement = select(func.count()).select_from(LogPermission)

    if status:
        statement = statement.where(LogPermission.status == status.lower())
        count_statement = count_statement.where(LogPermission.status == status.lower())

    if q:
        like = f"%{q.lower()}%"
        statement = statement.where(
            func.lower(LogPermission.action).like(like)
            | func.lower(LogPermission.email).like(like)
            | func.lower(func.coalesce(LogPermission.reason, "")).like(like)
            | func.lower(func.coalesce(LogPermission.granted_plan, "")).like(like)
        )
        count_statement = count_statement.where(
            func.lower(LogPermission.action).like(like)
            | func.lower(LogPermission.email).like(like)
            | func.lower(func.coalesce(LogPermission.reason, "")).like(like)
            | func.lower(func.coalesce(LogPermission.granted_plan, "")).like(like)
        )

    count = session.exec(count_statement).one()
    rows = session.exec(
        statement.order_by(col(LogPermission.created_at).desc()).offset(skip).limit(limit)
    ).all()
    return LogPermissionsPublic(data=list(rows), count=count)


@router.post("/", response_model=LogPermissionPublic)
def create_log_permission(
    *,
    session: SessionDep,
    current_user: CurrentUser,
    request_in: LogPermissionCreate,
) -> Any:
    if not current_user.is_superuser:
        raise HTTPException(status_code=403, detail="Not enough permissions")

    payload = request_in.model_dump()
    payload["status"] = payload["status"].lower()
    request_obj = LogPermission.model_validate(payload)
    session.add(request_obj)
    session.commit()
    session.refresh(request_obj)
    return request_obj


@router.patch("/{id}", response_model=LogPermissionPublic)
def update_log_permission(
    *,
    session: SessionDep,
    current_user: CurrentUser,
    id: uuid.UUID,
    request_in: LogPermissionUpdate,
) -> Any:
    if not current_user.is_superuser:
        raise HTTPException(status_code=403, detail="Not enough permissions")

    request_obj = session.get(LogPermission, id)
    if not request_obj:
        raise HTTPException(status_code=404, detail="LogPermission not found")

    update_data = request_in.model_dump(exclude_unset=True)
    if "status" in update_data and update_data["status"]:
        update_data["status"] = update_data["status"].lower()
    request_obj.sqlmodel_update(update_data)
    session.add(request_obj)
    session.commit()
    session.refresh(request_obj)
    return request_obj
