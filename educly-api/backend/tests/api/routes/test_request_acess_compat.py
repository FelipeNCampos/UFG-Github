from app.api.routes import request_acess


def test_request_acess_module_exports_log_permissions_router() -> None:
    assert request_acess.router.prefix == "/log-permissions"
