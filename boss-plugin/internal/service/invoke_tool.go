package service

import (
	"github.com/gin-gonic/gin"
	"github.com/boss-net/api/boss-plugin/internal/core/plugin_daemon"
	"github.com/boss-net/api/boss-plugin/internal/core/plugin_daemon/access_types"
	"github.com/boss-net/api/boss-plugin/internal/core/session_manager"
	"github.com/boss-net/api/boss-plugin/internal/utils/stream"
	"github.com/boss-net/api/boss-plugin/pkg/entities/plugin_entities"
	"github.com/boss-net/api/boss-plugin/pkg/entities/requests"
	"github.com/boss-net/api/boss-plugin/pkg/entities/tool_entities"
)

func InvokeTool(
	r *plugin_entities.InvokePluginRequest[requests.RequestInvokeTool],
	ctx *gin.Context,
	max_timeout_seconds int,
) {
	baseSSEWithSession(
		func(session *session_manager.Session) (*stream.Stream[tool_entities.ToolResponseChunk], error) {
			return plugin_daemon.InvokeTool(session, &r.Data)
		},
		access_types.PLUGIN_ACCESS_TYPE_TOOL,
		access_types.PLUGIN_ACCESS_ACTION_INVOKE_TOOL,
		r,
		ctx,
		max_timeout_seconds,
	)
}
