package datasource

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	redisV "github.com/vatsal278/go-redis-cache"
	"github.com/vatsal278/html-pdf-service/internal/config"
	"testing"
	"time"
)

func TestHealth(t *testing.T) {
	tests := []struct {
		name         string
		requestBody  string
		validateFunc func(string, error)
	}{
		{
			name:        "Success::Health",
			requestBody: "localhost:6379",
			validateFunc: func(s string, err error) {
				if err != nil {
					t.Errorf("want %v got %v", nil, err.Error())
				}
				if s != "PONG" {
					t.Errorf("want %v got %v", "PONG", s)
				}
			},
		},
		{
			name:        "Failure:: Health",
			requestBody: "abc",
			validateFunc: func(s string, err error) {
				if err == nil {
					t.Errorf("want %v got %v", "not nil", nil)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cacher := redisV.NewCacher(redisV.Config{Addr: tt.requestBody})
			x, err := cacher.Health()
			tt.validateFunc(x, err)
		})
	}

}

func TestSaveFile(t *testing.T) {

	tests := []struct {
		name         string
		requestBody  string
		expiry       time.Duration
		validateFunc func(DataSource, string, error)
	}{
		{
			name:        "Success:: Set",
			requestBody: "localhost:6379",
			validateFunc: func(cacher DataSource, data string, err error) {
				if err != nil {
					t.Errorf("want %v got %v", nil, err.Error())
				}
				x, err := cacher.GetFile("1")
				if err != nil {
					t.Errorf("want %v got %v", nil, err.Error())
				}
				if fmt.Sprintf("%s", x) != data {
					t.Errorf("want %v got %v", data, fmt.Sprintf("%s", x))
				}
			},
		},
		{
			name:        "Success:: Set:: With Expiry",
			requestBody: "localhost:6379",
			expiry:      2 * time.Second,
			validateFunc: func(cacher DataSource, data string, err error) {
				if err != nil {
					t.Errorf("want %v got %v", nil, err.Error())
				}
				time.Sleep(2 * time.Second)
				_, err = cacher.GetFile("1")
				if err != redis.Nil {
					t.Errorf("want %v got %v", redis.Nil, err)
				}
			},
		},
		{
			name:        "Failure:: Set",
			requestBody: "localhost:6",
			validateFunc: func(cacher DataSource, data string, err error) {
				if err == nil {
					t.Errorf("want %v got %v", "not nil", nil)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cacher := redisV.NewCacher(redisV.Config{Addr: tt.requestBody})
			cacherSvc := config.CacherSvc{Cacher: cacher}
			ds := NewRedisDs(&cacherSvc)
			data := tt.requestBody
			err := ds.SaveFile("1", data, tt.expiry)
			tt.validateFunc(ds, data, err)
		})
	}

}

func TestDelete(t *testing.T) {
	tests := []struct {
		name         string
		requestBody  string
		setupFunc    func(DataSource, string)
		validateFunc func(error)
	}{
		{
			name:        "Success:: Delete",
			requestBody: "localhost:6379",
			setupFunc: func(cacher DataSource, data string) {
				err := cacher.SaveFile("1", data, 0)
				if err != nil {
					t.Errorf("want %v got %v", nil, err)
				}
			},
			validateFunc: func(err error) {
				if err != nil {
					t.Errorf("want %v got %v", nil, err)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cacher := redisV.NewCacher(redisV.Config{Addr: tt.requestBody})
			cacherSvc := config.CacherSvc{Cacher: cacher}
			ds := NewRedisDs(&cacherSvc)
			tt.setupFunc(ds, "1")
			err := ds.DeleteFile("1")
			tt.validateFunc(err)
		})
	}

}

func TestGet(t *testing.T) {
	tests := []struct {
		name         string
		requestBody  string
		setupFunc    func(DataSource, string)
		validateFunc func([]byte, string, error)
	}{
		{
			name:        "Success:: Get",
			requestBody: "1",
			setupFunc: func(cacher DataSource, data string) {
				err := cacher.SaveFile("1", data, 0)
				if err != nil {
					t.Errorf("want %v got %v", nil, err.Error())
				}
			},
			validateFunc: func(s []byte, request string, err error) {
				if err != nil {
					t.Errorf("want %v got %v", nil, err.Error())
				}
				if fmt.Sprintf("%s", s) != request {
					t.Errorf("want %v got %v", request, fmt.Sprintf("%s", s))
				}
			},
		},
		{
			name:        "Failure:: Get",
			requestBody: "2",
			setupFunc: func(cahcer DataSource, data string) {

			},
			validateFunc: func(s []byte, request string, err error) {
				if err != nil {
					t.Errorf("want %v got %v", nil, err)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cacher := redisV.NewCacher(redisV.Config{Addr: "localhost:6379"})
			cacherSvc := config.CacherSvc{Cacher: cacher}
			ds := NewRedisDs(&cacherSvc)
			key := tt.requestBody
			tt.setupFunc(ds, "Hello")
			x, err := ds.GetFile(key)
			tt.validateFunc(x, "Hello", err)
		})
	}

}
