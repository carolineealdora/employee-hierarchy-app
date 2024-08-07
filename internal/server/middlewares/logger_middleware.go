package middlewares

import (
	"time"

	"github.com/carolineealdora/employee-hierarchy-app/internal/pkg/utils/logger"
	"github.com/gin-gonic/gin"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log := logger.Log
		startTime := time.Now()
		ctx.Next()
		endTime := time.Now()

		latencyTime := endTime.Sub(startTime)
		reqMethod := ctx.Request.Method
		reqUri := ctx.Request.RequestURI
		statusCode := ctx.Writer.Status()
		clientIP := ctx.ClientIP()

		if lastErr := ctx.Errors.Last(); lastErr != nil {
			log.WithFields(map[string]any{
				"METHOD":    reqMethod,
				"URI":       reqUri,
				"STATUS":    statusCode,
				"LATENCY":   latencyTime,
				"CLIENT_IP": clientIP,
			}).Error(ctx.Errors)
			return
		}

		log.WithFields(map[string]any{
			"METHOD":    reqMethod,
			"URI":       reqUri,
			"STATUS":    statusCode,
			"LATENCY":   latencyTime,
			"CLIENT_IP": clientIP,
		}).Infof("REQUEST %s %s SUCCESS", reqMethod, reqUri)
	}
}
