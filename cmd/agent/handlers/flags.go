package handlers

import (
	"errors"
	"flag"
	"strconv"
	"strings"
)

var MainURL = "localhost:8080"
var PollInterval = 2
var ReportInterval = 10

type NetAddress struct {
	Host string
	Port int
}

func (na NetAddress) String() string {
	return na.Host + ":" + strconv.Itoa(na.Port)
}

func (na *NetAddress) Set(s string) error {
	hp := strings.Split(s, ":")
	if len(hp) != 2 {
		return errors.New("need address in a form host:port")
	}
	port, err := strconv.Atoi(hp[1])
	if err != nil {
		return err
	}
	na.Host = hp[0]
	na.Port = port
	return nil
}

func ParseFlags() {
	flag.StringVar(&MainURL, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&ReportInterval, "r", 10, "report interval in seconds")
	flag.IntVar(&PollInterval, "p", 2, "poll interval in seconds")
	flag.Parse()
}
