"""add webhook event table

Revision ID: aa7d9c2b11ef
Revises: f2b1c3d4e5f6
Create Date: 2026-02-22 12:20:00.000000

"""
from alembic import op
import sqlalchemy as sa


# revision identifiers, used by Alembic.
revision = "aa7d9c2b11ef"
down_revision = "7c174f245ceb"
branch_labels = None
depends_on = None


def upgrade():
    op.create_table(
        "webhook_event",
        sa.Column("id", sa.Uuid(), nullable=False),
        sa.Column("provider", sa.String(length=100), nullable=False),
        sa.Column("url", sa.String(length=512), nullable=False),
        sa.Column("payload", sa.JSON(), server_default=sa.text("'{}'::json"), nullable=False),
        sa.Column("headers", sa.JSON(), server_default=sa.text("'{}'::json"), nullable=False),
        sa.Column("created_at", sa.DateTime(timezone=True), nullable=True),
        sa.PrimaryKeyConstraint("id"),
    )
    op.create_index(op.f("ix_webhook_event_provider"), "webhook_event", ["provider"], unique=False)


def downgrade():
    op.drop_index(op.f("ix_webhook_event_provider"), table_name="webhook_event")
    op.drop_table("webhook_event")
