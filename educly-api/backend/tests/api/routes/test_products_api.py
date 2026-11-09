import uuid

from fastapi.testclient import TestClient
from app.core.config import settings


def test_read_products_superuser(
    client: TestClient, superuser_token_headers: dict[str, str]
) -> None:
    response = client.get(
        f"{settings.API_V1_STR}/products/",
        headers=superuser_token_headers,
    )
    assert response.status_code == 200
    payload = response.json()
    assert "data" in payload
    assert "count" in payload
    assert payload["count"] >= 1


def test_read_products_forbidden_for_normal_user(
    client: TestClient, normal_user_token_headers: dict[str, str]
) -> None:
    response = client.get(
        f"{settings.API_V1_STR}/products/",
        headers=normal_user_token_headers,
    )
    assert response.status_code == 403


def test_create_update_delete_product_superuser(
    client: TestClient, superuser_token_headers: dict[str, str]
) -> None:
    unique_suffix = uuid.uuid4().hex[:8]
    create_payload = {
        "name": f"Produto Teste {unique_suffix}",
        "product_id": "FREE",
        "description": "Criado por teste",
        "product_color": "gray",
        "product_codes": [f"code-{unique_suffix}"],
        "valor": 999,
    }

    created_response = client.post(
        f"{settings.API_V1_STR}/products/",
        headers=superuser_token_headers,
        json=create_payload,
    )
    assert created_response.status_code == 200
    created = created_response.json()
    assert created["name"] == create_payload["name"]
    # valor and color follow reference mapping from product_id.
    assert created["valor"] == 1
    assert created["product_color"] == "gray"

    update_payload = {
        "name": f"Produto Atualizado {unique_suffix}",
        "description": "Atualizado por teste",
        "product_codes": [f"code-updated-{unique_suffix}"],
    }
    updated_response = client.put(
        f"{settings.API_V1_STR}/products/{created['id']}",
        headers=superuser_token_headers,
        json=update_payload,
    )
    assert updated_response.status_code == 200
    updated = updated_response.json()
    assert updated["name"] == update_payload["name"]
    assert updated["product_id"] == "FREE"
    assert updated["valor"] == 1
    assert updated["product_color"] == "gray"

    deleted_response = client.delete(
        f"{settings.API_V1_STR}/products/{created['id']}",
        headers=superuser_token_headers,
    )
    assert deleted_response.status_code == 200
    assert deleted_response.json()["message"] == "Product deleted successfully"


def test_create_product_invalid_reference_returns_400(
    client: TestClient, superuser_token_headers: dict[str, str]
) -> None:
    payload = {
        "name": "Produto invalido",
        "product_id": "PLANO_NAO_MAPEADO",
        "description": "deve falhar",
        "product_codes": ["invalid-reference-code"],
    }
    response = client.post(
        f"{settings.API_V1_STR}/products/",
        headers=superuser_token_headers,
        json=payload,
    )
    assert response.status_code == 400
    assert response.json()["detail"] == "Invalid product_id reference"


def test_update_product_not_found_returns_404(
    client: TestClient, superuser_token_headers: dict[str, str]
) -> None:
    response = client.put(
        f"{settings.API_V1_STR}/products/{uuid.uuid4()}",
        headers=superuser_token_headers,
        json={"name": "nao existe"},
    )
    assert response.status_code == 404
    assert response.json()["detail"] == "Product not found"


def test_update_product_duplicate_product_id_returns_400(
    client: TestClient, superuser_token_headers: dict[str, str]
) -> None:
    create_response = client.post(
        f"{settings.API_V1_STR}/products/",
        headers=superuser_token_headers,
        json={
            "name": "Produto Free Duplicate Test",
            "product_id": "FREE",
            "description": "temp",
            "product_codes": ["free-dup-code"],
        },
    )
    assert create_response.status_code == 200
    created = create_response.json()

    update_response = client.put(
        f"{settings.API_V1_STR}/products/{created['id']}",
        headers=superuser_token_headers,
        json={"product_id": "BASE"},
    )
    assert update_response.status_code == 400
    assert (
        update_response.json()["detail"]
        == "Product with this product_id already exists"
    )
