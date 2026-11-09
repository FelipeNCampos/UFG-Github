from types import SimpleNamespace

import app.initial_data as initial_data


def test_main_calls_init(monkeypatch) -> None:
    called = {"init": False}

    def fake_init() -> None:
        called["init"] = True

    monkeypatch.setattr(initial_data, "init", fake_init)
    initial_data.main()

    assert called["init"] is True


def test_init_uses_session_and_calls_init_db(monkeypatch) -> None:
    called = {"session_opened": False, "init_db": False}

    class DummySessionCtx:
        def __enter__(self):
            called["session_opened"] = True
            return SimpleNamespace()

        def __exit__(self, exc_type, exc, tb):
            return None

    def fake_session(_engine):
        return DummySessionCtx()

    def fake_init_db(_session):
        called["init_db"] = True

    monkeypatch.setattr(initial_data, "Session", fake_session)
    monkeypatch.setattr(initial_data, "init_db", fake_init_db)

    initial_data.init()

    assert called["session_opened"] is True
    assert called["init_db"] is True
