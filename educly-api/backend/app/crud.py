import uuid
from typing import Any

from sqlmodel import Session, select

from app.core.security import get_password_hash, verify_password
from app.models import Perfil, PerfilCreate, PerfilUpdate, User, UserCreate, UserUpdate


def _normalize_profile_bitmap(bitmap: int | None) -> int:
    normalized = int(bitmap or 1)
    if normalized < 1:
        normalized = 1
    # free (bit 1) cannot coexist with paid plans.
    if (normalized & ~1) != 0:
        normalized = normalized & ~1
    return normalized


def create_user(*, session: Session, user_create: UserCreate) -> User:
    db_obj = User.model_validate(
        user_create, update={"hashed_password": get_password_hash(user_create.password)}
    )
    session.add(db_obj)
    session.commit()
    session.refresh(db_obj)

    perfil = Perfil(
        user_id=db_obj.id,
        profile_bitmap=_normalize_profile_bitmap(user_create.profile_bitmap),
        personalized_trail=user_create.personalized_trail,
        avatar_url=user_create.avatar_url,
        cover_url=user_create.cover_url,
        bio=user_create.bio,
    )
    session.add(perfil)
    session.commit()

    return db_obj


def update_user(*, session: Session, db_user: User, user_in: UserUpdate) -> Any:
    user_data = user_in.model_dump(exclude_unset=True)
    extra_data = {}
    if "password" in user_data:
        password = user_data["password"]
        hashed_password = get_password_hash(password)
        extra_data["hashed_password"] = hashed_password
    db_user.sqlmodel_update(user_data, update=extra_data)
    session.add(db_user)
    session.commit()
    session.refresh(db_user)
    return db_user


def get_user_by_email(*, session: Session, email: str) -> User | None:
    statement = select(User).where(User.email == email)
    session_user = session.exec(statement).first()
    return session_user


# Dummy hash to use for timing attack prevention when user is not found
# This is an Argon2 hash of a random password, used to ensure constant-time comparison
DUMMY_HASH = "$argon2id$v=19$m=65536,t=3,p=4$MjQyZWE1MzBjYjJlZTI0Yw$YTU4NGM5ZTZmYjE2NzZlZjY0ZWY3ZGRkY2U2OWFjNjk"


def authenticate(*, session: Session, email: str, password: str) -> User | None:
    db_user = get_user_by_email(session=session, email=email)
    if not db_user:
        # Prevent timing attacks by running password verification even when user doesn't exist
        # This ensures the response time is similar whether or not the email exists
        verify_password(password, DUMMY_HASH)
        return None
    verified, updated_password_hash = verify_password(password, db_user.hashed_password)
    if not verified:
        return None
    if updated_password_hash:
        db_user.hashed_password = updated_password_hash
        session.add(db_user)
        session.commit()
        session.refresh(db_user)
    return db_user


def get_perfils(*, session: Session, skip: int = 0, limit: int = 100) -> list[Perfil]:
    statement = select(Perfil).offset(skip).limit(limit)
    return list(session.exec(statement).all())


def get_perfil(*, session: Session, perfil_id: uuid.UUID) -> Perfil | None:
    return session.get(Perfil, perfil_id)


def get_perfil_by_user_id(*, session: Session, user_id: uuid.UUID) -> Perfil | None:
    statement = select(Perfil).where(Perfil.user_id == user_id)
    return session.exec(statement).first()


def create_perfil(*, session: Session, perfil_in: PerfilCreate) -> Perfil:
    perfil = Perfil.model_validate(perfil_in)
    perfil.profile_bitmap = _normalize_profile_bitmap(perfil.profile_bitmap)
    session.add(perfil)
    session.commit()
    session.refresh(perfil)
    return perfil


def update_perfil(*, session: Session, db_perfil: Perfil, perfil_in: PerfilUpdate) -> Perfil:
    update_data = perfil_in.model_dump(exclude_unset=True)
    if "profile_bitmap" in update_data and update_data["profile_bitmap"] is not None:
        update_data["profile_bitmap"] = _normalize_profile_bitmap(update_data["profile_bitmap"])
    db_perfil.sqlmodel_update(update_data)
    session.add(db_perfil)
    session.commit()
    session.refresh(db_perfil)
    return db_perfil


def delete_perfil(*, session: Session, db_perfil: Perfil) -> None:
    session.delete(db_perfil)
    session.commit()


def set_user_profile_bitmap(
    *, session: Session, user_id: uuid.UUID, profile_bitmap: int
) -> Perfil:
    normalized_bitmap = _normalize_profile_bitmap(profile_bitmap)
    perfil = get_perfil_by_user_id(session=session, user_id=user_id)
    if perfil:
        perfil.profile_bitmap = normalized_bitmap
        session.add(perfil)
        session.commit()
        session.refresh(perfil)
        return perfil

    perfil_in = PerfilCreate(
        user_id=user_id,
        profile_bitmap=normalized_bitmap,
        personalized_trail=False,
        avatar_url=None,
        cover_url=None,
        bio=None,
    )
    return create_perfil(session=session, perfil_in=perfil_in)


