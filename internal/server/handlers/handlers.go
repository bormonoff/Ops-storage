package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"ops-storage/internal/server/storage"
	"strings"

	"github.com/gin-gonic/gin"
)

type errCodePair struct {
	code  int
	descr string
}

type updateQueryValidator struct {
	MType string
	Name  string
	Value string
}

func (v updateQueryValidator) String() string {
	return fmt.Sprint(v.MType, v.Name, v.Value)
}

func UpdateQueryMetric(c *gin.Context) {
	validator := updateQueryValidator{MType: c.Param("type"), Name: c.Param("name"), Value: c.Param("value")}
	errCode, valid := validateUpdateMetric(validator.MType, validator.Name, validator.Value)
	if !valid {
		c.String(errCode.code, errCode.descr)
		return
	}

	err := storage.StorageInstace().Insert(validator.MType, validator.Name, validator.Value)
	if err == storage.ErrIvalidMetric {
		c.String(http.StatusBadRequest, "parsing counter error")
		return
	}
	c.String(http.StatusOK, validator.String())
}

func GetMetricViaQuery(c *gin.Context) {
	validator := updateQueryValidator{MType: c.Param("type"), Name: c.Param("name")}
	errCode, valid := validateGetMetric(validator.MType, validator.Name)
	if !valid {
		c.String(errCode.code, errCode.descr)
		return
	}

	res, err := storage.StorageInstace().GetMetric(validator.MType, validator.Name)
	if err == storage.ErrNotFound {
		c.String(http.StatusNotFound, "parsing counter error")
		return
	}
	c.String(http.StatusOK, res)
}

type updateJSONValidator struct {
	MType   string      `json:"type"`
	Name    string      `json:"id"`
	Counter json.Number `json:"delta,omitempty"`
	Gauge   json.Number `json:"value,omitempty"`
}

func UpdateJSONMetric(c *gin.Context) {
	if c.ContentType() != gin.MIMEJSON {
		c.String(http.StatusBadRequest, "request should have an application/json type\n")
		return
	}

	var validator updateJSONValidator
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	err = json.Unmarshal(body, &validator)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	if validator.Counter != "" {
		err = storage.StorageInstace().Insert(validator.MType, validator.Name, string(validator.Counter))
	} else {
		err = storage.StorageInstace().Insert(validator.MType, validator.Name, string(validator.Gauge))
	}
	if err == storage.ErrIvalidMetric {
		c.String(http.StatusBadRequest, "parsing counter error")
		return
	}

	c.JSON(http.StatusOK, validator)
}

func GetMetricViaJSON(c *gin.Context) {
	if c.ContentType() != gin.MIMEJSON {
		c.String(http.StatusBadRequest, "request should have an application/json type\n")
		return
	}

	var validator updateJSONValidator
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Println(err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	err = json.Unmarshal(body, &validator)
	if err != nil {
		fmt.Println(err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	res, err := storage.StorageInstace().GetMetric(validator.MType, validator.Name)
	if err == storage.ErrNotFound {
		c.String(http.StatusNotFound, fmt.Sprintf("%s counter not found\n", validator.Name))
		return
	}
	if validator.MType == "gauge" {
		validator.Gauge = json.Number(res)
	} else {
		validator.Counter = json.Number(res)
	}
	c.JSON(http.StatusOK, validator)
}

func GetAllMetrics(c *gin.Context) {
	metrics := storage.StorageInstace().GetAllMetrics()

	var text strings.Builder
	for name, val := range *metrics {
		text.WriteString(fmt.Sprintf("%s = %s<br>", name, val))
	}
	c.Header("Content-Type", "text/html")

	c.String(http.StatusOK, text.String())
}
