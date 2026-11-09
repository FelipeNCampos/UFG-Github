import json
import logging
import os
from pathlib import Path
from typing import Any

from sqlmodel import Session, select

from app.models import ActivityType, Content, Topic, TopicBehavior

logger = logging.getLogger(__name__)

DEFAULT_SEED_FILE = "/app/backend/content/seeds/desafio_iniciante_ia/db_seed.json"


def _read_seed(path: Path) -> dict[str, Any]:
    raw = path.read_text(encoding="utf-8")
    # Accept legacy seed files that still contain trailing comment lines.
    clean = "\n".join(
        line for line in raw.splitlines() if not line.lstrip().startswith("//")
    )
    loaded = json.loads(clean)
    if not isinstance(loaded, dict):
        raise RuntimeError("Seed root must be an object")
    return loaded


def _topic_title_body(
    behavior: TopicBehavior,
    activity_type: ActivityType | None,
    payload: dict[str, Any],
) -> tuple[str | None, str | None]:
    if behavior == TopicBehavior.explicacao:
        return (
            str(payload.get("title", "")).strip() or None,
            str(payload.get("content", "")).strip() or None,
        )
    if behavior == TopicBehavior.atividade:
        if activity_type == ActivityType.MatchWords:
            return (
                str(payload.get("title", "")).strip() or None,
                str(payload.get("description", "")).strip() or None,
            )
        if activity_type == ActivityType.FillBlanks:
            return (
                str(payload.get("title", "")).strip() or None,
                str(payload.get("sentence", "")).strip() or None,
            )
        if activity_type == ActivityType.Quizz:
            return (
                str(payload.get("title", "")).strip() or None,
                str(payload.get("question", "")).strip() or None,
            )
        if activity_type == ActivityType.PromptTreiner:
            return (
                str(payload.get("title", "")).strip() or None,
                str(payload.get("description", "")).strip() or None,
            )
    return None, None


def bootstrap_content_seed(
    session: Session,
    *,
    seed_file: Path,
    replace_topics: bool = True,
) -> tuple[int, int, int, int]:
    seed = _read_seed(seed_file)
    raw_contents = seed.get("contents", [])
    raw_topics = seed.get("topics", [])
    if not isinstance(raw_contents, list) or not raw_contents:
        raise RuntimeError("Seed without contents")
    if not isinstance(raw_topics, list):
        raise RuntimeError("Seed topics must be a list")

    topics_by_day: dict[int, list[dict[str, Any]]] = {}
    for topic in raw_topics:
        if not isinstance(topic, dict):
            continue
        day_id = int(topic["day_id"])
        topics_by_day.setdefault(day_id, []).append(topic)
    for day_id in topics_by_day:
        topics_by_day[day_id].sort(key=lambda item: int(item.get("order", 1)))

    created_contents = 0
    updated_contents = 0
    deleted_topics = 0
    created_topics = 0

    contents = sorted(
        (c for c in raw_contents if isinstance(c, dict)),
        key=lambda c: (
            int(c.get("day_id", 0)),
            int(c.get("content_index", 1)),
            str(c.get("language", "")),
        ),
    )

    for content_in in contents:
        plan = str(content_in["plan"])
        language = str(content_in["language"])
        day_id = int(content_in["day_id"])
        content_index = int(content_in.get("content_index", 1))
        title = str(content_in["title"])
        content_text = str(content_in["content"])

        existing = session.exec(
            select(Content).where(
                Content.plan == plan,
                Content.language == language,
                Content.day_id == day_id,
                Content.content_index == content_index,
            )
        ).first()

        if existing is None:
            content = Content(
                plan=plan,
                language=language,
                day_id=day_id,
                content_index=content_index,
                title=title,
                content=content_text,
            )
            session.add(content)
            session.flush()
            created_contents += 1
        else:
            existing.title = title
            existing.content = content_text
            content = existing
            updated_contents += 1

        if replace_topics:
            existing_topics = session.exec(
                select(Topic).where(Topic.content_id == content.id)
            ).all()
            for old_topic in existing_topics:
                session.delete(old_topic)
                deleted_topics += 1
            session.flush()

        day_topics = topics_by_day.get(day_id, [])
        for topic_in in day_topics:
            behavior = TopicBehavior(str(topic_in["behavior"]))
            raw_activity = topic_in.get("activity_type")
            activity_type = (
                ActivityType(str(raw_activity))
                if raw_activity is not None
                else None
            )
            payload = topic_in.get("payload", {})
            if not isinstance(payload, dict):
                payload = {}
            title_value, body_value = _topic_title_body(
                behavior=behavior,
                activity_type=activity_type,
                payload=payload,
            )
            new_topic = Topic(
                content_id=int(content.id),
                behavior=behavior,
                activity_type=activity_type,
                payload=payload,
                title=title_value,
                body=body_value,
                order=int(topic_in.get("order", 1)),
            )
            session.add(new_topic)
            created_topics += 1

    session.commit()
    return created_contents, updated_contents, deleted_topics, created_topics


def main() -> None:
    logging.basicConfig(level=logging.INFO)
    enabled = os.getenv("BOOTSTRAP_CONTENT_SEED", "true").lower() in {
        "1",
        "true",
        "yes",
        "on",
    }
    if not enabled:
        logger.info("Content seed bootstrap disabled by BOOTSTRAP_CONTENT_SEED")
        return

    seed_file = Path(os.getenv("CONTENT_SEED_FILE", DEFAULT_SEED_FILE))
    replace_topics = os.getenv("CONTENT_SEED_REPLACE_TOPICS", "true").lower() in {
        "1",
        "true",
        "yes",
        "on",
    }
    if not seed_file.exists():
        logger.warning("Content seed file not found: %s", seed_file)
        return

    from app.core.db import engine

    with Session(engine) as session:
        created, updated, deleted_topics, created_topics = bootstrap_content_seed(
            session,
            seed_file=seed_file,
            replace_topics=replace_topics,
        )
    logger.info(
        "Content seed bootstrap done: contents(created=%s, updated=%s) topics(deleted=%s, created=%s)",
        created,
        updated,
        deleted_topics,
        created_topics,
    )


if __name__ == "__main__":
    main()
