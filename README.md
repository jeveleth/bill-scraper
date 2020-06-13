
This repo allows users to run a binary that will generate a CSV file that contains all relevant legislation for a given state.
Tested on a Mac, so it doesn't yet cross-compile for Windows or Linux.

To use:
```
export OPENSTATES_API_KEY=[your api key]
go run main.go
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