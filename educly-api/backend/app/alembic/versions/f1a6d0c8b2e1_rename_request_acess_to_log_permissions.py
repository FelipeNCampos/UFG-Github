"""rename request_acess to log_permissions

Revision ID: f1a6d0c8b2e1
Revises: e7f4a2c1b9d0
Create Date: 2026-02-22 20:05:00.000000

"""
from alembic import op
import sqlalchemy as sa


# revision identifiers, used by Alembic.
revision = "f1a6d0c8b2e1"
down_revision = "e7f4a2c1b9d0"
branch_labels = None
depends_on = None


def upgrade():
    op.rename_table("request_acess", "log_permissions")
    op.alter_column("log_permissions", "requester_name", new_column_name="action")
    op.alter_column("log_permissions", "requester_email", new_column_name="email")
    op.add_column("log_permissions", sa.Column("granted_plan", sa.String(length=255), nullable=True))
    op.execute("UPDATE log_permissions SET action = 'assign_access' WHERE action IS NULL OR action = ''")
    op.execute("UPDATE log_permissions SET email = 'unknown' WHERE email IS NULL OR email = ''")
    op.alter_column("log_permissions", "action", nullable=False)
    op.alter_column("log_permissions", "email", nullable=False)


def downgrade():
    op.drop_column("log_permissions", "granted_plan")
    op.alter_column("log_permissions", "email", new_column_name="requester_email")
    op.alter_column("log_permissions", "action", new_column_name="requester_name")
    op.rename_table("log_permissions", "request_acess")
