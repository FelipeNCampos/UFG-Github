import json
import uuid
from functools import lru_cache
from pathlib import Path
from typing import Any, Literal

from fastapi import APIRouter, Depends, HTTPException, Query
from pydantic import ValidationError
from sqlalchemy.exc import IntegrityError
from sqlmodel import col, func, select

from app import crud
from app.api.deps import CurrentUser, SessionDep, get_current_active_superuser
from app.models import (
    ActivityType,
    BaseContentPublic,
    Content,
    ContentCreate,
    ContentPublic,
    ContentsPublic,
    ContentUpdate,
    ExplanationPayload,
    FillBlanksPayload,
    MatchWordsPayload,
    Message,
    PromptTreinerPayload,
    QuizzPayload,
    Topic,
    TopicBehavior,
    TopicCreate,
    TopicPublic,
    TopicsPublic,
    TopicUpdate,
)

router = APIRouter(prefix="/content", tags=["content"])

BASE_PLAN_BIT = 2
CONTENT_DIR = Path(__file__).resolve().parents[3] / "content"


def _topic_422(detail: str) -> HTTPException:
    return HTTPException(status_code=422, detail=detail)


def _normalize_topic_fields(
    *,
    behavior: TopicBehavior,
    activity_type: ActivityType | None,
    payload: dict[str, Any],
) -> tuple[str, str]:
    if behavior == TopicBehavior.explicacao:
        if activity_type is not None:
            raise _topic_422("activity_type must be null when behavior is explicacao")
        try:
            explanation = ExplanationPayload.model_validate(payload)
        except ValidationError as exc:
            raise _topic_422(f"Invalid explicacao payload: {exc}") from exc
        return explanation.title, explanation.content

    if behavior == TopicBehavior.atividade:
        if activity_type is None:
            raise _topic_422("activity_type is required when behavior is atividade")

        try:
            if activity_type == ActivityType.MatchWords:
                activity = MatchWordsPayload.model_validate(payload)
                return activity.title, activity.description
            if activity_type == ActivityType.FillBlanks:
                activity = FillBlanksPayload.model_validate(payload)
                return activity.title, activity.sentence
            if activity_type == ActivityType.Quizz:
                activity = QuizzPayload.model_validate(payload)
                return activity.title, activity.question
            if activity_type == ActivityType.PromptTreiner:
                activity = PromptTreinerPayload.model_validate(payload)
                return activity.title, activity.description
        except ValidationError as exc:
            raise _topic_422(f"Invalid atividade payload: {exc}") from exc

        raise _topic_422(f"Unsupported activity_type: {activity_type}")

    raise _topic_422(f"Unsupported behavior: {behavior}")


@lru_cache(maxsize=32)
def load_plan_content(*, plan: str, language: str) -> dict[int, dict[str, str]]:
    file_path = CONTENT_DIR / plan / f"{language}.json"
    if not file_path.exists():
        raise HTTPException(status_code=404, detail="Content file not found")

    raw_data: Any
    try:
        with file_path.open("r", encoding="utf-8") as file:
            raw_data = json.load(file)
    except json.JSONDecodeError as exc:
        raise HTTPException(status_code=500, detail="Invalid content file format") from exc

    days = raw_data.get("days")
    if not isinstance(days, list):
        raise HTTPException(status_code=500, detail="Invalid content file structure")

    by_day: dict[int, dict[str, str]] = {}
    for item in days:
        if not isinstance(item, dict):
            continue
        day_id = item.get("id")
        title = item.get("title")
        content = item.get("content")
        if (
            isinstance(day_id, int)
            and isinstance(title, str)
            and isinstance(content, str)
        ):
            by_day[day_id] = {"title": title, "content": content}

    return by_day


@router.get("/base/{id}", response_model=BaseContentPublic)
def read_base_content_by_day(
    *,
    id: int,
    language: Literal["pt-br"] = Query(default="pt-br"),
    session: SessionDep,
    current_user: CurrentUser,
) -> BaseContentPublic:
    perfil = crud.get_perfil_by_user_id(session=session, user_id=current_user.id)
    user_bitmap = int(perfil.profile_bitmap) if perfil else 1
    if (user_bitmap & BASE_PLAN_BIT) != BASE_PLAN_BIT:
        raise HTTPException(status_code=403, detail="Base plan access required")

    plan_content = load_plan_content(plan="base", language=language)
    day_content = plan_content.get(id)
    if not day_content:
        raise HTTPException(status_code=404, detail="Content day not found")

    return BaseContentPublic(
        id=id,
        language=language,
        plan="base",
        title=day_content["title"],
        content=day_content["content"],
    )


@router.get(
    "/items",
    response_model=ContentsPublic,
    dependencies=[Depends(get_current_active_superuser)],
)
def read_contents(
    *,
    session: SessionDep,
    skip: int = 0,
    limit: int = 100,
    plan: str | None = Query(default=None),
    language: str | None = Query(default=None),
    search: str | None = Query(default=None),
) -> Any:
    filters = []
    if plan:
        filters.append(Content.plan == plan)
    if language:
        filters.append(Content.language == language)
    if search:
        normalized = f"%{search.strip().lower()}%"
        filters.append(
            (
                func.lower(col(Content.title)).like(normalized)
                | func.lower(col(Content.content)).like(normalized)
            )
        )

    count_statement = select(func.count()).select_from(Content)
    for condition in filters:
        count_statement = count_statement.where(condition)
    count = session.exec(count_statement).one()

    statement = select(Content).order_by(
        col(Content.plan).asc(),
        col(Content.language).asc(),
        col(Content.day_id).asc(),
        col(Content.content_index).asc(),
    )
    for condition in filters:
        statement = statement.where(condition)
    statement = statement.offset(skip).limit(limit)

    contents = session.exec(statement).all()
    return ContentsPublic(data=list(contents), count=count)


