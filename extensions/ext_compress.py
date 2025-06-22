from boss_app import BossApp
from configs import boss_config


def is_enabled() -> bool:
    return boss_config.API_COMPRESSION_ENABLED


def init_app(app: BossApp):
    from flask_compress import Compress  # type: ignore

    compress = Compress()
    compress.init_app(app)
