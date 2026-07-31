package repository

import (
	"context"

	"{{ .ProjectName }}/internal/model"
)

type {{ .StructName }}Repository interface {
	FindByID(context.Context, uint64) (*model.{{ .StructName }}, error)
}

type {{ .StructNameLowerFirst }}Repository struct {
	repository *Repository
}

func New{{ .StructName }}Repository(repository *Repository) {{ .StructName }}Repository {
	return &{{ .StructNameLowerFirst }}Repository{repository: repository}
}

func (r *{{ .StructNameLowerFirst }}Repository) FindByID(ctx context.Context, id uint64) (*model.{{ .StructName }}, error) {
	var result model.{{ .StructName }}
	if err := r.repository.DB(ctx).First(&result, id).Error; err != nil {
		return nil, err
	}
	return &result, nil
}
