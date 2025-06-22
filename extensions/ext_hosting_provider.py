from core.hosting_configuration import HostingConfiguration

hosting_configuration = HostingConfiguration()


from boss_app import BossApp


def init_app(app: BossApp):
    hosting_configuration.init_app(app)
