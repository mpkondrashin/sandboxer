package main

import (
	"examen/pkg/fifo"
	"examen/pkg/logging"
	"log"
)

type StringChannel chan string

func SubmitDispatch(ch StringChannel) {
	fifoReader, err := fifo.NewReader()
	if err != nil {
		log.Fatal(err)
	}
	for {
		var s string
		if err := fifoReader.Read(&s); err != nil {
			logging.Errorf("fifo read: %v", err)
		}
		logging.Infof("Got new path: %s", s)
		ch <- s
	}
}
