package datasource

import (
	"github.com/vatsal278/html-pdf-service/internal/config"
	"github.com/vatsal278/html-pdf-service/internal/model"
	"time"
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
func (r redisDs) GetFile(s string) ([]byte, error) {
	x := r.redisSvc.Cacher
	val, err := x.Get(s)
	if err != nil {
		return nil, err
	}
	return val, nil
}

func (r redisDs) SaveFile(key string, val interface{}, exp time.Duration) error {
	x := r.redisSvc.Cacher
	err := x.Set(key, val, exp)
	if err != nil {
		return err
	}
	return nil
}
func (r redisDs) DeleteFile(key string) error {
	x := r.redisSvc.Cacher
	err := x.Delete(key)
	if err != nil {
		return err
	}
	return nil
}
