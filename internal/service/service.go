package service

import "go.uber.org/zap"

type Service struct {
	Logger *zap.Logger
}

func New(logger *zap.Logger) *Service {
	return &Service{Logger: logger}
}
