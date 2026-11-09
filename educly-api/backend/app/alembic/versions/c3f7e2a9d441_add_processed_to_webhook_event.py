"""add processed to webhook_event

Revision ID: c3f7e2a9d441
Revises: b9f6d2a1c3e4
Create Date: 2026-02-22 14:10:00.000000

"""
from alembic import op
import sqlalchemy as sa


# revision identifiers, used by Alembic.
revision = "c3f7e2a9d441"
down_revision = "b9f6d2a1c3e4"
branch_labels = None
depends_on = None


def upgrade():
    op.add_column(
        "webhook_event",
        sa.Column(
            "processed",
            sa.Boolean(),
            nullable=False,
            server_default=sa.text("false"),
        ),
    )
    op.execute("UPDATE webhook_event SET processed = false WHERE processed IS NULL")


def downgrade():
    op.drop_column("webhook_event", "processed")
