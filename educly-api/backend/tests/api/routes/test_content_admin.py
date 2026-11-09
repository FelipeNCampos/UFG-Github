from fastapi.testclient import TestClient

from app.core.config import settings


def test_content_admin_crud_and_filters(
    client: TestClient, superuser_token_headers: dict[str, str]
) -> None:
    create_content_payload = {
        "plan": "base",
        "language": "pt-br",
        "day_id": 901,
        "title": "Conteudo Admin Teste",
        "content": "Texto principal do conteudo",
    }

    create_content_response = client.post(
        f"{settings.API_V1_STR}/content/items",
        headers=superuser_token_headers,
        json=create_content_payload,
    )
    assert create_content_response.status_code == 200
    created_content = create_content_response.json()
    content_id = created_content["id"]

    list_content_response = client.get(
        f"{settings.API_V1_STR}/content/items",
        headers=superuser_token_headers,
        params={"plan": "base", "language": "pt-br", "search": "admin teste"},
    )
    assert list_content_response.status_code == 200
    listed_contents = list_content_response.json()["data"]
    assert any(item["id"] == content_id for item in listed_contents)

    update_content_response = client.put(
        f"{settings.API_V1_STR}/content/items/{content_id}",
        headers=superuser_token_headers,
        json={"title": "Conteudo Atualizado Admin Teste"},
    )
    assert update_content_response.status_code == 200
    assert update_content_response.json()["title"] == "Conteudo Atualizado Admin Teste"

    create_topic_payload = {
        "content_id": content_id,
        "behavior": "explicacao",
        "payload": {
            "title": "Topico 1",
            "content": "Detalhes do topico",
        },
        "order": 1,
    }
    create_topic_response = client.post(
        f"{settings.API_V1_STR}/content/items/{content_id}/topics",
        headers=superuser_token_headers,
        json=create_topic_payload,
    )
    assert create_topic_response.status_code == 200
    created_topic = create_topic_response.json()
    topic_id = created_topic["id"]

    list_topic_response = client.get(
        f"{settings.API_V1_STR}/content/items/{content_id}/topics",
        headers=superuser_token_headers,
        params={"behavior": "explicacao", "search": "topico"},
    )
    assert list_topic_response.status_code == 200
    listed_topics = list_topic_response.json()["data"]
    assert any(item["id"] == topic_id for item in listed_topics)

    update_topic_response = client.put(
        f"{settings.API_V1_STR}/content/topics/{topic_id}",
        headers=superuser_token_headers,
        json={
            "payload": {
                "title": "Topico 1 Atualizado",
                "content": "Detalhes do topico",
            },
        },
    )
    assert update_topic_response.status_code == 200
    assert update_topic_response.json()["title"] == "Topico 1 Atualizado"

    delete_topic_response = client.delete(
        f"{settings.API_V1_STR}/content/topics/{topic_id}",
        headers=superuser_token_headers,
    )
    assert delete_topic_response.status_code == 200
    assert delete_topic_response.json()["message"] == "Topic deleted successfully"

    delete_content_response = client.delete(
        f"{settings.API_V1_STR}/content/items/{content_id}",
        headers=superuser_token_headers,
    )
    assert delete_content_response.status_code == 200
    assert delete_content_response.json()["message"] == "Content deleted successfully"


def test_content_admin_requires_superuser(
    client: TestClient, normal_user_token_headers: dict[str, str]
) -> None:
    response = client.get(
        f"{settings.API_V1_STR}/content/items",
        headers=normal_user_token_headers,
    )
    assert response.status_code == 403
    assert response.json()["detail"] == "The user doesn't have enough privileges"


def test_content_admin_allows_same_plan_day_with_different_content_index(
    client: TestClient, superuser_token_headers: dict[str, str]
) -> None:
    payload_base = {
        "plan": "base",
        "language": "pt-br",
        "day_id": 902,
        "title": "Conteudo indice 1",
        "content": "Texto do conteudo indice 1",
    }
    payload_second = {
        "plan": "base",
        "language": "pt-br",
        "day_id": 902,
        "content_index": 2,
        "title": "Conteudo indice 2",
        "content": "Texto do conteudo indice 2",
    }

    create_first = client.post(
        f"{settings.API_V1_STR}/content/items",
        headers=superuser_token_headers,
        json=payload_base,
    )
    assert create_first.status_code == 200
    first_id = create_first.json()["id"]

    create_second = client.post(
        f"{settings.API_V1_STR}/content/items",
        headers=superuser_token_headers,
        json=payload_second,
    )
    assert create_second.status_code == 200
    second_body = create_second.json()
    second_id = second_body["id"]
    assert second_body["content_index"] == 2

    list_response = client.get(
        f"{settings.API_V1_STR}/content/items",
        headers=superuser_token_headers,
        params={"plan": "base", "language": "pt-br"},
    )
    assert list_response.status_code == 200
    data = list_response.json()["data"]
    ids = {item["id"] for item in data}
    assert first_id in ids
    assert second_id in ids

    for content_id in (first_id, second_id):
        delete_response = client.delete(
            f"{settings.API_V1_STR}/content/items/{content_id}",
            headers=superuser_token_headers,
        )
        assert delete_response.status_code == 200


def test_content_admin_create_topic_atividade_matchwords(
    client: TestClient, superuser_token_headers: dict[str, str]
) -> None:
    content_payload = {
        "plan": "base",
        "language": "pt-br",
        "day_id": 903,
        "title": "Conteudo para atividade",
        "content": "Texto principal",
    }
    create_content_response = client.post(
        f"{settings.API_V1_STR}/content/items",
        headers=superuser_token_headers,
        json=content_payload,
    )
    assert create_content_response.status_code == 200
    content_id = create_content_response.json()["id"]

    topic_payload = {
        "content_id": content_id,
        "behavior": "atividade",
        "activity_type": "MatchWords",
        "payload": {
            "title": "Relacionar termos",
            "description": "Associe os pares corretos.",
            "pairs": [
                {"text": "API", "match": "Interface"},
                {"text": "DB", "match": "Banco de dados"},
            ],
        },
        "order": 1,
    }
    create_topic_response = client.post(
        f"{settings.API_V1_STR}/content/items/{content_id}/topics",
        headers=superuser_token_headers,
        json=topic_payload,
    )
    assert create_topic_response.status_code == 200
    topic_data = create_topic_response.json()
    assert topic_data["activity_type"] == "MatchWords"
    assert topic_data["payload"]["title"] == "Relacionar termos"

    topic_id = topic_data["id"]
    delete_topic_response = client.delete(
        f"{settings.API_V1_STR}/content/topics/{topic_id}",
        headers=superuser_token_headers,
    )
    assert delete_topic_response.status_code == 200

    delete_content_response = client.delete(
        f"{settings.API_V1_STR}/content/items/{content_id}",
        headers=superuser_token_headers,
    )
    assert delete_content_response.status_code == 200
