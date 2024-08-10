package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"ops-storage/internal/server/handlers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCheckIfPost(t *testing.T) {
	type input struct {
		method      string
		url         string
		contentType string
	}
	type expect struct {
		statusCode int
	}

	tests := []struct {
		name   string
		input  input
		expect expect
	}{
		{
			name: "valid request: expect 200 code",
			input: input{
				method:      http.MethodPost,
				url:         "/update/gauge/someMetric/1.0",
				contentType: "text/plain",
			},
			expect: expect{
				statusCode: http.StatusOK,
			},
		},
		{
			name: "invalid request: only POST method is allowed",
			input: input{
				method:      http.MethodGet,
				url:         "/update/gauge/someMetric/1.0",
				contentType: "text/plain",
			},
			expect: expect{
				statusCode: http.StatusMethodNotAllowed,
			},
		},
	}

	r := gin.New()
	r.Any("/update/:type/:name/:value", CheckIfPost,
		CheckContentTypeIsText,
		ValidateType,
		ValidateName,
		handlers.UpdateMetric)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.input.method, tt.input.url, nil)
			request.Header["Content-Type"] = []string{tt.input.contentType}

			w := httptest.NewRecorder()

			r.ServeHTTP(w, request)
			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, res.StatusCode, tt.expect.statusCode)
		})
	}
}

func TestContentTypeIsText(t *testing.T) {
	type input struct {
		method      string
		url         string
		contentType string
	}
	type expect struct {
		statusCode int
	}

	tests := []struct {
		name   string
		input  input
		expect expect
	}{
		{
			name: "valid request: expect 200 code",
			input: input{
				method:      http.MethodPost,
				url:         "/update/gauge/someMetric/1.0",
				contentType: "text/plain",
			},
			expect: expect{
				statusCode: http.StatusOK,
			},
		},
		{
			name: "invalid request: Content-Type header has to contain text/plain",
			input: input{
				method:      http.MethodPost,
				url:         "/update/gauge/someMetric/1.0",
				contentType: "application/json",
			},
			expect: expect{
				statusCode: http.StatusBadRequest,
			},
		},
	}

	r := gin.New()
	r.Any("/update/:type/:name/:value", CheckIfPost,
		CheckContentTypeIsText,
		ValidateType,
		ValidateName,
		handlers.UpdateMetric)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.input.method, tt.input.url, nil)
			request.Header["Content-Type"] = []string{tt.input.contentType}

			w := httptest.NewRecorder()

			r.ServeHTTP(w, request)
			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, res.StatusCode, tt.expect.statusCode)
		})
	}
}
