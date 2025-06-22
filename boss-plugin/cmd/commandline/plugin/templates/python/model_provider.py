import logging
from collections.abc import Mapping

from boss_plugin import ModelProvider
from boss_plugin.entities.model import ModelType
from boss_plugin.errors.model import CredentialsValidateFailedError

logger = logging.getLogger(__name__)


class {{ .PluginName | SnakeToCamel }}ModelProvider(ModelProvider):
    def validate_provider_credentials(self, credentials: Mapping) -> None:
        """
        Validate provider credentials
        if validate failed, raise exception

        :param credentials: provider credentials, credentials form defined in `provider_credential_schema`.
        """
        try:
            pass
        except CredentialsValidateFailedError as ex:
            raise ex
        except Exception as ex:
            logger.exception(
                f"{self.get_provider_schema().provider} credentials validate failed"
            )
            raise ex
