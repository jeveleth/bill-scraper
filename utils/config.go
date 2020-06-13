package utils

import "flag"

type Config struct {
	State        string
	SearchPhrase string
	NumBills     int
	Session      string
}

func MustLoadConfig() Config {
	config := Config{}
	flag.StringVar(&config.State, "state", "California", "the name of the state you're researching")
	flag.StringVar(&config.SearchPhrase, "phrase", "peace officer", "the text of the phrase you're looking up")
	flag.IntVar(&config.NumBills, "num-bills", 5, "the number of bills you wish to retrieve")
	flag.StringVar(&config.Session, "sess", "2019-2020", "The year(s) of the legislative session")
	flag.Parse()
	return config
}
