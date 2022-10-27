package main

import (
	"github.com/vatsal278/html-pdf-service/pkg/sdk"
	"io/ioutil"
	"log"
	"os"
)

type Cart struct {
	Item   string
	Amount float64
}
type Student struct {
	Name  string
	Marks int
	Id    string
}
type Class []Student

func main() {
	// after registering generate the pdf and then replace and again generate pdf
	controller := sdk.NewHtmlToPdfSvc("http://localhost:9090")
	fileBytes, err := os.ReadFile("./docs/Failure.html")
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
	var class Class
	// defining struct instance
	std1 := Student{"A", 90, "1"}
	std2 := Student{"B", 100, "2"}
	std3 := Student{"C", 88, "3"}
	std4 := Student{"D", 25, "4"}
	std5 := Student{"E", 35, "5"}
	class = append(class, std4, std2, std3, std1, std5)
	b, err := controller.GeneratePdf(map[string]interface{}{"data": class}, id)
	if err != nil {
		log.Print(err.Error())
		return
	}
	log.Print(string(b))
	err = ioutil.WriteFile("./docs/output.pdf", b, 0777)
	if err != nil {
		log.Fatalln(err)
	}

	fileBytes, err = os.ReadFile("./docs/Template.html")
	err = controller.Replace(fileBytes, id)
	if err != nil {
		log.Print(err.Error())
		return
	}
	b, err = controller.GeneratePdf(map[string]any{
		"Items": []Cart{
			{
				Item:   "Bread",
				Amount: 24,
			},
			{
				Item:   "Rice",
				Amount: 56.7,
			},
			{
				Item:   "Clothes",
				Amount: 150.45,
			},
			{
				Item:   "Water",
				Amount: 100,
			},
			{
				Item:   "Gas",
				Amount: 100.00,
			},
		},
		"Title": "Inventory list",
	}, id)
	if err != nil {
		log.Print(err.Error())
		return
	}
	err = ioutil.WriteFile("./docs/newOutput.pdf", b, 0777)
	if err != nil {
		log.Fatalln(err)
	}
}
