package handlers

import (
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
)

func isTextType(contentType string) bool {
	if contentType != gin.MIMEPlain {
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
	if mathed, _ := regexp.MatchString("^[[:digit:]]+[.]?[[:digit:]]$", val); !mathed {
		return false
	}
	return true
}
