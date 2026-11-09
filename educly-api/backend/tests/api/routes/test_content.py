from fastapi.testclient import TestClient
from sqlmodel import Session

from app import crud
from app.core.config import settings
from app.models import UserCreate
from tests.utils.user import user_authentication_headers
from tests.utils.utils import random_email, random_lower_string


def test_read_base_content_day_with_base_access(client: TestClient, db: Session) -> None:
    email = random_email()
    password = random_lower_string()
    crud.create_user(
        session=db,
        user_create=UserCreate(
            email=email,
            password=password,
            full_name="Base User",
            profile_bitmap=2,
        ),
    )
    headers = user_authentication_headers(client=client, email=email, password=password)

    r = client.get(f"{settings.API_V1_STR}/content/base/1?language=pt-br", headers=headers)
    assert r.status_code == 200
    data = r.json()
    assert data["id"] == 1
    assert data["language"] == "pt-br"
    assert data["plan"] == "base"
    assert data["title"]
    assert data["content"]


def test_read_base_content_day_without_base_access_forbidden(
    client: TestClient, db: Session
) -> None:
    email = random_email()
    password = random_lower_string()
    crud.create_user(
        session=db,
        user_create=UserCreate(
            email=email,
            password=password,
            full_name="Free User",
            profile_bitmap=1,
        ),
    )
    headers = user_authentication_headers(client=client, email=email, password=password)

    r = client.get(f"{settings.API_V1_STR}/content/base/1?language=pt-br", headers=headers)
    assert r.status_code == 403
    assert r.json()["detail"] == "Base plan access required"


def test_read_base_content_day_not_found(client: TestClient, db: Session) -> None:
    email = random_email()
    password = random_lower_string()
    crud.create_user(
        session=db,
        user_create=UserCreate(
            email=email,
            password=password,
            full_name="Base User",
            profile_bitmap=2,
        ),
    )
    headers = user_authentication_headers(client=client, email=email, password=password)

    r = client.get(f"{settings.API_V1_STR}/content/base/999?language=pt-br", headers=headers)
    assert r.status_code == 404
    assert r.json()["detail"] == "Content day not found"


def test_read_base_content_invalid_language_returns_422(
    client: TestClient, db: Session
) -> None:
    email = random_email()
    password = random_lower_string()
    crud.create_user(
        session=db,
        user_create=UserCreate(
            email=email,
            password=password,
            full_name="Base User",
            profile_bitmap=2,
        ),
    )
    headers = user_authentication_headers(client=client, email=email, password=password)

    r = client.get(f"{settings.API_V1_STR}/content/base/1?language=en-us", headers=headers)
    assert r.status_code == 422
