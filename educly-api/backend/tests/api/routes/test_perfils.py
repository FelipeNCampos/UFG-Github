import uuid

import pytest
from fastapi.testclient import TestClient
from sqlmodel import Session, delete, select

from app.core.config import settings
from app.models import Perfil
from tests.utils.user import create_random_user


def create_perfil_payload(*, user_id: uuid.UUID, profile_bitmap: int = 1) -> dict:
    return {
        "user_id": str(user_id),
        "profile_bitmap": profile_bitmap,
        "personalized_trail": False,
        "avatar_url": "https://example.com/avatar.png",
        "cover_url": "https://example.com/cover.png",
        "bio": "Perfil de teste",
    }


def create_user_without_perfil(db: Session):
    user = create_random_user(db)
    db.exec(delete(Perfil).where(Perfil.user_id == user.id))
    db.commit()
    return user

def test_create_perfil_superuser(
    client: TestClient, superuser_token_headers: dict[str, str], db: Session
) -> None:
    user = create_user_without_perfil(db)
    data = create_perfil_payload(user_id=user.id, profile_bitmap=2)

    r = client.post(
        f"{settings.API_V1_STR}/perfils/",
        headers=superuser_token_headers,
        json=data,
    )

    assert r.status_code == 200
    created_perfil = r.json()
    assert created_perfil["user_id"] == str(user.id)
    assert created_perfil["profile_bitmap"] == 2
    assert created_perfil["bio"] == "Perfil de teste"

    db_perfil = db.exec(select(Perfil).where(Perfil.user_id == user.id)).first()
    assert db_perfil
    assert db_perfil.profile_bitmap == 2


def test_create_perfil_user_not_found(
    client: TestClient, superuser_token_headers: dict[str, str]
) -> None:
    data = create_perfil_payload(user_id=uuid.uuid4(), profile_bitmap=4)

    r = client.post(
        f"{settings.API_V1_STR}/perfils/",
        headers=superuser_token_headers,
        json=data,
    )

    assert r.status_code == 404
    assert r.json()["detail"] == "User not found"


def test_create_perfil_duplicate_user(
    client: TestClient, superuser_token_headers: dict[str, str], db: Session
) -> None:
    user = create_user_without_perfil(db)
    data = create_perfil_payload(user_id=user.id, profile_bitmap=1)

    r = client.post(
        f"{settings.API_V1_STR}/perfils/",
        headers=superuser_token_headers,
        json=data,
    )
    assert r.status_code == 200

    r_dup = client.post(
        f"{settings.API_V1_STR}/perfils/",
        headers=superuser_token_headers,
        json=data,
    )
    assert r_dup.status_code == 400
    assert r_dup.json()["detail"] == "Perfil for this user already exists"


def test_read_perfils_superuser(
    client: TestClient, superuser_token_headers: dict[str, str], db: Session
) -> None:
    user = create_user_without_perfil(db)
    data = create_perfil_payload(user_id=user.id, profile_bitmap=8)
    client.post(
        f"{settings.API_V1_STR}/perfils/",
        headers=superuser_token_headers,
        json=data,
    )

    r = client.get(f"{settings.API_V1_STR}/perfils/", headers=superuser_token_headers)
    response = r.json()

    assert r.status_code == 200
    assert "count" in response
    assert isinstance(response["data"], list)
    assert response["count"] >= len(response["data"])


def test_read_perfils_normal_user_forbidden(
    client: TestClient, normal_user_token_headers: dict[str, str]
) -> None:
    r = client.get(f"{settings.API_V1_STR}/perfils/", headers=normal_user_token_headers)
    assert r.status_code == 403
    assert r.json()["detail"] == "Not enough permissions"


def test_read_perfil_by_id_superuser(
    client: TestClient, superuser_token_headers: dict[str, str], db: Session
) -> None:
    user = create_user_without_perfil(db)
    data = create_perfil_payload(user_id=user.id, profile_bitmap=4)
    create_response = client.post(
        f"{settings.API_V1_STR}/perfils/",
        headers=superuser_token_headers,
        json=data,
    )
    perfil_id = create_response.json()["id"]

    r = client.get(
        f"{settings.API_V1_STR}/perfils/{perfil_id}",
        headers=superuser_token_headers,
    )

    assert r.status_code == 200
    perfil = r.json()
    assert perfil["id"] == perfil_id
    assert perfil["user_id"] == str(user.id)


def test_read_perfil_by_id_not_found(
    client: TestClient, superuser_token_headers: dict[str, str]
) -> None:
    r = client.get(
        f"{settings.API_V1_STR}/perfils/{uuid.uuid4()}",
        headers=superuser_token_headers,
    )
    assert r.status_code == 404
    assert r.json()["detail"] == "Perfil not found"


def test_read_perfil_by_id_normal_user_forbidden(
    client: TestClient, normal_user_token_headers: dict[str, str]
) -> None:
    r = client.get(
        f"{settings.API_V1_STR}/perfils/{uuid.uuid4()}",
        headers=normal_user_token_headers,
    )
    assert r.status_code == 403
    assert r.json()["detail"] == "Not enough permissions"


