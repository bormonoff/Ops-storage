package main

import (
	"net/http"

	"ops-storage/internal/server/handlers"
	wr "ops-storage/internal/server/handlers/wrappers"
	"ops-storage/internal/server/logger"
	"ops-storage/internal/server/storage"

	"github.com/gin-gonic/gin"
)

func main() {
	logger.Initialize()

	opts := options{}
	Parse(&opts)
	
	storage.Init(opts.dsn, storage.RecoverConfig{
		RelPath:  opts.filePath,
		Interval: opts.storeInterval,
		Restore:  opts.restore,
	})

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.POST("/update", wr.LogWrapper(wr.CompressWrapper(handlers.UpdateJSONMetric)))
	router.POST("/update/:type/:name/:value",
		wr.LogWrapper(wr.CompressWrapper(handlers.UpdateQueryMetric)))

	router.POST("/value", wr.LogWrapper(wr.CompressWrapper(handlers.GetMetricViaJSON)))
	router.GET("/value/:type/:name",
		wr.LogWrapper(wr.CompressWrapper(handlers.GetMetricViaQuery)))

	router.GET("/", wr.LogWrapper(wr.CompressWrapper(handlers.GetAll)))
	router.GET("/ping", wr.LogWrapper(handlers.CheckDB))

	router.NoRoute(wr.LogWrapper(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": "Not found"})
	}))

	logger.MainLog.Infof("The server is running on %s", opts.endpoint)
	err := http.ListenAndServe(opts.endpoint, router)
	if err != nil {
		logger.MainLog.Panic(err)
	}
}
