## Overview:
This repository allows users to run a binary that will generate a CSV file that contains all relevant legislation for a given state.

## Requirements:
* [Docker](https://hub.docker.com/editions/community/docker-ce-desktop-mac/)
* [docker-compose](https://docs.docker.com/compose/install/)
* An OpenStates [API Key](https://openstates.org/accounts/login/).

## Usage:
### Setup:

```
export OPENSTATES_API_KEY=your-api-key
docker-compose up --build -d
make build # this might take a few minutes to pull all the necessary docker images
```
### Running the program:
To run the program, pass the search terms particular to your use case. For example, if you wish to find all bills in the 20192020 legislative session
for California taht contain the phrase "peace officer", run:
```
docker run --rm -e OPENSTATES_API_KEY=${OPENSTATES_API_KEY} -v \
    $(pwd):/tmp -ti bill-scraper_app:latest ./app \
    -state "California" \
    -session "20192020" \
    -phrase "peace officer" \
    -num-bills 100 \
    run-app
```

You will then see a file named `legislation.csv`, which you can import to a Google or Excel spreadsheet. Currently, it contains columns for:
```
Bill Number
Title
Chamber
Updated
OpenstatesURL
Session
Jurisdiction
Latest Action
Abstract
```


If you want to run the binary without docker (and you have Go set up), you can build and run by doing:
```
go build -v -o scraper .
./scraper -state "New York" -session "2019-2020" -phrase "Use of force" -num-bills 20
```

TODOs:
* format timestamp
* pre-load all search terms