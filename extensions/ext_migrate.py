from boss_app import BossApp


def init_app(app: BossApp):
    import flask_migrate  # type: ignore

    from extensions.ext_database import db

    flask_migrate.Migrate(app, db)
