package htmlToPdf

import (
	"io"
	"net/http/httptest"
	"os"
	"testing"
)

func TestGetJsonFromHtml(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing due to unavailability of testing environment")
	}
	htmlToPdf := NewWkHtmlToPdfSvc()
	tests := []struct {
		name         string
		setupFunc    func() []byte
		validateFunc func([]byte, error)
	}{
		{
			name: "Success",
			setupFunc: func() []byte {
				b, err := os.ReadFile("Failure.html")
				if err != nil {
					t.Error(err)
				}
				return b
			},
			validateFunc: func(jsBytes []byte, err error) {
				if err != nil {
					t.Error(err)
				}
				w := httptest.NewRecorder()
				err = htmlToPdf.GeneratePdf(w, jsBytes)
				if err != nil {
					t.Error(err)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			js, err := htmlToPdf.GetJsonFromHtml(tt.setupFunc())
			tt.validateFunc(js, err)

		})
	}
}
func TestGeneratePdf(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing due to unavailability of testing environment")
	}
	htmlToPdf := NewWkHtmlToPdfSvc()
	tests := []struct {
		name         string
		setupFunc    func() []byte
		validateFunc func(io.Writer, error)
	}{
		{
			name: "Success",
			setupFunc: func() []byte {
				b, err := os.ReadFile("Failure.html")
				if err != nil {
					t.Error(err)
				}
				js, err := htmlToPdf.GetJsonFromHtml(b)
				return js
			},
			validateFunc: func(w io.Writer, err error) {
				if err != nil {
					t.Error(err)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			err := htmlToPdf.GeneratePdf(w, tt.setupFunc())
			tt.validateFunc(w, err)

		})
	}
}
