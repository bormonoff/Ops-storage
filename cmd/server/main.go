package main

import (
	"net/http"

	"ops-storage/internal/server/handlers"
	"ops-storage/internal/server/handlers/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	opts := options{}
	Parse(&opts)

	router := gin.Default()

	router.POST("/update/:type/:name/:value",
		middleware.ValidateType,
		middleware.ValidateName,
		handlers.UpdateMetric)

	router.GET("value/:type/:name",
		middleware.ValidateType,
		middleware.ValidateName,
		handlers.GetMetric)

	router.GET("/", handlers.GetAllMetrics)

	err := http.ListenAndServe(opts.endpoint, router)
	if err != nil {
		panic(err)
	}
}
