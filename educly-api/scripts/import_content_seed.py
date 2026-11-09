#!/usr/bin/env python3
"""Importa dados de conteudo/topicos para a API admin.

Entrada esperada:
- db_seed.json gerado por scripts/convert_desafio_to_db_model.py

Fluxo:
1) Login em /api/v1/login/access-token (superuser)
2) Upsert de content/items por (plan, language, day_id, content_index)
3) Importacao de topicos por dia
   - modo padrao: remove topicos existentes do dia e recria todos
"""

from __future__ import annotations

import argparse
import json
import sys
from dataclasses import dataclass
from pathlib import Path
from typing import Any
from urllib import error, parse, request


@dataclass
class ApiClient:
    base_url: str
    token: str | None = None

    def _url(self, path: str) -> str:
        return f"{self.base_url.rstrip('/')}{path}"

    def _headers(self, *, json_body: bool = False) -> dict[str, str]:
        headers: dict[str, str] = {}
        if self.token:
            headers["Authorization"] = f"Bearer {self.token}"
        if json_body:
            headers["Content-Type"] = "application/json"
        return headers

    def request_json(
        self,
        method: str,
        path: str,
        *,
        body: dict[str, Any] | None = None,
        query: dict[str, Any] | None = None,
        form: dict[str, Any] | None = None,
    ) -> tuple[int, dict[str, Any] | list[Any] | str]:
        url = self._url(path)
        if query:
            normalized_query = {k: v for k, v in query.items() if v is not None}
            if normalized_query:
                url = f"{url}?{parse.urlencode(normalized_query)}"

        data: bytes | None = None
        headers = self._headers(json_body=body is not None)
        if body is not None:
            data = json.dumps(body, ensure_ascii=False).encode("utf-8")
        elif form is not None:
            data = parse.urlencode(form).encode("utf-8")
            headers["Content-Type"] = "application/x-www-form-urlencoded"

        req = request.Request(url, data=data, method=method, headers=headers)
        try:
            with request.urlopen(req, timeout=60) as resp:
                raw = resp.read().decode("utf-8")
                return resp.status, self._parse_json_or_text(raw)
        except error.HTTPError as exc:
            raw = exc.read().decode("utf-8", errors="replace")
            return exc.code, self._parse_json_or_text(raw)

    @staticmethod
    def _parse_json_or_text(raw: str) -> dict[str, Any] | list[Any] | str:
        text = raw.strip()
        if not text:
            return ""
        try:
            return json.loads(text)
        except json.JSONDecodeError:
            return text


def read_env_file(path: Path) -> dict[str, str]:
    env: dict[str, str] = {}
    if not path.exists():
        return env
    for line in path.read_text(encoding="utf-8").splitlines():
        stripped = line.strip()
        if not stripped or stripped.startswith("#") or "=" not in stripped:
            continue
        key, value = stripped.split("=", 1)
        env[key.strip()] = value.strip().strip('"').strip("'")
    return env


def ensure(condition: bool, message: str) -> None:
    if not condition:
        raise RuntimeError(message)


def build_content_key(content: dict[str, Any]) -> tuple[str, str, int, int]:
    return (
        str(content["plan"]),
        str(content["language"]),
        int(content["day_id"]),
        int(content.get("content_index", 1)),
    )


def group_topics_by_day(topics: list[dict[str, Any]]) -> dict[int, list[dict[str, Any]]]:
    grouped: dict[int, list[dict[str, Any]]] = {}
    for topic in topics:
        day_id = int(topic["day_id"])
        grouped.setdefault(day_id, []).append(topic)
    for day_id in grouped:
        grouped[day_id].sort(key=lambda t: int(t.get("order", 1)))
    return grouped


def login(client: ApiClient, email: str, password: str) -> str:
    status, payload = client.request_json(
        "POST",
        "/api/v1/login/access-token",
        form={"username": email, "password": password},
    )
    ensure(status == 200, f"Falha no login ({status}): {payload}")
    ensure(isinstance(payload, dict), f"Resposta de login invalida: {payload}")
    token = payload.get("access_token")
    ensure(isinstance(token, str) and token, "Token ausente na resposta de login")
    return token


