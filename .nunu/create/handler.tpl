package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"{{ .ProjectName }}/internal/service"
)

type {{ .StructName }}Handler struct {
	handler *Handler
	service service.{{ .StructName }}Service
}

func New{{ .StructName }}Handler(handler *Handler, service service.{{ .StructName }}Service) *{{ .StructName }}Handler {
	return &{{ .StructName }}Handler{handler: handler, service: service}
}

func (h *{{ .StructName }}Handler) Register(routes gin.IRoutes) {
	routes.GET("/{{ .StructNameSnakeCase }}/:id", h.get)
}

func (h *{{ .StructName }}Handler) get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "invalid_id", "message": "invalid id"})
		return
	}
	result, err := h.service.Get(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": "internal_error", "message": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}
