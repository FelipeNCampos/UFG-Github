import uuid

from fastapi.testclient import TestClient
from sqlmodel import Session, col, select

from app import crud
from app.core.config import settings
from app.models import LogPermission, Perfil, Product, UserCreate, WebhookEvent


def test_receive_webhook_persists_event(client: TestClient, db: Session) -> None:
    payload = {
        "event": "PURCHASE_APPROVED",
        "purchase_id": "abc-123",
        "buyer": {"email": "test@example.com"},
    }
    headers = {"x-provider-signature": "sig_123"}

    r = client.post(
        f"{settings.API_V1_STR}/webhooks/hotmart",
        json=payload,
        headers=headers,
    )

    assert r.status_code == 202
    assert r.json()["message"] == "Webhook received"

    event = db.exec(
        select(WebhookEvent)
        .where(WebhookEvent.provider == "hotmart")
        .order_by(col(WebhookEvent.created_at).desc())
    ).first()
    assert event
    assert event.payload["event"] == "PURCHASE_APPROVED"
    assert event.payload["purchase_id"] == "abc-123"
    assert event.processed is False


def test_read_webhooks_superuser(
    client: TestClient,
    superuser_token_headers: dict[str, str],
) -> None:
    # seed a webhook event
    client.post(
        f"{settings.API_V1_STR}/webhooks/hotmart",
        json={"event": "PING"},
    )

    r = client.get(
        f"{settings.API_V1_STR}/webhooks/?provider=hotmart&q=hotmart",
        headers=superuser_token_headers,
    )
    assert r.status_code == 200
    data = r.json()
    assert data["count"] >= 1


def test_receive_webhook_for_all_providers(client: TestClient, db: Session) -> None:
    providers = ["primer", "stripe-1", "stripe-2", "hotmart", "paddle"]

    for provider in providers:
        r = client.post(
            f"{settings.API_V1_STR}/webhooks/{provider}",
            json={"event": "PING", "provider_echo": provider},
        )
        assert r.status_code == 202

        event = db.exec(
            select(WebhookEvent)
            .where(WebhookEvent.provider == provider)
            .order_by(col(WebhookEvent.created_at).desc())
        ).first()
        assert event
        assert event.payload["provider_echo"] == provider


def test_receive_webhook_updates_user_bitmap_from_product_codes(
    client: TestClient, db: Session
) -> None:
    email = f"webhook-access-{uuid.uuid4()}@test.com"
    user = crud.create_user(
        session=db,
        user_create=UserCreate(email=email, password="testpass123", full_name="Webhook User"),
    )

    base_product = db.exec(select(Product).where(Product.product_id == "BASE")).first()
    assert base_product
    assert "7018201" in set(base_product.product_codes)

    payload = {
        "event": "PURCHASE_APPROVED",
        "data": {
            "buyer": {"email": email},
            "product": {"id": 7018201},
        },
    }

    r = client.post(
        f"{settings.API_V1_STR}/webhooks/hotmart",
        json=payload,
    )
    assert r.status_code == 202

    perfil = db.exec(select(Perfil).where(Perfil.user_id == user.id)).first()
    assert perfil
    assert perfil.profile_bitmap == 2

    event = db.exec(
        select(WebhookEvent)
        .where(WebhookEvent.provider == "hotmart")
        .order_by(col(WebhookEvent.created_at).desc())
    ).first()
    assert event
    assert event.processed is True

    access_log = db.exec(
        select(LogPermission)
        .where(LogPermission.email == email)
        .order_by(col(LogPermission.created_at).desc())
    ).first()
    assert access_log
    assert access_log.status == "success"
    assert access_log.reason == "access_assigned"
    assert access_log.granted_plan == "Base"
    assert access_log.created_at is not None


