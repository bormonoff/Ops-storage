package handlers

import (
	"net/http"

	"ops-storage/internal/server/core"

	"github.com/gin-gonic/gin"
)

func UpdateMetric(c *gin.Context) {
	err := core.GetStorageInstace().Insert(c.Param("type"), c.Param("name"), c.Param("value"))
	if err == core.ErrIvalidMetric {
		c.String(http.StatusBadRequest, "parsing counter error")
		return
	}
}

func GetMetric(c *gin.Context) {
	res, err := core.GetStorageInstace().GetMetric(c.Param("type"), c.Param("name"))
	if err == core.ErrNotFound {
		c.String(http.StatusNotFound, "parsing counter error")
		return
	}
	c.String(http.StatusOK, res)
}

func GetAllMetrics(c *gin.Context) {
	res := core.GetStorageInstace().GetActualMetrics()

	c.JSON(http.StatusOK, res)
}
