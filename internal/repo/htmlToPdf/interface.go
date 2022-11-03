package htmlToPdf

import "io"

//go:generate mockgen --build_flags=--mod=mod --destination=./../../../pkg/mock/mock_htmltopdf.go --package=mock github.com/vatsal278/html-pdf-service/internal/repo/htmlToPdf HtmlToPdf

type HtmlToPdf interface {
	HealthCheck() bool
	GeneratePdf(io.Writer, []byte) error
	GetJsonFromHtml([]byte) ([]byte, error)
}
