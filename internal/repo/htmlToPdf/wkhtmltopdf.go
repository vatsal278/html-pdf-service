package htmlToPdf

type wkHtmlToPdf struct {
}

func NewWkHtmlToPdfSvc() HtmlToPdf {
	return &wkHtmlToPdf{}
}

func (w wkHtmlToPdf) HealthCheck() bool {
	return true
}
