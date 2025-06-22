package model_entities

import "github.com/boss-net/api/boss-plugin/pkg/entities/plugin_entities"

type GetModelSchemasResponse struct {
	ModelSchema *plugin_entities.ModelDeclaration `json:"model_schema" validate:"omitempty"`
}
