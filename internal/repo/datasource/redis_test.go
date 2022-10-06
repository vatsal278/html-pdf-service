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

func Test_Health(t *testing.T) {
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

func Test_SaveFile(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	tests := []struct {
		name         string
		expiry       time.Duration
		setupFunc    func() *mocks.MockCacher
		validateFunc func(DataSource, string, error)
	}{
		{
			name: "Success:: Save File",
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
			name:   "Success:: Save File:: With Expiry",
			expiry: 2 * time.Second,
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
			name: "Failure:: Save File :: err saving file",
			setupFunc: func() *mocks.MockCacher {
				mockcacher := mocks.NewMockCacher(mockCtrl)
				mockcacher.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(errors.New("Failed to set cache"))
				return mockcacher
			},
			validateFunc: func(cacher DataSource, data string, err error) {
				if err == errors.New("Failed to set cache") {
					t.Errorf("want %v got %v", errors.New("Failed to set cache"), nil)
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

func Test_DeleteFile(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	tests := []struct {
		name         string
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
		{
			name: "Failure:: Delete File :: err deleting file",
			setupFunc: func() *mocks.MockCacher {
				mockcacher := mocks.NewMockCacher(mockCtrl)
				mockcacher.EXPECT().Delete(gomock.Any()).Times(1).Return(errors.New("failed to delete cache"))
				return mockcacher
			},
			validateFunc: func(err error) {
				t.Log(err)
				if err == errors.New("failed to delete cache") {
					t.Errorf("want %v got %v", errors.New("failed to delete cache"), err)
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

func Test_GetFile(t *testing.T) {
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
				mockcacher.EXPECT().Get(gomock.Any()).Times(1).Return(nil, errors.New("Failed to get cache"))
				return mockcacher
			},
			validateFunc: func(s []byte, request string, err error) {
				if err == errors.New("Failed to get cache") {
					t.Errorf("want %v got %v", errors.New("Failed to get cache"), err)
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
