package logic

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/vatsal278/html-pdf-service/internal/codes"
	"github.com/vatsal278/html-pdf-service/internal/model"
	"github.com/vatsal278/html-pdf-service/internal/repo/htmlToPdf"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/PereRohit/util/log"
	respModel "github.com/PereRohit/util/model"
	"github.com/vatsal278/html-pdf-service/internal/repo/datasource"
)

//go:generate mockgen --build_flags=--mod=mod --destination=./../../pkg/mock/mock_logic.go --package=mock github.com/vatsal278/html-pdf-service/internal/logic HtmlPdfServiceLogicIer

type HtmlPdfServiceLogicIer interface {
	Ping(*model.PingRequest) *respModel.Response
	HealthCheck() bool
	HtmlToPdf(w io.Writer, req *model.GenerateReq) *respModel.Response
	Upload(file io.Reader) *respModel.Response
	Replace(id string, file io.Reader) *respModel.Response
}

type htmlPdfServiceLogic struct {
	dsSvc datasource.DataSource
	htSvc htmlToPdf.HtmlToPdf
}

func NewHtmlPdfServiceLogic(ds datasource.DataSource, ht htmlToPdf.HtmlToPdf) HtmlPdfServiceLogicIer {
	return &htmlPdfServiceLogic{
		dsSvc: ds,
		htSvc: ht,
	}
}

func (l htmlPdfServiceLogic) Ping(req *model.PingRequest) *respModel.Response {
	// add business logic here
	res, err := l.dsSvc.Ping(&model.PingDs{
		Data: req.Data,
	})
	if err != nil {
		log.Error("datasource error", err)
		return &respModel.Response{
			Status:  http.StatusInternalServerError,
			Message: "",
			Data:    nil,
		}
	}
	return &respModel.Response{
		Status:  http.StatusOK,
		Message: "Pong",
		Data:    res,
	}
}

func (l htmlPdfServiceLogic) HealthCheck() bool {
	// check all internal services are working fine
	return l.dsSvc.HealthCheck() && l.htSvc.HealthCheck()

}

func (l htmlPdfServiceLogic) Upload(file io.Reader) *respModel.Response {
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Error(err)
		return &respModel.Response{
			Status:  http.StatusInternalServerError,
			Message: codes.GetErr(codes.ErrReadFileFail),
			Data:    nil,
		}
	}
	jb, err := l.htSvc.GetJsonFromHtml(fileBytes)
	if err != nil {
		return &respModel.Response{
			Status:  http.StatusInternalServerError,
			Message: codes.GetErr(codes.ErrFileConversionFail),
			Data:    nil,
		}
	}
	u := uuid.NewString()
	err = l.dsSvc.SaveFile(u, jb, 0)
	if err != nil {
		log.Error(err)
		return &respModel.Response{
			Status:  http.StatusInternalServerError,
			Message: codes.GetErr(codes.ErrFileStoreFail),
			Data:    nil,
		}
	}
	return &respModel.Response{
		Status:  http.StatusCreated,
		Message: "SUCCESS",
		Data: map[string]interface{}{
			"id": u,
		},
	}
}

func (l htmlPdfServiceLogic) Replace(id string, file io.Reader) *respModel.Response {
	_, err := l.dsSvc.GetFile(id)
	if err != nil {
		log.Error(err)
		return &respModel.Response{
			Status:  http.StatusInternalServerError,
			Message: codes.GetErr(codes.ErrKeyNotFound),
			Data:    nil,
		}
	}
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Error(err)
		return &respModel.Response{
			Status:  http.StatusInternalServerError,
			Message: codes.GetErr(codes.ErrReadFileFail),
			Data:    nil,
		}
	}
	jb, err := l.htSvc.GetJsonFromHtml(fileBytes)
	if err != nil {
		return &respModel.Response{
			Status:  http.StatusInternalServerError,
			Message: codes.GetErr(codes.ErrFileConversionFail),
			Data:    nil,
		}
	}
	err = l.dsSvc.SaveFile(id, jb, 0)
	if err != nil {
		log.Error(err)
		return &respModel.Response{
			Status:  http.StatusInternalServerError,
			Message: codes.GetErr(codes.ErrFileStoreFail),
			Data:    nil,
		}
	}
	return &respModel.Response{
		Status:  http.StatusOK,
		Message: "SUCCESS",
		Data: map[string]interface{}{
			"id": id,
		},
	}
}

func (l htmlPdfServiceLogic) HtmlToPdf(w io.Writer, req *model.GenerateReq) *respModel.Response {
	var z map[string]interface{}
	b, err := l.dsSvc.GetFile(req.Id)
	if err != nil {
		log.Error(err)
		return &respModel.Response{
			Status:  http.StatusInternalServerError,
			Message: codes.GetErr(codes.ErrFetchingFile),
			Data:    nil,
		}
	}
	err = json.NewDecoder(bytes.NewBuffer(b)).Decode(&z)
	if err != nil {
		log.Error("error unmarshalling JSON:" + err.Error())
		return &respModel.Response{
			Status:  http.StatusInternalServerError,
			Message: codes.GetErr(codes.ErrFileParseFail),
			Data:    nil,
		}
	}
	k, ok := z["Pages"].([]interface{})
	if !ok {
		log.Error("assertion for Pages failed")
		return &respModel.Response{
			Status:  http.StatusBadRequest,
			Message: codes.GetErr(codes.ErrDecodingData),
			Data:    nil,
		}
	}
	for i, p := range k {
		page, ok := p.(map[string]interface{})
		if !ok {
			log.Error("assertion for map[string]interface{} failed")
			return &respModel.Response{
				Status:  http.StatusBadRequest,
				Message: codes.GetErr(codes.ErrDecodingData),
				Data:    nil,
			}
		}
		l, ok := page["Base64PageData"].(string)
		if !ok {
			continue
		}
		buf, err := base64.StdEncoding.DecodeString(l)
		if err != nil {
			log.Error("error decoding base 64 input on page " + fmt.Sprint(i) + " " + err.Error())
			return &respModel.Response{
				Status:  http.StatusInternalServerError,
				Message: codes.GetErr(codes.ErrDecodingData),
				Data:    nil,
			}
		}
		//log.Error(fmt.Sprintf("%+v", *req))
		t, err := template.New(req.Id).Parse(string(buf))
		if err != nil {
			log.Error(err)
			return &respModel.Response{
				Status:  http.StatusInternalServerError,
				Message: codes.GetErr(codes.ErrFileParseFail),
				Data:    nil,
			}
		}
		buffer := bytes.NewBuffer(nil)
		err = t.Execute(buffer, req.Values)
		//empty id failure case
		if err != nil {
			log.Error(err)
			return &respModel.Response{
				Status:  http.StatusInternalServerError,
				Message: codes.GetErr(codes.ErrFileStoreFail),
				Data:    nil,
			}
		}
		page["Base64PageData"] = base64.StdEncoding.EncodeToString(buffer.Bytes())
	}
	buff := bytes.NewBuffer(nil)
	err = json.NewEncoder(buff).Encode(z)
	if err != nil {
		log.Error(err)
		return &respModel.Response{
			Status:  http.StatusInternalServerError,
			Message: codes.GetErr(codes.ErrEncodingFile),
			Data:    nil,
		}
	}
	err = l.htSvc.GeneratePdf(w, buff.Bytes())
	if err != nil {
		log.Error(err)
		return &respModel.Response{
			Status:  http.StatusInternalServerError,
			Message: codes.GetErr(codes.ErrConvertingToPdf),
			Data:    nil,
		}
	}
	return &respModel.Response{Status: http.StatusOK}
}
