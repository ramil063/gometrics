package main

import (
	"log"
)

func mulfunc(i int) (int, error) {
	return i * 2, nil
}

func main() {
	i, err := mulfunc(1)
	if err != nil {
		return
	}
	log.Println(i)
}
