package middleware

import (
	"aion/util"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Request /**
func Request() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set("traceId", util.GetUuid())
		ctx.Set("startTime", time.Now())
		ctx.Set("parentId", ctx.GetHeader("X-Ca-TraceId"))
		ctx.Next()
		log.Printf("%s %s %s", ctx.ClientIP(), ctx.Request.Method, ctx.Request.URL.Path)
	}
}
