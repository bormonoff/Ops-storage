package wrappers

import (
	"time"
	
	"ops-storage/internal/server/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func LogWrapper(fn func(cxt *gin.Context)) func(cxt *gin.Context) {
	res := func(ctx *gin.Context) {
		start := time.Now()
		reqUUID := uuid.NewString()
		logger.HandlerLog.Info("Request:",
			zap.String("req-ID", reqUUID),
			zap.String("method", ctx.Request.Method),
			zap.String("url", ctx.Request.URL.String()))

		fn(ctx)

		logger.HandlerLog.Info("Response:",
			zap.String("req-ID", reqUUID),
			zap.Int("code", ctx.Writer.Status()),
			zap.Int("length", ctx.Writer.Size()),
			zap.String("dur", time.Since(start).String()))
	}
	return res
}
