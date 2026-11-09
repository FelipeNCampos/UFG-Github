"""add product_color to product

Revision ID: e7f4a2c1b9d0
Revises: d4a1c6b2e9f3
Create Date: 2026-02-22 19:20:00.000000

"""

from alembic import op
import sqlalchemy as sa


# revision identifiers, used by Alembic.
revision = "e7f4a2c1b9d0"
down_revision = "d4a1c6b2e9f3"
branch_labels = None
depends_on = None


def upgrade():
    op.add_column("product", sa.Column("product_color", sa.String(length=20), nullable=True))
    op.execute(
        """
        UPDATE product
        SET product_color = CASE product_id
            WHEN 'BASE' THEN 'green'
            WHEN 'FREELANCER' THEN 'purple'
            WHEN 'HUB_IA' THEN 'blue'
            WHEN 'COMBO_FREELANCER_AI' THEN 'blue'
            WHEN 'FREE' THEN 'gray'
            ELSE 'green'
        END
        WHERE product_color IS NULL
        """
    )
    op.alter_column("product", "product_color", nullable=False, server_default="green")


def downgrade():
    op.drop_column("product", "product_color")
