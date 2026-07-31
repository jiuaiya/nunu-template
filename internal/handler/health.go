package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type DependencyCheck struct {
	Name  string
	Check func(context.Context) error
}

type Health struct {
	checks []DependencyCheck
}

func NewHealth(checks ...DependencyCheck) *Health {
	return &Health{checks: checks}
}

func (h *Health) Register(routes gin.IRoutes) {
	routes.GET("/health/live", h.live)
	routes.GET("/health/ready", h.ready)
}

func (h *Health) live(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *Health) ready(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second)
	defer cancel()

	for _, dependency := range h.checks {
		if err := dependency.Check(ctx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":     "unavailable",
				"dependency": dependency.Name,
			})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
