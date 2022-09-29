package logic

import (
	"encoding/base32"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/vatsal278/html-pdf-service/internal/codes"
	"github.com/vatsal278/html-pdf-service/internal/repo/htmlToPdf"
	"net/http"
	"strings"
	"testing"
	"time"

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
		setup func() (datasource.DataSource, htmlToPdf.HtmlToPdf)
		want  bool
	}{
		{
			name: "Success",
			setup: func() (datasource.DataSource, htmlToPdf.HtmlToPdf) {
				mockDs := mock.NewMockDataSource(mockCtrl)
				mockDs.EXPECT().HealthCheck().Times(1).
					Return(true)
				mockHt := mock.NewMockHtmlToPdf(mockCtrl)
				mockHt.EXPECT().HealthCheck().Times(1).Return(true)
				return mockDs, mockHt
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
		setup func() (datasource.DataSource, htmlToPdf.HtmlToPdf)
		give  *model.PingRequest
		want  *respModel.Response
	}{
		{
			name: "Success",
			setup: func() (datasource.DataSource, htmlToPdf.HtmlToPdf) {
				mockHt := mock.NewMockHtmlToPdf(mockCtrl)
				mockDs := mock.NewMockDataSource(mockCtrl)

				mockDs.EXPECT().Ping(&model.PingDs{
					Data: "ping",
				}).Times(1).
					Return(&model.DsResponse{
						Data: "pong",
					}, nil)

				return mockDs, mockHt
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
			setup: func() (datasource.DataSource, htmlToPdf.HtmlToPdf) {
				mockHt := mock.NewMockHtmlToPdf(mockCtrl)
				mockDs := mock.NewMockDataSource(mockCtrl)

				mockDs.EXPECT().Ping(&model.PingDs{
					Data: "ping",
				}).Times(1).
					Return(nil, errors.New("ds down"))

				return mockDs, mockHt
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
			ds, ht := tt.setup()
			rec := NewHtmlPdfServiceLogic(ds, ht)

			got := rec.Ping(tt.give)

			diff := testutil.Diff(got, tt.want)
			if diff != "" {
				t.Error(testutil.Callers(), diff)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	tests := []struct {
		name         string
		requestBody  string
		setupFunc    func() *htmlPdfServiceLogic
		validateFunc func(*respModel.Response)
	}{
		{
			name: "Success:: Update",
			setupFunc: func() *htmlPdfServiceLogic {
				mockHtmlsvc := mock.NewMockHtmlToPdf(mockCtrl)
				mockHtmlsvc.EXPECT().GetJsonFromHtml([]byte("abc")).Return([]byte("abc"), nil)
				mockDatasource := mock.NewMockDataSource(mockCtrl)
				mockDatasource.EXPECT().SaveFile(gomock.Any(), []byte("abc"), time.Duration(0)).Return(nil)
				rec := &htmlPdfServiceLogic{
					dsSvc: mockDatasource,
					htSvc: mockHtmlsvc,
				}
				return rec
			},
			validateFunc: func(x *respModel.Response) {
				if x.Status != http.StatusCreated {
					t.Errorf("want %v got %v", http.StatusCreated, x)
				}
			},
		},
		{
			name: "Failure:: Update",
			setupFunc: func() *htmlPdfServiceLogic {
				mockHtmlsvc := mock.NewMockHtmlToPdf(mockCtrl)
				mockHtmlsvc.EXPECT().GetJsonFromHtml([]byte("abc")).Return(nil, errors.New(""))
				mockDatasource := mock.NewMockDataSource(mockCtrl)
				rec := &htmlPdfServiceLogic{
					dsSvc: mockDatasource,
					htSvc: mockHtmlsvc,
				}
				return rec
			},
			validateFunc: func(x *respModel.Response) {
				if x.Status != http.StatusInternalServerError {
					t.Errorf("want %v got %v", http.StatusInternalServerError, x.Status)
				}
				if x.Message != codes.GetErr(codes.ErrFileConversionFail) {
					t.Errorf("want %v got %v", codes.GetErr(codes.ErrFileConversionFail), x.Message)
				}
			},
		},
		{
			name: "Failure:: Update",
			setupFunc: func() *htmlPdfServiceLogic {
				mockHtmlsvc := mock.NewMockHtmlToPdf(mockCtrl)
				mockHtmlsvc.EXPECT().GetJsonFromHtml([]byte("abc")).Return([]byte("abc"), nil)
				mockDatasource := mock.NewMockDataSource(mockCtrl)
				mockDatasource.EXPECT().SaveFile(gomock.Any(), []byte("abc"), time.Duration(0)).Return(errors.New(""))
				rec := &htmlPdfServiceLogic{
					dsSvc: mockDatasource,
					htSvc: mockHtmlsvc,
				}
				return rec
			},
			validateFunc: func(x *respModel.Response) {
				if x.Status != http.StatusInternalServerError {
					t.Errorf("want %v got %v", http.StatusInternalServerError, x.Status)
				}
				if x.Message != codes.GetErr(codes.ErrFileStoreFail) {
					t.Errorf("want %v got %v", codes.GetErr(codes.ErrFileStoreFail), x.Message)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := tt.setupFunc()
			resp := rec.Upload(strings.NewReader("abc"))
			tt.validateFunc(resp)
		})
	}

}

func TestReplace(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	tests := []struct {
		name         string
		requestBody  string
		setupFunc    func() *htmlPdfServiceLogic
		validateFunc func(*respModel.Response)
	}{
		{
			name: "Success:: Replace",
			setupFunc: func() *htmlPdfServiceLogic {
				mockHtmlsvc := mock.NewMockHtmlToPdf(mockCtrl)
				mockHtmlsvc.EXPECT().GetJsonFromHtml([]byte("abc")).Return([]byte("abc"), nil)
				mockDatasource := mock.NewMockDataSource(mockCtrl)
				mockDatasource.EXPECT().SaveFile(gomock.Any(), []byte("abc"), time.Duration(0)).Return(nil)
				mockDatasource.EXPECT().GetFile("1").Return([]byte(""), nil)
				rec := &htmlPdfServiceLogic{
					dsSvc: mockDatasource,
					htSvc: mockHtmlsvc,
				}
				return rec
			},
			validateFunc: func(x *respModel.Response) {
				if x.Status != http.StatusOK {
					t.Errorf("want %v got %v", http.StatusOK, x.Status)
				}
				if x.Message != "SUCCESS" {
					t.Errorf("want %v got %v", "SUCCESS", x.Message)
				}
			},
		},
		{
			name: "Failure:: Replace:: get file fail",
			setupFunc: func() *htmlPdfServiceLogic {
				mockHtmlsvc := mock.NewMockHtmlToPdf(mockCtrl)
				//mockHtmlsvc.EXPECT().GetJsonFromHtml([]byte("abc")).Return([]byte("abc"), nil)
				mockDatasource := mock.NewMockDataSource(mockCtrl)
				mockDatasource.EXPECT().GetFile("1").Return(nil, errors.New(""))
				rec := &htmlPdfServiceLogic{
					dsSvc: mockDatasource,
					htSvc: mockHtmlsvc,
				}
				return rec
			},
			validateFunc: func(x *respModel.Response) {
				if x.Status != http.StatusInternalServerError {
					t.Errorf("want %v got %v", http.StatusInternalServerError, x.Status)
				}
				if x.Message != codes.GetErr(codes.ErrKeyNotFound) {
					t.Errorf("want %v got %v", codes.GetErr(codes.ErrKeyNotFound), x.Message)
				}
			},
		},
		{
			name: "Failure:: Replace:: save file fail",
			setupFunc: func() *htmlPdfServiceLogic {
				mockHtmlsvc := mock.NewMockHtmlToPdf(mockCtrl)
				mockHtmlsvc.EXPECT().GetJsonFromHtml([]byte("abc")).Return([]byte("abc"), nil)
				mockDatasource := mock.NewMockDataSource(mockCtrl)
				mockDatasource.EXPECT().SaveFile(gomock.Any(), []byte("abc"), time.Duration(0)).Return(errors.New(""))
				mockDatasource.EXPECT().GetFile("1").Return([]byte(""), nil)
				rec := &htmlPdfServiceLogic{
					dsSvc: mockDatasource,
					htSvc: mockHtmlsvc,
				}
				return rec
			},
			validateFunc: func(x *respModel.Response) {
				if x.Status != http.StatusInternalServerError {
					t.Errorf("want %v got %v", http.StatusInternalServerError, x.Status)
				}
				if x.Message != codes.GetErr(codes.ErrFileStoreFail) {
					t.Errorf("want %v got %v", codes.GetErr(codes.ErrFileStoreFail), x.Message)
				}
			},
		},
		{
			name: "Failure:: Replace:: file conversion fail",
			setupFunc: func() *htmlPdfServiceLogic {
				mockHtmlsvc := mock.NewMockHtmlToPdf(mockCtrl)
				mockHtmlsvc.EXPECT().GetJsonFromHtml([]byte("abc")).Return(nil, errors.New(""))
				mockDatasource := mock.NewMockDataSource(mockCtrl)
				mockDatasource.EXPECT().GetFile("1").Return([]byte(""), nil)
				rec := &htmlPdfServiceLogic{
					dsSvc: mockDatasource,
					htSvc: mockHtmlsvc,
				}
				return rec
			},
			validateFunc: func(x *respModel.Response) {
				if x.Status != http.StatusInternalServerError {
					t.Errorf("want %v got %v", http.StatusInternalServerError, x.Status)
				}
				if x.Message != codes.GetErr(codes.ErrFileConversionFail) {
					t.Errorf("want %v got %v", codes.GetErr(codes.ErrFileConversionFail), x.Message)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := tt.setupFunc()
			resp := rec.Replace("1", strings.NewReader("abc"))
			tt.validateFunc(resp)
		})
	}

}

func TestHtmlToPdf(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	tests := []struct {
		name         string
		requestBody  string
		setupFunc    func() *htmlPdfServiceLogic
		validateFunc func(*respModel.Response)
	}{
		{
			name: "Success:: HtmlToPdf",
			setupFunc: func() *htmlPdfServiceLogic {
				//w := httptest.NewRecorder()
				js, _ := json.Marshal(map[string]interface{}{
					"Custom-Field": "hello",
					"Pages": []interface{}{
						"hello-world",
						map[string]interface{}{
							"Base64PageData": base64.StdEncoding.EncodeToString([]byte("abc")),
							"Custom-Data":    "world",
						},
					},
				})
				mockHtmlsvc := mock.NewMockHtmlToPdf(mockCtrl)
				mockHtmlsvc.EXPECT().GeneratePdf(gomock.Any(), gomock.Any()).Return(nil)
				mockDatasource := mock.NewMockDataSource(mockCtrl)
				mockDatasource.EXPECT().GetFile(gomock.Any()).Return(js, nil)
				rec := &htmlPdfServiceLogic{
					dsSvc: mockDatasource,
					htSvc: mockHtmlsvc,
				}
				return rec
			},
			validateFunc: func(x *respModel.Response) {
				expected := respModel.Response{}
				if x.Status != expected.Status {
					t.Errorf("want %v got %v", expected.Status, x.Status)
				}
			},
		},
		{
			name: "Failure:: HtmlToPdf:: get file fail",
			setupFunc: func() *htmlPdfServiceLogic {
				mockHtmlsvc := mock.NewMockHtmlToPdf(mockCtrl)
				mockDatasource := mock.NewMockDataSource(mockCtrl)
				mockDatasource.EXPECT().GetFile("1").Return(nil, errors.New(""))
				rec := &htmlPdfServiceLogic{
					dsSvc: mockDatasource,
					htSvc: mockHtmlsvc,
				}
				return rec
			},
			validateFunc: func(x *respModel.Response) {
				if x.Status != http.StatusInternalServerError {
					t.Errorf("want %v got %v", http.StatusInternalServerError, x.Status)
				}
				if x.Message != codes.GetErr(codes.ErrFetchingFile) {
					t.Errorf("want %v got %v", codes.GetErr(codes.ErrFetchingFile), x.Message)
				}
			},
		},
		{
			name: "Failure:: HtmlToPdf:: err unmarshalling json",
			setupFunc: func() *htmlPdfServiceLogic {
				mockHtmlsvc := mock.NewMockHtmlToPdf(mockCtrl)
				mockDatasource := mock.NewMockDataSource(mockCtrl)
				mockDatasource.EXPECT().GetFile("1").Return([]byte(""), nil)

				rec := &htmlPdfServiceLogic{
					dsSvc: mockDatasource,
					htSvc: mockHtmlsvc,
				}
				return rec
			},
			validateFunc: func(x *respModel.Response) {
				if x.Status != http.StatusInternalServerError {
					t.Errorf("want %v got %v", http.StatusInternalServerError, x.Status)
				}
				if x.Message != codes.GetErr(codes.ErrFileParseFail) {
					t.Errorf("want %v got %v", codes.GetErr(codes.ErrFileParseFail), x.Message)
				}
			},
		},
		{
			name: "Failure:: HtmlToPdf:: failed to decode base64 data",
			setupFunc: func() *htmlPdfServiceLogic {
				js, _ := json.Marshal(map[string]interface{}{
					"Custom-Field": "hello",
					"Pages": []interface{}{
						"hello-world",
						map[string]interface{}{
							"Base64PageData": base32.StdEncoding.EncodeToString([]byte("abc")),
							"Custom-Data":    "world",
						},
					},
				})
				mockHtmlsvc := mock.NewMockHtmlToPdf(mockCtrl)
				mockDatasource := mock.NewMockDataSource(mockCtrl)
				mockDatasource.EXPECT().GetFile("1").Return(js, nil)
				rec := &htmlPdfServiceLogic{
					dsSvc: mockDatasource,
					htSvc: mockHtmlsvc,
				}
				return rec
			},
			validateFunc: func(x *respModel.Response) {
				if x.Status != http.StatusInternalServerError {
					t.Errorf("want %v got %v", http.StatusInternalServerError, x.Status)
				}
				if x.Message != codes.GetErr(codes.ErrDecodingData) {
					t.Errorf("want %v got %v", codes.GetErr(codes.ErrDecodingData), x.Message)
				}
			},
		},
		{
			name: "Failure:: HtmlToPdf:: failed to generate pdf",
			setupFunc: func() *htmlPdfServiceLogic {
				js, _ := json.Marshal(map[string]interface{}{
					"Custom-Field": "hello",
					"Pages": []interface{}{
						"hello-world",
						map[string]interface{}{
							"Base64PageData": base64.StdEncoding.EncodeToString([]byte("abc")),
							"Custom-Data":    "world",
						},
					},
				})
				mockHtmlsvc := mock.NewMockHtmlToPdf(mockCtrl)
				mockHtmlsvc.EXPECT().GeneratePdf(gomock.Any(), gomock.Any()).Return(errors.New(""))
				mockDatasource := mock.NewMockDataSource(mockCtrl)
				mockDatasource.EXPECT().GetFile("1").Return(js, nil)
				rec := &htmlPdfServiceLogic{
					dsSvc: mockDatasource,
					htSvc: mockHtmlsvc,
				}
				return rec
			},
			validateFunc: func(x *respModel.Response) {
				if x.Status != http.StatusInternalServerError {
					t.Errorf("want %v got %v", http.StatusInternalServerError, x.Status)
				}
				if x.Message != codes.GetErr(codes.ErrConvertingToPdf) {
					t.Errorf("want %v got %v", codes.GetErr(codes.ErrConvertingToPdf), x.Message)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := tt.setupFunc()
			resp := rec.HtmlToPdf(nil, &model.GenerateReq{
				Values: nil,
				Id:     "1",
			})
			tt.validateFunc(resp)
		})
	}
}
