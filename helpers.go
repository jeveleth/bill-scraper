package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
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
	OpenstatesURL      string             `json:"openstatesUrl"`
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

// LegislativeSession contains the session and state of the legislature
type LegislativeSession struct {
	Identifier   string       `json:"identifier"`
	Jurisdiction Jurisdiction `json:"jurisdiction"`
}

// Jurisdiction contains the state of the legislature
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

// ResponseBody contains the data returned from the GraphQL API
type ResponseBody struct {
	Data Search `json:"data"`
}

// Search contains the search results
type Search struct {
	Search SearchResults `json:"search_1"`
}

// SearchResults contains the edge and paging information
type SearchResults struct {
	Edges    []Node   `json:"edges"`
	PageInfo PageInfo `json:"pageInfo"`
}

// PageInfo contains the paging information
type PageInfo struct {
	HasNextPage bool   `json:"hasNextPage"`
	EndCursor   string `json:"endCursor"`
}

var bills [][]string

func formatBill(res ResponseBody) []string {
	var fmtBill []string
	// if there is more data to get, run another API call from the cursor
	if res.Data.Search.PageInfo.HasNextPage == true {
		next = true
		cursor = res.Data.Search.PageInfo.EndCursor
	}
	for _, bill := range res.Data.Search.Edges {
		fmtBill = []string{bill.Node.Number,
			bill.Node.Title,
			bill.Node.FromOrganization.Name,
			bill.Node.UpdatedAt,
			bill.Node.OpenstatesURL,
			bill.Node.LegislativeSession.Identifier,
			bill.Node.LegislativeSession.Jurisdiction.Name,
		}
		// Not all legislatures have information about actions (e.g., CA)
		if len(bill.Node.Actions) > 0 {
			// Grab the most recent action for the bill
			action := bill.Node.Actions[len(bill.Node.Actions)-1]
			fmtBill = append(fmtBill, action.Description)
		}

		for _, abstract := range bill.Node.Abstracts {
			fmtBill = append(fmtBill, abstract.Abstract)
		}
		bills = append(bills, fmtBill)
	}
	return fmtBill
}

func writeToCSV(file io.Writer) {
	w := csv.NewWriter(file)
	defer w.Flush()

	// Add file headers only when they don't already exist
	fi, err := os.Open("legislation.csv")
	stat, err := fi.Stat()
	checkError("error getting file info", err)
	if stat.Size() <= 0 {
		w.Write([]string{
			"Bill Number",
			"Title",
			"Chamber",
			"Updated",
			"OpenstatesURL",
			"Session",
			"Jurisdiction",
			"Latest Action",
			"Bill Abstract",
		})
	}

	for _, bill := range bills {
		if err := w.Write(bill); err != nil {
			log.Fatalln("error writing bill to csv:", err)
		}
	}

	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
}

func getAndProcessBills(body []byte) {
	var err error

	os.Chdir("/tmp")
	file, err := os.OpenFile("legislation.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	path, err := os.Getwd()

	if err != nil {
		log.Println(err)
	}
	fmt.Println(path)

	checkError("Cannot create file", err)
	defer file.Close()
	res := ResponseBody{}

	json.Unmarshal(body, &res)
	formatBill(res)
	writeToCSV(file)
	for next == true {
		log.Println("There are more bills to retrieve. Grabbing them for processing.")
		body, err = query(state, searchPhrase, numBills, session, cursor)
		if err != nil {
			log.Panicf("Error calling API %v", err)
		}
		json.Unmarshal(body, &res)
		formatBill(res)
		writeToCSV(file)
		if res.Data.Search.PageInfo.HasNextPage == false {
			next = false
		}
	}
	log.Println("File processing complete. Check legislation.csv")
}

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}

func query(state string, searchTerm string, numResults int, session string, cursor string) ([]byte, error) {
	var apiQuery string
	apiQuery = fmt.Sprintf(`{
		search_1: bills(first: %d, jurisdiction: "%s", searchQuery: "%s", session:"%s", after:"%s") {
		  edges {
			node {
			  id
			  identifier
			  title
			  classification
			  updatedAt
			  openstatesUrl
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
		  pageInfo {
			hasNextPage
			endCursor
		  }
		}
	  }
	  `, numResults, state, searchTerm, session, cursor)
	if cursor == "" {
		apiQuery = fmt.Sprintf(`{
		search_1: bills(first: %d, jurisdiction: "%s", searchQuery: "%s", session:"%s") {
		  edges {
			node {
			  identifier
			  title
			  classification
			  updatedAt
			  openstatesUrl
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
		  pageInfo {
			hasNextPage
			endCursor
		  }
		}
	  }
	  `, numResults, state, searchTerm, session)
	}

	jsonData := map[string]string{"query": apiQuery}
	jsonValue, _ := json.Marshal(jsonData)
	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(jsonValue))
	req.Header.Add("content-type", "application/json")
	req.Header.Add("x-api-key", API_KEY)
	fmtAuth := fmt.Sprintf("Bearer %s", API_KEY)
	req.Header.Add("authorization", fmtAuth)
	client := &http.Client{Timeout: time.Second * 60}
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
