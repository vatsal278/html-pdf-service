package htmlToPdf

import (
	"bytes"
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"io"
)

type wkHtmlToPdf struct {
}

func NewWkHtmlToPdfSvc() HtmlToPdf {
	return &wkHtmlToPdf{}
}

func (w wkHtmlToPdf) HealthCheck() bool {
	return true
}

func (w wkHtmlToPdf) GetJsonFromHtml(b []byte) ([]byte, error) {
	pdfg := wkhtmltopdf.NewPDFPreparer()
	pdfg.AddPage(wkhtmltopdf.NewPageReader(bytes.NewReader(b)))
	jb, err := pdfg.ToJSON()
	if err != nil {
		return nil, err
	}
	return jb, nil
}
func (w wkHtmlToPdf) GeneratePdf(wr io.Writer, b []byte) error {
	pdfgFromJSON, err := wkhtmltopdf.NewPDFGeneratorFromJSON(bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	pdfgFromJSON.SetOutput(wr)
	err = pdfgFromJSON.Create()
	if err != nil {
		return err
	}
	return nil
}
