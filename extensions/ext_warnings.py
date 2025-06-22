from boss_app import BossApp


def init_app(app: BossApp):
    import warnings

    warnings.simplefilter("ignore", ResourceWarning)
