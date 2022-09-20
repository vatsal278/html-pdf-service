package datasource

import (
	"github.com/vatsal278/html-pdf-service/internal/model"
	"time"
)

//go:generate mockgen --build_flags=--mod=mod --destination=./../../../pkg/mock/mock_datasource.go --package=mock github.com/vatsal278/html-pdf-service/internal/repo/datasource DataSource

type DataSource interface {
	HealthCheck() bool
	Ping(*model.PingDs) (*model.DsResponse, error)
	Get(s string) ([]byte, error)
	Set(key string, val interface{}, exp time.Duration) error
	Delete(key string) error
}
