import uuid
from typing import Any

from fastapi import APIRouter, HTTPException
from sqlmodel import func, select

from app import crud
from app.api.deps import CurrentUser, SessionDep
from app.models import (
    Message,
    Perfil,
    PerfilBitmapUpdate,
    PerfilCreate,
    PerfilPublic,
    PerfilsPublic,
    PerfilUpdate,
    User,
)

router = APIRouter(prefix="/perfils", tags=["perfils"])


@router.get("/", response_model=PerfilsPublic)
def read_perfils(
    session: SessionDep, current_user: CurrentUser, skip: int = 0, limit: int = 100
) -> Any:
    """
    Retrieve perfils. Only superusers can access.
    """
    if not current_user.is_superuser:
        raise HTTPException(status_code=403, detail="Not enough permissions")

    count_statement = select(func.count()).select_from(Perfil)
    count = session.exec(count_statement).one()
    perfils = crud.get_perfils(session=session, skip=skip, limit=limit)
    return PerfilsPublic(data=perfils, count=count)


@router.get("/{id}", response_model=PerfilPublic)
def read_perfil(session: SessionDep, current_user: CurrentUser, id: uuid.UUID) -> Any:
    """
    Get perfil by ID.
    """
    if not current_user.is_superuser:
        raise HTTPException(status_code=403, detail="Not enough permissions")

    perfil = crud.get_perfil(session=session, perfil_id=id)
    if not perfil:
        raise HTTPException(status_code=404, detail="Perfil not found")
    return perfil


@router.post("/", response_model=PerfilPublic)
def create_perfil(
    *, session: SessionDep, current_user: CurrentUser, perfil_in: PerfilCreate
) -> Any:
    """
    Create new perfil. Only superusers can create perfils.
    """
    if not current_user.is_superuser:
        raise HTTPException(status_code=403, detail="Not enough permissions")

    user = session.get(User, perfil_in.user_id)
    if not user:
        raise HTTPException(status_code=404, detail="User not found")

    existing_perfil = crud.get_perfil_by_user_id(session=session, user_id=perfil_in.user_id)
    if existing_perfil:
        raise HTTPException(status_code=400, detail="Perfil for this user already exists")

    return crud.create_perfil(session=session, perfil_in=perfil_in)


@router.patch("/{id}", response_model=PerfilPublic)
def update_perfil(
    *,
    session: SessionDep,
    current_user: CurrentUser,
    id: uuid.UUID,
    perfil_in: PerfilUpdate,
) -> Any:
    """
    Update a perfil. Only superusers can update perfils.
    """
    if not current_user.is_superuser:
        raise HTTPException(status_code=403, detail="Not enough permissions")

    db_perfil = crud.get_perfil(session=session, perfil_id=id)
    if not db_perfil:
        raise HTTPException(status_code=404, detail="Perfil not found")

    return crud.update_perfil(session=session, db_perfil=db_perfil, perfil_in=perfil_in)


@router.delete("/{id}", response_model=Message)
def delete_perfil(
    session: SessionDep, current_user: CurrentUser, id: uuid.UUID
) -> Message:
    """
    Delete a perfil. Only superusers can delete perfils.
    """
    if not current_user.is_superuser:
        raise HTTPException(status_code=403, detail="Not enough permissions")

    db_perfil = crud.get_perfil(session=session, perfil_id=id)
    if not db_perfil:
        raise HTTPException(status_code=404, detail="Perfil not found")

    crud.delete_perfil(session=session, db_perfil=db_perfil)
    return Message(message="Perfil deleted successfully")


@router.patch("/users/{user_id}/bitmap", response_model=PerfilPublic)
def update_user_bitmap(
    *,
    session: SessionDep,
    current_user: CurrentUser,
    user_id: uuid.UUID,
    body: PerfilBitmapUpdate,
) -> Any:
    """
    Update user access bitmap (profile_bitmap). Creates perfil if it does not exist.
    """
    if not current_user.is_superuser:
        raise HTTPException(status_code=403, detail="Not enough permissions")

    user = session.get(User, user_id)
    if not user:
        raise HTTPException(status_code=404, detail="User not found")

    return crud.set_user_profile_bitmap(
        session=session, user_id=user_id, profile_bitmap=body.profile_bitmap
    )
