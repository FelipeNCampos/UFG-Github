from fastapi import APIRouter

from app.api.routes import (
    content,
    login,
    log_permissions,
    perfils,
    private,
    products,
    users,
    utils,
    webhooks,
)
from app.core.config import settings

api_router = APIRouter()
api_router.include_router(login.router)
api_router.include_router(users.router)
api_router.include_router(utils.router)
api_router.include_router(content.router)
api_router.include_router(products.router)
api_router.include_router(perfils.router)
api_router.include_router(webhooks.router)
api_router.include_router(log_permissions.router)


if settings.ENVIRONMENT == "local":
    api_router.include_router(private.router)
