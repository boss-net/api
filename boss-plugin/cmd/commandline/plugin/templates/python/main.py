from boss_plugin import Plugin, BossPluginEnv

plugin = Plugin(BossPluginEnv(MAX_REQUEST_TIMEOUT=120))

if __name__ == '__main__':
    plugin.run()
