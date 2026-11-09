import uuid
from datetime import datetime, timezone
from enum import Enum
from typing import Any

from pydantic import EmailStr
from sqlalchemy import DateTime
import sqlalchemy as sa
from sqlmodel import Field, SQLModel


def get_datetime_utc() -> datetime:
    return datetime.now(timezone.utc)


# Shared properties
class UserBase(SQLModel):
    full_name: str = Field(default="", max_length=255)
    email: EmailStr = Field(unique=True, index=True, max_length=255)
    is_active: bool = True
    is_superuser: bool = False
    created_at: datetime | None = Field(
        default_factory=get_datetime_utc,
        sa_type=DateTime(timezone=True),  # type: ignore
    )
    updated_at: datetime | None = Field(
        default_factory=get_datetime_utc,
        sa_type=DateTime(timezone=True),  # type: ignore
    )
    last_login: datetime | None = Field(
        default=None,
        sa_type=DateTime(timezone=True),  # type: ignore
    )


# Properties to receive via API on creation
class UserCreate(SQLModel):
    full_name: str = Field(default="", max_length=255)
    email: EmailStr = Field(max_length=255)
    password: str = Field(min_length=8, max_length=128)
    is_active: bool = True
    is_superuser: bool = False
    profile_bitmap: int = Field(default=1, ge=1, le=15)
    personalized_trail: bool = False
    avatar_url: str | None = Field(default=None, max_length=255)
    cover_url: str | None = Field(default=None, max_length=255)
    bio: str | None = Field(default=None, max_length=500)


class UserRegister(SQLModel):
    email: EmailStr = Field(max_length=255)
    password: str = Field(min_length=8, max_length=128)
    full_name: str = Field(min_length=1, max_length=255)
    profile_bitmap: int = Field(default=1, ge=1, le=15)
    personalized_trail: bool = False
    avatar_url: str | None = Field(default=None, max_length=255)
    cover_url: str | None = Field(default=None, max_length=255)
    bio: str | None = Field(default=None, max_length=500)


# Properties to receive via API on update, all are optional
class UserUpdate(SQLModel):
    full_name: str | None = Field(default=None, max_length=255)
    email: EmailStr | None = Field(default=None, max_length=255)
    password: str | None = Field(default=None, min_length=8, max_length=128)
    is_active: bool | None = None
    is_superuser: bool | None = None

class UserUpdateMe(SQLModel):
    full_name: str | None = Field(default=None, min_length=1, max_length=255)
    email: EmailStr | None = Field(default=None, max_length=255)

class UpdatePassword(SQLModel):
    current_password: str = Field(min_length=8, max_length=128)
    new_password: str = Field(min_length=8, max_length=128)


# Database model, database table inferred from class name
class User(UserBase, table=True):
    id: uuid.UUID = Field(default_factory=uuid.uuid4, primary_key=True)
    hashed_password: str = Field(max_length=255)


# Properties to return via API, id is always required
class UserPublic(UserBase):
    id: uuid.UUID


class UsersPublic(SQLModel):
    data: list[UserPublic]
    count: int

# Base perfil

class PerfilBase(SQLModel):
    # Bitmask: 1 = free, 2 = base, 4 = freelancer, 8 = hub_ia
    profile_bitmap: int = Field(default=1, ge=1, le=15)
    personalized_trail: bool = False
    avatar_url: str | None = Field(default=None, max_length=255)
    cover_url: str | None = Field(default=None, max_length=255)
    bio: str | None = Field(default=None, max_length=500)

# tabela perfil

class Perfil(PerfilBase, table=True):
    id: uuid.UUID = Field(default_factory=uuid.uuid4, primary_key=True)
    user_id: uuid.UUID = Field(foreign_key="user.id", index=True)

# Put Perfil Patch Perfil

class PerfilUpdate(SQLModel):
    profile_bitmap: int | None = Field(default=None, ge=1, le=15)
    personalized_trail: bool | None = Field(default=None)
    avatar_url: str | None = Field(default=None, max_length=255)
    cover_url: str | None = Field(default=None, max_length=255)
    bio: str | None = Field(default=None, max_length=500)


class PerfilBitmapUpdate(SQLModel):
    profile_bitmap: int = Field(ge=1, le=15)


# Patch proprio perfil
class PerfilUpdateMe(SQLModel):
    avatar_url: str | None = Field(default=None, max_length=255)
    cover_url: str | None = Field(default=None, max_length=255)
    bio: str | None = Field(default=None, max_length=500)


# resp api perfil
class PerfilPublic(PerfilBase):
    id: uuid.UUID
    user_id: uuid.UUID

