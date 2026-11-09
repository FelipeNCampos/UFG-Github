import uuid
from typing import Any

from fastapi import APIRouter, HTTPException
from sqlmodel import col, func, select

from app.api.deps import CurrentUser, SessionDep
from app.models import (
    Message,
    Product,
    ProductCreate,
    ProductPublic,
    ProductsPublic,
    ProductUpdate,
)

router = APIRouter(prefix="/products", tags=["products"])

PRODUCT_REFERENCE_VALUES = {
    "FREE": 1,
    "BASE": 2,
    "FREELANCER": 4,
    "HUB_IA": 8,
    "COMBO_FREELANCER_AI": 12,
}

PRODUCT_REFERENCE_COLORS = {
    "FREE": "gray",
    "BASE": "green",
    "FREELANCER": "purple",
    "HUB_IA": "blue",
    "COMBO_FREELANCER_AI": "blue",
}


def get_reference_value_or_400(product_id: str) -> int:
    value = PRODUCT_REFERENCE_VALUES.get(product_id)
    if value is None:
        raise HTTPException(status_code=400, detail="Invalid product_id reference")
    return value


def get_reference_color_or_400(product_id: str) -> str:
    color = PRODUCT_REFERENCE_COLORS.get(product_id)
    if color is None:
        raise HTTPException(status_code=400, detail="Invalid product_id reference")
    return color


@router.get("/", response_model=ProductsPublic)
def read_products(
    session: SessionDep, current_user: CurrentUser, skip: int = 0, limit: int = 100
) -> Any:
    """
    Retrieve products. Only superusers can access.
    """
    if not current_user.is_superuser:
        raise HTTPException(status_code=403, detail="Not enough permissions")

    count_statement = select(func.count()).select_from(Product)
    count = session.exec(count_statement).one()
    statement = (
        select(Product)
        .order_by(col(Product.created_at).desc())
        .offset(skip)
        .limit(limit)
    )
    products = session.exec(statement).all()

    return ProductsPublic(data=list(products), count=count)


@router.get("/{id}", response_model=ProductPublic)
def read_product(session: SessionDep, current_user: CurrentUser, id: uuid.UUID) -> Any:
    """
    Get product by ID.
    """
    if not current_user.is_superuser:
        raise HTTPException(status_code=403, detail="Not enough permissions")

    product = session.get(Product, id)
    if not product:
        raise HTTPException(status_code=404, detail="Product not found")
    return product


@router.post("/", response_model=ProductPublic)
def create_product(
    *, session: SessionDep, current_user: CurrentUser, product_in: ProductCreate
) -> Any:
    """
    Create new product. Only superusers can create products.
    """
    if not current_user.is_superuser:
        raise HTTPException(status_code=403, detail="Not enough permissions")

    # Check if product_id already exists
    statement = select(Product).where(Product.product_id == product_in.product_id)
    existing_product = session.exec(statement).first()
    if existing_product:
        raise HTTPException(
            status_code=400, detail="Product with this product_id already exists"
        )

    payload = product_in.model_dump()
    payload["valor"] = get_reference_value_or_400(product_in.product_id)
    payload["product_color"] = get_reference_color_or_400(product_in.product_id)
    product = Product.model_validate(payload)
    session.add(product)
    session.commit()
    session.refresh(product)
    return product


@router.put("/{id}", response_model=ProductPublic)
def update_product(
    *,
    session: SessionDep,
    current_user: CurrentUser,
    id: uuid.UUID,
    product_in: ProductUpdate,
) -> Any:
    """
    Update a product. Only superusers can update products.
    """
    if not current_user.is_superuser:
        raise HTTPException(status_code=403, detail="Not enough permissions")

    product = session.get(Product, id)
    if not product:
        raise HTTPException(status_code=404, detail="Product not found")

    # Check if product_id is being changed to one that already exists
    if product_in.product_id and product_in.product_id != product.product_id:
        statement = select(Product).where(Product.product_id == product_in.product_id)
        existing_product = session.exec(statement).first()
        if existing_product:
            raise HTTPException(
                status_code=400, detail="Product with this product_id already exists"
            )

    update_dict = product_in.model_dump(exclude_unset=True)
    if "product_id" in update_dict:
        update_dict["valor"] = get_reference_value_or_400(update_dict["product_id"])
        update_dict["product_color"] = get_reference_color_or_400(update_dict["product_id"])
    else:
        update_dict["valor"] = get_reference_value_or_400(product.product_id)
        if "product_color" not in update_dict:
            update_dict["product_color"] = get_reference_color_or_400(product.product_id)
    product.sqlmodel_update(update_dict)
    session.add(product)
    session.commit()
    session.refresh(product)
    return product


@router.delete("/{id}")
def delete_product(
    session: SessionDep, current_user: CurrentUser, id: uuid.UUID
) -> Message:
    """
    Delete a product. Only superusers can delete products.
    """
    if not current_user.is_superuser:
        raise HTTPException(status_code=403, detail="Not enough permissions")

    product = session.get(Product, id)
    if not product:
        raise HTTPException(status_code=404, detail="Product not found")
    session.delete(product)
    session.commit()
    return Message(message="Product deleted successfully")
