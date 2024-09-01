package main

import (
	"net/http"

	"ops-storage/internal/server/handlers"
	wr "ops-storage/internal/server/handlers/wrappers"
	"ops-storage/internal/server/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	logger.Initialize()

	opts := options{}
	Parse(&opts)

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.POST("/update", wr.LogWrapper(wr.CompressWrapper(handlers.UpdateJsonMetric)))
	router.POST("/update/:type/:name/:value",
		wr.LogWrapper(wr.CompressWrapper(handlers.UpdateQueryMetric)))

	router.POST("/value", wr.LogWrapper(wr.CompressWrapper(handlers.GetMetricViaJson)))
	router.GET("/value/:type/:name",
		wr.LogWrapper(wr.CompressWrapper(handlers.GetMetricViaQuery)))

	router.GET("/", wr.LogWrapper(wr.CompressWrapper(handlers.GetAllMetrics)))

	router.NoRoute(wr.LogWrapper(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": "Not found"})
	}))

	logger.MainLog.Infow("The server is running on localhost:8080")
	err := http.ListenAndServe(opts.endpoint, router)
	if err != nil {
		logger.MainLog.Panic(err)
	}
}