def test_update_perfil_superuser(
    client: TestClient, superuser_token_headers: dict[str, str], db: Session
) -> None:
    user = create_user_without_perfil(db)
    data = create_perfil_payload(user_id=user.id, profile_bitmap=1)
    create_response = client.post(
        f"{settings.API_V1_STR}/perfils/",
        headers=superuser_token_headers,
        json=data,
    )
    perfil_id = create_response.json()["id"]

    update_data = {"profile_bitmap": 8, "personalized_trail": True, "bio": "Atualizado"}
    r = client.patch(
        f"{settings.API_V1_STR}/perfils/{perfil_id}",
        headers=superuser_token_headers,
        json=update_data,
    )

    assert r.status_code == 200
    updated = r.json()
    assert updated["profile_bitmap"] == 8
    assert updated["personalized_trail"] is True
    assert updated["bio"] == "Atualizado"

    db_perfil = db.get(Perfil, uuid.UUID(perfil_id))
    assert db_perfil
    assert db_perfil.profile_bitmap == 8


def test_update_perfil_not_found(
    client: TestClient, superuser_token_headers: dict[str, str]
) -> None:
    r = client.patch(
        f"{settings.API_V1_STR}/perfils/{uuid.uuid4()}",
        headers=superuser_token_headers,
        json={"profile_bitmap": 2},
    )
    assert r.status_code == 404
    assert r.json()["detail"] == "Perfil not found"


def test_update_perfil_normal_user_forbidden(
    client: TestClient, normal_user_token_headers: dict[str, str]
) -> None:
    r = client.patch(
        f"{settings.API_V1_STR}/perfils/{uuid.uuid4()}",
        headers=normal_user_token_headers,
        json={"profile_bitmap": 4},
    )
    assert r.status_code == 403
    assert r.json()["detail"] == "Not enough permissions"


def test_update_user_bitmap_superuser_creates_perfil_when_missing(
    client: TestClient, superuser_token_headers: dict[str, str], db: Session
) -> None:
    user = create_user_without_perfil(db)

    r = client.patch(
        f"{settings.API_V1_STR}/perfils/users/{user.id}/bitmap",
        headers=superuser_token_headers,
        json={"profile_bitmap": 12},
    )

    assert r.status_code == 200
    payload = r.json()
    assert payload["user_id"] == str(user.id)
    assert payload["profile_bitmap"] == 12

    db_perfil = db.exec(select(Perfil).where(Perfil.user_id == user.id)).first()
    assert db_perfil
    assert db_perfil.profile_bitmap == 12


def test_update_user_bitmap_superuser_updates_existing_perfil(
    client: TestClient, superuser_token_headers: dict[str, str], db: Session
) -> None:
    user = create_random_user(db)

    r = client.patch(
        f"{settings.API_V1_STR}/perfils/users/{user.id}/bitmap",
        headers=superuser_token_headers,
        json={"profile_bitmap": 8},
    )
    assert r.status_code == 200
    assert r.json()["profile_bitmap"] == 8

    db_perfil = db.exec(select(Perfil).where(Perfil.user_id == user.id)).first()
    assert db_perfil
    assert db_perfil.profile_bitmap == 8


def test_update_user_bitmap_user_not_found(
    client: TestClient, superuser_token_headers: dict[str, str]
) -> None:
    r = client.patch(
        f"{settings.API_V1_STR}/perfils/users/{uuid.uuid4()}/bitmap",
        headers=superuser_token_headers,
        json={"profile_bitmap": 4},
    )
    assert r.status_code == 404
    assert r.json()["detail"] == "User not found"


def test_update_user_bitmap_normal_user_forbidden(
    client: TestClient, normal_user_token_headers: dict[str, str], db: Session
) -> None:
    user = create_random_user(db)
    r = client.patch(
        f"{settings.API_V1_STR}/perfils/users/{user.id}/bitmap",
        headers=normal_user_token_headers,
        json={"profile_bitmap": 4},
    )
    assert r.status_code == 403
    assert r.json()["detail"] == "Not enough permissions"


def test_delete_perfil_superuser(
    client: TestClient, superuser_token_headers: dict[str, str], db: Session
) -> None:
    user = create_random_user(db)
    perfil = db.exec(select(Perfil).where(Perfil.user_id == user.id)).first()
    assert perfil
    perfil_id = str(perfil.id)

    r = client.delete(
        f"{settings.API_V1_STR}/perfils/{perfil_id}",
        headers=superuser_token_headers,
    )

    assert r.status_code == 200
    assert r.json()["message"] == "Perfil deleted successfully"
    db.expire_all()
    db_perfil = db.get(Perfil, uuid.UUID(perfil_id))
    assert db_perfil is None


def test_delete_perfil_not_found(
    client: TestClient, superuser_token_headers: dict[str, str]
) -> None:
    r = client.delete(
        f"{settings.API_V1_STR}/perfils/{uuid.uuid4()}",
        headers=superuser_token_headers,
    )
    assert r.status_code == 404
    assert r.json()["detail"] == "Perfil not found"


def test_delete_perfil_normal_user_forbidden(
    client: TestClient, normal_user_token_headers: dict[str, str]
) -> None:
    r = client.delete(
        f"{settings.API_V1_STR}/perfils/{uuid.uuid4()}",
        headers=normal_user_token_headers,
    )
    assert r.status_code == 403
    assert r.json()["detail"] == "Not enough permissions"
