package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ryo-arima/locky/pkg/server/middleware"
	"github.com/stretchr/testify/assert"
)

func TestGetRequestID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	// Test without middleware
	requestID := middleware.GetRequestID(c)
	assert.Empty(t, requestID)

	// Test with middleware
	c.Set("request_id", "test-request-id-123")
	requestID = middleware.GetRequestID(c)
	assert.Equal(t, "test-request-id-123", requestID)
}

func TestRequestIDMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.RequestIDMiddleware())

	router.GET("/test", func(c *gin.Context) {
		requestID := middleware.GetRequestID(c)
		assert.NotEmpty(t, requestID)
		c.JSON(http.StatusOK, gin.H{"request_id": requestID})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
