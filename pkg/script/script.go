package script

import (
	"fmt"
	"os"
)

type Script struct {
	filePath string
	lines    []string
	header   string
}

func New(filePath, header string) *Script {
	return &Script{
		filePath: filePath,
		header:   header,
	}
}

func (s *Script) Save() error {
	f, err := os.Create(s.filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	fmt.Fprintf(f, "%s\n", s.header)
	for i := len(s.lines) - 1; i >= 0; i-- {
		fmt.Fprintf(f, "%s\n", s.lines[i])
	}
	return nil
}

func (s *Script) AddLine(line string) error {
	s.lines = append(s.lines, line)
	return s.Save()
}
