/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

genver.go

Generate version.go file
*/
package main

import (
	"fmt"
	"os"
	"text/template"
)

const versionGo = `// Code generated by genver. DO NOT EDIT

package globals
var (
	Version = "{{.Version}}"
	Build   = "{{.Build}}"
)
`

func main() {
	if len(os.Args) != 4 {
		fmt.Println("Usage genver <version> <build> <outputfile>")
		os.Exit(1)
	}
	data := make(map[string]string)
	data["Version"] = os.Args[1]
	data["Build"] = os.Args[2]
	tmpl, err := template.New("version").Parse(versionGo)
	if err != nil {
		panic(err)
	}
	file, err := os.Create(os.Args[3])
	if err != nil {
		panic(err)
	}
	if err := tmpl.Execute(file, data); err != nil {
		panic(err)
	}
}
