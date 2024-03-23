package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"sandboxer/pkg/config"
	"strings"
	"text/template"
)

type Data struct {
	Package string
	Fields  []Field
}

type Field struct {
	StructName string
	Setter     bool
	Getter     bool
	Name       string
	Type       string
}

func Process(a any) (result []Field, packageName string, structName string, err error) {
	//	val := reflect.ValueOf(a)
	//fmt.Println(os.Getenv("GOROOT"))

	val := reflect.TypeOf(a)
	fmt.Println("GOROOT:", runtime.GOROOT())
	fmt.Println("pkgpath:", val.PkgPath())
	//fmt.Println("path :", filepath.Join(runtime.GOROOT(), val.PkgPath()))
	parts := strings.Split(val.String(), ".")
	packageName = parts[0]
	structName = parts[1]
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		//var getter, setter bool
		tag := field.Tag.Get("gsetter")
		//name := val.Name()
		//name = cases.Title(language.English, cases.Compact).String(val.Name())
		typeString := field.Type.String()
		typeString = strings.TrimPrefix(typeString, packageName+".")
		result = append(result,
			Field{
				StructName: structName,
				Setter:     tag == "w" || tag == "rw" || tag == "",
				Getter:     tag == "r" || tag == "rw" || tag == "",
				Name:       field.Name,
				Type:       typeString,
			})
	}
	return
}

func addImports(filePath string) error {
	gopls, err := exec.LookPath("gopls")
	if err != nil {
		return err
	}
	cmd := []string{
		"imports",
		"-w",
		filePath,
	}

	return exec.Command(gopls, cmd...).Run()
}

var codeTemplate = `
package {{ .Package }}

{{range .Fields -}}
{{- if .Getter -}}
func (s *{{ .StructName }}) Get{{.Name}}() {{.Type}} {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return s.{{.Name}}
}

{{end -}}
{{- if .Setter -}}
func (s *{{ .StructName }}) Set{{.Name}}(value {{.Type}} ) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.{{.Name}} = value
}

{{end -}}
{{- end -}}
`

func Generate(s any, folder string) error {
	fields, packageName, structName, err := Process(s)
	if err != nil {
		return err
	}
	data := Data{
		Package: packageName,
		Fields:  fields,
	}
	t, err := template.New("code").Parse(codeTemplate)
	if err != nil {
		return err
	}
	fileName := strings.ToLower(structName) + "_gsetter.go"
	filePath := filepath.Join(folder, fileName)
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	err = t.Execute(f, data)
	if err != nil {
		return nil
	}
	err = addImports(filePath)
	if err != nil {
		return nil
	}
	return nil
}

func main() {
	//var x []string
	//fmt.Println(reflect.TypeOf(x).String())
	// return
	if err := Generate(config.Configuration{}, "../../pkg/config"); err != nil {
		panic(err)
	}
	if err := Generate(config.DDAn{}, "../../pkg/config"); err != nil {
		panic(err)
	}
	if err := Generate(config.VisionOne{}, "../../pkg/config"); err != nil {
		panic(err)
	}
	if err := Generate(config.Proxy{}, "../../pkg/config"); err != nil {
		panic(err)
	}
}
