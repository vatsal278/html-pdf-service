package main

import (
	"github.com/vatsal278/html-pdf-service/pkg/sdk"
	"io/ioutil"
	"log"
)

func main() {
	controller := sdk.NewHtmlToPdfSvc("http://localhost:9090")
	id, err := controller.Register("./test/Failure.html")
	if err != nil {
		log.Print(err.Error())
		return
	}
	log.Print(id)
	err = controller.Replace("./test/Failure.html", id)
	if err != nil {
		log.Print(err.Error())
		return
	}
	b, err := controller.GeneratePdf(map[string]interface{}{"id": "1"}, id)
	if err != nil {
		log.Print(err.Error())
		return
	}
	err = ioutil.WriteFile("output.pdf", b, 0777)
	if err != nil {
		log.Fatalln(err)
	}
}
