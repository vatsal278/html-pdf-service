package logic

import (
	"bytes"
	"encoding/base32"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/PereRohit/util/log"
	"github.com/vatsal278/html-pdf-service/internal/codes"
	"github.com/vatsal278/html-pdf-service/internal/repo/htmlToPdf"
	"io"
	"net/http"
	"reflect"
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

type Reader string

func (Reader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}

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

func Test_Upload(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tests := []struct {
		name         string
		requestBody  interface{}
		setupFunc    func() *htmlPdfServiceLogic
		validateFunc func(*respModel.Response)
	}{
		{
			name:        "Success:: Upload",
			requestBody: strings.NewReader("abc"),
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
				expected := respModel.Response{
					Status:  http.StatusCreated,
					Message: "SUCCESS",
					Data: map[string]interface{}{
						"id": gomock.Any(),
					},
				}
				if x.Status != expected.Status {
					t.Errorf("want %v got %v", expected, x)
				}
				if x.Message != expected.Message {
					t.Errorf("want %v got %v", expected, x)
				}
			},
		},
		{
			name:        "Failure:: Upload :: Read file failure",
			requestBody: Reader(""),
			setupFunc: func() *htmlPdfServiceLogic {
				mockHtmlsvc := mock.NewMockHtmlToPdf(mockCtrl)
				mockDatasource := mock.NewMockDataSource(mockCtrl)
				rec := &htmlPdfServiceLogic{
					dsSvc: mockDatasource,
					htSvc: mockHtmlsvc,
				}
				return rec
			},
			validateFunc: func(x *respModel.Response) {
				expected := respModel.Response{
					Status:  http.StatusInternalServerError,
					Message: codes.GetErr(codes.ErrReadFileFail),
					Data:    nil,
				}
				if !reflect.DeepEqual(x, &expected) {
					t.Errorf("want %v got %v", expected, x)
				}
			},
		},
		{
			name:        "Failure:: Upload :: GetJsonFromHtml failure",
			requestBody: strings.NewReader("abc"),
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
				expected := respModel.Response{
					Status:  http.StatusInternalServerError,
					Message: codes.GetErr(codes.ErrFileConversionFail),
					Data:    nil,
				}
				if !reflect.DeepEqual(x, &expected) {
					t.Errorf("want %v got %v", expected, x)
				}
			},
		},
		{
			name:        "Failure:: Upload :: SaveFile failure",
			requestBody: strings.NewReader("abc"),
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
				expected := respModel.Response{
					Status:  http.StatusInternalServerError,
					Message: codes.GetErr(codes.ErrFileStoreFail),
					Data:    nil,
				}
				if !reflect.DeepEqual(x, &expected) {
					t.Errorf("want %v got %v", expected, x)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := tt.setupFunc()
			resp := rec.Upload(tt.requestBody.(io.Reader))
			tt.validateFunc(resp)
		})
	}

}

