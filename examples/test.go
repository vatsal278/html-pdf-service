package main

import (
	"github.com/vatsal278/html-pdf-service/pkg/sdk"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	// after registering generate the pdf and then replace and again generate pdf
	controller := sdk.NewHtmlToPdfSvc("http://localhost:9090")
	fileBytes, err := os.ReadFile("./../docs/Failure.html")
	if err != nil {
		log.Print(err.Error())
		return
	}
	id, err := controller.Register(fileBytes)
	if err != nil {
		log.Print(err.Error())
		return
	}
	log.Print(id)
	b, err := controller.GeneratePdf(map[string]interface{}{"id": "1"}, id)
	if err != nil {
		log.Print(err.Error())
		return
	}
	err = ioutil.WriteFile("output.pdf", b, 0777)
	if err != nil {
		log.Fatalln(err)
	}

	fileBytes, err = os.ReadFile("./../docs/Template.html")
	err = controller.Replace(fileBytes, id)
	if err != nil {
		log.Print(err.Error())
		return
	}
	b, err = controller.GeneratePdf(map[string]interface{}{"id": "1"}, id)
	if err != nil {
		log.Print(err.Error())
		return
	}
	err = ioutil.WriteFile("newOutput.pdf", b, 0777)
	if err != nil {
		log.Fatalln(err)
	}
}
