from collections.abc import Generator
import re

import pytest
from fastapi.testclient import TestClient
from sqlalchemy import inspect, text
from sqlalchemy.engine import URL, make_url
from sqlmodel import SQLModel, Session, delete, create_engine

import app.core.db as core_db
from app.api.deps import get_db
from app.core.config import settings
from app.core.db import init_db
from app.main import app
from app.models import LogPermission, Perfil, Product, User, WebhookEvent
from tests.utils.user import authentication_token_from_email
from tests.utils.utils import get_superuser_token_headers


def _build_test_database_url() -> URL:
    return make_url(str(settings.SQLALCHEMY_DATABASE_URI))


TEST_SCHEMA = "pytest_schema"
if not re.fullmatch(r"[A-Za-z0-9_]+", TEST_SCHEMA):
    raise RuntimeError(
        f"Invalid test schema '{TEST_SCHEMA}'. Use only letters, numbers, underscore."
    )

TEST_DATABASE_URL = _build_test_database_url()
engine = create_engine(
    TEST_DATABASE_URL,
    connect_args={"options": f"-csearch_path={TEST_SCHEMA}"},
)
core_db.engine = engine


@pytest.fixture(scope="session", autouse=True)
def db() -> Generator[Session, None, None]:
    with engine.connect() as connection:
        connection.execute(text(f'CREATE SCHEMA IF NOT EXISTS "{TEST_SCHEMA}"'))
        connection.commit()
    SQLModel.metadata.create_all(engine)
    with Session(engine) as session:
        init_db(session)
        yield session
        inspector = inspect(session.get_bind())
        table_names = set(inspector.get_table_names())

        if "perfil" in table_names:
            session.execute(delete(Perfil))
        if "product" in table_names:
            session.execute(delete(Product))
        if "webhook_event" in table_names:
            session.execute(delete(WebhookEvent))
        if "log_permissions" in table_names:
            session.execute(delete(LogPermission))
        if "user" in table_names:
            session.execute(delete(User))
        session.commit()
    with engine.connect() as connection:
        connection.execute(text(f'DROP SCHEMA IF EXISTS "{TEST_SCHEMA}" CASCADE'))
        connection.commit()


@pytest.fixture(scope="module")
def client() -> Generator[TestClient, None, None]:
    def override_get_db() -> Generator[Session, None, None]:
        with Session(engine) as session:
            yield session

    app.dependency_overrides[get_db] = override_get_db
    with TestClient(app) as c:
        yield c
    app.dependency_overrides.clear()


@pytest.fixture(scope="module")
def superuser_token_headers(client: TestClient) -> dict[str, str]:
    return get_superuser_token_headers(client)


@pytest.fixture(scope="module")
def normal_user_token_headers(client: TestClient, db: Session) -> dict[str, str]:
    return authentication_token_from_email(
        client=client, email=settings.EMAIL_TEST_USER, db=db
    )
