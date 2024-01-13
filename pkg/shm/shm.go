package shm

import (
	"encoding/json"
	"os"

	"bitbucket.org/avd/go-ipc/mmf"
	"bitbucket.org/avd/go-ipc/shm"
)

const memObjectName = "examen"

type SHMWriter struct {
	wRegion *mmf.MemoryRegion
}

func NewSHMWriter(size int) (*SHMWriter, error) {
	s := &SHMWriter{}
	shm.DestroyMemoryObject(memObjectName)
	obj, err := shm.NewMemoryObject(memObjectName, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	if err := obj.Truncate(1024); err != nil {
		return nil, err
	}
	s.wRegion, err = mmf.NewMemoryRegion(obj, mmf.MEM_READWRITE, 0, size)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (s *SHMWriter) Close() {
	s.wRegion.Close()
}

func (s *SHMWriter) Write(data any) error {
	writer := mmf.NewMemoryRegionWriter(s.wRegion)
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = writer.WriteAt(jsonData, 0)
	return err
}

type SHMReader struct {
	rRegion *mmf.MemoryRegion
}

func NewSHMReader(size int) (*SHMReader, error) {
	s := &SHMReader{}
	obj, err := shm.NewMemoryObject(memObjectName, os.O_RDONLY, 0600)
	if err != nil {
		return nil, err
	}
	s.rRegion, err = mmf.NewMemoryRegion(obj, mmf.MEM_READ_ONLY, 0, size)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (s *SHMReader) Read(data any) error {
	//log.Println("reader got Data: ", s.rRegion.Data())
	reader := mmf.NewMemoryRegionReader(s.rRegion)
	return json.NewDecoder(reader).Decode(data)
	/*
	   // write data at the specified offset
	   // read data at the same offset via another region.
	   actual := make([]byte, 10)
	   read, err := reader.ReadAt(actual, 0)

	   	if err != nil {
	   		log.Fatal(err)
	   	}

	   log.Println("reader got: (%d) ", read, actual)
	*/
}

func (s *SHMReader) Close() {
	s.rRegion.Close()
}

/*
func main() {
	// cleanup previous objects
	shm.DestroyMemoryObject("obj")
	// create new object and resize it.
	obj, err := shm.NewMemoryObject("obj", os.O_CREATE|os.O_EXCL|os.O_RDWR, 0600)
	if err != nil {
		log.Fatal(err)
	}
	if err := obj.Truncate(1024); err != nil {
		log.Fatal(err)
	}
	// create two regions for reading and writing.
	rwRegion, err := mmf.NewMemoryRegion(obj, mmf.MEM_READWRITE, 0, 1024)
	if err != nil {
		log.Fatal(err)
	}
	defer rwRegion.Close()
	// for each region we create a reader and a writer, which is a better solution, than
	// using region.Data() bytes directly.
	writer := mmf.NewMemoryRegionWriter(rwRegion)
	// write data at the specified offset
	for i := 0; ; i++ {
		log.Println("writer: ", i)
		var data []byte
		for j := 0; j < 10; j++ {
			data = append(data, byte(i))
		}
		written, err := writer.WriteAt(data, 0)
		if err != nil {
			panic(err)
		}
		if written != len(data) {
			log.Fatal("written:", written)
		}
		time.Sleep(1 * time.Second)
	}
}
*/
