package handlers

import (
	"net/http"

	"ops-storage/internal/server/core"

	"github.com/gin-gonic/gin"
)

type errCodePair struct {
	code  int
	descr string
}

func UpdateMetric(c *gin.Context) {
	errCode, valid := validateUpdateMetric(c)
	if !valid {
		c.String(errCode.code, errCode.descr)
		return
	}

	err := core.GetStorageInstace().Insert(c.Param("type"), c.Param("name"), c.Param("value"))
	if err == core.ErrIvalidMetric {
		c.String(http.StatusBadRequest, "parsing counter error")
		return
	}
}

func validateUpdateMetric(c *gin.Context) (errCodePair, bool) {
	if !isTextType(c.ContentType()) {
		return errCodePair{
			code:  http.StatusBadRequest,
			descr: "request should be have a text/plain type\n",
		}, false
	}

	if !isTypeValid(c.Param("type")) {
		return errCodePair{
			code:  http.StatusBadRequest,
			descr: "parsing counter type error\n",
		}, false
	}

	if code, ok := isNameValid(c.Param("name")); !ok {
		if code == http.StatusNotFound {
			return errCodePair{
				code:  code,
				descr: "metric isn't found\n",
			}, false
		} else {
			return errCodePair{
				code:  code,
				descr: "parsing counter name error\n",
			}, false
		}

	}

	if !isValueValid(c.Param("value")) {
		return errCodePair{
			code:  http.StatusBadRequest,
			descr: "parsing counter value error\n",
		}, false
	}

	return errCodePair{}, true
}

func GetMetric(c *gin.Context) {
	errCode, valid := validateGetMetric(c)
	if !valid {
		c.String(errCode.code, errCode.descr)
		return
	}

	res, err := core.GetStorageInstace().GetMetric(c.Param("type"), c.Param("name"))
	if err == core.ErrNotFound {
		c.String(http.StatusNotFound, "parsing counter error")
		return
	}
	c.String(http.StatusOK, res)
}

func validateGetMetric(c *gin.Context) (errCodePair, bool) {
	if !isTextType(c.ContentType()) {
		return errCodePair{
			code:  http.StatusBadRequest,
			descr: "request should be have a text/plain type\n",
		}, false
	}

	if !isTypeValid(c.Param("type")) {
		return errCodePair{
			code:  http.StatusBadRequest,
			descr: "parsing counter type error\n",
		}, false
	}

	if code, ok := isNameValid(c.Param("name")); !ok {
		if code == http.StatusNotFound {
			return errCodePair{
				code:  code,
				descr: "metric isn't found\n",
			}, false
		} else {
			return errCodePair{
				code:  code,
				descr: "parsing counter name error\n",
			}, false
		}

	}

	return errCodePair{}, true
}

func GetAllMetrics(c *gin.Context) {
	errCode, valid := validateGetAllMetrics(c)
	if !valid {
		c.String(errCode.code, errCode.descr)
		return
	}
	res := core.GetStorageInstace().GetActualMetrics()

	c.JSON(http.StatusOK, res)
}

func validateGetAllMetrics(c *gin.Context) (errCodePair, bool) {
	if !isTextType(c.ContentType()) {
		return errCodePair{
			code:  http.StatusBadRequest,
			descr: "request should be have a text/plain type\n",
		}, false
	}

	return errCodePair{}, true
}
