package main

import (
	"github.com/vatsal278/html-pdf-service/pkg/sdk"
	"log"
)

func main() {
	controller := sdk.NewHtmlToPdfSvc("http://localhost:9090")
	id, err := controller.Register("./test/Failure.html")
	if err != nil {
		log.Print(err.Error())
		return
	}
	//controller := sdk.NewHtmlToPdfSvc("http://localhost:9090")
	//err = controller.Replace("./test/Failure.html", id)
	//if err != nil {
	//	log.Print(err.Error())
	//	return
	//}
	//4e80ebd2-dccf-4319-8a04-beac4acbe09c

	//controller := sdk.NewHtmlToPdfSvc("http://localhost:9090")
	_, err = controller.GeneratePdf(map[string]interface{}{"id": "1"}, id)
	if err != nil {
		log.Print(err.Error())
		return
	}
	//log.Print(b[1])

}
