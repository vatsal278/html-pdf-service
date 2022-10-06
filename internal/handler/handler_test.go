package handler

import (
	"bytes"
	"encoding/json"
	respModel "github.com/PereRohit/util/model"
	"github.com/PereRohit/util/testutil"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/vatsal278/html-pdf-service/internal/codes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vatsal278/html-pdf-service/internal/model"
	"github.com/vatsal278/html-pdf-service/internal/repo/datasource"
	"github.com/vatsal278/html-pdf-service/pkg/mock"
)

func Test_htmlPdfService_Ping(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tests := []struct {
		name     string
		setup    func() (HtmlPdfServiceHandler, http.ResponseWriter, *http.Request)
		validate func(w http.ResponseWriter)
	}{
		{
			name: "Success",
			setup: func() (HtmlPdfServiceHandler, http.ResponseWriter, *http.Request) {
				mockLogic := mock.NewMockHtmlPdfServiceLogicIer(mockCtrl)

				req := &model.PingRequest{
					Data: "hello-world",
				}

				mockLogic.EXPECT().Ping(req).Return(&respModel.Response{
					Status:  http.StatusOK,
					Message: "Ok",
					Data:    "pong",
				}).Times(1)

				rec := &htmlPdfService{
					logic: mockLogic,
				}

				reqB, err := json.Marshal(req)
				if err != nil {
					t.Errorf(err.Error())
					t.Fail()
				}
				r := httptest.NewRequest(http.MethodGet, "https://ping", bytes.NewReader(reqB))
				w := httptest.NewRecorder()

				return rec, w, r
			},
			validate: func(w http.ResponseWriter) {
				wIn := w.(*httptest.ResponseRecorder)

				diff := testutil.Diff(wIn.Code, http.StatusOK)
				if diff != "" {
					t.Error(testutil.Callers(), diff)
				}

				diff = testutil.Diff(wIn.Header().Get("Content-Type"), "application/json")
				if diff != "" {
					t.Error(testutil.Callers(), diff)
				}

				resp := respModel.Response{}
				err := json.NewDecoder(wIn.Body).Decode(&resp)
				diff = testutil.Diff(err, nil)
				if diff != "" {
					t.Error(testutil.Callers(), diff)
				}

				diff = testutil.Diff(resp, respModel.Response{
					Status:  http.StatusOK,
					Message: "Ok",
					Data:    "pong",
				})
				if diff != "" {
					t.Error(testutil.Callers(), diff)
				}
			},
		},
		{
			name: "Failure::request not as expected",
			setup: func() (HtmlPdfServiceHandler, http.ResponseWriter, *http.Request) {
				req := "hello world"

				rec := &htmlPdfService{}

				r := httptest.NewRequest(http.MethodGet, "https://ping", bytes.NewReader([]byte(req)))
				w := httptest.NewRecorder()

				return rec, w, r
			},
			validate: func(w http.ResponseWriter) {
				wIn := w.(*httptest.ResponseRecorder)

				diff := testutil.Diff(wIn.Code, http.StatusBadRequest)
				if diff != "" {
					t.Error(testutil.Callers(), diff)
				}

				diff = testutil.Diff(wIn.Header().Get("Content-Type"), "application/json")
				if diff != "" {
					t.Error(testutil.Callers(), diff)
				}

				resp := respModel.Response{}
				err := json.NewDecoder(wIn.Body).Decode(&resp)
				diff = testutil.Diff(err, nil)
				if diff != "" {
					t.Error(testutil.Callers(), diff)
				}

				// ignore specific message
				resp.Message = ""

				diff = testutil.Diff(resp, respModel.Response{
					Status:  http.StatusBadRequest,
					Message: "",
					Data:    nil,
				})
				if diff != "" {
					t.Error(testutil.Callers(), diff)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			receiver, w, r := tt.setup()

			receiver.Ping(w, r)

			tt.validate(w)
		})
	}
}

func Test_htmlPdfService_HealthCheck(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tests := []struct {
		name        string
		setup       func() HtmlPdfServiceHandler
		wantSvcName string
		wantMsg     string
		wantStat    bool
	}{
		{
			name: "Success",
			setup: func() HtmlPdfServiceHandler {
				mockLogic := mock.NewMockHtmlPdfServiceLogicIer(mockCtrl)

				mockLogic.EXPECT().HealthCheck().
					Return(true).Times(1)

				rec := &htmlPdfService{
					logic: mockLogic,
				}

				return rec
			},
			wantSvcName: HtmlPdfServiceName,
			wantMsg:     "",
			wantStat:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			receiver := tt.setup()

			svcName, msg, stat := receiver.HealthCheck()

			diff := testutil.Diff(svcName, tt.wantSvcName)
			if diff != "" {
				t.Error(testutil.Callers(), diff)
			}

			diff = testutil.Diff(msg, tt.wantMsg)
			if diff != "" {
				t.Error(testutil.Callers(), diff)
			}

			diff = testutil.Diff(stat, tt.wantStat)
			if diff != "" {
				t.Error(testutil.Callers(), diff)
			}
		})
	}
}

