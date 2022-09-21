package router

import (
	"github.com/vatsal278/html-pdf-service/internal/repo/htmlToPdf"
	"net/http"

	"github.com/PereRohit/util/constant"
	"github.com/PereRohit/util/middleware"
	"github.com/gorilla/mux"

	"github.com/vatsal278/html-pdf-service/internal/config"
	"github.com/vatsal278/html-pdf-service/internal/handler"
	"github.com/vatsal278/html-pdf-service/internal/repo/datasource"
)

func Register(svcCfg *config.SvcConfig) *mux.Router {
	m := mux.NewRouter()

	// group all routes for specific version. e.g.: /v1
	if svcCfg.ServiceRouteVersion != "" {
		m = m.PathPrefix("/" + svcCfg.ServiceRouteVersion).Subrouter()
	}

	m.StrictSlash(true)
	m.Use(middleware.RequestHijacker)
	m.Use(middleware.RecoverPanic)

	commons := handler.NewCommonSvc()
	m.HandleFunc(constant.HealthRoute, commons.HealthCheck).Methods(http.MethodGet)
	m.NotFoundHandler = http.HandlerFunc(commons.RouteNotFound)
	m.MethodNotAllowedHandler = http.HandlerFunc(commons.MethodNotAllowed)

	// attach routes for services below
	m = attachHtmlPdfServiceRoutes(m, svcCfg)

	return m
}

func attachHtmlPdfServiceRoutes(m *mux.Router, svcCfg *config.SvcConfig) *mux.Router {
	dataSource := datasource.NewRedisDs(&svcCfg.CacherSvc)
	htmlTopdfSvc := htmlToPdf.NewWkHtmlToPdfSvc()

	svc := handler.NewHtmlPdfService(dataSource, htmlTopdfSvc, svcCfg.MaxMemmory)

	m.HandleFunc("/ping", svc.Ping).Methods(http.MethodPost)
	m.HandleFunc("/register", svc.Upload).Methods(http.MethodPost)
	m.HandleFunc("/generate/{id}", svc.ConvertToPdf).Methods(http.MethodPost)
	m.HandleFunc("/register/{id}", svc.ReplaceHtml).Methods(http.MethodPut)
	return m
}
