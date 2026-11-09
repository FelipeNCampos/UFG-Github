"""merge heads log_permissions and content

Revision ID: 9e2f6c1a4b7d
Revises: 124824675783, a7c91d3e4f20
Create Date: 2026-02-22 16:20:00.000000

"""

# revision identifiers, used by Alembic.
revision = "9e2f6c1a4b7d"
down_revision = ("124824675783", "a7c91d3e4f20")
branch_labels = None
depends_on = None


def upgrade():
    # Merge revision: no schema changes, only unifies migration history heads.
    pass


def downgrade():
    # Split merged history back into two heads.
    pass