func TestNewHtmlPdfService(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tests := []struct {
		name     string
		setup    func() (datasource.DataSource, *mock.MockHtmlToPdf)
		wantStat bool
	}{
		{
			name: "Success",
			setup: func() (datasource.DataSource, *mock.MockHtmlToPdf) {
				mockDs := mock.NewMockDataSource(mockCtrl)

				mockDs.EXPECT().HealthCheck().Times(1).
					Return(false)
				mockHt := mock.NewMockHtmlToPdf(mockCtrl)

				return mockDs, mockHt
			},
			wantStat: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds, ht := tt.setup()
			rec := NewHtmlPdfService(ds, ht, 10204)

			_, _, stat := rec.HealthCheck()

			diff := testutil.Diff(stat, tt.wantStat)
			if diff != "" {
				t.Error(testutil.Callers(), diff)
			}
		})
	}
}
func TestUpload(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	tests := []struct {
		name         string
		requestBody  string
		setupFunc    func() (*http.Request, *htmlPdfService)
		validateFunc func(*httptest.ResponseRecorder)
	}{
		{
			name: "Success:: Upload",
			setupFunc: func() (*http.Request, *htmlPdfService) {
				b := new(bytes.Buffer)
				y := multipart.NewWriter(b)
				part, err := y.CreateFormFile("file", "some-file")
				if err != nil {
					return nil, nil
				}
				_, err = part.Write([]byte("abc"))
				if err != nil {
					return nil, nil
				}
				y.Close()
				r := httptest.NewRequest(http.MethodPost, "/v1/register", b)
				r.Header.Set("Content-Type", y.FormDataContentType())
				mockLogicier := mock.NewMockHtmlPdfServiceLogicIer(mockCtrl)
				mockLogicier.EXPECT().Upload(gomock.Any()).Times(1).
					DoAndReturn(func(f io.Reader) *respModel.Response {
						gotData, err := ioutil.ReadAll(f)
						if err != nil {
							t.Error(err)
							t.FailNow()
						}
						diff := testutil.Diff(gotData, []byte("abc"))
						if diff != "" {
							t.Error(testutil.Callers(), diff)
						}
						return &respModel.Response{
							Status:  http.StatusCreated,
							Message: "SUCCESS",
							Data: map[string]interface{}{
								"id": "1",
							},
						}
					})
				rec := &htmlPdfService{
					logic: mockLogicier,
				}
				return r, rec
			},
			validateFunc: func(x *httptest.ResponseRecorder) {
				if x.Code != http.StatusCreated {
					t.Errorf("want %v got %v", http.StatusCreated, x)
				}
				var r respModel.Response
				got, err := io.ReadAll(x.Body)
				if err != nil {
					t.Errorf(err.Error())
				}
				err = json.Unmarshal(got, &r)
				if err != nil {
					t.Error(err.Error())
				}
				diff := testutil.Diff(r, respModel.Response{
					Status:  http.StatusCreated,
					Message: "SUCCESS",
					Data: map[string]interface{}{
						"id": "1",
					},
				})
				if diff != "" {
					t.Error(testutil.Callers(), diff)
				}
				diff = testutil.Diff(err, nil)
				if diff != "" {
					t.Error(testutil.Callers(), diff)
				}
			},
		},
		{
			name: "Failure:: Upload:: ParseMultiForm failure",
			setupFunc: func() (*http.Request, *htmlPdfService) {
				b := new(bytes.Buffer)
				r := httptest.NewRequest(http.MethodPost, "/v1/register", b)
				mockLogicier := mock.NewMockHtmlPdfServiceLogicIer(mockCtrl)
				rec := &htmlPdfService{
					logic: mockLogicier,
				}
				return r, rec
			},
			validateFunc: func(x *httptest.ResponseRecorder) {
				if x.Code != http.StatusBadRequest {
					t.Errorf("want %v got %v", http.StatusBadRequest, x.Code)
				}
				var r respModel.Response
				got, err := io.ReadAll(x.Body)
				if err != nil {
					t.Errorf(err.Error())
				}
				err = json.Unmarshal(got, &r)
				if err != nil {
					t.Errorf(err.Error())
				}
				diff := testutil.Diff(r, respModel.Response{
					Status:  http.StatusBadRequest,
					Message: "request Content-Type isn't multipart/form-data",
					Data:    nil,
				})
				if diff != "" {
					t.Error(testutil.Callers(), diff)
				}
				diff = testutil.Diff(err, nil)
				if diff != "" {
					t.Error(testutil.Callers(), diff)
				}
			},
		},
		{
			name: "Failure:: Upload:: no file found",
			setupFunc: func() (*http.Request, *htmlPdfService) {
				b := new(bytes.Buffer)
				y := multipart.NewWriter(b)
				part, _ := y.CreateFormFile("f", "some-file")
				_, err := part.Write([]byte("abc"))
				if err != nil {
					return nil, nil
				}
				y.Close()
				r := httptest.NewRequest(http.MethodPost, "/v1/register", b)
				r.Header.Set("Content-Type", y.FormDataContentType())
				mockLogicier := mock.NewMockHtmlPdfServiceLogicIer(mockCtrl)
				rec := &htmlPdfService{
					logic: mockLogicier,
				}
				return r, rec
			},
			validateFunc: func(x *httptest.ResponseRecorder) {
				if x.Code != http.StatusBadRequest {
					t.Errorf("want %v got %v", http.StatusBadRequest, x.Code)
				}
				var r respModel.Response
				got, err := io.ReadAll(x.Body)
				if err != nil {
					t.Errorf(err.Error())
				}
				err = json.Unmarshal(got, &r)
				if err != nil {
					t.Errorf(err.Error())
				}
				diff := testutil.Diff(r, respModel.Response{
					Status:  http.StatusBadRequest,
					Message: codes.GetErr(codes.ErrFileParseFail),
					Data:    nil,
				})
				if diff != "" {
					t.Error(testutil.Callers(), diff)
				}
				diff = testutil.Diff(err, nil)
				if diff != "" {
					t.Error(testutil.Callers(), diff)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, rec := tt.setupFunc()
			w := httptest.NewRecorder()
			x := rec.Upload
			x(w, r)
			tt.validateFunc(w)
		})
	}
}

func TestConvertToPdf(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	tests := []struct {
		name         string
		requestBody  string
		setupFunc    func() (*http.Request, *htmlPdfService)
		validateFunc func(*httptest.ResponseRecorder)
	}{
		{
			name:        "Success:: ConvertToPdf",
			requestBody: "1",
			setupFunc: func() (*http.Request, *htmlPdfService) {
				var temp = model.GenerateReq{
					Values: map[string]interface{}{
						"Name":  "vatsal",
						"Marks": "90",
						"ID":    "1",
					}}
				temp.Id = "1"
				b, err := json.Marshal(temp)
				if err != nil {
					t.Error(err)
					return nil, nil
				}
				r := httptest.NewRequest(http.MethodPost, "/v1/generate/1", bytes.NewBuffer(b))
				r = mux.SetURLVars(r, map[string]string{"id": "1"})
				mockLogicier := mock.NewMockHtmlPdfServiceLogicIer(mockCtrl)
				mockLogicier.EXPECT().HtmlToPdf(gomock.Any(), &temp).Times(1).
					DoAndReturn(func(w io.Writer, req *model.GenerateReq) *respModel.Response {
						_, err = w.Write([]byte("hello-world"))
						if err != nil {
							t.Errorf(err.Error())
						}
						return &respModel.Response{
							Status: http.StatusOK,
						}
					})
				rec := &htmlPdfService{
					logic: mockLogicier,
				}
				return r, rec
			},
			validateFunc: func(x *httptest.ResponseRecorder) {
				if x.Code != http.StatusOK {
					t.Errorf("want %v got %v", http.StatusOK, x.Code)
				}
				var r respModel.Response
				got, err := io.ReadAll(x.Body)
				if err != nil {
					t.Errorf(err.Error())
				}
				if err != nil {
					t.Error(err)
					return
				}
				err = json.Unmarshal(got, &r)
				if err != nil {
					t.Error(err)
					return
				}
				diff := testutil.Diff(got, []byte("hello-world"))
				if diff != "" {
					t.Error(testutil.Callers(), diff)
				}
				diff = testutil.Diff(err, nil)
				if diff != "" {
					t.Error(testutil.Callers(), diff)
				}
			},
		},
		{
			name:        "Failure:: ConvertToPdf: id not found",
			requestBody: "1",
			setupFunc: func() (*http.Request, *htmlPdfService) {
				rec := &htmlPdfService{
					logic: nil,
				}
				return httptest.NewRequest(http.MethodPost, "/v1/generate", nil), rec
			},
			validateFunc: func(x *httptest.ResponseRecorder) {
				if x.Code != http.StatusBadRequest {
					t.Errorf("want %v got %v", http.StatusBadRequest, x.Code)
				}
				var r respModel.Response
				got, err := io.ReadAll(x.Body)
				if err != nil {
					t.Errorf(err.Error())
				}
				err = json.Unmarshal(got, &r)
				if err != nil {
					t.Error(err)
					return
				}
				diff := testutil.Diff(r, respModel.Response{
					Status:  http.StatusBadRequest,
					Message: codes.GetErr(codes.ErrIdNeeded),
					Data:    nil,
				})
				if diff != "" {
					t.Error(testutil.Callers(), diff)
				}
				diff = testutil.Diff(err, nil)
				if diff != "" {
					t.Error(testutil.Callers(), diff)
				}
			},
		},
		{
			name:        "Failure:: ConvertToPdf: json failure",
			requestBody: "1",
			setupFunc: func() (*http.Request, *htmlPdfService) {
				rec := &htmlPdfService{
					logic: nil,
				}
				r := httptest.NewRequest(http.MethodPost, "/v1/generate", nil)
				r = mux.SetURLVars(r, map[string]string{"id": "1"})
				return r, rec
			},
			validateFunc: func(x *httptest.ResponseRecorder) {
				if x.Code != http.StatusBadRequest {
					t.Errorf("want %v got %v", http.StatusBadRequest, x.Code)
				}
				var r respModel.Response
				got, err := io.ReadAll(x.Body)
				if err != nil {
					t.Errorf(err.Error())
				}
				err = json.Unmarshal(got, &r)
				if err != nil {
					t.Error(err)
					return
				}
				diff := testutil.Diff(r, respModel.Response{
					Status:  http.StatusBadRequest,
					Message: codes.GetErr(codes.ErrDecodingData),
					Data:    nil,
				})
				if diff != "" {
					t.Error(testutil.Callers(), diff)
				}
				diff = testutil.Diff(err, nil)
				if diff != "" {
					t.Error(testutil.Callers(), diff)
				}
			},
		},
		{
			name:        "Failure:: ConvertToPdf :: wkhtmltopdf failure",
			requestBody: "1",
			setupFunc: func() (*http.Request, *htmlPdfService) {
				var temp struct {
					Values struct {
						Name  string `json:"name"`
						Marks int    `json:"marks"`
						ID    string `json:"id"`
					} `json:"values"`
				}
				temp.Values.Name = "vatsal"
				temp.Values.Marks = 90
				b, err := json.Marshal(temp)
				if err != nil {
					t.Error(err)
				}
				r := httptest.NewRequest(http.MethodPost, "/v1/generate/1", bytes.NewBuffer(b))
				r = mux.SetURLVars(r, map[string]string{"id": "1"})
				mockLogicier := mock.NewMockHtmlPdfServiceLogicIer(mockCtrl)
				mockLogicier.EXPECT().HtmlToPdf(gomock.Any(), gomock.Any()).Times(1).
					DoAndReturn(func(w io.Writer, req *model.GenerateReq) *respModel.Response {
						_, err = w.Write([]byte("hello-world"))
						return &respModel.Response{
							Status:  http.StatusInternalServerError,
							Message: codes.GetErr(codes.ErrConvertingToPdf),
							Data:    nil,
						}
					})
				rec := &htmlPdfService{
					logic: mockLogicier,
				}
				return r, rec
			},
			validateFunc: func(x *httptest.ResponseRecorder) {
				var r respModel.Response
				got, err := io.ReadAll(x.Body)
				if err != nil {
					t.Errorf(err.Error())
				}
				err = json.Unmarshal(got, &r)
				if err != nil {
					t.Errorf(err.Error())
				}
				diff := testutil.Diff(r, respModel.Response{
					Status:  http.StatusInternalServerError,
					Message: codes.GetErr(codes.ErrConvertingToPdf),
					Data:    nil,
				})
				if diff != "" {
					t.Error(testutil.Callers(), diff)
				}
				diff = testutil.Diff(err, nil)
				if diff != "" {
					t.Error(testutil.Callers(), diff)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, rec := tt.setupFunc()
			w := httptest.NewRecorder()
			x := rec.ConvertToPdf
			x(w, r)
			tt.validateFunc(w)
		})
	}
}

func TestReplace(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	tests := []struct {
		name         string
		requestBody  string
		setupFunc    func() (*http.Request, *htmlPdfService)
		validateFunc func(*httptest.ResponseRecorder)
	}{
		{
			name:        "Success:: Replace",
			requestBody: "1",
			setupFunc: func() (*http.Request, *htmlPdfService) {
				b := new(bytes.Buffer)
				y := multipart.NewWriter(b)
				part, err := y.CreateFormFile("file", "some-file")
				if err != nil {
					t.Error(err)
				}
				_, err = part.Write([]byte("abc"))
				if err != nil {
					t.Errorf(err.Error())
				}
				y.Close()
				mockLogicier := mock.NewMockHtmlPdfServiceLogicIer(mockCtrl)
				mockLogicier.EXPECT().Replace(gomock.Any(), gomock.Any()).Times(1).
					DoAndReturn(func(id string, f io.Reader) *respModel.Response {
						gotData, err := ioutil.ReadAll(f)
						if err != nil {
							t.Errorf(err.Error())
						}
						diff := testutil.Diff(gotData, []byte("abc"))
						if diff != "" {
							t.Error(testutil.Callers(), diff)
						}
						return &respModel.Response{
							Status:  http.StatusOK,
							Message: "SUCCESS",
							Data: map[string]interface{}{
								"id": id,
							},
						}
					})
				rec := &htmlPdfService{
					logic: mockLogicier,
				}
				r := httptest.NewRequest(http.MethodPut, "/v1/register/1", b)
				r = mux.SetURLVars(r, map[string]string{"id": "1"})
				r.Header.Set("Content-Type", y.FormDataContentType())
				return r, rec
			},
			validateFunc: func(x *httptest.ResponseRecorder) {
				if x.Code != http.StatusOK {
					t.Errorf("want %v got %v", http.StatusOK, x.Code)
				}
				var r respModel.Response
				got, err := io.ReadAll(x.Body)
				if err != nil {
					t.Errorf(err.Error())
				}
				err = json.Unmarshal(got, &r)
				if err != nil {
					t.Errorf(err.Error())
				}
				diff := testutil.Diff(r, respModel.Response{
					Status:  http.StatusOK,
					Message: "SUCCESS",
					Data: map[string]interface{}{
						"id": "1",
					},
				})
				if diff != "" {
					t.Error(testutil.Callers(), diff)
				}
				diff = testutil.Diff(err, nil)
				if diff != "" {
					t.Error(testutil.Callers(), diff)
				}
			},
		},
		{
			name:        "Failure:: Replace:: id not found",
			requestBody: "1",
			setupFunc: func() (*http.Request, *htmlPdfService) {
				b := new(bytes.Buffer)
				mockLogicier := mock.NewMockHtmlPdfServiceLogicIer(mockCtrl)
				rec := &htmlPdfService{
					logic: mockLogicier,
				}
				return httptest.NewRequest(http.MethodPut, "/v1/register/", b), rec
			},
			validateFunc: func(x *httptest.ResponseRecorder) {
				if x.Code != http.StatusBadRequest {
					t.Errorf("want %v got %v", http.StatusBadRequest, x.Code)
				}
				var r respModel.Response
				got, err := io.ReadAll(x.Body)
				if err != nil {
					t.Errorf(err.Error())
				}
				err = json.Unmarshal(got, &r)
				if err != nil {
					t.Errorf(err.Error())
				}
				diff := testutil.Diff(r, respModel.Response{
					Status:  http.StatusBadRequest,
					Message: codes.GetErr(codes.ErrIdNeeded),
					Data:    nil,
				})
				if diff != "" {
					t.Error(testutil.Callers(), diff)
				}
				diff = testutil.Diff(err, nil)
				if diff != "" {
					t.Error(testutil.Callers(), diff)
				}
			},
		},
		{
			name:        "Failure:: Replace:: parse multipart failure",
			requestBody: "1",
			setupFunc: func() (*http.Request, *htmlPdfService) {
				b := new(bytes.Buffer)
				y := multipart.NewWriter(b)
				part, _ := y.CreateFormFile("file", "some-file")
				_, err := part.Write([]byte("abc"))
				if err != nil {
					t.Errorf(err.Error())
				}
				y.Close()
				mockLogicier := mock.NewMockHtmlPdfServiceLogicIer(mockCtrl)
				rec := &htmlPdfService{
					logic: mockLogicier,
				}
				r := httptest.NewRequest(http.MethodPut, "/v1/register/", b)
				r = mux.SetURLVars(r, map[string]string{"id": "1"})
				return r, rec
			},
			validateFunc: func(x *httptest.ResponseRecorder) {
				if x.Code != http.StatusBadRequest {
					t.Errorf("want %v got %v", http.StatusBadRequest, x)
				}
				var r respModel.Response
				got, err := io.ReadAll(x.Body)
				if err != nil {
					t.Errorf(err.Error())
				}
				err = json.Unmarshal(got, &r)
				if err != nil {
					t.Errorf(err.Error())
				}
				diff := testutil.Diff(r, respModel.Response{
					Status:  http.StatusBadRequest,
					Message: codes.GetErr(codes.ErrFileSizeExceeded),
					Data:    nil,
				})
				if diff != "" {
					t.Error(testutil.Callers(), diff)
				}
				diff = testutil.Diff(err, nil)
				if diff != "" {
					t.Error(testutil.Callers(), diff)
				}
			},
		}, {
			name:        "Failure:: Replace:: parse multipart failure",
			requestBody: "1",
			setupFunc: func() (*http.Request, *htmlPdfService) {
				b := new(bytes.Buffer)
				y := multipart.NewWriter(b)
				part, _ := y.CreateFormFile("f", "some-file")
				_, err := part.Write([]byte("abc"))
				if err != nil {
					return nil, nil
				}
				y.Close()
				mockLogicier := mock.NewMockHtmlPdfServiceLogicIer(mockCtrl)
				rec := &htmlPdfService{
					logic: mockLogicier,
				}
				r := httptest.NewRequest(http.MethodPut, "/v1/register/", b)
				r = mux.SetURLVars(r, map[string]string{"id": "1"})
				r.Header.Set("Content-Type", y.FormDataContentType())
				return r, rec
			},
			validateFunc: func(x *httptest.ResponseRecorder) {
				if x.Code != http.StatusBadRequest {
					t.Errorf("want %v got %v", http.StatusBadRequest, x)
				}
				var r respModel.Response
				got, err := io.ReadAll(x.Body)
				if err != nil {
					t.Errorf(err.Error())
				}
				err = json.Unmarshal(got, &r)
				if err != nil {
					t.Errorf(err.Error())
				}
				diff := testutil.Diff(r, respModel.Response{
					Status:  http.StatusBadRequest,
					Message: codes.GetErr(codes.ErrFileParseFail),
					Data:    nil,
				})
				if diff != "" {
					t.Error(testutil.Callers(), diff)
				}
				diff = testutil.Diff(err, nil)
				if diff != "" {
					t.Error(testutil.Callers(), diff)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, rec := tt.setupFunc()
			w := httptest.NewRecorder()
			x := rec.ReplaceHtml
			x(w, r)
			tt.validateFunc(w)
		})
	}
}
