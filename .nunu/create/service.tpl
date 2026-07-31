package service

import (
	"context"

	"{{ .ProjectName }}/internal/model"
	"{{ .ProjectName }}/internal/repository"
)

type {{ .StructName }}Service interface {
	Get(context.Context, uint64) (*model.{{ .StructName }}, error)
}

type {{ .StructNameLowerFirst }}Service struct {
	service    *Service
	repository repository.{{ .StructName }}Repository
}

func New{{ .StructName }}Service(service *Service, repository repository.{{ .StructName }}Repository) {{ .StructName }}Service {
	return &{{ .StructNameLowerFirst }}Service{service: service, repository: repository}
}

func (s *{{ .StructNameLowerFirst }}Service) Get(ctx context.Context, id uint64) (*model.{{ .StructName }}, error) {
	return s.repository.FindByID(ctx, id)
}
