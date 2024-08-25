package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	
	"ops-storage/internal/server/core"

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

func UpdateQueryMetric(c *gin.Context) {
	validator := updateQueryValidator{MType: c.Param("type"), Name: c.Param("name"), Value: c.Param("value")}
	errCode, valid := validateUpdateMetric(validator.MType, validator.Name, validator.Value)
	if !valid {
		c.String(errCode.code, errCode.descr)
		return
	}

	err := core.GetStorageInstace().Insert(validator.MType, validator.Name, validator.Value)
	if err == core.ErrIvalidMetric {
		c.String(http.StatusBadRequest, "parsing counter error")
		return
	}
}

func GetMetricViaQuery(c *gin.Context) {
	validator := updateQueryValidator{MType: c.Param("type"), Name: c.Param("name")}
	errCode, valid := validateGetMetric(validator.MType, validator.Name)
	if !valid {
		c.String(errCode.code, errCode.descr)
		return
	}

	res, err := core.GetStorageInstace().GetMetric(validator.MType, validator.Name)
	if err == core.ErrNotFound {
		c.String(http.StatusNotFound, "parsing counter error")
		return
	}
	c.String(http.StatusOK, res)
}

type updateJsonValidator struct {
	MType string      `json:"type"`
	Name  string      `json:"id"`
	Delta json.Number `json:"delta,omitempty"`
	Value json.Number `json:"value,omitempty"`
}

func UpdateJsonMetric(c *gin.Context) {
	if c.ContentType() != gin.MIMEJSON {
		c.String(http.StatusBadRequest, "request should have an application/json type\n")
		return
	}

	var validator updateJsonValidator
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

	if validator.Delta != "" {
		err = core.GetStorageInstace().Insert(validator.MType, validator.Name, string(validator.Delta))
	} else {
		err = core.GetStorageInstace().Insert(validator.MType, validator.Name, string(validator.Value))
	}
	if err == core.ErrIvalidMetric {
		c.String(http.StatusBadRequest, "parsing counter error")
		return
	}

	c.JSON(http.StatusOK, validator)
}

func GetMetricViaJson(c *gin.Context) {
	if c.ContentType() != gin.MIMEJSON {
		c.String(http.StatusBadRequest, "request should have an application/json type\n")
		return
	}

	var validator updateJsonValidator
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

	res, err := core.GetStorageInstace().GetMetric(validator.MType, validator.Name)
	if err == core.ErrNotFound {
		c.String(http.StatusNotFound, fmt.Sprintf("%s counter not found\n", validator.Name))
		return
	}
	if validator.MType == "gauge" {
		validator.Value = json.Number(res)
	} else {
		validator.Delta = json.Number(res)
	}
	c.JSON(http.StatusOK, validator)
}

func GetAllMetrics(c *gin.Context) {
	res := core.GetStorageInstace().GetActualMetrics()

	c.JSON(http.StatusOK, res)
}
