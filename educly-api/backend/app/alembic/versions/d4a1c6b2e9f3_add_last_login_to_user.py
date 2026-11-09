"""add last_login to user

Revision ID: d4a1c6b2e9f3
Revises: c3f7e2a9d441
Create Date: 2026-02-22 14:45:00.000000

"""
from alembic import op
import sqlalchemy as sa


# revision identifiers, used by Alembic.
revision = "d4a1c6b2e9f3"
down_revision = "c3f7e2a9d441"
branch_labels = None
depends_on = None


def upgrade():
    op.add_column("user", sa.Column("last_login", sa.DateTime(timezone=True), nullable=True))


def downgrade():
    op.drop_column("user", "last_login")
