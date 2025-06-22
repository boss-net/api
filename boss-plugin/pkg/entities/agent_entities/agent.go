package agent_entities

import "github.com/boss-net/api/boss-plugin/pkg/entities/tool_entities"

type AgentStrategyResponseChunk struct {
	tool_entities.ToolResponseChunk `json:",inline"`
}