def fetch_existing_contents(
    client: ApiClient, plan: str | None, language: str | None
) -> dict[tuple[str, str, int, int], dict[str, Any]]:
    status, payload = client.request_json(
        "GET",
        "/api/v1/content/items",
        query={"skip": 0, "limit": 500, "plan": plan, "language": language},
    )
    ensure(status == 200, f"Falha ao listar contents ({status}): {payload}")
    ensure(isinstance(payload, dict), f"Payload de contents invalido: {payload}")
    data = payload.get("data", [])
    ensure(isinstance(data, list), "Campo data de contents invalido")

    indexed: dict[tuple[str, str, int, int], dict[str, Any]] = {}
    for item in data:
        if not isinstance(item, dict):
            continue
        key = (
            str(item.get("plan", "")),
            str(item.get("language", "")),
            int(item.get("day_id", 0)),
            int(item.get("content_index", 1)),
        )
        indexed[key] = item
    return indexed


def fetch_topics(client: ApiClient, content_id: int) -> list[dict[str, Any]]:
    status, payload = client.request_json(
        "GET",
        f"/api/v1/content/items/{content_id}/topics",
        query={"skip": 0, "limit": 500},
    )
    ensure(status == 200, f"Falha ao listar topicos de {content_id} ({status}): {payload}")
    ensure(isinstance(payload, dict), "Payload de topicos invalido")
    data = payload.get("data", [])
    ensure(isinstance(data, list), "Campo data de topicos invalido")
    return [item for item in data if isinstance(item, dict)]


def delete_topic(client: ApiClient, topic_id: str) -> None:
    status, payload = client.request_json("DELETE", f"/api/v1/content/topics/{topic_id}")
    ensure(status == 200, f"Falha ao deletar topico {topic_id} ({status}): {payload}")


def create_topic(client: ApiClient, content_id: int, topic: dict[str, Any]) -> None:
    request_body = {
        "content_id": content_id,
        "behavior": topic["behavior"],
        "activity_type": topic["activity_type"],
        "payload": topic["payload"],
        "order": int(topic.get("order", 1)),
    }
    status, payload = client.request_json(
        "POST",
        f"/api/v1/content/items/{content_id}/topics",
        body=request_body,
    )
    ensure(status == 200, f"Falha ao criar topico (day {topic['day_id']}, order {topic.get('order')}) ({status}): {payload}")


def upsert_content(
    client: ApiClient,
    existing_by_key: dict[tuple[str, str, int, int], dict[str, Any]],
    content: dict[str, Any],
    *,
    allow_update: bool,
) -> tuple[int, str]:
    key = build_content_key(content)
    existing = existing_by_key.get(key)
    payload = {
        "plan": content["plan"],
        "language": content["language"],
        "day_id": int(content["day_id"]),
        "content_index": int(content.get("content_index", 1)),
        "title": content["title"],
        "content": content["content"],
    }

    if existing is None:
        status, response = client.request_json("POST", "/api/v1/content/items", body=payload)
        ensure(status == 200, f"Falha ao criar content day {payload['day_id']} ({status}): {response}")
        ensure(isinstance(response, dict) and isinstance(response.get("id"), int), f"Resposta invalida ao criar content: {response}")
        content_id = int(response["id"])
        existing_by_key[key] = response
        return content_id, "created"

    content_id = int(existing["id"])
    if allow_update:
        status, response = client.request_json(
            "PUT",
            f"/api/v1/content/items/{content_id}",
            body=payload,
        )
        ensure(status == 200, f"Falha ao atualizar content {content_id} ({status}): {response}")
        ensure(isinstance(response, dict), f"Resposta invalida ao atualizar content {content_id}: {response}")
        existing_by_key[key] = response
        return content_id, "updated"

    return content_id, "skipped"


