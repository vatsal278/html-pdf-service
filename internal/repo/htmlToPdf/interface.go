package htmlToPdf

import "io"

type HtmlToPdf interface {
	HealthCheck() bool
	PDFGenerator(io.Writer, []byte) error
	PDFPreparer([]byte) ([]byte, error)
}
