package logic

import (
	"errors"
	"github.com/vatsal278/html-pdf-service/internal/codes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"reflect"
	"testing"

	respModel "github.com/PereRohit/util/model"
	"github.com/PereRohit/util/testutil"
	"github.com/golang/mock/gomock"

	"github.com/vatsal278/html-pdf-service/internal/model"
	"github.com/vatsal278/html-pdf-service/internal/repo/datasource"
	"github.com/vatsal278/html-pdf-service/pkg/mock"
)

func Test_htmlPdfServiceLogic_HealthCheck(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tests := []struct {
		name  string
		setup func() datasource.DataSource
		want  bool
	}{
		{
			name: "Success",
			setup: func() datasource.DataSource {
				mockDs := mock.NewMockDataSource(mockCtrl)

				mockDs.EXPECT().HealthCheck().Times(1).
					Return(true)

				return mockDs
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := NewHtmlPdfServiceLogic(tt.setup())

			got := rec.HealthCheck()

			diff := testutil.Diff(got, tt.want)
			if diff != "" {
				t.Error(testutil.Callers(), diff)
			}
		})
	}
}

func Test_htmlPdfServiceLogic_Ping(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tests := []struct {
		name  string
		setup func() datasource.DataSource
		give  *model.PingRequest
		want  *respModel.Response
	}{
		{
			name: "Success",
			setup: func() datasource.DataSource {
				mockDs := mock.NewMockDataSource(mockCtrl)

				mockDs.EXPECT().Ping(&model.PingDs{
					Data: "ping",
				}).Times(1).
					Return(&model.DsResponse{
						Data: "pong",
					}, nil)

				return mockDs
			},
			give: &model.PingRequest{
				Data: "ping",
			},
			want: &respModel.Response{
				Status:  http.StatusOK,
				Message: "Pong",
				Data: &model.DsResponse{
					Data: "pong",
				},
			},
		},
		{
			name: "Failure::datasource error",
			setup: func() datasource.DataSource {
				mockDs := mock.NewMockDataSource(mockCtrl)

				mockDs.EXPECT().Ping(&model.PingDs{
					Data: "ping",
				}).Times(1).
					Return(nil, errors.New("ds down"))

				return mockDs
			},
			give: &model.PingRequest{
				Data: "ping",
			},
			want: &respModel.Response{
				Status:  http.StatusInternalServerError,
				Message: "",
				Data:    nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := NewHtmlPdfServiceLogic(tt.setup())

			got := rec.Ping(tt.give)

			diff := testutil.Diff(got, tt.want)
			if diff != "" {
				t.Error(testutil.Callers(), diff)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	tests := []struct {
		name         string
		requestBody  string
		setupFunc    func() multipart.File
		validateFunc func(*model.Response)
	}{
		{
			name:        "Success:: Update",
			requestBody: "1",
			setupFunc: func() multipart.File {
				body, _ := ioutil.TempFile(".", "example")
				return multipart.File(body)
			},
			validateFunc: func(x *model.Response) {
				if x.Status != http.StatusCreated {
					t.Errorf("want %v got %v", http.StatusCreated, x)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := tt.setupFunc()

			x := cacher.Upload(data)
			tt.validateFunc(x)
		})
	}

}

func TestHtmlToPdf(t *testing.T) {

	os.Setenv("Address", "0.0.0.0")
	appContainer := config.GetAppContainer()
	cacher := NewHtmltopdfsvcLogic(appContainer)
	tests := []struct {
		name         string
		requestBody  string
		setupFunc    func() modelV.GenerateReq
		validateFunc func(*model.Response)
	}{
		{
			name:        "Success:: Htmltopdf",
			requestBody: "1",
			setupFunc: func() modelV.GenerateReq {
				//body, _ := ioutil.TempFile(".", "example")
				var data modelV.GenerateReq
				data.Id = "ee5371c2-7200-45a1-b543-e4c5bd4c48ed"
				return data
			},
			validateFunc: func(x *model.Response) {
				if x.Status != http.StatusCreated {
					t.Errorf("want %v got %v", http.StatusCreated, x)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := tt.setupFunc()
			var w io.Writer
			x := cacher.HtmlToPdf(w, &data)
			tt.validateFunc(x)
		})
	}

}

func TestReplace(t *testing.T) {
	os.Setenv("Address", "0.0.0.0")
	appContainer := config.GetAppContainer()
	cacher := NewHtmltopdfsvcLogic(appContainer)
	tests := []struct {
		name         string
		requestBody  string
		setupFunc    func() (multipart.File, string)
		validateFunc func(*model.Response)
	}{
		{
			name:        "Success:: Replace",
			requestBody: "1",
			setupFunc: func() (multipart.File, string) {
				body, _ := ioutil.TempFile(".", "example")
				return multipart.File(body), "ee5371c2-7200-45a1-b543-e4c5bd4c48ed"
			},
			validateFunc: func(x *model.Response) {
				if x.Status != http.StatusCreated {
					t.Errorf("want %v got %v", http.StatusCreated, x)
				}
			},
		},
		{
			name:        "Failure:: Replace:: no key found",
			requestBody: "1",
			setupFunc: func() (multipart.File, string) {
				body, _ := ioutil.TempFile(".", "example")
				return multipart.File(body), ""
			},
			validateFunc: func(x *model.Response) {
				var temp *model.Response
				if !reflect.DeepEqual(x, &model.Response{
					Status:  http.StatusInternalServerError,
					Message: codes.GetErr(codes.ErrKeyNotFound),
					Data:    nil,
				}) {
					t.Errorf("want %v got %v", temp, x)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, id := tt.setupFunc()

			x := cacher.Replace(id, data)
			tt.validateFunc(x)
		})
	}
}
