package htmlToPdf

import (
	"bytes"
	"github.com/PereRohit/util/log"
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

func (w wkHtmlToPdf) PDFPreparer(b []byte) ([]byte, error) {
	pdfg := wkhtmltopdf.NewPDFPreparer()
	pdfg.AddPage(wkhtmltopdf.NewPageReader(bytes.NewReader(b)))
	jb, err := pdfg.ToJSON()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return jb, nil
}
func (w wkHtmlToPdf) PDFGenerator(wr io.Writer, b []byte) error {
	pdfgFromJSON, err := wkhtmltopdf.NewPDFGeneratorFromJSON(bytes.NewBuffer(b))
	if err != nil {
		log.Error(err)
		return err
	}
	pdfgFromJSON.SetOutput(wr)
	err = pdfgFromJSON.Create()

	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}
