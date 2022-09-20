package logic

import (
	"github.com/vatsal278/html-pdf-service/internal/repo/htmlToPdf"
	"io"
	"net/http"

	"github.com/PereRohit/util/log"
	respModel "github.com/PereRohit/util/model"

	"github.com/vatsal278/html-pdf-service/internal/model"
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

func (l htmlPdfServiceLogic) HtmlToPdf(w io.Writer, req *model.GenerateReq) *respModel.Response {
	//TODO implement me
	panic("implement me")
}

func (l htmlPdfServiceLogic) Upload(file io.Reader) *respModel.Response {
	//TODO implement me
	panic("implement me")
}

func (l htmlPdfServiceLogic) Replace(id string, file io.Reader) *respModel.Response {
	//TODO implement me
	panic("implement me")
}