func Test_Replace(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	tests := []struct {
		name         string
		requestBody  interface{}
		setupFunc    func() *htmlPdfServiceLogic
		validateFunc func(*respModel.Response)
	}{
		{
			name:        "Success:: Replace",
			requestBody: strings.NewReader("abc"),
			setupFunc: func() *htmlPdfServiceLogic {
				mockHtmlsvc := mock.NewMockHtmlToPdf(mockCtrl)
				mockHtmlsvc.EXPECT().GetJsonFromHtml([]byte("abc")).Return([]byte("abc"), nil)
				mockDatasource := mock.NewMockDataSource(mockCtrl)
				mockDatasource.EXPECT().SaveFile("1", []byte("abc"), time.Duration(0)).Return(nil)
				mockDatasource.EXPECT().GetFile("1").Return([]byte(""), nil)
				rec := &htmlPdfServiceLogic{
					dsSvc: mockDatasource,
					htSvc: mockHtmlsvc,
				}
				return rec
			},
			validateFunc: func(x *respModel.Response) {
				expected := respModel.Response{
					Status:  http.StatusOK,
					Message: "SUCCESS",
					Data:    map[string]interface{}{"id": "1"},
				}
				if !reflect.DeepEqual(x, &expected) {
					t.Errorf("want %v got %v", expected, x)
				}
			},
		},
		{
			name:        "Failure:: Replace:: GetFile fail",
			requestBody: strings.NewReader("abc"),
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
				expected := respModel.Response{
					Status:  http.StatusInternalServerError,
					Message: codes.GetErr(codes.ErrKeyNotFound),
					Data:    nil,
				}
				if !reflect.DeepEqual(x, &expected) {
					t.Errorf("want %v got %v", expected, x)
				}
			},
		},
		{
			name:        "Failure:: Replace:: ReadFile fail",
			requestBody: Reader(""),
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
				expected := respModel.Response{
					Status:  http.StatusInternalServerError,
					Message: codes.GetErr(codes.ErrReadFileFail),
					Data:    nil,
				}
				if !reflect.DeepEqual(x, &expected) {
					t.Errorf("want %v got %v", expected, x)
				}
			},
		},
		{
			name:        "Failure:: Replace:: SaveFile fail",
			requestBody: strings.NewReader("abc"),
			setupFunc: func() *htmlPdfServiceLogic {
				mockHtmlsvc := mock.NewMockHtmlToPdf(mockCtrl)
				mockHtmlsvc.EXPECT().GetJsonFromHtml([]byte("abc")).Return([]byte("abc"), nil)
				mockDatasource := mock.NewMockDataSource(mockCtrl)
				mockDatasource.EXPECT().SaveFile("1", []byte("abc"), time.Duration(0)).Return(errors.New(""))
				mockDatasource.EXPECT().GetFile("1").Return([]byte(""), nil)
				rec := &htmlPdfServiceLogic{
					dsSvc: mockDatasource,
					htSvc: mockHtmlsvc,
				}
				return rec
			},
			validateFunc: func(x *respModel.Response) {
				expected := respModel.Response{
					Status:  http.StatusInternalServerError,
					Message: codes.GetErr(codes.ErrFileStoreFail),
					Data:    nil,
				}
				if !reflect.DeepEqual(x, &expected) {
					t.Errorf("want %v got %v", expected, x)
				}
			},
		},
		{
			name:        "Failure:: Replace:: GetJsonFromHtml fail",
			requestBody: strings.NewReader("abc"),
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
				expected := respModel.Response{
					Status:  http.StatusInternalServerError,
					Message: codes.GetErr(codes.ErrFileConversionFail),
					Data:    nil,
				}
				if !reflect.DeepEqual(x, &expected) {
					t.Errorf("want %v got %v", expected, x)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := tt.setupFunc()
			resp := rec.Replace("1", tt.requestBody.(io.Reader))
			tt.validateFunc(resp)
		})
	}

}

