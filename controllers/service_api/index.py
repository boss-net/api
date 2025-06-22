from flask_restful import Resource

from configs import boss_config
from controllers.service_api import api


class IndexApi(Resource):
    def get(self):
        return {
            "welcome": "Boss OpenAPI",
            "api_version": "v1",
            "server_version": boss_config.CURRENT_VERSION,
        }


api.add_resource(IndexApi, "/")