def test_update_webhook_purchase_product_replaces_product_and_adjusts_bitmap(
    client: TestClient,
    db: Session,
    superuser_token_headers: dict[str, str],
) -> None:
    email = f"webhook-product-change-{uuid.uuid4()}@test.com"
    user = crud.create_user(
        session=db,
        user_create=UserCreate(email=email, password="testpass123", full_name="Webhook User"),
    )

    payload = {
        "event": "PURCHASE_APPROVED",
        "data": {
            "buyer": {"email": email},
            "product": {"id": 7018201},  # BASE -> 2
        },
    }
    create_webhook = client.post(
        f"{settings.API_V1_STR}/webhooks/hotmart",
        json=payload,
    )
    assert create_webhook.status_code == 202

    perfil_before = db.exec(select(Perfil).where(Perfil.user_id == user.id)).first()
    assert perfil_before
    assert perfil_before.profile_bitmap == 2

    event = db.exec(
        select(WebhookEvent)
        .where(WebhookEvent.provider == "hotmart")
        .order_by(col(WebhookEvent.created_at).desc())
    ).first()
    assert event

    freelancer_product = db.exec(
        select(Product).where(Product.product_id == "FREELANCER")
    ).first()
    assert freelancer_product
    assert freelancer_product.valor == 4

    update_response = client.patch(
        f"{settings.API_V1_STR}/webhooks/{event.id}/purchase-product",
        headers=superuser_token_headers,
        json={"product_id": str(freelancer_product.id)},
    )
    assert update_response.status_code == 200

    db.refresh(event)
    assert event.processed is True
    assert event.payload.get("manual_product_code") in set(freelancer_product.product_codes or [])

    db.expire_all()
    perfil_after = db.exec(select(Perfil).where(Perfil.user_id == user.id)).first()
    assert perfil_after
    assert perfil_after.profile_bitmap == 4


def test_retry_assign_webhook_access_success(
    client: TestClient,
    db: Session,
    superuser_token_headers: dict[str, str],
) -> None:
    email = f"webhook-retry-{uuid.uuid4()}@test.com"
    payload = {
        "event": "PURCHASE_APPROVED",
        "data": {
            "buyer": {"email": email},
            "product": {"id": 7018201},
        },
    }
    created = client.post(f"{settings.API_V1_STR}/webhooks/hotmart", json=payload)
    assert created.status_code == 202

    event = db.exec(
        select(WebhookEvent)
        .where(WebhookEvent.provider == "hotmart")
        .order_by(col(WebhookEvent.created_at).desc())
    ).first()
    assert event
    assert event.processed is False

    user = crud.create_user(
        session=db,
        user_create=UserCreate(email=email, password="testpass123", full_name="Retry User"),
    )
    perfil_before = db.exec(select(Perfil).where(Perfil.user_id == user.id)).first()
    assert perfil_before
    assert perfil_before.profile_bitmap == 1

    retry_response = client.post(
        f"{settings.API_V1_STR}/webhooks/{event.id}/retry-assign-access",
        headers=superuser_token_headers,
    )
    assert retry_response.status_code == 200
    assert retry_response.json()["message"] == "Access assigned successfully"

    db.expire_all()
    perfil_after = db.exec(select(Perfil).where(Perfil.user_id == user.id)).first()
    assert perfil_after
    assert perfil_after.profile_bitmap == 2


def test_receive_webhook_logs_error_when_user_does_not_exist(
    client: TestClient, db: Session
) -> None:
    email = f"webhook-missing-user-{uuid.uuid4()}@test.com"
    payload = {
        "event": "PURCHASE_APPROVED",
        "data": {
            "buyer": {"email": email},
            "product": {"id": 7018201},
        },
    }

    response = client.post(f"{settings.API_V1_STR}/webhooks/hotmart", json=payload)
    assert response.status_code == 202

    event = db.exec(
        select(WebhookEvent)
        .where(WebhookEvent.provider == "hotmart")
        .order_by(col(WebhookEvent.created_at).desc())
    ).first()
    assert event
    assert event.processed is False

    access_log = db.exec(
        select(LogPermission)
        .where(LogPermission.email == email)
        .order_by(col(LogPermission.created_at).desc())
    ).first()
    assert access_log
    assert access_log.status == "error"
    assert access_log.reason == "user_not_found"
    assert access_log.granted_plan == "Base"
    assert access_log.created_at is not None