func Test_HtmlToPdf(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	tests := []struct {
		name         string
		requestBody  string
		setupFunc    func() *htmlPdfServiceLogic
		validateFunc func(*respModel.Response)
	}{
		{
			name:        "Success:: HtmlToPdf",
			requestBody: "1",

			setupFunc: func() *htmlPdfServiceLogic {
				buff := bytes.NewBuffer(nil)
				err := json.NewEncoder(buff).Encode(map[string]interface{}{
					"Pages": []interface{}{
						map[string]interface{}{
							"Base64PageData": base64.StdEncoding.EncodeToString([]byte("{{range $index, $element := .Items}}\n<li>{{ $element }}</li>{{ end }}")),
						},
					},
				})
				if err != nil {
					t.Errorf("unable to encode data")
				}
				mockHtmlsvc := mock.NewMockHtmlToPdf(mockCtrl)
				mockHtmlsvc.EXPECT().GeneratePdf(gomock.Any(), gomock.Any()).Return(nil).Times(1).DoAndReturn(func(_ io.Writer, b []byte) error {
					t.Log(string(b))
					type Pages struct {
						P string `json:"Base64PageData"`
					}
					type Data struct {
						Page []Pages `json:"Pages"`
					}
					var data Data
					err = json.Unmarshal(b, &data)
					if err != nil {
						return err
					}
					decodedB, err := base64.StdEncoding.DecodeString(data.Page[0].P)
					if err != nil {
						return err
					}
					t.Log(string(decodedB))
					if !reflect.DeepEqual(decodedB, []byte(`
<li>{Bread 24}</li>
<li>{Rice 56.7}</li>
<li>{Clothes 150.45}</li>
<li>{Water 100}</li>
<li>{Gas 100}</li>`)) {
						t.Errorf("want %v got %v", `
<li>{Bread 24}</li>
<li>{Rice 56.7}</li>
<li>{Clothes 150.45}</li>
<li>{Water 100}</li>
<li>{Gas 100}</li>`, string(decodedB))
					}
					return nil
				})
				mockDatasource := mock.NewMockDataSource(mockCtrl)
				mockDatasource.EXPECT().GetFile("1").Return(buff.Bytes(), nil) //PR COMMENT STILL IN PROGRESS
				rec := &htmlPdfServiceLogic{
					dsSvc: mockDatasource,
					htSvc: mockHtmlsvc,
				}
				return rec
			},
			validateFunc: func(x *respModel.Response) {
				expected := respModel.Response{
					Status: http.StatusOK,
				}
				if !reflect.DeepEqual(x, &expected) {
					t.Errorf("want %v got %v", expected, x)
				}
			},
		},
		{
			name:        "Failure:: HtmlToPdf:: GetFile fail",
			requestBody: "1",
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
				expected := respModel.Response{
					Status:  http.StatusInternalServerError,
					Message: codes.GetErr(codes.ErrFetchingFile),
					Data:    nil,
				}
				if !reflect.DeepEqual(x, &expected) {
					t.Errorf("want %v got %v", expected, x)
				}
			},
		},
		{
			name:        "Failure:: HtmlToPdf:: err unmarshalling json",
			requestBody: "1",
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
				expected := respModel.Response{
					Status:  http.StatusInternalServerError,
					Message: codes.GetErr(codes.ErrFileParseFail),
					Data:    nil,
				}
				if !reflect.DeepEqual(x, &expected) {
					t.Errorf("want %v got %v", expected, x)
				}
			},
		},
		{
			name:        "Failure:: HtmlToPdf:: failed to decode base64 data",
			requestBody: "1",
			setupFunc: func() *htmlPdfServiceLogic {
				js, _ := json.Marshal(map[string]interface{}{
					"Pages": []interface{}{
						map[string]interface{}{
							"Base64PageData": base32.StdEncoding.EncodeToString([]byte("abc")),
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
				expected := respModel.Response{
					Status:  http.StatusInternalServerError,
					Message: codes.GetErr(codes.ErrDecodingData),
					Data:    nil,
				}
				if !reflect.DeepEqual(x, &expected) {
					t.Errorf("want %v got %v", expected, x)
				}
			},
		},
		{
			name:        "Failure:: HtmlToPdf:: assertion for Base64PageData failed",
			requestBody: "1",
			setupFunc: func() *htmlPdfServiceLogic {
				buff := bytes.NewBuffer(nil)
				json.NewEncoder(buff).Encode(map[string]interface{}{
					"Custom-Field": "hello",
					"Pages": []interface{}{
						map[string]interface{}{
							"Base64PageData": map[string]interface{}{},
						},
					},
				})
				mockHtmlsvc := mock.NewMockHtmlToPdf(mockCtrl)
				mockHtmlsvc.EXPECT().GeneratePdf(gomock.Any(), buff.Bytes()).Return(errors.New("Base64PageData is empty")).Times(1)
				mockDatasource := mock.NewMockDataSource(mockCtrl)
				mockDatasource.EXPECT().GetFile("1").Return(buff.Bytes(), nil)
				rec := &htmlPdfServiceLogic{
					dsSvc: mockDatasource,
					htSvc: mockHtmlsvc,
				}
				return rec
			},
			validateFunc: func(x *respModel.Response) {
				expected := respModel.Response{
					Status:  http.StatusInternalServerError,
					Message: codes.GetErr(codes.ErrConvertingToPdf),
					Data:    nil,
				}
				if !reflect.DeepEqual(x, &expected) {
					t.Errorf("want %v got %v", expected, x)
				}
			},
		},
		{
			name:        "Failure:: HtmlToPdf:: assertion for map[string]interface{} failed",
			requestBody: "1",
			setupFunc: func() *htmlPdfServiceLogic {
				js, _ := json.Marshal(map[string]interface{}{
					"Custom-Field": "hello",
					"Pages": []interface{}{
						"hello-world",
						map[string]interface{}{
							"Custom-Data": "world",
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
				expected := respModel.Response{
					Status:  http.StatusBadRequest,
					Message: codes.GetErr(codes.ErrDecodingData),
					Data:    nil,
				}
				if !reflect.DeepEqual(x, &expected) {
					t.Errorf("want %v got %v", expected, x)
				}
			},
		},
		{
			name:        "Failure:: HtmlToPdf:: assertion for Pages failed",
			requestBody: "1",
			setupFunc: func() *htmlPdfServiceLogic {
				js, _ := json.Marshal(map[string]interface{}{
					"Custom-Field": "hello",
				},
				)
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
				expected := respModel.Response{
					Status:  http.StatusBadRequest,
					Message: codes.GetErr(codes.ErrDecodingData),
					Data:    nil,
				}
				if !reflect.DeepEqual(x, &expected) {
					t.Errorf("want %v got %v", expected, x)
				}
			},
		},
		{
			name:        "Failure:: HtmlToPdf:: failed to create new template ",
			requestBody: "1",
			setupFunc: func() *htmlPdfServiceLogic {
				js, err := json.Marshal(map[string]interface{}{
					"Pages": []interface{}{
						map[string]interface{}{
							"Base64PageData": base64.StdEncoding.EncodeToString([]byte("{{ if le .Marks  50 }}")),
						},
					},
				})
				if err != nil {
					log.Error(err.Error())
				}
				mockHtmlsvc := mock.NewMockHtmlToPdf(mockCtrl)
				mockDatasource := mock.NewMockDataSource(mockCtrl)
				mockDatasource.EXPECT().GetFile(gomock.Any()).Return(js, nil)
				rec := &htmlPdfServiceLogic{
					dsSvc: mockDatasource,
					htSvc: mockHtmlsvc,
				}
				return rec
			},
			validateFunc: func(x *respModel.Response) {
				expected := respModel.Response{
					Status:  http.StatusInternalServerError,
					Message: codes.GetErr(codes.ErrFileParseFail),
					Data:    nil,
				}
				if !reflect.DeepEqual(x, &expected) {
					t.Errorf("want %v got %v", expected, x)
				}
			},
		},
		{
			name:        "Failure:: HtmlToPdf:: failed to execute template ",
			requestBody: "1",
			setupFunc: func() *htmlPdfServiceLogic {
				js, err := json.Marshal(map[string]interface{}{
					"Pages": []interface{}{
						map[string]interface{}{
							"Base64PageData": base64.StdEncoding.EncodeToString([]byte("{{ if le .Marks  50 }}{{ end }}")),
						},
					},
				})
				if err != nil {
					log.Error(err.Error())
				}
				mockHtmlsvc := mock.NewMockHtmlToPdf(mockCtrl)
				mockDatasource := mock.NewMockDataSource(mockCtrl)
				mockDatasource.EXPECT().GetFile(gomock.Any()).Return(js, nil)
				rec := &htmlPdfServiceLogic{
					dsSvc: mockDatasource,
					htSvc: mockHtmlsvc,
				}
				return rec
			},
			validateFunc: func(x *respModel.Response) {
				expected := respModel.Response{
					Status:  http.StatusInternalServerError,
					Message: codes.GetErr(codes.ErrFileStoreFail),
					Data:    nil,
				}
				if !reflect.DeepEqual(x, &expected) {
					t.Errorf("want %v got %v", expected, x)
				}
			},
		},
		{
			name:        "Failure:: HtmlToPdf:: failed to generate pdf",
			requestBody: "1",
			setupFunc: func() *htmlPdfServiceLogic {
				js, _ := json.Marshal(map[string]interface{}{
					"Pages": []interface{}{
						map[string]interface{}{
							"Base64PageData": base64.StdEncoding.EncodeToString([]byte("abc")),
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
				expected := respModel.Response{
					Status:  http.StatusInternalServerError,
					Message: codes.GetErr(codes.ErrConvertingToPdf),
					Data:    nil,
				}
				if !reflect.DeepEqual(x, &expected) {
					t.Errorf("want %v got %v", expected, x)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := tt.setupFunc()
			type Cart struct {
				Item   string
				Amount float64
			}
			value := map[string]any{
				"Items": []Cart{
					{
						Item:   "Bread",
						Amount: 24,
					},
					{
						Item:   "Rice",
						Amount: 56.7,
					},
					{
						Item:   "Clothes",
						Amount: 150.45,
					},
					{
						Item:   "Water",
						Amount: 100,
					},
					{
						Item:   "Gas",
						Amount: 100.00,
					},
				},
				"Title": "Inventory list",
			}
			resp := rec.HtmlToPdf(nil, &model.GenerateReq{

				Values: value,
				Id:     tt.requestBody,
			})
			tt.validateFunc(resp)
		})
	}
}
