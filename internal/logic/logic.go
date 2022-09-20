package logic

import (
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
}

type htmlPdfServiceLogic struct {
	dummyDsSvc datasource.DataSource
}

func NewHtmlPdfServiceLogic(ds datasource.DataSource) HtmlPdfServiceLogicIer {
	return &htmlPdfServiceLogic{
		dummyDsSvc: ds,
	}
}

func (l htmlPdfServiceLogic) Ping(req *model.PingRequest) *respModel.Response {
	// add business logic here
	res, err := l.dummyDsSvc.Ping(&model.PingDs{
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
	return l.dummyDsSvc.HealthCheck()
}
