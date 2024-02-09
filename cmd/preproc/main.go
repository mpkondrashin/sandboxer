/*
TunnelEffect (c) 2022 by Mikhail Kondrashin (mkondrashin@gmail.com)

main.go

Text preprocessor used to generate reame and administrator guide documents.
*/

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"text/template"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type preProc struct {
	Windows      bool
	Linux        bool
	Darwin       bool
	Date         string
	DateTime     string
	Version      string
	Build        string
	Platform     string
	Arch         string
	ManifestArch string
	OS           string
	GOOS         string
	GOARCH       string
	Exe          string
	Content      string
	//	Build    string
	//	Hash     string
}

func newPreProc() *preProc {
	p := preProc{
		Date:         time.Now().Format("02.01.2006"),
		DateTime:     time.Now().Format("02/01/2006 15:04:05"),
		Platform:     getPlatform(),
		OS:           getOS(),
		Arch:         getArch(),
		ManifestArch: getManifestArch(),
		GOOS:         os.Getenv("GOOS"),
		GOARCH:       os.Getenv("GOARCH"),
		//	Version:  os.Args[1],
		Windows: os.Getenv("GOOS") == "windows",
		Linux:   os.Getenv("GOOS") == "linux",
		Darwin:  os.Getenv("GOOS") == "darwin",
	}
	if p.Windows {
		p.Exe = ".exe"
	}
	var fileName string
	flag.StringVar(&p.Version, "version", "", "Version taken from git")
	flag.StringVar(&p.Build, "build", "", "Build — commit count taken from git")
	flag.StringVar(&fileName, "content", "", "Content of file provided")
	flag.Parse()
	if fileName != "" {
		data, err := os.ReadFile(fileName)
		if err != nil {
			panic(err)
		}
		p.Content = string(data)
	}
	return &p
}

func (p *preProc) Generate(templatePath, targetPath string) error {
	contents, err := os.ReadFile(templatePath)
	if err != nil {
		return err
	}
	t, err := template.New("preproc").Parse(string(contents))
	if err != nil {
		return err
	}
	out, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer out.Close()
	return t.Execute(out, p)
}

func getPlatform() string {
	return fmt.Sprintf("%s %s", getOS(), getArch())
}

func getArch() string {
	a := os.Getenv("GOARCH")
	if a == "" {
		panic(fmt.Errorf("%s is not set", "GOARCH"))
	}
	switch a {
	case "amd64":
		return "x64"
	case "386":
		return "x86"
	default:
		return strings.ToUpper(a)
	}
}

func getManifestArch() string {
	a := os.Getenv("GOARCH")
	if a == "" {
		panic(fmt.Errorf("%s is not set", "GOARCH"))
	}
	switch a {
	case "386":
		return "x86"
	default:
		return strings.ToLower(a)
	}
}

func getOS() string {
	o := os.Getenv("GOOS")
	if o == "" {
		panic(fmt.Errorf("%s is not set", "GOOS"))
	}
	switch o {
	case "darwin":
		return "macOS"
	default:
		caser := cases.Title(language.Und)
		return caser.String(os.Getenv("GOOS"))
	}
}

func usage() {
	fmt.Printf("Usage: %s <options> source target\n", os.Args[0])
	os.Exit(1)
}

func main() {
	p := newPreProc()
	if len(flag.Args()) != 2 {
		usage()
	}
	err := p.Generate(flag.Args()[0], flag.Args()[1])
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(3)
	}
}
