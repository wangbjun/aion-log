package zlog

import (
	"aion/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"time"
)

var Logger *zap.Logger

func WithContext(ctx *gin.Context) *zap.Logger {
	return Logger.With(getContext(ctx)...)
}

func Init() {
	var logLevel zapcore.Level
	if logLevel.UnmarshalText([]byte(config.GetAPP("LOG_LEVEL").String())) != nil {
		logLevel = zapcore.InfoLevel
	}
	logFile := config.GetAPP("LOG_FILE").String()
	writer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    500, // megabytes
		MaxBackups: 0,
		MaxAge:     28, // days
		LocalTime:  true,
	})
	sync := zapcore.AddSync(writer)

	jsonEncoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		LevelKey:       "level",
		NameKey:        "name",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stack",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})

	core := zapcore.NewTee(zapcore.NewCore(jsonEncoder, sync, logLevel))

	logger := zap.New(core, zap.AddCaller())
	defer logger.Sync()

	Logger = logger
}

func getContext(ctx *gin.Context) []zap.Field {
	var (
		now          = time.Now().Format(time.DateTime)
		startTime, _ = ctx.Get("startTime")
		duration     = float64(time.Now().Sub(startTime.(time.Time)).Nanoseconds()/1e4) / 100.0 //单位毫秒,保留2位小数
		serviceStart = startTime.(time.Time).Format(time.DateTime)
		request      = ctx.Request.RequestURI
		hostAddress  = ctx.Request.Host
		clientIp     = ctx.ClientIP()
		traceId      = ctx.GetString("traceId")
		parentId     = ctx.GetString("parentId")
		params       = ctx.Request.PostForm
	)
	return []zap.Field{
		zap.String("traceId", traceId),
		zap.String("serviceStart", serviceStart),
		zap.String("serviceEnd", now),
		zap.String("request", request),
		zap.String("params", params.Encode()),
		zap.String("hostAddress", hostAddress),
		zap.String("clientIp", clientIp),
		zap.String("parentId", parentId),
		zap.Float64("duration", duration)}
}
