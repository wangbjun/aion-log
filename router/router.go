package router

import (
	. "aion/controller"
	"aion/middleware"
	"github.com/gin-gonic/gin"
)

func Route(Router *gin.Engine) {
	Router.GET("/", BaseController.Index)

	api := Router.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			v1.GET("/logs", BattleController.GetAll).
				GET("/ranks", BattleController.GetRank).
				GET("/players", BattleController.GetPlayers)

			v1.Group("/user").
				POST("/register", UserController.Register). //用户注册
				POST("/login", UserController.Login).       //用户登录
				GET("/current", UserController.Current)     //用户登录

			v1.Group("/user").Use(middleware.Auth()).
				POST("/logout", UserController.Logout) //用户退出登录
		}
	}
}
