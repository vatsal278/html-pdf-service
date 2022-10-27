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
	"time"
)

type htmlToPdfSvc struct {
	svcUrl string
	client http.Client
}

func NewHtmlToPdfSvc(url string) HtmlToPdfSvcI {
	return &htmlToPdfSvc{
		svcUrl: url,
		client: http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

type HtmlToPdfSvcI interface {
	Register([]byte) (string, error)
	Replace([]byte, string) error
	GeneratePdf(map[string]interface{}, string) ([]byte, error)
}

func (h *htmlToPdfSvc) Register(fileBytes []byte) (string, error) {
	//take in the slice of bytes so that you can pass it directly to io.copy
	//take io.reader as argument
	file := bytes.NewReader(fileBytes)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "output")
	if err != nil {
		return "", err
	}
	_, err = io.Copy(part, file)
	writer.Close()
	r, err := h.client.Post(h.svcUrl+"/v1/register", writer.FormDataContentType(), body)
	if err != nil {
		return "", err
	}

	if r.StatusCode < 200 || r.StatusCode > 299 {
		return "", fmt.Errorf("non success status code received : %v", r.StatusCode)
	}
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

func (h *htmlToPdfSvc) Replace(fileBytes []byte, id string) error {
	file := bytes.NewReader(fileBytes)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "output")
	if err != nil {
		return err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}
	contType := writer.FormDataContentType()
	writer.Close()
	r, err := http.NewRequest(http.MethodPut, h.svcUrl+"/v1/register/"+id, body)
	if err != nil {
		return err
	}
	r.Header.Set("Content-Type", contType)
	r = mux.SetURLVars(r, map[string]string{"id": id})
	resp, err := h.client.Do(r)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("non success status code received : %v", resp.StatusCode)
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var response respModel.Response
	err = json.Unmarshal(bytes, &response)
	if err != nil {
		return err
	}
	data, ok := response.Data.(map[string]interface{})
	if !ok {
		errNew := errors.New("unable to parse response data")
		return errNew
	}
	if data["id"] != id {
		return fmt.Errorf("incorrect id received in response")
	}
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

	r, err := http.NewRequest(http.MethodPost, h.svcUrl+"/v1/generate/"+id, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	r.Header.Set("Content-Type", "application/json")
	resp, err := h.client.Do(r)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("non success status code received : %v", resp.StatusCode)
	}
	file := resp.Header.Get("Content-Disposition")
	log.Info(file)
	filebyte, err := json.Marshal(file)
	if err != nil {
		return nil, err
	}
	return filebyte, err
}
