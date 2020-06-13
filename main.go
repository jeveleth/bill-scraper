package main

import (
	"os"

	"github.com/jeveleth/bill-scraper/utils"
)

var URL = "https://openstates.org/graphql?="
var API_KEY = os.Getenv("OPENSTATES_API_KEY")

var thisConfig = utils.MustLoadConfig()

var state = thisConfig.State
var searchPhrase = thisConfig.SearchPhrase
var numBills = thisConfig.NumBills
var session = thisConfig.Session

func main() {
	data, err := query(state, searchPhrase, numBills, session)
	checkError("Error getting data from openstates api", err)
	parseData(data)
}
