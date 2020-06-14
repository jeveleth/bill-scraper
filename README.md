
This repo allows users to run a binary that will generate a CSV file that contains all relevant legislation for a given state.
Tested on a Mac, so it doesn't yet cross-compile for Windows or Linux.

Requirements:
* Docker (and docker-compose)
* Make
* An openstates.org API Key.

To use:

```
export OPENSTATES_API_KEY=[your api key]
make up # this might take a few minutes to pull all the necessary docker images
go build -v -o scraper .
```

make state="California" session="20192020" phrase="peace officer" num-bills=100 run-app

To run the program, do:
```
./scraper -state "New York" -session "2019-2020" -phrase "Use of force" -num-bills 20
```


You will then see a file named `legislation.csv`, which you can import to a Google or Excel spreadsheet. Currently, it contains columns for:
```
Bill Number
Title
Chamber
Updated
Session
Jurisdiction
Latest Action
Abstract
```

TODOs: 
* hook up to API [done]
* format timestamp
* Handle secrets: API Token [done]
* Parameterize input (so users can grab a state) [done]
* pre-load all search terms 
* pagination
* convert to a CLI tool
* make cross-platform compatible via a docker container