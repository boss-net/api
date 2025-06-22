package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/boss-net/api/boss-plugin/internal/manifest"
	"github.com/boss-net/api/boss-plugin/internal/types/app"
	"github.com/boss-net/api/boss-plugin/internal/utils/routine"
)

func HealthCheck(app *app.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":      "ok",
			"pool_status": routine.FetchRoutineStatus(),
			"version":     manifest.VersionX,
			"build_time":  manifest.BuildTimeX,
			"platform":    app.Platform,
		})
	}
}
