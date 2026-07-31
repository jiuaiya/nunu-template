package model

import "time"

type {{ .StructName }} struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ({{ .StructName }}) TableName() string {
	return "{{ .StructNameSnakeCase }}"
}
