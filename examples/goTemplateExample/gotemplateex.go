package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"
)

type Cart struct {
	Item   string
	Amount float64
}

func main() {

	t, err := template.New("").Parse("EOF")
	if err != nil {
		log.Fatal(err)
		return
	}
	buffer := bytes.NewBuffer(nil)
	err = t.Execute(buffer, nil)
	//empty id failure case
	if err != nil {
		log.Fatalln(err)
		return
	}

	d := map[string]any{
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
	}
	data, err := json.Marshal(d)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("%s", data)

}
