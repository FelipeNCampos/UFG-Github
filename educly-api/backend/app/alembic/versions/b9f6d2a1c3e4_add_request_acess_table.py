"""add request acess table

Revision ID: b9f6d2a1c3e4
Revises: 6403c515b6d0
Create Date: 2026-02-22 13:05:00.000000

"""
from alembic import op
import sqlalchemy as sa
import sqlmodel.sql.sqltypes


# revision identifiers, used by Alembic.
revision = "b9f6d2a1c3e4"
down_revision = "6403c515b6d0"
branch_labels = None
depends_on = None


def upgrade():
    op.create_table(
        "request_acess",
        sa.Column("requester_name", sqlmodel.sql.sqltypes.AutoString(length=255), nullable=False),
        sa.Column("requester_email", sqlmodel.sql.sqltypes.AutoString(length=255), nullable=False),
        sa.Column("reason", sqlmodel.sql.sqltypes.AutoString(length=500), nullable=True),
        sa.Column("status", sqlmodel.sql.sqltypes.AutoString(length=30), nullable=False, server_default="pending"),
        sa.Column("id", sa.Uuid(), nullable=False),
        sa.Column("created_at", sa.DateTime(timezone=True), nullable=True),
        sa.PrimaryKeyConstraint("id"),
    )


def downgrade():
    op.drop_table("request_acess")
