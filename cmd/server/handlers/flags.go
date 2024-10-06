package handlers

import (
	"flag"
)

var MainURL = "localhost:8080"

func ParseFlags() {
	flag.StringVar(&MainURL, "a", "localhost:8080", "address and port to run server")
	flag.Parse()
}
