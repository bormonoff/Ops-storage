package handlers

import (
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
)

func isTextType(cType string) bool {
	if cType != gin.MIMEPlain {
		return false
	}
	return true
}

func isTypeValid(kind string) bool {
	if mathed, _ := regexp.MatchString("^[[:alpha:]]+$", kind); !mathed {
		return false
	}
	return true
}

func isNameValid(name string) (int, bool) {
	if mathed, _ := regexp.MatchString("^[[:alnum:]]+$", name); !mathed {
		if name == "" {
			return http.StatusNotFound, false
		} else {
			return http.StatusBadRequest, false
		}
	}
	return 0, true
}

func isValueValid(val string) bool {
	intMatch, _ := regexp.MatchString("^[0-9]+$", val)
	floatMatch, _ := regexp.MatchString("^[[:digit:]]+[.]?[[:digit:]]+$", val)

	if !intMatch && !floatMatch {
		return false
	}
	return true
}
