"""make legacy product columns nullable

Revision ID: f2b1c3d4e5f6
Revises: e1a2b3c4d5e6
Create Date: 2026-02-22 11:55:00.000000

"""
from alembic import op


# revision identifiers, used by Alembic.
revision = "f2b1c3d4e5f6"
down_revision = "e1a2b3c4d5e6"
branch_labels = None
depends_on = None


def upgrade():
    # Legacy column from previous schema: keep it compatible while code no longer writes it.
    op.execute(
        """
        DO $$
        BEGIN
            IF EXISTS (
                SELECT 1
                FROM information_schema.columns
                WHERE table_name = 'product'
                  AND column_name = 'product_type'
            ) THEN
                ALTER TABLE product ALTER COLUMN product_type DROP NOT NULL;
                ALTER TABLE product ALTER COLUMN product_type SET DEFAULT '';
            END IF;
        END $$;
        """
    )


def downgrade():
    op.execute(
        """
        DO $$
        BEGIN
            IF EXISTS (
                SELECT 1
                FROM information_schema.columns
                WHERE table_name = 'product'
                  AND column_name = 'product_type'
            ) THEN
                UPDATE product
                SET product_type = COALESCE(product_type, '')
                WHERE product_type IS NULL;
                ALTER TABLE product ALTER COLUMN product_type SET NOT NULL;
            END IF;
        END $$;
        """
    )
