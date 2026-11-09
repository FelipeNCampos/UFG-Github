from sqlmodel import Session, select

from app.models import Product


def test_default_products_seeded(db: Session) -> None:
    products = db.exec(select(Product)).all()
    by_product_id = {product.product_id: product for product in products}

    assert "BASE" in by_product_id
    assert "FREELANCER" in by_product_id
    assert "HUB_IA" in by_product_id
    assert by_product_id["BASE"].product_color == "green"
    assert by_product_id["FREELANCER"].product_color == "purple"
    assert by_product_id["HUB_IA"].product_color == "blue"

    base_codes = set(by_product_id["BASE"].product_codes)
    freelancer_codes = set(by_product_id["FREELANCER"].product_codes)
    hub_codes = set(by_product_id["HUB_IA"].product_codes)

    assert "01KF3ZG2F9N10HA5FQDQ0FDXSB" in base_codes
    assert "product_01KEZ83ACWTYCKQ3EB24MMPFZM" in base_codes
    assert "7018201" in base_codes
    assert "7031196" in freelancer_codes
    assert "01KG88G787M4W70J3GGAVS6RM1" in hub_codes
