package sdk

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PereRohit/util/log"
	respModel "github.com/PereRohit/util/model"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type htmlToPdfSvc struct {
	svcUrl string
}

func NewHtmlToPdfSvc(url string) HtmlToPdfSvcI {
	return &htmlToPdfSvc{
		svcUrl: url,
	}
}

type HtmlToPdfSvcI interface {
	Register(string) (string, error)
	Replace(string, string) error
	GeneratePdf(map[string]interface{}, string) ([]byte, error)
}

func (h *htmlToPdfSvc) Register(filePath string) (string, error) {
	//take io.reader as argument
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return "", err
	}
	_, err = io.Copy(part, file)
	err = writer.Close()
	if err != nil {
		return "", err
	}
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	r, err := client.Post(h.svcUrl+"/v1/register", writer.FormDataContentType(), body)
	if err != nil {
		return "", err
	}

	if r.StatusCode < 200 || r.StatusCode > 299 {
		return "", fmt.Errorf("non success status code received : %v", r.StatusCode)
	}
	//parse id and return the id
	resp, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	var response respModel.Response
	err = json.Unmarshal(resp, &response)
	if err != nil {
		return "", err
	}
	x := response.Data
	m, ok := x.(map[string]interface{})
	if !ok {
		errNew := errors.New("unable to parse response data")
		return "", errNew
	}
	return fmt.Sprint(m["id"]), err
}

func (h *htmlToPdfSvc) Replace(filePath string, id string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	defer file.Close()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return err
	}
	_, err = io.Copy(part, file)
	err = writer.Close()
	if err != nil {
		return err
	}
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	r, err := http.NewRequest(http.MethodPut, h.svcUrl+"/v1/register/"+id, body)
	if err != nil {
		log.Error(err)
		return err
	}
	r = mux.SetURLVars(r, map[string]string{"id": id})
	resp, err := client.Do(r)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("non success status code received : %v", resp.StatusCode)
	}
	//validate the id
	return err
}

type GenPdfReq struct {
	Values map[string]interface{} `json:"values"`
}

func (h *htmlToPdfSvc) GeneratePdf(templateData map[string]interface{}, id string) ([]byte, error) {
	b, err := json.Marshal(GenPdfReq{
		Values: templateData,
	})
	if err != nil {
		return nil, err
	}
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	r, err := http.NewRequest(http.MethodPost, h.svcUrl+"/v1/generate/"+id, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(r)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("non success status code received : %v", resp.StatusCode)
	}
	by, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	return by, err
}
