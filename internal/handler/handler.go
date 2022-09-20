package handler

import (
	"fmt"
	"github.com/vatsal278/html-pdf-service/internal/repo/htmlToPdf"
	"net/http"

	"github.com/PereRohit/util/request"
	"github.com/PereRohit/util/response"

	"github.com/vatsal278/html-pdf-service/internal/logic"
	"github.com/vatsal278/html-pdf-service/internal/model"
	"github.com/vatsal278/html-pdf-service/internal/repo/datasource"
)

const HtmlPdfServiceName = "htmlPdfService"

//go:generate mockgen --build_flags=--mod=mod --destination=./../../pkg/mock/mock_handler.go --package=mock github.com/vatsal278/html-pdf-service/internal/handler HtmlPdfServiceHandler

type HtmlPdfServiceHandler interface {
	HealthChecker
	Ping(w http.ResponseWriter, r *http.Request)
}

type htmlPdfService struct {
	logic logic.HtmlPdfServiceLogicIer
}

func NewHtmlPdfService(ds datasource.DataSource, ht htmlToPdf.HtmlToPdf) HtmlPdfServiceHandler {
	svc := &htmlPdfService{
		logic: logic.NewHtmlPdfServiceLogic(ds, ht),
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
