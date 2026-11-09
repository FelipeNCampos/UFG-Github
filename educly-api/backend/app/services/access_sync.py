from __future__ import annotations

import re
from typing import Any

from sqlmodel import Session, select

from app import crud
from app.models import LogPermission, Product, WebhookEvent

EMAIL_REGEX = re.compile(r"^[^\s@]+@[^\s@]+\.[^\s@]+$")


def walk_values(value: Any, out: list[str]) -> None:
    if value is None:
        return
    if isinstance(value, str):
        out.append(value)
        return
    if isinstance(value, (int, float, bool)):
        out.append(str(value))
        return
    if isinstance(value, list):
        for item in value:
            walk_values(item, out)
        return
    if isinstance(value, dict):
        for nested in value.values():
            walk_values(nested, out)


def extract_email(payload: dict[str, Any]) -> str | None:
    candidates: list[str] = []
    walk_values(payload, candidates)
    for value in candidates:
        text = value.strip().lower()
        if EMAIL_REGEX.match(text):
            return text
    return None


def extract_product_code(payload: dict[str, Any], known_codes: set[str]) -> str | None:
    candidates: list[str] = []
    walk_values(payload, candidates)
    for value in candidates:
        code = value.strip()
        if code in known_codes:
            return code
    return None


def get_code_to_bitmap_map(*, session: Session) -> dict[str, int]:
    products = session.exec(select(Product)).all()
    code_to_bitmap: dict[str, int] = {}
    for product in products:
        for code in product.product_codes or []:
            normalized = str(code).strip()
            if normalized:
                code_to_bitmap[normalized] = product.valor
    return code_to_bitmap


def log_permission_attempt(
    *,
    session: Session,
    email: str,
    status: str,
    reason: str,
    granted_plan: str | None,
    action: str = "assign_access",
) -> None:
    log_obj = LogPermission(
        action=action,
        email=(email or "unknown")[:255],
        status=status.lower(),
        reason=reason[:500],
        granted_plan=(granted_plan[:255] if granted_plan else None),
    )
    session.add(log_obj)
    session.commit()


def apply_access_from_payload(
    *, payload: dict[str, Any], session: Session, action: str = "assign_access:webhook"
) -> bool:
    code_to_bitmap = get_code_to_bitmap_map(session=session)
    code_to_plan_name: dict[str, str] = {}
    products = session.exec(select(Product)).all()
    for product in products:
        for code in product.product_codes or []:
            normalized = str(code).strip()
            if normalized:
                code_to_plan_name[normalized] = product.name
    if not code_to_bitmap:
        log_permission_attempt(
            session=session,
            email="unknown",
            status="error",
            reason="no_product_mapping",
            granted_plan=None,
            action=action,
        )
        return False

    user_email = extract_email(payload)
    if not user_email:
        log_permission_attempt(
            session=session,
            email="unknown",
            status="error",
            reason="email_not_found_in_payload",
            granted_plan=None,
            action=action,
        )
        return False

    product_code = extract_product_code(payload, set(code_to_bitmap.keys()))
    attempted_plan = code_to_plan_name.get(product_code) if product_code else None

    user = crud.get_user_by_email(session=session, email=user_email)
    if not user:
        log_permission_attempt(
            session=session,
            email=user_email,
            status="error",
            reason="user_not_found",
            granted_plan=attempted_plan,
            action=action,
        )
        return False

    if not product_code:
        log_permission_attempt(
            session=session,
            email=user_email,
            status="error",
            reason="product_code_not_found_in_payload",
            granted_plan=None,
            action=action,
        )
        return False

    profile_bitmap = code_to_bitmap[product_code]
    perfil = crud.get_perfil_by_user_id(session=session, user_id=user.id)
    current_bitmap = int(perfil.profile_bitmap) if perfil else 1
    next_bitmap = current_bitmap | int(profile_bitmap)
    if next_bitmap < 1:
        next_bitmap = 1
    crud.set_user_profile_bitmap(
        session=session, user_id=user.id, profile_bitmap=next_bitmap
    )
    granted_plan = code_to_plan_name.get(product_code, f"bitmap:{profile_bitmap}")
    log_permission_attempt(
        session=session,
        email=user_email,
        status="success",
        reason="access_assigned",
        granted_plan=granted_plan,
        action=action,
    )
    return True


def replace_product_in_payload(*, payload: dict[str, Any], product_code: str) -> dict[str, Any]:
    next_payload = dict(payload)
    data = dict(next_payload.get("data") or {})
    product = dict(data.get("product") or {})
    # Keep webhook traceability while forcing Look/product matching by known code.
    product["id"] = product_code
    data["product"] = product
    next_payload["data"] = data
    next_payload["manual_product_code"] = product_code
    return next_payload


def update_user_bitmap_for_product_change(
    *, session: Session, payload: dict[str, Any], new_product: Product
) -> bool:
    code_to_bitmap = get_code_to_bitmap_map(session=session)
    if not code_to_bitmap:
        return False

    user_email = extract_email(payload)
    if not user_email:
        return False

    user = crud.get_user_by_email(session=session, email=user_email)
    if not user:
        return False

    perfil = crud.get_perfil_by_user_id(session=session, user_id=user.id)
    if not perfil:
        return False

    old_code = extract_product_code(payload, set(code_to_bitmap.keys()))
    old_bitmap = code_to_bitmap.get(old_code, 0) if old_code else 0
    current_bitmap = int(perfil.profile_bitmap or 1)
    without_old = current_bitmap & ~old_bitmap if old_bitmap else current_bitmap
    updated_bitmap = without_old | int(new_product.valor)
    if updated_bitmap < 1:
        updated_bitmap = 1

    crud.set_user_profile_bitmap(
        session=session, user_id=user.id, profile_bitmap=updated_bitmap
    )
    return True


def process_pending_webhooks_for_email(*, email: str, session: Session) -> None:
    normalized_email = email.strip().lower()
    if not normalized_email:
        return

    user = crud.get_user_by_email(session=session, email=normalized_email)
    if not user:
        return

    pending_events = session.exec(
        select(WebhookEvent).where(WebhookEvent.processed.is_(False))
    ).all()
    for event in pending_events:
        event_email = extract_email(event.payload)
        if event_email != normalized_email:
            continue
        if apply_access_from_payload(
            payload=event.payload, session=session, action="assign_access:pending_reprocess"
        ):
            event.processed = True
            session.add(event)

    session.commit()
