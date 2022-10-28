package handler

import (
	"encoding/json"
	"fmt"
	"github.com/PereRohit/util/log"
	"github.com/PereRohit/util/request"
	"github.com/PereRohit/util/response"
	"github.com/gorilla/mux"
	"github.com/vatsal278/html-pdf-service/internal/codes"
	"github.com/vatsal278/html-pdf-service/internal/repo/htmlToPdf"
	"net/http"

	"github.com/vatsal278/html-pdf-service/internal/logic"
	"github.com/vatsal278/html-pdf-service/internal/model"
	"github.com/vatsal278/html-pdf-service/internal/repo/datasource"
)

const HtmlPdfServiceName = "htmlPdfService"

//go:generate mockgen --build_flags=--mod=mod --destination=./../../pkg/mock/mock_handler.go --package=mock github.com/vatsal278/html-pdf-service/internal/handler HtmlPdfServiceHandler

type HtmlPdfServiceHandler interface {
	HealthChecker
	Ping(w http.ResponseWriter, r *http.Request)
	Upload(w http.ResponseWriter, r *http.Request)
	ConvertToPdf(w http.ResponseWriter, r *http.Request)
	ReplaceHtml(w http.ResponseWriter, r *http.Request)
}

type htmlPdfService struct {
	logic     logic.HtmlPdfServiceLogicIer
	maxMemory int64
}

func NewHtmlPdfService(ds datasource.DataSource, ht htmlToPdf.HtmlToPdf, mx int64) HtmlPdfServiceHandler {
	svc := &htmlPdfService{
		logic:     logic.NewHtmlPdfServiceLogic(ds, ht),
		maxMemory: mx,
	}
	AddHealthChecker(svc)
	return svc
}

func (svc htmlPdfService) HealthCheck() (svcName string, msg string, stat bool) {
	set := false
	defer func() {
		svcName = HtmlPdfServiceName
		if !set {
			msg = ""
			stat = true
		}
	}()
	stat = svc.logic.HealthCheck()
	set = true
	return
}
func (svc htmlPdfService) Ping(w http.ResponseWriter, r *http.Request) {
	req := &model.PingRequest{}

	suggestedCode, err := request.FromJson(r, req)
	if err != nil {
		response.ToJson(w, suggestedCode, fmt.Sprintf("FAILED: %s", err.Error()), nil)
		return
	}
	// call logic
	resp := svc.logic.Ping(req)
	response.ToJson(w, resp.Status, resp.Message, resp.Data)
	return
}
func (svc htmlPdfService) Upload(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(svc.maxMemory) //File size to come from config
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
	resp := svc.logic.Upload(file)
	response.ToJson(w, resp.Status, resp.Message, resp.Data)
}
func (svc htmlPdfService) ConvertToPdf(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	//we take id as a parameter from url path
	id, ok := vars["id"]
	if !ok {
		response.ToJson(w, http.StatusBadRequest, codes.GetErr(codes.ErrIdNeeded), nil)
		return
	}
	var data model.GenerateReq
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		response.ToJson(w, http.StatusBadRequest, codes.GetErr(codes.ErrDecodingData), nil)
		log.Error(err.Error())
		return
	}
	data.Id = id
	resp := svc.logic.HtmlToPdf(w, &data)
	if resp.Status != http.StatusOK {
		response.ToJson(w, resp.Status, resp.Message, resp.Data)
		log.Error(resp.Message)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+data.Id+".pdf")
	w.Header().Set("Content-Type", "application/pdf")
}
func (svc htmlPdfService) ReplaceHtml(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	//we take id as a parameter from url path
	id, ok := vars["id"]
	if !ok {
		response.ToJson(w, http.StatusBadRequest, codes.GetErr(codes.ErrIdNeeded), nil)
		return
	}
	err := r.ParseMultipartForm(svc.maxMemory) //File size to come from config

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
	resp := svc.logic.Replace(id, file)
	response.ToJson(w, resp.Status, resp.Message, resp.Data)
}
