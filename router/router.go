package router

import (
	. "aion/controller"
	"aion/middleware"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func Route(Router *gin.Engine) {
	Router.Use(middleware.Request())
	api := Router.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			v1.GET("/logs", BattleController.GetAll).
				GET("/ranks", BattleController.GetRank).
				GET("/players", BattleController.GetPlayers).
				GET("/timeline", BattleController.GetTimeline).
				GET("/classTop", BattleController.GetClassTop).
				GET("/logs/import/status", BattleController.GetImportStatus)
			v1.POST("/logs/import", BattleController.ImportLog)
		}
	}
	registerFrontend(Router)
}

func registerFrontend(Router *gin.Engine) {
	distDir := filepath.Join("frontend", "dist")
	indexFile := filepath.Join(distDir, "index.html")

	Router.GET("/", func(ctx *gin.Context) {
		if _, err := os.Stat(indexFile); err == nil {
			ctx.File(indexFile)
			return
		}
		BaseController.Index(ctx)
	})

	Router.NoRoute(func(ctx *gin.Context) {
		if strings.HasPrefix(ctx.Request.URL.Path, "/api/") {
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}

		cleanPath := filepath.Clean(strings.TrimPrefix(ctx.Request.URL.Path, "/"))
		if cleanPath != "." {
			filePath := filepath.Join(distDir, cleanPath)
			if info, err := os.Stat(filePath); err == nil && !info.IsDir() {
				ctx.File(filePath)
				return
			}
		}

		if _, err := os.Stat(indexFile); err == nil {
			ctx.File(indexFile)
			return
		}
		ctx.AbortWithStatus(http.StatusNotFound)
	})
}
