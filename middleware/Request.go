package middleware

import (
	"aion/util"
	"aion/zlog"
	"github.com/gin-gonic/gin"
	"time"
)

/**
 * 记录请求日志，加入traceId
 */
func Request() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set("traceId", util.GetUuid())
		ctx.Set("startTime", time.Now())
		ctx.Set("parentId", ctx.GetHeader("X-Ca-TraceId"))

		zlog.WithContext(ctx).Info("Before_Request")
		ctx.Next()
		zlog.WithContext(ctx).Info("After_Request")
	}
}