# create perfil

class PerfilCreate(PerfilBase):
    user_id: uuid.UUID


class PerfilsPublic(SQLModel):
    data: list[PerfilPublic]
    count: int


# Generic message
class Message(SQLModel):
    message: str


# JSON payload containing access token
class Token(SQLModel):
    access_token: str
    token_type: str = "bearer"


# Contents of JWT token
class TokenPayload(SQLModel):
    sub: str | None = None


class NewPassword(SQLModel):
    token: str
    new_password: str = Field(min_length=8, max_length=128)

# Produtos Hub IA | Freelancer | Base
class ProductBase(SQLModel):
    name: str = Field(min_length=1, max_length=255)
    product_id: str = Field(unique=True, index=True, max_length=255)
    description: str | None = Field(default=None, max_length=255)
    product_color: str = Field(default="green", max_length=20)
    product_codes: list[str] = Field(
        default_factory=list,
        sa_column=sa.Column(sa.JSON(), nullable=False, server_default="[]"),
    )
    valor: int = Field(default=0, ge=0)


class ProductCreate(ProductBase):
    pass


class ProductUpdate(SQLModel):
    name: str | None = Field(default=None, max_length=255)
    product_id: str | None = Field(default=None, max_length=255)
    description: str | None = Field(default=None, max_length=255)
    product_color: str | None = Field(default=None, max_length=20)
    product_codes: list[str] | None = Field(default=None)
    valor: int | None = Field(default=None, ge=0)

class ProductPublic(ProductBase):
    id: uuid.UUID


class ProductsPublic(SQLModel):
    data: list[ProductPublic]
    count: int


class Product(ProductBase, table=True):
    id: uuid.UUID = Field(default_factory=uuid.uuid4, primary_key=True)
    created_at: datetime | None = Field(
        default_factory=get_datetime_utc,
        sa_type=DateTime(timezone=True), # type: ignore
    )



class WebhookEvent(SQLModel, table=True):
    __tablename__ = "webhook_event"

    id: uuid.UUID = Field(default_factory=uuid.uuid4, primary_key=True)
    provider: str = Field(max_length=100, index=True)
    url: str = Field(max_length=512)
    payload: dict[str, Any] = Field(
        default_factory=dict,
        sa_column=sa.Column(sa.JSON(), nullable=False, server_default="{}"),
    )
    headers: dict[str, Any] = Field(
        default_factory=dict,
        sa_column=sa.Column(sa.JSON(), nullable=False, server_default="{}"),
    )
    processed: bool = Field(default=False)
    created_at: datetime | None = Field(
        default_factory=get_datetime_utc,
        sa_type=DateTime(timezone=True), # type: ignore
    )


class WebhookEventPublic(SQLModel):
    id: uuid.UUID
    provider: str
    url: str
    payload: dict[str, Any]
    headers: dict[str, Any]
    processed: bool
    created_at: datetime | None = None


class WebhookEventsPublic(SQLModel):
    data: list[WebhookEventPublic]
    count: int


class WebhookPurchaseProductUpdate(SQLModel):
    product_id: uuid.UUID


class BaseContentPublic(SQLModel):
    id: int
    language: str
    plan: str
    title: str
    content: str


class TopicBehavior(str, Enum):
    atividade = "atividade"
    explicacao = "explicacao"


class ActivityType(str, Enum):
    MatchWords = "MatchWords"
    FillBlanks = "FillBlanks"
    Quizz = "Quizz"
    PromptTreiner = "PromptTreiner"


class MatchWordPair(SQLModel):
    text: str = Field(min_length=1)
    match: str = Field(min_length=1)


class MatchWordsPayload(SQLModel):
    title: str = Field(min_length=1, max_length=255)
    description: str = Field(min_length=1)
    pairs: list[MatchWordPair] = Field(min_length=1)


class FillBlanksPayload(SQLModel):
    title: str = Field(min_length=1, max_length=255)
    sentence: str = Field(min_length=1)
    answer: str = Field(min_length=1)
    options: list[str] = Field(min_length=1)
    explanation: str = Field(min_length=1)


class QuizzOption(SQLModel):
    text: str = Field(min_length=1)
    is_correct: bool
    explanation: str = Field(min_length=1)


class QuizzPayload(SQLModel):
    title: str = Field(min_length=1, max_length=255)
    question: str = Field(min_length=1)
    options: list[QuizzOption] = Field(min_length=1)


class PromptTreinerPayload(SQLModel):
    title: str = Field(min_length=1, max_length=255)
    description: str = Field(min_length=1)
    AIToolContext: str = Field(min_length=1)
    successCriteria: str = Field(min_length=1)


