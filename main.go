package main

import (
	"log"

	"github.com/jeveleth/bill-scraper/utils"
)

var URL = "https://openstates.org/graphql?="

// var API_KEY = os.Getenv("OPENSTATES_API_KEY")

var thisConfig = utils.MustLoadConfig()

var state = thisConfig.State
var searchPhrase = thisConfig.SearchPhrase
var numBills = thisConfig.NumBills
var session = thisConfig.Session
var next = false
var cursor string

func main() {
	if len(API_KEY) == 0 {
		log.Panic("You need to supply your api key")
	}
	data, err := query(state, searchPhrase, numBills, session, cursor)
	checkError("Error getting data from openstates api", err)
	getAndProcessBills(data)
}
