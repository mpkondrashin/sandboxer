/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

script.go

uninstall script
*/
package script

import (
	"fmt"
	"os"
)

type Script struct {
	FilePath string
	lines    []string
	header   string
}

func New(filePath, header string) *Script {
	return &Script{
		FilePath: filePath,
		header:   header,
	}
}

func (s *Script) Save() error {
	f, err := os.Create(s.FilePath)
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