class ExplanationPayload(SQLModel):
    title: str = Field(min_length=1, max_length=255)
    content: str = Field(min_length=1)


class ContentBase(SQLModel):
    plan: str = Field(min_length=1, max_length=50)
    language: str = Field(default="pt-br", max_length=10)
    day_id: int = Field(ge=1)
    content_index: int = Field(default=1, ge=1)
    title: str = Field(min_length=1, max_length=255)
    content: str = Field(min_length=1)


class ContentCreate(ContentBase):
    pass


class ContentUpdate(SQLModel):
    plan: str | None = Field(default=None, min_length=1, max_length=50)
    language: str | None = Field(default=None, max_length=10)
    day_id: int | None = Field(default=None, ge=1)
    content_index: int | None = Field(default=None, ge=1)
    title: str | None = Field(default=None, min_length=1, max_length=255)
    content: str | None = Field(default=None, min_length=1)


class ContentPublic(ContentBase):
    id: int
    created_at: datetime | None = None


class ContentsPublic(SQLModel):
    data: list[ContentPublic]
    count: int


class Content(ContentBase, table=True):
    __tablename__ = "content"

    id: int | None = Field(
        default=None,
        sa_column=sa.Column(sa.Integer, primary_key=True, autoincrement=True),
    )
    created_at: datetime | None = Field(
        default_factory=get_datetime_utc,
        sa_type=DateTime(timezone=True),  # type: ignore
    )


class TopicBase(SQLModel):
    content_id: int = Field(foreign_key="content.id", index=True)
    behavior: TopicBehavior = Field(
        sa_column=sa.Column(
            sa.Enum(TopicBehavior, name="topic_behavior_enum"),
            nullable=False,
        )
    )
    activity_type: ActivityType | None = Field(
        default=None,
        sa_column=sa.Column(
            sa.Enum(ActivityType, name="activity_type_enum"),
            nullable=True,
        ),
    )
    payload: dict[str, Any] = Field(
        default_factory=dict,
        sa_column=sa.Column(sa.JSON(), nullable=False, server_default="{}"),
    )
    title: str | None = Field(default=None, max_length=255)
    body: str | None = Field(default=None)
    order: int = Field(
        default=1,
        ge=1,
        sa_column=sa.Column("order_index", sa.Integer, nullable=False, server_default="1"),
    )


class TopicCreate(TopicBase):
    pass


class TopicUpdate(SQLModel):
    content_id: int | None = Field(default=None)
    behavior: TopicBehavior | None = None
    activity_type: ActivityType | None = None
    payload: dict[str, Any] | None = None
    title: str | None = Field(default=None, min_length=1, max_length=255)
    body: str | None = Field(default=None, min_length=1)
    order: int | None = Field(default=None, ge=1)


class TopicPublic(TopicBase):
    id: uuid.UUID
    created_at: datetime | None = None


class TopicsPublic(SQLModel):
    data: list[TopicPublic]
    count: int


class Topic(TopicBase, table=True):
    __tablename__ = "topic"

    id: uuid.UUID = Field(default_factory=uuid.uuid4, primary_key=True)
    created_at: datetime | None = Field(
        default_factory=get_datetime_utc,
        sa_type=DateTime(timezone=True),  # type: ignore
    )


class LogPermissionBase(SQLModel):
    action: str = Field(default="assign_access", max_length=255)
    email: str = Field(default="unknown", max_length=255)
    status: str = Field(default="pending", max_length=30)
    reason: str | None = Field(default=None, max_length=500)
    granted_plan: str | None = Field(default=None, max_length=255)


class LogPermissionCreate(LogPermissionBase):
    pass


class LogPermissionUpdate(SQLModel):
    action: str | None = Field(default=None, max_length=255)
    email: str | None = Field(default=None, max_length=255)
    status: str | None = Field(default=None, max_length=30)
    reason: str | None = Field(default=None, max_length=500)
    granted_plan: str | None = Field(default=None, max_length=255)


class LogPermissionPublic(LogPermissionBase):
    id: uuid.UUID
    created_at: datetime | None = None


class LogPermissionsPublic(SQLModel):
    data: list[LogPermissionPublic]
    count: int


class LogPermission(LogPermissionBase, table=True):
    __tablename__ = "log_permissions"

    id: uuid.UUID = Field(default_factory=uuid.uuid4, primary_key=True)
    created_at: datetime | None = Field(
        default_factory=get_datetime_utc,
        sa_type=DateTime(timezone=True), # type: ignore
    )
