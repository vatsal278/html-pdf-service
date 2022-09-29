package datasource

import (
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/vatsal278/go-redis-cache/mocks"
	"github.com/vatsal278/html-pdf-service/internal/config"
	"testing"
	"time"
)

func TestHealth(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	tests := []struct {
		name         string
		setupFunc    func() *mocks.MockCacher
		validateFunc func(bool)
	}{
		{
			name: "Success::Health Check",
			setupFunc: func() *mocks.MockCacher {
				mockcacher := mocks.NewMockCacher(mockCtrl)
				mockcacher.EXPECT().Health().Times(1).Return("PONG", nil)
				return mockcacher
			},
			validateFunc: func(s bool) {
				if s != true {
					t.Errorf("want %v got %v", "not nil", nil)
				}
			},
		},
		{
			name: "Failure:: Health Check",
			setupFunc: func() *mocks.MockCacher {
				mockcacher := mocks.NewMockCacher(mockCtrl)
				mockcacher.EXPECT().Health().Times(1).Return("", errors.New(""))
				return mockcacher
			},
			validateFunc: func(s bool) {
				if s != false {
					t.Errorf("want %v got %v", "not nil", nil)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockcacher := tt.setupFunc()
			cacherSvc := config.CacherSvc{Cacher: mockcacher}
			ds := NewRedisDs(&cacherSvc)
			x := ds.HealthCheck()
			tt.validateFunc(x)
		})
	}

}

func TestSaveFile(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	tests := []struct {
		name         string
		requestBody  string
		expiry       time.Duration
		setupFunc    func() *mocks.MockCacher
		validateFunc func(DataSource, string, error)
	}{
		{
			name:        "Success:: Save File",
			requestBody: "localhost:6379",
			setupFunc: func() *mocks.MockCacher {
				mockcacher := mocks.NewMockCacher(mockCtrl)
				mockcacher.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(nil)
				return mockcacher
			},
			validateFunc: func(cacher DataSource, data string, err error) {
				if err != nil {
					t.Errorf("want %v got %v", nil, err.Error())
				}
			},
		},
		{
			name:        "Success:: Save File:: With Expiry",
			requestBody: "localhost:6379",
			expiry:      2 * time.Second,
			setupFunc: func() *mocks.MockCacher {
				mockcacher := mocks.NewMockCacher(mockCtrl)
				mockcacher.EXPECT().Set(gomock.Any(), gomock.Any(), 2*time.Second).Times(1).Return(nil)
				return mockcacher
			},
			validateFunc: func(cacher DataSource, data string, err error) {
				if err != nil {
					t.Errorf("want %v got %v", nil, err.Error())
				}
			},
		},
		{
			name: "Failure:: Save File",
			setupFunc: func() *mocks.MockCacher {
				mockcacher := mocks.NewMockCacher(mockCtrl)
				mockcacher.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(errors.New(""))
				return mockcacher
			},
			validateFunc: func(cacher DataSource, data string, err error) {
				if err == nil {
					t.Errorf("want %v got %v", "not nil", nil)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockcacher := tt.setupFunc()
			cacherSvc := config.CacherSvc{Cacher: mockcacher}
			ds := NewRedisDs(&cacherSvc)
			err := ds.SaveFile("1", []byte("abc"), tt.expiry)
			tt.validateFunc(ds, "abc", err)
		})
	}

}

func TestDelete(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	tests := []struct {
		name         string
		requestBody  string
		setupFunc    func() *mocks.MockCacher
		validateFunc func(error)
	}{
		{
			name: "Success:: Delete File",
			setupFunc: func() *mocks.MockCacher {
				mockcacher := mocks.NewMockCacher(mockCtrl)
				mockcacher.EXPECT().Delete(gomock.Any()).Times(1).Return(nil)
				return mockcacher
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
			mockcacher := tt.setupFunc()
			cacherSvc := config.CacherSvc{Cacher: mockcacher}
			ds := NewRedisDs(&cacherSvc)
			err := ds.DeleteFile("1")
			tt.validateFunc(err)
		})
	}

}

func TestGet(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	tests := []struct {
		name         string
		requestBody  string
		setupFunc    func() *mocks.MockCacher
		validateFunc func([]byte, string, error)
	}{
		{
			name:        "Success:: Get file",
			requestBody: "1",
			setupFunc: func() *mocks.MockCacher {
				mockcacher := mocks.NewMockCacher(mockCtrl)
				mockcacher.EXPECT().Get(gomock.Any()).Times(1).Return([]byte("Hello"), nil)
				return mockcacher
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
			name:        "Failure:: Get file",
			requestBody: "2",
			setupFunc: func() *mocks.MockCacher {
				mockcacher := mocks.NewMockCacher(mockCtrl)
				mockcacher.EXPECT().Get(gomock.Any()).Times(1).Return(nil, errors.New(""))
				return mockcacher
			},
			validateFunc: func(s []byte, request string, err error) {
				if err == nil {
					t.Errorf("want %v got %v", nil, err)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockcacher := tt.setupFunc()
			cacherSvc := config.CacherSvc{Cacher: mockcacher}
			ds := NewRedisDs(&cacherSvc)
			key := tt.requestBody
			x, err := ds.GetFile(key)
			tt.validateFunc(x, "Hello", err)
		})
	}

}
