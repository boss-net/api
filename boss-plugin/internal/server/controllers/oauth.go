package controllers

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/boss-net/api/boss-plugin/internal/service"
	"github.com/boss-net/api/boss-plugin/internal/types/app"
	"github.com/boss-net/api/boss-plugin/pkg/entities/plugin_entities"
	"github.com/boss-net/api/boss-plugin/pkg/entities/requests"
)

func OAuthGetAuthorizationURL(config *app.Config) gin.HandlerFunc {
	type request = plugin_entities.InvokePluginRequest[requests.RequestOAuthGetAuthorizationURL]

	return func(c *gin.Context) {
		BindPluginDispatchRequest(
			c,
			func(ipr request) {
				service.OAuthGetAuthorizationURL(
					&ipr,
					c,
					time.Duration(config.PluginMaxExecutionTimeout)*time.Second,
				)
			},
		)
	}
}

func OAuthGetCredentials(config *app.Config) gin.HandlerFunc {
	type request = plugin_entities.InvokePluginRequest[requests.RequestOAuthGetCredentials]

	return func(c *gin.Context) {
		BindPluginDispatchRequest(c, func(ipr request) {
			service.OAuthGetCredentials(
				&ipr,
				c,
				time.Duration(config.PluginMaxExecutionTimeout)*time.Second,
			)
		})
	}
}
