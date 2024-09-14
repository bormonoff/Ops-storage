package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"ops-storage/internal/server/storage"
	serror "ops-storage/internal/server/storage/error"

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

	err := storage.Instance().Insert(validator.MType, validator.Name, validator.Value)
	if err != nil {
		if errors.Is(err, serror.ErrIvalidMetric) {
			c.String(http.StatusBadRequest, "parse error")
			return
		}
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.String(http.StatusOK, validator.String())
}

func GetMetricViaQuery(c *gin.Context) {
	validator := updateQueryValidator{MType: c.Param("type"), Name: c.Param("name")}
	errCode, valid := validateGet(validator.MType, validator.Name)
	if !valid {
		c.String(errCode.code, errCode.descr)
		return
	}

	res, err := storage.Instance().Get(validator.MType, validator.Name)
	if err != nil {
		if errors.Is(err, serror.ErrNotFound) {
			c.String(http.StatusNotFound, "parse error")
			return
		}
		c.String(http.StatusInternalServerError, err.Error())
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
		err = storage.Instance().Insert(validator.MType, validator.Name, string(validator.Counter))
	} else {
		err = storage.Instance().Insert(validator.MType, validator.Name, string(validator.Gauge))
	}
	if err != nil {
		if errors.Is(err, serror.ErrIvalidMetric) {
			c.String(http.StatusBadRequest, "parse error")
			return
		}
		c.String(http.StatusInternalServerError, err.Error())
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

	res, err := storage.Instance().Get(validator.MType, validator.Name)
	if err != nil {
		if errors.Is(err, serror.ErrNotFound) {
			c.String(http.StatusNotFound, "parse error")
			return
		}
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	if validator.MType == "gauge" {
		validator.Gauge = json.Number(res)
	} else {
		validator.Counter = json.Number(res)
	}
	c.JSON(http.StatusOK, validator)
}

func GetAll(c *gin.Context) {
	metrics, err := storage.Instance().GetAll()

	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	var text strings.Builder
	for name, val := range *metrics {
		text.WriteString(fmt.Sprintf("%s = %s<br>", name, val))
	}
	c.Header("Content-Type", "text/html")

	c.String(http.StatusOK, text.String())
}

func CheckDB(c *gin.Context) {
	if storage.Instance().IsStorageAlive() {
		c.String(http.StatusOK, "Db is active")
		return
	}
	c.String(http.StatusInternalServerError, "Db is inactive")
}
