package datasource

import (
	"github.com/vatsal278/html-pdf-service/internal/config"
	"github.com/vatsal278/html-pdf-service/internal/model"
)

type redisDs struct {
	redisSvc *config.CacherSvc
}

func NewRedisDs(cacherSvc *config.CacherSvc) DataSource {
	return &redisDs{
		redisSvc: cacherSvc,
	}
}
func (r redisDs) HealthCheck() bool {
	_, err := r.redisSvc.Cacher.Health()
	if err != nil {
		return false
	}
	return true
}
func (r redisDs) Ping(ds *model.PingDs) (*model.DsResponse, error) {
	return &model.DsResponse{}, nil
}
