"""add content and topic models

Revision ID: a7c91d3e4f20
Revises: f1a6d0c8b2e1
Create Date: 2026-02-22 23:40:00.000000

"""

from alembic import op
import sqlalchemy as sa
from sqlalchemy.dialects import postgresql


# revision identifiers, used by Alembic.
revision = "a7c91d3e4f20"
down_revision = "f1a6d0c8b2e1"
branch_labels = None
depends_on = None


topic_behavior_enum = postgresql.ENUM(
    "atividade",
    "explicacao",
    name="topic_behavior_enum",
    create_type=False,
)


def upgrade():
    op.create_table(
        "content",
        sa.Column("id", sa.Uuid(), nullable=False),
        sa.Column("plan", sa.String(length=50), nullable=False),
        sa.Column("language", sa.String(length=10), nullable=False),
        sa.Column("day_id", sa.Integer(), nullable=False),
        sa.Column("content_index", sa.Integer(), nullable=False, server_default="1"),
        sa.Column("title", sa.String(length=255), nullable=False),
        sa.Column("content", sa.String(), nullable=False),
        sa.Column("created_at", sa.DateTime(timezone=True), nullable=True),
        sa.PrimaryKeyConstraint("id"),
        sa.UniqueConstraint(
            "plan",
            "language",
            "day_id",
            "content_index",
            name="uq_content_plan_language_day_index",
        ),
    )
    op.create_index(
        op.f("ix_content_plan"),
        "content",
        ["plan"],
        unique=False,
    )

    op.execute(
        """
        DO $$
        BEGIN
            IF NOT EXISTS (
                SELECT 1
                FROM pg_type
                WHERE typname = 'topic_behavior_enum'
            ) THEN
                CREATE TYPE topic_behavior_enum AS ENUM ('atividade', 'explicacao');
            END IF;
        END$$;
        """
    )
    op.create_table(
        "topic",
        sa.Column("id", sa.Uuid(), nullable=False),
        sa.Column("content_id", sa.Uuid(), nullable=False),
        sa.Column("title", sa.String(length=255), nullable=False),
        sa.Column("body", sa.String(), nullable=False),
        sa.Column("behavior", topic_behavior_enum, nullable=False),
        sa.Column("order_index", sa.Integer(), server_default="1", nullable=False),
        sa.Column("created_at", sa.DateTime(timezone=True), nullable=True),
        sa.ForeignKeyConstraint(["content_id"], ["content.id"]),
        sa.PrimaryKeyConstraint("id"),
    )
    op.create_index(op.f("ix_topic_content_id"), "topic", ["content_id"], unique=False)


def downgrade():
    op.drop_index(op.f("ix_content_plan"), table_name="content")
    op.drop_index(op.f("ix_topic_content_id"), table_name="topic")
    op.drop_table("topic")
    op.drop_table("content")
    op.execute(
        """
        DO $$
        BEGIN
            IF EXISTS (
                SELECT 1
                FROM pg_type
                WHERE typname = 'topic_behavior_enum'
            ) THEN
                DROP TYPE topic_behavior_enum;
            END IF;
        END$$;
        """
    )
