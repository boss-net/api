from boss_app import BossApp
from configs import boss_config


def init_app(app: BossApp):
    if boss_config.SENTRY_DSN:
        import openai
        import sentry_sdk
        from langfuse import parse_error  # type: ignore
        from sentry_sdk.integrations.celery import CeleryIntegration
        from sentry_sdk.integrations.flask import FlaskIntegration
        from werkzeug.exceptions import HTTPException

        from core.model_runtime.errors.invoke import InvokeRateLimitError

        def before_send(event, hint):
            if "exc_info" in hint:
                exc_type, exc_value, tb = hint["exc_info"]
                if parse_error.defaultErrorResponse in str(exc_value):
                    return None

            return event

        sentry_sdk.init(
            dsn=boss_config.SENTRY_DSN,
            integrations=[FlaskIntegration(), CeleryIntegration()],
            ignore_errors=[
                HTTPException,
                ValueError,
                FileNotFoundError,
                openai.APIStatusError,
                InvokeRateLimitError,
                parse_error.defaultErrorResponse,
            ],
            traces_sample_rate=boss_config.SENTRY_TRACES_SAMPLE_RATE,
            profiles_sample_rate=boss_config.SENTRY_PROFILES_SAMPLE_RATE,
            environment=boss_config.DEPLOY_ENV,
            release=f"boss-{boss_config.CURRENT_VERSION}-{boss_config.COMMIT_SHA}",
            before_send=before_send,
        )
