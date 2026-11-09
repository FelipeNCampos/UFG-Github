"""switch content primary key to integer

Revision ID: c9a2f7e4b1d8
Revises: b1d3f6a8c2e4
Create Date: 2026-02-23 17:20:00.000000

"""

from alembic import op
import sqlalchemy as sa


# revision identifiers, used by Alembic.
revision = "c9a2f7e4b1d8"
down_revision = "b1d3f6a8c2e4"
branch_labels = None
depends_on = None


def _uuid_from_int_sql(column_name: str) -> str:
    return (
        "(substr(lpad(to_hex({0}), 32, '0'), 1, 8) || '-' || "
        "substr(lpad(to_hex({0}), 32, '0'), 9, 4) || '-' || "
        "substr(lpad(to_hex({0}), 32, '0'), 13, 4) || '-' || "
        "substr(lpad(to_hex({0}), 32, '0'), 17, 4) || '-' || "
        "substr(lpad(to_hex({0}), 32, '0'), 21, 12))::uuid"
    ).format(column_name)


def upgrade():
    op.add_column("content", sa.Column("id_new", sa.Integer(), nullable=True))
    op.execute(sa.text("CREATE SEQUENCE IF NOT EXISTS content_id_new_seq"))
    op.execute(
        sa.text(
            "ALTER TABLE content "
            "ALTER COLUMN id_new SET DEFAULT nextval('content_id_new_seq')"
        )
    )
    op.execute(sa.text("UPDATE content SET id_new = nextval('content_id_new_seq') WHERE id_new IS NULL"))
    op.alter_column("content", "id_new", nullable=False)

    op.add_column("topic", sa.Column("content_id_new", sa.Integer(), nullable=True))
    op.execute(
        sa.text(
            "UPDATE topic "
            "SET content_id_new = content.id_new "
            "FROM content "
            "WHERE topic.content_id = content.id"
        )
    )
    op.alter_column("topic", "content_id_new", nullable=False)

    op.drop_constraint("topic_content_id_fkey", "topic", type_="foreignkey")
    op.drop_index(op.f("ix_topic_content_id"), table_name="topic")
    op.drop_constraint("content_pkey", "content", type_="primary")

    op.drop_column("topic", "content_id")
    op.drop_column("content", "id")

    op.alter_column("content", "id_new", new_column_name="id")
    op.alter_column("topic", "content_id_new", new_column_name="content_id")

    op.create_primary_key("content_pkey", "content", ["id"])
    op.create_foreign_key("topic_content_id_fkey", "topic", "content", ["content_id"], ["id"])
    op.create_index(op.f("ix_topic_content_id"), "topic", ["content_id"], unique=False)

    op.execute(
        sa.text(
            "SELECT setval('content_id_new_seq', "
            "COALESCE((SELECT MAX(id) FROM content), 1), true)"
        )
    )
    op.execute(
        sa.text(
            "ALTER TABLE content "
            "ALTER COLUMN id SET DEFAULT nextval('content_id_new_seq')"
        )
    )
    op.execute(sa.text("ALTER SEQUENCE content_id_new_seq OWNED BY content.id"))


def downgrade():
    op.add_column("content", sa.Column("id_uuid", sa.Uuid(), nullable=True))
    op.execute(sa.text(f"UPDATE content SET id_uuid = {_uuid_from_int_sql('id')}"))
    op.alter_column("content", "id_uuid", nullable=False)

    op.add_column("topic", sa.Column("content_id_uuid", sa.Uuid(), nullable=True))
    op.execute(
        sa.text(
            "UPDATE topic "
            "SET content_id_uuid = content.id_uuid "
            "FROM content "
            "WHERE topic.content_id = content.id"
        )
    )
    op.alter_column("topic", "content_id_uuid", nullable=False)

    op.drop_constraint("topic_content_id_fkey", "topic", type_="foreignkey")
    op.drop_index(op.f("ix_topic_content_id"), table_name="topic")
    op.drop_constraint("content_pkey", "content", type_="primary")

    op.drop_column("topic", "content_id")
    op.drop_column("content", "id")

    op.alter_column("content", "id_uuid", new_column_name="id")
    op.alter_column("topic", "content_id_uuid", new_column_name="content_id")

    op.create_primary_key("content_pkey", "content", ["id"])
    op.create_foreign_key("topic_content_id_fkey", "topic", "content", ["content_id"], ["id"])
    op.create_index(op.f("ix_topic_content_id"), "topic", ["content_id"], unique=False)

    op.execute(sa.text("DROP SEQUENCE IF EXISTS content_id_new_seq"))
