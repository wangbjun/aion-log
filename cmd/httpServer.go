package cmd

import (
	"aion/config"
	"aion/router"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"log"
)

func init() {
	rootCmd.AddCommand(httpServerCmd)
}

var httpServerCmd = &cobra.Command{
	Use:   "httpServer",
	Short: "Start A Http Server",
	Run: func(cmd *cobra.Command, args []string) {
		if config.GetAPP("DEBUG").String() != "true" {
			gin.SetMode(gin.ReleaseMode)
		}
		engine := gin.New()
		engine.Use(gin.Recovery())
		// 加载路由
		router.Route(engine)

		// 启动服务器
		log.Println("server started success")
		err := engine.Run(":" + config.GetAPP("PORT").String())
		if err != nil {
			log.Fatalf("server start failed, error: %s", err.Error())
		}
	},
}
