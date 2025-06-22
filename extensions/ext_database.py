from boss_app import BossApp
from models import db


def init_app(app: BossApp):
    db.init_app(app)
