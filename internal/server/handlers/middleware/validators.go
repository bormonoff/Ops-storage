package middleware

import (
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
)

func CheckIfPost(c *gin.Context) {
	if c.Request.Method != http.MethodPost {
		c.String(http.StatusMethodNotAllowed, "Only POST methods are allowed\n")
		c.Abort()
	}
}

func CheckContentTypeIsText(c *gin.Context) {
	contentType := c.ContentType()
	if contentType != gin.MIMEPlain {
		c.String(http.StatusBadRequest, "request should be have a text/plain type\n")
		c.Abort()
	}
}

func ValidateType(c *gin.Context) {
	kind := c.Param("type")
	if mathed, _ := regexp.MatchString("^[[:alpha:]]+$", kind); !mathed {
		c.String(http.StatusBadRequest, "parsing counter type error")
		c.Abort()
	}
}

func ValidateName(c *gin.Context) {
	name := c.Param("name")
	if mathed, _ := regexp.MatchString("^[[:alnum:]]+$", name); !mathed {
		if name == "" {
			c.String(http.StatusNotFound, "metric isn't found")
		} else {
			c.String(http.StatusBadRequest, "parsing counter name error")
		}
		c.Abort()
	}
}

func ValidateValue(c *gin.Context) {
	val := c.Param("value")
	if mathed, _ := regexp.MatchString("^([[:digit:]]|[.])+$", val); !mathed {
		c.String(http.StatusBadRequest, "parsing counter value error")
		c.Abort()
	}
}
