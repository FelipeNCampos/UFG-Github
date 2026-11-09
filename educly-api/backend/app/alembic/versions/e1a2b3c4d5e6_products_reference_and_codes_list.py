"""products reference and codes list

Revision ID: e1a2b3c4d5e6
Revises: 8d611ba50071
Create Date: 2026-02-22 11:20:00.000000

"""
from alembic import op


# revision identifiers, used by Alembic.
revision = "e1a2b3c4d5e6"
down_revision = "8d611ba50071"
branch_labels = None
depends_on = None


def upgrade():
    # Ensure product_id exists and can be used as reference key.
    op.execute("ALTER TABLE product ADD COLUMN IF NOT EXISTS product_id VARCHAR(255)")
    op.execute("UPDATE product SET product_id = id::text WHERE product_id IS NULL")
    op.execute("ALTER TABLE product ALTER COLUMN product_id SET NOT NULL")
    op.execute(
        "CREATE UNIQUE INDEX IF NOT EXISTS ix_product_product_id ON product (product_id)"
    )

    # Ensure valor exists.
    op.execute("ALTER TABLE product ADD COLUMN IF NOT EXISTS valor INTEGER")

    # If product_codes exists as non-JSON type, preserve it temporarily.
    op.execute(
        """
        DO $$
        BEGIN
            IF EXISTS (
                SELECT 1
                FROM information_schema.columns
                WHERE table_name = 'product'
                  AND column_name = 'product_codes'
                  AND data_type <> 'jsonb'
            ) THEN
                IF NOT EXISTS (
                    SELECT 1
                    FROM information_schema.columns
                    WHERE table_name = 'product'
                      AND column_name = 'product_codes_legacy'
                ) THEN
                    ALTER TABLE product RENAME COLUMN product_codes TO product_codes_legacy;
                END IF;
            END IF;
        END $$;
        """
    )

    # Create canonical JSON list column.
    op.execute("ALTER TABLE product ADD COLUMN IF NOT EXISTS product_codes JSONB")

    # Migrate legacy text/list representations.
    op.execute(
        """
        DO $$
        BEGIN
            IF EXISTS (
                SELECT 1
                FROM information_schema.columns
                WHERE table_name = 'product'
                  AND column_name = 'product_codes_legacy'
            ) THEN
                UPDATE product
                SET product_codes = to_jsonb(
                    array_remove(
                        ARRAY(
                            SELECT NULLIF(trim(value), '')
                            FROM unnest(string_to_array(product_codes_legacy, ',')) AS value
                        ),
                        NULL
                    )
                )
                WHERE product_codes_legacy IS NOT NULL
                  AND (
                    product_codes IS NULL
                    OR product_codes = '[]'::jsonb
                  );
            END IF;
        END $$;
        """
    )

    # Backfill from old product_type if present and product_codes still empty.
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
                SET product_codes = to_jsonb(ARRAY[product_type])
                WHERE product_type IS NOT NULL
                  AND (
                    product_codes IS NULL
                    OR product_codes = '[]'::jsonb
                  );
            END IF;
        END $$;
        """
    )

    # Normalize null and empty default.
    op.execute("UPDATE product SET product_codes = '[]'::jsonb WHERE product_codes IS NULL")
    op.execute("ALTER TABLE product ALTER COLUMN product_codes SET DEFAULT '[]'::jsonb")
    op.execute("ALTER TABLE product ALTER COLUMN product_codes SET NOT NULL")

    # Compute valor from reference product_id.
    op.execute(
        """
        UPDATE product
        SET valor = CASE product_id
            WHEN 'FREE' THEN 1
            WHEN 'BASE' THEN 2
            WHEN 'FREELANCER' THEN 4
            WHEN 'HUB_IA' THEN 8
            ELSE COALESCE(valor, 0)
        END
        """
    )
    op.execute("ALTER TABLE product ALTER COLUMN valor SET DEFAULT 0")
    op.execute("UPDATE product SET valor = 0 WHERE valor IS NULL")
    op.execute("ALTER TABLE product ALTER COLUMN valor SET NOT NULL")


def downgrade():
    # Convert JSON list back to text representation.
    op.execute("ALTER TABLE product ADD COLUMN IF NOT EXISTS product_codes_legacy VARCHAR(255)")
    op.execute(
        """
        UPDATE product
        SET product_codes_legacy = (
            SELECT string_agg(value, ',')
            FROM jsonb_array_elements_text(product_codes) AS value
        )
        WHERE product_codes IS NOT NULL
        """
    )

    op.execute("ALTER TABLE product DROP COLUMN IF EXISTS product_codes")
    op.execute("ALTER TABLE product RENAME COLUMN product_codes_legacy TO product_codes")
    op.execute("ALTER TABLE product DROP COLUMN IF EXISTS valor")
