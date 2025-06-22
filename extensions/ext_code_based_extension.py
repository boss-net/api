from boss_app import BossApp
from core.extension.extension import Extension


def init_app(app: BossApp):
    code_based_extension.init()


code_based_extension = Extension()
