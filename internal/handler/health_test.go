package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestHealthReadyReportsDependencyFailure(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	NewHealth(DependencyCheck{
		Name: "mysql",
		Check: func(context.Context) error {
			return errors.New("offline")
		},
	}).Register(router)

	request := httptest.NewRequest(http.MethodGet, "/health/ready", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusServiceUnavailable)
	}
}

func TestHealthLive(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	NewHealth().Register(router)

	request := httptest.NewRequest(http.MethodGet, "/health/live", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}
}
