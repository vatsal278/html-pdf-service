package datasource

import (
	"time"
)

//go:generate mockgen --build_flags=--mod=mod --destination=./../../../pkg/mock/mock_datasource.go --package=mock github.com/vatsal278/html-pdf-service/internal/repo/datasource DataSource

type DataSource interface {
	HealthCheck() bool
	GetFile(s string) ([]byte, error)
	SaveFile(key string, val interface{}, exp time.Duration) error
	DeleteFile(key string) error
}
