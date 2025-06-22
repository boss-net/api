from boss_app import BossApp
from configs import boss_config


def init_app(app: BossApp):
    app.secret_key = boss_config.SECRET_KEY
