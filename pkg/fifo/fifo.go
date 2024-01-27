package fifo

import (
	"encoding/json"
	"os"

	"bitbucket.org/avd/go-ipc/fifo"

	"sandboxer/pkg/globals"
)

type Writer struct {
	fifo fifo.Fifo
}

func NewWriter() (*Writer, error) {
	w := &Writer{}
	var err error
	w.fifo, err = fifo.New(globals.FIFOName /* os.O_CREATE|*/, os.O_WRONLY|fifo.O_NONBLOCK, 0600)
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (w *Writer) Close() error {
	return w.fifo.Close()
}
func (w *Writer) Write(data any) error {
	return json.NewEncoder(w.fifo).Encode(data)
}

type Reader struct {
	fifo fifo.Fifo
}

func NewReader() (*Reader, error) {
	r := &Reader{}
	var err error
	r.fifo, err = fifo.New(globals.FIFOName, os.O_CREATE|os.O_RDONLY, 0600)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (r *Reader) Close() error {
	return r.fifo.Close()
}

func (w *Reader) Read(data any) error {
	return json.NewDecoder(w.fifo).Decode(data)
}
