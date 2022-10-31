package config

import (
	"github.com/PereRohit/util/config"
	"github.com/vatsal278/go-redis-cache"
)

type Config struct {
	ServiceRouteVersion string              `json:"service_route_version"`
	ServerConfig        config.ServerConfig `json:"server_config"`
	// add custom config structs below for any internal services
	Cache     CacheCfg `json:"cache"`
	MaxMemory int64    `json:"max_memory"`
}

type CacheCfg struct {
	Port string `json:"port"`
	Host string `json:"host"`
}

type SvcConfig struct {
	cfg                 *Config
	ServiceRouteVersion string
	SvrCfg              config.ServerConfig
	// add internal services after init
	CacherSvc  CacherSvc
	MaxMemmory int64
}

type CacherSvc struct {
	Cacher redis.Cacher
}

func InitSvcConfig(cfg Config) *SvcConfig {
	// init required services and assign to the service struct fields
	cacher := redis.NewCacher(redis.Config{Addr: cfg.Cache.Host + ":" + cfg.Cache.Port})
	return &SvcConfig{
		cfg:                 &cfg,
		ServiceRouteVersion: cfg.ServiceRouteVersion,
		SvrCfg:              cfg.ServerConfig,
		CacherSvc:           CacherSvc{Cacher: cacher},
		MaxMemmory:          cfg.MaxMemory,
	}
}
