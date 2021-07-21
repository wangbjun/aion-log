package main

import (
	"aion/config"
	"aion/model"
	"aion/router"
	"aion/service"
	"aion/zlog"
	"flag"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	var conf string
	flag.StringVar(&conf, "c", "./app.ini", "custom conf file path")
	flag.Parse()

	config.Init(conf)
	zlog.Init()
	model.Init()

	go startDaemon()

	gin.SetMode(getMode())
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
}

func getMode() string {
	debug := config.GetAPP("DEBUG").String()
	if debug == "true" {
		return gin.DebugMode
	}
	return gin.ReleaseMode
}

func startDaemon() {
	go service.RankService{}.Start()
	go service.ClassfiyService{}.Start()
}
