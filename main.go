package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var URL = "https://openstates.org/graphql?="
var API_KEY = os.Getenv("OPENSTATES_API_KEY")

func main() {

	jsonData := map[string]string{
		"query": `{
				search_1: bills(first: 5, jurisdiction: "New York", searchQuery: "Peace officer") {
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
					  actions {
						date
						description
						classification
					  }
					  documents {
						date
						note
						links {
						  url
						}
					  }
					  versions {
						date
						note
						links {
						  url
						}
					  }
					  sources {
						url
						note
					  }
					}
				  }
				}
			  }
			  `,
	}

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
	parseData(data)
}
