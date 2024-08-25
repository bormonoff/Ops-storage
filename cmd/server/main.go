package main

import (
	"net/http"

	"ops-storage/internal/server/handlers"
	"ops-storage/internal/server/handlers/wrappers"
	"ops-storage/internal/server/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	logger.Initialize()

	opts := options{}
	Parse(&opts)

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.POST("/update", wrappers.LogWrapper(handlers.UpdateJsonMetric))
	router.POST("/update/:type/:name/:value", wrappers.LogWrapper(handlers.UpdateQueryMetric))

	router.POST("/value", wrappers.LogWrapper(handlers.GetMetricViaJson))
	router.GET("/value/:type/:name", wrappers.LogWrapper(handlers.GetMetricViaQuery))

	router.GET("/", wrappers.LogWrapper(handlers.GetAllMetrics))

	router.NoRoute(wrappers.LogWrapper(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": "Not found"})
	}))

	logger.MainLog.Infow("The server is running on localhost:8080")
	err := http.ListenAndServe(opts.endpoint, router)
	if err != nil {
		logger.MainLog.Panic(err)
	}
}
