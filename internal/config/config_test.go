package config

import (
	"github.com/vatsal278/go-redis-cache"
	"testing"

	"github.com/PereRohit/util/config"
	"github.com/PereRohit/util/testutil"
)

func TestInitSvcConfig(t *testing.T) {
	type args struct {
		cfg Config
	}
	tests := []struct {
		name string
		args args
		want *SvcConfig
	}{
		{
			name: "Success",
			args: args{
				cfg: Config{
					ServiceRouteVersion: "v2",
					ServerConfig:        config.ServerConfig{},
					DummyCfg: DummySvcCfg{
						DummyCfg: "dummy cfg",
					},
					Cache: CacheCfg{
						Port: "",
						Host: "",
					},
					MaxMemory: 1000,
				},
			},
			want: &SvcConfig{
				cfg: &Config{
					ServiceRouteVersion: "v2",
					ServerConfig:        config.ServerConfig{},
					DummyCfg: DummySvcCfg{
						DummyCfg: "dummy cfg",
					},
					Cache: CacheCfg{
						Port: "",
						Host: "",
					},
					MaxMemory: 1000,
				},
				ServiceRouteVersion: "v2",
				SvrCfg:              config.ServerConfig{},
				DummySvc:            DummyInternalSvc{},
				CacherSvc: func() CacherSvc {
					return CacherSvc{
						Cacher: redis.NewCacher(redis.Config{Addr: ":"})}
				}(),
				MaxMemmory: 1000,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := InitSvcConfig(tt.args.cfg)
			diff := testutil.Diff(got, tt.want)
			if diff != "" {
				t.Error(testutil.Callers(), diff)
			}
		})
	}
}
