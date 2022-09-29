package config

import (
	"encoding/json"
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
		want string
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
			want: func() string {
				b, err := json.Marshal(&SvcConfig{
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
				})
				if err != nil {
					t.Error(err)
				}
				return string(b)
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := InitSvcConfig(tt.args.cfg)
			b, err := json.Marshal(got)
			if err != nil {
				t.Error(err)
			}
			diff := testutil.Diff(string(b), tt.want)
			if diff != "" {
				t.Error(testutil.Callers(), diff)
			}
		})
	}
}
