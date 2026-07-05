package main

import (
	"aion/model"
	"aion/router"
	"aion/service"
	"log"
	"net"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	model.Init()

	engine := gin.New()
	engine.Use(gin.Recovery())
	router.Route(engine)

	cacheService := service.NewCacheService()
	if err := cacheService.Load(); err != nil {
		log.Fatalf("load cache failed: %s", err)
	}
	log.Println("load cache success")

	addrToListen := os.Getenv("AION_ADDR")
	if addrToListen == "" {
		addrToListen = "127.0.0.1:18080"
		if os.Getenv("AION_PORT_FILE") != "" {
			addrToListen = "127.0.0.1:0"
		}
	}

	listener, err := net.Listen("tcp", addrToListen)
	if err != nil {
		log.Fatalf("server listen failed: %s", err)
	}

	addr := listener.Addr().String()
	if portFile := os.Getenv("AION_PORT_FILE"); portFile != "" {
		if err := os.WriteFile(portFile, []byte(addr), 0600); err != nil {
			log.Fatalf("write server address failed: %s", err)
		}
	}

	log.Printf("server listening on http://%s", addr)
	if err := engine.RunListener(listener); err != nil {
		log.Fatalf("server start failed: %s", err)
	}
}