def main() -> int:
    parser = argparse.ArgumentParser(
        description="Importa db_seed.json para o modelo atual de content/topic via API admin."
    )
    parser.add_argument("--seed-file", default="tmp/desafio_iniciante_ia/db_seed.json")
    parser.add_argument("--api-base-url", default="http://localhost:8000")
    parser.add_argument("--email", default=None, help="Superuser email. Default: FIRST_SUPERUSER do .env")
    parser.add_argument("--password", default=None, help="Superuser password. Default: FIRST_SUPERUSER_PASSWORD do .env")
    parser.add_argument(
        "--upsert",
        action="store_true",
        help="Se content existir, atualiza title/content antes de importar topicos.",
    )
    parser.add_argument(
        "--replace-topics",
        action="store_true",
        default=True,
        help="Remove topicos existentes do dia antes de recriar (padrao: true).",
    )
    parser.add_argument(
        "--no-replace-topics",
        action="store_true",
        help="Nao remove topicos existentes. Apenas adiciona novos topicos do seed.",
    )
    args = parser.parse_args()

    seed_path = Path(args.seed_file)
    ensure(seed_path.exists(), f"Arquivo de seed nao encontrado: {seed_path}")
    seed = json.loads(seed_path.read_text(encoding="utf-8"))
    ensure(isinstance(seed, dict), "Formato invalido do seed: raiz deve ser objeto")

    contents = seed.get("contents", [])
    topics = seed.get("topics", [])
    ensure(isinstance(contents, list) and contents, "Seed sem contents")
    ensure(isinstance(topics, list), "Seed com topics invalido")

    env = read_env_file(Path(".env"))
    email = args.email or env.get("FIRST_SUPERUSER")
    password = args.password or env.get("FIRST_SUPERUSER_PASSWORD")
    ensure(email is not None and email.strip(), "Informe --email ou FIRST_SUPERUSER no .env")
    ensure(password is not None and password.strip(), "Informe --password ou FIRST_SUPERUSER_PASSWORD no .env")

    client = ApiClient(base_url=args.api_base_url.rstrip("/"))
    client.token = login(client, email=email.strip(), password=password.strip())

    first_content = contents[0]
    plan = str(first_content.get("plan", ""))
    language = str(first_content.get("language", ""))
    existing_by_key = fetch_existing_contents(client, plan=plan, language=language)
    topics_by_day = group_topics_by_day([t for t in topics if isinstance(t, dict)])

    replace_topics = not args.no_replace_topics

    created_count = 0
    updated_count = 0
    skipped_count = 0
    deleted_topics_count = 0
    created_topics_count = 0

    contents_sorted = sorted(
        (c for c in contents if isinstance(c, dict)),
        key=lambda c: (int(c.get("day_id", 0)), int(c.get("content_index", 1))),
    )

    for content in contents_sorted:
        day_id = int(content["day_id"])
        content_id, action = upsert_content(
            client,
            existing_by_key,
            content,
            allow_update=args.upsert,
        )
        if action == "created":
            created_count += 1
        elif action == "updated":
            updated_count += 1
        else:
            skipped_count += 1

        day_topics = topics_by_day.get(day_id, [])
        if replace_topics:
            existing_topics = fetch_topics(client, content_id)
            for existing_topic in existing_topics:
                topic_id = existing_topic.get("id")
                if isinstance(topic_id, str) and topic_id:
                    delete_topic(client, topic_id)
                    deleted_topics_count += 1

        for topic in day_topics:
            create_topic(client, content_id, topic)
            created_topics_count += 1

        print(
            f"[day {day_id:02d}] content={action} topics_imported={len(day_topics)}"
        )

    print("Importacao concluida.")
    print(
        "Resumo: "
        f"contents(created={created_count}, updated={updated_count}, skipped={skipped_count}) | "
        f"topics(deleted={deleted_topics_count}, created={created_topics_count})"
    )
    return 0


if __name__ == "__main__":
    try:
        raise SystemExit(main())
    except RuntimeError as exc:
        print(f"Erro: {exc}", file=sys.stderr)
        raise SystemExit(1)
