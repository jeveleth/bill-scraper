package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

// Bill is a record of a bill
type Bill struct {
	Number             string             `json:"identifier"`
	Title              string             `json:"title"`
	UpdatedAt          string             `json:"updatedAt"`
	FromOrganization   Organization       `json:"fromOrganization"`
	Abstracts          []Abstract         `json:"abstracts"`
	Actions            []Action           `json:"actions"`
	LegislativeSession LegislativeSession `json:"legislativeSession"`
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

type LegislativeSession struct {
	Identifier   string       `json:"identifier"`
	Jurisdiction Jurisdiction `json:"jurisdiction"`
}
type Jurisdiction struct {
	Name string `json:"name"`
}

// Action represents an action taken on a bill
type Action struct {
	Date           string `json:"date"`
	Description    string `json:"description"`
	Classification string `json:"classification"`
}

// Node contains a bill
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
			bill.Node.LegislativeSession.Identifier,
			bill.Node.LegislativeSession.Jurisdiction.Name,
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

	w.Write([]string{"Bill Number", "Title", "Chamber", "Updated", "Session", "Jurisdiction", "Latest Action", "Bill Abstract"})
	for _, bill := range bills {
		log.Printf("bill is %v", bill)
		if err := w.Write(bill); err != nil {
			log.Fatalln("error writing bill to csv:", err)
		}
	}

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

func query(state string, searchTerm string, numResults int, session string) ([]byte, error) {
	apiQuery := fmt.Sprintf(`{
	search_1: bills(first: %d, jurisdiction: "%s", searchQuery: "%s", session:"%s") {
	  edges {
		node {
		  id
		  identifier
		  title
		  classification
		  updatedAt
		  createdAt
		  legislativeSession {
			identifier
			jurisdiction {
			  name
			}
		  }
		  fromOrganization {
			name
			classification
		  }
		  actions {
			date
			description
			classification
		  }
		  abstracts {
			abstract
		  }
		}
	  }
	}
  }
  `, numResults, state, searchTerm, session)

	jsonData := map[string]string{"query": apiQuery}
	jsonValue, _ := json.Marshal(jsonData)
	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(jsonValue))
	req.Header.Add("content-type", "application/json")
	req.Header.Add("x-api-key", API_KEY)
	fmtAuth := fmt.Sprintf("Bearer %s", API_KEY)
	req.Header.Add("authorization", fmtAuth)
	client := &http.Client{Timeout: time.Second * 10}
	res, err := client.Do(req)
	if err != nil {
		log.Panicf("error getting response %v", err)
	}
	defer res.Body.Close()
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	}
	data, _ := ioutil.ReadAll(res.Body)
	return data, nil
}
