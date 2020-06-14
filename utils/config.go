package utils

import "flag"

type Config struct {
	State        string
	SearchPhrase string
	NumBills     int
	Session      string
	Cursor       string
}

func MustLoadConfig() Config {
	config := Config{}
	flag.StringVar(&config.State, "state", "California", "The name of the state you're researching")
	flag.StringVar(&config.SearchPhrase, "phrase", "peace officer", "The text of the phrase you're looking up")
	flag.IntVar(&config.NumBills, "num-bills", 10, "The number of bills you wish to retrieve")
	flag.StringVar(&config.Session, "session", "20192020", "The year(s) of the legislative session")
	flag.StringVar(&config.Cursor, "cursor", "", "The cursor after which to start your query")
	flag.Parse()
	return config
}
