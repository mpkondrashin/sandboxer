package submit

import (
	"examen/pkg/fifo"
	"examen/pkg/logging"
	"log"
)

type StringChannel chan string

func SubmitDispatch(ch StringChannel) {
	for {
		fifoReader, err := fifo.NewReader()
		if err != nil {
			log.Fatal(err)
		}
		var s string
		if err := fifoReader.Read(&s); err != nil {
			if err.Error() != "EOF" {
				logging.Errorf("read FIFO: %v", err)
			}
			continue
		}
		fifoReader.Close()
		logging.Infof("Got new path: %s", s)
		ch <- s
	}
}
