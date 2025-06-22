from boss_app import BossApp


def init_app(app: BossApp):
    from events import event_handlers  # noqa: F401
