/*
TunnelEffect (c) 2022 by Mikhail Kondrashin (mkondrashin@gmail.com)

rotated.go

Implement unix log files rotation.
*/
package logging

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// File - comply to io.Writer to be used as
// output file for logging, providing log rotating
// feature
type File struct {
	file     *os.File
	flag     int
	perm     fs.FileMode
	folder   string
	fileName string
	size     int
	keep     int
	maxSize  int
}

func OpenRotated(folder, fileName string, perm fs.FileMode, maxSize, keep int) (*File, error) {
	r := NewRotated(folder, fileName, perm, maxSize, keep)
	if err := r.Open(); err != nil {
		return nil, err
	}
	return r, nil
}

// NewRotated - create new rotated file
func NewRotated(folder, fileName string, perm fs.FileMode, maxSize, keep int) *File {
	return &File{
		nil,
		os.O_RDWR | os.O_CREATE | os.O_APPEND,
		perm,
		folder,
		fileName,
		0,
		keep,
		maxSize,
	}
}

// Open - open file
func (r *File) Open() (err error) {
	filePath := filepath.Join(r.folder, r.fileName)
	r.file, err = os.OpenFile(filePath, r.flag, r.perm)
	return err
}

// Write - write to file
func (r *File) Write(data []byte) (n int, err error) {
	//	fmt.Printf("Write start: %v\n", r.file)
	r.size += len(data)
	if r.size > r.maxSize {
		err := r.rotate()
		if err != nil {
			return 0, err
		}
		r.size = len(data)
	}

	//	fmt.Printf("Write end: %v\n", r.file)
	return r.file.Write(data)
}

func (r *File) rotate() error {
	r.Close()
	// a.log -> a.log.0 a.log.0 -> a.log.1  a.log.1 -> a.log.2
	// keep = 3
	for i := r.keep - 1; i >= 0; i-- {
		fileName := fmt.Sprintf("%s.%d", r.fileName, i)
		filePath := filepath.Join(r.folder, fileName)
		var fileNameNext string
		if i > 0 {
			fileNameNext = fmt.Sprintf("%s.%d", r.fileName, i-1)
		} else {
			fileNameNext = r.fileName
		}
		filePathNext := filepath.Join(r.folder, fileNameNext)
		err := os.Rename(filePathNext, filePath)
		if err != nil {
			continue
		}
	}
	return r.Open()
}

// Close - close file
func (r *File) Close() error {
	if r.file == nil {
		return nil
	}
	return r.file.Close()
}
