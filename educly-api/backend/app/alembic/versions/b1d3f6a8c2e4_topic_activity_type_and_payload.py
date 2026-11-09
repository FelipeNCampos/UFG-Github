"""topic activity_type and payload

Revision ID: b1d3f6a8c2e4
Revises: 9e2f6c1a4b7d
Create Date: 2026-02-22 17:05:00.000000

"""
from alembic import op
import sqlalchemy as sa
from sqlalchemy.dialects import postgresql


# revision identifiers, used by Alembic.
revision = "b1d3f6a8c2e4"
down_revision = "9e2f6c1a4b7d"
branch_labels = None
depends_on = None


def upgrade():
    op.execute(
        """
        DO $$
        BEGIN
            IF NOT EXISTS (
                SELECT 1
                FROM pg_type
                WHERE typname = 'activity_type_enum'
            ) THEN
                CREATE TYPE activity_type_enum AS ENUM (
                    'MatchWords',
                    'FillBlanks',
                    'Quizz',
                    'PromptTreiner'
                );
            END IF;
        END$$;
        """
    )

    op.add_column(
        "topic",
        sa.Column(
            "activity_type",
            postgresql.ENUM(
                "MatchWords",
                "FillBlanks",
                "Quizz",
                "PromptTreiner",
                name="activity_type_enum",
                create_type=False,
            ),
            nullable=True,
        ),
    )
    op.add_column(
        "topic",
        sa.Column(
            "payload",
            postgresql.JSONB(astext_type=sa.Text()),
            nullable=False,
            server_default=sa.text("'{}'::jsonb"),
        ),
    )
    op.alter_column("topic", "title", existing_type=sa.String(length=255), nullable=True)
    op.alter_column("topic", "body", existing_type=sa.String(), nullable=True)


def downgrade():
    op.execute("UPDATE topic SET title = COALESCE(title, 'Sem titulo') WHERE title IS NULL")
    op.execute("UPDATE topic SET body = COALESCE(body, 'Sem conteudo') WHERE body IS NULL")
    op.alter_column("topic", "title", existing_type=sa.String(length=255), nullable=False)
    op.alter_column("topic", "body", existing_type=sa.String(), nullable=False)
    op.drop_column("topic", "payload")
    op.drop_column("topic", "activity_type")

    op.execute(
        """
        DO $$
        BEGIN
            IF EXISTS (
                SELECT 1
                FROM pg_type
                WHERE typname = 'activity_type_enum'
            ) THEN
                DROP TYPE activity_type_enum;
            END IF;
        END$$;
        """
    )
