from fastapi.testclient import TestClient

from app.core.config import settings


def test_create_and_read_log_permissions(
    client: TestClient,
    superuser_token_headers: dict[str, str],
) -> None:
    create_payload = {
        "action": "assign_access:webhook:primer",
        "email": "fulano@example.com",
        "reason": "usuario sem conta",
        "status": "pending",
        "granted_plan": "Base",
    }

    create_r = client.post(
        f"{settings.API_V1_STR}/log-permissions/",
        headers=superuser_token_headers,
        json=create_payload,
    )
    assert create_r.status_code == 200
    created = create_r.json()
    assert created["email"] == "fulano@example.com"

    read_r = client.get(
        f"{settings.API_V1_STR}/log-permissions/?status=pending&q=fulano",
        headers=superuser_token_headers,
    )
    assert read_r.status_code == 200
    data = read_r.json()
    assert data["count"] >= 1


def test_update_log_permission_not_found(
    client: TestClient,
    superuser_token_headers: dict[str, str],
) -> None:
    response = client.patch(
        f"{settings.API_V1_STR}/log-permissions/7d8df298-0f77-4b66-8b6a-3f6f7023af0a",
        headers=superuser_token_headers,
        json={"status": "error"},
    )
    assert response.status_code == 404
    assert response.json()["detail"] == "LogPermission not found"


def test_log_permissions_forbidden_for_normal_user(
    client: TestClient,
    normal_user_token_headers: dict[str, str],
) -> None:
    read_response = client.get(
        f"{settings.API_V1_STR}/log-permissions/",
        headers=normal_user_token_headers,
    )
    assert read_response.status_code == 403

    create_response = client.post(
        f"{settings.API_V1_STR}/log-permissions/",
        headers=normal_user_token_headers,
        json={
            "action": "assign_access:webhook",
            "email": "normal@example.com",
            "reason": "forbidden-test",
            "status": "pending",
            "granted_plan": "Base",
        },
    )
    assert create_response.status_code == 403
