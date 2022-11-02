package sdk

import (
	"encoding/json"
	"errors"
	"github.com/PereRohit/util/model"
	"github.com/PereRohit/util/response"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

type GenerateReq struct {
	Values map[string]interface{} `json:"values"`
	Id     string                 `json:"-"`
}

func testServer(url string, method string, f func(w http.ResponseWriter, r *http.Request)) *httptest.Server {
	router := mux.NewRouter()
	router.HandleFunc(url, f).Methods(method)
	svr := httptest.NewServer(router)
	return svr
}

type Reader string

func (Reader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}
func Test_Register(t *testing.T) {
	tests := []struct {
		name              string
		setupFunc         func() *httptest.Server
		mockServerHandler func(w http.ResponseWriter, r *http.Request)
		ValidateFunc      func(id string, err error)
		cleanupFunc       func(*httptest.Server)
		expectedResponse  model.Response
	}{
		{
			name: "Success:: Register",
			setupFunc: func() *httptest.Server {
				svr := testServer("/v1/register", http.MethodPost, func(w http.ResponseWriter, r *http.Request) {
					err := r.ParseMultipartForm(10000) //File size to come from config
					if err != nil {
						response.ToJson(w, http.StatusBadRequest, "file size exceeded", nil)
						t.Error(err.Error())
						return
					}
					file, _, err := r.FormFile("file")
					if err != nil {
						response.ToJson(w, http.StatusBadRequest, "file parse fail", nil)
						t.Error(err.Error())
						return
					}
					fileBytes, err := ioutil.ReadAll(file)
					if err != nil {
						response.ToJson(w, http.StatusBadRequest, "file read fail", nil)
						t.Error(err.Error())
						return
					}
					if !reflect.DeepEqual(fileBytes, []byte("abc")) {
						t.Errorf("Want: %v, Got: %v", []byte("abc"), fileBytes)
					}
					defer file.Close()
					cType := r.Header.Get("Content-Type")
					x := strings.Split(cType, ";")
					if !reflect.DeepEqual(x[0], "multipart/form-data") {
						t.Errorf("Want: %v, Got: %v", "multipart/form-data", x[0])
					}
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
			name: "Failure :: Register :: id not found in response",
			setupFunc: func() *httptest.Server {
				svr := testServer("/v1/register", http.MethodPost, func(w http.ResponseWriter, r *http.Request) {
					response.ToJson(w, http.StatusCreated, "SUCCESS", map[string]interface{}{
						"xyz": "1",
					})
				})
				return svr
			},
			ValidateFunc: func(id string, err error) {
				if err.Error() != errors.New("id not found in response").Error() {
					t.Errorf("Want: %v, Got: %v", "id not found in response", err.Error())
				}
			},
			cleanupFunc: func(svr *httptest.Server) {
				svr.Close()
			},
		},
		{
			name: "Failure:: Register:: failed to assert response data",
			setupFunc: func() *httptest.Server {
				svr := testServer("/v1/register", http.MethodPost, func(w http.ResponseWriter, r *http.Request) {
					response.ToJson(w, http.StatusCreated, "SUCCESS", "1")
				})
				return svr
			},
			ValidateFunc: func(id string, err error) {
				if err.Error() != "unable to parse response data" {
					t.Errorf("Want: %v, Got: %v", "unable to parse response data", err.Error())
				}
			},
			cleanupFunc: func(svr *httptest.Server) {
				svr.Close()
			},
		},
		{
			name: "Failure:: Register:: json error",
			setupFunc: func() *httptest.Server {
				svr := testServer("/v1/register", http.MethodPost, func(w http.ResponseWriter, r *http.Request) {
					json.NewEncoder(w).Encode(Reader(""))
				})
				return svr
			},
			ValidateFunc: func(id string, err error) {
				if err.Error() != "json: cannot unmarshal string into Go value of type model.Response" {
					t.Errorf("Want: %v, Got: %v", errors.New("json: cannot unmarshal string into Go value of type model.Response").Error(), err.Error())
				}
			},
			cleanupFunc: func(svr *httptest.Server) {
				svr.Close()
			},
		},
		{
			name: "Failure:: Register :: incorrect status code received",
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
			name: "Failure:: Register :: client.do failure",
			setupFunc: func() *httptest.Server {
				svr := testServer("new", http.MethodPost, func(w http.ResponseWriter, r *http.Request) {
					response.ToJson(w, http.StatusCreated, "SUCCESS", nil)
				})
				svr.Close()
				return svr
			},
			ValidateFunc: func(id string, err error) {
				if !strings.Contains(err.Error(), "No connection could be made because the target machine actively refused it") {
					t.Errorf("Want: %v, Got: %v", "No connection could be made because the target machine actively refused it.", err.Error())
				}
				t.Log(err.Error())
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
			id, err := calls.Register([]byte("abc"))

			tt.ValidateFunc(id, err)
		})
	}
}
func Test_Replace(t *testing.T) {
	tests := []struct {
		name              string
		id                string
		setupFunc         func() *httptest.Server
		mockServerHandler func(w http.ResponseWriter, r *http.Request)
		ValidateFunc      func(err error)
		cleanupFunc       func(*httptest.Server)
		expectedResponse  model.Response
	}{
		{
			name: "Success:: Replace",
			id:   "1",
			setupFunc: func() *httptest.Server {
				svr := testServer("/v1/register/{id}", http.MethodPut, func(w http.ResponseWriter, r *http.Request) {
					vars := mux.Vars(r)
					id, ok := vars["id"]
					if !ok {
						response.ToJson(w, http.StatusBadRequest, "id not found", nil)
						return
					}
					err := r.ParseMultipartForm(10000)
					if err != nil {
						response.ToJson(w, http.StatusBadRequest, "file size exceeded", nil)
						t.Error(err.Error())
						return
					}
					file, _, err := r.FormFile("file")
					if err != nil {
						response.ToJson(w, http.StatusBadRequest, "file parse fail", nil)
						t.Error(err.Error())
						return
					}
					fileBytes, err := ioutil.ReadAll(file)
					if err != nil {
						response.ToJson(w, http.StatusBadRequest, "file read fail", nil)
						t.Error(err.Error())
						return
					}
					if !reflect.DeepEqual(id, "1") {
						t.Errorf("Want: %v, Got: %v", "1", id)
					}
					if !reflect.DeepEqual(fileBytes, []byte("abc")) {
						t.Errorf("Want: %v, Got: %v", []byte("abc"), fileBytes)
					}
					defer file.Close()
					cType := r.Header.Get("Content-Type")
					x := strings.Split(cType, ";")
					if !reflect.DeepEqual(x[0], "multipart/form-data") {
						t.Errorf("Want: %v, Got: %v", "multipart/form-data", x[0])
					}
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
			name: "Failure:: Replace :: type inference fail",
			id:   "1",
			setupFunc: func() *httptest.Server {
				svr := testServer("/v1/register/{id}", http.MethodPut, func(w http.ResponseWriter, r *http.Request) {
					response.ToJson(w, http.StatusOK, "SUCCESS", Reader(""))
				})
				return svr
			},
			ValidateFunc: func(err error) {
				if err.Error() != "unable to parse response data" {
					t.Errorf("Want: %v, Got: %v", "unable to parse response data", err)
				}
			},
			cleanupFunc: func(svr *httptest.Server) {
				svr.Close()
			},
		},
		{
			name: "Failure:: Replace :: json failure",
			id:   "1",
			setupFunc: func() *httptest.Server {
				svr := testServer("/v1/register/{id}", http.MethodPut, func(w http.ResponseWriter, r *http.Request) {

					json.NewEncoder(w).Encode("")
				})
				return svr
			},
			ValidateFunc: func(err error) {
				if err.Error() != "json: cannot unmarshal string into Go value of type model.Response" {
					t.Errorf("Want: %v, Got: %v", "json: cannot unmarshal string into Go value of type model.Response", err)
				}
			},
			cleanupFunc: func(svr *httptest.Server) {
				svr.Close()
			},
		},
		{
			name: "Failure:: Replace :: incorrect id in response",
			id:   "1",
			setupFunc: func() *httptest.Server {
				svr := testServer("/v1/register/{id}", http.MethodPut, func(w http.ResponseWriter, r *http.Request) {
					response.ToJson(w, http.StatusOK, "SUCCESS", map[string]interface{}{
						"id": "gfhgv",
					})
				})
				return svr
			},
			ValidateFunc: func(err error) {
				if err.Error() != "incorrect id received in response" {
					t.Errorf("Want: %v, Got: %v", "incorrect id received in response", err)
				}
			},
			cleanupFunc: func(svr *httptest.Server) {
				svr.Close()
			},
		},
		{
			name: "Failure:: Replace :: incorrect status code",
			setupFunc: func() *httptest.Server {
				svr := testServer("", http.MethodPut, func(w http.ResponseWriter, r *http.Request) {})
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
			name: "Failure:: Replace :: client.do failure",
			setupFunc: func() *httptest.Server {
				svr := testServer("/v1/register{id}", http.MethodPut, func(w http.ResponseWriter, r *http.Request) {})
				svr.Close()
				return svr
			},
			ValidateFunc: func(err error) {
				if !strings.Contains(err.Error(), "No connection could be made because the target machine actively refused it") {
					t.Errorf("Want: %v, Got: %v", "No connection could be made because the target machine actively refused it.", err.Error())
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
			t.Log(svr.URL)
			calls := NewHtmlToPdfSvc(svr.URL)

			err := calls.Replace([]byte("abc"), tt.id)

			tt.ValidateFunc(err)
		})
	}
}
func Test_GeneratePdf(t *testing.T) {
	type Student struct {
		Name  string
		Marks int
		Id    string
	}
	type Class []Student
	var class Class
	// defining struct instance
	std1 := Student{"A", 90, "1"}
	std2 := Student{"B", 100, "2"}
	std3 := Student{"C", 88, "3"}
	std4 := Student{"D", 25, "4"}
	std5 := Student{"E", 35, "5"}
	class = append(class, std4, std2, std3, std1, std5)
	tests := []struct {
		name              string
		id                string
		data              map[string]interface{}
		setupFunc         func() *httptest.Server
		mockServerHandler func(w http.ResponseWriter, r *http.Request)
		ValidateFunc      func(b []byte, err error)
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
						response.ToJson(w, http.StatusBadRequest, "id not found", nil)
						return
					}
					var data GenerateReq
					err := json.NewDecoder(r.Body).Decode(&data)
					if err != nil {
						response.ToJson(w, http.StatusBadRequest, "error decoding data", nil)
						t.Error(err.Error())
						return
					}
					data.Id = id
					tempData := GenerateReq{
						Values: map[string]interface{}{"id": "1"},
						Id:     "1",
					}
					if !reflect.DeepEqual(data, tempData) {
						t.Errorf("Want: %v, Got: %v", tempData, data)
					}
					if !reflect.DeepEqual(r.Header.Get("Content-Type"), "application/json") {
						t.Errorf("Want: %v, Got: %v", "application/json", r.Header.Get("Content-Type"))
					}
					response.ToJson(w, http.StatusOK, "", nil)
				})
				return svr
			},
			ValidateFunc: func(b []byte, err error) {
				if err != nil {
					t.Errorf("Want: %v, Got: %v", nil, err.Error())
				}
				var tempResp model.Response
				err = json.Unmarshal(b, &tempResp)
				if err != nil {
					t.Error(err)
					return
				}
				if !reflect.DeepEqual(tempResp, model.Response{
					Status:  http.StatusOK,
					Message: "",
					Data:    nil,
				}) {
					t.Errorf("Want: %v, Got: %v", model.Response{
						Status:  http.StatusOK,
						Message: "",
						Data:    nil,
					}, tempResp)
				}
			},
			cleanupFunc: func(svr *httptest.Server) {
				svr.Close()
			},
		},
		{
			name: "Failure:: GeneratePdf :: marshall error",
			id:   "1",
			data: map[string]interface{}{
				"foo": make(chan int),
			},
			setupFunc: func() *httptest.Server {
				svr := testServer("/v1/generate/{id}", http.MethodPost, func(w http.ResponseWriter, r *http.Request) {
				})
				return svr
			},
			ValidateFunc: func(b []byte, err error) {
				if !strings.Contains(err.Error(), "json") {
					t.Errorf("Want: %v, Got: %v", "json: unsupported type: chan int", err.Error())
				}
			},
			cleanupFunc: func(svr *httptest.Server) {
				svr.Close()
			},
		},
		{
			name: "Failure:: GeneratePdf :: client.do failure",
			id:   "1",
			data: map[string]interface{}{"id": "1"},
			setupFunc: func() *httptest.Server {
				svr := testServer("/v1/generate/{id}", http.MethodPost, func(w http.ResponseWriter, r *http.Request) {})
				svr.Close()
				return svr
			},
			ValidateFunc: func(b []byte, err error) {
				if !strings.Contains(err.Error(), "No connection could be made because the target machine actively refused it") {
					t.Errorf("Want: %v, Got: %v", "No connection could be made because the target machine actively refused it", err.Error())
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
					response.ToJson(w, http.StatusNotFound, "", nil)
				})
				return svr
			},
			ValidateFunc: func(b []byte, err error) {
				if !reflect.DeepEqual(err.Error(), "non success status code received : 404") {
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
			b, err := calls.GeneratePdf(tt.data, tt.id)

			tt.ValidateFunc(b, err)
		})
	}
}
