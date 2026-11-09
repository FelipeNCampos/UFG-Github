"""fix

Revision ID: 8d611ba50071
Revises: 3733336cd56c
Create Date: 2026-02-21 16:23:21.195832

"""
from alembic import op
import sqlalchemy as sa
import sqlmodel.sql.sqltypes


# revision identifiers, used by Alembic.
revision = '8d611ba50071'
down_revision = '3733336cd56c'
branch_labels = None
depends_on = None


def upgrade():
    bind = op.get_bind()
    inspector = sa.inspect(bind)
    table_names = set(inspector.get_table_names())

    if 'perfil' not in table_names:
        op.create_table('perfil',
        sa.Column('profile_bitmap', sa.Integer(), nullable=False),
        sa.Column('personalized_trail', sa.Boolean(), nullable=False),
        sa.Column('avatar_url', sqlmodel.sql.sqltypes.AutoString(length=255), nullable=True),
        sa.Column('cover_url', sqlmodel.sql.sqltypes.AutoString(length=255), nullable=True),
        sa.Column('bio', sqlmodel.sql.sqltypes.AutoString(length=500), nullable=True),
        sa.Column('id', sa.Uuid(), nullable=False),
        sa.Column('user_id', sa.Uuid(), nullable=False),
        sa.ForeignKeyConstraint(['user_id'], ['user.id'], ),
        sa.PrimaryKeyConstraint('id')
        )

    if 'perfil' in set(inspector.get_table_names()):
        index_names = {idx['name'] for idx in inspector.get_indexes('perfil')}
        index_name = op.f('ix_perfil_user_id')
        if index_name not in index_names:
            op.create_index(index_name, 'perfil', ['user_id'], unique=False)


def downgrade():
    bind = op.get_bind()
    inspector = sa.inspect(bind)
    table_names = set(inspector.get_table_names())

    if 'perfil' in table_names:
        index_names = {idx['name'] for idx in inspector.get_indexes('perfil')}
        index_name = op.f('ix_perfil_user_id')
        if index_name in index_names:
            op.drop_index(index_name, table_name='perfil')
        op.drop_table('perfil')
