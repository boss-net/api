package serverless_runtime

import (
	"net/http"

	"github.com/boss-net/api/boss-plugin/internal/core/plugin_manager/basic_runtime"
	"github.com/boss-net/api/boss-plugin/internal/utils/mapping"
	"github.com/boss-net/api/boss-plugin/pkg/entities"
	"github.com/boss-net/api/boss-plugin/pkg/entities/plugin_entities"
)

type ServerlessPluginRuntime struct {
	basic_runtime.BasicChecksum
	plugin_entities.PluginRuntime

	// access url for the lambda function
	LambdaURL  string
	LambdaName string

	// listeners mapping session id to the listener
	listeners mapping.Map[string, *entities.Broadcast[plugin_entities.SessionMessage]]

	client *http.Client

	PluginMaxExecutionTimeout int // in seconds
}
