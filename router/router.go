package router

import (
	. "aion/controller"
	"aion/middleware"
	"github.com/gin-gonic/gin"
)

func Route(Router *gin.Engine) {
	Router.Use(middleware.Request())
	Router.GET("/", BaseController.Index)
	api := Router.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			v1.GET("/logs", BattleController.GetAll).
				GET("/ranks", BattleController.GetRank).
				GET("/players", BattleController.GetPlayers).
				GET("/timeline", BattleController.GetTimeline).
				GET("/classTop", BattleController.GetClassTop)
		}
	}
}
