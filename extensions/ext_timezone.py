import os
import time

from boss_app import BossApp


def init_app(app: BossApp):
    os.environ["TZ"] = "UTC"
    # windows platform not support tzset
    if hasattr(time, "tzset"):
        time.tzset()
