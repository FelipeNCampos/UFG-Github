from typing import Any
import uuid

from fastapi import APIRouter, HTTPException, Query, Request
from sqlmodel import col, func, select

from app.api.deps import CurrentUser, SessionDep
from app.models import (
    Message,
    Product,
    WebhookEvent,
    WebhookEventPublic,
    WebhookEventsPublic,
    WebhookPurchaseProductUpdate,
)
from app.services.access_sync import (
    apply_access_from_payload,
    replace_product_in_payload,
    update_user_bitmap_for_product_change,
)

router = APIRouter(prefix="/webhooks", tags=["webhooks"])

PROVIDER_PRIMER = "primer"
PROVIDER_STRIPE_1 = "stripe-1"
PROVIDER_STRIPE_2 = "stripe-2"
PROVIDER_HOTMART = "hotmart"
PROVIDER_PADDLE = "paddle"


@router.get("/", response_model=WebhookEventsPublic)
def read_webhooks(
    session: SessionDep,
    current_user: CurrentUser,
    skip: int = 0,
    limit: int = 100,
    provider: str | None = None,
    q: str | None = Query(default=None, max_length=255),
) -> Any:
    if not current_user.is_superuser:
        raise HTTPException(status_code=403, detail="Not enough permissions")

    statement = select(WebhookEvent)
    count_statement = select(func.count()).select_from(WebhookEvent)

    if provider:
        statement = statement.where(WebhookEvent.provider == provider.lower())
        count_statement = count_statement.where(WebhookEvent.provider == provider.lower())

    if q:
        like = f"%{q.lower()}%"
        statement = statement.where(
            func.lower(WebhookEvent.provider).like(like) | func.lower(WebhookEvent.url).like(like)
        )
        count_statement = count_statement.where(
            func.lower(WebhookEvent.provider).like(like) | func.lower(WebhookEvent.url).like(like)
        )

    count = session.exec(count_statement).one()
    rows = session.exec(
        statement.order_by(col(WebhookEvent.created_at).desc()).offset(skip).limit(limit)
    ).all()
    return WebhookEventsPublic(data=list(rows), count=count)


async def save_webhook_event(
    provider: str,
    request: Request,
    session: SessionDep,
) -> Message:
    payload: Any
    try:
        payload = await request.json()
    except Exception:
        payload = {"raw_body": (await request.body()).decode("utf-8", errors="replace")}

    if not isinstance(payload, dict):
        payload = {"data": payload}

    processed = apply_access_from_payload(
        payload=payload, session=session, action=f"assign_access:webhook:{provider}"
    )

    event = WebhookEvent(
        provider=provider,
        url=str(request.url),
        payload=payload,
        headers=dict(request.headers),
        processed=processed,
    )
    session.add(event)
    session.commit()

    return Message(message="Webhook received")


@router.post("/primer", status_code=202, response_model=Message)
async def receive_primer_webhook(request: Request, session: SessionDep) -> Message:
    return await save_webhook_event(PROVIDER_PRIMER, request, session)


@router.post("/stripe-1", status_code=202, response_model=Message)
async def receive_stripe_1_webhook(request: Request, session: SessionDep) -> Message:
    return await save_webhook_event(PROVIDER_STRIPE_1, request, session)


@router.post("/stripe-2", status_code=202, response_model=Message)
async def receive_stripe_2_webhook(request: Request, session: SessionDep) -> Message:
    return await save_webhook_event(PROVIDER_STRIPE_2, request, session)


@router.post("/hotmart", status_code=202, response_model=Message)
async def receive_hotmart_webhook(request: Request, session: SessionDep) -> Message:
    return await save_webhook_event(PROVIDER_HOTMART, request, session)


@router.post("/paddle", status_code=202, response_model=Message)
async def receive_paddle_webhook(request: Request, session: SessionDep) -> Message:
    return await save_webhook_event(PROVIDER_PADDLE, request, session)


@router.patch("/{webhook_id}/purchase-product", response_model=WebhookEventPublic)
def update_webhook_purchase_product(
    *,
    session: SessionDep,
    current_user: CurrentUser,
    webhook_id: uuid.UUID,
    body: WebhookPurchaseProductUpdate,
) -> Any:
    if not current_user.is_superuser:
        raise HTTPException(status_code=403, detail="Not enough permissions")

    event = session.get(WebhookEvent, webhook_id)
    if not event:
        raise HTTPException(status_code=404, detail="Webhook event not found")

    product = session.get(Product, body.product_id)
    if not product:
        raise HTTPException(status_code=404, detail="Product not found")

    available_codes = [str(code).strip() for code in (product.product_codes or []) if str(code).strip()]
    if not available_codes:
        raise HTTPException(status_code=400, detail="Product has no product_codes configured")

    bitmap_updated = update_user_bitmap_for_product_change(
        session=session, payload=event.payload, new_product=product
    )
    if not bitmap_updated:
        raise HTTPException(
            status_code=400,
            detail="Could not update user bitmap for this webhook",
        )

    event.payload = replace_product_in_payload(
        payload=event.payload, product_code=available_codes[0]
    )
    event.processed = True
    session.add(event)
    session.commit()
    session.refresh(event)
    return event


@router.post("/{webhook_id}/retry-assign-access", response_model=Message)
def retry_assign_webhook_access(
    *,
    session: SessionDep,
    current_user: CurrentUser,
    webhook_id: uuid.UUID,
) -> Message:
    if not current_user.is_superuser:
        raise HTTPException(status_code=403, detail="Not enough permissions")

    event = session.get(WebhookEvent, webhook_id)
    if not event:
        raise HTTPException(status_code=404, detail="Webhook event not found")

    processed = apply_access_from_payload(
        payload=event.payload, session=session, action="assign_access:retry_manual"
    )
    if not processed:
        raise HTTPException(
            status_code=400,
            detail="Could not assign access for this webhook",
        )

    event.processed = True
    session.add(event)
    session.commit()
    return Message(message="Access assigned successfully")
