package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/boss-net/api/boss-plugin/internal/service"
	"github.com/boss-net/api/boss-plugin/pkg/entities/requests"
)

func GetRemoteDebuggingKey(c *gin.Context) {
	BindRequest(
		c, func(request requests.RequestGetRemoteDebuggingKey) {
			c.JSON(200, service.GetRemoteDebuggingKey(request.TenantID))
		},
	)
}
