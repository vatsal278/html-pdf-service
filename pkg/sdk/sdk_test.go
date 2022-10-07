package sdk

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PereRohit/util/log"
	"github.com/PereRohit/util/model"
	"github.com/PereRohit/util/response"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/vatsal278/html-pdf-service/internal/codes"
	modelV "github.com/vatsal278/html-pdf-service/internal/model"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

func testServer(url string, method string, f func(w http.ResponseWriter, r *http.Request)) *httptest.Server {
	router := mux.NewRouter()
	router.HandleFunc(url, f).Methods(method)
	svr := httptest.NewServer(router)
	return svr
}

func Test_Register(t *testing.T) {
	tests := []struct {
		name              string
		filePath          string
		setupFunc         func() *httptest.Server
		mockServerHandler func(w http.ResponseWriter, r *http.Request)
		ValidateFunc      func(id string, err error)
		cleanupFunc       func(*httptest.Server)
		expectedResponse  model.Response
	}{
		{
			name:     "Success:: Register",
			filePath: "./../../docs/Failure.html",
			setupFunc: func() *httptest.Server {
				svr := testServer("/v1/register", http.MethodPost, func(w http.ResponseWriter, r *http.Request) {
					err := r.ParseMultipartForm(10000) //File size to come from config
					if err != nil {
						response.ToJson(w, http.StatusBadRequest, codes.GetErr(codes.ErrFileSizeExceeded), nil)
						log.Error(err.Error())
						return
					}
					file, _, err := r.FormFile("file")
					if err != nil {
						response.ToJson(w, http.StatusBadRequest, codes.GetErr(codes.ErrFileParseFail), nil)
						log.Error(err.Error())
						return
					}
					defer file.Close()

					response.ToJson(w, http.StatusCreated, "SUCCESS", map[string]interface{}{
						"id": "1",
					})
				})
				return svr
			},
			ValidateFunc: func(id string, err error) {
				if err != nil {
					t.Errorf("Want: %v, Got: %v", nil, err.Error())
				}
				if id != "1" {
					t.Errorf("Want: %v, Got: %v", "not nil", "")
				}
			},
			cleanupFunc: func(svr *httptest.Server) {
				svr.Close()
			},
		},
		{
			name:     "Failure:: Register :: incorrect status code received",
			filePath: "./../../docs/Failure.html",
			setupFunc: func() *httptest.Server {
				svr := testServer("/v1/register", http.MethodPost, func(w http.ResponseWriter, r *http.Request) {
					response.ToJson(w, http.StatusBadRequest, "Failure", nil)
				})
				return svr
			},
			ValidateFunc: func(id string, err error) {
				if err.Error() != "non success status code received : 400" {
					t.Errorf("Want: %v, Got: %v", "non success status code received : 400", err.Error())
				}
			},
			cleanupFunc: func(svr *httptest.Server) {
				svr.Close()
			},
		},
		{
			name:     "Failure:: Register :: incorrect test server path", //doubt
			filePath: "./../../go.mod",
			setupFunc: func() *httptest.Server {
				svr := testServer("new", http.MethodPost, func(w http.ResponseWriter, r *http.Request) {
					response.ToJson(w, http.StatusCreated, "SUCCESS", nil)
				})
				return svr
			},
			ValidateFunc: func(id string, err error) {
				if err.Error() != fmt.Errorf("non success status code received : %v", http.StatusNotFound).Error() {
					t.Errorf("Want: %v, Got: %v", fmt.Errorf("non success status code received : %v", http.StatusNotFound), err)
				}
			},
			cleanupFunc: func(svr *httptest.Server) {
				svr.Close()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svr := tt.setupFunc()
			defer tt.cleanupFunc(svr)

			calls := NewHtmlToPdfSvc(svr.URL)
			fileBytes, _ := os.ReadFile(tt.filePath)
			id, err := calls.Register(fileBytes)

			tt.ValidateFunc(id, err)
		})
	}
}
func Test_Replace(t *testing.T) {
	tests := []struct {
		name              string
		filePath          string
		id                string
		setupFunc         func() *httptest.Server
		mockServerHandler func(w http.ResponseWriter, r *http.Request)
		ValidateFunc      func(err error)
		cleanupFunc       func(*httptest.Server)
		expectedResponse  model.Response
	}{
		{
			name:     "Success:: Replace",
			id:       "1",
			filePath: "./../../docs/Failure.html",
			setupFunc: func() *httptest.Server {
				svr := testServer("/v1/register/{id}", http.MethodPut, func(w http.ResponseWriter, r *http.Request) {
					vars := mux.Vars(r)
					id, ok := vars["id"]
					if !ok {
						response.ToJson(w, http.StatusBadRequest, codes.GetErr(codes.ErrIdNeeded), nil)
						return
					}
					err := r.ParseMultipartForm(10000)
					if err != nil {
						response.ToJson(w, http.StatusBadRequest, codes.GetErr(codes.ErrFileSizeExceeded), nil)
						log.Error(err.Error())
						return
					}
					file, _, err := r.FormFile("file")
					if err != nil {
						response.ToJson(w, http.StatusBadRequest, codes.GetErr(codes.ErrFileParseFail), nil)
						log.Error(err.Error())
						return
					}
					defer file.Close()

					response.ToJson(w, http.StatusOK, "SUCCESS", map[string]interface{}{
						"id": id,
					})
				})
				return svr
			},
			ValidateFunc: func(err error) {
				if err != nil {
					t.Errorf("Want: %v, Got: %v", nil, err)
				}
			},
			cleanupFunc: func(svr *httptest.Server) {
				svr.Close()
			},
		},
		{
			name:     "Failure:: Replace :: incorrect status code",
			filePath: "./../../docs/Failure.html",
			setupFunc: func() *httptest.Server {
				svr := testServer("", http.MethodPut, func(w http.ResponseWriter, r *http.Request) {
					err := r.ParseMultipartForm(10000) //File size to come from config
					if err != nil {
						response.ToJson(w, http.StatusBadRequest, codes.GetErr(codes.ErrFileSizeExceeded), nil)
						log.Error(err.Error())
						return
					}
					file, _, err := r.FormFile("file")
					if err != nil {
						response.ToJson(w, http.StatusBadRequest, codes.GetErr(codes.ErrFileParseFail), nil)
						log.Error(err.Error())
						return
					}
					defer file.Close()

					response.ToJson(w, http.StatusCreated, "SUCCESS", map[string]interface{}{
						"id": gomock.Any(),
					})
				})
				return svr
			},
			ValidateFunc: func(err error) {
				if err.Error() != "non success status code received : 404" {
					t.Errorf("Want: %v, Got: %v", "non success status code received : 404", err.Error())
				}
			},
			cleanupFunc: func(svr *httptest.Server) {
				svr.Close()
			},
		},
		{
			name:     "Failure:: Replace :: incorrect file type", //doubt
			filePath: "./../../Dockerfile",
			setupFunc: func() *httptest.Server {
				svr := testServer("/v1/register{id}", http.MethodPut, func(w http.ResponseWriter, r *http.Request) {
					err := r.ParseMultipartForm(10000) //File size to come from config
					if err != nil {
						response.ToJson(w, http.StatusBadRequest, codes.GetErr(codes.ErrFileSizeExceeded), nil)
						log.Error(err.Error())
						return
					}
					file, _, err := r.FormFile("file")
					if err != nil {
						response.ToJson(w, http.StatusBadRequest, codes.GetErr(codes.ErrFileParseFail), nil)
						log.Error(err.Error())
						return
					}
					defer file.Close()

					response.ToJson(w, http.StatusCreated, "SUCCESS", map[string]interface{}{
						"id": gomock.Any(),
					})
				})
				return svr
			},
			ValidateFunc: func(err error) {
				if err.Error() != "non success status code received : 404" {
					t.Errorf("Want: %v, Got: %v", "non success status code received : 404", err.Error())
				}
			},
			cleanupFunc: func(svr *httptest.Server) {
				svr.Close()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svr := tt.setupFunc()
			defer tt.cleanupFunc(svr)

			calls := NewHtmlToPdfSvc(svr.URL)
			fileBytes, err := os.ReadFile(tt.filePath)

			err = calls.Replace(fileBytes, tt.id)

			tt.ValidateFunc(err)
		})
	}
}
func Test_GeneratePdf(t *testing.T) {
	tests := []struct {
		name              string
		id                string
		data              map[string]interface{}
		setupFunc         func() *httptest.Server
		mockServerHandler func(w http.ResponseWriter, r *http.Request)
		ValidateFunc      func(err error)
		cleanupFunc       func(*httptest.Server)
		expectedResponse  model.Response
	}{
		{
			name: "Success:: GeneratePdf",
			id:   "1",
			data: map[string]interface{}{"id": "1"},
			setupFunc: func() *httptest.Server {
				svr := testServer("/v1/generate/{id}", http.MethodPost, func(w http.ResponseWriter, r *http.Request) {
					vars := mux.Vars(r)
					//we take id as a parameter from url path
					id, ok := vars["id"]
					if !ok {
						response.ToJson(w, http.StatusBadRequest, codes.GetErr(codes.ErrIdNeeded), nil)
						return
					}
					var data modelV.GenerateReq
					err := json.NewDecoder(r.Body).Decode(&data)
					if err != nil {
						response.ToJson(w, http.StatusBadRequest, codes.GetErr(codes.ErrDecodingData), nil)
						log.Error(err.Error())
						return
					}
					data.Id = id

					response.ToJson(w, http.StatusOK, "", nil)
				})
				return svr
			},
			ValidateFunc: func(err error) {
				if err != nil {
					t.Errorf("Want: %v, Got: %v", nil, err.Error())
				}
			},
			cleanupFunc: func(svr *httptest.Server) {
				svr.Close()
			},
		},
		{
			name: "Failure:: GeneratePdf :: id required",
			id:   "",
			data: map[string]interface{}{"id": "1"},
			setupFunc: func() *httptest.Server {
				svr := testServer("/v1/generate/{id}", http.MethodPut, func(w http.ResponseWriter, r *http.Request) {
					vars := mux.Vars(r)
					//we take id as a parameter from url path
					id, ok := vars["id"]
					if !ok {
						log.Error("here")
						response.ToJson(w, http.StatusBadRequest, codes.GetErr(codes.ErrIdNeeded), nil)
						return
					}
					var data modelV.GenerateReq
					err := json.NewDecoder(r.Body).Decode(&data)
					if err != nil {
						response.ToJson(w, http.StatusBadRequest, codes.GetErr(codes.ErrDecodingData), nil)
						log.Error(err.Error())
						return
					}
					data.Id = id

					response.ToJson(w, http.StatusOK, "", nil)
				})
				return svr
			},
			ValidateFunc: func(err error) {
				if err.Error() != errors.New("non success status code received : 404").Error() {
					t.Errorf("Want: %v, Got: %v", "non success status code received : 404", err)
				}
			},
			cleanupFunc: func(svr *httptest.Server) {
				svr.Close()
			},
		},
		{
			name: "Failure:: GeneratePdf :: incorrect status code",
			id:   "1",
			data: map[string]interface{}{"id": "1"},
			setupFunc: func() *httptest.Server {
				svr := testServer("/v1/generate/{id}", http.MethodPost, func(w http.ResponseWriter, r *http.Request) {
					vars := mux.Vars(r)
					//we take id as a parameter from url path
					id, ok := vars["id"]
					if !ok {
						log.Error("here")
						response.ToJson(w, http.StatusBadRequest, codes.GetErr(codes.ErrIdNeeded), nil)
						return
					}
					var data modelV.GenerateReq
					err := json.NewDecoder(r.Body).Decode(&data)
					if err != nil {
						response.ToJson(w, http.StatusBadRequest, codes.GetErr(codes.ErrDecodingData), nil)
						log.Error(err.Error())
						return
					}
					data.Id = id

					response.ToJson(w, http.StatusNotFound, "", nil)
				})
				return svr
			},
			ValidateFunc: func(err error) {
				if !reflect.DeepEqual(err.Error(), errors.New("non success status code received : 404").Error()) {
					t.Errorf("Want: %v, Got: %v", "non success status code received : 404", err)
				}
			},
			cleanupFunc: func(svr *httptest.Server) {
				svr.Close()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svr := tt.setupFunc()
			defer tt.cleanupFunc(svr)

			calls := NewHtmlToPdfSvc(svr.URL)
			_, err := calls.GeneratePdf(tt.data, tt.id)

			tt.ValidateFunc(err)
		})
	}
}
