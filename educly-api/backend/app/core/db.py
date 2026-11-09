from sqlmodel import Session, create_engine, select

from app import crud
from app.core.config import settings
from app.models import Product, User, UserCreate

engine = create_engine(str(settings.SQLALCHEMY_DATABASE_URI))

PRODUCT_CODE_PLAN_MAP = {
    "product_01KEZ83ACWTYCKQ3EB24MMPFZM": "base",
    "product_01KEZ9HH94AJM2NQXQPYPBFHWF": "base",
    "product_01KEZQQZ898SN6M45C1E3WCRPY": "base",
    "7018201": "base",
    "01KFH6TYKXMKQ5AJQJ4PZZ1G43": "base",
    "01KFHKR3GQBV995WFEGPAXYTPQ": "base",
    "01KF3ZG2F9N10HA5FQDQ0FDXSB": "base",
    "01KF4029KCA1QPQXXCMGN69DV9": "freelancer",
    "7031196": "freelancer",
    "01KG88G787M4W70J3GGAVS6RM1": "ai_hub",
    "7179838": "ai_hub",
    "7181563": "combo_freelancer_ai",
}

DEFAULT_PRODUCTS = [
    {
        "product_id": "BASE",
        "name": "Base",
        "valor": 2,
        "product_color": "green",
        "description": "Plano Base",
        "plan_key": "base",
    },
    {
        "product_id": "FREELANCER",
        "name": "Freelancer",
        "valor": 4,
        "product_color": "purple",
        "description": "Plano Freelancer",
        "plan_key": "freelancer",
    },
    {
        "product_id": "HUB_IA",
        "name": "Hub IA",
        "valor": 8,
        "product_color": "blue",
        "description": "Plano Hub IA",
        "plan_key": "ai_hub",
    },
    {
        "product_id": "COMBO_FREELANCER_AI",
        "name": "Combo Freelancer + Hub IA",
        "valor": 12,
        "product_color": "blue",
        "description": "Plano combinado Freelancer + Hub IA",
        "plan_key": "combo_freelancer_ai",
    },
]


def _codes_for_plan(plan_key: str) -> list[str]:
    return [code for code, mapped_plan in PRODUCT_CODE_PLAN_MAP.items() if mapped_plan == plan_key]


def _seed_default_products(session: Session) -> None:
    for product_data in DEFAULT_PRODUCTS:
        product = session.exec(
            select(Product).where(Product.product_id == product_data["product_id"])
        ).first()
        seed_codes = _codes_for_plan(product_data["plan_key"])

        if not product:
            product = Product(
                name=product_data["name"],
                product_id=product_data["product_id"],
                description=product_data["description"],
                product_color=product_data["product_color"],
                product_codes=seed_codes,
                valor=product_data["valor"],
            )
            session.add(product)
            continue

        existing_codes = product.product_codes or []
        merged_codes = sorted(set(existing_codes + seed_codes))
        changed = False
        if product.name != product_data["name"]:
            product.name = product_data["name"]
            changed = True
        if product.description != product_data["description"]:
            product.description = product_data["description"]
            changed = True
        if product.product_color != product_data["product_color"]:
            product.product_color = product_data["product_color"]
            changed = True
        if product.valor != product_data["valor"]:
            product.valor = product_data["valor"]
            changed = True
        if merged_codes != existing_codes:
            product.product_codes = merged_codes
            changed = True
        if changed:
            session.add(product)

    session.commit()


# make sure all SQLModel models are imported (app.models) before initializing DB
# otherwise, SQLModel might fail to initialize relationships properly
# for more details: https://github.com/fastapi/full-stack-fastapi-template/issues/28


def init_db(session: Session) -> None:
    # Tables should be created with Alembic migrations
    # But if you don't want to use migrations, create
    # the tables un-commenting the next lines
    # from sqlmodel import SQLModel

    # This works because the models are already imported and registered from app.models
    # SQLModel.metadata.create_all(engine)

    user = session.exec(
        select(User).where(User.email == settings.FIRST_SUPERUSER)
    ).first()
    if not user:
        user_in = UserCreate(
            email=settings.FIRST_SUPERUSER,
            password=settings.FIRST_SUPERUSER_PASSWORD,
            is_superuser=True,
            full_name="Admin",
        )
        user = crud.create_user(session=session, user_create=user_in)
    else:
        # Ensure bootstrap admin keeps superuser permissions even after test runs/manual edits.
        needs_update = not user.is_superuser or not user.is_active
        if needs_update:
            user.is_superuser = True
            user.is_active = True
            session.add(user)
            session.commit()

    _seed_default_products(session)
