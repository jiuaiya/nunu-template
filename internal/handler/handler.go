package handler

import "go.uber.org/zap"

type Handler struct {
	Logger *zap.Logger
}

func New(logger *zap.Logger) *Handler {
	return &Handler{Logger: logger}
}