@router.post(
    "/items",
    response_model=ContentPublic,
    dependencies=[Depends(get_current_active_superuser)],
)
def create_content(*, session: SessionDep, content_in: ContentCreate) -> Any:
    content = Content.model_validate(content_in)
    session.add(content)
    try:
        session.commit()
    except IntegrityError as exc:
        session.rollback()
        raise HTTPException(
            status_code=409,
            detail="Content with this plan/language/day/index already exists",
        ) from exc
    session.refresh(content)
    return content


@router.put(
    "/items/{content_id}",
    response_model=ContentPublic,
    dependencies=[Depends(get_current_active_superuser)],
)
def update_content(
    *,
    session: SessionDep,
    content_id: int,
    content_in: ContentUpdate,
) -> Any:
    content = session.get(Content, content_id)
    if not content:
        raise HTTPException(status_code=404, detail="Content not found")

    update_data = content_in.model_dump(exclude_unset=True)
    content.sqlmodel_update(update_data)
    session.add(content)
    try:
        session.commit()
    except IntegrityError as exc:
        session.rollback()
        raise HTTPException(
            status_code=409,
            detail="Content with this plan/language/day/index already exists",
        ) from exc
    session.refresh(content)
    return content


@router.delete(
    "/items/{content_id}",
    response_model=Message,
    dependencies=[Depends(get_current_active_superuser)],
)
def delete_content(*, session: SessionDep, content_id: int) -> Message:
    content = session.get(Content, content_id)
    if not content:
        raise HTTPException(status_code=404, detail="Content not found")

    related_topics = session.exec(
        select(Topic).where(Topic.content_id == content_id)
    ).all()
    for topic in related_topics:
        session.delete(topic)

    session.delete(content)
    session.commit()
    return Message(message="Content deleted successfully")


@router.get(
    "/items/{content_id}/topics",
    response_model=TopicsPublic,
    dependencies=[Depends(get_current_active_superuser)],
)
def read_topics_by_content(
    *,
    session: SessionDep,
    content_id: int,
    skip: int = 0,
    limit: int = 100,
    behavior: TopicBehavior | None = Query(default=None),
    search: str | None = Query(default=None),
) -> Any:
    content = session.get(Content, content_id)
    if not content:
        raise HTTPException(status_code=404, detail="Content not found")

    filters = [Topic.content_id == content_id]
    if behavior:
        filters.append(Topic.behavior == behavior)
    if search:
        normalized = f"%{search.strip().lower()}%"
        filters.append(
            (
                func.lower(func.coalesce(col(Topic.title), "")).like(normalized)
                | func.lower(func.coalesce(col(Topic.body), "")).like(normalized)
            )
        )

    count_statement = select(func.count()).select_from(Topic)
    for condition in filters:
        count_statement = count_statement.where(condition)
    count = session.exec(count_statement).one()

    statement = select(Topic).order_by(col(Topic.order).asc(), col(Topic.created_at).asc())
    for condition in filters:
        statement = statement.where(condition)
    statement = statement.offset(skip).limit(limit)

    topics = session.exec(statement).all()
    return TopicsPublic(data=list(topics), count=count)


@router.post(
    "/items/{content_id}/topics",
    response_model=TopicPublic,
    dependencies=[Depends(get_current_active_superuser)],
)
def create_topic(
    *,
    session: SessionDep,
    content_id: int,
    topic_in: TopicCreate,
) -> Any:
    content = session.get(Content, content_id)
    if not content:
        raise HTTPException(status_code=404, detail="Content not found")
    if topic_in.content_id != content_id:
        raise HTTPException(
            status_code=400,
            detail="Topic content_id must match URL content_id",
        )

    normalized_title, normalized_body = _normalize_topic_fields(
        behavior=topic_in.behavior,
        activity_type=topic_in.activity_type,
        payload=topic_in.payload,
    )
    topic = Topic.model_validate(topic_in)
    topic.title = normalized_title
    topic.body = normalized_body
    session.add(topic)
    session.commit()
    session.refresh(topic)
    return topic


@router.put(
    "/topics/{topic_id}",
    response_model=TopicPublic,
    dependencies=[Depends(get_current_active_superuser)],
)
def update_topic(
    *,
    session: SessionDep,
    topic_id: uuid.UUID,
    topic_in: TopicUpdate,
) -> Any:
    topic = session.get(Topic, topic_id)
    if not topic:
        raise HTTPException(status_code=404, detail="Topic not found")

    update_data = topic_in.model_dump(exclude_unset=True)
    if "content_id" in update_data and update_data["content_id"] is not None:
        if not session.get(Content, update_data["content_id"]):
            raise HTTPException(status_code=404, detail="Target content not found")

    next_behavior = update_data.get("behavior", topic.behavior)
    next_activity_type = update_data.get("activity_type", topic.activity_type)
    next_payload = update_data.get("payload", topic.payload)
    normalized_title, normalized_body = _normalize_topic_fields(
        behavior=next_behavior,
        activity_type=next_activity_type,
        payload=next_payload,
    )
    update_data["title"] = normalized_title
    update_data["body"] = normalized_body

    topic.sqlmodel_update(update_data)
    session.add(topic)
    session.commit()
    session.refresh(topic)
    return topic


@router.delete(
    "/topics/{topic_id}",
    response_model=Message,
    dependencies=[Depends(get_current_active_superuser)],
)
def delete_topic(*, session: SessionDep, topic_id: uuid.UUID) -> Message:
    topic = session.get(Topic, topic_id)
    if not topic:
        raise HTTPException(status_code=404, detail="Topic not found")
    session.delete(topic)
    session.commit()
    return Message(message="Topic deleted successfully")
