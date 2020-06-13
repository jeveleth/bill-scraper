package main

import (
	"encoding/csv"
	"encoding/json"
	"log"
	"os"
)

// Bill is a record of a bill
type Bill struct {
	Number           string       `json:"identifier"`
	Title            string       `json:"title"`
	UpdatedAt        string       `json:"updatedAt"`
	FromOrganization Organization `json:"fromOrganization"`
	Abstracts        []Abstract   `json:"abstracts"`
	Actions          []Action     `json:"actions"`
}

// Organization represents the legislative body
type Organization struct {
	Name           string `json:"name"`
	Classification string `json:"classification"`
}

// Abstract represents an official abstract for a bill
type Abstract struct {
	Abstract string `json:"abstract"`
	Note     string `json:"note"`
	Date     string `json:"date"`
}

type Action struct {
	Date           string `json:"date"`
	Description    string `json:"description"`
	Classification string `json:"classification"`
}
type Node struct {
	Node Bill `json:"node"`
}
type ResponseBody struct {
	Data Search `json:"data"`
}

type Search struct {
	Search Edges `json:"search_1"`
}

type Edges struct {
	Edges []Node `json:"edges"`
}

var bills [][]string

func parseData(body []byte) {
	res := ResponseBody{}
	json.Unmarshal(body, &res)
	var fmtBill []string
	for _, bill := range res.Data.Search.Edges {
		fmtBill = []string{bill.Node.Number,
			bill.Node.Title,
			bill.Node.FromOrganization.Name,
			bill.Node.UpdatedAt,
		}
		action := bill.Node.Actions[len(bill.Node.Actions)-1]
		fmtBill = append(fmtBill, action.Description)
		for _, abstract := range bill.Node.Abstracts {
			fmtBill = append(fmtBill, abstract.Abstract)
		}
		bills = append(bills, fmtBill)
	}

	file, err := os.Create("legislation.csv")
	checkError("Cannot create file", err)
	defer file.Close()

	w := csv.NewWriter(file)
	defer w.Flush()

	w.Write([]string{"Bill Number", "Title", "Chamber", "Updated", "Action", "Abstract"})
	for _, bill := range bills {
		log.Printf("bill is %v", bill)
		if err := w.Write(bill); err != nil {
			log.Fatalln("error writing bill to csv:", err)
		}
	}

	// Write any buffered data to the underlying writer (standard output).
	w.Flush()

	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
}

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}
